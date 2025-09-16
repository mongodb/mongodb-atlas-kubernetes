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

package translate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	APIMAppingsAnnotation = "api-mappings"

	SecretProperySelector = "$.data.#"
)

type Mapper struct {
	deps   DependencyRepo
	expand bool
}

func (m *Mapper) mapProperties(path []string, props, obj map[string]any) error {
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		subPath := append(path, key)
		if isReference(mapping) {
			err := m.mapReference(subPath, key, mapping, obj)
			if err != nil {
				return fmt.Errorf("failed to process reference: %w", err)
			}
			continue
		}
		rawField, ok, err := unstructured.NestedFieldNoCopy(obj, key)
		if !ok {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to access %q: %w", key, err)
		}
		if arrayField, ok := (rawField).([]any); ok {
			return m.mapArray(subPath, mapping, arrayField)
		}
		subSpec, ok := (rawField).(map[string]any)
		if !ok {
			return fmt.Errorf("unsupported mapping of type %T", rawField)
		}
		if err := m.mapObject(subPath, key, mapping, subSpec); err != nil {
			return fmt.Errorf("failed to process mapping %q: %w", key, err)
		}
	}
	return nil
}

func (m *Mapper) mapArray(path []string, mapping map[string]any, list []any) error {
	mapItems, err := accessField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", base(path), err)
	}
	for mapName, mapItem := range mapItems {
		mapping, ok := (mapItem).(map[string]any)
		if !ok {
			return fmt.Errorf("expected field %q at %v to be a map but was: %T", mapName, path, mapItem)
		}
		key, entry := entryMatchingMapping(mapName, mapping, list, m.expand)
		if entry == nil {
			continue
		}
		subPath := append(path, key)
		if err := m.mapObject(subPath, mapName, mapping, entry); err != nil {
			return fmt.Errorf("failed to map property from array item %q at %v: %w", key, path, err)
		}
	}
	return nil
}

func (m *Mapper) mapObject(path []string, mapName string, mapping, obj map[string]any) error {
	if mapping["properties"] != nil {
		props, err := accessField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("failed to access properties at %q: %w", path, err)
		}
		return m.mapProperties(path, props, obj)
	}
	if isReference(mapping) {
		return m.mapReference(path, mapName, mapping, obj)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, fieldsOf(mapping))
}

func (m *Mapper) mapReference(path []string, mappingName string, mapping, obj map[string]any) error {
	rm := refMapping{}
	if err := fromUnstructured(&rm, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}
	ref := newRef(mappingName, &rm)
	if m.expand {
		return ref.Expand(m.deps, path, obj)
	}
	return ref.Collapse(m.deps, path, obj)
}

func entryMatchingMapping(mapName string, mapping map[string]any, list []any, expand bool) (string, map[string]any) {
	key := mapName
	if expand {
		refMap := refMapping{}
		if err := fromUnstructured(&refMap, mapping); err != nil {
			return "", nil // not a ref, cannot reverse mapping dfrom API property name
		}
		path := resolveXPath(refMap.XOpenAPIMapping.Property)
		key = base(path)
	}
	return key, findByExistingKey(list, key)
}

func findByExistingKey(list []any, key string) map[string]any {
	for _, item := range list {
		obj, ok := (item).(map[string]any)
		if !ok {
			continue
		}
		if _, ok := obj[key]; ok {
			return obj
		}
	}
	return nil
}
