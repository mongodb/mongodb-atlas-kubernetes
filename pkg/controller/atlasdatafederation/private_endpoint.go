package atlasdatafederation

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDataFederationReconciler) ensurePrivateEndpoints(ctx *workflow.Context, project *akov2.AtlasProject, dataFederation *akov2.AtlasDataFederation, service datafederation.DatafederationPrivateEndpointService) workflow.Result {
	projectID := project.ID()
	specPEs := make([]*datafederation.DatafederationPrivateEndpointEntry, 0, len(dataFederation.Spec.PrivateEndpoints))
	for _, pe := range dataFederation.Spec.PrivateEndpoints {
		specPEs = append(specPEs, datafederation.NewDatafederationPrivateEndpointEntry(&pe, projectID))
	}

	//NewDatafederationPrivateEndpointEntry
	atlasPEs, err := getAllDataFederationPEs(ctx.Context, service, projectID)
	if err != nil {
		ctx.Log.Debugw("getAllDataFederationPEs error", "err", err.Error())
	}

	result := syncPrivateEndpointsWithAtlas(ctx, service, projectID, specPEs, atlasPEs)
	if !result.IsOk() {
		ctx.SetConditionFromResult(api.DataFederationPEReadyType, result)
		return result
	}

	return workflow.OK()
}

func syncPrivateEndpointsWithAtlas(ctx *workflow.Context, service datafederation.DatafederationPrivateEndpointService, projectID string, specPEs, atlasPEs []*datafederation.DatafederationPrivateEndpointEntry) workflow.Result {
	endpointsToCreate := set.Difference(specPEs, atlasPEs)
	ctx.Log.Debugw("Data Federation PEs to Create", "endpoints", endpointsToCreate)
	for _, e := range endpointsToCreate {
		endpoint := e.(*datafederation.DatafederationPrivateEndpointEntry)
		if err := service.Create(ctx.Context, endpoint); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	endpointsToDelete := set.Difference(atlasPEs, specPEs)
	ctx.Log.Debugw("Data Federation PEs to Delete", "endpoints", endpointsToDelete)
	for _, item := range endpointsToDelete {
		endpoint := item.(*datafederation.DatafederationPrivateEndpointEntry)
		if err := service.Delete(ctx.Context, endpoint); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	return workflow.OK()
}

func getAllDataFederationPEs(ctx context.Context, service datafederation.DatafederationPrivateEndpointService, projectID string) (endpoints []*datafederation.DatafederationPrivateEndpointEntry, err error) {
	endpoints, err = service.List(ctx, projectID)
	if endpoints == nil {
		endpoints = make([]*datafederation.DatafederationPrivateEndpointEntry, 0)
	}
	return
}
