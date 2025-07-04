package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// NoOp is a struct that implements the Plugin interface that does nothing.
// It can be embedded in plugins that do not need to implement all methods of Plugin interface.
type NoOp struct {
	Plugin
}

func (n *NoOp) ProcessMapping(g Generator, mapping *configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error {
	return nil
}

func (n *NoOp) ProcessProperty(g Generator, mapping *configv1alpha1.FieldMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) {
	return
}
