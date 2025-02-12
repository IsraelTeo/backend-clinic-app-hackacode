package validate

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

func ValidateDNI(patient *model.Patient) error {
	if CheckDNIExists[model.Patient](patient.DNI, patient) {
		log.Printf("validation: Error checking if patient exists by DNI: %s", patient.DNI)
		return response.ErrorPatientExistsDNI
	}

	return nil
}

func ValidatePhoneNumber(patient *model.Patient) error {
	if CheckPhoneNumberExists[model.Patient](patient.PhoneNumber, patient) {
		log.Printf("validation: Error checking if patient exists by phone number: %s", patient.PhoneNumber)
		return response.ErrorPatientExistsPhoneNumber
	}

	return nil
}

func ValidateEmail(patient *model.Patient) error {
	if CheckEmailExists[model.Patient](patient.Email, patient) {
		log.Printf("validation: Error checking if patient exists by email: %s", patient.Email)
		return response.ErrorPatientExistsEmail
	}

	return nil
}

func ValidateBirthDate(birthDateStr string) error {
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
