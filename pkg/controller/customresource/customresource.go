package customresource

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// GetResource queries the Custom Resource 'request.NamespacedName' and populates the 'resource' pointer.
// Note the logic: any reconcile result different from nil should be considered as "terminal" and will stop reconciliation
// right away (the pointer will be empty). Otherwise the pointer 'resource' will always reference the existing resource
func GetResource(client client.Client, request reconcile.Request, resource runtime.Object, log *zap.SugaredLogger) *reconcile.Result {
	err := client.Get(context.Background(), request.NamespacedName, resource)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Debugf("Object %s doesn't exist, was it deleted after reconcile request?", request.NamespacedName)
			return &reconcile.Result{}
		}
		// Error reading the object - requeue the request.
		log.Errorf("Failed to query object %s: %s", request.NamespacedName, err)
		return &reconcile.Result{RequeueAfter: workflow.DefaultRetry}
	}
	return nil
}
