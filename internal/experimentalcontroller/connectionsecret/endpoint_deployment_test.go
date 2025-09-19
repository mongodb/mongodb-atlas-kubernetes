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

func createDummyDeployment(t *testing.T) *akov2.AtlasDeployment {
	t.Helper()

	depl := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-depl",
			Namespace: "test-ns",
		},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "cluster1"},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name:      "test-project",
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

func TestDeploymentEndpoint_GetName(t *testing.T) {
	eNil := DeploymentEndpoint{obj: nil}
	assert.Equal(t, "", eNil.GetName())
	dep := createDummyDeployment(t)
	e := DeploymentEndpoint{obj: dep}
	assert.Equal(t, "cluster1", e.GetName())
}

func TestDeploymentEndpoint_IsReady(t *testing.T) {
	eNil := DeploymentEndpoint{obj: nil}
	assert.False(t, eNil.IsReady())

	notReady := &akov2.AtlasDeployment{
		Status: status.AtlasDeploymentStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: "False"}},
			},
		},
	}
	assert.False(t, DeploymentEndpoint{obj: notReady}.IsReady())

	ready := &akov2.AtlasDeployment{
		Status: status.AtlasDeploymentStatus{
			Common: api.Common{
				Conditions: []api.Condition{{Type: api.ReadyType, Status: "True"}},
			},
		},
	}
	assert.True(t, DeploymentEndpoint{obj: ready}.IsReady())
}

func TestDeploymentEndpoint_GetScopeType(t *testing.T) {
	e := DeploymentEndpoint{}
	assert.Equal(t, akov2.DeploymentScopeType, e.GetScopeType())
}

func TestDeploymentEndpoint_GetProjectID(t *testing.T) {
	r := createDummyEnv(t, nil)
	depl := createDummyDeployment(t)
	deplsdk := createDummyDeploymentSDK(t)

	tests := map[string]struct {
		endpoint DeploymentEndpoint
		want     string
		wantErr  bool
	}{
		"fail: nil deployment": {
			endpoint: DeploymentEndpoint{
				obj:             nil,
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"fail: project ref missing": {
			endpoint: DeploymentEndpoint{
				obj:             &akov2.AtlasDeployment{Spec: akov2.AtlasDeploymentSpec{}},
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"fail: k8s project ref but project not found": {
			endpoint: DeploymentEndpoint{
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
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			wantErr: true,
		},
		"success: external project ID": {
			endpoint: DeploymentEndpoint{
				obj:             deplsdk,
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			want: "test-project-id",
		},
		"success: k8s project ref": {
			endpoint: DeploymentEndpoint{
				obj:             depl,
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			},
			want: "test-project-id",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.endpoint.GetProjectID(context.Background())
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeploymentEndpoint_ListObj(t *testing.T) {
	e := DeploymentEndpoint{}
	l := e.ListObj()
	_, ok := l.(*akov2.AtlasDeploymentList)
	assert.True(t, ok)
}

func TestDeploymentEndpoint_SelectorByProject(t *testing.T) {
	e := DeploymentEndpoint{}
	s := e.SelectorByProject("p-1")
	assert.True(t, s.Matches(fields.Set{indexer.AtlasDeploymentByProject: "p-1"}))
	assert.False(t, s.Matches(fields.Set{indexer.AtlasDeploymentByProject: "other"}))
}

func TestDeploymentEndpoint_SelectorByProjectAndName(t *testing.T) {
	e := DeploymentEndpoint{}
	ids := &ConnSecretIdentifiers{ProjectID: "pX", ClusterName: "cY"}
	s := e.SelectorByProjectAndName(ids)
	assert.True(t, s.Matches(fields.Set{indexer.AtlasDeploymentBySpecNameAndProjectID: "pX-cY"}))
	assert.False(t, s.Matches(fields.Set{indexer.AtlasDeploymentBySpecNameAndProjectID: "pX-cZ"}))
}

func TestDeploymentEndpoint_ExtractList(t *testing.T) {
	r := createDummyEnv(t, nil)

	list := &akov2.AtlasDeploymentList{
		Items: []akov2.AtlasDeployment{
			{Spec: akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "a"}}},
			{Spec: akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "b"}}},
		},
	}

	e := DeploymentEndpoint{
		k8s:             r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
	}
	out, err := e.ExtractList(list)
	assert.NoError(t, err)
	if assert.Len(t, out, 2) {
		assert.Equal(t, "a", out[0].GetName())
		assert.Equal(t, "b", out[1].GetName())
	}

	_, err = e.ExtractList(&akov2.AtlasProjectList{})
	assert.Error(t, err)
}

func TestDeploymentEndpoint_BuildConnData(t *testing.T) {
	r := createDummyEnv(t, nil)
	depl := createDummyDeployment(t)
	user := createDummyUser(t, "test-user")

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
		endpoint *akov2.AtlasDeployment
		user     *akov2.AtlasDatabaseUser
		want     ConnSecretData
		wantErr  bool
	}{
		"fail: nil endpoint and user": {
			endpoint: nil,
			user:     nil,
			wantErr:  true,
		},
		"fail: missing password": {
			endpoint: depl,
			user:     userNoPass,
			wantErr:  true,
		},
		"success: builds from deployment connection strings": {
			endpoint: depl,
			user:     user,
			want: ConnSecretData{
				DBUserName: "admin",
				Password:   "secret",
				ConnURL:    "mongodb://std:27017",
				SrvConnURL: "mongodb+srv://std",
				PrivateConnURLs: []PrivateLinkConnURLs{
					{PvtConnURL: "mongodb://priv:27017", PvtSrvConnURL: "mongodb+srv://priv"},
					{PvtConnURL: "mongodb://pe1:27017", PvtSrvConnURL: "mongodb+srv://pe1", PvtShardConnURL: "mongodb+srv://pe1-shard"},
					{PvtConnURL: "mongodb://pe2:27017", PvtSrvConnURL: "mongodb+srv://pe2", PvtShardConnURL: "mongodb+srv://pe2-shard"},
				},
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e := DeploymentEndpoint{
				obj:             tc.endpoint,
				k8s:             r.Client,
				provider:        r.AtlasProvider,
				globalSecretRef: r.GlobalSecretRef,
				log:             r.Log,
			}
			got, err := e.BuildConnData(context.Background(), tc.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want.DBUserName, got.DBUserName)
			assert.Equal(t, tc.want.Password, got.Password)
			assert.Equal(t, tc.want.ConnURL, got.ConnURL)
			assert.Equal(t, tc.want.SrvConnURL, got.SrvConnURL)
			assert.Equal(t, tc.want.PrivateConnURLs, got.PrivateConnURLs)
		})
	}
}
