package timeutil

import (
	"testing"
)

func TestParseISO8601(t *testing.T) {
	tests := []struct {
		name     string
		dateTime string
		wantErr  bool
	}{
		{name: "Correct Date (long timezone)", dateTime: "2020-11-02T20:04:05-0700", wantErr: false},
		{name: "Correct Date (short timezone)", dateTime: "2016-12-02T20:04:05-07", wantErr: false},
		{name: "Correct Date (timezone with dot)", dateTime: "2021-11-30T15:04:05+08:00", wantErr: false},
		{name: "Correct Date (UTC)", dateTime: "2021-02-07T21:39:31Z", wantErr: false},
		{name: "Correct Date (no timezone)", dateTime: "2021-11-30T15:04:05", wantErr: false},
		{name: "Correct Date (no time)", dateTime: "2021-11-30", wantErr: false},
		{name: "Incorrect Date (date)", dateTime: "2021/11/30T15:04:05", wantErr: true},
		{name: "Incorrect Date (time)", dateTime: "2021-11-30T15-04:05", wantErr: true},
		{name: "Incorrect Date (timezone)", dateTime: "2021-11-30T15:04:05-8", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseISO8601(tt.dateTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
