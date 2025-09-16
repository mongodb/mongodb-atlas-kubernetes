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

// PrimitiveHookFn converts an OpenAPI primitive type to a GoType
func PrimitiveHookFn(td *gotype.TypeDict, hooks []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if !crd.IsPrimitive(crdType) {
		return nil, fmt.Errorf("%s is not a primitive type: %w", crdType.Schema.Type, crd.ErrNotProcessed)
	}
	kind := crdType.Schema.Type
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return gotype.NewPrimitive(goTypeName, goTypeName), nil
}

// openAPIKindtoGoType converts an OpenAPI kind to a Go type
func openAPIKindtoGoType(kind string) (string, error) {
	switch kind {
	case crd.OpenAPIString:
		return gotype.StringKind, nil
	case crd.OpenAPIInteger:
		return gotype.IntKind, nil
	case crd.OpenAPINumber:
		return gotype.FloatKind, nil
	case crd.OpenAPIBoolean:
		return gotype.BoolKind, nil
	default:
		return "", fmt.Errorf("unsupported Open API kind %s", kind)
	}
}
