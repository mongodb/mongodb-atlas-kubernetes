package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

// NoOp is a struct that implements the Plugin interface that does nothing.
// It can be embedded in plugins that do not need to implement all methods of Plugin interface.
type NoOp struct {
	Plugin
}

func (n *NoOp) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	return nil
}

func (n *NoOp) ProcessProperty(g Generator, propertyConfig *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps {
	return props
}

func (n *NoOp) ProcessCRD(g Generator, crdConfig *configv1alpha1.CRDConfig) error {
	return nil
}
