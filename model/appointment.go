package model

import "time"

//Cita médica
type Appointment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	PatientID   uint      `json:"patient_id"`
	Patient     Patient   `json:"patient" validate:"required"`
	DoctorID    uint      `json:"doctor_id" validate:"required"`
	ServiceID   uint      `json:"service_id" validate:"required"`
	PackageID   uint      `json:"package_id" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	Paid        bool      `json:"paid" default:"false"`
	TotalAmount float64   `json:"total_amount"`
}

//al registrar  se da de respuesta

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
//factura PDF
