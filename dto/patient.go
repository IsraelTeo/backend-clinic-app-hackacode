package dto

// PatientRequest represents a request to create or update a patient.
type PatientRequest struct {
	PersonRequest
	HealthInsurance bool `json:"health_insurance"`
}

// PatientResponse represents the response for a patient.
type PatientResponse struct {
	PersonResponse
	HealthInsurance bool `json:"health_insurance"`
}

// PatientAppointmentResponse represents a response for a patient's appointment.
type PatientAppointmentResponse struct {
	FullName string `json:"full_name"`
	DNI      string `json:"dni"`
}
