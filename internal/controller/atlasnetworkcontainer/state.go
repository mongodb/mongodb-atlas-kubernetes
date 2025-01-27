package atlasnetworkcontainer

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func (r *AtlasNetworkContainerReconciler) handle(ctx context.Context, ipAccessList *akov2.AtlasNetworkContainer) (ctrl.Result, error) {
	return reconcile.Result{}, nil
}
