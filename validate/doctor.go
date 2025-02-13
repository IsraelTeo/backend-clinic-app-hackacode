package validate

import (
	"log"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

func DNIDoctor(doctor *model.Doctor) error {
	if CheckDNIExists[model.Doctor](doctor.DNI, doctor) {
		log.Printf("validation: Error checking if doctor exists by DNI: %s", doctor.DNI)
		return response.ErrorDoctorExistsDNI
	}

	return nil
}

func PhoneNumberDoctor(doctor *model.Doctor) error {
	if CheckPhoneNumberExists[model.Doctor](doctor.PhoneNumber, doctor) {
		log.Printf("validation: Error checking if patient exists by phone number: %s", doctor.PhoneNumber)
		return response.ErrorDoctorExistsPhoneNumber
	}

	return nil
}

func EmailDoctor(doctor *model.Doctor) error {
	if CheckEmailExists[model.Doctor](doctor.Email, doctor) {
		log.Printf("validation: Error checking if doctor exists by email: %s", doctor.Email)
		return response.ErrorDoctorExistsEmail
	}

	return nil
}

func BirthDateDoctor(birthDateStr string) (*time.Time, error) {
	birthDate, err := ParseDate(birthDateStr)
	if err != nil {
		log.Printf("validation: Error parsing birthdate: %v", birthDateStr)
		return &birthDate, response.ErrorDoctorInvalidDateFormat
	}

	if !IsDateInPast(birthDate) {
		log.Printf("validation: Error birthdate is future: %v", birthDateStr)
		return nil, response.ErrorDoctorBirthDateIsFuture
	}

	return &birthDate, nil
}
