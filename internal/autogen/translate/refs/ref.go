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

package refs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/samples/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/unstructured"
)

const (
	xKubeMappingKey        = "x-kubernetes-mapping"
	xOpenAPIMappingKey     = "x-openapi-mapping"
	refName                = "name"
	refKey                 = "key"
	propertySelectorSuffix = ".#"
)

// PtrClientObj is a pointer type implementing client.Object
type PtrClientObj[T any] interface {
	*T
	client.Object
}

var ErrNoMatchingPropertySelector = errors.New("no matching property selector found to set value")

type EncodeDecodeFunc func(any) (any, error)

type referenceResolver struct {
	decoders           map[string]EncodeDecodeFunc
	encoders           map[string]EncodeDecodeFunc
	optionalExpansions map[string]struct{} // A set is better for lookups
	kubeObjectRegistry map[string]func(map[string]any) (client.Object, error)
}

func newReferenceResolver() *referenceResolver {
	return &referenceResolver{

		decoders: map[string]EncodeDecodeFunc{
			"v1/secrets": func(in any) (any, error) {
				s, ok := in.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string for secret decoding, but got %T", in)
				}
				return secretDecode(s)
			},
		},

		encoders: map[string]EncodeDecodeFunc{
			"v1/secrets": func(in any) (any, error) {
				s, ok := in.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string for secret encoding, but got %T", in)
				}
				return secretEncode(s), nil
			},
		},

		optionalExpansions: map[string]struct{}{"groupRef": {}},

		kubeObjectRegistry: map[string]func(map[string]any) (client.Object, error){
			"v1/secrets":                            newKubeObjectFactory[corev1.Secret](),
			"atlas.generated.mongodb.com/v1/groups": newKubeObjectFactory[v1.Group](),
		},
	}
}

func newKubeObjectFactory[T any, P PtrClientObj[T]]() func(map[string]any) (client.Object, error) {
	return func(unstructured map[string]any) (client.Object, error) {
		obj := new(T)
		initializedObj, err := initObject(obj, unstructured)
		if err != nil {
			return nil, err
		}
		// Assert the concrete pointer type (*P) to the interface type.
		// This is guaranteed to be safe because of the PtrClientObj constraint
		return any(initializedObj).(client.Object), nil
	}
}

func initObject[T any](obj *T, unstructuredObj map[string]any) (*T, error) {
	if unstructuredObj != nil {
		if err := unstructured.FromUnstructured(obj, unstructuredObj); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

type refMapping struct {
	XKubernetesMapping kubeMapping    `json:"x-kubernetes-mapping"`
	XOpenAPIMapping    openAPIMapping `json:"x-openapi-mapping"`
}

type namedRef struct {
	*refMapping
	name string
}

func isReference(obj map[string]any) bool {
	return obj[xKubeMappingKey] != nil && obj[xOpenAPIMappingKey] != nil
}

func newRef(name string, rm *refMapping) *namedRef {
	return &namedRef{name: name, refMapping: rm}
}

func (ref *namedRef) Expand(mc *context, pathHint []string, obj map[string]any) error {
	path := ref.pathToExpand(pathHint)
	rawValue, err := unstructured.AccessField[any](obj, unstructured.Base(path))
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}
	refSolver := newReferenceResolver()
	value, err := ref.XKubernetesMapping.Encode(refSolver, rawValue)
	if err != nil {
		return fmt.Errorf("failed to encode value at path %v: %w", path, err)
	}
	gvr := ref.XKubernetesMapping.GVR()
	depUnstructured, err := unstructuredKubeObjectFor(refSolver, gvr)
	if err != nil {
		return fmt.Errorf("failed to populate unstructured dependency: %w", err)
	}
	dep, err := ref.XKubernetesMapping.setAtPropertySelectors(
		refSolver, gvr, depUnstructured, ref.XOpenAPIMapping.Property, value)
	if err != nil {
		if _, ok := refSolver.optionalExpansions[ref.name]; ok {
			return nil
		}
		return fmt.Errorf("failed to populate final dependency object: %w", err)
	}
	name := ref.Name(mc.main.GetName(), path)
	if mc.has(name) {
		return nil
	}
	dep.SetName(name)
	dep.SetNamespace(mc.main.GetNamespace())
	refData := map[string]any{refName: dep.GetName()}
	if ref.XOpenAPIMapping.Property != "" {
		path := resolveXPath(ref.XOpenAPIMapping.Property)
		refData[refKey] = unstructured.Base(path)
	}
	obj[ref.name] = refData
	mc.add(dep)
	return nil
}

// The pathHint points to the reference field itself (e.g., ["spec", "passwordSecretRef"]).
// We need to find the raw value in the OpenAPI object, which is at a different path
// (e.g., ["spec", "password"]). This function replaces the reference field name
// with the target OpenAPI property name to get the correct path for expansion.
func (ref *namedRef) pathToExpand(pathHint []string) []string {
	path := make([]string, len(pathHint))
	copy(path, pathHint)
	path[len(path)-1] = unstructured.Base(resolveXPath(ref.XOpenAPIMapping.Property))
	return path
}

func (ref *namedRef) Name(prefix string, path []string) string {
	if path[0] == "entry" {
		path = path[1:]
	}
	return PrefixedName(prefix, path[0], path[1:]...)
}

func (ref *namedRef) Collapse(mc *context, path []string, obj map[string]any) error {
	reference, err := unstructured.AccessField[map[string]any](obj, unstructured.Base(path))
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", path, err)
	}
	if len(reference) == 0 {
		return nil
	}

	targetPath := ref.XOpenAPIMapping.TargetPath()
	key, ok := reference[refKey].(string)
	if !ok || key == "" {
		key = unstructured.Base(targetPath)
	}
	value, err := ref.XKubernetesMapping.FetchReferencedValue(mc, key, reference)
	if err != nil {
		return fmt.Errorf("failed to fetch referenced value %s: %w", key, err)
	}
	return unstructured.CreateField(obj, value, targetPath...)
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

func (km kubeMapping) FetchReferencedValue(mc *context, target string, reference map[string]any) (any, error) {
	refPath := km.NameSelector
	if refPath == "" {
		return nil, fmt.Errorf("cannot solve reference without a %s.nameSelector", xKubeMappingKey)
	}
	refName, err := unstructured.AccessField[string](reference, unstructured.AsPath(refPath)...)
	if err != nil {
		return nil, fmt.Errorf("failed to access field %q at %v: %w", refPath, reference, err)
	}
	resource := mc.find(refName)
	if resource == nil {
		return nil, fmt.Errorf("failed to find Kubernetes resource %q: %w", refName, err)
	}
	gvk := resource.GetObjectKind().GroupVersionKind()
	if km.Type.Kind != "" && !km.Equal(gvk) {
		return nil, fmt.Errorf("resource %q had to be a %q but got %q", refName, km.GVK(), gvk)
	}
	resourceMap, err := unstructured.ToUnstructured(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to turn resource %q into an unstuctued map: %w", refName, err)
	}
	value, err := km.fetchFromProperties(resourceMap)
	if err != nil && !errors.Is(err, unstructured.ErrNotFound) {
		return nil, fmt.Errorf("failed to resolve reference properties: %w", err)
	}
	if errors.Is(err, unstructured.ErrNotFound) {
		var err error
		value, err = km.fetchFromPropertySelectors(resourceMap, target)
		if errors.Is(err, unstructured.ErrNotFound) {
			return nil, fmt.Errorf("failed to resolve reference properties or property selectors: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to resolve reference property selectors: %w", err)
		}
	}
	refSolver := newReferenceResolver()
	return km.Decode(refSolver, value)
}

func (km kubeMapping) Decode(refSolver *referenceResolver, value any) (any, error) {
	decode := refSolver.decoders[km.GVR()]
	if decode != nil {
		return decode(value)
	}
	return value, nil
}

func (km kubeMapping) Encode(refSolver *referenceResolver, value any) (any, error) {
	encode := refSolver.encoders[km.GVR()]
	if encode != nil {
		return encode(value)
	}
	return value, nil
}

func (km kubeMapping) fetchFromProperties(resource map[string]any) (any, error) {
	for _, prop := range km.Properties {
		path := resolveXPath(prop)
		value, err := unstructured.AccessField[any](resource, path...)
		if errors.Is(err, unstructured.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to access property as %v: %w", path, err)
		}
		return value, nil
	}
	return nil, unstructured.ErrNotFound
}

func (km kubeMapping) fetchFromPropertySelectors(resource map[string]any, target string) (any, error) {
	for _, selector := range km.PropertySelectors {
		prop := selector
		if strings.HasSuffix(prop, propertySelectorSuffix) {
			prop = fmt.Sprintf("%s.%s", prop[:len(prop)-2], target)
		}
		path := resolveXPath(prop)
		value, err := unstructured.AccessField[any](resource, path...)
		if errors.Is(err, unstructured.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to access selected property as %v: %w", path, err)
		}
		return value, nil
	}
	return nil, unstructured.ErrNotFound
}

func (km kubeMapping) setAtPropertySelectors(refSolver *referenceResolver, gvr string, obj map[string]any, target string, value any) (client.Object, error) {
	for _, selector := range km.PropertySelectors {
		prop := selector
		if strings.HasSuffix(prop, propertySelectorSuffix) {
			targetPath := resolveXPath(target)
			prop = fmt.Sprintf("%s.%s", prop[:len(prop)-2], unstructured.Base(targetPath))
		}
		path := resolveXPath(prop)
		if err := unstructured.CreateField(obj, value, path...); err != nil {
			return nil, fmt.Errorf("failed to set value at %q: %w", path, err)
		}
		obj, err := initializedKubeObjectFor(refSolver, gvr, obj)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Kubernetes object: %w", err)
		}
		unstructuredCopy, err := unstructured.ToUnstructured(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to read Kubernetes object contents: %w", err)
		}
		valueCopy, err := unstructured.AccessField[any](unstructuredCopy, path...)
		if reflect.DeepEqual(value, valueCopy) {
			return obj, nil
		}
		if err != nil && !errors.Is(err, unstructured.ErrNotFound) {
			return nil, fmt.Errorf("failed to check Kubernetes object contents: %w", err)
		}
	}
	return nil, ErrNoMatchingPropertySelector
}

type openAPIMapping struct {
	Property string `json:"property"`
	Type     string `json:"type"`
}

func (oam openAPIMapping) TargetPath() []string {
	return resolveXPath(oam.Property)
}

func resolveXPath(xpath string) []string {
	if strings.HasPrefix(xpath, "$.") {
		return unstructured.AsPath(xpath[1:])
	}
	return unstructured.AsPath(xpath)
}

func unstructuredKubeObjectFor(refSolver *referenceResolver, gvr string) (map[string]any, error) {
	objCopy, err := kubeObjectFor(refSolver, gvr)
	if err != nil {
		return nil, fmt.Errorf("failed to get unstructured kube object for GVR %q: %w", gvr, err)
	}
	return unstructured.ToUnstructured(objCopy)
}

func kubeObjectFor(refSolver *referenceResolver, gvr string) (client.Object, error) {
	return initializedKubeObjectFor(refSolver, gvr, nil)
}

func initializedKubeObjectFor(refSolver *referenceResolver, gvr string, unstructuredData map[string]any) (client.Object, error) {
	objFn, ok := refSolver.kubeObjectRegistry[gvr]
	if !ok {
		return nil, fmt.Errorf("unsupported kube object for GVR %q", gvr)
	}
	return objFn(unstructuredData)
}
