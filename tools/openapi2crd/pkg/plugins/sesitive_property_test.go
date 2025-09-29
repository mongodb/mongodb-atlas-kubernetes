package plugins

import (
	"testing"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestSensitivePropertyName(t *testing.T) {
	p := &SensitiveProperty{}
	assert.Equal(t, "sensitive_property", p.Name())
}

func TestSensitivePropertyProcess(t *testing.T) {
	stringJson := apiextensions.JSON(".data.password")

	tests := map[string]struct {
		request            *PropertyProcessorRequest
		expectedProps      *apiextensions.JSONSchemaProps
		expectedExtensions *openapi3.SchemaRef
		expectedError      error
	}{
		"do nothing when property config is nil": {
			request: &PropertyProcessorRequest{
				Property: &apiextensions.JSONSchemaProps{
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				Properties: map[string]apiextensions.JSONSchemaProps{},
			},
			expectedError: nil,
		},
		"convert sensitive property": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Filters: configv1alpha1.Filters{
						SensitiveProperties: []string{"$.password"},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
				ExtensionsSchema: &openapi3.SchemaRef{Value: &openapi3.Schema{}},
				OpenAPISchema:    &openapi3.Schema{Type: &openapi3.Types{"string"}, Description: "the password"},
				Path:             []string{"$", "password"},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				ID:          "passwordSecretRef",
				Type:        "object",
				Description: "SENSITIVE FIELD\n\nReference to a secret containing data for the \"password\" field:\n\nthe password",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type:        "string",
						Description: `Name of the secret containing the sensitive field value.`,
					},
					"key": {
						Type:        "string",
						Default:     &stringJson,
						Description: `Key of the secret data containing the sensitive field value, defaults to "password".`,
					},
				},
			},
			expectedExtensions: &openapi3.SchemaRef{Value: &openapi3.Schema{
				Extensions: map[string]interface{}{
					"x-kubernetes-mapping": map[string]interface{}{
						"type": map[string]interface{}{
							"kind":     "Secret",
							"resource": v1.ResourceSecrets,
							"version":  "v1",
						},
						"nameSelector":      ".name",
						"propertySelectors": []string{"$.data.#"},
					},
					"x-openapi-mapping": map[string]interface{}{
						"property": ".password",
						"type":     &openapi3.Types{"string"},
					},
				},
			}},
			expectedError: nil,
		},
		"convert sensitive property in nested object": {
			request: &PropertyProcessorRequest{
				PropertyConfig: &configv1alpha1.PropertyMapping{
					Filters: configv1alpha1.Filters{
						SensitiveProperties: []string{"$.credentials[*].password"},
					},
				},
				Property: &apiextensions.JSONSchemaProps{
					Required:   nil,
					Properties: map[string]apiextensions.JSONSchemaProps{},
				},
				ExtensionsSchema: &openapi3.SchemaRef{Value: &openapi3.Schema{}},
				OpenAPISchema:    &openapi3.Schema{Type: &openapi3.Types{"string"}, Description: "the credentials password"},
				Path:             []string{"$", "credentials[*]", "password"},
			},
			expectedProps: &apiextensions.JSONSchemaProps{
				ID:          "passwordSecretRef",
				Type:        "object",
				Description: "SENSITIVE FIELD\n\nReference to a secret containing data for the \"password\" field:\n\nthe credentials password",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {
						Type:        "string",
						Description: `Name of the secret containing the sensitive field value.`,
					},
					"key": {
						Type:        "string",
						Default:     &stringJson,
						Description: `Key of the secret data containing the sensitive field value, defaults to "password".`,
					},
				},
			},
			expectedExtensions: &openapi3.SchemaRef{Value: &openapi3.Schema{
				Extensions: map[string]interface{}{
					"x-kubernetes-mapping": map[string]interface{}{
						"type": map[string]interface{}{
							"kind":     "Secret",
							"resource": v1.ResourceSecrets,
							"version":  "v1",
						},
						"nameSelector":      ".name",
						"propertySelectors": []string{"$.data.#"},
					},
					"x-openapi-mapping": map[string]interface{}{
						"property": ".password",
						"type":     &openapi3.Types{"string"},
					},
				},
			}},
			expectedError: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &SensitiveProperty{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedProps, test.request.Property)
			assert.Equal(t, test.expectedExtensions, test.request.ExtensionsSchema)
		})
	}
}
