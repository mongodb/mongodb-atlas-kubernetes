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
)

var (
	// ErrNotFound when a field is not found in an object
	ErrNotFound = errors.New("not found")

	// ErrNilObject when an object is unexpectedly nil
	ErrNilObject = errors.New("nil object")

	// ErrNotObject when a field is not an object
	ErrNotObject = errors.New("not an object")

	// ErrNotArray when a field is not an array
	ErrNotArray = errors.New("not an array")
)

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
