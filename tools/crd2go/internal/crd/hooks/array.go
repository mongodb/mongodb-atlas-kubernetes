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

	"github.com/josvazg/crd2go/internal/crd"
	"github.com/josvazg/crd2go/internal/gotype"
)

// ArrayHookFn converts an OpenAPI array schema to a GoType array
func ArrayHookFn(td *gotype.TypeDict, hooks []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != crd.OpenAPIArray {
		return nil, fmt.Errorf("%s is not an array: %w", crdType.Schema.Type, crd.ErrNotProcessed)
	}
	if crdType.Schema.Items == nil {
		return nil, fmt.Errorf("array %s has no items", crdType.Name)
	}
	if crdType.Schema.Items.Schema == nil {
		return nil, fmt.Errorf("array %s has no items schema", crdType.Name)
	}
	elementType, err := crd.FromOpenAPIType(td, hooks, &crd.CRDType{
		Name:   crdType.Name,
		Schema: crdType.Schema.Items.Schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse array %s element type: %w", crdType.Name, err)
	}
	if err := td.RenameType(crdType.Parents, elementType); err != nil {
		return nil, fmt.Errorf("failed to rename element type under %s: %w", crdType.Name, err)
	}
	return gotype.NewArray(elementType), nil
}
