package hooks

import "github.com/josvazg/crd2go/internal/crd"

var Hooks = []crd.OpenAPI2GoHook{
	UnstructuredHookFn,
	DictHookFn,
	DatetimeHookFn,
	PrimitiveHookFn,
	StructHookFn,
	ArrayHookFn,
}
