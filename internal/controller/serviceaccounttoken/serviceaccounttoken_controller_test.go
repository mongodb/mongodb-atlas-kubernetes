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

package serviceaccounttoken_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/accesstoken"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/serviceaccounttoken"
)

type fakeTokenProvider struct {
	token  string
	expiry time.Time
	err    error
	calls  int
}

func (f *fakeTokenProvider) FetchToken(_ context.Context, _, _ string) (string, time.Time, error) {
	f.calls++
	return f.token, f.expiry, f.err
}

func newScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
}

func newReconciler(t *testing.T, tp serviceaccounttoken.TokenProvider, objs ...client.Object) (*serviceaccounttoken.ServiceAccountTokenReconciler, client.Client) {
	t.Helper()
	scheme := newScheme(t)
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	r := &serviceaccounttoken.ServiceAccountTokenReconciler{
		Client:        k8sClient,
		Scheme:        scheme,
		Log:           zap.NewNop().Sugar(),
		TokenProvider: tp,
	}
	return r, k8sClient
}

func TestReconcile_SkipsNonServiceAccountSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "api-key-secret", Namespace: "ns"},
		Data: map[string][]byte{
			"orgId":         []byte("org-123"),
			"publicApiKey":  []byte("pub"),
			"privateApiKey": []byte("priv"),
		},
	}
	tp := &fakeTokenProvider{}
	r, _ := newReconciler(t, tp, secret)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "api-key-secret", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.Equal(t, ctrl.Result{}, result)
	assert.Equal(t, 0, tp.calls)
}

func TestReconcile_CreatesTokenSecretOnFirstRun(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			UID:       types.UID("test-uid"),
		},
		Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("my-client-id"),
			"clientSecret": []byte("my-client-secret"),
		},
	}
	expiry := time.Now().Add(1 * time.Hour)
	tp := &fakeTokenProvider{token: "access-token-value", expiry: expiry}
	r, k8sClient := newReconciler(t, tp, secret)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.True(t, result.RequeueAfter > 0)
	assert.Equal(t, 1, tp.calls)

	expectedTokenName, _ := accesstoken.DeriveSecretName(secret.Namespace, secret.Name)

	tokenSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: expectedTokenName, Namespace: "ns"}, tokenSecret))

	assert.Equal(t, "access-token-value", string(tokenSecret.Data["accessToken"]))
	assert.NotEmpty(t, tokenSecret.Data["expiry"])
	assert.NotEmpty(t, tokenSecret.Data["credentialsHash"],
		"created token must record a hash of the source credentials for staleness detection")
	assert.Equal(t, secretservice.CredLabelVal, tokenSecret.Labels[secretservice.TypeLabelKey])

	require.Len(t, tokenSecret.OwnerReferences, 1)
	assert.Equal(t, "sa-creds", tokenSecret.OwnerReferences[0].Name)
	assert.Equal(t, types.UID("test-uid"), tokenSecret.OwnerReferences[0].UID)

	updatedSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: "sa-creds", Namespace: "ns"}, updatedSecret))
	assert.Empty(t, updatedSecret.Annotations,
		"credential secret must not be mutated by the controller")
}

func TestReconcile_RefreshesExpiredToken(t *testing.T) {
	expiredExpiry := time.Now().Add(-10 * time.Minute)
	tokenSecretName, _ := accesstoken.DeriveSecretName("ns", "sa-creds")

	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName,
			Namespace: "ns",
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"accessToken": []byte("old-token"),
			"expiry":      []byte(expiredExpiry.Format(time.RFC3339)),
		},
	}
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			UID:       types.UID("test-uid"),
		},
		Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("client-id"),
			"clientSecret": []byte("client-secret"),
		},
	}

	newExpiry := time.Now().Add(1 * time.Hour)
	tp := &fakeTokenProvider{token: "new-token", expiry: newExpiry}
	r, k8sClient := newReconciler(t, tp, credSecret, tokenSecret)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.True(t, result.RequeueAfter > 0)
	assert.Equal(t, 1, tp.calls)

	updatedToken := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: tokenSecretName, Namespace: "ns"}, updatedToken))
	assert.Equal(t, "new-token", string(updatedToken.Data["accessToken"]))
	expectedHash, err := accesstoken.CredentialsHash("client-id", "client-secret")
	require.NoError(t, err)
	assert.Equal(t, expectedHash, string(updatedToken.Data["credentialsHash"]),
		"refresh must write the credentials hash alongside the new bearer token")
}

func TestReconcile_SkipsRefreshWhenTokenStillValid(t *testing.T) {
	futureExpiry := time.Now().Add(1 * time.Hour)
	tokenSecretName, _ := accesstoken.DeriveSecretName("ns", "sa-creds")

	matchingHash, err := accesstoken.CredentialsHash("client-id", "client-secret")
	require.NoError(t, err)

	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName,
			Namespace: "ns",
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"accessToken":     []byte("valid-token"),
			"expiry":          []byte(futureExpiry.Format(time.RFC3339)),
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

	tp := &fakeTokenProvider{}
	r, _ := newReconciler(t, tp, credSecret, tokenSecret)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.True(t, result.RequeueAfter > 0)
	assert.Equal(t, 0, tp.calls)
}

func TestReconcile_RefreshesWhenCredentialsChange(t *testing.T) {
	futureExpiry := time.Now().Add(1 * time.Hour)
	tokenSecretName, _ := accesstoken.DeriveSecretName("ns", "sa-creds")

	staleHash, err := accesstoken.CredentialsHash("old-client-id", "old-client-secret")
	require.NoError(t, err)

	// Token Secret was minted from old credentials ("old-client-id"/"old-client-secret").
	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName,
			Namespace: "ns",
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"accessToken":     []byte("stale-token"),
			"expiry":          []byte(futureExpiry.Format(time.RFC3339)),
			"credentialsHash": []byte(staleHash),
		},
	}
	// Credential Secret now holds rotated credentials.
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
		},
		Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("new-client-id"),
			"clientSecret": []byte("new-client-secret"),
		},
	}

	newExpiry := time.Now().Add(1 * time.Hour)
	tp := &fakeTokenProvider{token: "fresh-token", expiry: newExpiry}
	r, k8sClient := newReconciler(t, tp, credSecret, tokenSecret)

	_, err = r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, tp.calls, "unexpired token with stale credentials hash must be refreshed")

	updated := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: tokenSecretName, Namespace: "ns"}, updated))
	assert.Equal(t, "fresh-token", string(updated.Data["accessToken"]))
	expectedHash, err := accesstoken.CredentialsHash("new-client-id", "new-client-secret")
	require.NoError(t, err)
	assert.Equal(t,
		expectedHash,
		string(updated.Data["credentialsHash"]),
		"hash must reflect the new credentials after refresh")
}

func TestReconcile_HandleFetchTokenError(t *testing.T) {
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
	tp := &fakeTokenProvider{err: fmt.Errorf("oauth error")}
	r, _ := newReconciler(t, tp, secret)

	// Returning err (not Result{RequeueAfter}) lets controller-runtime apply
	// its own exponential backoff rather than a constant 10s retry cadence.
	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.Error(t, err)
	assert.Equal(t, ctrl.Result{}, result)
	assert.Equal(t, 1, tp.calls)
}

func TestReconcile_SecretNotFound(t *testing.T) {
	tp := &fakeTokenProvider{}
	r, _ := newReconciler(t, tp)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "nonexistent", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.Equal(t, ctrl.Result{}, result)
	assert.Equal(t, 0, tp.calls)
}

func TestReconcile_RecreatesTokenWhenMissing(t *testing.T) {
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			UID:       types.UID("test-uid"),
		},
		Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("client-id"),
			"clientSecret": []byte("client-secret"),
		},
	}

	newExpiry := time.Now().Add(1 * time.Hour)
	tp := &fakeTokenProvider{token: "fresh-token", expiry: newExpiry}
	r, k8sClient := newReconciler(t, tp, credSecret)

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.NoError(t, err)
	assert.True(t, result.RequeueAfter > 0)
	assert.Equal(t, 1, tp.calls)

	expectedTokenName, _ := accesstoken.DeriveSecretName("ns", "sa-creds")
	tokenSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: expectedTokenName, Namespace: "ns"}, tokenSecret))
	assert.Equal(t, "fresh-token", string(tokenSecret.Data["accessToken"]))
}

func TestReconcile_IsIdempotentOnDuplicateEvent(t *testing.T) {
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			UID:       types.UID("test-uid"),
		},
		Data: map[string][]byte{
			"orgId":        []byte("org-123"),
			"clientId":     []byte("client-id"),
			"clientSecret": []byte("client-secret"),
		},
	}
	expiry := time.Now().Add(1 * time.Hour)
	tp := &fakeTokenProvider{token: "token-v1", expiry: expiry}
	r, k8sClient := newReconciler(t, tp, credSecret)

	// First reconcile creates the token.
	_, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})
	require.NoError(t, err)

	// Second reconcile sees the existing token and does not error.
	_, err = r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})
	require.NoError(t, err, "duplicate reconcile must be idempotent")

	// There must be exactly one token Secret at the derived name.
	expectedTokenName, _ := accesstoken.DeriveSecretName("ns", "sa-creds")
	tokenSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: expectedTokenName, Namespace: "ns"}, tokenSecret))
}
