package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/calculation"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type appointmentServiceID struct {
	repositoryAppointment repository.Repository[model.Appointment]
	repositoryService     repository.Repository[model.Service]
}

type AppointmentServiceID interface {
	IsServiceIDEXists(ID uint, hasInsurance bool) (*model.FinalServicePrice, error)
}

func NewAppointmentServiceID(repositoryAppointment repository.Repository[model.Appointment]) AppointmentServiceID {
	return &appointmentServiceID{repositoryAppointment: repositoryAppointment}
}

func (l *appointmentServiceID) IsServiceIDEXists(ID uint, hasInsurance bool) (*model.FinalServicePrice, error) {
	log.Println("appointment-service-id-logic -> method: IsServiceIDExists: received")

	service, err := l.repositoryService.GetByID(ID)
	if err != nil {
		log.Printf("appointment-service-id-logic -> method: IsServiceIDExists: Service with ID %d not exists: %v", ID, err)
		return nil, response.ErrorServiceNotFound
	}

	finalServicePrice := calculation.TotalServiceAmount(service.Price, hasInsurance)

	return finalServicePrice, nil
}
