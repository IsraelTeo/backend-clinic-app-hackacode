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

func (l *appointmentTime) hasTimeConflict(appointment *model.Appointment, startTimeAppointment, endTimeAppointment, appointmentDate time.Time) bool {
	log.Printf("appointment-times-logic -> method:  hasTimeConflict: received")

	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil {
		log.Printf("appointment-times-logic -> method:  hasTimeConflict: Error doctor fetching with ID: %d error: %v", appointment.DoctorID, err)
		return true
	}

	appointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctor(appointment.DoctorID)
	if err != nil {
		log.Printf("appointment-times-logic -> method:  hasTimeConflict: Error doctor fetching appoinments: error: %v", err)
		return true
	}

	workingDays := strings.Split(doctor.Days, ",")

	appointmentWeekDay := validate.TranslateDay(time.Weekday(appointmentDate.Weekday()).String())

	if !validate.IsDayAvailable(appointmentWeekDay, workingDays) {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Conflict detected - Doctor does not work on this day: %v", appointmentWeekDay)
		return true
	}

	doctorStartTime, doctorEndTime, err := l.parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Error parsing doctor's working hours: %v", err)
		return true
	}

	if !validate.IsWithinTimeRange(startTimeAppointment, endTimeAppointment, doctorStartTime, doctorEndTime) {
		log.Printf("appointment-times-logic -> method: hasTimeConflict: Conflict detected - Appointment time (%v - %v) is outside doctor's working hours (%v - %v)",
			startTimeAppointment, endTimeAppointment, doctorStartTime, doctorEndTime)
		return true
	}

	for _, a := range appointments {
		if a.Date != appointmentDate.Format("2006-01-02") {
			continue // Ignoramos citas de otros días
		}

		parsedStartTime, err1 := validate.ParseTime(a.StartTime)
		parsedEndTime, err2 := validate.ParseTime(a.EndTime)
		if err1 != nil || err2 != nil {
			log.Println("❌ Error al parsear los horarios de las citas existentes")
			continue
		}

		// ✅ Comparar solo la hora, ignorando la fecha
		if startTimeAppointment.Hour() < parsedEndTime.Hour() && endTimeAppointment.Hour() > parsedStartTime.Hour() {
			log.Printf("⛔ Conflict: Overlapping appointment on %s (%v - %v) with existing appointment (%v - %v)",
				appointmentDate.Format("2006-01-02"), startTimeAppointment, endTimeAppointment, parsedStartTime, parsedEndTime)
			return true
		}
	}

	return false
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

	if l.hasTimeConflict(appointment, startTime, endTime, appointmentDate) {
		log.Printf("appointment-times-logic -> method: ValidateAppointmentTime: Appointment time conflicts with existing appointments")
		return response.ErrorAppointmentTimeConflict
	}

	return nil
}
