package appointment

import (
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

func NewAppointmentServiceID(repositoryAppointment repository.Repository[model.Appointment], repositoryService repository.Repository[model.Service]) AppointmentServiceID {
	return &appointmentServiceID{repositoryAppointment: repositoryAppointment, repositoryService: repositoryService}
}

func (l *appointmentServiceID) IsServiceIDEXists(ID uint, hasInsurance bool) (*model.FinalServicePrice, error) {
	service, err := l.repositoryService.GetByID(ID)
	if err != nil {
		return nil, response.ErrorServiceNotFound
	}

	finalServicePrice := calculation.TotalServiceAmount(service.Price, hasInsurance)

	return finalServicePrice, nil
}
