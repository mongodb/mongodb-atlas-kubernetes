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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

func TestSyncCloudProviderIntegration(t *testing.T) {
	t.Run("should fail when atlas is unavailable", func(t *testing.T) {
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			ListCloudProviderAccess(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessExecute(mock.Anything).
			Return(nil, nil, errors.New("service unavailable"))
		atlasClient := admin.APIClient{
			CloudProviderAccessApi: cpa,
		}
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &atlasClient,
			},
			Context: context.Background(),
		}
		result, err := syncCloudProviderIntegration(workflowCtx, "projectID", []akov2.CloudProviderIntegration{})
		assert.EqualError(t, err, "unable to fetch cloud provider access from Atlas: service unavailable")
		assert.False(t, result)
	})

	t.Run("should synchronize all operations without reach ready status", func(t *testing.T) {
		cpas := []akov2.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
			},
			{
				ProviderName: "AWS",
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(time.Now().Add(5 * time.Minute)),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-4"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-4"),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my-role-4"),
				CreatedDate:                pointer.MakePtr(time.Now()),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-4"),
			},
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			ListCloudProviderAccess(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			CreateCloudProviderAccess(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRole{
					AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-3"),
					AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-3"),
					CreatedDate:                pointer.MakePtr(time.Now()),
					ProviderName:               "AWS",
					RoleId:                     pointer.MakePtr("role-3"),
				},
				&http.Response{},
				nil,
			)
		cpa.EXPECT().
			AuthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
					atlasCPA := admin.CloudProviderAccessRole{
						AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
						AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
						CreatedDate:                pointer.MakePtr(time.Now()),
						AuthorizedDate:             pointer.MakePtr(time.Now().Add(5 * time.Minute)),
						IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
						ProviderName:               "AWS",
						RoleId:                     pointer.MakePtr("role-2"),
					}

					return &atlasCPA, &http.Response{}, nil
				},
			)
		cpa.EXPECT().
			DeauthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, nil)
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Context: context.Background(),
		}

		result, err := syncCloudProviderIntegration(workflowCtx, "projectID", cpas)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should synchronize all operations and reach ready status", func(t *testing.T) {
		cpas := []akov2.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(time.Now().Add(5 * time.Minute)),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			ListCloudProviderAccess(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
					atlasCPA := admin.CloudProviderAccessRole{
						AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
						AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
						CreatedDate:                pointer.MakePtr(time.Now()),
						AuthorizedDate:             pointer.MakePtr(time.Now().Add(5 * time.Minute)),
						IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
						ProviderName:               "AWS",
						RoleId:                     pointer.MakePtr("role-2"),
					}

					return &atlasCPA, &http.Response{}, nil
				},
			)
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Context: context.Background(),
		}

		result, err := syncCloudProviderIntegration(workflowCtx, "projectID", cpas)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should synchronize operations with errors", func(t *testing.T) {
		cpas := []akov2.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(time.Now().Add(5 * time.Minute)),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(time.Now()),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			ListCloudProviderAccess(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeProviderAccessRoleExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Log:     zaptest.NewLogger(t).Sugar(),
			Context: context.Background(),
		}

		result, err := syncCloudProviderIntegration(workflowCtx, "projectID", cpas)
		assert.EqualError(t, err, "not all items were synchronized successfully")
		assert.False(t, result)
	})
}

func TestInitiateStatus(t *testing.T) {
	t.Run("should create a cloud provider status as new", func(t *testing.T) {
		expected := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		spec := []akov2.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
			},
			{
				ProviderName: "AWS",
			},
		}

		assert.Equal(t, expected, initiateStatuses(spec))
	})
}

func TestEnrichStatuses(t *testing.T) {
	t.Run("all statuses are new", func(t *testing.T) {
		expected := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, []admin.CloudProviderAccessAWSIAMRole{}))
	})

	t.Run("one new and one authorized statuses", func(t *testing.T) {
		createdAt := time.Now()
		authorizedAt := createdAt.Add(5 * time.Minute)
		expected := []*status.CloudProviderIntegration{
			{
				AtlasAWSAccountArn:         "atlas-account-arn",
				AtlasAssumedRoleExternalID: "atlas-external-role-id",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderIntegrationStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id"),
				CreatedDate:                pointer.MakePtr(createdAt),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one new, one created and one authorized statuses", func(t *testing.T) {
		createdAt := time.Now()
		authorizedAt := createdAt.Add(5 * time.Minute)
		expected := []*status.CloudProviderIntegration{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderIntegrationStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				CreatedDate:                pointer.MakePtr(createdAt),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one new, one created, one authorized, and one authorized to remove statuses", func(t *testing.T) {
		createdAt := time.Now()
		authorizedAt := createdAt.Add(5 * time.Minute)
		expected := []*status.CloudProviderIntegration{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderIntegrationStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderIntegrationStatusDeAuthorize,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-3"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-3"),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-3"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-3"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one created with empty ARN, one created, and one authorized statuses", func(t *testing.T) {
		createdAt := time.Now()
		authorizedAt := createdAt.Add(5 * time.Minute)
		expected := []*status.CloudProviderIntegration{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderIntegrationStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-3"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-3"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-3"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one created with empty ARN, one created, one authorized, and one to be removed statuses", func(t *testing.T) {
		createdAt := time.Now()
		authorizedAt := createdAt.Add(5 * time.Minute)
		expected := []*status.CloudProviderIntegration{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderIntegrationStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderIntegrationStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-4",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-4",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				ProviderName:               "AWS",
				RoleID:                     "role-4",
				Status:                     status.CloudProviderIntegrationStatusDeAuthorize,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				AuthorizedDate:             pointer.MakePtr(authorizedAt),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-1"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-3"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-3"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-3"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				IamAssumedRoleArn:          pointer.MakePtr("aws:arn/my_role-2"),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-4"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-4"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-4"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("match two status with empty ARN and two existing on Atlas", func(t *testing.T) {
		createdAt := time.Now()
		expected := []*status.CloudProviderIntegration{
			{
				ProviderName:               "AWS",
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				RoleID:                     "role-1",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				Status:                     status.CloudProviderIntegrationStatusCreated,
			},
			{
				ProviderName:               "AWS",
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				RoleID:                     "role-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				Status:                     status.CloudProviderIntegrationStatusCreated,
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("match two status with empty ARN and update them with ARN", func(t *testing.T) {
		createdAt := time.Now()
		expected := []*status.CloudProviderIntegration{
			{
				ProviderName:               "AWS",
				IamAssumedRoleArn:          "was:arn/role-1",
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				RoleID:                     "role-1",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				Status:                     status.CloudProviderIntegrationStatusCreated,
			},
			{
				ProviderName:               "AWS",
				IamAssumedRoleArn:          "was:arn/role-2",
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				RoleID:                     "role-2",
				CreatedDate:                timeutil.FormatISO8601(createdAt),
				Status:                     status.CloudProviderIntegrationStatusCreated,
			},
		}
		statuses := []*status.CloudProviderIntegration{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "was:arn/role-1",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "was:arn/role-2",
				Status:            status.CloudProviderIntegrationStatusNew,
			},
		}
		atlasCPAs := []admin.CloudProviderAccessAWSIAMRole{
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-1"),
			},
			{
				AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-2"),
				AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-2"),
				CreatedDate:                pointer.MakePtr(createdAt),
				ProviderName:               "AWS",
				RoleId:                     pointer.MakePtr("role-2"),
			},
		}

		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})
}

func TestCreateCloudProviderIntegration(t *testing.T) {
	t.Run("should create cloud provider integration successfully", func(t *testing.T) {
		createdAt := time.Now()
		expected := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                timeutil.FormatISO8601(createdAt),
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusCreated,
			ErrorMessage:               "",
		}
		cpaStatus := &status.CloudProviderIntegration{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderIntegrationStatusNew,
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			CreateCloudProviderAccess(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRole{
					AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
					AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
					CreatedDate:                pointer.MakePtr(createdAt),
					ProviderName:               "AWS",
					RoleId:                     pointer.MakePtr("role-1"),
				},
				&http.Response{},
				nil,
			)
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Context: context.Background(),
		}

		assert.Equal(t, expected, createCloudProviderAccess(workflowCtx, "projectID", cpaStatus))
	})

	t.Run("should fail to create cloud provider integration", func(t *testing.T) {
		expected := &status.CloudProviderIntegration{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderIntegrationStatusFailedToCreate,
			ErrorMessage:      "service unavailable",
		}
		cpaStatus := &status.CloudProviderIntegration{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderIntegrationStatusNew,
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			CreateCloudProviderAccess(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Log:     zaptest.NewLogger(t).Sugar(),
			Context: context.Background(),
		}

		assert.Equal(t, expected, createCloudProviderAccess(workflowCtx, "projectID", cpaStatus))
	})
}

func TestAuthorizeCloudProviderIntegration(t *testing.T) {
	createdAt := time.Now()
	authorizedAt := createdAt.Add(5 * time.Minute)
	t.Run("should authorize cloud provider integration successfully", func(t *testing.T) {
		expected := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                timeutil.FormatISO8601(createdAt),
			AuthorizedDate:             timeutil.FormatISO8601(authorizedAt),
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusAuthorized,
			ErrorMessage:               "",
		}
		cpaStatus := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                timeutil.FormatISO8601(createdAt),
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusNew,
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			AuthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeProviderAccessRoleExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRole{
					AtlasAWSAccountArn:         pointer.MakePtr("atlas-account-arn-1"),
					AtlasAssumedRoleExternalId: pointer.MakePtr("atlas-external-role-id-1"),
					CreatedDate:                pointer.MakePtr(createdAt),
					AuthorizedDate:             pointer.MakePtr(authorizedAt),
					ProviderName:               "AWS",
					RoleId:                     pointer.MakePtr("role-1"),
				},
				&http.Response{},
				nil,
			)
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Context: context.Background(),
		}

		assert.Equal(t, expected, authorizeCloudProviderAccess(workflowCtx, "projectID", cpaStatus))
	})

	t.Run("should fail to authorize cloud provider integration", func(t *testing.T) {
		expected := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusFailedToAuthorize,
			ErrorMessage:               "service unavailable",
		}
		cpaStatus := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusCreated,
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			AuthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeProviderAccessRoleExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Log:     zaptest.NewLogger(t).Sugar(),
			Context: context.Background(),
		}

		assert.Equal(t, expected, authorizeCloudProviderAccess(workflowCtx, "projectID", cpaStatus))
	})
}

func TestDeleteCloudProviderIntegration(t *testing.T) {
	t.Run("should delete cloud provider integration successfully", func(t *testing.T) {
		cpaStatus := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			AuthorizedDate:             "authorized-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusFailedToDeAuthorize,
			ErrorMessage:               "",
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			DeauthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, nil)
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Context: context.Background(),
		}

		deleteCloudProviderAccess(workflowCtx, "projectID", cpaStatus)
		assert.Empty(t, cpaStatus.ErrorMessage)
	})

	t.Run("should fail to delete cloud provider integration", func(t *testing.T) {
		cpaStatus := &status.CloudProviderIntegration{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderIntegrationStatusFailedToDeAuthorize,
			ErrorMessage:               "",
		}
		cpa := mockadmin.NewCloudProviderAccessApi(t)
		cpa.EXPECT().
			DeauthorizeProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					CloudProviderAccessApi: cpa,
				},
			},
			Log:     zaptest.NewLogger(t).Sugar(),
			Context: context.Background(),
		}

		deleteCloudProviderAccess(workflowCtx, "projectID", cpaStatus)
		assert.Equal(t, "service unavailable", cpaStatus.ErrorMessage)
	})
}
