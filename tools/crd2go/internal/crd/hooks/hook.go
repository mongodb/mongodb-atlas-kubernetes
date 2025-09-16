package hooks

import "github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"

var Hooks = []crd.OpenAPI2GoHook{
	UnstructuredHookFn,
	DictHookFn,
	DatetimeHookFn,
	PrimitiveHookFn,
	StructHookFn,
	ArrayHookFn,
}
