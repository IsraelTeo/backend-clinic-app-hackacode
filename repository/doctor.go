package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gorm.io/gorm"
)

type DoctorRepository interface {
	GetDoctorByDNI(DNI string) (*model.Doctor, error)
}

type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

func (r *doctorRepository) GetDoctorByDNI(DNI string) (*model.Doctor, error) {
	var doctor model.Doctor

	err := r.db.
		Where("dni = ?", DNI).
		First(&doctor).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, response.ErrorDoctorNotFoundDNI
		}

		return nil, err
	}

	return &doctor, nil
}
