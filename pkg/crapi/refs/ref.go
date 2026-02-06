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
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/objmap"
)

const (
	xKubeMappingKey        = "x-kubernetes-mapping"
	xOpenAPIMappingKey     = "x-openapi-mapping"
	refName                = "name"
	refKey                 = "key"
	propertySelectorSuffix = ".#"
)

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
func ExpandAll(scheme *runtime.Scheme, mappings []*Mapping, main client.Object, deps []client.Object, cr map[string]any) ([]client.Object, error) {
	ks := newKubeset(scheme, main, deps)
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
func CollapseAll(scheme *runtime.Scheme, mappings []*Mapping, main client.Object, deps []client.Object, req map[string]any) error {
	ks := newKubeset(scheme, main, deps)
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
	if err := objmap.FromObjectMap(&oam, oamMap); err != nil {
		return nil, fmt.Errorf("failed to parse a reference mapping: %w", err)
	}

	kmMap, ok := kubeExt.(map[string]any)
	if !ok {
		return nil,
			fmt.Errorf("failed to coerce Kubernetes mapping extension type, expected map[string]any got %T", kubeExt)
	}
	km := KubeMapping{}
	if err := objmap.FromObjectMap(&km, kmMap); err != nil {
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
	holder, err := objmap.GetFieldObject(obj, collapsedPath...)
	if errors.Is(err, objmap.ErrNotFound) {
		return nil
	}
	rawValue, found := holder[objmap.Base(collapsedPath)]
	if !found {
		return fmt.Errorf("failed to access rawValue at %v: %w",
			collapsedPath, objmap.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed accessing value at path %v: %w", mapping.path, err)
	}

	// Check if any existing dependency matches the value in its properties
	existingDep, err := mapping.findMatchingDependency(ks, rawValue)
	if err != nil {
		return fmt.Errorf("failed to find matching dependency: %w", err)
	}
	if existingDep != nil {
		refData := map[string]any{refName: existingDep.GetName()}
		holder[objmap.Base(mapping.path)] = refData
		// Remove the original API field (e.g., groupId) since we're using the reference
		delete(holder, objmap.Base(collapsedPath))
		return nil
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
		refData[refKey] = objmap.Base(path)
	}
	holder[objmap.Base(mapping.path)] = refData
	ks.add(dep)
	return nil
}

// findMatchingDependency searches through existing dependencies to find one
// whose property value matches the given raw value from the API.
// This is used to map API fields (like groupId) to existing Kubernetes resources
// (like a Group) that were passed as dependencies.
func (mapping *Mapping) findMatchingDependency(ks *kubeset, rawValue any) (client.Object, error) {
	for _, dep := range ks.m {
		// Check if this dependency's GVK matches the mapping's expected type
		gvk := dep.GetObjectKind().GroupVersionKind()
		if gvk.Kind == "" || gvk.GroupVersion().String() == "" {
			gvks, _, err := ks.scheme.ObjectKinds(dep)
			// Skip if type is not registered in the scheme (expected case for unknown types)
			// or if no GVKs are returned. This is a best-effort search through dependencies.
			if runtime.IsNotRegisteredError(err) || len(gvks) == 0 {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("failed to get ObjectKinds for dependency %q: %w", dep.GetName(), err)
			}
			gvk = gvks[0]
		}
		if !mapping.xKubernetesMapping.equal(gvk) {
			continue
		}

		// If no properties are defined, we can't match by property value - skip this dependency
		// (some mappings use propertySelectors instead of properties)
		if len(mapping.xKubernetesMapping.Properties) == 0 {
			continue
		}

		// Convert the dependency to an object map to access its properties
		depMap, err := objmap.ToObjectMap(dep)
		if err != nil {
			return nil, fmt.Errorf("failed to convert dependency %q to object map: %w", dep.GetName(), err)
		}

		// Check if any of the mapping's properties match the raw value
		for _, prop := range mapping.xKubernetesMapping.Properties {
			path := resolveXPath(prop)
			value, err := objmap.GetField[any](depMap, path...)
			// Skip if the property path doesn't exist in this dependency (expected case)
			if errors.Is(err, objmap.ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("failed to get field %v from dependency %q: %w", path, dep.GetName(), err)
			}
			// Compare values - handle both string and other types
			if valuesMatch(value, rawValue) {
				return dep, nil
			}
		}
	}
	return nil, nil
}

// valuesMatch compares two values for equality, handling type conversions
func valuesMatch(a, b any) bool {
	// Direct equality check
	if a == b {
		return true
	}
	// Handle string comparisons with pointer types
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && bIsStr {
		return aStr == bStr
	}
	// Handle pointer to string
	if aPtr, ok := a.(*string); ok && aPtr != nil && bIsStr {
		return *aPtr == bStr
	}
	if bPtr, ok := b.(*string); ok && bPtr != nil && aIsStr {
		return aStr == *bPtr
	}
	return false
}

// collapse processes the unstructured (API request) object at the given path to
// follow a reference and extract its value onto the object at the expected path
func (mapping *Mapping) collapse(ks *kubeset, obj map[string]any) error {
	holder, err := objmap.GetFieldObject(obj, mapping.path...)
	if errors.Is(err, objmap.ErrNotFound) {
		return nil
	}
	rawValue, found := holder[objmap.Base(mapping.path)]
	if !found {
		return fmt.Errorf("failed to access rawValue at %v: %w",
			mapping.path, objmap.ErrNotFound)
	}
	if errors.Is(err, objmap.ErrNotFound) {
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
		key = objmap.Base(targetPath)
	}
	value, err := mapping.xKubernetesMapping.fetchReferencedValue(ks, key, reference)
	if err != nil {
		return fmt.Errorf("failed to fetch referenced value %s: %w", key, err)
	}
	holder[objmap.Base(targetPath)] = value
	return nil
}

func (mapping *Mapping) collapsedPath() []string {
	path := make([]string, len(mapping.path))
	copy(path, mapping.path)
	return append(objmap.Dir(path), objmap.AsPath(strings.TrimLeft(mapping.xOpenAPIMapping.Property, "$"))...)
}

func (mapping *Mapping) nameFor(prefix string, path []string) string {
	if path[0] == "entry" {
		path = path[1:]
	}
	return prefixedName(prefix, path[0], path[1:]...)
}
