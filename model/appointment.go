package model

import "time"

// Appointment represents a medical appointment in the system.
type Appointment struct {
	ID               uint            `gorm:"primaryKey;autoIncrement"`
	DoctorID         uint            `gorm:"not null"`
	Doctor           *Doctor         `gorm:"foreignKey:DoctorID"`
	PatientID        uint            `gorm:"not null"`
	Patient          *Patient        `gorm:"foreignKey:PatientID"`
	MedicalServiceID *uint           `gorm:"default:null"`
	MedicalService   *MedicalService `gorm:"foreignKey:MedicalServiceID"`
	PackageID        *uint           `gorm:"default:null"`
	MedicalPackage   *MedicalPackage `gorm:"foreignKey:PackageID"`
	Date             time.Time       `gorm:"not null"`
	StartTime        time.Time       `gorm:"not null"`
	EndTime          time.Time       `gorm:"not null"`
	TotalAmount      float64         `gorm:"not null"`
	CreatedAt        time.Time       `gorm:"not null"`
}

/*
// Payment represents a payment for an appointment.
type Payment struct {
	AppoimentID uint        `json:"appoiment_id" validate:"required"`
	Paid        bool        `json:"paid" validate:"required"`
	TotalAmount float64     `json:"total_amount"`
	PaymentType PaymentType `json:"payment_type" validate:"required"`
}

// Método de pago
type PaymentType string

const (
	Cash        PaymentType = "efectivo"
	Card        PaymentType = "tarjeta"
	Application PaymentType = "applicativo"
)

// Respuesta al realizar el pago
type PaymentResponse struct {
	QRCode     string `json:"qr_code"`
	PDFReceipt string `json:"pdf_receipt"`
}*/
