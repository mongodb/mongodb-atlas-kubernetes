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

package connectionsecret

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	admin "go.mongodb.org/atlas-sdk/v20250312002/admin"
	"go.mongodb.org/atlas-sdk/v20250312002/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestConnectionSecretReconcile(t *testing.T) {
	type testCase struct {
		reqName            string
		deployment         *akov2.AtlasDeployment
		user               *akov2.AtlasDatabaseUser
		project            *akov2.AtlasProject
		secrets            []client.Object
		expectedSecretName string
		expectedResult     ctrl.Result
		expectErr          bool
	}

	const ns = "default"

	tests := map[string]testCase{
		"fail: invalid secret name format": {
			reqName:   "invalid-format",
			expectErr: true,
		},
		"fail: K8s format secret with missing secret": {
			reqName:   "myproject-cluster1-admin",
			expectErr: true,
		},
		"fail: missing deployment": {
			reqName: "test-project-id$cluster1$admin",
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: ns},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: "DatabaseUserReady", Status: "True"}},
					},
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
					Data:       map[string][]byte{"password": []byte("test-pass")},
				},
			},
			expectErr: true,
		},
		"fail: K8s format secret with no labels": {
			reqName: "myproject-cluster1-admin",
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: ns,
						Labels:    map[string]string{},
					},
				},
			},
			expectErr: true,
		},
		"requque: resources are not ready yet": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDeploymentStatus{}, // Not ready
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: ns},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
					Data:       map[string][]byte{"password": []byte("test-pass")},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "sdk-creds", Namespace: ns},
					Data: map[string][]byte{
						"orgId":         []byte("test-pass"),
						"publicApiKey":  []byte("test-pass"),
						"privateApiKey": []byte("test-pass"),
					},
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: 10 * time.Second},
			expectErr:      false,
		},
		"success: deployment uses ProjectRef, user uses ExternalProjectRef with internal format for request": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-atlas-project",
							Namespace: ns,
						},
					},
				},
				Status: status.AtlasDeploymentStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
					ConnectionStrings: &status.ConnectionStrings{
						Standard:    "mongodb+srv://cluster1.mongodb.net",
						StandardSrv: "mongodb://cluster1.mongodb.net",
					},
				},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: ns},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "my-atlas-project", Namespace: ns},
				Spec: akov2.AtlasProjectSpec{
					Name: "MyProject",
				},
				Status: status.AtlasProjectStatus{
					ID: "test-project-id",
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
					Data:       map[string][]byte{"password": []byte("test-pass")},
				},
			},
			expectedSecretName: "myproject-cluster1-admin",
			expectedResult:     ctrl.Result{RequeueAfter: 30 * time.Second},
			expectErr:          false,
		},
		"success: both deployment and user use ProjectRef with internal format for request": {
			reqName: "myproject-cluster1-admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dep1",
					Namespace: ns,
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-atlas-project",
							Namespace: ns,
						},
					},
				},
				Status: status.AtlasDeploymentStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
					ConnectionStrings: &status.ConnectionStrings{
						Standard:    "mongodb+srv://cluster1.mongodb.net",
						StandardSrv: "mongodb://cluster1.mongodb.net",
					},
				},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "admin",
					Namespace: ns,
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-atlas-project",
							Namespace: ns,
						},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-atlas-project",
					Namespace: ns,
				},
				Spec:   akov2.AtlasProjectSpec{Name: "myproject"},
				Status: status.AtlasProjectStatus{ID: "test-project-id"},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "admin-password",
						Namespace: ns,
					},
					Data: map[string][]byte{"password": []byte("test-pass")},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: ns,
						Labels: map[string]string{
							ProjectLabelKey: "test-project-id",
							ClusterLabelKey: "cluster1",
							TypeLabelKey:    "connection",
						},
					},
				},
			},
			expectedSecretName: "myproject-cluster1-admin",
			expectedResult:     ctrl.Result{RequeueAfter: 30 * time.Second},
			expectErr:          false,
		},
		"success: both deployment and user use ExternalRef (SDK required) with internal format for request": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDeploymentStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
					ConnectionStrings: &status.ConnectionStrings{
						Standard:    "mongodb+srv://cluster1.mongodb.net",
						StandardSrv: "mongodb://cluster1.mongodb.net",
					},
				},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: ns},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
					Data:       map[string][]byte{"password": []byte("test-pass")},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "sdk-creds", Namespace: ns},
					Data: map[string][]byte{
						"orgId":         []byte("test-pass"),
						"publicApiKey":  []byte("test-pass"),
						"privateApiKey": []byte("test-pass"),
					},
				},
			},
			expectedSecretName: "myproject-cluster1-admin",
			expectedResult:     ctrl.Result{RequeueAfter: 30 * time.Second},
			expectErr:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))
			require.NoError(t, akov2.AddToScheme(scheme))

			logger := zaptest.NewLogger(t)
			ctx := context.Background()

			objects := make([]client.Object, 0, 3)
			if tc.deployment != nil {
				objects = append(objects, tc.deployment)
			}
			if tc.user != nil {
				objects = append(objects, tc.user)
			}
			if tc.project != nil {
				objects = append(objects, tc.project)
			}
			objects = append(objects, tc.secrets...)

			preClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				Build()

			indexerDeployment := indexer.NewAtlasDeploymentBySpecNameIndexer(ctx, preClient, logger)
			indexerDatabaseUser := indexer.NewAtlasDatabaseUserBySpecUsernameIndexer(ctx, preClient, logger)

			client := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				WithIndex(indexerDeployment.Object(), indexerDeployment.Name(), indexerDeployment.Keys).
				WithIndex(indexerDatabaseUser.Object(), indexerDatabaseUser.Name(), indexerDatabaseUser.Keys).
				Build()

			atlasProvider := &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsApi(t)

					projectAPI.EXPECT().
						GetProject(mock.Anything, "test-project-id").
						Return(admin.GetProjectApiRequest{ApiService: projectAPI})

					projectAPI.EXPECT().
						GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(&admin.Group{
							Id:   pointer.MakePtr("test-project-id"),
							Name: "MyProject",
						}, nil, nil)

					return &atlas.ClientSet{
						SdkClient20250312002: &admin.APIClient{
							ProjectsApi: projectAPI,
						},
					}, nil
				},
				IsSupportedFunc: func() bool { return true },
				IsCloudGovFunc:  func() bool { return false },
			}

			r := &ConnectionSecretReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:          client,
					Log:             logger.Sugar(),
					GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: ns},
					AtlasProvider:   atlasProvider,
				},
				EventRecorder: record.NewFakeRecorder(10),
			}

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      tc.reqName,
				},
			}

			res, err := r.Reconcile(ctx, req)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)

				if tc.expectedSecretName != "" {
					var outputSecret corev1.Secret
					err := client.Get(ctx, types.NamespacedName{
						Namespace: ns,
						Name:      tc.expectedSecretName,
					}, &outputSecret)
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestConnectionSecretReconcile_MultiDeploymentMultiUser(t *testing.T) {
	const ns = "default"

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Prepare deployments (2)
	deployments := []*akov2.AtlasDeployment{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDeploymentStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
				ConnectionStrings: &status.ConnectionStrings{
					Standard:    "mongodb+srv://cluster1.mongodb.net",
					StandardSrv: "mongodb://cluster1.mongodb.net",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "dep2", Namespace: ns},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster2"},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDeploymentStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
				ConnectionStrings: &status.ConnectionStrings{
					Standard:    "mongodb+srv://cluster2.mongodb.net",
					StandardSrv: "mongodb://cluster2.mongodb.net",
				},
			},
		},
	}

	// Prepare users (3)
	users := []*akov2.AtlasDatabaseUser{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username:       "admin",
				PasswordSecret: &common.ResourceRef{Name: "admin-password"},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDatabaseUserStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "user2", Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username:       "user2",
				PasswordSecret: &common.ResourceRef{Name: "user2-password"},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDatabaseUserStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "user3", Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username:       "user3",
				PasswordSecret: &common.ResourceRef{Name: "user3-password"},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDatabaseUserStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
			},
		},
	}

	// Project resource
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "my-atlas-project", Namespace: ns},
		Spec:       akov2.AtlasProjectSpec{Name: "MyProject"},
		Status:     status.AtlasProjectStatus{ID: "test-project-id"},
	}

	// Secrets for user passwords
	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
			Data:       map[string][]byte{"password": []byte("adminpass")},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "user2-password", Namespace: ns},
			Data:       map[string][]byte{"password": []byte("user2pass")},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "user3-password", Namespace: ns},
			Data:       map[string][]byte{"password": []byte("user3pass")},
		},
	}

	objs := make([]client.Object, 0, 6)
	for _, d := range deployments {
		objs = append(objs, d)
	}
	for _, u := range users {
		objs = append(objs, u)
	}
	objs = append(objs, project)
	objs = append(objs, secrets...)

	preClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		Build()

	indexerDeployment := indexer.NewAtlasDeploymentBySpecNameIndexer(ctx, preClient, logger)
	indexerDatabaseUser := indexer.NewAtlasDatabaseUserBySpecUsernameIndexer(ctx, preClient, logger)

	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithIndex(indexerDeployment.Object(), indexerDeployment.Name(), indexerDeployment.Keys).
		WithIndex(indexerDatabaseUser.Object(), indexerDatabaseUser.Name(), indexerDatabaseUser.Keys).
		Build()

	atlasProvider := &atlasmock.TestProvider{
		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
			projectAPI := mockadmin.NewProjectsApi(t)

			projectAPI.EXPECT().
				GetProject(mock.Anything, "test-project-id").
				Return(admin.GetProjectApiRequest{ApiService: projectAPI})

			projectAPI.EXPECT().
				GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
				Return(&admin.Group{
					Id:   pointer.MakePtr("test-project-id"),
					Name: "MyProject",
				}, nil, nil)

			return &atlas.ClientSet{
				SdkClient20250312002: &admin.APIClient{
					ProjectsApi: projectAPI,
				},
			}, nil
		},
		IsSupportedFunc: func() bool { return true },
		IsCloudGovFunc:  func() bool { return false },
	}

	r := &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          client,
			Log:             logger.Sugar(),
			GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: ns},
			AtlasProvider:   atlasProvider,
		},
		EventRecorder: record.NewFakeRecorder(10),
	}

	for _, d := range deployments {
		for _, u := range users {
			reqName := "test-project-id$" + d.Spec.DeploymentSpec.Name + "$" + u.Spec.Username
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      reqName,
				},
			}

			res, err := r.Reconcile(ctx, req)
			assert.NoError(t, err, "Reconcile failed for %s", reqName)
			assert.Equal(t, ctrl.Result{RequeueAfter: 30 * time.Second}, res, "Unexpected result for %s", reqName)

			expectedSecretName := "myproject-" + d.Spec.DeploymentSpec.Name + "-" + u.Spec.Username
			var outputSecret corev1.Secret
			err = client.Get(ctx, types.NamespacedName{
				Namespace: ns,
				Name:      expectedSecretName,
			}, &outputSecret)
			assert.NoError(t, err, "Secret not found for %s", reqName)
		}
	}
}

func TestGenerateConnectionSecretRequests(t *testing.T) {
	const ns = "default"
	const projectID = "test-project-id"

	deployment := func(name string) akov2.AtlasDeployment {
		return akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: name},
			},
		}
	}

	user := func(username string, scopes ...string) akov2.AtlasDatabaseUser {
		resScopes := make([]akov2.ScopeSpec, 0, len(scopes))
		for _, s := range scopes {
			resScopes = append(resScopes, akov2.ScopeSpec{
				Type: akov2.DeploymentScopeType,
				Name: s,
			})
		}
		return akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{Name: username, Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username: username,
				Scopes:   resScopes,
			},
		}
	}

	tests := map[string]struct {
		deployments []akov2.AtlasDeployment
		users       []akov2.AtlasDatabaseUser
		expected    []reconcile.Request
	}{
		"no deployments or users": {
			deployments: nil,
			users:       nil,
			expected:    nil,
		},
		"deployment but no users": {
			deployments: []akov2.AtlasDeployment{deployment("cluster1")},
			users:       nil,
			expected:    nil,
		},
		"users and deployments but all scopes mismatched": {
			deployments: []akov2.AtlasDeployment{
				deployment("cluster1"),
				deployment("cluster2"),
			},
			users: []akov2.AtlasDatabaseUser{
				user("user1", "other1"),
				user("user2", "other2"),
			},
			expected: nil,
		},
		"users and deployments with valid scopes (including global)": {
			deployments: []akov2.AtlasDeployment{
				deployment("cluster1"),
				deployment("cluster2"),
				deployment("cluster3"),
			},
			users: []akov2.AtlasDatabaseUser{
				user("admin", "cluster1", "cluster2"),
				user("user2", "cluster1"),
				user("user3", "cluster2"),
				user("user4", "other"),
				user("global"),
			},
			expected: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster1", "admin"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster2", "admin"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster1", "user2"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster2", "user3"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster1", "global"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster2", "global"),
				}},
				{NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      CreateInternalFormat(projectID, "cluster3", "global"),
				}},
			},
		},
	}

	r := &ConnectionSecretReconciler{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.generateConnectionSecretRequests(projectID, tc.deployments, tc.users)
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}

func TestNewDeploymentMapFunc(t *testing.T) {
	const ns = "default"
	const projectID = "test-project-id"

	scheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(scheme))

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	deployment := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{Name: "my-project", Namespace: ns},
			},
		},
		Status: status.AtlasDeploymentStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
			},
		},
	}

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user1",
			Scopes: []akov2.ScopeSpec{
				{Name: "cluster1", Type: akov2.DeploymentScopeType},
			},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{Name: "my-project", Namespace: ns},
			},
		},
	}

	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: ns},
		Status:     status.AtlasProjectStatus{ID: projectID},
	}

	objects := []client.Object{deployment, user, project}

	preClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build()

	userIndexer := indexer.NewAtlasDatabaseUserByProjectIndexer(ctx, preClient, logger)
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		WithIndex(userIndexer.Object(), userIndexer.Name(), userIndexer.Keys).
		Build()

	r := &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: client,
			Log:    logger.Sugar(),
		},
	}

	reqs := r.newDeploymentMapFunc(ctx, deployment)
	require.Len(t, reqs, 1)
	assert.Equal(t, types.NamespacedName{
		Namespace: ns,
		Name:      CreateInternalFormat(projectID, "cluster1", "user1"),
	}, reqs[0].NamespacedName)
}

func TestNewDatabaseUserMapFunc(t *testing.T) {
	const ns = "default"
	const projectID = "test-project-id"

	scheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(scheme))

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user1",
			Scopes: []akov2.ScopeSpec{
				{Name: "cluster1", Type: akov2.DeploymentScopeType},
			},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{Name: "my-project", Namespace: ns},
			},
		},
	}

	deployment := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{Name: "my-project", Namespace: ns},
			},
		},
	}

	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: ns},
		Status:     status.AtlasProjectStatus{ID: projectID},
	}

	objects := []client.Object{deployment, user, project}

	preClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build()

	depIndexer := indexer.NewAtlasDeploymentByProjectIndexer(ctx, preClient, logger)
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		WithIndex(depIndexer.Object(), depIndexer.Name(), depIndexer.Keys).
		Build()

	r := &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: client,
			Log:    logger.Sugar(),
		},
	}

	reqs := r.newDatabaseUserMapFunc(ctx, user)
	require.Len(t, reqs, 1)
	assert.Equal(t, types.NamespacedName{
		Namespace: ns,
		Name:      CreateInternalFormat(projectID, "cluster1", "user1"),
	}, reqs[0].NamespacedName)
}
