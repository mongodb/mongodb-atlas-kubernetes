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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
)

func TestCreateK8sFormat(t *testing.T) {
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

func TestLoadRequestIdentifiers(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))

	secretValid := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myproj-mycluster-admin",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "proj123",
				ClusterLabelKey: "mycluster",
			},
		},
	}
	secretMissingLabels := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "missing-mycluster-admin",
			Namespace: "default",
			Labels:    map[string]string{},
		},
	}
	secretEmptyProject := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "emptyproject-mycluster-admin",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "",
				ClusterLabelKey: "mycluster",
			},
		},
	}
	secretEmptyCluster := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myproj--admin",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "proj123",
				ClusterLabelKey: "",
			},
		},
	}
	secretBadSplit := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "-mycluster-admin",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "proj123",
				ClusterLabelKey: "mycluster",
			},
		},
	}
	secretInvalidSep := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "invalid-separator",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "proj123",
				ClusterLabelKey: "unknown",
			},
		},
	}

	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			secretValid,
			secretMissingLabels,
			secretEmptyProject,
			secretEmptyCluster,
			secretBadSplit,
			secretInvalidSep,
		).
		Build()

	tests := map[string]struct {
		name        string
		namespace   string
		expected    ConnSecretIdentifiers
		expectedErr error
	}{
		"valid internal format": {
			name: "proj123$mycluster$admin",
			expected: ConnSecretIdentifiers{
				ProjectID:        "proj123",
				ClusterName:      "mycluster",
				DatabaseUsername: "admin",
			},
		},
		"internal format with too few parts": {
			name:        "proj123$clusterOnly",
			expectedErr: ErrInternalFormatPartsInvalid,
		},
		"internal format with empty part": {
			name:        "proj123$$admin",
			expectedErr: ErrInternalFormatPartEmpty,
		},
		"valid k8s format": {
			name:      "myproj-mycluster-admin",
			namespace: "default",
			expected: ConnSecretIdentifiers{
				ProjectID:        "proj123",
				ProjectName:      "myproj",
				ClusterName:      "mycluster",
				DatabaseUsername: "admin",
			},
		},
		"k8s format with missing labels": {
			name:        "missing-mycluster-admin",
			namespace:   "default",
			expectedErr: ErrK8sLabelsMissing,
		},
		"k8s format with empty project label": {
			name:        "emptyproject-mycluster-admin",
			namespace:   "default",
			expectedErr: ErrK8sLabelEmpty,
		},
		"k8s format with empty cluster label": {
			name:        "myproj--admin",
			namespace:   "default",
			expectedErr: ErrK8sLabelEmpty,
		},
		"k8s format with invalid name separator": {
			name:        "invalid-separator",
			namespace:   "default",
			expectedErr: ErrK8sNameSplitInvalid,
		},
		"k8s format with empty value after split": {
			name:        "-mycluster-admin",
			namespace:   "default",
			expectedErr: ErrK8sNameSplitEmpty,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			ids, err := LoadRequestIdentifiers(
				context.Background(),
				client,
				types.NamespacedName{Name: tc.name, Namespace: tc.namespace},
			)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, ids)
		})
	}
}

func TestPair_IsReady(t *testing.T) {
	t.Run("Both ready", func(t *testing.T) {
		p := &ConnSecretPair{
			Deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{Name: "dep"},
				Status: status.AtlasDeploymentStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			User: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "user"},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			ProjectID: "proj123",
		}
		ok, notReady := p.IsReady()
		assert.True(t, ok)
		assert.Empty(t, notReady)
	})

	t.Run("One not ready", func(t *testing.T) {
		p := &ConnSecretPair{
			Deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{Name: "dep"},
				Status: status.AtlasDeploymentStatus{
					Common: api.Common{
						Conditions: []api.Condition{{Type: api.ReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			User: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "user"},
				Status: status.AtlasDatabaseUserStatus{
					Common: api.Common{
						// Intentionally not ReadyType to simulate "not ready"
						Conditions: []api.Condition{{Type: api.DatabaseUserReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
			ProjectID: "proj123",
		}
		ok, notReady := p.IsReady()
		assert.False(t, ok)
		assert.Equal(t, []string{"AtlasDatabaseUser/user"}, notReady)
	})
}

func TestPair_LoadPairedResources(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(akov2.AddToScheme(scheme))

	const (
		ns             = "default"
		projectID      = "proj123"
		otherprojectID = "proj456"
	)

	tests := map[string]struct {
		clusterName      string
		databaseUsername string
		deployments      []client.Object
		users            []client.Object
		expectedErr      error
	}{
		"no deployments found overall": {
			clusterName:      "clusterA",
			databaseUsername: "admin",
			users: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
			},
			expectedErr: ErrNoDeploymentFound,
		},
		"no deployments found due to missing index": {
			clusterName:      "clusterB",
			databaseUsername: "admin",
			deployments: []client.Object{
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterB"},
					},
				},
			},
			users: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
			},
			expectedErr: ErrNoDeploymentFound,
		},
		"multiple users found": {
			clusterName:      "clusterA",
			databaseUsername: "admin",
			deployments: []client.Object{
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterA"},
					},
				},
			},
			users: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user2", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
			},
			expectedErr: ErrManyUsers,
		},
		"successfully finds one deployment and one user": {
			clusterName:      "clusterA",
			databaseUsername: "admin",
			deployments: []client.Object{
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: ns},
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{Name: "clusterA"},
					},
				},
			},
			users: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
			},
			expectedErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			allObjects := append(tt.deployments, tt.users...)

			cl := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(allObjects...).
				WithIndex(&akov2.AtlasDeployment{}, indexer.AtlasDeploymentBySpecNameAndProjectID, func(obj client.Object) []string {
					// Simulate an index only for projectID + "clusterA"
					return []string{projectID + "-" + "clusterA"}
				}).
				WithIndex(&akov2.AtlasDatabaseUser{}, indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, func(obj client.Object) []string {
					// Simulate an index only for projectID + "admin"
					return []string{projectID + "-" + "admin"}
				}).
				Build()

			ids := ConnSecretIdentifiers{
				ProjectID:        projectID,
				ClusterName:      tt.clusterName,
				DatabaseUsername: tt.databaseUsername,
			}

			pair, err := LoadPairedResources(context.Background(), cl, ids, ns)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, pair)
				assert.NotNil(t, pair.Deployment)
				assert.NotNil(t, pair.User)
				assert.Equal(t, tt.clusterName, pair.Deployment.GetDeploymentName())
				assert.Equal(t, tt.databaseUsername, pair.User.Spec.Username)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			// When the projectID doesn't match the indexed keys, BOTH resources are missing -> special error.
			failIDs := ConnSecretIdentifiers{
				ProjectID:        otherprojectID,
				ClusterName:      tt.clusterName,
				DatabaseUsername: tt.databaseUsername,
			}

			failPair, failErr := LoadPairedResources(context.Background(), cl, failIDs, ns)
			assert.Error(t, failErr)
			assert.Nil(t, failPair)
			assert.ErrorIs(t, failErr, ErrNoPairedResourcesFound)
		})
	}
}

func TestPair_BuildConnectionData(t *testing.T) {
	const (
		username      = "admin"
		passwordValue = "p@ssw0rd"
	)

	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "admin-password",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"password": []byte(passwordValue),
		},
	}

	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "admin",
			Namespace: "default",
		},
		Spec: akov2.AtlasDatabaseUserSpec{
			Username: username,
			PasswordSecret: &common.ResourceRef{
				Name: "admin-password",
			},
		},
	}

	deployment := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dep1",
			Namespace: "default",
		},
		Status: status.AtlasDeploymentStatus{
			ConnectionStrings: &status.ConnectionStrings{
				Standard:    "mongodb+srv://cluster.mongodb.net",
				StandardSrv: "mongodb://cluster.mongodb.net",
				Private:     "mongodb://private.mongodb.net",
				PrivateSrv:  "mongodb+srv://private.mongodb.net",
				PrivateEndpoint: []status.PrivateEndpoint{
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
				},
			},
		},
	}

	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(secret, user, deployment).
		Build()

	p := &ConnSecretPair{
		Deployment: deployment,
		User:       user,
		ProjectID:  "proj123",
	}

	data, err := p.BuildConnectionData(context.Background(), client)
	assert.NoError(t, err)
	assert.Equal(t, username, data.DBUserName)
	assert.Equal(t, passwordValue, data.Password)
	assert.Equal(t, "mongodb+srv://cluster.mongodb.net", data.ConnURL)
	assert.Equal(t, "mongodb://cluster.mongodb.net", data.SrvConnURL)
	assert.Len(t, data.PrivateConnURLs, 3)

	assert.Equal(t, "mongodb://private.mongodb.net", data.PrivateConnURLs[0].PvtConnURL)
	assert.Equal(t, "mongodb+srv://private.mongodb.net", data.PrivateConnURLs[0].PvtSrvConnURL)

	assert.Equal(t, "mongodb://pe1.mongodb.net", data.PrivateConnURLs[1].PvtConnURL)
	assert.Equal(t, "mongodb+srv://pe1.mongodb.net", data.PrivateConnURLs[1].PvtSrvConnURL)
	assert.Equal(t, "mongodb+srv://pe1-shard.mongodb.net", data.PrivateConnURLs[1].PvtShardConnURL)

	assert.Equal(t, "mongodb://pe2.mongodb.net", data.PrivateConnURLs[2].PvtConnURL)
	assert.Equal(t, "mongodb+srv://pe2.mongodb.net", data.PrivateConnURLs[2].PvtSrvConnURL)
	assert.Equal(t, "mongodb+srv://pe2-shard.mongodb.net", data.PrivateConnURLs[2].PvtShardConnURL)
}
