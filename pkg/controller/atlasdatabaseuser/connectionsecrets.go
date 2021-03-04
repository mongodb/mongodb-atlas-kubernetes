package atlasdatabaseuser

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func createOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) error {
	clusters, _, err := ctx.Client.Clusters.List(context.Background(), project.ID(), &mongodbatlas.ListOptions{})
	if err != nil {
		// TODO ignore the 404 exception in case no clusters exist by this time
		return err
	}

	secretNames := make(map[string]string)
	for _, cluster := range clusters {
		password, err := dbUser.ReadPassword(k8sClient)
		if err != nil {
			return err
		}
		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			ConnURL:    cluster.ConnectionStrings.Standard,
			SrvConnURL: cluster.ConnectionStrings.StandardSrv,
			Password:   password,
		}
		var secretName string
		if secretName, err = connectionsecret.Ensure(k8sClient, dbUser.Namespace, project.Spec.Name, project.ID(), cluster.Name, data); err != nil {
			return err
		}
		ctx.Log.Debugw("Ensured connection Secret up-to-date", "name", secretName)
		secretNames[cluster.Name] = secretName
	}

	// TODO we need to remove old secrets in case the dbuser name has changed

	ctx.EnsureStatusOption(secretNames)
	return nil
}
