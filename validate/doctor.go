package validate

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/dto"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type DoctorValidator interface {
	ValidateDoctor(doctorRequest dto.DoctorRequest) (*model.Doctor, error)
}

type doctorValidator struct{}

func DNIDoctor(doctor *model.Doctor) error {
	if CheckDNIExists[model.Doctor](doctor.DNI, doctor) {
		return response.ErrorDoctorExistsDNI
	}

	return nil
}

func PhoneNumberDoctor(doctor *model.Doctor) error {
	if CheckPhoneNumberExists[model.Doctor](doctor.PhoneNumber, doctor) {
		return response.ErrorDoctorExistsPhoneNumber
	}

	return nil
}

func EmailDoctor(doctor *model.Doctor) error {
	if CheckEmailExists[model.Doctor](doctor.Email, doctor) {
		return response.ErrorDoctorExistsEmail
	}

	return nil
}

func BirthDateDoctor(birthDateStr string) (string, error) {
	birthDate, err := ParseDate(birthDateStr)
	if err != nil {
		return "", response.ErrorDoctorInvalidDateFormat
	}

	if !IsDateInPast(birthDate) {
		return "", response.ErrorDoctorBirthDateIsFuture
	}

	return birthDateStr, nil
}
