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

package child

import (
	"context"
	"testing"

	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

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

func TestGetDependencies(t *testing.T) {
	tests := []struct {
		name            string
		existingObjects []client.Object
		child           *v1.Child
		wantDepsCount   int
		wantDepNames    []string
		wantErr         bool
		wantErrContains []string
	}{
		{
			name: "resolves parent when ParentRef is set",
			existingObjects: []client.Object{
				&v1.Parent{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-parent",
						Namespace: "default",
					},
				},
			},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "my-parent"},
					},
				},
			},
			wantDepsCount: 1,
			wantDepNames:  []string{"my-parent"},
			wantErr:       false,
		},
		{
			name:            "returns error when parent not found",
			existingObjects: []client.Object{},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "nonexistent-parent"},
					},
				},
			},
			wantDepsCount:   0,
			wantErr:         true,
			wantErrContains: []string{"failed to get Parent", "nonexistent-parent"},
		},
		{
			name:            "returns empty deps when no ParentRef",
			existingObjects: []client.Object{},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentId: strPtr("direct-parent-id"),
					},
				},
			},
			wantDepsCount: 0,
			wantErr:       false,
		},
		{
			name:            "returns empty deps when V20250312 spec is nil",
			existingObjects: []client.Object{},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "default",
				},
				Spec: v1.ChildSpec{},
			},
			wantDepsCount: 0,
			wantErr:       false,
		},
		{
			name: "looks up parent in same namespace",
			existingObjects: []client.Object{
				&v1.Parent{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-parent",
						Namespace: "production",
					},
				},
			},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "production",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "my-parent"},
					},
				},
			},
			wantDepsCount: 1,
			wantDepNames:  []string{"my-parent"},
			wantErr:       false,
		},
		{
			name: "fails when parent exists in different namespace",
			existingObjects: []client.Object{
				&v1.Parent{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-parent",
						Namespace: "default",
					},
				},
			},
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-child",
					Namespace: "production",
				},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "my-parent"},
					},
				},
			},
			wantDepsCount:   0,
			wantErr:         true,
			wantErrContains: []string{"failed to get Parent"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			scheme := newScheme(t)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.existingObjects...).
				Build()

			handler := NewHandlerv20250312(fakeClient, nil, nil, false)
			deps, err := handler.getDependencies(ctx, tc.child)

			if tc.wantErr {
				require.Error(t, err)
				for _, errStr := range tc.wantErrContains {
					assert.Contains(t, err.Error(), errStr)
				}
				assert.Empty(t, deps)
				return
			}

			require.NoError(t, err)
			assert.Len(t, deps, tc.wantDepsCount)

			for i, name := range tc.wantDepNames {
				assert.Equal(t, name, deps[i].GetName())
			}
		})
	}
}
