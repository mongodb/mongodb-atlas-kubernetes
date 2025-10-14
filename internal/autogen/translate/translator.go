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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/santhosh-tekuri/jsonschema/v5"
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
	majorVersion string
	jsonSchema   *jsonschema.Schema
	annotations  map[string]string
}

// Request holds common parameters for all translation request
type Request struct {
	Translator   *Translator
	Log          logr.Logger
	Dependencies []client.Object
}

// APIImporter can translate itself into Kubernetes Objects.
// Use to customize or accelerate translations ad-hoc
type APIImporter[T any, P PtrClientObj[T]] interface {
	FromAPI(tr *Request, target P) ([]client.Object, error)
}

// APIExporter can translate itself to an API Object.
// Use to customize or accelerate translations ad-hoc
type APIExporter[T any] interface {
	ToAPI(tr *Request, target *T) error
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
func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, majorVersion string) (*Translator, error) {
	specVersion := selectVersion(&crd.Spec, crdVersion)
	if err := assertMajorVersion(specVersion, crd.Spec.Names.Kind, majorVersion); err != nil {
		return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", majorVersion, err)
	}
	schema, err := compileCRDSchema(specVersion.Schema.OpenAPIV3Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return &Translator{
		majorVersion: majorVersion,
		jsonSchema:   schema,
		annotations:  crd.Annotations,
	}, nil
}

// ToAPI translates a source Kubernetes spec into a target API structure
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
	value, err := unstructured.AccessField[map[string]any](unstructuredSrc, "spec", r.Translator.majorVersion)
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}

	if err := CollapseMappings(r, value, source); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	rawEntry := value["entry"]
	if entry, ok := rawEntry.(map[string]any); ok {
		unstructured.CopyFields(targetUnstructured, unstructured.SkipKeys(value, "entry"))
		entryPathInTarget := []string{}
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

// FromAPI translates a source API structure into a Kubernetes object.
func FromAPI[S any, T any, P PtrClientObj[T]](r *Request, target P, source *S) ([]client.Object, error) {
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

	versionedSpec, err := unstructured.AccessOrCreateField(
		targetUnstructured, map[string]any{}, "spec", r.Translator.majorVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure versioned spec object in unstructured target: %w", err)
	}
	unstructured.CopyFields(versionedSpec, sourceUnstructured)

	versionedSpecEntry := map[string]any{}
	unstructured.CopyFields(versionedSpecEntry, sourceUnstructured)
	versionedSpec["entry"] = versionedSpecEntry

	versionedStatus := map[string]any{}
	unstructured.CopyFields(versionedStatus, sourceUnstructured)
	if err := unstructured.CreateField(targetUnstructured, versionedStatus, "status", r.Translator.majorVersion); err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in unsstructured target: %w", err)
	}

	extraObjects, err := ExpandMappings(r, targetUnstructured, target)
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

// Validate checks whether or not an unstructured object value conforms to the
// translator's CRD schema
func (t *Translator) Validate(unstructuredObj map[string]any) error {
	if err := t.jsonSchema.Validate(unstructuredObj); err != nil {
		return fmt.Errorf("object validation failed against CRD schema: %w", err)
	}
	return nil
}

func assertMajorVersion(specVersion *apiextensionsv1.CustomResourceDefinitionVersion, kind string, majorVersion string) error {
	props, err := getOpenAPIProperties(kind, specVersion)
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD schema properties: %w", err)
	}
	specProps, err := getSpecPropertiesFor(kind, props, "spec")
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD spec properties: %w", err)
	}
	_, ok := specProps[majorVersion]
	if !ok {
		return fmt.Errorf("failed to match the CRD spec version %q in schema", majorVersion)
	}
	return nil
}

func compileCRDSchema(openAPISchema *apiextensionsv1.JSONSchemaProps) (*jsonschema.Schema, error) {
	schemaBytes, err := json.Marshal(openAPISchema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CRD schema to JSON: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(schemaBytes)); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return schema, nil
}
