package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type appointmentPatientID struct {
	repositoryPatient repository.Repository[model.Patient]
}

type AppointmentPatientID interface {
	IsPatientIDExists(ID uint) (*model.Patient, error)
}

func NewAppointmentPatientID(repositoryPatient repository.Repository[model.Patient]) AppointmentPatientID {
	return &appointmentPatientID{repositoryPatient: repositoryPatient}
}

func (l *appointmentPatientID) IsPatientIDExists(ID uint) (*model.Patient, error) {
	log.Println("appointment-patient-id-logic -> method: isPatientIDexists: received")

	patient, err := l.repositoryPatient.GetByID(ID)
	if err != nil {
		log.Printf("appointment-patient-id-logic -> method: isPatientIDexists: Error fetching patient by ID: %v", err)
		return nil, response.ErrorPatientNotFoundID
	}

	return patient, nil
}

/*func (l *appointmentPatientID) HandlePatientIDCreation(appointment *model.Appointment) error {
	log.Println("appointment-patient-id-logic -> method: HandlePatientIDCreation: received")

	patient, err := l.isPatientIDExists(appointment.PatientID)
	if err != nil {
		log.Printf("appointment-patient-id-logic -> method: HandlePatientIDCreation: Patient not found with ID: %d", appointment.PatientID)
		return err
	}

	appointment.Patient = *patient

	return nil
}*/
