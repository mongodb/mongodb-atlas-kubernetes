package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type Plugin interface {
	Name() string
	ProcessCRD(g Generator, crdConfig *configv1alpha1.CRDConfig) error
	ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error
	ProcessProperty(g Generator, propertyConfig *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps
}

type Generator interface {
	ConvertProperty(schema, extensionsSchema *openapi3.SchemaRef, propertyConfig *configv1alpha1.PropertyMapping, path ...string) *apiextensions.JSONSchemaProps
}
