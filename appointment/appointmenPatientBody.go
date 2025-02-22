package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentPatientBody interface {
	HandlePatientBodyCreation(appointment *model.Appointment) error
	HandlePatientBodyUpdate(appointment *model.Appointment) error
}

type appointmentPatientBody struct {
	patientLogic logic.PatientLogic
}

// Constructor corregido para inyectar PatientLogic en lugar del repositorio directamente
func NewAppointmentPatientBody(patientLogic logic.PatientLogic) AppointmentPatientBody {
	return &appointmentPatientBody{patientLogic: patientLogic}
}

func (l *appointmentPatientBody) isPatientBodyExists(appointment *model.Appointment) bool {
	log.Printf("appointment-patient-body-logic -> method: isPatientBodyExists: received")
	return appointment.Patient != nil
}

func (l *appointmentPatientBody) HandlePatientBodyCreation(appointment *model.Appointment) error {
	if !l.isPatientBodyExists(appointment) {
		log.Printf("appointment-patient-body-logic -> method: handlerPatientBodyCreation: Body patient is empty")
		return response.ErrorBodyPatientEmpty
	}

	if err := l.patientLogic.CreatePatient(appointment.Patient); err != nil {
		log.Printf("appointment-patient-body-logic -> method: handlerPatientBodyCreation: Error creating patient: %v", err)
		return err
	}

	return nil
}

func (l *appointmentPatientBody) HandlePatientBodyUpdate(appointment *model.Appointment) error {
	if !l.isPatientBodyExists(appointment) {
		log.Printf("appointment-patient-body-logic -> method: handlerPatientBodyCreation: Body patient is empty")
		return response.ErrorBodyPatientEmpty
	}

	if err := l.patientLogic.UpdatePatient(appointment.PatientID, appointment.Patient); err != nil {
		log.Printf("appointment-patient-body-logic -> method: handlerPatientBodyCreation: Error creating patient: %v", err)
		return err
	}

	return nil
}
