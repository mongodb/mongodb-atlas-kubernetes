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
//

package checkerr

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckErr(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := log.Writer()
	originalFlags := log.Flags()

	log.SetOutput(&buf)
	log.SetFlags(0)

	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	tests := map[string]struct {
		msg         string
		f           funcErrs
		expectedLog string
	}{
		"no error": {
			msg: "Operation",
			f: funcErrs(func() error {
				return nil
			}),
			expectedLog: "",
		},
		"with error": {
			msg: "Operation",
			f: funcErrs(func() error {
				return assert.AnError
			}),
			expectedLog: "Operation failed: assert.AnError general error for testing\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			CheckErr(tt.msg, tt.f)
			logged := buf.String()
			buf.Reset()
			assert.Equal(t, tt.expectedLog, logged)
		})
	}
}
