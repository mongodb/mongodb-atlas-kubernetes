package translate

import (
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type refMapping struct {
	XKubernetesMapping kubeMapping    `json:"x-kubernetes-mapping"`
	XOpenAPIMapping    openAPIMapping `json:"x-openapi-mapping"`
}

type kubeMapping struct {
	NameSelector      string   `json:"nameSelector"`
	PropertySelectors []string `json:"propertySelectors"`
	Type              kubeType `json:"type"`
}

type kubeType struct {
	Kind     string `json:"kind"`
	Group    string `json:"group,omitempty"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
}

type openAPIMapping struct {
	Property string `json:"property"`
	Type     string `json:"type"`
}

func (km kubeMapping) GVK() string {
	if km.Type.Group == "" {
		return fmt.Sprintf("%s, Kind=%s", km.Type.Version, km.Type.Kind)
	}
	return fmt.Sprintf("%s/%s, Kind=%s", km.Type.Group, km.Type.Version, km.Type.Kind)
}

func (km kubeMapping) GVR() string {
	if km.Type.Group == "" {
		return fmt.Sprintf("%s/%s", km.Type.Version, km.Type.Resource)
	}
	return fmt.Sprintf("%s/%s/%s", km.Type.Group, km.Type.Version, km.Type.Resource)
}

func (km kubeMapping) Equal(gvk schema.GroupVersionKind) bool {
	return km.Type.Group == gvk.Group && km.Type.Version == gvk.Version && km.Type.Kind == gvk.Kind
}

func isReference(obj map[string]any) bool {
	return obj["x-kubernetes-mapping"] != nil && obj["x-openapi-mapping"] != nil
}

func processReference(path []string, namespace string, mapping, spec map[string]any, deps DependencyFinder) error {
	reference, err := accessField[map[string]any](spec, base(path))
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}
	if len(reference) == 0 {
		return nil
	}
	refMap := refMapping{}
	if err := toStructured(&refMap, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	if refMap.XKubernetesMapping.GVR() == "atlas.generated.mongodb.com/v1/groups" {
		// TODO: implement group refs
		return nil
	}

	if refMap.XOpenAPIMapping.Type != "string" {
		return fmt.Errorf("unsupported referenced value type %q (refMap=%v)",
			refMap.XOpenAPIMapping.Type, refMap)
	}

	return processSecretReference(path, namespace, &refMap, reference, spec, deps)
}

func solveReferencedDependency(path []string, namespace string, reference map[string]any, refMap *refMapping, deps DependencyFinder) (map[string]any, error) {
	referenceValue, err := accessField[string](reference, asPath(refMap.XKubernetesMapping.NameSelector)...)
	if err != nil {
		return nil, fmt.Errorf("failed accessing reference value for mapping at %v: %w", path, err)
	}
	dep := findReferencedDep(deps, &refMap.XKubernetesMapping, referenceValue, namespace)
	if dep == nil {
		return nil, fmt.Errorf("kubernetes dependency of type %q not found with name %q",
			refMap.XKubernetesMapping.GVK(), referenceValue)
	}

	depUnstructured, err := toUnstructured(dep)
	if err != nil {
		return nil, fmt.Errorf("failed to translate referenced kubernetes type %q to unstructured: %w",
			refMap.XKubernetesMapping.GVK(), err)
	}
	return depUnstructured, nil
}

func findReferencedDep(deps DependencyFinder, kubeMap *kubeMapping, name, namespace string) client.Object {
	dep := deps.Find(name, namespace)
	if dep != nil && kubeMap.Equal(dep.GetObjectKind().GroupVersionKind()) {
		return dep
	}
	log.Printf("NOT FOUND %q at %v", name, deps)
	return nil
}
