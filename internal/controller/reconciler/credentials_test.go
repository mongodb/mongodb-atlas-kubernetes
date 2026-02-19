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

	t.Run("service account secret without annotation returns error", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds", Namespace: "ns"},
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
		assert.Contains(t, err.Error(), "missing the atlas.mongodb.com/access-token annotation")
	})

	t.Run("service account secret with valid token", func(t *testing.T) {
		tokenSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds-token-abc", Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken": []byte("bearer-token-value"),
				"expiry":      []byte("2099-01-01T00:00:00Z"),
			},
		}
		credSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-creds",
				Namespace: "ns",
				Annotations: map[string]string{
					AccessTokenAnnotation: "sa-creds-token-abc",
				},
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
		tokenSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "sa-creds-token-abc", Namespace: "ns"},
			Data: map[string][]byte{
				"accessToken": []byte(""),
				"expiry":      []byte("2099-01-01T00:00:00Z"),
			},
		}
		credSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-creds",
				Namespace: "ns",
				Annotations: map[string]string{
					AccessTokenAnnotation: "sa-creds-token-abc",
				},
			},
			Data: map[string][]byte{
				"orgId":        []byte("org-123"),
				"clientId":     []byte("client-id"),
				"clientSecret": []byte("client-secret"),
			},
		}
		k8sClient := newFakeKubeClient(t, credSecret, tokenSecret)
		ref := client.ObjectKey{Name: "sa-creds", Namespace: "ns"}

		_, err := GetConnectionConfig(ctx, k8sClient, &ref, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty accessToken")
	})
}

func TestValidateConnectionConfig(t *testing.T) {
	t.Run("should be invalid and all missing data", func(t *testing.T) {
		missing, ok := validate(nil)
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"orgId", "publicApiKey", "privateApiKey"})
	})

	t.Run("should be invalid and organization id is missing", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"orgId"})
	})

	t.Run("should be invalid and public key id is missing", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PrivateKey: "secret"}}})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"publicApiKey"})
	})

	t.Run("should be invalid and private key id is missing", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local"}}})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"privateApiKey"})
	})

	t.Run("should be valid", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{OrgID: "some", Credentials: &atlas.Credentials{APIKeys: &atlas.APIKeys{PublicKey: "local", PrivateKey: "secret"}}})
		assert.True(t, ok)
		assert.Empty(t, missing)
	})

	t.Run("should be valid with service account token", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{
			OrgID: "some",
			Credentials: &atlas.Credentials{
				ServiceAccount: &atlas.ServiceAccountToken{BearerToken: "token"},
			},
		})
		assert.True(t, ok)
		assert.Empty(t, missing)
	})

	t.Run("should be invalid with service account token missing bearer", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{
			OrgID: "some",
			Credentials: &atlas.Credentials{
				ServiceAccount: &atlas.ServiceAccountToken{},
			},
		})
		assert.False(t, ok)
		assert.Contains(t, missing, "accessToken")
	})

	t.Run("should be invalid with service account token missing orgId", func(t *testing.T) {
		missing, ok := validate(&atlas.ConnectionConfig{
			Credentials: &atlas.Credentials{
				ServiceAccount: &atlas.ServiceAccountToken{BearerToken: "token"},
			},
		})
		assert.False(t, ok)
		assert.Contains(t, missing, "orgId")
	})
}

func newFakeKubeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}
