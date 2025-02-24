package appointment

import (
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
	_, err := l.repositoryDoctor.GetByID(doctorID)
	if err != nil {
		return false
	}

	return true
}
