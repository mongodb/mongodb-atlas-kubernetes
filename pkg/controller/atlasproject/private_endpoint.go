package atlasproject

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func (r *AtlasProjectReconciler) ensurePrivateEndpoint(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	peList := project.Spec.DeepCopy().PrivateEndpoints
	if result := createOrDeletePEInAtlas(ctx.Client, projectID, peList, ctx.Log); !result.IsOk() {
		return result
	}

	// ctx.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(expired))
	return workflow.OK()
}

func createOrDeletePEInAtlas(client mongodbatlas.Client, projectID string, operatorPrivateEndpoints []project.PrivateEndpoint, log *zap.SugaredLogger) workflow.Result {
	provider := "AWS"
	atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.ProjectPrivateEndpointNotCreatedInAtlas, err.Error())
	}

	// Making a new slice with synonyms as Atlas IP Access list to enable usage of 'Identifiable'
	atlasPrivateEndpoints := make([]atlasProjectPrivateEndpoint, len(atlasPeConnections))
	for i, r := range atlasPeConnections {
		atlasPrivateEndpoints[i] = atlasProjectPrivateEndpoint(r)
	}

	difference := set.Difference(atlasPrivateEndpoints, operatorPrivateEndpoints)

	if err := deletePrivateEndpointsFromAtlas(client, projectID, difference, log); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}

	if result := createPrivateEndpointsInAtlas(client, projectID, operatorPrivateEndpoints); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

type atlasProjectPrivateEndpoint mongodbatlas.PrivateEndpointConnection

func (pe atlasProjectPrivateEndpoint) Identifier() interface{} {
	return pe.ProviderName + pe.Region
}

func createPrivateEndpointsInAtlas(client mongodbatlas.Client, projectID string, privateEndpoints []project.PrivateEndpoint) workflow.Result {
	for _, pe := range privateEndpoints {
		if _, _, err := client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		}); err != nil {
			return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
		}
	}

	return workflow.OK()
}

func deletePrivateEndpointsFromAtlas(client mongodbatlas.Client, projectID string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) error {
	for _, l := range listsToRemove {
		endpoint := l.(atlasProjectPrivateEndpoint)
		if _, err := client.PrivateEndpoints.Delete(context.Background(), projectID, endpoint.ProviderName, endpoint.ID); err != nil {
			return err
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", l.Identifier())
	}
	return nil
}
