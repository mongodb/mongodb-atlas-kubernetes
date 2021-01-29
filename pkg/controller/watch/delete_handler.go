package watch

import (
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

type DeleteEventHandler struct {
	*handler.EnqueueRequestForObject
	Controller interface {
		Delete(runtime.Object) error
	}
}

func (d *DeleteEventHandler) Delete(e event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	objectKey := kube.ObjectKey(e.Meta.GetNamespace(), e.Meta.GetName())
	log := zap.S().With("resource", objectKey)

	if err := d.Controller.Delete(e.Object); err != nil {
		log.Errorf("Resource %s removed from Kubernetes, but failed to clean some state in Atlas: %s", e.Object.GetObjectKind(), err)
		return
	}
}
