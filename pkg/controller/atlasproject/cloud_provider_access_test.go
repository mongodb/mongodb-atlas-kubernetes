package atlasproject

import (
	"context"
	"errors"
	"testing"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

type cloudProviderAccessClient struct {
	ListRolesFunc       func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error)
	CreateRoleFunc      func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error)
	AuthorizeRoleFunc   func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error)
	DeauthorizeRoleFunc func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error)
}

func (c *cloudProviderAccessClient) ListRoles(_ context.Context, projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
	return c.ListRolesFunc(projectID)
}

func (c *cloudProviderAccessClient) GetRole(_ context.Context, _ string, _ string) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *cloudProviderAccessClient) CreateRole(_ context.Context, projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	return c.CreateRoleFunc(projectID, cpa)
}

func (c *cloudProviderAccessClient) AuthorizeRole(_ context.Context, projectID string, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	return c.AuthorizeRoleFunc(projectID, roleID, cpa)
}
func (c *cloudProviderAccessClient) DeauthorizeRole(_ context.Context, cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
	return c.DeauthorizeRoleFunc(cpa)
}

func TestSyncCloudProviderAccess(t *testing.T) {
	t.Run("should fail when atlas is unavailable", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("service unavailable")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result, err := syncCloudProviderAccess(context.TODO(), workflowCtx, "projectID", []mdbv1.CloudProviderAccessRole{})
		assert.EqualError(t, err, "unable to fetch cloud provider access from Atlas: service unavailable")
		assert.False(t, result)
	})

	t.Run("should synchronize all operations without reach ready status", func(t *testing.T) {
		cpas := []mdbv1.CloudProviderAccessRole{
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
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-4",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-4",
				IAMAssumedRoleARN:          "aws:arn/my-role-4",
				CreatedDate:                "created-date-4",
				ProviderName:               "AWS",
				RoleID:                     "role-4",
			},
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: atlasCPAs,
					}, &mongodbatlas.Response{}, nil
				},
				CreateRoleFunc: func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRole{
						AtlasAWSAccountARN:         "atlas-account-arn-3",
						AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
						CreatedDate:                "created-date-3",
						ProviderName:               "AWS",
						RoleID:                     "role-3",
					}, &mongodbatlas.Response{}, nil
				},
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					atlasCPA := atlasCPAs[1]
					atlasCPA.AuthorizedDate = "authorized-date-2"

					return &atlasCPA, &mongodbatlas.Response{}, nil
				},
				DeauthorizeRoleFunc: func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
					return &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		result, err := syncCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpas)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should synchronize all operations and reach ready status", func(t *testing.T) {
		cpas := []mdbv1.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: atlasCPAs,
					}, &mongodbatlas.Response{}, nil
				},
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					atlasCPA := atlasCPAs[1]
					atlasCPA.AuthorizedDate = "authorized-date-2"

					return &atlasCPA, &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		result, err := syncCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpas)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should synchronize operations with errors", func(t *testing.T) {
		cpas := []mdbv1.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: atlasCPAs,
					}, &mongodbatlas.Response{}, nil
				},
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("service unavailable")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := syncCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpas)
		assert.EqualError(t, err, "not all items were synchronized successfully")
		assert.False(t, result)
	})
}

func TestInitiateStatus(t *testing.T) {
	t.Run("should create a cloud provider status as new", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		spec := []mdbv1.CloudProviderAccessRole{
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
		expected := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, []mongodbatlas.CloudProviderAccessRole{}))
	})

	t.Run("one new and one authorized statuses", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				AtlasAWSAccountArn:         "atlas-account-arn",
				AtlasAssumedRoleExternalID: "atlas-external-role-id",
				AuthorizedDate:             "authorized-date",
				CreatedDate:                "created-date",
				IamAssumedRoleArn:          "aws:arn/my_role",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderAccessStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn",
				AtlasAssumedRoleExternalID: "atlas-external-role-id",
				AuthorizedDate:             "authorized-date",
				CreatedDate:                "created-date",
				IAMAssumedRoleARN:          "aws:arn/my_role",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one new, one created and one authorized statuses", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderAccessStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one new, one created, one authorized, and one authorized to remove statuses", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderAccessStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				AuthorizedDate:             "authorized-date-3",
				CreatedDate:                "created-date-3",
				IamAssumedRoleArn:          "aws:arn/my_role-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderAccessStatusDeAuthorize,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				AuthorizedDate:             "authorized-date-3",
				CreatedDate:                "created-date-3",
				IAMAssumedRoleARN:          "aws:arn/my_role-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one created with empty ARN, one created, and one authorized statuses", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderAccessStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                "created-date-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                "created-date-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("one created with empty ARN, one created, one authorized, and one to be removed statuses", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IamAssumedRoleArn:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
				Status:                     status.CloudProviderAccessStatusAuthorized,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IamAssumedRoleArn:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                "created-date-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
				Status:                     status.CloudProviderAccessStatusCreated,
				ErrorMessage:               "",
			},
			{
				AtlasAWSAccountArn:         "atlas-account-arn-4",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-4",
				CreatedDate:                "created-date-4",
				ProviderName:               "AWS",
				RoleID:                     "role-4",
				Status:                     status.CloudProviderAccessStatusDeAuthorize,
				ErrorMessage:               "",
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-1",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "aws:arn/my_role-2",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-3",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
				CreatedDate:                "created-date-3",
				ProviderName:               "AWS",
				RoleID:                     "role-3",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-4",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-4",
				CreatedDate:                "created-date-4",
				ProviderName:               "AWS",
				RoleID:                     "role-4",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("match two status with empty ARN and two existing on Atlas", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				ProviderName:               "AWS",
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				RoleID:                     "role-1",
				CreatedDate:                "created-date-1",
				Status:                     status.CloudProviderAccessStatusCreated,
			},
			{
				ProviderName:               "AWS",
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				RoleID:                     "role-2",
				CreatedDate:                "created-date-2",
				Status:                     status.CloudProviderAccessStatusCreated,
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName: "AWS",
				Status:       status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				CreatedDate:                "created-date-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})

	t.Run("match two status with empty ARN and update them with ARN", func(t *testing.T) {
		expected := []*status.CloudProviderAccessRole{
			{
				ProviderName:               "AWS",
				IamAssumedRoleArn:          "was:arn/role-1",
				AtlasAWSAccountArn:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				RoleID:                     "role-1",
				CreatedDate:                "created-date-1",
				Status:                     status.CloudProviderAccessStatusCreated,
			},
			{
				ProviderName:               "AWS",
				IamAssumedRoleArn:          "was:arn/role-2",
				AtlasAWSAccountArn:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				RoleID:                     "role-2",
				CreatedDate:                "created-date-2",
				Status:                     status.CloudProviderAccessStatusCreated,
			},
		}
		statuses := []*status.CloudProviderAccessRole{
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "was:arn/role-1",
				Status:            status.CloudProviderAccessStatusNew,
			},
			{
				ProviderName:      "AWS",
				IamAssumedRoleArn: "was:arn/role-2",
				Status:            status.CloudProviderAccessStatusNew,
			},
		}
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				CreatedDate:                "created-date-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}

		assert.Equal(t, expected, enrichStatuses(statuses, atlasCPAs))
	})
}

func TestCreateCloudProviderAccess(t *testing.T) {
	t.Run("should create cloud provider access successfully", func(t *testing.T) {
		expected := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusCreated,
			ErrorMessage:               "",
		}
		cpa := &status.CloudProviderAccessRole{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderAccessStatusNew,
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				CreateRoleFunc: func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRole{
						AtlasAWSAccountARN:         "atlas-account-arn-1",
						AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
						CreatedDate:                "created-date-1",
						ProviderName:               "AWS",
						RoleID:                     "role-1",
					}, &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		assert.Equal(t, expected, createCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa))
	})

	t.Run("should fail to create cloud provider access", func(t *testing.T) {
		expected := &status.CloudProviderAccessRole{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderAccessStatusFailedToCreate,
			ErrorMessage:      "service unavailable",
		}
		cpa := &status.CloudProviderAccessRole{
			ProviderName:      "AWS",
			IamAssumedRoleArn: "aws:arn/my_role-1",
			Status:            status.CloudProviderAccessStatusNew,
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				CreateRoleFunc: func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("service unavailable")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		assert.Equal(t, expected, createCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa))
	})
}

func TestAuthorizeCloudProviderAccess(t *testing.T) {
	t.Run("should authorize cloud provider access successfully", func(t *testing.T) {
		expected := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			AuthorizedDate:             "authorized-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusAuthorized,
			ErrorMessage:               "",
		}
		cpa := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusNew,
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRole{
						AtlasAWSAccountARN:         "atlas-account-arn-1",
						AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
						CreatedDate:                "created-date-1",
						AuthorizedDate:             "authorized-date-1",
						ProviderName:               "AWS",
						RoleID:                     "role-1",
					}, &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		assert.Equal(t, expected, authorizeCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa))
	})

	t.Run("should fail to authorize cloud provider access", func(t *testing.T) {
		expected := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusFailedToAuthorize,
			ErrorMessage:               "service unavailable",
		}
		cpa := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusCreated,
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("service unavailable")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		assert.Equal(t, expected, authorizeCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa))
	})
}

func TestDeleteCloudProviderAccess(t *testing.T) {
	t.Run("should delete cloud provider access successfully", func(t *testing.T) {
		cpa := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			AuthorizedDate:             "authorized-date-1",
			IamAssumedRoleArn:          "aws:arn/my_role-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusFailedToDeAuthorize,
			ErrorMessage:               "",
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				DeauthorizeRoleFunc: func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
					return &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		deleteCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa)
		assert.Empty(t, cpa.ErrorMessage)
	})

	t.Run("should fail to delete cloud provider access", func(t *testing.T) {
		cpa := &status.CloudProviderAccessRole{
			AtlasAWSAccountArn:         "atlas-account-arn-1",
			AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
			CreatedDate:                "created-date-1",
			ProviderName:               "AWS",
			RoleID:                     "role-1",
			Status:                     status.CloudProviderAccessStatusFailedToDeAuthorize,
			ErrorMessage:               "",
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				DeauthorizeRoleFunc: func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
					return &mongodbatlas.Response{}, errors.New("service unavailable")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		deleteCloudProviderAccess(context.TODO(), workflowCtx, "projectID", cpa)
		assert.Equal(t, "service unavailable", cpa.ErrorMessage)
	})
}

func TestCanCloudProviderAccessReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canCloudProviderAccessReconcile(context.TODO(), mongodbatlas.Client{}, false, &mdbv1.AtlasProject{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), mongodbatlas.Client{}, true, akoProject)
		assert.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		assert.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.EqualError(t, err, "failed to retrieve data")
		assert.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn1",
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderAccessRoles\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn1",
							},
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn2",
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when access was created but not authorized yet", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{
							{
								ProviderName: "AWS",
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
					{
						ProviderName:      "AWS",
						IamAssumedRoleArn: "arn1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderAccessRoles\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when unable to reconcile Cloud Provider Access", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn1",
							},
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn2",
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderAccessRoles\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		assert.NoError(t, err)
		assert.False(t, result)
	})
}

func TestEnsureCloudProviderAccess(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, true)

		assert.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.CloudProviderAccessRole{
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn1",
							},
							{
								ProviderName:      "AWS",
								IAMAssumedRoleARN: "arn2",
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"cloudProviderAccessRoles\":[{\"providerName\":\"AWS\",\"iamAssumedRoleArn\":\"arn1\"}]}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, true)

		assert.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})

	t.Run("should return earlier when there are not items to operate", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		workflowCtx := &workflow.Context{}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.OK(),
			result,
		)
	})

	t.Run("should fail to reconcile", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "unable to fetch cloud provider access from Atlas: failed to retrieve data"),
			result,
		)
	})

	t.Run("should reconcile without reach ready status", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-4",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-4",
				IAMAssumedRoleARN:          "aws:arn/my-role-4",
				CreatedDate:                "created-date-4",
				ProviderName:               "AWS",
				RoleID:                     "role-4",
			},
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: atlasCPAs,
					}, &mongodbatlas.Response{}, nil
				},
				CreateRoleFunc: func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRole{
						AtlasAWSAccountARN:         "atlas-account-arn-3",
						AtlasAssumedRoleExternalID: "atlas-external-role-id-3",
						CreatedDate:                "created-date-3",
						ProviderName:               "AWS",
						RoleID:                     "role-3",
					}, &mongodbatlas.Response{}, nil
				},
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					atlasCPA := atlasCPAs[1]
					atlasCPA.AuthorizedDate = "authorized-date-2"

					return &atlasCPA, &mongodbatlas.Response{}, nil
				},
				DeauthorizeRoleFunc: func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
					return &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.InProgress(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "not all entries are authorized"),
			result,
		)
	})

	t.Run("should reconcile and reach ready status", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
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
		atlasCPAs := []mongodbatlas.CloudProviderAccessRole{
			{
				AtlasAWSAccountARN:         "atlas-account-arn-1",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-1",
				AuthorizedDate:             "authorized-date-1",
				CreatedDate:                "created-date-1",
				IAMAssumedRoleARN:          "aws:arn/my_role-1",
				ProviderName:               "AWS",
				RoleID:                     "role-1",
			},
			{
				AtlasAWSAccountARN:         "atlas-account-arn-2",
				AtlasAssumedRoleExternalID: "atlas-external-role-id-2",
				CreatedDate:                "created-date-2",
				IAMAssumedRoleARN:          "aws:arn/my_role-2",
				ProviderName:               "AWS",
				RoleID:                     "role-2",
			},
		}
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: atlasCPAs,
					}, &mongodbatlas.Response{}, nil
				},
				AuthorizeRoleFunc: func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
					atlasCPA := atlasCPAs[1]
					atlasCPA.AuthorizedDate = "authorized-date-2"

					return &atlasCPA, &mongodbatlas.Response{}, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProviderAccessStatus(context.TODO(), workflowCtx, akoProject, false)
		assert.Equal(
			t,
			workflow.OK(),
			result,
		)
	})
}
