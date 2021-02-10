package watch

import (
	"go.uber.org/zap"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// EventHandlerWithDelete is an extension of EnqueueRequestForObject that will _not_ trigger a reconciliation for a Delete event.
// Instead, it will call an external controller's Delete() method and pass the event argument unchanged.
type EventHandlerWithDelete struct {
	handler.EnqueueRequestForObject
	Controller interface {
		Delete(e event.DeleteEvent) error
	}
}

func (d *EventHandlerWithDelete) Delete(e event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	objectKey := kube.ObjectKeyFromObject(e.Meta)
	log := zap.S().With("resource", objectKey)

	if err := d.Controller.Delete(e); err != nil {
		log.Errorf("Object (%s) removed from Kubernetes, but controller could not delete it: %s", e.Object.GetObjectKind(), err)
		return
	}
}
