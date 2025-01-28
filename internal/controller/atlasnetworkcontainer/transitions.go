package atlasnetworkcontainer

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

func (r *AtlasNetworkContainerReconciler) ready(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer, service networkcontainer.NetworkContainerService) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotSet, err), nil
	}

	workflowCtx.SetConditionTrue(api.NetworkContainerReady).
		SetConditionTrue(api.ReadyType)
		// TODO: add .EnsureStatusOption(networkContainer.NewNetworkContainerStatus(service)) ?

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
