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

package generator

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGeneratorMajorVersions(t *testing.T) {
	tests := map[string]struct {
		config         v1alpha1.CRDConfig
		expectedResult []string
	}{
		"mapping with single version": {
			config: v1alpha1.CRDConfig{
				Mappings: []v1alpha1.CRDMapping{
					{
						MajorVersion: "v1",
					},
				},
			},
			expectedResult: []string{"- v1"},
		},
		"mapping with multiple versions": {
			config: v1alpha1.CRDConfig{
				Mappings: []v1alpha1.CRDMapping{
					{
						MajorVersion: "v1",
					},
					{
						MajorVersion: "v2",
					},
				},
			},
			expectedResult: []string{"- v1", "- v2"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &Generator{}
			result := g.majorVersions(tt.config)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestClearPropertiesWithoutExtensions(t *testing.T) {
	tests := map[string]struct {
		schema         *openapi3.Schema
		expectedResult bool
	}{
		"schema with properties and extensions": {
			schema: &openapi3.Schema{
				Properties: map[string]*openapi3.SchemaRef{
					"prop1": {Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
					"prop2": {Value: &openapi3.Schema{}},
				},
			},
			expectedResult: true,
		},
		"schema with properties without extensions": {
			schema: &openapi3.Schema{
				Properties: map[string]*openapi3.SchemaRef{
					"prop1": {Value: &openapi3.Schema{}},
					"prop2": {Value: &openapi3.Schema{}},
				},
			},
			expectedResult: false,
		},
		"schema with nested properties and extensions": {
			schema: &openapi3.Schema{
				Properties: map[string]*openapi3.SchemaRef{
					"prop1": {Value: &openapi3.Schema{
						Properties: map[string]*openapi3.SchemaRef{
							"nestedProp": {Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
						},
					}},
					"prop2": {Value: &openapi3.Schema{}},
				},
			},
			expectedResult: true,
		},
		"schema with additionalProperties and extensions": {
			schema: &openapi3.Schema{
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
				},
			},
			expectedResult: true,
		},
		"schema with items and extensions": {
			schema: &openapi3.Schema{
				Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
			},
			expectedResult: true,
		},
		"schema with allOf and extensions": {
			schema: &openapi3.Schema{
				AllOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
				},
			},
			expectedResult: true,
		},
		"schema with anyOf and extensions": {
			schema: &openapi3.Schema{
				AnyOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
				},
			},
			expectedResult: true,
		},
		"schema with oneOf and extensions": {
			schema: &openapi3.Schema{
				OneOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{Extensions: map[string]interface{}{"x-kubernetes-preserve-unknown-fields": true}}},
				},
			},
			expectedResult: true,
		},
		"nil schema": {
			schema:         nil,
			expectedResult: false,
		},
		"empty schema": {
			schema:         &openapi3.Schema{},
			expectedResult: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := clearPropertiesWithoutExtensions(tt.schema)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGeneratorGenerate(t *testing.T) {
	tests := map[string]struct {
		//openapi        *openapi3.T
		apiDefinitions map[string]v1alpha1.OpenAPIDefinition
		config         *v1alpha1.CRDConfig
		expectedResult *apiextensions.CustomResourceDefinition
		expectError    bool
	}{
		"generate with valid openapi and config": {
			apiDefinitions: map[string]v1alpha1.OpenAPIDefinition{
				"Pet": {
					Name: "Pet",
					Path: "testdata/openapi.yaml",
				},
			},
			config: &v1alpha1.CRDConfig{
				Mappings: []v1alpha1.CRDMapping{
					{
						OpenAPIRef: v1alpha1.LocalObjectReference{
							Name: "Pet",
						},
					},
				},
			},
			expectedResult: &apiextensions.CustomResourceDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name: "examples.test.com",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "test.com",
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:     "examples",
						Singular:   "example",
						Kind:       "Example",
						ListKind:   "ExampleList",
						ShortNames: []string{"ex"},
						Categories: []string{"test"},
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "v1",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.NamespaceScoped,
					Validation: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"name": {
											Type:        "string",
											Description: "Name of the resource",
										},
									},
								},
								"status": {
									Type:        "object",
									Description: "Most recently observed status of the example.",
								},
							},
							Required: []string{"spec"},
						},
					},
					PreserveUnknownFields: nil,
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"v1"},
				},
			},
		},
		"duplicate major versions": {
			apiDefinitions: map[string]v1alpha1.OpenAPIDefinition{
				"Pet": {
					Name: "Pet",
					Path: "testdata/openapi.yaml",
				},
			},
			config: &v1alpha1.CRDConfig{
				Mappings: []v1alpha1.CRDMapping{
					{
						OpenAPIRef: v1alpha1.LocalObjectReference{
							Name: "Pet",
						},
						MajorVersion: "v1",
					},
					{
						OpenAPIRef: v1alpha1.LocalObjectReference{
							Name: "Pet",
						},
						MajorVersion: "v1",
					},
				},
			},
			expectError: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			openapiLoader := config.NewLoaderMock(t)
			openapiLoader.EXPECT().Load(context.Background(), "testdata/openapi.yaml").Return(&openapi3.T{}, nil)

			atlasLoader := config.NewLoaderMock(t)

			crdPlugin := plugins.NewCRDPluginMock(t)
			crdPlugin.EXPECT().Process(mock.AnythingOfType("*plugins.CRDProcessorRequest")).
				RunAndReturn(func(request *plugins.CRDProcessorRequest) error {
					baseCRD(request.CRD)
					return nil
				})
			mappingPlugin := plugins.NewMappingPluginMock(t)
			mappingPlugin.EXPECT().Process(mock.AnythingOfType("*plugins.MappingProcessorRequest")).
				RunAndReturn(func(request *plugins.MappingProcessorRequest) error {
					request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties["name"] = apiextensions.JSONSchemaProps{
						Type:        "string",
						Description: "Name of the resource",
					}

					return nil
				})
			extensionPlugin := plugins.NewExtensionPluginMock(t)
			extensionPlugin.EXPECT().Process(mock.AnythingOfType("*plugins.ExtensionProcessorRequest")).Return(nil)

			g := &Generator{
				definitions: tt.apiDefinitions,
				pluginSet: &plugins.Set{
					CRD:       []plugins.CRDPlugin{crdPlugin},
					Mapping:   []plugins.MappingPlugin{mappingPlugin},
					Extension: []plugins.ExtensionPlugin{extensionPlugin},
				},
				openapiLoader: openapiLoader,
				atlasLoader:   atlasLoader,
			}
			result, err := g.Generate(context.Background(), tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func baseCRD(crd *apiextensions.CustomResourceDefinition) {
	crd.ObjectMeta = v1.ObjectMeta{
		Name: "examples.test.com",
	}
	crd.Spec = apiextensions.CustomResourceDefinitionSpec{
		Group: "test.com",
		Names: apiextensions.CustomResourceDefinitionNames{
			Plural:     "examples",
			Singular:   "example",
			Kind:       "Example",
			ListKind:   "ExampleList",
			ShortNames: []string{"ex"},
			Categories: []string{"test"},
		},
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    "v1",
				Served:  true,
				Storage: true,
			},
		},
		Validation: &apiextensions.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"spec": {
						Type:       "object",
						Properties: map[string]apiextensions.JSONSchemaProps{},
					},
					"status": {
						Type:        "object",
						Description: "Most recently observed status of the example.",
					},
				},
				Required: []string{"spec"},
			},
		},
		Scope: apiextensions.NamespaceScoped,
	}
	crd.Status = apiextensions.CustomResourceDefinitionStatus{
		StoredVersions: []string{"v1"},
	}
}
