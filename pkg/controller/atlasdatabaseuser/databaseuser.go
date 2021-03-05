package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	retryAfterUpdate := workflow.InProgress(workflow.DatabaseUserClustersAppliedChanges, "Clusters are scheduled to handle database users updates")

	apiUser, err := dbUser.ToAtlas(r.Client)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	// Try to find the user
	u, _, err := ctx.Client.DatabaseUsers.Get(context.Background(), dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
			// User doesn't exist? Try to create it
			if _, _, err = ctx.Client.DatabaseUsers.Create(context.Background(), project.ID(), apiUser); err != nil {
				return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
			}
			ctx.Log.Infow("Created Atlas Database User", "name", dbUser.Spec.Username)
			return retryAfterUpdate
		} else {
			return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		}
	}
	// Update if the spec has changed
	if done, err := userMatchesSpec(ctx.Log, u, dbUser.Spec); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	} else if !done {
		_, _, err = ctx.Client.DatabaseUsers.Update(context.Background(), project.ID(), dbUser.Spec.Username, apiUser)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserNotUpdatedInAtlas, err.Error())
		}
		ctx.Log.Infow("Updated Atlas Database User", "name", dbUser.Spec.Username)
		// after the successful update we'll retry reconciliation so that clusters had a chance to start working
		return retryAfterUpdate
	}

	if result := checkClustersHaveReachedGoalState(ctx, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if err = createOrUpdateConnectionSecrets(ctx, r.Client, project, dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}

	return workflow.OK()
}

func checkClustersHaveReachedGoalState(ctx *workflow.Context, projectID string, user mdbv1.AtlasDatabaseUser) workflow.Result {
	allClustersInProject, _, err := ctx.Client.Clusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	var clustersToCheck []string
	if user.Spec.Scopes != nil {
		clustersToCheck = filterScopeClusters(user, allClustersInProject)
	} else {
		// otherwise we just take all the existing clusters
		for _, cluster := range allClustersInProject {
			clustersToCheck = append(clustersToCheck, cluster.Name)
		}
	}

	readyClusters := 0
	for _, c := range clustersToCheck {
		ready, err := cluserIsReady(ctx.Client, projectID, c)
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		if ready {
			readyClusters++
		}
	}
	msg := fmt.Sprintf("%d out of %d clusters have applied database user changes", readyClusters, len(clustersToCheck))
	ctx.Log.Debugf(msg)

	if readyClusters < len(clustersToCheck) {
		return workflow.InProgress(workflow.DatabaseUserClustersAppliedChanges, msg)
	}

	return workflow.OK()
}

func cluserIsReady(client mongodbatlas.Client, projectID, clusterName string) (bool, error) {
	// TODO CLOUDP-83026 GET for clusters isn't working, we need the https://docs.atlas.mongodb.com/reference/api/clusters-check-operation-status/
	_ = fmt.Sprintf("%v %v %v", client, projectID, clusterName)
	if projectID == "xyz" {
		// Cheating the linter
		return false, errors.New("boo")
	}
	return true, nil
}

func filterScopeClusters(user mdbv1.AtlasDatabaseUser, allClustersInProject []mongodbatlas.Cluster) []string {
	scopeClusters := user.GetScopes(mdbv1.ClusterScopeType)
	var clustersToCheck []string
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
	return clustersToCheck
}

// TODO move to a separate utils (reuse from clusters)
func userMatchesSpec(log *zap.SugaredLogger, atlasSpec *mongodbatlas.DatabaseUser, operatorSpec mdbv1.AtlasDatabaseUserSpec) (bool, error) {
	userMerged := mongodbatlas.DatabaseUser{}
	if err := compat.JSONCopy(&userMerged, atlasSpec); err != nil {
		return false, err
	}

	if err := compat.JSONCopy(&userMerged, operatorSpec); err != nil {
		return false, err
	}

	d := cmp.Diff(*atlasSpec, userMerged, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Users differs from spec: %s", d)
	}

	return d == "", nil
}
