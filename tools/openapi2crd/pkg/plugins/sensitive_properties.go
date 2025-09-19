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

	"github.com/getkin/kin-openapi/openapi3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type SensitiveProperties struct {
	NoOp
}

func NewSensitivePropertiesPlugin() *SensitiveProperties {
	return &SensitiveProperties{}
}

func (s *SensitiveProperties) Name() string {
	return "sensitive_properties"
}

func (n *SensitiveProperties) ProcessProperty(g Generator, propertyConfig *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps {
	if !isSensitiveField(path, propertyConfig) {
		return props
	}

	props.ID = path[len(path)-1] + "SecretRef"

	if extensionsSchema.Value.Extensions == nil {
		extensionsSchema.Value.Extensions = map[string]interface{}{}
	}

	extensionsSchema.Value.Extensions["x-kubernetes-mapping"] = map[string]interface{}{
		"type": map[string]interface{}{
			"kind":     "Secret",
			"resource": v1.ResourceSecrets,
			"version":  "v1",
		},
		"nameSelector":      ".name",
		"propertySelectors": []string{"$.data.#"},
	}

	extensionsSchema.Value.Extensions["x-openapi-mapping"] = map[string]interface{}{
		"property": "." + path[len(path)-1],
		"type":     propertySchema.Type,
	}

	props.Type = "object"
	props.Description = fmt.Sprintf("SENSITIVE FIELD\n\nReference to a secret containing data for the %q field:\n\n%v", path[len(path)-1], propertySchema.Description)
	defaultKey := apiextensions.JSON(".data." + path[len(path)-1])
	props.Properties = map[string]apiextensions.JSONSchemaProps{
		"name": {
			Type:        "string",
			Description: `Name of the secret containing the sensitive field value.`,
		},
		"key": {
			Type:        "string",
			Default:     &defaultKey,
			Description: fmt.Sprintf(`Key of the secret data containing the sensitive field value, defaults to %q.`, path[len(path)-1]),
		},
	}

	return props
}

func isSensitiveField(path []string, mapping *configv1alpha1.PropertyMapping) bool {
	if mapping == nil {
		return false
	}

	p := jsonPath(path)

	for _, sensitiveField := range mapping.Filters.SensitiveProperties {
		if sensitiveField == p {
			return true
		}
	}

	return false
}
