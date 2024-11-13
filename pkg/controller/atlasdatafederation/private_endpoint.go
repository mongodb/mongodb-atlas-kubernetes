package atlasdatafederation

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compare"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDataFederationReconciler) ensurePrivateEndpoints(ctx *workflow.Context, projectID string, dataFederation *akov2.AtlasDataFederation) workflow.Result {
	r.privateEndpointService = datafederation.NewPrivateEndpointService(ctx.SdkClient.DataFederationApi)

	if err := r.reconcilePrivateEndpoints(ctx, projectID, dataFederation.Spec.PrivateEndpoints); err != nil {
		ctx.SetConditionFalseMsg(api.DataFederationPEReadyType, err.Error())

		return workflow.Terminate(workflow.Internal, err.Error())
	}

	ctx.SetConditionTrue(api.DataFederationPEReadyType)

	if len(dataFederation.Spec.PrivateEndpoints) == 0 {
		ctx.UnsetCondition(api.DataFederationPEReadyType)
	}

	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) reconcilePrivateEndpoints(ctx *workflow.Context, projectID string, privateEndpoints []akov2.DataFederationPE) error {
	specPrivateEndpoints, err := datafederation.NewPrivateEndpoints(projectID, privateEndpoints)
	if err != nil {
		return fmt.Errorf("failed to parse private endpoint specifications: %w", err)
	}

	atlasPrivateEndpoints, err := r.privateEndpointService.List(ctx.Context, projectID)
	if err != nil {
		return err
	}

	for _, specPrivateEndpoint := range specPrivateEndpoints {
		if compare.ContainsDeepEqual(atlasPrivateEndpoints, specPrivateEndpoint) {
			continue
		}

		if err = r.privateEndpointService.Create(ctx.Context, &specPrivateEndpoint); err != nil {
			return err
		}
	}

	for _, atlasPrivateEndpoint := range atlasPrivateEndpoints {
		if compare.ContainsDeepEqual(specPrivateEndpoints, atlasPrivateEndpoint) {
			continue
		}

		if err = r.privateEndpointService.Delete(ctx.Context, &atlasPrivateEndpoint); err != nil {
			return err
		}
	}

	return nil
}
