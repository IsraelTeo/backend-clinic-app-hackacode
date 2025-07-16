package dto

// AppointmentRequest represents a request to create or update a medical appointment.
type AppointmentRequest struct {
	DoctorID  uint   `json:"doctor_id" validate:"required"`
	PatientID uint   `json:"patient_id" validate:"required"`
	ServiceID uint   `json:"service_id"`
	PackageID uint   `json:"package_id"`
	Date      string `json:"date" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

// AppointmentResponse represents the response for a medical appointment.
type AppointmentResponse struct {
	ID             uint                               `json:"id"`
	Doctor         DoctorResponse                     `json:"doctor"`
	Patient        PatientResponse                    `json:"patient"`
	MedicalService *MedicalServiceAppointmentResponse `json:"service,omitempty"`
	MedicalPackage *PackageResponse                   `json:"package,omitempty"`
	Date           string                             `json:"date"`
	StartTime      string                             `json:"start_time"`
	EndTime        string                             `json:"end_time"`
	TotalAmount    float64                            `json:"total_amount"`
}
