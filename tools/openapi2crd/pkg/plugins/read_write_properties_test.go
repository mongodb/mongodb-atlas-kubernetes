package plugins

import (
	"testing"
	"tools/openapi2crd/pkg/apis/config/v1alpha1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestReadWritePropertyName(t *testing.T) {
	p := &ReadWriteProperties{}
	assert.Equal(t, "read_write_property", p.Name())
}

func TestReadWritePropertyProcess(t *testing.T) {
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
		"do nothing when read write only filter is false": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadWriteOnly: false,
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
		"remove entire property when schema is read only": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadWriteOnly: true,
					},
				},
				OpenAPISchema: &openapi3.Schema{
					ReadOnly: true,
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   []string{"a", "b"},
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: nil,
			expectedError: nil,
		},
		"remove read only properties from required list and keep others": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &v1alpha1.PropertyMapping{
					Filters: v1alpha1.Filters{
						ReadWriteOnly: true,
					},
				},
				OpenAPISchema: &openapi3.Schema{
					Required: []string{"a", "b", "c"},
					Properties: map[string]*openapi3.SchemaRef{
						"a": {Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
						"b": {Value: &openapi3.Schema{Type: &openapi3.Types{"string"}, ReadOnly: true}},
						"c": {Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required: []string{"a", "b", "c"},
					Properties: map[string]apiextensions.JSONSchemaProps{
						"a": {Type: "string"},
						"b": {Type: "string"},
						"c": {Type: "string"},
					},
				},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Required: []string{"a", "c"},
				Properties: map[string]apiextensions.JSONSchemaProps{
					"a": {Type: "string"},
					"b": {Type: "string"},
					"c": {Type: "string"},
				},
			},
			expectedError: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &ReadWriteProperties{}
			err := p.Process(tc.request)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedProps, tc.request.Property)
		})
	}
}
