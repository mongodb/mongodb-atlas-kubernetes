package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type Plugin interface {
	Name() string
	ProcessMapping(g Generator, mapping *configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error
	ProcessProperty(g Generator, mapping *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps
}

type Generator interface {
	ConvertProperty(schema, extensionsSchema *openapi3.SchemaRef, mapping *configv1alpha1.PropertyMapping, path ...string) *apiextensions.JSONSchemaProps
}
