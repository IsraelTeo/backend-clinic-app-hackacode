package logic

import (
	"log"

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
	if err := l.repository.Create(doctor); err != nil {
		log.Printf("doctor: Error saving doctor: %v", err)
		return response.ErrorToCreatedDoctor
	}
	return nil
}

func (l *doctorLogic) UpdateDoctor(ID uint, doctor *model.Doctor) error {
	doctorUpdate, err := l.GetDoctorByID(ID)
	if err != nil {
		log.Printf("doctor: Error fetching doctor with ID %d: %v to update", ID, err)
		return response.ErrorDoctorNotFound
	}

	doctorUpdate.Name = doctor.Name
	doctorUpdate.LastName = doctor.LastName
	doctorUpdate.Especialty = doctor.Especialty
	doctorUpdate.Day = doctor.Day
	doctorUpdate.StartTime = doctor.StartTime
	doctorUpdate.EndTime = doctor.EndTime
	doctorUpdate.Salary = doctor.Salary

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
