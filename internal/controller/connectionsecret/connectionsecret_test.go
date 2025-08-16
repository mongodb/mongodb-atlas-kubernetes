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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	admin "go.mongodb.org/atlas-sdk/v20250312002/admin"
	"go.mongodb.org/atlas-sdk/v20250312002/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func Test_resolveProjectName(t *testing.T) {
	type expectedResult struct {
		expectedProjectName string
		expectedError       error
	}

	projectName := "project-name"
	projectID := "test-project-id"

	type testCase struct {
		ids     ConnSecretIdentifiers
		pair    ConnSecretPair
		project *akov2.AtlasProject
		secrets []client.Object
		result  expectedResult
	}

	tests := map[string]testCase{
		"fail: missing deployment and missing user": {
			result: expectedResult{
				expectedProjectName: "",
				expectedError:       fmt.Errorf("unable to resolve ProjectName"),
			},
		},
		"fail: missing connectionSecret on deployment and user": {
			pair: ConnSecretPair{
				ProjectID: projectID,
				Deployment: &akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dep1",
						Namespace: "default",
					},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						},
					},
				},
				User: &akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "admin",
						Namespace: "default",
					},
					Spec: akov2.AtlasDatabaseUserSpec{
						Username:       "admin",
						PasswordSecret: &common.ResourceRef{Name: "admin-password"},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{ID: "test-project-id"},
						},
					},
				},
			},
			result: expectedResult{
				expectedProjectName: "",
				expectedError:       fmt.Errorf("error getting credentials from project reference: failed to read Atlas API credentials from the secret default/global-secret: secrets \"global-secret\" not found"),
			},
		},
		"success: projectName is already present": {
			ids: ConnSecretIdentifiers{
				ProjectName: projectName,
			},

			result: expectedResult{
				expectedProjectName: projectName,
				expectedError:       nil,
			},
		},
		"success: resolve project name via deployment and project": {
			pair: ConnSecretPair{
				Deployment: &akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "admin",
						Namespace: "default",
					},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "my-project",
								Namespace: "default",
							},
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
					Name: projectName,
				},
			},
			result: expectedResult{
				expectedProjectName: projectName,
				expectedError:       nil,
			},
		},
		"success: resolve project name via user and project": {
			pair: ConnSecretPair{
				User: &akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "admin",
						Namespace: "default",
					},
					Spec: akov2.AtlasDatabaseUserSpec{
						Username: "admin",
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "my-project",
								Namespace: "default",
							},
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
					Name: projectName,
				},
			},
			result: expectedResult{
				expectedProjectName: projectName,
				expectedError:       nil,
			},
		},
		"success: resolve via deployment SDK": {
			pair: ConnSecretPair{
				ProjectID: projectID,
				Deployment: &akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dep1",
						Namespace: "default",
					},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{ID: projectID},
							ConnectionSecret:   &api.LocalObjectReference{Name: "sdk-creds"},
						},
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
			},
			result: expectedResult{
				expectedProjectName: projectName,
				expectedError:       nil,
			},
		},
		"success: resolve via user SDK": {
			pair: ConnSecretPair{
				ProjectID: projectID,
				User: &akov2.AtlasDatabaseUser{
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
			},
			result: expectedResult{
				expectedProjectName: projectName,
				expectedError:       nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			utilruntime.Must(corev1.AddToScheme(scheme))
			utilruntime.Must(akov2.AddToScheme(scheme))

			objs := []client.Object{}
			if tc.project != nil {
				objs = append(objs, tc.project)
			}
			if tc.pair.User != nil {
				objs = append(objs, tc.pair.User)
			}
			if tc.pair.Deployment != nil {
				objs = append(objs, tc.pair.Deployment)
			}
			if tc.secrets != nil {
				objs = append(objs, tc.secrets...)
			}

			logger := zaptest.NewLogger(t)
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()

			atlasProvider := &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsApi(t)

					projectAPI.EXPECT().
						GetProject(mock.Anything, projectID).
						Return(admin.GetProjectApiRequest{ApiService: projectAPI})

					projectAPI.EXPECT().
						GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(&admin.Group{
							Id:   pointer.MakePtr(projectID),
							Name: projectName,
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
					Client:          fakeClient,
					Log:             logger.Sugar(),
					GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: "default"},
					AtlasProvider:   atlasProvider,
				},
				EventRecorder: record.NewFakeRecorder(10),
			}

			gotName, err := r.resolveProjectName(context.Background(), &tc.ids, &tc.pair)

			require.Equal(t, tc.result.expectedProjectName, gotName)
			if tc.result.expectedError != nil {
				require.EqualError(t, err, tc.result.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_handleDelete(t *testing.T) {
	type expectedResult struct {
		expectedResult ctrl.Result
		expectedError  error
	}

	const (
		ns          = "default"
		cluster     = "cluster1"
		username    = "admin"
		projectID   = "test-project-id"
		projectName = "myproject"
	)

	type testCase struct {
		ids     ConnSecretIdentifiers
		pair    ConnSecretPair
		project *akov2.AtlasProject
		secrets []client.Object
		result  expectedResult
	}

	tests := map[string]testCase{
		"fail: unresolved project name": {
			ids: ConnSecretIdentifiers{
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  fmt.Errorf("project name is empty"),
			},
		},
		"success: no secret present": {
			ids: ConnSecretIdentifiers{
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
		"success: delete existing secret without resolution": {
			ids: ConnSecretIdentifiers{
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      CreateK8sFormat(projectName, cluster, username),
						Namespace: ns,
					},
				},
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
		"success: delete project with resolution": {
			ids: ConnSecretIdentifiers{
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User: &akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: username, Namespace: ns},
					Spec: akov2.AtlasDatabaseUserSpec{
						Username: username,
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "my-project",
								Namespace: ns,
							},
						},
					},
				},
				Deployment: &akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cluster,
						Namespace: ns,
					},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: cluster},
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "my-project",
								Namespace: ns,
							},
						},
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: ns},
				Spec:       akov2.AtlasProjectSpec{Name: projectName},
			},
			secrets: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      CreateK8sFormat(projectName, cluster, username),
						Namespace: ns,
					},
				},
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			utilruntime.Must(corev1.AddToScheme(scheme))
			utilruntime.Must(akov2.AddToScheme(scheme))

			objects := make([]client.Object, 0, 4)
			if tc.project != nil {
				objects = append(objects, tc.project)
			}
			if tc.pair.User != nil {
				objects = append(objects, tc.pair.User)
			}
			if tc.pair.Deployment != nil {
				objects = append(objects, tc.pair.Deployment)
			}
			objects = append(objects, tc.secrets...)

			logger := zaptest.NewLogger(t)
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()

			r := &ConnectionSecretReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:          fakeClient,
					Log:             logger.Sugar(),
					GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: ns},
				},
				EventRecorder: record.NewFakeRecorder(10),
			}

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{Namespace: ns, Name: "any"},
			}

			res, err := r.handleDelete(context.Background(), req, &tc.ids, &tc.pair)
			assert.Equal(t, tc.result.expectedResult, res)
			if tc.result.expectedError != nil {
				require.EqualError(t, err, tc.result.expectedError.Error())
				return
			}
			require.NoError(t, err)

			if tc.ids.ClusterName != "" && tc.ids.DatabaseUsername != "" {
				projectName := tc.ids.ProjectName
				if projectName == "" && tc.project != nil {
					projectName = tc.project.Spec.Name
				}
				if projectName != "" {
					name := CreateK8sFormat(projectName, tc.ids.ClusterName, tc.ids.DatabaseUsername)
					var s corev1.Secret
					getErr := fakeClient.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: name}, &s)
					require.True(t, apierrors.IsNotFound(getErr), "expected secret %s to be deleted", name)
				}
			}
		})
	}
}

func Test_handleUpsert(t *testing.T) {
	type expectedResult struct {
		expectedResult ctrl.Result
		expectedError  error
	}

	const (
		ns          = "default"
		cluster     = "cluster1"
		username    = "admin"
		projectID   = "test-project-id"
		projectName = "myproject"
	)

	newDeployment := func() *akov2.AtlasDeployment {
		return &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: cluster},
			},
			Status: status.AtlasDeploymentStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
				ConnectionStrings: &status.ConnectionStrings{
					Standard:    "mongodb://cluster1.mongodb.net/?authSource=admin",
					StandardSrv: "mongodb+srv://cluster1.mongodb.net/?authSource=admin",
					PrivateEndpoint: []status.PrivateEndpoint{
						{
							ConnectionString:                  "mongodb://pe1.mongodb.net",
							SRVConnectionString:               "mongodb+srv://pe1.mongodb.net",
							SRVShardOptimizedConnectionString: "mongodb+srv://pe1-shard.mongodb.net",
						},
						{
							ConnectionString:                  "mongodb://pe2.mongodb.net",
							SRVConnectionString:               "mongodb+srv://pe2.mongodb.net",
							SRVShardOptimizedConnectionString: "mongodb+srv://pe2-shard.mongodb.net",
						},
					},
				},
			},
		}
	}

	newUser := func() *akov2.AtlasDatabaseUser {
		return &akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{Name: username, Namespace: ns},
			Spec: akov2.AtlasDatabaseUserSpec{
				Username:       username,
				PasswordSecret: &common.ResourceRef{Name: "admin-password"},
			},
			Status: status.AtlasDatabaseUserStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
				},
			},
		}
	}

	newPasswordSecret := func() *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "admin-password", Namespace: ns},
			Data:       map[string][]byte{passwordKey: []byte("test-pass")},
		}
	}

	newExistingConnSecret := func() *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      CreateK8sFormat(projectName, cluster, username),
				Namespace: ns,
				Labels: map[string]string{
					TypeLabelKey:    CredLabelVal,
					ProjectLabelKey: projectID,
					ClusterLabelKey: cluster,
				},
			},
			Data: map[string][]byte{
				userNameKey:    []byte("beforeusername"),
				passwordKey:    []byte("beforepassword"),
				standardKey:    []byte("mongodb://cluster1.mongodb.net/?authSource=admin"),
				standardKeySrv: []byte("mongodb+srv://cluster1.mongodb.net/?authSource=admin"),
			},
		}
	}

	tests := map[string]struct {
		ids     ConnSecretIdentifiers
		pair    ConnSecretPair
		project *akov2.AtlasProject
		secrets []client.Object
		result  expectedResult
	}{
		"fail: unresolved project name": {
			ids: ConnSecretIdentifiers{
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  fmt.Errorf("project name is empty"),
			},
		},
		"success: test create": {
			ids: ConnSecretIdentifiers{
				ProjectID:        projectID,
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				Deployment: newDeployment(),
				User:       newUser(),
			},
			secrets: []client.Object{newPasswordSecret()},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
		"success: test update": {
			ids: ConnSecretIdentifiers{
				ProjectID:        projectID,
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				Deployment: newDeployment(),
				User:       newUser(),
			},
			secrets: []client.Object{
				newPasswordSecret(),
				newExistingConnSecret(),
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			utilruntime.Must(corev1.AddToScheme(scheme))
			utilruntime.Must(akov2.AddToScheme(scheme))

			objects := make([]client.Object, 0, len(tc.secrets)+3)
			if tc.project != nil {
				objects = append(objects, tc.project)
			}
			if tc.pair.User != nil {
				objects = append(objects, tc.pair.User)
			}
			if tc.pair.Deployment != nil {
				objects = append(objects, tc.pair.Deployment)
			}
			objects = append(objects, tc.secrets...)

			logger := zaptest.NewLogger(t)
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()

			r := &ConnectionSecretReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:          fakeClient,
					Log:             logger.Sugar(),
					GlobalSecretRef: types.NamespacedName{Name: "global-secret", Namespace: ns},
				},
				Scheme:        scheme,
				EventRecorder: record.NewFakeRecorder(10),
			}

			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: ns, Name: "any"}}

			res, err := r.handleUpsert(context.Background(), req, &tc.ids, &tc.pair)
			assert.Equal(t, tc.result.expectedResult, res)
			if tc.result.expectedError != nil {
				require.EqualError(t, err, tc.result.expectedError.Error())
				return
			}
			require.NoError(t, err)

			secretName := CreateK8sFormat(projectName, cluster, username)
			var s corev1.Secret
			getErr := fakeClient.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: secretName}, &s)
			require.NoError(t, getErr)

			require.Equal(t, CredLabelVal, s.Labels[TypeLabelKey])
			require.Equal(t, projectID, s.Labels[ProjectLabelKey])
			require.Equal(t, cluster, s.Labels[ClusterLabelKey])

			require.Equal(t, username, string(s.Data[userNameKey]))
			require.Equal(t, "test-pass", string(s.Data[passwordKey]))

			// Verify all connection string variants
			urlsToCheck := map[string]string{
				standardKey:    "mongodb://cluster1.mongodb.net/?authSource=admin",
				standardKeySrv: "mongodb+srv://cluster1.mongodb.net/?authSource=admin",
			}

			privateEndpoints := []status.PrivateEndpoint{
				{
					ConnectionString:                  "mongodb://pe1.mongodb.net",
					SRVConnectionString:               "mongodb+srv://pe1.mongodb.net",
					SRVShardOptimizedConnectionString: "mongodb+srv://pe1-shard.mongodb.net",
				},
				{
					ConnectionString:                  "mongodb://pe2.mongodb.net",
					SRVConnectionString:               "mongodb+srv://pe2.mongodb.net",
					SRVShardOptimizedConnectionString: "mongodb+srv://pe2-shard.mongodb.net",
				},
			}

			for i, pe := range privateEndpoints {
				var suffix string
				if i != 0 {
					suffix = fmt.Sprint(i)
				}

				urlsToCheck[fmt.Sprintf("%s%s", privateKey, suffix)] = pe.ConnectionString
				urlsToCheck[fmt.Sprintf("%s%s", privateSrvKey, suffix)] = pe.SRVConnectionString
				urlsToCheck[fmt.Sprintf("%s%s", privateShardKey, suffix)] = pe.SRVShardOptimizedConnectionString
			}

			for key, baseURL := range urlsToCheck {
				want, _ := CreateURL(baseURL, username, "test-pass")
				require.Equal(t, want, string(s.Data[key]), "mismatch for %s", key)
			}
		})
	}
}
