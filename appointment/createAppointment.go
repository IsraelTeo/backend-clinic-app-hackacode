package appointment

import (
	"errors"
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

// struct que contiene las dependencias necesarias para la creación de una cita
/*type CreateAppointmentDeps struct {
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

type AppointmentCreate interface {
	CreateAppointment(appointment *model.Appointment) (interface{}, error)
}

type appointmentCreate struct {
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

func NewAppointmentCreate(
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
) AppointmentCreate {
	return &appointmentCreate{
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

func (l *appointmentCreate) CreateAppointment(appointment *model.Appointment) (interface{}, error) {
	log.Printf("appointment-create-logic -> method: CreateAppointment: received")

	if appointment == nil {
		log.Println("appointment-create-logic -> ERROR: appointment es nil")
		return nil, errors.New("internal server error")
	}

	if !l.appointmentDoctor.IsDoctorExists(appointment.DoctorID) {
		log.Printf("appointment-create-logic -> method: CreateAppointment: Doctor not found with ID: %d", appointment.DoctorID)
		return nil, response.ErrorDoctorNotFoundID
	}

	appointment.Patient.ID = appointment.PatientID

	patient, err := l.getPatient(appointment)
	if err != nil {
		log.Printf("appointment-create-logic -> method: CreateAppointment: Error getting patient: %v", err)
		return nil, err
	}

	if err := l.appointmentTime.ValidateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	priceDetails, err := l.getPriceDetails(appointment, patient)
	if err != nil {
		log.Printf("appointment-create-logic -> method: CreateAppointment: Error getting price details: %v", err)
		return nil, err
	}

	appointmentCreated := l.buildAppointment(appointment, patient)

	if err := l.repositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment-create-logic -> method: CreateAppointment: Error creating appointment: %v", err)
		return nil, err
	}

	return priceDetails, nil
}

func (l *appointmentCreate) getPatient(appointment *model.Appointment) (*model.Patient, error) {
	log.Printf("appointment-create-logic -> method: getPatientt: received")

	if appointment.Patient.ID != 0 {
		patient, err := l.appointmentPatientID.IsPatientIDExists(appointment.Patient.ID)
		if err != nil {
			log.Printf("appointment-create-logic -> method: getPatient: Error getting patient by ID: %v", err)
			return nil, err
		}

		return patient, nil
	}

	err := l.appointmentPatientBody.HandlePatientBodyCreation(appointment)
	if err != nil {
		log.Printf("appointment-create-logic -> method: getPatient: Error creating patient: %v", err)
		return nil, err
	}

	return &appointment.Patient, nil
}

// Método para obtener el precio de un servicio o paquete
func (l *appointmentCreate) getPriceDetails(appointment *model.Appointment, patient *model.Patient) (interface{}, error) {
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

// Método para construir la cita
func (l *appointmentCreate) buildAppointment(appointment *model.Appointment, patient *model.Patient) *model.Appointment {
	log.Printf("appointment-create-logic -> method: buildAppointment: received")

	return &model.Appointment{
		Patient:     *patient,
		DoctorID:    appointment.DoctorID,
		ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: appointment.TotalAmount,
	}
}

/*func (l *appointmentLogic) CreateAppointmentWithService(appointment *model.Appointment) (*model.FinalServicePrice, error) {
	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		return nil, err
	}

	if err := l.handlePatientIDCreation(appointment); err != nil {
		return nil, err
	}

	if err := l.handlePatientBodyCreation(appointment); err != nil {
		return nil, err
	}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalServicePrice, err := l.isServiceIDEXists(appointment.ServiceID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentCreated := &model.Appointment{
		Patient:     appointment.Patient,
		DoctorID:    appointment.DoctorID,
		ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: appointment.TotalAmount,
	}

	if err := l.repo.RepositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment-logic: Error al crear la cita: %v", err)
		return &model.FinalServicePrice{}, err
	}

	return &model.FinalServicePrice{
		TotalAmount:       finalServicePrice.TotalAmount,
		InsuranceDiscount: finalServicePrice.InsuranceDiscount,
		FinalPrice:        finalServicePrice.FinalPrice,
	}, nil
}

func (l *appointmentLogic) CreateAppointmentWithPackage(appointment *model.Appointment) (*model.FinalPackagePriceWithInsegurance, error) {
	if err := l.validateDoctorAvailability(appointment.DoctorID); err != nil {
		log.Printf("appointment-logic: Error fetching doctor with ID %d: %v", appointment.DoctorID, err)
		return nil, err
	}

	//tratando de crear el paciente con body
	//if err := l.handlePatientBodyCreation(appointment); err != nil {
	//	return nil, err
	//}

	if err := l.validateAppointmentTime(appointment); err != nil {
		return nil, err
	}

	finalPkgPrice, err := l.isPackageIDExists(appointment.PackageID, appointment.Patient.Insurance)
	if err != nil {
		return nil, err
	}

	appointmentCreated := &model.Appointment{
		//PatientID:   appointment.PatientID,
		DoctorID: appointment.DoctorID,
		//ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: appointment.TotalAmount,
	}
	//tratando de crear el paciente con ID
	if err := l.handlePatientIDCreation(appointment); err != nil {
		return nil, err
	}

	//tratando de crear el paciente con ID
	if err := l.handlePatientIDCreation(appointmentCreated); err != nil {
		return nil, err
	}

	if err := l.repo.RepositoryAppointment.Create(appointmentCreated); err != nil {
		log.Printf("appointment-logic: Error al crear la cita: %v", err)
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
*/
