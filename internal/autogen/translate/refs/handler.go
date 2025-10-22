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

package refs

import (
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/unstructured"
)

const (
	APIMAppingsAnnotation = "api-mappings"

	SecretProperySelector = "$.data.#"
)

// Handler holds the context needed to expand or collapse the references on an
// Kubernetes object translaion to and from API data
type Handler struct {
	*context
	expand bool
}

// NewHandler creates the basic context to expand or collapse references, such
// context requires the Kubernetes object being translated and its dependencies
func NewHandler(main client.Object, deps []client.Object) *Handler {
	return &Handler{context: newMapContext(main, deps)}
}

// ExpandReferences uses the handler context to expand references on a given
// unstructured value matching the main object being translated, using the
// given reference mappings and acting on a particular path og the value object
func (h *Handler) ExpandReferences(obj, mappings map[string]any, path ...string) error {
	h.expand = true

	props, err := accessMappingPropsAt(mappings, path)
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to access mappings to expand references: %w", err)
	}

	field, err := unstructured.GetField[map[string]any](obj, path...)
	if err != nil {
		return fmt.Errorf("failed to access object's %v: %w", path, err)
	}
	if err := h.scanProperties([]string{}, props, field); err != nil {
		return fmt.Errorf("failed to expand references at %v: %w", path, err)
	}
	return nil
}

// CollapseReferences uses the handler context to collapse references on a given
// unstructured value matching the main object being translated, using the
// given reference mappings and acting on a particular path of the value object
func (h *Handler) CollapseReferences(obj, mappings map[string]any, path ...string) error {
	h.expand = false

	props, err := accessMappingPropsAt(mappings, path)
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to access mappings to collapse references: %w", err)
	}

	field, err := unstructured.GetField[map[string]any](obj, path...)
	if err != nil {
		return fmt.Errorf("failed to access object's %v: %w", path, err)
	}

	if err := h.scanProperties([]string{}, props, field); err != nil {
		return fmt.Errorf("failed to collapse references at %v: %w", path, err)
	}
	return nil
}

// Added returns any kubernetes objects created as references by ExpandReferences
func (h *Handler) Added() []client.Object {
	return h.added
}

// accessMappingPropsAt reads the mappings object at a given path
func accessMappingPropsAt(mappings map[string]any, path []string) (map[string]any, error) {
	expandedPath := []string{"properties"}
	for _, field := range path {
		expandedPath = append(expandedPath, field, "properties")
	}
	props, err := unstructured.GetField[map[string]any](mappings, expandedPath...)
	if err != nil {
		return nil, fmt.Errorf("failed to access the API mapping properties for %v: %w", expandedPath, err)
	}
	return props, nil
}

// scanProperties checks an object value path position holding field properties
// against reference mappings that may apply at that path
func (m *Handler) scanProperties(path []string, props, obj map[string]any) error {
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		subPath := append(path, key)
		if isReference(mapping) {
			if err := m.processReference(subPath, key, mapping, obj); err != nil {
				return fmt.Errorf("failed to process reference: %w", err)
			}
			continue
		}
		rawField, err := unstructured.GetField[any](obj, key)
		if errors.Is(err, unstructured.ErrNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to access %q: %w", key, err)
		}
		if arrayField, ok := (rawField).([]any); ok {
			if err := m.scanArray(subPath, mapping, arrayField); err != nil {
				return fmt.Errorf("failed to process array mapping %q: %w", key, err)
			}
			continue
		}
		subSpec, ok := (rawField).(map[string]any)
		if !ok {
			return fmt.Errorf("unsupported mapping of type %T", rawField)
		}
		if err := m.scanObject(subPath, key, mapping, subSpec); err != nil {
			return fmt.Errorf("failed to process object mapping %q: %w", key, err)
		}
	}
	return nil
}

// scanArray checks an object value path position holding an array against
// reference mappings that may apply at that path
func (m *Handler) scanArray(path []string, mapping map[string]any, list []any) error {
	mapItems, err := unstructured.GetField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", unstructured.Base(path), err)
	}
	for mapName, mapItem := range mapItems {
		mapping, ok := (mapItem).(map[string]any)
		if !ok {
			return fmt.Errorf("expected field %q at %v to be a map but was: %T", mapName, path, mapItem)
		}
		key, entry, err := entryMatchingMapping(mapName, mapping, list, m.expand)
		if err != nil {
			return fmt.Errorf("failed to match mapping within array item %q at %v: %w", mapItem, path, err)
		}
		if entry == nil {
			continue
		}
		subPath := append(path, key)
		if err := m.scanObject(subPath, mapName, mapping, entry); err != nil {
			return fmt.Errorf("failed to map property from array item %q at %v: %w", key, path, err)
		}
	}
	return nil
}

// scanObject checks an object value path position holding an object against
// reference mappings that may apply at that path
func (m *Handler) scanObject(path []string, mapName string, mapping, obj map[string]any) error {
	if mapping["properties"] != nil {
		props, err := unstructured.GetField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("failed to access properties at %q: %w", path, err)
		}
		return m.scanProperties(path, props, obj)
	}
	if isReference(mapping) {
		return m.processReference(path, mapName, mapping, obj)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, unstructured.FieldsOf(mapping))
}

// processReference kicks of a reference expansion or collapse
func (m *Handler) processReference(path []string, mappingName string, mapping, obj map[string]any) error {
	rm := refMapping{}
	if err := unstructured.FromUnstructured(&rm, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}
	ref := newRef(mappingName, &rm)
	if m.expand {
		return ref.Expand(m.context, path, obj)
	}
	return ref.Collapse(m.context, path, obj)
}

// entryMatchingMapping returns the entry name key and value from an array that
// matches the reference ebing collapsed or expanded
func entryMatchingMapping(mapName string, mapping map[string]any, list []any, expand bool) (string, map[string]any, error) {
	key := mapName
	if expand {
		refMap := refMapping{}
		if err := unstructured.FromUnstructured(&refMap, mapping); err != nil {
			return "", nil, fmt.Errorf("not a reference, cannot reverse mapping from API property name")
		}
		path := resolveXPath(refMap.XOpenAPIMapping.Property)
		key = unstructured.Base(path)
	}
	m, err := findByExistingUniqueKey(list, key)
	return key, m, err
}

// findByExistingUniqueKey returns the value of an array holding a given key
func findByExistingUniqueKey(list []any, key string) (map[string]any, error) {
	candidates := []map[string]any{}
	for _, item := range list {
		obj, ok := (item).(map[string]any)
		if !ok {
			continue
		}
		if _, ok := obj[key]; ok {
			candidates = append(candidates, obj)
		}
	}
	if len(candidates) == 1 {
		return candidates[0], nil
	}
	if len(candidates) > 1 {
		return nil, fmt.Errorf("too many matches for key %q: %v", key, candidates)
	}
	return nil, nil
}
