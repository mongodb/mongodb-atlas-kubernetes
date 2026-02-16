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

package parent

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

	indexer "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/indexers"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

func newScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	err := v1.AddToScheme(scheme)
	require.NoError(t, err)
	return scheme
}

func TestGetDependents(t *testing.T) {
	tests := []struct {
		name              string
		children          []client.Object
		parent            *v1.Parent
		wantRequestsCount int
		wantRequestNames  []string
	}{
		{
			name: "returns children that reference the parent",
			children: []client.Object{
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-1",
						Namespace: "default",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "my-parent"},
						},
					},
				},
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-2",
						Namespace: "default",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "my-parent"},
						},
					},
				},
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-other",
						Namespace: "default",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "other-parent"},
						},
					},
				},
			},
			parent: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-parent",
					Namespace: "default",
				},
			},
			wantRequestsCount: 2,
			wantRequestNames:  []string{"child-1", "child-2"},
		},
		{
			name: "returns empty when no children reference the parent",
			children: []client.Object{
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-other",
						Namespace: "default",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "other-parent"},
						},
					},
				},
			},
			parent: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lonely-parent",
					Namespace: "default",
				},
			},
			wantRequestsCount: 0,
			wantRequestNames:  nil,
		},
		{
			name: "only returns children in the same namespace",
			children: []client.Object{
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-same-ns",
						Namespace: "default",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "my-parent"},
						},
					},
				},
				&v1.Child{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "child-diff-ns",
						Namespace: "production",
					},
					Spec: v1.ChildSpec{
						V20250312: &v1.V20250312{
							ParentRef: &k8s.LocalReference{Name: "my-parent"},
						},
					},
				},
			},
			parent: &v1.Parent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-parent",
					Namespace: "default",
				},
			},
			wantRequestsCount: 1,
			wantRequestNames:  []string{"child-same-ns"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			scheme := newScheme(t)
			logger := zaptest.NewLogger(t)
			idx := indexer.NewChildByParentIndexer(logger)

			allObjects := append(tc.children, tc.parent)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(allObjects...).
				WithIndex(idx.Object(), idx.Name(), idx.Keys).
				Build()

			handler := NewHandlerintegrations(fakeClient, nil, nil, false)
			requests := handler.getDependents(ctx, tc.parent)

			require.Len(t, requests, tc.wantRequestsCount)

			if tc.wantRequestNames != nil {
				names := make([]string, len(requests))
				for i, req := range requests {
					names[i] = req.Name
				}
				for _, wantName := range tc.wantRequestNames {
					assert.Contains(t, names, wantName)
				}
			}
		})
	}
}
