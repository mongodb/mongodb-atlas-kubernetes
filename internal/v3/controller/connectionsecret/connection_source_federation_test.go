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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
)

func createDummyFederation(t *testing.T) *akov2.AtlasDataFederation {
	t.Helper()

	df := &akov2.AtlasDataFederation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-df",
			Namespace: "test-ns",
		},
		Spec: akov2.DataFederationSpec{
			Name: "my-df-name",
			Project: common.ResourceRefNamespaced{
				Name:      "test-project",
				Namespace: "test-ns",
			},
		},
		Status: status.DataFederationStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
			},
		},
	}

	return df
}

func runFederationProjectTest[T any](t *testing.T, method func(DataFederationConnectionSource) (T, error), wantField string) {
	r := createDummyEnv(t, nil)
	df := createDummyFederation(t)

	tests := map[string]struct {
		connectionSource DataFederationConnectionSource
		want             string
		wantErr          bool
	}{
		"fail: nil federation": {
			connectionSource: DataFederationConnectionSource{
				obj:             nil,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"fail: missing project ref": {
			connectionSource: DataFederationConnectionSource{
				obj: &akov2.AtlasDataFederation{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-df",
						Namespace: "test-ns",
					},
					Spec: akov2.DataFederationSpec{
						Name: "mising-proj",
					},
				},
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"success": {
			connectionSource: DataFederationConnectionSource{
				obj:             df,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			want: wantField,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := method(tc.connectionSource)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFederationConnectionSource_GetName(t *testing.T) {
	eNil := DataFederationConnectionSource{obj: nil}
	assert.Equal(t, "", eNil.GetName())
	e := DataFederationConnectionSource{obj: createDummyFederation(t)}
	assert.Equal(t, "my-df-name", e.GetName())
}

func TestFederationConnectionSource_IsReady(t *testing.T) {
	eNil := DataFederationConnectionSource{obj: nil}
	assert.False(t, eNil.IsReady())

	eNotReady := DataFederationConnectionSource{
		obj: &akov2.AtlasDataFederation{
			Status: status.DataFederationStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: "False"}},
				},
			},
		},
	}
	assert.False(t, eNotReady.IsReady())

	eReady := DataFederationConnectionSource{
		obj: &akov2.AtlasDataFederation{
			Status: status.DataFederationStatus{
				Common: api.Common{
					Conditions: []api.Condition{{Type: api.ReadyType, Status: "True"}},
				},
			},
		},
	}
	assert.True(t, eReady.IsReady())
}

func TestFederationConnectionSource_GetScopeType(t *testing.T) {
	e := DataFederationConnectionSource{}
	assert.Equal(t, akov2.DataLakeScopeType, e.GetScopeType())
}

func TestFederationConnectionSource_GetProjectID(t *testing.T) {
	runFederationProjectTest(t,
		func(fe DataFederationConnectionSource) (string, error) {
			return fe.GetProjectID(context.Background())
		},
		"test-project-id",
	)
}

func TestFederationConnectionSource_SelectorByProject(t *testing.T) {
	e := DataFederationConnectionSource{}
	s := e.SelectorByProjectID("p123")
	assert.True(t, s.Matches(fields.Set{indexer.AtlasDataFederationByProjectID: "p123"}))
	assert.False(t, s.Matches(fields.Set{indexer.AtlasDataFederationByProjectID: "other"}))
}

func TestFederationConnectionSource_SelectorByProjectIDAndClusterName(t *testing.T) {
	e := DataFederationConnectionSource{}
	ids := &ConnectionSecretIdentifiers{ProjectID: "pX", ClusterName: "dfY"}
	s := e.SelectorByProjectIDAndClusterName(ids)
	assert.True(t, s.Matches(fields.Set{indexer.AtlasDataFederationBySpecNameAndProjectID: "pX-dfY"}))
	assert.False(t, s.Matches(fields.Set{indexer.AtlasDataFederationBySpecNameAndProjectID: "pX-dfZ"}))
}

func TestFederationConnectionSource_BuildConnData(t *testing.T) {
	r := createDummyEnv(t, nil)
	df := createDummyFederation(t)
	user := createDummyUser(t, "test-user")

	userNoPass := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "test-user-nopass", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			PasswordSecret: &common.ResourceRef{
				Name: "missing-secret",
			},
			Username: "theuser",
		},
	}

	dfNoProject := &akov2.AtlasDataFederation{
		ObjectMeta: metav1.ObjectMeta{Name: "df", Namespace: "test-ns"},
		Spec:       akov2.DataFederationSpec{Name: "df"},
	}

	tests := map[string]struct {
		objs             []client.Object
		override         func(*ConnSecretReconciler)
		connectionSource *akov2.AtlasDataFederation
		user             *akov2.AtlasDatabaseUser
		wantURL          string
		wantErr          bool
	}{
		"fail: nil connectionSource and nil user": {
			connectionSource: nil,
			user:             nil,
			wantErr:          true,
		},
		"fail: password is missing": {
			connectionSource: dfNoProject,
			user:             userNoPass,
			wantErr:          true,
		},
		"fail: connectionSource exists but project missing": {
			connectionSource: dfNoProject,
			user:             user,
			wantErr:          true,
		},
		"success: builds URL from DF hostnames": {
			override: func(r *ConnSecretReconciler) {
				dfAPI := mockadmin.NewDataFederationApi(t)

				dfAPI.EXPECT().
					GetFederatedDatabase(mock.Anything, "test-project-id", "my-df-name").
					Return(admin.GetFederatedDatabaseApiRequest{ApiService: dfAPI})

				dfAPI.EXPECT().
					GetFederatedDatabaseExecute(mock.AnythingOfType("admin.GetFederatedDatabaseApiRequest")).
					Return(&admin.DataLakeTenant{
						Hostnames: &[]string{"h1.example.net", "h2.example.net"},
					}, nil, nil)

				r.AtlasProvider = &atlasmock.TestProvider{
					SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
						return &atlas.ClientSet{
							SdkClient20250312006: &admin.APIClient{
								DataFederationApi: dfAPI,
							},
						}, nil
					},
					IsSupportedFunc: func() bool { return true },
					IsCloudGovFunc:  func() bool { return false },
				}
			},
			connectionSource: df,
			user:             user,
			wantURL:          "mongodb://h1.example.net,h2.example.net/?ssl=true",
			wantErr:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.override != nil {
				tc.override(r)
			}
			e := DataFederationConnectionSource{
				obj:             tc.connectionSource,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			}
			got, err := e.BuildConnectionData(context.Background(), tc.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "admin", got.DBUserName)
			assert.Equal(t, "secret", got.Password)
			assert.Equal(t, tc.wantURL, got.ConnectionURL)
		})
	}
}
