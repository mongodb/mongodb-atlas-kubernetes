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

package atlasprivateendpoint

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	"go.mongodb.org/atlas-sdk/v20250312012/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

func TestReconcile(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		request        reconcile.Request
		expectedResult reconcile.Result
		expectedLogs   []string
	}{
		"failed to prepare resource": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "pe2"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasPrivateEndpoint reconciliation",
				"Object default/pe2 doesn't exist, was it deleted after reconcile request?",
			},
		},
		"prepare resource for reconciliation": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "pe1"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasPrivateEndpoint reconciliation",
				"-> Skipping AtlasPrivateEndpoint reconciliation as annotation mongodb.com/atlas-reconciliation-policy=skip",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pe := &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
					},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			}

			core, logs := observer.New(zap.DebugLevel)
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(pe).
				Build()
			r := &AtlasPrivateEndpointReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    zap.New(core).Sugar(),
				},
			}
			result, _ := r.Reconcile(ctx, tt.request)
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, len(tt.expectedLogs), logs.Len())
			for i, logMsg := range tt.expectedLogs {
				assert.Equal(t, logMsg, logs.All()[i].Message)
			}
		})
	}
}

func TestEnsureCustomResource(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		provider             atlas.Provider
		expectedResult       reconcile.Result
		wantErr              bool
		expectedLogs         []string
	}{
		"skip custom resource reconciliation": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
					},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Skipping AtlasPrivateEndpoint reconciliation as annotation mongodb.com/atlas-reconciliation-policy=skip",
			},
		},
		"custom resource version is invalid": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
					Labels: map[string]string{
						customresource.ResourceVersion: "wrong",
					},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			wantErr: true,
			expectedLogs: []string{
				"resource version for 'pe1' is invalid",
				"AtlasPrivateEndpoint is invalid: {true 10000000000 wrong is not a valid semver version for label mongodb.com/atlas-resource-version AtlasResourceVersionIsInvalid true false}",
				"Status update",
			},
		},
		"custom resource is not supported": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			provider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return true
				},
				IsSupportedFunc: func() bool {
					return false
				},
			},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"Status update",
			},
		},
		"failed to get project from atlas": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{

					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			provider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return true
				},
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return nil, errors.New("failed to create sdk client")
				},
			},
			wantErr: true,
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"resource *v1.AtlasPrivateEndpoint(default/pe1) failed on condition Ready: failed to create sdk client",
				"Status update",
			},
		},
		"failed to get project from cluster": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			provider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return true
				},
				IsSupportedFunc: func() bool {
					return true
				},
			},
			wantErr: true,
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"resource *v1.AtlasPrivateEndpoint(default/pe1) failed on condition Ready: error resolving project reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-project\" not found",
				"Status update",
			},
		},
		"custom resource is ready for reconciliation": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "INITIATING",
				},
			},
			provider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return true
				},
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetGroup(mock.Anything, projectID).Return(admin.GetGroupApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetGroupExecute(mock.AnythingOfType("admin.GetGroupApiRequest")).
						Return(&admin.Group{Id: &projectID}, nil, nil)

					peAPI := mockadmin.NewPrivateEndpointServicesApi(t)
					peAPI.EXPECT().GetPrivateEndpointService(mock.Anything, projectID, "AWS", "pe-service-id").
						Return(admin.GetPrivateEndpointServiceApiRequest{ApiService: peAPI})
					peAPI.EXPECT().GetPrivateEndpointServiceExecute(mock.AnythingOfType("admin.GetPrivateEndpointServiceApiRequest")).
						Return(
							&admin.EndpointService{
								Id:            pointer.MakePtr("pe-service-id"),
								CloudProvider: "AWS",
								RegionName:    pointer.MakePtr("US_EAST_1"),
								Status:        pointer.MakePtr("INITIATING"),
							},
							nil,
							nil,
						)

					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: projectAPI, PrivateEndpointServicesApi: peAPI},
					}, nil
				},
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"Status update",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasPrivateEndpoint, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("orgId"),
						"publicApiKey":  []byte("publicApiKey"),
						"privateApiKey": []byte("privateApiKey"),
					},
				}).
				WithStatusSubresource(tt.atlasPrivateEndpoint).
				Build()
			r := &AtlasPrivateEndpointReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        fakeClient,
					Log:           zap.New(core).Sugar(),
					AtlasProvider: tt.provider,
				},
				EventRecorder: record.NewFakeRecorder(10),
			}
			result, err := r.ensureCustomResource(ctx, tt.atlasPrivateEndpoint)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, len(tt.expectedLogs), logs.Len())
			for i, logMsg := range tt.expectedLogs {
				assert.True(t, strings.Contains(logs.All()[i].Message, logMsg[:len(logMsg)-1]),
					"log:'%s'\ndoesn't contain entry:\n'%s", logs.All()[i].Message, logMsg)
			}
		})
	}
}

func TestFailManageFinalizer(t *testing.T) {
	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))
	pe := &akov2.AtlasPrivateEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pe1",
			Namespace: "default",
		},
		Spec: akov2.AtlasPrivateEndpointSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: "project-id",
				},
				ConnectionSecret: &api.LocalObjectReference{},
			},
			Provider: "AWS",
			Region:   "US_EAST_1",
		},
		Status: status.AtlasPrivateEndpointStatus{
			ServiceID:     "pe-service-id",
			ServiceStatus: "AVAILABLE",
		},
	}
	atlasPE := &privateendpoint.AWSService{
		CommonEndpointService: privateendpoint.CommonEndpointService{
			ID:            "pe-service-id",
			CloudRegion:   "US_EAST_1",
			ServiceStatus: "AVAILABLE",
		},
		ServiceName: "aws/service/name",
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(pe).
		WithStatusSubresource(pe).
		WithInterceptorFuncs(interceptor.Funcs{
			Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				return errors.New("failed to manage finalizer")
			},
		}).
		Build()

	logger := zaptest.NewLogger(t).Sugar()
	workflowCtx := &workflow.Context{
		Log: logger,
	}

	tests := map[string]struct {
		transition func(r *AtlasPrivateEndpointReconciler) (reconcile.Result, error)
	}{
		"failed to manage finalizer when transitioning to in progress": {
			transition: func(r *AtlasPrivateEndpointReconciler) (reconcile.Result, error) {
				return r.inProgress(
					workflowCtx,
					pe,
					atlasPE,
					api.PrivateEndpointServiceReady,
					workflow.PrivateEndpointServiceInitializing,
					"testing transition",
				)
			},
		},
		"failed to manage finalizer when transitioning to ready": {
			transition: func(r *AtlasPrivateEndpointReconciler) (reconcile.Result, error) {
				return r.ready(workflowCtx, pe, atlasPE)
			},
		},
		"failed to manage finalizer when transitioning to waiting configuration": {
			transition: func(r *AtlasPrivateEndpointReconciler) (reconcile.Result, error) {
				return r.waitForConfiguration(workflowCtx, pe, atlasPE)
			},
		},
		"failed to manage finalizer when transitioning to unmanage": {
			transition: func(r *AtlasPrivateEndpointReconciler) (reconcile.Result, error) {
				return r.unmanage(workflowCtx, pe)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &AtlasPrivateEndpointReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    zaptest.NewLogger(t).Sugar(),
				},
			}

			_, err := tt.transition(r)
			assert.Error(t, err)
		})
	}
}
