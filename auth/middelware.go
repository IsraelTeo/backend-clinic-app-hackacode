package auth

import (
	"log"
	"net/http"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/labstack/echo/v4"
)

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := validateToken(c)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return response.WriteError(c, "Invalid token.", http.StatusUnauthorized)
		}
		return next(c)
	}
}
