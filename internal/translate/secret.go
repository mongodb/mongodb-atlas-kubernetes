package translate

import (
	"encoding/base64"
	"fmt"
)

func processSecretReference(path []string, namespace string, refMap *refMapping, reference, spec map[string]any, deps DependencyFinder) error {
	if refMap.XKubernetesMapping.GVR() != "v1/secrets" {
		return fmt.Errorf("unsupported GVR %q", refMap.XKubernetesMapping.GVR())
	}
	dep, err := solveReferencedDependency(path, namespace, reference, refMap, deps)
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

func fetchReferencedSecretValue(refMap *refMapping, dep map[string]any) (string, error) {
	if !in(refMap.XKubernetesMapping.PropertySelectors, SecretProperySelector) {
		return "", fmt.Errorf("unsupported property selectors for secret value: %v",
			refMap.XKubernetesMapping.PropertySelectors)
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
