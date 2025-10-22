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

package unstructured

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var ErrNotFound = errors.New("not found")

var ErrNilObject = errors.New("nil object")

// ToUnstructured returns an unstructured map holding the public field values
// from the original input obj value
func ToUnstructured(obj any) (map[string]any, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object into JSON: %w", err)
	}
	result := map[string]any{}
	if err := json.Unmarshal(js, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object JSON onto a map: %w", err)
	}
	return result, nil
}

// FromUnstructured fills a target value with the field values from an
// unstructured map
func FromUnstructured[T any](target *T, source map[string]any) error {
	js, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal map into JSON: %w", err)
	}
	if err := json.Unmarshal(js, target); err != nil {
		return fmt.Errorf("failed to unmarshal map JSON onto object: %w", err)
	}
	return nil
}

// GetField gets the value of a named path within the given unstructured map
func GetField[T any](obj map[string]any, fields ...string) (T, error) {
	var zeroValue T
	if obj == nil {
		return zeroValue, ErrNilObject
	}
	rawValue, ok, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if !ok {
		return zeroValue, fmt.Errorf("path %v %w", fields, ErrNotFound)
	}
	if err != nil {
		return zeroValue, fmt.Errorf("failed to access field path %v: %w", fields, err)
	}
	value, ok := (rawValue).(T)
	if !ok {
		return zeroValue, fmt.Errorf("field path %v expected %T but was %T", fields, zeroValue, rawValue)
	}
	return value, nil
}

// CreateField creates a field with a value at the given path, it does create
// the parent path objects as needed
func CreateField[T any](obj map[string]any, value T, fields ...string) error {
	current := obj
	path := []string{}
	for i := 0; i < len(fields)-1; i++ {
		path = append(path, fields[i])
		if rawNext, exists := current[fields[i]]; exists {
			next, typeOk := rawNext.(map[string]any)
			if !typeOk {
				return fmt.Errorf("intermediate path %v exists but is of type %T", path, rawNext)
			}
			current = next
		} else {
			next := map[string]any{}
			current[fields[i]] = next
			current = next
		}
	}
	lastField := fields[len(fields)-1]
	path = append(path, lastField)
	if previousValue, exists := current[lastField]; exists {
		return fmt.Errorf("path %v is already set to value %v", path, previousValue)
	}
	current[lastField] = value
	return nil
}

// GetOrCreateField access a field at the given path, it creates the
// field with the given defaultValue if it did not exist
func GetOrCreateField[T any](obj map[string]any, defaultValue T, fields ...string) (T, error) {
	value, err := GetField[T](obj, fields...)
	if err == nil {
		return value, nil
	}
	var emptyValue T
	if errors.Is(err, ErrNotFound) {
		if err := CreateField(obj, defaultValue, fields...); err != nil {
			return emptyValue, fmt.Errorf("failed to create field at path %v: %w", fields, err)
		}
		return defaultValue, nil
	}
	return emptyValue, fmt.Errorf("failed to check for field at path %v: %w", fields, err)
}

// AsPath translates the given simplified xpath expression into a sequence of
// path entries. This ia very shallow xpath formatter, not full xpath compliant
func AsPath(xpath string) []string {
	if strings.HasPrefix(xpath, ".") {
		return AsPath(xpath[1:])
	}
	return strings.Split(xpath, ".")
}

// CopyFields copies all unstructured fields from an source to a target
func CopyFields(target, source map[string]any) {
	for field, value := range source {
		target[field] = value
	}
}

// SkipKeys returns a copy of the origin obj without the given skips keys
func SkipKeys(obj map[string]any, skips ...string) map[string]any {
	result := map[string]any{}
	for field, value := range obj {
		if in(skips, field) {
			continue
		}
		result[field] = value
	}
	return result
}

// Base returns the base of the given path, namely the last name in the array
func Base(path []string) string {
	if len(path) == 0 {
		return ""
	}
	lastIndex := len(path) - 1
	return path[lastIndex]
}

// FieldsOf returns the names of the fields at the given obj value
func FieldsOf(obj map[string]any) []string {
	fields := make([]string, 0, len(obj))
	for field := range obj {
		fields = append(fields, field)
	}
	return fields
}

func in[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
