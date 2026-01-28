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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	fakeID = "fake-id"

	fakeURL = "https://example.com/fake"

	fakeSecret = "fake-secret"

	matchingFakeHash = "5842a3ca8c9cdff1f658a278930fb49d2d5a8c493380f28e39bc7ebbbe21872c"
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
		Name: "webhook-secret",
		Annotations: map[string]string{
			AnnotationContentHash: matchingFakeHash,
		},
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

	unannotatedSampleWebhookSecret := sampleWebhookSecret.DeepCopy()
	unannotatedSampleWebhookSecret.Annotations = nil

	for _, tc := range []struct {
		name           string
		state          state.ResourceState
		provider       atlas.Provider
		serviceBuilder serviceBuilderFunc
		input          *akov2.AtlasThirdPartyIntegration
		objects        []client.Object
		want           ctrlstate.Result
		wantErr        string
	}{
		{
			name:  "initial creates",
			state: state.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Created Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "initial updates",
			state: state.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				return integrationsService
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{unannotatedSampleWebhookSecret},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "initial get fails",
			state: state.StateInitial,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
			want:    ctrlstate.Result{NextState: "Initial"},
			wantErr: "Error getting WEBHOOK Atlas Integration for project testProjectID: unexpected error",
		},

		{
			name:  "created creates",
			state: state.StateCreated,
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
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Created Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "updated updates due unhashed secret",
			state: state.StateUpdated,
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				integrationsService := mocks.NewThirdPartyIntegrationServiceMock(t)
				integrationsService.EXPECT().Get(mock.Anything, "testProjectID", "WEBHOOK").
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				integrationsService.EXPECT().Update(mock.Anything, "testProjectID", mock.Anything).
					Return(&thirdpartyintegration.ThirdPartyIntegration{
						AtlasThirdPartyIntegrationSpec: sampleWebhookIntegration.Spec,
						ID:                             fakeID,
					}, nil)
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{unannotatedSampleWebhookSecret},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for WEBHOOK.",
			},
		},

		{
			name:  "imported does not update",
			state: state.StateInitial,
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
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: ctrlstate.Result{
				NextState: "Created",
				StateMsg:  "Synced WEBHOOK Atlas Third Party Integration for testProjectID.",
			},
		},

		{
			name:  "id diff does not update",
			state: state.StateUpdated,
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
				return integrationsService
			},
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Synced WEBHOOK Atlas Third Party Integration for testProjectID.",
			},
		},

		{
			name:  "different integration updates",
			state: state.StateUpdated,
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
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
					Name: "slack-secret",
					Annotations: map[string]string{
						AnnotationContentHash: "b651a9a05aa41a0555c56f3757d1ab44ae08a7cde4d8fe4477f729e3961cebc6",
					},
				},
				Data: map[string][]byte{
					"apiToken": ([]byte)("fake-token"),
				},
			}},
			want: ctrlstate.Result{
				NextState: "Updated",
				StateMsg:  "Updated Atlas Third Party Integration for SLACK.",
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
				deletionProtection: false,
				serviceBuilder:     tc.serviceBuilder,
			}

			handle := h.HandleInitial
			switch tc.state {
			case state.StateInitial:
				handle = h.HandleInitial
			case state.StateCreated:
				handle = h.HandleCreated
			case state.StateUpdated:
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
		want               ctrlstate.Result
		wantErr            string
	}{
		{
			name:               "deletion deletes",
			deletionProtection: false,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
			want: ctrlstate.Result{
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
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				return mocks.NewThirdPartyIntegrationServiceMock(t)
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: ctrlstate.Result{
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
						SdkClient20250312012: &admin.APIClient{ProjectsApi: mockFindFakeParentProject(t)},
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
			want:    ctrlstate.Result{NextState: "DeletionRequested"},
			wantErr: "Error deleting WEBHOOK Atlas Integration for project testProjectID: unexpected error",
		},

		{
			name:               "deletion deleted when it fails to find project",
			deletionProtection: false,
			provider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetGroupByName(mock.Anything, "fake-project").
						Return(admin.GetGroupByNameApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetGroupByNameExecute(mock.Anything).
						Return(nil, nil, fmt.Errorf("unexpected project fetch error"))
					return &atlas.ClientSet{
						SdkClient20250312012: &admin.APIClient{ProjectsApi: projectAPI},
					}, nil
				},
			},
			serviceBuilder: func(_ *atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService {
				return mocks.NewThirdPartyIntegrationServiceMock(t)
			},
			input:   &sampleWebhookIntegration,
			objects: []client.Object{sampleWebhookSecret},
			want: ctrlstate.Result{
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
