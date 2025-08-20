package plugins

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type ParametersPlugin struct{}

func (p *ParametersPlugin) Name() string {
	return "parameters"
}

func (p *ParametersPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	crd := i.CRD
	openApiSpec := i.OpenAPISpec

	majorVersionSpec := crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion]

	if mappingConfig.ParametersMapping.Path.Name != "" {
		var operation *openapi3.Operation

		pathItem := openApiSpec.Paths.Find(mappingConfig.ParametersMapping.Path.Name)
		if pathItem == nil {
			return fmt.Errorf("OpenAPI path %v does not exist", mappingConfig.ParametersMapping)
		}

		switch mappingConfig.ParametersMapping.Path.Verb {
		case "post":
			operation = pathItem.Post
		case "put":
			operation = pathItem.Put
		case "patch":
			operation = pathItem.Patch
		default:
			return fmt.Errorf("verb %q unsupported", mappingConfig.ParametersMapping.Path.Verb)
		}

		for _, param := range operation.Parameters {
			switch param.Value.Name {
			case "includeCount":
			case "itemsPerPage":
			case "pageNum":
			case "envelope":
			case "pretty":
			default:
				props := i.Converter.Convert(param.Value.Schema, openapi3.NewSchemaRef("", openapi3.NewSchema()), &mappingConfig.ParametersMapping, 0, "$", param.Value.Name)
				props.Description = param.Value.Description
				props.XValidations = apiextensions.ValidationRules{
					{
						Rule:    "self == oldSelf",
						Message: fmt.Sprintf("%s cannot be modified after creation", param.Value.Name),
					},
				}
				majorVersionSpec.Properties[param.Value.Name] = *props
				majorVersionSpec.Required = append(majorVersionSpec.Required, param.Value.Name)
			}
		}
	}

	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = majorVersionSpec

	return nil
}

func NewParametersPlugin() *ParametersPlugin {
	return &ParametersPlugin{}
}
