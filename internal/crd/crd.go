package crd

import (
	"errors"
	"fmt"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/josvazg/crd2go/internal/gotype"
)

const (
	FirstVersion = ""
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
	// ErrNotProcessed means the hook did nothing as the CRDType did not apply to this hook
	ErrNotProcessed = errors.New("hook not")
)

type OpenAPI2GoHook func(td *gotype.TypeDict, hooks []OpenAPI2GoHook, crdType *CRDType) (*gotype.GoType, error)

type CRDType struct {
	Name    string
	Parents []string
	Schema  *apiextensionsv1.JSONSchemaProps
}

// FromOpenAPIType converts an OpenAPI schema to a GoType
func FromOpenAPIType(td *gotype.TypeDict, hooks []OpenAPI2GoHook, crdType *CRDType) (*gotype.GoType, error) {
	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		gt, err := hook(td, hooks, crdType)
		if errors.Is(err, ErrNotProcessed) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("hook failed: %w", err)
		}
		return gt, nil
	}
	return nil, fmt.Errorf("unsupported Open API type %q", crdType.Name)
}

func Kind2Filename(kind string) string {
	return fmt.Sprintf("%s.go", strings.ToLower(kind))
}

func IsPrimitive(crdType *CRDType) bool {
	return matchesAny(crdType.Schema.Type, OpenAPIString, OpenAPIInteger, OpenAPINumber, OpenAPIBoolean)
}

func IsDateTimeFormat(crdType *CRDType) bool {
	return matchesAny(crdType.Schema.Format, "datetime", "date-time")
}

func matchesAny(s string, options ...string) bool {
	for _, opt := range options {
		if s == opt {
			return true
		}
	}
	return false
}
