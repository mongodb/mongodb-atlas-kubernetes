package hooks

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
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
