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
		log.Printf("patient: Error fetching patient with ID %d: %v", ID, err)
		return nil, response.ErrorPatientNotFoundID
	}

	return patient, nil
}

func (l *patientLogic) GetPatientByDNI(DNI string) (*model.Patient, error) {
	patient, err := l.repositoryPatientMain.GetPatientByDNI(DNI)
	if err != nil {
		log.Printf("patient: Error fetching patient with DNI %s: %v", DNI, err)
		return nil, response.ErrorPatientNotFoundDNI
	}

	return patient, nil
}

func (l *patientLogic) GetAllPatients() ([]model.Patient, error) {
	patients, err := l.repositoryPatient.GetAll()
	if err != nil {
		log.Printf("patient: Error fetching patients: %v", err)
		return nil, response.ErrorPatientsNotFound
	}

	if len(patients) == 0 {
		log.Println("patient: No patients found")
		return []model.Patient{}, response.ErrorListPatientsEmpty
	}

	return patients, nil
}

func (l *patientLogic) CreatePatient(patient *model.Patient) error {
	if validate.CheckDNIExists[model.Patient](patient.DNI, patient) {
		log.Printf("Error checking if patient exists by DNI: %s", patient.DNI)
		return response.ErrorPatientExistsDNI
	}

	if validate.CheckPhoneNumberExists[model.Patient](patient.PhoneNumber, patient) {
		log.Printf("Error checking if patient exists by phone number: %s", patient.PhoneNumber)
		return response.ErrorPatientExistsPhoneNumber
	}

	if validate.CheckEmailExists[model.Patient](patient.Email, patient) {
		log.Printf("Error checking if patient exists by email: %s", patient.Email)
		return response.ErrorPatientExistsEmail
	}

	birthDate, err := validate.ParseDate(patient.BirthDate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", patient.BirthDate)
		return response.ErrorPatientInvalidDateFormat
	}

	if !validate.IsDateInPast(birthDate) {
		log.Printf("Error birthdate is past: %v", patient.BirthDate)
		return response.ErrorPatientBrithDateIsFuture
	}

	if err := l.repositoryPatient.Create(patient); err != nil {
		log.Printf("patient: Error saving patient: %v", err)
		return response.ErrorToCreatedPatient
	}

	return nil
}

func (l *patientLogic) UpdatePatient(ID uint, patient *model.Patient) error {
	patientUpdate, err := l.repositoryPatient.GetByID(ID)
	if err != nil {
		log.Printf("patient: Error fetching patient with ID %d: %v to update", ID, err)
		return response.ErrorPatientNotFoundID
	}

	if validate.CheckDNIExists[model.Patient](patient.DNI, patient) {
		log.Printf("Error checking if patient exists by DNI: %s", patient.DNI)
		return response.ErrorPatientExistsDNI
	}

	if validate.CheckPhoneNumberExists[model.Patient](patient.PhoneNumber, patient) {
		log.Printf("Error checking if patient exists by phone number: %s", patient.PhoneNumber)
		return response.ErrorPatientExistsPhoneNumber
	}

	birthDate, err := validate.ParseDate(patient.BirthDate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", patient.BirthDate)
		return response.ErrorPatientInvalidDateFormat
	}

	if !validate.IsDateInPast(birthDate) {
		log.Printf("Error birthdate is past: %v", patient.BirthDate)
		return response.ErrorPatientBrithDateIsFuture
	}

	if validate.CheckEmailExists[model.Patient](patient.Email, patient) {
		log.Printf("Error checking if patient exists by email: %s", patient.Email)
		return response.ErrorPatientExistsEmail
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
		log.Printf("patient: Error updating patient with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPatient
	}

	return nil
}

func (l *patientLogic) DeletePatient(ID uint) error {
	if _, err := l.repositoryPatient.GetByID(ID); err != nil {
		log.Printf("patient: Error fetching patient with ID %d: %v to deleting", ID, err)
		return response.ErrorPatientNotFoundID
	}

	if err := l.repositoryAppointmentMain.UnlinkPatientAppointments(ID); err != nil {
		log.Printf("patient: Error unlinking appointments for patient with ID %d: %v", ID, err)
		return response.ErrorUnlinkingAppointments
	}

	if err := l.repositoryPatient.Delete(ID); err != nil {
		log.Printf("patient: Error deleting patient with ID %d: %v", ID, err)
		return response.ErrorToDeletedPatient
	}

	return nil
}
