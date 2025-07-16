package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// DoctorRepository defines the interface for doctor operations
type DoctorRepository interface {
	GetByID(ID uint) (*model.Doctor, error)
	Get() ([]model.Doctor, error)
	Create(doctor *model.Doctor) error
	Update(doctor *model.Doctor) error
	Delete(ID uint) error
	GetByDNI(DNI string) (*model.Doctor, error)
}

// doctorRepository implements the DoctorRepository interface
type doctorRepository struct {
	db *gorm.DB
}

// Dependency injection for the doctor repository
func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

// GetByID retrieves a doctor by its ID
func (r *doctorRepository) GetByID(ID uint) (*model.Doctor, error) {
	var doctor model.Doctor

	if err := r.db.First(&doctor, ID).Error; err != nil {
		return nil, err
	}

	return &doctor, nil
}

// Get retrieves all doctors
func (r *doctorRepository) Get() ([]model.Doctor, error) {
	var doctors []model.Doctor

	if err := r.db.Find(&doctors).Error; err != nil {
		return nil, err
	}

	return doctors, nil
}

// Create adds a new doctor to the database
func (r *doctorRepository) Create(doctor *model.Doctor) error {
	if err := r.db.Create(doctor).Error; err != nil {
		return err
	}

	return nil
}

// Update modifies an existing doctor
func (r *doctorRepository) Update(doctor *model.Doctor) error {
	if err := r.db.Save(doctor).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes a doctor by its ID
func (r *doctorRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.Doctor{}, ID).Error; err != nil {
		return err
	}

	return nil
}

// GetByDNI retrieves a doctor by their DNI
func (r *doctorRepository) GetByDNI(DNI string) (*model.Doctor, error) {
	var doctor model.Doctor

	err := r.db.
		Where("dni = ?", DNI).
		First(&doctor).
		Error

	if err != nil {
		return nil, err
	}

	return &doctor, nil
}
