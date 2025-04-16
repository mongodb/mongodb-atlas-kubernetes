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
	"errors"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/maintenancewindow"
)

// ensureMaintenanceWindow ensures that the state of the Atlas Maintenance Window matches the
// state of the Maintenance Window specified in the project CR. If a Maintenance Window exists
// in Atlas but is not specified in the CR, it is deleted.
func (r *AtlasProjectReconciler) ensureMaintenanceWindow(workflowCtx *workflow.Context, atlasProject *akov2.AtlasProject, maintenanceService maintenancewindow.MaintenanceWindowService) workflow.Result {
	if isEmptyWindow(atlasProject.Spec.MaintenanceWindow) {
		if condition, found := workflowCtx.GetCondition(api.MaintenanceWindowReadyType); found {
			workflowCtx.Log.Debugw("Window is empty, deleting in Atlas")
			if err := maintenanceService.Reset(workflowCtx.Context, atlasProject.ID()); err != nil {
				result := workflow.Terminate(workflow.ProjectWindowNotDeletedInAtlas, err)
				workflowCtx.SetConditionFromResult(condition.Type, result)
				return result
			}
			workflowCtx.UnsetCondition(condition.Type)
		}

		return workflow.OK()
	}

	if result := r.syncAtlasWithSpec(workflowCtx, atlasProject.ID(), atlasProject.Spec.MaintenanceWindow, maintenanceService); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.MaintenanceWindowReadyType, result)
		return result
	}

	workflowCtx.SetConditionTrue(api.MaintenanceWindowReadyType)
	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAtlasWithSpec(ctx *workflow.Context, projectID string, windowSpec project.MaintenanceWindow, maintenanceService maintenancewindow.MaintenanceWindowService) workflow.Result {
	ctx.Log.Debugw("Validate the maintenance window")
	if err := validateMaintenanceWindow(windowSpec); err != nil {
		return workflow.Terminate(workflow.ProjectWindowInvalid, err)
	}

	ctx.Log.Debugw("Checking if window needs update")
	windowInAtlas, err := maintenanceService.Get(ctx.Context, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotObtainedFromAtlas, err)
	}

	windowInAKO := maintenancewindow.NewMaintenanceWindow(&windowSpec)

	if daysOrHoursAreDifferent(*windowInAtlas, *windowInAKO) {
		ctx.Log.Debugw("Creating or updating window")
		// We set startASAP to false because the operator takes care of calling the API a second time if both
		// startASAP and the new maintenance time-slots are defined
		if err = maintenanceService.Update(ctx.Context, projectID, windowInAKO.WithStartASAP(false)); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotCreatedInAtlas, err)
		}
	} else if windowInAtlas.AutoDefer != windowInAKO.AutoDefer {
		// If autoDefer flag is different in Atlas, and we haven't updated the window previously, we toggle the flag
		ctx.Log.Debugw("Toggling autoDefer")
		if err = maintenanceService.ToggleAutoDefer(ctx.Context, projectID); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotAutoDeferredInAtlas, err)
		}
	}

	if windowSpec.StartASAP {
		ctx.Log.Debugw("Starting maintenance ASAP")
		// To avoid any unexpected behavior, we send a request to the API containing only the StartASAP flag,
		// although the API should ignore other fields in that case
		if err = maintenanceService.Update(ctx.Context, projectID,
			maintenancewindow.NewMaintenanceWindow(&project.MaintenanceWindow{StartASAP: true})); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotCreatedInAtlas, err)
		}
		return workflow.OK()
	}

	if windowSpec.Defer {
		ctx.Log.Debugw("Deferring scheduled maintenance")
		if err = maintenanceService.Defer(ctx.Context, projectID); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotDeferredInAtlas, err)
		}
		// Nothing else should be done after deferring
		return workflow.OK()
	}

	return workflow.OK()
}

func isEmpty(i int) bool {
	return i == 0
}

func isEmptyWindow(window project.MaintenanceWindow) bool {
	return isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.StartASAP && !window.Defer && !window.AutoDefer
}

func windowSpecified(window project.MaintenanceWindow) bool {
	return !isEmpty(window.DayOfWeek)
}

func maxOneFlag(window project.MaintenanceWindow) bool {
	return !(window.StartASAP && window.Defer)
}

func daysOrHoursAreDifferent(inAtlas, inAKO maintenancewindow.MaintenanceWindow) bool {
	return inAtlas.DayOfWeek != inAKO.DayOfWeek || inAtlas.HourOfDay != inAKO.HourOfDay
}

// validateMaintenanceWindow performs validation of the Maintenance Window. Note, that we intentionally don't validate
// that hour of day and day of week are in the bounds - this will be done by Atlas.
func validateMaintenanceWindow(window project.MaintenanceWindow) error {
	if windowSpecified(window) && maxOneFlag(window) {
		return nil
	}
	errorString := "projectMaintenanceWindow must respect the following constraints, or be empty : " +
		"1) dayOfWeek must be specified (hourOfDay is 0 by default, autoDeferral is false by default) " +
		"2) only one of (startASAP, defer) is true"
	return errors.New(errorString)
}
