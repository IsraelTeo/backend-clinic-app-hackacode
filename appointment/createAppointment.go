package appointment

import (
	"errors"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentCreate interface {
	CreateAppointment(appointment *model.Appointment) (model.PriceDetails, error)
}

type appointmentCreate struct {
	repositoryAppointment repository.Repository[model.Appointment]
	repositoryDoctor      repository.Repository[model.Doctor]
	repositoryService     repository.Repository[model.Service]
	repositoryPackage     repository.Repository[model.Package]
	repositoryPatient     repository.Repository[model.Patient]
	repositoryPatientMain repository.PatientRepository
	appointmentDoctor     AppointmentDoctorID
	appointmentPackageID  AppointmentPackageID
	appointmentServiceID  AppointmentServiceID
	appointmentTime       AppointmentTime
}

func NewAppointmentCreate(
	repositoryAppointment repository.Repository[model.Appointment],
	repositoryDoctor repository.Repository[model.Doctor],
	repositoryService repository.Repository[model.Service],
	repositoryPackage repository.Repository[model.Package],
	repositoryPatient repository.Repository[model.Patient],
	repositoryPatientMain repository.PatientRepository,
	appointmentDoctor AppointmentDoctorID,
	appointmentPackageID AppointmentPackageID,
	appointmentServiceID AppointmentServiceID,
	appointmentTime AppointmentTime,
) AppointmentCreate {
	return &appointmentCreate{
		repositoryAppointment: repositoryAppointment,
		repositoryDoctor:      repositoryDoctor,
		repositoryService:     repositoryService,
		repositoryPackage:     repositoryPackage,
		repositoryPatient:     repositoryPatient,
		repositoryPatientMain: repositoryPatientMain,
		appointmentDoctor:     appointmentDoctor,
		appointmentPackageID:  appointmentPackageID,
		appointmentServiceID:  appointmentServiceID,
		appointmentTime:       appointmentTime,
	}
}

func (l *appointmentCreate) CreateAppointment(appointment *model.Appointment) (model.PriceDetails, error) {
	if appointment == nil {
		return nil, errors.New("internal server error")
	}

	//Verifica la existencia del médico
	if !l.appointmentDoctor.IsDoctorExists(appointment.DoctorID) {
		return nil, response.ErrorDoctorNotFoundID
	}

	patientFound, err := l.isPatientDNIExists(appointment.PatientDNI)
	if err != nil {
		return nil, err
	}

	err = l.appointmentTime.ValidateAppointmentTime(appointment)
	if err != nil {
		return nil, err
	}

	//obtener precios de paquete o servicio
	priceDetails, err := l.getPriceDetails(appointment, patientFound)
	if err != nil {
		return nil, err
	}

	appointmentCreated := l.buildAppointment(appointment, patientFound, priceDetails)

	err = l.repositoryAppointment.Create(appointmentCreated)
	if err != nil {
		return nil, err
	}

	return priceDetails, nil
}

func (l *appointmentCreate) isPatientDNIExists(DNI string) (*model.Patient, error) {
	patient, err := l.repositoryPatientMain.GetPatientByDNI(DNI)
	if err != nil {
		return nil, response.ErrorPatientNotFoundDNI
	}

	return patient, nil
}

func (l *appointmentCreate) getPriceDetails(appointment *model.Appointment, patient *model.Patient) (model.PriceDetails, error) {
	if appointment.ServiceID != 0 {
		finalServicePrice, err := l.appointmentServiceID.IsServiceIDEXists(appointment.ServiceID, patient.Insurance)
		if err != nil {
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
func (l *appointmentCreate) buildAppointment(appointment *model.Appointment, patient *model.Patient, priceDetails model.PriceDetails) *model.Appointment {
	return &model.Appointment{
		PatientID:   patient.ID,
		DoctorID:    appointment.DoctorID,
		ServiceID:   appointment.ServiceID,
		PackageID:   appointment.PackageID,
		Date:        appointment.Date,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Paid:        false,
		TotalAmount: priceDetails.GetFinalPrice(),
	}
}
