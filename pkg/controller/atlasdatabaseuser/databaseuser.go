package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

func (r *AtlasDatabaseUserReconciler) handleDatabaseUser(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser) ctrl.Result {
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

	var atlasProject *project.Project
	if atlasDatabaseUser.Spec.ExternalProjectRef != nil {
		atlasProject, err = r.getProjectFromAtlas(ctx, atlasDatabaseUser)
	} else {
		atlasProject, err = r.getProjectFromKube(ctx, atlasDatabaseUser)
	}
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}

	return r.dbuLifeCycle(ctx, atlasDatabaseUser, atlasProject)
}

func (r *AtlasDatabaseUserReconciler) dbuLifeCycle(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser, atlasProject *project.Project) ctrl.Result {
	databaseUserInAtlas, err := r.dbUserService.Get(ctx.Context, atlasDatabaseUser.Spec.DatabaseName, atlasProject.ID, atlasDatabaseUser.Spec.Username)
	if err != nil && !errors.Is(err, dbuser.ErrorNotFound) {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	expired, err := isExpired(atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserInvalidSpec, false, err)
	}
	if expired {
		err = connectionsecret.RemoveStaleSecretsByUserName(ctx.Context, r.Client, atlasProject.ID, atlasDatabaseUser.Spec.Username, *atlasDatabaseUser, r.Log)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotDeleted, true, err)
		}

		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserExpired, false, errors.New("an expired user cannot be managed"))
	}

	scopesAreValid, err := r.areDeploymentScopesValid(ctx, atlasProject.ID, atlasDatabaseUser)
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
		return r.create(ctx, atlasProject.ID, atlasDatabaseUser)
	case dbUserExists && !wasDeleted:
		return r.update(ctx, atlasProject, atlasDatabaseUser, databaseUserInAtlas)
	case dbUserExists && wasDeleted:
		return r.delete(ctx, atlasProject.ID, atlasDatabaseUser)
	default:
		return r.unmanage(ctx, atlasProject.ID, atlasDatabaseUser)
	}
}

func (r *AtlasDatabaseUserReconciler) create(ctx *workflow.Context, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) ctrl.Result {
	if !canManageOIDC(r.FeaturePreviewOIDCAuthEnabled, atlasDatabaseUser.Spec.OIDCAuthType) {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, false, ErrOIDCNotEnabled)
	}

	userPassword, passwordVersion, err := r.readPassword(ctx.Context, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	databaseUserInAKO, err := dbuser.NewUser(atlasDatabaseUser.Spec.DeepCopy(), projectID, userPassword)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	err = r.dbUserService.Create(ctx.Context, databaseUserInAKO)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotCreatedInAtlas, true, err)
	}

	if wasRenamed(atlasDatabaseUser) {
		err = connectionsecret.RemoveStaleSecretsByUserName(ctx.Context, r.Client, projectID, atlasDatabaseUser.Status.UserName, *atlasDatabaseUser, r.Log)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotDeleted, true, err)
		}

		ctx.Log.Infow("'spec.username' has changed - removing the old user from Atlas", "newUserName", atlasDatabaseUser.Spec.Username, "oldUserName", atlasDatabaseUser.Status.UserName)
		if err = r.removeOldUser(ctx.Context, projectID, atlasDatabaseUser); err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
		}
	}

	return r.inProgress(ctx, atlasDatabaseUser, passwordVersion, "Clusters are scheduled to handle database users updates")
}

func (r *AtlasDatabaseUserReconciler) update(ctx *workflow.Context, atlasProject *project.Project, atlasDatabaseUser *akov2.AtlasDatabaseUser, databaseUserInAtlas *dbuser.User) ctrl.Result {
	if !canManageOIDC(r.FeaturePreviewOIDCAuthEnabled, atlasDatabaseUser.Spec.OIDCAuthType) {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, false, ErrOIDCNotEnabled)
	}

	userPassword, passwordVersion, err := r.readPassword(ctx.Context, atlasDatabaseUser)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	databaseUserInAKO, err := dbuser.NewUser(atlasDatabaseUser.Spec.DeepCopy(), atlasProject.ID, userPassword)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	if !hasChanged(databaseUserInAKO, databaseUserInAtlas, atlasDatabaseUser.Status.PasswordVersion, passwordVersion) {
		return r.readiness(ctx, atlasProject, atlasDatabaseUser, passwordVersion)
	}

	r.Log.Debug(dbuser.DiffSpecs(databaseUserInAKO, databaseUserInAtlas))
	err = r.dbUserService.Update(ctx.Context, databaseUserInAKO)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotUpdatedInAtlas, true, err)
	}

	return r.inProgress(ctx, atlasDatabaseUser, passwordVersion, "Clusters are scheduled to handle database users updates")
}

func (r *AtlasDatabaseUserReconciler) delete(ctx *workflow.Context, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) ctrl.Result {
	if customresource.IsResourcePolicyKeepOrDefault(atlasDatabaseUser, r.ObjectDeletionProtection) {
		r.Log.Info("Not removing Atlas database user from Atlas as per configuration")

		return r.unmanage(ctx, projectID, atlasDatabaseUser)
	}

	err := r.dbUserService.Delete(ctx.Context, atlasDatabaseUser.Spec.DatabaseName, projectID, atlasDatabaseUser.Spec.Username)
	if err != nil {
		if !errors.Is(err, dbuser.ErrorNotFound) {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserNotDeletedInAtlas, true, err)
		}

		r.Log.Info("Database user doesn't exist or is already deleted")
	}

	return r.unmanage(ctx, projectID, atlasDatabaseUser)
}

func (r *AtlasDatabaseUserReconciler) readiness(ctx *workflow.Context, atlasProject *project.Project, atlasDatabaseUser *akov2.AtlasDatabaseUser, passwordVersion string) ctrl.Result {
	allDeploymentNames, err := r.deploymentService.ListClusterNames(ctx.Context, atlasProject.ID)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	deploymentsToCheck := allDeploymentNames
	if atlasDatabaseUser.Spec.Scopes != nil {
		deploymentsToCheck = filterScopeDeployments(atlasDatabaseUser, allDeploymentNames)
	}

	readyDeployments := 0
	for _, c := range deploymentsToCheck {
		ready, err := r.deploymentService.DeploymentIsReady(ctx.Context, atlasProject.ID, c)
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
	result := connectionsecret.CreateOrUpdateConnectionSecrets(ctx, r.Client, r.deploymentService, r.EventRecorder, atlasProject, *atlasDatabaseUser)
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

func (r *AtlasDatabaseUserReconciler) areDeploymentScopesValid(ctx *workflow.Context, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) (bool, error) {
	for _, s := range atlasDatabaseUser.GetScopes(akov2.DeploymentScopeType) {
		exists, err := r.deploymentService.ClusterExists(ctx.Context, projectID, s)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}

	return true, nil
}

func (r *AtlasDatabaseUserReconciler) removeOldUser(ctx context.Context, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) error {
	deleteAttempts := 3
	var err error
	for i := 1; i <= deleteAttempts; i++ {
		err = r.dbUserService.Delete(ctx, atlasDatabaseUser.Spec.DatabaseName, projectID, atlasDatabaseUser.Status.UserName)
		if err == nil || errors.Is(err, dbuser.ErrorNotFound) {
			return nil
		}

		// There may be some rare errors due to the databaseName change or maybe the user has already been removed - this
		// is not-critical (the stale connection secret has already been removed) and we shouldn't retry to avoid infinite retries
		r.Log.Errorf("Failed to remove user %s from Atlas (attempt %d/%d): %s", atlasDatabaseUser.Status.UserName, i, deleteAttempts, err)
	}

	return err
}

func (r *AtlasDatabaseUserReconciler) getProjectFromAtlas(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser) (*project.Project, error) {
	sdkClient, _, err := r.AtlasProvider.SdkClient(
		ctx.Context,
		&client.ObjectKey{Namespace: atlasDatabaseUser.Namespace, Name: atlasDatabaseUser.Credentials().Name},
		r.Log,
	)
	if err != nil {
		return nil, err
	}

	projectService := project.NewProjectAPIService(sdkClient.ProjectsApi)
	r.dbUserService = dbuser.NewAtlasUsers(sdkClient.DatabaseUsersApi)
	r.deploymentService = deployment.NewAtlasDeployments(sdkClient.ClustersApi, sdkClient.ServerlessInstancesApi, r.AtlasProvider.IsCloudGov())

	atlasProject, err := projectService.GetProject(ctx.Context, atlasDatabaseUser.Spec.ExternalProjectRef.ID)
	if err != nil {
		return nil, err
	}

	return atlasProject, nil
}

func (r *AtlasDatabaseUserReconciler) getProjectFromKube(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser) (*project.Project, error) {
	atlasProject := &akov2.AtlasProject{}
	if err := r.Client.Get(ctx.Context, atlasDatabaseUser.AtlasProjectObjectKey(), atlasProject); err != nil {
		return nil, err
	}

	credentialsSecret, err := customresource.ComputeSecret(atlasProject, atlasDatabaseUser)
	if err != nil {
		return nil, err
	}

	sdkClient, orgID, err := r.AtlasProvider.SdkClient(ctx.Context, credentialsSecret, r.Log)
	if err != nil {
		return nil, err
	}

	r.dbUserService = dbuser.NewAtlasUsers(sdkClient.DatabaseUsersApi)
	r.deploymentService = deployment.NewAtlasDeployments(sdkClient.ClustersApi, sdkClient.ServerlessInstancesApi, r.AtlasProvider.IsCloudGov())

	return project.NewProject(atlasProject, orgID), nil
}

func canManageOIDC(isEnabled bool, oidcType string) bool {
	if !isEnabled && (oidcType != "" && oidcType != "NONE") {
		return false
	}

	return true
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
