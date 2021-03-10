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

func (r *AtlasClusterReconciler) ensureClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (c *mongodbatlas.Cluster, _ workflow.Result) {
	c, resp, err := ctx.Client.Clusters.Get(context.Background(), project.Status.ID, cluster.Spec.Name)
	if err != nil {
		if resp == nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		c, err = cluster.Spec.Cluster()
		if err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Cluster %s doesn't exist in Atlas - creating", cluster.Spec.Name)
		c, _, err = ctx.Client.Clusters.Create(context.Background(), project.Status.ID, c)
		if err != nil {
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch c.StateName {
	case "IDLE":
		resultingCluster, err := mergedCluster(*c, cluster.Spec)
		if err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		if done := clustersEqual(ctx.Log, *c, resultingCluster); done {
			return c, workflow.OK()
		}

		if cluster.Spec.Paused != nil {
			if c.Paused == nil || *c.Paused != *cluster.Spec.Paused {
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

		c, _, err = ctx.Client.Clusters.Update(context.Background(), project.Status.ID, cluster.Spec.Name, &resultingCluster)
		if err != nil {
			return c, workflow.Terminate(workflow.ClusterNotUpdatedInAtlas, err.Error())
		}

		return c, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	case "CREATING":
		return c, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return c, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return c, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", c.StateName))
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

// mergedCluster will return the result of merging AtlasClusterSpec with Atlas Cluster
func mergedCluster(cluster mongodbatlas.Cluster, spec mdbv1.AtlasClusterSpec) (result mongodbatlas.Cluster, err error) {
	if err = compat.JSONCopy(&result, cluster); err != nil {
		return
	}

	if err = compat.JSONCopy(&result, spec); err != nil {
		return
	}

	// TODO: might need to do this with other slices
	if err = compat.JSONSliceMerge(&result.ReplicationSpecs, cluster.ReplicationSpecs); err != nil {
		return
	}

	if err = compat.JSONSliceMerge(&result.ReplicationSpecs, spec.ReplicationSpecs); err != nil {
		return
	}

	return
}

// clustersEqual compares two Atlas Clusters
func clustersEqual(log *zap.SugaredLogger, clusterA mongodbatlas.Cluster, clusterB mongodbatlas.Cluster) bool {
	d := cmp.Diff(clusterA, clusterB, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Clusters are different: %s", d)
	}

	return d == ""
}
