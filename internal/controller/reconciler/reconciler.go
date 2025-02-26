package reconciler

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

type AtlasReconciler struct {
	Client          client.Client
	Log             *zap.SugaredLogger
	GlobalSecretRef client.ObjectKey
}

func (r *AtlasReconciler) Skip(ctx context.Context, typeName string, resource api.AtlasCustomResource, spec any) (ctrl.Result, error) {
	msg := fmt.Sprintf("-> Skipping %s reconciliation as annotation %s=%s",
		typeName, customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip)
	r.Log.Infow(msg, "spec", spec)
	if !resource.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, resource, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult(), nil
		}
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasReconciler) Invalidate(typeName string, invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("%T is invalid: %v", typeName, invalid)
	return invalid.ReconcileResult(), nil
}

func (r *AtlasReconciler) Unsupport(ctx *workflow.Context, typeName string) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported,
		fmt.Errorf("the %s is not supported by Atlas for government", typeName),
	).WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}
