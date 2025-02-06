package atlasnetworkpeering

import (
	"context"
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

const (
	typeName = "AtlasNetworkPeering"
)

type reconcileRequest struct {
	service          networkpeering.NetworkPeeringService
	containerService networkcontainer.NetworkContainerService
	projectID        string
	networkPeering   *akov2.AtlasNetworkPeering
}

func (r *AtlasNetworkPeeringReconciler) handleCustomResource(ctx context.Context, networkPeering *akov2.AtlasNetworkPeering) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(networkPeering) {
		return r.Skip(ctx, typeName, networkPeering, &networkPeering.Spec)
	}

	conditions := api.InitCondition(networkPeering, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, networkPeering)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, networkPeering)

	isValid := customresource.ValidateResourceVersion(workflowCtx, networkPeering, r.Log)
	if !isValid.IsOk() {
		return r.Invalidate(typeName, isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(networkPeering) {
		return r.Unsupport(workflowCtx, typeName)
	}

	credentials, err := r.ResolveCredentials(ctx, networkPeering)
	if err != nil {
		return r.release(workflowCtx, networkPeering, err)
	}
	sdkClientSet, _, err := r.AtlasProvider.SdkClientSet(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, networkPeering, workflow.NetworkPeeringNotConfigured, err)
	}
	project, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20231115008, networkPeering)
	if err != nil {
		return r.release(workflowCtx, networkPeering, err)
	}
	return r.handle(workflowCtx, &reconcileRequest{
		service:          networkpeering.NewNetworkPeeringServiceFromClientSet(sdkClientSet),
		containerService: networkcontainer.NewNetworkContainerServiceFromClientSet(sdkClientSet),
		projectID:        project.ID,
		networkPeering:   networkPeering,
	})
}

func (r *AtlasNetworkPeeringReconciler) handle(workflowCtx *workflow.Context, req *reconcileRequest) (ctrl.Result, error) {
	r.Log.Infow("handling network peering reconcile request",
		"service set", (req.service != nil), "projectID", req.projectID, "networkPeering", req.networkPeering)
	container, err := r.getContainer(workflowCtx.Context, req)
	if err != nil {
		return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	if container == nil {
		err := fmt.Errorf("container not found for reference %v", req.networkPeering.Spec.ContainerRef)
		return r.terminate(workflowCtx, req.networkPeering, workflow.NetworkPeeringMissingContainer, err)
	}
	var atlasPeer *networkpeering.NetworkPeer
	if req.networkPeering.Status.ID != "" {
		peer, err := req.service.Get(workflowCtx.Context, req.projectID, req.networkPeering.Status.ID)
		if err != nil && !errors.Is(err, networkpeering.ErrNotFound) {
			return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, err)
		}
		atlasPeer = peer
	}
	inAtlas := atlasPeer != nil
	deleted := req.networkPeering.DeletionTimestamp != nil

	switch {
	case !deleted && !inAtlas:
		return r.create(workflowCtx, req, container)
	case !deleted && inAtlas:
		return r.sync(workflowCtx, req, atlasPeer, container)
	case deleted && inAtlas:
		return r.delete(workflowCtx, req, atlasPeer, container)
	default: // deleted && !inAtlas
		return r.unmanage(workflowCtx, req)
	}
}

func (r *AtlasNetworkPeeringReconciler) getContainer(ctx context.Context, req *reconcileRequest) (*networkcontainer.NetworkContainer, error) {
	id := req.networkPeering.Spec.ContainerRef.ID
	if req.networkPeering.Spec.ContainerRef.ID == "" { // Name should be non nil instead
		var err error
		id, err = getContainerIDFromKubernetes(ctx, r.Client, req.networkPeering)
		if err != nil {
			return nil, fmt.Errorf("failed to solve Network Container id from Kubernetes: %w", err)
		}
		if id == "" {
			return nil, fmt.Errorf("container %s has no id, is it still to be created?",
				req.networkPeering.Spec.ContainerRef.Name)
		}
	}
	container, err := req.containerService.Get(ctx, req.projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Network Container %s from Atlas by id: %w", id, err)
	}
	return container, nil
}

func getContainerIDFromKubernetes(ctx context.Context, k8sClient client.Client, networkPeering *akov2.AtlasNetworkPeering) (string, error) {
	k8sContainer := akov2.AtlasNetworkContainer{}
	key := client.ObjectKey{
		Name:      networkPeering.Spec.ContainerRef.Name,
		Namespace: networkPeering.Namespace,
	}
	err := k8sClient.Get(ctx, key, &k8sContainer)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the Kubernetes Network Container %s info: %w", key.Name, err)
	}
	id := k8sContainer.Spec.ID
	if id == "" {
		id = k8sContainer.Status.ID
	}
	return id, nil
}
