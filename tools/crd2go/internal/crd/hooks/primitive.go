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

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

// PrimitiveHookFn converts an OpenAPI primitive type to a GoType
func PrimitiveHookFn(_ *gotype.TypeDict, _ []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if !crd.IsPrimitive(crdType) {
		return nil, fmt.Errorf("%s is not a primitive type: %w", crdType.Schema.Type, crd.ErrNotProcessed)
	}
	kind := crdType.Schema.Type
	goTypeName := openAPIkindGoType(kind)

	return gotype.NewPrimitive(goTypeName, goTypeName), nil
}

// openAPIkindGoType converts an OpenAPI kind to a Go type
func openAPIkindGoType(kind string) string {
	switch kind {
	case crd.OpenAPIInteger:
		return gotype.IntKind
	case crd.OpenAPINumber:
		return gotype.FloatKind
	case crd.OpenAPIBoolean:
		return gotype.BoolKind
	}

	// if none of the above, string type is assumed
	return gotype.StringKind
}
