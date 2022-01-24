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

func ensureClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasDeployment) (atlasCluster *mongodbatlas.Cluster, _ workflow.Result) {
	atlasCluster, resp, err := ctx.Client.Clusters.Get(context.Background(), project.Status.ID, cluster.Spec.DeploymentSpec.Name)
	if err != nil {
		if resp == nil {
			return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		atlasCluster, err = cluster.Spec.Cluster()
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Cluster %s doesn't exist in Atlas - creating", cluster.Spec.DeploymentSpec.Name)
		atlasCluster, _, err = ctx.Client.Clusters.Create(context.Background(), project.Status.ID, atlasCluster)
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch atlasCluster.StateName {
	case "IDLE":

		return regularClusterIdle(ctx, project, cluster, atlasCluster)
	case "CREATING":
		return atlasCluster, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	case "DELETING":
		return atlasCluster, workflow.InProgress(workflow.ClusterDeleting, "cluster is being deleted")

	case "DELETED":
		return atlasCluster, workflow.InProgress(workflow.ClusterDeleted, "cluster has been deleted")
	default:
		return atlasCluster, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", atlasCluster.StateName))
	}
}

func regularClusterIdle(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasDeployment, atlasCluster *mongodbatlas.Cluster) (*mongodbatlas.Cluster, workflow.Result) {
	resultingCluster, err := MergedCluster(*atlasCluster, cluster.Spec)
	if err != nil {
		return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
	}

	if done := ClustersEqual(ctx.Log, *atlasCluster, resultingCluster); done {
		return atlasCluster, workflow.OK()
	}

	if cluster.Spec.DeploymentSpec.Paused != nil {
		if atlasCluster.Paused == nil || *atlasCluster.Paused != *cluster.Spec.DeploymentSpec.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			resultingCluster = mongodbatlas.Cluster{
				Paused: cluster.Spec.DeploymentSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			resultingCluster.Paused = nil
		}
	}

	resultingCluster = cleanupCluster(resultingCluster)

	// Handle shared (M0,M2,M5) cluster to non-shared cluster upgrade
	scheduled, err := handleSharedClusterUpgrade(ctx, atlasCluster, &resultingCluster)
	if err != nil {
		return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
	}
	if scheduled {
		return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is upgrading")
	}

	atlasCluster, _, err = ctx.Client.Clusters.Update(context.Background(), project.Status.ID, cluster.Spec.DeploymentSpec.Name, &resultingCluster)
	if err != nil {
		return atlasCluster, workflow.Terminate(workflow.ClusterNotUpdatedInAtlas, err.Error())
	}

	return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")
}

// cleanupCluster will unset some fields that cannot be changed via API or are deprecated.
func cleanupCluster(cluster mongodbatlas.Cluster) mongodbatlas.Cluster {
	cluster.ID = ""
	cluster.MongoDBVersion = ""
	cluster.MongoURI = ""
	cluster.MongoURIUpdated = ""
	cluster.MongoURIWithOptions = ""
	cluster.SrvAddress = ""
	cluster.StateName = ""
	cluster.ReplicationFactor = nil
	cluster.ReplicationSpec = nil
	cluster.ConnectionStrings = nil
	cluster = removeOutdatedFields(&cluster, nil)

	return cluster
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
			result.ProviderSettings.AutoScaling.Compute = &mongodbatlas.Compute{}
		}

		if lookAt.AutoScaling.DiskGBEnabled != nil && *lookAt.AutoScaling.DiskGBEnabled {
			result.DiskSizeGB = nil
		}
	}

	return result
}

// MergedCluster will return the result of merging AtlasDeploymentSpec with Atlas Cluster
func MergedCluster(atlasCluster mongodbatlas.Cluster, spec mdbv1.AtlasDeploymentSpec) (result mongodbatlas.Cluster, err error) {
	if err = compat.JSONCopy(&result, atlasCluster); err != nil {
		return
	}

	if err = compat.JSONCopy(&result, spec.DeploymentSpec); err != nil {
		return
	}

	mergeRegionConfigs(result.ReplicationSpecs, spec.DeploymentSpec.ReplicationSpecs)

	// According to the docs for 'providerSettings.regionName' (https://docs.atlas.mongodb.com/reference/api/clusters-create-one/):
	// "Don't specify this parameter when creating a multi-region cluster using the replicationSpec object or a Global
	// Cluster with the replicationSpecs array."
	// The problem is that Atlas API accepts the create/update request but then returns the 'ProviderSettings.RegionName' empty in GET request
	// So we need to consider this while comparing (to avoid perpetual updates)
	if len(result.ReplicationSpecs) > 0 && atlasCluster.ProviderSettings.RegionName == "" {
		result.ProviderSettings.RegionName = ""
	}

	return
}

// mergeRegionConfigs removes replicationSpecs[i].RegionsConfigs[key] from Atlas Cluster that are absent in Operator.
// Dev idea: this could have been added into some more generic method like `JSONCopy` or something wrapping it to make
// sure any Atlas map get redundant keys removed. So far there's only one map in Cluster ('RegionsConfig') so we'll do this
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

// ClustersEqual compares two Atlas Clusters
func ClustersEqual(log *zap.SugaredLogger, clusterAtlas mongodbatlas.Cluster, clusterOperator mongodbatlas.Cluster) bool {
	clusterAtlas = removeOutdatedFields(&clusterAtlas, &clusterOperator)
	clusterOperator = removeOutdatedFields(&clusterOperator, nil)

	d := cmp.Diff(clusterAtlas, clusterOperator, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Clusters are different: %s", d)
	}

	return d == ""
}

func (r *AtlasDeploymentReconciler) ensureConnectionSecrets(ctx *workflow.Context, project *mdbv1.AtlasProject, name string, connectionStrings *mongodbatlas.ConnectionStrings, clusterResource *mdbv1.AtlasDeployment) workflow.Result {
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

		scopes := dbUser.GetScopes(mdbv1.ClusterScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, name) {
			continue
		}

		password, err := dbUser.ReadPassword(r.Client)
		if err != nil {
			return workflow.Terminate(workflow.ClusterConnectionSecretsNotCreated, err.Error())
		}

		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			ConnURL:    connectionStrings.Standard,
			SrvConnURL: connectionStrings.StandardSrv,
			Password:   password,
		}

		secretName, err := connectionsecret.Ensure(r.Client, project.Namespace, project.Spec.Name, project.ID(), name, data)
		if err != nil {
			return workflow.Terminate(workflow.ClusterConnectionSecretsNotCreated, err.Error())
		}
		secrets = append(secrets, secretName)
	}

	if len(secrets) > 0 {
		r.EventRecorder.Eventf(clusterResource, "Normal", "ConnectionSecretsEnsured", "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	return workflow.OK()
}

func handleSharedClusterUpgrade(ctx *workflow.Context, current *mongodbatlas.Cluster, new *mongodbatlas.Cluster) (scheduled bool, _ error) {
	baseErr := "can not perform cluster upgrade. ERR: %v"
	if !clusterShouldBeUpgraded(current, new) {
		ctx.Log.Debug("cluster shouldn't be upgraded")
		return false, nil
	}

	// Remove backingProviderName
	new.ProviderSettings.BackingProviderName = ""
	ctx.Log.Infof("performing cluster upgrade from %s, to %s",
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

func clusterShouldBeUpgraded(current *mongodbatlas.Cluster, new *mongodbatlas.Cluster) bool {
	if isSharedCluster(current.ProviderSettings.InstanceSizeName) && !isSharedCluster(new.ProviderSettings.InstanceSizeName) {
		return true
	}
	return false
}

func isSharedCluster(instanceSizeName string) bool {
	switch strings.ToUpper(instanceSizeName) {
	case "M0", "M2", "M5":
		return true
	}
	return false
}
