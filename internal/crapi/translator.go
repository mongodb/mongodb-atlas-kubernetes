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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/objmap"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
)

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
		gvk:           schema.GroupVersionKind{Group: crd.Spec.Group, Version: crdVersion, Kind: crd.Spec.Names.Kind},
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
func NewPerVersionTranslators(scheme *runtime.Scheme, crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, versions ...string) (map[string]Translator, error) {
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
			scheme:        scheme,
			majorVersion:  version,
			gvk:           schema.GroupVersionKind{Group: crd.Spec.Group, Version: crdVersion, Kind: crd.Spec.Names.Kind},
			mappingSchema: &openapi3.SchemaRef{Value: &mappingSchema},
		}
	}
	return translators, nil
}

// translator implements Translator to translate from a given CRD to and from
// a given SDK version using the same upstream OpenAPI schema
type translator struct {
	scheme        *runtime.Scheme
	majorVersion  string
	gvk           schema.GroupVersionKind
	mappingSchema *openapi3.SchemaRef
}

func (tr *translator) ToAPI(target any, source client.Object, objs ...client.Object) error {
	if isNil(source) {
		return fmt.Errorf("source is nil")
	}
	if isNil(target) {
		return fmt.Errorf("target is nil")
	}
	if err := checkGVK(tr.scheme, source, tr.gvk); err != nil {
		return fmt.Errorf("Source GVK check failed: %w", err)
	}
	objMapSrc, err := objmap.ToObjectMap(source)
	if err != nil {
		return fmt.Errorf("failed to convert k8s source value to object map: %w", err)
	}
	targetObjMap := map[string]any{}

	if err := collapseReferences(tr, objMapSrc, source, objs); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	value, err := objmap.GetField[map[string]any](objMapSrc, "spec", tr.MajorVersion())
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}
	objmap.CopyFields(targetObjMap, value)
	rawEntry := value["entry"]
	if entry, ok := rawEntry.(map[string]any); ok {
		objmap.CopyFields(targetObjMap, entry)
	}
	if err := objmap.FromObjectMap(target, targetObjMap); err != nil {
		return fmt.Errorf("failed to set structured value from object map: %w", err)
	}
	return nil
}

func (tr *translator) FromAPI(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
	if isNil(source) {
		return nil, fmt.Errorf("source is nil")
	}
	if isNil(target) {
		return nil, fmt.Errorf("target is nil")
	}
	if err := checkGVK(tr.scheme, target, tr.gvk); err != nil {
		return nil, fmt.Errorf("Target GVK check failed: %w", err)
	}
	sourceObjMap, err := objmap.ToObjectMap(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert API source value to object map: %w", err)
	}

	targetObjMap, err := objmap.ToObjectMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to convert target value to object map: %w", err)
	}

	versionedSpec, err := objmap.GetOrCreateField(
		targetObjMap, map[string]any{}, "spec", tr.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to ensure versioned spec object in target: %w", err)
	}
	versionedStatus, err := objmap.GetOrCreateField(
		targetObjMap, map[string]any{}, "status", tr.MajorVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in target: %w", err)
	}

	objmap.CopyFields(versionedSpec, sourceObjMap)
	versionedSpecEntry := map[string]any{}
	objmap.CopyFields(versionedSpecEntry, sourceObjMap)
	versionedSpec["entry"] = versionedSpecEntry
	objmap.CopyFields(versionedStatus, sourceObjMap)

	extraObjects, err := expandReferences(tr, targetObjMap, target, objs)
	if err != nil {
		return nil, fmt.Errorf("failed to process API mappings: %w", err)
	}
	if err := objmap.FromObjectMap(target, targetObjMap); err != nil {
		return nil, fmt.Errorf("failed set structured kubernetes object from object map: %w", err)
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

// isNil properly checks if an interface{} value is nil, including nil pointers
// assigned to interfaces (which are not == nil in Go)
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}

func checkGVK(scheme *runtime.Scheme, target client.Object, gvk schema.GroupVersionKind) error {
	actualGVK := target.GetObjectKind().GroupVersionKind()
	if actualGVK.Kind == "" || actualGVK.GroupVersion().String() == "" {
		gvks, _, err := scheme.ObjectKinds(target)
		if err != nil || len(gvks) == 0 {
			return fmt.Errorf("failed to infer GroupVersionKind for resource from scheme: %w", err)
		}
		actualGVK = gvks[0]
	}
	if actualGVK.Kind != gvk.Kind || actualGVK.GroupVersion().String() != gvk.GroupVersion().String() {
		return fmt.Errorf("target must be a %q but got %q", gvk, actualGVK)
	}
	return nil
}
