package stringutil

import (
	"testing"
	"time"
)

func TestStringToTime(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string

		want    time.Time
		wantErr string
	}{
		{
			name:  "valid time",
			input: "2023-07-18T16:12:23Z",
			want: time.Date(
				2023, 7, 18,
				16, 12, 23, 0,
				time.UTC,
			),
		},
		{
			name:  "valid time with millis",
			input: "2023-07-18T16:12:23.456Z",
			want: time.Date(
				2023, 7, 18,
				16, 12, 23, 456_000_000,
				time.UTC,
			),
		},
		{
			name:    "invalid time",
			input:   "invalid",
			wantErr: `parsing time "invalid" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "invalid" as "2006"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := ""
			got, err := StringToTime(tc.input)
			if err != nil {
				gotErr = err.Error()
			}

			if gotErr != tc.wantErr {
				t.Errorf("want error %q, got %q", tc.wantErr, gotErr)
			}
			if got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
