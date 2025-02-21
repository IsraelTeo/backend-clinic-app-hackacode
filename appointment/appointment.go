package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentLogic interface {
	GetAppointmentByID(ID uint) (*model.Appointment, error)
	GetAllAppointments() ([]model.Appointment, error)
	CreateAppointment(appointment *model.Appointment) (interface{}, error)
	UpdateAppointment(ID uint, appointment *model.Appointment) (interface{}, error)
	DeleteAppointment(ID uint) error
}

type appointmentLogic struct {
	repositoryAppointment     repository.Repository[model.Appointment]
	repositoryAppointmentMain repository.AppointmentRepository
	repositoryDoctor          repository.Repository[model.Doctor]
	repositoryPatient         repository.Repository[model.Patient]
	repositoryService         repository.Repository[model.Service]
	repositoryPackageMain     repository.PackageRepository
	repositoryPackage         repository.Repository[model.Package]
	logicAppointmentCreate    AppointmentCreate
	logicAppointmentUpdate    AppointmentUpdate
}

func NewAppointmentLogic(
	repositoryAppointment repository.Repository[model.Appointment],
	repositoryAppointmentMain repository.AppointmentRepository,
	repositoryDoctor repository.Repository[model.Doctor],
	repositoryPatient repository.Repository[model.Patient],
	repositoryService repository.Repository[model.Service],
	repositoryPackageMain repository.PackageRepository,
	repositoryPackage repository.Repository[model.Package],
	logicAppointmentCreate AppointmentCreate,
	logicAppointmentUpdate AppointmentUpdate) AppointmentLogic {
	return &appointmentLogic{
		repositoryAppointment:     repositoryAppointment,
		repositoryAppointmentMain: repositoryAppointmentMain,
		repositoryDoctor:          repositoryDoctor,
		repositoryPatient:         repositoryPatient,
		repositoryService:         repositoryService,
		repositoryPackageMain:     repositoryPackageMain,
		repositoryPackage:         repositoryPackage,
		logicAppointmentCreate:    logicAppointmentCreate,
		logicAppointmentUpdate:    logicAppointmentUpdate,
	}
}

/*type appointmentRepositories struct {
	RepositoryAppointment     repository.Repository[model.Appointment]
	RepositoryAppointmentMain repository.AppointmentRepository
	RepositoryDoctor          repository.Repository[model.Doctor]
	RepositoryPatient         repository.Repository[model.Patient]
	RepositoryService         repository.Repository[model.Service]
	RepositoryPackageMain     repository.PackageRepository
	RepositoryPackage         repository.Repository[model.Package]
}*/
/*
	func NewAppointmentLogic(repo AppointmentRepositories, logicAppointmentCreate CreateAppointment, logicAppointmentUpdate AppointmentUpdate) AppointmentLogic {
		log.Println("RepoAppointmentMain:", repo.RepositoryAppointmentMain)
		log.Println("LogicAppointmentCreate:", logicAppointmentCreate)
		return &appointmentLogic{
			repo:                   repo,
			logicAppointmentCreate: logicAppointmentCreate,
			logicAppointmentUpdate: logicAppointmentUpdate,
		}
	}
*/
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
		log.Printf("appointment-logic: Error fetching appointments: %v", err)
		return nil, response.ErrorAppointmetsNotFound
	}

	return appointments, nil
}

func (l *appointmentLogic) CreateAppointment(appointment *model.Appointment) (interface{}, error) {
	finalPrice, err := l.logicAppointmentCreate.CreateAppointment(appointment)
	if err != nil {
		log.Printf("appointment-logic -> method: CreateAppointment: Error to create: %v", err)
		return nil, err
	}

	return finalPrice, nil
}

func (l *appointmentLogic) UpdateAppointment(ID uint, appointment *model.Appointment) (interface{}, error) {
	finalPrice, err := l.logicAppointmentUpdate.UpdateAppointment(ID, appointment)
	if err != nil {
		log.Printf("appointment-logic -> method: UpdateAppointment: Error to update: %v", err)
		return nil, err
	}

	return finalPrice, nil
}

func (l *appointmentLogic) DeleteAppointment(ID uint) error {
	if _, err := l.GetAppointmentByID(ID); err != nil {
		return response.ErrorAppointmentNotFound
	}

	if err := l.repositoryAppointment.Delete(ID); err != nil {
		return response.ErrorToDeletedAppointment
	}

	return nil
}

/*func (l *appointmentLogic) CreateAppointmentWithPackage(appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error) {
	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		log.Printf("appointment-logic: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return nil, err
	}

	//tratando de crear el paciente con body
	//if err := l.handlePatientBodyCreation(appointment); err != nil {
	//	return nil, err
	//}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalPkgPrice, err := l.isPackageIDExists(appointment.PackageID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentCreated := &model.Appointment{
		//PatientID:   appointment.PatientID,
		DoctorID: appointment.DoctorID,
		//ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: appointment.TotalAmount,
	}
	//tratando de crear el paciente con ID
	if err := l.handlePatientIDCreation(appointment); err != nil {
		return nil, err
	}

	//tratando de crear el paciente con ID
	if err := l.handlePatientIDCreation(appointmentCreated); err != nil {
		return nil, err
	}

	if err := l.repo.RepositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment-logic: Error al crear la cita: %v", err)
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
	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	if err := l.handlePatientIDCreation(appointment); err != nil {
		return nil, err
	}

	if err := l.handlePatientBodyCreation(appointment); err != nil {
		return nil, err
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalServicePrice, err := l.isServiceIDEXists(appointment.ServiceID, appointment.Patient.Insurance)
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

	if err := l.repo.RepositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment-logic: Error al crear la cita: %v", err)
		return &model.FinalServicePrice{}, err
	}

	return &model.FinalServicePrice{
		TotalAmount:       finalServicePrice.TotalAmount,
		InsuranceDiscount: finalServicePrice.InsuranceDiscount,
		FinalPrice:        finalServicePrice.FinalPrice,
	}, nil
}

func (l *appointmentLogic) UpdateAppointmentWithPackage(ID uint, appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error) {
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return nil, err
	}

	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	// Si el paciente tiene ID, solo se asocia. Si no, se registra uno nuevo.
	if appointment.Patient.ID == 0 {
		if err := l.handlePatientBodyCreation(appointment); err != nil {
			return nil, err
		}
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalPkgPrice, err := l.isPackageIDExists(appointment.PackageID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentUpdate.Patient = appointment.Patient
	appointmentUpdate.DoctorID = appointment.DoctorID
	appointmentUpdate.PackageID = appointment.PackageID
	appointmentUpdate.Date = appointment.Date
	appointmentUpdate.StartTime = appointment.StartTime
	appointmentUpdate.EndTime = appointment.EndTime
	appointmentUpdate.TotalAmount = appointment.TotalAmount

	if err := l.repo.RepositoryAppointment.Update(appointmentUpdate); err != nil {
		log.Printf("appointment-logic: Error update appointment con ID %d: %v", ID, err)
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

func (l *appointmentLogic) UpdateAppointmentWithService(ID uint, appointment *model.Appointment) (*model.FinalServicePrice, error) {
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return nil, err
	}

	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	if appointment.Patient.ID != 0 {
		if err := l.handlePatientIDCreation(appointment); err != nil {
			return nil, err
		}
	}

	if appointment.Patient != (model.Patient{}) {
		if err := l.handlePatientBodyCreation(appointment); err != nil {
			return nil, err
		}
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalServicePrice, err := l.isServiceIDEXists(appointment.ServiceID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentUpdate.Patient = appointment.Patient
	appointmentUpdate.DoctorID = appointment.DoctorID
	appointmentUpdate.ServiceID = appointment.ServiceID
	appointmentUpdate.Date = appointment.Date
	appointmentUpdate.StartTime = appointment.StartTime
	appointmentUpdate.EndTime = appointment.EndTime
	appointmentUpdate.TotalAmount = appointment.TotalAmount

	if err := l.repo.RepositoryAppointment.Update(appointmentUpdate); err != nil {
		log.Printf("appointment-logic: Error actualizando la cita con ID %d: %v", ID, err)
		return nil, err
	}

	return &model.FinalServicePrice{
		TotalAmount:       finalServicePrice.TotalAmount,
		InsuranceDiscount: finalServicePrice.InsuranceDiscount,
		FinalPrice:        finalServicePrice.FinalPrice,
	}, nil
}

/*func (l *appointmentLogic) parseStartAndEndTime(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	startTime, err := validate.ParseTime(startTimeStr)
	if err != nil {
		log.Printf("appointment-logic: Invalid start time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	endTime, err := validate.ParseTime(endTimeStr)
	if err != nil {
		log.Printf("appointment-logic: Invalid end time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}

func (l *appointmentLogic) parseTimesAndDate(appointment *model.Appointment) (time.Time, time.Time, time.Time, error) {
	startTime, endTime, err := l.parseStartAndEndTime(appointment.StartTime, appointment.EndTime)
	if err != nil {
		log.Printf("appointment-logic: Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, response.ErrorAppointmentInvalidDateFormat
	}

	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		log.Printf("appointment-logic: Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, response.ErrorAppointmentInvalidDateFormat
	}

	return startTime, endTime, appointmentDate, nil
}

/*func (l *appointmentLogic) ValidatePatient(patient *model.Patient) error {
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

/*func (l *appointmentLogic) ValidateUpdatedPatientFields(patient *model.Patient, patientUpdate *model.Patient) error {
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

/////////////////////////

// Verifica que el paciente esté en la cita, como ID o como cuerpo
/*func (l *appointmentLogic) isPatientExistsInAppointment(appointment *model.Appointment) bool {
	//verifica que esté en ID
	if patientExists := l.isPatientIDExists(appointment.PatientID); patientExists != nil {
		//appointment.Patient = *patientExists
		return true
	}

	//verifica que esté en el cuerpo
	if l.isPatientBodyExists(appointment) {
		return true
	}

	return false
}*/

/*func (l *appointmentLogic) isPatientIDExists(ID uint) (*model.Patient, error) {
	log.Println("appointment-logic: the request brought an ID patient")
	patient, err := l.repo.RepositoryPatient.GetByID(ID)
	if err != nil {
		log.Printf("appointment-logic: Error fetching patient by ID: %v", err)
		return nil, err
	}

	return patient, nil
}

func (l *appointmentLogic) isPatientBodyExists(appointment *model.Appointment) bool {
	return appointment.Patient == (model.Patient{})
}

func (l *appointmentLogic) handlePatientIDCreation(appointment *model.Appointment) error {
	patientFound, err := l.isPatientIDExists(appointment.PatientID)
	if err != nil {
		log.Printf("appointment-logic -> method handlePatientIDCreation: Patient not found with ID: %d", appointment.PatientID)
		return response.ErrorPatientNotFoundID
	}

	appointment.Patient = *patientFound

	return nil
}

func (l *appointmentLogic) handlePatientBodyCreation(appointment *model.Appointment) error {
	if !l.isPatientBodyExists(appointment) {
		return response.ErrorBodyPatientEmpty
	}

	if err := l.repo.RepositoryPatient.Create(&appointment.Patient); err != nil {
		log.Printf("appointment-logic: Error creating patient: %v", err)
		return err
	}

	return nil
}

/////////////////////////

/*func (l *appointmentLogic) isServiceIDEXists(ID uint, hasInsurance bool) (*model.FinalServicePrice, error) {
	if ID == 0 {
		return &model.FinalServicePrice{}, response.ErrorInvalidID
	}

	log.Println("appointment-logic: the request brought an ID service")

	service, err := l.repo.RepositoryService.GetByID(ID)
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

/*func (l *appointmentLogic) isPackageIDExists(ID uint, hasInsurance bool) (*model.FinalPackagePriceWithInsegurance, error) {
	pkg, err := l.repo.RepositoryPackageMain.GetByID(ID)
	if err != nil || pkg == nil {
		log.Printf("appointment-logic: The package with ID %d not exists: %v", ID, err)
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
}*/

/*func (l *appointmentLogic) validateDoctorAvailability(doctorID uint) error {
	doctor, err := l.repo.RepositoryDoctor.GetByID(doctorID)
	if err != nil || doctor == nil {
		return response.ErrorDoctorNotFoundID
	}

	return nil
}*/

/*func (l *appointmentLogic) validateAppointmentTime(appointment *model.Appointment) error {
	startTime, endTime, appointmentDate, err := l.parseTimesAndDate(appointment)
	if err != nil {
		return err
	}

	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment-logic: Appointment date is in the past")
		return response.ErrorAppointmentDateInPast
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment-logic: Start time is not before end time")
		return response.ErrorInvalidAppointmentTimeRange
	}

	existingAppointments, err := l.repo.RepositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointment.Date)
	if err != nil {
		log.Printf("appointment-logic: Error fetching existing appointments: %v", err)
		return err
	}

	if validate.HasTimeConflict(existingAppointments, startTime, endTime) {
		log.Printf("appointment-logic: Appointment time conflicts with existing appointments")
		return response.ErrorAppointmentTimeConflict
	}

	return nil
}*/
