package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// AppointmentRepository defines the interface for appointment operations
type AppointmentRepository interface {
	GetByID(ID uint) (*model.Appointment, error)
	Get() ([]model.Appointment, error)
	Create(appointment *model.Appointment) error
	Update(appointment *model.Appointment) error
	Delete(ID uint) error
	//GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error)
	//GetAppointmentsByDoctorAndDate(doctorID uint, date time.Time) ([]model.Appointment, error)
	//UnlinkPatientAppointments(patientID uint) error
}

// appointmentRepository implements the AppointmentRepository interface
type appointmentRepository struct {
	db *gorm.DB
}

// Dependency injection for the appointment repository
func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

// GetByID retrieves an appointment by its ID
func (r *appointmentRepository) GetByID(ID uint) (*model.Appointment, error) {
	var appointment model.Appointment

	if err := r.db.First(&appointment, ID).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

// Get retrieves appointments
func (r *appointmentRepository) Get() ([]model.Appointment, error) {
	var appointments []model.Appointment

	if err := r.db.Preload("Patient").Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

// Create adds a new appointment to the database
func (r *appointmentRepository) Create(appointment *model.Appointment) error {
	if err := r.db.Create(appointment).Error; err != nil {
		return err
	}

	return nil
}

// Update modifies an existing appointment
func (r *appointmentRepository) Update(appointment *model.Appointment) error {
	if err := r.db.Save(appointment).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes an appointment by its ID
func (r *appointmentRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.Appointment{}, ID).Error; err != nil {
		return err
	}

	return nil
}

/*
func (r *appointmentRepository) GetAll(limit, offset int) ([]model.Appointment, error) {
	var appointments []model.Appointment
	query := r.db.Preload("Patient")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}*/

/*
func (r *appointmentRepository) GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment

	err := r.db.
		Where("doctor_id = ?", doctorID).
		Find(&appointments).
		Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepository) GetAppointmentsByDoctorAndDate(doctorID uint, date time.Time) ([]model.Appointment, error) {
	var appointments []model.Appointment

	dateStr := validate.FormatDate(date)

	err := r.db.
		Where("doctor_id = ? AND date = ?", doctorID, dateStr).
		Find(&appointments).
		Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepository) UpdatePaid(appointmentID uint) error {
	appointment := model.Appointment{}

	err := r.db.First(&appointment, appointmentID).Error
	if err != nil {
		return err
	}

	appointment.Paid = true

	err = r.db.Save(&appointment).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *appointmentRepository) UnlinkPatientAppointments(patientID uint) error {
	err := r.db.
		Model(&model.Appointment{}).
		Where("patient_id = ?", patientID).
		Update("patient_id", nil).
		Error
	if err != nil {
		return err
	}

	return nil
}
*/
