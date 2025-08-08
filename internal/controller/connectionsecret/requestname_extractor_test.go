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
	"fmt"
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

func TestLoadRequestNameParts(t *testing.T) {
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
	secretEmptyLabel := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "emptylabel-mycluster-admin",
			Namespace: "default",
			Labels: map[string]string{
				ProjectLabelKey: "",
				ClusterLabelKey: "mycluster",
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
			secretEmptyLabel,
			secretBadSplit,
			secretInvalidSep,
		).
		Build()

	tests := map[string]struct {
		name      string
		namespace string
		expected  RequestNameParts
		errSubstr string
	}{
		"valid internal format": {
			name: "proj123$mycluster$admin",
			expected: RequestNameParts{
				ProjectID:        "proj123",
				ClusterName:      "mycluster",
				DatabaseUsername: "admin",
			},
		},
		"internal format with too few parts": {
			name:      "proj123$clusterOnly",
			errSubstr: "expected 3 parts",
		},
		"internal format with empty part": {
			name:      "proj123$$admin",
			errSubstr: "empty value in one or more parts",
		},
		"valid k8s format": {
			name:      "myproj-mycluster-admin",
			namespace: "default",
			expected: RequestNameParts{
				ProjectID:        "proj123",
				ProjectName:      "myproj",
				ClusterName:      "mycluster",
				DatabaseUsername: "admin",
			},
		},
		"k8s format with missing secret": {
			name:      "nonexistent-secret",
			namespace: "default",
			errSubstr: "unable to retrieve Secret",
		},
		"k8s format with missing labels": {
			name:      "missing-mycluster-admin",
			namespace: "default",
			errSubstr: "missing required label",
		},
		"k8s format with empty label value": {
			name:      "emptylabel-mycluster-admin",
			namespace: "default",
			errSubstr: "has empty value for label",
		},
		"k8s format with invalid name separator": {
			name:      "invalid-separator",
			namespace: "default",
			errSubstr: "expected separator",
		},
		"k8s format with empty value after split": {
			name:      "-mycluster-admin",
			namespace: "default",
			errSubstr: "empty value in one or more parts",
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			ids, err := LoadRequestNameParts(context.Background(), client, types.NamespacedName{
				Name:      tc.name,
				Namespace: tc.namespace,
			})

			if tc.errSubstr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, ids)
			}
		})
	}
}

func TestPair_NeedsSDKProjectResolution(t *testing.T) {
	p := &ConnectionPair{
		Deployment: &akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ExternalProjectRef: &akov2.ExternalProjectReference{ID: "abc"},
				},
			},
		},
		User: &akov2.AtlasDatabaseUser{
			Spec: akov2.AtlasDatabaseUserSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ExternalProjectRef: &akov2.ExternalProjectReference{ID: "abc"},
				},
			},
		},
	}
	assert.True(t, p.NeedsSDKProjectResolution())

	p.User.Spec.ExternalProjectRef = nil
	assert.False(t, p.NeedsSDKProjectResolution())
}

func TestPair_AreResourcesReady(t *testing.T) {
	t.Run("Both ready", func(t *testing.T) {
		p := &ConnectionPair{
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
		}
		ok, notReady := p.AreResourcesReady()
		assert.True(t, ok)
		assert.Empty(t, notReady)
	})

	t.Run("One not ready", func(t *testing.T) {
		p := &ConnectionPair{
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
						Conditions: []api.Condition{{Type: api.DatabaseUserReadyType, Status: corev1.ConditionTrue}},
					},
				},
			},
		}
		ok, notReady := p.AreResourcesReady()
		assert.False(t, ok)
		assert.Equal(t, []string{"AtlasDatabaseUser/user"}, notReady)
	})
}

func TestLoadDeploymentAndUser(t *testing.T) {
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
		expectedErr      string
	}{
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
		},
		"no deployments found overall": {
			clusterName:      "clusterA",
			databaseUsername: "admin",
			users: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: ns},
					Spec:       akov2.AtlasDatabaseUserSpec{Username: "admin"},
				},
			},
			expectedErr: `expected 1 AtlasDeployment for "proj123-clusterA", found 0`,
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
			expectedErr: `expected 1 AtlasDeployment for "proj123-clusterB", found 0`,
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
			expectedErr: `expected 1 AtlasDatabaseUser for "proj123-admin", found 2`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			allObjects := append(tt.deployments, tt.users...)

			client := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(allObjects...).
				WithIndex(&akov2.AtlasDeployment{}, indexer.AtlasDeploymentBySpecNameAndProjectID, func(obj client.Object) []string {
					return []string{projectID + "-" + "clusterA"}
				}).
				WithIndex(&akov2.AtlasDatabaseUser{}, indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, func(obj client.Object) []string {
					return []string{projectID + "-" + "admin"}
				}).
				Build()

			ids := RequestNameParts{
				ProjectID:        projectID,
				ClusterName:      tt.clusterName,
				DatabaseUsername: tt.databaseUsername,
			}

			pair, err := LoadPairedResources(context.Background(), client, ids, ns)

			if tt.expectedErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, pair.Deployment)
				assert.NotNil(t, pair.User)
				assert.Equal(t, tt.clusterName, pair.Deployment.GetDeploymentName())
				assert.Equal(t, tt.databaseUsername, pair.User.Spec.Username)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}

			failIDs := RequestNameParts{
				ProjectID:        otherprojectID,
				ClusterName:      tt.clusterName,
				DatabaseUsername: tt.databaseUsername,
			}

			failPair, failErr := LoadPairedResources(context.Background(), client, failIDs, ns)
			assert.Error(t, failErr)
			assert.Nil(t, failPair)
			assert.Contains(t, failErr.Error(), fmt.Sprintf(`expected 1 AtlasDeployment for "%s-%s"`, otherprojectID, tt.clusterName))
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

	p := &ConnectionPair{
		Deployment: deployment,
		User:       user,
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
