// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasproject

import (
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

func ensureProjectSettings(workflowCtx *workflow.Context, project *akov2.AtlasProject) (result workflow.DeprecatedResult) {
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

func syncProjectSettings(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) workflow.DeprecatedResult {
	spec := project.Spec.Settings

	atlas, err := fetchSettings(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectSettingsReady, err)
	}

	if !areSettingsInSync(atlas, spec) {
		if err := patchSettings(ctx, projectID, spec); err != nil {
			return workflow.Terminate(workflow.ProjectSettingsReady, err)
		}
	}

	return workflow.OK()
}

func areSettingsInSync(atlas, spec *akov2.ProjectSettings) bool {
	return isOneContainedInOther(spec, atlas)
}

func patchSettings(ctx *workflow.Context, projectID string, spec *akov2.ProjectSettings) error {
	specAsAtlas := spec.ToAtlas()

	_, _, err := ctx.SdkClientSet.SdkClient20250312009.ProjectsApi.UpdateGroupSettings(ctx.Context, projectID, specAsAtlas).Execute()
	return err
}

func fetchSettings(ctx *workflow.Context, projectID string) (*akov2.ProjectSettings, error) {
	data, _, err := ctx.SdkClientSet.SdkClient20250312009.ProjectsApi.GetGroupSettings(ctx.Context, projectID).Execute()
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
