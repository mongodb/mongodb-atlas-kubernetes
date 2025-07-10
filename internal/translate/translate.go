package translate

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TypeInfo struct {
	CRDVersion string
	SDKVersion string
	CRD        *apiextensionsv1.CustomResourceDefinition
}

func ToAPI[T any, S any](typeInfo *TypeInfo, target T, spec S, deps ...client.Object) error {
	specVersion := selectVersion(&typeInfo.CRD.Spec, typeInfo.CRDVersion)
	kind := typeInfo.CRD.Spec.Names.Kind
	props, err := getOpenAPIProperties(kind, specVersion)
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD schema properties: %w", err)
	}
	specProps, err := getSpecPropertiesFor(kind, props, "spec")
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD spec properties: %w", err)
	}
	if _, ok := specProps[typeInfo.SDKVersion]; !ok {
		return fmt.Errorf("failed to match the CRD spec version %q in schema", typeInfo.SDKVersion)
	}
	unstructuredSpec, err := toUnstructured(spec)
	if err != nil {
		return fmt.Errorf("failed to convert spec value to unstructured: %w", err)
	}
	specValue, err := accessField[map[string]any](unstructuredSpec, typeInfo.SDKVersion)
	if err != nil {
		return fmt.Errorf("failed to access version %q spec value: %w", typeInfo.SDKVersion, err)
	}
	if err := processMappings(typeInfo, specValue, deps...); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}
	log.Printf("%s", jsonize(specValue))
	entryFields, ok, err := unstructured.NestedMap(specValue, "entry")
	if !ok {
		return fmt.Errorf("failed to extract the CRD spec entry fields value: %w", err)
	}
	log.Printf("%s", jsonize(entryFields))
	targetUnstructured := map[string]any{}
	copyFields(targetUnstructured, skipKeys(specValue, "entry"))
	copyFields(targetUnstructured, entryFields)
	if err := toStructured(target, targetUnstructured); err != nil {
		return fmt.Errorf("failed to set structured value from unstructured: %w", err)
	}
	return nil
}

// selectVersion returns the version from the CRD spec that matches the given version string
func selectVersion(spec *apiextensionsv1.CustomResourceDefinitionSpec, version string) *apiextensionsv1.CustomResourceDefinitionVersion {
	if len(spec.Versions) == 0 {
		return nil
	}
	if version == "" {
		return &spec.Versions[0]
	}
	for _, v := range spec.Versions {
		if v.Name == version {
			return &v
		}
	}
	return nil
}

func getOpenAPIProperties(kind string, version *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]apiextensionsv1.JSONSchemaProps, error) {
	if version == nil {
		return nil, fmt.Errorf("missing version %q from %v spec", version, kind)
	}
	if version.Schema == nil {
		return nil, fmt.Errorf("missing version %q schema from %v spec", version, kind)
	}
	if version.Schema.OpenAPIV3Schema == nil {
		return nil, fmt.Errorf("missing version %q OpenAPI Schema from %v spec", version, kind)
	}
	if version.Schema.OpenAPIV3Schema.Properties == nil {
		return nil, fmt.Errorf("missing version %q OpenAPI Properties from %v spec", version, kind)
	}
	return version.Schema.OpenAPIV3Schema.Properties, nil
}

func getSpecPropertiesFor(kind string, props map[string]apiextensionsv1.JSONSchemaProps, field string) (map[string]apiextensionsv1.JSONSchemaProps, error) {
	prop, ok := props[field]
	if !ok {
		return nil, fmt.Errorf(" kind %q spec is missing field %q on", field, kind)
	}
	if prop.Type != "object" {
		return nil, fmt.Errorf("kind %q field %q expected to be object but is %v", kind, field, prop.Type)
	}
	return prop.Properties, nil
}

func toUnstructured(obj any) (map[string]any, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object to JSON: %w", err)
	}
	unstructuredObject := map[string]any{}
	if err := json.Unmarshal(js, &unstructuredObject); err != nil {
		return nil, fmt.Errorf("failed to convert object JSON to unstructured: %w", err)
	}
	return unstructuredObject, nil
}

func toStructured[T any](target T, source map[string]any) error {
	js, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to convert source to JSON: %w", err)
	}
	if err := json.Unmarshal(js, target); err != nil {
		return fmt.Errorf("failed to convert source JSON into object: %w", err)
	}
	return nil
}

func accessField[T any](obj map[string]any, fields ...string) (T, error) {
	var zeroValue T
	rawValue, ok, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if !ok {
		return zeroValue, nil
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

func asPath(path string) []string {
	if strings.HasPrefix(path, ".") {
		return asPath(path[1:])
	}
	return strings.Split(path, ".")
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

func jsonize(obj any) string {
	js, err := json.MarshalIndent(obj, "  ", "  ")
	if err != nil {
		return err.Error()
	}
	return string(js)
}

