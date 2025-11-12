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

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestReferenceName(t *testing.T) {
	p := &References{}
	assert.Equal(t, "reference", p.Name())
}

func TestReferenceProcess(t *testing.T) {
	tests := map[string]struct {
		request       *MappingProcessorRequest
		expectedProps map[string]apiextensions.JSONSchemaProps
		expectedError error
	}{
		"do nothing when no references": {
			request:       groupMappingRequest(t, groupBaseCRDWithMajorVersion(t), nil, nil),
			expectedProps: map[string]apiextensions.JSONSchemaProps{},
			expectedError: nil,
		},
		"add reference to group": {
			request: mappingRequestWithReferences(t, dataFederationCRD(t)),
			expectedProps: map[string]apiextensions.JSONSchemaProps{
				"groupRef": {
					Type:        "object",
					Description: "A reference to a \"Group\" resource.\nThe value of \"$.status.v20250312.id\" will be used to set \"groupId\".\nMutually exclusive with the \"groupId\" property.",
					Properties: map[string]apiextensions.JSONSchemaProps{
						"name": {
							Type:        "string",
							Description: `Name of the "Group" resource.`,
						},
					},
				},
			},
			expectedError: nil,
		},
		"error when reference target has no properties": {
			request: &MappingProcessorRequest{
				CRD: dataFederationCRD(t),
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					OpenAPIRef: configv1alpha1.LocalObjectReference{
						Name: "v20250312",
					},
					ParametersMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Target: configv1alpha1.Target{},
							},
						},
					},
				},
			},
			expectedProps: map[string]apiextensions.JSONSchemaProps{},
			expectedError: errors.New("reference target must have at least one property defined"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &References{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			spec := test.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[test.request.MappingConfig.MajorVersion]
			assert.Equal(t, test.expectedProps, spec.Properties)
		})
	}
}

func dataFederationCRD(t *testing.T) *apiextensions.CustomResourceDefinition {
	t.Helper()

	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "datafederations.atlas.generated.mongodb.com",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: "atlas.generated.mongodb.com",
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     "datafederations",
				Singular:   "datafederation",
				Kind:       "DataFederation",
				ShortNames: []string{"adf"},
				Categories: []string{"atlas"},
			},
			Scope: apiextensions.NamespaceScoped,
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
			PreserveUnknownFields: ptr.To(false),
			Validation: &apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Type:        "object",
					Description: "A group, managed by the MongoDB Kubernetes Atlas Operator.",
					Properties: map[string]apiextensions.JSONSchemaProps{
						"spec": {
							Type:        "object",
							Description: "Specification of the group supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"v20250312": {
									Type:        "object",
									Description: "The spec of the group resource for version v20250312.",
									Properties:  map[string]apiextensions.JSONSchemaProps{},
								},
							},
						},
						"status": {
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
					},
				},
			},
			Subresources: &apiextensions.CustomResourceSubresources{
				Status: &apiextensions.CustomResourceSubresourceStatus{},
			},
		},
	}
}

func mappingRequestWithReferences(
	t *testing.T,
	crd *apiextensions.CustomResourceDefinition,
) *MappingProcessorRequest {
	t.Helper()

	return &MappingProcessorRequest{
		CRD: crd,
		MappingConfig: &configv1alpha1.CRDMapping{
			MajorVersion: "v20250312",
			OpenAPIRef: configv1alpha1.LocalObjectReference{
				Name: "v20250312",
			},
			ParametersMapping: configv1alpha1.PropertyMapping{
				Path: configv1alpha1.PropertyPath{
					Name: "/api/atlas/v2/groups/{groupId}/dataFederation",
					Verb: "post",
				},
				References: []configv1alpha1.Reference{
					{
						Name:     "groupRef",
						Property: "$.groupId",
						Target: configv1alpha1.Target{
							Type: configv1alpha1.Type{
								Kind:     "Group",
								Resource: "Groups",
								Group:    "atlas.generated.mongodb.com",
								Version:  "v1",
							},
							Properties: []string{"$.status.v20250312.id"},
						},
					},
				},
			},
		},
	}
}
