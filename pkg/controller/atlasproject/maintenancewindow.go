package atlasproject

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// ensureMaintenanceWindow ensures that the state of the Atlas Maintenance Window matches the
// state of the Maintenance Window specified in the project CR. If a Maintenance Window exists
// in Atlas but is not specified in the CR, it is deleted.
func ensureMaintenanceWindow(ctx *workflow.Context, projectID string, atlasProject *mdbv1.AtlasProject) workflow.Result {
	windowSpec := atlasProject.Spec.ProjectMaintenanceWindow
	if err := validateMaintenanceWindow(windowSpec); err != nil {
		return workflow.Terminate(workflow.ProjectWindowInvalid, err.Error())
	}

	ctx.Log.Debugw(fmt.Sprintf("%s%t", "Checking if projectMaintenanceWindow field is empty or undefined : ", isEmptyWindow(windowSpec)))
	if isEmptyWindow(windowSpec) {
		ctx.Log.Debugw("Deleting in Atlas")
		if result := deleteInAtlas(ctx.Client, projectID); !result.IsOk() {
			return result
		}
		return workflow.OK()
	}

	if needUpdate(windowSpec) {
		ctx.Log.Debugw("Updating in Atlas")
		// We set startASAP to false because the operator takes care of calling the API a second time if both
		// startASAP and the new maintenance timeslots are defined
		if result := createOrUpdateInAtlas(ctx.Client, projectID, windowSpec.WithStartASAP(false)); !result.IsOk() {
			return result
		}
	}

	if windowSpec.StartASAP {
		ctx.Log.Debugw("Starting maintenance ASAP")
		// To avoid any conflict, we send a request to the API containing only the StartASAP flag, although the API
		// should ignore other fields in that case
		if result := createOrUpdateInAtlas(ctx.Client, projectID, project.NewMaintenanceWindow().WithStartASAP(true)); !result.IsOk() {
			return result
		}
		// Nothing else should be done after sending a StartASAP request
		return workflow.OK()
	}

	if windowSpec.Defer {
		ctx.Log.Debugw("Deferring maintenance")
		if result := deferInAtlas(ctx.Client, projectID); !result.IsOk() {
			return result
		}
		// Nothing else should be done after deferring
		return workflow.OK()
	}

	if windowSpec.AutoDefer {
		ctx.Log.Debugw("Auto-deferring maintenance")
		if result := autoDeferInAtlas(ctx.Client, projectID); !result.IsOk() {
			return result
		}
		// Nothing else should be done after auto-deferring
		return workflow.OK()
	}
	return workflow.OK()
}

func isEmptyWindow(window project.MaintenanceWindow) bool {
	return isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.AutoDeferOnceEnabled && !window.StartASAP && notDeferred(window)
}

func needUpdate(window project.MaintenanceWindow) bool {
	return !isEmpty(window.DayOfWeek) && !isEmpty(window.HourOfDay)
}

func notDeferred(window project.MaintenanceWindow) bool {
	return !window.AutoDefer && !window.Defer
}

func noOtherFieldsThanDefer(window project.MaintenanceWindow) bool {
	return isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.AutoDeferOnceEnabled && !window.StartASAP
}

// validateMaintenanceWindow performs validation of the Maintenance Window. Note, that we intentionally don't validate
// that hour of day and day of week are in the bounds - this will be done by Atlas.
func validateMaintenanceWindow(window project.MaintenanceWindow) error {
	switch {
	case isEmptyWindow(window):
		return nil
	case needUpdate(window) && notDeferred(window) && !window.StartASAP:
		return nil
	case needUpdate(window) && notDeferred(window) && window.StartASAP && !window.AutoDeferOnceEnabled:
		return nil
	case window.StartASAP && notDeferred(window) && !window.AutoDeferOnceEnabled:
		return nil
	case window.Defer && noOtherFieldsThanDefer(window) && !window.AutoDefer:
		return nil
	case window.AutoDefer && noOtherFieldsThanDefer(window) && !window.Defer:
		return nil
	default:
		return errors.New(`
			projectMaintenanceWindow must respect one of the following constraints :
				1) both hourOfDay and dayOfWeek are specified (!= 0), deferral fields are empty,
				   only one or none of startASAP and autoDeferOnceEnabled is true
				2) startASAP is true, deferral fields are empty, autoDeferOnceEnabled is false or empty
				3) defer is true, all other fields are empty
				4) autoDefer is true, all other fields are empty
				5) all fields are empty (will delete the window if it exists)
		`)
	}
}

func isEmpty(i int) bool {
	return i == 0
}

// operatorToAtlasMaintenanceWindow converts the maintenanceWindow specified in the project CR to the format
// expected by the Atlas API.
func operatorToAtlasMaintenanceWindow(maintenanceWindow project.MaintenanceWindow) (*mongodbatlas.MaintenanceWindow, workflow.Result) {
	operatorWindow, err := maintenanceWindow.ToAtlas()
	if err != nil {
		return nil, workflow.Terminate(workflow.Internal, err.Error())
	}
	return operatorWindow, workflow.OK()
}

func createOrUpdateInAtlas(client mongodbatlas.Client, projectID string, maintenanceWindow project.MaintenanceWindow) workflow.Result {
	operatorWindow, status := operatorToAtlasMaintenanceWindow(maintenanceWindow)
	if !status.IsOk() {
		return status
	}

	if _, err := client.MaintenanceWindows.Update(context.Background(), projectID, operatorWindow); err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotCreatedInAtlas, err.Error())
	}
	return workflow.OK()
}

func deleteInAtlas(client mongodbatlas.Client, projectID string) workflow.Result {
	if _, err := client.MaintenanceWindows.Reset(context.Background(), projectID); err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotDeletedInAtlas, err.Error())
	}
	return workflow.OK()
}

func deferInAtlas(client mongodbatlas.Client, projectID string) workflow.Result {
	if _, err := client.MaintenanceWindows.Defer(context.Background(), projectID); err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotDeferredInAtlas, err.Error())
	}
	return workflow.OK()
}

func autoDeferInAtlas(client mongodbatlas.Client, projectID string) workflow.Result {
	if _, err := client.MaintenanceWindows.AutoDefer(context.Background(), projectID); err != nil {
		return workflow.Terminate(workflow.ProjectWindowNotAutoDeferredInAtlas, err.Error())
	}
	return workflow.OK()
}
