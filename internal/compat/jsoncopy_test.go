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

package compat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
)

// JSONCopy will copy src to dst via JSON serialization/deserialization.
func TestJSONCopy(t *testing.T) {
	old := struct {
		Field1 string `json:"field1"`
		Field3 string `json:"field3"`
	}{
		"old field1",
		"old field3",
	}

	new := struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
	}{
		"new field1",
		"new field2",
	}

	assert := assert.New(t)
	err := JSONCopy(&old, new)
	assert.NoError(err, "JSONCopy should not fail")

	assert.Equal("new field1", old.Field1, "Field1 should be overwritten")
	assert.Equal("old field3", old.Field3, "Field3 should be left unchanged")
}
