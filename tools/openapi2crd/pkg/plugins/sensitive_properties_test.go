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
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestSensitivePropertyName(t *testing.T) {
	p := &SensitiveProperties{}
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
				Required:    []string{"name"},
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
				Required:    []string{"name"},
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
			p := &SensitiveProperties{}
			err := p.Process(test.request)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedProps, test.request.Property)
			assert.Equal(t, test.expectedExtensions, test.request.ExtensionsSchema)
		})
	}
}
