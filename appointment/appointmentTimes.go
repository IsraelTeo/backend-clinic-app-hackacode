package appointment

import (
	"log"
	"strings"
	"time"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

type AppointmentTime interface {
	ValidateAppointmentTime(appointment *model.Appointment) error
}

type appointmentTime struct {
	repositoryAppointmentMain repository.AppointmentRepository
	repositoryDoctor          repository.Repository[model.Doctor]
}

func NewAppointmentTime(repositoryAppointmentMain repository.AppointmentRepository, repositoryDoctor repository.Repository[model.Doctor]) AppointmentTime {
	return &appointmentTime{repositoryAppointmentMain: repositoryAppointmentMain, repositoryDoctor: repositoryDoctor}
}

func (l *appointmentTime) parseStartAndEndTime(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	log.Printf("appointment-times-logic -> method: parseStartAndTime: received")

	startTime, err := validate.ParseTime(startTimeStr)
	if err != nil {
		log.Printf("appointment-times-logic -> method: ParseStartAndTime: Invalid start time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	endTime, err := validate.ParseTime(endTimeStr)
	if err != nil {
		log.Printf("appointment-times-logic -> method: parseStartAndTime: Invalid end time format: %v", err)
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}

func (l *appointmentTime) parseTimesAndDate(appointment *model.Appointment) (time.Time, time.Time, time.Time, error) {
	log.Printf("appointment-times-logic -> method: parseTimesAndDate: received")

	startTime, endTime, err := l.parseStartAndEndTime(appointment.StartTime, appointment.EndTime)
	if err != nil {
		log.Printf("appointment-times-logic -> method: parseTimesAndDate: Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		log.Printf("appointment-times-logic v-> method: parseTimesAndDate : Invalid appointment date format: %v", err)
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	return startTime, endTime, appointmentDate, nil
}

func (l *appointmentTime) hasTimeConflict(appointment *model.Appointment, startTimeAppointment, endTimeAppointment, appointmentDate time.Time) error {
	log.Printf("appointment-times-logic -> method:  hasTimeConflict: received")

	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil {
		log.Printf("appointment-times-logic -> method:  hasTimeConflict: Error doctor fetching with ID: %d error: %v", appointment.DoctorID, err)
		return response.ErrorDoctorNotFoundID
	}

	workingDays := strings.Split(doctor.Days, ",")

	appointmentWeekDaySpanish := validate.DayToGolang[appointmentDate.Weekday()]

	//verifica que el día de la cita sea un día que el médico tenga turno
	if !validate.IsDayAvailable(appointmentWeekDaySpanish, workingDays) {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Conflict detected - Doctor does not work on this day: %v", appointment.Date)
		return response.ErrorAppointmentDayNotAvailable
	}

	doctorStartTime, doctorEndTime, err := l.parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Error parsing doctor's working hours: %v", err)
		return err
	}

	//verifica que el horario de la cita esté dentro del turno del doctor
	if !validate.IsWithinTimeRange(startTimeAppointment, endTimeAppointment, doctorStartTime, doctorEndTime, true, true) {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Conflict detected - Appointment time (%v - %v) is outside doctor's working hours (%v - %v)",
			startTimeAppointment, endTimeAppointment, doctorStartTime, doctorEndTime)
		return response.ErrorAppointmentTimeConflict
	}

	doctorAppointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointmentDate)
	if err != nil {
		log.Printf("appointment-times-logic -> method:  hasTimeConflict: Error doctor fetching appoinments: error: %v", err)
		return response.ErrorFetchingAppointments
	}

	//para validar que la cita nueva no cruce horario con otra cita del médico
	/*for _, doctorAppointment := range doctorAppointments {
		if doctorAppointment.Date != appointmentDate.Format("2006-01-02") {
			continue
		}

		parsedDoctorStartTime, parsedDoctorEndTime, err := l.parseStartAndEndTime(doctorAppointment.StartTime, doctorAppointment.EndTime)
		if err != nil {
			log.Println("❌ Error al parsear los horarios de las citas existentes")
			return err
		}

		appointmentStartHour := startTimeAppointment.Hour()
		appointmentEndHour := endTimeAppointment.Hour()
		doctorStartHour := parsedDoctorStartTime.Hour()
		doctorEndHour := parsedDoctorEndTime.Hour()

		if appointmentStartHour < doctorEndHour || appointmentEndHour > doctorStartHour {
			return response.ErrorAppointmentTimeConflict
		}
	}*/

	for _, doctorAppointment := range doctorAppointments {
		if doctorAppointment.Date != appointmentDate.Format("2006-01-02") {
			continue
		}

		parsedDoctorStartTime, parsedDoctorEndTime, err := l.parseStartAndEndTime(doctorAppointment.StartTime, doctorAppointment.EndTime)
		if err != nil {
			return err
		}

		appointmentStartTime := startTimeAppointment
		appointmentEndTime := endTimeAppointment

		if appointmentStartTime.Before(parsedDoctorEndTime) && appointmentEndTime.After(parsedDoctorStartTime) {
			return response.ErrorAppointmentTimeConflict
		}
	}

	return nil
}

func (l *appointmentTime) ValidateAppointmentTime(appointment *model.Appointment) error {
	log.Printf("appointment-times-logic -> method: ValidateAppointmentTime: received")

	startTime, endTime, appointmentDate, err := l.parseTimesAndDate(appointment)
	if err != nil {
		return err
	}

	log.Printf("Nueva cita: Fecha: %v, Inicio: %v, Fin: %v", appointment.Date, startTime, endTime)

	if validate.IsDateInPast(appointmentDate) {
		log.Printf("appointment-times-logic -> method: ValidateAppointmentTime: Appointment date is in the past: %v", appointmentDate)
		return response.ErrorAppointmentDateInPast
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		log.Printf("appointment-times-logic -> method: ValidateAppointmentTime: Start time is not before end time: %v and %v", startTime, endTime)
		return response.ErrorInvalidAppointmentTimeRange
	}

	if err := l.hasTimeConflict(appointment, startTime, endTime, appointmentDate); err != nil {
		log.Printf("appointment-times-logic -> method: ValidateAppointmentTime: Appointment time conflicts with existing appointments")
		return err
	}

	return nil
}
