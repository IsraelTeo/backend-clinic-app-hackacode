package validate

import (
	"strconv"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ParseID(c echo.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, response.ErrorInvalidID
	}
	return uint(id), nil
}

// CustomValidator es una estructura que envuelve el validador.
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate implementa la interfaz de validación de Echo.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func Init() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}
