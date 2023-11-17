package atlasdatafederation

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/set"
)

func (r *AtlasDataFederationReconciler) ensurePrivateEndpoints(ctx *workflow.Context, project *mdbv1.AtlasProject, dataFederation *mdbv1.AtlasDataFederation) workflow.Result {
	clientDF := NewClient(ctx.Client, r.AtlasDomain)

	projectID := project.ID()
	specPEs := dataFederation.Spec.PrivateEndpoints

	atlasPEs, err := getAllDataFederationPEs(clientDF, projectID)
	if err != nil {
		ctx.Log.Debugw("getAllDataFederationPEs error", "err", err.Error())
	}

	result := syncPrivateEndpointsWithAtlas(ctx, clientDF, projectID, specPEs, atlasPEs)
	if !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationPEReadyType, result)
		return result
	}

	return workflow.OK()
}

func syncPrivateEndpointsWithAtlas(ctx *workflow.Context, clientDF *DataFederationServiceOp, projectID string, specPEs, atlasPEs []mdbv1.DataFederationPE) workflow.Result {
	endpointsToCreate := set.Difference(specPEs, atlasPEs)
	ctx.Log.Debugw("Data Federation PEs to Create", "endpoints", endpointsToCreate)
	for _, e := range endpointsToCreate {
		endpoint := e.(mdbv1.DataFederationPE)
		if _, _, err := clientDF.CreateOnePrivateEndpoint(context.Background(), projectID, endpoint); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	endpointsToDelete := set.Difference(atlasPEs, specPEs)
	ctx.Log.Debugw("Data Federation PEs to Delete", "endpoints", endpointsToDelete)
	for _, item := range endpointsToDelete {
		endpoint := item.(mdbv1.DataFederationPE)
		if _, _, err := clientDF.DeleteOnePrivateEndpoint(context.Background(), projectID, endpoint.EndpointID); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	return workflow.OK()
}

func getAllDataFederationPEs(client *DataFederationServiceOp, projectID string) (endpoints []mdbv1.DataFederationPE, err error) {
	endpoints, _, err = client.GetAllPrivateEndpoints(context.Background(), projectID)
	if endpoints == nil {
		endpoints = make([]mdbv1.DataFederationPE, 0)
	}
	return
}
