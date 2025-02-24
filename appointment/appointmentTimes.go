package appointment

import (
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
	startTime, err := validate.ParseTime(startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := validate.ParseTime(endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}

func (l *appointmentTime) parseTimesAndDate(appointment *model.Appointment) (time.Time, time.Time, time.Time, error) {
	startTime, endTime, err := l.parseStartAndEndTime(appointment.StartTime, appointment.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	appointmentDate, err := validate.ParseDate(appointment.Date)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	return startTime, endTime, appointmentDate, nil
}

func (l *appointmentTime) hasTimeConflict(appointment *model.Appointment, startTimeAppointment, endTimeAppointment, appointmentDate time.Time) error {
	doctor, err := l.repositoryDoctor.GetByID(appointment.DoctorID)
	if err != nil {
		return response.ErrorDoctorNotFoundID
	}

	workingDays := strings.Split(doctor.Days, ",")

	appointmentWeekDaySpanish := validate.DayToGolang[appointmentDate.Weekday()]

	//verifica que el día de la cita sea un día que el médico tenga turno
	if !validate.IsDayAvailable(appointmentWeekDaySpanish, workingDays) {
		return response.ErrorAppointmentDayNotAvailable
	}

	doctorStartTime, doctorEndTime, err := l.parseStartAndEndTime(doctor.StartTime, doctor.EndTime)
	if err != nil {
		return err
	}

	//verifica que el horario de la cita esté dentro del turno del doctor
	if !validate.IsWithinTimeRange(startTimeAppointment, endTimeAppointment, doctorStartTime, doctorEndTime, true, true) {
		return response.ErrorInvalidAppointmentTime
	}

	doctorAppointments, err := l.repositoryAppointmentMain.GetAppointmentsByDoctorAndDate(appointment.DoctorID, appointmentDate)
	if err != nil {
		return response.ErrorFetchingAppointments
	}

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

		// Verifica si la cita tiene cruce con otra cita existente
		if appointmentStartTime.Before(parsedDoctorEndTime) && appointmentEndTime.After(parsedDoctorStartTime) {
			return response.ErrorAppointmentTimeConflict
		}

	}

	return nil
}

func (l *appointmentTime) ValidateAppointmentTime(appointment *model.Appointment) error {
	startTime, endTime, appointmentDate, err := l.parseTimesAndDate(appointment)
	if err != nil {
		return err
	}

	if validate.IsDateInPast(appointmentDate) {
		return response.ErrorAppointmentDateInPast
	}

	if !validate.IsStartBeforeEnd(startTime, endTime) {
		return response.ErrorInvalidAppointmentTimeRange
	}

	err = l.hasTimeConflict(appointment, startTime, endTime, appointmentDate)
	if err != nil {
		return err
	}

	return nil
}
