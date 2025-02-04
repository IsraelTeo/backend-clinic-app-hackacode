package logic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type DoctorLogic interface {
	GetDoctorByID(ID uint) (*model.Doctor, error)
	GetAllDoctors() ([]model.Doctor, error)
	CreateDoctor(doctor *model.Doctor) error
	UpdateDoctor(ID uint, doctor *model.Doctor) error
	DeleteDoctor(ID uint) error
}

type doctorLogic struct {
	repository repository.Repository[model.Doctor]
}

func NewDoctorLogic(repository repository.Repository[model.Doctor]) DoctorLogic {
	return &doctorLogic{repository: repository}
}

func (l *doctorLogic) GetDoctorByID(ID uint) (*model.Doctor, error) {
	doctor, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("doctor: Error fetching doctor with ID %d: %v", ID, err)
		return nil, response.ErrorDoctorNotFound
	}
	return doctor, nil
}

func (l *doctorLogic) GetAllDoctors() ([]model.Doctor, error) {
	doctors, err := l.repository.GetAll()
	if err != nil {
		log.Printf("doctor: Error fetching doctors: %v", err)
		return nil, response.ErrorDoctorNotFound
	}

	if len(doctors) == 0 {
		log.Println("doctor: No doctors found")
		return []model.Doctor{}, response.ErrorListDoctorsEmpty
	}
	return doctors, nil
}

func (l *doctorLogic) CreateDoctor(doctor *model.Doctor) error {
	// Normalizar los días ingresados
	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err // Retornar el error si los días no son válidos
	}
	doctor.Days = normalizedDays

	// Convertir las horas y la fecha de nacimiento
	startTime, endTime, birthDate, err := parseTimeAndDate(doctor.StartTime, doctor.EndTime, doctor.BirthDate)
	if err != nil {
		return err // Retornar el error si el formato de hora o fecha no es válido
	}

	// Combinando la fecha y la hora de inicio
	now := time.Now()
	startTimeMain := combineDateAndTime(startTime, now)

	// Combinando la fecha y la hora de fin
	endTimeMain := combineDateAndTime(endTime, now)

	if endTimeMain.Before(startTimeMain) {
		return fmt.Errorf("end time must be after start time")
	}

	// Construir un nuevo struct Doctor con los valores procesados
	newDoctor := model.Doctor{
		Person: model.Person{
			Name:        doctor.Name,
			LastName:    doctor.LastName,
			DNI:         doctor.DNI,
			BirthDate:   birthDate.String(),
			Email:       doctor.Email,
			PhoneNumber: doctor.PhoneNumber,
			Address:     doctor.Address,
		},
		Especialty: doctor.Especialty,
		Days:       normalizedDays,
		StartTime:  startTimeMain.Format("15:04"),
		EndTime:    endTimeMain.Format("15:04"),
		Salary:     doctor.Salary,
	}

	if err := l.repository.Create(&newDoctor); err != nil {
		log.Printf("doctor: Error saving doctor: %v", err)
		return response.ErrorToCreatedDoctor
	}
	return nil
}

func (l *doctorLogic) UpdateDoctor(ID uint, doctor *model.Doctor) error {
	// Obtener el doctor por ID
	doctorUpdate, err := l.GetDoctorByID(ID)
	if err != nil {
		log.Printf("doctor: Error fetching doctor with ID %d: %v to update", ID, err)
		return response.ErrorDoctorNotFound
	}

	// Normalizar los días ingresados
	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err // Retornar el error si los días no son válidos
	}
	doctor.Days = normalizedDays

	// Convertir las horas y la fecha de nacimiento
	startTime, endTime, birthDate, err := parseTimeAndDate(doctor.StartTime, doctor.EndTime, doctor.BirthDate)
	if err != nil {
		return err // Retornar el error si el formato de hora o fecha no es válido
	}

	// Combinando la fecha y la hora de inicio
	now := time.Now()
	startTimeMain := combineDateAndTime(startTime, now)

	// Combinando la fecha y la hora de fin
	endTimeMain := combineDateAndTime(endTime, now)

	// Validar que la hora de fin no sea antes de la hora de inicio
	if endTimeMain.Before(startTimeMain) {
		return fmt.Errorf("end time must be after start time")
	}

	// Actualizar los campos de doctorUpdate con los valores nuevos
	doctorUpdate.Name = doctor.Name
	doctorUpdate.LastName = doctor.LastName
	doctorUpdate.Especialty = doctor.Especialty
	doctorUpdate.Salary = doctor.Salary
	doctorUpdate.StartTime = startTimeMain.Format("15:04")
	doctorUpdate.EndTime = endTimeMain.Format("15:04")
	doctorUpdate.BirthDate = birthDate.String()
	doctorUpdate.PhoneNumber = doctor.PhoneNumber
	doctorUpdate.Email = doctor.Email
	doctorUpdate.Address = doctor.Address

	// Realizar la actualización en el repositorio
	if err = l.repository.Update(doctorUpdate); err != nil {
		log.Printf("doctor: Error updating doctor with ID %d: %v", ID, err)
		return response.ErrorToUpdatedDoctor
	}
	return nil
}

func (l *doctorLogic) DeleteDoctor(ID uint) error {
	if err := l.repository.Delete(ID); err != nil {
		log.Printf("doctor: Error deleting doctor with ID %d: %v", ID, err)
		return response.ErrorToDeletedDoctor
	}

	return nil
}

func normalizeDays(days string) (string, error) {
	validDays := map[string]string{
		"lunes":     string(model.Moonday),
		"martes":    string(model.Tuesday),
		"miercoles": string(model.Wednesday),
		"miércoles": string(model.Wednesday),
		"jueves":    string(model.Thursday),
		"viernes":   string(model.Friday),
	}

	// Dividir los días en una lista
	daysArray := strings.Split(days, ",")
	normalizedDays := []string{}

	for _, day := range daysArray {
		normalizedDay, exists := validDays[strings.ToLower(strings.TrimSpace(day))]
		if !exists {
			return "", fmt.Errorf("invalid day: %s, only Monday to Friday are allowed", day)
		}
		normalizedDays = append(normalizedDays, normalizedDay)
	}

	// Serializar los días a una cadena separada por comas
	return strings.Join(normalizedDays, ","), nil
}

// Función para convertir las horas y la fecha
func parseTimeAndDate(startTimeStr, endTimeStr, birthDateStr string) (time.Time, time.Time, time.Time, error) {
	// Convertir las horas a time.Time (puedes usar una fecha arbitraria, como 1970-01-01)
	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid end time format: %v", err)
	}

	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid birth date format: %v", err)
	}

	return startTime, endTime, birthDate, nil
}

// Función para combinar la fecha y la hora
func combineDateAndTime(timeObj time.Time, referenceDate time.Time) time.Time {
	return time.Date(
		referenceDate.Year(), referenceDate.Month(), referenceDate.Day(),
		timeObj.Hour(), timeObj.Minute(), 0, 0, referenceDate.Location(),
	)
}
