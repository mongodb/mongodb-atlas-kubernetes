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

package atlasorgsettings

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/atlasorgsettings"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

//nolint:gosec
const (
	fakeOrgID     = "fake-org-id"
	fakeAPIKey    = "fake-api-key"
	fakeAPISecret = "fake-api-secret"
)

var fakeAtlasSecret = corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "atlas-credentials",
		Namespace: "default",
	},
	Data: map[string][]byte{
		"orgId":         []byte(fakeOrgID),
		"publicApiKey":  []byte(fakeAPIKey),
		"privateApiKey": []byte(fakeAPISecret),
	},
}

var sampleAtlasOrgSettings = akov2.AtlasOrgSettings{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "sample-org-settings",
		Namespace: "default",
	},
	Spec: akov2.AtlasOrgSettingsSpec{
		OrgID: fakeOrgID,
		ConnectionSecretRef: &api.LocalObjectReference{
			Name: "atlas-credentials",
		},
		ApiAccessListRequired:                  pointer.MakePtr(true),
		GenAIFeaturesEnabled:                   pointer.MakePtr(true),
		MaxServiceAccountSecretValidityInHours: pointer.MakePtr(10),
		MultiFactorAuthRequired:                pointer.MakePtr(true),
		RestrictEmployeeAccess:                 pointer.MakePtr(true),
		SecurityContact:                        pointer.MakePtr("123@mongodb.com"),
		StreamsCrossGroupEnabled:               pointer.MakePtr(true),
	},
	Status: status.AtlasOrgSettingsStatus{},
}

func createSuccessfulProvider() atlas.Provider {
	return &atlasmock.TestProvider{
		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
			return &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{OrganizationsApi: &admin.OrganizationsApiService{}},
			}, nil
		},
	}
}

func createFailingProvider(errorMsg string) atlas.Provider {
	return &atlasmock.TestProvider{
		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
			return nil, errors.New(errorMsg)
		},
	}
}

// Helper functions for creating common service builders
func createServiceBuilder(t *testing.T, getCurrentReturn *atlasorgsettings.AtlasOrgSettings, getCurrentErr error,
	updateReturn *atlasorgsettings.AtlasOrgSettings, updateErr error, expectUpdate bool) func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
	return func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
		service := mocks.NewAtlasOrgSettingsServiceMock(t)
		service.EXPECT().Get(mock.Anything, fakeOrgID).Return(getCurrentReturn, getCurrentErr)

		if expectUpdate {
			service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).Return(updateReturn, updateErr)
		}
		return service
	}
}

func TestNewReconcileContext(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	tests := []struct {
		name           string
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		globalSecret   client.ObjectKey
		wantErr        string
		expectService  bool
	}{
		{
			name:     "successful context creation with connection secret",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:         &sampleAtlasOrgSettings,
			objects:       []client.Object{&fakeAtlasSecret},
			expectService: true,
		},
		{
			name:     "successful context creation with global secret",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-connection-secret",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: fakeOrgID,
					// No ConnectionSecretRef - should use global secret
				},
			},
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-atlas-credentials",
						Namespace: "atlas-system",
					},
					Data: map[string][]byte{
						"orgId":         []byte(fakeOrgID),
						"publicApiKey":  []byte(fakeAPIKey),
						"privateApiKey": []byte(fakeAPISecret),
					},
				},
			},
			globalSecret: client.ObjectKey{
				Name:      "global-atlas-credentials",
				Namespace: "atlas-system",
			},
			expectService: true,
		},
		{
			name:     "context creation with nil connection secret ref",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nil-connection-secret",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID:               fakeOrgID,
					ConnectionSecretRef: nil,
				},
			},
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-atlas-credentials",
						Namespace: "atlas-system",
					},
					Data: map[string][]byte{
						"orgId":         []byte(fakeOrgID),
						"publicApiKey":  []byte(fakeAPIKey),
						"privateApiKey": []byte(fakeAPISecret),
					},
				},
			},
			globalSecret: client.ObjectKey{
				Name:      "global-atlas-credentials",
				Namespace: "atlas-system",
			},
			expectService: true,
		},
		{
			name:     "context creation with empty connection secret name",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "empty-secret-name",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: fakeOrgID,
					ConnectionSecretRef: &api.LocalObjectReference{
						Name: "", // Empty name
					},
				},
			},
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "global-atlas-credentials",
						Namespace: "atlas-system",
					},
					Data: map[string][]byte{
						"orgId":         []byte(fakeOrgID),
						"publicApiKey":  []byte(fakeAPIKey),
						"privateApiKey": []byte(fakeAPISecret),
					},
				},
			},
			globalSecret: client.ObjectKey{
				Name:      "global-atlas-credentials",
				Namespace: "atlas-system",
			},
			expectService: true,
		},
		{
			name:     "missing connection secret",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-secret-org-settings",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: fakeOrgID,
					ConnectionSecretRef: &api.LocalObjectReference{
						Name: "non-existent-secret",
					},
				},
			},
			objects: []client.Object{},
			wantErr: "secrets \"non-existent-secret\" not found",
		},
		{
			name:     "atlas provider sdk client error",
			provider: createFailingProvider("SDK initialization failed"),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			wantErr: "SDK initialization failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tt.objects, tt.input)...).
				WithStatusSubresource(tt.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:          k8sClient,
					AtlasProvider:   tt.provider,
					Log:             zap.NewNop().Sugar(),
					GlobalSecretRef: tt.globalSecret,
				},
				serviceBuilder: tt.serviceBuilder,
			}

			reconcileCtx, err := h.newReconcileRequest(ctx, tt.input)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Nil(t, reconcileCtx)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, reconcileCtx)
				if tt.expectService {
					assert.NotNil(t, reconcileCtx.svc)
				}
				assert.Equal(t, tt.input, reconcileCtx.aos)
			}
		})
	}
}

func TestUpsert(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	tests := []struct {
		name           string
		currentState   state.ResourceState
		nextState      state.ResourceState
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		want           ctrlstate.Result
		wantErr        string
	}{
		{
			name:         "successful upsert with settings different - update needed",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                 fakeOrgID,
						ApiAccessListRequired: pointer.MakePtr(false), // Different from sample
					},
				}, nil,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                 fakeOrgID,
						ApiAccessListRequired: pointer.MakePtr(true),
					},
				}, nil, true),
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Updated.",
			},
		},
		{
			name:         "successful upsert with identical settings - no update needed",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                                  fakeOrgID,
						ApiAccessListRequired:                  sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
						GenAIFeaturesEnabled:                   sampleAtlasOrgSettings.Spec.GenAIFeaturesEnabled,
						MaxServiceAccountSecretValidityInHours: sampleAtlasOrgSettings.Spec.MaxServiceAccountSecretValidityInHours,
						MultiFactorAuthRequired:                sampleAtlasOrgSettings.Spec.MultiFactorAuthRequired,
						RestrictEmployeeAccess:                 sampleAtlasOrgSettings.Spec.RestrictEmployeeAccess,
						SecurityContact:                        sampleAtlasOrgSettings.Spec.SecurityContact,
						StreamsCrossGroupEnabled:               sampleAtlasOrgSettings.Spec.StreamsCrossGroupEnabled,
					},
				}, nil, nil, nil, false), // No update expected
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Ready.",
			},
		},
		{
			name:         "successful upsert with nil current atlas settings",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t, nil, nil,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                 fakeOrgID,
						ApiAccessListRequired: sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
					},
				}, nil, true),
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Updated.",
			},
		},
		{
			name:         "failed reconcile context creation",
			currentState: state.StateCreated,
			nextState:    state.StateUpdated,
			provider:     createFailingProvider("connection error"),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Created"},
			wantErr: "failed to create reconcile context: connection error",
		},
		{
			name:           "get current settings error",
			currentState:   state.StateInitial,
			nextState:      state.StateCreated,
			provider:       createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t, nil, errors.New("get failed"), nil, nil, false),
			input:          &sampleAtlasOrgSettings,
			objects:        []client.Object{&fakeAtlasSecret},
			want:           ctrlstate.Result{NextState: "Initial"},
			wantErr:        "failed to get current org settings from Atlas: get failed",
		},
		{
			name:         "update error after successful get",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                 fakeOrgID,
						ApiAccessListRequired: pointer.MakePtr(false), // Different from sample
					},
				}, nil, nil, errors.New("update failed"), true),
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "update failed",
		},
		{
			name:         "nil response from atlas update service",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: createServiceBuilder(t,
				&atlasorgsettings.AtlasOrgSettings{
					AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
						OrgID:                 fakeOrgID,
						ApiAccessListRequired: pointer.MakePtr(false), // Different from sample
					},
				}, nil, nil, nil, true), // Update returns nil
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "atlas returned OrgSettings which is nil after update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tt.objects, tt.input)...).
				WithStatusSubresource(tt.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tt.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tt.serviceBuilder,
			}

			got, err := h.upsert(ctx, tt.currentState, tt.nextState, tt.input)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnmanage(t *testing.T) {
	h := &AtlasOrgSettingsHandler{}

	tests := []struct {
		name     string
		orgID    string
		expected ctrlstate.Result
	}{
		{
			name:  "unmanage with standard org id",
			orgID: fakeOrgID,
			expected: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  fmt.Sprintf("Unmanaged AtlasOrgSettings for orgID %s.", fakeOrgID),
			},
		},
		{
			name:  "unmanage with different org id",
			orgID: "another-org-id",
			expected: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  "Unmanaged AtlasOrgSettings for orgID another-org-id.",
			},
		},
		{
			name:  "unmanage with empty org id",
			orgID: "",
			expected: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  "Unmanaged AtlasOrgSettings for orgID .",
			},
		},
		{
			name:  "unmanage with special characters in org id",
			orgID: "org-with-special-chars!@#$%",
			expected: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  "Unmanaged AtlasOrgSettings for orgID org-with-special-chars!@#$%.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.unmanage(tt.orgID)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func setupHandlerTest(t *testing.T, scheme *runtime.Scheme) (*AtlasOrgSettingsHandler, context.Context) {
	provider := createSuccessfulProvider()
	serviceBuilder := createServiceBuilder(t,
		&atlasorgsettings.AtlasOrgSettings{
			AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
				OrgID:                 fakeOrgID,
				ApiAccessListRequired: pointer.MakePtr(false), // Different from sample
			},
		}, nil,
		&atlasorgsettings.AtlasOrgSettings{
			AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
				OrgID:                 fakeOrgID,
				ApiAccessListRequired: pointer.MakePtr(true),
			},
		}, nil, true)

	k8sClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&fakeAtlasSecret, &sampleAtlasOrgSettings).
		WithStatusSubresource(&sampleAtlasOrgSettings).Build()

	h := &AtlasOrgSettingsHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:        k8sClient,
			AtlasProvider: provider,
			Log:           zap.NewNop().Sugar(),
		},
		serviceBuilder: serviceBuilder,
	}

	return h, context.Background()
}

func TestHandlerMethods(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	tests := []struct {
		name           string
		handlerFunc    func(*AtlasOrgSettingsHandler, context.Context, *akov2.AtlasOrgSettings) (ctrlstate.Result, error)
		expectedResult ctrlstate.Result
	}{
		{
			name:        "HandleInitial",
			handlerFunc: (*AtlasOrgSettingsHandler).HandleInitial,
			expectedResult: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Updated AtlasOrgSettings.",
			},
		},
		{
			name:        "HandleUpdated",
			handlerFunc: (*AtlasOrgSettingsHandler).HandleUpdated,
			expectedResult: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Updated.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ctx := setupHandlerTest(t, scheme)
			got, err := tt.handlerFunc(h, ctx, &sampleAtlasOrgSettings)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, got)
		})
	}

	// Test HandleDeletionRequested
	t.Run("HandleDeletionRequested", func(t *testing.T) {
		h := &AtlasOrgSettingsHandler{}
		ctx := context.Background()

		got, err := h.HandleDeletionRequested(ctx, &sampleAtlasOrgSettings)
		require.NoError(t, err)
		assert.Equal(t, ctrlstate.Result{
			NextState: "Deleted",
			StateMsg:  fmt.Sprintf("Unmanaged AtlasOrgSettings for orgID %s.", fakeOrgID),
		}, got)
	})
}

func TestEqualMethodBehaviorInUpsert(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	// Test case where Equal method returns false due to nil current settings
	t.Run("Equal with nil current settings - update needed", func(t *testing.T) {
		provider := createSuccessfulProvider()
		serviceBuilder := createServiceBuilder(t, nil, nil, // nil current settings
			&atlasorgsettings.AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                 fakeOrgID,
					ApiAccessListRequired: pointer.MakePtr(true),
				},
			}, nil, true)

		k8sClient := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(&fakeAtlasSecret, &sampleAtlasOrgSettings).
			WithStatusSubresource(&sampleAtlasOrgSettings).Build()

		h := &AtlasOrgSettingsHandler{
			AtlasReconciler: reconciler.AtlasReconciler{
				Client:        k8sClient,
				AtlasProvider: provider,
				Log:           zap.NewNop().Sugar(),
			},
			serviceBuilder: serviceBuilder,
		}

		got, err := h.upsert(ctx, state.StateInitial, state.StateCreated, &sampleAtlasOrgSettings)
		require.NoError(t, err)
		assert.Equal(t, ctrlstate.Result{
			NextState: "Created",
			StateMsg:  "Updated.",
		}, got)
	})

	// Test case with different state transitions
	t.Run("Updated to Updated state transition", func(t *testing.T) {
		provider := createSuccessfulProvider()
		serviceBuilder := createServiceBuilder(t,
			&atlasorgsettings.AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                 fakeOrgID,
					ApiAccessListRequired: pointer.MakePtr(false), // Different from sample
				},
			}, nil,
			&atlasorgsettings.AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                 fakeOrgID,
					ApiAccessListRequired: pointer.MakePtr(true),
				},
			}, nil, true)

		k8sClient := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(&fakeAtlasSecret, &sampleAtlasOrgSettings).
			WithStatusSubresource(&sampleAtlasOrgSettings).Build()

		h := &AtlasOrgSettingsHandler{
			AtlasReconciler: reconciler.AtlasReconciler{
				Client:        k8sClient,
				AtlasProvider: provider,
				Log:           zap.NewNop().Sugar(),
			},
			serviceBuilder: serviceBuilder,
		}

		got, err := h.upsert(ctx, state.StateUpdated, state.StateUpdated, &sampleAtlasOrgSettings)
		require.NoError(t, err)
		assert.Equal(t, ctrlstate.Result{
			NextState: "Updated",
			StateMsg:  "Updated.",
		}, got)
	})
}
