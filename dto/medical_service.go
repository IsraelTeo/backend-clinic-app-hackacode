package dto

// MedicalServiceRequest represents a request to create or update a medical service.
type MedicalServiceRequest struct {
	Name        string  `json:"name" gorm:"size:50;not null" validate:"required,max=50"`
	Description string  `json:"description" gorm:"size:250;not null" validate:"required,max=250"`
	Price       float64 `json:"price" validate:"min=0,numeric"`
}

// MedicalServiceResponse represents the response for a medical service.
type MedicalServiceResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// MedicalServiceAppointmentResponse represents a request for a medical service appointment.
type MedicalServiceAppointmentResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
