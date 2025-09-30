package plugins

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"tools/openapi2crd/pkg/converter"
)

// Status plugin adds the status schema to the CRD if specified in the mapping configuration.
// It requires the base plugin to be run first.
type Status struct{}

func (p *Status) Name() string {
	return "status"
}

func (p *Status) Process(req *MappingProcessorRequest) error {
	if req.MappingConfig.StatusMapping.Schema == "" {
		return nil
	}

	statusSchema, ok := req.OpenAPISpec.Components.Schemas[req.MappingConfig.StatusMapping.Schema]
	if !ok {
		return fmt.Errorf("status schema %q not found in openapi spec", req.MappingConfig.StatusMapping.Schema)
	}

	statusProps := req.Converter.Convert(
		converter.PropertyConvertInput{
			Schema:              statusSchema,
			ExtensionsSchemaRef: openapi3.NewSchemaRef("", openapi3.NewSchema()),
			PropertyConfig:      &req.MappingConfig.StatusMapping,
			Depth:               0,
			Path:                nil,
		},
	)
	if statusProps != nil {
		statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", req.CRD.Spec.Names.Singular, req.MappingConfig.MajorVersion)
		req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[req.MappingConfig.MajorVersion] = *statusProps
	}

	return nil
}
