package plugins

import (
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/util/sets"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type ReadOnlyProperties struct {
	NoOp
}

var _ Plugin = &ReadOnlyProperties{}

func NewReadOnlyPropertiesPlugin() *ReadOnlyProperties {
	return &ReadOnlyProperties{}
}

func (s *ReadOnlyProperties) Name() string {
	return "read_only_properties"
}

func (n *ReadOnlyProperties) ProcessProperty(g Generator, propertyConfig *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps {
	if propertyConfig == nil || !propertyConfig.Filters.ReadOnly {
		return props
	}

	if propertySchema.ReadOnly {
		return props
	}

	required := sets.New(propertySchema.Required...)
	for name, p := range propertySchema.Properties {
		if !p.Value.ReadOnly {
			required.Delete(name)
		}
	}
	props.Required = required.UnsortedList()
	slices.Sort(props.Required)

	// ignore root
	if len(path) == 1 {
		return props
	}

	return nil
}
