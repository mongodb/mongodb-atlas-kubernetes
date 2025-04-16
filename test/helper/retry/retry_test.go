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
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestRetryUpdateOnConflict(t *testing.T) {
	for _, tc := range []struct {
		name             string
		key              client.ObjectKey
		objects          []client.Object
		interceptorFuncs interceptor.Funcs

		want    *akov2.AtlasProject
		wantErr string
	}{
		{
			name: "fail immediately if not found",
			key: types.NamespacedName{
				Name:      "foo",
				Namespace: "bar",
			},
			want:    &akov2.AtlasProject{},
			wantErr: "atlasprojects.atlas.mongodb.com \"foo\" not found",
		},
		{
			name: "succeed if found",
			key: types.NamespacedName{
				Name:      "foo",
				Namespace: "bar",
			},
			objects: []client.Object{
				&akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
			},
			want: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
		},
		{
			name: "exhaust on conflict",
			key: types.NamespacedName{
				Name:      "foo",
				Namespace: "bar",
			},
			objects: []client.Object{
				&akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
			},
			interceptorFuncs: interceptor.Funcs{
				Update: func(context.Context, client.WithWatch, client.Object, ...client.UpdateOption) error {
					return &apierrors.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonConflict, Message: "conflict"}}
				},
			},
			want:    &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
			wantErr: "conflict",
		},
		{
			name: "fail on any other update error",
			key: types.NamespacedName{
				Name:      "foo",
				Namespace: "bar",
			},
			objects: []client.Object{
				&akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
			},
			interceptorFuncs: interceptor.Funcs{
				Update: func(context.Context, client.WithWatch, client.Object, ...client.UpdateOption) error {
					return errors.New("boom")
				},
			},
			want:    &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}},
			wantErr: "boom",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.objects...).
				WithInterceptorFuncs(tc.interceptorFuncs).
				Build()

			got, err := RetryUpdateOnConflict(context.Background(), k8sClient, tc.key, func(*akov2.AtlasProject) {})
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}

			if gotErr != tc.wantErr {
				t.Errorf("want error %q, got %q", tc.wantErr, gotErr)
			}

			// ignore unnecessary fields
			got.ResourceVersion = ""

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want AtlasProject %+v, got %+v", tc.want, got)
			}
		})
	}
}
