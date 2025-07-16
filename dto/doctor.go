package dto

// DoctorRequest represents a request to create or update a doctor.
type DoctorRequest struct {
	PersonRequest
	Especialty string  `json:"especialty" validate:"required,max=50"`
	Days       string  `json:"days" validate:"required"`
	StartTime  string  `json:"start_time" validate:"required"`
	EndTime    string  `json:"end_time" validate:"required"`
	Salary     float64 `json:"salary" validate:"required,numeric"`
}

// DoctorResponse represents the response for a doctor.
type DoctorResponse struct {
	PersonResponse
	Especialty string  `json:"especialty"`
	Days       string  `json:"days"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	Salary     float64 `json:"salary"`
}

// DoctorAppointmentRequest represents a request for a doctor's appointment.
type DoctorAppointmentResponse struct {
	FullName string `json:"full_name"`
}
