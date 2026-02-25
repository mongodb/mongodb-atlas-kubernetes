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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func TestAllowsByScopes(t *testing.T) {
	for _, tc := range []struct {
		name     string
		user     *generatedv1.DatabaseUser
		epName   string
		epType   string
		expected bool
	}{
		{
			name:     "nil user spec allows all",
			user:     &generatedv1.DatabaseUser{},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: true,
		},
		{
			name: "empty scopes allows all",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: true,
		},
		{
			name: "matching scope allows",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "my-cluster", Type: "CLUSTER"},
							},
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: true,
		},
		{
			name: "non-matching scope denies",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "other-cluster", Type: "CLUSTER"},
							},
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: false,
		},
		{
			name: "scope type mismatch denies",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "my-cluster", Type: "DATA_LAKE"},
							},
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: false,
		},
		{
			name: "multiple scopes with one matching allows",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "cluster-a", Type: "CLUSTER"},
								{Name: "my-cluster", Type: "CLUSTER"},
								{Name: "data-lake", Type: "DATA_LAKE"},
							},
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: true,
		},
		{
			name: "multiple scopes none matching denies",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "cluster-a", Type: "CLUSTER"},
								{Name: "cluster-b", Type: "CLUSTER"},
							},
						},
					},
				},
			},
			epName:   "my-cluster",
			epType:   "CLUSTER",
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := allowsByScopes(tc.user, tc.epName, tc.epType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetScopes(t *testing.T) {
	for _, tc := range []struct {
		name      string
		user      *generatedv1.DatabaseUser
		scopeType string
		expected  []string
	}{
		{
			name:      "nil spec returns nil",
			user:      &generatedv1.DatabaseUser{},
			scopeType: "CLUSTER",
			expected:  nil,
		},
		{
			name: "nil scopes returns nil",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
						},
					},
				},
			},
			scopeType: "CLUSTER",
			expected:  nil,
		},
		{
			name: "filters by scope type",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "cluster-a", Type: "CLUSTER"},
								{Name: "cluster-b", Type: "CLUSTER"},
								{Name: "data-lake", Type: "DATA_LAKE"},
							},
						},
					},
				},
			},
			scopeType: "CLUSTER",
			expected:  []string{"cluster-a", "cluster-b"},
		},
		{
			name: "returns empty for no matching type",
			user: &generatedv1.DatabaseUser{
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{
						Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
							Username: "test-user",
							Scopes: &[]generatedv1.Scopes{
								{Name: "cluster-a", Type: "CLUSTER"},
							},
						},
					},
				},
			},
			scopeType: "DATA_LAKE",
			expected:  nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := getScopes(tc.user, tc.scopeType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCheckAndDeleteStaleSecret(t *testing.T) {
	const (
		testNamespace = "default"
		testProjectID = "project-123"
	)

	for _, tc := range []struct {
		name          string
		secret        *corev1.Secret
		targets       []target.ConnectionTargetInstance
		expectDeleted bool
		wantErr       string
	}{
		{
			name: "secret matches existing target - not deleted",
			secret: newConnectionSecret("my-secret", testNamespace, map[string]string{
				TypeLabelKey:         CredLabelVal,
				ProjectLabelKey:      testProjectID,
				TargetLabelKey:       "my-cluster",
				DatabaseUserLabelKey: "admin",
			}, map[string]string{
				ConnectionTypelKey: "cluster",
			}),
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "my-cluster", targetType: "cluster"},
			},
			expectDeleted: false,
		},
		{
			name: "secret does not match any target - deleted",
			secret: newConnectionSecret("stale-secret", testNamespace, map[string]string{
				TypeLabelKey:         CredLabelVal,
				ProjectLabelKey:      testProjectID,
				TargetLabelKey:       "deleted-cluster",
				DatabaseUserLabelKey: "admin",
			}, map[string]string{
				ConnectionTypelKey: "cluster",
			}),
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "my-cluster", targetType: "cluster"},
			},
			expectDeleted: true,
		},
		{
			name: "secret target type mismatch - deleted",
			secret: newConnectionSecret("type-mismatch", testNamespace, map[string]string{
				TypeLabelKey:         CredLabelVal,
				ProjectLabelKey:      testProjectID,
				TargetLabelKey:       "my-cluster",
				DatabaseUserLabelKey: "admin",
			}, map[string]string{
				ConnectionTypelKey: "flexcluster",
			}),
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "my-cluster", targetType: "cluster"},
			},
			expectDeleted: true,
		},
		{
			name: "no targets - secret deleted",
			secret: newConnectionSecret("orphan-secret", testNamespace, map[string]string{
				TypeLabelKey:         CredLabelVal,
				ProjectLabelKey:      testProjectID,
				TargetLabelKey:       "any-cluster",
				DatabaseUserLabelKey: "admin",
			}, map[string]string{
				ConnectionTypelKey: "cluster",
			}),
			targets:       []target.ConnectionTargetInstance{},
			expectDeleted: true,
		},
		{
			name: "multiple targets with one match - not deleted",
			secret: newConnectionSecret("multi-target", testNamespace, map[string]string{
				TypeLabelKey:         CredLabelVal,
				ProjectLabelKey:      testProjectID,
				TargetLabelKey:       "cluster-b",
				DatabaseUserLabelKey: "admin",
			}, map[string]string{
				ConnectionTypelKey: "cluster",
			}),
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "cluster-a", targetType: "cluster"},
				&fakeTargetInstance{name: "cluster-b", targetType: "cluster"},
				&fakeTargetInstance{name: "cluster-c", targetType: "cluster"},
			},
			expectDeleted: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.secret).
				Build()

			r := &ConnectionSecretReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			ctx := context.Background()
			err := r.checkAndDeleteStaleSecret(ctx, tc.secret, tc.targets)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)

			// Check if secret still exists
			secret := &corev1.Secret{}
			getErr := fakeClient.Get(ctx, client.ObjectKey{
				Namespace: tc.secret.Namespace,
				Name:      tc.secret.Name,
			}, secret)

			if tc.expectDeleted {
				assert.Error(t, getErr, "secret should have been deleted")
			} else {
				assert.NoError(t, getErr, "secret should still exist")
			}
		})
	}
}

func TestCleanupStaleSecrets(t *testing.T) {
	const (
		testNamespace = "default"
		testProjectID = "project-123"
	)

	for _, tc := range []struct {
		name            string
		secrets         []*corev1.Secret
		targets         []target.ConnectionTargetInstance
		expectRemaining []string
		wantErr         string
	}{
		{
			name:            "no secrets - nothing to cleanup",
			secrets:         []*corev1.Secret{},
			targets:         []target.ConnectionTargetInstance{},
			expectRemaining: []string{},
		},
		{
			name: "all secrets match targets - none deleted",
			secrets: []*corev1.Secret{
				newConnectionSecret("secret-a", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "cluster-a",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
				newConnectionSecret("secret-b", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "cluster-b",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "cluster-a", targetType: "cluster"},
				&fakeTargetInstance{name: "cluster-b", targetType: "cluster"},
			},
			expectRemaining: []string{"secret-a", "secret-b"},
		},
		{
			name: "some secrets stale - stale ones deleted",
			secrets: []*corev1.Secret{
				newConnectionSecret("valid-secret", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "active-cluster",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
				newConnectionSecret("stale-secret", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "deleted-cluster",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{name: "active-cluster", targetType: "cluster"},
			},
			expectRemaining: []string{"valid-secret"},
		},
		{
			name: "all secrets stale - all deleted",
			secrets: []*corev1.Secret{
				newConnectionSecret("stale-1", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "old-cluster-1",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
				newConnectionSecret("stale-2", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      testProjectID,
					TargetLabelKey:       "old-cluster-2",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
			},
			targets:         []target.ConnectionTargetInstance{},
			expectRemaining: []string{},
		},
		{
			name: "secrets from different project - not affected",
			secrets: []*corev1.Secret{
				newConnectionSecret("other-project", testNamespace, map[string]string{
					TypeLabelKey:         CredLabelVal,
					ProjectLabelKey:      "other-project-id",
					TargetLabelKey:       "some-cluster",
					DatabaseUserLabelKey: "admin",
				}, map[string]string{
					ConnectionTypelKey: "cluster",
				}),
			},
			targets:         []target.ConnectionTargetInstance{},
			expectRemaining: []string{"other-project"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))

			clientBuilder := fake.NewClientBuilder().WithScheme(scheme)
			for _, s := range tc.secrets {
				clientBuilder = clientBuilder.WithObjects(s)
			}
			fakeClient := clientBuilder.Build()

			r := &ConnectionSecretReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			ctx := context.Background()
			err := r.cleanupStaleSecrets(ctx, testNamespace, tc.targets, testProjectID)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)

			// Check remaining secrets
			secretList := &corev1.SecretList{}
			require.NoError(t, fakeClient.List(ctx, secretList, client.InNamespace(testNamespace)))

			remainingNames := make([]string, 0, len(secretList.Items))
			for _, s := range secretList.Items {
				remainingNames = append(remainingNames, s.Name)
			}

			assert.ElementsMatch(t, tc.expectRemaining, remainingNames)
		})
	}
}

//nolint:unparam
func newConnectionSecret(name, namespace string, labels, annotations map[string]string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: map[string][]byte{
			"username": []byte("test-user"),
			"password": []byte("test-pass"),
		},
	}
}

// fakeTargetInstance implements target.ConnectionTargetInstance for testing
type fakeTargetInstance struct {
	name       string
	targetType string
	scopeType  string
	projectID  string
	ready      bool
	connData   *data.ConnectionSecret
}

func (f *fakeTargetInstance) GetConnectionTargetType() string {
	return f.targetType
}

func (f *fakeTargetInstance) GetName() string {
	return f.name
}

func (f *fakeTargetInstance) IsReady() bool {
	return f.ready
}

func (f *fakeTargetInstance) GetScopeType() string {
	return f.scopeType
}

func (f *fakeTargetInstance) GetProjectID(ctx context.Context) string {
	return f.projectID
}

func (f *fakeTargetInstance) BuildConnectionData(ctx context.Context) *data.ConnectionSecret {
	if f.connData != nil {
		return f.connData
	}
	// Return default connection data for ready targets
	if f.ready {
		return &data.ConnectionSecret{
			ConnectionURL:    "mongodb://cluster.mongodb.net",
			SrvConnectionURL: "mongodb+srv://cluster.mongodb.net",
		}
	}
	return nil
}

func TestHandleBatchUpsert(t *testing.T) {
	const (
		testNamespace   = "default"
		testProjectID   = "62b6e34b3d91647abb20e7b8"
		testProjectName = "my-project"
		testUserName    = "test-user"
		testClusterName = "my-cluster"
	)

	for _, tc := range []struct {
		name                string
		user                *generatedv1.DatabaseUser
		objects             []client.Object
		targets             []target.ConnectionTargetInstance
		atlasProvider       atlas.Provider
		wantErr             string
		wantSecretCreated   bool
		wantConditionStatus metav1.ConditionStatus
	}{
		{
			name: "user with nil spec entry returns error",
			user: &generatedv1.DatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testUserName,
					Namespace: testNamespace,
				},
				Spec: generatedv1.DatabaseUserSpec{
					V20250312: &generatedv1.DatabaseUserSpecV20250312{},
				},
			},
			targets: []target.ConnectionTargetInstance{},
			wantErr: "user spec has no entry",
		},
		{
			name: "ready target creates connection secret",
			user: newTestDatabaseUser(testUserName, testNamespace, testProjectID),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       testClusterName,
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      true,
				},
			},
			atlasProvider:       newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated:   true,
			wantConditionStatus: metav1.ConditionTrue,
		},
		{
			name: "target not ready skips secret creation",
			user: newTestDatabaseUser(testUserName, testNamespace, testProjectID),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       testClusterName,
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      false,
				},
			},
			atlasProvider:     newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated: false,
		},
		{
			name: "user with non-matching scope skips secret creation",
			user: newTestDatabaseUserWithScopes(testUserName, testNamespace, testProjectID, []generatedv1.Scopes{
				{Name: "other-cluster", Type: "CLUSTER"},
			}),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       testClusterName,
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      true,
				},
			},
			atlasProvider:     newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated: false,
		},
		{
			name: "user with matching scope creates secret",
			user: newTestDatabaseUserWithScopes(testUserName, testNamespace, testProjectID, []generatedv1.Scopes{
				{Name: testClusterName, Type: "CLUSTER"},
			}),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       testClusterName,
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      true,
				},
			},
			atlasProvider:       newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated:   true,
			wantConditionStatus: metav1.ConditionTrue,
		},
		{
			name: "multiple targets with mixed readiness",
			user: newTestDatabaseUser(testUserName, testNamespace, testProjectID),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       "cluster-a",
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      true,
				},
				&fakeTargetInstance{
					name:       "cluster-b",
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      false,
				},
			},
			atlasProvider:       newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated:   true,
			wantConditionStatus: metav1.ConditionTrue,
		},
		{
			name: "expired user deletes secret instead of creating",
			user: newTestDatabaseUserWithExpiry(testUserName, testNamespace, testProjectID, "2020-01-01T00:00:00Z"),
			objects: []client.Object{
				newTestPasswordSecret(testNamespace),
				newTestCredentialsSecret("atlas-credentials", testNamespace),
			},
			targets: []target.ConnectionTargetInstance{
				&fakeTargetInstance{
					name:       testClusterName,
					targetType: "cluster",
					scopeType:  "CLUSTER",
					projectID:  testProjectID,
					ready:      true,
				},
			},
			atlasProvider:     newMockAtlasProvider(t, testProjectID, testProjectName),
			wantSecretCreated: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))
			require.NoError(t, generatedv1.AddToScheme(scheme))

			clientBuilder := fake.NewClientBuilder().WithScheme(scheme)
			if tc.user != nil {
				clientBuilder = clientBuilder.WithObjects(tc.user).WithStatusSubresource(tc.user)
			}
			for _, obj := range tc.objects {
				clientBuilder = clientBuilder.WithObjects(obj)
			}
			fakeClient := clientBuilder.Build()

			logger := zaptest.NewLogger(t)
			r := &ConnectionSecretReconciler{
				Client:        fakeClient,
				Scheme:        scheme,
				Logger:        logger,
				AtlasProvider: tc.atlasProvider,
				GlobalSecretRef: client.ObjectKey{
					Name:      "atlas-credentials",
					Namespace: testNamespace,
				},
			}

			ctx := context.Background()
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      testUserName,
					Namespace: testNamespace,
				},
			}

			result, err := r.handleBatchUpsert(ctx, req, tc.user, testProjectID, tc.targets)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, ctrl.Result{}, result)

			// Check if connection secret was created
			secretList := &corev1.SecretList{}
			require.NoError(t, fakeClient.List(ctx, secretList, client.InNamespace(testNamespace)))

			// Filter for connection secrets (not password secrets)
			var connSecrets []corev1.Secret
			for _, s := range secretList.Items {
				if _, ok := s.Labels[TypeLabelKey]; ok {
					connSecrets = append(connSecrets, s)
				}
			}

			if tc.wantSecretCreated {
				assert.NotEmpty(t, connSecrets, "expected connection secret to be created")
			} else {
				assert.Empty(t, connSecrets, "expected no connection secret to be created")
			}

			// Check condition on user
			if tc.wantConditionStatus != "" {
				updatedUser := &generatedv1.DatabaseUser{}
				require.NoError(t, fakeClient.Get(ctx, client.ObjectKey{
					Namespace: testNamespace,
					Name:      testUserName,
				}, updatedUser))

				if updatedUser.Status.Conditions != nil {
					var found bool
					for _, c := range *updatedUser.Status.Conditions {
						if c.Type == ConnectionSecretReady {
							assert.Equal(t, tc.wantConditionStatus, c.Status)
							found = true
							break
						}
					}
					if tc.wantConditionStatus == metav1.ConditionTrue {
						assert.True(t, found, "expected ConnectionSecretReady condition")
					}
				}
			}
		})
	}
}

// Test helper functions

func newTestDatabaseUser(name, namespace, projectID string) *generatedv1.DatabaseUser {
	return &generatedv1.DatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: generatedv1.DatabaseUserSpec{
			V20250312: &generatedv1.DatabaseUserSpecV20250312{
				GroupId: &projectID,
				Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
					Username:     name,
					DatabaseName: "admin",
					PasswordSecretRef: &generatedv1.PasswordSecretRef{
						Name: "password-secret",
						Key:  pointer.MakePtr("password"),
					},
				},
			},
		},
		Status: generatedv1.DatabaseUserStatus{
			Conditions: &[]metav1.Condition{
				{
					Type:               state.ReadyCondition,
					Status:             metav1.ConditionTrue,
					LastTransitionTime: metav1.Now(),
					Reason:             "Ready",
				},
			},
		},
	}
}

func newTestDatabaseUserWithScopes(name, namespace, projectID string, scopes []generatedv1.Scopes) *generatedv1.DatabaseUser {
	user := newTestDatabaseUser(name, namespace, projectID)
	user.Spec.V20250312.Entry.Scopes = &scopes
	return user
}

func newTestDatabaseUserWithExpiry(name, namespace, projectID, deleteAfterDate string) *generatedv1.DatabaseUser {
	user := newTestDatabaseUser(name, namespace, projectID)
	user.Spec.V20250312.Entry.DeleteAfterDate = &deleteAfterDate
	return user
}

//nolint:unparam
func newTestPasswordSecret(namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "password-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"password": []byte("test-password"),
		},
	}
}

//nolint:unparam
func newTestCredentialsSecret(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"orgId":         []byte("test-org-id"),
			"publicApiKey":  []byte("test-public-key"),
			"privateApiKey": []byte("test-private-key"),
		},
	}
}

//nolint:unparam
func newMockAtlasProvider(t *testing.T, projectID, projectName string) atlas.Provider {
	projectsAPI := mockadmin.NewProjectsApi(t)
	projectsAPI.EXPECT().GetGroup(mock.Anything, projectID).
		Return(admin.GetGroupApiRequest{ApiService: projectsAPI}).Maybe()
	projectsAPI.EXPECT().GetGroupExecute(mock.Anything).
		Return(&admin.Group{
			Id:   &projectID,
			Name: projectName,
		}, nil, nil).Maybe()

	return &atlasmock.TestProvider{
		IsSupportedFunc: func() bool { return true },
		IsCloudGovFunc:  func() bool { return false },
		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
			return &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					ProjectsApi: projectsAPI,
				},
			}, nil
		},
	}
}
