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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

func TestFromAKO(t *testing.T) {
	tests := map[string]struct {
		integrations     []project.Integration
		wantErr          bool
		wantIntegrations []*thirdpartyintegration.ThirdPartyIntegration
	}{
		"empty integrations": {
			integrations:     []project.Integration{},
			wantErr:          false,
			wantIntegrations: []*thirdpartyintegration.ThirdPartyIntegration{},
		},
		"multiple integrations": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
				{
					Type:                     "MICROSOFT_TEAMS",
					MicrosoftTeamsWebhookURL: "https://example.com/webhook",
				},
				{
					Type:          "NEW_RELIC",
					AccountID:     "1234567890",
					LicenseKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
					ReadTokenRef:  common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
					WriteTokenRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
				{
					Type:      "OPS_GENIE",
					Region:    "US",
					APIKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
				{
					Type:          "PAGER_DUTY",
					Region:        "US",
					ServiceKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
				{
					Type:             "PROMETHEUS",
					Enabled:          true,
					ServiceDiscovery: "http",
					UserName:         "user",
					PasswordRef:      common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
				{
					Type:        "SLACK",
					ChannelName: "ako",
					APITokenRef: common.ResourceRefNamespaced{Name: "integration-secret"},
				},
				{
					Type:          "VICTOR_OPS",
					APIKeyRef:     common.ResourceRefNamespaced{Name: "integration-secret"},
					RoutingKeyRef: common.ResourceRefNamespaced{Name: "integration-secret"},
				},
				{
					Type:      "WEBHOOK",
					URL:       "https://example.com/webhook",
					SecretRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
			},
			wantErr: false,
			wantIntegrations: []*thirdpartyintegration.ThirdPartyIntegration{
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "DATADOG",
						Datadog: &akov2.DatadogIntegration{
							Region:                       "EU",
							SendCollectionLatencyMetrics: pointer.MakePtr("disabled"),
							SendDatabaseMetrics:          pointer.MakePtr("disabled"),
						},
					},
					DatadogSecrets: &thirdpartyintegration.DatadogSecrets{
						APIKey: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type:           "MICROSOFT_TEAMS",
						MicrosoftTeams: &akov2.MicrosoftTeamsIntegration{},
					},
					MicrosoftTeamsSecrets: &thirdpartyintegration.MicrosoftTeamsSecrets{
						WebhookUrl: "https://example.com/webhook",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type:     "NEW_RELIC",
						NewRelic: &akov2.NewRelicIntegration{},
					},
					NewRelicSecrets: &thirdpartyintegration.NewRelicSecrets{
						AccountID:  "1234567890",
						LicenseKey: "my-secret-password",
						ReadToken:  "my-secret-password",
						WriteToken: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "OPS_GENIE",
						OpsGenie: &akov2.OpsGenieIntegration{
							Region: "US",
						},
					},
					OpsGenieSecrets: &thirdpartyintegration.OpsGenieSecrets{
						APIKey: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type:      "PAGER_DUTY",
						PagerDuty: &akov2.PagerDutyIntegration{Region: "US"},
					},
					PagerDutySecrets: &thirdpartyintegration.PagerDutySecrets{
						ServiceKey: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "PROMETHEUS",
						Prometheus: &akov2.PrometheusIntegration{
							Enabled:          pointer.MakePtr("enabled"),
							ServiceDiscovery: "http",
						},
					},
					PrometheusSecrets: &thirdpartyintegration.PrometheusSecrets{
						Username: "user",
						Password: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "SLACK",
						Slack: &akov2.SlackIntegration{
							ChannelName: "ako",
						},
					},
					SlackSecrets: &thirdpartyintegration.SlackSecrets{
						APIToken: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "VICTOR_OPS",
						VictorOps: &akov2.VictorOpsIntegration{
							RoutingKey: "my-secret-password",
						},
					},
					VictorOpsSecrets: &thirdpartyintegration.VictorOpsSecrets{
						APIKey: "my-secret-password",
					},
				},
				{
					AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
						Type:    "WEBHOOK",
						Webhook: &akov2.WebhookIntegration{},
					},
					WebhookSecrets: &thirdpartyintegration.WebhookSecrets{
						URL:    "https://example.com/webhook",
						Secret: "my-secret-password",
					},
				},
			},
		},
		"failed to convert": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "invalid-secret", Namespace: "default"},
				},
			},
			wantErr:          true,
			wantIntegrations: []*thirdpartyintegration.ThirdPartyIntegration{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			atlasProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name:         "My Project",
					Integrations: tt.integrations,
				},
			}
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "integration-secret",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"password": []byte("my-secret-password"),
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			reconciler := &AtlasProjectReconciler{
				Client: fake.NewClientBuilder().
					WithScheme(testScheme).
					WithObjects(atlasProject, secret).
					Build(),
			}
			integrations, err := reconciler.fromAKO(context.Background(), atlasProject)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, integrations)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantIntegrations, integrations)
			}
		})
	}
}

func TestEnsureIntegration(t *testing.T) {
	tests := map[string]struct {
		integrations           []project.Integration
		lastAppliedIntegration string
		apiMock                func() admin.ThirdPartyIntegrationsApi
		wantOk                 bool
		wantConditions         []api.Condition
	}{
		"empty integrations, condition unset": {
			integrations:           []project.Integration{},
			lastAppliedIntegration: "{}",
			apiMock: func() admin.ThirdPartyIntegrationsApi {
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListThirdPartyIntegrations(context.Background(), "0123456789").
					Return(admin.ListThirdPartyIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListThirdPartyIntegrationsExecute(mock.AnythingOfType("admin.ListThirdPartyIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				return integrationsApi
			},
			wantOk:         true,
			wantConditions: []api.Condition{},
		},
		"failed to convert integrations": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "invalid-secret", Namespace: "default"},
				},
			},
			lastAppliedIntegration: "{}",
			apiMock: func() admin.ThirdPartyIntegrationsApi {
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)

				return integrationsApi
			},
			wantOk: false,
			wantConditions: []api.Condition{
				api.FalseCondition(api.IntegrationReadyType).
					WithReason(string(workflow.ProjectIntegrationInternal)).
					WithMessageRegexp("failed to convert integrations from AKO: failed to read API key for Datadog integration: secrets \"invalid-secret\" not found"),
			},
		},
		"failed to map lastApplied config": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
			},
			lastAppliedIntegration: "{aaaa",
			apiMock: func() admin.ThirdPartyIntegrationsApi {
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)

				return integrationsApi
			},
			wantOk: false,
			wantConditions: []api.Condition{
				api.FalseCondition(api.IntegrationReadyType).
					WithReason(string(workflow.ProjectIntegrationInternal)).
					WithMessageRegexp("failed to map last applied integrations: error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]: invalid character 'a' looking for beginning of object key string"),
			},
		},
		"failed to reconcile": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
			},
			lastAppliedIntegration: "{}",
			apiMock: func() admin.ThirdPartyIntegrationsApi {
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListThirdPartyIntegrations(context.Background(), "0123456789").
					Return(admin.ListThirdPartyIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListThirdPartyIntegrationsExecute(mock.AnythingOfType("admin.ListThirdPartyIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				integrationsApi.EXPECT().CreateThirdPartyIntegration(context.Background(), "DATADOG", "0123456789", mock.AnythingOfType("*admin.ThirdPartyIntegration")).
					Return(admin.CreateThirdPartyIntegrationApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().CreateThirdPartyIntegrationExecute(mock.AnythingOfType("admin.CreateThirdPartyIntegrationApiRequest")).
					Return(nil, nil, errors.New("failed to create integration"))

				return integrationsApi
			},
			wantOk: false,
			wantConditions: []api.Condition{
				api.FalseCondition(api.IntegrationReadyType).
					WithReason(string(workflow.ProjectIntegrationRequest)).
					WithMessageRegexp("failed to create integration DATADOG: failed to create integration from config: failed to create integration"),
			},
		},
		"successfully reconcile": {
			integrations: []project.Integration{
				{
					Type:      "DATADOG",
					Region:    "EU",
					APIKeyRef: common.ResourceRefNamespaced{Name: "integration-secret", Namespace: "default"},
				},
			},
			lastAppliedIntegration: "{}",
			apiMock: func() admin.ThirdPartyIntegrationsApi {
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListThirdPartyIntegrations(context.Background(), "0123456789").
					Return(admin.ListThirdPartyIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListThirdPartyIntegrationsExecute(mock.AnythingOfType("admin.ListThirdPartyIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				integrationsApi.EXPECT().CreateThirdPartyIntegration(context.Background(), "DATADOG", "0123456789", mock.AnythingOfType("*admin.ThirdPartyIntegration")).
					Return(admin.CreateThirdPartyIntegrationApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().CreateThirdPartyIntegrationExecute(mock.AnythingOfType("admin.CreateThirdPartyIntegrationApiRequest")).
					Return(
						&admin.PaginatedIntegration{
							Results: &[]admin.ThirdPartyIntegration{
								{
									Type: pointer.MakePtr("DATADOG"),
								},
							},
							TotalCount: pointer.MakePtr(1),
						},
						nil,
						nil,
					)

				return integrationsApi
			},
			wantOk: true,
			wantConditions: []api.Condition{
				api.TrueCondition(api.IntegrationReadyType),
			},
		},
		//{
		//	name: "fromAKO returns error, condition set false",
		//	setupProject: func() *akov2.AtlasProject {
		//		return &akov2.AtlasProject{
		//			Spec: akov2.AtlasProjectSpec{
		//				Integrations: []akov2.AtlasThirdPartyIntegration{
		//					{Type: "DATADOG", APIKeyRef: common.ResourceRefNamespaced{Name: "invalid-secret"}},
		//				},
		//			},
		//		}
		//	},
		//	setupReconciler: func(t *testing.T) *AtlasProjectReconciler {
		//		return &AtlasProjectReconciler{
		//			Client: fake.NewClientBuilder().Build(),
		//		}
		//	},
		//	wantOk:        false,
		//	wantCondition: pointer.MakePtr(false),
		//	wantCode:      workflow.ProjectIntegrationInternal,
		//	wantErrMsg:    "failed to convert integrations from AKO",
		//},
		//{
		//	name: "mapLastAppliedProjectIntegrations returns error, condition set false",
		//	setupProject: func() *akov2.AtlasProject {
		//		p := &akov2.AtlasProject{}
		//		// Simulate error by setting invalid last applied annotation
		//		p.Annotations = map[string]string{"mongodb.com/last-applied-configuration": "{"}
		//		return p
		//	},
		//	setupReconciler: func(t *testing.T) *AtlasProjectReconciler {
		//		return &AtlasProjectReconciler{}
		//	},
		//	wantOk:        false,
		//	wantCondition: pointer.MakePtr(false),
		//	wantCode:      workflow.ProjectIntegrationInternal,
		//	wantErrMsg:    "failed to map last applied integrations",
		//},
		//{
		//	name: "reconcile returns error, condition set false",
		//	setupProject: func() *akov2.AtlasProject {
		//		return &akov2.AtlasProject{
		//			Spec: akov2.AtlasProjectSpec{
		//				Integrations: []akov2.AtlasThirdPartyIntegration{{Type: "DATADOG"}},
		//			},
		//		}
		//	},
		//	setupReconciler: func(t *testing.T) *AtlasProjectReconciler {
		//		r := &AtlasProjectReconciler{}
		//		// Patch fromAKO to return valid, but patch NewIntegrationReconciler to return a reconciler that fails
		//		orig := NewIntegrationReconciler
		//		NewIntegrationReconciler = func(ctx *workflow.Context, project *akov2.AtlasProject, integrations []*integration.ThirdPartyIntegration, lastApplied map[string]struct{}) *integrationReconciler {
		//			return &integrationReconciler{
		//				project: project,
		//				lasAppliedIntegrationsTypes: lastApplied,
		//				integrationsInAKO: map[string]*integration.ThirdPartyIntegration{
		//					"DATADOG": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
		//				},
		//				service: &integration.MockThirdPartyIntegrationService{
		//					ListFunc: func(ctx context.Context, projectID string) ([]*integration.ThirdPartyIntegration, error) {
		//						return nil, fmt.Errorf("fail")
		//					},
		//				},
		//			}
		//		}
		//		t.Cleanup(func() { NewIntegrationReconciler = orig })
		//		return r
		//	},
		//	wantOk:        false,
		//	wantCondition: pointer.MakePtr(false),
		//},
		//{
		//	name: "integrations exist, condition set true",
		//	setupProject: func() *akov2.AtlasProject {
		//		return &akov2.AtlasProject{
		//			Spec: akov2.AtlasProjectSpec{
		//				Integrations: []akov2.AtlasThirdPartyIntegration{{Type: "DATADOG"}},
		//			},
		//		}
		//	},
		//	setupReconciler: func(t *testing.T) *AtlasProjectReconciler {
		//		r := &AtlasProjectReconciler{}
		//		// Patch NewIntegrationReconciler to always return OK
		//		orig := NewIntegrationReconciler
		//		NewIntegrationReconciler = func(ctx *workflow.Context, project *akov2.AtlasProject, integrations []*integration.ThirdPartyIntegration, lastApplied map[string]struct{}) *integrationReconciler {
		//			return &integrationReconciler{
		//				project: project,
		//				lasAppliedIntegrationsTypes: lastApplied,
		//				integrationsInAKO: map[string]*integration.ThirdPartyIntegration{
		//					"DATADOG": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
		//				},
		//				service: &integration.MockThirdPartyIntegrationService{
		//					ListFunc: func(ctx context.Context, projectID string) ([]*integration.ThirdPartyIntegration, error) {
		//						return []*integration.ThirdPartyIntegration{}, nil
		//					},
		//				},
		//			}
		//		}
		//		t.Cleanup(func() { NewIntegrationReconciler = orig })
		//		return r
		//	},
		//	wantOk:        true,
		//	wantCondition: pointer.MakePtr(true),
		//},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			atlasProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: tt.lastAppliedIntegration,
					},
				},
				Spec: akov2.AtlasProjectSpec{
					Name:         "My Project",
					Integrations: tt.integrations,
				},
				Status: status.AtlasProjectStatus{
					ID: "0123456789",
				},
			}
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "integration-secret",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"password": []byte("my-secret-password"),
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312006: admin.NewAPIClient(&admin.Configuration{Host: "cloud-qa.mongodb.com"}),
				},
			}
			workflowCtx.SdkClientSet.SdkClient20250312006.ThirdPartyIntegrationsApi = tt.apiMock()
			reconciler := &AtlasProjectReconciler{
				Client: fake.NewClientBuilder().
					WithScheme(testScheme).
					WithObjects(atlasProject, secret).
					WithStatusSubresource(atlasProject).
					Build(),
			}

			result := reconciler.ensureIntegration(workflowCtx, atlasProject)
			assert.Equal(t, tt.wantOk, result.IsOk())
			assert.True(
				t,
				cmp.Equal(
					tt.wantConditions,
					workflowCtx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
			t.Log(
				cmp.Diff(
					tt.wantConditions,
					workflowCtx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestIntegrationReconcile(t *testing.T) {
	tests := map[string]struct {
		lasAppliedIntegrations map[string]struct{}
		integrationsInAKO      map[string]*thirdpartyintegration.ThirdPartyIntegration
		service                func() thirdpartyintegration.ThirdPartyIntegrationService
		wantOk                 bool
		wantPrometheusStatus   bool
	}{
		"handles empty atlas integrations": {
			lasAppliedIntegrations: map[string]struct{}{},
			integrationsInAKO:      map[string]*thirdpartyintegration.ThirdPartyIntegration{},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").Return([]*thirdpartyintegration.ThirdPartyIntegration{}, nil)

				return s
			},
			wantOk: true,
		},
		"fail listing integrations": {
			lasAppliedIntegrations: map[string]struct{}{},
			integrationsInAKO:      map[string]*thirdpartyintegration.ThirdPartyIntegration{},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").Return(nil, errors.New("some error"))

				return s
			},
			wantOk: false,
		},
		"delete owned integrations": {
			lasAppliedIntegrations: map[string]struct{}{"DATADOG": {}},
			integrationsInAKO:      map[string]*thirdpartyintegration.ThirdPartyIntegration{},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return(
						[]*thirdpartyintegration.ThirdPartyIntegration{
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "NEW_RELIC"}},
						},
						nil,
					)
				s.EXPECT().Delete(context.Background(), "0123456789", "DATADOG").Return(nil)

				return s
			},
			wantOk: true,
		},
		"failed to delete integrations": {
			lasAppliedIntegrations: map[string]struct{}{"DATADOG": {}},
			integrationsInAKO:      map[string]*thirdpartyintegration.ThirdPartyIntegration{},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return(
						[]*thirdpartyintegration.ThirdPartyIntegration{
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
						},
						nil,
					)
				s.EXPECT().Delete(context.Background(), "0123456789", "DATADOG").Return(errors.New("error"))

				return s
			},
			wantOk: false,
		},
		"creates new integrations": {
			lasAppliedIntegrations: map[string]struct{}{},
			integrationsInAKO: map[string]*thirdpartyintegration.ThirdPartyIntegration{
				"SLACK": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "SLACK"}},
			},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return([]*thirdpartyintegration.ThirdPartyIntegration{}, nil)
				s.EXPECT().Create(
					context.Background(),
					"0123456789",
					mock.AnythingOfType("*thirdpartyintegration.ThirdPartyIntegration"),
				).Return(&thirdpartyintegration.ThirdPartyIntegration{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "SLACK"}}, nil)

				return s
			},
			wantOk: true,
		},
		"faild creating new integrations": {
			lasAppliedIntegrations: map[string]struct{}{},
			integrationsInAKO: map[string]*thirdpartyintegration.ThirdPartyIntegration{
				"SLACK": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "SLACK"}},
			},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return([]*thirdpartyintegration.ThirdPartyIntegration{}, nil)
				s.EXPECT().Create(
					context.Background(),
					"0123456789",
					mock.AnythingOfType("*thirdpartyintegration.ThirdPartyIntegration"),
				).Return(nil, errors.New("error"))

				return s
			},
			wantOk: false,
		},
		"updates existing integrations": {
			lasAppliedIntegrations: map[string]struct{}{"DATADOG": {}},
			integrationsInAKO: map[string]*thirdpartyintegration.ThirdPartyIntegration{
				"DATADOG": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
			},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return(
						[]*thirdpartyintegration.ThirdPartyIntegration{
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
						},
						nil,
					)
				s.EXPECT().Update(
					context.Background(),
					"0123456789",
					mock.AnythingOfType("*thirdpartyintegration.ThirdPartyIntegration"),
				).Return(&thirdpartyintegration.ThirdPartyIntegration{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}}, nil)

				return s
			},
			wantOk: true,
		},
		"fails updating integrations": {
			lasAppliedIntegrations: map[string]struct{}{"DATADOG": {}},
			integrationsInAKO: map[string]*thirdpartyintegration.ThirdPartyIntegration{
				"DATADOG": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
			},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return(
						[]*thirdpartyintegration.ThirdPartyIntegration{
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "DATADOG"}},
						},
						nil,
					)
				s.EXPECT().Update(
					context.Background(),
					"0123456789",
					mock.AnythingOfType("*thirdpartyintegration.ThirdPartyIntegration"),
				).Return(nil, errors.New("error"))

				return s
			},
			wantOk: false,
		},
		"handles prometheus integration": {
			lasAppliedIntegrations: map[string]struct{}{"PROMETHEUS": {}},
			integrationsInAKO: map[string]*thirdpartyintegration.ThirdPartyIntegration{
				"PROMETHEUS": {AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "PROMETHEUS"}},
			},
			service: func() thirdpartyintegration.ThirdPartyIntegrationService {
				s := translation.NewThirdPartyIntegrationServiceMock(t)
				s.EXPECT().List(context.Background(), "0123456789").
					Return(
						[]*thirdpartyintegration.ThirdPartyIntegration{
							{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "PROMETHEUS"}},
						},
						nil,
					)
				s.EXPECT().Update(
					context.Background(),
					"0123456789",
					mock.AnythingOfType("*thirdpartyintegration.ThirdPartyIntegration"),
				).Return(&thirdpartyintegration.ThirdPartyIntegration{AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{Type: "PROMETHEUS"}}, nil)

				return s
			},
			wantOk:               true,
			wantPrometheusStatus: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			atlasProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "My Project",
				},
				Status: status.AtlasProjectStatus{
					ID: "0123456789",
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312006: admin.NewAPIClient(&admin.Configuration{Host: "cloud-qa.mongodb.com"}),
				},
			}
			reconciler := IntegrationReconciler{
				project:                     atlasProject,
				lasAppliedIntegrationsTypes: tt.lasAppliedIntegrations,
				integrationsInAKO:           tt.integrationsInAKO,
				service:                     tt.service(),
			}
			result := reconciler.reconcile(ctx)
			assert.Equal(t, tt.wantOk, result.IsOk())

			if tt.wantPrometheusStatus {
				v := ctx.StatusOptions()[0].(status.AtlasProjectStatusOption)
				v(&atlasProject.Status)
				assert.Equal(
					t,
					&status.Prometheus{Scheme: "https", DiscoveryURL: "https://cloud-qa.mongodb.com/prometheus/v1.0/groups/0123456789/discovery"},
					atlasProject.Status.Prometheus,
				)
			}
		})
	}
}
