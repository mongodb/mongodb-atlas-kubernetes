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
			err := m.mapReference(subPath, mapping, obj)
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
		if err := m.mapObject(subPath, mapping, subSpec); err != nil {
			return fmt.Errorf("failed to process mapping %q: %w", key, err)
		}
	}
	return nil
}

func (m *Mapper) mapArray(path []string, mapping map[string]any, obj []any) error {
	items, err := accessField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", base(path), err)
	}
	for key, item := range items {
		spec := findByExistingKey(obj, key)
		if spec == nil {
			continue
		}
		mapping, ok := (item).(map[string]any)
		if !ok {
			return fmt.Errorf("expected field %q at %v to be a map but was: %T", key, path, item)
		}
		subPath := append(path, key)
		if err := m.mapObject(subPath, mapping, spec); err != nil {
			return fmt.Errorf("failed to map property from array item %q at %v: %w", key, path, err)
		}
	}
	return nil
}

func (m *Mapper) mapObject(path []string, mapping, obj map[string]any) error {
	if mapping["properties"] != nil {
		props, err := accessField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("faild to access properties at %q: %w", path, err)
		}
		return m.mapProperties(path, props, obj)
	}
	if isReference(mapping) {
		return m.mapReference(path, mapping, obj)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, fieldsOf(mapping))
}

func (m *Mapper) mapReference(path []string, mapping, obj map[string]any) error {
	if m.expand {
		return expandReference(m.deps, path, mapping, obj)
	}
	return collapseReference(m.deps, path, mapping, obj)
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
