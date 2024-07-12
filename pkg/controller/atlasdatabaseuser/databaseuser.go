package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, project akov2.AtlasProject, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	apiUser, err := dbUser.ToAtlas(ctx.Context, r.Client)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if result := checkUserExpired(ctx.Context, ctx.Log, r.Client, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if err = validateScopes(ctx, project.ID(), dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserInvalidSpec, err.Error())
	}

	if result := performUpdateInAtlas(ctx, r.Client, project, dbUser, apiUser); !result.IsOk() {
		return result
	}

	if result := checkDeploymentsHaveReachedGoalState(ctx, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if result := connectionsecret.CreateOrUpdateConnectionSecrets(ctx, r.Client, r.EventRecorder, project, dbUser); !result.IsOk() {
		return result
	}

	// We need to remove the old Atlas User right after all the connection secrets are ensured if username has changed.
	if result := handleUserNameChange(ctx, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	// We mark the status.Username only when everything is finished including connection secrets
	ctx.EnsureStatusOption(status.AtlasDatabaseUserNameOption(dbUser.Spec.Username))

	return workflow.OK()
}

func handleUserNameChange(ctx *workflow.Context, projectID string, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	if dbUser.Spec.Username != dbUser.Status.UserName && dbUser.Status.UserName != "" {
		ctx.Log.Infow("'spec.username' has changed - removing the old user from Atlas", "newUserName", dbUser.Spec.Username, "oldUserName", dbUser.Status.UserName)

		deleteAttempts := 3
		for i := 1; i <= deleteAttempts; i++ {
			_, err := ctx.Client.DatabaseUsers.Delete(ctx.Context, dbUser.Spec.DatabaseName, projectID, dbUser.Status.UserName)
			if err == nil {
				break
			}

			// There may be some rare errors due to the databaseName change or maybe the user has already been removed - this
			// is not-critical (the stale connection secret has already been removed) and we shouldn't retry to avoid infinite retries
			ctx.Log.Errorf("Failed to remove user %s from Atlas (attempt %d/%d): %s", dbUser.Status.UserName, i, deleteAttempts, err)
		}
	}
	return workflow.OK()
}

func checkUserExpired(ctx context.Context, log *zap.SugaredLogger, k8sClient client.Client, projectID string, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	if dbUser.Spec.DeleteAfterDate == "" {
		return workflow.OK()
	}

	deleteAfter, err := timeutil.ParseISO8601(dbUser.Spec.DeleteAfterDate)
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserInvalidSpec, err.Error()).WithoutRetry()
	}
	if deleteAfter.Before(time.Now()) {
		if err = connectionsecret.RemoveStaleSecretsByUserName(ctx, k8sClient, projectID, dbUser.Spec.Username, dbUser, log); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		return workflow.Terminate(workflow.DatabaseUserExpired, "The database user is expired and has been removed from Atlas").WithoutRetry()
	}
	return workflow.OK()
}

func performUpdateInAtlas(ctx *workflow.Context, k8sClient client.Client, project akov2.AtlasProject, dbUser akov2.AtlasDatabaseUser, apiUser *mongodbatlas.DatabaseUser) workflow.Result {
	log := ctx.Log

	secret := &corev1.Secret{}
	passwordKey := dbUser.PasswordSecretObjectKey()
	var currentPasswordResourceVersion string
	if passwordKey != nil {
		if err := k8sClient.Get(ctx.Context, *passwordKey, secret); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		currentPasswordResourceVersion = secret.ResourceVersion
	}

	retryAfterUpdate := workflow.InProgress(workflow.DatabaseUserDeploymentAppliedChanges, "Clusters are scheduled to handle database users updates")

	// Try to find the user
	u, _, err := ctx.Client.DatabaseUsers.Get(ctx.Context, dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
			log.Debugw("User doesn't exist. Create new user", "apiUser", apiUser)
			if _, _, err = ctx.Client.DatabaseUsers.Create(ctx.Context, project.ID(), apiUser); err != nil {
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
		_, _, err = ctx.Client.DatabaseUsers.Update(ctx.Context, project.ID(), dbUser.Spec.Username, apiUser)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserNotUpdatedInAtlas, err.Error())
		}
		// Update the status password resource version so that next time no API update call happened
		ctx.EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(currentPasswordResourceVersion))

		ctx.Log.Infow("Updated Atlas Database User", "name", dbUser.Spec.Username)
		// after the successful update we'll retry reconciliation so that deployments had a chance to start working
		return retryAfterUpdate
	}

	return workflow.OK()
}

func validateScopes(ctx *workflow.Context, projectID string, user akov2.AtlasDatabaseUser) error {
	for _, s := range user.GetScopes(akov2.DeploymentScopeType) {
		var apiError *mongodbatlas.ErrorResponse
		_, _, advancedErr := ctx.Client.AdvancedClusters.Get(ctx.Context, projectID, s)
		if errors.As(advancedErr, &apiError) && apiError.ErrorCode == atlas.ClusterNotFound {
			return fmt.Errorf(`"scopes" field references deployment named "%s" but such deployment doesn't exist in Atlas'`, s)
		}
	}
	return nil
}

func checkDeploymentsHaveReachedGoalState(ctx *workflow.Context, projectID string, user akov2.AtlasDatabaseUser) workflow.Result {
	allDeploymentNames, err := atlasdeployment.GetAllDeploymentNames(ctx.Context, ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	var deploymentsToCheck []string
	if user.Spec.Scopes != nil {
		deploymentsToCheck = filterScopeDeployments(user, allDeploymentNames)
	} else {
		// otherwise we just take all the existing deployments
		deploymentsToCheck = allDeploymentNames
	}

	readyDeployments := 0
	for _, c := range deploymentsToCheck {
		ready, err := deploymentIsReady(ctx.Context, ctx.Client, projectID, c)
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		if ready {
			readyDeployments++
		}
	}
	msg := fmt.Sprintf("%d out of %d deployments have applied database user changes", readyDeployments, len(deploymentsToCheck))
	ctx.Log.Debugf(msg)

	if readyDeployments < len(deploymentsToCheck) {
		return workflow.InProgress(workflow.DatabaseUserDeploymentAppliedChanges, msg)
	}

	return workflow.OK()
}

func deploymentIsReady(ctx context.Context, client *mongodbatlas.Client, projectID, deploymentName string) (bool, error) {
	resourceStatus, _, err := client.Clusters.Status(ctx, projectID, deploymentName)
	if err != nil {
		return false, err
	}
	return resourceStatus.ChangeStatus == mongodbatlas.ChangeStatusApplied, nil
}

func filterScopeDeployments(user akov2.AtlasDatabaseUser, allDeploymentsInProject []string) []string {
	scopeDeployments := user.GetScopes(akov2.DeploymentScopeType)
	var deploymentsToCheck []string
	if len(scopeDeployments) > 0 {
		// filtering the scope deployments by the ones existing in Atlas
		for _, scopeDep := range scopeDeployments {
			for _, projectDep := range allDeploymentsInProject {
				if projectDep == scopeDep {
					deploymentsToCheck = append(deploymentsToCheck, scopeDep)
					break
				}
			}
		}
	}
	return deploymentsToCheck
}

func shouldUpdate(log *zap.SugaredLogger, atlasSpec *mongodbatlas.DatabaseUser, operatorDBUser akov2.AtlasDatabaseUser, currentPasswordResourceVersion string) (bool, error) {
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

// TODO move to a separate utils (reuse from deployments)
func userMatchesSpec(log *zap.SugaredLogger, atlasSpec *mongodbatlas.DatabaseUser, operatorSpec akov2.AtlasDatabaseUserSpec) (bool, error) {
	userMerged := mongodbatlas.DatabaseUser{}
	if err := compat.JSONCopy(&userMerged, atlasSpec); err != nil {
		return false, err
	}

	if err := compat.JSONCopy(&userMerged, operatorSpec); err != nil {
		return false, err
	}

	// performing some normalization of dates
	if atlasSpec.DeleteAfterDate != "" {
		atlasDeleteDate, err := timeutil.ParseISO8601(atlasSpec.DeleteAfterDate)
		if err != nil {
			return false, err
		}
		atlasSpec.DeleteAfterDate = timeutil.FormatISO8601(atlasDeleteDate)
	}
	if operatorSpec.DeleteAfterDate != "" {
		operatorDeleteDate, err := timeutil.ParseISO8601(operatorSpec.DeleteAfterDate)
		if err != nil {
			return false, err
		}
		userMerged.DeleteAfterDate = timeutil.FormatISO8601(operatorDeleteDate)
	}
	d := cmp.Diff(*atlasSpec, userMerged, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Users differs from spec: %s", d)
	}

	return d == "", nil
}
