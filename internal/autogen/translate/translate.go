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

package translate

import (
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/refs"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/unstructured"
)

const (
	APIMAppingsAnnotation = "api-mappings"
)

// Translator allows to translate back and forth between a CRD schema
// and SDK API structures of a certain version.
// A translator is an immutable configuration object, it can be safely shared
// across threads
type Translator interface {
	// MajorVersion returns the pinned SDK major version
	MajorVersion() string

	// Annotation gives access to the value of the given annotation on the translated schema
	Annotation(annotationKey string) string

	// Validate checks the given unsttructured object complies with the translated schema
	Validate(unstructuredObj map[string]any) error
}

// Request holds common parameters for all translation request
type Request struct {
	Translator   Translator
	Log          logr.Logger
	Dependencies []client.Object
}

// APIImporter can translate itself into Kubernetes Objects.
// Use to customize or accelerate translations ad-hoc
type APIImporter[T any, P refs.PtrClientObj[T]] interface {
	FromAPI(tr *Request, target P) ([]client.Object, error)
}

// APIExporter can translate itself to an API Object.
// Use to customize or accelerate translations ad-hoc
type APIExporter[T any] interface {
	ToAPI(tr *Request, target *T) error
}

// ToAPI translates a source Kubernetes spec into a target API structure.
// It uses the spec only to populate ethe API request, nothing from the status.
func ToAPI[T any](r *Request, target *T, source client.Object) error {
	exporter, ok := (source).(APIExporter[T])
	if ok {
		return exporter.ToAPI(r, target)
	}
	unstructuredSrc, err := unstructured.ToUnstructured(source)
	if err != nil {
		return fmt.Errorf("failed to convert k8s source value to unstructured: %w", err)
	}
	if err := r.Translator.Validate(unstructuredSrc); err != nil {
		return fmt.Errorf("failed to validate unstructured object input: %w", err)
	}
	targetUnstructured := map[string]any{}

	if err := collapseReferences(r, unstructuredSrc, source); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	value, err := unstructured.GetField[map[string]any](unstructuredSrc, "spec", r.Translator.MajorVersion())
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
func FromAPI[S any, T any, P refs.PtrClientObj[T]](r *Request, target P, source *S) ([]client.Object, error) {
	importer, ok := any(source).(APIImporter[T, P])
	if ok {
		return importer.FromAPI(r, target)
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
		targetUnstructured, map[string]any{}, "spec", r.Translator.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to ensure versioned spec object in unstructured target: %w", err)
	}
	versionedStatus, err := unstructured.GetOrCreateField(
		targetUnstructured, map[string]any{}, "status", r.Translator.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in unstructured target: %w", err)
	}

	unstructured.CopyFields(versionedSpec, sourceUnstructured)
	versionedSpecEntry := map[string]any{}
	unstructured.CopyFields(versionedSpecEntry, sourceUnstructured)
	versionedSpec["entry"] = versionedSpecEntry
	unstructured.CopyFields(versionedStatus, sourceUnstructured)

	extraObjects, err := expandReferences(r, targetUnstructured, target)
	if err != nil {
		return nil, fmt.Errorf("failed to process API mappings: %w", err)
	}
	if err := r.Translator.Validate(targetUnstructured); err != nil {
		return nil, fmt.Errorf("failed to validate unstructured object output: %w", err)
	}
	if err := unstructured.FromUnstructured(target, targetUnstructured); err != nil {
		return nil, fmt.Errorf("failed set structured kubernetes object from unstructured: %w", err)
	}
	return append([]client.Object{target}, extraObjects...), nil
}

// collapseReferences finds all Kubernetes references, solves them and collapses
// them by replacing their values from the reference (e.g kubernetes secret or
// group), into the corresponding API field
func collapseReferences(r *Request, obj map[string]any, main client.Object) error {
	h := refs.NewHandler(main, r.Dependencies)
	mappingsYML := r.Translator.Annotation(APIMAppingsAnnotation)
	if mappingsYML == "" {
		return nil
	}
	mappings := map[string]any{}
	if err := yaml.Unmarshal([]byte(mappingsYML), mappings); err != nil {
		return fmt.Errorf("failed to unmarshal mappings YAML: %w", err)
	}

	return h.CollapseReferences(obj, mappings, "spec", r.Translator.MajorVersion())
}

// expandReferences finds all API fields that must be referenced, and expand
// such reference, moving the value (e.g. sensitive field or group id) to a
// referenced Kubernetes object (e.g. kubernets secret or Atlas Group)
func expandReferences(r *Request, obj map[string]any, main client.Object) ([]client.Object, error) {
	h := refs.NewHandler(main, r.Dependencies)
	mappingsYML := r.Translator.Annotation(APIMAppingsAnnotation)
	if mappingsYML == "" {
		return []client.Object{}, nil
	}
	mappings := map[string]any{}
	if err := yaml.Unmarshal([]byte(mappingsYML), mappings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mappings YAML: %w", err)
	}

	majorVersion := r.Translator.MajorVersion()
	for _, entry := range []struct {
		title string
		path  []string
	}{
		{title: "spec", path: []string{"spec", majorVersion}},
		{title: "spec entry", path: []string{"spec", majorVersion, "entry"}},
		{title: "status", path: []string{"status", majorVersion}},
	} {
		if err := h.ExpandReferences(obj, mappings, entry.path...); err != nil {
			return nil, fmt.Errorf("failed to map properties of %q from API to Kubernetes: %w", entry.title, err)
		}
	}
	return h.Added(), nil
}
