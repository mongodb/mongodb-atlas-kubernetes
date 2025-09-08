package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var ErrNotFound = errors.New("not found")

func toUnstructured(obj any) (map[string]any, error) {
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

func fromUnstructured[T any](target *T, source map[string]any) error {
	js, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal map into JSON: %w", err)
	}
	if err := json.Unmarshal(js, target); err != nil {
		return fmt.Errorf("failed to unmarshal map JSON onto object: %w", err)
	}
	return nil
}

func accessField[T any](obj map[string]any, fields ...string) (T, error) {
	var zeroValue T
	rawValue, ok, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if !ok {
		return zeroValue, fmt.Errorf("path %v %w", fields, ErrNotFound)
	}
	if err != nil {
		return zeroValue, fmt.Errorf("failed to access field path %v: %w", fields, err)
	}
	value, ok := (rawValue).(T)
	if !ok {
		return zeroValue, fmt.Errorf("field path %v is not an object map", fields)
	}
	return value, nil
}

func createField[T any](obj map[string]any, value T, fields ...string) error {
	current := obj
	path := []string{}
	for i := 0; i < len(fields)-1; i++ {
		path = append(path, fields[i])
		if rawNext, exists := current[fields[i]]; exists {
			next, typeOk := rawNext.(map[string]any)
			if !typeOk {
				return fmt.Errorf("intermediate path %v exists but is of type %T", path, next)
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

func asPath(xpath string) []string {
	if strings.HasPrefix(xpath, ".") {
		return asPath(xpath[1:])
	}
	return strings.Split(xpath, ".")
}

func base(path []string) string {
	if len(path) == 0 {
		return ""
	}
	lastIndex := len(path) - 1
	return path[lastIndex]
}

func copyFields(target, source map[string]any) {
	for field, value := range source {
		target[field] = value
	}
}

func fieldsOf(obj map[string]any) []string {
	fields := make([]string, 0, len(obj))
	for field := range obj {
		fields = append(fields, field)
	}
	return fields
}

func skipKeys(obj map[string]any, skips ...string) map[string]any {
	result := map[string]any{}
	for field, value := range obj {
		if in(skips, field) {
			continue
		}
		result[field] = value
	}
	return result
}

func in[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
