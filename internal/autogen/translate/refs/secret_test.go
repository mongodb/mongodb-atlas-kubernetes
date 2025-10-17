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

package refs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretDecode(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      string
		expectAnError bool
	}{
		{
			name:          "valid base64 string",
			input:         "aGVsbG8td29ybGQ=",
			expected:      "hello-world",
			expectAnError: false,
		},
		{
			name:          "empty string",
			input:         "",
			expected:      "",
			expectAnError: false,
		},
		{
			name:          "invalid base64 string with illegal characters",
			input:         "not-valid-!", // "!" is not a valid base64 character
			expected:      "",
			expectAnError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Action
			result, err := secretDecode(tc.input)

			// Assertions
			if tc.expectAnError {
				assert.Error(t, err, "expected an error but got none")
				assert.Empty(t, result, "result should be empty on error")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
				assert.Equal(t, tc.expected, result, "decoded string does not match expected value")
			}
		})
	}
}

func TestEncodeDecodeSymmetry(t *testing.T) {
	originalStrings := []string{
		"my-secret-password",
		"12345",
		"!@#$%^&*()_+",
		"another long string with spaces and punctuation.",
	}

	for _, original := range originalStrings {
		t.Run(original, func(t *testing.T) {
			encoded := secretEncode(original)
			decoded, err := secretDecode(encoded)

			require.NoError(t, err, "decoding should not fail for a string we just encoded")
			assert.Equal(t, original, decoded, "decoded string must match the original")
		})
	}
}
