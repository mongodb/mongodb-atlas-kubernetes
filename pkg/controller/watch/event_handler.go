package watch

import (
	"context"

	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/kube"
)

// EventHandlerWithDelete is an extension of EnqueueRequestForObject that will _not_ trigger a reconciliation for a Delete event.
// Instead, it will call an external controller's Delete() method and pass the event argument unchanged.
type EventHandlerWithDelete struct {
	handler.EnqueueRequestForObject
	Controller interface {
		Delete(ctx context.Context, e event.DeleteEvent) error
	}
}

func (d *EventHandlerWithDelete) Delete(ctx context.Context, e event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	objectKey := kube.ObjectKeyFromObject(e.Object)
	log := zap.S().With("resource", objectKey)

	if err := d.Controller.Delete(ctx, e); err != nil && k8serrors.IsNotFound(err) {
		log.Errorf("Object (%s) removed from Kubernetes, but controller could not delete it: %s", e.Object.GetObjectKind(), err)
	}
}
