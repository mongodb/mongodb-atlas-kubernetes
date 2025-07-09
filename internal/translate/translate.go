package translate

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	APIMAppingsAnnotation = "api-mappings"
)

type kubeMapping struct {
	gvr              string
	nameSelector     string
	propertySelector string
	support          *gvrSupport
}

type openAPIMapping struct {
	property string
	typeName string
}

type valueFilterFunc func(string) string

func base64DecodeFn(value string) string {
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type gvrSupport struct {
	kind          string
	version       string
	valueFilterFn valueFilterFunc
}

var supportedGVRs = map[string]gvrSupport{
	"secrets/v1": {
		kind:          "Secret",
		version:       "v1",
		valueFilterFn: base64DecodeFn,
	},
}

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

func processMappings(typeInfo *TypeInfo, spec map[string]any, deps ...client.Object) error {
	mappingsYML := typeInfo.CRD.ObjectMeta.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return nil
	}
	mappings := map[string]any{}
	yaml.Unmarshal([]byte(mappingsYML), mappings)
	props, err := accessField[map[string]any](mappings,
		"properties", "spec", "properties", typeInfo.SDKVersion, "properties")
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for the spec: %w", err)
	}
	return processMappingProperties([]string{}, props, spec, deps...)
}

func processMappingProperties(path []string, props, spec map[string]any, deps ...client.Object) error {
	for key, prop := range props {
		mapping, ok := (prop).(map[string]any)
		if !ok {
			continue
		}
		rawField, ok, err := unstructured.NestedFieldNoCopy(spec, key)
		if !ok {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to access %q: %w", key, err)
		}
		subPath := append(path, key)
		if arrayField, ok := (rawField).([]any); ok {
			return processMappingArray(subPath, mapping, arrayField, deps...)
		}
		subSpec, ok := (rawField).(map[string]any)
		if !ok {
			return fmt.Errorf("unsupported mapping of type %T", rawField)
		}
		if err := processMappingProperty(subPath, mapping, subSpec, deps...); err != nil {
			return fmt.Errorf("failed to process mapping %q: %w", key, err)
		}
	}
	return nil
}

func processMappingProperty(path []string, mapping, spec map[string]any, deps ...client.Object) error {
	if mapping["properties"] != nil {
		props, err := accessField[map[string]any](mapping, "properties")
		if err != nil {
			return fmt.Errorf("faild to access properties at %q: %w", path, err)
		}
		return processMappingProperties(path, props, spec, deps...)
	}
	if mapping["x-kubernetes-mapping"] != nil && mapping["x-openapi-mapping"] != nil {
		return processReference(path, mapping, spec, deps...)
	}
	return fmt.Errorf("unsupported extension at %v with fields %v", path, fieldsOf(mapping))
}

func processMappingArray(path []string, mapping map[string]any, specs []any, deps ...client.Object) error {
	items, err := accessField[map[string]any](mapping, "items", "properties")
	if err != nil {
		return fmt.Errorf("failed to access %q: %w", base(path), err)
	}
	for key, item := range items {
		spec := findByKey(specs, key)
		if spec == nil {
			continue
		}
		mapping, ok := (item).(map[string]any)
		if !ok {
			return fmt.Errorf("expected field %q at %v to be a map but was: %T", key, path, item)
		}
		subPath := append(path, key)
		if err := processMappingProperty(subPath, mapping, spec, deps...); err != nil {
			return fmt.Errorf("failed to map property from array item %q at %v: %w", key, path, err)
		}
	}
	return nil
}

func findByKey(list []any, key string) map[string]any {
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

func processReference(path []string, mapping, spec map[string]any, deps ...client.Object) error {
	reference, err := accessField[map[string]any](spec, base(path))
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}

	rawKubeMap, err := accessField[map[string]any](mapping, "x-kubernetes-mapping")
	if err != nil {
		return fmt.Errorf("failed accessing the Kubernetes mapping at %v: %w", path, err)
	}
	kubeMap, err := fromRawKubeMap(rawKubeMap)
	if err != nil {
		return fmt.Errorf("failed to setup kubernetes mapping at %v: %w", path, err)
	}

	rawOpenAPIMap, err := accessField[map[string]any](mapping, "x-openapi-mapping")
	if err != nil {
		return fmt.Errorf("failed accessing the Open API mapping at %v: %w", path, err)
	}
	openAPIMap := fromRawOpenAPIMap(rawOpenAPIMap)
	if openAPIMap.typeName != "string" {
		return fmt.Errorf("unsupported referenced value type %v", openAPIMap.typeName)
	}

	dep, err := solveReferencedDependency(path, reference, kubeMap, deps...)
	if err != nil {
		return fmt.Errorf("failed solving referenced kubernetes dependency: %w", err)
	}
	value, err := fetchReferencedValue(kubeMap, dep)
	if err != nil {
		return fmt.Errorf("failed fetching referenced value: %w", err)
	}
	property := base(asPath(openAPIMap.property))
	spec[property] = value
	return nil
}

func solveReferencedDependency(path []string, reference map[string]any, kubeMap *kubeMapping, deps ...client.Object) (map[string]any, error) {
	referenceValue, err := accessField[string](reference, asPath(kubeMap.nameSelector)...)
	if err != nil {
		return nil, fmt.Errorf("failed accessing reference value for mapping at %v: %w", path, err)
	}
	dep := findReferencedDep(deps, kubeMap, referenceValue)
	if dep == nil {
		return nil, fmt.Errorf("kubernetes dependency of type %q not found with name %q", kubeMap.gvr, referenceValue)
	}

	depUnstructured, err := toUnstructured(dep)
	if err != nil {
		return nil, fmt.Errorf("failed to translate referenced kubernetes type %q to unstructured: %w", kubeMap.gvr, err)
	}
	return depUnstructured, nil
}

func fetchReferencedValue(kubeMap *kubeMapping, dep map[string]any) (string, error) {
	propertySelectorPath := asPath(kubeMap.propertySelector)
	// TODO: remove fix up for secrets
	if kubeMap.gvr == "secrets/v1" && propertySelectorPath[0] != "data" {
		propertySelectorPath = append([]string{"data"}, propertySelectorPath...)
	}
	value, err := accessField[string](dep, propertySelectorPath...)
	if err != nil {
		return "", fmt.Errorf("failed to access referenced value at %v: %w", propertySelectorPath, err)
	}
	if kubeMap.support.valueFilterFn != nil {
		value = kubeMap.support.valueFilterFn(value)
	}
	return value, nil
}

func fromRawKubeMap(rawKubeMap map[string]any) (*kubeMapping, error) {
	gvr := stringify(rawKubeMap["gvr"])
	gvrSupport, ok := supportedGVRs[gvr]
	if !ok {
		return nil, fmt.Errorf("unsupported Group Version Resource %q", gvr)
	}
	return &kubeMapping{
		gvr:              gvr,
		nameSelector:     stringify(rawKubeMap["nameSelector"]),
		propertySelector: stringify(rawKubeMap["propertySelector"]),
		support:          &gvrSupport,
	}, nil
}

func fromRawOpenAPIMap(rawOpenAPIMap map[string]any) *openAPIMapping {
	return &openAPIMapping{
		property: stringify(rawOpenAPIMap["property"]),
		typeName: stringify(rawOpenAPIMap["type"]),
	}
}

func findReferencedDep(deps []client.Object, kubeMap *kubeMapping, name string) client.Object {
	for _, dep := range deps {
		gvk := dep.GetObjectKind().GroupVersionKind()
		if kubeMap.support.kind == gvk.Kind && kubeMap.support.version == gvk.Version &&
			dep.GetName() == name {
			return dep
		}
	}
	return nil
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

func stringify(obj any) string {
	s, ok := (obj).(string)
	if !ok {
		return fmt.Sprintf("failed to cast %v (type %T) into a string", obj, obj)
	}
	return s
}
