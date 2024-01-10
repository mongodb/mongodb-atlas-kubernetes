package timeutil

import "time"

// ParseISO8601 parses string as ISO8601 date. Note that the format allows different timezones, also time is optional,
// so we have to list all possible formats
func ParseISO8601(dateTime string) (time.Time, error) {
	parse, err := time.Parse("2006-01-02T15:04:05-07", dateTime)
	if err == nil {
		return parse, nil
	}
	parse, err = time.Parse("2006-01-02T15:04:05-07:00", dateTime)
	if err == nil {
		return parse, nil
	}
	parse, err = time.Parse("2006-01-02T15:04:05", dateTime)
	if err == nil {
		return parse, nil
	}
	parse, err = time.Parse("2006-01-02", dateTime)
	if err == nil {
		return parse, nil
	}
	parse, err = time.Parse("2006-01-02T15:04:05-0700", dateTime)
	if err == nil {
		return parse, nil
	}
	// This is the default format as returned by Atlas (UTC time - marked by 'Z') - so let's show this in the error
	// if any
	parse, err = time.Parse("2006-01-02T15:04:05.999Z", dateTime)
	if err == nil {
		return parse, nil
	}
	return parse, err
}

// MustParseISO8601 returns time or panics. Mostly needed for tests.
func MustParseISO8601(dateTime string) time.Time {
	iso8601, err := ParseISO8601(dateTime)
	if err != nil {
		panic(err.Error())
	}
	return iso8601
}

// FormatISO8601 returns the ISO8601 string format for the dateTime.
func FormatISO8601(dateTime time.Time) string {
	return dateTime.Format("2006-01-02T15:04:05.999Z")
}
