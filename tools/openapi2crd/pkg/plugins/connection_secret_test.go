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
)

func TestConnectionSecretName(t *testing.T) {
	p := &ConnectionSecret{}
	assert.Equal(t, "connection_secret", p.Name())
}

func TestConnectionSecretProcess(t *testing.T) {
	tests := map[string]struct {
		request            *MappingProcessorRequest
		expectedProperty   apiextensions.JSONSchemaProps
		expectedValidation apiextensions.ValidationRules
		expectedErr        error
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
						XValidations: apiextensions.ValidationRules{
							apiextensions.ValidationRule{
								Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
								Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
							},
						},
					},
				},
				XValidations: apiextensions.ValidationRules{
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250312.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250312.groupId))",
						Message: "connectionSecretRef must be set if groupId is set for version v20250312",
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
						XValidations: apiextensions.ValidationRules{
							apiextensions.ValidationRule{
								Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
								Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
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
						XValidations: apiextensions.ValidationRules{
							apiextensions.ValidationRule{
								Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
								Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
							},
						},
					},
				},
				XValidations: apiextensions.ValidationRules{
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250219.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250219.groupId))",
						Message: "connectionSecretRef must be set if groupId is set for version v20250219",
					},
					apiextensions.ValidationRule{
						Rule:    "(has(self.v20250312.groupId) && has(self.connectionSecretRef)) || (!has(self.v20250312.groupId))",
						Message: "connectionSecretRef must be set if groupId is set for version v20250312",
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
			Message: "connectionSecretRef must be set if groupId is set for version v20250219",
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
		XValidations: apiextensions.ValidationRules{
			apiextensions.ValidationRule{
				Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
				Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
			},
		},
	}
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = spec

	return crd
}
