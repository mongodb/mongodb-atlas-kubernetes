package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type ParametersPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &ParametersPlugin{}

func NewParametersPlugin(crd *apiextensions.CustomResourceDefinition) *ParametersPlugin {
	return &ParametersPlugin{
		crd: crd,
	}
}

func (s *ParametersPlugin) Name() string {
	return "parameters"
}

func (s *ParametersPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	majorVersionSpec := s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion]

	if mappingConfig.ParametersMapping.Path.Name != "" {
		var operation *openapi3.Operation

		pathItem := openApiSpec.Paths.Find(mappingConfig.ParametersMapping.Path.Name)
		if pathItem == nil {
			return fmt.Errorf("OpenAPI path %q does not exist", mappingConfig.ParametersMapping)
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

		for _, p := range operation.Parameters {
			switch p.Value.Name {
			case "includeCount":
			case "itemsPerPage":
			case "pageNum":
			case "envelope":
			case "pretty":
			default:
				props := g.ConvertProperty(p.Value.Schema, openapi3.NewSchemaRef("", openapi3.NewSchema()), &mappingConfig.ParametersMapping, 0, "$", p.Value.Name)
				props.Description = p.Value.Description
				props.XValidations = apiextensions.ValidationRules{
					{
						Rule:    "self == oldSelf",
						Message: fmt.Sprintf("%s cannot be modified after creation", p.Value.Name),
					},
				}
				majorVersionSpec.Properties[p.Value.Name] = *props
				majorVersionSpec.Required = append(majorVersionSpec.Required, p.Value.Name)
			}
		}
	}

	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = majorVersionSpec

	return nil
}
