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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/unstructured"
)

// translator implements Translator to translate from a given CRD to and from
// a given SDK version using the same upstream OpenAPI schema
type translator struct {
	scheme        *runtime.Scheme
	majorVersion  string
	mappingSchema *openapi3.SchemaRef
}

func (tr *translator) ToAPI(target any, source client.Object, objs ...client.Object) error {
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

func (tr *translator) FromAPI(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
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

// Annotation returns the annotation value from the pinned translated schema
// CRD version
// FindAllMappings traverses the mapping schema and returns all found extensions.
func (tr *translator) Mappings() ([]*refs.Mapping, error) {
	noMappings := []*refs.Mapping{}
	if tr.mappingSchema == nil || tr.mappingSchema.Value == nil {
		return noMappings, nil
	}
	spec := propertyValueOrNil(tr.mappingSchema.Value, "spec")
	if spec == nil {
		return noMappings, nil
	}
	version := propertyValueOrNil(spec, tr.majorVersion)
	if version == nil {
		return noMappings, nil
	}
	return refs.FindMappings(version, []string{"spec", tr.majorVersion})
}

// MajorVersion returns the CRD pinned version
func (tr *translator) MajorVersion() string {
	return tr.majorVersion
}

// Scheme returns the Kubernetes scheme used to translate the CRD.
func (tr *translator) Scheme() *runtime.Scheme {
	return tr.scheme
}

// NewTranslator creates a translator for a particular CRD version. It is also
// locked into a particular API majorVersion.
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
func NewTranslator(scheme *runtime.Scheme, crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, majorVersion string) (Translator, error) {
	specVersion := crds.SelectVersion(&crd.Spec, crdVersion)
	if err := crds.AssertMajorVersion(specVersion, crd.Spec.Names.Kind, majorVersion); err != nil {
		return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", majorVersion, err)
	}
	var mappingSchema openapi3.Schema
	mappingString, ok := crd.Annotations["api-mappings"]
	if ok && mappingString != "" {
		jsonBytes, err := yaml.YAMLToJSON([]byte(mappingString))
		if err != nil {
			return nil, fmt.Errorf("failed to convert 'api-mappings' YAML to JSON: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &mappingSchema); err != nil {
			return nil, fmt.Errorf("failed to unmarshal 'api-mappings' JSON into schema: %w", err)
		}
	}

	return &translator{
		scheme:        scheme,
		majorVersion:  majorVersion,
		mappingSchema: &openapi3.SchemaRef{Value: &mappingSchema},
	}, nil
}

// NewPerVersionTranslators creates a set of translators indexed by SDK versions
//
// Given the following example resource:
//
//		apiVersion: atlas.generated.mongodb.com/v1
//		kind: SearchIndex
//		metadata:
//		  name: search-index
//		spec:
//		  v20250312:
//	    ...
//		  v20250810:
//
// In the above case crdVersion is "v1" and versions can be "v20250312"
// and/or "v20250810".
func NewPerVersionTranslators(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, versions ...string) (map[string]Translator, error) {
	translators := map[string]Translator{}
	specVersion := crds.SelectVersion(&crd.Spec, crdVersion)
	for _, version := range versions {
		if err := crds.AssertMajorVersion(specVersion, crd.Spec.Names.Kind, version); err != nil {
			return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", version, err)
		}
		var mappingSchema openapi3.Schema
		mappingString, ok := crd.Annotations["api-mappings"]
		if ok && mappingString != "" {
			jsonBytes, err := yaml.YAMLToJSON([]byte(mappingString))
			if err != nil {
				return nil, fmt.Errorf("failed to convert 'api-mappings' YAML to JSON: %w", err)
			}
			if err := json.Unmarshal(jsonBytes, &mappingSchema); err != nil {
				return nil, fmt.Errorf("failed to unmarshal 'api-mappings' JSON into schema: %w", err)
			}
		}

		translators[version] = &translator{
			majorVersion:  version,
			mappingSchema: &openapi3.SchemaRef{Value: &mappingSchema},
		}
	}
	return translators, nil
}

// collapseReferences finds all Kubernetes references, solves them and collapses
// them by replacing their values from the reference (e.g Kubernetes secret or
// group), into the corresponding API request value
func collapseReferences(tr Translator, req map[string]any, main client.Object, objs []client.Object) error {
	mappings, err := tr.Mappings()
	if err != nil {
		return fmt.Errorf("failed to extract mappings to collapse: %w", err)
	}
	return refs.CollapseAll(tr.Scheme(), mappings, main, objs, req)
}

// expandReferences finds all API fields that must be referenced, and expand
// such reference, moving the value (e.g. sensitive field or group id) to a
// referenced Kubernetes object (e.g. Kubernetes secret or Atlas Group)
func expandReferences(tr Translator, cr map[string]any, main client.Object, objs []client.Object) ([]client.Object, error) {
	mappings, err := tr.Mappings()
	if err != nil {
		return nil, fmt.Errorf("failed to extract mappings to expand: %w", err)
	}
	return refs.ExpandAll(tr.Scheme(), mappings, main, objs, cr)
}

func propertyValueOrNil(schema *openapi3.Schema, propertyName string) *openapi3.Schema {
	if schema.Properties != nil &&
		schema.Properties[propertyName] != nil && schema.Properties[propertyName].Value != nil {
		return schema.Properties[propertyName].Value
	}
	return nil
}
