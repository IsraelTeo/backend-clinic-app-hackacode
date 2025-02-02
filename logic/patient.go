package logic

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type PatientLogic interface {
	GetPatientByID(ID uint) (*model.Patient, error)
	GetAllPatients() ([]model.Patient, error)
	CreatePatient(patient *model.Patient) error
	UpdatePatient(ID uint, patient *model.Patient) error
	DeletePatient(ID uint) error
}

type patientLogic struct {
	repository repository.Repository[model.Patient]
}

func NewPatientLogic(repository repository.Repository[model.Patient]) PatientLogic {
	return &patientLogic{repository: repository}
}

func (l *patientLogic) GetPatientByID(ID uint) (*model.Patient, error) {
	patient, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("patient: Error fetching patient with ID %d: %v", ID, err)
		return nil, response.ErrorPatientNotFound
	}
	return patient, nil
}

func (l *patientLogic) GetAllPatients() ([]model.Patient, error) {
	patients, err := l.repository.GetAll()
	if err != nil {
		log.Printf("patient: Error fetching patients: %v", err)
		return nil, response.ErrorPatientNotFound
	}

	if len(patients) == 0 {
		log.Println("patient: No patients found")
		return []model.Patient{}, response.ErrorListPatientsEmpty
	}
	return patients, nil
}

func (l *patientLogic) CreatePatient(patient *model.Patient) error {
	if err := l.repository.Create(patient); err != nil {
		log.Printf("patient: Error saving patient: %v", err)
		return response.ErrorToCreatedPatient
	}
	return nil
}

func (l *patientLogic) UpdatePatient(ID uint, patient *model.Patient) error {
	patientUpdate, err := l.GetPatientByID(ID)
	if err != nil {
		log.Printf("patient: Error fetching patient with ID %d: %v to update", ID, err)
		return response.ErrorPatientNotFound
	}

	patientUpdate.Name = patient.Name
	patientUpdate.LastName = patient.LastName
	patientUpdate.DNI = patient.DNI
	patientUpdate.BirthDate = patient.BirthDate
	patientUpdate.Email = patient.Email
	patientUpdate.PhoneNumber = patient.PhoneNumber
	patientUpdate.Address = patient.Address
	patientUpdate.Insurance = patient.Insurance

	if err = l.repository.Update(patientUpdate); err != nil {
		log.Printf("patient: Error updating patient with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPatient
	}

	return nil
}

func (l *patientLogic) DeletePatient(ID uint) error {
	if err := l.repository.Delete(ID); err != nil {
		log.Printf("patient: Error deleting patient with ID %d: %v", ID, err)
		return response.ErrorToDeletedPatient
	}

	return nil
}
