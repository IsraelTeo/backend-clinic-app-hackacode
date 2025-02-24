package appointment

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type AppointmentUpdate interface {
	UpdateAppointment(appointmentID uint, updatedAppointment *model.Appointment) (model.PriceDetails, error)
}

type appointmentUpdate struct {
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

func NewAppointmentUpdate(
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
) AppointmentUpdate {
	return &appointmentUpdate{
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

func (l *appointmentUpdate) UpdateAppointment(ID uint, updatedAppointment *model.Appointment) (model.PriceDetails, error) {
	existingAppointment, err := l.repositoryAppointment.GetByID(ID)
	if err != nil {
		return nil, response.ErrorAppointmentNotFound
	}

	if !l.appointmentDoctor.IsDoctorExists(updatedAppointment.DoctorID) {
		return nil, response.ErrorDoctorNotFoundID
	}

	patientFound, err := l.isPatientDNIExists(updatedAppointment.PatientDNI)
	if err != nil {
		return nil, err
	}

	err = l.appointmentTime.ValidateAppointmentTime(updatedAppointment)
	if err != nil {
		return nil, err
	}

	priceDetails, err := l.getPriceDetails(updatedAppointment, patientFound)
	if err != nil {
		return nil, err
	}

	// Construir la cita actualizada
	updatedAppointmentData := l.buildUpdatedAppointment(existingAppointment, updatedAppointment, patientFound, priceDetails)

	err = l.repositoryAppointment.Update(updatedAppointmentData)
	if err != nil {
		return nil, err
	}

	return priceDetails, nil
}

func (l *appointmentUpdate) isPatientDNIExists(DNI string) (*model.Patient, error) {
	patient, err := l.repositoryPatientMain.GetPatientByDNI(DNI)
	if err != nil {
		return nil, response.ErrorPatientNotFoundDNI
	}

	return patient, nil
}

func (l *appointmentUpdate) getPriceDetails(appointment *model.Appointment, patient *model.Patient) (model.PriceDetails, error) {
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

// Método para construir la cita actualizada
func (l *appointmentUpdate) buildUpdatedAppointment(existingAppointment, updatedAppointment *model.Appointment, patient *model.Patient, priceDetails model.PriceDetails) *model.Appointment {
	return &model.Appointment{
		ID:          existingAppointment.ID,
		PatientID:   patient.ID,
		DoctorID:    updatedAppointment.DoctorID,
		ServiceID:   updatedAppointment.ServiceID,
		PackageID:   updatedAppointment.PackageID,
		Date:        updatedAppointment.Date,
		StartTime:   updatedAppointment.StartTime,
		EndTime:     updatedAppointment.EndTime,
		Paid:        existingAppointment.Paid,
		TotalAmount: priceDetails.GetFinalPrice(),
	}
}
