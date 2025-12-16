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

package objmap

import (
	"errors"
	"fmt"
)

// GetField gets the value of a named path within the given unstructured map
func GetField[T any](obj map[string]any, fields ...string) (T, error) {
	var zeroValue T
	if obj == nil {
		return zeroValue, ErrNilObject
	}
	rawValue, _, err := nestedField(obj, fields, 0)
	if err != nil {
		return zeroValue, fmt.Errorf("failed to access field path %v: %w", fields, err)
	}
	value, ok := (rawValue).(T)
	if !ok {
		return zeroValue, fmt.Errorf("field path %v expected %T but was %T", fields, zeroValue, rawValue)
	}
	return value, nil
}

// GetFieldObject returns the object holding the given field path
func GetFieldObject(obj map[string]any, fields ...string) (map[string]any, error) {
	if obj == nil {
		return nil, ErrNilObject
	}
	field, dir := Base(fields), Dir(fields)
	value, _, err := nestedField(obj, dir, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to access field location %v: %w", dir, err)
	}
	array, isArray := (value).([]any)
	if isArray {
		for _, item := range array {
			holder, isObj := (item).(map[string]any)
			if !isObj {
				continue
			}
			if holder[field] != nil {
				return holder, nil
			}
		}
		return nil, fmt.Errorf("failed to access field %q in array: %w", field, ErrNotFound)
	}
	holder, isObj := (value).(map[string]any)
	if !isObj {
		return nil, fmt.Errorf("holder at %v is a %T: %w", dir, value, ErrNotObject)
	}
	if holder[field] == nil {
		return nil, fmt.Errorf("failed to access field %q in object: %w", field, ErrNotFound)
	}
	return holder, nil
}

// CreateField creates a field with a value at the given path
func CreateField[T any](obj map[string]any, value T, fields ...string) error {
	if len(fields) == 0 {
		return errors.New("expected one or more fields")
	}
	if fields[0] == "[]" {
		return errors.New("root must be an object not an array")
	}
	dir := Dir(fields)
	if len(dir) == 0 {
		return createObjectAt(obj, value, fields[0])
	}
	if Base(dir) == "[]" {
		return appendToArrayAt(obj, value, Dir(dir)...)
	}
	rawParent, _, err := nestedField(obj, dir, 0)
	if err != nil {
		return fmt.Errorf("failed to access parent to create field: %w", err)
	}
	parent, isObj := (rawParent).(map[string]any)
	if !isObj {
		return fmt.Errorf("parent at %v is not array nor object but a %T", dir, rawParent)
	}
	return createObjectAt(parent, value, Base(fields))
}

// RecursiveCreateField creates a field with a value at the given path, it does create
// the parent path objects as needed
func RecursiveCreateField[T any](obj map[string]any, value T, fields ...string) error {
	_, pos, err := nestedField(obj, fields, 0)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("failed to probe path %v: %w", fields, err)
	}
	for i := pos; i < len(fields)-1; i++ {
		path := fields[:i+1]
		if err := CreateField(obj, map[string]any{}, path...); err != nil {
			return fmt.Errorf("failed to create intermediate field path %v: %w", path, err)
		}
	}
	return CreateField(obj, value, fields...)
}

func appendToArrayAt[T any](obj map[string]any, value T, arrayPath ...string) error {
	parentObjectPath := Dir(arrayPath)
	objParentRaw, _, err := nestedField(obj, parentObjectPath, 0)
	if err != nil {
		return fmt.Errorf("failed to access the array's parent object: %w", err)
	}
	parentObj, isObj := (objParentRaw).(map[string]any)
	if !isObj {
		return fmt.Errorf("array parent must be object but got %T", objParentRaw)
	}
	arrayField := Base(arrayPath)
	arrayRaw, found := parentObj[arrayField]
	if !found {
		return fmt.Errorf("array not found at %q", arrayField)
	}
	array, isArray := (arrayRaw).([]any)
	if !isArray {
		return fmt.Errorf("array not found at %q", arrayField)
	}
	array = append(array, value)
	return createObjectAt(parentObj, array, arrayField)
}

func createObjectAt[T any](obj map[string]any, value T, field string) error {
	obj[field] = value
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
		if err := RecursiveCreateField(obj, defaultValue, fields...); err != nil {
			return emptyValue, fmt.Errorf("failed to create field at path %v: %w", fields, err)
		}
		return defaultValue, nil
	}
	return emptyValue, fmt.Errorf("failed to check for field at path %v: %w", fields, err)
}

// CopyFields copies all unstructured fields from an source to a target
func CopyFields(target, source map[string]any) {
	for field, value := range source {
		target[field] = value
	}
}

// FieldsOf returns the names of the fields at the given obj value
func FieldsOf(obj map[string]any) []string {
	fields := make([]string, 0, len(obj))
	for field := range obj {
		fields = append(fields, field)
	}
	return fields
}

func nestedField(value any, fields []string, pos int) (any, int, error) {
	if len(fields) == 0 {
		return value, pos, nil
	}
	if fields[pos] == "[]" {
		array, ok := (value).([]any)
		if !ok {
			return nil, pos, fmt.Errorf("field %q (type %T): %w", fields[0], value, ErrNotArray)
		}
		if pos+1 == len(fields) {
			return array, pos, nil
		}
		return nestedArrayField(array, fields, pos+1)
	}
	obj, ok := (value).(map[string]any)
	if !ok {
		return nil, pos, fmt.Errorf("field %q: %w", rootObjIfEmpty(fields, pos-1), ErrNotObject)
	}
	return nestedObjectField(obj, fields, pos)
}

func nestedObjectField(obj map[string]any, fields []string, pos int) (any, int, error) {
	field := fields[pos]
	value, ok := obj[field]
	if !ok {
		return nil, pos, fmt.Errorf("field %q: %w", field, ErrNotFound)
	}
	if pos+1 == len(fields) {
		return value, pos, nil
	}
	if value == nil {
		return nil, pos, fmt.Errorf("field %q: %w", field, ErrNilObject)
	}
	return nestedField(value, fields, pos+1)
}

func nestedArrayField(array []any, fields []string, pos int) (any, int, error) {
	for _, item := range array {
		value, branchParents, err := nestedField(item, fields, pos)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, pos, fmt.Errorf("failed searching array: %w", err)
		}
		return value, branchParents, nil
	}
	return nil, pos, ErrNotFound
}

func rootObjIfEmpty(fields []string, pos int) string {
	if pos < 0 {
		return "(obj)"
	}
	return fields[pos]
}
