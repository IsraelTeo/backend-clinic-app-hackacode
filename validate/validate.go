package validate

import (
	"fmt"
	"strconv"
	"time"

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

// recibe un tiempo inicial y un tiempo final en formato de texto
func Times(startTime, endTime string) error {
	now := time.Now()             //tiempo actual
	currentHour := now.Hour()     //Hora actual
	currentMinute := now.Minute() //Minuto actual

	// Transforma el tiempo inicial en un formato texto a un formato de tiempo
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return fmt.Errorf("start_time tiene un formato inválido")
	}

	// Transforma el tiempo final en un formato texto a un formato de tiempo
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return fmt.Errorf("end_time tiene un formato inválido")
	}

	// Verificar si start_time está en el pasado
	if start.Hour() < currentHour {
		// Si la hora de inicio es menor que la hora actual
		return fmt.Errorf("start_time debe ser una hora futura")
	}

	if start.Hour() == currentHour {
		// Si está en la misma hora, verificar los minutos
		if start.Minute() <= currentMinute {
			return fmt.Errorf("start_time debe ser una hora futura")
		}
	}

	// Verificar si end_time está en el pasado
	if end.Hour() < currentHour {
		// Si la hora de finalización es menor que la hora actual
		return fmt.Errorf("end_time debe ser una hora futura")
	}

	if end.Hour() == currentHour {
		// Si está en la misma hora, verificar los minutos
		if end.Minute() <= currentMinute {
			return fmt.Errorf("end_time debe ser una hora futura")
		}
	}
	// Validar que EndTime sea después de StartTime
	if end.Before(start) {
		return fmt.Errorf("end_time debe ser después de start_time")
	}

	return nil
}
