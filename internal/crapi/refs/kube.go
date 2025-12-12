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

package refs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/unstructured"
)

// ErrNoMatchingPropertySelector when no property selectors matched the value
var ErrNoMatchingPropertySelector = errors.New("no matching property selector found to set value")

type KubeMapping struct {
	NameSelector      string   `json:"nameSelector"`
	PropertySelectors []string `json:"propertySelectors"`
	Properties        []string `json:"properties"`
	Type              KubeType `json:"type"`
}

type KubeType struct {
	Kind     string `json:"kind"`
	Group    string `json:"group,omitempty"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
}

func (km KubeMapping) gvk() string {
	if km.Type.Group == "" {
		return fmt.Sprintf("%s, Kind=%s", km.Type.Version, km.Type.Kind)
	}
	return fmt.Sprintf("%s/%s, Kind=%s", km.Type.Group, km.Type.Version, km.Type.Kind)
}

func (km KubeMapping) gvr() string {
	if km.Type.Group == "" {
		return fmt.Sprintf("%s/%s", km.Type.Version, km.Type.Resource)
	}
	return fmt.Sprintf("%s/%s/%s", km.Type.Group, km.Type.Version, km.Type.Resource)
}

func (km KubeMapping) equal(gvk schema.GroupVersionKind) bool {
	return km.Type.Group == gvk.Group && km.Type.Version == gvk.Version && km.Type.Kind == gvk.Kind
}

func (km KubeMapping) fetchReferencedValue(mc *kubeset, target string, reference map[string]any) (any, error) {
	refPath := km.NameSelector
	if refPath == "" {
		return nil, fmt.Errorf("cannot solve reference without a %s.nameSelector", xKubeMappingKey)
	}
	refName, err := unstructured.GetField[string](reference, unstructured.AsPath(refPath)...)
	if err != nil {
		return nil, fmt.Errorf("failed to access field %q at %v: %w", refPath, reference, err)
	}
	resource := mc.find(refName)
	if resource == nil {
		return nil, fmt.Errorf("failed to find Kubernetes resource %q: %w", refName, err)
	}
	gvk := resource.GetObjectKind().GroupVersionKind()
	if gvk.Kind == "" || gvk.GroupVersion().String() == "" {
		gvks, _, err := mc.scheme.ObjectKinds(resource)
		if err != nil || len(gvks) == 0 {
			return nil, fmt.Errorf("failed to infer GroupVersionKind for resource %q from scheme: %w", refName, err)
		}
		gvk = gvks[0]
	}
	if km.Type.Kind != "" && !km.equal(gvk) {
		return nil, fmt.Errorf("resource %q had to be a %q but got %q", refName, km.gvk(), gvk)
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
	return km.decode(refSolver, value)
}

func (km KubeMapping) decode(refSolver *resolver, value any) (any, error) {
	decode := refSolver.decoders[km.gvr()]
	if decode != nil {
		return decode(value)
	}
	return value, nil
}

func (km KubeMapping) encode(refSolver *resolver, value any) (any, error) {
	encode := refSolver.encoders[km.gvr()]
	if encode != nil {
		return encode(value)
	}
	return value, nil
}

func (km KubeMapping) fetchFromProperties(resource map[string]any) (any, error) {
	for _, prop := range km.Properties {
		path := resolveXPath(prop)
		value, err := unstructured.GetField[any](resource, path...)
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

func (km KubeMapping) fetchFromPropertySelectors(resource map[string]any, target string) (any, error) {
	for _, selector := range km.PropertySelectors {
		prop := selector
		if strings.HasSuffix(prop, propertySelectorSuffix) {
			prop = fmt.Sprintf("%s.%s", prop[:len(prop)-2], target)
		}
		path := resolveXPath(prop)
		value, err := unstructured.GetField[any](resource, path...)
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

func (km KubeMapping) setAtPropertySelectors(refSolver *resolver, gvr string, obj map[string]any, target string, value any) (client.Object, error) {
	for _, selector := range km.PropertySelectors {
		prop := selector
		if strings.HasSuffix(prop, propertySelectorSuffix) {
			targetPath := resolveXPath(target)
			prop = fmt.Sprintf("%s.%s", prop[:len(prop)-2], unstructured.Base(targetPath))
		}
		path := resolveXPath(prop)
		if err := unstructured.RecursiveCreateField(obj, value, path...); err != nil {
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
		valueCopy, err := unstructured.GetField[any](unstructuredCopy, path...)
		if reflect.DeepEqual(value, valueCopy) {
			return obj, nil
		}
		if err != nil && !errors.Is(err, unstructured.ErrNotFound) {
			return nil, fmt.Errorf("failed to check Kubernetes object contents: %w", err)
		}
	}
	return nil, ErrNoMatchingPropertySelector
}
