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
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/plugins"
)

func TestGeneratorConvert(t *testing.T) {
	example := apiextensions.JSON(nil)
	trueVar := true

	tests := map[string]struct {
		input         converter.PropertyConvertInput
		expectedProps *apiextensions.JSONSchemaProps
	}{
		"standard schema": {
			input: converter.PropertyConvertInput{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Schema: "Pet",
				},
				Schema:              regularSchemaRef(),
				ExtensionsSchemaRef: openapi3.NewSchemaRef("", openapi3.NewSchema()),
				Path:                []string{},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Type: "object",
				AllOf: []apiextensions.JSONSchemaProps{
					{
						Type: "object",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"id": {
								Type:       "integer",
								Properties: map[string]apiextensions.JSONSchemaProps{},
								Example:    &example,
							},
						},
						Example: &example,
					},
					{
						Type: "object",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"active": {
								Type:       "boolean",
								Properties: map[string]apiextensions.JSONSchemaProps{},
								Example:    &example,
							},
						},
						Example: &example,
					},
				},
				Properties: map[string]apiextensions.JSONSchemaProps{
					"tags": {
						Type: "array",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type:       "string",
								Properties: map[string]apiextensions.JSONSchemaProps{},
								Example:    &example,
							},
						},
						Properties: map[string]apiextensions.JSONSchemaProps{},
						Example:    &example,
					},
					"name": {
						Type:       "string",
						Properties: map[string]apiextensions.JSONSchemaProps{},
						Example:    &example,
					},
					"attributes": {
						Type:       "object",
						Properties: map[string]apiextensions.JSONSchemaProps{},
						AdditionalProperties: &apiextensions.JSONSchemaPropsOrBool{
							Allows: true,
							Schema: &apiextensions.JSONSchemaProps{
								Type:       "string",
								Properties: map[string]apiextensions.JSONSchemaProps{},
								Example:    &example,
							},
						},
						Example: &example,
					},
				},
				AdditionalProperties: &apiextensions.JSONSchemaPropsOrBool{
					Allows: true,
					Schema: &apiextensions.JSONSchemaProps{
						Type:       "string",
						Properties: map[string]apiextensions.JSONSchemaProps{},
						Example:    &example,
					},
				},
				Example: &example,
			},
		},
		"oneOf schema": {
			input: converter.PropertyConvertInput{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Schema: "Pet",
				},
				Schema:              oneOfSchemaRef(),
				ExtensionsSchemaRef: openapi3.NewSchemaRef("", openapi3.NewSchema()),
				Path:                []string{},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Type: "object",
				OneOf: []apiextensions.JSONSchemaProps{
					{
						Type:       "string",
						Properties: map[string]apiextensions.JSONSchemaProps{},
						Example:    &example,
					},
					{
						Type:       "integer",
						Properties: map[string]apiextensions.JSONSchemaProps{},
						Example:    &example,
					},
				},
				Properties:             map[string]apiextensions.JSONSchemaProps{},
				XPreserveUnknownFields: &trueVar,
				Example:                &example,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			propertyPlugin := plugins.NewPropertyPluginMock(t)
			propertyPlugin.EXPECT().Process(mock.AnythingOfType("*plugins.PropertyProcessorRequest")).
				Return(nil)

			g := &Generator{
				pluginSet: &plugins.Set{
					Property: []plugins.PropertyPlugin{propertyPlugin},
				},
			}
			p := g.Convert(tt.input)
			assert.Equal(t, tt.expectedProps, p)
		})
	}
}

func regularSchemaRef() *openapi3.SchemaRef {
	return &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: map[string]*openapi3.SchemaRef{
				"name": {
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
				"tags": {
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"array"},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
							},
						},
					},
				},
				"attributes": {
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"object"},
						AdditionalProperties: openapi3.AdditionalProperties{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"string"},
								},
							},
						},
					},
				},
			},
			AllOf: openapi3.SchemaRefs{
				{
					Value: &openapi3.Schema{
						Properties: map[string]*openapi3.SchemaRef{
							"id": {
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"integer"},
								},
							},
						},
					},
				},
				{
					Value: &openapi3.Schema{
						Properties: map[string]*openapi3.SchemaRef{
							"active": {
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"boolean"},
								},
							},
						},
					},
				},
			},
			AdditionalProperties: openapi3.AdditionalProperties{
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
			},
		},
	}
}

func oneOfSchemaRef() *openapi3.SchemaRef {
	return &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			OneOf: openapi3.SchemaRefs{
				{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
					Ref: "#/components/schemas/CustomObject",
				},
				{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"integer"},
					},
				},
			},
		},
	}
}
