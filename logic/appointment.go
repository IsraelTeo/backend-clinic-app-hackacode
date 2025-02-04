package logic

import (
	"fmt"
	"log"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentLogic interface {
	GetAppointmentByID(ID uint) (*model.Appointment, error)
	GetAllAppointments() ([]model.Appointment, error)
	CreateAppointment(appointment *model.Appointment) error
	UpdateAppointment(ID uint, appointment *model.Appointment) error
	DeleteAppointment(ID uint) error
}

type appointmentLogic struct {
	repositoryAppointment     repository.Repository[model.Appointment]
	repositoryAppointmentMain repository.AppointmentRepository
	repositoryDoctor          repository.Repository[model.Doctor]
	repositoryPatient         repository.Repository[model.Patient]
}

func NewAppointmentLogic(
	repositoryAppointment repository.Repository[model.Appointment],
	repositoryDoctor repository.Repository[model.Doctor],
	repositoryPatient repository.Repository[model.Patient],
	repositoryAppointmentMain repository.AppointmentRepository,
) AppointmentLogic {
	return &appointmentLogic{
		repositoryAppointment:     repositoryAppointment,
		repositoryDoctor:          repositoryDoctor,
		repositoryPatient:         repositoryPatient,
		repositoryAppointmentMain: repositoryAppointmentMain,
	}
}

func (l *appointmentLogic) GetAppointmentByID(ID uint) (*model.Appointment, error) {
	appointment, err := l.repositoryAppointment.GetByID(ID)
	if err != nil {
		log.Printf("appointment: Error fetching appointment with ID %d: %v", ID, err)
		return nil, response.ErrorAppointmentNotFound
	}
	return appointment, nil
}

func (l *appointmentLogic) GetAllAppointments() ([]model.Appointment, error) {
	appointments, err := l.repositoryAppointment.GetAll()
	if err != nil {
		log.Printf("appointment: Error fetching appointments: %v", err)
		return nil, response.ErrorAppointmetsNotFound
	}

	if len(appointments) == 0 {
		log.Println("appointment : No appointments found")
		return []model.Appointment{}, response.ErrorListAppointmentsEmpty
	}

	return appointments, nil
}

func (l *appointmentLogic) CreateAppointment(appointment *model.Appointment) error {
	// Verificar que el doctor exista
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return response.ErrorDoctorNotFound
	}

	// Verificar que la cita sea en el futuro
	if appointment.Date.Before(time.Now()) {
		return fmt.Errorf("La cita debe ser en una fecha futura")
	}

	// Obtener el día de la semana de la cita
	appointmentDay := appointment.Date.Weekday().String()

	// Verificar si el día está dentro de los días que trabaja el doctor
	isAvailable := false
	for _, day := range doctor.Days {
		if string(day) == appointmentDay {
			isAvailable = true
			break
		}
	}
	if !isAvailable {
		return fmt.Errorf("El doctor no atiende el día seleccionado: %s", appointmentDay)
	}

	// Obtener citas existentes del doctor
	existingAppointmentsDoctor, err := l.repositoryAppointmentMain.GetAppointmentsByDoctor(appointment.DoctorID)
	if err != nil {
		return fmt.Errorf("Error obteniendo citas del doctor: %v", err)
	}

	// Verificar conflictos de horario
	for _, existing := range existingAppointmentsDoctor {
		if existing.Date.Equal(appointment.Date) && (appointment.StartTime.Before(existing.EndTime) && appointment.EndTime.After(existing.StartTime)) {
			return fmt.Errorf("El doctor ya tiene una cita en este horario")
		}
	}

	// Crear la cita
	if err := l.repositoryAppointment.Create(appointment); err != nil {
		log.Printf("appointment: Error creating appointment: %v", err)
		return err
	}

	return nil
}

func (l *appointmentLogic) UpdateAppointment(ID uint, appointment *model.Appointment) error {
	// Obtener la cita actual por ID
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return err
	}

	// Verificar que el doctor exista
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil || doctor == nil {
		log.Printf("appointment: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return response.ErrorDoctorNotFound
	}

	// Verificar que la cita sea en el futuro
	if appointment.Date.Before(time.Now()) {
		return fmt.Errorf("La cita debe ser en una fecha futura")
	}

	// Obtener el día de la semana de la cita
	appointmentDay := appointment.Date.Weekday().String()

	// Verificar si el día está dentro de los días que trabaja el doctor
	isAvailable := false
	for _, day := range doctor.Days {
		if string(day) == appointmentDay {
			isAvailable = true
			break
		}
	}
	if !isAvailable {
		return fmt.Errorf("El doctor no atiende el día seleccionado: %s", appointmentDay)
	}

	// Obtener citas existentes del doctor
	existingAppointmentsDoctor, err := l.repositoryAppointmentMain.GetAppointmentsByDoctor(appointment.DoctorID)
	if err != nil {
		return fmt.Errorf("Error obteniendo citas del doctor: %v", err)
	}

	// Verificar conflictos de horario
	for _, existing := range existingAppointmentsDoctor {
		if existing.Date.Equal(appointment.Date) && (appointment.StartTime.Before(existing.EndTime) && appointment.EndTime.After(existing.StartTime)) {
			return fmt.Errorf("El doctor ya tiene una cita en este horario")
		}
	}
	// Actualizar la cita
	appointment.ID = appointmentUpdate.ID
	if err := l.repositoryAppointment.Update(appointment); err != nil {
		log.Printf("appointment: Error updating appointment with ID %d: %v", ID, err)
		return err
	}

	return nil
}

func (l *appointmentLogic) DeleteAppointment(ID uint) error {
	if err := l.repositoryAppointment.Delete(ID); err != nil {
		return err
	}

	return nil
}
