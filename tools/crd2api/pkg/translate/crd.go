package translate

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type CRDInfo struct {
	definition *apiextensionsv1.CustomResourceDefinition
	version    string
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
