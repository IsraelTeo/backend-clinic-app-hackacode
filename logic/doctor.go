package logic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

type DoctorLogic interface {
	GetDoctorByID(ID uint) (*model.Doctor, error)
	GetDoctorByDNI(DNI string) (*model.Doctor, error)
	GetAllDoctors() ([]model.Doctor, error)
	CreateDoctor(doctor *model.Doctor) error
	UpdateDoctor(ID uint, doctor *model.Doctor) error
	DeleteDoctor(ID uint) error
}

type doctorLogic struct {
	repositoryDoctor     repository.Repository[model.Doctor]
	repositoryDoctorMain repository.DoctorRepository
}

func NewDoctorLogic(repositoryDoctor repository.Repository[model.Doctor], repositoryDoctorMain repository.DoctorRepository) DoctorLogic {
	return &doctorLogic{repositoryDoctor: repositoryDoctor, repositoryDoctorMain: repositoryDoctorMain}
}

func (l *doctorLogic) GetDoctorByID(ID uint) (*model.Doctor, error) {
	doctor, err := l.repositoryDoctor.GetByID(ID)
	if err != nil {
		log.Printf("doctor: Error fetching doctor with ID %d: %v", ID, err)
		return nil, response.ErrorDoctorNotFound
	}
	return doctor, nil
}

func (l *doctorLogic) GetDoctorByDNI(DNI string) (*model.Doctor, error) {
	patient, err := l.repositoryDoctorMain.GetDoctorByDNI(DNI)
	if err != nil {
		log.Printf("patient: Error fetching patient with DNI %s: %v", DNI, err)
		return nil, response.ErrorPatientNotFound
	}

	return patient, nil
}

func (l *doctorLogic) GetAllDoctors() ([]model.Doctor, error) {
	doctors, err := l.repositoryDoctor.GetAll()
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
		return err //falta averiguar de que es este error
	}

	doctor.Days = normalizedDays

	// Convertir las horas y la fecha de nacimiento

	startTime, err := validate.ParseTime(doctor.StartTime)
	if err != nil {
		return err
	}

	endTime, err := validate.ParseTime(doctor.EndTime)
	if err != nil {
		return err
	}

	birthDate, err := validate.ParseDate(doctor.BirthDate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", doctor.BirthDate)
		return response.ErrorDoctorInvalidDateFormat
	}

	if validate.CheckDNIExists[model.Doctor](doctor.DNI, doctor) {
		log.Printf("Error checking if doctorexists by DNI: %s", doctor.DNI)
		return response.ErrorDoctorExistsDNI
	}

	if validate.CheckPhoneNumberExists[model.Doctor](doctor.PhoneNumber, doctor) {
		log.Printf("Error checking if patient exists by phone number: %s", doctor.PhoneNumber)
		return response.ErrorDoctorExistsPhoneNumber
	}

	if validate.CheckEmailExists[model.Doctor](doctor.Email, doctor) {
		log.Printf("Error checking if patient exists by email: %s", doctor.Email)
		return response.ErrorDoctorExistsEmail
	}

	if !validate.IsDateInPast(birthDate) {
		log.Printf("Error birthdate is past: %v", doctor.BirthDate)
		return response.ErrorDoctorBrithDateIsFuture
	}

	// Combinando la fecha y la hora de inicio
	startTimeMain := combineDateAndTime(startTime, time.Now())

	// Combinando la fecha y la hora de fin
	endTimeMain := combineDateAndTime(endTime, time.Now())

	//validar que la hora inicial no sea en futuro de la hora final
	if validate.IsStartBeforeEnd(startTime, endTimeMain) {
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

	if err := l.repositoryDoctor.Create(&newDoctor); err != nil {
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
	startTime, err := validate.ParseTime(doctor.StartTime)
	if err != nil {
		return err
	}

	endTime, err := validate.ParseTime(doctor.EndTime)
	if err != nil {
		return err
	}

	birthDate, err := validate.ParseDate(doctor.BirthDate)
	if err != nil {
		return err
	}

	// Combinando la fecha y la hora de inicio
	startTimeMain := combineDateAndTime(startTime, time.Now())

	// Combinando la fecha y la hora de fin
	endTimeMain := combineDateAndTime(endTime, time.Now())

	//validar que la hora inicial no sea en futuro de la hora final
	if validate.IsStartBeforeEnd(startTime, endTimeMain) {
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
	if err = l.repositoryDoctor.Update(doctorUpdate); err != nil {
		log.Printf("doctor: Error updating doctor with ID %d: %v", ID, err)
		return response.ErrorToUpdatedDoctor
	}
	return nil
}

func (l *doctorLogic) DeleteDoctor(ID uint) error {
	if err := l.repositoryDoctor.Delete(ID); err != nil {
		log.Printf("doctor: Error deleting doctor with ID %d: %v", ID, err)
		return response.ErrorToDeletedDoctor
	}

	return nil
}

func normalizeDays(days string) (string, error) {
	//Días válidos
	validDays := map[string]string{
		"lunes":     string(model.Moonday),
		"martes":    string(model.Tuesday),
		"miercoles": string(model.Wednesday),
		"miércoles": string(model.Wednesday),
		"jueves":    string(model.Thursday),
		"viernes":   string(model.Friday),
	}

	// Dividir los días en una lista, es una lista de string
	daysArray := strings.Split(days, ",")

	///declarar una lista de dias normalizados
	normalizedDays := []string{}

	//Recorremos los días de lunes a viernes
	for _, day := range daysArray {

		// Normalizar el día: convertirlo a minúsculas y eliminar espacios adicionales
		normalizedDay := strings.ToLower(strings.TrimSpace(day))

		// Verificar si el día normalizado está en el mapa de días válidos
		validDay, exists := validDays[normalizedDay]
		if !exists {
			return "", fmt.Errorf("invalid day: %s, only Monday to Friday are allowed", day)
		}

		// Agregar el día válido a la lista de días normalizados
		normalizedDays = append(normalizedDays, validDay)
	}

	// Serializar los días a una cadena separada por comas
	return strings.Join(normalizedDays, ","), nil
}

func combineDateAndTime(timeObj time.Time, referenceDate time.Time) time.Time {
	return time.Date(
		referenceDate.Year(),
		referenceDate.Month(),
		referenceDate.Day(),
		timeObj.Hour(),
		timeObj.Minute(),
		0,
		0,
		referenceDate.Location(),
	)
}
