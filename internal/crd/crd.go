package crd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/josvazg/crd2go/internal/gotype"
)

const (
	OpenAPIObject  = "object"
	OpenAPIArray   = "array"
	OpenAPIString  = "string"
	OpenAPIInteger = "integer"
	OpenAPINumber  = "number"
	OpenAPIBoolean = "boolean"
)

var (
	// ErrNotApplied means the hook did nothing as the CRDType did not apply to this hook
	ErrNotApplied = errors.New("hook does not apply")
)

type FromOpenAPITypeFunc func(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error)

var Hooks = []FromOpenAPITypeFunc{
	UnstructuredHookFn,
	DictHookFn,
	DatetimeHookFn,
	PrimitiveHookFn,
	StructHookFn,
	ArrayHookFn,
}

type CRDType struct {
	Name    string
	Parents []string
	Schema  *apiextensionsv1.JSONSchemaProps
}

// FromOpenAPIType converts an OpenAPI schema to a GoType
func FromOpenAPIType(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		gt, err := hook(td, hooks, crdType)
		if errors.Is(err, ErrNotApplied) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("hook failed: %w", err)
		}
		return gt, nil
	}
	return nil, fmt.Errorf("unsupported Open API type %q", crdType.Name)
}

// orderedkeys returns a sorted slice of keys from the given map
func orderedkeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func CRD2Filename(crd *apiextensionsv1.CustomResourceDefinition) string {
	return fmt.Sprintf("%s.go", strings.ToLower(crd.Spec.Names.Kind))
}

func oneOf(s string, options ...string) bool {
	for _, opt := range options {
		if s == opt {
			return true
		}
	}
	return false
}
