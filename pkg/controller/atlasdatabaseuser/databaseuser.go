package atlasdatabaseuser

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	apiUser, err := dbUser.ToAtlas(r.Client)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	// Try to find the user
	_, _, err = ctx.Client.DatabaseUsers.Get(context.Background(), dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
			// User doesn't exist? Try to create it
			if _, _, err = ctx.Client.DatabaseUsers.Create(context.Background(), project.ID(), apiUser); err != nil {
				return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
			}
			ctx.Log.Infow("Created Atlas Database User", "name", dbUser.Spec.Username)
			return workflow.OK()
		} else {
			return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		}
	}
	// Update
	_, _, err = ctx.Client.DatabaseUsers.Update(context.Background(), project.ID(), dbUser.Spec.Username, apiUser)
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserNotUpdatedInAtlas, err.Error())
	}

	return waitClustersToHandleChanges(ctx, project.ID(), dbUser)
}

func waitClustersToHandleChanges(ctx *workflow.Context, projectID string, user mdbv1.AtlasDatabaseUser) workflow.Result {
	allClustersInProject, _, err := ctx.Client.Clusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	var clustersToCheck []string
	if user.Spec.Scopes != nil {
		var scopeClusters []string
		for _, scope := range user.Spec.Scopes {
			if scope.Type == mdbv1.ClusterScopeType {
				scopeClusters = append(scopeClusters, scope.Name)
			}
		}
		if len(scopeClusters) > 0 {
			// filtering the scope clusters by the ones existing in Atlas
			for _, c := range scopeClusters {
				for _, a := range allClustersInProject {
					if a.Name == c {
						clustersToCheck = append(clustersToCheck, c)
						break
					}
				}
			}
		}
	}

	if len(clustersToCheck) == 0 {
		// if no clusters
		ctx.Log.Debug("foo")
	}

	return workflow.OK()
}
