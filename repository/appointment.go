package repository

import (
	"errors"
	"time"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetByID(ID uint) (*model.Appointment, error)
	GetAll(limit, offset int) ([]model.Appointment, error)
	GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error)
	GetAppointmentsByDoctorAndDate(doctorID uint, date time.Time) ([]model.Appointment, error)
	UpdatePaid(appointmentID uint) error
	UnlinkPatientAppointments(patientID uint) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) GetByID(ID uint) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.
		Preload("Patient").
		Where("id = ?", ID).
		First(&appointment).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorAppointmentNotFound
		}

		return nil, err
	}

	return &appointment, nil
}

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
}

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
