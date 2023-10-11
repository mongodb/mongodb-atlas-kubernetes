package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureProjectSettings(workflowCtx *workflow.Context, project *v1.AtlasProject, protected bool) (result workflow.Result) {
	canReconcile, err := canProjectSettingsReconcile(workflowCtx, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.ProjectSettingsReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.ProjectSettingsReadyType, result)

		return result
	}

	if result = syncProjectSettings(workflowCtx, project.ID(), project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.ProjectSettingsReadyType, result)
		return result
	}

	if project.Spec.Settings == nil {
		workflowCtx.UnsetCondition(status.ProjectSettingsReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(status.ProjectSettingsReadyType)
	return workflow.OK()
}

func syncProjectSettings(ctx *workflow.Context, projectID string, project *v1.AtlasProject) workflow.Result {
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

func areSettingsInSync(atlas, spec *v1.ProjectSettings) bool {
	return isOneContainedInOther(spec, atlas)
}

func patchSettings(ctx *workflow.Context, projectID string, spec *v1.ProjectSettings) error {
	specAsAtlas, err := spec.ToAtlas()
	if err != nil {
		return err
	}

	_, _, err = ctx.Client.Projects.UpdateProjectSettings(context.Background(), projectID, specAsAtlas)
	return err
}

func fetchSettings(ctx *workflow.Context, projectID string) (*v1.ProjectSettings, error) {
	data, _, err := ctx.Client.Projects.GetProjectSettings(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugw("Got Project Settings", "data", data)

	settings := v1.ProjectSettings(*data)
	return &settings, nil
}

func isOneContainedInOther(one, other *v1.ProjectSettings) bool {
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

func canProjectSettingsReconcile(workflowCtx *workflow.Context, protected bool, akoProject *v1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &v1.AtlasProjectSpec{}
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

func areSettingsEqual(operator *v1.ProjectSettings, atlas *mongodbatlas.ProjectSettings) bool {
	if operator == nil && atlas == nil {
		return true
	}

	if operator == nil {
		operator = &v1.ProjectSettings{}
	}

	if operator.IsCollectDatabaseSpecificsStatisticsEnabled == nil {
		operator.IsCollectDatabaseSpecificsStatisticsEnabled = toptr.MakePtr(true)
	}

	if operator.IsDataExplorerEnabled == nil {
		operator.IsDataExplorerEnabled = toptr.MakePtr(true)
	}

	if operator.IsExtendedStorageSizesEnabled == nil {
		operator.IsExtendedStorageSizesEnabled = toptr.MakePtr(false)
	}

	if operator.IsPerformanceAdvisorEnabled == nil {
		operator.IsPerformanceAdvisorEnabled = toptr.MakePtr(true)
	}

	if operator.IsRealtimePerformancePanelEnabled == nil {
		operator.IsRealtimePerformancePanelEnabled = toptr.MakePtr(true)
	}

	if operator.IsSchemaAdvisorEnabled == nil {
		operator.IsSchemaAdvisorEnabled = toptr.MakePtr(true)
	}

	return *operator.IsCollectDatabaseSpecificsStatisticsEnabled == *atlas.IsCollectDatabaseSpecificsStatisticsEnabled &&
		*operator.IsDataExplorerEnabled == *atlas.IsDataExplorerEnabled &&
		*operator.IsExtendedStorageSizesEnabled == *atlas.IsExtendedStorageSizesEnabled &&
		*operator.IsPerformanceAdvisorEnabled == *atlas.IsPerformanceAdvisorEnabled &&
		*operator.IsRealtimePerformancePanelEnabled == *atlas.IsRealtimePerformancePanelEnabled &&
		*operator.IsSchemaAdvisorEnabled == *atlas.IsSchemaAdvisorEnabled
}
