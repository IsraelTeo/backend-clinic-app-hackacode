package validate

import (
	"time"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

func ParseTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, response.ErrorAppointmentTimeFormat
	}

	return parsedTime, nil
}

func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, response.ErrorAppointmentInvalidDateFormat
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
	return (startTime.Equal(rangeStart) || startTime.After(rangeStart)) && (endTime.Equal(rangeEnd) || endTime.Before(rangeEnd))
}
