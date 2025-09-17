// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package plugins

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/util/sets"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type References struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &References{}

func NewReferencesPlugin(crd *apiextensions.CustomResourceDefinition) *References {
	return &References{
		crd: crd,
	}
}

func (r *References) Name() string {
	return "references"
}

func (r *References) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	majorVersionSpec := r.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion]

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
		r.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = majorVersionSpec

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
