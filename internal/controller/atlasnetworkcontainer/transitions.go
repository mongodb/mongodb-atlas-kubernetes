package atlasnetworkcontainer

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	ctrl "sigs.k8s.io/controller-runtime"
)

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
