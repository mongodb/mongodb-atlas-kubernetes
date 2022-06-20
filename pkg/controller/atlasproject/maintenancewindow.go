package atlasproject

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureMaintenanceWindow(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	if err := validateMaintenanceWindow(project.Spec.ProjectMaintenanceWindow); err != nil {
		return workflow.Terminate(workflow.ProjectWindowInvalid, err.Error())
	}

	if result := createOrUpdateInAtlas(ctx.Client, projectID, project.Spec.ProjectMaintenanceWindow); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

func validateMaintenanceWindow(window project.MaintenanceWindow) error {
	// TODO verify we make correct checks here
	if window.StartASAP {
		if noneSpecified := isEmpty(window.DayOfWeek) && isEmpty(window.HourOfDay) && !window.AutoDeferOnceEnabled; !noneSpecified {
			return errors.New("none of 'dayOfWeek', 'hourOfDay' and 'autoDeferOnceEnabled' should be specified if" +
				" 'startASAP is true")
		}
	}

	if isEmpty(window.DayOfWeek) || isEmpty(window.HourOfDay) {
		return errors.New("both 'dayOfWeek' and 'hourOfDay' must be specified")
	}

	// Atlas will check if dayOfWeek and hourOfDay are in the bounds
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

func isEmpty(i int) bool {
	return i == 0
}
