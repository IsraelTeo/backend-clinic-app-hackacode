package validate

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/db"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/unicode/norm"
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

func ParseTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func IsDateInPast(date time.Time) bool {
	now := time.Now()
	return date.Before(now)
}

func IsStartBeforeEnd(start, end time.Time) bool {
	return start.Before(end)
}

func IsWithinTimeRange(startTime, endTime, rangeStart, rangeEnd time.Time) bool {
	return !(startTime.Before(rangeStart) || endTime.After(rangeEnd))
}

func HasTimeConflict(existingAppointments []model.Appointment, newStartTime, newEndTime time.Time) bool {

	for _, appointment := range existingAppointments {
		parsedEndTime, _ := ParseTime(appointment.EndTime)
		parsedStartTime, _ := ParseTime(appointment.StartTime)
		// Verificar si el horario de la nueva cita tiene conflicto con una cita existente
		if newStartTime.Before(parsedEndTime) && newEndTime.After(parsedStartTime) {
			return true // Hay un conflicto
		}
	}
	return false // No hay conflicto
}

// Función para traducir días en español a inglés
var dayTranslations = map[string]string{
	"Lunes":     "Monday",
	"Martes":    "Tuesday",
	"Miércoles": "Wednesday",
	"Jueves":    "Thursday",
	"Viernes":   "Friday",
	"Sábado":    "Saturday",
	"Domingo":   "Sunday",
}

// Función para traducir días en inglés a español
func TranslateDayToSpanish(englishDay string) string {
	daysTranslation := map[string]string{
		"Sunday":    "domingo",
		"Monday":    "lunes",
		"Tuesday":   "martes",
		"Wednesday": "miércoles",
		"Thursday":  "jueves",
		"Friday":    "viernes",
		"Saturday":  "sábado",
	}
	return daysTranslation[englishDay]
}

// Función para traducir días en español a inglés
func TranslateDay(day string) string {
	if translated, exists := dayTranslations[day]; exists {
		return translated
	}
	return day // Si no está en el mapa, devolver tal cual
}

func IsDayAvailable(appointmentDay string, validDays []string) bool {
	// Convertimos el día de la cita a minúsculas y lo normalizamos
	appointmentDay = normalizeDay(appointmentDay)

	// Verificamos si el día de la cita está en la lista de días disponibles
	for _, validDay := range validDays {
		if validDay == appointmentDay {
			return true
		}
	}
	return false
}

// Normaliza las cadenas eliminando tildes (acentos) y pasando a minúsculas.
func normalizeDay(day string) string {
	// Convertir a minúsculas
	day = strings.ToLower(day)
	// Eliminar acentos/tildes
	return removeAccents(day)
}

// Elimina los acentos/tildes de una cadena
func removeAccents(input string) string {
	// Normalizar a forma NFD (Descompuesta), y luego eliminar caracteres diacríticos
	output := norm.NFD.String(input)
	result := []rune{}
	for _, r := range output {
		if unicode.Is(unicode.Mn, r) {
			// Eliminar los caracteres que son marcas de acento
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func CheckEmailExists(email string) bool {
	var user model.User
	if err := db.GDB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}

func CheckDNIExists(DNI string) bool {
	var patient model.Patient
	if err := db.GDB.Where("DNI = ?", DNI).First(&patient).Error; err != nil {
		return false
	}

	return true
}
