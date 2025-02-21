package calculation

import "github.com/IsraelTeo/clinic-backend-hackacode-app/model"

func TotalServiceAmount(servicePrice float64, hasInsurance bool) *model.FinalServicePrice {

	var discount float64 = 0.0

	if hasInsurance {
		discount = servicePrice * 0.20
	}

	finalPrice := servicePrice - discount

	return &model.FinalServicePrice{
		TotalAmount:       servicePrice,
		InsuranceDiscount: discount,
		FinalPrice:        finalPrice,
	}
}

func TotalServicePackageAmount(services []model.Service) *model.FinalPackagePrice {
	//verifica que la lista de servicios no esté vacía
	if len(services) == 0 {
		return &model.FinalPackagePrice{
			TotalAmount:     0,
			DiscountPackage: 0,
			FinalPrice:      0,
		}
	}

	var totalAmount float64
	//calcula el precio total
	for _, service := range services {
		totalAmount += service.Price
	}

	discountPackage := totalAmount * 0.15       //descuento por paquete
	finalPrice := totalAmount - discountPackage //precio total con descuento de paquete

	return &model.FinalPackagePrice{
		TotalAmount:     totalAmount,
		DiscountPackage: discountPackage,
		FinalPrice:      finalPrice,
	}
}

func TotalServicePackageAmountToAppointment(services []model.Service, hasInsurance bool) *model.FinalPackagePriceWithInsegurance {
	//verifica que la lista de servicios no esté vacía
	if len(services) == 0 {
		return &model.FinalPackagePriceWithInsegurance{}
	}

	var totalAmount float64
	//calcula el precio total
	for _, service := range services {
		totalAmount += service.Price
	}

	discountPackage := totalAmount * 0.15                      //descuento por paquete
	priceAfterPackageDiscount := totalAmount - discountPackage //precio total con descuento de paquete

	insuranceDiscount := 0.0
	if hasInsurance {
		insuranceDiscount = priceAfterPackageDiscount * 0.20 //descuento por seguro médico
	}

	finalPrice := priceAfterPackageDiscount - insuranceDiscount //precio total con descuento de seguro médico

	return &model.FinalPackagePriceWithInsegurance{
		InsuranceDiscount: insuranceDiscount,
		FinalPackagePrice: model.FinalPackagePrice{
			TotalAmount:     totalAmount,
			DiscountPackage: discountPackage,
			FinalPrice:      finalPrice,
		},
	}
}

func CalculateTotalAmount(services []model.Service, isPackage bool, hasInsurance bool) float64 {
	var totalAmount float64
	for _, service := range services {
		totalAmount += service.Price
	}

	if isPackage {
		totalAmount *= 0.85
	}

	if hasInsurance {
		totalAmount *= 0.80
	}

	return totalAmount
}
