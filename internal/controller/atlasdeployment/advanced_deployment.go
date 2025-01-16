package atlasdeployment

import (
	"fmt"
	"reflect"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
)

const FreeTier = "M0"

func (r *AtlasDeploymentReconciler) handleAdvancedDeployment(ctx *workflow.Context, projectService project.ProjectService, deploymentService deployment.AtlasDeploymentsService, deploymentInAKO, deploymentInAtlas *deployment.Cluster) (ctrl.Result, error) {
	if deploymentInAtlas == nil {
		ctx.Log.Infof("Advanced Deployment %s doesn't exist in Atlas - creating", deploymentInAKO.GetName())
		newDeployment, err := deploymentService.CreateDeployment(ctx.Context, deploymentInAKO)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		deploymentInAtlas = newDeployment.(*deployment.Cluster)
	}

	switch deploymentInAtlas.GetState() {
	case status.StateIDLE:
		if changes, occurred := deployment.ComputeChanges(deploymentInAKO, deploymentInAtlas); occurred {
			updatedDeployment, err := deploymentService.UpdateDeployment(ctx.Context, changes)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), updatedDeployment, workflow.DeploymentUpdating, "deployment is updating")
		}

		transition := r.ensureBackupScheduleAndPolicy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource())
		if transition != nil {
			return transition(workflow.Internal)
		}

		transition = r.ensureAdvancedOptions(ctx, deploymentService, deploymentInAKO, deploymentInAtlas)
		if transition != nil {
			return transition(workflow.DeploymentAdvancedOptionsReady)
		}

		err := r.ensureConnectionSecrets(ctx, projectService, deploymentInAKO, deploymentInAtlas.GetConnection())
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentConnectionSecretsNotCreated, err)
		}

		if !r.AtlasProvider.IsCloudGov() {
			searchNodeResult := handleSearchNodes(ctx, deploymentInAKO.GetCustomResource(), deploymentInAKO.GetProjectID())
			if transition = r.transitionFromResult(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), searchNodeResult); transition != nil {
				return transition(workflow.Internal)
			}
		}

		searchService := searchindex.NewSearchIndexes(ctx.SdkClientSet.SdkClient20241113001.AtlasSearchApi)
		result := handleSearchIndexes(ctx, r.Client, searchService, deploymentInAKO.GetCustomResource(), deploymentInAKO.GetProjectID())
		if transition = r.transitionFromResult(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), result); transition != nil {
			return transition(workflow.Internal)
		}

		result = r.ensureCustomZoneMapping(
			ctx,
			deploymentService,
			deploymentInAKO.GetProjectID(),
			deploymentInAKO.GetCustomResource().Spec.DeploymentSpec.CustomZoneMapping,
			deploymentInAKO.GetName(),
		)
		if transition = r.transitionFromResult(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), result); transition != nil {
			return transition(workflow.Internal)
		}

		result = r.ensureManagedNamespaces(
			ctx,
			deploymentService,
			deploymentInAKO.GetProjectID(),
			deploymentInAKO.ClusterType,
			deploymentInAKO.GetCustomResource().Spec.DeploymentSpec.ManagedNamespaces,
			deploymentInAKO.GetName(),
		)
		if transition = r.transitionFromResult(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), result); transition != nil {
			return transition(workflow.Internal)
		}

		err = customresource.ApplyLastConfigApplied(ctx.Context, deploymentInAKO.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas)
	case status.StateCREATING:
		return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
	case status.StateDELETING, status.StateDELETED:
		return r.deleted()
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", deploymentInAtlas.GetState()))
	}
}

func (r *AtlasDeploymentReconciler) ensureConnectionSecrets(ctx *workflow.Context, projectService project.ProjectService, deploymentInAKO deployment.Deployment, connection *status.ConnectionStrings) error {
	databaseUsers := &akov2.AtlasDatabaseUserList{}
	listOpts := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserByProject, deploymentInAKO.GetProjectID()),
	}
	err := r.Client.List(ctx.Context, databaseUsers, listOpts)
	if err != nil {
		return err
	}

	secrets := make([]string, 0)
	for _, dbUser := range databaseUsers.Items {
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
		if len(scopes) != 0 && !stringutil.Contains(scopes, deploymentInAKO.GetName()) {
			continue
		}

		password, err := dbUser.ReadPassword(ctx.Context, r.Client)
		if err != nil {
			return err
		}

		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    connection.Standard,
			SrvConnURL: connection.StandardSrv,
		}
		if connection.Private != "" {
			data.PrivateConnURLs = append(data.PrivateConnURLs, connectionsecret.PrivateLinkConnURLs{
				PvtConnURL:    connection.Private,
				PvtSrvConnURL: connection.PrivateSrv,
			})
		}

		for _, pe := range connection.PrivateEndpoint {
			data.PrivateConnURLs = append(data.PrivateConnURLs, connectionsecret.PrivateLinkConnURLs{
				PvtConnURL:      pe.ConnectionString,
				PvtSrvConnURL:   pe.SRVConnectionString,
				PvtShardConnURL: pe.SRVShardOptimizedConnectionString,
			})
		}

		project, err := projectService.GetProject(ctx.Context, deploymentInAKO.GetProjectID())
		if err != nil {
			return err
		}

		ctx.Log.Debugw("Creating a connection Secret", "data", data)
		secretName, err := connectionsecret.Ensure(ctx.Context, r.Client, dbUser.Namespace, project.Name, deploymentInAKO.GetProjectID(), deploymentInAKO.GetName(), data)
		if err != nil {
			return err
		}
		secrets = append(secrets, secretName)
	}

	if len(secrets) > 0 {
		r.EventRecorder.Eventf(deploymentInAKO.GetCustomResource(), "Normal", "ConnectionSecretsEnsured", "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	return nil
}

func (r *AtlasDeploymentReconciler) ensureAdvancedOptions(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, deploymentInAKO, deploymentInAtlas *deployment.Cluster) transitionFn {
	if deploymentInAKO.IsTenant() {
		return nil
	}

	err := deploymentService.ClusterWithProcessArgs(ctx.Context, deploymentInAtlas)
	if err != nil {
		return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), err)
	}

	if deploymentInAKO.ProcessArgs != nil && !reflect.DeepEqual(deploymentInAKO.ProcessArgs, deploymentInAtlas.ProcessArgs) {
		err = deploymentService.UpdateProcessArgs(ctx.Context, deploymentInAKO)
		if err != nil {
			return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), err)
		}

		return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), nil)
	}

	return nil
}
