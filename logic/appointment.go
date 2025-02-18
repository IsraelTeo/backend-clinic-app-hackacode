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
	CreateAppointmentWithPackage(appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error)
	CreateAppointmentWithService(appointment *model.Appointment) (*model.FinalServicePrice, error)
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
		log.Printf("appointment-logic: Error fetching appointment with ID %d: %v", ID, err)
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

func (l *appointmentLogic) CreateAppointmentWithPackage(appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error) {
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment-logic: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return nil, response.ErrorDoctorNotFoundID
	}

	if l.isPatientExistsInAppointment(appointment) {
		if err := validate.PatientToCreate(&appointment.Patient); err != nil {
			return nil, err
		}

		if err := l.repositoryPatient.Create(&appointment.Patient); err != nil {
			log.Printf("appointment: Error creating patient: %v", err)
			return nil, err
		}
	}

	startTime, endTime, appointmentDate, err := l.parseTimesAndDate(appointment)
	if err != nil {
		return nil, err
	}

	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment: Appointment date is in the past")
		return nil, response.ErrorAppointmentDateInPast
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment: Start time is not before end time")
		return nil, response.ErrorInvalidAppointmentTimeRange
	}

	doctorStartTime, doctorEndTime, err := l.parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid doctor end time format: %v", err)
		return nil, err
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
		return nil, response.ErrorAppointmentDayNotAvailable
	}

	if !validate.IsWithinTimeRange(startTime, endTime, doctorStartTime, doctorEndTime) {
		log.Printf("appointment: Appointment time is outside the allowed time range")
		return nil, response.ErrorInvalidAppointmentTime
	}

	existingAppointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointment.Date)
	if err != nil {
		log.Printf("appointment: Error fetching existing appointments: %v", err)
		return nil, err
	}

	if validate.HasTimeConflict(existingAppointments, startTime, endTime) {
		log.Printf("appointment: Appointment time conflicts with existing appointments")
		return nil, response.ErrorAppointmentTimeConflict
	}

	hasInsurance := appointment.Patient.Insurance

	finalPkgPrice, err := l.isPackageIDExists(appointment.PackageID, hasInsurance)
	if err != nil {
		return nil, err
	}

	//creamos la cita
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
		return nil, err
	}

	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: finalPkgPrice.InsuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     finalPkgPrice.TotalAmount,
			DiscountPackage: finalPkgPrice.DiscountPackage,
			FinalPrice:      finalPkgPrice.FinalPrice,
		},
	}, nil
}

func (l *appointmentLogic) CreateAppointmentWithService(appointment *model.Appointment) (*model.FinalServicePrice, error) {
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return &model.FinalServicePrice{}, response.ErrorDoctorNotFoundID
	}

	if l.isPatientExistsInAppointment(appointment) {
		if err := validate.PatientToCreate(&appointment.Patient); err != nil {
			return &model.FinalServicePrice{}, err
		}

		if err := l.repositoryPatient.Create(&appointment.Patient); err != nil {
			log.Printf("appointment: Error creating patient: %v", err)
			return &model.FinalServicePrice{}, err
		}
	}

	startTime, endTime, appointmentDate, err := l.parseTimesAndDate(appointment)
	if err != nil {
		return &model.FinalServicePrice{}, err
	}

	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment: Appointment date is in the past")
		return nil, response.ErrorAppointmentDateInPast
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment: Start time is not before end time")
		return nil, response.ErrorInvalidAppointmentTimeRange
	}

	doctorStartTime, doctorEndTime, err := l.parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment: Invalid doctor end time format: %v", err)
		return nil, err
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
		return nil, response.ErrorAppointmentDayNotAvailable
	}

	if !validate.IsWithinTimeRange(startTime, endTime, doctorStartTime, doctorEndTime) {
		log.Printf("appointment: Appointment time is outside the allowed time range")
		return nil, response.ErrorInvalidAppointmentTime
	}

	existingAppointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointment.Date)
	if err != nil {
		log.Printf("appointment: Error fetching existing appointments: %v", err)
		return nil, err
	}

	if validate.HasTimeConflict(existingAppointments, startTime, endTime) {
		log.Printf("appointment: Appointment time conflicts with existing appointments")
		return nil, response.ErrorAppointmentTimeConflict
	}

	hasInsurance := appointment.Patient.Insurance

	finalServicePrice, err := l.isServiceIDEXists(appointment.ServiceID, hasInsurance)
	if err != nil {
		return nil, err
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
		return &model.FinalServicePrice{}, err
	}

	return &model.FinalServicePrice{
		TotalAmount:       finalServicePrice.TotalAmount,
		InsuranceDiscount: finalServicePrice.InsuranceDiscount,
		FinalPrice:        finalServicePrice.FinalPrice,
	}, nil
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

func (l *appointmentLogic) parseStartAndEndTime(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
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

func (l *appointmentLogic) parseTimesAndDate(appointment *model.Appointment) (time.Time, time.Time, time.Time, error) {
	startTime, endTime, err := l.parseStartAndEndTime(appointment.EndTime, appointment.EndTime)
	if err != nil {
		log.Printf("appointment-logic: Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, response.ErrorAppointmentInvalidDateFormat
	}

	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		log.Printf("appointment: Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, response.ErrorAppointmentInvalidDateFormat
	}

	return startTime, endTime, appointmentDate, nil
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

func (l *appointmentLogic) ValidatePatient(patient *model.Patient) error {
	if err := validate.DNIPatient(patient); err != nil {
		return err
	}

	if err := validate.PhoneNumberPatient(patient); err != nil {
		return err
	}

	if err := validate.EmailPatient(patient); err != nil {
		return err
	}

	if err := validate.BirthDatePatient(patient.BirthDate); err != nil {
		return err
	}

	return nil
}

func (l *appointmentLogic) ValidateUpdatedPatientFields(patient *model.Patient, patientUpdate *model.Patient) error {
	if patient.DNI != patientUpdate.DNI {
		if err := validate.DNIPatient(patient); err != nil {
			return err
		}
	}

	if patient.PhoneNumber != patientUpdate.PhoneNumber {
		if err := validate.PhoneNumberPatient(patient); err != nil {
			return err
		}
	}

	if patient.Email != patientUpdate.Email {
		if err := validate.EmailPatient(patient); err != nil {
			return err
		}
	}

	return nil
}

func (l *appointmentLogic) isPatientExistsInAppointment(appointment *model.Appointment) bool {
	if patientExists := l.isPatientIDExists(appointment.PatientID); patientExists != nil {
		appointment.Patient = *patientExists
		return true
	}

	if l.isPatientBodyExists(appointment) {
		return true
	}

	return false
}

func (l *appointmentLogic) isPatientIDExists(ID uint) *model.Patient {
	if ID == 0 {
		return nil
	}

	log.Println("appointment-logic: the request brought an ID patient")
	patient, err := l.repositoryPatient.GetByID(ID)
	if err != nil {
		log.Printf("appointment-logic: Error fetching patient by ID: %v", err)
		return nil
	}

	return patient
}

func (l *appointmentLogic) isPatientBodyExists(appointment *model.Appointment) bool {
	return appointment.Patient != (model.Patient{})
}

func (l *appointmentLogic) isServiceIDEXists(ID uint, hasInsurance bool) (*model.FinalServicePrice, error) {
	if ID == 0 {
		return &model.FinalServicePrice{}, response.ErrorInvalidID
	}

	log.Println("appointment-logic: the request brought an ID service")

	service, err := l.repositoryService.GetByID(ID)
	if err != nil || service == nil {
		log.Printf("appointment: El servicio con ID %d no existe: %v", ID, err)
		return &model.FinalServicePrice{}, response.ErrorServiceNotFound
	}

	var discount float64 = 0.0

	if hasInsurance {
		discount = service.Price * 0.20
	}

	finalPrice := service.Price - discount

	return &model.FinalServicePrice{
		TotalAmount:       service.Price,
		InsuranceDiscount: discount,
		FinalPrice:        finalPrice,
	}, nil
}

func (l *appointmentLogic) isPackageIDExists(ID uint, hasInsurance bool) (*model.FinalPackagePriceWithInsegurance, error) {
	pkg, err := l.repositoryPackageMain.GetByID(ID)
	if err != nil || pkg == nil {
		log.Printf("appointment: El paquete con ID %d no existe: %v", ID, err)
		return &model.FinalPackagePriceWithInsegurance{}, response.ErrorPackageNotFound
	}

	finalPricePkg := calculation.TotalServicePackageAmountToAppointment(pkg.Services, hasInsurance)
	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: finalPricePkg.InsuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     finalPricePkg.TotalAmount,
			DiscountPackage: finalPricePkg.DiscountPackage,
			FinalPrice:      finalPricePkg.FinalPrice,
		},
	}, nil
}

/*//nos dieron un ID de paquete en el JSON?
if appointment.PackageID != 0 {
	//verificamos que el paquete exista
	pkg, err := l.repositoryPackageMain.GetByID(appointment.PackageID)
	if err != nil || pkg == nil {
		log.Printf("appointment: El paquete con ID %d no existe: %v", appointment.PackageID, err)
		return &model.FinalPackagePriceWithInsegurance{}, response.ErrorPackageNotFound
	}

	log.Printf("appointment: Servicios en el paquete: %+v", pkg.Services)

	finalPricePkg = calculation.TotalServicePackageAmountToAppointment(pkg.Services, hasInsurance)

	//nos dieron el ID  de un servicio en el JSON?
} else if appointment.ServiceID != 0 {
	//verificamos que el servicio exista
	service, err := l.repositoryService.GetByID(appointment.ServiceID)
	if err != nil || service == nil {
		log.Printf("appointment: El servicio con ID %d no existe: %v", appointment.ServiceID, err)
		return nil, response.ErrorServiceNotFound
	}

	//precio original del servicio
	originalPrice = service.Price

	if hasInsurance {
		discount = originalPrice * 0.20
	}

	finalPrice = originalPrice - discount //descuento del servicio si hay seuro
	appointment.TotalAmount = finalPrice

} else {
	//es necesario un paquete o un servicio
	log.Println("appointment: Debe especificar un paquete o un servicio.")
	return &model.FinalPackagePriceWithInsegurance{}, response.ErrorPackageAndServiceEmpty
}*/
