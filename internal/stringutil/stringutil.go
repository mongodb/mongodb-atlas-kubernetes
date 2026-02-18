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

import "slices"

import "time"

// Contains returns true if there is at least one string in `slice`
// that is equal to `s`.
func Contains(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

// StringToTime parses the given string and returns the resulting time.
// The expected format is identical to the format returned by Atlas API, documented as ISO 8601 timestamp format in UTC.
// Example formats: "2023-07-18T16:12:23Z", "2023-07-18T16:12:23.456Z"
func StringToTime(val string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, val)
}
