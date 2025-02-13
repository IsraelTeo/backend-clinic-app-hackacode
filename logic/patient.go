package logic

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

type PatientLogic interface {
	GetPatientByID(ID uint) (*model.Patient, error)
	GetPatientByDNI(DNI string) (*model.Patient, error)
	GetAllPatients() ([]model.Patient, error)
	CreatePatient(patient *model.Patient) error
	UpdatePatient(ID uint, patient *model.Patient) error
	DeletePatient(ID uint) error
}

type patientLogic struct {
	repositoryPatient         repository.Repository[model.Patient]
	repositoryPatientMain     repository.PatientRepository
	repositoryAppointmentMain repository.AppointmentRepository
}

func NewPatientLogic(repositoryPatient repository.Repository[model.Patient],
	repositoryPatientMain repository.PatientRepository,
	repositoryAppointmentMain repository.AppointmentRepository,
) PatientLogic {
	return &patientLogic{
		repositoryPatient:         repositoryPatient,
		repositoryPatientMain:     repositoryPatientMain,
		repositoryAppointmentMain: repositoryAppointmentMain,
	}
}

func (l *patientLogic) GetPatientByID(ID uint) (*model.Patient, error) {
	patient, err := l.repositoryPatient.GetByID(ID)
	if err != nil {
		return nil, response.ErrorPatientNotFoundID
	}

	return patient, nil
}

func (l *patientLogic) GetPatientByDNI(DNI string) (*model.Patient, error) {
	patient, err := l.repositoryPatientMain.GetPatientByDNI(DNI)
	if err != nil {
		log.Printf("patient-logic: Error fetching patient with DNI %s: %v", DNI, err)
		return nil, response.ErrorPatientNotFoundDNI
	}

	return patient, nil
}

func (l *patientLogic) GetAllPatients() ([]model.Patient, error) {
	patients, err := l.repositoryPatient.GetAll()
	if err != nil {
		log.Printf("patient-logic: Error fetching patients: %v", err)
		return nil, response.ErrorPatientsNotFound
	}

	if len(patients) == 0 {
		log.Println("patient-logic: No patients found")
		return []model.Patient{}, nil
	}

	return patients, nil
}

func (l *patientLogic) CreatePatient(patient *model.Patient) error {
	if err := l.validatePatient(patient); err != nil {
		return err
	}

	if err := l.repositoryPatient.Create(patient); err != nil {
		log.Printf("patient-logic: Error saving patient: %v", err)
		return response.ErrorToCreatedPatient
	}

	return nil
}

func (l *patientLogic) UpdatePatient(ID uint, patient *model.Patient) error {
	patientUpdate, err := l.repositoryPatient.GetByID(ID)
	if err != nil {
		log.Printf("patient-logic: Error fetching patient with ID %d: %v", ID, err)
		return response.ErrorPatientNotFoundID
	}

	if err := l.validateUpdatedPatientFields(patient, patientUpdate); err != nil {
		return err
	}

	patientUpdate.Name = patient.Name
	patientUpdate.LastName = patient.LastName
	patientUpdate.DNI = patient.DNI
	patientUpdate.BirthDate = patient.BirthDate
	patientUpdate.Email = patient.Email
	patientUpdate.PhoneNumber = patient.PhoneNumber
	patientUpdate.Address = patient.Address
	patientUpdate.Insurance = patient.Insurance

	if err = l.repositoryPatient.Update(patientUpdate); err != nil {
		log.Printf("patient-logic: Error updating patient with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPatient
	}

	return nil
}

func (l *patientLogic) DeletePatient(ID uint) error {
	if _, err := l.repositoryPatient.GetByID(ID); err != nil {
		log.Printf("patient-logic: Error fetching patient with ID %d: %v to deleting", ID, err)
		return response.ErrorPatientNotFoundID
	}

	if err := l.repositoryAppointmentMain.UnlinkPatientAppointments(ID); err != nil {
		log.Printf("patient-logic: Error unlinking appointments for patient with ID %d: %v", ID, err)
		return response.ErrorUnlinkingAppointments
	}

	if err := l.repositoryPatient.Delete(ID); err != nil {
		log.Printf("patient-logic: Error deleting patient with ID %d: %v", ID, err)
		return response.ErrorToDeletedPatient
	}

	return nil
}

func (l *patientLogic) validatePatient(patient *model.Patient) error {
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

func (l *patientLogic) validateUpdatedPatientFields(patient *model.Patient, patientUpdate *model.Patient) error {
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
