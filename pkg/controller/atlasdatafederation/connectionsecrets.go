package atlasdatafederation

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDataFederationReconciler) ensureConnectionSecrets(ctx *workflow.Context, federationService datafederation.DataFederationService, project *akov2.AtlasProject, df *akov2.AtlasDataFederation) workflow.Result {
	databaseUsers := akov2.AtlasDatabaseUserList{}
	err := r.Client.List(ctx.Context, &databaseUsers, &client.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	atlasDF, err := federationService.Get(ctx.Context, project.ID(), df.Spec.Name)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	connectionHosts := atlasDF.Hostnames

	secrets := make([]string, 0)
	for i := range databaseUsers.Items {
		dbUser := databaseUsers.Items[i]

		if !dbUserBelongsToProject(&dbUser, project) {
			continue
		}

		found := false
		for _, c := range dbUser.Status.Conditions {
			if c.Type == api.ReadyType && c.Status == v1.ConditionTrue {
				found = true
				break
			}
		}

		if !found {
			ctx.Log.Debugw("AtlasDatabaseUser not ready - not creating connection secret", "user.name", dbUser.Name)
			continue
		}

		scopes := dbUser.GetScopes(akov2.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, df.Spec.Name) {
			continue
		}

		password, err := dbUser.ReadPassword(ctx.Context, r.Client)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}

		var connURLs []string
		for _, host := range connectionHosts {
			connURLs = append(connURLs, fmt.Sprintf("mongodb://%s:%s@%s?ssl=true", dbUser.Spec.Username, password, host))
		}

		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    strings.Join(connURLs, ","),
		}

		ctx.Log.Debugw("Creating a connection Secret", "data", data)

		secretName, err := connectionsecret.Ensure(ctx.Context, r.Client, dbUser.Namespace, project.Spec.Name, project.ID(), df.Spec.Name, data)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}
		secrets = append(secrets, secretName)
	}

	if len(secrets) > 0 {
		r.EventRecorder.Eventf(df, "Normal", "ConnectionSecretsEnsured", "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	return workflow.OK()
}

func dbUserBelongsToProject(dbUser *akov2.AtlasDatabaseUser, project *akov2.AtlasProject) bool {
	if dbUser.Spec.Project.Name != project.Name {
		return false
	}

	if dbUser.Spec.Project.Namespace == "" && dbUser.Namespace != project.Namespace {
		return false
	}

	if dbUser.Spec.Project.Namespace != "" && dbUser.Spec.Project.Namespace != project.Namespace {
		return false
	}

	return true
}
