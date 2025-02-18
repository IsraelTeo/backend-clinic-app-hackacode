package auth

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/labstack/echo/v4"
)

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := validateToken(c)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return response.WriteError(&response.WriteResponse{
				C:       c,
				Message: err.Error(),
				Status:  http.StatusUnauthorized,
				Data:    nil,
			})
		}

		return next(c)
	}
}
