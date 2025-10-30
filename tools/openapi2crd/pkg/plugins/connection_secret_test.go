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
)

func TestConnectionSecretName(t *testing.T) {
	p := &ConnectionSecret{}
	assert.Equal(t, "connection_secret", p.Name())
}

func TestConnectionSecretProcess(t *testing.T) {
	tests := map[string]struct {
		request          *MappingProcessorRequest
		expectedProperty apiextensions.JSONSchemaProps
		expectedErr      error
	}{
		"add connectionSecretRef property and validations": {
			request: mappingRequestWithReferences(t, groupBaseCRDWithParameters(t)),
			expectedProperty: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "Specification of the group supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"connectionSecretRef": {
						Type:        "object",
						Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"name": {
								Type:        "string",
								Description: "Name of the secret containing the Atlas credentials.",
							},
						},
					},
					"v20250312": {
						Type:        "object",
						Description: "The spec of the group resource for version v20250312.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"groupId": {
								Type: "string",
							},
							"groupRef": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
								},
							},
						},
					},
				},
				XValidations: apiextensions.ValidationRules{
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250312.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250312.groupId))",
						Message: "spec.connectionSecretRef must be set if spec.v20250312.groupId is set.",
					},
				},
			},
		},
		"add connectionSecretRef property and validations to multiple versions": {
			request: mappingRequestWithReferences(t, groupBaseCRDMultiVersionWithParameters(t)),
			expectedProperty: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "Specification of the group supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"connectionSecretRef": {
						Type:        "object",
						Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"name": {
								Type:        "string",
								Description: "Name of the secret containing the Atlas credentials.",
							},
						},
					},
					"v20250312": {
						Type:        "object",
						Description: "The spec of the group resource for version v20250312.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"groupId": {
								Type: "string",
							},
							"groupRef": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
								},
							},
						},
					},
					"v20250219": {
						Type:        "object",
						Description: "The spec of the group resource for version v20250219.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"groupId": {
								Type: "string",
							},
							"groupRef": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
								},
							},
						},
					},
				},
				XValidations: apiextensions.ValidationRules{
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250219.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250219.groupId))",
						Message: "spec.connectionSecretRef must be set if spec.v20250219.groupId is set.",
					},
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250312.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250312.groupId))",
						Message: "spec.connectionSecretRef must be set if spec.v20250312.groupId is set.",
					},
				},
			},
		},
		"version but not groupId in CRD mappings": {
			request: mappingRequestWithReferences(t, orgBaseCRDWithParameters(t)),
			expectedProperty: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "Specification of the organization supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"connectionSecretRef": {
						Type:        "object",
						Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"name": {
								Type:        "string",
								Description: "Name of the secret containing the Atlas credentials.",
							},
						},
					},
					"v20250312": {
						Type:        "object",
						Description: "The spec of the organization resource for version v20250312.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"orgID": {
								Type: "string",
							},
						},
					},
				},
			},
		},
		"no versions in CRD mappings": {
			request: mappingRequestWithReferences(t, groupBaseCRD(t)),
			expectedProperty: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "Specification of the group supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"connectionSecretRef": {
						Type:        "object",
						Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"name": {
								Type:        "string",
								Description: "Name of the secret containing the Atlas credentials.",
							},
						},
					},
				},
			},
			expectedErr: errors.New("version v20250312 not found in spec"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &ConnectionSecret{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedProperty, tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"])
		})
	}
}

func groupBaseCRDWithParameters(t *testing.T) *apiextensions.CustomResourceDefinition {
	crd := groupBaseCRD(t)
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties["v20250312"] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250312.",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"groupId": {
				Type: "string",
			},
			"groupRef": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type: "string",
					},
				},
			},
		},
	}

	return crd
}

func groupBaseCRDMultiVersionWithParameters(t *testing.T) *apiextensions.CustomResourceDefinition {
	crd := groupBaseCRD(t)
	spec := crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"]
	spec.Properties["connectionSecretRef"] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"name": {
				Type:        "string",
				Description: "Name of the secret containing the Atlas credentials.",
			},
		},
	}
	spec.XValidations = apiextensions.ValidationRules{
		apiextensions.ValidationRule{
			Rule:    "(has(self.v20250219.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250219.groupId))",
			Message: "spec.connectionSecretRef must be set if spec.v20250219.groupId is set.",
		},
	}
	spec.Properties["v20250312"] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250312.",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"groupId": {
				Type: "string",
			},
			"groupRef": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type: "string",
					},
				},
			},
		},
	}
	spec.Properties["v20250219"] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250219.",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"groupId": {
				Type: "string",
			},
			"groupRef": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type: "string",
					},
				},
			},
		},
	}
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = spec

	return crd
}

func orgBaseCRDWithParameters(t *testing.T) *apiextensions.CustomResourceDefinition {
	t.Helper()

	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "groups.atlas.generated.mongodb.com",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: "atlas.generated.mongodb.com",
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:       "Organization",
				ListKind:   "OrganizationList",
				Plural:     "Organizations",
				Singular:   "organization",
				ShortNames: []string{"ao"},
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
					Description: "A organization, managed by the MongoDB Kubernetes Atlas Operator.",
					Properties: map[string]apiextensions.JSONSchemaProps{
						"spec": {
							Type:        "object",
							Description: "Specification of the organization supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"v20250312": {
									Type:        "object",
									Description: "The spec of the organization resource for version v20250312.",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"orgID": {
											Type: "string",
										},
									},
								},
							},
						},
						"status": {
							Type:        "object",
							Description: `Most recently observed read-only status of the organization for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`,
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
