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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

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
	dep := createDummyDeployment(t, "test-depl", "test-project", "cluster1")
	dbuser := createDummyUser(t, "test-user", "admin", "dummy-uid")
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
	dbUser := createDummyUser(t, "test-user", "admin", "dummy-uid")
	dep := createDummyDeployment(t, "test-depl", "test-project", "cluster1")

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
