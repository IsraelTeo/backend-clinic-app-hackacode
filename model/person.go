package model

type Person struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `json:"name" validate:"required,max=50"`
	LastName    string `json:"last_name" validate:"required,max=80"`
	DNI         string `json:"dni" validate:"required,max=20"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Email       string `json:"email" validate:"required,email,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
	Address     string `json:"address" validate:"required,max=200"`
}

type Doctor struct {
	Person
	Especialty string  `json:"especialty" validate:"required,max=50"`
	Days       string  `json:"days" validate:"required"`
	StartTime  string  `json:"start_time" validate:"required,regex=^[0-9]{2}:[0-9]{2}$"`
	EndTime    string  `json:"end_time" validate:"required,regex=^[0-9]{2}:[0-9]{2}$"`
	Salary     float64 `json:"salary" validate:"required,numeric"`
}

type Patient struct {
	Person
	Insurance bool `json:"health_insurance" validate:"required"`
}

type Day string

const (
	Moonday   Day = "Lunes"
	Tuesday   Day = "Martes"
	Wednesday Day = "Miercoles"
	Thursday  Day = "Jueves"
	Friday    Day = "Viernes"
	Saturday  Day = "Sabado"
	Sunday    Day = "Domingo"
)
