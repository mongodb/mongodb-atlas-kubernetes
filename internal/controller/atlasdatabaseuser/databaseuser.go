// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

func (r *AtlasDatabaseUserReconciler) handleDatabaseUser(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser) (ctrl.Result, error) {
	valid, err := customresource.ResourceVersionIsValid(atlasDatabaseUser)
	switch {
	case err != nil:
		return r.terminate(ctx, atlasDatabaseUser, api.ResourceVersionStatus, workflow.AtlasResourceVersionIsInvalid, true, err)
	case !valid:
		return r.terminate(
			ctx,
			atlasDatabaseUser,
			api.ResourceVersionStatus,
			workflow.AtlasResourceVersionMismatch,
			true,
			fmt.Errorf("version of the resource '%s' is higher than the operator version '%s'", atlasDatabaseUser.GetName(), version.Version),
		)
	default:
		ctx.SetConditionTrue(api.ResourceVersionStatus).
			SetConditionTrue(api.ValidationSucceeded)
	}

	if !r.AtlasProvider.IsResourceSupported(atlasDatabaseUser) {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasGovUnsupported, false, fmt.Errorf("the %T is not supported by Atlas for government", atlasDatabaseUser))
	}

	connectionConfig, err := r.ResolveConnectionConfig(ctx.Context, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}
	sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx.Context, connectionConfig.Credentials, r.Log)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}
	dbUserService := dbuser.NewAtlasUsers(sdkClientSet.SdkClient20250312009.DatabaseUsersApi)
	deploymentService := deployment.NewAtlasDeployments(sdkClientSet.SdkClient20250312009.ClustersApi, sdkClientSet.SdkClient20250312009.GlobalClustersApi, sdkClientSet.SdkClient20250312009.FlexClustersApi, r.AtlasProvider.IsCloudGov())
	atlasProject, err := r.ResolveProject(ctx.Context, sdkClientSet.SdkClient20250312009, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}

	return r.dbuLifeCycle(ctx, dbUserService, deploymentService, atlasDatabaseUser, atlasProject)
}

func (r *AtlasDatabaseUserReconciler) dbuLifeCycle(ctx *workflow.Context, dbUserService dbuser.AtlasUsersService,
	deploymentService deployment.AtlasDeploymentsService, atlasDatabaseUser *akov2.AtlasDatabaseUser,
	atlasProject *project.Project) (ctrl.Result, error) {
	databaseUserInAtlas, err := dbUserService.Get(ctx.Context, atlasDatabaseUser.Spec.DatabaseName, atlasProject.ID, atlasDatabaseUser.Spec.Username)
	if err != nil && !errors.Is(err, dbuser.ErrorNotFound) {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	expired, err := isExpired(atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserInvalidSpec, false, err)
	}
	if expired {
		err = RemoveStaleSecretsByUserName(ctx.Context, r.Client, atlasProject.ID, atlasDatabaseUser.Spec.Username, *atlasDatabaseUser, r.Log)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotDeleted, true, err)
		}

		ctx.SetConditionFromResult(api.DatabaseUserReadyType, workflow.Terminate(workflow.DatabaseUserExpired, errors.New("an expired user cannot be managed")))
		return r.unmanage(ctx, atlasProject.ID, atlasDatabaseUser)
	}

	scopesAreValid, err := r.areDeploymentScopesValid(ctx, deploymentService, atlasProject.ID, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserInvalidSpec, false, err)
	}
	if !scopesAreValid {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserInvalidSpec, false, errors.New("\"scopes\" field refer to one or more deployments that don't exist"))
	}

	dbUserExists := databaseUserInAtlas != nil
	wasDeleted := !atlasDatabaseUser.DeletionTimestamp.IsZero()

	switch {
	case !dbUserExists && !wasDeleted:
		return r.create(ctx, dbUserService, atlasProject.ID, atlasDatabaseUser)
	case dbUserExists && !wasDeleted:
		return r.update(ctx, dbUserService, deploymentService, atlasProject, atlasDatabaseUser, databaseUserInAtlas)
	case dbUserExists && wasDeleted:
		return r.delete(ctx, dbUserService, atlasProject.ID, atlasDatabaseUser)
	default:
		return r.unmanage(ctx, atlasProject.ID, atlasDatabaseUser)
	}
}

func (r *AtlasDatabaseUserReconciler) create(ctx *workflow.Context, dbUserService dbuser.AtlasUsersService,
	projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) (ctrl.Result, error) {
	userPassword, passwordVersion, err := r.readPassword(ctx.Context, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	databaseUserInAKO, err := dbuser.NewUser(atlasDatabaseUser.Spec.DeepCopy(), projectID, userPassword)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	err = dbUserService.Create(ctx.Context, databaseUserInAKO)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotCreatedInAtlas, true, err)
	}

	if wasRenamed(atlasDatabaseUser) {
		err = RemoveStaleSecretsByUserName(ctx.Context, r.Client, projectID, atlasDatabaseUser.Status.UserName, *atlasDatabaseUser, r.Log)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotDeleted, true, err)
		}

		ctx.Log.Infow("'spec.username' has changed - removing the old user from Atlas", "newUserName", atlasDatabaseUser.Spec.Username, "oldUserName", atlasDatabaseUser.Status.UserName)
		if err = r.removeOldUser(ctx.Context, dbUserService, projectID, atlasDatabaseUser); err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
		}
	}

	return r.inProgress(ctx, atlasDatabaseUser, passwordVersion, "Clusters are scheduled to handle database users updates")
}

func (r *AtlasDatabaseUserReconciler) update(ctx *workflow.Context, dbUserService dbuser.AtlasUsersService,
	deploymentService deployment.AtlasDeploymentsService, atlasProject *project.Project,
	atlasDatabaseUser *akov2.AtlasDatabaseUser, databaseUserInAtlas *dbuser.User) (ctrl.Result, error) {
	userPassword, passwordVersion, err := r.readPassword(ctx.Context, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	databaseUserInAKO, err := dbuser.NewUser(atlasDatabaseUser.Spec.DeepCopy(), atlasProject.ID, userPassword)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	if !hasChanged(databaseUserInAKO, databaseUserInAtlas, atlasDatabaseUser.Status.PasswordVersion, passwordVersion) {
		return r.readiness(ctx, deploymentService, atlasProject, atlasDatabaseUser, passwordVersion)
	}

	r.Log.Debug(dbuser.DiffSpecs(databaseUserInAKO, databaseUserInAtlas))
	err = dbUserService.Update(ctx.Context, databaseUserInAKO)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotUpdatedInAtlas, true, err)
	}

	return r.inProgress(ctx, atlasDatabaseUser, passwordVersion, "Clusters are scheduled to handle database users updates")
}

func (r *AtlasDatabaseUserReconciler) delete(ctx *workflow.Context, dbUserService dbuser.AtlasUsersService, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(atlasDatabaseUser, r.ObjectDeletionProtection) {
		r.Log.Info("Not removing Atlas database user from Atlas as per configuration")

		return r.unmanage(ctx, projectID, atlasDatabaseUser)
	}

	err := dbUserService.Delete(ctx.Context, atlasDatabaseUser.Spec.DatabaseName, projectID, atlasDatabaseUser.Spec.Username)
	if err != nil {
		if !errors.Is(err, dbuser.ErrorNotFound) {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotDeletedInAtlas, true, err)
		}

		r.Log.Info("Database user doesn't exist or is already deleted")
	}

	return r.unmanage(ctx, projectID, atlasDatabaseUser)
}

func (r *AtlasDatabaseUserReconciler) readiness(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService,
	atlasProject *project.Project, atlasDatabaseUser *akov2.AtlasDatabaseUser, passwordVersion string) (ctrl.Result, error) {
	allDeploymentNames, err := deploymentService.ListDeploymentNames(ctx.Context, atlasProject.ID)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	removedOrphanSecrets, err := ReapOrphanConnectionSecrets(
		ctx.Context, r.Client, atlasProject.ID, atlasDatabaseUser.Namespace, allDeploymentNames)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}
	if len(removedOrphanSecrets) > 0 {
		r.Log.Debugw("Removed orphan secrets bound to an non existent deployment",
			"project", atlasProject.Name, "removed", len(removedOrphanSecrets))
		for _, orphan := range removedOrphanSecrets {
			r.Log.Debugw("Removed orphan", "secret", orphan)
		}
	}

	deploymentsToCheck := allDeploymentNames
	if atlasDatabaseUser.Spec.Scopes != nil {
		deploymentsToCheck = filterScopeDeployments(atlasDatabaseUser, allDeploymentNames)
	}

	readyDeployments := 0
	for _, c := range deploymentsToCheck {
		ready, err := deploymentService.DeploymentIsReady(ctx.Context, atlasProject.ID, c)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
		}
		if ready {
			readyDeployments++
		}
	}

	if readyDeployments < len(deploymentsToCheck) {
		return r.inProgress(
			ctx,
			atlasDatabaseUser,
			passwordVersion,
			fmt.Sprintf("%d out of %d deployments have applied database user changes", readyDeployments, len(deploymentsToCheck)),
		)
	}

	// TODO refactor connectionsecret package to follow state machine approach
	result := CreateOrUpdateConnectionSecrets(ctx, r.Client, deploymentService, r.EventRecorder, atlasProject, *atlasDatabaseUser)
	if !result.IsOk() {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotCreated, true, errors.New(result.GetMessage()))
	}

	return r.ready(ctx, atlasDatabaseUser, passwordVersion)
}

func (r *AtlasDatabaseUserReconciler) readPassword(ctx context.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser) (string, string, error) {
	if atlasDatabaseUser.Spec.PasswordSecret == nil {
		return "", "", nil
	}

	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, *atlasDatabaseUser.PasswordSecretObjectKey(), secret); err != nil {
		return "", "", err
	}

	p, exist := secret.Data["password"]
	switch {
	case !exist:
		return "", "", fmt.Errorf("secret %s is invalid: it doesn't contain 'password' field", secret.Name)
	case len(p) == 0:
		return "", "", fmt.Errorf("secret %s is invalid: the 'password' field is empty", secret.Name)
	default:
		return string(p), secret.ResourceVersion, nil
	}
}

func (r *AtlasDatabaseUserReconciler) areDeploymentScopesValid(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) (bool, error) {
	for _, s := range atlasDatabaseUser.GetScopes(akov2.DeploymentScopeType) {
		exists, err := deploymentService.ClusterExists(ctx.Context, projectID, s)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}

	return true, nil
}

func (r *AtlasDatabaseUserReconciler) removeOldUser(ctx context.Context, dbUserService dbuser.AtlasUsersService, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) error {
	deleteAttempts := 3
	var err error
	for i := 1; i <= deleteAttempts; i++ {
		err = dbUserService.Delete(ctx, atlasDatabaseUser.Spec.DatabaseName, projectID, atlasDatabaseUser.Status.UserName)
		if err == nil || errors.Is(err, dbuser.ErrorNotFound) {
			return nil
		}

		// There may be some rare errors due to the databaseName change or maybe the user has already been removed - this
		// is not-critical (the stale connection secret has already been removed) and we shouldn't retry to avoid infinite retries
		r.Log.Errorf("Failed to remove user %s from Atlas (attempt %d/%d): %s", atlasDatabaseUser.Status.UserName, i, deleteAttempts, err)
	}

	return err
}

func isExpired(atlasDatabaseUser *akov2.AtlasDatabaseUser) (bool, error) {
	if atlasDatabaseUser.Spec.DeleteAfterDate == "" {
		return false, nil
	}

	deleteAfter, err := timeutil.ParseISO8601(atlasDatabaseUser.Spec.DeleteAfterDate)
	if err != nil {
		return false, err
	}

	if !deleteAfter.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func hasChanged(databaseUserInAKO, databaseUserInAtlas *dbuser.User, currentPassVersion, passVersion string) bool {
	return !dbuser.EqualSpecs(databaseUserInAKO, databaseUserInAtlas) || currentPassVersion != passVersion
}

func wasRenamed(atlasDatabaseUser *akov2.AtlasDatabaseUser) bool {
	return atlasDatabaseUser.Status.UserName != "" && atlasDatabaseUser.Spec.Username != atlasDatabaseUser.Status.UserName
}

func filterScopeDeployments(user *akov2.AtlasDatabaseUser, allDeploymentsInProject []string) []string {
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
