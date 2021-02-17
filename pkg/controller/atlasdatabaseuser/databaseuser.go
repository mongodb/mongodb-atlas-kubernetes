package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, project *mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	// Try to find the user
	p, _, err := ctx.Client.DatabaseUsers.Get(context.Background(), project.ID(), dbUser.Spec.DatabaseName, dbUser.Spec.Username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.NotInGroup {
			// User doesn't exist? Try to create it
			apiUser, err := dbUser.ToAtlas(r.Client)
			if err != nil {
				return workflow.Terminate(workflow.Internal, err.Error())
			}
			if p, _, err = ctx.Client.DatabaseUsers.Create(context.Background(), project.ID(), apiUser); err != nil {
				return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
			}
			ctx.Log.Infow("Created Atlas Database User", "name", dbUser.Spec.Username)
		} else {
			return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		}
	}
	fmt.Printf("%v", p)
	return workflow.OK()
}
