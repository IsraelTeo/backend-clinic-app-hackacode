package appointment

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/calculation"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type appointmentPackageID struct {
	repositoryPackageMain repository.PackageRepository
}

type AppointmentPackageID interface {
	IsPackageIDExists(ID uint, hasInsurance bool) (*model.FinalPackagePriceWithInsegurance, error)
}

func NewAppointmentPackageID(repositoryPackageMain repository.PackageRepository) AppointmentPackageID {
	return &appointmentPackageID{repositoryPackageMain: repositoryPackageMain}
}

func (l *appointmentPackageID) IsPackageIDExists(ID uint, hasInsurance bool) (*model.FinalPackagePriceWithInsegurance, error) {
	pkg, err := l.repositoryPackageMain.GetByID(ID)
	if err != nil || pkg == nil {
		return nil, response.ErrorPackageNotFound
	}

	finalPricePkg := calculation.TotalServicePackageAmountToAppointment(pkg.Services, hasInsurance)

	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: finalPricePkg.InsuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     finalPricePkg.TotalAmount,
			DiscountPackage: finalPricePkg.DiscountPackage,
			FinalPrice:      finalPricePkg.FinalPrice,
		},
	}, nil
}
