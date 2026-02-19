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

package thirdpartyintegration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	integration "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	testProjectID = "fake-project"

	testID = "fake-id"

	testID2 = "fake-id-2"

	testRegion = "fake-region"

	testIntegrationType = "PAGER_DUTY"

	testServiceKey = "fake-service-key"

	testAccount = "fake-account-id"

	testLicenseKey = "fake-license-key"

	testReadToken = "fake-read-token"

	testWriteToken = "fake-write-token"
)

var (
	ErrFakeFailure = errors.New("fake failure")
)

func TestIntegrationsCreate(t *testing.T) {
	testAPIKey := utils.RandomName("fake-apy-key")
	for _, tc := range []struct {
		title         string
		integration   *integration.ThirdPartyIntegration
		api           admin.ThirdPartyIntegrationsApi
		expected      *integration.ThirdPartyIntegration
		expectedError error
	}{
		{
			title: "successful api create",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			api: testCreateIntegrationAPI(
				[]admin.ThirdPartyIntegration{
					{
						Id:                           pointer.MakePtr(string(testID)),
						Type:                         pointer.MakePtr("DATADOG"),
						ApiKey:                       pointer.MakePtr(testAPIKey),
						Region:                       pointer.MakePtr(string(testRegion)),
						SendCollectionLatencyMetrics: pointer.MakePtr(true),
					},
				},
				nil,
			),
			expected: &integration.ThirdPartyIntegration{
				ID: testID,
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			expectedError: nil,
		},

		{
			title: "API failure gets passed through",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			api: testCreateIntegrationAPI(
				nil,
				ErrFakeFailure,
			),
			expected:      nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse config returns before calling API",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "BLAH",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			expected:      nil,
			expectedError: integration.ErrUnsupportedIntegrationType,
		},

		{
			title: "failure to parse API reply",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			api: testCreateIntegrationAPI(
				[]admin.ThirdPartyIntegration{
					{
						Id:                           pointer.MakePtr(string(testID)),
						Type:                         pointer.MakePtr("BLAH"),
						ApiKey:                       pointer.MakePtr(testAPIKey),
						Region:                       pointer.MakePtr(string(testRegion)),
						SendCollectionLatencyMetrics: pointer.MakePtr(true),
					},
				},
				nil,
			),
			expected:      nil,
			expectedError: integration.ErrUnsupportedIntegrationType,
		},

		{
			title: "failure to extract matching type from API reply",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			api: testCreateIntegrationAPI(
				[]admin.ThirdPartyIntegration{
					{
						Id:          pointer.MakePtr(string(testID)),
						Type:        pointer.MakePtr("SLACK"),
						ApiToken:    pointer.MakePtr("fake-token"),
						ChannelName: pointer.MakePtr("channel"),
						TeamName:    pointer.MakePtr("team"),
					},
					{
						Id:   pointer.MakePtr(string(testID2)),
						Type: pointer.MakePtr("WEBHOOK"),
						Url:  pointer.MakePtr("http://example.com/fake"),
					},
				},
				nil,
			),
			expected:      nil,
			expectedError: integration.ErrNotFound,
		},

		{
			title: "extracts matching type from API reply",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
						SendDatabaseMetrics:          pointer.MakePtr("disabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			api: testCreateIntegrationAPI(
				[]admin.ThirdPartyIntegration{
					{
						Id:          pointer.MakePtr(string(testID)),
						Type:        pointer.MakePtr("SLACK"),
						ApiToken:    pointer.MakePtr("fake-token"),
						ChannelName: pointer.MakePtr("channel"),
						TeamName:    pointer.MakePtr("team"),
					},
					{
						Id:                           pointer.MakePtr(string(testID2)),
						Type:                         pointer.MakePtr("DATADOG"),
						SendCollectionLatencyMetrics: pointer.MakePtr(false),
						SendDatabaseMetrics:          pointer.MakePtr(true),
						ApiKey:                       pointer.MakePtr(testAPIKey),
					},
				},
				nil,
			),
			expected: &integration.ThirdPartyIntegration{
				ID: testID2,
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						SendCollectionLatencyMetrics: pointer.MakePtr("disabled"),
						SendDatabaseMetrics:          pointer.MakePtr("enabled"),
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: testAPIKey,
				},
			},
			expectedError: nil,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := integration.NewThirdPartyIntegrationService(tc.api)
			newIntegrations, err := s.Create(ctx, testProjectID, tc.integration)
			assert.Equal(t, tc.expected, newIntegrations)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestIntegrationsgGet(t *testing.T) {
	for _, tc := range []struct {
		title         string
		api           admin.ThirdPartyIntegrationsApi
		expected      *integration.ThirdPartyIntegration
		expectedError error
	}{
		{
			title: "successful api get returns success",
			api: testGetIntegrationAPI(
				&admin.ThirdPartyIntegration{
					Id:                           pointer.MakePtr(string(testID)),
					Type:                         pointer.MakePtr(string(testIntegrationType)),
					ServiceKey:                   pointer.MakePtr(testServiceKey),
					Region:                       pointer.MakePtr(string(testRegion)),
					SendCollectionLatencyMetrics: pointer.MakePtr(true),
				},
				nil,
			),
			expected: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: testIntegrationType,
					PagerDuty: &akov2.PagerDutyIntegration{
						Region: testRegion,
					},
				},
				ID: testID,
				PagerDutySecrets: &integration.PagerDutySecrets{
					ServiceKey: testServiceKey,
				},
			},
			expectedError: nil,
		},

		{
			title: "generic API failure passes though",
			api: testGetIntegrationAPI(
				nil,
				ErrFakeFailure,
			),
			expected:      nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse API reply",
			api: testGetIntegrationAPI(
				&admin.ThirdPartyIntegration{
					Id:                           pointer.MakePtr(string(testID)),
					Type:                         pointer.MakePtr("BLAH"),
					ServiceKey:                   pointer.MakePtr(testServiceKey),
					Region:                       pointer.MakePtr(string(testRegion)),
					SendCollectionLatencyMetrics: pointer.MakePtr(true),
				},
				nil,
			),
			expected:      nil,
			expectedError: integration.ErrUnsupportedIntegrationType,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := integration.NewThirdPartyIntegrationService(tc.api)
			integrations, err := s.Get(ctx, testProjectID, testIntegrationType)
			assert.Equal(t, tc.expected, integrations)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestIntegrationsUpdate(t *testing.T) {
	for _, tc := range []struct {
		title         string
		integration   *integration.ThirdPartyIntegration
		api           admin.ThirdPartyIntegrationsApi
		expected      *integration.ThirdPartyIntegration
		expectedError error
	}{
		{
			title: "successful api update returns success",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type:     "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{},
				},
				NewRelicSecrets: &integration.NewRelicSecrets{
					AccountID:  testAccount,
					LicenseKey: testLicenseKey,
					ReadToken:  testReadToken,
					WriteToken: testWriteToken,
				},
			},
			api: testUpdateIntegrationAPI(
				&admin.ThirdPartyIntegration{
					Id:         pointer.MakePtr(string(testID)),
					Type:       pointer.MakePtr("NEW_RELIC"),
					AccountId:  pointer.MakePtr(string(testAccount)),
					LicenseKey: pointer.MakePtr(string(testLicenseKey)),
					ReadToken:  pointer.MakePtr(string(testReadToken)),
					WriteToken: pointer.MakePtr(string(testWriteToken)),
				},
				nil,
			),
			expected: &integration.ThirdPartyIntegration{
				ID: testID,
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type:     "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{},
				},
				NewRelicSecrets: &integration.NewRelicSecrets{
					AccountID:  testAccount,
					LicenseKey: testLicenseKey,
					ReadToken:  testReadToken,
					WriteToken: testWriteToken,
				},
			},
			expectedError: nil,
		},

		{
			title: "API failure gets passed through",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type:     "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{},
				},
				NewRelicSecrets: &integration.NewRelicSecrets{
					AccountID:  testAccount,
					LicenseKey: testLicenseKey,
					ReadToken:  testReadToken,
					WriteToken: testWriteToken,
				},
			},
			api: testUpdateIntegrationAPI(
				nil,
				ErrFakeFailure,
			),
			expected:      nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse config returns before calling API",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "BLAH",
					Datadog: &akov2.DatadogIntegration{
						Region:                       testRegion,
						SendCollectionLatencyMetrics: pointer.MakePtr("true"),
						SendDatabaseMetrics:          nil,
					},
				},
				DatadogSecrets: &integration.DatadogSecrets{
					APIKey: "",
				},
			},
			expected:      nil,
			expectedError: integration.ErrUnsupportedIntegrationType,
		},

		{
			title: "failure to parse API reply",
			integration: &integration.ThirdPartyIntegration{
				AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
					Type:     "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{},
				},
				NewRelicSecrets: &integration.NewRelicSecrets{
					AccountID:  testAccount,
					LicenseKey: testLicenseKey,
					ReadToken:  testReadToken,
					WriteToken: testWriteToken,
				},
			},
			api: testUpdateIntegrationAPI(
				&admin.ThirdPartyIntegration{
					Id:                           pointer.MakePtr(string(testID)),
					Type:                         pointer.MakePtr("BLAH"),
					Region:                       pointer.MakePtr(string(testRegion)),
					SendCollectionLatencyMetrics: pointer.MakePtr(true),
				},
				nil,
			),
			expected:      nil,
			expectedError: integration.ErrUnsupportedIntegrationType,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := integration.NewThirdPartyIntegrationService(tc.api)
			updatedIntegrations, err := s.Update(ctx, testProjectID, tc.integration)
			assert.Equal(t, tc.expected, updatedIntegrations)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestIntegrationDelete(t *testing.T) {
	for _, tc := range []struct {
		title         string
		api           admin.ThirdPartyIntegrationsApi
		expectedError error
	}{
		{
			title:         "successful api delete returns success",
			api:           testDeleteIntegrationAPI(nil),
			expectedError: nil,
		},

		{
			title:         "generic API failure passes though",
			api:           testDeleteIntegrationAPI(ErrFakeFailure),
			expectedError: ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := integration.NewThirdPartyIntegrationService(tc.api)
			err := s.Delete(ctx, testProjectID, testIntegrationType)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testCreateIntegrationAPI(integrations []admin.ThirdPartyIntegration, err error) admin.ThirdPartyIntegrationsApi {
	var apiMock mockadmin.ThirdPartyIntegrationsApi

	apiMock.EXPECT().CreateGroupIntegration(
		mock.Anything, mock.Anything, testProjectID, mock.Anything,
	).Return(admin.CreateGroupIntegrationApiRequest{ApiService: &apiMock})

	paginatedIntegration := &admin.PaginatedIntegration{}
	paginatedIntegration.Results = &integrations
	apiMock.EXPECT().CreateGroupIntegrationExecute(
		mock.AnythingOfType("admin.CreateGroupIntegrationApiRequest"),
	).Return(paginatedIntegration, nil, err)
	return &apiMock
}

func testGetIntegrationAPI(integration *admin.ThirdPartyIntegration, err error) admin.ThirdPartyIntegrationsApi {
	var apiMock mockadmin.ThirdPartyIntegrationsApi

	apiMock.EXPECT().GetGroupIntegration(
		mock.Anything, testProjectID, testIntegrationType,
	).Return(admin.GetGroupIntegrationApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().GetGroupIntegrationExecute(
		mock.AnythingOfType("admin.GetGroupIntegrationApiRequest"),
	).Return(integration, nil, err)
	return &apiMock
}

func testUpdateIntegrationAPI(integration *admin.ThirdPartyIntegration, err error) admin.ThirdPartyIntegrationsApi {
	var apiMock mockadmin.ThirdPartyIntegrationsApi

	apiMock.EXPECT().UpdateGroupIntegration(
		mock.Anything, mock.Anything, testProjectID, mock.Anything,
	).Return(admin.UpdateGroupIntegrationApiRequest{ApiService: &apiMock})

	paginatedIntegration := &admin.PaginatedIntegration{}
	if integration != nil {
		paginatedIntegration.Results = &[]admin.ThirdPartyIntegration{
			*integration,
		}
	}
	apiMock.EXPECT().UpdateGroupIntegrationExecute(
		mock.AnythingOfType("admin.UpdateGroupIntegrationApiRequest"),
	).Return(paginatedIntegration, nil, err)
	return &apiMock
}

func testDeleteIntegrationAPI(err error) admin.ThirdPartyIntegrationsApi {
	var apiMock mockadmin.ThirdPartyIntegrationsApi

	apiMock.EXPECT().DeleteGroupIntegration(
		mock.Anything, testIntegrationType, testProjectID,
	).Return(admin.DeleteGroupIntegrationApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().DeleteGroupIntegrationExecute(
		mock.AnythingOfType("admin.DeleteGroupIntegrationApiRequest"),
	).Return(nil, err)
	return &apiMock
}
