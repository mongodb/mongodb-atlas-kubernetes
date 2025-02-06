package atlasnetworkpeering

import (
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

func (r *AtlasNetworkPeeringReconciler) create(workflowCtx *workflow.Context, req *reconcileRequest, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	newPeer, err := req.service.Create(
		workflowCtx.Context,
		req.projectID,
		container.ID,
		&req.networkPeering.Spec.AtlasNetworkPeeringConfig,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create peering connection: %w", err)
		return r.terminate(workflowCtx, req.networkPeering, workflow.NetworkPeeringNotConfigured, wrappedErr)
	}
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, req.networkPeering, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, req.networkPeering, workflow.AtlasFinalizerNotSet, err)
	}
	return r.inProgress(workflowCtx, workflow.NetworkPeeringConnectionCreating, newPeer, container)
}

func (r *AtlasNetworkPeeringReconciler) sync(workflowCtx *workflow.Context, req *reconcileRequest, atlasPeer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	switch {
	case atlasPeer.Failed():
		err := fmt.Errorf("peering connection failed: %s", atlasPeer.ErrorMessage)
		return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, err)
	case !atlasPeer.Available():
		return r.inProgress(workflowCtx, workflow.NetworkPeeringConnectionPending, atlasPeer, container)
	}
	specPeer := networkpeering.NewNetworkPeer(atlasPeer.ID, &req.networkPeering.Spec.AtlasNetworkPeeringConfig)
	if !networkpeering.CompareConfigs(atlasPeer, specPeer) {
		return r.update(workflowCtx, req, container)
	}
	return r.ready(workflowCtx, req, atlasPeer, container)
}

func (r *AtlasNetworkPeeringReconciler) update(workflowCtx *workflow.Context, req *reconcileRequest, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	updatedPeer, err := req.service.Update(workflowCtx.Context, req.projectID, req.networkPeering.Status.ID, container.ID, &req.networkPeering.Spec.AtlasNetworkPeeringConfig)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update peering connection: %w", err)
		return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
	}
	return r.inProgress(workflowCtx, workflow.NetworkPeeringConnectionUpdating, updatedPeer, container)
}

func (r *AtlasNetworkPeeringReconciler) delete(workflowCtx *workflow.Context, req *reconcileRequest, atlasPeer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	id := req.networkPeering.Status.ID
	peer := atlasPeer
	if id != "" && !atlasPeer.Closing() {
		if err := req.service.Delete(workflowCtx.Context, req.projectID, id); err != nil {
			wrappedErr := fmt.Errorf("failed to delete peer connection %s: %w", id, err)
			return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
		}
		closingPeer, err := req.service.Get(workflowCtx.Context, req.projectID, id)
		if err != nil && !errors.Is(err, networkpeering.ErrNotFound) {
			wrappedErr := fmt.Errorf("failed to get closing peer connection %s: %w", id, err)
			return r.terminate(workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
		}
		peer = closingPeer
	}
	if peer == nil {
		return r.unmanage(workflowCtx, req)
	}
	return r.inProgress(workflowCtx, workflow.NetworkPeeringConnectionClosing, peer, container)
}

func (r *AtlasNetworkPeeringReconciler) unmanage(workflowCtx *workflow.Context, req *reconcileRequest) (ctrl.Result, error) {
	workflowCtx.EnsureStatusOption(clearPeeringStatusOption())
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, req.networkPeering, customresource.UnsetFinalizer); err != nil {
		return r.terminate(workflowCtx, req.networkPeering, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.Deleted().ReconcileResult(), nil
}

func (r *AtlasNetworkPeeringReconciler) inProgress(workflowCtx *workflow.Context, reason workflow.ConditionReason, peer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	statusMsg := fmt.Sprintf("Network Peering Connection %s is %s", peer.ID, peer.Status)
	workflowCtx.EnsureStatusOption(updatePeeringStatusOption(peer, container))
	workflowCtx.SetConditionFalseMsg(api.NetworkPeerReadyType, statusMsg)
	workflowCtx.SetConditionFalse(api.ReadyType)

	return workflow.InProgress(reason, statusMsg).ReconcileResult(), nil
}

func (r *AtlasNetworkPeeringReconciler) ready(workflowCtx *workflow.Context, req *reconcileRequest, peer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, req.networkPeering, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, req.networkPeering, workflow.AtlasFinalizerNotSet, err)
	}

	workflowCtx.EnsureStatusOption(updatePeeringStatusOption(peer, container))
	workflowCtx.SetConditionTrue(api.NetworkPeerReadyType)
	workflowCtx.SetConditionTrue(api.ReadyType)

	if req.networkPeering.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult(), nil
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasNetworkPeeringReconciler) release(workflowCtx *workflow.Context, networkPeering *akov2.AtlasNetworkPeering, err error) (ctrl.Result, error) {
	if errors.Is(err, reconciler.ErrMissingKubeProject) {
		if finalizerErr := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkPeering, customresource.UnsetFinalizer); finalizerErr != nil {
			err = errors.Join(err, finalizerErr)
		}
	}
	return r.terminate(workflowCtx, networkPeering, workflow.NetworkPeeringNotConfigured, err)
}

func (r *AtlasNetworkPeeringReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) (ctrl.Result, error) {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).SetConditionFromResult(condition, result)

	return result.ReconcileResult(), nil
}

func updatePeeringStatusOption(peer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		applyPeeringStatus(peeringStatus, peer, container)
	}
}

func applyPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus, peer *networkpeering.NetworkPeer, container *networkcontainer.NetworkContainer) {
	switch provider.ProviderName(peer.Provider) {
	case provider.ProviderAWS:
		if peer.AWSStatus != nil {
			if peeringStatus.AWSStatus == nil {
				peeringStatus.AWSStatus = &status.AWSPeeringStatus{}
			}
			peeringStatus.AWSStatus.ConnectionID = peer.AWSStatus.ConnectionID
			peeringStatus.AWSStatus.VpcID = container.AWSStatus.VpcID
		}
	case provider.ProviderAzure:
		if container.AzureStatus != nil {
			if peeringStatus.AzureStatus == nil {
				peeringStatus.AzureStatus = &status.AzurePeeringStatus{}
			}
			peeringStatus.AzureStatus.AzureSubscriptionID = container.AzureStatus.AzureSubscriptionID
			peeringStatus.AzureStatus.VnetName = container.AzureStatus.VnetName
		}
	case provider.ProviderGCP:
		if container.GCPStatus != nil {
			if peeringStatus.GCPStatus == nil {
				peeringStatus.GCPStatus = &status.GCPPeeringStatus{}
			}
			peeringStatus.GCPStatus.GCPProjectID = container.GCPStatus.GCPProjectID
			peeringStatus.GCPStatus.NetworkName = container.GCPStatus.NetworkName
		}
	default:
		peeringStatus.Status = fmt.Sprintf("unsupported provider: %q", peer.Provider)
		return
	}
	peeringStatus.ID = peer.ID
	peeringStatus.Status = peer.Status
}

func clearPeeringStatusOption() status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		clearPeeringStatus(peeringStatus)
	}
}

func clearPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus) {
	peeringStatus.ID = ""
	peeringStatus.Status = ""
	if peeringStatus.AWSStatus != nil {
		peeringStatus.AWSStatus.ConnectionID = ""
	}
}
