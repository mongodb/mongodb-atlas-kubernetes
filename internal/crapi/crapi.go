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

package crapi

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/unstructured"
)

// Translator allows to translate back and forth between a CRD schema
// and SDK API structures of a certain version.
// A translator is an immutable configuration object, it can be safely shared
// across goroutines
type Translator interface {
	// MajorVersion returns the pinned SDK major version
	MajorVersion() string

	// Mappings returns all the OpenAPi custom reference extensions, or an error
	Mappings() ([]*refs.Mapping, error)
}

// Request is deprecated do not use
//
// Deprecated: request is no longer used in the ToAPI and FromAPI calls
type Request struct {
	Translator   Translator
	Dependencies []client.Object
}

// APIImporter can translate itself into Kubernetes Objects.
// Use to customize or accelerate translations ad-hoc
type APIImporter[T any, P refs.PtrClientObj[T]] interface {
	FromAPI(tr Translator, target P, objs ...client.Object) ([]client.Object, error)
}

// APIExporter can translate itself to an API Object.
// Use to customize or accelerate translations ad-hoc
type APIExporter[T any] interface {
	ToAPI(tr Translator, target *T, objs ...client.Object) error
}

// ToAPI translates a source Kubernetes spec into a target API structure.
// It uses the spec only to populate ethe API request, nothing from the status.
// The target is set to a API request struct to be filled.
// The source is set to the Kubernetes CR value. Only the spec data is used here.
// The request includes the translator and the dependencies associated with the
// source CR, usually Kubernetes secrets.
func ToAPI[T any](tr Translator, target *T, source client.Object, objs ...client.Object) error {
	exporter, ok := (source).(APIExporter[T])
	if ok {
		return exporter.ToAPI(tr, target)
	}
	unstructuredSrc, err := unstructured.ToUnstructured(source)
	if err != nil {
		return fmt.Errorf("failed to convert k8s source value to unstructured: %w", err)
	}
	targetUnstructured := map[string]any{}

	if err := collapseReferences(tr, unstructuredSrc, source, objs); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	value, err := unstructured.GetField[map[string]any](unstructuredSrc, "spec", tr.MajorVersion())
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}
	unstructured.CopyFields(targetUnstructured, value)
	rawEntry := value["entry"]
	if entry, ok := rawEntry.(map[string]any); ok {
		unstructured.CopyFields(targetUnstructured, entry)
	}
	if err := unstructured.FromUnstructured(target, targetUnstructured); err != nil {
		return fmt.Errorf("failed to set structured value from unstructured: %w", err)
	}
	return nil
}

// FromAPI translates a source API structure into a Kubernetes object.
// The API source is used to populate the Kubernetes spec, including the
// spec.entry and status as well.
// The target is set to CR value to be filled. Both spec and stuatus are filled.
// The source is set to API response.
// The request includes the translator and any dependencies associated with the
// source CR.
// Returns any extra objects extracted from the response as separate Kubernetes
// objects, such as Kubernetes secrets, for instance. This list does not include
// the mutated target, and will be empty if nothing else was extracted off the ç
// response.
func FromAPI[S any, T any, P refs.PtrClientObj[T]](tr Translator, target P, source *S, objs ...client.Object) ([]client.Object, error) {
	importer, ok := any(source).(APIImporter[T, P])
	if ok {
		return importer.FromAPI(tr, target, objs...)
	}
	sourceUnstructured, err := unstructured.ToUnstructured(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert API source value to unstructured: %w", err)
	}

	targetUnstructured, err := unstructured.ToUnstructured(target)
	if err != nil {
		return nil, fmt.Errorf("failed to convert target value to unstructured: %w", err)
	}

	versionedSpec, err := unstructured.GetOrCreateField(
		targetUnstructured, map[string]any{}, "spec", tr.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to ensure versioned spec object in unstructured target: %w", err)
	}
	versionedStatus, err := unstructured.GetOrCreateField(
		targetUnstructured, map[string]any{}, "status", tr.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in unstructured target: %w", err)
	}

	unstructured.CopyFields(versionedSpec, sourceUnstructured)
	versionedSpecEntry := map[string]any{}
	unstructured.CopyFields(versionedSpecEntry, sourceUnstructured)
	versionedSpec["entry"] = versionedSpecEntry
	unstructured.CopyFields(versionedStatus, sourceUnstructured)

	extraObjects, err := expandReferences(tr, targetUnstructured, target, objs)
	if err != nil {
		return nil, fmt.Errorf("failed to process API mappings: %w", err)
	}
	if err := unstructured.FromUnstructured(target, targetUnstructured); err != nil {
		return nil, fmt.Errorf("failed set structured kubernetes object from unstructured: %w", err)
	}
	return extraObjects, nil
}

// collapseReferences finds all Kubernetes references, solves them and collapses
// them by replacing their values from the reference (e.g Kubernetes secret or
// group), into the corresponding API request value
func collapseReferences(tr Translator, req map[string]any, main client.Object, objs []client.Object) error {
	mappings, err := tr.Mappings()
	if err != nil {
		return fmt.Errorf("failed to extract mappings to collapse: %w", err)
	}
	return refs.CollapseAll(mappings, main, objs, req)
}

// expandReferences finds all API fields that must be referenced, and expand
// such reference, moving the value (e.g. sensitive field or group id) to a
// referenced Kubernetes object (e.g. Kubernetes secret or Atlas Group)
func expandReferences(tr Translator, cr map[string]any, main client.Object, objs []client.Object) ([]client.Object, error) {
	mappings, err := tr.Mappings()
	if err != nil {
		return nil, fmt.Errorf("failed to extract mappings to expand: %w", err)
	}
	return refs.ExpandAll(mappings, main, objs, cr)
}

func propertyValueOrNil(schema *openapi3.Schema, propertyName string) *openapi3.Schema {
	if schema.Properties != nil &&
		schema.Properties[propertyName] != nil && schema.Properties[propertyName].Value != nil {
		return schema.Properties[propertyName].Value
	}
	return nil
}
