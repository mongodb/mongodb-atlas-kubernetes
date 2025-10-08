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
	"errors"
	"fmt"
	"reflect"

	"github.com/stretchr/testify/assert/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/v3/translate/unstructured"
)

const (
	APIMAppingsAnnotation = "api-mappings"

	SecretProperySelector = "$.data.#"
)

type mapContext struct {
	main  client.Object
	m     map[client.ObjectKey]client.Object
	added []client.Object
}

func newMapContext(main client.Object, objs ...client.Object) *mapContext {
	m := map[client.ObjectKey]client.Object{}
	for _, obj := range objs {
		m[client.ObjectKeyFromObject(obj)] = obj
	}
	return &mapContext{main: main, m: m}
}

func (mc *mapContext) find(name string) client.Object {
	key := client.ObjectKey{Name: name, Namespace: mc.main.GetNamespace()}
	return mc.m[key]
}

func (mc *mapContext) has(name string) bool {
	return mc.find(name) != nil
}

func (mc *mapContext) add(obj client.Object) {
	mc.m[client.ObjectKeyFromObject(obj)] = obj
	mc.added = append(mc.added, obj)
}

type mapper struct {
	*mapContext
	expand bool
}

func newExpanderMapper(main client.Object, objs ...client.Object) *mapper {
	return newMapper(true, main, objs...)
}

func newCollarserMapper(main client.Object, objs ...client.Object) *mapper {
	return newMapper(false, main, objs...)
}

func newMapper(expand bool, main client.Object, objs ...client.Object) *mapper {
	return &mapper{
		mapContext: newMapContext(main, objs...),
		expand:     expand,
	}
}

func ExpandMappings(t *Translator, obj map[string]any, main client.Object, objs ...client.Object) ([]client.Object, error) {
	em := newExpanderMapper(main, objs...)
	mappingsYML := t.crd.definition.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return []client.Object{}, nil
	}
	mappings := map[string]any{}
	if err := yaml.Unmarshal([]byte(mappingsYML), mappings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mappings YAML: %w", err)
	}

	if err := em.expandMappingsAt(obj, mappings, "spec", t.majorVersion); err != nil {
		return nil, fmt.Errorf("failed to map properties of spec from API to Kubernetes: %w", err)
	}
	if err := em.expandMappingsAt(obj, mappings, "spec", t.majorVersion, "entry"); err != nil {
		return nil, fmt.Errorf("failed to map properties of spec from API to Kubernetes: %w", err)
	}
	if err := em.expandMappingsAt(obj, mappings, "status", t.majorVersion); err != nil {
		return nil, fmt.Errorf("failed to map properties of status from API to Kubernetes: %w", err)
	}
	return em.added, nil
}

func CollapseMappings(t *Translator, spec map[string]any, main client.Object, objs ...client.Object) error {
	cm := newCollarserMapper(main, objs...)
	mappingsYML := t.crd.definition.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return nil
	}
	mappings := map[string]any{}
	if err := yaml.Unmarshal([]byte(mappingsYML), mappings); err != nil {
		return fmt.Errorf("failed to unmarshal mappings YAML: %w", err)
	}
	props, err := unstructured.AccessField[map[string]any](mappings,
		"properties", "spec", "properties", t.majorVersion, "properties")
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for the spec: %w", err)
	}
	return cm.mapProperties([]string{}, props, spec)
}

func findEntryPathInTarget(targetType reflect.Type) []string {
	if targetType.String() == "admin.CreateAlertConfigurationApiParams" {
		return []string{"GroupAlertsConfig"}
	}
	return []string{}
}

func (m *mapper) expandMappingsAt(obj, mappings map[string]any, fields ...string) error {
	expandedPath := []string{"properties"}
	for _, field := range fields {
		expandedPath = append(expandedPath, field, "properties")
	}
	props, err := unstructured.AccessField[map[string]any](mappings, expandedPath...)
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for %v: %w", expandedPath, err)
	}
	field, err := unstructured.AccessField[map[string]any](obj, fields...)
	if err != nil {
		return fmt.Errorf("failed to access object's %v: %w", fields, err)
	}
	if err := m.mapProperties([]string{}, props, field); err != nil {
		return fmt.Errorf("failed to process properties from API into %v: %w", fields, err)
	}
	return nil
}

func (m *mapper) mapProperties(path []string, props, obj map[string]any) error {
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		subPath := append(path, key)
		if isReference(mapping) {
			if err := m.mapReference(subPath, key, mapping, obj); err != nil {
				return fmt.Errorf("failed to process reference: %w", err)
			}
			continue
		}
		rawField, err := unstructured.AccessField[any](obj, key) // unstructured.NestedFieldNoCopy(obj, key)
		if errors.Is(err, unstructured.ErrNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to access %q: %w", key, err)
		}
		if arrayField, ok := (rawField).([]any); ok {
			if err := m.mapArray(subPath, mapping, arrayField); err != nil {
				return fmt.Errorf("failed to process array mapping %q: %w", key, err)
			}
			continue
		}
		subSpec, ok := (rawField).(map[string]any)
		if !ok {
			return fmt.Errorf("unsupported mapping of type %T", rawField)
		}
		if err := m.mapObject(subPath, key, mapping, subSpec); err != nil {
			return fmt.Errorf("failed to process object mapping %q: %w", key, err)
		}
	}
	return nil
}

func (m *mapper) mapArray(path []string, mapping map[string]any, list []any) error {
	mapItems, err := unstructured.AccessField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", unstructured.Base(path), err)
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

func (m *mapper) mapObject(path []string, mapName string, mapping, obj map[string]any) error {
	if mapping["properties"] != nil {
		props, err := unstructured.AccessField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("failed to access properties at %q: %w", path, err)
		}
		return m.mapProperties(path, props, obj)
	}
	if isReference(mapping) {
		return m.mapReference(path, mapName, mapping, obj)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, unstructured.FieldsOf(mapping))
}

func (m *mapper) mapReference(path []string, mappingName string, mapping, obj map[string]any) error {
	rm := refMapping{}
	if err := unstructured.FromUnstructured(&rm, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}
	ref := newRef(mappingName, &rm)
	if m.expand {
		return ref.Expand(m.mapContext, path, obj)
	}
	return ref.Collapse(m.mapContext, path, obj)
}

func entryMatchingMapping(mapName string, mapping map[string]any, list []any, expand bool) (string, map[string]any) {
	key := mapName
	if expand {
		refMap := refMapping{}
		if err := unstructured.FromUnstructured(&refMap, mapping); err != nil {
			return "", nil // not a ref, cannot reverse mapping dfrom API property name
		}
		path := resolveXPath(refMap.XOpenAPIMapping.Property)
		key = unstructured.Base(path)
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
