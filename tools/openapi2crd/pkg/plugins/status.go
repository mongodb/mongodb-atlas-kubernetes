package plugins

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/processor"
)

type StatusPlugin struct{}

func (p *StatusPlugin) Name() string {
	return "status"
}

func (p *StatusPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	crd := i.CRD
	openApiSpec := i.OpenAPISpec

	if mappingConfig.StatusMapping.Schema == "" {
		return nil
	}

	statusSchema, ok := openApiSpec.Components.Schemas[mappingConfig.StatusMapping.Schema]
	if !ok {
		return fmt.Errorf("status schema %q not found in openapi spec", mappingConfig.StatusMapping.Schema)
	}

	statusProps := i.Converter.Convert(statusSchema, openapi3.NewSchemaRef("", openapi3.NewSchema()), &mappingConfig.StatusMapping, 0)
	if statusProps != nil {
		statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", crd.Spec.Names.Singular, mappingConfig.MajorVersion)
		crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[mappingConfig.MajorVersion] = *statusProps
	}

	return nil
}

func NewStatusPlugin() *StatusPlugin {
	return &StatusPlugin{}
}
