package atlasproject

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func ensureProjectSettings(workflowCtx *workflow.Context, project *akov2.AtlasProject, protected bool) (result workflow.Result) {
	canReconcile, err := canProjectSettingsReconcile(workflowCtx, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(api.ProjectSettingsReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(api.ProjectSettingsReadyType, result)

		return result
	}

	if result = syncProjectSettings(workflowCtx, project.ID(), project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.ProjectSettingsReadyType, result)
		return result
	}

	if project.Spec.Settings == nil {
		workflowCtx.UnsetCondition(api.ProjectSettingsReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(api.ProjectSettingsReadyType)
	return workflow.OK()
}

func syncProjectSettings(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) workflow.Result {
	spec := project.Spec.Settings

	atlas, err := fetchSettings(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectSettingsReady, err.Error())
	}

	if !areSettingsInSync(atlas, spec) {
		if err := patchSettings(ctx, projectID, spec); err != nil {
			return workflow.Terminate(workflow.ProjectSettingsReady, err.Error())
		}
	}

	return workflow.OK()
}

func areSettingsInSync(atlas, spec *akov2.ProjectSettings) bool {
	return isOneContainedInOther(spec, atlas)
}

func patchSettings(ctx *workflow.Context, projectID string, spec *akov2.ProjectSettings) error {
	specAsAtlas, err := spec.ToAtlas()
	if err != nil {
		return err
	}

	_, _, err = ctx.Client.Projects.UpdateProjectSettings(ctx.Context, projectID, specAsAtlas)
	return err
}

func fetchSettings(ctx *workflow.Context, projectID string) (*akov2.ProjectSettings, error) {
	data, _, err := ctx.Client.Projects.GetProjectSettings(ctx.Context, projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugw("Got Project Settings", "data", data)

	settings := akov2.ProjectSettings(*data)
	return &settings, nil
}

func isOneContainedInOther(one, other *akov2.ProjectSettings) bool {
	if one == nil {
		return true
	}

	if other == nil {
		return false
	}

	oneVal := reflect.ValueOf(one).Elem()
	otherVal := reflect.ValueOf(other).Elem()

	for i := 0; i < oneVal.NumField(); i++ {
		if oneVal.Field(i).IsNil() {
			continue
		}

		oneBool := oneVal.Field(i).Elem().Bool()
		otherBool := otherVal.Field(i).Elem().Bool()

		if oneBool != otherBool {
			return false
		}
	}

	return true
}

func canProjectSettingsReconcile(workflowCtx *workflow.Context, protected bool, akoProject *akov2.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &akov2.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	settings, _, err := workflowCtx.Client.Projects.GetProjectSettings(workflowCtx.Context, akoProject.ID())
	if err != nil {
		return false, err
	}

	if settings == nil {
		return true, nil
	}

	return areSettingsEqual(latestConfig.Settings, settings) ||
		areSettingsEqual(akoProject.Spec.Settings, settings), nil
}

func areSettingsEqual(operator *akov2.ProjectSettings, atlas *mongodbatlas.ProjectSettings) bool {
	if operator == nil && atlas == nil {
		return true
	}

	if operator == nil {
		operator = &akov2.ProjectSettings{}
	}

	if operator.IsCollectDatabaseSpecificsStatisticsEnabled == nil {
		operator.IsCollectDatabaseSpecificsStatisticsEnabled = pointer.MakePtr(true)
	}

	if operator.IsDataExplorerEnabled == nil {
		operator.IsDataExplorerEnabled = pointer.MakePtr(true)
	}

	if operator.IsExtendedStorageSizesEnabled == nil {
		operator.IsExtendedStorageSizesEnabled = pointer.MakePtr(false)
	}

	if operator.IsPerformanceAdvisorEnabled == nil {
		operator.IsPerformanceAdvisorEnabled = pointer.MakePtr(true)
	}

	if operator.IsRealtimePerformancePanelEnabled == nil {
		operator.IsRealtimePerformancePanelEnabled = pointer.MakePtr(true)
	}

	if operator.IsSchemaAdvisorEnabled == nil {
		operator.IsSchemaAdvisorEnabled = pointer.MakePtr(true)
	}

	return *operator.IsCollectDatabaseSpecificsStatisticsEnabled == *atlas.IsCollectDatabaseSpecificsStatisticsEnabled &&
		*operator.IsDataExplorerEnabled == *atlas.IsDataExplorerEnabled &&
		*operator.IsExtendedStorageSizesEnabled == *atlas.IsExtendedStorageSizesEnabled &&
		*operator.IsPerformanceAdvisorEnabled == *atlas.IsPerformanceAdvisorEnabled &&
		*operator.IsRealtimePerformancePanelEnabled == *atlas.IsRealtimePerformancePanelEnabled &&
		*operator.IsSchemaAdvisorEnabled == *atlas.IsSchemaAdvisorEnabled
}
