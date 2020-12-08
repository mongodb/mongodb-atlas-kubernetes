package atlasproject

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
)

func ensureProjectExists(connection atlas.Connection, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	client, err := atlas.Client(connection, log)
	if err != nil {
		return err
	}
	p := &mongodbatlas.Project{
		OrgID: connection.OrgID,
		Name:  project.Spec.Name,
	}
	if _, _, err := client.Projects.Create(context.Background(), p); err != nil {
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
		return err
	}
	return nil
}
