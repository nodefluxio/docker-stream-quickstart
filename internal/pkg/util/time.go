package util

import (
	"time"
)

// ConvertTimeNowToTzAndUtc is functin for format time to timezone and format again to UTC
func ConvertTimeNowToTzAndUtc(timezone string) (string, error) {
	t := time.Now()

	tzLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	TzTime := t.In(tzLocation)
	t1 := time.Date(TzTime.Year(), TzTime.Month(), TzTime.Day(), 0, 0, 0, TzTime.Nanosecond(), TzTime.Location())

	utcLocation, err := time.LoadLocation("UTC")
	if err != nil {
		return "", err
	}
	utcTime := t1.In(utcLocation)
	return utcTime.Format(time.RFC3339), nil
}
