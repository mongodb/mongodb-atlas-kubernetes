package connectionsecret

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"
)

const ConnectionSecretsEnsuredEvent = "ConnectionSecretsEnsured"

func CreateOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, recorder record.EventRecorder, project mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	advancedDeployments, _, err := ctx.Client.AdvancedClusters.List(context.Background(), project.ID(), &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}

	var deploymentSecrets []deploymentSecret
	for _, c := range advancedDeployments.Results {
		deploymentSecrets = append(deploymentSecrets, deploymentSecret{
			name:              c.Name,
			connectionStrings: c.ConnectionStrings,
		})
	}

	serverlessDeployments, err := GetAllServerless(ctx, project.ID())
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err.Error())
	}
	for _, c := range serverlessDeployments {
		found := false

		for _, advancedDeployment := range advancedDeployments.Results {
			if advancedDeployment.Name == c.Name {
				found = true
				break
			}
		}

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
		data := ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    ds.connectionStrings.Standard,
			SrvConnURL: ds.connectionStrings.StandardSrv,
		}
		FillPrivateConnStrings(ds.connectionStrings, &data)

		var secretName string
		if secretName, err = Ensure(k8sClient, dbUser.Namespace, project.Spec.Name, project.ID(), ds.name, data); err != nil {
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
		return RemoveStaleSecretsByUserName(k8sClient, projectID, user.Status.UserName, user, ctx.Log)
	}
	return nil
}

// removeStaleByScope removes the secrets that are not relevant due to changes to 'scopes' field for the AtlasDatabaseUser.
func removeStaleByScope(ctx *workflow.Context, k8sClient client.Client, projectID string, user mdbv1.AtlasDatabaseUser) error {
	scopes := user.GetScopes(mdbv1.DeploymentScopeType)
	if len(scopes) == 0 {
		return nil
	}
	secrets, err := ListByUserName(k8sClient, user.Namespace, projectID, user.Spec.Username)
	if err != nil {
		return err
	}
	for i, s := range secrets {
		deployment, ok := s.Labels[ClusterLabelKey]
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

// RemoveStaleSecretsByUserName removes the stale secrets when the database user name changes (as it's used as a part of Secret name)
func RemoveStaleSecretsByUserName(k8sClient client.Client, projectID, userName string, user mdbv1.AtlasDatabaseUser, log *zap.SugaredLogger) error {
	secrets, err := ListByUserName(k8sClient, user.Namespace, projectID, userName)
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

func FillPrivateConnStrings(connStrings *mongodbatlas.ConnectionStrings, data *ConnectionData) {
	if connStrings.Private != "" {
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:    connStrings.Private,
			PvtSrvConnURL: connStrings.PrivateSrv,
		})
	}

	for _, pe := range connStrings.PrivateEndpoint {
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:      pe.ConnectionString,
			PvtSrvConnURL:   pe.SRVConnectionString,
			PvtShardConnURL: pe.SRVShardOptimizedConnectionString,
		})
	}
}

func GetAllServerless(ctx *workflow.Context, projectID string) ([]*mongodbatlas.Cluster, error) {
	serverless, _, err := ctx.Client.ServerlessInstances.List(context.Background(), projectID, nil)
	if err != nil {
		if !IsCloudGovDomain(ctx) {
			return nil, fmt.Errorf("error getting serverless: %w", err)
		} else {
			return make([]*mongodbatlas.Cluster, 0), nil
		}
	}
	return serverless.Results, nil
}

func IsCloudGovDomain(ctx *workflow.Context) bool {
	domains := []string{
		"cloudgov.mongodb.com",
		"cloud.mongodbgov.com",
		"cloud-dev.mongodbgov.com",
		"cloud-qa.mongodbgov.com",
	}

	for _, domain := range domains {
		if strings.HasPrefix(ctx.Client.BaseURL.Host, domain) {
			return true
		}
	}

	return false
}
