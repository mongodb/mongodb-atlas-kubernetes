package atlasprivateendpoint

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{},
					Provider:              "AWS",
					Region:                "US_EAST_1",
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
				Client: fakeClient,
				Log:    zap.New(core).Sugar(),
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

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		provider             atlas.Provider
		expectedResult       reconcile.Result
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{},
					Provider:              "AWS",
					Region:                "US_EAST_1",
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{},
					Provider:              "AWS",
					Region:                "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{},
					Provider:              "AWS",
					Region:                "US_EAST_1",
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
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
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to create sdk client")
				},
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"resource *v1.AtlasPrivateEndpoint(default/pe1) failed on condition Ready: failed to create Atlas SDK client: failed to create sdk client",
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
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
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
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedLogs: []string{
				"resource 'pe1' version is valid",
				"resource *v1.AtlasPrivateEndpoint(default/pe1) failed on condition Ready: failed to retrieve project custom resource: atlasprojects.atlas.mongodb.com \"my-project\" not found",
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
					ExternalProject: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
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
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProject(ctx, projectID).Return(admin.GetProjectApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(&admin.Group{Id: &projectID}, nil, nil)

					peAPI := mockadmin.NewPrivateEndpointServicesApi(t)
					peAPI.EXPECT().GetPrivateEndpointService(ctx, projectID, "AWS", "pe-service-id").
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

					return &admin.APIClient{ProjectsApi: projectAPI, PrivateEndpointServicesApi: peAPI}, "", nil
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
				WithObjects(tt.atlasPrivateEndpoint).
				WithStatusSubresource(tt.atlasPrivateEndpoint).
				Build()
			r := &AtlasPrivateEndpointReconciler{
				Client:        fakeClient,
				AtlasProvider: tt.provider,
				EventRecorder: record.NewFakeRecorder(10),
				Log:           zap.New(core).Sugar(),
			}
			result, err := r.ensureCustomResource(ctx, tt.atlasPrivateEndpoint)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, len(tt.expectedLogs), logs.Len())
			for i, logMsg := range tt.expectedLogs {
				assert.Equal(t, logMsg, logs.All()[i].Message)
			}
		})
	}
}

func TestGetProjectFromKube(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		provider             atlas.Provider
		expectedProject      *project.Project
		expectedErr          error
	}{
		"failed to resolve secret": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-missing-project",
						Namespace: "default",
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			expectedErr: fmt.Errorf(
				"failed to retrieve project custom resource: %w",
				&k8serrors.StatusError{
					ErrStatus: metav1.Status{
						Status:  "Failure",
						Message: "atlasprojects.atlas.mongodb.com \"my-missing-project\" not found",
						Reason:  "NotFound",
						Details: &metav1.StatusDetails{
							Name:  "my-missing-project",
							Group: "atlas.mongodb.com",
							Kind:  "atlasprojects",
						},
						Code: 404,
					},
				},
			),
		},
		"failed to create sdk client": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
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
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to create sdk client")
				},
			},
			expectedErr: fmt.Errorf("failed to create Atlas SDK client: %w", errors.New("failed to create sdk client")),
		},
		"get project from cluster": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
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
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{}, "org-id", nil
				},
			},
			expectedProject: &project.Project{
				OrgID: "org-id",
				ID:    "project-id",
				Name:  "My Project",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			akoProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "My Project",
				},
				Status: status.AtlasProjectStatus{
					ID: projectID,
				},
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(akoProject).
				Build()
			r := &AtlasPrivateEndpointReconciler{
				Client:        fakeClient,
				AtlasProvider: tt.provider,
			}
			p, err := r.getProjectFromKube(ctx, tt.atlasPrivateEndpoint)
			assert.Equal(t, tt.expectedProject, p)
			assert.Equal(t, tt.expectedErr, err)
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
			ExternalProject: &akov2.ExternalProjectReference{
				ID: "project-id",
			},
			LocalCredentialHolder: api.LocalCredentialHolder{},
			Provider:              "AWS",
			Region:                "US_EAST_1",
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
				Client: fakeClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}

			result, err := tt.transition(r)
			assert.NoError(t, err)
			assert.Equal(t, reconcile.Result{RequeueAfter: workflow.DefaultRetry}, result)
		})
	}
}
