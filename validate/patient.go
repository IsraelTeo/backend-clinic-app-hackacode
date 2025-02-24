package validate

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

func DNIPatient(patient *model.Patient) error {
	if CheckDNIExists[model.Patient](patient.DNI, patient) {
		return response.ErrorPatientExistsDNI
	}

	return nil
}

func PhoneNumberPatient(patient *model.Patient) error {
	if CheckPhoneNumberExists[model.Patient](patient.PhoneNumber, patient) {
		return response.ErrorPatientExistsPhoneNumber
	}

	return nil
}

func EmailPatient(patient *model.Patient) error {
	if CheckEmailExists[model.Patient](patient.Email, patient) {
		return response.ErrorPatientExistsEmail
	}

	return nil
}

func BirthDatePatient(birthDateStr string) error {
	birthDate, err := ParseDate(birthDateStr)
	if err != nil {
		return response.ErrorPatientInvalidDateFormat
	}

	if !IsDateInPast(birthDate) {
		return response.ErrorPatientBrithDateIsFuture
	}

	return nil
}

func PatientToCreate(patient *model.Patient) error {
	err := DNIPatient(patient)
	if err != nil {
		return err
	}

	err = PhoneNumberPatient(patient)
	if err != nil {
		return err
	}

	err = EmailPatient(patient)
	if err != nil {
		return err
	}

	err = BirthDatePatient(patient.BirthDate)
	if err != nil {
		return err
	}

	return nil
}

func PatientToUpdate(patient *model.Patient, patientUpdate *model.Patient) error {
	if patient.DNI != patientUpdate.DNI {
		err := DNIPatient(patientUpdate)
		if err != nil {
			return err
		}
	}

	if patient.PhoneNumber != patientUpdate.PhoneNumber {
		err := PhoneNumberPatient(patientUpdate)
		if err != nil {
			return err
		}
	}

	if patient.Email != patientUpdate.Email {
		err := EmailPatient(patientUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}
