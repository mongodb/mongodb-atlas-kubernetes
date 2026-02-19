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

package serviceaccount_test

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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/serviceaccount"
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

func newReconciler(t *testing.T, tp serviceaccount.TokenProvider, objs ...client.Object) (*serviceaccount.ServiceAccountReconciler, client.Client) {
	t.Helper()
	scheme := newScheme(t)
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	r := &serviceaccount.ServiceAccountReconciler{
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

	updatedSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: "sa-creds", Namespace: "ns"}, updatedSecret))
	tokenName := updatedSecret.Annotations[reconciler.AccessTokenAnnotation]
	assert.NotEmpty(t, tokenName)

	tokenSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: tokenName, Namespace: "ns"}, tokenSecret))
	assert.Equal(t, "access-token-value", string(tokenSecret.Data["accessToken"]))
	assert.NotEmpty(t, tokenSecret.Data["expiry"])

	assert.Equal(t, secretservice.CredLabelVal, tokenSecret.Labels[secretservice.TypeLabelKey])

	require.Len(t, tokenSecret.OwnerReferences, 1)
	assert.Equal(t, "sa-creds", tokenSecret.OwnerReferences[0].Name)
	assert.Equal(t, types.UID("test-uid"), tokenSecret.OwnerReferences[0].UID)
}

func TestReconcile_RefreshesExpiredToken(t *testing.T) {
	expiredExpiry := time.Now().Add(-10 * time.Minute)
	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds-token-xyz",
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
			Annotations: map[string]string{
				reconciler.AccessTokenAnnotation: "sa-creds-token-xyz",
			},
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
		types.NamespacedName{Name: "sa-creds-token-xyz", Namespace: "ns"}, updatedToken))
	assert.Equal(t, "new-token", string(updatedToken.Data["accessToken"]))
}

func TestReconcile_SkipsRefreshWhenTokenStillValid(t *testing.T) {
	futureExpiry := time.Now().Add(1 * time.Hour)
	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds-token-xyz",
			Namespace: "ns",
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"accessToken": []byte("valid-token"),
			"expiry":      []byte(futureExpiry.Format(time.RFC3339)),
		},
	}
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			Annotations: map[string]string{
				reconciler.AccessTokenAnnotation: "sa-creds-token-xyz",
			},
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

func TestReconcile_HandlesFetchTokenError(t *testing.T) {
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

	result, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Name: "sa-creds", Namespace: "ns"},
	})

	require.Error(t, err)
	assert.True(t, result.RequeueAfter > 0)
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

func TestReconcile_CreatesNewTokenWhenAnnotatedSecretMissing(t *testing.T) {
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sa-creds",
			Namespace: "ns",
			UID:       types.UID("test-uid"),
			Annotations: map[string]string{
				reconciler.AccessTokenAnnotation: "deleted-token-secret",
			},
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

	updatedSecret := &corev1.Secret{}
	require.NoError(t, k8sClient.Get(context.Background(),
		types.NamespacedName{Name: "sa-creds", Namespace: "ns"}, updatedSecret))
	tokenName := updatedSecret.Annotations[reconciler.AccessTokenAnnotation]
	assert.NotEqual(t, "deleted-token-secret", tokenName)
}
