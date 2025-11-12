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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestBaseName(t *testing.T) {
	p := &Base{}
	assert.Equal(t, "base", p.Name())
}

func TestBaseProcess(t *testing.T) {
	tests := map[string]struct {
		request     *CRDProcessorRequest
		expectedCrd *apiextensions.CustomResourceDefinition
		expectedErr error
	}{
		"add the base of the CRD": {
			request:     groupCRDRequest(t, &apiextensions.CustomResourceDefinition{}),
			expectedCrd: groupBaseCRD(t),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Base{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedCrd, tt.request.CRD)
		})
	}
}

func TestGuessKindToResource(t *testing.T) {
	tests := map[string]struct {
		gvk              metav1.GroupVersionKind
		expectedPlural   schema.GroupVersionResource
		expectedSingular schema.GroupVersionResource
	}{
		"regular case": {
			gvk: metav1.GroupVersionKind{
				Group:   "atlas.generated.mongodb.com",
				Version: "v1",
				Kind:    "Group",
			},
			expectedPlural: schema.GroupVersionResource{
				Group:    "atlas.generated.mongodb.com",
				Version:  "v1",
				Resource: "groups",
			},
			expectedSingular: schema.GroupVersionResource{
				Group:    "atlas.generated.mongodb.com",
				Version:  "v1",
				Resource: "group",
			},
		},
		"kind ending with s": {
			gvk: metav1.GroupVersionKind{
				Group:   "example.com",
				Version: "v1",
				Kind:    "Bus",
			},
			expectedPlural: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "buses",
			},
			expectedSingular: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "bus",
			},
		},
		"kind ending with x": {
			gvk: metav1.GroupVersionKind{
				Group:   "example.com",
				Version: "v1",
				Kind:    "Box",
			},
			expectedPlural: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "boxes",
			},
			expectedSingular: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "box",
			},
		},
		"kind ending with y": {
			gvk: metav1.GroupVersionKind{
				Group:   "example.com",
				Version: "v1",
				Kind:    "City",
			},
			expectedPlural: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "cities",
			},
			expectedSingular: schema.GroupVersionResource{
				Group:    "example.com",
				Version:  "v1",
				Resource: "city",
			},
		},
		"empty kind": {
			gvk:              metav1.GroupVersionKind{},
			expectedPlural:   schema.GroupVersionResource{},
			expectedSingular: schema.GroupVersionResource{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			plural, singular := guessKindToResource(tt.gvk)
			assert.Equal(t, tt.expectedPlural, plural)
			assert.Equal(t, tt.expectedSingular, singular)
		})
	}
}

func groupCRDRequest(t *testing.T, crd *apiextensions.CustomResourceDefinition) *CRDProcessorRequest {
	t.Helper()

	return &CRDProcessorRequest{
		CRD: crd,
		CRDConfig: &configv1alpha1.CRDConfig{
			GVK: metav1.GroupVersionKind{
				Group:   "atlas.generated.mongodb.com",
				Version: "v1",
				Kind:    "Group",
			},
			Categories: []string{"atlas"},
			Mappings: []configv1alpha1.CRDMapping{
				{
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
			},
			ShortNames: []string{"ag"},
		},
	}
}

func groupBaseCRD(t *testing.T) *apiextensions.CustomResourceDefinition {
	t.Helper()

	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "groups.atlas.generated.mongodb.com",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: "atlas.generated.mongodb.com",
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:       "Group",
				ListKind:   "GroupList",
				Plural:     "groups",
				Singular:   "group",
				ShortNames: []string{"ag"},
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
							Properties:  map[string]apiextensions.JSONSchemaProps{},
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
		Status: apiextensions.CustomResourceDefinitionStatus{
			StoredVersions: []string{"v1"},
		},
	}
}
