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

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/util/sets"
)

// References adds reference properties to the CRD OpenAPI schema based on the mapping configuration.
// It requires base and major version schemas to be already processed.
type References struct{}

func (r *References) Name() string {
	return "reference"
}

func (r *References) Process(req *MappingProcessorRequest) error {
	majorVersionSpec := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion]

	for _, ref := range req.MappingConfig.ParametersMapping.References {
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
		req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion] = majorVersionSpec
	}

	return nil
}
