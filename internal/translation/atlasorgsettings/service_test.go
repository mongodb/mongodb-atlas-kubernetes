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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	testOrgID = "fake-org-id"
)

var (
	ErrFakeAPIFailure = errors.New("fake API failure")
)

func TestNewAtlasOrgSettingsService(t *testing.T) {
	mockAPI := &mockadmin.OrganizationsApi{}

	service := NewAtlasOrgSettingsService(mockAPI)

	assert.NotNil(t, service)
	assert.IsType(t, &AtlasOrgSettingsServiceImpl{}, service)
}

func TestAtlasOrgSettingsService_Get(t *testing.T) {
	for _, tc := range []struct {
		title         string
		orgID         string
		api           admin.OrganizationsApi
		expected      *AtlasOrgSettings
		expectedError error
	}{
		{
			title: "successful get organization settings",
			orgID: testOrgID,
			api: testGetOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
				http.StatusOK,
				nil,
			),
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  testOrgID,
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
			},
			expectedError: nil,
		},
		{
			title: "successful get with partial settings",
			orgID: testOrgID,
			api: testGetOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(false),
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                pointer.MakePtr(false),
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        pointer.MakePtr("admin@example.com"),
					StreamsCrossGroupEnabled:               nil,
				},
				http.StatusOK,
				nil,
			),
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  testOrgID,
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(false),
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                pointer.MakePtr(false),
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        pointer.MakePtr("admin@example.com"),
					StreamsCrossGroupEnabled:               nil,
				},
			},
			expectedError: nil,
		},
		{
			title: "API failure gets passed through",
			orgID: testOrgID,
			api: testGetOrgSettingsAPI(
				nil,
				http.StatusInternalServerError,
				ErrFakeAPIFailure,
			),
			expected:      nil,
			expectedError: fmt.Errorf("failed to get AtlasOrgSettings: %w", ErrFakeAPIFailure),
		},
		{
			title: "non-200 status code with no error returns error",
			orgID: testOrgID,
			api: testGetOrgSettingsAPI(
				&admin.OrganizationSettings{},
				http.StatusNotFound,
				nil,
			),
			expected:      nil,
			expectedError: nil, // The service returns the original error (nil) for non-200 status
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			service := NewAtlasOrgSettingsService(tc.api)

			result, err := service.Get(context.Background(), tc.orgID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				if tc.expected != nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				} else {
					// For the non-200 status case, we expect nil result but may have nil error too
					assert.Nil(t, result)
				}
			}
		})
	}
}

func TestAtlasOrgSettingsService_Update(t *testing.T) {
	inputSettings := &AtlasOrgSettings{
		AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
			OrgID:                                  testOrgID,
			ApiAccessListRequired:                  pointer.MakePtr(true),
			GenAIFeaturesEnabled:                   pointer.MakePtr(false),
			MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
			MultiFactorAuthRequired:                pointer.MakePtr(true),
			RestrictEmployeeAccess:                 pointer.MakePtr(false),
			SecurityContact:                        pointer.MakePtr("security@example.com"),
			StreamsCrossGroupEnabled:               pointer.MakePtr(true),
		},
	}

	// this is a common function to make golangci-lint happy
	createSuccessfulUpdateTest := func(title string, statusCode int) struct {
		title         string
		orgID         string
		settings      *AtlasOrgSettings
		api           admin.OrganizationsApi
		expected      *AtlasOrgSettings
		expectedError error
	} {
		return struct {
			title         string
			orgID         string
			settings      *AtlasOrgSettings
			api           admin.OrganizationsApi
			expected      *AtlasOrgSettings
			expectedError error
		}{
			title:    title,
			orgID:    testOrgID,
			settings: inputSettings,
			api: testUpdateOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
				},
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
				statusCode,
				nil,
			),
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  testOrgID,
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
			},
			expectedError: nil,
		}
	}

	for _, tc := range []struct {
		title         string
		orgID         string
		settings      *AtlasOrgSettings
		api           admin.OrganizationsApi
		expected      *AtlasOrgSettings
		expectedError error
	}{
		createSuccessfulUpdateTest("successful update organization settings", http.StatusOK),
		createSuccessfulUpdateTest("successful update with 201 status code", http.StatusCreated),
		{
			title:    "API failure gets passed through",
			orgID:    testOrgID,
			settings: inputSettings,
			api: testUpdateOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
				},
				nil,
				http.StatusInternalServerError,
				ErrFakeAPIFailure,
			),
			expected:      nil,
			expectedError: fmt.Errorf("failed to update AtlasOrgSettings: %w", ErrFakeAPIFailure),
		},
		{
			title:         "nil settings returns nil",
			orgID:         testOrgID,
			settings:      nil,
			api:           &mockadmin.OrganizationsApi{}, // No expectations set since it shouldn't be called
			expected:      nil,
			expectedError: nil,
		},
		{
			title: "settings that convert to nil atlas settings return nil",
			orgID: testOrgID,
			settings: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID: testOrgID,
					// All other fields are nil, so ToAtlas might return nil
				},
			},
			api: testUpdateOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
				},
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
				http.StatusOK,
				nil,
			),
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  testOrgID,
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
			},
			expectedError: nil,
		},
		{
			title:    "non-200/201 status code with no error returns error",
			orgID:    testOrgID,
			settings: inputSettings,
			api: testUpdateOrgSettingsAPI(
				&admin.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
				},
				&admin.OrganizationSettings{},
				http.StatusBadRequest,
				nil,
			),
			expected:      nil,
			expectedError: nil, // The service returns the original error (nil) for non-200/201 status
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			service := NewAtlasOrgSettingsService(tc.api)

			result, err := service.Update(context.Background(), tc.orgID, tc.settings)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				if tc.expected != nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				} else {
					// For cases where we expect nil result
					assert.Nil(t, result)
				}
			}
		})
	}
}

func testGetOrgSettingsAPI(response *admin.OrganizationSettings, statusCode int, err error) admin.OrganizationsApi {
	mockAPI := &mockadmin.OrganizationsApi{}

	request := admin.GetOrganizationSettingsApiRequest{ApiService: mockAPI}
	mockAPI.EXPECT().GetOrganizationSettings(mock.Anything, testOrgID).Return(request)
	mockAPI.EXPECT().GetOrganizationSettingsExecute(
		mock.AnythingOfType("admin.GetOrganizationSettingsApiRequest")).Return(response, &http.Response{StatusCode: statusCode}, err)
	return mockAPI
}

func testUpdateOrgSettingsAPI(input *admin.OrganizationSettings, response *admin.OrganizationSettings, statusCode int, err error) admin.OrganizationsApi {
	mockAPI := &mockadmin.OrganizationsApi{}

	request := admin.UpdateOrganizationSettingsApiRequest{ApiService: mockAPI}

	if input != nil {
		mockAPI.EXPECT().UpdateOrganizationSettings(mock.Anything, testOrgID, input).Return(request)
		mockAPI.EXPECT().UpdateOrganizationSettingsExecute(
			mock.AnythingOfType("admin.UpdateOrganizationSettingsApiRequest")).Return(response, &http.Response{StatusCode: statusCode}, err)
	} else {
		mockAPI.EXPECT().UpdateOrganizationSettings(mock.Anything, testOrgID, mock.Anything).Return(request).Maybe()
	}

	return mockAPI
}
