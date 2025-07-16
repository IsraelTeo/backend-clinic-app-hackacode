package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// PatientRepository defines the interface for patient operations
type PatientRepository interface {
	GetByID(ID uint) (*model.Patient, error)
	Get() ([]model.Patient, error)
	Create(patient *model.Patient) error
	Update(patient *model.Patient) error
	Delete(ID uint) error
	GetByDNI(DNI string) (*model.Patient, error)
}

// patientRepository implements the PatientRepository interface
type patientRepository struct {
	db *gorm.DB
}

// Dependency injection for the patient repository
func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

// GetByID retrieves a patient by its ID
func (r *patientRepository) GetByID(ID uint) (*model.Patient, error) {
	var patient model.Patient

	if err := r.db.First(&patient, ID).Error; err != nil {
		return nil, err
	}

	return &patient, nil
}

// Get retrieves patients
func (r *patientRepository) Get() ([]model.Patient, error) {
	var patients []model.Patient

	if err := r.db.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

// Create adds a new patient to the database
func (r *patientRepository) Create(patient *model.Patient) error {
	if err := r.db.Create(patient).Error; err != nil {
		return err
	}

	return nil
}

// Update modifies an existing patient
func (r *patientRepository) Update(patient *model.Patient) error {
	if err := r.db.Save(patient).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes a patient by its ID
func (r *patientRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.Patient{}, ID).Error; err != nil {
		return err
	}

	return nil
}

// GetByDNI retrieves a patient by their DNI
func (r *patientRepository) GetByDNI(DNI string) (*model.Patient, error) {
	var patient model.Patient

	err := r.db.
		Where("dni = ?", DNI).
		First(&patient).
		Error
	if err != nil {
		return nil, err
	}

	return &patient, nil
}
