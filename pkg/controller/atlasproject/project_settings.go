package atlasproject

import (
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func ensureProjectSettings(workflowCtx *workflow.Context, project *akov2.AtlasProject) (result workflow.Result) {
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
	specAsAtlas := spec.ToAtlas()

	_, _, err := ctx.SdkClient.ProjectsApi.UpdateProjectSettings(ctx.Context, projectID, specAsAtlas).Execute()
	return err
}

func fetchSettings(ctx *workflow.Context, projectID string) (*akov2.ProjectSettings, error) {
	data, _, err := ctx.SdkClient.ProjectsApi.GetProjectSettings(ctx.Context, projectID).Execute()
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugw("Got Project Settings", "data", data)

	return akov2.ProjectSettingsFromAtlas(data), nil
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
