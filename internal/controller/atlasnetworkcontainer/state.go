package atlasnetworkcontainer

import (
	"context"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

type reconcileRequest struct {
	projectID        string
	networkContainer *akov2.AtlasNetworkContainer
	service          networkcontainer.NetworkContainerService
}

func (r *AtlasNetworkContainerReconciler) handleCustomResource(ctx context.Context, networkContainer *akov2.AtlasNetworkContainer) (ctrl.Result, error) {
	typeName := reflect.TypeOf(*networkContainer).Name()
	if customresource.ReconciliationShouldBeSkipped(networkContainer) {
		return r.Skip(ctx, typeName, networkContainer, networkContainer.Spec)
	}

	conditions := api.InitCondition(networkContainer, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, networkContainer)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, networkContainer)

	isValid := customresource.ValidateResourceVersion(workflowCtx, networkContainer, r.Log)
	if !isValid.IsOk() {
		return r.Invalidate(typeName, isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(networkContainer) {
		return r.Unsupport(workflowCtx, typeName)
	}

	credentials, err := r.ResolveCredentials(ctx, networkContainer)
	if err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.NetworkContainerNotConfigured, err), nil
	}
	sdkClientSet, orgID, err := r.AtlasProvider.SdkClientSet(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.NetworkContainerNotConfigured, err), nil
	}
	project, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20231115008, networkContainer, orgID)
	if err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.NetworkContainerNotConfigured, err), nil
	}
	return r.handle(workflowCtx, &reconcileRequest{
		projectID:        project.ID,
		networkContainer: networkContainer,
		service:          networkcontainer.NewNetworkContainerServiceFromClientSet(sdkClientSet),
	})
}

func (r *AtlasNetworkContainerReconciler) handle(workflowCtx *workflow.Context, req *reconcileRequest) (ctrl.Result, error) {
	return r.ready(workflowCtx, req.networkContainer, req.service)
}
