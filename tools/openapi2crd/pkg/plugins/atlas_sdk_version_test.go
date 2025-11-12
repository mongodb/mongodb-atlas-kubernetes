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

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestAtlasSdkVersionName(t *testing.T) {
	p := &AtlasSdkVersionPlugin{}
	assert.Equal(t, "atlas_sdk_version", p.Name())
}

func TestAtlasSdkVersionPluginProcess(t *testing.T) {
	tests := map[string]struct {
		request            *ExtensionProcessorRequest
		expectedExtensions map[string]any
		expectedErr        error
	}{
		"add atlas sdk version extension to the CRD": {
			request: extensionRequest(t, map[string]configv1alpha1.OpenAPIDefinition{
				"v20250312": {
					Name:    "v20250312",
					Package: "go.mongodb.org/atlas-sdk/v20250312005/admin",
				},
			}),
			expectedExtensions: map[string]any{"x-atlas-sdk-version": "go.mongodb.org/atlas-sdk/v20250312005/admin"},
		},
		"no sdk packgae defined": {
			request:            extensionRequest(t, nil),
			expectedExtensions: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &AtlasSdkVersionPlugin{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			extensions := tt.request.ExtensionsSchema.Properties["spec"].Value.Properties["v20250312"].Value.Extensions
			assert.Equal(t, tt.expectedExtensions, extensions)
		})
	}
}

func extensionRequest(t *testing.T, apiDefinitions map[string]configv1alpha1.OpenAPIDefinition) *ExtensionProcessorRequest {
	t.Helper()

	extensionsSchema := openapi3.NewSchema()
	extensionsSchema.Properties = map[string]*openapi3.SchemaRef{
		"spec": {Value: &openapi3.Schema{Properties: map[string]*openapi3.SchemaRef{
			"v20250312": {Value: &openapi3.Schema{Properties: map[string]*openapi3.SchemaRef{}}},
		}}},
	}

	return &ExtensionProcessorRequest{
		ExtensionsSchema: extensionsSchema,
		ApiDefinitions:   apiDefinitions,
		MappingConfig: &configv1alpha1.CRDMapping{
			MajorVersion: "v20250312",
			OpenAPIRef: configv1alpha1.LocalObjectReference{
				Name: "v20250312",
			},
			ParametersMapping: configv1alpha1.PropertyMapping{
				Path: configv1alpha1.PropertyPath{
					Name: "/api/atlas/v2/groups",
					Verb: "post",
				},
			},
			EntryMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadWriteOnly: true,
				},
			},
			StatusMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadOnly:       true,
					SkipProperties: []string{"$.links"},
				},
			},
		},
	}
}
