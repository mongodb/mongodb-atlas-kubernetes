package integrations

import (
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"

	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

func (r *AtlasThirdPartyIntegrationsReconciler) release(workflowCtx *workflow.Context, integration *akov2next.AtlasThirdPartyIntegration, err error) (ctrl.Result, error) {
	if errors.Is(err, reconciler.ErrMissingKubeProject) {
		if finalizerErr := customresource.ManageFinalizer(workflowCtx.Context, r.Client, integration, customresource.UnsetFinalizer); finalizerErr != nil {
			err = errors.Join(err, finalizerErr)
		}
	}
	return r.terminate(workflowCtx, integration, workflow.NetworkPeeringNotConfigured, err)
}

func (r *AtlasThirdPartyIntegrationsReconciler) terminate(
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
