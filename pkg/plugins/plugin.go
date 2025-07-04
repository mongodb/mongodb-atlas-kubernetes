package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type Plugin interface {
	Name() string
	ProcessMapping(g Generator, mapping *configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error
	ProcessProperty(g Generator, mapping *configv1alpha1.FieldMapping, props *apiextensions.JSONSchemaProps, propertySchema, extensionsSchema *openapi3.Schema, path ...string)
	ProcessPropertyName(mapping *configv1alpha1.FieldMapping, path []string) string
}

type Generator interface {
	ConvertProperty(propertyName string, schemaRef *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, extensionsSchema *openapi3.Schema, path ...string) *apiextensions.JSONSchemaProps
}
