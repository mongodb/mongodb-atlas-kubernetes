package atlasproject

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

type ipAccessListClient struct {
	ListFunc   func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error)
	CreateFunc func(projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error)
	DeleteFunc func(projectID, entry string) (*mongodbatlas.Response, error)
}

func (c *ipAccessListClient) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
	return c.ListFunc(projectID)
}

func (c *ipAccessListClient) Get(_ context.Context, _ string, _ string) (*mongodbatlas.ProjectIPAccessList, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *ipAccessListClient) Create(_ context.Context, projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
	return c.CreateFunc(projectID, ipAccessLists)
}

func (c *ipAccessListClient) Delete(_ context.Context, projectID, entry string) (*mongodbatlas.Response, error) {
	return c.DeleteFunc(projectID, entry)
}

func TestFilterActiveIPAccessLists(t *testing.T) {
	t.Run("One expired, one active", func(t *testing.T) {
		dateBefore := time.Now().UTC().Add(time.Hour * -1).Format("2006-01-02T15:04:05.999Z")
		dateAfter := time.Now().UTC().Add(time.Hour * 5).Format("2006-01-02T15:04:05.999Z")
		ipAccessExpired := project.IPAccessList{DeleteAfterDate: dateBefore}
		ipAccessActive := project.IPAccessList{DeleteAfterDate: dateAfter}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive, ipAccessExpired})
		assert.Equal(t, []project.IPAccessList{ipAccessActive}, active)
		assert.Equal(t, []project.IPAccessList{ipAccessExpired}, expired)
	})
	t.Run("Two active", func(t *testing.T) {
		dateAfter1 := time.Now().UTC().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		dateAfter2 := time.Now().UTC().Add(time.Hour * 5).Format("2006-01-02T15:04:05")
		ipAccessActive1 := project.IPAccessList{DeleteAfterDate: dateAfter1}
		ipAccessActive2 := project.IPAccessList{DeleteAfterDate: dateAfter2}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Equal(t, []project.IPAccessList{ipAccessActive2, ipAccessActive1}, active)
		assert.Empty(t, expired)
	})
}

func TestCanIPAccessListReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canIPAccessListReconcile(context.TODO(), mongodbatlas.Client{}, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canIPAccessListReconcile(context.TODO(), mongodbatlas.Client{}, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canIPAccessListReconcile(context.TODO(), atlasClient, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{TotalCount: 0}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canIPAccessListReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								GroupID:   "123456",
								CIDRBlock: "192.168.0.0/24",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
					},
					{
						CIDRBlock: "10.0.0.0/24",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"projectIpAccessList\":[{\"cidrBlock\":\"192.168.0.0/24\"}]}"})
		result, err := canIPAccessListReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								GroupID:   "123456",
								CIDRBlock: "192.168.0.0/24",
							},
							{
								GroupID:   "123456",
								CIDRBlock: "10.0.0.0/24",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
					},
					{
						CIDRBlock: "10.0.0.0/24",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"projectIpAccessList\":[{\"cidrBlock\":\"192.168.0.0/24\"}]}"})
		result, err := canIPAccessListReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile IP Access List", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								GroupID:   "123456",
								CIDRBlock: "192.168.0.0/24",
							},
							{
								GroupID:   "123456",
								CIDRBlock: "10.0.0.0/24",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
					},
					{
						CIDRBlock: "10.1.0.0/24",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"projectIpAccessList\":[{\"cidrBlock\":\"192.168.0.0/24\"}]}"})
		result, err := canIPAccessListReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureIPAccessList(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureIPAccessList(context.TODO(), workflowCtx, atlas.CustomIPAccessListStatus(&atlasClient), akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								GroupID:   "123456",
								CIDRBlock: "192.168.0.0/24",
							},
							{
								GroupID:   "123456",
								CIDRBlock: "10.1.0.0/24",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
					},
					{
						CIDRBlock: "10.0.0.0/24",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"projectIpAccessList\":[{\"cidrBlock\":\"192.168.0.0/24\"}]}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureIPAccessList(context.TODO(), workflowCtx, atlas.CustomIPAccessListStatus(&atlasClient), akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile IP Access List due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})

	t.Run("should reconcile successfully", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				ListFunc: func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								GroupID:   "123456",
								CIDRBlock: "192.168.0.10/32",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
				CreateFunc: func(projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectIPAccessLists{
						Results: []mongodbatlas.ProjectIPAccessList{
							{
								CIDRBlock: "192.168.0.0/24",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
				DeleteFunc: func(projectID, entry string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
					},
					{
						CIDRBlock:       "10.0.0.0/24",
						DeleteAfterDate: "2022-12-25T14:30:15",
					},
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureIPAccessList(
			context.TODO(),
			workflowCtx,
			func(ctx context.Context, projectID, entryValue string) (string, error) {
				return "ACTIVE", nil
			},
			akoProject,
			false,
		)

		assert.Equal(t, workflow.OK(), result)
	})
}

func TestSyncIPAccessList(t *testing.T) {
	t.Run("should fail to perform deletion", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				CIDRBlock: "10.0.0.0/24",
			},
		}
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				DeleteFunc: func(projectID, entry string) (*mongodbatlas.Response, error) {
					return nil, errors.New("failed")
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		assert.ErrorContains(t, syncIPAccessList(context.TODO(), workflowCtx, "projectID", current, desired), "failed")
	})

	t.Run("should fail to perform creation", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				CIDRBlock: "10.0.0.0/24",
			},
		}
		atlasClient := mongodbatlas.Client{
			ProjectIPAccessList: &ipAccessListClient{
				CreateFunc: func(projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed")
				},
				DeleteFunc: func(projectID, entry string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			},
		}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		assert.ErrorContains(t, syncIPAccessList(context.TODO(), workflowCtx, "projectID", current, desired), "failed")
	})

	t.Run("should succeed when there are no changes", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		atlasClient := mongodbatlas.Client{}
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}

		assert.NoError(t, syncIPAccessList(context.TODO(), workflowCtx, "projectID", current, desired))
	})
}
