package plugins

import (
	"fmt"

	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type MajorVersionSpecPlugin struct{}

func (s *MajorVersionSpecPlugin) Name() string {
	return "major_version"
}

func (s *MajorVersionSpecPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	crd := i.CRD

	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", crd.Spec.Names.Singular, mappingConfig.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	return nil
}

func NewMajorVersionPlugin() *MajorVersionSpecPlugin {
	return &MajorVersionSpecPlugin{}
}
