package atlasnetworkcontainer

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

func (r *AtlasNetworkContainerReconciler) inProgress(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotSet, err), nil
	}
	result := workflow.InProgress(
		workflow.NetworkContainerProvisioning,
		fmt.Sprintf("Network Container %s is being provisioned", container.ID),
	)
	workflowCtx.SetConditionFalse(api.ReadyType).SetConditionFromResult(api.ReadyType, result).
		EnsureStatusOption(updateNetworkContainerStatusOption(container))

	return result.ReconcileResult(), nil
}

func (r *AtlasNetworkContainerReconciler) unmanage(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.UnsetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotRemoved, err), nil
	}
	return workflow.Deleted().ReconcileResult(), nil
}

func (r *AtlasNetworkContainerReconciler) ready(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotSet, err), nil
	}

	workflowCtx.SetConditionTrueMsg(api.NetworkContainerReady, fmt.Sprintf("Network Container %s is ready", container.ID)).
		SetConditionTrue(api.ReadyType).EnsureStatusOption(updateNetworkContainerStatusOption(container))

	if networkContainer.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult(), nil
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasNetworkContainerReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) ctrl.Result {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}

func updateNetworkContainerStatusOption(container *networkcontainer.NetworkContainer) status.AtlasNetworkContainerStatusOption {
	return func(containerStatus *status.AtlasNetworkContainerStatus) {
		networkcontainer.ApplyNetworkContainerStatus(containerStatus, container)
	}
}
