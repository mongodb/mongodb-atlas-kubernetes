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
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
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
		ConnectionSecret: &common.ResourceRef{
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
				SdkClient20250312006: &admin.APIClient{OrganizationsApi: &admin.OrganizationsApiService{}},
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

func TestHandleInitial(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name           string
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		want           ctrlstate.Result
		wantErr        string
	}{
		{
			name:     "successful initial update",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(&atlasorgsettings.AtlasOrgSettings{
						AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
							OrgID:                                  fakeOrgID,
							ConnectionSecret:                       nil,
							ApiAccessListRequired:                  sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
							GenAIFeaturesEnabled:                   sampleAtlasOrgSettings.Spec.GenAIFeaturesEnabled,
							MaxServiceAccountSecretValidityInHours: sampleAtlasOrgSettings.Spec.MaxServiceAccountSecretValidityInHours,
							MultiFactorAuthRequired:                sampleAtlasOrgSettings.Spec.MultiFactorAuthRequired,
							RestrictEmployeeAccess:                 sampleAtlasOrgSettings.Spec.RestrictEmployeeAccess,
							SecurityContact:                        sampleAtlasOrgSettings.Spec.SecurityContact,
							StreamsCrossGroupEnabled:               sampleAtlasOrgSettings.Spec.StreamsCrossGroupEnabled,
						},
					}, nil)
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Initialized.",
			},
		},
		{
			name:     "connection config error",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-org-settings",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: fakeOrgID,
					ConnectionSecret: &common.ResourceRef{
						Name: "non-existent-secret",
					},
				},
			},
			objects: []client.Object{},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "failed to create reconcile context",
		},
		{
			name:     "atlas provider sdk client error",
			provider: createFailingProvider("sdk client error"),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "failed to create reconcile context: sdk client error",
		},
		{
			name:     "org settings service update error",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(nil, fmt.Errorf("atlas api error"))
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "atlas api error",
		},
		{
			name:     "nil response from atlas service",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(nil, nil)
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "atlas returned OrgSettings which is nil after update",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tc.serviceBuilder,
			}

			got, err := h.HandleInitial(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHandleCreated(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name           string
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		want           ctrlstate.Result
		wantErr        string
	}{
		{
			name:     "successful created update",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(&atlasorgsettings.AtlasOrgSettings{
						AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
							OrgID:                                  fakeOrgID,
							ConnectionSecret:                       nil,
							ApiAccessListRequired:                  sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
							GenAIFeaturesEnabled:                   sampleAtlasOrgSettings.Spec.GenAIFeaturesEnabled,
							MaxServiceAccountSecretValidityInHours: sampleAtlasOrgSettings.Spec.MaxServiceAccountSecretValidityInHours,
							MultiFactorAuthRequired:                sampleAtlasOrgSettings.Spec.MultiFactorAuthRequired,
							RestrictEmployeeAccess:                 sampleAtlasOrgSettings.Spec.RestrictEmployeeAccess,
							SecurityContact:                        sampleAtlasOrgSettings.Spec.SecurityContact,
							StreamsCrossGroupEnabled:               sampleAtlasOrgSettings.Spec.StreamsCrossGroupEnabled,
						},
					}, nil)
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Initialized.",
			},
		},
		{
			name:     "org settings service update error in created state",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(nil, fmt.Errorf("update failed"))
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Created"},
			wantErr: "update failed",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tc.serviceBuilder,
			}

			got, err := h.HandleCreated(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHandleUpdated(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name           string
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		want           ctrlstate.Result
		wantErr        string
	}{
		{
			name:     "successful updated state",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(&atlasorgsettings.AtlasOrgSettings{
						AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
							OrgID:                                  fakeOrgID,
							ConnectionSecret:                       nil,
							ApiAccessListRequired:                  sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
							GenAIFeaturesEnabled:                   sampleAtlasOrgSettings.Spec.GenAIFeaturesEnabled,
							MaxServiceAccountSecretValidityInHours: sampleAtlasOrgSettings.Spec.MaxServiceAccountSecretValidityInHours,
							MultiFactorAuthRequired:                sampleAtlasOrgSettings.Spec.MultiFactorAuthRequired,
							RestrictEmployeeAccess:                 sampleAtlasOrgSettings.Spec.RestrictEmployeeAccess,
							SecurityContact:                        sampleAtlasOrgSettings.Spec.SecurityContact,
							StreamsCrossGroupEnabled:               sampleAtlasOrgSettings.Spec.StreamsCrossGroupEnabled,
						},
					}, nil)
				return service
			},
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "updated-org-settings",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: fakeOrgID,
					ConnectionSecret: &common.ResourceRef{
						Name: "atlas-credentials",
					},
					MultiFactorAuthRequired: pointer.MakePtr(false),
					RestrictEmployeeAccess:  pointer.MakePtr(true),
					ApiAccessListRequired:   pointer.MakePtr(false),
				},
			},
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Initialized.",
			},
		},
		{
			name:     "org settings service update error in updated state",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(nil, fmt.Errorf("atlas permissions error"))
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want:    ctrlstate.Result{NextState: "Updated"},
			wantErr: "atlas permissions error",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tc.serviceBuilder,
			}

			got, err := h.HandleUpdated(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHandleDeletionRequested(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name    string
		input   *akov2.AtlasOrgSettings
		objects []client.Object
		want    ctrlstate.Result
		wantErr string
	}{
		{
			name:    "successful deletion unmanage",
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  fmt.Sprintf("unmanaged is AtlasOrgSettings for orgID %s.", fakeOrgID),
			},
		},
		{
			name: "deletion with different org id",
			input: &akov2.AtlasOrgSettings{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "different-org-settings",
					Namespace: "default",
				},
				Spec: akov2.AtlasOrgSettingsSpec{
					OrgID: "different-org-id",
					ConnectionSecret: &common.ResourceRef{
						Name: "atlas-credentials",
					},
				},
			},
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  "unmanaged is AtlasOrgSettings for orgID different-org-id.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    zap.NewNop().Sugar(),
				},
			}

			got, err := h.HandleDeletionRequested(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewReconcileContext(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name           string
		provider       atlas.Provider
		serviceBuilder func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService
		input          *akov2.AtlasOrgSettings
		objects        []client.Object
		wantErr        string
	}{
		{
			name:     "successful reconcile context creation",
			provider: createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
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
					ConnectionSecret: &common.ResourceRef{
						Name: "non-existent-secret",
					},
				},
			},
			objects: []client.Object{},
			wantErr: "secrets \"non-existent-secret\" not found",
		},
		{
			name:     "atlas provider error",
			provider: createFailingProvider("invalid credentials"),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				return mocks.NewAtlasOrgSettingsServiceMock(t)
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			wantErr: "invalid credentials",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tc.serviceBuilder,
			}

			reconcileCtx, err := h.newReconcileContext(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.NotNil(t, reconcileCtx)
				assert.NotNil(t, reconcileCtx.svc)
				assert.Equal(t, tc.input, reconcileCtx.aos)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Nil(t, reconcileCtx)
			}
		})
	}
}

func TestUpsert(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
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
			name:         "successful upsert from initial to created",
			currentState: state.StateInitial,
			nextState:    state.StateCreated,
			provider:     createSuccessfulProvider(),
			serviceBuilder: func(_ *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
				service := mocks.NewAtlasOrgSettingsServiceMock(t)
				service.EXPECT().Update(mock.Anything, fakeOrgID, mock.Anything).
					Return(&atlasorgsettings.AtlasOrgSettings{
						AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
							OrgID:                                  fakeOrgID,
							ConnectionSecret:                       nil,
							ApiAccessListRequired:                  sampleAtlasOrgSettings.Spec.ApiAccessListRequired,
							GenAIFeaturesEnabled:                   sampleAtlasOrgSettings.Spec.GenAIFeaturesEnabled,
							MaxServiceAccountSecretValidityInHours: sampleAtlasOrgSettings.Spec.MaxServiceAccountSecretValidityInHours,
							MultiFactorAuthRequired:                sampleAtlasOrgSettings.Spec.MultiFactorAuthRequired,
							RestrictEmployeeAccess:                 sampleAtlasOrgSettings.Spec.RestrictEmployeeAccess,
							SecurityContact:                        sampleAtlasOrgSettings.Spec.SecurityContact,
							StreamsCrossGroupEnabled:               sampleAtlasOrgSettings.Spec.StreamsCrossGroupEnabled,
						},
					}, nil)
				return service
			},
			input:   &sampleAtlasOrgSettings,
			objects: []client.Object{&fakeAtlasSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Initialized.",
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, tc.input)...).
				WithStatusSubresource(tc.input).Build()

			h := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
					Log:           zap.NewNop().Sugar(),
				},
				serviceBuilder: tc.serviceBuilder,
			}

			got, err := h.upsert(ctx, tc.currentState, tc.nextState, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUnmanage(t *testing.T) {
	h := &AtlasOrgSettingsHandler{}

	for _, tc := range []struct {
		name  string
		orgID string
		want  ctrlstate.Result
	}{
		{
			name:  "unmanage with standard org id",
			orgID: fakeOrgID,
			want: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  fmt.Sprintf("unmanaged is AtlasOrgSettings for orgID %s.", fakeOrgID),
			},
		},
		{
			name:  "unmanage with different org id",
			orgID: "another-org-id",
			want: ctrlstate.Result{
				NextState: "Deleted",
				StateMsg:  "unmanaged is AtlasOrgSettings for orgID another-org-id.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := h.unmanage(tc.orgID)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
