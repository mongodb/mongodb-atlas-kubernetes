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
	// Note, that the full format is intentionally left as the last one as it will be included into the error, shown
	// to the user, so let's show the full ISO8601
	parse, err = time.Parse("2006-01-02T15:04:05-0700", dateTime)
	if err == nil {
		return parse, nil
	}
	return parse, err
}
