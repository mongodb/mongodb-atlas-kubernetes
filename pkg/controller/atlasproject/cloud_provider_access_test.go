package atlasproject

import (
	"context"
	"errors"
	"testing"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

type cloudProviderAccessClient struct {
	ListRolesFunc func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error)
}

func (c *cloudProviderAccessClient) ListRoles(_ context.Context, projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
	return c.ListRolesFunc(projectID)
}

func (c *cloudProviderAccessClient) GetRole(_ context.Context, _ string, _ string) (*mongodbatlas.AWSIAMRole, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *cloudProviderAccessClient) CreateRole(_ context.Context, _ string, _ *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.AWSIAMRole, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *cloudProviderAccessClient) AuthorizeRole(_ context.Context, _ string, _ string, _ *mongodbatlas.CloudProviderAuthorizationRequest) (*mongodbatlas.AWSIAMRole, *mongodbatlas.Response, error) {
	return nil, nil, nil
}
func (c *cloudProviderAccessClient) DeauthorizeRole(_ context.Context, _ *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
	return nil, nil
}

func TestCanCloudProviderAccessReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canCloudProviderAccessReconcile(context.TODO(), mongodbatlas.Client{}, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), mongodbatlas.Client{}, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
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

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.AWSIAMRole{},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canCloudProviderAccessReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.AWSIAMRole{
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

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.AWSIAMRole{
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

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile Cloud Provider Access", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.AWSIAMRole{
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

		require.NoError(t, err)
		require.False(t, result)
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

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CloudProviderAccess: &cloudProviderAccessClient{
				ListRolesFunc: func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderAccessRoles{
						AWSIAMRoles: []mongodbatlas.AWSIAMRole{
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

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
