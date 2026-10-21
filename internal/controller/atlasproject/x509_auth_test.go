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

package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

const (
	testProjectID = "project-id"
	testCA        = "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"
)

func userSecurityWithCAs(cas string) *admin.UserSecurity {
	if cas == "" {
		return &admin.UserSecurity{}
	}
	return &admin.UserSecurity{
		CustomerX509: &admin.DBUserTLSX509Settings{Cas: &cas},
	}
}

func projectWithoutCertRef(authModes authmode.AuthModes) *akov2.AtlasProject {
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "project", Namespace: "default"},
	}
	project.Status.ID = testProjectID
	project.Status.AuthModes = authModes
	return project
}

func x509TestContext(t *testing.T, ldapAPI admin.LDAPConfigurationAPI, x509API admin.X509AuthenticationAPI) *workflow.Context {
	t.Helper()

	apiClient := &admin.APIClient{}
	if ldapAPI != nil {
		apiClient.LDAPConfigurationAPI = ldapAPI
	}
	if x509API != nil {
		apiClient.X509AuthenticationAPI = x509API
	}

	return &workflow.Context{
		Context:      context.Background(),
		Log:          zaptest.NewLogger(t).Sugar(),
		SdkClientSet: &atlas.ClientSet{SdkClient20250312: apiClient},
	}
}

func ldapAPIReturningCAs(t *testing.T, cas string) *mockadmin.LDAPConfigurationAPI {
	t.Helper()

	ldapAPI := mockadmin.NewLDAPConfigurationAPI(t)
	ldapAPI.EXPECT().GetUserSecurity(mock.Anything, testProjectID).
		Return(admin.GetUserSecurityApiRequest{ApiService: ldapAPI})
	ldapAPI.EXPECT().GetUserSecurityExecute(mock.Anything).
		Return(userSecurityWithCAs(cas), nil, nil)

	return ldapAPI
}

// TestEnsureX509_DisablesWhenAtlasHasCertAndStatusDoesNot is the regression test
// for the bug: with the certificate reference removed and the status not listing
// X509 (a status that diverged from Atlas), the CA used to stay trusted forever
// because the disable path was only reachable through the status.
func TestEnsureX509_DisablesWhenAtlasHasCertAndStatusDoesNot(t *testing.T) {
	ldapAPI := ldapAPIReturningCAs(t, testCA)

	x509API := mockadmin.NewX509AuthenticationAPI(t)
	x509API.EXPECT().DisableSecurityCustomerX509(mock.Anything, testProjectID).
		Return(admin.DisableSecurityCustomerX509ApiRequest{ApiService: x509API})
	x509API.EXPECT().DisableSecurityCustomerX509Execute(mock.Anything).
		Return(nil, nil, nil)

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram})
	reconciler := &AtlasProjectReconciler{Log: zaptest.NewLogger(t).Sugar()}

	result := reconciler.ensureX509(x509TestContext(t, ldapAPI, x509API), project)

	assert.True(t, result.IsOk())
	assert.False(t, project.Status.AuthModes.CheckAuthMode(authmode.X509))
}

func TestEnsureX509_DisablesWhenStatusListsX509(t *testing.T) {
	ldapAPI := ldapAPIReturningCAs(t, testCA)

	x509API := mockadmin.NewX509AuthenticationAPI(t)
	x509API.EXPECT().DisableSecurityCustomerX509(mock.Anything, testProjectID).
		Return(admin.DisableSecurityCustomerX509ApiRequest{ApiService: x509API})
	x509API.EXPECT().DisableSecurityCustomerX509Execute(mock.Anything).
		Return(nil, nil, nil)

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram, authmode.X509})
	reconciler := &AtlasProjectReconciler{Log: zaptest.NewLogger(t).Sugar()}

	result := reconciler.ensureX509(x509TestContext(t, ldapAPI, x509API), project)

	assert.True(t, result.IsOk())
	assert.False(t, project.Status.AuthModes.CheckAuthMode(authmode.X509))
}

// TestEnsureX509_NoDisableCallWhenAtlasHasNoCert makes sure we do not call the
// disable endpoint for the many projects that never used X.509. The mock has no
// X509AuthenticationApi expectations, so any call would fail the test.
func TestEnsureX509_NoDisableCallWhenAtlasHasNoCert(t *testing.T) {
	ldapAPI := ldapAPIReturningCAs(t, "")

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram})
	reconciler := &AtlasProjectReconciler{Log: zaptest.NewLogger(t).Sugar()}

	result := reconciler.ensureX509(x509TestContext(t, ldapAPI, mockadmin.NewX509AuthenticationAPI(t)), project)

	assert.True(t, result.IsOk())
	assert.False(t, project.Status.AuthModes.CheckAuthMode(authmode.X509))
}

// TestEnsureX509_ClearsStaleStatusEntry covers a status that claims X509 while
// Atlas has nothing configured: the status must converge on Atlas.
func TestEnsureX509_ClearsStaleStatusEntry(t *testing.T) {
	ldapAPI := ldapAPIReturningCAs(t, "")

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram, authmode.X509})
	reconciler := &AtlasProjectReconciler{Log: zaptest.NewLogger(t).Sugar()}

	result := reconciler.ensureX509(x509TestContext(t, ldapAPI, mockadmin.NewX509AuthenticationAPI(t)), project)

	assert.True(t, result.IsOk())
	assert.False(t, project.Status.AuthModes.CheckAuthMode(authmode.X509))
	assert.True(t, project.Status.AuthModes.CheckAuthMode(authmode.Scram))
}

// TestEnsureX509_SkipsUpdateWhenCertAlreadyMatches relies on the cert already
// read from Atlas, so no second GetUserSecurity and no UpdateUserSecurity call.
func TestEnsureX509_SkipsUpdateWhenCertAlreadyMatches(t *testing.T) {
	ldapAPI := ldapAPIReturningCAs(t, testCA)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "x509-cert", Namespace: "default"},
		Data:       map[string][]byte{"ca.crt": []byte(testCA)},
	}

	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram})
	project.Spec.X509CertRef = &common.ResourceRefNamespaced{Name: "x509-cert", Namespace: "default"}

	reconciler := &AtlasProjectReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build(),
		Log:    zaptest.NewLogger(t).Sugar(),
	}

	result := reconciler.ensureX509(x509TestContext(t, ldapAPI, nil), project)

	assert.True(t, result.IsOk())
	assert.True(t, project.Status.AuthModes.CheckAuthMode(authmode.X509))
}

func TestEnsureX509_TerminatesWhenAtlasReadFails(t *testing.T) {
	ldapAPI := mockadmin.NewLDAPConfigurationAPI(t)
	ldapAPI.EXPECT().GetUserSecurity(mock.Anything, testProjectID).
		Return(admin.GetUserSecurityApiRequest{ApiService: ldapAPI})
	ldapAPI.EXPECT().GetUserSecurityExecute(mock.Anything).
		Return(nil, nil, errors.New("atlas is unavailable"))

	project := projectWithoutCertRef(authmode.AuthModes{authmode.Scram})
	reconciler := &AtlasProjectReconciler{Log: zaptest.NewLogger(t).Sugar()}

	ctx := x509TestContext(t, ldapAPI, nil)
	result := reconciler.ensureX509(ctx, project)

	assert.False(t, result.IsOk())
}

// TestReadX509CertFromSecretDoesNotLeakValues is the regression test for the secret
// dump: the error used to embed the whole secret.Data map, and that error travels
// verbatim into the AtlasProject status and the operator log.
func TestReadX509CertFromSecretDoesNotLeakValues(t *testing.T) {
	const privateKeyValue = "super-secret-private-api-key"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "atlas-creds", Namespace: "operator-ns"},
		Data: map[string][]byte{
			"orgId":         []byte("5f4c1b2a3d4e5f6a7b8c9d0e"),
			"publicApiKey":  []byte("abcdefgh"),
			"privateApiKey": []byte(privateKeyValue),
		},
	}

	kubeClient := fake.NewClientBuilder().
		WithScheme(testScheme(t)).
		WithObjects(secret).
		Build()

	_, err := readX509CertFromSecret(
		context.Background(),
		kubeClient,
		client.ObjectKey{Name: "atlas-creds", Namespace: "operator-ns"},
		zaptest.NewLogger(t).Sugar(),
	)

	require.Error(t, err)
	assertNoSecretValues(t, err.Error(), secret.Data)

	// the key names stay, they are what makes the error actionable
	assert.Contains(t, err.Error(), "privateApiKey")
	assert.Contains(t, err.Error(), "3 entries")
}

// TestEnsureX509_ConditionMessageDoesNotLeakSecret covers the full sink: a referenced
// secret that is not a certificate must not end up readable in the project status.
func TestEnsureX509_ConditionMessageDoesNotLeakSecret(t *testing.T) {
	const privateKeyValue = "super-secret-private-api-key"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "atlas-creds", Namespace: "operator-ns"},
		Data: map[string][]byte{
			"orgId":         []byte("5f4c1b2a3d4e5f6a7b8c9d0e"),
			"privateApiKey": []byte(privateKeyValue),
		},
	}

	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "project", Namespace: "tenant-ns"},
		Spec: akov2.AtlasProjectSpec{
			X509CertRef: &common.ResourceRefNamespaced{Name: "atlas-creds", Namespace: "operator-ns"},
		},
	}
	project.Status.ID = testProjectID

	reconciler := &AtlasProjectReconciler{
		Client: fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(secret).Build(),
		Log:    zaptest.NewLogger(t).Sugar(),
	}

	workflowCtx := x509TestContext(t, nil, nil)
	result := reconciler.ensureX509(workflowCtx, project)

	require.False(t, result.IsOk())
	condition := workflowCtx.LastCondition()
	require.NotNil(t, condition)
	assertNoSecretValues(t, condition.Message, secret.Data)
}

func assertNoSecretValues(t *testing.T, message string, data map[string][]byte) {
	t.Helper()

	for key, value := range data {
		assert.NotContains(t, message, string(value),
			"message leaks the plaintext value of %q", key)

		// %v on a map[string][]byte renders values as decimal byte lists
		assert.NotContains(t, message, fmt.Sprintf("%v", value),
			"message leaks the byte representation of %q", key)
	}
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	return scheme
}
