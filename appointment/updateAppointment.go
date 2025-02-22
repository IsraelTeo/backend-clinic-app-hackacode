package appointment

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

/*type UpdateAppointmentDeps struct {
	RepositoryAppointment repository.Repository[model.Appointment]
	RepositoryDoctor      repository.Repository[model.Doctor]
	RepositoryService     repository.Repository[model.Service]
	RepositoryPackage     repository.Repository[model.Package]
	RepositoryPatient     repository.Repository[model.Patient]

	AppointmentPatientBody AppointmentPatientBody
	AppointmentDoctor      AppointmentDoctor
	AppointmentPackageID   AppointmentPackageID
	AppointmentPatientID   AppointmentPatientID
	AppointmentServiceID   AppointmentServiceID
	AppointmentTime        AppointmentTime
}*/

type AppointmentUpdate interface {
	UpdateAppointment(appointmentID uint, updatedAppointment *model.Appointment) (interface{}, error)
}

type appointmentUpdate struct {
	repositoryAppointment repository.Repository[model.Appointment]
	repositoryDoctor      repository.Repository[model.Doctor]
	repositoryService     repository.Repository[model.Service]
	repositoryPackage     repository.Repository[model.Package]
	repositoryPatient     repository.Repository[model.Patient]

	appointmentPatientBody AppointmentPatientBody
	appointmentDoctor      AppointmentDoctorID
	appointmentPackageID   AppointmentPackageID
	appointmentPatientID   AppointmentPatientID
	appointmentServiceID   AppointmentServiceID
	appointmentTime        AppointmentTime
}

func NewAppointmentUpdate(
	repositoryAppointment repository.Repository[model.Appointment],
	repositoryDoctor repository.Repository[model.Doctor],
	repositoryService repository.Repository[model.Service],
	repositoryPackage repository.Repository[model.Package],
	repositoryPatient repository.Repository[model.Patient],
	appointmentPatientBody AppointmentPatientBody,
	appointmentDoctor AppointmentDoctorID,
	appointmentPackageID AppointmentPackageID,
	appointmentPatientID AppointmentPatientID,
	appointmentServiceID AppointmentServiceID,
	appointmentTime AppointmentTime,
) AppointmentUpdate {
	return &appointmentUpdate{
		repositoryAppointment: repositoryAppointment,
		repositoryDoctor:      repositoryDoctor,
		repositoryService:     repositoryService,
		repositoryPackage:     repositoryPackage,
		repositoryPatient:     repositoryPatient,

		appointmentPatientBody: appointmentPatientBody,
		appointmentDoctor:      appointmentDoctor,
		appointmentPackageID:   appointmentPackageID,
		appointmentPatientID:   appointmentPatientID,
		appointmentServiceID:   appointmentServiceID,
		appointmentTime:        appointmentTime,
	}
}

/*func NewUpdateAppointment(deps UpdateAppointmentDeps) AppointmentUpdate {
	return &appointmentUpdate{
		repositoryAppointment: deps.RepositoryAppointment,
		repositoryDoctor:      deps.RepositoryDoctor,
		repositoryService:     deps.RepositoryService,
		repositoryPackage:     deps.RepositoryPackage,
		repositoryPatient:     deps.RepositoryPatient,

		appointmentPatientBody: deps.AppointmentPatientBody,
		appointmentDoctor:      deps.AppointmentDoctor,
		appointmentPackageID:   deps.AppointmentPackageID,
		appointmentPatientID:   deps.AppointmentPatientID,
		appointmentServiceID:   deps.AppointmentServiceID,
		appointmentTime:        deps.AppointmentTime,
	}
}*/

func (l *appointmentUpdate) UpdateAppointment(appointmentID uint, updatedAppointment *model.Appointment) (interface{}, error) {
	log.Printf("appointment-update-logic -> method: UpdateAppointment: received")

	// Buscar la cita existente
	existingAppointment, err := l.repositoryAppointment.GetByID(appointmentID)
	if err != nil {
		log.Printf("appointment-update-logic -> method: UpdateAppointment: Appointment not found with ID: %d", appointmentID)
		return nil, response.ErrorAppointmentNotFound
	}

	// Validar si el doctor existe
	if !l.appointmentDoctor.IsDoctorExists(updatedAppointment.DoctorID) {
		log.Printf("appointment-update-logic -> method: UpdateAppointment: Doctor not found with ID: %d", updatedAppointment.DoctorID)
		return nil, response.ErrorDoctorNotFoundID
	}

	// Obtener o actualizar el paciente
	patient, err := l.getPatientToUpdate(updatedAppointment)
	if err != nil {
		log.Printf("appointment-update-logic -> method: UpdateAppointment: Error getting patient: %v", err)
		return nil, err
	}

	// Obtener los detalles de precio actualizados
	priceDetails, err := l.getPriceDetails(updatedAppointment, patient)
	if err != nil {
		log.Printf("appointment-update-logic -> method: UpdateAppointment: Error getting price details: %v", err)
		return nil, err
	}

	// Construir la cita actualizada
	updatedAppointmentData := l.buildUpdatedAppointment(existingAppointment, updatedAppointment, patient)

	// Guardar la actualización en la BD
	if err := l.repositoryAppointment.Update(updatedAppointmentData); err != nil {
		log.Printf("appointment-update-logic -> method: UpdateAppointment: Error updating appointment: %v", err)
		return nil, err
	}

	return priceDetails, nil
}

func (l *appointmentUpdate) getPatientToUpdate(appointment *model.Appointment) (*model.Patient, error) {
	log.Printf("appointment-create-logic -> method: getPatientt: received")

	if appointment.Patient.ID != 0 {
		patient, err := l.appointmentPatientID.IsPatientIDExists(appointment.Patient.ID)
		if err != nil {
			log.Printf("appointment-create-logic -> method: getPatient: Error getting patient by ID: %v", err)
			return nil, err
		}

		return patient, nil
	}

	err := l.appointmentPatientBody.HandlePatientBodyUpdate(appointment)
	if err != nil {
		log.Printf("appointment-create-logic -> method: getPatient: Error creating patient: %v", err)
		return nil, err
	}

	return appointment.Patient, nil
}

// Método para obtener el precio de un servicio o paquete
func (l *appointmentUpdate) getPriceDetails(appointment *model.Appointment, patient *model.Patient) (interface{}, error) {
	log.Printf("appointment-create-logic -> getPriceDetails: received")

	if appointment.ServiceID != 0 {
		finalServicePrice, err := l.appointmentServiceID.IsServiceIDEXists(appointment.ServiceID, patient.Insurance)
		if err != nil {
			log.Printf("appointment-create-logic -> getPriceDetails: Error getting service price: %v", err)
			return nil, err
		}

		return &model.FinalServicePrice{
			TotalAmount:       finalServicePrice.TotalAmount,
			InsuranceDiscount: finalServicePrice.InsuranceDiscount,
			FinalPrice:        finalServicePrice.FinalPrice,
		}, nil
	}

	finalPkgPrice, err := l.appointmentPackageID.IsPackageIDExists(appointment.PackageID, patient.Insurance)
	if err != nil {
		log.Printf("appointment-create-logic -> getPriceDetails: Error getting package price: %v", err)
		return nil, err
	}

	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: finalPkgPrice.InsuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     finalPkgPrice.TotalAmount,
			DiscountPackage: finalPkgPrice.DiscountPackage,
			FinalPrice:      finalPkgPrice.FinalPrice,
		},
	}, nil
}

// Método para construir la cita actualizada
func (l *appointmentUpdate) buildUpdatedAppointment(existingAppointment, updatedAppointment *model.Appointment, patient *model.Patient) *model.Appointment {
	log.Printf("appointment-update-logic -> method: buildUpdatedAppointment: received")

	return &model.Appointment{
		ID:          existingAppointment.ID, // Mantiene el mismo ID
		Patient:     patient,
		DoctorID:    updatedAppointment.DoctorID,
		ServiceID:   updatedAppointment.ServiceID,
		PackageID:   updatedAppointment.PackageID,
		Date:        updatedAppointment.Date,
		StartTime:   updatedAppointment.StartTime,
		EndTime:     updatedAppointment.EndTime,
		Paid:        existingAppointment.Paid, // Mantiene el estado de pago
		TotalAmount: updatedAppointment.TotalAmount,
	}
}

/*func (l *appointmentLogic) UpdateAppointmentWithService(ID uint, appointment *model.Appointment) (*model.FinalServicePrice, error) {
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return nil, err
	}

	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	if appointment.Patient.ID != 0 {
		if err := l.handlePatientIDCreation(appointment); err != nil {
			return nil, err
		}
	}

	if appointment.Patient != (model.Patient{}) {
		if err := l.handlePatientBodyCreation(appointment); err != nil {
			return nil, err
		}
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalServicePrice, err := l.isServiceIDEXists(appointment.ServiceID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentUpdate.Patient = appointment.Patient
	appointmentUpdate.DoctorID = appointment.DoctorID
	appointmentUpdate.ServiceID = appointment.ServiceID
	appointmentUpdate.Date = appointment.Date
	appointmentUpdate.StartTime = appointment.StartTime
	appointmentUpdate.EndTime = appointment.EndTime
	appointmentUpdate.TotalAmount = appointment.TotalAmount

	if err := l.repo.RepositoryAppointment.Update(appointmentUpdate); err != nil {
		log.Printf("appointment-logic: Error actualizando la cita con ID %d: %v", ID, err)
		return nil, err
	}

	return &model.FinalServicePrice{
		TotalAmount:       finalServicePrice.TotalAmount,
		InsuranceDiscount: finalServicePrice.InsuranceDiscount,
		FinalPrice:        finalServicePrice.FinalPrice,
	}, nil
}


func (l *appointmentLogic) UpdateAppointmentWithPackage(ID uint, appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error) {
	appointmentUpdate, err := l.GetAppointmentByID(ID)
	if err != nil {
		return nil, err
	}

	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	// Si el paciente tiene ID, solo se asocia. Si no, se registra uno nuevo.
	if appointment.Patient.ID == 0 {
		if err := l.handlePatientBodyCreation(appointment); err != nil {
			return nil, err
		}
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalPkgPrice, err := l.isPackageIDExists(appointment.PackageID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentUpdate.Patient = appointment.Patient
	appointmentUpdate.DoctorID = appointment.DoctorID
	appointmentUpdate.PackageID = appointment.PackageID
	appointmentUpdate.Date = appointment.Date
	appointmentUpdate.StartTime = appointment.StartTime
	appointmentUpdate.EndTime = appointment.EndTime
	appointmentUpdate.TotalAmount = appointment.TotalAmount

	if err := l.repo.RepositoryAppointment.Update(appointmentUpdate); err != nil {
		log.Printf("appointment-logic: Error update appointment con ID %d: %v", ID, err)
		return nil, err
	}

	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: finalPkgPrice.InsuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     finalPkgPrice.TotalAmount,
			DiscountPackage: finalPkgPrice.DiscountPackage,
			FinalPrice:      finalPkgPrice.FinalPrice,
		},
	}, nil
}*/
