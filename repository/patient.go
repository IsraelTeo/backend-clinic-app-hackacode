package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gorm.io/gorm"
)

type PatientRepository interface {
	GetPatientByDNI(DNI string) (*model.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) GetPatientByDNI(DNI string) (*model.Patient, error) {
	var patient model.Patient

	err := r.db.
		Where("dni = ?", DNI).
		First(&patient).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, response.ErrorPatientNotFoundDNI
		}

		return nil, err
	}

	return &patient, nil
}
