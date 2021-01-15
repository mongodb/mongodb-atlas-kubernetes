package atlascluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.mongodb.org/atlas/mongodbatlas"
)

func ensureClusterState(wctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (c *mongodbatlas.Cluster, _ workflow.Result) {
	ctx := context.Background()

	client, err := atlas.Client(connection, wctx.Log)
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

		c, err := cluster.Spec.Cluster()
		if err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		c, _, err = client.Clusters.Create(ctx, project.Status.ID, c)
		if err != nil {
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch c.StateName {
	case "IDLE":
		spec, err := cluster.Spec.Cluster()
		if err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		if done, err := clusterMatchesSpec(wctx, c, cluster.Spec); err != nil {
			return c, workflow.Terminate(workflow.Internal, err.Error())
		} else if done {
			return c, workflow.OK()
		}

		c, _, err = client.Clusters.Update(ctx, project.Status.ID, cluster.Spec.Name, spec)
		if err != nil {
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		return c, workflow.InProgress(workflow.ClusterNotUpToDate, "cluster update started")

	case "CREATING":
		return c, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return c, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return c, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", c.StateName))
	}
}

// clusterMatchesSpec will merge everything from the Spec into existing Cluster and use that to detect change
func clusterMatchesSpec(ctx *workflow.Context, cluster *mongodbatlas.Cluster, spec mdbv1.AtlasClusterSpec) (bool, error) {
	clusterMerged := mongodbatlas.Cluster{}
	if err := jsonCopy(&clusterMerged, cluster); err != nil {
		return false, err
	}

	if err := jsonCopy(&clusterMerged, spec); err != nil {
		return false, err
	}

	d := cmp.Diff(*cluster, clusterMerged)
	if d != "" {
		ctx.Log.Debugf("Cluster differs from spec: %s", d)
	}

	return d == "", nil
}

func jsonCopy(dst, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dst)
	if err != nil {
		return err
	}

	return nil
}
