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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

// ensureRegionalizedPrivateEndpointMode ensures that if the AtlasProject spec
// defines a regionalized private endpoint setting, it is reflected in Atlas.
func (r *AtlasProjectReconciler) ensureRegionalizedPrivateEndpointMode(workflowCtx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.DeprecatedResult {
	if atlasProject.Spec.RegionalizedPrivateEndpoint == nil {
		workflowCtx.UnsetCondition(api.RegionalizedPrivateEndpointReadyType)
		return workflow.OK()
	}

	expectedMode := atlasProject.Spec.RegionalizedPrivateEndpoint.Enabled

	peApi := privateendpoint.NewPrivateEndpointAPI(workflowCtx.SdkClientSet.SdkClient20250312002.PrivateEndpointServicesApi)
	currentMode, err := peApi.GetRegionalizedPrivateEndpointSetting(workflowCtx.Context, atlasProject.ID())
	if err != nil {
		result := workflow.Terminate(workflow.ProjectRegionalizedEndpointModeIsNotReadyInAtlas, err)
		workflowCtx.SetConditionFromResult(api.RegionalizedPrivateEndpointReadyType, result)
		return result
	}

	if currentMode != expectedMode {
		if _, err := peApi.ToggleRegionalizedPrivateEndpointSetting(workflowCtx.Context, atlasProject.ID(), expectedMode); err != nil {
			result := workflow.Terminate(workflow.ProjectRegionalizedEndpointModeIsNotReadyInAtlas, err)
			workflowCtx.SetConditionFromResult(api.RegionalizedPrivateEndpointReadyType, result)
			return result
		}
	}

	workflowCtx.SetConditionTrue(api.RegionalizedPrivateEndpointReadyType)
	return workflow.OK()
}
