package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, dus dbuser.AtlasUsersService, ds deployment.AtlasDeploymentsService, project akov2.AtlasProject, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	password, err := dbUser.ReadPassword(ctx.Context, r.Client)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	apiUser, err := dbuser.NewUser(dbUser.Spec.DeepCopy(), project.ID(), password)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if result := checkUserExpired(ctx.Context, ctx.Log, r.Client, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if err := validateScopes(ctx, ds, project.ID(), dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserInvalidSpec, err.Error())
	}

	if result := performUpdateInAtlas(ctx, r.Client, dus, project, &dbUser, apiUser); !result.IsOk() {
		return result
	}

	if result := checkDeploymentsHaveReachedGoalState(ctx, ds, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	if result := connectionsecret.CreateOrUpdateConnectionSecrets(ctx, r.Client, ds, r.EventRecorder, project, dbUser); !result.IsOk() {
		return result
	}

	// We need to remove the old Atlas User right after all the connection secrets are ensured if username has changed.
	if result := handleUserNameChange(ctx, dus, project.ID(), dbUser); !result.IsOk() {
		return result
	}

	// We mark the status.Username only when everything is finished including connection secrets
	ctx.EnsureStatusOption(status.AtlasDatabaseUserNameOption(dbUser.Spec.Username))

	return workflow.OK()
}

func handleUserNameChange(ctx *workflow.Context, dus dbuser.AtlasUsersService, projectID string, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	if dbUser.Spec.Username != dbUser.Status.UserName && dbUser.Status.UserName != "" {
		ctx.Log.Infow("'spec.username' has changed - removing the old user from Atlas", "newUserName", dbUser.Spec.Username, "oldUserName", dbUser.Status.UserName)

		deleteAttempts := 3
		var err error
		for i := 1; i <= deleteAttempts; i++ {
			err = dus.Delete(ctx.Context, dbUser.Spec.DatabaseName, projectID, dbUser.Status.UserName)
			if err == nil || errors.Is(err, dbuser.ErrorNotFound) {
				return workflow.OK()
			}

			// There may be some rare errors due to the databaseName change or maybe the user has already been removed - this
			// is not-critical (the stale connection secret has already been removed) and we shouldn't retry to avoid infinite retries
			ctx.Log.Errorf("Failed to remove user %s from Atlas (attempt %d/%d): %s", dbUser.Status.UserName, i, deleteAttempts, err)
		}
		return workflow.Terminate(workflow.Internal, err.Error())
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

func performUpdateInAtlas(ctx *workflow.Context, k8sClient client.Client, dus dbuser.AtlasUsersService, project akov2.AtlasProject, dbUser *akov2.AtlasDatabaseUser, apiUser *dbuser.User) workflow.Result {
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
	au, err := dus.Get(ctx.Context, dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
	if errors.Is(err, dbuser.ErrorNotFound) {
		log.Debugw("User doesn't exist. Create new user", "dbUser", dbUser)
		if err = dus.Create(ctx.Context, apiUser); err != nil {
			return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		}
		ctx.EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(currentPasswordResourceVersion))

		ctx.Log.Infow("Created Atlas Database User", "name", dbUser.Spec.Username)
		return retryAfterUpdate
	}
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
	}
	// Update if the spec has changed
	if shouldUpdate, err := shouldUpdate(ctx.Log, au, dbUser, currentPasswordResourceVersion); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	} else if shouldUpdate {
		err = dus.Update(ctx.Context, apiUser)
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

func validateScopes(ctx *workflow.Context, ds deployment.AtlasDeploymentsService, projectID string, user akov2.AtlasDatabaseUser) error {
	for _, s := range user.GetScopes(akov2.DeploymentScopeType) {
		exists, err := ds.ClusterExists(ctx.Context, projectID, s)
		if !exists {
			return fmt.Errorf(`"scopes" field references deployment named "%s" but such deployment doesn't exist in Atlas'`, s)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func checkDeploymentsHaveReachedGoalState(ctx *workflow.Context, ds deployment.AtlasDeploymentsService, projectID string, user akov2.AtlasDatabaseUser) workflow.Result {
	allDeploymentNames, err := ds.ListClusterNames(ctx.Context, projectID)
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
		ready, err := ds.DeploymentIsReady(ctx.Context, projectID, c)
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

func shouldUpdate(log *zap.SugaredLogger, atlasUser *dbuser.User, operatorUser *akov2.AtlasDatabaseUser, currentPasswordResourceVersion string) (bool, error) {
	userSpecCopy, err := dbuser.NewUser(operatorUser.Spec.DeepCopy(), "", "") // project and password will not be compared
	if err != nil {
		return false, fmt.Errorf("failed to convert User Spec: %w", err)
	}
	if !dbuser.EqualSpecs(userSpecCopy, atlasUser) {
		if log.Level() == zapcore.DebugLevel {
			log.Debugf("Atlas Database user differs from Spec: %v", dbuser.DiffSpecs(userSpecCopy, atlasUser))
		}
		return true, nil
	}
	// We need to check if the password has changed since the last time
	passwordsChanged := operatorUser.Status.PasswordVersion != currentPasswordResourceVersion
	if passwordsChanged {
		log.Debug("Database User password has changed - making the request to Atlas")
	}
	return passwordsChanged, nil
}
