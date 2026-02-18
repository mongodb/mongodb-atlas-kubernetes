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

package compare

import "slices"

func IsEqualWithoutOrder[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[T]bool, len(a))
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; !ok {
			return false
		}
	}
	return true
}

func PtrValuesEqual[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func Contains[T comparable](a []T, b T) bool {
	return slices.Contains(a, b)
}
