package retry

import (
	"context"

	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RetryUpdateOnConflict is a wrapper around client-go/util/retry.RetryOnConflict,
// adding the following often repeated actions:
//
// 1. client.Get a resource for the given key
// 2. mutate the retrieved object using the given mutator function.
// 3. client.Update the updated resource and retry on conflict
// using the client-go/util/retry.DefaultRetry strategy.
func RetryUpdateOnConflict[T any](ctx context.Context, k8s client.Client, key client.ObjectKey, mutator func(*T)) (*T, error) {
	var obj T
	clientObj := any(&obj).(client.Object)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := k8s.Get(ctx, key, clientObj); err != nil {
			return err
		}
		mutator(&obj)
		return k8s.Update(ctx, clientObj)
	})
	return &obj, err
}
