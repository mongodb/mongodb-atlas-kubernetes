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

package hooks

import (
	"fmt"
	"slices"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

// StructHookFn converts and OpenAPI object to a GoType struct
func StructHookFn(td *gotype.TypeDict, hooks []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != crd.OpenAPIObject {
		return nil, fmt.Errorf("%s is not an object: %w", crdType.Schema.Type, crd.ErrNotProcessed)
	}
	fields := []*gotype.GoField{}
	fieldsParents := append(crdType.Parents, crdType.Name)
	for _, key := range orderedKeys(crdType.Schema.Properties) {
		props := crdType.Schema.Properties[key]
		fieldType, err := crd.FromOpenAPIType(td, hooks, &crd.CRDType{
			Name:    key,
			Parents: fieldsParents,
			Schema:  &props,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s type: %w", key, err)
		}
		field := gotype.NewGoFieldWithKey(key, key, fieldType)
		field.Comment = props.Description
		field.Required = slices.Contains(crdType.Schema.Required, key)
		if err := td.RenameField(field, fieldsParents); err != nil {
			return nil, fmt.Errorf("failed to rename field %v: %w", field, err)
		}
		fields = append(fields, field)
	}
	return gotype.NewStruct(crdType.Name, fields), nil
}

// orderedKeys returns a sorted slice of keys from the given map
func orderedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
