package atlasdatabaseuser

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"
)

func CreateOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	clusters, _, err := ctx.Client.Clusters.List(context.Background(), project.ID(), &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}

	secretNames := make(map[string]string)
	requeue := false
	for _, cluster := range clusters {
		scopes := dbUser.GetScopes(mdbv1.ClusterScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, cluster.Name) {
			continue
		}
		// Cluster may be not ready yet, so no connection urls - skipping
		// Note, that Atlas usually returns the not-nil connection strings with empty fields in it
		if cluster.ConnectionStrings == nil || cluster.ConnectionStrings.StandardSrv == "" {
			ctx.Log.Debugw("Cluster is not ready yet - not creating a connection Secret", "cluster", cluster.Name)
			requeue = true
			continue
		}
		// Cluster may be not ready yet, so no connection urls - skipping
		// Note, that Atlas usually returns the not-nil connection strings with empty fields in it
		if cluster.ConnectionStrings == nil || cluster.ConnectionStrings.StandardSrv == "" {
			ctx.Log.Debugw("Cluster is not ready yet - not creating a connection Secret", "cluster", cluster.Name)
			requeue = true
			continue
		}
		password, err := dbUser.ReadPassword(k8sClient)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
		}
		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			ConnURL:    cluster.ConnectionStrings.Standard,
			SrvConnURL: cluster.ConnectionStrings.StandardSrv,
			Password:   password,
		}
		var secretName string
		if secretName, err = connectionsecret.Ensure(k8sClient, dbUser.Namespace, project.Spec.Name, project.ID(), cluster.Name, data); err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
		}
		ctx.Log.Debugw("Ensured connection Secret up-to-date", "secretname", secretName)
		secretNames[cluster.Name] = secretName
	}

	if err := cleanupStaleSecrets(ctx, k8sClient, project.ID(), dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserStaleConnectionSecrets, err.Error())
	}

	ctx.EnsureStatusOption(status.AtlasDatabaseUserSecretsOption(secretNames))
	if requeue {
		return workflow.InProgress(workflow.DatabaseUserConnectionSecretsNotCreated, "Waiting for clusters to get created/updated")
	}
	return workflow.OK()
}

func cleanupStaleSecrets(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	if err := removeStaleByScope(ctx, k8sClient, projectID, user); err != nil {
		return err
	}
	if err := removeStaleByUserName(ctx, k8sClient, projectID, user); err != nil {
		return err
	}
	return nil
}

// removeStaleByScope removes the secrets that are not relevant due to changes to 'scopes' field for the AtlasDatabaseUser.
func removeStaleByScope(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	scopes := user.GetScopes(mdbv1.ClusterScopeType)
	if len(scopes) == 0 {
		return nil
	}
	secrets, err := connectionsecret.ListByUserName(k8sClient, user.Namespace, projectID, user.Spec.Username)
	if err != nil {
		return err
	}
	for i, s := range secrets {
		cluster, ok := s.Labels[connectionsecret.ClusterLabelKey]
		if !ok {
			continue
		}
		if !stringutil.Contains(scopes, cluster) {
			if err = k8sClient.Delete(context.Background(), &secrets[i]); err != nil {
				return err
			}
			ctx.Log.Debugw("Removed connection Secret as it's not referenced by the AtlasDatabaseUser anymore", "secretname", s.Name)
		}
	}
	return nil
}

// removeStaleByUserName removes the stale secrets when the database user name changes (as it's used as a part of Secret name)
func removeStaleByUserName(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	if user.Status.UserName == user.Spec.Username {
		return nil
	}
	secrets, err := connectionsecret.ListByUserName(k8sClient, user.Namespace, projectID, user.Status.UserName)
	if err != nil {
		return err
	}
	for i, s := range secrets {
		if err = k8sClient.Delete(context.Background(), &secrets[i]); err != nil {
			return err
		}
		ctx.Log.Debugw("Removed connection Secret as the database user name has changed", "secretname", s.Name)
	}
	return nil
}
