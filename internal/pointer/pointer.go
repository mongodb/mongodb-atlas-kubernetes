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

package pointer

func GetOrNilIfEmpty[T any](val []T) *[]T {
	if len(val) == 0 {
		return nil
	}
	return &val
}

// SetOrNil returns the address of the given value or nil if it equals defaultValue
func SetOrNil[T comparable](val T, defaultValue T) *T {
	if val == defaultValue {
		return nil
	}
	return &val
}

// GetOrDefault returns the value of a pointer or a default value
func GetOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// GetOrPointerToDefault returns the value of a pointer or a pointer to default value
func GetOrPointerToDefault[T any](ptr *T, defaultValue T) *T {
	if ptr != nil {
		return ptr
	}
	return &defaultValue
}

// NonZeroOrDefault returns the address of the given value or defaultValue when it has zero value
func NonZeroOrDefault[T comparable](val T, defaultValue T) *T {
	if val == *new(T) {
		return &defaultValue
	}
	return &val
}

// MakePtr returns a pointer to the given value
func MakePtr[T any](value T) *T {
	return &value
}

// MakePtrOrNil returns a pointer only when value is not empty.
// Otherwise Atlas versioned API interprets a pointer to an empty value as not empty.
func MakePtrOrNil[T comparable](value T) *T {
	if value == *new(T) {
		return nil
	}
	return &value
}
