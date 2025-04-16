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

package atlasdatafederation

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

func TestDeleteConnectionSecrets(t *testing.T) {
	for _, tc := range []struct {
		name              string
		service           func(serviceMock *translation.DataFederationServiceMock) datafederation.DataFederationService
		atlasProject      *akov2.AtlasProject
		dataFederation    *akov2.AtlasDataFederation
		connectionSecrets []*corev1.Secret

		wantDataFederation *akov2.AtlasDataFederation
		wantSecrets        []corev1.Secret
		wantResult         workflow.Result
	}{
		{
			name: "no finalizer",
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			},
			wantSecrets: []corev1.Secret{},
			wantResult:  workflow.OK(),
		},
		{
			name: "finalizer set and deletion protection is enabled",
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers:  []string{customresource.FinalizerLabel},
					Annotations: map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep},
				},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
					Annotations: map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep},
				},
			},
			wantSecrets: []corev1.Secret{},
			wantResult:  workflow.OK(),
		},
		{
			name: "federation object without secrets",
			service: func(serviceMock *translation.DataFederationServiceMock) datafederation.DataFederationService {
				serviceMock.EXPECT().Delete(context.Background(), mock.Anything, mock.Anything).Return(nil)
				return serviceMock
			},
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fooProject", Namespace: "bar"},
			},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.DataFederationSpec{Project: common.ResourceRefNamespaced{Name: "fooProject"}},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
				},
				Spec: akov2.DataFederationSpec{Project: common.ResourceRefNamespaced{Name: "fooProject"}},
			},
			wantSecrets: []corev1.Secret{},
			wantResult:  workflow.OK(),
		},
		{
			name: "federation object with secrets",
			service: func(serviceMock *translation.DataFederationServiceMock) datafederation.DataFederationService {
				serviceMock.EXPECT().Delete(context.Background(), mock.Anything, mock.Anything).Return(nil)
				return serviceMock
			},
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fooProject", Namespace: "bar"},
				Status:     status.AtlasProjectStatus{ID: "123"},
			},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.DataFederationSpec{
					Name:    "data-federation-name",
					Project: common.ResourceRefNamespaced{Name: "fooProject"},
				},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
				},
				Spec: akov2.DataFederationSpec{
					Name:    "data-federation-name",
					Project: common.ResourceRefNamespaced{Name: "fooProject"},
				},
			},
			connectionSecrets: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "fooSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "data-federation-name",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keepSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "some-cluster",
						},
					},
				},
			},
			wantSecrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keepSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "some-cluster",
						},
					},
				},
			},
			wantResult: workflow.OK(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			var objects []client.Object
			if tc.atlasProject != nil {
				objects = append(objects, tc.atlasProject)
			}
			if tc.dataFederation != nil {
				objects = append(objects, tc.dataFederation)
			}
			for _, s := range tc.connectionSecrets {
				objects = append(objects, s)
			}
			scheme := runtime.NewScheme()
			utilruntime.Must(corev1.AddToScheme(scheme))
			utilruntime.Must(akov2.AddToScheme(scheme))
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				Build()
			project := &akov2.AtlasProject{}

			r := &AtlasDataFederationReconciler{
				Client: fakeClient,
				Log:    logger,
			}

			var svc datafederation.DataFederationService
			if tc.service != nil {
				svc = tc.service(translation.NewDataFederationServiceMock(t))
			}
			gotResult := r.handleDelete(ctx, logger, tc.dataFederation, project, svc)
			assert.Equal(t, tc.wantResult, gotResult)

			gotDataFederation := &akov2.AtlasDataFederation{}
			err := fakeClient.Get(ctx.Context, client.ObjectKeyFromObject(tc.dataFederation), gotDataFederation)
			assert.NoError(t, err)
			gotDataFederation.ResourceVersion = ""
			assert.Equal(t, tc.wantDataFederation, gotDataFederation)

			gotSecrets := &corev1.SecretList{}
			err = fakeClient.List(ctx.Context, gotSecrets)
			for i := range gotSecrets.Items {
				gotSecrets.Items[i].ResourceVersion = ""
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantSecrets, gotSecrets.Items)
		})
	}
}

func TestFindAtlasDataFederationForProjects(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "same namespace",
			obj: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "project", Namespace: "ns"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDataFederation{
					ObjectMeta: metav1.ObjectMeta{Name: "adf1", Namespace: "ns"},
					Spec: akov2.DataFederationSpec{
						Project: common.ResourceRefNamespaced{Name: "project"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "adf1", Namespace: "ns"}},
			},
		},
		{
			name: "different namespace",
			obj: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "project", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDataFederation{
					ObjectMeta: metav1.ObjectMeta{Name: "adf1", Namespace: "ns"},
					Spec: akov2.DataFederationSpec{
						Project: common.ResourceRefNamespaced{Name: "project"},
					},
				},
			},
			want: []reconcile.Request{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			idx := indexer.NewAtlasDataFederationByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(idx.Object(), idx.Name(), idx.Keys).
				Build()
			reconciler := &AtlasDataFederationReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findAtlasDataFederationForProjects(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}
