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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func TestProjectSettingsReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canProjectSettingsReconcile(context.TODO(), mongodbatlas.Client{}, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canProjectSettingsReconcile(context.TODO(), mongodbatlas.Client{}, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canProjectSettingsReconcile(context.TODO(), atlasClient, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when configuration is empty in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return nil, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canProjectSettingsReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
						IsDataExplorerEnabled:                       toptr.MakePtr(true),
						IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
						IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Settings: &mdbv1.ProjectSettings{
					IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
					IsDataExplorerEnabled:                       toptr.MakePtr(true),
					IsExtendedStorageSizesEnabled:               toptr.MakePtr(true),
					IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
					IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
					IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{
"settings": {
	"isCollectDatabaseSpecificsStatisticsEnabled": true,
	"isDataExplorerEnabled": true,
	"isPerformanceAdvisorEnabled": true,
	"isRealtimePerformancePanelEnabled": true,
	"isSchemaAdvisorEnabled": true
}}`,
			},
		)
		result, err := canProjectSettingsReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
						IsDataExplorerEnabled:                       toptr.MakePtr(true),
						IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
						IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Settings: &mdbv1.ProjectSettings{
					IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
					IsDataExplorerEnabled:                       toptr.MakePtr(true),
					IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
					IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
					IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
					IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{
"settings": {
	"isCollectDatabaseSpecificsStatisticsEnabled": true,
	"isDataExplorerEnabled": true,
	"isExtendedStorageSizesEnabled": true,
	"isPerformanceAdvisorEnabled": true,
	"isRealtimePerformancePanelEnabled": true,
	"isSchemaAdvisorEnabled": true
}}`,
			},
		)
		result, err := canProjectSettingsReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile Project Settings", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
						IsDataExplorerEnabled:                       toptr.MakePtr(true),
						IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
						IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Settings: &mdbv1.ProjectSettings{
					IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
					IsDataExplorerEnabled:                       toptr.MakePtr(false),
					IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
					IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
					IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
					IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{
"settings": {
	"isCollectDatabaseSpecificsStatisticsEnabled": true,
	"isDataExplorerEnabled": true,
	"isExtendedStorageSizesEnabled": true,
	"isPerformanceAdvisorEnabled": true,
	"isRealtimePerformancePanelEnabled": true,
	"isSchemaAdvisorEnabled": true
}}`,
			},
		)
		result, err := canProjectSettingsReconcile(context.TODO(), atlasClient, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureProjectSettings(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProjectSettings(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectSettingsFunc: func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
					return &mongodbatlas.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
						IsDataExplorerEnabled:                       toptr.MakePtr(true),
						IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
						IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Settings: &mdbv1.ProjectSettings{
					IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
					IsDataExplorerEnabled:                       toptr.MakePtr(false),
					IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
					IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
					IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
					IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{
"settings": {
	"isCollectDatabaseSpecificsStatisticsEnabled": true,
	"isDataExplorerEnabled": true,
	"isExtendedStorageSizesEnabled": true,
	"isPerformanceAdvisorEnabled": true,
	"isRealtimePerformancePanelEnabled": true,
	"isSchemaAdvisorEnabled": true
}}`,
			},
		)
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensureProjectSettings(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}

func TestAreSettingsInSync(t *testing.T) {
	atlas := &mdbv1.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
		IsDataExplorerEnabled:                       toptr.MakePtr(true),
		IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
		IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
		IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
	}
	spec := &mdbv1.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
		IsDataExplorerEnabled:                       toptr.MakePtr(true),
	}

	areEqual := areSettingsInSync(atlas, spec)
	assert.True(t, areEqual, "Only fields which are set should be compared")

	spec.IsPerformanceAdvisorEnabled = toptr.MakePtr(false)
	areEqual = areSettingsInSync(atlas, spec)
	assert.False(t, areEqual, "Field values should be the same ")
}
