package model

//Cita médica
type Appoiment struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"-"`
	PatientID   uint    `json:"patient_id" validate:"required"`
	Patient     Patient `json:"patient"`
	DoctorID    uint    `json:"doctor_id" validate:"required"`
	ServiceID   uint    `json:"service_id" validate:"required"`
	PackageID   uint    `json:"package_id" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Time        string  `json:"start_time" validate:"required"`
	Paid        bool    `json:"paid" default:"false"`
	TotalAmount float64 `json:"total_amount"`
}

//Pago
type Payment struct {
	AppoimentID uint        `json:"appoiment_id" validate:"required"`
	PaymentType PaymentType `json:"payment_type" validate:"required"`
}

//Método de pago
type PaymentType string

const (
	Cash        PaymentType = "efectivo"
	Card        PaymentType = "card"
	Application PaymentType = "applicativo"
)

//responder el QR
//factura PDF
