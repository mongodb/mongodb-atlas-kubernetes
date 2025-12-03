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

package atlasdatabaseuser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

func TestHandleDatabaseUser(t *testing.T) {
	tests := map[string]struct {
		dbUserInAKO        *akov2.AtlasDatabaseUser
		dbUserSecret       *corev1.Secret
		atlasProject       *akov2.AtlasProject
		atlasProvider      atlas.Provider
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"user spec is invalid": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "invalid-version",
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			atlasProvider: &atlasmock.TestProvider{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ResourceVersionStatus).
					WithReason(string(workflow.AtlasResourceVersionIsInvalid)).
					WithMessageRegexp("invalid-version is not a valid semver version for label mongodb.com/atlas-resource-version"),
			},
		},
		"user spec is mismatch": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "1000.0.1",
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			atlasProvider: &atlasmock.TestProvider{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ResourceVersionStatus).
					WithReason(string(workflow.AtlasResourceVersionMismatch)).
					WithMessageRegexp("version of the resource 'user1' is higher than the operator version '2.4.1'"),
			},
		},
		"resource is not supported": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "2.4.1",
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.ValidationSucceeded),
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.AtlasGovUnsupported)).
					WithMessageRegexp("the *v1.AtlasDatabaseUser is not supported by Atlas for government"),
			},
		},
		"manage user with independent configuration": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "project-id",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				IsCloudGovFunc: func() bool {
					return false
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProject(context.Background(), "project-id").
						Return(admin.GetProjectApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(&admin.Group{Id: pointer.MakePtr("project-id")}, nil, nil)

					userAPI := mockadmin.NewDatabaseUsersApi(t)
					userAPI.EXPECT().GetDatabaseUser(context.Background(), "project-id", "admin", "user1").
						Return(admin.GetDatabaseUserApiRequest{ApiService: userAPI})
					userAPI.EXPECT().GetDatabaseUserExecute(mock.AnythingOfType("admin.GetDatabaseUserApiRequest")).
						Return(nil, nil, nil)
					userAPI.EXPECT().CreateDatabaseUser(context.Background(), "project-id", mock.AnythingOfType("*admin.CloudDatabaseUser")).
						Return(admin.CreateDatabaseUserApiRequest{ApiService: userAPI})
					userAPI.EXPECT().CreateDatabaseUserExecute(mock.AnythingOfType("admin.CreateDatabaseUserApiRequest")).
						Return(&admin.CloudDatabaseUser{}, nil, nil)

					clusterAPI := mockadmin.NewClustersApi(t)

					return &atlas.ClientSet{
						SdkClient20250312006: &admin.APIClient{ProjectsApi: projectAPI, ClustersApi: clusterAPI, DatabaseUsersApi: userAPI},
					}, nil
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.ValidationSucceeded),
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
		"manage user with linked configuration": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "2.4.1",
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			atlasProvider:  DefaultTestProvider(t),
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.ValidationSucceeded),
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUserInAKO).
				WithStatusSubresource(tt.dbUserInAKO)

			if tt.atlasProject != nil {
				k8sClient.WithObjects(tt.atlasProject)
			}

			if tt.dbUserSecret != nil {
				k8sClient.WithObjects(tt.dbUserSecret)
			}
			k8sClient.WithObjects(&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"orgId":         []byte("orgId"),
					"publicApiKey":  []byte("publicApiKey"),
					"privateApiKey": []byte("privateApiKey"),
				},
			})

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient.Build(),
					Log:    logger,
					GlobalSecretRef: client.ObjectKey{
						Namespace: "default",
						Name:      "secret",
					},
					AtlasProvider: tt.atlasProvider,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			version.Version = "2.4.1"

			result, err := r.handleDatabaseUser(ctx, tt.dbUserInAKO)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			logger.Infof("conditions", ctx.Conditions())

			diff := cmp.Diff(
				tt.expectedConditions,
				ctx.Conditions(),
				cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
			)
			if diff != "" {
				t.Errorf("status conditions differ: %v", diff)
			}
		})
	}
}

func TestDbuLifeCycle(t *testing.T) {
	deletionTime := metav1.Now()

	tests := map[string]struct {
		dbUserInAKO        *akov2.AtlasDatabaseUser
		dbUserSecret       *corev1.Secret
		dbUserService      func() dbuser.AtlasUsersService
		dService           func() deployment.AtlasDeploymentsService
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"failed to get user in atlas": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").
					Return(nil, errors.New("failed to get user"))

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get user"),
			},
		},
		"failed to check user is expired": { //nolint:dupl
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName:    "admin",
					DeleteAfterDate: "wrong-date",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").Return(nil, nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserInvalidSpec)).
					WithMessageRegexp("parsing time \"wrong-date\" as \"2006-01-02T15:04:05.999Z\": cannot parse \"wrong-date\" as \"2006\""),
			},
		},
		"user is expired": { //nolint:dupl
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName:    "admin",
					DeleteAfterDate: "2021-05-30T13:45:30Z",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").Return(nil, nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserExpired)).
					WithMessageRegexp("an expired user cannot be managed"),
			},
		},
		"failed to validate scope": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster1",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").Return(nil, nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ClusterExists(context.Background(), "", "cluster1").
					Return(false, errors.New("failed to check cluster exists"))

				return service
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserInvalidSpec)).
					WithMessageRegexp("failed to check cluster exists"),
			},
		},
		"deployment scope is invalid": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster1",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").Return(nil, nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ClusterExists(context.Background(), "", "cluster1").
					Return(false, nil)

				return service
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserInvalidSpec)).
					WithMessageRegexp("\"scopes\" field refer to one or more deployments that don't exist"),
			},
		},
		"create an user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").Return(nil, nil)
				service.EXPECT().Create(context.Background(), mock.AnythingOfType("*dbuser.User")).Return(nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
		"update an user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "999",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").
					Return(
						&dbuser.User{
							AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
								Username: "user1",
								PasswordSecret: &common.ResourceRef{
									Name: "user-pass",
								},
								DatabaseName: "admin",
								Scopes:       []akov2.ScopeSpec{},
							},
						},
						nil,
					)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").Return([]string{}, nil)
				service.EXPECT().ListDeploymentConnections(context.Background(), "").Return([]deployment.Connection{}, nil)

				return service
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.DatabaseUserReadyType),
			},
		},
		"delete an user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "user1",
					Namespace:         "default",
					Finalizers:        []string{"mongodbatlas/finalizer"},
					DeletionTimestamp: &deletionTime,
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "999",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").
					Return(
						&dbuser.User{
							AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
								Username: "user1",
								PasswordSecret: &common.ResourceRef{
									Name: "user-pass",
								},
								DatabaseName: "admin",
							},
						},
						nil,
					)
				service.EXPECT().Delete(context.Background(), "admin", "", "user1").Return(nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{},
		},
		"unmanage an user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "user1",
					Namespace:         "default",
					Finalizers:        []string{"mongodbatlas/finalizer"},
					DeletionTimestamp: &deletionTime,
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "999",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Get(context.Background(), "admin", "", "user1").
					Return(nil, nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUserInAKO).
				WithStatusSubresource(tt.dbUserInAKO)

			if tt.dbUserSecret != nil {
				k8sClient.WithObjects(tt.dbUserSecret)
			}

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient.Build(),
					Log:    logger,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := r.dbuLifeCycle(ctx, tt.dbUserService(), tt.dService(), tt.dbUserInAKO, &project.Project{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		dbUserInAKO        *akov2.AtlasDatabaseUser
		dbUserSecret       *corev1.Secret
		dbUserService      func() dbuser.AtlasUsersService
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"failed to read user password": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("secrets \"user-pass\" not found"),
			},
		},
		"failed to convert to internal user": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DeleteAfterDate: "wrong-date",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to create internal user type: failed to parse \"wrong-date\" to an ISO date: parsing time \"wrong-date\" as \"2006-01-02T15:04:05.999Z\": cannot parse \"wrong-date\" as \"2006\""),
			},
		},
		"failed to create user": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Create(context.Background(), mock.AnythingOfType("*dbuser.User")).
					Return(errors.New("failed to create user"))

				return service
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserNotCreatedInAtlas)).
					WithMessageRegexp("failed to create user"),
			},
		},
		"renamed user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user-renamed",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					UserName:        "user1",
					PasswordVersion: "999",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Create(context.Background(), mock.AnythingOfType("*dbuser.User")).Return(nil)
				service.EXPECT().Delete(context.Background(), "admin", "project-id", "user1").Return(nil)

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
		"create user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Create(context.Background(), mock.AnythingOfType("*dbuser.User")).Return(nil)

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUserInAKO).
				WithStatusSubresource(tt.dbUserInAKO)

			if tt.dbUserSecret != nil {
				k8sClient.WithObjects(tt.dbUserSecret)
			}

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient.Build(),
					Log:    logger,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := r.create(ctx, tt.dbUserService(), "project-id", tt.dbUserInAKO)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		dbUserInAKO        *akov2.AtlasDatabaseUser
		dbUserSecret       *corev1.Secret
		dbUserInAtlas      *dbuser.User
		dbUserService      func() dbuser.AtlasUsersService
		dService           func() deployment.AtlasDeploymentsService
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"failed to read user password": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("secrets \"user-pass\" not found"),
			},
		},
		"failed to convert to internal user": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DeleteAfterDate: "wrong-date",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to create internal user type: failed to parse \"wrong-date\" to an ISO date: parsing time \"wrong-date\" as \"2006-01-02T15:04:05.999Z\": cannot parse \"wrong-date\" as \"2006\""),
			},
		},
		"user hasn't change": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "999",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
					Scopes:       []akov2.ScopeSpec{},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").Return([]string{}, nil)
				service.EXPECT().ListDeploymentConnections(context.Background(), "").Return([]deployment.Connection{}, nil)

				return service
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.DatabaseUserReadyType),
			},
		},
		"failed to update user": {
			wantErr: true,
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Update(context.Background(), mock.AnythingOfType("*dbuser.User")).
					Return(errors.New("failed to update user"))

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserNotUpdatedInAtlas)).
					WithMessageRegexp("failed to update user"),
			},
		},
		"update user": {
			dbUserInAKO: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					PasswordVersion: "888",
				},
			},
			dbUserSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Update(context.Background(), mock.AnythingOfType("*dbuser.User")).Return(nil)

				return service
			},
			dService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("Clusters are scheduled to handle database users updates"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUserInAKO).
				WithStatusSubresource(tt.dbUserInAKO)

			if tt.dbUserSecret != nil {
				k8sClient.WithObjects(tt.dbUserSecret)
			}

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient.Build(),
					Log:    logger,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := r.update(ctx, tt.dbUserService(), tt.dService(), &project.Project{}, tt.dbUserInAKO, tt.dbUserInAtlas)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		dbUser             *akov2.AtlasDatabaseUser
		dbUserService      func() dbuser.AtlasUsersService
		deletionProtection bool
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"don't delete resource on atlas when deletion protection is enabled": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			deletionProtection: true,
			expectedResult:     ctrl.Result{},
			expectedConditions: nil,
		},
		"don't delete resource on atlas when is annotated to keep resource": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				return translation.NewAtlasUsersServiceMock(t)
			},
			deletionProtection: false,
			expectedResult:     ctrl.Result{},
			expectedConditions: nil,
		},
		"failed to delete resource": {
			wantErr: true,
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Delete(context.Background(), "admin", "project-id", "user1").
					Return(errors.New("failed to delete user"))

				return service
			},
			deletionProtection: false,
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserNotDeletedInAtlas)).
					WithMessageRegexp("failed to delete user"),
			},
		},
		"user was already deleted": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Delete(context.Background(), "admin", "project-id", "user1").
					Return(dbuser.ErrorNotFound)

				return service
			},
			deletionProtection: false,
			expectedResult:     ctrl.Result{},
			expectedConditions: nil,
		},
		"delete user": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dbUserService: func() dbuser.AtlasUsersService {
				service := translation.NewAtlasUsersServiceMock(t)
				service.EXPECT().Delete(context.Background(), "admin", "project-id", "user1").
					Return(nil)

				return service
			},
			deletionProtection: false,
			expectedResult:     ctrl.Result{},
			expectedConditions: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUser).
				WithStatusSubresource(tt.dbUser).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
				ObjectDeletionProtection: tt.deletionProtection,
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := r.delete(ctx, tt.dbUserService(), "project-id", tt.dbUser)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestReadiness(t *testing.T) {
	tests := map[string]struct {
		dbUser             *akov2.AtlasDatabaseUser
		dService           func() deployment.AtlasDeploymentsService
		expectedResult     ctrl.Result
		wantErr            bool
		expectedConditions []api.Condition
	}{
		"failed to list cluster names": {
			wantErr: true,
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").
					Return(nil, errors.New("failed to list cluster names"))

				return service
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to list cluster names"),
			},
		},
		"failed to check deployment status": {
			wantErr: true,
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").
					Return([]string{"cluster1", "cluster2"}, nil)
				service.EXPECT().DeploymentIsReady(context.Background(), "", "cluster2").
					Return(false, errors.New("failed to check status"))

				return service
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to check status"),
			},
		},
		"deployments are in progress": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").
					Return([]string{"cluster1", "cluster2"}, nil)
				service.EXPECT().DeploymentIsReady(context.Background(), "", "cluster2").
					Return(false, nil)

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)).
					WithMessageRegexp("0 out of 1 deployments have applied database user changes"),
			},
		},
		"failed to create connection secrets": {
			wantErr: true,
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").
					Return([]string{"cluster1", "cluster2"}, nil)
				service.EXPECT().DeploymentIsReady(context.Background(), "", "cluster2").
					Return(true, nil)
				service.EXPECT().ListDeploymentConnections(context.Background(), "").
					Return(nil, errors.New("failed to list cluster connections"))

				return service
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.DatabaseUserConnectionSecretsNotCreated)).
					WithMessageRegexp("failed to list cluster connections"),
			},
		},
		"resource is ready": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: "cluster2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			dService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ListDeploymentNames(context.Background(), "").
					Return([]string{"cluster1", "cluster2"}, nil)
				service.EXPECT().DeploymentIsReady(context.Background(), "", "cluster2").
					Return(true, nil)
				service.EXPECT().ListDeploymentConnections(context.Background(), "").
					Return([]deployment.Connection{}, nil)

				return service
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.DatabaseUserReadyType),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUser).
				WithStatusSubresource(tt.dbUser).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := r.readiness(ctx, tt.dService(), &project.Project{}, tt.dbUser, "999")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestReadPassword(t *testing.T) {
	tests := map[string]struct {
		dbUser           *akov2.AtlasDatabaseUser
		secret           *corev1.Secret
		expectedPassword string
		expectedVersion  string
		expectedErr      error
	}{
		"return empty if password secret is unset": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{},
			},
		},
		"read password from secret": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			expectedPassword: "Passw0rd!",
			expectedVersion:  "999",
		},
		"err when secret is not found": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			expectedErr: &k8serrors.StatusError{ErrStatus: metav1.Status{Status: "Failure", Message: "secrets \"user-pass\" not found", Code: 404, Details: &metav1.StatusDetails{Name: "user-pass", Kind: "secrets"}, Reason: "NotFound"}},
		},
		"error when password is not present in the secret": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"no-password": []byte("Passw0rd!"),
				},
			},
			expectedErr: errors.New("secret user-pass is invalid: it doesn't contain 'password' field"),
		},
		"error when password is empty in the secret": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": {},
				},
			},
			expectedErr: errors.New("secret user-pass is invalid: the 'password' field is empty"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme)

			if tt.secret != nil {
				k8sClient.WithObjects(tt.secret)
			}

			r := AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient.Build(),
				},
			}

			pass, passVersion, err := r.readPassword(context.Background(), tt.dbUser)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPassword, pass)
			assert.Equal(t, tt.expectedVersion, passVersion)
		})
	}
}

func TestAreDeploymentScopesValid(t *testing.T) {
	tests := map[string]struct {
		dbUser   *akov2.AtlasDatabaseUser
		call     func(ctx context.Context, s string, s2 string) (bool, error)
		expected bool
		err      error
	}{
		"scope is empty": {
			dbUser:   &akov2.AtlasDatabaseUser{},
			expected: true,
		},
		"scope doesn't contains deployments": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{
							Name: "dl1",
							Type: akov2.DataLakeScopeType,
						},
						{
							Name: "dl2",
							Type: akov2.DataLakeScopeType,
						},
					},
				},
			},
			expected: true,
		},
		"fail to evaluate scope": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{
							Name: "d1",
							Type: akov2.DeploymentScopeType,
						},
						{
							Name: "d2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			call: func(ctx context.Context, s string, s2 string) (bool, error) {
				return false, errors.New("failed to request")
			},
			expected: false,
			err:      errors.New("failed to request"),
		},
		"return false when deployment doesn't exist": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{
							Name: "d1",
							Type: akov2.DeploymentScopeType,
						},
						{
							Name: "d2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			call: func(ctx context.Context, s string, s2 string) (bool, error) {
				return false, nil
			},
			expected: false,
		},
		"return true when all deployment exist": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{
							Name: "d1",
							Type: akov2.DeploymentScopeType,
						},
						{
							Name: "d2",
							Type: akov2.DeploymentScopeType,
						},
					},
				},
			},
			call: func(ctx context.Context, s string, s2 string) (bool, error) {
				return true, nil
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			deploymentService := translation.NewAtlasDeploymentsServiceMock(t)
			if tt.call != nil {
				deploymentService.EXPECT().ClusterExists(context.Background(), "project-id", mock.AnythingOfType("string")).
					RunAndReturn(tt.call)
			}
			r := AtlasDatabaseUserReconciler{}
			ctx := &workflow.Context{
				Context: context.Background(),
			}
			valid, err := r.areDeploymentScopesValid(ctx, deploymentService, "project-id", tt.dbUser)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, valid)
		})
	}

	scopeSpecs := []akov2.ScopeSpec{{
		Name: "dbLake",
		Type: akov2.DataLakeScopeType,
	}, {
		Name: "cluster1",
		Type: akov2.DeploymentScopeType,
	}, {
		Name: "cluster2",
		Type: akov2.DeploymentScopeType,
	}}
	clusters := []string{"cluster1", "cluster4", "cluster5"}
	scopeClusters := filterScopeDeployments(&akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: scopeSpecs}}, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}

func TestRemoveOldUser(t *testing.T) {
	failedFirst := false

	tests := map[string]struct {
		call func(ctx context.Context, db string, project string, username string) error
		err  error
	}{
		"delete old user": {
			call: func(ctx context.Context, db string, project string, username string) error {
				return nil
			},
		},
		"user was already deleted": {
			call: func(ctx context.Context, db string, project string, username string) error {
				return dbuser.ErrorNotFound
			},
		},
		"failed on first attempt and then succeed": {
			call: func(ctx context.Context, db string, project string, username string) error {
				if failedFirst {
					return nil
				}

				failedFirst = true
				return errors.New("failed first")
			},
		},
		"failed to delete old user": {
			call: func(ctx context.Context, db string, project string, username string) error {
				return errors.New("fail always")
			},
			err: errors.New("fail always"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dbUserService := translation.NewAtlasUsersServiceMock(t)
			dbUserService.EXPECT().Delete(context.Background(), "admin", "project-id", "old-name").
				RunAndReturn(tt.call)
			r := &AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Log: zaptest.NewLogger(t).Sugar(),
				},
			}

			user := &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
				},
				Status: status.AtlasDatabaseUserStatus{
					UserName: "old-name",
				},
			}
			assert.Equal(t, tt.err, r.removeOldUser(context.Background(), dbUserService, "project-id", user))
		})
	}
}

func TestIsExpired(t *testing.T) {
	after := time.Now().UTC().Add(time.Hour * 24).Format("2006-01-02T15:04:05")

	tests := map[string]struct {
		dbUser   *akov2.AtlasDatabaseUser
		expected bool
		err      error
	}{
		"user has no expiration date": {
			dbUser:   akov2.DefaultDBUser("ns", "theuser", ""),
			expected: false,
			err:      nil,
		},
		"user has invalid expiration date": {
			dbUser:   akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("foo"),
			expected: false,
			err:      &time.ParseError{Layout: "2006-01-02T15:04:05.999Z", Value: "foo", LayoutElem: "2006", ValueElem: "foo"},
		},
		"user has non expired date": {
			dbUser:   akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate(after),
			expected: false,
			err:      nil,
		},
		"user has an expired date": {
			dbUser:   akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("2021-11-30T15:04:05"),
			expected: true,
			err:      nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			expired, err := isExpired(tt.dbUser)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, expired)
		})
	}
}

func TestHasChanged(t *testing.T) {
	tests := map[string]struct {
		dbUserInAKO    *dbuser.User
		dbUserInAtlas  *dbuser.User
		currentVersion string
		version        string
		expected       bool
	}{
		"user and password haven't changed": {
			dbUserInAKO: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "NONE",
				},
				ProjectID: "project-id",
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "NONE",
				},
				ProjectID: "project-id",
			},
			currentVersion: "123456",
			version:        "123456",
			expected:       false,
		},
		"user has changed and password doesn't": {
			dbUserInAKO: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "MANAGED",
				},
				ProjectID: "project-id",
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "NONE",
				},
				ProjectID: "project-id",
			},
			currentVersion: "123456",
			version:        "123456",
			expected:       true,
		},
		"user hasn't changed and password does": {
			dbUserInAKO: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "NONE",
				},
				ProjectID: "project-id",
			},
			dbUserInAtlas: &dbuser.User{
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					Username:     "user1",
					DatabaseName: "admin",
					OIDCAuthType: "NONE",
					AWSIAMType:   "NONE",
					X509Type:     "NONE",
				},
				ProjectID: "project-id",
			},
			currentVersion: "123456",
			version:        "654321",
			expected:       true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hasChanged(tt.dbUserInAKO, tt.dbUserInAtlas, tt.currentVersion, tt.version))
		})
	}
}

func TestWasRenamed(t *testing.T) {
	tests := map[string]struct {
		dbUser   *akov2.AtlasDatabaseUser
		expected bool
	}{
		"the user is new": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
				},
			},
			expected: false,
		},
		"the user hasn't been renamed": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
				},
				Status: status.AtlasDatabaseUserStatus{
					UserName: "user1",
				},
			},
			expected: false,
		},
		"the user was renamed": {
			dbUser: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "user1",
				},
				Status: status.AtlasDatabaseUserStatus{
					UserName: "user2",
				},
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, wasRenamed(tt.dbUser))
		})
	}
}

func TestFilterScopeDeployments(t *testing.T) {
	scopeSpecs := []akov2.ScopeSpec{{
		Name: "dbLake",
		Type: akov2.DataLakeScopeType,
	}, {
		Name: "cluster1",
		Type: akov2.DeploymentScopeType,
	}, {
		Name: "cluster2",
		Type: akov2.DeploymentScopeType,
	}}
	clusters := []string{"cluster1", "cluster4", "cluster5"}
	scopeClusters := filterScopeDeployments(&akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: scopeSpecs}}, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}

func DefaultTestProvider(t *testing.T) *atlasmock.TestProvider {
	t.Helper()

	return &atlasmock.TestProvider{
		IsSupportedFunc: func() bool {
			return true
		},
		IsCloudGovFunc: func() bool {
			return false
		},
		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
			userAPI := mockadmin.NewDatabaseUsersApi(t)
			userAPI.EXPECT().GetDatabaseUser(context.Background(), "my-project", "admin", "user1").
				Return(admin.GetDatabaseUserApiRequest{ApiService: userAPI})
			userAPI.EXPECT().GetDatabaseUserExecute(mock.AnythingOfType("admin.GetDatabaseUserApiRequest")).
				Return(nil, nil, nil)
			userAPI.EXPECT().CreateDatabaseUser(context.Background(), "my-project", mock.AnythingOfType("*admin.CloudDatabaseUser")).
				Return(admin.CreateDatabaseUserApiRequest{ApiService: userAPI})
			userAPI.EXPECT().CreateDatabaseUserExecute(mock.AnythingOfType("admin.CreateDatabaseUserApiRequest")).
				Return(&admin.CloudDatabaseUser{}, nil, nil)

			clusterAPI := mockadmin.NewClustersApi(t)

			projectAPI := mockadmin.NewProjectsApi(t)
			projectAPI.EXPECT().GetProjectByName(mock.Anything, "my-project").
				Return(admin.GetProjectByNameApiRequest{ApiService: projectAPI})
			projectAPI.EXPECT().GetProjectByNameExecute(mock.Anything).
				Return(&admin.Group{Id: pointer.MakePtr("my-project")}, nil, nil)

			return &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ProjectsApi:      projectAPI,
					ClustersApi:      clusterAPI,
					DatabaseUsersApi: userAPI,
				},
			}, nil
		},
	}
}
