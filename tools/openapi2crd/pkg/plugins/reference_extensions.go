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

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type ReferenceExtensions struct{}

func (r *ReferenceExtensions) Name() string {
	return "reference_metadata"
}

func (r *ReferenceExtensions) Process(req *ExtensionProcessorRequest) error {
	for _, ref := range req.MappingConfig.ParametersMapping.References {
		schemaRef, err := r.newMappingExtensions(ref)
		if err != nil {
			return err
		}

		req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Properties[ref.Name] = schemaRef
	}

	for _, ref := range req.MappingConfig.EntryMapping.References {
		schemaRef, err := r.newMappingExtensions(ref)
		if err != nil {
			return err
		}

		req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Properties["entry"].Value.Properties[ref.Name] = schemaRef
	}

	return nil
}

func (r *ReferenceExtensions) newMappingExtensions(ref configv1alpha1.Reference) (*openapi3.SchemaRef, error) {
	if len(ref.Target.Properties) == 0 {
		return nil, errors.New("reference target must have at least one property defined")
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
	schemaRef := openapi3.NewSchemaRef("", schema)

	return schemaRef, nil
}
