package validate

import (
	"time"
)

type Validator interface {
	ParseTime(timeStr string) (time.Time, error)
	ParseDate(dateStr string) (time.Time, error)
	ParseTimeAndDate(timeStr, dateStr string) (time.Time, time.Time, error)
}

type defaultValidator struct{}

func (v *defaultValidator) ParseTime(timeStr string) (time.Time, error) {
	timeParsed, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return timeParsed, nil
}

func (v *defaultValidator) ParseDate(dateStr string) (time.Time, error) {
	dateParsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return dateParsed, nil
}

func (v *defaultValidator) ParseTimeAndDate(timeStr, dateStr string) (time.Time, time.Time, error) {
	time, err := v.ParseTime(timeStr)
	if err != nil {
		return time.Time{}, time.Time{}, nil
	}

	date, err := v.ParseDate(dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, nil
	}

	return time, date, nil
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func IsDateInPast(date time.Time) bool {
	now := time.Now()
	return date.Before(now)
}

func IsStartBeforeEnd(start, end time.Time) bool {
	return start.Before(end)
}

func IsWithinTimeRange(start, end, rangeStart, rangeEnd time.Time, includeStart, includeEnd bool) bool {
	if rangeStart.After(rangeEnd) {
		return false
	}

	startOK := start.After(rangeStart) || (includeStart && start.Equal(rangeStart))
	endOK := end.Before(rangeEnd) || (includeEnd && end.Equal(rangeEnd))

	return startOK && endOK
}
