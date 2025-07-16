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

/*
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
*/
