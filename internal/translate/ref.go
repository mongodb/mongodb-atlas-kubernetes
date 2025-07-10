package translate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type refMapping struct {
	XKubernetesMapping kubeMapping    `json:"x-kubernetes-mapping"`
	XOpenAPIMapping    openAPIMapping `json:"x-openapi-mapping"`
}

type kubeMapping struct {
	GVR              string `json:"gvr"`
	NameSelector     string `json:"nameSelector"`
	PropertySelector string `json:"propertySelector"`
}

type openAPIMapping struct {
	Property string `json:"property"`
	Type     string `json:"type"`
}

func isReference(obj map[string]any) bool {
	return obj["x-kubernetes-mapping"] != nil && obj["x-openapi-mapping"] != nil
}

func processReference(path []string, mapping, spec map[string]any, deps ...client.Object) error {
	reference, err := accessField[map[string]any](spec, base(path))
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}
	refMap := refMapping{}
	if err := toStructured(&refMap, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	if refMap.XOpenAPIMapping.Type != "string" {
		return fmt.Errorf("unsupported referenced value type %v", refMap.XOpenAPIMapping.Type)
	}

	return processSecretReference(path, &refMap, reference, spec, deps...)
}

func solveReferencedDependency(path []string, reference map[string]any, refMap *refMapping, gvk schema.GroupVersionKind, deps ...client.Object) (map[string]any, error) {
	referenceValue, err := accessField[string](reference, asPath(refMap.XKubernetesMapping.NameSelector)...)
	if err != nil {
		return nil, fmt.Errorf("failed accessing reference value for mapping at %v: %w", path, err)
	}
	dep := findReferencedDep(deps, gvk, referenceValue)
	if dep == nil {
		return nil, fmt.Errorf("kubernetes dependency of type %q not found with name %q",
			refMap.XKubernetesMapping.GVR, referenceValue)
	}

	depUnstructured, err := toUnstructured(dep)
	if err != nil {
		return nil, fmt.Errorf("failed to translate referenced kubernetes type %q to unstructured: %w",
			refMap.XKubernetesMapping.GVR, err)
	}
	return depUnstructured, nil
}

func findReferencedDep(deps []client.Object, gvk schema.GroupVersionKind, name string) client.Object {
	for _, dep := range deps {
		if equalGroupVersionKind(gvk, dep.GetObjectKind().GroupVersionKind()) && dep.GetName() == name {
			return dep
		}
	}
	return nil
}

func equalGroupVersionKind(a, b schema.GroupVersionKind) bool {
	return a.Group == b.Group && a.Version == b.Version && a.Kind == b.Kind
}
