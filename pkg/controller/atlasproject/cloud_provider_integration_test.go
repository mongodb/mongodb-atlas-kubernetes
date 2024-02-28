package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"go.uber.org/zap/zaptest"
)

func TestSyncCloudProviderIntegration(t *testing.T) {
	t.Run("should fail when atlas is unavailable", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(nil, nil, errors.New("service unavailable"))
		atlasClient := admin.APIClient{
			CloudProviderAccessApi: cpa,
		}
		workflowCtx := &workflow.Context{
			SdkClient: &atlasClient,
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			CreateCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessRoleExecute(mock.Anything).
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
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeCloudProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
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
			DeauthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, nil)
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeCloudProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
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
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Log:       zaptest.NewLogger(t).Sugar(),
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			CreateCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessRoleExecute(mock.Anything).
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
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			CreateCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessRoleExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Log:       zaptest.NewLogger(t).Sugar(),
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
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
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Log:       zaptest.NewLogger(t).Sugar(),
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, nil)
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, errors.New("service unavailable"))
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Log:       zaptest.NewLogger(t).Sugar(),
			Context:   context.Background(),
		}

		deleteCloudProviderAccess(workflowCtx, "projectID", cpaStatus)
		assert.Equal(t, "service unavailable", cpaStatus.ErrorMessage)
	})
}

func TestCanCloudProviderIntegrationReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, false, &akov2.AtlasProject{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{}
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{},
			Context:   context.Background(),
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)
		assert.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		assert.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("failed to retrieve data"))
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.EqualError(t, err, "failed to retrieve data")
		assert.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{}}, &http.Response{}, nil)
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn1"),
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderIntegrations\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn1"),
						},
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn2"),
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn2",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderIntegrations\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when access was created but not authorized yet", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName: "AWS",
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderIntegrations\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when unable to reconcile cloud provider integration", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn1"),
						},
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn2"),
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn3",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderIntegrations\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when migrating configuration but spec are equal", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn1"),
						},
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn2"),
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn2",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderAccessRoles\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result, err := canCloudProviderIntegrationReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})
}

func TestEnsureCloudProviderIntegration(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("failed to retrieve data"))
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, true)

		assert.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(
				&admin.CloudProviderAccessRoles{
					AwsIamRoles: &[]admin.CloudProviderAccessAWSIAMRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn1"),
						},
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: pointer.MakePtr("arn2"),
						},
					},
				},
				&http.Response{},
				nil,
			)
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn3",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderIntegrations\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, true)

		assert.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Cloud Provider Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})

	t.Run("should return earlier when there are not items to operate", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{}
		workflowCtx := &workflow.Context{Context: context.Background()}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.OK(),
			result,
		)
	})

	t.Run("should fail to reconcile", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-2",
					},
				},
			},
		}
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("failed to retrieve data"))
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.Terminate(workflow.ProjectCloudIntegrationsIsNotReadyInAtlas, "unable to fetch cloud provider access from Atlas: failed to retrieve data"),
			result,
		)
	})

	t.Run("should reconcile without reach ready status", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
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
				},
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			CreateCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			CreateCloudProviderAccessRoleExecute(mock.Anything).
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
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeCloudProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
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
			DeauthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeauthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			DeauthorizeCloudProviderAccessRoleExecute(mock.Anything).
			Return(&http.Response{}, nil)
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.InProgress(workflow.ProjectCloudIntegrationsIsNotReadyInAtlas, "not all entries are authorized"),
			result,
		)
	})

	t.Run("should reconcile and reach ready status", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderIntegrations: []akov2.CloudProviderIntegration{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-2",
					},
				},
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeCloudProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
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
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.OK(),
			result,
		)
		assert.Len(t, workflowCtx.Conditions(), 1)
		assert.Equal(t, status.CloudProviderIntegrationReadyType, workflowCtx.Conditions()[0].Type)
		assert.Equal(t, "True", string(workflowCtx.Conditions()[0].Status))
		assert.Empty(t, workflowCtx.Conditions()[0].Message)
	})

	t.Run("should reconcile and reach ready status using deprecated configuration", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				CloudProviderAccessRoles: []akov2.CloudProviderAccessRole{ //nolint:staticcheck SA1019
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-1",
					},
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "aws:arn/my_role-2",
					},
				},
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
		cpa := atlasmock.NewCloudProviderAccessApiMock(t)
		cpa.EXPECT().
			ListCloudProviderAccessRoles(mock.Anything, mock.Anything).
			Return(admin.ListCloudProviderAccessRolesApiRequest{ApiService: cpa})
		cpa.EXPECT().
			ListCloudProviderAccessRolesExecute(mock.Anything).
			Return(&admin.CloudProviderAccessRoles{AwsIamRoles: &atlasCPAs}, &http.Response{}, nil)
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRole(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.AuthorizeCloudProviderAccessRoleApiRequest{ApiService: cpa})
		cpa.EXPECT().
			AuthorizeCloudProviderAccessRoleExecute(mock.Anything).
			RunAndReturn(
				func(request admin.AuthorizeCloudProviderAccessRoleApiRequest) (*admin.CloudProviderAccessRole, *http.Response, error) {
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
		workflowCtx := workflow.Context{
			SdkClient: &admin.APIClient{CloudProviderAccessApi: cpa},
			Context:   context.Background(),
		}
		result := ensureCloudProviderIntegration(&workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.OK(),
			result,
		)
		assert.Len(t, workflowCtx.Conditions(), 1)
		assert.Equal(t, status.CloudProviderIntegrationReadyType, workflowCtx.Conditions()[0].Type)
		assert.Equal(t, "True", string(workflowCtx.Conditions()[0].Status))
		assert.Equal(
			t,
			"The CloudProviderAccessRole has been deprecated, please move your configuration under CloudProviderIntegration.",
			workflowCtx.Conditions()[0].Message,
		)
	})
}
