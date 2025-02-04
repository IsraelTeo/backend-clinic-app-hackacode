package repository

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error)
}
type appointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository crea una nueva instancia del repositorio
func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

// Obtener citas de un doctor específico
func (r *appointmentRepository) GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&appointments).Error; err != nil {
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
