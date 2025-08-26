package crd

import (
	"fmt"
	"slices"

	"github.com/josvazg/crd2go/internal/gotype"
)

// StructHookFn converts and OpenAPI object to a GoType struct
func StructHookFn(td *gotype.TypeDict, hooks []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != OpenAPIObject {
		return nil, fmt.Errorf("%s is not an object: %w", crdType.Schema.Type, ErrNotApplied)
	}
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
