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

package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyWithSkip(t *testing.T) {
	t.Run("should return the collection without the skip item", func(t *testing.T) {
		c := []string{"a", "b", "c", "d", "e"}
		assert.Equal(t, []string{"a", "b", "d", "e"}, CopyWithSkip(c, "c"))
	})

	t.Run("should return the same collection when the skip item is not present", func(t *testing.T) {
		c := []string{"a", "b", "c", "d", "e"}
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, CopyWithSkip(c, "f"))
	})
}

func TestMapDiff(t *testing.T) {
	tests := map[string]struct {
		a        map[string]int
		b        map[string]int
		expected map[string]int
	}{
		"Disjoint maps": {
			a:        map[string]int{"a": 1, "b": 2},
			b:        map[string]int{"c": 3, "d": 4},
			expected: map[string]int{"a": 1, "b": 2},
		},
		"Partially overlapping maps": {
			a:        map[string]int{"a": 1, "b": 2, "c": 3},
			b:        map[string]int{"b": 2, "d": 4},
			expected: map[string]int{"a": 1, "c": 3},
		},
		"Fully overlapping maps": {
			a:        map[string]int{"a": 1, "b": 2},
			b:        map[string]int{"a": 1, "b": 2},
			expected: map[string]int{},
		},
		"Empty map a": {
			a:        map[string]int{},
			b:        map[string]int{"a": 1},
			expected: map[string]int{},
		},
		"Empty map b": {
			a:        map[string]int{"a": 1, "b": 2},
			b:        map[string]int{},
			expected: map[string]int{"a": 1, "b": 2},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, MapDiff(tt.a, tt.b))
		})
	}
}
