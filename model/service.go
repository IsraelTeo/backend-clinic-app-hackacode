package model

//Servicio médico
type Service struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"size:50;not null" validate:"required,max=50"`
	Description string  `json:"description" gorm:"size:250;not null" validate:"required,max=250"`
	Price       float64 `json:"price" validate:"min=0,numeric"`
}

//Paquete de servicios médicos
type Package struct {
	ID       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string    `json:"name"`
	Services []Service `json:"services" gorm:"many2many:package_services;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Price    float64   `json:"price"`
}

//Creación de paquete médico
type CreatePackageRequest struct {
	Name       string `json:"name" validate:"required,max=50"`
	ServiceIDs []uint `json:"service_ids" validate:"required"`
}

type PriceDetails interface {
	GetFinalPrice() float64
}

//Precio final de servicio médico con descuento por seguro médico del paciente
type FinalServicePrice struct {
	TotalAmount       float64
	InsuranceDiscount float64
	FinalPrice        float64
}

func (f *FinalServicePrice) GetFinalPrice() float64 {
	return f.FinalPrice
}

//Precio final de paquete médico con descuento por paquete
type FinalPackagePrice struct {
	TotalAmount     float64
	DiscountPackage float64
	FinalPrice      float64
}

func (f *FinalPackagePrice) GetFinalPrice() float64 {
	return f.FinalPrice
}

// Precio final de paquete médico con descuento por paquete y con descuento de seguro médico del paciente
type FinalPackagePriceWithInsegurance struct {
	InsuranceDiscount float64
	FinalPackagePrice
}

func (f *FinalPackagePriceWithInsegurance) GetFinalPrice() float64 {
	return f.FinalPackagePrice.FinalPrice
}
