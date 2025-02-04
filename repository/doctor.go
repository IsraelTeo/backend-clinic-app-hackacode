package repository

import "gorm.io/gorm"

type DoctorRepository interface{}

type doctorRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository crea una nueva instancia del repositorio
func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}
