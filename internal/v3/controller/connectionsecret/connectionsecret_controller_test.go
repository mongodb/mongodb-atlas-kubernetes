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

package connectionsecret

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	admin "go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"
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
		deployment       *[]akov2.AtlasDeployment
		user             *[]akov2.AtlasDatabaseUser
		expectedDeletion bool
		expectedUpdate   bool
		expectedResult   func() (ctrl.Result, error)
	}

	depl := createDummyDeployment(t)
	user := createDummyUser(t, "test-user")
	secondUser := createDummyUser(t, "test-user-second")

	tests := map[string]testCase{
		"fail: could not load identifiers 1": {
			reqName: "my-project$cluster",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate("InvalidConnectionSecretName", ErrInternalFormatErr).ReconcileResult()
			},
		},
		"fail: could not load identifiers 2": {
			reqName: "my-project$cluster$theuser",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate("InvalidConnectionSecretName", ErrInternalFormatErr).ReconcileResult()
			},
		},
		"success: could not find secret with k8s format; assume deleted": {
			reqName: "connection-442423",
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: only one available resource from the pair, trigger delete": {
			reqName:          "test-project-id$cluster1$admin$deployment",
			user:             &[]akov2.AtlasDatabaseUser{*user},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"fail: ambigous pair": {
			reqName:          "test-project-id$cluster1$admin$deployment",
			user:             &[]akov2.AtlasDatabaseUser{*user, *secondUser},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.Terminate("ConnSecretPairNotLoaded", ErrAmbiguousPairing).
					ReconcileResult()
			},
		},
		"success: invalid scopes; trigger delete": {
			reqName:    "test-project-id$cluster1$admin$deployment",
			deployment: &[]akov2.AtlasDeployment{*depl},
			user: &[]akov2.AtlasDatabaseUser{
				{
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
				}},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"success: expired user; trigger delete": {
			reqName:    "test-project-id$cluster1$admin$deployment",
			deployment: &[]akov2.AtlasDeployment{*depl},
			user: &[]akov2.AtlasDatabaseUser{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "admin",
						Namespace: "test-ns",
					},
					Spec: akov2.AtlasDatabaseUserSpec{
						Username:        "admin",
						DeleteAfterDate: time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339),
					},
				},
			},
			expectedDeletion: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
		"re-enque: resources are not ready yet": {
			reqName:    "test-project-id$cluster1$admin$deployment",
			deployment: &[]akov2.AtlasDeployment{*depl},
			user: &[]akov2.AtlasDatabaseUser{
				{
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
			},
			expectedResult: func() (ctrl.Result, error) {
				return workflow.InProgress(workflow.ConnectionSecretNotReady, "resources not ready").ReconcileResult()
			},
		},
		"success: pair ready; trigger upsert": {
			reqName:        "test-project-id$cluster1$admin$deployment",
			deployment:     &[]akov2.AtlasDeployment{*depl},
			user:           &[]akov2.AtlasDatabaseUser{*user},
			expectedUpdate: true,
			expectedResult: func() (ctrl.Result, error) {
				return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var all []client.Object

			// If `tc.deployment` is not nil, append its contents to `all`
			if tc.deployment != nil {
				for _, dep := range *tc.deployment { // Dereferencing the pointer to iterate over the slice
					all = append(all, &dep) // Adding individual deployments to `all`
				}
			}

			// If `tc.user` is not nil, append its contents to `all`
			if tc.user != nil {
				for _, usr := range *tc.user { // Dereferencing the pointer to iterate over the slice
					all = append(all, &usr) // Adding individual users to `all`
				}
			}

			r := createDummyEnv(t, all)
			r.ConnectionTargetKinds = []ConnectionTarget{
				DeploymentConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
				DataFederationConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
			}

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

				expectedName := K8sConnectionSecretName(ids.ProjectID, ids.TargetName, ids.DatabaseUsername, ids.ConnectionType)
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

				expectedName := K8sConnectionSecretName(ids.ProjectID, ids.TargetName, ids.DatabaseUsername, ids.ConnectionType)
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
		"deny: only data lake scopes present for deployment connectionTarget": {
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
		projectID         string
		connectionTargets []ConnectionTarget
		users             []akov2.AtlasDatabaseUser
		expect            []types.NamespacedName
	}

	const (
		projectID = "proj-1"
		ns1       = "ns-1"
		ns2       = "ns-2"
	)

	r := createDummyEnv(t, nil)

	depA := DeploymentConnectionTarget{
		client:          r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
		obj: &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: "test-depl", Namespace: "test-ns"},
			Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "my-depl-name"}},
		},
	}
	df1 := DataFederationConnectionTarget{
		client:          r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
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
		"no scopes; all connectionTargets allowed": {
			projectID:         projectID,
			connectionTargets: []ConnectionTarget{depA, df1},
			users:             []akov2.AtlasDatabaseUser{userNoScopes},
			expect: []types.NamespacedName{
				{Namespace: ns1, Name: "proj-1$my-depl-name$user1$deployment"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user1$data-federation"},
			},
		},
		"deployment scoping filters correctly": {
			projectID:         projectID,
			connectionTargets: []ConnectionTarget{depA, df1},
			users:             []akov2.AtlasDatabaseUser{userDepScopedMatch, userDepScopedNoMatch},
			expect: []types.NamespacedName{
				{Namespace: ns2, Name: "proj-1$my-depl-name$user2$deployment"},
			},
		},
		"data lake scoping filters correctly with mixed users": {
			projectID:         projectID,
			connectionTargets: []ConnectionTarget{depA, df1},
			users:             []akov2.AtlasDatabaseUser{userNoScopes, userDfScopedMatch},
			expect: []types.NamespacedName{
				{Namespace: ns1, Name: "proj-1$my-depl-name$user1$deployment"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user1$data-federation"},
				{Namespace: ns1, Name: "proj-1$my-df-name$user4$data-federation"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := r.generateConnectionSecretRequests(tc.projectID, tc.connectionTargets, tc.users)

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

func Test_newConnectionTargetMapFunc(t *testing.T) {
	type testCase struct {
		objs   []client.Object
		obj    client.Object
		expect []types.NamespacedName
	}

	const (
		projectID = "test-project-id"
	)

	depl := createDummyDeployment(t)
	df := createDummyFederation(t)
	userNoScopes := createDummyUser(t, "test-user")

	userScopedDep := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u2", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user2",
			Scopes:   []akov2.ScopeSpec{{Type: akov2.DeploymentScopeType, Name: "cluster1"}},
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{ID: projectID},
			},
		},
	}
	userScopedDf := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u3", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user3",
			Scopes:   []akov2.ScopeSpec{{Type: akov2.DataLakeScopeType, Name: "my-df-name"}},
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{ID: projectID},
			},
		},
	}
	userOtherProject := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u4", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user4",
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{ID: "proj-OTHER"},
			},
		},
	}

	tests := map[string]testCase{
		"deployment maps only users in same project and allowed by scopes": {
			objs: []client.Object{depl, userNoScopes, userScopedDep, userScopedDf, userOtherProject},
			obj:  depl,
			expect: []types.NamespacedName{
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "cluster1", "admin", "deployment")},
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "cluster1", "user2", "deployment")},
			},
		},
		"datafederation maps only users in same project and allowed by scopes": {
			objs: []client.Object{df, userNoScopes, userScopedDf, userScopedDep, userOtherProject},
			obj:  df,
			expect: []types.NamespacedName{
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "my-df-name", "admin", "data-federation")},
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "my-df-name", "user3", "data-federation")},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := createDummyEnv(t, tc.objs)
			r.ConnectionTargetKinds = []ConnectionTarget{
				DeploymentConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
				DataFederationConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
			}

			reqs := r.newConnectionTargetMapFunc(context.Background(), tc.obj)

			require.Len(t, reqs, len(tc.expect))
			got := make(map[types.NamespacedName]struct{}, len(reqs))
			for _, r := range reqs {
				got[r.NamespacedName] = struct{}{}
			}
			for _, e := range tc.expect {
				_, ok := got[e]
				assert.Truef(t, ok, "missing expected request %s/%s", e.Namespace, e.Name)
			}
		})
	}
}

func Test_newDatabaseUserMapFunc(t *testing.T) {
	type testCase struct {
		objs   []client.Object
		user   *akov2.AtlasDatabaseUser
		expect []types.NamespacedName
	}

	const (
		projectID = "test-project-id"
	)

	depl := createDummyDeployment(t)
	df := createDummyFederation(t)
	userNoScopes := createDummyUser(t, "test-user")

	userScoped := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: "u2", Namespace: "test-ns"},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: "user",
			Scopes: []akov2.ScopeSpec{
				{Type: akov2.DeploymentScopeType, Name: "cluster1"},
				{Type: akov2.DataLakeScopeType, Name: "my-df-name"},
			},
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{ID: projectID},
			},
		},
	}

	tests := map[string]testCase{
		"user without scopes; all connectionTargets in same project": {
			objs: []client.Object{depl, df},
			user: userNoScopes,
			expect: []types.NamespacedName{
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "cluster1", "admin", "deployment")},
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "my-df-name", "admin", "data-federation")},
			},
		},
		"user with scopes; only matching connectionTargets in same project": {
			objs: []client.Object{depl, df},
			user: userScoped,
			expect: []types.NamespacedName{
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "cluster1", "user", "deployment")},
				{Namespace: "test-ns", Name: NewConnectionSecretRequestName(projectID, "my-df-name", "user", "data-federation")},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			all := append([]client.Object{tc.user}, tc.objs...)
			r := createDummyEnv(t, all)
			r.ConnectionTargetKinds = []ConnectionTarget{
				DeploymentConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
				DataFederationConnectionTarget{
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
			}

			reqs := r.newDatabaseUserMapFunc(context.Background(), tc.user)

			require.Len(t, reqs, len(tc.expect))
			got := make(map[types.NamespacedName]struct{}, len(reqs))
			for _, rq := range reqs {
				got[rq.NamespacedName] = struct{}{}
			}
			for _, e := range tc.expect {
				_, ok := got[e]
				assert.Truef(t, ok, "missing expected request %s/%s", e.Namespace, e.Name)
			}
		})
	}
}
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
				TargetLabelKey:  "cluster1",
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

func createDummyUser(t *testing.T, kubernetesUsername string) *akov2.AtlasDatabaseUser {
	t.Helper()

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubernetesUsername,
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
