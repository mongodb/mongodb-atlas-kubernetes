package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
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
	secret := &corev1.Secret{}
	if err := r.Client.Get(context.Background(), *dbUser.PasswordSecretObjectKey(), secret); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	currentPasswordResourceVersion := secret.ResourceVersion

	if err = validateScopes(ctx, project.ID(), dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserInvalidSpec, err.Error())
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
			ctx.EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(currentPasswordResourceVersion))

			ctx.Log.Infow("Created Atlas Database User", "name", dbUser.Spec.Username)
			return retryAfterUpdate
		} else {
			return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		}
	}
	// Update if the spec has changed
	if shouldUpdate, err := shouldUpdate(ctx.Log, u, dbUser, currentPasswordResourceVersion); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	} else if shouldUpdate {
		_, _, err = ctx.Client.DatabaseUsers.Update(context.Background(), project.ID(), dbUser.Spec.Username, apiUser)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserNotUpdatedInAtlas, err.Error())
		}
		// Update the status password resource version so that next time no API update call happened
		ctx.EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(currentPasswordResourceVersion))

		ctx.Log.Infow("Updated Atlas Database User", "name", dbUser.Spec.Username)
		// after the successful update we'll retry reconciliation so that clusters had a chance to start working
		return retryAfterUpdate
	}

	if result := checkClustersHaveReachedGoalState(ctx, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if result := createOrUpdateConnectionSecrets(ctx, r.Client, project, dbUser); !result.IsOk() {
		return result
	}

	// We mark the status.Username only when everything is finished including connection secrets
	ctx.EnsureStatusOption(status.AtlasDatabaseUserNameOption(dbUser.Spec.Username))

	return workflow.OK()
}

func validateScopes(ctx *workflow.Context, projectID string, user mdbv1.AtlasDatabaseUser) error {
	for _, s := range user.GetScopes(mdbv1.ClusterScopeType) {
		var apiError *mongodbatlas.ErrorResponse
		_, _, err := ctx.Client.Clusters.Get(context.Background(), projectID, s)
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.ClusterNotFound {
			return fmt.Errorf(`"scopes" field references cluster named "%s" but such cluster doesn't exist in Atlas'`, s)
		}
	}
	return nil
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
		ready, err := clusterIsReady(ctx.Client, projectID, c)
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

func clusterIsReady(client mongodbatlas.Client, projectID, clusterName string) (bool, error) {
	status, _, err := client.Clusters.Status(context.Background(), projectID, clusterName)
	if err != nil {
		return false, err
	}
	return status.ChangeStatus == mongodbatlas.ChangeStatusApplied, nil
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

func shouldUpdate(log *zap.SugaredLogger, atlasSpec *mongodbatlas.DatabaseUser, operatorDBUser mdbv1.AtlasDatabaseUser, currentPasswordResourceVersion string) (bool, error) {
	matches, err := userMatchesSpec(log, atlasSpec, operatorDBUser.Spec)
	if err != nil {
		return false, err
	}
	if !matches {
		return true, nil
	}
	// We need to check if the password has changed since the last time
	passwordsChanged := operatorDBUser.Status.PasswordVersion != currentPasswordResourceVersion
	if passwordsChanged {
		log.Debug("Database User password has changed - making the request to Atlas")
	}
	return passwordsChanged, nil
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
