package response

import (
	"errors"
	"fmt"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/labstack/echo/v4"
)

// Mensaje de error generales
var (
	ErrorInvalidID  = errors.New("el ID debe ser un número positivo")
	ErrorBadRequest = errors.New("algún campo está faltando")
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

// Mensajes de error para generación y validación de tokens
var (
	ErrorEnvVariablesInvalid    = errors.New("las variables de entorno JWT_EXP y API_SECRET son nulas o inválidas")
	ErrorJWTExpInvalid          = errors.New("el valor de JWT_EXP es inválido, se utilizará el valor predeterminado de 1 hora (3600 segundos)")
	ErrorSigningToken           = errors.New("error al firmar el token")
	ErrorAuthorizationHeader    = errors.New("encabezado de autorización no encontrado")
	ErrorAuthorizationHeaderFmt = errors.New("formato del encabezado de autorización inválido")
	ErrorTokenInvalid           = errors.New("el token no es válido")
	ErrorTokenClaimsInvalid     = errors.New("no se pudo recuperar la información del payload o el token es inválido")
	ErrorTokenEmailInvalid      = errors.New("el campo 'email' está ausente o no es válido en los claims del token")
	ErrorSigningMethodInvalid   = errors.New("el método de firma del token no es válido")
	ErrorTokenMissingInRequest  = errors.New("no se encontró un token en la solicitud")
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
	ErrorPackageNotFound    = errors.New("el paquete no fue encontrado")
	ErrorPackagesNotFound   = errors.New("no fueron encontrados paquetes")
	ErrorListPackagesEmpty  = errors.New("no fueron encontrados paquetes")
	ErrorBadRequestPackage  = errors.New("el cuerpo de la solicitud no es válido para el paquete")
	ErrorToCreatedPackage   = errors.New("no se pudo crear el paquete")
	ErrorToUpdatedPackage   = errors.New("no se pudo actualizar el paquete")
	ErrorToDeletedPackage   = errors.New("no se pudo eliminar el paquete")
	ErrorNoServicesProvided = errors.New("no se proporcionaron servicios para el paquete")
	ErrorFetchingServices   = errors.New("no se pudieron obtener los servicios para el paquete")
)

// Mensajes de exito de doctores

const (
	SuccessDoctorFound   = "¡Médico encontrado exitosamente!"
	SuccessDoctorUpdated = "¡Médico actualizado exitosamente!"
	SuccessDoctorsFound  = "¡Médicos encontrados exitosamente!"
	SuccessDoctorCreated = "¡Médico creado exitosamente!"
	SuccessDoctorDeleted = "¡Médico eliminado exitosamente!"
)

// Mensajes de error para doctores
var (
	ErrorDoctorNotFoundID        = errors.New("el médico no fue encontrado con el ID proporcionado")
	ErrorDoctorDNIRequired       = errors.New("el médico no fue encontrado con el DNI proporcionado")
	ErrorDoctorNotFoundDNI       = errors.New("el médico no fue encontrado con el DNI proporcionado")
	ErrorDoctorsNotFound         = errors.New("no fueron encontrados médicos")
	ErrorListDoctorsEmpty        = errors.New("no fueron encontrados médicos")
	ErrorToCreatedDoctor         = errors.New("no se pudo crear el médico")
	ErrorToUpdatedDoctor         = errors.New("no se pudo actualizar el médico")
	ErrorToDeletedDoctor         = errors.New("no se pudo eliminar el médico")
	ErrorInvalidDoctorTime       = errors.New("fecha inválida")
	ErrorDoctorExistsDNI         = errors.New("ya existe un médico con el DNI ingresado")
	ErrorDoctorExistsPhoneNumber = errors.New("el número telefónico ingresado ya existe")
	ErrorDoctorExistsEmail       = errors.New("el email ingresado ya existe")
	ErrorDoctorInvalidDateFormat = errors.New("ingrese el formato adecuado para la fecha de nacimiento del médico")
	ErrorDoctorBrithDateIsFuture = errors.New("la fecha de cumpleaños debe ser en tiempo pasado")
)

// Mensajes de exito de pacientes

const (
	SuccessPatientFound   = "¡Paciente encontrado exitosamente!"
	SuccessPatientUpdated = "¡Paciente actualizado exitosamente!"
	SuccessPatientsFound  = "¡Pacientes encontrados exitosamente!"
	SuccessPatientCreated = "¡Paciente registrado exitosamente!"
	SuccessPatientDeleted = "¡Paciente eliminado exitosamente!"
)

// Mensajes de error para pacientes
var (
	ErrorPatientNotFoundID        = errors.New("el paciente no fue encontrado con el ID proporcionado")
	ErrorPatientNotFoundDNI       = errors.New("el paciente no fue encontrado con el DNI proporcionado")
	ErrorPatientExistsDNI         = errors.New("ya existe un paciente con el DNI ingresado")
	ErrorPatientExistsPhoneNumber = errors.New("el número telefónico ingresado ya existe")
	ErrorPatientExistsEmail       = errors.New("el email ingresado ya existe")
	ErrorPatientsNotFound         = errors.New("no fueron encontrados pacientes")
	ErrorListPatientsEmpty        = errors.New("no fueron encontrados pacientes")
	ErrorToCreatedPatient         = errors.New("no se pudo crear el paciente")
	ErrorToUpdatedPatient         = errors.New("no se pudo actualizar el paciente")
	ErrorToDeletedPatient         = errors.New("no se pudo eliminar el paciente")
	ErrorPatientInvalidDateFormat = errors.New("ingrese el formato adecuado para la fecha de nacimiento del paciente")
	ErrorPatientBrithDateIsFuture = errors.New("la fecha de cumpleaños debe ser en tiempo pasado")
	ErrorUnlinkingAppointments    = errors.New("no se pudo desvincular las citas del paciente")
	ErrorPatientDNIRequired       = errors.New("el DNI es requerido")
)

// Mensajes de exito de citas

const (
	SuccessAppointmentFound   = "¡Cita encontrada exitosamente!"
	SuccessAppointmentUpdated = "¡Cita actualizada exitosamente!"
	SuccessAppointmentsFound  = "¡Citas encontradas exitosamente!"
	SuccessAppointmentCreated = "¡Cita registrada exitosamente, proceda a realizar el pago!"
	SuccessAppointmentDeleted = "¡Cita eliminada exitosamente!"
)

// Mensajes de error para citas
var (
	ErrorAppointmentNotFound          = errors.New("no se encontró la cita especificada en el sistema")
	ErrorAppointmetsNotFound          = errors.New("no se encontraron citas registradas")
	ErrorListAppointmentsEmpty        = errors.New("no se encontró ninguna cita disponible")
	ErrorToCreatedAppointment         = errors.New("hubo un error al intentar registrar la cita; por favor, verifique los datos ingresados")
	ErrorToUpdatedAppointment         = errors.New("no se pudo actualizar la información de la cita; intente nuevamente")
	ErrorToDeletedAppointment         = errors.New("hubo un error al intentar eliminar la cita; intente nuevamente")
	ErrorAppointmentDateInPast        = errors.New("la fecha de la cita no puede ser en el pasado; por favor, elija una fecha futura")
	ErrorAppointmentInvalidDateFormat = errors.New("el formato de la fecha ingresada no es válido; use el formato AAAA-MM-DD")
	ErrorAppointmentDayNotAvailable   = errors.New("el médico no tiene disponibilidad para el día seleccionado")
	ErrorInvalidAppointmentTime       = errors.New("el horario de la cita no coincide con el horario laboral del médico")
	ErrorAppointmentTimeConflict      = errors.New("el horario de la cita tiene conflictos con otra cita programada para el mismo médico")
	ErrorInvalidAppointmentTimeRange  = errors.New("el rango de tiempo especificado para la cita no es válido; asegúrese de que la hora de inicio sea anterior a la de finalización")
	ErrorAppointmentTimeFormat        = errors.New("el formato de hora ingresado no es válido; use el formato HH:MM")
	ErrorPatientExists                = errors.New("el paciente ya fue registrado anteriormente, solo ingresa su id")
	ErrorPatientDataRequired          = errors.New("se requieren datos del paciente")
	ErrorInvalidAppointment           = errors.New("debe seleccionar al menos un paquete o servicio para la cita")
	ErrorPackageAndServiceEmpty       = errors.New("ni paquete ni servicio especificado")
)

// Mensajes de éxito de pago realizado
const (
	SuccessPaymentRegister = "Pago registrado exitosamente"
)

// Mensajes de error del pago
var (
	ErrorPaidNotTrue           = errors.New("el pago debe ser confirmado")
	ErrorTotalAmountEmpty      = errors.New("por favor ingresar la cantidad de dinero")
	ErrorTotalAmountBadRequest = errors.New("por favor ingresar la cantidad de dinero adecuada")
	ErrorToUpdatePaid          = errors.New("error al actualizar el estado del pago")
	ErrorGeneratingQRCode      = errors.New("error al generar el código QR")
	ErrorGeneratingPDF         = errors.New("error al generar la boleta en formato pdf")
	ErrorInvalidPaymentType    = errors.New("el tipo de pago es inválido, ingrese: efectivo, pago por aplicación o pago con tarjeta")
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

func WriteSuccessAppointmentDesc(c echo.Context, message string, status int, originalPrice, packageDiscount, insuranceDiscount, finalPrice float64, hasInsurance bool) error {
	return c.JSON(status, map[string]interface{}{
		"El descuento por paquete es de: $/.": fmt.Sprintf("%.2f", packageDiscount),
		"El descuento por seguro es de: $/.":  fmt.Sprintf("%.2f", insuranceDiscount),
		"El precio de la cita es: $/.":        fmt.Sprintf("%.2f", originalPrice),
		"El precio final de la cita es: $/.":  fmt.Sprintf("%.2f", finalPrice),
		"tiene seguro":                        hasInsurance,
		"message":                             message,
		"status":                              status,
	})
}

func WriteSuccessAppointment(c echo.Context, message string, status int, originalPrice float64) error {
	return c.JSON(status, map[string]interface{}{
		"status":                       status,
		"message":                      message,
		"El precio de la cita es: $/.": originalPrice,
	})
}

func WriteSuccessPayment(c echo.Context, message string, status int, paymentResponse model.PaymentResponse) error {
	return c.JSON(status, map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    paymentResponse,
	})
}
