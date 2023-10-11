package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func TestValidateMaintenanceWindow(t *testing.T) {
	testCases := []struct {
		in    project.MaintenanceWindow
		valid bool
	}{
		// All fields empty, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: false,
				Defer:     false,
				AutoDefer: false,
			},
			valid: false,
		},

		// Only dayOfWeek specified, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 1, // Sunday
				HourOfDay: 0, // Will default to midnight
				StartASAP: false,
				Defer:     false,
				AutoDefer: false,
			},
			valid: true,
		},

		// Specify window and autoDefer, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: false,
				Defer:     false,
				AutoDefer: true,
			},
			valid: true,
		},

		// Specify window, autoDefer, and startASAP, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: true,
				Defer:     false,
				AutoDefer: true,
			},
			valid: true,
		},

		// Specify window, autoDefer, and defer, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: false,
				Defer:     true,
				AutoDefer: true,
			},
			valid: true,
		},

		// StartASAP only, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: true,
				Defer:     false,
				AutoDefer: false,
			},
			valid: false,
		},

		// AutoDefer only, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: false,
				Defer:     false,
				AutoDefer: true,
			},
			valid: false,
		},

		// Both startASAP and Defer specified, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: true,
				Defer:     true,
				AutoDefer: true,
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			err := validateMaintenanceWindow(testCase.in)
			if !testCase.valid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCanMaintenanceWindowReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when configuration is empty in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return &mongodbatlas.MaintenanceWindow{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return &mongodbatlas.MaintenanceWindow{
						DayOfWeek: 1,
						HourOfDay: toptr.MakePtr(1),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				MaintenanceWindow: project.MaintenanceWindow{
					DayOfWeek: 7,
					HourOfDay: 20,
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"maintenanceWindow\":{\"dayOfWeek\":1,\"hourOfDay\":1}}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return &mongodbatlas.MaintenanceWindow{
						DayOfWeek: 1,
						HourOfDay: toptr.MakePtr(1),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				MaintenanceWindow: project.MaintenanceWindow{
					DayOfWeek: 1,
					HourOfDay: 1,
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"maintenanceWindow\":{\"dayOfWeek\":7,\"hourOfDay\":20}}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile IP Access List", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return &mongodbatlas.MaintenanceWindow{
						DayOfWeek: 1,
						HourOfDay: toptr.MakePtr(1),
						StartASAP: toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				MaintenanceWindow: project.MaintenanceWindow{
					DayOfWeek: 7,
					HourOfDay: 20,
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"maintenanceWindow\":{\"dayOfWeek\":1,\"hourOfDay\":1}}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canMaintenanceWindowReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureMaintenanceWindow(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result := ensureMaintenanceWindow(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			MaintenanceWindows: &atlas.MaintenanceWindowClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
					return &mongodbatlas.MaintenanceWindow{
						DayOfWeek:            1,
						HourOfDay:            toptr.MakePtr(1),
						StartASAP:            toptr.MakePtr(true),
						AutoDeferOnceEnabled: toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				MaintenanceWindow: project.MaintenanceWindow{
					DayOfWeek: 1,
					HourOfDay: 1,
					StartASAP: true,
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"maintenanceWindow\":{\"dayOfWeek\":1,\"hourOfDay\":20,\"startASAP\":true,\"autoDefer\":true}}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result := ensureMaintenanceWindow(workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
