// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flexcluster

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// newFakeClientWithGVK creates a fake client that preserves GVK on retrieved objects.
// The fake client from controller-runtime doesn't preserve TypeMeta by default, so we wrap it
// to automatically set GVK after Get operations using the scheme.
func newFakeClientWithGVK(scheme *runtime.Scheme, objects []client.Object, statusSubresource client.Object) client.Client {
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		WithStatusSubresource(statusSubresource).
		Build()
	return &gvkPreservingClient{Client: fakeClient, scheme: scheme}
}

// gvkPreservingClient wraps a client.Client to set GVK on objects after retrieval.
// This is needed because the fake client doesn't preserve TypeMeta when retrieving objects.
type gvkPreservingClient struct {
	client.Client
	scheme *runtime.Scheme
}

func (c *gvkPreservingClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if err := c.Client.Get(ctx, key, obj, opts...); err != nil {
		return err
	}
	return setGVK(c.scheme, obj)
}

// setGVK sets the GroupVersionKind on an object using the scheme.
// This is necessary because fake clients don't preserve TypeMeta when retrieving objects.
func setGVK(scheme *runtime.Scheme, obj runtime.Object) error {
	gvks, _, err := scheme.ObjectKinds(obj)
	if err != nil {
		return err
	}
	if len(gvks) == 0 {
		return nil
	}
	objectKind, ok := obj.(schema.ObjectKind)
	if !ok {
		return nil
	}
	objectKind.SetGroupVersionKind(gvks[0])
	return nil
}
