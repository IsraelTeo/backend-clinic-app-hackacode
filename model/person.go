package model

type Person struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"` //
	Name        string `json:"name" validate:"required,max=50"`
	LastName    string `json:"last_name" validate:"required,max=80"`
	DNI         string `json:"dni" validate:"required,max=20"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Email       string `json:"email" validate:"required,email,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
	Address     string `json:"address" validate:"required,max=200"`
}

// Médico
type Doctor struct {
	Person
	Especialty string  `json:"especialty" validate:"required,max=50"`
	Days       string  `json:"days" validate:"required"`
	StartTime  string  `json:"start_time" validate:"required"`
	EndTime    string  `json:"end_time" validate:"required"`
	Salary     float64 `json:"salary" validate:"required,numeric"`
}

// Paciente
type Patient struct {
	Person
	Insurance bool `json:"health_insurance"`
}

// Días válidos para que trabaje el doctor
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
