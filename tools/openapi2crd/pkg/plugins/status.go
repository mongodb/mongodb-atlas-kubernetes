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
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
)

// Status plugin adds the status schema to the CRD if specified in the mapping configuration.
// It requires the base plugin to be run first.
type Status struct{}

func (p *Status) Name() string {
	return "status"
}

func (p *Status) Process(req *MappingProcessorRequest) error {
	if req.MappingConfig.StatusMapping.Schema == "" {
		return nil
	}

	statusSchema, ok := req.OpenAPISpec.Components.Schemas[req.MappingConfig.StatusMapping.Schema]
	if !ok {
		return fmt.Errorf("status schema %q not found in openapi spec", req.MappingConfig.StatusMapping.Schema)
	}

	statusProps := req.Converter.Convert(
		converter.PropertyConvertInput{
			Schema:              statusSchema,
			ExtensionsSchemaRef: openapi3.NewSchemaRef("", openapi3.NewSchema()),
			PropertyConfig:      &req.MappingConfig.StatusMapping,
			Depth:               0,
			Path:                nil,
		},
	)
	if statusProps != nil {
		statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", req.CRD.Spec.Names.Singular, req.MappingConfig.MajorVersion)
		req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[req.MappingConfig.MajorVersion] = *statusProps
	}

	return nil
}
