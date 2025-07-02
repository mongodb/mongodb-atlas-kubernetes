package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type ParametersPlugin struct {
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &ParametersPlugin{}

func NewParametersPlugin(crd *apiextensions.CustomResourceDefinition) *ParametersPlugin {
	return &ParametersPlugin{
		crd: crd,
	}
}

func (s *ParametersPlugin) Name() string {
	return "parameters_plugin"
}

func (s *ParametersPlugin) ProcessMapping(g Generator, mapping configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error {
	majorVersionSpec := s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion]

	if mapping.ParametersMapping.FieldPath.Name != "" {
		var operation *openapi3.Operation

		pathItem, ok := openApiSpec.Paths[mapping.ParametersMapping.FieldPath.Name]
		if !ok {
			return fmt.Errorf("OpenAPI path %q does not exist", mapping.ParametersMapping)
		}

		switch mapping.ParametersMapping.FieldPath.Verb {
		case "post":
			operation = pathItem.Post
		case "put":
			operation = pathItem.Put
		default:
			return fmt.Errorf("verb %q unsupported", mapping.ParametersMapping.FieldPath.Verb)
		}

		for _, p := range operation.Parameters {
			switch p.Value.Name {
			case "includeCount":
			case "itemsPerPage":
			case "pageNum":
			case "envelope":
			case "pretty":
			default:
				props := g.SchemaPropsToJSONProps(p.Value.Schema, nil, openapi3.NewSchema())
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

	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = majorVersionSpec

	return nil
}

func (s *ParametersPlugin) ProcessField(g Generator, path []string, props *apiextensions.JSONSchemaProps) error {
	return nil
}
