package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type StatusPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &StatusPlugin{}

func NewStatusPlugin(crd *apiextensions.CustomResourceDefinition) *StatusPlugin {
	return &StatusPlugin{
		crd: crd,
	}
}

func (s *StatusPlugin) Name() string {
	return "status"
}

func (s *StatusPlugin) ProcessMapping(g Generator, mapping *configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error {
	if mapping.StatusMapping.Schema == "" {
		return nil
	}

	statusSchema, ok := openApiSpec.Components.Schemas[mapping.StatusMapping.Schema]
	if !ok {
		return fmt.Errorf("status schema %q not found in openapi spec", mapping.StatusMapping.Schema)
	}

	statusProps := g.ConvertProperty(statusSchema, openapi3.NewSchemaRef("", openapi3.NewSchema()), &mapping.StatusMapping)
	statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", s.crd.Spec.Names.Singular, mapping.MajorVersion)
	if statusProps != nil {
		s.crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[mapping.MajorVersion] = *statusProps
	}

	return nil
}
