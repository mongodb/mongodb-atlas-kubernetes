package plugins

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// MajorVersion is a plugin that adds the major version schema to the CRD.
// It requires the base plugin to be run first.
type MajorVersion struct{}

func (s *MajorVersion) Name() string {
	return "major_version"
}

func (s *MajorVersion) Process(req *MappingProcessorRequest) error {
	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", req.CRD.Spec.Names.Singular, req.MappingConfig.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}

	return nil
}
