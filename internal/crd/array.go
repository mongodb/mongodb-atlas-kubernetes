package crd

import (
	"fmt"

	"github.com/josvazg/crd2go/internal/gotype"
)

// ArrayHookFn converts an OpenAPI array schema to a GoType array
func ArrayHookFn(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != OpenAPIArray {
		return nil, fmt.Errorf("%s is not an array: %w", crdType.Schema.Type, ErrNotApplied)
	}
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
