package validate

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

func DNIPatient(patient *model.Patient) error {
	if CheckDNIExists[model.Patient](patient.DNI, patient) {
		log.Printf("validation: Error checking if patient exists by DNI: %s", patient.DNI)
		return response.ErrorPatientExistsDNI
	}

	return nil
}

func PhoneNumberPatient(patient *model.Patient) error {
	if CheckPhoneNumberExists[model.Patient](patient.PhoneNumber, patient) {
		log.Printf("validation: Error checking if patient exists by phone number: %s", patient.PhoneNumber)
		return response.ErrorPatientExistsPhoneNumber
	}

	return nil
}

func EmailPatient(patient *model.Patient) error {
	if CheckEmailExists[model.Patient](patient.Email, patient) {
		log.Printf("validation: Error checking if patient exists by email: %s", patient.Email)
		return response.ErrorPatientExistsEmail
	}

	return nil
}

func BirthDatePatient(birthDateStr string) error {
	birthDate, err := ParseDate(birthDateStr)
	if err != nil {
		log.Printf("validation: Error parsing birthdate: %v", birthDateStr)
		return response.ErrorPatientInvalidDateFormat
	}

	if !IsDateInPast(birthDate) {
		log.Printf("validation: Error birthdate is future: %v", birthDateStr)
		return response.ErrorPatientBrithDateIsFuture
	}

	return nil
}

func PatientToCreate(patient *model.Patient) error {
	if err := DNIPatient(patient); err != nil {
		return err
	}

	if err := PhoneNumberPatient(patient); err != nil {
		return err
	}

	if err := EmailPatient(patient); err != nil {
		return err
	}

	if err := BirthDatePatient(patient.BirthDate); err != nil {
		return err
	}

	return nil
}

func PatientFieldsToUpdate(patient *model.Patient, patientUpdate *model.Patient) error {
	if patient.DNI != patientUpdate.DNI {
		if err := DNIPatient(patient); err != nil {
			return err
		}
	}

	if patient.PhoneNumber != patientUpdate.PhoneNumber {
		if err := PhoneNumberPatient(patient); err != nil {
			return err
		}
	}

	if patient.Email != patientUpdate.Email {
		if err := EmailPatient(patient); err != nil {
			return err
		}
	}

	return nil
}
