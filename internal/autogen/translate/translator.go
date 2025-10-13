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

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/unstructured"
)

// PtrClientObj is a pointer type implementing client.Object
type PtrClientObj[T any] interface {
	*T
	client.Object
}

// Translator allows to translate back and forth between a CRD schema
// and SDK API structures of a certain version.
// A translator is an immutable configuration object, it can be safely shared
// across threads
type Translator struct {
	crd          CRDInfo
	majorVersion string
}

// APIImporter can translate itself into Kubernetes Objects.
// Use to customize or accelerate translations ad-hoc
type APIImporter[T any, P PtrClientObj[T]] interface {
	FromAPI(t *Translator, target P) ([]client.Object, error)
}

// APIExporter can translate itself to an API Object.
// Use to customize or accelerate translations ad-hoc
type APIExporter[T any] interface {
	ToAPI(t *Translator, target *T) error
}

// NewTranslator creates a translator for a particular CRD and major version pairs,
// and with a particular set of known Kubernetes object dependencies.
//
// Given the following example resource:
//
//	apiVersion: atlas.generated.mongodb.com/v1
//	kind: SearchIndex
//	metadata:
//	  name: search-index
//	spec:
//	  v20250312:
//
// In the above case crdVersion is "v1" and majorVersion is "v20250312".
func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, majorVersion string) *Translator {
	return &Translator{
		crd:          CRDInfo{definition: crd, version: crdVersion},
		majorVersion: majorVersion,
	}
}

// FromAPI translaters a source API structure into a Kubernetes object.
func FromAPI[S any, T any, P PtrClientObj[T]](t *Translator, target P, source *S, objs ...client.Object) ([]client.Object, error) {
	importer, ok := any(source).(APIImporter[T, P])
	if ok {
		return importer.FromAPI(t, target)
	}
	sourceUnstructured, err := unstructured.ToUnstructured(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert API source value to unstructured: %w", err)
	}

	targetUnstructured := map[string]any{}

	versionedSpec := map[string]any{}
	unstructured.CopyFields(versionedSpec, sourceUnstructured)
	if err := unstructured.CreateField(targetUnstructured, versionedSpec, "spec", t.majorVersion); err != nil {
		return nil, fmt.Errorf("failed to create versioned spec object in unstructured target: %w", err)
	}
	versionedSpecEntry := map[string]any{}
	unstructured.CopyFields(versionedSpecEntry, sourceUnstructured)
	versionedSpec["entry"] = versionedSpecEntry

	versionedStatus := map[string]any{}
	unstructured.CopyFields(versionedStatus, sourceUnstructured)
	if err := unstructured.CreateField(targetUnstructured, versionedStatus, "status", t.majorVersion); err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in unsstructured target: %w", err)
	}

	extraObjects, err := ExpandMappings(t, targetUnstructured, target, objs...)
	if err != nil {
		return nil, fmt.Errorf("failed to process API mappings: %w", err)
	}
	if err := unstructured.FromUnstructured(target, targetUnstructured); err != nil {
		return nil, fmt.Errorf("failed set structured kubernetes object from unstructured: %w", err)
	}
	return append([]client.Object{target}, extraObjects...), nil
}

// ToAPI translates a source Kubernetes spec into a target API structure
func ToAPI[T any](t *Translator, target *T, source client.Object, objs ...client.Object) error {
	exporter, ok := (source).(APIExporter[T])
	if ok {
		return exporter.ToAPI(t, target)
	}
	specVersion := selectVersion(&t.crd.definition.Spec, t.crd.version)
	kind := t.crd.definition.Spec.Names.Kind
	props, err := getOpenAPIProperties(kind, specVersion)
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD schema properties: %w", err)
	}
	specProps, err := getSpecPropertiesFor(kind, props, "spec")
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD spec properties: %w", err)
	}
	if _, ok := specProps[t.majorVersion]; !ok {
		return fmt.Errorf("failed to match the CRD spec version %q in schema", t.majorVersion)
	}
	unstructuredSrc, err := unstructured.ToUnstructured(source)
	if err != nil {
		return fmt.Errorf("failed to convert k8s source value to unstructured: %w", err)
	}
	targetUnstructured := map[string]any{}
	value, err := unstructured.AccessField[map[string]any](unstructuredSrc, "spec", t.majorVersion)
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}

	if err := CollapseMappings(t, value, source, objs...); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	rawEntry := value["entry"]
	if entry, ok := rawEntry.(map[string]any); ok {
		unstructured.CopyFields(targetUnstructured, unstructured.SkipKeys(value, "entry"))
		entryPathInTarget := findEntryPathInTarget(targetType)
		dst := targetUnstructured
		if len(entryPathInTarget) > 0 {
			newValue := map[string]any{}
			if err = unstructured.CreateField(targetUnstructured, newValue, entryPathInTarget...); err != nil {
				return fmt.Errorf("failed to set target copy destination to path %v: %w", entryPathInTarget, err)
			}
			dst = newValue
		}
		unstructured.CopyFields(dst, entry)
	} else {
		unstructured.CopyFields(targetUnstructured, value)
	}
	delete(targetUnstructured, "groupref")
	if err := unstructured.FromUnstructured(target, targetUnstructured); err != nil {
		return fmt.Errorf("failed to set structured value from unstructured: %w", err)
	}
	return nil
}
