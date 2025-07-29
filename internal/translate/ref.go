package translate

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type DecoderFunc func(any) (any, error)

var decoders = map[string]DecoderFunc{
	"v1/secrets": func(in any) (any, error) {
		return base64Decode((in).(string))
	},
}

type refMapping struct {
	XKubernetesMapping kubeMapping    `json:"x-kubernetes-mapping"`
	XOpenAPIMapping    openAPIMapping `json:"x-openapi-mapping"`
}

type kubeMapping struct {
	NameSelector      string   `json:"nameSelector"`
	PropertySelectors []string `json:"propertySelectors"`
	Properties        []string `json:"properties"`
	Type              kubeType `json:"type"`
}

type kubeType struct {
	Kind     string `json:"kind"`
	Group    string `json:"group,omitempty"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
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

func (km kubeMapping) FetchReferencedValue(target string, reference map[string]any, deps DependencyFinder) (any, error) {
	refPath := km.NameSelector
	if refPath == "" {
		return nil, errors.New("cannot solve reference without a x-kubernetes-mapping.nameSelector")
	}
	refName, err := accessField[string](reference, asPath(refPath)...)
	if err != nil {
		return nil, fmt.Errorf("failed to access field %q at %v: %w", refPath, reference, err)
	}
	resource := deps.Find(refName, SetFallbackNamespace)
	if resource == nil {
		return nil, fmt.Errorf("failed to find Kubernetes resource %q: %w", refName, err)
	}
	gvk := resource.GetObjectKind().GroupVersionKind()
	if km.Type.Kind != "" && !km.Equal(gvk) {
		return nil, fmt.Errorf("resource %q had to be a %q but got %q", refName, km.GVK(), gvk)
	}
	resourceMap, err := toUnstructured(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to turn resource %q into an unestuctued map: %w", refName, err)
	}
	value, err := km.fetchFromProperties(resourceMap)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("failed to resolve reference properties: %w", err)
	}
	if err == ErrNotFound {
		var err error
		value, err = km.fetchFromPropertySelectors(resourceMap, target)
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("failed to resolve reference properties or property selectors: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to resolve reference property selectors: %w", err)
		}
	}
	return km.Decode(value)
}

func (km kubeMapping) Decode(value any) (any, error) {
	decode := decoders[km.GVR()]
	if decode != nil {
		return decode(value)
	}
	return value, nil
}

func (km kubeMapping) fetchFromProperties(resource map[string]any) (any, error) {
	for _, prop := range km.Properties {
		path := resolveXPath(prop)
		value, err := accessField[any](resource, path...)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to access property as %v: %w", path, err)
		}
		return value, nil
	}
	return nil, ErrNotFound
}

func (km kubeMapping) fetchFromPropertySelectors(resource map[string]any, target string) (any, error) {
	for _, selector := range km.PropertySelectors {
		prop := selector
		if strings.HasSuffix(prop, ".#") {
			prop = fmt.Sprintf("%s.%s", prop[:len(prop)-2], target)
		}
		path := resolveXPath(prop)
		value, err := accessField[any](resource, path...)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to access selected property as %v: %w", path, err)
		}
		return value, nil
	}
	return nil, ErrNotFound
}

type openAPIMapping struct {
	Property string `json:"property"`
	Type     string `json:"type"`
}

func (oam openAPIMapping) TargetPath() []string {
	return resolveXPath(oam.Property)
}

func isReference(obj map[string]any) bool {
	return obj["x-kubernetes-mapping"] != nil && obj["x-openapi-mapping"] != nil
}

func processReference(path []string, mapping, spec map[string]any, deps DependencyFinder) error {
	reference, err := accessField[map[string]any](spec, base(path))
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}
	if len(reference) == 0 {
		return nil
	}
	refMap := refMapping{}
	if err := fromUnstructured(&refMap, mapping); err != nil {
		return fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	targetPath := refMap.XOpenAPIMapping.TargetPath()
	key, ok := reference["key"].(string)
	if !ok || key == "" {
		key = base(targetPath)
	}
	value, err := refMap.XKubernetesMapping.FetchReferencedValue(key, reference, deps)
	return createField(spec, value, targetPath...)
}

func resolveXPath(xpath string) []string {
	if strings.HasPrefix(xpath, "$.") {
		return asPath(xpath[1:])
	}
	return asPath(xpath)
}
