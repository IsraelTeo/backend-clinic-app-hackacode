package logic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/dto"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/mapper"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
)

// DoctorService defines the interface for doctor service operations
type DoctorService interface {
	GetDoctorByID(ID uint) (*dto.DoctorResponse, error)
	GetDoctors() ([]dto.DoctorResponse, error)
	CreateDoctor(request *dto.DoctorRequest) (*dto.DoctorResponse, error)
	UpdateDoctor(ID uint, doctorRequest *dto.DoctorRequest) (*dto.DoctorResponse, error)
	DeleteDoctor(ID uint) error
	GetDoctorByDNI(DNI string) (*dto.DoctorResponse, error)
}

// doctorLogic implements the DoctorService interface
type doctorService struct {
	repositoryDoctor repository.DoctorRepository
	mapper           mapper.DoctorMapper
}

// Dependency injection for the doctor service
func NewDoctorService(repositoryDoctor repository.DoctorRepository, mapper mapper.DoctorMapper) DoctorService {
	return &doctorService{
		repositoryDoctor: repositoryDoctor,
		mapper:           mapper,
	}
}

// GetDoctorByID retrieves a doctor by its ID
func (s *doctorService) GetDoctorByID(ID uint) (*dto.DoctorResponse, error) {
	doctor, err := s.repositoryDoctor.GetByID(ID)
	if err != nil {
		log.Printf("doctor-service: Error fetching doctor with ID %d: %v", ID, err)
		return nil, response.ErrorDoctorNotFoundID
	}

	doctorDTO, err := s.mapper.ModelToResponse(doctor)
	if err != nil {
		log.Printf("doctor-service: Error mapping doctor to DTO: %v", err)
		return nil, response.ErrorToDoctorResponseConversion
	}

	return doctorDTO, nil
}

// GetDoctors retrieves a list of doctors
func (s *doctorService) GetDoctors() ([]dto.DoctorResponse, error) {
	doctors, err := s.repositoryDoctor.Get()
	if err != nil {
		log.Printf("doctor-service: Error fetching doctors: %v", err)
		return nil, response.ErrorDoctorsNotFound
	}

	doctorsDTOS, err := s.mapper.ModelListToResponseList(doctors)
	if err != nil {
		log.Printf("doctor-service: Error mapping doctors to DTOs: %v", err)
		return nil, response.ErrorToDoctorResponseConversion
	}

	return doctorsDTOS, nil
}

// CreateDoctor creates a new doctor
func (s *doctorService) CreateDoctor( doctorRequest *dto.DoctorRequest) (*dto.DoctorResponse, error) {
	doctor, err := s.mapper.RequestToModel( doctorRequest)
	if err != nil {
		log.Printf("doctor-service: Error mapping DTO to doctor model: %v", err)
		return nil, response.ErrorToDoctorResponseConversion
	}

	/*
		birthDate, err := s.validateDoctor(doctor)
		if err != nil {
			log.Printf("doctor-service: Validation error: %v", err)
			return nil, err
		}

		normalizedDays, err := normalizeDays(doctor.Days)
		if err != nil {
			return err
		}

		doctor.Days = normalizedDays*/

	newDoctor := model.Doctor{
		Person: model.Person{
			Name:        doctor.Name,
			LastName:    doctor.LastName,
			DNI:         doctor.DNI,
			BirthDate:   doctorDTO.BirthDate,
			Email:       doctor.Email,
			PhoneNumber: doctor.PhoneNumber,
			Address:     doctor.Address,
		},
		Especialty: doctor.Especialty,
		//Days:       normalizedDays,
		//StartTime: doctor.StartTime,
		//EndTime:   doctor.EndTime,
		Salary:    doctor.Salary,
		CreatedAt: time.Now(),
	}

	err = s.repositoryDoctor.Create(&newDoctor)
	if err != nil {
		log.Printf("doctor-service: Error saving doctor: %v", err)
		return nil, response.ErrorToCreatedDoctor
	}

	doctorResponse, err := s.mapper.ModelToResponse(&newDoctor)
	if err != nil {
		log.Printf("doctor-service: Error mapping created doctor to DTO: %v", err)
		return nil, response.ErrorToDoctorResponseConversion
	}

	return doctorResponse, nil
}

// UpdateDoctor updates an existing doctor
func (s *doctorService) UpdateDoctor(ID uint, request *dto.DoctorRequest) (*dto.DoctorResponse, error) {
	doctorFind, err := s.repositoryDoctor.GetByID(ID)
	if err != nil {
		log.Printf("doctor-service: Error fetching doctor with ID %d: %v", ID, err)
		return nil, response.ErrorDoctorNotFoundID
	}
	/*
		birthDate, err := s.validateDoctor(doctor)
		if err != nil {
			log.Printf("doctor-service: Validation error: %v", err)
			return nil, err
		}

		normalizedDays, err := normalizeDays(doctor.Days)
		if err != nil {
			return err
		}

		doctor.Days = normalizedDays*/

		doctorFind.Name = doctorDTO.Name
		doctorFind.LastName = doctorDTO.LastName
		doctorFind.DNI = doctorDTO.DNI
		doctorFind.BirthDate = doctorDTO.BirthDate
		doctorFind.Email = doctorDTO.Email
		doctorFind.PhoneNumber = doctorDTO.PhoneNumber
		doctorFind.Address = doctorDTO.PhoneNumber
		doctorFind.Especialty = doctorDTO.Especialty
		//doctorFind.StartTime = doctorDTO.StartTime
		//doctorFind.EndTime = doctorDTO.EndTime
		doctorFind.Salary = doctorDTO.Salary


	}

	err = s.repositoryDoctor.Create(&newDoctor)
	if err != nil {
		log.Printf("doctor-service: Error saving doctor: %v", err)
		return nil, response.ErrorToCreatedDoctor
	}

	doctorResponse, err := s.mapper.ModelToResponse(&newDoctor)
	if err != nil {
		log.Printf("doctor-service: Error mapping created doctor to DTO: %v", err)
		return nil, response.ErrorToDoctorResponseConversion
	}

	return doctorResponse, nil

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
