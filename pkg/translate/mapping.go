package translate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	APIMAppingsAnnotation = "api-mappings"

	SecretProperySelector = "$.data.#"
)

func processKubeProperties(path []string, props, spec map[string]any, deps *depsBuilder) error {
	//extras := []client.Object{}
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		subPath := append(path, key)
		if isReference(mapping) {
			err := expandReference(subPath, mapping, spec, deps)
			if err != nil {
				return fmt.Errorf("failed to process reference: %w", err)
			}
			continue
		}
	}
	return nil
}

func processAPIProperties(path []string, props, spec map[string]any, deps DependencyFinder) error {
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		subPath := append(path, key)
		if isReference(mapping) {
			err := processReference(subPath, mapping, spec, deps)
			if err != nil {
				return fmt.Errorf("failed to process reference: %w", err)
			}
			continue
		}
		rawField, ok, err := unstructured.NestedFieldNoCopy(spec, key)
		if !ok {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to access %q: %w", key, err)
		}
		if arrayField, ok := (rawField).([]any); ok {
			return processArrayMapping(subPath, mapping, arrayField, deps)
		}
		subSpec, ok := (rawField).(map[string]any)
		if !ok {
			return fmt.Errorf("unsupported mapping of type %T", rawField)
		}
		if err := processObjectMapping(subPath, mapping, subSpec, deps); err != nil {
			return fmt.Errorf("failed to process mapping %q: %w", key, err)
		}
	}
	return nil
}

func processArrayMapping(path []string, mapping map[string]any, specs []any, deps DependencyFinder) error {
	items, err := accessField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", base(path), err)
	}
	for key, item := range items {
		spec := findByExistingKey(specs, key)
		if spec == nil {
			continue
		}
		mapping, ok := (item).(map[string]any)
		if !ok {
			return fmt.Errorf("expected field %q at %v to be a map but was: %T", key, path, item)
		}
		subPath := append(path, key)
		if err := processObjectMapping(subPath, mapping, spec, deps); err != nil {
			return fmt.Errorf("failed to map property from array item %q at %v: %w", key, path, err)
		}
	}
	return nil
}

func processObjectMapping(path []string, mapping, spec map[string]any, deps DependencyFinder) error {
	if mapping["properties"] != nil {
		props, err := accessField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("faild to access properties at %q: %w", path, err)
		}
		return processAPIProperties(path, props, spec, deps)
	}
	if isReference(mapping) {
		return processReference(path, mapping, spec, deps)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, fieldsOf(mapping))
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
