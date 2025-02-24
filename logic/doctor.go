package logic

import (
	"fmt"
	"log"
	"strings"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

type DoctorLogic interface {
	GetDoctorByID(ID uint) (*model.Doctor, error)
	GetDoctorByDNI(DNI string) (*model.Doctor, error)
	GetAllDoctors(limit, offset int) ([]model.Doctor, error)
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
		log.Printf("doctor-logic: Error fetching doctor with ID %d: %v", ID, err)
		return nil, response.ErrorDoctorNotFoundID
	}

	return doctor, nil
}

func (l *doctorLogic) GetDoctorByDNI(DNI string) (*model.Doctor, error) {
	patient, err := l.repositoryDoctorMain.GetDoctorByDNI(DNI)
	if err != nil {
		log.Printf("doctor-logic: Error fetching doctor with DNI %s: %v", DNI, err)
		return nil, response.ErrorDoctorNotFoundDNI
	}

	return patient, nil
}

func (l *doctorLogic) GetAllDoctors(limit, offset int) ([]model.Doctor, error) {
	doctors, err := l.repositoryDoctor.GetAll(limit, offset)
	if err != nil {
		log.Printf("doctor-logic: Error fetching doctors: %v", err)
		return nil, response.ErrorDoctorsNotFound
	}

	return doctors, nil
}

func (l *doctorLogic) CreateDoctor(doctor *model.Doctor) error {
	birthDate, err := l.validateDoctor(doctor)
	if err != nil {
		return err
	}

	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err
	}

	doctor.Days = normalizedDays

	newDoctor := model.Doctor{
		Person: model.Person{
			Name:        doctor.Name,
			LastName:    doctor.LastName,
			DNI:         doctor.DNI,
			BirthDate:   birthDate,
			Email:       doctor.Email,
			PhoneNumber: doctor.PhoneNumber,
			Address:     doctor.Address,
		},
		Especialty: doctor.Especialty,
		Days:       normalizedDays,
		StartTime:  doctor.StartTime,
		EndTime:    doctor.EndTime,
		Salary:     doctor.Salary,
	}

	err = l.repositoryDoctor.Create(&newDoctor)
	if err != nil {
		log.Printf("doctor-logic: Error saving doctor: %v", err)
		return response.ErrorToCreatedDoctor
	}

	return nil
}

func (l *doctorLogic) UpdateDoctor(ID uint, doctor *model.Doctor) error {
	doctorUpdate, err := l.GetDoctorByID(ID)
	if err != nil {
		log.Printf("doctor-logic: Error fetching doctor with ID %d: %v to update", ID, err)
		return response.ErrorDoctorNotFoundID
	}

	birthDate, err := l.validateUpdatedDoctorFields(doctor, doctorUpdate)
	if err != nil {
		return err
	}

	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err
	}

	doctor.Days = normalizedDays

	doctorUpdate.Name = doctor.Name
	doctorUpdate.LastName = doctor.LastName
	doctorUpdate.Especialty = doctor.Especialty
	doctorUpdate.Salary = doctor.Salary
	doctorUpdate.StartTime = doctor.StartTime
	doctorUpdate.EndTime = doctor.EndTime
	doctorUpdate.BirthDate = birthDate
	doctorUpdate.PhoneNumber = doctor.PhoneNumber
	doctorUpdate.Email = doctor.Email
	doctorUpdate.Address = doctor.Address

	err = l.repositoryDoctor.Update(doctorUpdate)
	if err != nil {
		log.Printf("doctor-logic:: Error updating doctor with ID %d: %v", ID, err)
		return response.ErrorToUpdatedDoctor
	}

	return nil
}

func (l *doctorLogic) DeleteDoctor(ID uint) error {
	_, err := l.GetDoctorByID(ID)
	if err != nil {
		log.Printf("doctor-logic: Error fetching doctor with ID %d: %v to update", ID, err)
		return response.ErrorDoctorNotFoundID
	}

	err = l.repositoryDoctor.Delete(ID)
	if err != nil {
		log.Printf("doctor-logic: Error deleting doctor with ID %d: %v", ID, err)
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

	daysArray := strings.Split(days, ",")

	normalizedDays := []string{}

	for _, day := range daysArray {
		normalizedDay := strings.ToLower(strings.TrimSpace(day))
		validDay, exists := validDays[normalizedDay]
		if !exists {
			return "", fmt.Errorf("día inválido: %s, solo Lunes a Domingos son días permitidos", day)
		}

		normalizedDays = append(normalizedDays, validDay)
	}

	return strings.Join(normalizedDays, ","), nil
}

func (l *doctorLogic) validateDoctor(doctor *model.Doctor) (string, error) {
	err := validate.DNIDoctor(doctor)
	if err != nil {
		return "", err
	}

	err = validate.PhoneNumberDoctor(doctor)
	if err != nil {
		return "", err
	}

	err = validate.EmailDoctor(doctor)
	if err != nil {
		return "", err
	}

	birthDate, err := validate.BirthDateDoctor(doctor.BirthDate)
	if err != nil {
		return "", err
	}

	return birthDate, nil
}

func (l *doctorLogic) validateUpdatedDoctorFields(doctor *model.Doctor, doctorUpdate *model.Doctor) (string, error) {
	if doctor.DNI != doctorUpdate.DNI {
		err := validate.DNIDoctor(doctor)
		if err != nil {
			return "", err
		}
	}

	if doctor.PhoneNumber != doctorUpdate.PhoneNumber {
		err := validate.PhoneNumberDoctor(doctor)
		if err != nil {
			return "", err
		}
	}

	if doctor.Email != doctorUpdate.Email {
		err := validate.EmailDoctor(doctor)
		if err != nil {
			return "", err
		}
	}

	birthDate, err := validate.BirthDateDoctor(doctor.BirthDate)
	if err != nil {
		return "", err
	}

	return birthDate, nil
}
