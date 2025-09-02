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

package atlasdatafederation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCredentialsToConnectionURL(t *testing.T) {
	testCases := []struct {
		name          string
		connURL       string
		userName      string
		password      string
		expected      string
		expectedError string
	}{
		{
			name:          "Valid Connection URL with Credentials",
			connURL:       "mongodb://mongodb0.example.com:27017/?ssl=true",
			userName:      "testUser",
			password:      "testPassword",
			expected:      "mongodb://testUser:testPassword@mongodb0.example.com:27017/?ssl=true",
			expectedError: "",
		},
		{
			name:          "Malformed URL",
			connURL:       "://invalid",
			userName:      "user",
			password:      "pass",
			expected:      "",
			expectedError: "missing protocol scheme",
		},
		{
			name:          "Special Characters in Credentials",
			connURL:       "mongodb://mongodb.example.com:27017/?ssl=true",
			userName:      "user@name",
			password:      "pass#word",
			expected:      "mongodb://user%40name:pass%23word@mongodb.example.com:27017/?ssl=true",
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := AddCredentialsToConnectionURL(tc.connURL, tc.userName, tc.password)

			if tc.expectedError != "" {
				assert.Error(t, err, "Expected an error for test case: %s", tc.name)
				assert.Contains(t, err.Error(), tc.expectedError, "Error message mismatch for test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for test case: %s", tc.name)
				assert.Equal(t, tc.expected, actual, "Unexpected output for test case: %s", tc.name)
			}
		})
	}
}
