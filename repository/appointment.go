package repository

import (
	"errors"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetByID(ID uint) (*model.Appointment, error)
	GetAll() ([]model.Appointment, error)
	GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error)
	GetAppointmentsByDoctorAndDate(doctorID uint, date string) ([]model.Appointment, error)
}
type appointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository crea una nueva instancia del repositorio
func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) GetByID(ID uint) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.
		Preload("Patient").  // Cargar información del paciente
		Where("id = ?", ID). // Filtro explícito por el campo ID
		First(&appointment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorAppointmentNotFound
		}
		return nil, err
	}

	return &appointment, nil
}

func (r *appointmentRepository) GetAll() ([]model.Appointment, error) {
	var appointments []model.Appointment
	err := r.db.
		Preload("Patient"). // Cargar información del paciente
		Find(&appointments).Error

	if err != nil {
		return nil, err
	}

	return appointments, nil
}

// Obtener citas de un doctor específico
func (r *appointmentRepository) GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

// Obtener citas de un médico en una fecha específica
func (r *appointmentRepository) GetAppointmentsByDoctorAndDate(doctorID uint, date string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := r.db.Where("doctor_id = ? AND date = ?", doctorID, date).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

/*
// GetByPatientID obtiene todas las citas asociadas a un paciente
func (r *AppointmentRepository) GetByPatientID(patientID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment

	// Consulta a la base de datos
	err := r.db.Where("patient_id = ?", patientID).Find(&appointments).Error
	if err != nil {
		log.Printf("repository: Error fetching appointments for patient ID %d: %v", patientID, err)
		return nil, err
	}

	return appointments, nil
}

// GetByDoctorID obtiene todas las citas asociadas a un doctor
func (r *AppointmentRepository) GetByDoctorID(doctorID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment

	// Consulta a la base de datos
	err := r.db.Where("doctor_id = ?", doctorID).Find(&appointments).Error
	if err != nil {
		log.Printf("repository: Error fetching appointments for doctor ID %d: %v", doctorID, err)
		return nil, err
	}

	return appointments, nil
}
*/
