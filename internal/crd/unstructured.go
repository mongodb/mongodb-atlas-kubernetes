package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/josvazg/crd2go/internal/gotype"
)

func UnstructuredHookFn(td *gotype.TypeDict, _ []FromOpenAPITypeFunc, crdType *CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type == OpenAPIObject && isUnstructured(crdType.Schema) {
		return gotype.JSONType, nil
	}
	return nil, fmt.Errorf("%s is not unstructured (has %d properties and x-preserve-unknown-fields is %v): %w",
		crdType.Schema.Type, len(crdType.Schema.Properties), crdType.Schema.XPreserveUnknownFields, ErrNotApplied)
}

func isUnstructured(schema *apiextensionsv1.JSONSchemaProps) bool {
	return (len(schema.Properties) == 0 && schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields == true)
}
