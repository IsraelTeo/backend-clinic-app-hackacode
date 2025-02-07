package repository

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type PatientRepository interface {
	GetPatientByDNI(dni string) (*model.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository crea una nueva instancia del repositorio
func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

// Buscar paciente por DNI
func (r *patientRepository) GetPatientByDNI(dni string) (*model.Patient, error) {
	var patient model.Patient
	if err := r.db.Where("dni = ?", dni).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}
