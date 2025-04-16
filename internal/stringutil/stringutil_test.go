// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
