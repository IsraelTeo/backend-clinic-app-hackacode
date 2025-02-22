package model

//Cita médica
type Appointment struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID   uint     `json:"patient_id" gorm:"constraint:OnDelete:CASCADE"`
	Patient     *Patient `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID    uint     `json:"doctor_id" validate:"required"`
	ServiceID   uint     `json:"service_id"`
	PackageID   uint     `json:"package_id"`
	Date        string   `json:"date" validate:"required"`
	StartTime   string   `json:"start_time" validate:"required"`
	EndTime     string   `json:"end_time" validate:"required"`
	Paid        bool     `json:"paid"`
	TotalAmount float64  `json:"total_amount"`
}

//Pago
type Payment struct {
	AppoimentID uint `json:"appoiment_id" validate:"required"`
	Paid        bool
	TotalAmount float64     `json:"total_amount"`
	PaymentType PaymentType `json:"payment_type" validate:"required"`
}

//Método de pago
type PaymentType string

const (
	Cash        PaymentType = "efectivo"
	Card        PaymentType = "tarjeta"
	Application PaymentType = "applicativo"
)

//Respuesta al realizar el pago
type PaymentResponse struct {
	QRCode     string `json:"qr_code"`
	PDFReceipt string `json:"pdf_receipt"`
}
