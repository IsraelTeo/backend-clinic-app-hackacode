package validate

import (
	"log"
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

type CustomValidator struct {
	Validator *validator.Validate
}

func Init() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func (c *CustomValidator) Validate(i interface{}) error {
	err := c.Validator.Struct(i)
	if err != nil {
		log.Println("Validation Error:", err)
	}

	return err
}

func ParseID(c echo.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, response.ErrorInvalidID
	}

	return uint(id), nil
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
		if newStartTime.Before(parsedEndTime) && newEndTime.After(parsedStartTime) {
			return true
		}
	}

	return false
}

var dayTranslations = map[string]string{
	"Lunes":     "Monday",
	"Martes":    "Tuesday",
	"Miércoles": "Wednesday",
	"Jueves":    "Thursday",
	"Viernes":   "Friday",
	"Sábado":    "Saturday",
	"Domingo":   "Sunday",
}

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

func TranslateDay(day string) string {
	if translated, exists := dayTranslations[day]; exists {
		return translated
	}

	return day
}

func IsDayAvailable(appointmentDay string, validDays []string) bool {
	appointmentDay = normalizeDay(appointmentDay)
	for _, validDay := range validDays {
		if validDay == appointmentDay {
			return true
		}
	}

	return false
}

func normalizeDay(day string) string {
	day = strings.ToLower(day)
	return removeAccents(day)
}

func removeAccents(input string) string {
	output := norm.NFD.String(input)
	result := []rune{}
	for _, r := range output {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		result = append(result, r)
	}

	return string(result)
}

func CheckEmailExists[T any](email string, entity *T) bool {
	if err := db.GDB.Where("email = ?", email).First(&entity).Error; err != nil {
		return false
	}

	return true
}

func CheckDNIExists[T any](DNI string, entity *T) bool {
	if err := db.GDB.Where("DNI = ?", DNI).First(&entity).Error; err != nil {
		return false
	}

	return true
}

func CheckPhoneNumberExists[T any](phoneNumber string, entity *T) bool {
	if err := db.GDB.Where("phone_number = ?", phoneNumber).First(&entity).Error; err != nil {
		return false
	}

	return true
}
