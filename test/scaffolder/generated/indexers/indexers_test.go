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

package indexer

import (
	"context"
	"testing"

	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

// newScheme returns a scheme with all generated types registered.
func newScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	err := v1.AddToScheme(scheme)
	require.NoError(t, err)
	return scheme
}

func strPtr(s string) *string {
	return &s
}

// TestChildByParentIndexer_Keys verifies the indexer extracts correct keys
// from Child resources referencing a Parent.
func TestChildByParentIndexer_Keys(t *testing.T) {
	logger := zaptest.NewLogger(t)
	idx := NewChildByParentIndexer(logger)

	// Verify indexer metadata
	assert.Equal(t, "child.parentRef", idx.Name())
	assert.IsType(t, &v1.Child{}, idx.Object())

	tests := []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "wrong type returns nil",
			object:   &v1.Parent{},
			wantKeys: nil,
		},
		{
			name: "child with parentRef returns key",
			object: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "my-parent"},
					},
				},
			},
			wantKeys: []string{"default/my-parent"},
		},
		{
			name: "child without V20250312 returns nil",
			object: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{},
			},
			wantKeys: nil,
		},
		{
			name: "child with nil parentRef returns nil",
			object: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentId: strPtr("some-id"),
						// ParentRef is nil
					},
				},
			},
			wantKeys: nil,
		},
		{
			name: "child with empty parentRef name returns nil",
			object: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: ""},
					},
				},
			},
			wantKeys: nil,
		},
		{
			name: "child in different namespace",
			object: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-child",
					Namespace: "production",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "prod-parent"},
					},
				},
			},
			wantKeys: []string{"production/prod-parent"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			keys := idx.Keys(tc.object)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}

// TestParentBySecretIndexer_Keys verifies the indexer extracts keys from
// nested arrays (integrations -> credentials -> secretRef).
func TestParentBySecretIndexer_Keys(t *testing.T) {
	logger := zaptest.NewLogger(t)
	idx := NewParentBySecretIndexer(logger)

	// Verify indexer metadata
	assert.Equal(t, "parent.secretRef", idx.Name())
	assert.IsType(t, &v1.Parent{}, idx.Object())

	tests := []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "wrong type returns nil",
			object:   &v1.Child{},
			wantKeys: nil,
		},
		{
			name: "parent without integrations returns nil",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{},
			},
			wantKeys: nil,
		},
		{
			name: "parent with empty integrations returns nil",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{},
				},
			},
			wantKeys: nil,
		},
		{
			name: "parent with single secret ref",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Name: strPtr("integration-1"),
							Credentials: &[]v1.Credentials{
								{SecretRef: &k8s.LocalReference{Name: "secret-1"}},
							},
						},
					},
				},
			},
			wantKeys: []string{"default/secret-1"},
		},
		{
			name: "parent with multiple secrets across integrations",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Name: strPtr("integration-1"),
							Credentials: &[]v1.Credentials{
								{SecretRef: &k8s.LocalReference{Name: "secret-1"}},
								{SecretRef: &k8s.LocalReference{Name: "secret-2"}},
							},
						},
						{
							Name: strPtr("integration-2"),
							Credentials: &[]v1.Credentials{
								{SecretRef: &k8s.LocalReference{Name: "secret-3"}},
							},
						},
					},
				},
			},
			wantKeys: []string{"default/secret-1", "default/secret-2", "default/secret-3"},
		},
		{
			name: "parent with nil credentials slice returns nil",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Name: strPtr("integration-1"),
							// Credentials is nil
						},
					},
				},
			},
			wantKeys: nil,
		},
		{
			name: "parent with nil secretRef in credentials returns nil",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Name: strPtr("integration-1"),
							Credentials: &[]v1.Credentials{
								{SecretRef: nil},
							},
						},
					},
				},
			},
			wantKeys: nil,
		},
		{
			name: "parent with empty secretRef name returns nil",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "default",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Name: strPtr("integration-1"),
							Credentials: &[]v1.Credentials{
								{SecretRef: &k8s.LocalReference{Name: ""}},
							},
						},
					},
				},
			},
			wantKeys: nil,
		},
		{
			name: "parent in different namespace",
			object: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-parent",
					Namespace: "production",
				},
				Spec: v1.ParentSpec{
					Integrations: &[]v1.Integrations{
						{
							Credentials: &[]v1.Credentials{
								{SecretRef: &k8s.LocalReference{Name: "prod-secret"}},
							},
						},
					},
				},
			},
			wantKeys: []string{"production/prod-secret"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			keys := idx.Keys(tc.object)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}

// TestChildByParentMapFunc verifies the MapFunc returns reconcile requests
// for all Children referencing a given Parent.
func TestChildByParentMapFunc(t *testing.T) {
	ctx := context.Background()
	scheme := newScheme(t)

	// Create test data: one parent with two children referencing it
	parent := &v1.Parent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-parent",
			Namespace: "default",
		},
	}

	child1 := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "child-1",
			Namespace: "default",
		},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{
				ParentRef: &k8s.LocalReference{Name: "my-parent"},
			},
		},
	}

	child2 := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "child-2",
			Namespace: "default",
		},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{
				ParentRef: &k8s.LocalReference{Name: "my-parent"},
			},
		},
	}

	// Child referencing a different parent
	childOther := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "child-other",
			Namespace: "default",
		},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{
				ParentRef: &k8s.LocalReference{Name: "other-parent"},
			},
		},
	}

	logger := zaptest.NewLogger(t)
	idx := NewChildByParentIndexer(logger)

	// Build fake client with field index
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(parent, child1, child2, childOther).
		WithIndex(idx.Object(), idx.Name(), idx.Keys).
		Build()

	mapFunc := NewChildByParentMapFunc(fakeClient)
	requests := mapFunc(ctx, parent)

	// Should return reconcile requests for child-1 and child-2
	require.Len(t, requests, 2)

	names := []string{requests[0].Name, requests[1].Name}
	assert.Contains(t, names, "child-1")
	assert.Contains(t, names, "child-2")
}

// TestParentBySecretMapFunc verifies the MapFunc returns reconcile requests
// for all Parents referencing a given Secret.
func TestParentBySecretMapFunc(t *testing.T) {
	ctx := context.Background()
	scheme := newScheme(t)

	// Create test data: two parents referencing the same secret
	parent1 := &v1.Parent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "parent-1",
			Namespace: "default",
		},
		Spec: v1.ParentSpec{
			Integrations: &[]v1.Integrations{
				{
					Credentials: &[]v1.Credentials{
						{SecretRef: &k8s.LocalReference{Name: "shared-secret"}},
					},
				},
			},
		},
	}

	parent2 := &v1.Parent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "parent-2",
			Namespace: "default",
		},
		Spec: v1.ParentSpec{
			Integrations: &[]v1.Integrations{
				{
					Credentials: &[]v1.Credentials{
						{SecretRef: &k8s.LocalReference{Name: "shared-secret"}},
					},
				},
			},
		},
	}

	// Parent referencing a different secret
	parentOther := &v1.Parent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "parent-other",
			Namespace: "default",
		},
		Spec: v1.ParentSpec{
			Integrations: &[]v1.Integrations{
				{
					Credentials: &[]v1.Credentials{
						{SecretRef: &k8s.LocalReference{Name: "other-secret"}},
					},
				},
			},
		},
	}

	logger := zaptest.NewLogger(t)
	idx := NewParentBySecretIndexer(logger)

	// Build fake client with field index
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(parent1, parent2, parentOther).
		WithIndex(idx.Object(), idx.Name(), idx.Keys).
		Build()

	// Simulate a Secret object being the trigger
	secretObj := &metav1.PartialObjectMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "shared-secret",
			Namespace: "default",
		},
	}

	mapFunc := NewParentBySecretMapFunc(fakeClient)
	requests := mapFunc(ctx, secretObj)

	// Should return reconcile requests for parent-1 and parent-2
	require.Len(t, requests, 2)

	names := []string{requests[0].Name, requests[1].Name}
	assert.Contains(t, names, "parent-1")
	assert.Contains(t, names, "parent-2")
}
