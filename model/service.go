package model

type Service struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"size:50;not null" validate:"required,max=50"`
	Description string  `json:"description" validate:"required,max=250"`
	Price       float64 `json:"price" validate:"min=0,numeric"`
}

type Package struct {
	ID       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string    `json:"name" gorm:"size:50;not null" validate:"required,max=50"`
	Services []Service `json:"services" gorm:"many2many:package_services;"`
	Price    float64   `json:"price"`
}

type CreatePackageRequest struct {
	Name       string `json:"name" validate:"required,max=50"`
	ServiceIDs []uint `json:"service_ids" validate:"required"`
}
