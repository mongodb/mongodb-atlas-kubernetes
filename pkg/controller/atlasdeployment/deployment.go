package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"
)

func ensureDeploymentState(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment) (atlasDeployment *mongodbatlas.Cluster, _ workflow.Result) {
	atlasDeployment, resp, err := ctx.Client.Clusters.Get(context.Background(), project.Status.ID, deployment.Spec.DeploymentSpec.Name)
	if err != nil {
		if resp == nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return atlasDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}

		atlasDeployment, err = deployment.Spec.Deployment()
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Deployment %s doesn't exist in Atlas - creating", deployment.Spec.DeploymentSpec.Name)
		atlasDeployment, _, err = ctx.Client.Clusters.Create(context.Background(), project.Status.ID, atlasDeployment)
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}
	}

	result := EnsureCustomZoneMapping(ctx, project.ID(), deployment.Spec.DeploymentSpec.CustomZoneMapping, atlasDeployment.Name)
	if !result.IsOk() {
		return atlasDeployment, result
	}

	result = EnsureManagedNamespaces(ctx, project.ID(), string(deployment.Spec.DeploymentSpec.ClusterType), deployment.Spec.DeploymentSpec.ManagedNamespaces, atlasDeployment.Name)
	if !result.IsOk() {
		return atlasDeployment, result
	}

	switch atlasDeployment.StateName {
	case status.StateIDLE:

		return regularDeploymentIdle(ctx, project, deployment, atlasDeployment)
	case status.StateCREATING:
		return atlasDeployment, workflow.InProgress(workflow.DeploymentCreating, "deployment is provisioning")

	case "UPDATING", "REPAIRING":
		return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return atlasDeployment, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown deployment state %q", atlasDeployment.StateName))
	}
}

func regularDeploymentIdle(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, atlasDeployment *mongodbatlas.Cluster) (*mongodbatlas.Cluster, workflow.Result) {
	resultingDeployment, err := MergedDeployment(*atlasDeployment, deployment.Spec)
	if err != nil {
		return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
	}

	if done := DeploymentsEqual(ctx.Log, *atlasDeployment, resultingDeployment); done {
		return atlasDeployment, workflow.OK()
	}

	if deployment.Spec.DeploymentSpec.Paused != nil {
		if atlasDeployment.Paused == nil || *atlasDeployment.Paused != *deployment.Spec.DeploymentSpec.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			resultingDeployment = mongodbatlas.Cluster{
				Paused: deployment.Spec.DeploymentSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			resultingDeployment.Paused = nil
		}
	}

	resultingDeployment = cleanupDeployment(resultingDeployment)

	// Handle shared (M0,M2,M5) deployment to non-shared deployment upgrade
	scheduled, err := handleSharedDeploymentUpgrade(ctx, atlasDeployment, &resultingDeployment)
	if err != nil {
		return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
	}
	if scheduled {
		return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is upgrading")
	}

	atlasDeployment, _, err = ctx.Client.Clusters.Update(context.Background(), project.Status.ID, deployment.Spec.DeploymentSpec.Name, &resultingDeployment)
	if err != nil {
		return atlasDeployment, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

// cleanupDeployment will unset some fields that cannot be changed via API or are deprecated.
func cleanupDeployment(deployment mongodbatlas.Cluster) mongodbatlas.Cluster {
	deployment.ID = ""
	deployment.MongoDBVersion = ""
	deployment.MongoURI = ""
	deployment.MongoURIUpdated = ""
	deployment.MongoURIWithOptions = ""
	deployment.SrvAddress = ""
	deployment.StateName = ""
	deployment.ReplicationFactor = nil
	deployment.ReplicationSpec = nil
	deployment.ConnectionStrings = nil
	deployment = removeOutdatedFields(&deployment, nil)
	if IsFreeTierCluster(&deployment) {
		deployment.DiskSizeGB = nil
	}

	if deployment.AutoScaling != nil {
		deployment.AutoScaling.AutoIndexingEnabled = nil
	}

	if deployment.ProviderSettings != nil && deployment.ProviderSettings.AutoScaling != nil {
		deployment.ProviderSettings.AutoScaling.AutoIndexingEnabled = nil
	}

	return deployment
}

func IsFreeTierCluster(deployment *mongodbatlas.Cluster) bool {
	if deployment != nil && deployment.ProviderSettings != nil && deployment.ProviderSettings.InstanceSizeName == "M0" {
		return true
	}
	return false
}

// removeOutdatedFields unsets fields which are should be empty based on flags
func removeOutdatedFields(removeFrom *mongodbatlas.Cluster, lookAt *mongodbatlas.Cluster) mongodbatlas.Cluster {
	if lookAt == nil {
		lookAt = removeFrom
	}

	result := *removeFrom
	if lookAt.AutoScaling != nil && lookAt.AutoScaling.Compute != nil {
		if *lookAt.AutoScaling.Compute.Enabled {
			result.ProviderSettings.InstanceSizeName = ""
		} else {
			if result.ProviderSettings == nil {
				result.ProviderSettings = &mongodbatlas.ProviderSettings{}
			}
			if result.ProviderSettings.AutoScaling == nil {
				result.ProviderSettings.AutoScaling = &mongodbatlas.AutoScaling{}
			}
			result.ProviderSettings.AutoScaling.Compute = &mongodbatlas.Compute{}
			result.ProviderSettings.AutoScaling.AutoIndexingEnabled = nil
		}
	}

	if lookAt.AutoScaling != nil {
		result.AutoScaling.AutoIndexingEnabled = nil

		if lookAt.AutoScaling.DiskGBEnabled != nil && *lookAt.AutoScaling.DiskGBEnabled {
			result.DiskSizeGB = nil
		}
	}

	return result
}

// MergedDeployment will return the result of merging AtlasDeploymentSpec with Atlas Deployment
func MergedDeployment(atlasDeployment mongodbatlas.Cluster, spec mdbv1.AtlasDeploymentSpec) (result mongodbatlas.Cluster, err error) {
	if err = compat.JSONCopy(&result, atlasDeployment); err != nil {
		return
	}

	if err = compat.JSONCopy(&result, spec.DeploymentSpec); err != nil {
		return
	}

	mergeRegionConfigs(result.ReplicationSpecs, spec.DeploymentSpec.ReplicationSpecs)

	// According to the docs for 'providerSettings.regionName' (https://docs.atlas.mongodb.com/reference/api/clusters-create-one/):
	// "Don't specify this parameter when creating a multi-region deployment using the replicationSpec object or a Global
	// Deployment with the replicationSpecs array."
	// The problem is that Atlas API accepts the create/update request but then returns the 'ProviderSettings.RegionName' empty in GET request
	// So we need to consider this while comparing (to avoid perpetual updates)
	if len(result.ReplicationSpecs) > 0 && atlasDeployment.ProviderSettings.RegionName == "" {
		result.ProviderSettings.RegionName = ""
	}

	return
}

// mergeRegionConfigs removes replicationSpecs[i].RegionsConfigs[key] from Atlas Deployment that are absent in Operator.
// Dev idea: this could have been added into some more generic method like `JSONCopy` or something wrapping it to make
// sure any Atlas map get redundant keys removed. So far there's only one map in Deployment ('RegionsConfig') so we'll do this
// explicitly - but may make sense to refactor this later if more maps are added (and all follow the same logic).
func mergeRegionConfigs(atlasSpecs []mongodbatlas.ReplicationSpec, operatorSpecs []mdbv1.ReplicationSpec) {
	for i, operatorSpec := range operatorSpecs {
		if len(operatorSpec.RegionsConfig) == 0 {
			// Edge case: if the operator doesn't specify regions configs - Atlas will put the default ones. We shouldn't
			// remove it in this case.
			continue
		}
		atlasSpec := atlasSpecs[i]
		for key := range atlasSpec.RegionsConfig {
			if _, ok := operatorSpec.RegionsConfig[key]; !ok {
				delete(atlasSpec.RegionsConfig, key)
			}
		}
	}
}

// DeploymentsEqual compares two Atlas Deployments
func DeploymentsEqual(log *zap.SugaredLogger, deploymentAtlas mongodbatlas.Cluster, deploymentOperator mongodbatlas.Cluster) bool {
	deploymentAtlas = removeOutdatedFields(&deploymentAtlas, &deploymentOperator)
	deploymentOperator = removeOutdatedFields(&deploymentOperator, nil)

	d := cmp.Diff(deploymentAtlas, deploymentOperator, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Deployments are different: %s", d)
	}

	return d == ""
}

func (r *AtlasDeploymentReconciler) ensureConnectionSecrets(ctx *workflow.Context, project *mdbv1.AtlasProject, name string, connectionStrings *mongodbatlas.ConnectionStrings, deploymentResource *mdbv1.AtlasDeployment) workflow.Result {
	databaseUsers := mdbv1.AtlasDatabaseUserList{}
	err := r.Client.List(context.TODO(), &databaseUsers, client.InNamespace(project.Namespace))
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	secrets := make([]string, 0)
	for _, dbUser := range databaseUsers.Items {
		found := false
		for _, c := range dbUser.Status.Conditions {
			if c.Type == status.ReadyType && c.Status == v1.ConditionTrue {
				found = true
				break
			}
		}

		if !found {
			ctx.Log.Debugw("AtlasDatabaseUser not ready - not creating connection secret", "user.name", dbUser.Name)
			continue
		}

		scopes := dbUser.GetScopes(mdbv1.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, name) {
			continue
		}

		password, err := dbUser.ReadPassword(r.Client)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}

		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    connectionStrings.Standard,
			SrvConnURL: connectionStrings.StandardSrv,
		}
		connectionsecret.FillPrivateConnStrings(connectionStrings, &data)

		ctx.Log.Debugw("Creating a connection Secret", "data", data)

		secretName, err := connectionsecret.Ensure(r.Client, project.Namespace, project.Spec.Name, project.ID(), name, data)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}
		secrets = append(secrets, secretName)
	}

	if len(secrets) > 0 {
		r.EventRecorder.Eventf(deploymentResource, "Normal", "ConnectionSecretsEnsured", "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	return workflow.OK()
}

func handleSharedDeploymentUpgrade(ctx *workflow.Context, current *mongodbatlas.Cluster, new *mongodbatlas.Cluster) (scheduled bool, _ error) {
	baseErr := "can not perform deployment upgrade. ERR: %v"
	if !deploymentShouldBeUpgraded(current, new) {
		ctx.Log.Debug("deployment shouldn't be upgraded")
		return false, nil
	}

	// Remove backingProviderName
	new.ProviderSettings.BackingProviderName = ""
	ctx.Log.Infof("performing deployment upgrade from %s, to %s",
		current.ProviderSettings.InstanceSizeName, new.ProviderSettings.InstanceSizeName)

	// TODO: Replace with the go-atlas-client when this method will be added to go-atlas-client
	atlasClient := ctx.Client
	urlStr := fmt.Sprintf("/api/atlas/v1.0/groups/%s/clusters/tenantUpgrade", current.GroupID)
	req, err := atlasClient.NewRequest(context.Background(), http.MethodPost, urlStr, new)
	if err != nil {
		return false, fmt.Errorf(baseErr, err)
	}

	_, err = atlasClient.Do(context.Background(), req, &new)
	if err != nil {
		return false, fmt.Errorf(baseErr, err)
	}

	return true, nil
}

func deploymentShouldBeUpgraded(current *mongodbatlas.Cluster, new *mongodbatlas.Cluster) bool {
	if isSharedDeployment(current.ProviderSettings.InstanceSizeName) && !isSharedDeployment(new.ProviderSettings.InstanceSizeName) {
		return true
	}
	return false
}

func isSharedDeployment(instanceSizeName string) bool {
	switch strings.ToUpper(instanceSizeName) {
	case "M0", "M2", "M5":
		return true
	}
	return false
}
