package atlascluster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

func (r *AtlasClusterReconciler) ensureClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (atlasCluster *mongodbatlas.Cluster, _ workflow.Result) {
	atlasCluster, resp, err := ctx.Client.Clusters.Get(context.Background(), project.Status.ID, cluster.Spec.Name)
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

		ctx.Log.Infof("Cluster %s doesn't exist in Atlas - creating", cluster.Spec.Name)
		atlasCluster, _, err = ctx.Client.Clusters.Create(context.Background(), project.Status.ID, atlasCluster)
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch atlasCluster.StateName {
	case "IDLE":
		resultingCluster, err := MergedCluster(*atlasCluster, cluster.Spec)
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		if done := ClustersEqual(ctx.Log, *atlasCluster, resultingCluster); done {
			return atlasCluster, workflow.OK()
		}

		if cluster.Spec.Paused != nil {
			if atlasCluster.Paused == nil || *atlasCluster.Paused != *cluster.Spec.Paused {
				// paused is different from Atlas
				// we need to first send a special (un)pause request before reconciling everything else
				resultingCluster = mongodbatlas.Cluster{
					Paused: cluster.Spec.Paused,
				}
			} else {
				// otherwise, don't send the paused field
				resultingCluster.Paused = nil
			}
		}

		resultingCluster = cleanupCluster(resultingCluster)

		atlasCluster, _, err = ctx.Client.Clusters.Update(context.Background(), project.Status.ID, cluster.Spec.Name, &resultingCluster)
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotUpdatedInAtlas, err.Error())
		}

		return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	case "CREATING":
		return atlasCluster, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return atlasCluster, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", atlasCluster.StateName))
	}
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
	return cluster
}

// MergedCluster will return the result of merging AtlasClusterSpec with Atlas Cluster
func MergedCluster(atlasCluster mongodbatlas.Cluster, spec mdbv1.AtlasClusterSpec) (result mongodbatlas.Cluster, err error) {
	if err = compat.JSONCopy(&result, atlasCluster); err != nil {
		return
	}

	if err = compat.JSONCopy(&result, spec); err != nil {
		return
	}

	mergeRegionConfigs(result.ReplicationSpecs, spec.ReplicationSpecs)

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
	d := cmp.Diff(clusterAtlas, clusterOperator, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Clusters are different: %s", d)
	}

	return d == ""
}
