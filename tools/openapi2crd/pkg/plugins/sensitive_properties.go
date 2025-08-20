package plugins

import (
	"fmt"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/processor"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type SensitiveProperties struct{}

func (p *SensitiveProperties) Name() string {
	return "sensitive_properties"
}

func (p *SensitiveProperties) Process(input processor.Input) error {
	i, ok := input.(processor.PropertyInput)
	if !ok {
		return nil
	}
	propertyConfig := i.PropertyConfig
	props := i.KubeSchema
	propertySchema := i.OpenAPISchema
	extensionsSchema := i.ExtensionsSchema
	path := i.Path

	if !isSensitiveField(path, propertyConfig) {
		return nil
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

	return nil
}

func NewSensitivePropertiesPlugin() *SensitiveProperties {
	return &SensitiveProperties{}
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
