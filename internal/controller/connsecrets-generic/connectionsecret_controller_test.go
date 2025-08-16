// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connsecretsgeneric

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

func TestConnectionSecretReconcile(t *testing.T) {
	type testCase struct {
		reqName          string
		deployment       *akov2.AtlasDeployment
		federation       *akov2.AtlasDataFederation
		user             *akov2.AtlasDatabaseUser
		expectedDeletion bool
		expectedUpdate   bool
		expectedResult   func() (ctrl.Result, error)
	}

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

	tests := map[string]testCase{
		"fail: could not load identifiers": {
			reqName: "my-project$cluster",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate("InvalidConnectionSecretName", ErrInternalFormatErr).ReconcileResult()
			},
		},
		"success: could not find secret with k8s format; assume deleted": {
			reqName: "test-project-id-cluster1-admin",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: missing both parent resources; garbage collect secret": {
			reqName:          "test-project-id$cluster1$admin",
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate(workflow.ConnSecretUnresolvedProjectName, ErrUnresolvedProjectName).ReconcileResult()
			},
		},
		"success: only one available resource from the pair, trigger delete": {
			reqName:          "test-project-id$cluster1$admin",
			user:             user,
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: invalid scopes; trigger delete": {
			reqName:    "test-project-id$cluster1$admin",
			deployment: depl,
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username: "admin",
					Scopes: []akov2.ScopeSpec{
						{
							Name: "df",
							Type: akov2.DataLakeScopeType,
						},
					},
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: expired user; trigger delete": {
			reqName:    "test-project-id$cluster1$admin",
			deployment: depl,
			user: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "admin",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Username:        "admin",
					DeleteAfterDate: time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339),
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"requque: resources are not ready yet": {
			reqName:    "test-project-id$cluster1$admin",
			deployment: depl,
			user: &akov2.AtlasDatabaseUser{
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
			},
			expectedResult: func() (ctrl.Result, error) {
				return workflow.InProgress(workflow.ConnSecretNotReady, "resources not ready").ReconcileResult()
			},
		},
		"success: pair ready; trigger upsert": {
			reqName:        "test-project-id$cluster1$admin",
			deployment:     depl,
			user:           user,
			expectedUpdate: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var all []client.Object
			if tc.deployment != nil {
				all = append(all, tc.deployment)
			}
			if tc.federation != nil {
				all = append(all, tc.federation)
			}
			if tc.user != nil {
				all = append(all, tc.user)
			}

			r := createDummyEnv(t, all)
			r.EndpointKinds = []Endpoint{DeploymentEndpoint{r: r}, FederationEndpoint{r: r}}

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "test-ns",
					Name:      tc.reqName,
				},
			}

			res, err := r.Reconcile(context.Background(), req)
			expRes, expErr := tc.expectedResult()

			assert.Equal(t, expRes, res)
			if expErr != nil {
				assert.EqualError(t, err, expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedUpdate {
				ids, err := r.loadIdentifiers(context.Background(), req.NamespacedName)
				require.NoError(t, err)
				ids.ProjectName = "my-project-name"

				expectedName := CreateK8sFormat(ids.ProjectName, ids.ClusterName, ids.DatabaseUsername)
				var outputSecret corev1.Secret
				getErr := r.Client.Get(context.Background(), types.NamespacedName{
					Namespace: "test-ns",
					Name:      expectedName,
				}, &outputSecret)
				assert.NoError(t, getErr, "expected secret %q to exist", expectedName)
			}

			if tc.expectedDeletion {
				ids, err := r.loadIdentifiers(context.Background(), req.NamespacedName)
				require.NoError(t, err)

				expectedName := CreateK8sFormat(ids.ProjectName, ids.ClusterName, ids.DatabaseUsername)
				var check corev1.Secret
				getErr := r.Client.Get(context.Background(), types.NamespacedName{
					Namespace: "test-ns",
					Name:      expectedName,
				}, &check)
				assert.True(t, apiErrors.IsNotFound(getErr), "expected secret %q to be deleted", expectedName)
			}
		})
	}
}

func Test_allowsByScopes(t *testing.T) {
	type args struct {
		epName string
		epType akov2.ScopeType
	}
	tests := map[string]struct {
		user *akov2.AtlasDatabaseUser
		args args
		want bool
	}{
		"allow: no scopes field (nil)": {
			user: &akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: nil}},
			args: args{epName: "clusterA", epType: akov2.DeploymentScopeType},
			want: true,
		},
		"allow: empty scopes slice": {
			user: &akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: []akov2.ScopeSpec{}}},
			args: args{epName: "clusterA", epType: akov2.DeploymentScopeType},
			want: true,
		},
		"allow: deployment scope matches name": {
			user: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{{Type: akov2.DeploymentScopeType, Name: "clusterA"}},
				},
			},
			args: args{epName: "clusterA", epType: akov2.DeploymentScopeType},
			want: true,
		},
		"deny: only data lake scopes present for deployment endpoint": {
			user: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{Type: akov2.DeploymentScopeType, Name: "clusterB"},
						{Type: akov2.DataLakeScopeType, Name: "clusterA"},
						{Type: akov2.DataLakeScopeType, Name: "df1"},
						{Type: akov2.DataLakeScopeType, Name: "df2"},
					},
				},
			},
			args: args{epName: "clusterA", epType: akov2.DeploymentScopeType},
			want: false,
		},
		"allow: multiple scopes where one matches deployment name": {
			user: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Scopes: []akov2.ScopeSpec{
						{Type: akov2.DeploymentScopeType, Name: "clusterX"},
						{Type: akov2.DeploymentScopeType, Name: "clusterA"},
						{Type: akov2.DataLakeScopeType, Name: "df1"},
					},
				},
			},
			args: args{epName: "clusterA", epType: akov2.DeploymentScopeType},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := allowsByScopes(tc.user, tc.args.epName, tc.args.epType)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_generateConnectionSecretRequests(t *testing.T) {
	type testCase struct {
		projectID string
		endpoints []Endpoint
		users     []akov2.AtlasDatabaseUser
		expect    []types.NamespacedName
	}

	const (
		projectID = "proj-1"
		ns1       = "ns-1"
		ns2       = "ns-2"
	)

	r := createDummyEnv(t, nil)

	depA := DeploymentEndpoint{
		r: r,
		obj: &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: "test-depl", Namespace: "test-ns"},
			Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "my-depl-name"}},
		},
	}
	df1 := FederationEndpoint{
		r: r,
		obj: &akov2.AtlasDataFederation{
			ObjectMeta: metav1.ObjectMeta{Name: "test-df", Namespace: "test-ns"},
			Spec:       akov2.DataFederationSpec{Name: "my-df-name"},
		},
	}

	userNoScopes := akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u1", Namespace: ns1},
		Spec:       akov2.AtlasDatabaseUserSpec{Username: "user1"},
	}
	userDepScopedMatch := akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u2", Namespace: ns2},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user2",
			Scopes:   []akov2.ScopeSpec{{Type: akov2.DeploymentScopeType, Name: "my-depl-name"}},
		},
	}
	userDepScopedNoMatch := akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u3", Namespace: ns1},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user3",
			Scopes:   []akov2.ScopeSpec{{Type: akov2.DeploymentScopeType, Name: "missing-depl"}},
		},
	}
	userDfScopedMatch := akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u4", Namespace: ns1},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user4",
			Scopes:   []akov2.ScopeSpec{{Type: akov2.DataLakeScopeType, Name: "my-df-name"}},
		},
	}

	tests := map[string]testCase{
		"no scopes; all endpoints allowed": {
			projectID: projectID,
			endpoints: []Endpoint{depA, df1},
			users:     []akov2.AtlasDatabaseUser{userNoScopes},
			expect: []types.NamespacedName{
				{Namespace: ns1, Name: "proj-1$my-depl-name$user1"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user1"},
			},
		},
		"deployment scoping filters correctly": {
			projectID: projectID,
			endpoints: []Endpoint{depA, df1},
			users:     []akov2.AtlasDatabaseUser{userDepScopedMatch, userDepScopedNoMatch},
			expect: []types.NamespacedName{
				{Namespace: ns2, Name: "proj-1$my-depl-name$user2"},
			},
		},
		"data lake scoping filters correctly with mixed users": {
			projectID: projectID,
			endpoints: []Endpoint{depA, df1},
			users:     []akov2.AtlasDatabaseUser{userNoScopes, userDfScopedMatch},
			expect: []types.NamespacedName{
				{Namespace: ns1, Name: "proj-1$my-depl-name$user1"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user1"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user4"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := r.generateConnectionSecretRequests(tc.projectID, tc.endpoints, tc.users)

			require := require.New(t)
			assert := assert.New(t)

			require.Len(got, len(tc.expect), "unexpected number of requests")

			gotSet := map[types.NamespacedName]struct{}{}
			for _, req := range got {
				gotSet[req.NamespacedName] = struct{}{}
			}
			for _, e := range tc.expect {
				_, ok := gotSet[e]
				assert.Truef(ok, "missing expected request %s/%s", e.Namespace, e.Name)
			}
		})
	}
}
