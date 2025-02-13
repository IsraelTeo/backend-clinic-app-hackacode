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
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return 0, 0, 0, 0, response.ErrorDoctorNotFoundID
	}

	if appointment.PatientID != 0 {
		log.Println("the request brought an ID")

		existingPatient, err := l.repositoryPatient.GetByID(appointment.PatientID)
		if err != nil {
			log.Printf("appointment: Error fetching patient by ID: %v", err)
			return 0, 0, 0, 0, response.ErrorPatientNotFoundID
		}

		appointment.Patient = *existingPatient

	} else if appointment.Patient != (model.Patient{}) {
		existingPatient, _ := l.repositoryPatientMain.GetPatientByDNI(appointment.Patient.DNI)
		if existingPatient != nil {
			log.Printf("appointment: Patient with DNI already exists: %s", appointment.Patient.DNI)
			return 0, 0, 0, 0, response.ErrorPatientExists
		}

		parsedBirthDate, err := validate.ParseDate(appointment.Patient.BirthDate)
		if err != nil {
			log.Printf("appointment: Invalid patient birth date format: %v", err)
			return 0, 0, 0, 0, response.ErrorPatientInvalidDateFormat
		}

		appointment.Patient.BirthDate = parsedBirthDate.Format("2006-01-02")

		if err := l.repositoryPatient.Create(&appointment.Patient); err != nil {
			log.Printf("appointment: Error creating patient: %v", err)
			return 0, 0, 0, 0, err
		}
	} else {
		log.Println("appointment: Patient data or patient_id is required")
		return 0, 0, 0, 0, response.ErrorPatientDataRequired
	}

	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		log.Printf("appointment: Invalid appointment date format: %v", err)
		return 0, 0, 0, 0, response.ErrorAppointmentInvalidDateFormat
	}

	startTime, endTime, err := parseStartAndEndTime(appointment.StartTime, appointment.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid appointment date format: %v", err)
		return 0, 0, 0, 0, response.ErrorAppointmentInvalidDateFormat
	}

	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment: Appointment date is in the past")
		return 0, 0, 0, 0, response.ErrorAppointmentDateInPast
	}

	appointmentDay := validate.TranslateDayToSpanish(appointmentDate.Weekday().String())
	log.Printf("appointment: Día de la cita en español: %s", appointmentDay)

	validDays := strings.Split(doctor.Days, ",")
	for i := range validDays {
		validDays[i] = strings.ToLower(strings.TrimSpace(validDays[i]))
	}

	log.Printf("appointment: Días disponibles del doctor: %v", validDays)

	appointmentDay = strings.ToLower(appointmentDay)
	log.Printf("appointment: Día de la cita en minúsculas: %s", appointmentDay)

	if !validate.IsDayAvailable(appointmentDay, validDays) {
		log.Printf("appointment: El médico no está disponible el día: %s", appointmentDay)
		return 0, 0, 0, 0, response.ErrorAppointmentDayNotAvailable
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment: Start time is not before end time")
		return 0, 0, 0, 0, response.ErrorInvalidAppointmentTimeRange
	}

	doctorStartTime, doctorEndTime, err := parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid doctor end time format: %v", err)
		return 0, 0, 0, 0, err
	}

	if !validate.IsWithinTimeRange(startTime, endTime, doctorStartTime, doctorEndTime) {
		log.Printf("appointment: Appointment time is outside the allowed time range")
		return 0, 0, 0, 0, response.ErrorInvalidAppointmentTime
	}

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

	var (
		originalPrice, finalPrice, discount, packageDiscount, insuranceDiscount float64
	)

	if appointment.PackageID != 0 {
		pkg, err := l.repositoryPackageMain.GetByID(appointment.PackageID)
		if err != nil || pkg == nil {
			log.Printf("appointment: El paquete con ID %d no existe: %v", appointment.PackageID, err)
			return 0, 0, 0, 0, response.ErrorPackageNotFound
		}

		log.Printf("appointment: Servicios en el paquete: %+v", pkg.Services)

		originalPrice, packageDiscount, finalPrice = calculation.TotalServicePackageAmount(pkg.Services, hasInsurance)

		if hasInsurance {
			insuranceDiscount = (originalPrice - packageDiscount) * 0.20
		}

		finalPrice = (originalPrice - packageDiscount) - insuranceDiscount

		log.Printf("appointment: Precios calculados - OriginalPrice: %.2f, PackageDiscount: %.2f, InsuranceDiscount: %.2f, FinalPrice: %.2f", originalPrice, packageDiscount, insuranceDiscount, finalPrice)

		appointment.TotalAmount = finalPrice

	} else if appointment.ServiceID != 0 {
		service, err := l.repositoryService.GetByID(appointment.ServiceID)
		if err != nil || service == nil {
			log.Printf("appointment: El servicio con ID %d no existe: %v", appointment.ServiceID, err)
			return 0, 0, 0, 0, response.ErrorServiceNotFound
		}

		originalPrice = service.Price

		if hasInsurance {
			discount = originalPrice * 0.20
		}

		finalPrice = originalPrice - discount
		appointment.TotalAmount = finalPrice

	} else {
		log.Println("appointment: Debe especificar un paquete o un servicio.")
		return 0, 0, 0, 0, response.ErrorPackageAndServiceEmpty
	}

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
	_, err := l.repositoryAppointment.GetByID(ID)
	if err != nil {
		return err
	}

	if err := l.repositoryAppointment.Delete(ID); err != nil {
		return err
	}

	return nil
}

func parseStartAndEndTime(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
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

	return startTime, endTime, nil
}

func (l *appointmentLogic) validateDoctor(doctorID uint) error {
	doctor, err := l.repositoryDoctor.GetByID(doctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", doctorID, err)
		return response.ErrorDoctorNotFoundID
	}
	return nil
}

func (l *appointmentLogic) handlePatientByID(appointment *model.Appointment) error {
	if appointment.PatientID == 0 {
		return nil
	}

	log.Println("the request brought an ID")
	existingPatient, err := l.repositoryPatient.GetByID(appointment.PatientID)
	if err != nil {
		log.Printf("appointment: Error fetching patient by ID: %v", err)
		return response.ErrorPatientNotFoundID
	}

	appointment.Patient = *existingPatient
	return nil
}
