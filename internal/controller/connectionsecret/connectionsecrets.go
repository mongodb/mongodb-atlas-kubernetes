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

package connectionsecret

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

const ConnectionSecretsEnsuredEvent = "ConnectionSecretsEnsured"

func ReapOrphanConnectionSecrets(ctx context.Context, k8sClient client.Client, projectID, namespace string, projectDeploymentNames []string) ([]string, error) {
	secretList := &corev1.SecretList{}
	labelSelector := labels.SelectorFromSet(labels.Set{TypeLabelKey: CredLabelVal, ProjectLabelKey: projectID})
	err := k8sClient.List(context.Background(), secretList, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed listing possible orphan secrets: %w", err)
	}

	removedOrphanSecrets := []string{}
	for _, secret := range secretList.Items {
		clusterName, ok := secret.Labels[ClusterLabelKey]
		if !ok {
			continue
		}
		if clusterExists := stringutil.Contains(projectDeploymentNames, clusterName); clusterExists {
			continue
		}
		if err := k8sClient.Delete(ctx, &secret); err != nil {
			return nil, fmt.Errorf("failed to remove orphan connection Secret: %w", err)
		} else {
			removedOrphanSecrets = append(removedOrphanSecrets, fmt.Sprintf("%s/%s", namespace, secret.Name))
		}
	}
	return removedOrphanSecrets, nil
}

func CreateOrUpdateConnectionSecrets(ctx *workflow.Context, k8sClient client.Client, ds deployment.AtlasDeploymentsService, recorder record.EventRecorder, project *project.Project, dbUser akov2.AtlasDatabaseUser) workflow.Result {
	conns, err := ds.ListDeploymentConnections(ctx.Context, project.ID)
	if err != nil {
		return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err)
	}

	// ensure secrets for both deployments and advanced deployment.
	if result := createOrUpdateConnectionSecretsFromDeploymentSecrets(ctx, k8sClient, recorder, project, dbUser, conns); !result.IsOk() {
		return result
	}

	return workflow.OK()
}

func createOrUpdateConnectionSecretsFromDeploymentSecrets(ctx *workflow.Context, k8sClient client.Client, recorder record.EventRecorder, project *project.Project, dbUser akov2.AtlasDatabaseUser, conns []deployment.Connection) workflow.Result {
	requeue := false
	secrets := make([]string, 0)

	for _, di := range conns {
		scopes := dbUser.GetScopes(akov2.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, di.Name) {
			continue
		}
		// Deployment may be not ready yet, so no connection urls - skipping
		// Note, that Atlas usually returns the not-nil connection strings with empty fields in it
		if di.SrvConnURL == "" {
			ctx.Log.Debugw("Deployment is not ready yet - not creating a connection Secret", "deployment", di.Name)
			requeue = true
			continue
		}
		password, err := dbUser.ReadPassword(ctx.Context, k8sClient)
		if err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err)
		}
		data := ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    di.ConnURL,
			SrvConnURL: di.SrvConnURL,
		}
		FillPrivateConns(di, &data)

		var secretName string
		if secretName, err = Ensure(ctx.Context, k8sClient, dbUser.Namespace, project.Name, project.ID, di.Name, data); err != nil {
			return workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotCreated, err)
		}
		secrets = append(secrets, secretName)
		ctx.Log.Debugw("Ensured connection Secret up-to-date", "secretname", secretName)
	}

	if len(secrets) > 0 {
		recorder.Eventf(&dbUser, "Normal", ConnectionSecretsEnsuredEvent, "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	if err := cleanupStaleSecrets(ctx, k8sClient, project.ID, dbUser); err != nil {
		return workflow.Terminate(workflow.DatabaseUserStaleConnectionSecrets, err)
	}

	if requeue {
		return workflow.InProgress(workflow.DatabaseUserConnectionSecretsNotCreated, "Waiting for deployments to get created/updated")
	}
	return workflow.OK()
}

func cleanupStaleSecrets(ctx *workflow.Context, k8sClient client.Client, projectID string, user akov2.AtlasDatabaseUser) error {
	if err := removeStaleByScope(ctx, k8sClient, projectID, user); err != nil {
		return err
	}
	// Performing the cleanup of old secrets only if the username has changed
	if user.Status.UserName != user.Spec.Username {
		// Note, that we pass the username from the status, not from the spec
		return RemoveStaleSecretsByUserName(ctx.Context, k8sClient, projectID, user.Status.UserName, user, ctx.Log)
	}
	return nil
}

// removeStaleByScope removes the secrets that are not relevant due to changes to 'scopes' field for the AtlasDatabaseUser.
func removeStaleByScope(ctx *workflow.Context, k8sClient client.Client, projectID string, user akov2.AtlasDatabaseUser) error {
	scopes := user.GetScopes(akov2.DeploymentScopeType)
	if len(scopes) == 0 {
		return nil
	}
	secrets, err := ListByUserName(ctx.Context, k8sClient, user.Namespace, projectID, user.Spec.Username)
	if err != nil {
		return err
	}
	for i, s := range secrets {
		deployment, ok := s.Labels[ClusterLabelKey]
		if !ok {
			continue
		}
		if !stringutil.Contains(scopes, deployment) {
			if err = k8sClient.Delete(ctx.Context, &secrets[i]); err != nil {
				return err
			}
			ctx.Log.Debugw("Removed connection Secret as it's not referenced by the AtlasDatabaseUser anymore", "secretname", s.Name)
		}
	}
	return nil
}

// RemoveStaleSecretsByUserName removes the stale secrets when the database user name changes (as it's used as a part of Secret name)
func RemoveStaleSecretsByUserName(ctx context.Context, k8sClient client.Client, projectID, userName string, user akov2.AtlasDatabaseUser, log *zap.SugaredLogger) error {
	secrets, err := ListByUserName(ctx, k8sClient, user.Namespace, projectID, userName)
	if err != nil {
		return err
	}
	var lastError error
	removed := 0
	for i := range secrets {
		if err = k8sClient.Delete(ctx, &secrets[i]); err != nil {
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

func FillPrivateConns(conn deployment.Connection, data *ConnectionData) {
	if conn.PrivateURL != "" {
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:    conn.PrivateURL,
			PvtSrvConnURL: conn.SrvPrivateURL,
		})
	}

	if conn.Serverless {
		for _, pe := range conn.PrivateEndpoints {
			data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
				PvtSrvConnURL: pe.ServerURL,
			})
		}
	} else {
		for _, pe := range conn.PrivateEndpoints {
			data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
				PvtConnURL:      pe.URL,
				PvtSrvConnURL:   pe.ServerURL,
				PvtShardConnURL: pe.ShardURL,
			})
		}
	}
}

// FillPrivateConnStrings fills private conn urls from connection strings
// TODO: (CLOUDP-253951) remove once all usages move over to FillPrivateConns instead
// Right now only advanced deployment is using this one
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
	serverless, _, err := ctx.Client.ServerlessInstances.List(ctx.Context, projectID, nil)
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
