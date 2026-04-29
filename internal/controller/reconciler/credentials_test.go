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

package reconciler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func TestResolveConnectionConfig(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		title         string
		objects       []client.Object
		input         project.ProjectReferrerObject
		expected      *atlas.ConnectionConfig
		expectedError error
	}{
		{
			title: "fallback to global secret",
			// given an empty project reference
			input: &akov2.AtlasIPAccessList{},
			// we expect the credentials to match the global fallback secret
			expected: &atlas.ConnectionConfig{OrgID: "global", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "global", PrivateKey: "global"}}},
		},
		{
			title: "local connection secret reference",
			// given an AtlasIPAccessList referencing a local connection secret directly
			input: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-list",
					Namespace: "project-namespace",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "some-secret"},
					},
				},
			},
			// and a local secret in the same namespace
			objects: []client.Object{&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-secret",
					Namespace: "project-namespace",
				},
				Data: map[string][]byte{
					"orgId": []byte("some"), "publicApiKey": []byte("local"), "privateApiKey": []byte("secret"),
				},
			}},
			// we expect the credentials to match the local secret
			expected: &atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}},
		},
		{
			title: "project reference",
			// given an AtlasIPAccessList referencing an AtlasProject resource
			input: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-list",
					Namespace: "project-namespace",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "project",
							Namespace: "project-ns",
						},
					},
				},
			},
			objects: []client.Object{
				// and the AtlasProject resource referencing a local connection secret
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "project",
						Namespace: "project-ns",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "project-secret",
							Namespace: "project-ns",
						},
					},
				},
				// and a local secret
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "project-secret",
						Namespace: "project-ns",
					},
					Data: map[string][]byte{
						"orgId":         []byte("some"),
						"publicApiKey":  []byte("local"),
						"privateApiKey": []byte("secret"),
					},
				},
			},
			// we expect the credentials to match the local secret
			expected: &atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}},
		},
		{
			title: "project reference without namespace",
			// given an AtlasIPAccessList referencing an AtlasProject without specifying a namespace
			input: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-list",
					Namespace: "project-namespace",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "project",
						},
					},
				},
			},
			objects: []client.Object{
				// and the AtlasProject in the "project-namespace" namespace referencing a local connection secret
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "project",
						Namespace: "project-namespace",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "project-secret",
							Namespace: "project-ns",
						},
					},
				},
				// and a local secret
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "project-secret",
						Namespace: "project-ns",
					},
					Data: map[string][]byte{
						"orgId":         []byte("some"),
						"publicApiKey":  []byte("local"),
						"privateApiKey": []byte("secret"),
					},
				},
			},
			// we expect the credentials to match the local secret
			expected: &atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}},
		},
		{
			title: "project reference to non-existing project",
			// given an AtlasIPAccessList referencing an AtlasProject that does not exist
			input: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-list",
					Namespace: "project-namespace",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "project",
						},
					},
				},
			},
			// we expect an ErrAtlasProjectProjectNotFound error
			expectedError: ErrMissingKubeProject,
		},
		{
			title: "local connection secret reference if external project reference is set",
			input: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-role",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "testProjectID"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			objects: []client.Object{
				// and a local secret
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "local-secret",
						Namespace: "test-ns",
					},
					Data: map[string][]byte{
						"orgId":         []byte("some"),
						"publicApiKey":  []byte("local"),
						"privateApiKey": []byte("secret"),
					},
				},
			},
			expected: &atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}},
		},
		{
			title: "favor local connection secret over project reference",
			input: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-role",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "foo",
							Namespace: "bar",
						},
						ConnectionSecret: &api.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			objects: []client.Object{
				// and a local secret
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "local-secret",
						Namespace: "test-ns",
					},
					Data: map[string][]byte{
						"orgId":         []byte("some"),
						"publicApiKey":  []byte("local"),
						"privateApiKey": []byte("secret"),
					},
				},
			},
			// we expect the local secret to be used
			expected: &atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := newFakeKubeClient(t, append(tc.objects, tc.input, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"orgId":         []byte("global"),
					"publicApiKey":  []byte("global"),
					"privateApiKey": []byte("global"),
				},
			})...)

			r := AtlasReconciler{
				Client: fakeClient,
				GlobalSecretRef: client.ObjectKey{
					Namespace: "default",
					Name:      "secret",
				},
			}
			cfg, err := r.ResolveConnectionConfig(ctx, tc.input)
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}

func TestGetConnectionConfig_ServiceAccount(t *testing.T) {
	ctx := context.Background()

	// Pinned literal of accesstoken.DeriveSecretName("ns", "sa-creds").
	const tokenSecretName = "atlas-access-token-sa-creds-6cd4c4d5f7d8d84ff"
	// Pinned literal of accesstoken.CredentialsHash("client-id", "client-secret").
	const matchingHash = "3974328787184052522"
	// Pinned literal of accesstoken.CredentialsHash("old-client-id", "old-client-secret").
	const staleHash = "17764957043874091622"

	saCredSecret := func(clientID, clientSecret string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds", Namespace: "ns"},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte(clientID),
				"clientSecret": []byte(clientSecret),
			},
		}
	}
	tokenSecret := func(accessToken, expiry, hash string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: tokenSecretName, Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken":     []byte(accessToken),
				"expiry":          []byte(expiry),
				"credentialsHash": []byte(hash),
			},
		}
	}

	tokenRefStr := "ns/" + tokenSecretName
	_, parseExpiryErr := time.Parse(time.RFC3339, "not-a-timestamp")
	require.Error(t, parseExpiryErr)

	for _, tc := range []struct {
		name          string
		credSecret    *corev1.Secret
		tokenSecret   *corev1.Secret
		expectedError string
	}{
		{
			name:          "no access token secret yet",
			credSecret:    saCredSecret("client-id", "client-secret"),
			expectedError: fmt.Sprintf("access token secret %s does not exist yet", tokenRefStr),
		},
		{
			name:        "valid token returns service account credentials",
			credSecret:  saCredSecret("client-id", "client-secret"),
			tokenSecret: tokenSecret("bearer-token-value", "2099-01-01T00:00:00Z", matchingHash),
		},
		{
			name:          "stale token after credential rotation",
			credSecret:    saCredSecret("new-client-id", "new-client-secret"),
			tokenSecret:   tokenSecret("stale-bearer-token", "2099-01-01T00:00:00Z", staleHash),
			expectedError: fmt.Sprintf("access token secret %s is stale (credentials rotated); waiting for the service-account-token controller to refresh", tokenRefStr),
		},
		{
			name: "secret with both API keys and service account credentials is rejected",
			credSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "bad-creds", Namespace: "ns"},
				Data: map[string][]byte{
					"orgId":         []byte("org-123"),
					"publicApiKey":  []byte("pub"),
					"privateApiKey": []byte("priv"),
					"clientId":      []byte("client-id"),
					"clientSecret":  []byte("client-secret"),
				},
			},
			expectedError: "invalid connection secret ns/bad-creds: secret contains both API key and service account credentials; only one type is allowed",
		},
		{
			name:          "already-expired token is rejected",
			credSecret:    saCredSecret("client-id", "client-secret"),
			tokenSecret:   tokenSecret("bearer-token-value", "2000-01-01T00:00:00Z", matchingHash),
			expectedError: fmt.Sprintf("access token secret %s is expired (expiry: 2000-01-01T00:00:00Z); waiting for the service-account-token controller to refresh", tokenRefStr),
		},
		{
			name:          "empty expiry is rejected",
			credSecret:    saCredSecret("client-id", "client-secret"),
			tokenSecret:   tokenSecret("bearer-token-value", "", matchingHash),
			expectedError: fmt.Sprintf("access token secret %s has an empty expiry field", tokenRefStr),
		},
		{
			name:          "unparseable expiry is rejected",
			credSecret:    saCredSecret("client-id", "client-secret"),
			tokenSecret:   tokenSecret("bearer-token-value", "not-a-timestamp", matchingHash),
			expectedError: fmt.Sprintf("access token secret %s has an invalid expiry field %q: %s", tokenRefStr, "not-a-timestamp", parseExpiryErr),
		},
		{
			name:          "empty accessToken is rejected",
			credSecret:    saCredSecret("client-id", "client-secret"),
			tokenSecret:   tokenSecret("", "2099-01-01T00:00:00Z", matchingHash),
			expectedError: fmt.Sprintf("access token secret %s has an empty accessToken field", tokenRefStr),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			objects := []client.Object{tc.credSecret}
			if tc.tokenSecret != nil {
				objects = append(objects, tc.tokenSecret)
			}
			k8sClient := newFakeKubeClient(t, objects...)
			cfg, err := GetConnectionConfig(ctx, k8sClient, new(client.ObjectKeyFromObject(tc.credSecret)), nil)

			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, string(tc.credSecret.Data["orgId"]), cfg.OrgID)
			require.NotNil(t, cfg.Credentials.ServiceAccount)
			assert.Equal(t, string(tc.tokenSecret.Data["accessToken"]), cfg.Credentials.ServiceAccount.BearerToken)
			assert.Nil(t, cfg.Credentials.APIKeys)
		})
	}
}

func TestValidate(t *testing.T) {
	t.Run("should fail when secret has no data at all", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "orgId")
		assert.Contains(t, err.Error(), "publicApiKey")
		assert.Contains(t, err.Error(), "privateApiKey")
	})

	t.Run("should fail when orgId is missing", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"publicApiKey":  []byte("pub"),
			"privateApiKey": []byte("priv"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "orgId")
	})

	t.Run("should fail when publicApiKey is missing", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":         []byte("org-123"),
			"privateApiKey": []byte("priv"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "publicApiKey")
	})

	t.Run("should fail when privateApiKey is missing", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"publicApiKey": []byte("pub"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "privateApiKey")
	})

	t.Run("should succeed with complete API keys", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":         []byte("org-123"),
			"publicApiKey":  []byte("pub"),
			"privateApiKey": []byte("priv"),
		}})
		assert.NoError(t, err)
	})

	t.Run("should succeed with complete service account credentials", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("client-id"),
			"clientSecret": []byte("client-secret"),
		}})
		assert.NoError(t, err)
	})

	t.Run("should fail when clientId is missing", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientSecret": []byte("client-secret"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "clientId")
	})

	t.Run("should fail when clientSecret is missing", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":    []byte("org-123"),
			"clientId": []byte("client-id"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "clientSecret")
	})

	t.Run("should fail when service account secret is missing orgId", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"clientId":     []byte("client-id"),
			"clientSecret": []byte("client-secret"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "orgId")
	})

	t.Run("should fail when secret contains both API key and service account credentials", func(t *testing.T) {
		err := validateConnectionSecret(&corev1.Secret{Data: map[string][]byte{
			"orgId":         []byte("org-123"),
			"publicApiKey":  []byte("pub"),
			"privateApiKey": []byte("priv"),
			"clientId":      []byte("client-id"),
			"clientSecret":  []byte("client-secret"),
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "both API key and service account credentials")
	})
}

func newFakeKubeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}
