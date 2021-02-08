package timeutil

import "time"

func ParseISO8601(dateTime string) (time.Time, error) {
	return time.Parse(dateTime, "2006-01-02T15:04:05-0700")
}
