package plugins

import (
	"fmt"

	"tools/openapi2crd/pkg/converter"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// Parameters adds parameters from the OpenAPI spec to the CRD schema.
// It requires base and major version plugins to be run before.
type Parameters struct{}

func (p *Parameters) Name() string {
	return "parameters"
}

func (p *Parameters) Process(req *MappingProcessorRequest) error {
	majorVersionSpec := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion]

	if req.MappingConfig.ParametersMapping.Path.Name != "" {
		var operation *openapi3.Operation

		pathItem := req.OpenAPISpec.Paths.Find(req.MappingConfig.ParametersMapping.Path.Name)
		if pathItem == nil {
			return fmt.Errorf("OpenAPI path %v does not exist", req.MappingConfig.ParametersMapping)
		}

		switch req.MappingConfig.ParametersMapping.Path.Verb {
		case "post":
			operation = pathItem.Post
		case "put":
			operation = pathItem.Put
		case "patch":
			operation = pathItem.Patch
		default:
			return fmt.Errorf("verb %q unsupported", req.MappingConfig.ParametersMapping.Path.Verb)
		}

		for _, param := range operation.Parameters {
			switch param.Value.Name {
			case "includeCount":
			case "itemsPerPage":
			case "pageNum":
			case "envelope":
			case "pretty":
			default:
				props := req.Converter.Convert(
					converter.PropertyConvertInput{
						Schema:              param.Value.Schema,
						ExtensionsSchemaRef: openapi3.NewSchemaRef("", openapi3.NewSchema()),
						PropertyConfig:      &req.MappingConfig.ParametersMapping,
						Depth:               0,
						Path:                []string{"$", param.Value.Name},
					},
				)
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

	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion] = majorVersionSpec

	return nil
}
