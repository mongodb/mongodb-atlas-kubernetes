package plugins

import (
	"testing"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestSkippedPropertyName(t *testing.T) {
	p := &SkippedProperty{}
	assert.Equal(t, "skipped_property", p.Name())
}

func TestSkippedPropertyProcess(t *testing.T) {
	tests := map[string]struct {
		request       *PropertyProcessorRequest
		expectedProps *apiextensions.JSONSchemaProps
		expectedError error
	}{
		"do nothing when skipped property config is empty": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &configv1alpha1.PropertyMapping{},
				Property: &apiextensions.JSONSchemaProps{
					Properties: map[string]apiextensions.JSONSchemaProps{
						"description": {
							Type: "string",
						},
					},
					Required: []string{"description"},
				},
				OpenAPISchema: &openapi3.Schema{Type: &openapi3.Types{"object"}, Required: []string{"description"}},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Properties: map[string]apiextensions.JSONSchemaProps{
					"description": {
						Type: "string",
					},
				},
				Required: []string{"description"},
			},
			expectedError: nil,
		},
		"skip property": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Filters: configv1alpha1.Filters{
						SkipProperties: []string{"$.description"},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Properties: map[string]apiextensions.JSONSchemaProps{
						"description": {
							Type: "string",
						},
					},
					Required: []string{"description"},
				},
				Path: []string{"$", "description"},
			},
			expectedProps: nil,
			expectedError: nil,
		},
		"remove required property set to skip": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Filters: configv1alpha1.Filters{
						SkipProperties: []string{"$.details[*].description"},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Properties: map[string]apiextensions.JSONSchemaProps{
						"details": {
							Type: "array",
							Items: &apiextensions.JSONSchemaPropsOrArray{
								Schema: &apiextensions.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"name": {
											Type: "string",
										},
										"description": {
											Type: "string",
										},
									},
								},
							},
						},
					},
					Required: []string{"name", "description"},
				},
				OpenAPISchema: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{
						"details": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"array"},
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"object"},
										Properties: map[string]*openapi3.SchemaRef{
											"name": {
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"string"},
												},
											},
											"description": {
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"string"},
												},
											},
										},
									},
								},
							},
						},
					},
					Required: []string{"name", "description"},
				},
				Path: []string{"$", "details[*]"},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Properties: map[string]apiextensions.JSONSchemaProps{
					"details": {
						Type: "array",
						Items: &apiextensions.JSONSchemaPropsOrArray{
							Schema: &apiextensions.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
									"description": {
										Type: "string",
									},
								},
							},
						},
					},
				},
				Required: []string{"name"},
			},
			expectedError: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &SkippedProperty{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedProps, test.request.Property)
		})
	}
}
