package dryrun

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

// key defines the package local unique type for context values.
// In conjunction with context keys it forms a unique identifier.
// It is unexported to prevent collisions.
type key int

const (
	// runtimeObjectKey is the context key for the runtime object to be used for dry run recording.
	runtimeObjectKey key = iota
)

// WithRuntimeObject returns a new context that wraps the provided context and contains the provided runtime object.
func WithRuntimeObject(ctx context.Context, obj runtime.Object) context.Context {
	return context.WithValue(ctx, runtimeObjectKey, obj)
}

func runtimeObjectFrom(ctx context.Context) (runtime.Object, bool) {
	recorder, ok := ctx.Value(runtimeObjectKey).(runtime.Object)
	return recorder, ok
}
