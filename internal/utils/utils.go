package utils

import (
	"regexp"
	"time"
)

var hoursRegex = regexp.MustCompile(`^(0?\d|1[0-9]|2[0-4]).(00|30)$`)

func Contains[T comparable](targets []T, element T) bool {
	for idx := range targets {
		if targets[idx] == element {
			return true
		}
	}

	return false
}

// ValidHour returns true if the hour is between the range [0, 23] and minutes are 0 or 30.
// In any other case, false is returned
func ValidHour(hour string) bool {
	return hoursRegex.MatchString(hour)
}

// DateToString transforms the input in a string with format yyyy-mm-dd
func DateToString(date time.Time) string {
	return date.Format(time.DateOnly)
}
