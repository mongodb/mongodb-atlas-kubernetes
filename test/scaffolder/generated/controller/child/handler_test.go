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

package child

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/crd2go/constate"
	"github.com/crd2go/crapi"
	"github.com/crd2go/crapi/refs"
	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312026/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

// mockProvider implements atlas.Provider for testing.
type mockProvider struct {
	clientSet *atlas.ClientSet
	err       error
}

func (m *mockProvider) SdkClientSet(_ context.Context, _ *atlas.Credentials, _ *zap.SugaredLogger) (*atlas.ClientSet, error) {
	return m.clientSet, m.err
}

func (m *mockProvider) IsCloudGov() bool                                   { return false }
func (m *mockProvider) IsResourceSupported(_ api.AtlasCustomResource) bool { return true }

var _ atlas.Provider = (*mockProvider)(nil)

// mockTranslator implements crapi.Translator for handler dispatch tests.
type mockTranslator struct{}

func (m *mockTranslator) Scheme() *runtime.Scheme            { return nil }
func (m *mockTranslator) MajorVersion() string               { return "v20250312" }
func (m *mockTranslator) Mappings() ([]*refs.Mapping, error) { return nil, nil }
func (m *mockTranslator) ToAPI(_ any, _ client.Object, _ ...client.Object) error {
	return nil
}
func (m *mockTranslator) FromAPI(_ client.Object, _ any, _ ...client.Object) ([]client.Object, error) {
	return nil, nil
}

func newSchemeWithCoreV1(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, v1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
}

func newCredentialSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"orgId":         []byte("test-org"),
			"publicApiKey":  []byte("test-public-key"),
			"privateApiKey": []byte("test-private-key"),
		},
	}
}

// buildHandler creates a Handler with the given configuration for testing the dispatch layer.
func buildHandler(
	kubeClient client.Client,
	provider atlas.Provider,
	globalSecretRef client.ObjectKey,
	translators map[string]crapi.Translator,
) *Handler {
	logger := zap.NewNop()
	return &Handler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   provider,
			Client:          kubeClient,
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Sugar(),
		},
		deletionProtection: false,
		translators:        translators,
		handlerv20250312:   handlerv20250312Func,
	}
}

// TestGetHandlerForResource tests the version dispatch logic in handler.go.
func TestGetHandlerForResource(t *testing.T) {
	ctx := context.Background()
	scheme := newSchemeWithCoreV1(t)

	globalSecret := newCredentialSecret("global-secret")
	globalSecretRef := client.ObjectKey{Name: "global-secret", Namespace: "default"}

	tests := []struct {
		name        string
		child       *v1.Child
		translators map[string]crapi.Translator
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "selects v20250312 handler when V20250312 spec is set",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{
						ParentRef: &k8s.LocalReference{Name: "my-parent"},
					},
				},
			},
			translators: map[string]crapi.Translator{
				"v20250312": &mockTranslator{},
			},
			wantErr: false,
		},
		{
			name: "returns error when no spec version is set",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec:       v1.ChildSpec{},
			},
			translators: map[string]crapi.Translator{
				"v20250312": &mockTranslator{},
			},
			wantErr:    true,
			wantErrMsg: "no resource spec version specified",
		},
		{
			name: "returns error when translator not found for version",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec: v1.ChildSpec{
					V20250312: &v1.V20250312{},
				},
			},
			translators: map[string]crapi.Translator{},
			wantErr:     true,
			wantErrMsg:  "unsupported version v20250312",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(globalSecret).
				Build()

			provider := &mockProvider{
				clientSet: &atlas.ClientSet{
					SdkClient20250312: &v20250312sdk.APIClient{},
				},
			}

			handler := buildHandler(fakeClient, provider, globalSecretRef, tc.translators)
			result, err := handler.getHandlerForResource(ctx, tc.child)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

// TestGetSDKClientSet tests credential resolution from connection secrets.
func TestGetSDKClientSet(t *testing.T) {
	ctx := context.Background()
	scheme := newSchemeWithCoreV1(t)

	globalSecret := newCredentialSecret("global-secret")
	globalSecretRef := client.ObjectKey{Name: "global-secret", Namespace: "default"}

	perResourceSecret := newCredentialSecret("resource-secret")

	expectedClientSet := &atlas.ClientSet{
		SdkClient20250312: &v20250312sdk.APIClient{},
	}

	tests := []struct {
		name       string
		child      *v1.Child
		objects    []client.Object
		provider   *mockProvider
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "uses global secret when no connectionSecretRef",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec:       v1.ChildSpec{},
			},
			objects:  []client.Object{globalSecret},
			provider: &mockProvider{clientSet: expectedClientSet},
			wantErr:  false,
		},
		{
			name: "uses per-resource secret when connectionSecretRef is set",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec: v1.ChildSpec{
					ConnectionSecretRef: &k8s.LocalReference{Name: "resource-secret"},
				},
			},
			objects:  []client.Object{globalSecret, perResourceSecret},
			provider: &mockProvider{clientSet: expectedClientSet},
			wantErr:  false,
		},
		{
			name: "returns error when secret not found",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec: v1.ChildSpec{
					ConnectionSecretRef: &k8s.LocalReference{Name: "nonexistent-secret"},
				},
			},
			objects:    []client.Object{globalSecret},
			provider:   &mockProvider{clientSet: expectedClientSet},
			wantErr:    true,
			wantErrMsg: "failed to resolve Atlas credentials",
		},
		{
			name: "returns error when provider fails",
			child: &v1.Child{
				ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
				Spec:       v1.ChildSpec{},
			},
			objects: []client.Object{globalSecret},
			provider: &mockProvider{
				err: assert.AnError,
			},
			wantErr:    true,
			wantErrMsg: "failed to setup Atlas SDK client",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).
				Build()

			handler := buildHandler(fakeClient, tc.provider, globalSecretRef, nil)
			clientSet, err := handler.getSDKClientSet(ctx, tc.child)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, clientSet)
		})
	}
}

// TestHandlerStateTransitions tests that the Handler dispatch layer correctly
// delegates to the version-specific handler for each constate.
func TestHandlerStateTransitions(t *testing.T) {
	ctx := context.Background()
	scheme := newSchemeWithCoreV1(t)

	globalSecret := newCredentialSecret("global-secret")
	globalSecretRef := client.ObjectKey{Name: "global-secret", Namespace: "default"}

	child := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{},
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(globalSecret).
		Build()

	provider := &mockProvider{
		clientSet: &atlas.ClientSet{
			SdkClient20250312: &v20250312sdk.APIClient{},
		},
	}

	translators := map[string]crapi.Translator{
		"v20250312": &mockTranslator{},
	}

	handler := buildHandler(fakeClient, provider, globalSecretRef, translators)

	type stateFunc func(context.Context, *v1.Child) (constate.Result, error)

	stateTests := []struct {
		name      string
		fn        stateFunc
		wantState constate.ResourceState
	}{
		{"HandleInitial", handler.HandleInitial, constate.StateUpdated},
		{"HandleImportRequested", handler.HandleImportRequested, constate.StateImported},
		{"HandleImported", handler.HandleImported, constate.StateUpdated},
		{"HandleCreating", handler.HandleCreating, constate.StateCreated},
		{"HandleCreated", handler.HandleCreated, constate.StateUpdated},
		{"HandleUpdating", handler.HandleUpdating, constate.StateUpdated},
		{"HandleUpdated", handler.HandleUpdated, constate.StateUpdated},
		{"HandleDeletionRequested", handler.HandleDeletionRequested, constate.StateDeleting},
		{"HandleDeleting", handler.HandleDeleting, constate.StateDeleted},
	}

	for _, tc := range stateTests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.fn(ctx, child)
			require.NoError(t, err)
			assert.Equal(t, tc.wantState, result.NextState)
		})
	}
}

// TestHandlerStateTransitions_NoVersion tests that the Handler dispatch returns
// an error for each state when no version is set.
func TestHandlerStateTransitions_NoVersion(t *testing.T) {
	ctx := context.Background()
	scheme := newSchemeWithCoreV1(t)

	globalSecret := newCredentialSecret("global-secret")
	globalSecretRef := client.ObjectKey{Name: "global-secret", Namespace: "default"}

	child := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
		Spec:       v1.ChildSpec{},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(globalSecret).
		Build()

	provider := &mockProvider{
		clientSet: &atlas.ClientSet{
			SdkClient20250312: &v20250312sdk.APIClient{},
		},
	}

	translators := map[string]crapi.Translator{
		"v20250312": &mockTranslator{},
	}

	handler := buildHandler(fakeClient, provider, globalSecretRef, translators)

	type stateFunc func(context.Context, *v1.Child) (constate.Result, error)

	stateTests := []struct {
		name         string
		fn           stateFunc
		wantErrState constate.ResourceState
	}{
		{"HandleInitial", handler.HandleInitial, constate.StateInitial},
		{"HandleImportRequested", handler.HandleImportRequested, constate.StateImportRequested},
		{"HandleImported", handler.HandleImported, constate.StateImported},
		{"HandleCreating", handler.HandleCreating, constate.StateCreating},
		{"HandleCreated", handler.HandleCreated, constate.StateCreated},
		{"HandleUpdating", handler.HandleUpdating, constate.StateUpdating},
		{"HandleUpdated", handler.HandleUpdated, constate.StateUpdated},
		{"HandleDeletionRequested", handler.HandleDeletionRequested, constate.StateDeletionRequested},
		{"HandleDeleting", handler.HandleDeleting, constate.StateDeleting},
	}

	for _, tc := range stateTests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.fn(ctx, child)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "no resource spec version specified")
			assert.Equal(t, tc.wantErrState, result.NextState)
		})
	}
}

// TestGetDependents_Child verifies that the child handler's getDependents
// returns an empty slice (child has no dependents).
func TestGetDependents_Child(t *testing.T) {
	ctx := context.Background()

	child := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{
				ParentRef: &k8s.LocalReference{Name: "my-parent"},
			},
		},
	}

	handler := NewHandlerv20250312(nil, nil, nil, false)
	dependents, err := handler.getDependents(ctx, child)
	require.NoError(t, err)
	assert.Empty(t, dependents)
}

// loadTestCRD parses the test CRD YAML and returns the CRD matching the given kind.
func loadTestCRD(t *testing.T, kind string) *apiextensionsv1.CustomResourceDefinition {
	t.Helper()
	data, err := os.ReadFile("../../../testdata/crds.yaml")
	require.NoError(t, err, "failed to read testdata/crds.yaml")

	parsed, err := crds.ParseCRDs(bufio.NewScanner(bytes.NewBuffer(data)))
	require.NoError(t, err, "failed to parse CRDs from testdata")

	for _, crd := range parsed {
		if crd.Spec.Names.Kind == kind {
			return crd
		}
	}
	t.Fatalf("CRD %q not found in testdata/crds.yaml", kind)
	return nil
}

// TestHandlerWithRealTranslator validates the full constructor wiring path:
// CRD parsing -> translator creation -> handler dispatch -> state transition.
// This exercises the same code path as NewChildReconciler but using test CRDs
// instead of production embedded CRDs.
func TestHandlerWithRealTranslator(t *testing.T) {
	ctx := context.Background()
	scheme := newSchemeWithCoreV1(t)

	crd := loadTestCRD(t, "Child")
	translators, err := crapi.NewPerVersionTranslators(scheme, crd, crdVersion, sdkVersions...)
	require.NoError(t, err, "NewPerVersionTranslators should succeed for Child CRD")
	require.Contains(t, translators, "v20250312", "translators should contain v20250312")

	globalSecret := newCredentialSecret("global-secret")
	globalSecretRef := client.ObjectKey{Name: "global-secret", Namespace: "default"}

	parent := &v1.Parent{
		ObjectMeta: metav1.ObjectMeta{Name: "my-parent", Namespace: "default"},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(globalSecret, parent).
		Build()

	provider := &mockProvider{
		clientSet: &atlas.ClientSet{
			SdkClient20250312: &v20250312sdk.APIClient{},
		},
	}

	handler := buildHandler(fakeClient, provider, globalSecretRef, translators)

	child := &v1.Child{
		ObjectMeta: metav1.ObjectMeta{Name: "test-child", Namespace: "default"},
		Spec: v1.ChildSpec{
			V20250312: &v1.V20250312{
				ParentRef: &k8s.LocalReference{Name: "my-parent"},
			},
		},
	}

	// Verify the full dispatch chain works with real translators
	result, err := handler.HandleInitial(ctx, child)
	require.NoError(t, err)
	assert.Equal(t, constate.StateUpdated, result.NextState)

	// Verify version dispatch selects the correct handler
	versionHandler, err := handler.getHandlerForResource(ctx, child)
	require.NoError(t, err)
	assert.NotNil(t, versionHandler)
}
