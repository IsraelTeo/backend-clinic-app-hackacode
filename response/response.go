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

const (
	// Mensajes de éxito para paquetes
	SuccessPackageFound   = "¡Paquete encontrado exitosamente!"
	SuccessPackageUpdated = "¡Paquete actualizado exitosamente!"
	SuccessPackagesFound  = "¡Paquetes encontrados exitosamente!"
	SuccessPackageCreated = "¡Paquete creado exitosamente!"
	SuccessPackageDeleted = "¡Paquete eliminado exitosamente!"
)

var (
	//Mensaje de error generales
	ErrorInvalidID = errors.New("el ID debe ser un número positivo")
)
var (
	// Mensajes de error para servicios médicos
	ErrorServiceNotFound   = errors.New("el servicio médico no fue encontrado")
	ErrorServicesNotFound  = errors.New("no fueron encontrados servicios médicos")
	ErrorListServicesEmpty = errors.New("no fueron encontrados servicios médicos")
	ErrorBadRequestService = errors.New("el cuerpo de la solicitud no es válido para el servicio médico")
	ErrorToCreatedService  = errors.New("no se pudo crear el servicio médico")
	ErrorToUpdatedService  = errors.New("no se pudo actualizar el servicio médico")
	ErrorToDeletedService  = errors.New("no se pudo eliminar el servicio médico")
)

var (
	// Mensajes de error para paquetes

	ErrorPackageNotFound   = errors.New("el paquete no fue encontrado")
	ErrorPackagesNotFound  = errors.New("no fueron encontrados paquetes")
	ErrorListPackagesEmpty = errors.New("no fueron encontrados paquetes")
	ErrorBadRequestPackage = errors.New("el cuerpo de la solicitud no es válido para el paquete")
	ErrorToCreatedPackage  = errors.New("no se pudo crear el paquete")
	ErrorToUpdatedPackage  = errors.New("no se pudo actualizar el paquete")
	ErrorToDeletedPackage  = errors.New("no se pudo eliminar el paquete")
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
