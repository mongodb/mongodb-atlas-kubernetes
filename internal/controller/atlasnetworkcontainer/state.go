package atlasnetworkcontainer

import (
	"context"
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

const (
	typeName = "AtlasNetworkContainer"
)

type reconcileRequest struct {
	projectID        string
	networkContainer *akov2.AtlasNetworkContainer
	service          networkcontainer.NetworkContainerService
}

func (r *AtlasNetworkContainerReconciler) handleCustomResource(ctx context.Context, networkContainer *akov2.AtlasNetworkContainer) (ctrl.Result, error) {
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
		return r.release(workflowCtx, networkContainer, err), nil
	}
	sdkClientSet, _, err := r.AtlasProvider.SdkClientSet(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.NetworkContainerNotConfigured, err), nil
	}
	project, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20231115008, networkContainer)
	if err != nil {
		return r.release(workflowCtx, networkContainer, err), nil
	}
	return r.handle(workflowCtx, &reconcileRequest{
		projectID:        project.ID,
		networkContainer: networkContainer,
		service:          networkcontainer.NewNetworkContainerServiceFromClientSet(sdkClientSet),
	})
}

func (r *AtlasNetworkContainerReconciler) handle(workflowCtx *workflow.Context, req *reconcileRequest) (ctrl.Result, error) {
	atlasContainer, err := discover(workflowCtx.Context, req)
	if err != nil {
		return r.terminate(workflowCtx, req.networkContainer, workflow.NetworkContainerNotConfigured, err), nil
	}
	inAtlas := atlasContainer != nil
	deleted := req.networkContainer.DeletionTimestamp != nil
	switch {
	case !deleted && !inAtlas:
		return r.create(workflowCtx, req)
	case !deleted && inAtlas:
		return r.sync(workflowCtx, req, atlasContainer)
	case deleted && inAtlas:
		return r.delete(workflowCtx, req, atlasContainer)
	default: // deleted && !inAtlas:
		return r.unmanage(workflowCtx, req.networkContainer)
	}
}

func discover(ctx context.Context, req *reconcileRequest) (*networkcontainer.NetworkContainer, error) {
	id := req.networkContainer.Spec.ID
	if id == "" {
		id = req.networkContainer.Status.ID
	}
	if id != "" {
		container, err := req.service.Get(ctx, req.projectID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get container %s from project %s: %w", id, req.projectID, err)
		}
		return container, nil
	}
	cfg := networkcontainer.NewNetworkContainerConfig(
		req.networkContainer.Spec.Provider, &req.networkContainer.Spec.AtlasNetworkContainerConfig)
	container, err := req.service.Find(ctx, req.projectID, cfg)
	if err != nil && !errors.Is(err, networkcontainer.ErrNotFound) {
		return nil, fmt.Errorf("failed to find container from project %s: %w", req.projectID, err)
	}
	return container, nil
}
