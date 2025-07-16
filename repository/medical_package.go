package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// MedicalPackageRepository defines the interface for medical package operations
type MedicalPackageRepository interface {
	GetByID(ID uint) (*model.MedicalPackage, error)
	Get() ([]model.MedicalPackage, error)
	Create(medicalPackage *model.MedicalPackage) error
	Update(medicalPackage *model.MedicalPackage) error
	Delete(ID uint) error
}

// medicalPackageRepository implements the MedicalPackageRepository interface
type medicalPackageRepository struct {
	db *gorm.DB
}

// Dependency injection for the medical package repository
func NewMedicalPackageRepository(db *gorm.DB) MedicalPackageRepository {
	return &medicalPackageRepository{db: db}
}

// GetByID retrieves a medical package by its ID
func (r *medicalPackageRepository) GetByID(ID uint) (*model.MedicalPackage, error) {
	var medicalPackage model.MedicalPackage

	if err := r.db.First(&medicalPackage, ID).Error; err != nil {
		return nil, err
	}

	return &medicalPackage, nil
}

// Get retrieves medical packages
func (r *medicalPackageRepository) Get() ([]model.MedicalPackage, error) {
	var medicalPackages []model.MedicalPackage

	if err := r.db.Find(&medicalPackages).Error; err != nil {
		return nil, err
	}

	return medicalPackages, nil
}

// Create adds a new medical package to the database
func (r *medicalPackageRepository) Create(medicalPackage *model.MedicalPackage) error {
	if err := r.db.Create(medicalPackage).Error; err != nil {
		return err
	}

	return nil
}

// Update modifies an existing medical package
func (r *medicalPackageRepository) Update(medicalPackage *model.MedicalPackage) error {
	if err := r.db.Save(medicalPackage).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes a medical package by its ID
func (r *medicalPackageRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.MedicalPackage{}, ID).Error; err != nil {
		return err
	}

	return nil
}
