package atlasproject

import (
	"context"
	"reflect"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureProjectSettings(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) (result workflow.Result) {
	if result = syncProjectSettings(ctx, projectID, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ProjectSettingsReadyType, result)
		return result
	}

	if areProjectSettingsEmpty(project.Spec.Settings) {
		ctx.UnsetCondition(status.ProjectSettingsReadyType)
		return workflow.OK()
	}

	ctx.SetConditionTrue(status.ProjectSettingsReadyType)
	return workflow.OK()
}

func syncProjectSettings(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
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

func areProjectSettingsEmpty(settings *mdbv1.ProjectSettings) bool {
	return settings == nil
}

func areSettingsInSync(atlas, spec *v1.ProjectSettings) bool {
	return reflect.DeepEqual(atlas, spec)
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

	settings := v1.ProjectSettings(*data)
	return &settings, nil
}
