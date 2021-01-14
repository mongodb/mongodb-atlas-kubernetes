package atlascluster

import (
	"context"
	"fmt"
	"net/http"

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
		return c, workflow.OK()

	case "CREATING":
		return c, workflow.InProgress(workflow.ClusterCreating, "Cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return c, workflow.InProgress(workflow.ClusterUpdating, "Cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return c, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", c.StateName))
	}
}
