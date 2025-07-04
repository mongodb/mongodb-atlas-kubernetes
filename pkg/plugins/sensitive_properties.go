package plugins

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"strings"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
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

func (n *SensitiveProperties) ProcessProperty(g Generator, mapping *configv1alpha1.FieldMapping, props *apiextensions.JSONSchemaProps, propertySchema, extensionsSchema *openapi3.Schema, path ...string) {
	if !isSensitiveField(path, mapping) {
		return
	}

	if extensionsSchema.Extensions == nil {
		extensionsSchema.Extensions = map[string]interface{}{}
	}

	extensionsSchema.Extensions["x-kubernetes-mapping"] = map[string]interface{}{
		"gvr":              "secrets/v1",
		"nameSelector":     ".name",
		"propertySelector": ".key",
	}

	extensionsSchema.Extensions["x-openapi-mapping"] = map[string]interface{}{
		"property": "." + path[len(path)-1],
		"type":     propertySchema.Type,
	}

	props.Type = "object"
	props.Description = fmt.Sprintf("SENSITIVE FIELD\n\nReference to a secret containing data for the %q field:\n\n%v", path[len(path)-1], propertySchema.Description)
	defaultKey := apiextensions.JSON(".data." + path[len(path)-1])
	props.Properties = map[string]apiextensions.JSONSchemaProps{
		"name": {
			Type:        "string",
			Description: fmt.Sprintf(`Name of the secret containing the sensitive field value.`),
		},
		"key": {
			Type:        "string",
			Default:     &defaultKey,
			Description: fmt.Sprintf(`Key of the secret data containing the sensitive field value, defaults to %q.`, path[len(path)-1]),
		},
	}

	return
}

func (s *SensitiveProperties) ProcessPropertyName(mapping *configv1alpha1.FieldMapping, path []string) string {
	if isSensitiveField(path, mapping) {
		return path[len(path)-1] + "SecretRef"
	}

	return path[len(path)-1]
}

func isSensitiveField(path []string, mapping *configv1alpha1.FieldMapping) bool {
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

func jsonPath(path []string) string {
	result := strings.Join(path, ".")
	return strings.ReplaceAll(result, ".[*]", "[*]")
}
