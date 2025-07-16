package dto

// PersonRequest represents a request to create or update a person.
type PersonRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	LastName    string `json:"last_name" validate:"required,max=80"`
	DNI         string `json:"dni" validate:"required,max=20"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Email       string `json:"email" validate:"required,email,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
	Address     string `json:"address" validate:"required,max=200"`
}

// PersonResponse represents the response for a person.
type PersonResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	DNI         string `json:"dni"`
	BirthDate   string `json:"birth_date"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}
