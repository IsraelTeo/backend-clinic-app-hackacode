package logic

import (
	"log"
	"strings"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/calculation"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

type AppointmentLogic interface {
	GetAppointmentByID(ID uint) (*model.Appointment, error)
	GetAllAppointments() ([]model.Appointment, error)
	CreateAppointment(appointment *model.Appointment) (float64, float64, float64, float64, error)
	UpdateAppointment(ID uint, appointment *model.Appointment) error
	DeleteAppointment(ID uint) error
}

type appointmentLogic struct {
	repositoryAppointment     repository.Repository[model.Appointment]
	repositoryAppointmentMain repository.AppointmentRepository
	repositoryDoctor          repository.Repository[model.Doctor]
	repositoryPatient         repository.Repository[model.Patient]
	repositoryPackage         repository.Repository[model.Package]
	repositoryService         repository.Repository[model.Service]
	logicPatient              PatientLogic
	repositoryPatientMain     repository.PatientRepository
	repositoryPackageMain     repository.PackageRepository
}

func NewAppointmentLogic(
	repositoryAppointment repository.Repository[model.Appointment],
	repositoryDoctor repository.Repository[model.Doctor],
	repositoryPatient repository.Repository[model.Patient],
	repositoryAppointmentMain repository.AppointmentRepository,
	repositoryPackage repository.Repository[model.Package],
	repositoryService repository.Repository[model.Service],
	logicPatient PatientLogic,
	repositoryPatientMain repository.PatientRepository,
	repositoryPackageMain repository.PackageRepository,
) AppointmentLogic {
	return &appointmentLogic{
		repositoryAppointment:     repositoryAppointment,
		repositoryDoctor:          repositoryDoctor,
		repositoryPatient:         repositoryPatient,
		repositoryAppointmentMain: repositoryAppointmentMain,
		repositoryPackage:         repositoryPackage,
		repositoryService:         repositoryService,
		logicPatient:              logicPatient,
		repositoryPatientMain:     repositoryPatientMain,
		repositoryPackageMain:     repositoryPackageMain,
	}
}

func (l *appointmentLogic) GetAppointmentByID(ID uint) (*model.Appointment, error) {
	appointment, err := l.repositoryAppointmentMain.GetByID(ID)
	if err != nil {
		log.Printf("appointment: Error fetching appointment with ID %d: %v", ID, err)
		return nil, response.ErrorAppointmentNotFound
	}
	return appointment, nil
}

func (l *appointmentLogic) GetAllAppointments() ([]model.Appointment, error) {
	appointments, err := l.repositoryAppointmentMain.GetAll()
	if err != nil {
		log.Printf("appointment: Error fetching appointments: %v", err)
		return nil, response.ErrorAppointmetsNotFound
	}

	if len(appointments) == 0 {
		log.Println("appointment : No appointments found")
		return []model.Appointment{}, response.ErrorListAppointmentsEmpty
	}

	return appointments, nil
}

func (l *appointmentLogic) CreateAppointment(appointment *model.Appointment) (float64, float64, float64, float64, error) {
	// Verificar que el doctor exista para la cita
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return 0, 0, 0, 0, response.ErrorDoctorNotFound
	}

	// Verificar si el ID del paciente está presente en el request
	if appointment.PatientID != 0 { // Si el patient_id es proporcionado
		// Buscar paciente por ID y asignar a la cita
		existingPatient, err := l.repositoryPatient.GetByID(appointment.PatientID)
		if err != nil {
			log.Printf("appointment: Error fetching patient by ID: %v", err)
			return 0, 0, 0, 0, response.ErrorPatientNotFound
		}
		appointment.Patient = *existingPatient // Asignamos el paciente encontrado a la cita
	} else if appointment.Patient != (model.Patient{}) { // Si no hay patient_id, pero el objeto de paciente está presente
		// Crear un nuevo paciente con los datos proporcionados
		existingPatient, _ := l.repositoryPatientMain.GetPatientByDNI(appointment.Patient.DNI)
		if existingPatient != nil {
			log.Printf("appointment: Patient with DNI already exists: %s", appointment.Patient.DNI)
			return 0, 0, 0, 0, response.ErrorPatientExists
		}

		// Validar y formatear la fecha de nacimiento
		parsedBirthDate, err := validate.ParseDate(appointment.Patient.BirthDate)
		if err != nil {
			log.Printf("appointment: Invalid patient birth date format: %v", err)
			return 0, 0, 0, 0, response.ErrorPatientInvalidDateFormat
		}
		appointment.Patient.BirthDate = parsedBirthDate.Format("2006-01-02")

		// Crear el nuevo paciente
		if err := l.repositoryPatient.Create(&appointment.Patient); err != nil {
			log.Printf("appointment: Error creating patient: %v", err)
			return 0, 0, 0, 0, err
		}
	} else {
		// Si no se proporciona ni un patient_id ni los datos del paciente, lanzar un error
		log.Println("appointment: Patient data or patient_id is required")
		return 0, 0, 0, 0, response.ErrorPatientDataRequired
	}

	// Parsear la fecha de la cita
	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		log.Printf("appointment: Invalid appointment date format: %v", err)
		return 0, 0, 0, 0, response.ErrorAppointmentInvalidDateFormat
	}

	// Parsear la fecha de la cita y el horario
	startTime, endTime, err := parseStartAndEndTime(appointment.StartTime, appointment.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid appointment date format: %v", err)
		return 0, 0, 0, 0, response.ErrorAppointmentInvalidDateFormat
	}

	// Verificar que la fecha de la cita sea en el futuro
	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment: Appointment date is in the past")
		return 0, 0, 0, 0, response.ErrorAppointmentDateInPast
	}

	// Obtener el día de la semana (en inglés) y traducirlo a español
	appointmentDay := validate.TranslateDayToSpanish(appointmentDate.Weekday().String())
	log.Printf("appointment: Día de la cita en español: %s", appointmentDay) // Verificar el valor de appointmentDay

	// Dividir los días de la semana que trabaja el doctor (en español)
	validDays := strings.Split(doctor.Days, ",")
	for i := range validDays {
		validDays[i] = strings.ToLower(strings.TrimSpace(validDays[i])) // Convertir a minúsculas y eliminar espacios
	}
	log.Printf("appointment: Días disponibles del doctor: %v", validDays) // Verificar los días disponibles del doctor

	// Convertir el día de la cita a minúsculas para hacer la comparación
	appointmentDay = strings.ToLower(appointmentDay)
	log.Printf("appointment: Día de la cita en minúsculas: %s", appointmentDay) // Verificar el valor de appointmentDay en minúsculas

	// Verificar si el médico está disponible el día seleccionado
	if !validate.IsDayAvailable(appointmentDay, validDays) {
		log.Printf("appointment: El médico no está disponible el día: %s", appointmentDay)
		return 0, 0, 0, 0, response.ErrorAppointmentDayNotAvailable
	}

	//verificamos que el horario inicial sea antes del horario final
	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment: Start time is not before end time")
		return 0, 0, 0, 0, response.ErrorInvalidAppointmentTimeRange
	}

	doctorStartTime, doctorEndTime, err := parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid doctor end time format: %v", err)
		return 0, 0, 0, 0, response.ErrorInvalidDoctorTime
	}

	if !validate.IsWithinTimeRange(startTime, endTime, doctorStartTime, doctorEndTime) {
		log.Printf("appointment: Appointment time is outside the allowed time range")
		return 0, 0, 0, 0, response.ErrorInvalidAppointmentTime
	}

	// Validar que no haya conflicto con otras citas del doctor
	existingAppointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointment.Date)
	if err != nil {
		log.Printf("appointment: Error fetching existing appointments: %v", err)
		return 0, 0, 0, 0, err
	}

	if validate.HasTimeConflict(existingAppointments, startTime, endTime) {
		log.Printf("appointment: Appointment time conflicts with existing appointments")
		return 0, 0, 0, 0, response.ErrorAppointmentTimeConflict
	}

	hasInsurance := appointment.Patient.Insurance

	// Declarar las variables de manera consistente
	var (
		originalPrice, finalPrice, discount, packageDiscount, insuranceDiscount float64
	)

	// Verificar si la cita es para un paquete
	if appointment.PackageID != 0 {
		// Obtener el paquete por su ID
		pkg, err := l.repositoryPackageMain.GetByID(appointment.PackageID)
		if err != nil || pkg == nil {
			log.Printf("appointment: El paquete con ID %d no existe: %v", appointment.PackageID, err)
			return 0, 0, 0, 0, response.ErrorPackageNotFound
		}

		// Log para verificar los servicios en el paquete
		log.Printf("appointment: Servicios en el paquete: %+v", pkg.Services)

		// Calcular precios para el paquete
		originalPrice, packageDiscount, finalPrice = calculation.TotalServicePackageAmount(pkg.Services, hasInsurance)

		// Verificar si hay seguro y calcular el descuento por seguro
		if hasInsurance {
			// El descuento por seguro se aplica al precio después del descuento del paquete
			insuranceDiscount = (originalPrice - packageDiscount) * 0.20
		}

		// Aplicar todos los descuentos al precio final
		finalPrice = (originalPrice - packageDiscount) - insuranceDiscount

		// Log para validar los cálculos
		log.Printf("appointment: Precios calculados - OriginalPrice: %.2f, PackageDiscount: %.2f, InsuranceDiscount: %.2f, FinalPrice: %.2f", originalPrice, packageDiscount, insuranceDiscount, finalPrice)

		// Asignar el precio final a la cita
		appointment.TotalAmount = finalPrice

	} else if appointment.ServiceID != 0 {
		// Cita para un servicio único
		service, err := l.repositoryService.GetByID(appointment.ServiceID)
		if err != nil || service == nil {
			log.Printf("appointment: El servicio con ID %d no existe: %v", appointment.ServiceID, err)
			return 0, 0, 0, 0, response.ErrorServiceNotFound
		}

		// Precio original del servicio
		originalPrice = service.Price

		// Si hay seguro, calcular el descuento
		if hasInsurance {
			discount = originalPrice * 0.20 // 20% de descuento por seguro
		}

		// Calcular el precio final
		finalPrice = originalPrice - discount
		appointment.TotalAmount = finalPrice

	} else {
		log.Println("appointment: Debe especificar un paquete o un servicio.")
		return 0, 0, 0, 0, response.ErrorPackageAndServiceEmpty
	}

	// Crear la cita
	appointmentCreated := &model.Appointment{
		Patient:     appointment.Patient,
		DoctorID:    appointment.DoctorID,
		ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: appointment.TotalAmount,
	}

	if err := l.repositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment: Error al crear la cita: %v", err)
		return 0, 0, 0, 0, err
	}

	log.Printf("OriginalPrice antes del cálculo: %f", originalPrice)
	log.Printf("PackageDiscount antes del cálculo: %f", packageDiscount)
	log.Printf("InsuranceDiscount antes del cálculo: %f", insuranceDiscount)
	log.Printf("FinalPrice antes del cálculo: %f", finalPrice)

	return originalPrice, packageDiscount, finalPrice, insuranceDiscount, nil
}

func (l *appointmentLogic) UpdateAppointment(ID uint, appointment *model.Appointment) error {
	// Obtener la cita actual por ID
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return err
	}
	if err := l.repositoryAppointment.Update(appointmentUpdate); err != nil {
		log.Printf("appointment: Error updating appointment with ID %d: %v", ID, err)
		return err
	}

	return nil
}

func (l *appointmentLogic) DeleteAppointment(ID uint) error {
	// Verificar si el ID existe
	_, err := l.repositoryAppointment.GetByID(ID)
	if err != nil {
		return err
	}

	// Intentar eliminar el registro
	if err := l.repositoryAppointment.Delete(ID); err != nil {
		return err
	}

	return nil
}

func parseStartAndEndTime(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	// Parsear horarios de la cita y validarlos
	startTime, err := validate.ParseTime(startTimeStr)
	if err != nil {
		log.Printf("appointment: Invalid start time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	endTime, err := validate.ParseTime(endTimeStr)
	if err != nil {
		log.Printf("appointment: Invalid end time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	// Retornar la fecha de la cita y los horarios de inicio y fin
	return startTime, endTime, nil
}
