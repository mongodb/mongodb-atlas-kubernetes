package atlasproject

import (
	"errors"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/maintenancewindow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

// ensureMaintenanceWindow ensures that the state of the Atlas Maintenance Window matches the
// state of the Maintenance Window specified in the project CR. If a Maintenance Window exists
// in Atlas but is not specified in the CR, it is deleted.
func (r *AtlasProjectReconciler) ensureMaintenanceWindow(workflowCtx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.Result {
	if isEmptyWindow(atlasProject.Spec.MaintenanceWindow) {
		if condition, found := workflowCtx.GetCondition(api.MaintenanceWindowReadyType); found {
			workflowCtx.Log.Debugw("Window is empty, deleting in Atlas")
			if err := r.maintenanceService.Reset(workflowCtx.Context, atlasProject.ID()); err != nil {
				result := workflow.Terminate(workflow.ProjectWindowNotDeletedInAtlas, err.Error())
				workflowCtx.SetConditionFromResult(condition.Type, result)
				return result
			}
			workflowCtx.UnsetCondition(condition.Type)
		}

		return workflow.OK()
	}

	if result := r.syncAtlasWithSpec(workflowCtx, atlasProject.ID(), atlasProject.Spec.MaintenanceWindow); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.MaintenanceWindowReadyType, result)
		return result
	}

	workflowCtx.SetConditionTrue(api.MaintenanceWindowReadyType)
	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAtlasWithSpec(ctx *workflow.Context, projectID string, windowSpec project.MaintenanceWindow) workflow.Result {
	ctx.Log.Debugw("Validate the maintenance window")
	if err := validateMaintenanceWindow(windowSpec); err != nil {
		return workflow.Terminate(workflow.ProjectWindowInvalid, err.Error())
	}

	ctx.Log.Debugw("Checking if window needs update")
	windowInAtlas, err := r.maintenanceService.Get(ctx.Context, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotObtainedFromAtlas, err.Error())
	}

	windowInAKO := maintenancewindow.NewMaintenanceWindow(&windowSpec)

	if daysOrHoursAreDifferent(*windowInAtlas, *windowInAKO) {
		ctx.Log.Debugw("Creating or updating window")
		// We set startASAP to false because the operator takes care of calling the API a second time if both
		// startASAP and the new maintenance time-slots are defined
		if err = r.maintenanceService.Update(ctx.Context, projectID, windowInAKO.WithStartASAP(false)); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotCreatedInAtlas, err.Error())
		}
	} else if windowInAtlas.AutoDefer != windowInAKO.AutoDefer {
		// If autoDefer flag is different in Atlas, and we haven't updated the window previously, we toggle the flag
		ctx.Log.Debugw("Toggling autoDefer")
		if err = r.maintenanceService.ToggleAutoDefer(ctx.Context, projectID); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotAutoDeferredInAtlas, err.Error())
		}
	}

	if windowSpec.StartASAP {
		ctx.Log.Debugw("Starting maintenance ASAP")
		// To avoid any unexpected behavior, we send a request to the API containing only the StartASAP flag,
		// although the API should ignore other fields in that case
		if err = r.maintenanceService.Update(ctx.Context, projectID,
			maintenancewindow.NewMaintenanceWindow(&project.MaintenanceWindow{StartASAP: true})); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotCreatedInAtlas, err.Error())
		}
		return workflow.OK()
	}

	if windowSpec.Defer {
		ctx.Log.Debugw("Deferring scheduled maintenance")
		if err = r.maintenanceService.Defer(ctx.Context, projectID); err != nil {
			return workflow.Terminate(workflow.ProjectWindowNotDeferredInAtlas, err.Error())
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
