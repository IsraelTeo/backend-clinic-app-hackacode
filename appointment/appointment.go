package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentLogic interface {
	GetAppointmentByID(ID uint) (*model.Appointment, error)
	GetAllAppointments(limit, offset int) ([]model.Appointment, error)
	CreateAppointment(appointment *model.Appointment) (model.PriceDetails, error)
	UpdateAppointment(ID uint, appointment *model.Appointment) (model.PriceDetails, error)
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

func (l *appointmentLogic) GetAppointmentByID(ID uint) (*model.Appointment, error) {
	appointment, err := l.repositoryAppointmentMain.GetByID(ID)
	if err != nil {
		log.Printf("appointment-logic: Error fetching appointment with ID %d: %v", ID, err)
		return nil, response.ErrorAppointmentNotFound
	}

	return appointment, nil
}

func (l *appointmentLogic) GetAllAppointments(limit, offset int) ([]model.Appointment, error) {
	appointments, err := l.repositoryAppointmentMain.GetAll(limit, offset)
	if err != nil {
		log.Printf("appointment-logic: Error fetching appointments: %v", err)
		return nil, response.ErrorAppointmetsNotFound
	}

	return appointments, nil
}

func (l *appointmentLogic) CreateAppointment(appointment *model.Appointment) (model.PriceDetails, error) {
	finalPrice, err := l.logicAppointmentCreate.CreateAppointment(appointment)
	if err != nil {
		log.Printf("appointment-logic -> method: CreateAppointment: Error to create: %v", err)
		return nil, err
	}

	return finalPrice, nil
}

func (l *appointmentLogic) UpdateAppointment(ID uint, appointment *model.Appointment) (model.PriceDetails, error) {
	finalPrice, err := l.logicAppointmentUpdate.UpdateAppointment(ID, appointment)
	if err != nil {
		log.Printf("appointment-logic -> method: UpdateAppointment: Error to update: %v", err)
		return nil, err
	}

	return finalPrice, nil
}

func (l *appointmentLogic) DeleteAppointment(ID uint) error {
	_, err := l.GetAppointmentByID(ID)
	if err != nil {
		return response.ErrorAppointmentNotFound
	}

	err = l.repositoryAppointment.Delete(ID)
	if err != nil {
		return response.ErrorToDeletedAppointment
	}

	return nil
}
