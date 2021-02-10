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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

func (r *AtlasClusterReconciler) ensureClusterState(log *zap.SugaredLogger, connection atlas.Connection, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (c *mongodbatlas.Cluster, _ workflow.Result) {
	ctx := context.Background()

	client, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return c, workflow.Terminate(workflow.Internal, err.Error())
	}

	c, resp, err := client.Clusters.Get(ctx, project.Status.ID, cluster.Spec.Name)
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

		log.Infof("Cluster %s doesn't exist in Atlas - creating", cluster.Spec.Name)
		c, _, err = client.Clusters.Create(ctx, project.Status.ID, c)
		if err != nil {
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch c.StateName {
	case "IDLE":
		if done, err := clusterMatchesSpec(log, c, cluster.Spec); err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		} else if done {
			return c, workflow.OK()
		}

		spec, err := cluster.Spec.Cluster()
		if err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		if cluster.Spec.Paused != nil {
			if c.Paused == nil || *c.Paused != *cluster.Spec.Paused {
				// paused is different from Atlas
				// we need to first send a special (un)pause request before reconciling everything else
				spec = &mongodbatlas.Cluster{
					Paused: cluster.Spec.Paused,
				}
			} else {
				// otherwise, don't send the paused field
				spec.Paused = nil
			}
		}

		c, _, err = client.Clusters.Update(ctx, project.Status.ID, cluster.Spec.Name, spec)
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

// clusterMatchesSpec will merge everything from the Spec into existing Cluster and use that to detect change.
// Direct comparison is not feasible because Atlas will set a lot of fields to default values, so we need to apply our changes on top of that.
func clusterMatchesSpec(log *zap.SugaredLogger, cluster *mongodbatlas.Cluster, spec mdbv1.AtlasClusterSpec) (bool, error) {
	clusterMerged := mongodbatlas.Cluster{}
	if err := compat.JSONCopy(&clusterMerged, cluster); err != nil {
		return false, err
	}

	if err := compat.JSONCopy(&clusterMerged, spec); err != nil {
		return false, err
	}

	d := cmp.Diff(*cluster, clusterMerged, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Cluster differs from spec: %s", d)
	}

	return d == "", nil
}
