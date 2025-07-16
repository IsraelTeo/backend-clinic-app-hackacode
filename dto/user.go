package dto

//
type UserRegister struct {
	UserName string `json:"user_name" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}
