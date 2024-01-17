package atlasproject

import (
	"go.mongodb.org/atlas-sdk/v20231001002/admin"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

// ensureProjectExists creates the project if it doesn't exist yet. Returns the project ID
func (r *AtlasProjectReconciler) ensureProjectExists(ctx *workflow.Context, project *mdbv1.AtlasProject) (string, workflow.Result) {
	// Try to find the project
	p, _, err := ctx.SdkClient.ProjectsApi.GetProjectByName(ctx.Context, project.Spec.Name).Execute()
	if err != nil {
		ctx.Log.Infow("Error", "err", err.Error())
		if admin.IsErrorCode(err, atlas.NotInGroup) || admin.IsErrorCode(err, atlas.ResourceNotFound) {
			// Project doesn't exist? Try to create it
			p = &admin.Group{
				OrgId:                     ctx.OrgID,
				Name:                      project.Spec.Name,
				WithDefaultAlertsSettings: &project.Spec.WithDefaultAlertsSettings,
				RegionUsageRestrictions:   &project.Spec.RegionUsageRestrictions,
			}
			p, _, err = ctx.SdkClient.ProjectsApi.CreateProject(ctx.Context, p).Execute()
			if err != nil {
				return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
			}
			ctx.Log.Infow("Created Atlas Project", "name", project.Spec.Name, "id", p.GetId())
		} else {
			return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
		}
	}

	if p == nil || p.GetId() == "" {
		ctx.Log.Error("Project or its project ID are empty")
		return "", workflow.Terminate(workflow.Internal, "")
	}

	return p.GetId(), workflow.OK()
}
