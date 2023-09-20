package atlasproject

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// ensureProjectExists creates the project if it doesn't exist yet. Returns the project ID
func (r *AtlasProjectReconciler) ensureProjectExists(ctx *workflow.Context, project *mdbv1.AtlasProject) (string, workflow.Result) {
	// Try to find the project
	p, _, err := ctx.Client.Projects.GetOneProjectByName(context.Background(), project.Spec.Name)
	if err != nil {
		ctx.Log.Infow("Error", "err", err.Error())
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.NotInGroup || apiError.ErrorCode == atlas.ResourceNotFound) {
			// Project doesn't exist? Try to create it
			p = &mongodbatlas.Project{
				OrgID:                     ctx.Connection.OrgID,
				Name:                      project.Spec.Name,
				WithDefaultAlertsSettings: &project.Spec.WithDefaultAlertsSettings,
				RegionUsageRestrictions:   project.Spec.RegionUsageRestrictions,
			}
			if p, _, err = ctx.Client.Projects.Create(context.Background(), p, &mongodbatlas.CreateProjectOptions{}); err != nil {
				return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
			}
			ctx.Log.Infow("Created Atlas Project", "name", project.Spec.Name, "id", p.ID)
		} else {
			return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
		}
	}

	if p == nil || p.ID == "" {
		ctx.Log.Error("Project or its project ID are empty")
		return "", workflow.Terminate(workflow.Internal, "")
	}

	return p.ID, workflow.OK()
}
