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
	"fmt"
	"strings"
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
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestConnectionSecretReconcile(t *testing.T) {
	type testCase struct {
		reqName          string
		deployment       *akov2.AtlasDeployment
		user             *akov2.AtlasDatabaseUser
		project          *akov2.AtlasProject
		secrets          []client.Object
		expectedDeletion bool
		expectedUpdate   bool
		expectedResult   func() (ctrl.Result, error)
	}

	tests := map[string]testCase{
		"fail: could not load identifiers": {
			reqName: "my-project$cluster",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate("InvalidConnectionSecretName", ErrInternalFormatPartsInvalid).ReconcileResult()
			},
		},
		"success: could not find secret with k8s format": {
			reqName: "test-project-id-cluster1-admin",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: missing deployment and missing user; garbage collect secret": {
			reqName: "test-project-id$cluster1$admin",
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: "default",
						Labels: map[string]string{
							ProjectLabelKey: "test-project-id",
							ClusterLabelKey: "cluster1",
							TypeLabelKey:    "connection",
						},
					},
				},
			},

			// Deletion will internally be done by Kube via ownerRefernce GC
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: only one available resource from the pair, other non-existent": {
			reqName: "test-project-id$cluster1$admin",
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "admin",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "project",
				},
			},
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"requque: resources are not ready yet": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
				},
				Status: status.AtlasDeploymentStatus{},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "default",
				},
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
			expectedResult: func() (ctrl.Result, error) {
				notReady := []string{"AtlasDeployment/deployment"}
				return workflow.InProgress("ConnectionSecretNotReady", fmt.Sprintf("Not ready: %s", strings.Join(notReady, ", "))).ReconcileResult()
			},
		},
		"success: deployment missing triggers handleDelete()": {
			reqName: "test-project-id$cluster1$admin",
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "default",
				},
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
					ObjectMeta: metav1.ObjectMeta{
						Name:      "sdk-creds",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("test-pass"),
						"publicApiKey":  []byte("test-pass"),
						"privateApiKey": []byte("test-pass"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: "default",
						Labels: map[string]string{
							ProjectLabelKey: "test-project-id",
							ClusterLabelKey: "cluster1",
							TypeLabelKey:    "connection",
						},
					},
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: invalid scopes trigger handleDelete()": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
				},
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
				},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "admin",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: "other-cluster",
							Type: "CLUSTER",
						},
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
					ObjectMeta: metav1.ObjectMeta{
						Name:      "sdk-creds",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("test-pass"),
						"publicApiKey":  []byte("test-pass"),
						"privateApiKey": []byte("test-pass"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: "default",
						Labels: map[string]string{
							ProjectLabelKey: "test-project-id",
							ClusterLabelKey: "cluster1",
							TypeLabelKey:    "connection",
						},
					},
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: expired dbuser triggers handleDelete()": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
				},
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
				},
			},
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "admin",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:       "admin",
					PasswordSecret: &common.ResourceRef{Name: "admin-password"},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
					},
					DeleteAfterDate: time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339),
				},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "sdk-creds",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("test-pass"),
						"publicApiKey":  []byte("test-pass"),
						"privateApiKey": []byte("test-pass"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myproject-cluster1-admin",
						Namespace: "default",
						Labels: map[string]string{
							ProjectLabelKey: "test-project-id",
							ClusterLabelKey: "cluster1",
							TypeLabelKey:    "connection",
						},
					},
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: pair ready will call handleUpdate()": {
			reqName: "test-project-id$cluster1$admin",
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-atlas-project",
							Namespace: "default",
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
					Namespace: "default",
				},
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
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-atlas-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "MyProject",
				},
				Status: status.AtlasProjectStatus{
					ID: "test-project-id",
				},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: "default"},
					Data:       map[string][]byte{"password": []byte("test-pass")},
				},
			},
			expectedUpdate: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
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

			compositeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				WithIndex(&akov2.AtlasDeployment{}, indexer.AtlasDeploymentBySpecNameAndProjectID, func(obj client.Object) []string {
					d := obj.(*akov2.AtlasDeployment)
					if d.Spec.DeploymentSpec == nil || d.Spec.DeploymentSpec.Name == "" {
						return nil
					}
					return []string{"test-project-id-" + d.Spec.DeploymentSpec.Name}
				}).
				WithIndex(&akov2.AtlasDatabaseUser{}, indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, func(obj client.Object) []string {
					u := obj.(*akov2.AtlasDatabaseUser)
					if u.Spec.Username == "" {
						return nil
					}
					return []string{"test-project-id-" + u.Spec.Username}
				}).
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
					Client: compositeClient,
					Log:    logger.Sugar(),
					GlobalSecretRef: types.NamespacedName{
						Name:      "global-secret",
						Namespace: "default",
					},
					AtlasProvider: atlasProvider,
				},
				Scheme:        scheme,
				EventRecorder: record.NewFakeRecorder(10),
			}

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      tc.reqName,
				},
			}

			res, err := r.Reconcile(ctx, req)
			expRes, expErr := tc.expectedResult()

			assert.Equal(t, expRes, res)
			if expErr != nil {
				assert.EqualError(t, err, expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedUpdate {
				ids, err := LoadRequestIdentifiers(ctx, compositeClient, req.NamespacedName)
				require.NoError(t, err)
				ids.ProjectName = "myproject"

				expectedName := CreateK8sFormat(ids.ProjectName, ids.ClusterName, ids.DatabaseUsername)
				var outputSecret corev1.Secret
				getErr := compositeClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      expectedName,
				}, &outputSecret)
				assert.NoError(t, getErr, "expected secret %q to exist", expectedName)
			}

			if tc.expectedDeletion {
				ids, err := LoadRequestIdentifiers(ctx, compositeClient, req.NamespacedName)
				require.NoError(t, err)

				expectedName := CreateK8sFormat(ids.ProjectName, ids.ClusterName, ids.DatabaseUsername)
				var check corev1.Secret
				getErr := compositeClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      expectedName,
				}, &check)
				assert.True(t, apiErrors.IsNotFound(getErr), "expected secret %q to be deleted", expectedName)
			}
		})
	}
}

func TestConnectionSecretReconcile_MultiDeploymentMultiUser(t *testing.T) {
	const ns = "default"

	newDeployment := func(name, cluster string) *akov2.AtlasDeployment {
		return &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: cluster},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDeploymentStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
				ConnectionStrings: &status.ConnectionStrings{
					Standard:    fmt.Sprintf("mongodb+srv://%s.mongodb.net", cluster),
					StandardSrv: fmt.Sprintf("mongodb://%s.mongodb.net", cluster),
				},
			},
		}
	}

	newUser := func(username, passwordSecret string) *akov2.AtlasDatabaseUser {
		return &akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{Name: username, Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username:       username,
				PasswordSecret: &common.ResourceRef{Name: passwordSecret},
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "my-atlas-project", Namespace: ns},
				},
			},
			Status: status.AtlasDatabaseUserStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
			},
		}
	}

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Deployments (2)
	deployments := []*akov2.AtlasDeployment{
		newDeployment("dep1", "cluster1"),
		newDeployment("dep2", "cluster2"),
	}

	// Users (3)
	users := []*akov2.AtlasDatabaseUser{
		newUser("admin", "admin-password"),
		newUser("user2", "user2-password"),
		newUser("user3", "user3-password"),
	}

	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "my-atlas-project", Namespace: ns},
		Spec:       akov2.AtlasProjectSpec{Name: "MyProject"},
		Status:     status.AtlasProjectStatus{ID: "test-project-id"},
	}

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

	objs := make([]client.Object, 0, len(deployments)+len(users)+1+len(secrets))
	for _, d := range deployments {
		objs = append(objs, d)
	}
	for _, u := range users {
		objs = append(objs, u)
	}
	objs = append(objs, project)
	objs = append(objs, secrets...)

	clientWithIndexes := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithIndex(&akov2.AtlasDeployment{}, indexer.AtlasDeploymentBySpecNameAndProjectID, func(obj client.Object) []string {
			d := obj.(*akov2.AtlasDeployment)
			if d.Spec.DeploymentSpec == nil || d.Spec.DeploymentSpec.Name == "" {
				return nil
			}
			return []string{"test-project-id-" + d.Spec.DeploymentSpec.Name}
		}).
		WithIndex(&akov2.AtlasDatabaseUser{}, indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, func(obj client.Object) []string {
			u := obj.(*akov2.AtlasDatabaseUser)
			if u.Spec.Username == "" {
				return nil
			}
			return []string{"test-project-id-" + u.Spec.Username}
		}).
		Build()

	r := &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          clientWithIndexes,
			Log:             logger.Sugar(),
			GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: ns},
		},
		Scheme:        scheme,
		EventRecorder: record.NewFakeRecorder(10),
	}

	for _, d := range deployments {
		for _, u := range users {
			reqName := fmt.Sprintf("test-project-id$%s$%s", d.Spec.DeploymentSpec.Name, u.Spec.Username)
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ns,
					Name:      reqName,
				},
			}

			res, err := r.Reconcile(ctx, req)
			assert.NoError(t, err, "Reconcile failed for %s", reqName)
			assert.Equal(t, ctrl.Result{}, res, "Unexpected result for %s", reqName)

			expectedSecretName := fmt.Sprintf("myproject-%s-%s", d.Spec.DeploymentSpec.Name, u.Spec.Username)
			var outputSecret corev1.Secret
			err = clientWithIndexes.Get(ctx, types.NamespacedName{
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
