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
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate/unstructured"
)

const (
	xKubeMappingKey        = "x-kubernetes-mapping"
	xOpenAPIMappingKey     = "x-openapi-mapping"
	refName                = "name"
	refKey                 = "key"
	propertySelectorSuffix = ".#"
)

// PtrClientObj is a pointer type implementing client.Object.
// It represents a Kubernetes object for the controller-runtime library.
type PtrClientObj[T any] interface {
	*T
	client.Object
}

// Mapping encodes a reference mapping consisting of two extensions,
// a x-kubernetes-mapping and a x-openapi-mapping
type Mapping struct {
	path               []string
	propertyName       string
	xKubernetesMapping KubeMapping
	xOpenAPIMapping    OpenAPIMapping
}

// NewMapping creates a reference mapping from its basic information: path,
// property name, kubernetes and openapi extension mappings
func NewMapping(path []string, propertyName string, kubeMapping KubeMapping, openAPIMapping OpenAPIMapping) *Mapping {
	return &Mapping{
		path:               path,
		propertyName:       propertyName,
		xKubernetesMapping: kubeMapping,
		xOpenAPIMapping:    openAPIMapping,
	}
}

// FindMappings lists all reference mappings from the given schema at a path
func FindMappings(schema *openapi3.Schema, path []string) ([]*Mapping, error) {
	mappings := []*Mapping{}
	return findMappings(mappings, schema, path)
}

// ExpandAll takes a list of mappings and expands the ones found in the given
// CR unstructured object. CR corresponds to the main typed object and deps
// are other kubernetes objects related with such main object
func ExpandAll(mappings []*Mapping, main client.Object, deps []client.Object, cr map[string]any) ([]client.Object, error) {
	ks := newKubeset(main, deps)
	for _, mapping := range mappings {
		if err := mapping.expand(ks, cr); err != nil {
			return nil, fmt.Errorf("failed to expand reference for %q: %w", mapping.propertyName, err)
		}
	}
	return ks.added, nil
}

// CollapseAll takes a list of mappings and collapses the ones found in the
// given request unstructured object. The SDK request must map to the main typed
//
//	object and deps are other kubernetes objects related with such main object
func CollapseAll(mappings []*Mapping, main client.Object, deps []client.Object, req map[string]any) error {
	ks := newKubeset(main, deps)
	for _, mapping := range mappings {
		if err := mapping.collapse(ks, req); err != nil {
			return fmt.Errorf("failed to expand reference for %q: %w", mapping.propertyName, err)
		}
	}
	return nil
}

func findMappings(mappings []*Mapping, schema *openapi3.Schema, path []string) ([]*Mapping, error) {
	if schema == nil {
		return mappings, nil
	}

	mapping, err := extractMapping(schema, path)
	if err != nil {
		return nil, fmt.Errorf("failed to extract a mapping at %v: %w", path, err)
	}
	if mapping != nil {
		mappings = append(mappings, mapping)
		return mappings, nil
	}

	keys := []string{}
	for propName := range schema.Properties {
		keys = append(keys, propName)
	}
	sort.Strings(keys)
	for _, propName := range keys {
		propSchemaRef := schema.Properties[propName]
		if propSchemaRef != nil && propSchemaRef.Value != nil {
			var err error
			mappings, err = findMappings(mappings, propSchemaRef.Value, append(path, propName))
			if err != nil {
				return nil, fmt.Errorf("failed to parse a object property %q: %w", propName, err)
			}
		}
	}

	if schema.Items != nil && schema.Items.Value != nil {
		var err error
		mappings, err = findMappings(mappings, schema.Items.Value, append(path, "[]"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse a items value: %w", err)
		}
	}

	return mappings, nil
}

func extractMapping(schema *openapi3.Schema, path []string) (*Mapping, error) {
	oamExt, hasOpenAPIMapping := schema.Extensions[xOpenAPIMappingKey]
	if !hasOpenAPIMapping {
		return nil, nil
	}
	kubeExt, hashKubeMapping := schema.Extensions[xKubeMappingKey]
	if !hashKubeMapping {
		return nil, nil
	}

	oamMap, ok := oamExt.(map[string]any)
	if !ok {
		return nil,
			fmt.Errorf("failed to coerce open api extension type, expected map[string]any got %T", oamExt)
	}
	oam := OpenAPIMapping{}
	if err := unstructured.FromUnstructured(&oam, oamMap); err != nil {
		return nil, fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	kmMap, ok := kubeExt.(map[string]any)
	if !ok {
		return nil,
			fmt.Errorf("failed to coerce Kubernetes mapping extension type, expected map[string]any got %T", kubeExt)
	}
	km := KubeMapping{}
	if err := unstructured.FromUnstructured(&km, kmMap); err != nil {
		return nil, fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	propertyName := path[len(path)-1]
	currentPath := make([]string, len(path))
	copy(currentPath, path)
	return NewMapping(currentPath, propertyName, km, oam), nil
}

// expand proceses the unstructured (Kubernetes CR) object at the given path to
// create a reference and insert it as a new Kubernetes Object in the kubeset
func (mapping *Mapping) expand(ks *kubeset, obj map[string]any) error {
	collapsedPath := mapping.collapsedPath()
	holder, err := unstructured.GetFieldObject(obj, collapsedPath...)
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	rawValue, found := holder[unstructured.Base(collapsedPath)]
	if !found {
		return fmt.Errorf("failed to access rawValue at %v: %w",
			collapsedPath, unstructured.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", mapping.path, err)
	}
	refSolver := newReferenceResolver()
	value, err := mapping.xKubernetesMapping.encode(refSolver, rawValue)
	if err != nil {
		return fmt.Errorf("failed to encode value at path %v: %w", mapping.path, err)
	}
	gvr := mapping.xKubernetesMapping.gvr()
	depUnstructured, err := unstructuredKubeObjectFor(refSolver, gvr)
	if err != nil {
		return fmt.Errorf("failed to populate unstructured dependency: %w", err)
	}
	dep, err := mapping.xKubernetesMapping.setAtPropertySelectors(
		refSolver, gvr, depUnstructured, mapping.xOpenAPIMapping.Property, value)
	if err != nil {
		if _, ok := refSolver.optionalExpansions[mapping.propertyName]; ok {
			return nil
		}
		return fmt.Errorf("failed to populate final dependency object: %w", err)
	}
	name := mapping.nameFor(ks.main.GetName(), mapping.path)
	if ks.has(name) {
		return nil
	}
	dep.SetName(name)
	dep.SetNamespace(ks.main.GetNamespace())
	refData := map[string]any{refName: dep.GetName()}
	if mapping.xOpenAPIMapping.Property != "" {
		path := resolveXPath(mapping.xOpenAPIMapping.Property)
		refData[refKey] = unstructured.Base(path)
	}
	holder[unstructured.Base(mapping.path)] = refData
	ks.add(dep)
	return nil
}

// collapse processes the unstructured (API request) object at the given path to
// follow a reference and extract its value onto the object at the expected path
func (mapping *Mapping) collapse(ks *kubeset, obj map[string]any) error {
	holder, err := unstructured.GetFieldObject(obj, mapping.path...)
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	rawValue, found := holder[unstructured.Base(mapping.path)]
	if !found {
		return fmt.Errorf("failed to access rawValue at %v: %w",
			mapping.path, unstructured.ErrNotFound)
	}
	if errors.Is(err, unstructured.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", mapping.path, err)
	}
	reference, isObj := (rawValue).(map[string]any)
	if !isObj || len(reference) == 0 {
		return nil
	}
	targetPath := mapping.xOpenAPIMapping.targetPath()
	key, ok := reference[refKey].(string)
	if !ok || key == "" {
		key = unstructured.Base(targetPath)
	}
	value, err := mapping.xKubernetesMapping.fetchReferencedValue(ks, key, reference)
	if err != nil {
		return fmt.Errorf("failed to fetch referenced value %s: %w", key, err)
	}
	holder[unstructured.Base(targetPath)] = value
	return nil
}

func (mapping *Mapping) collapsedPath() []string {
	path := make([]string, len(mapping.path))
	copy(path, mapping.path)
	return append(unstructured.Dir(path), unstructured.AsPath(mapping.xOpenAPIMapping.Property)...)
}

func (mapping *Mapping) nameFor(prefix string, path []string) string {
	if path[0] == "entry" {
		path = path[1:]
	}
	return prefixedName(prefix, path[0], path[1:]...)
}
