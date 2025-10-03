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
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestMajorVersionName(t *testing.T) {
	p := &MajorVersion{}
	assert.Equal(t, "major_version", p.Name())
}

func TestMajorVersionProcess(t *testing.T) {
	tests := map[string]struct {
		request             *MappingProcessorRequest
		expectedVersionSpec apiextensions.JSONSchemaProps
		expectedErr         error
	}{
		"add major version schema to the CRD": {
			request:             groupMappingRequest(t, groupBaseCRD(t), majorVersionInitialExtensionsSchema(t), nil),
			expectedVersionSpec: majorVersionSchema(t),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &MajorVersion{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			versionSpec := tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[tt.request.MappingConfig.MajorVersion]
			assert.Equal(t, tt.expectedVersionSpec, versionSpec)
		})
	}
}

func majorVersionSchema(t *testing.T) apiextensions.JSONSchemaProps {
	t.Helper()

	return apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250312.",
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
}

func majorVersionInitialExtensionsSchema(t *testing.T) *openapi3.Schema {
	t.Helper()

	return &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: map[string]*openapi3.SchemaRef{
			"spec": {
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{
						"v20250312": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"entry": {Value: &openapi3.Schema{}},
								},
							},
						},
					},
				},
			},
		},
	}
}
