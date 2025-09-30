package plugins

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	"tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestReadOnlyPropertyName(t *testing.T) {
	p := &ReadOnlyProperties{}
	assert.Equal(t, "read_only_properties", p.Name())
}

func TestReadOnlyPropertyProcess(t *testing.T) {
	tests := map[string]struct {
		request       *PropertyProcessorRequest
		expectedProps *apiextensions.JSONSchemaProps
		expectedError error
	}{
		"do nothing when property config is nil": {
			request: &PropertyProcessorRequest{
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Required:   nil,
				Properties: map[string]apiextensions.JSONSchemaProps{},
			},
			expectedError: nil,
		},
		"do nothing when read only filter is false": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadOnly: false,
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Required:   nil,
				Properties: map[string]apiextensions.JSONSchemaProps{},
			},
			expectedError: nil,
		},
		"do nothing when schema is read only": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadOnly: true,
					},
				},
				OpenAPISchema: &openapi3.Schema{
					ReadOnly: true,
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Required:   nil,
				Properties: map[string]apiextensions.JSONSchemaProps{},
			},
			expectedError: nil,
		},
		"remove non-read-only properties from required list": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadOnly: true,
					},
				},
				OpenAPISchema: &openapi3.Schema{
					ReadOnly: false,
					Required: []string{"a", "b", "c"},
					Properties: map[string]*openapi3.SchemaRef{
						"a": {Value: &openapi3.Schema{ReadOnly: true}},
						"b": {Value: &openapi3.Schema{ReadOnly: false}},
						"c": {Value: &openapi3.Schema{ReadOnly: true}},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: nil,
			expectedError: nil,
		},
		"do nothing when path is root": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadOnly: true,
					},
				},
				OpenAPISchema: &openapi3.Schema{
					ReadOnly: false,
					Required: []string{"a", "b", "c"},
					Properties: map[string]*openapi3.SchemaRef{
						"a": {Value: &openapi3.Schema{ReadOnly: true}},
						"b": {Value: &openapi3.Schema{ReadOnly: false}},
						"c": {Value: &openapi3.Schema{ReadOnly: true}},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
				Path: []string{"$"},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Required:   []string{"a", "c"},
				Properties: map[string]apiextensions.JSONSchemaProps{},
			},
			expectedError: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &ReadOnlyProperties{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedProps, test.request.Property)
		})
	}
}
