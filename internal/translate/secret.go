package translate

import (
	"encoding/base64"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func processSecretReference(path []string, refMap *refMapping, reference, spec map[string]any, deps ...client.Object) error {
	if refMap.XKubernetesMapping.GVR != "secrets/v1" {
		return fmt.Errorf("unsupported GVR %q", refMap.XKubernetesMapping.GVR)
	}
	dep, err := solveSecretReferencedDependency(path, reference, refMap, deps...)
	if err != nil {
		return fmt.Errorf("failed solving referenced kubernetes dependency: %w", err)
	}
	value, err := fetchReferencedSecretValue(refMap, dep)
	if err != nil {
		return fmt.Errorf("failed fetching referenced value: %w", err)
	}
	property := base(asPath(refMap.XOpenAPIMapping.Property))
	spec[property] = value
	return nil
}

func solveSecretReferencedDependency(path []string, reference map[string]any, refMap *refMapping, deps ...client.Object) (map[string]any, error) {
	secretGVK := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}
	return solveReferencedDependency(path, reference, refMap, secretGVK, deps...)
}

func fetchReferencedSecretValue(refMap *refMapping, dep map[string]any) (string, error) {
	if refMap.XKubernetesMapping.PropertySelector != SecretProperySelector {
		return "", fmt.Errorf("unsupported property selector for secret value: %v",
			refMap.XKubernetesMapping.PropertySelector)
	}
	propertyPath := asPath(refMap.XOpenAPIMapping.Property)
	propertyPath = append([]string{"data"}, propertyPath...)
	value, err := accessField[string](dep, propertyPath...)
	if err != nil {
		return "", fmt.Errorf("failed to access referenced value at %v: %w", propertyPath, err)
	}
	decodedValue, err := base64Decode(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret base64 encoding at %v: %w", propertyPath, err)
	}
	return decodedValue, nil
}

func base64Decode(value string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}
	return string(bytes), nil
}