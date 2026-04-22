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
	"strings"
	"testing"

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

	t.Run("service account secret with no access token secret yet", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-creds",
				Namespace: "ns",
			},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte("client-id"),
				"clientSecret": []byte("client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, secret)
		ref := client.ObjectKey{Name: "sa-creds", Namespace: "ns"}

		_, err := GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist yet")
	})

	t.Run("service account secret with valid token", func(t *testing.T) {
		tokenSecretName, _ := DeriveAccessTokenSecretName("ns", "sa-creds")
		matchingHash, err := CredentialsHash("client-id", "client-secret")
		require.NoError(t, err)
		tokenSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: tokenSecretName, Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken":     []byte("bearer-token-value"),
				"expiry":          []byte("2099-01-01T00:00:00Z"),
				"credentialsHash": []byte(matchingHash),
			},
		}
		credSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-creds",
				Namespace: "ns",
			},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte("client-id"),
				"clientSecret": []byte("client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, credSecret, tokenSecret)
		ref := client.ObjectKey{Name: "sa-creds", Namespace: "ns"}

		cfg, err := GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.NoError(t, err)
		assert.Equal(t, "org-123", cfg.OrgID)
		require.NotNil(t, cfg.Credentials.ServiceAccount)
		assert.Equal(t, "bearer-token-value", cfg.Credentials.ServiceAccount.BearerToken)
		assert.Nil(t, cfg.Credentials.APIKeys)
	})

	t.Run("service account secret with stale token after credential rotation", func(t *testing.T) {
		tokenSecretName, _ := DeriveAccessTokenSecretName("ns", "sa-creds")
		staleHash, err := CredentialsHash("old-client-id", "old-client-secret")
		require.NoError(t, err)
		tokenSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: tokenSecretName, Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken":     []byte("stale-bearer-token"),
				"expiry":          []byte("2099-01-01T00:00:00Z"),
				"credentialsHash": []byte(staleHash),
			},
		}
		credSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds", Namespace: "ns"},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte("new-client-id"),
				"clientSecret": []byte("new-client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, credSecret, tokenSecret)
		ref := client.ObjectKey{Name: "sa-creds", Namespace: "ns"}

		_, err = GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is stale")
		assert.Contains(t, err.Error(), "credentials rotated")
	})

	t.Run("secret with both API keys and service account is rejected", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "bad-creds", Namespace: "ns"},
			Data: map[string][]byte{
				"orgId":         []byte("org-123"),
				"publicApiKey":  []byte("pub"),
				"privateApiKey": []byte("priv"),
				"clientId":      []byte("client-id"),
				"clientSecret":  []byte("client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, secret)
		ref := client.ObjectKey{Name: "bad-creds", Namespace: "ns"}

		_, err := GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "both API key and service account credentials")
	})

	t.Run("token secret with empty accessToken is rejected", func(t *testing.T) {
		tokenSecretName, _ := DeriveAccessTokenSecretName("ns", "sa-creds")
		matchingHash, err := CredentialsHash("client-id", "client-secret")
		require.NoError(t, err)
		tokenSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: tokenSecretName, Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken":     []byte(""),
				"expiry":          []byte("2099-01-01T00:00:00Z"),
				"credentialsHash": []byte(matchingHash),
			},
		}
		credSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds", Namespace: "ns"},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte("client-id"),
				"clientSecret": []byte("client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, credSecret, tokenSecret)
		ref := client.ObjectKey{Name: "sa-creds", Namespace: "ns"}

		_, err = GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty accessToken")
	})
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

func TestDeriveAccessTokenSecretName_Deterministic(t *testing.T) {
	const ns = "atlas-operator"
	const name = "my-sa-creds"
	first, _ := DeriveAccessTokenSecretName(ns, name)
	for i := 0; i < 20; i++ {
		second, _ := DeriveAccessTokenSecretName(ns, name)
		require.Equal(t, first, second,
			"function must return the same output for the same inputs on every call")
	}
	assert.True(t, strings.HasPrefix(first, "atlas-access-token-"))
	assert.Contains(t, first, name)
}

func TestDeriveAccessTokenSecretName_NamespaceSensitive(t *testing.T) {
	a, _ := DeriveAccessTokenSecretName("ns-a", "creds")
	b, _ := DeriveAccessTokenSecretName("ns-b", "creds")
	assert.NotEqual(t, a, b, "same name in different namespaces must yield different outputs")
}

func TestDeriveAccessTokenSecretName_NameSensitive(t *testing.T) {
	a, _ := DeriveAccessTokenSecretName("ns", "creds-a")
	b, _ := DeriveAccessTokenSecretName("ns", "creds-b")
	assert.NotEqual(t, a, b, "different names in same namespace must yield different outputs")
}

func TestDeriveAccessTokenSecretName_LengthFarPastLimit(t *testing.T) {
	longName := strings.Repeat("x", 500)
	result, _ := DeriveAccessTokenSecretName("ns", longName)
	assert.LessOrEqual(t, len(result), 253, "result must fit in DNS-1123 subdomain limit")
	assert.True(t, strings.HasPrefix(result, "atlas-access-token-"))
}

func TestDeriveAccessTokenSecretName_LengthAtBoundary(t *testing.T) {
	const ns = "ns"

	// A very long input forces truncation; the result must be exactly 253.
	veryLong, _ := DeriveAccessTokenSecretName(ns, strings.Repeat("a", 500))
	assert.Equal(t, 253, len(veryLong),
		"when truncation is forced, result length must be exactly 253 — guards off-by-one in maxNameLen")
	assert.True(t, strings.HasPrefix(veryLong, accessTokenSecretPrefix),
		"prefix must be preserved even under maximal truncation")

	// A name whose output fits in the 253-char budget must not be truncated
	// and must appear literally in the result.
	short, _ := DeriveAccessTokenSecretName(ns, "short-name")
	assert.LessOrEqual(t, len(short), 253)
	assert.Contains(t, short, "short-name",
		"short names are preserved literally")
}

func newFakeKubeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}
