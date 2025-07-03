package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// Just a standard implementation of Plugin interface that does nothing.
// Can be embedded in plugins that do not need to implement all methods of Plugin interface.
type NoOp struct {
	Plugin
}

func (n *NoOp) ProcessMapping(g Generator, mapping configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error {
	return nil
}

func (n *NoOp) ProcessProperty(g Generator, path []string, props *apiextensions.JSONSchemaProps) error {
	return nil
}

func (n *NoOp) ProcessPropertyName(mapping *configv1alpha1.FieldMapping, path []string) string {
	return path[len(path)-1]
}
