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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestPatcher(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)

	baseDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 5,
			Conditions: []appsv1.DeploymentCondition{
				{
					Status: corev1.ConditionTrue,
					Reason: "Ready",
				},
			},
		},
	}

	tests := []struct {
		name         string
		obj          client.Object
		setupPatcher func(*Patcher)
		interceptors *interceptor.Funcs
		wantErr      string
		wantObj      client.Object
	}{
		{
			name: "no changes - does nothing",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				// Don't call any update methods
			},
			wantErr: "",
			wantObj: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					Namespace:   "default",
					Annotations: nil, // unset
				},
				Status: appsv1.DeploymentStatus{}, // unset
			},
		},
		{
			name: "patches object with state tracker",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				p.UpdateStateTracker()
			},
			wantErr: "",
			wantObj: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationStateTracker: "58c96d66dbcfb4d4546f",
					},
				},
				Status: appsv1.DeploymentStatus{}, // unset
			},
		},
		{
			name: "patches status only",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				p.UpdateStatus()
			},
			wantErr: "",
			wantObj: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					Namespace:   "default",
					Annotations: nil, // unset
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 5,
					Conditions:    nil, // not set by patcher
				},
			},
		},
		{
			name: "patches both object and status",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				p.UpdateStateTracker().UpdateStatus()
			},
			wantErr: "",
			wantObj: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationStateTracker: "58c96d66dbcfb4d4546f",
					},
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 5,
					Conditions:    nil, // removed by patcher
				},
			},
		},
		{
			name: "returns pre-existing error without patching",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				p.err = errors.New("pre-existing error")
				p.objectChanged = true
			},
			wantErr: "pre-existing error",
		},
		{
			name: "returns error on status patch failure",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				p.UpdateStatus()
			},
			interceptors: &interceptor.Funcs{
				SubResourcePatch: func(ctx context.Context, c client.Client, subResourceName string, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
					return errors.New("status patch failed")
				},
			},
			wantErr: "status patch failed",
		},
		{
			name: "status patch error stops object patch",
			obj:  baseDeployment.DeepCopy(),
			setupPatcher: func(p *Patcher) {
				// Enable both status and object changes
				p.UpdateStatus().UpdateStateTracker()
			},
			interceptors: &interceptor.Funcs{
				SubResourcePatch: func(ctx context.Context, c client.Client, subResourceName string, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
					return errors.New("status patch failed first")
				},
			},
			// Status patch fails first, object patch should be skipped
			wantErr: "status patch failed first",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme)
			builder = builder.WithObjects(tc.obj)
			builder = builder.WithStatusSubresource(tc.obj)
			if tc.interceptors != nil {
				builder = builder.WithInterceptorFuncs(*tc.interceptors)
			}
			c := builder.Build()

			p := NewPatcher(tc.obj)
			if tc.setupPatcher != nil {
				tc.setupPatcher(p)
			}

			err := p.Patch(context.Background(), c)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				assert.NoError(t, err)

				gotObj := &appsv1.Deployment{}
				require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(p.patchedObj.Object, gotObj))
				require.True(t, reflect.DeepEqual(gotObj, tc.wantObj))
			}
		})
	}
}
