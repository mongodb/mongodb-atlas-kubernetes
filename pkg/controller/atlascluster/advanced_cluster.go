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

func (r *AtlasClusterReconciler) ensureAdvancedClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	advancedClusterSpec := cluster.Spec.AdvancedClusterSpec

	advancedCluster, resp, err := ctx.Client.AdvancedClusters.Get(context.Background(), project.Status.ID, advancedClusterSpec.Name)

	if err != nil {
		if resp == nil {
			return advancedCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return advancedCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		advancedCluster, err = advancedClusterSpec.AdvancedCluster()
		if err != nil {
			return advancedCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Advanced Cluster %s doesn't exist in Atlas - creating", advancedClusterSpec.Name)
		advancedCluster, _, err = ctx.Client.AdvancedClusters.Create(context.Background(), project.Status.ID, advancedCluster)
		if err != nil {
			return advancedCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch advancedCluster.StateName {
	case "IDLE":
		return advancedClusterIdle(ctx, project, cluster, advancedCluster)

	case "CREATING":
		return advancedCluster, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return advancedCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return advancedCluster, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", advancedCluster.StateName))
	}
}

func advancedClusterIdle(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, advancedCluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	resultingCluster, err := MergedAdvancedCluster(*advancedCluster, cluster.Spec)
	if err != nil {
		return advancedCluster, workflow.Terminate(workflow.Internal, err.Error())
	}

	if done := AdvancedClustersEqual(ctx.Log, *advancedCluster, resultingCluster); done {
		return advancedCluster, workflow.OK()
	}

	if cluster.Spec.AdvancedClusterSpec.Paused != nil {
		if advancedCluster.Paused == nil || *advancedCluster.Paused != *cluster.Spec.AdvancedClusterSpec.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			resultingCluster = mongodbatlas.AdvancedCluster{
				Paused: cluster.Spec.AdvancedClusterSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			resultingCluster.Paused = nil
		}
	}

	resultingCluster = cleanupAdvancedCluster(resultingCluster)

	advancedCluster, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, cluster.Spec.AdvancedClusterSpec.Name, &resultingCluster)
	if err != nil {
		return advancedCluster, workflow.Terminate(workflow.ClusterNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")
}

func cleanupAdvancedCluster(cluster mongodbatlas.AdvancedCluster) mongodbatlas.AdvancedCluster {
	cluster.ID = ""
	cluster.MongoDBVersion = ""
	cluster.StateName = ""
	cluster.ConnectionStrings = nil
	return cluster
}

// MergedAdvancedCluster will return the result of merging AtlasClusterSpec with Atlas Advanced Cluster
func MergedAdvancedCluster(advancedCluster mongodbatlas.AdvancedCluster, spec mdbv1.AtlasClusterSpec) (mongodbatlas.AdvancedCluster, error) {
	result := mongodbatlas.AdvancedCluster{}
	if err := compat.JSONCopy(&result, advancedCluster); err != nil {
		return result, err
	}

	if err := compat.JSONCopy(&result, spec.AdvancedClusterSpec); err != nil {
		return result, err
	}

	for i, replicationSpec := range advancedCluster.ReplicationSpecs {
		for k, v := range replicationSpec.RegionConfigs {
			// the response does not return backing provider names in some situations.
			// if this is the case, we want to strip these fields so they do not cause a bad comparison.
			if v.BackingProviderName == "" {
				result.ReplicationSpecs[i].RegionConfigs[k].BackingProviderName = ""
			}
		}
	}
	return result, nil
}

// AdvancedClustersEqual compares two Atlas Advanced Clusters
func AdvancedClustersEqual(log *zap.SugaredLogger, clusterAtlas mongodbatlas.AdvancedCluster, clusterOperator mongodbatlas.AdvancedCluster) bool {
	d := cmp.Diff(clusterAtlas, clusterOperator, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Clusters are different: %s", d)
	}

	return d == ""
}

// GetAllClusterNames returns all cluster names including regular and advanced clusters.
func GetAllClusterNames(client mongodbatlas.Client, projectID string) ([]string, error) {
	var clusterNames []string
	clusters, _, err := client.Clusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	advancedClusters, _, err := client.AdvancedClusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, c := range clusters {
		clusterNames = append(clusterNames, c.Name)
	}

	for _, c := range advancedClusters.Results {
		// based on configuration settings, some advanced clusters also show up in the regular clusters API.
		// For these clusters, we don't want to duplicate the secret so we skip them.
		found := false
		for _, regularCluster := range clusters {
			if regularCluster.Name == c.Name {
				found = true
				break
			}
		}

		// we only include cluster names which have not been handled by the regular cluster API.
		if !found {
			clusterNames = append(clusterNames, c.Name)
		}
	}
	return clusterNames, nil
}
