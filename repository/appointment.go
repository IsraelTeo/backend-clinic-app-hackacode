package repository

import (
	"errors"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetByID(ID uint) (*model.Appointment, error)
	GetAll() ([]model.Appointment, error)
	GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error)
	GetAppointmentsByDoctorAndDate(doctorID uint, date string) ([]model.Appointment, error)
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
	err := r.db.Preload("Patient").Where("id = ?", ID).First(&appointment).Error

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
	err := r.db.Preload("Patient").Find(&appointments).Error

	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepository) GetAppointmentsByDoctor(doctorID uint) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepository) GetAppointmentsByDoctorAndDate(doctorID uint, date string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := r.db.Where("doctor_id = ? AND date = ?", doctorID, date).Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepository) UpdatePaid(appointmentID uint) error {
	appointment := model.Appointment{}
	if err := r.db.First(&appointment, appointmentID).Error; err != nil {
		return err
	}

	appointment.Paid = true

	if err := r.db.Save(&appointment).Error; err != nil {
		return err
	}

	return nil
}

func (r *appointmentRepository) UnlinkPatientAppointments(patientID uint) error {
	if err := r.db.Model(&model.Appointment{}).Where("patient_id = ?", patientID).Update("patient_id", nil).Error; err != nil {
		return err
	}

	return nil
}
