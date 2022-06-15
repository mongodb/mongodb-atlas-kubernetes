package atlasdatabaseuser

import (
	"context"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"
)

const ConnectionSecretsEnsuredEvent = "ConnectionSecretsEnsured"

func CreateOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, recorder record.EventRecorder, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	deployments, _, err := ctx.Client.Clusters.List(context.Background(), project.ID(), &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}

	advancedDeployments, _, err := ctx.Client.AdvancedClusters.List(context.Background(), project.ID(), &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}

	var deploymentSecrets []deploymentSecret
	for _, c := range deployments {
		deploymentSecrets = append(deploymentSecrets, deploymentSecret{
			name:              c.Name,
			connectionStrings: c.ConnectionStrings,
		})
	}

	for _, c := range advancedDeployments.Results {
		// based on configuration settings, some advanced deployments also show up in the regular deployments API.
		// For these deployments, we don't want to duplicate the secret so we skip them.
		found := false
		for _, regularDeployment := range deployments {
			if regularDeployment.Name == c.Name {
				found = true
				break
			}
		}

		// we only include secrets which have not been handled by the regular deployment API.
		if !found {
			deploymentSecrets = append(deploymentSecrets, deploymentSecret{
				name:              c.Name,
				connectionStrings: c.ConnectionStrings,
			})
		}
	}

	// ensure secrets for both deployments and advanced deployment.
	if result := createOrUpdateConnectionSecretsFromDeploymentSecrets(ctx, k8sClient, recorder, project, dbUser, deploymentSecrets); !result.IsOk() {
		return result
	}

	return workflow.OK()
}

// deploymentSecret holds the information required to ensure a secret for a user in a given deployment.
type deploymentSecret struct {
	name              string
	connectionStrings *mongodbatlas.ConnectionStrings
}

func createOrUpdateConnectionSecretsFromDeploymentSecrets(ctx *workflow.Context, k8sClient client.Client, recorder record.EventRecorder, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser, deploymentSecrets []deploymentSecret) workflow.Result {
	requeue := false
	secrets := make([]string, 0)

	for _, ds := range deploymentSecrets {
		scopes := dbUser.GetScopes(mdbv1.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, ds.name) {
			continue
		}
		// Deployment may be not ready yet, so no connection urls - skipping
		// Note, that Atlas usually returns the not-nil connection strings with empty fields in it
		if ds.connectionStrings == nil || ds.connectionStrings.StandardSrv == "" {
			ctx.Log.Debugw("Deployment is not ready yet - not creating a connection Secret", "deployment", ds.name)
			requeue = true
			continue
		}
		password, err := dbUser.ReadPassword(k8sClient)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
		}
		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    ds.connectionStrings.Standard,
			SrvConnURL: ds.connectionStrings.StandardSrv,
		}
		fillPrivateConnStrings(ds.connectionStrings, &data)

		var secretName string
		if secretName, err = connectionsecret.Ensure(k8sClient, dbUser.Namespace, project.Spec.Name, project.ID(), ds.name, data); err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
		}
		secrets = append(secrets, secretName)
		ctx.Log.Debugw("Ensured connection Secret up-to-date", "secretname", secretName)
	}

	if len(secrets) > 0 {
		recorder.Eventf(&dbUser, "Normal", ConnectionSecretsEnsuredEvent, "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	if err := cleanupStaleSecrets(ctx, k8sClient, project.ID(), dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserStaleConnectionSecrets, err.Error())
	}

	if requeue {
		return workflow.InProgress(workflow.DatabaseUserConnectionSecretsNotCreated, "Waiting for deployments to get created/updated")
	}
	return workflow.OK()
}

func cleanupStaleSecrets(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	if err := removeStaleByScope(ctx, k8sClient, projectID, user); err != nil {
		return err
	}
	// Performing the cleanup of old secrets only if the username has changed
	if user.Status.UserName != user.Spec.Username {
		// Note, that we pass the username from the status, not from the spec
		return removeStaleSecretsByUserName(k8sClient, projectID, user.Status.UserName, user, ctx.Log)
	}
	return nil
}

// removeStaleByScope removes the secrets that are not relevant due to changes to 'scopes' field for the AtlasDatabaseUser.
func removeStaleByScope(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	scopes := user.GetScopes(mdbv1.DeploymentScopeType)
	if len(scopes) == 0 {
		return nil
	}
	secrets, err := connectionsecret.ListByUserName(k8sClient, user.Namespace, projectID, user.Spec.Username)
	if err != nil {
		return err
	}
	for i, s := range secrets {
		deployment, ok := s.Labels[connectionsecret.ClusterLabelKey]
		if !ok {
			continue
		}
		if !stringutil.Contains(scopes, deployment) {
			if err = k8sClient.Delete(context.Background(), &secrets[i]); err != nil {
				return err
			}
			ctx.Log.Debugw("Removed connection Secret as it's not referenced by the AtlasDatabaseUser anymore", "secretname", s.Name)
		}
	}
	return nil
}

// removeStaleSecretsByUserName removes the stale secrets when the database user name changes (as it's used as a part of Secret name)
func removeStaleSecretsByUserName(k8sClient client.Client, projectID, userName string, user mdbv1.AtlasDatabaseUser, log *zap.SugaredLogger) error {
	secrets, err := connectionsecret.ListByUserName(k8sClient, user.Namespace, projectID, userName)
	if err != nil {
		return err
	}
	var lastError error
	removed := 0
	for i := range secrets {
		if err = k8sClient.Delete(context.Background(), &secrets[i]); err != nil {
			log.Errorf("Failed to remove connection Secret: %v", err)
			lastError = err
		} else {
			log.Debugw("Removed connection Secret", "secret", kube.ObjectKeyFromObject(&secrets[i]))
			removed++
		}
	}
	if removed > 0 {
		log.Infof("Removed %d connection secrets", removed)
	}
	return lastError
}

func fillPrivateConnStrings(connStrings *mongodbatlas.ConnectionStrings, data *connectionsecret.ConnectionData) {
	if connStrings.Private != "" {
		data.PvtConnURL = connStrings.Private
		data.PvtSrvConnURL = connStrings.PrivateSrv
	}

	if len(connStrings.PrivateEndpoint) == 1 {
		pe := connStrings.PrivateEndpoint[0]
		data.PvtConnURL = pe.ConnectionString
		data.PvtSrvConnURL = pe.SRVConnectionString
	}
}
