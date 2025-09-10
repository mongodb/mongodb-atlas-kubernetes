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

package experimentalconnectionsecret

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

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

// createDummyEnv creates a dummy environment with some objects already setup
func createDummyEnv(t *testing.T, objs []client.Object) *ConnSecretReconciler {
	scheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(scheme))
	assert.NoError(t, corev1.AddToScheme(scheme))

	logger := zaptest.NewLogger(t)

	// Contains the project
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-project",
			Namespace: "test-ns",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "My Project Name",
			ConnectionSecret: &common.ResourceRefNamespaced{
				Name: "sdk-creds",
			},
		},
		Status: status.AtlasProjectStatus{
			ID: "test-project-id",
		},
	}

	// SDK credentials
	sdkSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdk-creds",
			Namespace: "test-ns",
		},
		Data: map[string][]byte{
			"orgId":         []byte("test-pass"),
			"publicApiKey":  []byte("test-pass"),
			"privateApiKey": []byte("test-pass"),
		},
	}

	// Connection Secret
	connSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-project-name-cluster1-admin",
			Namespace: "test-ns",
			Labels: map[string]string{
				ProjectLabelKey: "test-project-id",
				ClusterLabelKey: "cluster1",
				TypeLabelKey:    "connection",
			},
		},
	}

	// User password
	userSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "user-pass",
			Namespace: "test-ns",
		},
		Data: map[string][]byte{"password": []byte("secret")},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(project, sdkSecret, connSecret, userSecret).
		WithObjects(objs...).
		Build()

	indexer1 := indexer.NewAtlasDatabaseUserByProjectIndexer(context.Background(), cl, logger)
	indexer2 := indexer.NewAtlasDataFederationByProjectIDIndexer(context.Background(), cl, logger)
	indexer3 := indexer.NewAtlasDeploymentByProjectIndexer(context.Background(), cl, logger)

	cl = fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(project, sdkSecret, connSecret, userSecret).
		WithObjects(objs...).
		WithIndex(&akov2.AtlasDeployment{}, indexer.AtlasDeploymentBySpecNameAndProjectID, func(obj client.Object) []string {
			d := obj.(*akov2.AtlasDeployment)
			return []string{"test-project-id" + "-" + d.Spec.DeploymentSpec.Name}
		}).
		WithIndex(&akov2.AtlasDataFederation{}, indexer.AtlasDataFederationBySpecNameAndProjectID, func(obj client.Object) []string {
			df := obj.(*akov2.AtlasDataFederation)
			return []string{"test-project-id" + "-" + df.Spec.Name}
		}).
		WithIndex(&akov2.AtlasDatabaseUser{}, indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, func(obj client.Object) []string {
			u := obj.(*akov2.AtlasDatabaseUser)
			return []string{"test-project-id" + "-" + u.Spec.Username}
		}).
		WithIndex(indexer1.Object(), indexer1.Name(), indexer1.Keys).
		WithIndex(indexer2.Object(), indexer2.Name(), indexer2.Keys).
		WithIndex(indexer3.Object(), indexer3.Name(), indexer3.Keys).
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
					Name: "My Project Name",
				}, nil, nil)

			return &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ProjectsApi: projectAPI,
				},
			}, nil
		},
		IsSupportedFunc: func() bool { return true },
		IsCloudGovFunc:  func() bool { return false },
	}

	r := &ConnSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:        cl,
			AtlasProvider: atlasProvider,
			Log:           logger.Sugar(),
		},
		Scheme:        scheme,
		EventRecorder: record.NewFakeRecorder(10),
	}

	return r
}

func createDummyUser(t *testing.T) *akov2.AtlasDatabaseUser {
	t.Helper()

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-user",
			Namespace: "test-ns",
		},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username:       "admin",
			PasswordSecret: &common.ResourceRef{Name: "user-pass"},
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
	}

	return user
}

func createDummyUserSDK(t *testing.T) *akov2.AtlasDatabaseUser {
	t.Helper()

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-user",
			Namespace: "test-ns",
		},
		Spec: akov2.AtlasDatabaseUserSpec{
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
	}

	return user
}

func Test_getUserProjectName(t *testing.T) {
	r := createDummyEnv(t, []client.Object{})
	user := createDummyUser(t)
	usersdk := createDummyUserSDK(t)

	type testCase struct {
		user     *akov2.AtlasDatabaseUser
		wantName string
		wantErr  bool
	}

	tests := map[string]testCase{
		"fail: nil user returns error": {
			user:    nil,
			wantErr: true,
		},
		"fail: k8s project ref not found returns error": {
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "missing-proj",
							Namespace: "test-ns",
						},
					},
				},
			},
			wantErr: true,
		},
		"fail: no project ref and nil receiver falls back to not available error": {
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{},
			},
			wantErr: true,
		},
		"success: k8s project ref success returns normalized name by reference": {
			user:     user,
			wantName: "My Project Name",
			wantErr:  false,
		},
		"success: k8s project ref success returns normalized name by sdk": {
			user:     usersdk,
			wantName: "My Project Name",
			wantErr:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			receiver := r

			got, err := receiver.getUserProjectName(context.Background(), tc.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantName, got)
		})
	}
}
