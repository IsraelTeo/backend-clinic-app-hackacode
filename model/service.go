package model

type Service struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"-"`
	Name        string  `gorm:"size:50;not null" json:"name" validate:"required,max=50"`
	Description string  `json:"description" validate:"required,max=250"`
	Price       float64 `json:"price" validate:"min=0,numeric"`
}

type Package struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	Name     string    `gorm:"size:50;not null" json:"name" validate:"required,max=50"`
	Services []Service `gorm:"many2many:package_services;" json:"services"`
	Price    float64   `json:"-"`
}

func (p *Package) CalculatePrice() {
	total := 0.0
	for _, service := range p.Services {
		total += service.Price
	}
	p.Price = total
}
