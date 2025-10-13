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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
)

func createDummyDeployment(t *testing.T, deploymentName string, deploymentProjectName string, deploymentClusterName string) *akov2.AtlasDeployment {
	t.Helper()

	depl := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: "test-ns",
		},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: deploymentClusterName},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name:      deploymentProjectName,
					Namespace: "test-ns",
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
	}

	return depl
}

func createDummyDeploymentSDK(t *testing.T) *akov2.AtlasDeployment {
	t.Helper()

	depl := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-depl",
			Namespace: "test-ns",
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
			ConnectionStrings: &status.ConnectionStrings{
				Standard:    "mongodb+srv://cluster1.mongodb.net",
				StandardSrv: "mongodb://cluster1.mongodb.net",
			},
		},
	}

	return depl
}

func TestDeploymentConnectionTarget_GetName(t *testing.T) {
	eNil := DeploymentConnectionTarget{obj: nil}
	assert.Equal(t, "", eNil.GetName())
	dep := createDummyDeployment(t, "test-depl", "test-project-second", "cluster1")
	e := DeploymentConnectionTarget{obj: dep}
	assert.Equal(t, "cluster1", e.GetName())
}

func TestDeploymentConnectionTarget_IsReady(t *testing.T) {
	eNil := DeploymentConnectionTarget{obj: nil}
	assert.False(t, eNil.IsReady())

	notReady := &akov2.AtlasDeployment{
		Status: status.AtlasDeploymentStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: "False"}},
			},
		},
	}
	assert.False(t, DeploymentConnectionTarget{obj: notReady}.IsReady())

	ready := &akov2.AtlasDeployment{
		Status: status.AtlasDeploymentStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: "True"}},
			},
		},
	}
	assert.True(t, DeploymentConnectionTarget{obj: ready}.IsReady())
}

func TestDeploymentConnectionTarget_GetScopeType(t *testing.T) {
	e := DeploymentConnectionTarget{}
	assert.Equal(t, akov2.DeploymentScopeType, e.GetScopeType())
}

func TestDeploymentConnectionTarget_GetProjectID(t *testing.T) {
	r := createDummyEnv(t, nil)
	depl := createDummyDeployment(t, "test-depl", "test-project", "cluster1")
	deplsdk := createDummyDeploymentSDK(t)

	tests := map[string]struct {
		connectionTarget DeploymentConnectionTarget
		want             string
		wantErr          bool
	}{
		"fail: nil deployment": {
			connectionTarget: DeploymentConnectionTarget{
				obj:             nil,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"fail: project ref missing": {
			connectionTarget: DeploymentConnectionTarget{
				obj:             &akov2.AtlasDeployment{Spec: akov2.AtlasDeploymentSpec{}},
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"fail: k8s project ref but project not found": {
			connectionTarget: DeploymentConnectionTarget{
				obj: &akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Namespace: "test-ns"},
					Spec: akov2.AtlasDeploymentSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "missing-proj",
								Namespace: "test-ns",
							},
						},
					},
				},
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"success: external project ID": {
			connectionTarget: DeploymentConnectionTarget{
				obj:             deplsdk,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			want: "test-project-id",
		},
		"success: k8s project ref": {
			connectionTarget: DeploymentConnectionTarget{
				obj:             depl,
				client:          r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			want: "test-project-id",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.connectionTarget.GetProjectID(context.Background())
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeploymentConnectionTarget_SelectorByProject(t *testing.T) {
	e := DeploymentConnectionTarget{}
	s := e.SelectorByProjectID("p-1")
	assert.True(t, s.Matches(fields.Set{indexer.AtlasDeploymentByProject: "p-1"}))
	assert.False(t, s.Matches(fields.Set{indexer.AtlasDeploymentByProject: "other"}))
}

func TestDeploymentConnectionTarget_BuildConnData(t *testing.T) {
	r := createDummyEnv(t, nil)
	depl := createDummyDeployment(t, "test-depl", "test-project", "cluster1")
	user := createDummyUser(t, "test-user", "admin", "dummy-uid")

	userNoPass := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "test-user-nopass", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "theuser",
			PasswordSecret: &common.ResourceRef{
				Name: "missing-proj",
			},
		},
	}

	depl.Status.ConnectionStrings = &status.ConnectionStrings{
		Standard:    "mongodb://std:27017",
		StandardSrv: "mongodb+srv://std",
		Private:     "mongodb://priv:27017",
		PrivateSrv:  "mongodb+srv://priv",
		PrivateEndpoint: []status.PrivateEndpoint{
			{
				ConnectionString:                  "mongodb://pe1:27017",
				SRVConnectionString:               "mongodb+srv://pe1",
				SRVShardOptimizedConnectionString: "mongodb+srv://pe1-shard",
			},
			{
				ConnectionString:                  "mongodb://pe2:27017",
				SRVConnectionString:               "mongodb+srv://pe2",
				SRVShardOptimizedConnectionString: "mongodb+srv://pe2-shard",
			},
		},
	}

	tests := map[string]struct {
		connectionTarget *akov2.AtlasDeployment
		user             *akov2.AtlasDatabaseUser
		want             ConnectionSecretData
		wantErr          bool
	}{
		"fail: nil connectionTarget and user": {
			connectionTarget: nil,
			user:             nil,
			wantErr:          true,
		},
		"fail: missing password": {
			connectionTarget: depl,
			user:             userNoPass,
			wantErr:          true,
		},
		"success: builds from deployment connection strings": {
			connectionTarget: depl,
			user:             user,
			want: ConnectionSecretData{
				DBUserName:       "admin",
				Password:         "secret",
				ConnectionURL:    "mongodb://std:27017",
				SrvConnectionURL: "mongodb+srv://std",
				PrivateConnectionURLs: []PrivateLinkConnectionURLs{
					{ConnectionURL: "mongodb://priv:27017", SrvConnectionURL: "mongodb+srv://priv"},
					{ConnectionURL: "mongodb://pe1:27017", SrvConnectionURL: "mongodb+srv://pe1", ShardConnectionURL: "mongodb+srv://pe1-shard"},
					{ConnectionURL: "mongodb://pe2:27017", SrvConnectionURL: "mongodb+srv://pe2", ShardConnectionURL: "mongodb+srv://pe2-shard"},
				},
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e := DeploymentConnectionTarget{
				obj:             tc.connectionTarget,
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
			assert.Equal(t, tc.want.DBUserName, got.DBUserName)
			assert.Equal(t, tc.want.Password, got.Password)
			assert.Equal(t, tc.want.ConnectionURL, got.ConnectionURL)
			assert.Equal(t, tc.want.SrvConnectionURL, got.SrvConnectionURL)
			assert.Equal(t, tc.want.PrivateConnectionURLs, got.PrivateConnectionURLs)
		})
	}
}
