package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// MedicalServiceRepository defines the interface for medical service operations
type MedicalServiceRepository interface {
	GetByID(ID uint) (*model.MedicalService, error)
	Get() ([]model.MedicalService, error)
	Create(medicalService *model.MedicalService) error
	Update(medicalService *model.MedicalService) error
	Delete(ID uint) error
}

// medicalServiceRepository implements the MedicalServiceRepository interface
type medicalServiceRepository struct {
	db *gorm.DB
}

// Dependency injection for the medical service repository
func NewMedicalServiceRepository(db *gorm.DB) MedicalServiceRepository {
	return &medicalServiceRepository{db: db}
}

// GetByID retrieves a medical service by its ID
func (r *medicalServiceRepository) GetByID(ID uint) (*model.MedicalService, error) {
	var medicalService model.MedicalService

	if err := r.db.First(&medicalService, ID).Error; err != nil {
		return nil, err
	}

	return &medicalService, nil
}

// Get retrieves medical services
func (r *medicalServiceRepository) Get() ([]model.MedicalService, error) {
	var medicalServices []model.MedicalService

	if err := r.db.Find(&medicalServices).Error; err != nil {
		return nil, err
	}

	return medicalServices, nil
}

// Create adds a new medical service to the database
func (r *medicalServiceRepository) Create(medicalService *model.MedicalService) error {
	if err := r.db.Create(medicalService).Error; err != nil {
		return err
	}
	return nil
}

// Update modifies an existing medical service
func (r *medicalServiceRepository) Update(medicalService *model.MedicalService) error {
	if err := r.db.Save(medicalService).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes a medical service by its ID
func (r *medicalServiceRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.MedicalService{}, ID).Error; err != nil {
		return err
	}

	return nil
}
