package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type SkippedProperties struct {
	NoOp
}

func NewSkippedPropertiesPlugin() *SkippedProperties {
	return &SkippedProperties{}
}

func (s *SkippedProperties) Name() string {
	return "skipped_properties"
}

func (n *SkippedProperties) ProcessProperty(g Generator, mapping *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps {
	if isSkippedField(path, mapping) {
		return nil
	}

	if mapping == nil || len(mapping.Filters.SkipProperties) == 0 {
		return props
	}

	props.Required = filterSlice(propertySchema.Required, mapping.Filters.SkipProperties)

	return props
}

func isSkippedField(path []string, mapping *configv1alpha1.PropertyMapping) bool {
	if mapping == nil {
		return false
	}

	p := jsonPath(path)

	for _, skippedField := range mapping.Filters.SkipProperties {
		if skippedField == p {
			return true
		}
	}

	return false
}

func filterSlice(source, by []string) []string {
	filtered := []string{}
	for _, s := range source {
		if !contains(by, "$."+s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
