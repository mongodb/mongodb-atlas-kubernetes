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
func ensureMaintenanceWindow(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	if err := validateMaintenanceWindow(project.Spec.ProjectMaintenanceWindow); err != nil {
		return workflow.Terminate(workflow.ProjectWindowInvalid, err.Error())
	}

	ctx.Log.Debugw(fmt.Sprintf("%s%t", "Checking if Window is empty : ", isEmptyWindow(project.Spec.ProjectMaintenanceWindow)))
	if isEmptyWindow(project.Spec.ProjectMaintenanceWindow) {
		ctx.Log.Debugw("Deleting in Atlas")
		if result := deleteInAtlas(ctx.Client, projectID); !result.IsOk() {
			return result
		}
	} else {
		ctx.Log.Debugw("Updating in Atlas")
		if result := createOrUpdateInAtlas(ctx.Client, projectID, project.Spec.ProjectMaintenanceWindow); !result.IsOk() {
			return result
		}
	}

	return workflow.OK()
}

func isEmptyWindow(window project.MaintenanceWindow) bool {
	return isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.AutoDeferOnceEnabled && !window.StartASAP
}

// validateMaintenanceWindow performs validation of the Maintenance Window. Note, that we intentionally don't validate
// hour of day and day of week - this will be done by Atlas.
func validateMaintenanceWindow(window project.MaintenanceWindow) error {
	// If StartASAP is specified, it should be the only field
	if window.StartASAP {
		if noneSpecified := isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.AutoDeferOnceEnabled; !noneSpecified {
			return errors.New("none of 'dayOfWeek', 'hourOfDay' and 'autoDeferOnceEnabled' should be specified if" +
				" 'startASAP is true")
		}
	} else {
		// Query is valid if either all fields are empty, or both day and hour are specified
		// Atlas will check if dayOfWeek and hourOfDay are in the bounds
		if !isEmptyWindow(window) && (isEmpty(window.DayOfWeek) || isEmpty(window.HourOfDay)) {
			return errors.New("both 'dayOfWeek' and 'hourOfDay' must be specified")
		}
	}
	return nil
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

func isEmpty(i int) bool {
	return i == 0
}
