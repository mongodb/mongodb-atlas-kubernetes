package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type MajorVersionPlugin struct {
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &MajorVersionPlugin{}

func NewMajorVersionPlugin(crd *apiextensions.CustomResourceDefinition) *MajorVersionPlugin {
	return &MajorVersionPlugin{
		crd: crd,
	}
}

func (s *MajorVersionPlugin) Name() string {
	return "string_plugin"
}

func (s *MajorVersionPlugin) ProcessMapping(g Generator, mapping configv1alpha1.CRDMapping, spec *openapi3.T) error {
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", s.crd.Spec.Names.Singular, mapping.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	return nil
}

func (s *MajorVersionPlugin) ProcessField(g Generator, path []string, props *apiextensions.JSONSchemaProps) error {
	return nil
}
