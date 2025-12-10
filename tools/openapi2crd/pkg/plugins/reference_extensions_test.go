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
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestReferenceMetadataName(t *testing.T) {
	p := &ReferenceExtensions{}
	assert.Equal(t, "reference_metadata", p.Name())
}

func TestReferenceMetadataProcess(t *testing.T) {
	tests := map[string]struct {
		request            *ExtensionProcessorRequest
		expectedExtensions map[string]interface{}
		expectedErr        error
		refName            string
		isEntryRef         bool
	}{
		"do nothing when no references": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: nil,
			expectedErr:        nil,
			refName:            "",
			isEntryRef:         false,
		},
		"add reference metadata": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Name:     "myRef",
								Property: "spec.myRef",
								Target: configv1alpha1.Target{
									Type: configv1alpha1.Type{
										Kind:     "MyKind",
										Group:    "mygroup.example.com",
										Version:  "v1",
										Resource: "myresources",
									},
									Properties: []string{"status.id"},
								},
							},
						},
					},
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: map[string]interface{}{
				"x-kubernetes-mapping": map[string]interface{}{
					"type": map[string]interface{}{
						"kind":     "MyKind",
						"group":    "mygroup.example.com",
						"version":  "v1",
						"resource": "myresources",
					},
					"nameSelector": ".name",
					"properties": []string{
						"status.id",
					},
				},
				"x-openapi-mapping": map[string]interface{}{
					"property": "spec.myRef",
				},
			},
			expectedErr: nil,
			refName:     "myRef",
			isEntryRef:  false,
		},
		"error when reference target has no properties": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Name:     "myRef",
								Property: "spec.myRef",
								Target: configv1alpha1.Target{
									Type: configv1alpha1.Type{
										Kind:     "MyKind",
										Group:    "mygroup.example.com",
										Version:  "v1",
										Resource: "myresources",
									},
									Properties: []string{},
								},
							},
						},
					},
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: nil,
			expectedErr:        errors.New("reference target must have at least one property defined"),
			refName:            "",
			isEntryRef:         false,
		},
		"add reference metadata to entry": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					EntryMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Name:     "clusterRef",
								Property: "spec.clusterName",
								Target: configv1alpha1.Target{
									Type: configv1alpha1.Type{
										Kind:     "Cluster",
										Group:    "atlas.generated.mongodb.com",
										Version:  "v1",
										Resource: "clusters",
									},
									Properties: []string{"status.name"},
								},
							},
						},
					},
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type: &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{
												"entry": {
													Value: &openapi3.Schema{
														Type:       &openapi3.Types{"object"},
														Properties: map[string]*openapi3.SchemaRef{},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: map[string]interface{}{
				"x-kubernetes-mapping": map[string]interface{}{
					"type": map[string]interface{}{
						"kind":     "Cluster",
						"group":    "atlas.generated.mongodb.com",
						"version":  "v1",
						"resource": "clusters",
					},
					"nameSelector": ".name",
					"properties": []string{
						"status.name",
					},
				},
				"x-openapi-mapping": map[string]interface{}{
					"property": "spec.clusterName",
				},
			},
			expectedErr: nil,
			refName:     "clusterRef",
			isEntryRef:  true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &ReferenceExtensions{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedExtensions != nil {
				majorVersionProps := tt.request.ExtensionsSchema.Properties["spec"].Value.Properties[tt.request.MappingConfig.MajorVersion].Value.Properties
				if tt.isEntryRef {
					assert.Equal(t, tt.expectedExtensions, majorVersionProps["entry"].Value.Properties[tt.refName].Value.Extensions)
				} else {
					assert.Equal(t, tt.expectedExtensions, majorVersionProps[tt.refName].Value.Extensions)
				}
			} else if tt.refName == "" {
				if tt.isEntryRef {
					assert.Empty(t, tt.request.ExtensionsSchema.Properties["spec"].Value.Properties[tt.request.MappingConfig.MajorVersion].Value.Properties["entry"].Value.Properties)
				} else {
					assert.Empty(t, tt.request.ExtensionsSchema.Properties["spec"].Value.Properties[tt.request.MappingConfig.MajorVersion].Value.Properties)
				}
			}
		})
	}
}
