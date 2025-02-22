package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
)

type AppointmentDoctorID interface {
	IsDoctorExists(doctorID uint) bool
}
type appointmentDoctorID struct {
	repositoryDoctor repository.Repository[model.Doctor]
}

func NewAppointmentDoctorID(repositoryDoctor repository.Repository[model.Doctor]) AppointmentDoctorID {
	return &appointmentDoctorID{repositoryDoctor: repositoryDoctor}
}

func (l *appointmentDoctorID) IsDoctorExists(doctorID uint) bool {
	log.Println("appointment-doctor-logic -> method: IsDoctorExists: received")

	_, err := l.repositoryDoctor.GetByID(doctorID)
	if err != nil {
		log.Printf("appointment-doctor-logic -> method: IsDoctorExists: Error fetching doctor with ID %d: %v", doctorID, err)
		return false
	}

	return true
}
