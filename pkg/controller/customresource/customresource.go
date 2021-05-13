package customresource

import (
	"context"

	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

const (
	ResourcePolicyAnnotation       = "mongodb.com/atlas-resource-policy"
	ReconciliationPolicyAnnotation = "mongodb.com/atlas-reconciliation-policy"

	ResourcePolicyKeep       = "keep"
	ReconciliationPolicySkip = "skip"
)

// PrepareResource queries the Custom Resource 'request.NamespacedName' and populates the 'resource' pointer.
func PrepareResource(client client.Client, request reconcile.Request, resource mdbv1.AtlasCustomResource, log *zap.SugaredLogger) workflow.Result {
	return GetResource(client, request.Namespace, request.Name, resource, log)
}

// GetResource queries the Custom Resource key and populates the 'resource' pointer.
func GetResource(client client.Client, namespace, name string, resource client.Object, log *zap.SugaredLogger) workflow.Result {
	key := types.NamespacedName{Namespace: namespace, Name: name}
	err := client.Get(context.Background(), key, resource)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Debugf("Object %s doesn't exist, was it deleted after reconcile request?", key)
			return workflow.TerminateSilently().WithoutRetry()
		}
		// Error reading the object - requeue the request. Note, that we don't intend to update resource status
		// as most of all it will fail as well.
		log.Errorf("Failed to query object %s: %s", key, err)
		return workflow.TerminateSilently()
	}

	return workflow.OK()
}

// MarkReconciliationStarted updates the status of the Atlas Resource to indicate that the Operator has started working on it.
// Internally this will also update the 'observedGeneration' field that notify clients that the resource is being worked on
func MarkReconciliationStarted(client client.Client, resource mdbv1.AtlasCustomResource, log *zap.SugaredLogger) *workflow.Context {
	updatedConditions := status.EnsureConditionExists(status.FalseCondition(status.ReadyType), resource.GetStatus().GetConditions())

	ctx := workflow.NewContext(log, updatedConditions)
	statushandler.Update(ctx, client, nil, resource)

	return ctx
}

// ResourceShouldBeLeftInAtlas returns 'true' if the resource should not be removed from Atlas on K8s resource removal.
func ResourceShouldBeLeftInAtlas(resource mdbv1.AtlasCustomResource) bool {
	if v, ok := resource.GetAnnotations()[ResourcePolicyAnnotation]; ok {
		return v == ResourcePolicyKeep
	}
	return false
}

// ReconciliationShouldBeSkipped returns 'true' if reconciliation should be skipped for this resource.
func ReconciliationShouldBeSkipped(resource mdbv1.AtlasCustomResource) bool {
	if v, ok := resource.GetAnnotations()[ReconciliationPolicyAnnotation]; ok {
		return v == ReconciliationPolicySkip
	}
	return false
}
