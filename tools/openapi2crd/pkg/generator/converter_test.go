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
			input: converter.NewPropertyConvertInput(
				regularSchemaRef(),
				openapi3.NewSchemaRef("", openapi3.NewSchema()),
				&configv1alpha1.PropertyMapping{Schema: "Pet"},
				[]string{},
			),
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
			input: converter.NewPropertyConvertInput(
				oneOfSchemaRef(),
				openapi3.NewSchemaRef("", openapi3.NewSchema()),
				&configv1alpha1.PropertyMapping{Schema: "Pet"},
				[]string{},
			),
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

func TestGeneratorConvertRecursiveSchemas(t *testing.T) {
	example := apiextensions.JSON(nil)
	trueVar := true

	newGenerator := func(t *testing.T) *Generator {
		t.Helper()
		propertyPlugin := plugins.NewPropertyPluginMock(t)
		propertyPlugin.EXPECT().Process(mock.AnythingOfType("*plugins.PropertyProcessorRequest")).
			Return(nil)
		return &Generator{
			pluginSet: &plugins.Set{
				Property: []plugins.PropertyPlugin{propertyPlugin},
			},
		}
	}

	t.Run("direct self-referencing schema", func(t *testing.T) {
		// Schema A has a property "child" that references A itself.
		// Without cycle detection this causes infinite recursion.
		schemaA := &openapi3.Schema{
			Type: &openapi3.Types{"object"},
		}
		schemaA.Properties = map[string]*openapi3.SchemaRef{
			"name": {Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
			"child": {Value: schemaA},
		}

		g := newGenerator(t)
		p := g.Convert(converter.NewPropertyConvertInput(
			&openapi3.SchemaRef{Value: schemaA},
			openapi3.NewSchemaRef("", openapi3.NewSchema()),
			&configv1alpha1.PropertyMapping{Schema: "A"},
			[]string{},
		))

		assert.NotNil(t, p)
		// The recursive "child" property should be terminated with x-kubernetes-preserve-unknown-fields.
		assert.Equal(t, &trueVar, p.Properties["child"].XPreserveUnknownFields)
	})

	t.Run("indirect cycle A -> B -> A", func(t *testing.T) {
		// A references B which references A back.
		schemaA := &openapi3.Schema{
			Type: &openapi3.Types{"object"},
		}
		schemaB := &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: map[string]*openapi3.SchemaRef{
				"parent": {Value: schemaA},
			},
		}
		schemaA.Properties = map[string]*openapi3.SchemaRef{
			"b": {Value: schemaB},
		}

		g := newGenerator(t)
		p := g.Convert(converter.NewPropertyConvertInput(
			&openapi3.SchemaRef{Value: schemaA},
			openapi3.NewSchemaRef("", openapi3.NewSchema()),
			&configv1alpha1.PropertyMapping{Schema: "A"},
			[]string{},
		))

		assert.NotNil(t, p)
		// B's "parent" property references A which is an ancestor — should be terminated.
		assert.Equal(t, &trueVar, p.Properties["b"].Properties["parent"].XPreserveUnknownFields)
	})

	t.Run("diamond reference is not a false cycle", func(t *testing.T) {
		// A has properties B and C, both reference shared schema D.
		// D should be fully expanded in both branches — not treated as a cycle.
		schemaD := &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: map[string]*openapi3.SchemaRef{
				"value": {Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
			},
		}
		schemaA := &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: map[string]*openapi3.SchemaRef{
				"b": {Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{
						"d": {Value: schemaD},
					},
				}},
				"c": {Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{
						"d": {Value: schemaD},
					},
				}},
			},
		}

		g := newGenerator(t)
		p := g.Convert(converter.NewPropertyConvertInput(
			&openapi3.SchemaRef{Value: schemaA},
			openapi3.NewSchemaRef("", openapi3.NewSchema()),
			&configv1alpha1.PropertyMapping{Schema: "A"},
			[]string{},
		))

		assert.NotNil(t, p)
		// D should be fully expanded in both branches (not terminated as a cycle).
		expectedD := apiextensions.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextensions.JSONSchemaProps{
				"value": {
					Type:       "string",
					Properties: map[string]apiextensions.JSONSchemaProps{},
					Example:    &example,
				},
			},
			Example: &example,
		}
		assert.Equal(t, expectedD, p.Properties["b"].Properties["d"])
		assert.Equal(t, expectedD, p.Properties["c"].Properties["d"])
		// Verify neither branch was falsely terminated as a cycle.
		assert.Nil(t, p.Properties["b"].Properties["d"].XPreserveUnknownFields)
		assert.Nil(t, p.Properties["c"].Properties["d"].XPreserveUnknownFields)
	})
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
