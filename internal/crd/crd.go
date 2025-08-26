package crd

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/k8s"
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
	switch crdType.Schema.Type {
	case OpenAPIObject:
		return fromOpenAPIStruct(td, hooks, crdType)
	case OpenAPIArray:
		return fromOpenAPIArray(td, hooks, crdType)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", crdType.Name)
	}
}

// fromOpenAPIStruct converts and OpenAPI object to a GoType struct
func fromOpenAPIStruct(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	fields := []*gotype.GoField{}
	fieldsParents := append(crdType.Parents, crdType.Name)
	for _, key := range orderedkeys(crdType.Schema.Properties) {
		props := crdType.Schema.Properties[key]
		fieldType, err := FromOpenAPIType(td, hooks, &CRDType{
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

// fromOpenAPIArray converts an OpenAPI array schema to a GoType array
func fromOpenAPIArray(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Items == nil {
		return nil, fmt.Errorf("array %s has no items", crdType.Name)
	}
	if crdType.Schema.Items.Schema == nil {
		return nil, fmt.Errorf("array %s has no items schema", crdType.Name)
	}
	elementType, err := FromOpenAPIType(td, hooks, &CRDType{
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

// orderedkeys returns a sorted slice of keys from the given map
func orderedkeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func KnownTypes() []*gotype.GoType {
	return []*gotype.GoType{
		gotype.MustTypeFrom(reflect.TypeOf(k8s.LocalReference{})),
		gotype.MustTypeFrom(reflect.TypeOf(k8s.Reference{})),
		gotype.SetAlias(gotype.MustTypeFrom(reflect.TypeOf(metav1.Condition{})), "metav1"),
	}
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
