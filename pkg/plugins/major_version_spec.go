package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type MajorVersionSpecPlugin struct {
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &MajorVersionSpecPlugin{}

func NewMajorVersionPlugin(crd *apiextensions.CustomResourceDefinition) *MajorVersionSpecPlugin {
	return &MajorVersionSpecPlugin{
		crd: crd,
	}
}

func (s *MajorVersionSpecPlugin) Name() string {
	return "string_plugin"
}

func (s *MajorVersionSpecPlugin) ProcessMapping(g Generator, mapping configv1alpha1.CRDMapping, spec *openapi3.T) error {
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", s.crd.Spec.Names.Singular, mapping.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	return nil
}

func (s *MajorVersionSpecPlugin) ProcessField(g Generator, path []string, props *apiextensions.JSONSchemaProps) error {
	return nil
}
