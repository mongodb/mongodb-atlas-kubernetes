package plugins

import (
	"fmt"
	"strings"

	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type SensitiveProperty struct{}

func (p *SensitiveProperty) Name() string {
	return "sensitive_property"
}

func (p *SensitiveProperty) Process(req *PropertyProcessorRequest) error {
	if !isSensitiveField(req.Path, req.PropertyConfig) {
		return nil
	}

	req.Property.ID = req.Path[len(req.Path)-1] + "SecretRef"

	if req.ExtensionsSchema.Value.Extensions == nil {
		req.ExtensionsSchema.Value.Extensions = map[string]interface{}{}
	}

	req.ExtensionsSchema.Value.Extensions["x-kubernetes-mapping"] = map[string]interface{}{
		"type": map[string]interface{}{
			"kind":     "Secret",
			"resource": v1.ResourceSecrets,
			"version":  "v1",
		},
		"nameSelector":      ".name",
		"propertySelectors": []string{"$.data.#"},
	}

	req.ExtensionsSchema.Value.Extensions["x-openapi-mapping"] = map[string]interface{}{
		"property": "." + req.Path[len(req.Path)-1],
		"type":     req.OpenAPISchema.Type,
	}

	req.Property.Type = "object"
	req.Property.Description = fmt.Sprintf("SENSITIVE FIELD\n\nReference to a secret containing data for the %q field:\n\n%v", req.Path[len(req.Path)-1], req.OpenAPISchema.Description)
	defaultKey := apiextensions.JSON(".data." + req.Path[len(req.Path)-1])
	req.Property.Properties = map[string]apiextensions.JSONSchemaProps{
		"name": {
			Type:        "string",
			Description: `Name of the secret containing the sensitive field value.`,
		},
		"key": {
			Type:        "string",
			Default:     &defaultKey,
			Description: fmt.Sprintf(`Key of the secret data containing the sensitive field value, defaults to %q.`, req.Path[len(req.Path)-1]),
		},
	}

	return nil
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

func jsonPath(path []string) string {
	result := strings.Join(path, ".")
	return strings.ReplaceAll(result, ".[*]", "[*]")
}
