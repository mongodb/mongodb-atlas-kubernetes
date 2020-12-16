package atlasproject

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
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
	p := &mongodbatlas.Project{
		OrgID: connection.OrgID,
		Name:  project.Spec.Name,
	}
	if p, _, err = client.Projects.Create(context.Background(), p); err != nil {
		// Is it an API error?
		// TODO check the errorCode and return OK in case it's either 'GROUP_ALREADY_EXISTS' or 'NOT_ORG_GROUP_CREATOR'
		// (the latter could be ok if the API key has project level)
		/*		var e *mongodbatlas.ErrorResponse
				if ok := errors.As(err, &e); ok {
					switch e.errorCode {
						case atlas.GroupExistsApiErrorCode:

					}
				}

			}*/
		return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
	}
	ctx.EnsureStatusOption(status.NewIDOption(p.ID))

	return p.ID, workflow.OK()
}
