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

func (l *doctorLogic) GetAllDoctors() ([]model.Doctor, error) {
	doctors, err := l.repositoryDoctor.GetAll()
	if err != nil {
		log.Printf("doctor-logic: Error fetching doctors: %v", err)
		return nil, response.ErrorDoctorsNotFound
	}

	if len(doctors) == 0 {
		log.Println("doctor-logic: No doctors found")
		return []model.Doctor{}, response.ErrorListDoctorsEmpty
	}

	return doctors, nil
}

func (l *doctorLogic) CreateDoctor(doctor *model.Doctor) error {
	//Valimos el DNI, DNI, número telefónico no sean duplicads y que la fecha de nacimiento que no sea en pasadodel doctor
	birthDate, err := l.validateDoctor(doctor)
	if err != nil {
		return err
	}

	//Normalizamos días, es decir: Obtenemos los días que trabaja el doctor
	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err
	}

	doctor.Days = normalizedDays

	//Parseamos el horario inicial y final del torno, es decir: Su turno
	startTime, endTime, err := parseShiftDoctor(doctor)
	if err != nil {
		return err
	}

	//Obtener la fecha y horario del horario inicial y horario final
	startTimeMain := combineDateAndTime(startTime, time.Now())
	endTimeMain := combineDateAndTime(endTime, time.Now())

	//verificamos que el tiempo final fuera en futuro del tiempo inicial
	if !validate.IsStartBeforeEnd(*startTime, endTimeMain) {
		return response.ErrorInvalidEndTimeInPastDoctor
	}

	//construye el doctor
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

	//Valimos el DNI, DNI, número telefónico no sean duplicads y que la fecha de nacimiento que no sea en pasadodel doctor
	birthDate, err := l.validateUpdatedDoctorFields(doctor, doctorUpdate)
	if err != nil {
		return err
	}

	//Normalizamos días, es decir: Obtenemos los días que trabaja el doctor
	normalizedDays, err := normalizeDays(doctor.Days)
	if err != nil {
		return err
	}

	doctor.Days = normalizedDays

	//Parseamos el horario inicial y final del torno, es decir: Su turno
	startTime, endTime, err := parseShiftDoctor(doctor)
	if err != nil {
		return err
	}

	//Obtener la fecha y horario del horario inicial y horario final
	startTimeMain := combineDateAndTime(startTime, time.Now())
	endTimeMain := combineDateAndTime(endTime, time.Now())

	//verificamos que el tiempo final fuera en futuro del tiempo inicial
	if !validate.IsStartBeforeEnd(startTimeMain, endTimeMain) {
		return response.ErrorInvalidEndTimeInPastDoctor
	}

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

	if err = l.repositoryDoctor.Update(doctorUpdate); err != nil {
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

	if err := l.repositoryDoctor.Delete(ID); err != nil {
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

func parseShiftDoctor(doctor *model.Doctor) (*time.Time, *time.Time, error) {
	startTime, err := validate.ParseTime(doctor.StartTime)
	if err != nil {
		return nil, nil, response.ErrorInvalidStartTimeDoctor
	}

	endTime, err := validate.ParseTime(doctor.EndTime)
	if err != nil {
		return nil, nil, response.ErrorInvalidEndTimeDoctor
	}

	return &startTime, &endTime, nil
}

func combineDateAndTime(timeObj *time.Time, referenceDate time.Time) time.Time {
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

func (l *doctorLogic) validateDoctor(doctor *model.Doctor) (*time.Time, error) {
	if err := validate.DNIDoctor(doctor); err != nil {
		return nil, err
	}

	if err := validate.PhoneNumberDoctor(doctor); err != nil {
		return nil, err
	}

	if err := validate.EmailDoctor(doctor); err != nil {
		return nil, err
	}

	birthDate, err := validate.BirthDateDoctor(doctor.BirthDate)
	if err != nil {
		return nil, err
	}

	return birthDate, nil
}

func (l *doctorLogic) validateUpdatedDoctorFields(doctor *model.Doctor, doctorUpdate *model.Doctor) (*time.Time, error) {
	if doctor.DNI != doctorUpdate.DNI {
		if err := validate.DNIDoctor(doctor); err != nil {
			return nil, err
		}
	}

	if doctor.PhoneNumber != doctorUpdate.PhoneNumber {
		if err := validate.PhoneNumberDoctor(doctor); err != nil {
			return nil, err
		}
	}

	if doctor.Email != doctorUpdate.Email {
		if err := validate.EmailDoctor(doctor); err != nil {
			return nil, err
		}
	}

	birthDate, err := validate.BirthDateDoctor(doctor.BirthDate)
	if err != nil {
		return nil, err
	}

	return birthDate, nil
}
