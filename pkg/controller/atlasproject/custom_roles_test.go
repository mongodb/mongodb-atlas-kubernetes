package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestCalculateChanges(t *testing.T) {
	desired := []mdbv1.CustomRole{
		{
			Name: "cr-1",
		},
		{
			Name: "cr-3",
			InheritedRoles: []mdbv1.Role{
				{
					Name:     "admin",
					Database: "test",
				},
			},
		},
		{
			Name: "cr-4",
		},
	}
	current := []mdbv1.CustomRole{
		{
			Name: "cr-1",
		},
		{
			Name: "cr-2",
		},
		{
			Name: "cr-3",
		},
	}

	assert.Equal(
		t,
		CustomRolesOperations{
			Create: map[string]mdbv1.CustomRole{
				"cr-4": {
					Name: "cr-4",
				},
			},
			Update: map[string]mdbv1.CustomRole{
				"cr-3": {
					Name: "cr-3",
					InheritedRoles: []mdbv1.Role{
						{
							Name:     "admin",
							Database: "test",
						},
					},
				},
			},
			Delete: map[string]mdbv1.CustomRole{
				"cr-2": {
					Name: "cr-2",
				},
			},
		},
		calculateChanges(current, desired),
	)
}

func TestSyncCustomRolesStatus(t *testing.T) {
	t.Run("sync status when all operations were done successfully", func(t *testing.T) {
		desired := []mdbv1.CustomRole{
			{
				Name: "cr-1",
			},
			{
				Name: "cr-3",
				InheritedRoles: []mdbv1.Role{
					{
						Name:     "admin",
						Database: "test",
					},
				},
			},
			{
				Name: "cr-4",
			},
		}
		created := map[string]status.CustomRole{
			"cr-4": {
				Name:   "cr-4",
				Status: status.CustomRoleStatusOK,
			},
		}
		updated := map[string]status.CustomRole{
			"cr-3": {
				Name:   "cr-3",
				Status: status.CustomRoleStatusOK,
			},
		}
		deleted := map[string]status.CustomRole{
			"cr-2": {
				Name:   "cr-2",
				Status: status.CustomRoleStatusOK,
			},
		}
		ctx := workflow.NewContext(zap.S(), []status.Condition{}, nil)

		assert.Equal(
			t,
			workflow.OK(),
			syncCustomRolesStatus(ctx, desired, created, updated, deleted),
		)

		option := ctx.StatusOptions()[0].(status.AtlasProjectStatusOption)
		projectStatus := status.AtlasProjectStatus{}
		option(&projectStatus)
		assert.Equal(
			t,
			[]status.CustomRole{
				{
					Name:   "cr-1",
					Status: status.CustomRoleStatusOK,
				},
				{
					Name:   "cr-3",
					Status: status.CustomRoleStatusOK,
				},
				{
					Name:   "cr-4",
					Status: status.CustomRoleStatusOK,
				},
			},
			projectStatus.CustomRoles,
		)
	})

	t.Run("sync status when a operation fails", func(t *testing.T) {
		desired := []mdbv1.CustomRole{
			{
				Name: "cr-1",
			},
			{
				Name: "cr-3",
				InheritedRoles: []mdbv1.Role{
					{
						Name:     "admin",
						Database: "test",
					},
				},
			},
			{
				Name: "cr-4",
			},
		}
		created := map[string]status.CustomRole{
			"cr-4": {
				Name:   "cr-4",
				Status: status.CustomRoleStatusOK,
			},
		}
		updated := map[string]status.CustomRole{
			"cr-3": {
				Name:   "cr-3",
				Status: status.CustomRoleStatusFailed,
				Error:  "server failed",
			},
		}
		deleted := map[string]status.CustomRole{
			"cr-2": {
				Name:   "cr-2",
				Status: status.CustomRoleStatusOK,
			},
		}
		ctx := workflow.NewContext(zap.S(), []status.Condition{}, nil)

		assert.Equal(
			t,
			workflow.Terminate(workflow.ProjectCustomRolesReady, "failed to apply changes to custom roles: server failed"),
			syncCustomRolesStatus(ctx, desired, created, updated, deleted),
		)

		option := ctx.StatusOptions()[0].(status.AtlasProjectStatusOption)
		projectStatus := status.AtlasProjectStatus{}
		option(&projectStatus)
		assert.Equal(
			t,
			[]status.CustomRole{
				{
					Name:   "cr-1",
					Status: status.CustomRoleStatusOK,
				},
				{
					Name:   "cr-3",
					Status: status.CustomRoleStatusFailed,
					Error:  "server failed",
				},
				{
					Name:   "cr-4",
					Status: status.CustomRoleStatusOK,
				},
			},
			projectStatus.CustomRoles,
		)
	})
}

func TestCanCustomRolesReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(&workflowCtx, false, &mdbv1.AtlasProject{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)
		assert.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		assert.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.EqualError(t, err, "failed to retrieve data")
		assert.False(t, result)
	})

	t.Run("should return true when return nil from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return nil, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when return empty list from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return &[]mongodbatlas.CustomDBRole{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return &[]mongodbatlas.CustomDBRole{
						{
							RoleName:       "testRole1",
							InheritedRoles: nil,
							Actions: []mongodbatlas.Action{
								{
									Action: "INSERT",
									Resources: []mongodbatlas.Resource{
										{
											DB:         toptr.MakePtr("testDB"),
											Collection: toptr.MakePtr("testCollection"),
										},
									},
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CustomRoles: []mdbv1.CustomRole{
					{
						Name:           "testRole",
						InheritedRoles: nil,
						Actions: []mdbv1.Action{
							{
								Name: "INSERT",
								Resources: []mdbv1.Resource{
									{
										Database:   toptr.MakePtr("testDB"),
										Collection: toptr.MakePtr("testCollection"),
									},
								},
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"customRoles":[{"name":"testRole1","actions":[{"name":"INSERT","resources":[{"database":"testDB","collection":"testCollection"}]}]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return &[]mongodbatlas.CustomDBRole{
						{
							RoleName:       "testRole",
							InheritedRoles: nil,
							Actions: []mongodbatlas.Action{
								{
									Action: "INSERT",
									Resources: []mongodbatlas.Resource{
										{
											DB:         toptr.MakePtr("testDB"),
											Collection: toptr.MakePtr("testCollection"),
										},
									},
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CustomRoles: []mdbv1.CustomRole{
					{
						Name:           "testRole",
						InheritedRoles: nil,
						Actions: []mdbv1.Action{
							{
								Name: "INSERT",
								Resources: []mdbv1.Resource{
									{
										Database:   toptr.MakePtr("testDB"),
										Collection: toptr.MakePtr("testCollection"),
									},
								},
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"customRoles":[{"name":"testRole1","actions":[{"name":"INSERT","resources":[{"database":"testDB","collection":"testCollection"}]}]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when unable to reconcile custom roles", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return &[]mongodbatlas.CustomDBRole{
						{
							RoleName:       "testRole",
							InheritedRoles: nil,
							Actions: []mongodbatlas.Action{
								{
									Action: "INSERT",
									Resources: []mongodbatlas.Resource{
										{
											Cluster:    toptr.MakePtr(false),
											DB:         toptr.MakePtr("testDB"),
											Collection: toptr.MakePtr("testCollection"),
										},
									},
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CustomRoles: []mdbv1.CustomRole{
					{
						Name:           "testRole2",
						InheritedRoles: nil,
						Actions: []mdbv1.Action{
							{
								Name: "INSERT",
								Resources: []mdbv1.Resource{
									{
										Cluster:    toptr.MakePtr(false),
										Database:   toptr.MakePtr("testDB"),
										Collection: toptr.MakePtr("testCollection"),
									},
								},
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"customRoles":[{"name":"testRole1","actions":[{"name":"INSERT","resources":[{"cluster":false,"database":"testDB","collection":"testCollection"}]}]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canCustomRolesReconcile(workflowCtx, true, akoProject)

		assert.NoError(t, err)
		assert.False(t, result)
	})
}

func TestEnsureCustomRoles(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result := ensureCustomRoles(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			CustomDBRoles: &atlas.CustomRolesClientMock{
				ListFunc: func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
					return &[]mongodbatlas.CustomDBRole{
						{
							RoleName:       "testRole",
							InheritedRoles: nil,
							Actions: []mongodbatlas.Action{
								{
									Action: "INSERT",
									Resources: []mongodbatlas.Resource{
										{
											Cluster:    toptr.MakePtr(false),
											DB:         toptr.MakePtr("testDB"),
											Collection: toptr.MakePtr("testCollection"),
										},
									},
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				CustomRoles: []mdbv1.CustomRole{
					{
						Name:           "testRole2",
						InheritedRoles: nil,
						Actions: []mdbv1.Action{
							{
								Name: "INSERT",
								Resources: []mdbv1.Resource{
									{
										Cluster:    toptr.MakePtr(false),
										Database:   toptr.MakePtr("testDB"),
										Collection: toptr.MakePtr("testCollection"),
									},
								},
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"customRoles":[{"name":"testRole1","actions":[{"name":"INSERT","resources":[{"cluster":false,"database":"testDB","collection":"testCollection"}]}]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result := ensureCustomRoles(workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
