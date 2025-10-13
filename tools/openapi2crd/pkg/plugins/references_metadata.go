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

	"github.com/getkin/kin-openapi/openapi3"
)

type ReferencesMetadata struct{}

func (r *ReferencesMetadata) Name() string {
	return "reference_metadata"
}

func (r *ReferencesMetadata) Process(req *ExtensionProcessorRequest) error {
	for _, ref := range req.MappingConfig.ParametersMapping.References {
		if len(ref.Target.Properties) == 0 {
			return errors.New("reference target must have at least one property defined")
		}

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

		req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Properties[ref.Name] = openapi3.NewSchemaRef("", schema)
	}

	return nil
}
