package calculation

import "gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"

func TotalServicePackageAmount(services []model.Service, hasInsurance bool) (float64, float64, float64) {
	if len(services) == 0 {
		return 0, 0, 0
	}

	var totalAmount float64
	for _, service := range services {
		totalAmount += service.Price
	}

	discountPackage := totalAmount * 0.15
	priceAfterPackageDiscount := totalAmount - discountPackage

	insuranceDiscount := 0.0
	if hasInsurance {
		insuranceDiscount = priceAfterPackageDiscount * 0.20
	}

	finalPrice := priceAfterPackageDiscount - insuranceDiscount

	return totalAmount, insuranceDiscount, finalPrice
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
