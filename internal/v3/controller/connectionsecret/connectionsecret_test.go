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
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func TestNewConnectionSecretRequestName(t *testing.T) {
	tests := map[string]struct {
		projectID        string
		targetName       string
		databaseUsername string
		connectionType   string
		expected         string
	}{
		"normal values": {
			projectID:        "proj123",
			targetName:       "ClusterOne",
			databaseUsername: "DBUser",
			connectionType:   "deployment",
			expected:         "proj123$clusterone$dbuser$deployment",
		},
		"cluster and user already normalized": {
			projectID:        "id456",
			targetName:       "cluster",
			databaseUsername: "user",
			connectionType:   "deployment",
			expected:         "id456$cluster$user$deployment",
		},
		"values with spaces": {
			projectID:        "id789",
			targetName:       "CL X",
			databaseUsername: "U X",
			connectionType:   "data-federation",
			expected:         "id789$cl-x$u-x$data-federation",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewConnectionSecretRequestName(tc.projectID, tc.targetName, tc.databaseUsername, tc.connectionType)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_loadIdentifiers(t *testing.T) {
	uniqueID := strings.ToLower(uuid.New().String()[0:6])
	type want struct {
		projectID        string
		targetName       string
		databaseUsername string
		connectionType   string
		err              error
	}

	tests := map[string]struct {
		reqName string
		ns      string
		secret  *corev1.Secret
		want    want
	}{
		"fail: internal format-invalid parts count 1": {
			reqName: "only" + InternalSeparator + "two",
			ns:      "default",
			want:    want{err: ErrInternalFormatErr},
		},
		"fail: internal format-invalid parts count 2": {
			reqName: "only" + InternalSeparator + "only" + InternalSeparator + "three",
			ns:      "default",
			want:    want{err: ErrInternalFormatErr},
		},
		"fail: internal format-empty part": {
			reqName: "p" + InternalSeparator + InternalSeparator + "u",
			ns:      "default",
			want:    want{err: ErrInternalFormatErr},
		},
		"success: internal format": {
			reqName: uniqueID + InternalSeparator + "mycluster" + InternalSeparator + "theuser" + InternalSeparator + "deployment",
			ns:      "default",
			want: want{
				projectID:        uniqueID,
				targetName:       "mycluster",
				databaseUsername: "theuser",
				connectionType:   "deployment",
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
						TargetLabelKey:  "",
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
						TargetLabelKey:  "clusterX",
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
						TargetLabelKey:  "clusterY",
					},
				},
			},
			want: want{err: ErrK8SFormatErr},
		},
		"success: k8s format": {
			reqName: "connection-42424",
			ns:      "test-ns",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "connection-42424",
					Namespace: "test-ns",
					Labels: map[string]string{
						ProjectLabelKey:      uniqueID,
						TargetLabelKey:       "mycluster",
						DatabaseUserLabelKey: "theuser",
					},
					Annotations: map[string]string{
						ConnectionTypelKey: "deployment",
					},
				},
			},
			want: want{
				projectID:        uniqueID,
				targetName:       "mycluster",
				databaseUsername: "theuser",
				connectionType:   "deployment",
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
				assert.Equal(t, tc.want.targetName, got.TargetName)
				assert.Equal(t, tc.want.databaseUsername, got.DatabaseUsername)
				assert.Equal(t, tc.want.connectionType, got.ConnectionType)
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
		connectionTargetObjs []client.Object
		users                []*akov2.AtlasDatabaseUser
	}

	tests := map[string]struct {
		targetName       string
		databaseUsername string
		fields           fields
		expectedErr      error
		expectedPairNil  bool
		expectUserNil    bool
		expectEpNil      bool
	}{
		"fail: ambiguous-multiple users": {
			targetName:       "clusterA",
			databaseUsername: "admin",
			fields: fields{
				connectionTargetObjs: []client.Object{
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
		"fail: ambiguous-multiple connectionTargets (2 deployments)": {
			targetName:       "clusterB",
			databaseUsername: "root",
			fields: fields{
				connectionTargetObjs: []client.Object{
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
		"fail: ambiguous-multiple connectionTargets (deployment and federation share name)": {
			targetName:       "clusterC",
			databaseUsername: "admin",
			fields: fields{
				connectionTargetObjs: []client.Object{
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
			targetName:       "clusterD",
			databaseUsername: "andrpac",
			fields: fields{
				connectionTargetObjs: nil,
				users:                nil,
			},
			expectedErr:     ErrMissingPairing,
			expectedPairNil: true,
			expectUserNil:   true,
			expectEpNil:     true,
		},
		"fail: user present but connectionTarget missing": {
			targetName:       "missing",
			databaseUsername: "admin",
			fields: fields{
				connectionTargetObjs: nil,
				users: []*akov2.AtlasDatabaseUser{
					{ObjectMeta: metav1.ObjectMeta{Name: "u-only", Namespace: ns}, Spec: akov2.AtlasDatabaseUserSpec{Username: "admin"}},
				},
			},
			expectedErr: ErrMissingPairing,
			expectEpNil: true,
		},
		"fail: user absent but connectionTarget present": {
			targetName:       "clusterE",
			databaseUsername: "missing",
			fields: fields{
				connectionTargetObjs: []client.Object{
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
		"success: exactly one user and one connectionTarget": {
			targetName:       "clusterF",
			databaseUsername: "admin",
			fields: fields{
				connectionTargetObjs: []client.Object{
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
			all = append(all, tc.fields.connectionTargetObjs...)
			for _, u := range tc.fields.users {
				all = append(all, u)
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

			ids := &ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       tc.targetName,
				DatabaseUsername: tc.databaseUsername,
			}

			user, connectionTarget, err := r.loadPair(context.Background(), ids)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedPairNil {
				assert.Nil(t, user)
				assert.Nil(t, connectionTarget)
				return
			}

			if tc.expectUserNil {
				assert.Nil(t, user)
			} else {
				if assert.NotNil(t, user) {
					assert.Equal(t, tc.databaseUsername, user.Spec.Username)
				}
			}
			if tc.expectEpNil {
				assert.Nil(t, connectionTarget)
			} else {
				assert.NotNil(t, connectionTarget)
			}
			// assert.Equal(t, projectID, pair.ProjectID)

			missIDs := &ConnectionSecretIdentifiers{
				ProjectID:        otherProjectID,
				TargetName:       tc.targetName,
				DatabaseUsername: tc.databaseUsername,
			}
			missUser, missConnectionTarget, missErr := r.loadPair(context.Background(), missIDs)
			assert.ErrorIs(t, missErr, ErrMissingPairing)
			assert.Nil(t, missUser)
			assert.Nil(t, missConnectionTarget)
		})
	}
}

func Test_handleDelete(t *testing.T) {
	type expectedResult struct {
		expectedResult ctrl.Result
		expectedError  error
	}

	const (
		cluster        = "cluster1"
		username       = "admin"
		projectID      = "test-project-id"
		connectionType = "deployment"
	)

	type testCase struct {
		ids              ConnectionSecretIdentifiers
		result           expectedResult
		user             *akov2.AtlasDatabaseUser
		connectionTarget ConnectionTarget
	}

	r := createDummyEnv(t, nil)

	tests := map[string]testCase{
		"success: no secret present beforehand": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        "missing-proj",
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  nil,
			},
		},
		"success: delete existing secret": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
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

			res, err := r.handleDelete(context.Background(), req, &tc.ids)
			assert.Equal(t, tc.result.expectedResult, res)

			if tc.result.expectedError != nil {
				require.ErrorIs(t, err, tc.result.expectedError)
				return
			}
			require.NoError(t, err)

			if tc.connectionTarget == nil && tc.user == nil {
				return
			}

			var s corev1.Secret
			secretName := K8sConnectionSecretName(tc.ids.ProjectID, tc.ids.TargetName, tc.ids.DatabaseUsername, tc.ids.ConnectionType)
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
		ns             = "test-ns"
		cluster        = "cluster1"
		username       = "admin"
		projectID      = "test-project-id"
		connectionType = "deployment"
	)

	type testCase struct {
		ids              ConnectionSecretIdentifiers
		result           expectedResult
		user             *akov2.AtlasDatabaseUser
		connectionTarget ConnectionTarget
	}

	r := createDummyEnv(t, nil)
	dep := createDummyDeployment(t)
	dbuser := createDummyUser(t, "test-user")
	depConnectionTarget := DeploymentConnectionTarget{
		client:          r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
		obj:             dep,
	}

	tests := map[string]testCase{
		"fail: cannot build data": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			user:             nil,
			connectionTarget: depConnectionTarget,
			result: expectedResult{
				expectedResult: ctrl.Result{},
				expectedError:  ErrMissingPairing,
			},
		},
		"success: upsert secret": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			user:             dbuser,
			connectionTarget: depConnectionTarget,
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

			res, err := r.handleUpsert(context.Background(), req, &tc.ids, tc.user, tc.connectionTarget)
			assert.Equal(t, tc.result.expectedResult, res)

			if tc.result.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.result.expectedError.Error())
				return
			}
			require.NoError(t, err)

			if tc.connectionTarget == nil || tc.user == nil {
				return
			}

			var s corev1.Secret
			secretName := K8sConnectionSecretName(tc.ids.ProjectID, tc.ids.TargetName, tc.ids.DatabaseUsername, tc.ids.ConnectionType)
			require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: secretName}, &s))

			require.Equal(t, CredLabelVal, s.Labels[TypeLabelKey])
			require.Equal(t, projectID, s.Labels[ProjectLabelKey])
			require.Equal(t, cluster, s.Labels[TargetLabelKey])

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
		ns             = "test-ns"
		cluster        = "cluster1"
		username       = "admin"
		projectID      = "test-project-id"
		connectionType = "deployment"
	)

	r := createDummyEnv(t, nil)
	dbUser := createDummyUser(t, "test-user")
	dep := createDummyDeployment(t)

	connData := ConnectionSecretData{
		DBUserName:       username,
		Password:         "newpassword",
		ConnectionURL:    "mongodb://cluster1.mongodb.net/?authSource=admin",
		SrvConnectionURL: "mongodb+srv://cluster1.mongodb.net/?authSource=admin",
		PrivateConnectionURLs: []PrivateLinkConnectionURLs{
			{
				ConnectionURL:      "mongodb://pe1.mongodb.net",
				SrvConnectionURL:   "mongodb+srv://pe1.mongodb.net",
				ShardConnectionURL: "mongodb+srv://pe1-shard.mongodb.net",
			},
			{
				ConnectionURL:      "mongodb://pe2.mongodb.net",
				SrvConnectionURL:   "mongodb+srv://pe2.mongodb.net",
				ShardConnectionURL: "mongodb+srv://pe2-shard.mongodb.net",
			},
		},
	}

	depConnectionTarget := DeploymentConnectionTarget{
		client:          r.Client,
		provider:        r.AtlasProvider,
		globalSecretRef: r.GlobalSecretRef,
		log:             r.Log,
		obj:             dep,
	}

	tests := map[string]struct {
		ids              ConnectionSecretIdentifiers
		secrets          []client.Object
		data             ConnectionSecretData
		result           expectedResult
		user             *akov2.AtlasDatabaseUser
		connectionTarget ConnectionTarget
	}{
		"fail: invalid URL bubbles up and prevents creation": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			user:             dbUser,
			connectionTarget: depConnectionTarget,
			data: ConnectionSecretData{
				DBUserName:    username,
				Password:      "test-pass",
				ConnectionURL: "://\x00",
			},
			result: expectedResult{expectedError: fmt.Errorf("parse \"://\\x00\": net/url: invalid control character in URL")},
		},
		"success: create with private connectionTargets": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			user:             dbUser,
			connectionTarget: depConnectionTarget,
			data:             connData,
			result:           expectedResult{expectedError: nil},
		},
		"success: update existing secret": {
			ids: ConnectionSecretIdentifiers{
				ProjectID:        projectID,
				TargetName:       cluster,
				DatabaseUsername: username,
				ConnectionType:   connectionType,
			},
			user:             dbUser,
			connectionTarget: depConnectionTarget,
			data:             connData,
			result:           expectedResult{expectedError: nil},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := r.ensureSecret(context.Background(), &tc.ids, tc.user, tc.connectionTarget, tc.data)
			if tc.result.expectedError != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			secretName := K8sConnectionSecretName(tc.ids.ProjectID, tc.ids.TargetName, tc.ids.DatabaseUsername, tc.ids.ConnectionType)
			var s corev1.Secret
			getErr := r.Client.Get(context.Background(), types.NamespacedName{Namespace: ns, Name: secretName}, &s)
			require.NoError(t, getErr)

			require.Equal(t, CredLabelVal, s.Labels[TypeLabelKey])
			require.Equal(t, projectID, s.Labels[ProjectLabelKey])
			require.Equal(t, cluster, s.Labels[TargetLabelKey])

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
