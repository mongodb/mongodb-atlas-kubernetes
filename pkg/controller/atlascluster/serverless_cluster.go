package atlascluster

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasClusterReconciler) ensureServerlessClusterState(ctx *workflow.Context, project *mdbv1.AtlasProject, serverlessSpec *mdbv1.ServerlessSpec) (atlasCluster *mongodbatlas.Cluster, _ workflow.Result) {
	atlasCluster, resp, err := ctx.Client.ServerlessInstances.Get(context.Background(), project.Status.ID, serverlessSpec.Name)
	if err != nil {
		if resp == nil {
			return atlasCluster, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return atlasCluster, workflow.Terminate(workflow.ClusterNotCreatedInAtlas, err.Error())
		}

		ctx.Log.Infof("Serverless Instance %s doesn't exist in Atlas - creating", serverlessSpec.Name)
		atlasCluster, _, err = ctx.Client.ServerlessInstances.Create(context.Background(), project.Status.ID, &mongodbatlas.ServerlessCreateRequestParams{
			Name: serverlessSpec.Name,
			ProviderSettings: &mongodbatlas.ServerlessProviderSettings{
				BackingProviderName: serverlessSpec.ProviderSettings.BackingProviderName,
				ProviderName:        string(serverlessSpec.ProviderSettings.ProviderName),
				RegionName:          serverlessSpec.ProviderSettings.RegionName,
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
