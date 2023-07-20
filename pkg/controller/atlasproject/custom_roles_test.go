package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestCalculateChanges(t *testing.T) {
	desired := []v1.CustomRole{
		{
			Name: "cr-1",
		},
		{
			Name: "cr-3",
			InheritedRoles: []v1.Role{
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
	current := []v1.CustomRole{
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
			Create: map[string]v1.CustomRole{
				"cr-4": {
					Name: "cr-4",
				},
			},
			Update: map[string]v1.CustomRole{
				"cr-3": {
					Name: "cr-3",
					InheritedRoles: []v1.Role{
						{
							Name:     "admin",
							Database: "test",
						},
					},
				},
			},
			Delete: map[string]v1.CustomRole{
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
		desired := []v1.CustomRole{
			{
				Name: "cr-1",
			},
			{
				Name: "cr-3",
				InheritedRoles: []v1.Role{
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
		ctx := workflow.NewContext(zap.S(), []status.Condition{})

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
		desired := []v1.CustomRole{
			{
				Name: "cr-1",
			},
			{
				Name: "cr-3",
				InheritedRoles: []v1.Role{
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
		ctx := workflow.NewContext(zap.S(), []status.Condition{})

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
