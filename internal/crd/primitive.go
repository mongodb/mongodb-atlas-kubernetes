package crd

import (
	"fmt"

	"github.com/josvazg/crd2go/internal/gotype"
)

// PrimitiveHookFn converts an OpenAPI primitive type to a GoType
func PrimitiveHookFn(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	kind := crdType.Schema.Type
	if !oneOf(kind, OpenAPIString, OpenAPIInteger, OpenAPINumber, OpenAPIBoolean) {
		return nil, fmt.Errorf("%s is not a primitive type: %w", kind, ErrNotApplied)
	}
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return gotype.NewPrimitive(goTypeName, goTypeName), nil
}

// openAPIKindtoGoType converts an OpenAPI kind to a Go type
func openAPIKindtoGoType(kind string) (string, error) {
	switch kind {
	case OpenAPIString:
		return gotype.StringKind, nil
	case OpenAPIInteger:
		return gotype.IntKind, nil
	case OpenAPINumber:
		return gotype.FloatKind, nil
	case OpenAPIBoolean:
		return gotype.BoolKind, nil
	default:
		return "", fmt.Errorf("unsupported Open API kind %s", kind)
	}
}
