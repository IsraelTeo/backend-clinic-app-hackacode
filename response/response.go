package response

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// Mensaje de error generales
var (
	ErrorInvalidID = errors.New("el ID debe ser un número positivo")
)

// Mensajes de exito para autenticación
const (
	SuccessLogin = "¡Inicio de sesión exitoso!"
)

// Mensajes de error para autenticación
var (
	ErrorBadRequestUser  = errors.New("el cuerpo de la solicitud no es válido para el usuario")
	ErrorBadCretendials  = errors.New("credenciales inválidas")
	ErrorInvalidEmail    = errors.New("el email no es válido o no está registrado")
	ErrorGeneratingToken = errors.New("no se pudo generar el token de autenticación")
)

// Mensajes de éxito para servicios médicos
const (
	SuccessServiceFound   = "¡Servicio médico encontrado exitosamente!"
	SuccessServiceUpdated = "¡Servicio médico actualizado exitosamente!"
	SuccessServicesFound  = "¡Servicios médicos encontrados exitosamente!"
	SuccessServiceCreated = "¡Servicio médico creado exitosamente!"
	SuccessServiceDeleted = "¡Servicio médico eliminado exitosamente!"
)

// Mensajes de éxito para paquetes
const (
	SuccessPackageFound   = "¡Paquete encontrado exitosamente!"
	SuccessPackageUpdated = "¡Paquete actualizado exitosamente!"
	SuccessPackagesFound  = "¡Paquetes encontrados exitosamente!"
	SuccessPackageCreated = "¡Paquete creado exitosamente!"
	SuccessPackageDeleted = "¡Paquete eliminado exitosamente!"
)

// Mensajes de error para servicios médicos
var (
	ErrorServiceNotFound   = errors.New("el servicio médico no fue encontrado")
	ErrorServicesNotFound  = errors.New("no fueron encontrados servicios médicos")
	ErrorListServicesEmpty = errors.New("no fueron encontrados servicios médicos")
	ErrorBadRequestService = errors.New("el cuerpo de la solicitud no es válido para el servicio médico")
	ErrorToCreatedService  = errors.New("no se pudo crear el servicio médico")
	ErrorToUpdatedService  = errors.New("no se pudo actualizar el servicio médico")
	ErrorToDeletedService  = errors.New("no se pudo eliminar el servicio médico")
)

// Mensajes de error para paquetes
var (
	ErrorPackageNotFound   = errors.New("el paquete no fue encontrado")
	ErrorPackagesNotFound  = errors.New("no fueron encontrados paquetes")
	ErrorListPackagesEmpty = errors.New("no fueron encontrados paquetes")
	ErrorBadRequestPackage = errors.New("el cuerpo de la solicitud no es válido para el paquete")
	ErrorToCreatedPackage  = errors.New("no se pudo crear el paquete")
	ErrorToUpdatedPackage  = errors.New("no se pudo actualizar el paquete")
	ErrorToDeletedPackage  = errors.New("no se pudo eliminar el paquete")
)

// Mensajes de exito de doctores

const (
	SuccessDoctorFound   = "¡Doctor encontrado exitosamente!"
	SuccessDoctorUpdated = "¡Doctor actualizado exitosamente!"
	SuccessDoctorsFound  = "¡Doctores encontrados exitosamente!"
	SuccessDoctorCreated = "¡Doctor creado exitosamente!"
	SuccessDoctorDeleted = "¡Doctor eliminado exitosamente!"
)

// Mensajes de error para doctores
var (
	ErrorDoctorNotFound   = errors.New("el doctor no fue encontrado")
	ErrorListDoctorsEmpty = errors.New("no fueron encontrados doctores")
	ErrorToCreatedDoctor  = errors.New("no se pudo crear el doctor")
	ErrorToUpdatedDoctor  = errors.New("no se pudo actualizar el doctor")
	ErrorToDeletedDoctor  = errors.New("no se pudo eliminar el doctor")
)

// Mensajes de exito de pacientes

const (
	SuccessPatientFound   = "¡Paciente encontrado exitosamente!"
	SuccessPatientUpdated = "¡Paciente actualizado exitosamente!"
	SuccessPatientsFound  = "¡Paciente encontrados exitosamente!"
	SuccessPatientCreated = "¡Paciente creado exitosamente!"
	SuccessPatientDeleted = "¡Pacienteeliminado exitosamente!"
)

// Mensajes de error para pacientes
var (
	ErrorPatientNotFound   = errors.New("el paciente no fue encontrado")
	ErrorListPatientsEmpty = errors.New("no fueron encontrados pacientes")
	ErrorToCreatedPatient  = errors.New("no se pudo crear el paciente")
	ErrorToUpdatedPatient  = errors.New("no se pudo actualizar el paciente")
	ErrorToDeletedPatient  = errors.New("no se pudo eliminar el paciente")
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
