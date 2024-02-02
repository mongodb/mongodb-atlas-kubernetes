package stringutil

import "time"

// Contains returns true if there is at least one string in `slice`
// that is equal to `s`.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// StringToTime parses the given string and returns the resulting time.
// The expected format is identical to the format returned by Atlas API, documented as ISO 8601 timestamp format in UTC.
// Example formats: "2023-07-18T16:12:23Z", "2023-07-18T16:12:23.456Z"
func StringToTime(val string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, val)
}
