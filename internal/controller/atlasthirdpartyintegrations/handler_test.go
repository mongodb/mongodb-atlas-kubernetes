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

package integrations

import (
	"context"
	"fmt"
	"testing"

	"github.com/crd2go/constate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

const (
	fakeID = "fake-id"

	fakeURL = "https://example.com/fake"

	fakeSecret = "fake-secret"
)

var sampleWebhookIntegration = akov2.AtlasThirdPartyIntegration{
	ObjectMeta: metav1.ObjectMeta{Name: "webhook-test"},
	Spec: akov2.AtlasThirdPartyIntegrationSpec{
		ProjectDualReference: referenceFakeProject,
		Type:                 "WEBHOOK",
		Webhook: &akov2.WebhookIntegration{
			URLSecretRef: api.LocalObjectReference{
				Name: "webhook-secret",
			},
		},
	},
	Status: status.AtlasThirdPartyIntegrationStatus{
		ID: fakeID,
	},
}

var sampleWebhookSecret = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:            "webhook-secret",
		ResourceVersion: "secret-rv",
	},
	Data: map[string][]byte{
		"url":    ([]byte)(fakeURL),
		"secret": ([]byte)(fakeSecret),
	},
}

func TestHandleUpsert(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name             string
		state            constate.ResourceState
		provider         atlas.Provider
		serviceBuilder   serviceBuilderFunc
		input            *akov2.AtlasThirdPartyIntegration
		objects          []client.Object
		reapplyTimestamp string
		want             constate.Result
		wantErr          string
	}{
		{
			name:  "initial creates",
			state: constate.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(nil, thirdpartyintegration.ErrNotFound)
				integrationsService.EXPECT().Create(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Created",
				StateMsg:  "Created Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "initial updates",
			state: constate.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{ID: fakeID}, nil)
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "initial get fails",
			state: constate.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(nil, fmt.Errorf("unexpected error"))
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want:    constate.Result{NextState: "Initial"},
			wantErr: "Error getting WEBHOOK Atlas Integration for project testProjectID: unexpected error",
		},

		{
			name:  "created creates",
			state: constate.StateCreated,
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(nil, thirdpartyintegration.ErrNotFound)
				integrationsService.EXPECT().Create(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Created",
				StateMsg:  "Created Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:             "updated updates when integration matches",
			state:            constate.StateUpdated,
			reapplyTimestamp: "2025-01-02T03:04:05Z",
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{ID: fakeID}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "imported updates",
			state: constate.StateInitial,
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
						WebhookSecrets: &thirdpartyintegration.WebhookSecrets{
							URL:    fakeURL,
							Secret: fakeSecret,
						},
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{ID: fakeID}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "id diff updates",
			state: constate.StateUpdated,
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             "other-fake-id",
						WebhookSecrets: &thirdpartyintegration.WebhookSecrets{
							URL:    fakeURL,
							Secret: fakeSecret,
						},
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{ID: fakeID}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "different integration updates",
			state: constate.StateUpdated,
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "SLACK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
							ProjectDualReference: referenceFakeProject,
							Type:                 "SLACK",
							Slack: &akov2.SlackIntegration{
								ChannelName: "other-channel",
								TeamName:    "other-team",
							},
						},
						ID: fakeID,
						SlackSecrets: &thirdpartyintegration.SlackSecrets{
							APIToken: "fake-token",
						},
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
							ProjectDualReference: referenceFakeProject,
							Type:                 "SLACK",
							Slack: &akov2.SlackIntegration{
								ChannelName: "spec-channel",
								TeamName:    "spec-team",
							},
						},
						ID: fakeID,
					}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: metav1.ObjectMeta{Name: "slack-integration"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					ProjectDualReference: referenceFakeProject,
					Type:                 "SLACK",
					Slack: &akov2.SlackIntegration{
						APITokenSecretRef: api.LocalObjectReference{
							Name: "slack-secret",
						},
						ChannelName: "spec-channel",
						TeamName:    "spec-team",
					},
				},
			},
			objects: []client.Object{&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "slack-secret",
					ResourceVersion: "secret-rv",
				},
				Data: map[string][]byte{
					"apiToken": ([]byte)("fake-token"),
				},
			}},
			want: constate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for SLACK.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.input = tc.input.DeepCopy()
			if tc.reapplyTimestamp != "" {
				tc.input.Annotations = map[string]string{constate.AnnotationReapplyTimestamp: tc.reapplyTimestamp}
			}
			tc.input.Generation = 7
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, &fakeAtlasSecret, &fakeProject, tc.input)...).
				WithStatusSubresource(tc.input).Build()
			h := AtlasThirdPartyIntegrationHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
				},
				deletionProtection: false,
				serviceBuilder:     tc.serviceBuilder,
			}

			handle := h.HandleInitial
			switch tc.state {
			case constate.StateInitial:
				handle = h.HandleInitial
			case constate.StateCreated:
				handle = h.HandleCreated
			case constate.StateUpdated:
				handle = h.HandleUpdated
			default:
				panic(fmt.Errorf("unsupported state %v for test", tc.state))
			}
			got, err := handle(ctx, tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)

			persisted := &akov2.AtlasThirdPartyIntegration{}
			require.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(tc.input), persisted))
			if err == nil && (got.NextState == constate.StateCreated || got.NextState == constate.StateUpdated) {
				assert.Equal(t, sampleWebhookSecret.ResourceVersion, persisted.Annotations[customresource.ReferencedSecretResourceVersion])
				assert.Equal(t, tc.input.Annotations[constate.AnnotationReapplyTimestamp], persisted.Annotations[annotationLastReconciledReapplyTimestamp])
				if got.NextState == constate.StateUpdated {
					assert.Equal(t, tc.input.Generation, persisted.Status.ObservedGeneration)
				}
			} else if err != nil {
				assert.NotContains(t, persisted.Annotations, customresource.ReferencedSecretResourceVersion)
				assert.NotContains(t, persisted.Annotations, annotationLastReconciledReapplyTimestamp)
			}
		})
	}
}

func TestUpsertSkipsOnlyWhenSecretVersionAndGenerationMatch(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name               string
		generation         int64
		observedGeneration int64
		annotation         string
		reapplyTimestamp   string
		handledTimestamp   string
		wantSkip           bool
	}{
		{name: "matching version and observed generation", generation: 3, observedGeneration: 3, annotation: "secret-rv", reapplyTimestamp: "1", handledTimestamp: "1", wantSkip: true},
		{name: "stale observed generation", generation: 3, observedGeneration: 2, annotation: "secret-rv"},
		{name: "changed annotation", generation: 3, observedGeneration: 3, annotation: "old-rv"},
		{name: "missing annotation", generation: 3, observedGeneration: 3},
		{name: "changed reapply timestamp", generation: 3, observedGeneration: 3, annotation: "secret-rv", reapplyTimestamp: "2", handledTimestamp: "1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			integration := sampleWebhookIntegration.DeepCopy()
			integration.Generation = tc.generation
			integration.Status.ObservedGeneration = tc.observedGeneration
			integration.Annotations = nil
			if tc.annotation != "" || tc.reapplyTimestamp != "" || tc.handledTimestamp != "" {
				integration.Annotations = map[string]string{
					customresource.ReferencedSecretResourceVersion: tc.annotation,
					constate.AnnotationReapplyTimestamp:            tc.reapplyTimestamp,
					annotationLastReconciledReapplyTimestamp:       tc.handledTimestamp,
				}
			}
			secret := sampleWebhookSecret.DeepCopy()
			secret.ResourceVersion = "secret-rv"
			k8sClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(integration, secret).Build()
			h := AtlasThirdPartyIntegrationHandler{AtlasReconciler: reconciler.AtlasReconciler{Client: k8sClient}}

			result, err := h.upsert(ctx, constate.StateInitial, integration)
			if tc.wantSkip {
				require.NoError(t, err)
				assert.Contains(t, result.StateMsg, "Spec and referenced Secret unchanged")
				return
			}
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "resource version is unchanged")
		})
	}
}

func TestHandleDeletion(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	ctx := context.Background()

	for _, tc := range []struct {
		name               string
		deletionProtection bool
		provider           atlas.Provider
		serviceBuilder     serviceBuilderFunc
		input              *akov2.AtlasThirdPartyIntegration
		objects            []client.Object
		want               constate.Result
		wantErr            string
	}{
		{
			name:               "deletion deletes",
			deletionProtection: false,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Delete(mock.Anything, "testProjectID", "WEBHOOK").
					Return(nil)
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Deleted",
				StateMsg:  "Deleted Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:               "deletion with protection unmanaged",
			deletionProtection: true,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				return mocks.NewThirdPartyIntegrationServiceMock(t)
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Deleted",
				StateMsg:  "Deleted Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:               "deletion fails",
			deletionProtection: false,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Delete(mock.Anything, "testProjectID", "WEBHOOK").
					Return(fmt.Errorf("unexpected error"))
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want:    constate.Result{NextState: "DeletionRequested"},
			wantErr: "Error deleting WEBHOOK Atlas Integration for project testProjectID: unexpected error",
		},

		{
			name:               "deletion deleted when it fails to find project",
			deletionProtection: false,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsAPI(t)
					projectAPI.EXPECT().GetGroupByName(mock.Anything, "fake-project").
						Return(admin.GetGroupByNameApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetGroupByNameExecute(mock.Anything).
						Return(nil, nil, fmt.Errorf("unexpected project fetch error"))
					return &atlas.ClientSet{
						SdkClient20250312: &admin.APIClient{ProjectsAPI: projectAPI},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				return mocks.NewThirdPartyIntegrationServiceMock(t)
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: constate.Result{
				NextState: "Deleted",
				StateMsg:  "Deleted Atlas Third Party Integration for WEBHOOK.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(append(tc.objects, &fakeAtlasSecret, &fakeProject, tc.input)...).
				WithStatusSubresource(tc.input).Build()
			h := AtlasThirdPartyIntegrationHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client:        k8sClient,
					AtlasProvider: tc.provider,
				},
				deletionProtection: tc.deletionProtection,
				serviceBuilder:     tc.serviceBuilder,
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
