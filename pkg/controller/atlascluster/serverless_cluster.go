package atlascluster

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasClusterReconciler) ensureServerlessClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) (atlasCluster *mongodbatlas.Cluster, _ workflow.Result) {
	atlasCluster, resp, err := ctx.Client.ServerlessInstances.Get(context.Background(), project.Status.ID, cluster.Spec.ClusterSpec.Name)
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

		ctx.Log.Infof("Cluster %s doesn't exist in Atlas - creating", cluster.Spec.ClusterSpec.Name)
		atlasCluster, _, err = ctx.Client.ServerlessInstances.Create(context.Background(), project.Status.ID, &mongodbatlas.ServerlessCreateRequestParams{
			Name: cluster.Spec.ClusterSpec.Name,
			ProviderSettings: &mongodbatlas.ServerlessProviderSettings{
				BackingProviderName: cluster.Spec.ClusterSpec.ProviderSettings.BackingProviderName,
				ProviderName:        string(cluster.Spec.ClusterSpec.ProviderSettings.ProviderName),
				RegionName:          cluster.Spec.ClusterSpec.ProviderSettings.RegionName,
			},
		})
		if err != nil {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}
	}

	switch atlasCluster.StateName {
	case "CREATING":
		return atlasCluster, workflow.InProgress(workflow.ClusterCreating, "cluster is provisioning")

	case "UPDATING", "REPAIRING":
		return atlasCluster, workflow.InProgress(workflow.ClusterUpdating, "cluster is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return atlasCluster, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown cluster state %q", atlasCluster.StateName))
	}
}
