package hooks

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func UnstructuredHookFn(td *gotype.TypeDict, _ []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != crd.OpenAPIObject || !isUnstructured(crdType.Schema) {
		return nil, fmt.Errorf("%s is not unstructured (has %d properties and x-preserve-unknown-fields is %v): %w",
			crdType.Schema.Type, len(crdType.Schema.Properties), crdType.Schema.XPreserveUnknownFields, crd.ErrNotProcessed)
	}
	return gotype.JSONType, nil
}

func isUnstructured(schema *apiextensionsv1.JSONSchemaProps) bool {
	return (len(schema.Properties) == 0 && schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields)
}
