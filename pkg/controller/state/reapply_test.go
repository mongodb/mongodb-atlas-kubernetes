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

package state

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestReapplyPeriod(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
		want        time.Duration
		wantOk      bool
		wantErr     string
	}{
		{
			name:        "valid period",
			annotations: map[string]string{"mongodb.com/reapply-period": "2h"},
			want:        2 * time.Hour,
			wantOk:      true,
		},
		{
			name:        "period missing",
			annotations: map[string]string{},
			want:        0,
			wantOk:      false,
		},
		{
			name:        "period invalid format",
			annotations: map[string]string{"mongodb.com/reapply-period": "not-a-period"},
			want:        0,
			wantOk:      false,
			wantErr:     "invalid duration",
		},
		{
			name:        "period too short",
			annotations: map[string]string{"mongodb.com/reapply-period": "30s"},
			want:        0,
			wantOk:      false,
			wantErr:     "must be greater than 60m",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := newUnstructuredObj(tc.annotations)
			got, ok, err := ReapplyPeriod(obj)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestReapplyTimestamp(t *testing.T) {
	now := time.Now().UnixMilli()
	tests := []struct {
		name        string
		annotations map[string]string
		want        int64
		wantOk      bool
		wantErr     string
	}{
		{
			name:        "valid timestamp",
			annotations: map[string]string{AnnotationReapplyTimestamp: strconv.FormatInt(now, 10)},
			want:        now,
			wantOk:      true,
		},
		{
			name:        "timestamp missing",
			annotations: map[string]string{},
			want:        0,
			wantOk:      false,
		},
		{
			name:        "invalid timestamp",
			annotations: map[string]string{AnnotationReapplyTimestamp: "not-a-number"},
			want:        0,
			wantOk:      false,
			wantErr:     "parsing \"not-a-number\": invalid syntax",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := newUnstructuredObj(tc.annotations)
			got, ok, err := ReapplyTimestamp(obj)
			assertErrContains(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				assert.Equal(t, tc.want, got.UnixMilli())
			}
		})
	}
}

func TestShouldReapply(t *testing.T) {
	past := time.Now().Add(-2 * time.Hour).UnixMilli()
	future := time.Now().Add(2 * time.Hour).UnixMilli()

	tests := []struct {
		name        string
		annotations map[string]string
		want        bool
		wantErr     string
	}{
		{
			name: "should reapply (past+1h < now)",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   strconv.FormatInt(past, 10),
				"mongodb.com/reapply-period": "1h",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "should not reapply (future+1h > now)",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   strconv.FormatInt(future, 10),
				"mongodb.com/reapply-period": "1h",
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "missing period",
			annotations: map[string]string{
				AnnotationReapplyTimestamp: strconv.FormatInt(past, 10),
			},
			want:    false,
			wantErr: "",
		},
		{
			name:        "missing timestamp",
			annotations: map[string]string{"mongodb.com/reapply-period": "1h"},
			want:        false,
			wantErr:     "",
		},
		{
			name: "invalid period",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   strconv.FormatInt(past, 10),
				"mongodb.com/reapply-period": "bad",
			},
			want:    false,
			wantErr: "invalid duration",
		},
		{
			name: "invalid timestamp",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   "bad",
				"mongodb.com/reapply-period": "1h",
			},
			want:    false,
			wantErr: "invalid syntax",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := newUnstructuredObj(tc.annotations)
			got, err := ShouldReapply(obj)
			assertErrContains(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPatchReapplyTimestamp(t *testing.T) {
	now := time.Now()
	pastMillis := strconv.FormatInt(now.Add(-2*time.Hour).UnixMilli(), 10)

	tests := []struct {
		name        string
		annotations map[string]string
		patchErr    error
		want        time.Duration
		wantErr     string // substring to match in the error message
		wantPatched bool   // true if we expect the annotation to be updated
	}{
		{
			name: "patch performed",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   pastMillis,
				"mongodb.com/reapply-period": "1h",
			},
			want:        time.Hour,
			wantErr:     "",
			wantPatched: true,
		},
		{
			name:        "patch not needed (no period)",
			annotations: map[string]string{},
			want:        0,
			wantErr:     "",
			wantPatched: false,
		},
		{
			name: "patch error",
			annotations: map[string]string{
				AnnotationReapplyTimestamp:   pastMillis,
				"mongodb.com/reapply-period": "1h",
			},
			patchErr:    errors.New("fail"),
			want:        0,
			wantErr:     "fail",
			wantPatched: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "dummy",
					Namespace:   "default",
					Annotations: tc.annotations,
				},
			}
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			patchFn := func(_ context.Context, _ client.WithWatch, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return tc.patchErr
			}
			if tc.patchErr == nil {
				patchFn = nil
			}
			c := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(obj.DeepCopy()).
				WithInterceptorFuncs(interceptor.Funcs{Patch: patchFn}).
				Build()
			ctx := context.Background()

			period, err := PatchReapplyTimestamp(ctx, c, obj)
			assertErrContains(t, tc.wantErr, err)
			assert.Equal(t, tc.want, period)

			fetched := &corev1.Pod{}
			_ = c.Get(ctx, client.ObjectKeyFromObject(obj), fetched)

			annot := fetched.GetAnnotations()
			_, patched := annot[AnnotationReapplyTimestamp]

			assert.Equal(t, tc.wantPatched, patched, "Annotation patched?")
		})
	}
}

// Helper to create an Unstructured object with annotations.
func newUnstructuredObj(annotations map[string]string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAnnotations(annotations)
	return obj
}

func assertErrContains(t *testing.T, wantErr string, err error) {
	if wantErr == "" {
		assert.NoError(t, err)
	} else {
		assert.ErrorContains(t, err, wantErr)
	}
}
