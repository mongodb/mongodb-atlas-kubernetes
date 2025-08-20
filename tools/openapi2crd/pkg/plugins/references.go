package plugins

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/util/sets"
)

type References struct{}

func (r *References) Name() string {
	return "references"
}

func (r *References) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	crd := i.CRD
	extensionsSchema := i.ExtensionsSchema

	majorVersionSpec := crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion]

	for _, ref := range mappingConfig.ParametersMapping.References {
		var refProp apiextensions.JSONSchemaProps

		openApiPropertyPath := strings.Split(ref.Property, ".")
		openApiProperty := openApiPropertyPath[len(openApiPropertyPath)-1]
		refProp.Type = "object"

		switch len(ref.Target.Properties) {
		case 0:
			return errors.New("reference target must have at least one property defined")
		case 1:
			refProp.Description = fmt.Sprintf("A reference to a %q resource.\nThe value of %q will be used to set %q.\nMutually exclusive with the %q property.", ref.Target.Type.Kind, ref.Target.Properties[0], openApiProperty, openApiProperty)
		default:
			bulleted := "- " + strings.Join(ref.Target.Properties, "\n- ")
			refProp.Description = fmt.Sprintf("A reference to a %q resource.\nOne of the following mutually exclusive values will be used to retrieve the %q value:\n\n%s\n\nMutually exclusive with the %q property.", ref.Target.Type.Kind, openApiProperty, bulleted, openApiProperty)
		}

		refProp.Properties = map[string]apiextensions.JSONSchemaProps{
			"name": {
				Type:        "string",
				Description: fmt.Sprintf(`Name of the %q resource.`, ref.Target.Type.Kind),
			},
		}

		required := sets.New(majorVersionSpec.Required...)
		required.Delete(openApiProperty)
		majorVersionSpec.Required = required.UnsortedList()
		slices.Sort(majorVersionSpec.Required)

		majorVersionSpec.Properties[ref.Name] = refProp
		crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = majorVersionSpec

		schema := openapi3.NewSchema()
		schema.Extensions = map[string]interface{}{}
		schema.Extensions["x-kubernetes-mapping"] = map[string]interface{}{
			"type":         map[string]interface{}{"kind": ref.Target.Type.Kind, "group": ref.Target.Type.Group, "version": ref.Target.Type.Version, "resource": ref.Target.Type.Resource},
			"nameSelector": ".name",
			"properties":   ref.Target.Properties,
		}

		schema.Extensions["x-openapi-mapping"] = map[string]interface{}{
			"property": ref.Property,
		}

		extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Properties[ref.Name] = openapi3.NewSchemaRef("", schema)
	}

	return nil
}

func NewReferencesPlugin() *References {
	return &References{}
}
