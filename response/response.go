package response

import "github.com/labstack/echo/v4"

// Definir el tipo personalizado para Servicios Médicos

type MedicalServiceSuccess string

const (
	// Mensajes de éxito para servicios médicos
	SuccessServiceFound   MedicalServiceSuccess = "¡Servicio médico encontrado exitosamente!"
	SuccessServiceUpdated MedicalServiceSuccess = "¡Servicio médico actualizado exitosamente!"
	SuccessServicesFound  MedicalServiceSuccess = "¡Servicios médicos encontrados exitosamente!"
	SuccessServiceCreated MedicalServiceSuccess = "¡Servicio médico creado exitosamente!"
	SuccessServiceDeleted MedicalServiceSuccess = "¡Servicio médico eliminado exitosamente!"
)

type MedicalServiceError string

const (
	// Mensajes de error para servicios médicos
	ErrorInvalidId         MedicalServiceError = "El ID debe ser un número positivo."
	ErrorServiceNotFound   MedicalServiceError = "El servicio médico no fue encontrado."
	ErrorServicesNotFound  MedicalServiceError = "No fueron encontrados servicios médicos."
	ErrorListServicesEmpty MedicalServiceError = "No fueron encontrados servicios médicos."
	ErrorBadRequest        MedicalServiceError = "El cuerpo de la solicitud no es válido para el servicio médico."
	ErrorToCreated         MedicalServiceError = "No se pudo crear el servicio médico."
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
