package dto

// PackageCreateRequest represents a request to create a medical package.
type PackageCreateRequest struct {
	Name              string `json:"name" validate:"required,max=50"`
	MedicalServiceIDs []uint `json:"medical_service_ids" validate:"required"`
}

// PackageResponse represents the response for a medical package.
type PackageResponse struct {
	ID              uint                     `json:"id"`
	Name            string                   `json:"name"`
	MedicalServiceS []MedicalServiceResponse `json:"medical_services"`
}

// PackageAppointmentResponse represents a request for a medical package appointment.
type PackageAppointmentResponse struct {
	Name string `json:"name"`
}
