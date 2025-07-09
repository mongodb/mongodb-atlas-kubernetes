package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type MajorVersionSpecPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &MajorVersionSpecPlugin{}

func NewMajorVersionPlugin(crd *apiextensions.CustomResourceDefinition) *MajorVersionSpecPlugin {
	return &MajorVersionSpecPlugin{
		crd: crd,
	}
}

func (s *MajorVersionSpecPlugin) Name() string {
	return "major_version"
}

func (s *MajorVersionSpecPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", s.crd.Spec.Names.Singular, mappingConfig.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	return nil
}
