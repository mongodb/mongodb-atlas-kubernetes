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
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
)

func TestEntryName(t *testing.T) {
	p := &Entry{}
	assert.Equal(t, "entry", p.Name())
}

func TestEntryProcess(t *testing.T) {
	tests := map[string]struct {
		request       *MappingProcessorRequest
		expectedEntry apiextensions.JSONSchemaProps
		expectedErr   error
	}{
		"add entry schema to the CRD mapped to schema": {
			request: groupMappingRequest(t, groupBaseCRDWithMajorVersion(t), entryInitialExtensionsSchema(t), entryConverterMock(t)),
			expectedEntry: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "The entry fields of the group resource spec. These fields can be set for creating and updating groups.",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type:        "string",
						Description: "Human-readable label that identifies this group.",
					},
					"orgId": {
						Type:        "string",
						Description: "Unique 24-hexadecimal digit string that identifies the organization to which this group belongs.",
					},
					"teamIds": {
						Type:        "array",
						Description: "List of unique 24-hexadecimal digit strings that identify the teams to which this group belongs.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "string",
							},
						},
					},
					"labels": {
						Type:        "array",
						Description: "List of key-value pairs that can be attached to a group.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"key": {
										Type:        "string",
										Description: "Label key.",
									},
									"value": {
										Type:        "string",
										Description: "Label value.",
									},
								},
							},
						},
					},
				},
			},
		},
		"add entry schema to the CRD mapped to path": {
			request: groupMappingRequestWithPath(t, groupBaseCRDWithMajorVersion(t), entryInitialExtensionsSchema(t), entryConverterMock(t)),
			expectedEntry: apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: "The entry fields of the group resource spec. These fields can be set for creating and updating groups.",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type:        "string",
						Description: "Human-readable label that identifies this group.",
					},
					"orgId": {
						Type:        "string",
						Description: "Unique 24-hexadecimal digit string that identifies the organization to which this group belongs.",
					},
					"teamIds": {
						Type:        "array",
						Description: "List of unique 24-hexadecimal digit strings that identify the teams to which this group belongs.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "string",
							},
						},
					},
					"labels": {
						Type:        "array",
						Description: "List of key-value pairs that can be attached to a group.",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"key": {
										Type:        "string",
										Description: "Label key.",
									},
									"value": {
										Type:        "string",
										Description: "Label value.",
									},
								},
							},
						},
					},
				},
			},
		},
		"error when schema not found": {
			request:     groupMappingWithNonExistingSchema(t, groupBaseCRDWithMajorVersion(t), entryInitialExtensionsSchema(t), entryConverterMock(t)),
			expectedErr: fmt.Errorf("entry schema %q not found in openapi spec", "NonExistentSchema"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Entry{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			entry := tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[tt.request.MappingConfig.MajorVersion].Properties["entry"]
			assert.Equal(t, tt.expectedEntry, entry)
		})
	}
}

func groupMappingRequestWithPath(
	t *testing.T,
	crd *apiextensions.CustomResourceDefinition,
	extensionsSchema *openapi3.Schema,
	converterFunc converterFuncMock,
) *MappingProcessorRequest {
	req := groupMappingRequest(t, crd, extensionsSchema, converterFunc)
	req.MappingConfig.EntryMapping.Path = configv1alpha1.PropertyPath{
		Name: "/api/atlas/v2/groups",
		Verb: "post",
		RequestBody: configv1alpha1.RequestBody{
			MimeType: "application/vnd.atlas.2025-03-12+json",
		},
	}
	req.MappingConfig.EntryMapping.Schema = ""

	return req
}

func groupMappingWithNonExistingSchema(
	t *testing.T,
	crd *apiextensions.CustomResourceDefinition,
	extensionsSchema *openapi3.Schema,
	converterFunc converterFuncMock,
) *MappingProcessorRequest {
	req := groupMappingRequest(t, crd, extensionsSchema, converterFunc)
	req.MappingConfig.EntryMapping.Schema = "NonExistentSchema"

	return req
}

func groupBaseCRDWithMajorVersion(t *testing.T) *apiextensions.CustomResourceDefinition {
	crd := groupBaseCRD(t)
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties["v20250312"] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250312.",
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}

	return crd
}

func entryInitialExtensionsSchema(t *testing.T) *openapi3.Schema {
	t.Helper()

	return &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: map[string]*openapi3.SchemaRef{
			"spec": {
				Value: &openapi3.Schema{
					Type:       &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{},
				},
			},
		},
	}
}

func entryConverterMock(t *testing.T) converterFuncMock {
	t.Helper()

	return func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
		return &apiextensions.JSONSchemaProps{
			Type:        "object",
			Description: "The entry fields of the group resource spec. These fields can be set for creating and updating groups.",
			Properties: map[string]apiextensions.JSONSchemaProps{
				"name": {
					Type:        "string",
					Description: "Human-readable label that identifies this group.",
				},
				"orgId": {
					Type:        "string",
					Description: "Unique 24-hexadecimal digit string that identifies the organization to which this group belongs.",
				},
				"teamIds": {
					Type:        "array",
					Description: "List of unique 24-hexadecimal digit strings that identify the teams to which this group belongs.",
					Items: &apiextensions.JSONSchemaPropsOrArray{
						Schema: &apiextensions.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"labels": {
					Type:        "array",
					Description: "List of key-value pairs that can be attached to a group.",
					Items: &apiextensions.JSONSchemaPropsOrArray{
						Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"key": {
									Type:        "string",
									Description: "Label key.",
								},
								"value": {
									Type:        "string",
									Description: "Label value.",
								},
							},
						},
					},
				},
			},
		}
	}
}
