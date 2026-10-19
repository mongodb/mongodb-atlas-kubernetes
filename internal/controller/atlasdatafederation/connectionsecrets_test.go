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

package atlasdatafederation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAddCredentialsToConnectionURL(t *testing.T) {
	testCases := []struct {
		name          string
		connURL       string
		userName      string
		password      string
		expected      string
		expectedError string
	}{
		{
			name:          "Valid Connection URL with Credentials",
			connURL:       "mongodb://mongodb0.example.com:27017/?ssl=true",
			userName:      "testUser",
			password:      "testPassword",
			expected:      "mongodb://testUser:testPassword@mongodb0.example.com:27017/?ssl=true",
			expectedError: "",
		},
		{
			name:          "Malformed URL",
			connURL:       "://invalid",
			userName:      "user",
			password:      "pass",
			expected:      "",
			expectedError: "missing protocol scheme",
		},
		{
			name:          "Special Characters in Credentials",
			connURL:       "mongodb://mongodb.example.com:27017/?ssl=true",
			userName:      "user@name",
			password:      "pass#word",
			expected:      "mongodb://user%40name:pass%23word@mongodb.example.com:27017/?ssl=true",
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := AddCredentialsToConnectionURL(tc.connURL, tc.userName, tc.password)

			if tc.expectedError != "" {
				assert.Error(t, err, "Expected an error for test case: %s", tc.name)
				assert.Contains(t, err.Error(), tc.expectedError, "Error message mismatch for test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for test case: %s", tc.name)
				assert.Equal(t, tc.expected, actual, "Unexpected output for test case: %s", tc.name)
			}
		})
	}
}

// TestDbUserBelongsToProjectWithExternalProjectRef is the regression test for the nil
// pointer panic: AtlasDatabaseUsers are listed cluster-wide, so a single user pointing
// at an external Atlas project used to panic every data federation reconcile in the
// cluster. externalProjectRef is a supported feature, so this needs no attacker.
func TestDbUserBelongsToProjectWithExternalProjectRef(t *testing.T) {
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: "my-ns"},
	}

	for _, tc := range []struct {
		name   string
		dbUser *akov2.AtlasDatabaseUser
		want   bool
	}{
		{
			name: "external project reference does not belong to any kube project",
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "user", Namespace: "my-ns"},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "5f4c1b2a3d4e5f6a7b8c9d0e"},
					},
				},
			},
			want: false,
		},
		{
			name: "kube reference in the same namespace belongs to the project",
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "user", Namespace: "my-ns"},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{Name: "my-project"},
					},
				},
			},
			want: true,
		},
		{
			name: "kube reference to another namespace does not belong to the project",
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Name: "user", Namespace: "my-ns"},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{Name: "my-project", Namespace: "other-ns"},
					},
				},
			},
			want: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				assert.Equal(t, tc.want, dbUserBelongsToProject(tc.dbUser, project))
			})
		})
	}
}

// TestAtlasProjectObjectKeyWithoutProjectRef pins the zero-value contract of the
// deprecated accessors: they must not panic when the resource uses an external project.
func TestAtlasProjectObjectKeyWithoutProjectRef(t *testing.T) {
	external := akov2.ProjectDualReference{
		ExternalProjectRef: &akov2.ExternalProjectReference{ID: "5f4c1b2a3d4e5f6a7b8c9d0e"},
	}
	meta := metav1.ObjectMeta{Name: "resource", Namespace: "my-ns"}

	t.Run("AtlasDatabaseUser", func(t *testing.T) {
		user := akov2.AtlasDatabaseUser{
			ObjectMeta: meta,
			Spec:       akov2.AtlasDatabaseUserSpec{ProjectDualReference: external},
		}
		assert.Equal(t, client.ObjectKey{}, user.AtlasProjectObjectKey())
	})

	t.Run("AtlasDeployment", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			ObjectMeta: meta,
			Spec:       akov2.AtlasDeploymentSpec{ProjectDualReference: external},
		}
		assert.Equal(t, client.ObjectKey{}, deployment.AtlasProjectObjectKey())
	})

	t.Run("AtlasPrivateEndpoint", func(t *testing.T) {
		pe := akov2.AtlasPrivateEndpoint{
			ObjectMeta: meta,
			Spec:       akov2.AtlasPrivateEndpointSpec{ProjectDualReference: external},
		}
		assert.Equal(t, client.ObjectKey{}, pe.AtlasProjectObjectKey())
	})
}
