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

// createOrUpdateConnectionSecrets
func createOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
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

	// TODO CLOUDP-84202 we need to remove old secrets in case the dbuser name has changed

	// TODO 2 CLOUDP-84202 : we need to remove the secrets that don't match the scope anymore

	ctx.EnsureStatusOption(status.AtlasDatabaseUserSecretsOption(secretNames))
	if requeue {
		return workflow.InProgress(workflow.DatabaseUserConnectionSecretsNotCreated, "Waiting for clusters to get created/updated")
	}
	return workflow.OK()
}
