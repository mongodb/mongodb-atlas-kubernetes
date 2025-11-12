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

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/utils/ptr"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
)

func TestStatusName(t *testing.T) {
	p := &Status{}
	assert.Equal(t, "status", p.Name())
}

func TestStatusProcess(t *testing.T) {
	tests := map[string]struct {
		request       *MappingProcessorRequest
		expectedProps apiextensions.JSONSchemaProps
		expectedError error
	}{
		"do nothing when no status mapping": {
			request: &MappingProcessorRequest{
				CRD:           groupBaseCRD(t),
				MappingConfig: &configv1alpha1.CRDMapping{},
			},
			expectedProps: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: `Most recently observed read-only status of the group for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`,
				Properties: map[string]apiextensions.JSONSchemaProps{
					"conditions": {
						Type:        "array",
						Description: "Represents the latest available observations of a resource's current state.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"type": {
										Type:        "string",
										Description: "Type of condition.",
									},
									"status": {
										Type:        "string",
										Description: "Status of the condition, one of True, False, Unknown.",
									},
									"lastTransitionTime": {
										Type:        "string",
										Format:      "date-time",
										Description: "Last time the condition transitioned from one status to another.",
									},
									"reason": {
										Type:        "string",
										Description: "The reason for the condition's last transition.",
									},
									"message": {
										Type:        "string",
										Description: "A human readable message indicating details about the transition.",
									},
									"observedGeneration": {
										Type:        "integer",
										Description: "observedGeneration represents the .metadata.generation that the condition was set based upon.",
									},
								},
								Required: []string{"type", "status"},
							},
						},
						XListMapKeys: []string{"type"},
						XListType:    ptr.To("map"),
					},
				},
			},
			expectedError: nil,
		},
		"add status schema to the CRD mapped to path": {
			request: groupMappingRequest(t, groupBaseCRD(t), nil, groupStatusConvertMock(t)),
			expectedProps: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: `Most recently observed read-only status of the group for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`,
				Properties: map[string]apiextensions.JSONSchemaProps{
					"v20250312": {
						Type:        "object",
						Description: "The last observed Atlas state of the group resource for version v20250312.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"id": {
								Type:        "string",
								Description: "Unique 24-hexadecimal digit string that identifies the group",
							},
						},
					},
					"conditions": {
						Type:        "array",
						Description: "Represents the latest available observations of a resource's current state.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"type": {
										Type:        "string",
										Description: "Type of condition.",
									},
									"status": {
										Type:        "string",
										Description: "Status of the condition, one of True, False, Unknown.",
									},
									"lastTransitionTime": {
										Type:        "string",
										Format:      "date-time",
										Description: "Last time the condition transitioned from one status to another.",
									},
									"reason": {
										Type:        "string",
										Description: "The reason for the condition's last transition.",
									},
									"message": {
										Type:        "string",
										Description: "A human readable message indicating details about the transition.",
									},
									"observedGeneration": {
										Type:        "integer",
										Description: "observedGeneration represents the .metadata.generation that the condition was set based upon.",
									},
								},
								Required: []string{"type", "status"},
							},
						},
						XListMapKeys: []string{"type"},
						XListType:    ptr.To("map"),
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Status{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			props := test.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["status"]
			assert.Equal(t, test.expectedProps, props)
		})
	}
}

func groupStatusConvertMock(t *testing.T) converterFuncMock {
	t.Helper()

	return func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
		return &apiextensions.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextensions.JSONSchemaProps{
				"id": {
					Type:        "string",
					Description: "Unique 24-hexadecimal digit string that identifies the group",
				},
			},
		}
	}
}
