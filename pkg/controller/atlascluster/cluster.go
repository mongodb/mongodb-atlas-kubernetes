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

func ensureClusterState(ctx context.Context, wctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (c *mongodbatlas.Cluster, _ workflow.Result) {
	client, err := atlas.Client(connection, wctx.Log)
	if err != nil {
		wctx.Log.Errorf("Failed to read Atlas Connection details: %s", err.Error())
		return c, workflow.Terminate(workflow.Internal, err.Error())
	}

	c, resp, err := client.Clusters.Get(ctx, project.Status.ID, cluster.Spec.Name)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			wctx.Log.Errorf("Cannot get cluster %q: %w", cluster.Spec.Name, err)
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		c, err := cluster.Spec.Cluster()
		if err != nil {
			wctx.Log.Errorf("Cannot convert Spec to atlascluster: %w", cluster.Spec.Name, err)
			return c, workflow.Terminate(workflow.Internal, err.Error())
		}

		c, _, err = client.Clusters.Create(ctx, project.Status.ID, c)
		if err != nil {
			wctx.Log.Errorf("Cannot create cluster %q: %w", cluster.Spec.Name, err)
			return c, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch c.StateName {
	case "IDLE":
		return c, workflow.OK()

	case "CREATING":
		return c, workflow.InProgress("Cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return c, workflow.InProgress("Cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		wctx.Log.Errorf("Unknown cluster state %q", c.StateName)
		return c, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", c.StateName))
	}
}
