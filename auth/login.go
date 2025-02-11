package auth

import (
	"log"
	"net/http"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/labstack/echo/v4"
)

type LoginService interface {
	Login(c echo.Context) error
}

type loginService struct {
	repository repository.UserRepository
}

type userLogin struct {
	Email    string `json:"email" gorm:"size:50;not null" validate:"required, email"`
	Password string `json:"password" gorm:"size:100;not null" validate:"required"`
}

func NewLoginService(repository repository.UserRepository) LoginService {
	return &loginService{repository: repository}
}

func (s *loginService) Login(c echo.Context) error {
	log.Println("login-service: Request received in Login")

	user := userLogin{}
	if err := c.Bind(&user); err != nil {
		log.Printf("login-service: Error binding user input: %v", err)
		return response.WriteError(c, response.ErrorBadRequestUser.Error(), http.StatusBadRequest)
	}

	userFound, err := s.repository.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("login-service: User not found with email: %s, error: %v", user.Email, err)
		return response.WriteError(c, response.ErrorInvalidEmail.Error(), http.StatusUnauthorized)
	}

	if !comparePassword(user.Password, userFound.Password) {
		log.Printf("login-service: Password mismatch for user: %s", user.Email)
		return response.WriteError(c, response.ErrorBadCretendials.Error(), http.StatusUnauthorized)
	}

	log.Printf("[INFO] LoginService: Generating token for user: %s", user.Email)
	token, err := generateToken(userFound)
	if err != nil {
		log.Printf("login-service: Error generating token: %v", err)
		return response.WriteError(c, response.ErrorGeneratingToken.Error(), http.StatusInternalServerError)
	}

	log.Printf("login-service: Login successful for user: %s", user.Email)
	return response.WriteSuccess(c, response.SuccessLogin, http.StatusOK, token)
}
