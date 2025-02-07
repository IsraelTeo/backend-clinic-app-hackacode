package model

//Cita médica
type Appointment struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID   uint    `json:"patient_id" gorm:"not null;constraint:OnDelete:CASCADE"`
	Patient     Patient `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID    uint    `json:"doctor_id" validate:"required"`
	ServiceID   uint    `json:"service_id"`
	PackageID   uint    `json:"package_id"`
	Date        string  `json:"date" validate:"required"`
	StartTime   string  `json:"start_time" validate:"required"`
	EndTime     string  `json:"end_time" validate:"required"`
	Paid        bool    `json:"paid" default:"false"`
	TotalAmount float64 `json:"total_amount"`
}

//Pago
type Payment struct {
	AppoimentID uint        `json:"appoiment_id" validate:"required"`
	Paid        bool        //bool FALSE -> ELIMINO LA CITA LA BD Y SI ME LLEGA TRUE EL VALOR DE CITA PAID TRUE
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
//comprobante de pago PDF
