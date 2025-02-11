package repository

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type DoctorRepository interface {
	GetDoctorByDNI(DNI string) (*model.Doctor, error)
}

type doctorRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository crea una nueva instancia del repositorio
func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

func (r *doctorRepository) GetDoctorByDNI(DNI string) (*model.Doctor, error) {
	var doctor model.Doctor
	if err := r.db.Where("dni = ?", DNI).First(&doctor).Error; err != nil {
		return nil, err
	}

	return &doctor, nil
}
