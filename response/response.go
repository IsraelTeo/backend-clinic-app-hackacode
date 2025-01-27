package response

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// Definir el tipo personalizado para Servicios Médicos

const (
	// Mensajes de éxito para servicios médicos
	SuccessServiceFound   = "¡Servicio médico encontrado exitosamente!"
	SuccessServiceUpdated = "¡Servicio médico actualizado exitosamente!"
	SuccessServicesFound  = "¡Servicios médicos encontrados exitosamente!"
	SuccessServiceCreated = "¡Servicio médico creado exitosamente!"
	SuccessServiceDeleted = "¡Servicio médico eliminado exitosamente!"
)

var (
	// Mensajes de error para servicios médicos

	ErrorInvalidId         = errors.New("el ID debe ser un número positivo")
	ErrorServiceNotFound   = errors.New("el servicio médico no fue encontrado")
	ErrorServicesNotFound  = errors.New("no fueron encontrados servicios médicos")
	ErrorListServicesEmpty = errors.New("no fueron encontrados servicios médicos")
	ErrorBadRequest        = errors.New("el cuerpo de la solicitud no es válido para el servicio médico")
	ErrorToCreated         = errors.New("no se pudo crear el servicio médico")
	ErrorToUpdated         = errors.New("no se pudo actualizar el servicio médico")
	ErrorToDeleted         = errors.New("no se pudo eliminar el servicio médico")
)

func WriteSuccess(c echo.Context, message string, status int, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func WriteError(c echo.Context, message string, status int) error {
	return c.JSON(status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}
