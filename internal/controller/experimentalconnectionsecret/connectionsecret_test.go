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
package experimentalconnectionsecret

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func Test_createK8sFormat(t *testing.T) {
	tests := map[string]struct {
		projectName      string
		clusterName      string
		databaseUsername string
		expected         string
	}{
		"normal values": {
			projectName:      "MyProject",
			clusterName:      "MyCluster",
			databaseUsername: "AdminUser",
			expected:         "myproject-mycluster-adminuser",
		},
		"already normalized": {
			projectName:      "proj",
			clusterName:      "cluster",
			databaseUsername: "user",
			expected:         "proj-cluster-user",
		},
		"values with spaces and caps": {
			projectName:      "Proj A",
			clusterName:      "Cluster B",
			databaseUsername: "Admin X",
			expected:         "proj-a-cluster-b-admin-x",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := CreateK8sFormat(tc.projectName, tc.clusterName, tc.databaseUsername)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCreateInternalFormat(t *testing.T) {
	tests := map[string]struct {
		projectID        string
		clusterName      string
		databaseUsername string
		expected         string
	}{
		"normal values": {
			projectID:        "proj123",
			clusterName:      "ClusterOne",
			databaseUsername: "DBUser",
			expected:         "proj123$clusterone$dbuser",
		},
		"cluster and user already normalized": {
			projectID:        "id456",
			clusterName:      "cluster",
			databaseUsername: "user",
			expected:         "id456$cluster$user",
		},
		"values with spaces": {
			projectID:        "id789",
			clusterName:      "CL X",
			databaseUsername: "U X",
			expected:         "id789$cl-x$u-x",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := CreateInternalFormat(tc.projectID, tc.clusterName, tc.databaseUsername)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_loadIdentifiers(t *testing.T) {
	type want struct {
		projectID        string
		projectName      string
		clusterName      string
		databaseUsername string
		err              error
	}

	tests := map[string]struct {
		reqName string
		ns      string
		secret  *corev1.Secret
		want    want
	}{
		"fail: internal format-invalid parts count": {
			reqName: "only" + InternalSeparator + "two",
			ns:      "default",
			want:    want{err: ErrInternalFormatErr},
		},
		"fail: internal format-empty part": {
			reqName: "p" + InternalSeparator + InternalSeparator + "u",
			ns:      "default",
			want:    want{err: ErrInternalFormatErr},
		},
		"success: internal format": {
			reqName: "proj123" + InternalSeparator + "mycluster" + InternalSeparator + "theuser",
			ns:      "default",
			want: want{
				projectID:        "proj123",
				projectName:      "",
				clusterName:      "mycluster",
				databaseUsername: "theuser",
				err:              nil,
			},
		},
		"fail: k8s format-missing labels": {
			reqName: "p-c-u",
			ns:      "ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "p-c-u",
					Namespace: "ns",
					Labels:    map[string]string{},
				},
			},
			want: want{err: ErrK8SFormatErr},
		},
		"fail: k8s format-empty labels": {
			reqName: "p-c-u",
			ns:      "ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "p-c-u",
					Namespace: "ns",
					Labels: map[string]string{
						ProjectLabelKey: "",
						ClusterLabelKey: "",
					},
				},
			},
			want: want{err: ErrK8SFormatErr},
		},
		"fail: k8s format-name split invalid": {
			reqName: "proj-notmatchingsep-user",
			ns:      "ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "proj-notmatchingsep-user",
					Namespace: "ns",
					Labels: map[string]string{
						ProjectLabelKey: "pid-1",
						ClusterLabelKey: "clusterX",
					},
				},
			},
			want: want{err: ErrK8SFormatErr},
		},
		"fail: k8s format-name split empty": {
			reqName: "-clusterY-user",
			ns:      "ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "-clusterY-user",
					Namespace: "ns",
					Labels: map[string]string{
						ProjectLabelKey: "pid-2",
						ClusterLabelKey: "clusterY",
					},
				},
			},
			want: want{err: ErrK8SFormatErr},
		},
		"success: k8s format": {
			reqName: "myproj-mycluster-admin",
			ns:      "test-ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myproj-mycluster-admin",
					Namespace: "test-ns",
					Labels: map[string]string{
						ProjectLabelKey: "test-project-id",
						ClusterLabelKey: "mycluster",
					},
				},
			},
			want: want{
				projectID:        "test-project-id",
				projectName:      "myproj",
				clusterName:      "mycluster",
				databaseUsername: "admin",
				err:              nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var objs []client.Object
			if tc.secret != nil {
				objs = append(objs, tc.secret)
			}
			r := createDummyEnv(t, objs)

			got, err := r.loadIdentifiers(context.Background(), types.NamespacedName{
				Name:      tc.reqName,
				Namespace: tc.ns,
			})

			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			if assert.NotNil(t, got) {
				assert.Equal(t, tc.want.projectID, got.ProjectID)
				assert.Equal(t, tc.want.projectName, got.ProjectName)
				assert.Equal(t, tc.want.clusterName, got.ClusterName)
				assert.Equal(t, tc.want.databaseUsername, got.DatabaseUsername)
			}
		})
	}
}

func Test_loadPair(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(scheme))

	const (
		ns             = "test-ns"
		projectID      = "test-project-id"
		otherProjectID = "proj456"
	)

	type fields struct {
		endpointObjs []client.Object
		users        []*akov2.AtlasDatabaseUser
	}

	tests := map[string]struct {
		clusterName      string
		databaseUsername string
		fields           fields
		expectedErr      error
		expectedPairNil  bool
		expectUserNil    bool
		expectEpNil      bool
	}{
		"fail: ambiguous-multiple users": {
			clusterName:      "clusterA",
			databaseUsername: "admin",
			fields: fields{
				endpointObjs: []client.Object{
					&akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
						Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterA"}},
					},
				},
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "u1", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "u2", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
				},
			},
			expectedErr:     ErrAmbiguousPairing,
			expectedPairNil: true,
		},
		"fail: ambiguous-multiple endpoints (2 deployments)": {
			clusterName:      "clusterB",
			databaseUsername: "root",
			fields: fields{
				endpointObjs: []client.Object{
					&akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep-a", Namespace: ns},
						Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterB"}},
					},
					&akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep-b", Namespace: ns},
						Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterB"}},
					},
				},
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "u", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "root"}},
				},
			},
			expectedErr:     ErrAmbiguousPairing,
			expectedPairNil: true,
		},
		"fail: ambiguous-multiple endpoints (deployment and federation share name)": {
			clusterName:      "clusterC",
			databaseUsername: "admin",
			fields: fields{
				endpointObjs: []client.Object{
					&akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep-a", Namespace: ns},
						Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterC"}},
					},
					&akov2.AtlasDataFederation{
						ObjectMeta: metav1.ObjectMeta{Name: "df-a", Namespace: ns},
						Spec:       akov2.DataFederationSpec{Name: "clusterC"},
					},
				},
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "u", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
				},
			},
			expectedErr:     ErrAmbiguousPairing,
			expectedPairNil: true,
		},
		"fail: both missing": {
			clusterName:      "clusterD",
			databaseUsername: "andrpac",
			fields: fields{
				endpointObjs: nil,
				users:        nil,
			},
			expectedErr:     ErrMissingPairing,
			expectedPairNil: true,
			expectUserNil:   true,
			expectEpNil:     true,
		},
		"fail: user present but endpoint missing": {
			clusterName:      "missing",
			databaseUsername: "admin",
			fields: fields{
				endpointObjs: nil,
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "u-only", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
				},
			},
			expectedErr: ErrMissingPairing,
			expectEpNil: true,
		},
		"fail: user absent but endpoint present": {
			clusterName:      "clusterE",
			databaseUsername: "missing",
			fields: fields{
				endpointObjs: []client.Object{
					&akov2.AtlasDataFederation{
						ObjectMeta: metav1.ObjectMeta{Name: "df", Namespace: ns},
						Spec:       akov2.DataFederationSpec{Name: "clusterE"},
					},
				},
				users: nil,
			},
			expectedErr:   ErrMissingPairing,
			expectUserNil: true,
		},
		"success: exactly one user and one endpoint": {
			clusterName:      "clusterF",
			databaseUsername: "admin",
			fields: fields{
				endpointObjs: []client.Object{
					&akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep", Namespace: ns},
						Spec:       akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterF"}},
					},
				},
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "uu", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
				},
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var all []client.Object
			all = append(all, tc.fields.endpointObjs...)
			for _, u := range tc.fields.users {
				all = append(all, u)
			}

			r := createDummyEnv(t, all)
			r.EndpointKinds = []Endpoint{
				DeploymentEndpoint{
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
				FederationEndpoint{
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				},
			}

			ids := &ConnSecretIdentifiers{
				ProjectID:        projectID,
				ClusterName:      tc.clusterName,
				DatabaseUsername: tc.databaseUsername,
			}

			pair, err := r.loadPair(context.Background(), ids)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedPairNil {
				assert.Nil(t, pair)
				return
			}

			if tc.expectUserNil {
				assert.Nil(t, pair.User)
			} else {
				if assert.NotNil(t, pair.User) {
					assert.Equal(t, tc.databaseUsername, pair.User.Spec.Username)
				}
			}
			if tc.expectEpNil {
				assert.Nil(t, pair.Endpoint)
			} else {
				assert.NotNil(t, pair.Endpoint)
			}
			assert.Equal(t, projectID, pair.ProjectID)

			missIDs := &ConnSecretIdentifiers{
				ProjectID:        otherProjectID,
				ClusterName:      tc.clusterName,
				DatabaseUsername: tc.databaseUsername,
			}
			missPair, missErr := r.loadPair(context.Background(), missIDs)
			assert.ErrorIs(t, missErr, ErrMissingPairing)
			assert.Nil(t, missPair)
		})
	}
}

func Test_resolveProjectName(t *testing.T) {
	const (
		ns        = "test-ns"
		projectID = "test-project-id"
	)

	type want struct {
		projectName string
		err         error
	}

	r := createDummyEnv(t, nil)
	var notFoundErr = apiErrors.NewNotFound(
		schema.GroupResource{Group: "atlas.mongodb.com", Resource: "atlasprojects"},
		"missing-proj",
	)

	tests := map[string]struct {
		ids  *ConnSecretIdentifiers
		pair *ConnSecretPair
		want want
	}{
		"fail: nil pair and ids without projectName": {
			ids:  &ConnSecretIdentifiers{ProjectID: projectID},
			pair: nil,
			want: want{
				projectName: "",
				err:         ErrUnresolvedProjectName,
			},
		},
		"fail: cannot resolve from endpoint; no user": {
			ids: &ConnSecretIdentifiers{ProjectID: projectID},
			pair: &ConnSecretPair{
				ProjectID: projectID,
				Endpoint: DeploymentEndpoint{
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
					obj: &akov2.AtlasDeployment{
						ObjectMeta: metav1.ObjectMeta{Name: "dep", Namespace: ns},
						Spec: akov2.AtlasDeploymentSpec{
							ProjectDualReference: akov2.ProjectDualReference{
								ProjectRef: &common.ResourceRefNamespaced{Name: "missing-proj"},
							},
						},
					},
				},
				User: nil,
			},
			want: want{
				projectName: "",
				err:         notFoundErr,
			},
		},
		"fail: cannot resolve from user; no endpoint": {
			ids: &ConnSecretIdentifiers{ProjectID: projectID},
			pair: &ConnSecretPair{
				ProjectID: projectID,
				User: &akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "u", Namespace: ns},
					Spec: akov2.AtlasDatabaseUserSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: "missing-proj"},
						},
					},
				},
				Endpoint: nil,
			},
			want: want{
				projectName: "",
				err:         notFoundErr,
			},
		},
		"success: ids carries projectName": {
			ids: &ConnSecretIdentifiers{
				ProjectID:   projectID,
				ProjectName: "my-project-name",
			},
			pair: nil,
			want: want{
				projectName: "my-project-name",
				err:         nil,
			},
		},
		"success: resolve from endpoint": {
			ids: &ConnSecretIdentifiers{ProjectID: projectID},
			pair: &ConnSecretPair{
				ProjectID: projectID,
				Endpoint: DeploymentEndpoint{
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
					obj:             createDummyDeployment(t),
				},
			},
			want: want{
				projectName: "my-project-name",
				err:         nil,
			},
		},
		"success: resolve from user": {
			ids: &ConnSecretIdentifiers{ProjectID: projectID},
			pair: &ConnSecretPair{
				ProjectID: projectID,
				User:      createDummyUser(t),
				Endpoint:  nil,
			},
			want: want{
				projectName: "my-project-name",
				err:         nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := r.resolveProjectName(context.Background(), tc.ids, tc.pair)

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err, err)
				assert.Equal(t, "", got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.projectName, got)
		})
	}
}

func Test_handleDelete(t *testing.T) {
	type expectedResult struct {
		expectedResult ctrl.Result
		expectedError  error
	}

	const (
		cluster     = "cluster1"
		username    = "admin"
		projectID   = "test-project-id"
		projectName = "my-project-name"
	)

	type testCase struct {
		ids    ConnSecretIdentifiers
		pair   ConnSecretPair
		result expectedResult
	}

	r := createDummyEnv(t, nil)
	dep := createDummyDeployment(t)
	user := createDummyUser(t)
	depEndpoint := DeploymentEndpoint{
		k8s:             r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
		obj:             dep,
	}

	tests := map[string]testCase{
		"fail: projects with no parents cannot be directly deleted": {
			ids: ConnSecretIdentifiers{
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      nil,
				Endpoint:  nil,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  ErrUnresolvedProjectName,
			},
		},
		"success: no secret present beforehand": {
			ids: ConnSecretIdentifiers{
				ProjectName:      "missing-proj",
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      user,
				Endpoint:  depEndpoint,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
		"success: delete existing secret": {
			ids: ConnSecretIdentifiers{
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      user,
				Endpoint:  depEndpoint,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "test-ns",
					Name:      "any",
				},
			}

			res, err := r.handleDelete(context.Background(), req, &tc.ids, &tc.pair)
			assert.Equal(t, tc.result.expectedResult, res)

			if tc.result.expectedError != nil {
				require.ErrorIs(t, err, tc.result.expectedError)
				return
			}
			require.NoError(t, err)

			if tc.pair.Endpoint == nil && tc.pair.User == nil {
				return
			}

			resolvedProjectName := tc.ids.ProjectName
			if resolvedProjectName == "" {
				resolvedProjectName, _ = tc.pair.Endpoint.GetProjectName(context.Background())
			}

			var s corev1.Secret
			secretName := CreateK8sFormat(resolvedProjectName, tc.ids.ClusterName, tc.ids.DatabaseUsername)
			getErr := r.Client.Get(context.Background(), types.NamespacedName{Namespace: "test-ns", Name: secretName}, &s)
			require.True(t, apiErrors.IsNotFound(getErr), "expected secret %s to be deleted", secretName)
		})
	}
}

func Test_handleUpsert(t *testing.T) {
	type expectedResult struct {
		expectedResult ctrl.Result
		expectedError  error
	}

	const (
		ns          = "test-ns"
		cluster     = "cluster1"
		username    = "admin"
		projectID   = "test-project-id"
		projectName = "my-project-name"
	)

	type testCase struct {
		ids    ConnSecretIdentifiers
		pair   ConnSecretPair
		result expectedResult
	}

	r := createDummyEnv(t, nil)
	dep := createDummyDeployment(t)
	user := createDummyUser(t)
	depEndpoint := DeploymentEndpoint{
		k8s:             r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
		obj:             dep,
	}

	tests := map[string]testCase{
		"fail: upserting requires project resolution": {
			ids: ConnSecretIdentifiers{
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      nil,
				Endpoint:  nil,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  ErrUnresolvedProjectName,
			},
		},
		"fail: cannot build data": {
			ids: ConnSecretIdentifiers{
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      nil,
				Endpoint:  depEndpoint,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  ErrMissingPairing,
			},
		},
		"success: upsert secret": {
			ids: ConnSecretIdentifiers{
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
				ProjectID:        projectID,
			},
			pair: ConnSecretPair{
				ProjectID: projectID,
				User:      user,
				Endpoint:  depEndpoint,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{Namespace: ns, Name: "any"},
			}

			res, err := r.handleUpsert(context.Background(), req, &tc.ids, &tc.pair)
			assert.Equal(t, tc.result.expectedResult, res)

			if tc.result.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.result.expectedError.Error())
				return
			}
			require.NoError(t, err)

			if tc.pair.Endpoint == nil || tc.pair.User == nil {
				return
			}
			resolvedProjectName := tc.ids.ProjectName
			if resolvedProjectName == "" {
				resolvedProjectName, _ = tc.pair.Endpoint.GetProjectName(context.Background())
			}

			var s corev1.Secret
			secretName := CreateK8sFormat(resolvedProjectName, tc.ids.ClusterName, tc.ids.DatabaseUsername)
			require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: secretName}, &s))

			require.Equal(t, CredLabelVal, s.Labels[TypeLabelKey])
			require.Equal(t, projectID, s.Labels[ProjectLabelKey])
			require.Equal(t, cluster, s.Labels[ClusterLabelKey])

			require.Equal(t, username, string(s.Data[userNameKey]))
			require.Equal(t, "secret", string(s.Data[passwordKey]))
		})
	}
}

func Test_ensureSecret(t *testing.T) {
	type expectedResult struct {
		expectedError error
	}

	const (
		ns          = "test-ns"
		cluster     = "cluster1"
		username    = "admin"
		projectID   = "test-project-id"
		projectName = "my-project-name"
	)

	r := createDummyEnv(t, nil)
	user := createDummyUser(t)

	connData := ConnSecretData{
		DBUserName: username,
		Password:   "newpassword",
		ConnURL:    "mongodb://cluster1.mongodb.net/?authSource=admin",
		SrvConnURL: "mongodb+srv://cluster1.mongodb.net/?authSource=admin",
		PrivateConnURLs: []PrivateLinkConnURLs{
			{
				PvtConnURL:      "mongodb://pe1.mongodb.net",
				PvtSrvConnURL:   "mongodb+srv://pe1.mongodb.net",
				PvtShardConnURL: "mongodb+srv://pe1-shard.mongodb.net",
			},
			{
				PvtConnURL:      "mongodb://pe2.mongodb.net",
				PvtSrvConnURL:   "mongodb+srv://pe2.mongodb.net",
				PvtShardConnURL: "mongodb+srv://pe2-shard.mongodb.net",
			},
		},
	}

	tests := map[string]struct {
		ids     ConnSecretIdentifiers
		pair    ConnSecretPair
		secrets []client.Object
		data    ConnSecretData
		result  expectedResult
	}{
		"fail: invalid URL bubbles up and prevents creation": {
			ids: ConnSecretIdentifiers{
				ProjectID:        projectID,
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair: ConnSecretPair{User: user},
			data: ConnSecretData{
				DBUserName: username,
				Password:   "test-pass",
				ConnURL:    "://\x00",
			},
			result: expectedResult{expectedError: fmt.Errorf("parse \"://\\x00\": net/url: invalid control character in URL")},
		},
		"success: create with private endpoints": {
			ids: ConnSecretIdentifiers{
				ProjectID:        projectID,
				ProjectName:      "new-project-name",
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair:   ConnSecretPair{User: user},
			data:   connData,
			result: expectedResult{expectedError: nil},
		},
		"success: update existing secret": {
			ids: ConnSecretIdentifiers{
				ProjectID:        projectID,
				ProjectName:      projectName,
				ClusterName:      cluster,
				DatabaseUsername: username,
			},
			pair:   ConnSecretPair{User: user},
			data:   connData,
			result: expectedResult{expectedError: nil},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := r.ensureSecret(context.Background(), &tc.ids, &tc.pair, tc.data)
			if tc.result.expectedError != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			secretName := CreateK8sFormat(tc.ids.ProjectName, tc.ids.ClusterName, tc.ids.DatabaseUsername)
			var s corev1.Secret
			getErr := r.Client.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: secretName}, &s)
			require.NoError(t, getErr)

			require.Equal(t, CredLabelVal, s.Labels[TypeLabelKey])
			require.Equal(t, projectID, s.Labels[ProjectLabelKey])
			require.Equal(t, cluster, s.Labels[ClusterLabelKey])

			require.Equal(t, username, string(s.Data[userNameKey]))
			require.Equal(t, tc.data.Password, string(s.Data[passwordKey]))

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
				want, _ := CreateURL(baseURL, username, tc.data.Password)
				require.Equal(t, want, string(s.Data[key]), "mismatch for %s", key)
			}
		})
	}
}
