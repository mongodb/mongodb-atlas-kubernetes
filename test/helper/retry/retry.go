// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
