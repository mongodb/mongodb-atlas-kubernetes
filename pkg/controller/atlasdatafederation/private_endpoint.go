package atlasdatafederation

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDataFederationReconciler) ensurePrivateEndpoints(ctx *workflow.Context, project *akov2.AtlasProject, dataFederation *akov2.AtlasDataFederation) workflow.Result {
	clientDF := NewClient(ctx.Client)

	projectID := project.ID()
	specPEs := dataFederation.Spec.PrivateEndpoints

	atlasPEs, err := getAllDataFederationPEs(ctx.Context, clientDF, projectID)
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

func syncPrivateEndpointsWithAtlas(ctx *workflow.Context, clientDF *DataFederationServiceOp, projectID string, specPEs, atlasPEs []akov2.DataFederationPE) workflow.Result {
	endpointsToCreate := set.Difference(specPEs, atlasPEs)
	ctx.Log.Debugw("Data Federation PEs to Create", "endpoints", endpointsToCreate)
	for _, e := range endpointsToCreate {
		endpoint := e.(akov2.DataFederationPE)
		if _, _, err := clientDF.CreateOnePrivateEndpoint(ctx.Context, projectID, endpoint); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	endpointsToDelete := set.Difference(atlasPEs, specPEs)
	ctx.Log.Debugw("Data Federation PEs to Delete", "endpoints", endpointsToDelete)
	for _, item := range endpointsToDelete {
		endpoint := item.(akov2.DataFederationPE)
		if _, _, err := clientDF.DeleteOnePrivateEndpoint(ctx.Context, projectID, endpoint.EndpointID); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	return workflow.OK()
}

func getAllDataFederationPEs(ctx context.Context, client *DataFederationServiceOp, projectID string) (endpoints []akov2.DataFederationPE, err error) {
	endpoints, _, err = client.GetAllPrivateEndpoints(ctx, projectID)
	if endpoints == nil {
		endpoints = make([]akov2.DataFederationPE, 0)
	}
	return
}
