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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestEnsureAlertConfigurations(t *testing.T) {
	tests := []struct {
		name               string
		project            *akov2.AtlasProject
		expectedConditions []api.Condition
		expectOKResult     bool
		expectSecretErr    bool
		expectSyncErr      bool
		setupSecrets       func() []runtime.Object
	}{
		{
			name: "Alert configuration sync disabled",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					AlertConfigurationSyncEnabled: false,
				},
			},
			expectedConditions: []api.Condition{},
			expectOKResult:     true,
		},
		{
			name: "Empty alert configurations with sync enabled",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					AlertConfigurationSyncEnabled: true,
					AlertConfigurations:           []akov2.AlertConfiguration{},
				},
			},
			expectedConditions: []api.Condition{},
			expectOKResult:     true,
		},
		{
			name: "Alert configuration with missing secret should fail",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					AlertConfigurationSyncEnabled: true,
					AlertConfigurations: []akov2.AlertConfiguration{
						{
							EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
							Enabled:       true,
							Notifications: []akov2.Notification{
								{
									TypeName:    "SLACK",
									APITokenRef: common.ResourceRefNamespaced{Name: "missing-token"},
									ChannelName: "alerts",
								},
							},
						},
					},
				},
			},
			setupSecrets: func() []runtime.Object {
				return []runtime.Object{}
			},
			expectOKResult:  false,
			expectSecretErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))
			require.NoError(t, akov2.AddToScheme(scheme))

			var clientObjs []runtime.Object
			if tt.setupSecrets != nil {
				clientObjs = tt.setupSecrets()
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(clientObjs...).
				Build()

			reconciler := &AtlasProjectReconciler{
				Client: fakeClient,
			}

			logger := zaptest.NewLogger(t).Sugar()
			ctx := context.Background()

			workflowCtx := &workflow.Context{
				Context: ctx,
				Log:     logger,
			}

			result := reconciler.ensureAlertConfigurations(workflowCtx, tt.project)

			if tt.expectOKResult {
				assert.True(t, result.IsOk())
			} else {
				assert.False(t, result.IsOk())
			}
		})
	}
}

func TestReadAlertConfigurationsSecretsData(t *testing.T) {
	tests := []struct {
		name         string
		project      *akov2.AtlasProject
		alertConfigs []akov2.AlertConfiguration
		secrets      []runtime.Object
		expectError  bool
		errorMsg     string
	}{
		{
			name: "Successfully read API token secret",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
			},
			alertConfigs: []akov2.AlertConfiguration{
				{
					Notifications: []akov2.Notification{
						{
							APITokenRef: common.ResourceRefNamespaced{Name: "api-token"},
						},
					},
				},
			},
			secrets: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "api-token",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"APIToken": []byte("test-api-token"),
					},
				},
			},
			expectError: false,
		},
		{
			name: "Secret not found",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
			},
			alertConfigs: []akov2.AlertConfiguration{
				{
					Notifications: []akov2.Notification{
						{
							APITokenRef: common.ResourceRefNamespaced{Name: "missing-secret"},
						},
					},
				},
			},
			secrets:     []runtime.Object{},
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name: "Successfully read multiple secret types",
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "default",
				},
			},
			alertConfigs: []akov2.AlertConfiguration{
				{
					Notifications: []akov2.Notification{
						{
							DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "datadog-secret"},
						},
						{
							FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "flowdock-secret"},
						},
						{
							OpsGenieAPIKeyRef: common.ResourceRefNamespaced{Name: "opsgenie-secret"},
						},
						{
							ServiceKeyRef: common.ResourceRefNamespaced{Name: "service-secret"},
						},
						{
							VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "victorops-secret"},
						},
					},
				},
			},
			secrets: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "datadog-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"DatadogAPIKey": []byte("dd-key"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "flowdock-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"FlowdockAPIToken": []byte("flow-token"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "opsgenie-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"OpsGenieAPIKey": []byte("ops-key"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"ServiceKey": []byte("service-key"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "victorops-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"VictorOpsAPIKey":     []byte("victor-api-key"),
						"VictorOpsRoutingKey": []byte("victor-routing-key"),
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(tt.secrets...).
				Build()

			reconciler := &AtlasProjectReconciler{
				Client: fakeClient,
			}

			logger := zaptest.NewLogger(t).Sugar()
			ctx := context.Background()

			workflowCtx := &workflow.Context{
				Context: ctx,
				Log:     logger,
			}

			err := reconciler.readAlertConfigurationsSecretsData(tt.project, workflowCtx, tt.alertConfigs)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadNotificationSecret(t *testing.T) {
	tests := []struct {
		name            string
		secretRef       common.ResourceRefNamespaced
		parentNamespace string
		fieldName       string
		secret          *corev1.Secret
		expectedValue   string
		expectError     bool
		errorMsg        string
	}{
		{
			name: "Successfully read secret from same namespace",
			secretRef: common.ResourceRefNamespaced{
				Name: "test-secret",
			},
			parentNamespace: "default",
			fieldName:       "APIToken",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"APIToken": []byte("test-token-value"),
				},
			},
			expectedValue: "test-token-value",
			expectError:   false,
		},
		{
			name: "Successfully read secret from different namespace",
			secretRef: common.ResourceRefNamespaced{
				Name:      "test-secret",
				Namespace: "other-namespace",
			},
			parentNamespace: "default",
			fieldName:       "APIToken",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "other-namespace",
				},
				Data: map[string][]byte{
					"APIToken": []byte("cross-ns-token"),
				},
			},
			expectedValue: "cross-ns-token",
			expectError:   false,
		},
		{
			name: "Secret field does not exist",
			secretRef: common.ResourceRefNamespaced{
				Name: "test-secret",
			},
			parentNamespace: "default",
			fieldName:       "MissingField",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"APIToken": []byte("test-token-value"),
				},
			},
			expectError: true,
			errorMsg:    "doesn't contain 'MissingField' parameter",
		},
		{
			name: "Secret field is empty",
			secretRef: common.ResourceRefNamespaced{
				Name: "test-secret",
			},
			parentNamespace: "default",
			fieldName:       "EmptyField",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"EmptyField": []byte(""),
				},
			},
			expectError: true,
			errorMsg:    "contains an empty value for 'EmptyField' parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, corev1.AddToScheme(scheme))

			var objects []runtime.Object
			if tt.secret != nil {
				objects = append(objects, tt.secret)
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(objects...).
				Build()

			ctx := context.Background()

			value, err := readNotificationSecret(ctx, fakeClient, tt.secretRef, tt.parentNamespace, tt.fieldName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

func TestSyncAlertConfigurations(t *testing.T) {
	tests := []struct {
		name                 string
		groupID              string
		alertSpecs           []akov2.AlertConfiguration
		existingAlertConfigs []admin.GroupAlertsConfig
		mockAlertConfigsAPI  func() *mockadmin.AlertConfigurationsApi
		expectOKResult       bool
		expectedCreateCount  int
		expectedDeleteCount  int
	}{
		{
			name:                 "No alert configurations to sync",
			groupID:              "test-group-id",
			alertSpecs:           []akov2.AlertConfiguration{},
			existingAlertConfigs: []admin.GroupAlertsConfig{},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)
				apiMock.EXPECT().ListAlertConfigurations(mock.Anything, "test-group-id").
					Return(admin.ListAlertConfigurationsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListAlertConfigurationsExecute(mock.Anything).
					Return(&admin.PaginatedAlertConfig{
						Results: &[]admin.GroupAlertsConfig{},
					}, &http.Response{StatusCode: 200}, nil)
				return apiMock
			},
			expectOKResult:      true,
			expectedCreateCount: 0,
			expectedDeleteCount: 0,
		},
		{
			name:    "Create new alert configuration",
			groupID: "test-group-id",
			alertSpecs: []akov2.AlertConfiguration{
				{
					EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
					Enabled:       true,
					Notifications: []akov2.Notification{
						{
							TypeName:     "EMAIL",
							EmailAddress: "test@example.com",
						},
					},
				},
			},
			existingAlertConfigs: []admin.GroupAlertsConfig{},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)
				apiMock.EXPECT().ListAlertConfigurations(mock.Anything, "test-group-id").
					Return(admin.ListAlertConfigurationsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListAlertConfigurationsExecute(mock.Anything).
					Return(&admin.PaginatedAlertConfig{
						Results: &[]admin.GroupAlertsConfig{},
					}, &http.Response{StatusCode: 200}, nil)

				createdConfig := admin.GroupAlertsConfig{
					Id:            pointer.MakePtr("new-alert-id"),
					EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
					Enabled:       pointer.MakePtr(true),
				}

				apiMock.EXPECT().CreateAlertConfiguration(mock.Anything, "test-group-id", mock.Anything).
					Return(admin.CreateAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateAlertConfigurationExecute(mock.Anything).
					Return(&createdConfig, &http.Response{StatusCode: 201}, nil)

				return apiMock
			},
			expectOKResult:      true,
			expectedCreateCount: 1,
			expectedDeleteCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			ctx := context.Background()

			mockAPIClient := &admin.APIClient{}
			if tt.mockAlertConfigsAPI != nil {
				mockAPIClient.AlertConfigurationsApi = tt.mockAlertConfigsAPI()
			}

			atlasClientSet := &atlas.ClientSet{
				SdkClient20250312006: mockAPIClient,
			}

			workflowCtx := &workflow.Context{
				Context:      ctx,
				Log:          logger,
				SdkClientSet: atlasClientSet,
			}

			result := syncAlertConfigurations(workflowCtx, tt.groupID, tt.alertSpecs)

			if tt.expectOKResult {
				assert.True(t, result.IsOk())
			} else {
				assert.False(t, result.IsOk())
			}
		})
	}
}

func TestCheckAlertConfigurationStatuses(t *testing.T) {
	tests := []struct {
		name           string
		statuses       []status.AlertConfiguration
		expectOKResult bool
	}{
		{
			name:           "Empty statuses should return OK",
			statuses:       []status.AlertConfiguration{},
			expectOKResult: true,
		},
		{
			name: "All successful statuses should return OK",
			statuses: []status.AlertConfiguration{
				{
					ID:           "config-1",
					ErrorMessage: "",
				},
				{
					ID:           "config-2",
					ErrorMessage: "",
				},
			},
			expectOKResult: true,
		},
		{
			name: "Status with error should return failure",
			statuses: []status.AlertConfiguration{
				{
					ID:           "config-1",
					ErrorMessage: "",
				},
				{
					ID:           "config-2",
					ErrorMessage: "Failed to create alert configuration",
				},
			},
			expectOKResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkAlertConfigurationStatuses(tt.statuses)

			if tt.expectOKResult {
				assert.True(t, result.IsOk())
			} else {
				assert.False(t, result.IsOk())
			}
		})
	}
}

func TestDeleteAlertConfigs(t *testing.T) {
	tests := []struct {
		name                string
		groupID             string
		alertConfigIDs      []string
		mockAlertConfigsAPI func() *mockadmin.AlertConfigurationsApi
		expectError         bool
	}{
		{
			name:           "No alert configurations to delete",
			groupID:        "test-group-id",
			alertConfigIDs: []string{},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				return mockadmin.NewAlertConfigurationsApi(t)
			},
			expectError: false,
		},
		{
			name:           "Successfully delete alert configurations",
			groupID:        "test-group-id",
			alertConfigIDs: []string{"config-1", "config-2"},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)
				apiMock.EXPECT().DeleteAlertConfiguration(mock.Anything, "test-group-id", "config-1").
					Return(admin.DeleteAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().DeleteAlertConfigurationExecute(mock.Anything).
					Return(&http.Response{StatusCode: 204}, nil)
				apiMock.EXPECT().DeleteAlertConfiguration(mock.Anything, "test-group-id", "config-2").
					Return(admin.DeleteAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().DeleteAlertConfigurationExecute(mock.Anything).
					Return(&http.Response{StatusCode: 204}, nil)
				return apiMock
			},
			expectError: false,
		},
		{
			name:           "Error deleting alert configuration",
			groupID:        "test-group-id",
			alertConfigIDs: []string{"config-1"},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)
				apiMock.EXPECT().DeleteAlertConfiguration(mock.Anything, "test-group-id", "config-1").
					Return(admin.DeleteAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().DeleteAlertConfigurationExecute(mock.Anything).
					Return(nil, errors.New("API error"))
				return apiMock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			ctx := context.Background()

			mockAPIClient := &admin.APIClient{
				AlertConfigurationsApi: tt.mockAlertConfigsAPI(),
			}

			atlasClientSet := &atlas.ClientSet{
				SdkClient20250312006: mockAPIClient,
			}

			workflowCtx := &workflow.Context{
				Context:      ctx,
				Log:          logger,
				SdkClientSet: atlasClientSet,
			}

			err := deleteAlertConfigs(workflowCtx, tt.groupID, tt.alertConfigIDs)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateAlertConfigs(t *testing.T) {
	tests := []struct {
		name                string
		groupID             string
		alertSpecs          []akov2.AlertConfiguration
		mockAlertConfigsAPI func() *mockadmin.AlertConfigurationsApi
		expectedStatusCount int
	}{
		{
			name:       "No alert configurations to create",
			groupID:    "test-group-id",
			alertSpecs: []akov2.AlertConfiguration{},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				return mockadmin.NewAlertConfigurationsApi(t)
			},
			expectedStatusCount: 0,
		},
		{
			name:    "Successfully create alert configuration",
			groupID: "test-group-id",
			alertSpecs: []akov2.AlertConfiguration{
				{
					EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
					Enabled:       true,
					Notifications: []akov2.Notification{
						{
							TypeName:     "EMAIL",
							EmailAddress: "test@example.com",
						},
					},
				},
			},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)

				createdConfig := admin.GroupAlertsConfig{
					Id:            pointer.MakePtr("new-alert-id"),
					EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
					Enabled:       pointer.MakePtr(true),
				}

				apiMock.EXPECT().CreateAlertConfiguration(mock.Anything, "test-group-id", mock.Anything).
					Return(admin.CreateAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateAlertConfigurationExecute(mock.Anything).
					Return(&createdConfig, &http.Response{StatusCode: 201}, nil)

				return apiMock
			},
			expectedStatusCount: 1,
		},
		{
			name:    "Error creating alert configuration",
			groupID: "test-group-id",
			alertSpecs: []akov2.AlertConfiguration{
				{
					EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
					Enabled:       true,
				},
			},
			mockAlertConfigsAPI: func() *mockadmin.AlertConfigurationsApi {
				apiMock := mockadmin.NewAlertConfigurationsApi(t)

				apiMock.EXPECT().CreateAlertConfiguration(mock.Anything, "test-group-id", mock.Anything).
					Return(admin.CreateAlertConfigurationApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateAlertConfigurationExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 400}, errors.New("API error"))

				return apiMock
			},
			expectedStatusCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			ctx := context.Background()

			mockAPIClient := &admin.APIClient{
				AlertConfigurationsApi: tt.mockAlertConfigsAPI(),
			}

			atlasClientSet := &atlas.ClientSet{
				SdkClient20250312006: mockAPIClient,
			}

			workflowCtx := &workflow.Context{
				Context:      ctx,
				Log:          logger,
				SdkClientSet: atlasClientSet,
			}

			statuses := createAlertConfigs(workflowCtx, tt.groupID, tt.alertSpecs)

			assert.Len(t, statuses, tt.expectedStatusCount)
		})
	}
}

func TestSortAlertConfigs(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name                      string
		alertConfigSpecs          []akov2.AlertConfiguration
		atlasAlertConfigs         []admin.GroupAlertsConfig
		expectedCreateCount       int
		expectedDeleteCount       int
		expectedCreateStatusCount int
	}{
		{
			name:                      "Empty specs and atlas configs",
			alertConfigSpecs:          []akov2.AlertConfiguration{},
			atlasAlertConfigs:         []admin.GroupAlertsConfig{},
			expectedCreateCount:       0,
			expectedDeleteCount:       0,
			expectedCreateStatusCount: 0,
		},
		{
			name: "New alert config to create",
			alertConfigSpecs: []akov2.AlertConfiguration{
				{
					EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
					Enabled:       true,
				},
			},
			atlasAlertConfigs:         []admin.GroupAlertsConfig{},
			expectedCreateCount:       1,
			expectedDeleteCount:       0,
			expectedCreateStatusCount: 0,
		},
		{
			name:             "Alert config to delete",
			alertConfigSpecs: []akov2.AlertConfiguration{},
			atlasAlertConfigs: []admin.GroupAlertsConfig{
				{
					Id:            pointer.MakePtr("config-to-delete"),
					EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
					Enabled:       pointer.MakePtr(true),
				},
			},
			expectedCreateCount:       0,
			expectedDeleteCount:       1,
			expectedCreateStatusCount: 0,
		},
		{
			name: "Matching alert configs",
			alertConfigSpecs: []akov2.AlertConfiguration{
				{
					EventTypeName: "HOST_DOWN",
					Enabled:       true,
					Notifications: []akov2.Notification{},
					Matchers:      []akov2.Matcher{},
				},
			},
			atlasAlertConfigs: []admin.GroupAlertsConfig{
				{
					Id:            pointer.MakePtr("existing-config"),
					EventTypeName: pointer.MakePtr("HOST_DOWN"),
					Enabled:       pointer.MakePtr(true),
					Notifications: &[]admin.AlertsNotificationRootForGroup{},
					Matchers:      &[]admin.StreamsMatcher{},
				},
			},
			expectedCreateCount:       0,
			expectedDeleteCount:       0,
			expectedCreateStatusCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := sortAlertConfigs(logger, tt.alertConfigSpecs, tt.atlasAlertConfigs)

			assert.Len(t, diff.Create, tt.expectedCreateCount)
			assert.Len(t, diff.Delete, tt.expectedDeleteCount)
			assert.Len(t, diff.CreateStatus, tt.expectedCreateStatusCount)
		})
	}
}

func TestIsAlertConfigSpecEqualToAtlas(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name             string
		alertConfigSpec  akov2.AlertConfiguration
		atlasAlertConfig admin.GroupAlertsConfig
		expectedEqual    bool
	}{
		{
			name: "Different event type names",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("HOST_DOWN"),
			},
			expectedEqual: false,
		},
		{
			name: "Different enabled status",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(false),
			},
			expectedEqual: false,
		},
		{
			name: "Atlas enabled is nil",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       nil,
			},
			expectedEqual: false,
		},
		{
			name: "Different severity override",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName:    "OUTSIDE_METRIC_THRESHOLD",
				Enabled:          true,
				SeverityOverride: "CRITICAL",
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName:    pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:          pointer.MakePtr(true),
				SeverityOverride: pointer.MakePtr("WARNING"),
			},
			expectedEqual: false,
		},
		{
			name: "Different threshold",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				Threshold: &akov2.Threshold{
					Operator:  "GREATER_THAN",
					Threshold: "80",
					Units:     "PERCENT",
				},
				Notifications: []akov2.Notification{},
				Matchers:      []akov2.Matcher{},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{},
				Matchers:      &[]admin.StreamsMatcher{},
			},
			expectedEqual: false,
		},
		{
			name: "Different metric threshold",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				MetricThreshold: &akov2.MetricThreshold{
					MetricName: "CPU_USER",
					Operator:   "GREATER_THAN",
					Threshold:  "80",
					Units:      "PERCENT",
					Mode:       "AVERAGE",
				},
				Notifications: []akov2.Notification{},
				Matchers:      []akov2.Matcher{},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{},
				Matchers:      &[]admin.StreamsMatcher{},
			},
			expectedEqual: false,
		},
		{
			name: "Different notification count",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				Notifications: []akov2.Notification{
					{
						TypeName:     "EMAIL",
						EmailAddress: "test@example.com",
					},
				},
				Matchers: []akov2.Matcher{},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{},
				Matchers:      &[]admin.StreamsMatcher{},
			},
			expectedEqual: false,
		},
		{
			name: "Different matcher count",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				Notifications: []akov2.Notification{},
				Matchers: []akov2.Matcher{
					{
						FieldName: "HOSTNAME",
						Operator:  "EQUALS",
						Value:     "test-host",
					},
				},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{},
				Matchers:      &[]admin.StreamsMatcher{},
			},
			expectedEqual: false,
		},
		{
			name: "Notification same count but different content",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				Notifications: []akov2.Notification{
					{
						TypeName:     "EMAIL",
						EmailAddress: "test1@example.com",
					},
				},
				Matchers: []akov2.Matcher{},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{
					{
						TypeName:     pointer.MakePtr("EMAIL"),
						EmailAddress: pointer.MakePtr("test2@example.com"),
					},
				},
				Matchers: &[]admin.StreamsMatcher{},
			},
			expectedEqual: false,
		},
		{
			name: "Matcher same count but different content",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName: "OUTSIDE_METRIC_THRESHOLD",
				Enabled:       true,
				Notifications: []akov2.Notification{},
				Matchers: []akov2.Matcher{
					{
						FieldName: "HOSTNAME",
						Operator:  "EQUALS",
						Value:     "test-host-1",
					},
				},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName: pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:       pointer.MakePtr(true),
				Notifications: &[]admin.AlertsNotificationRootForGroup{},
				Matchers: &[]admin.StreamsMatcher{
					{
						FieldName: "HOSTNAME",
						Operator:  "EQUALS",
						Value:     "test-host-2",
					},
				},
			},
			expectedEqual: false,
		},
		{
			name: "Matching simple configuration",
			alertConfigSpec: akov2.AlertConfiguration{
				EventTypeName:    "OUTSIDE_METRIC_THRESHOLD",
				Enabled:          true,
				SeverityOverride: "CRITICAL",
				Notifications:    []akov2.Notification{},
				Matchers:         []akov2.Matcher{},
			},
			atlasAlertConfig: admin.GroupAlertsConfig{
				EventTypeName:    pointer.MakePtr("OUTSIDE_METRIC_THRESHOLD"),
				Enabled:          pointer.MakePtr(true),
				SeverityOverride: pointer.MakePtr("CRITICAL"),
				Notifications:    &[]admin.AlertsNotificationRootForGroup{},
				Matchers:         &[]admin.StreamsMatcher{},
			},
			expectedEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAlertConfigSpecEqualToAtlas(logger, tt.alertConfigSpec, tt.atlasAlertConfig)
			assert.Equal(t, tt.expectedEqual, result)
		})
	}
}
