package atlasproject

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.mongodb.org/atlas/mongodbatlas"
)

// ensureProjectExists creates the project if it doesn't exist yet. Returns the project ID
func ensureProjectExists(ctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject) (string, workflow.Result) {
	client, err := atlas.Client(connection, ctx.Log)
	if err != nil {
		return "", workflow.Terminate(workflow.Internal, err.Error())
	}
	// Try to find the project
	p, _, err := client.Projects.GetOneProjectByName(context.Background(), project.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
	}
	if p.ID != "" {
		ctx.Log.Debugw("Found Atlas Project", "id", p.ID)
		return p.ID, workflow.OK()
	}

	// Otherwise try to create it
	p = &mongodbatlas.Project{
		OrgID: connection.OrgID,
		Name:  project.Spec.Name,
	}

	if p, _, err = client.Projects.Create(context.Background(), p); err != nil {
		return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
	}
	ctx.Log.Infow("Created Atlas Project", "name", project.Spec.Name, "id", p.ID)

	return p.ID, workflow.OK()
}
