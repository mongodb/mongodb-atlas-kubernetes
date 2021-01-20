package workflow

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type DeleteEventHandler struct {
	*handler.EnqueueRequestForObject
	log    *zap.SugaredLogger
	Parent interface {
		Delete(runtime.Object) error
	}
}

func (d *DeleteEventHandler) Delete(e event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	objectKey := kube.ObjectKey(e.Meta.GetNamespace(), e.Meta.GetName())
	log := d.log.With("resource", objectKey)

	log.Infow("Cleaning up Atlas resource", "resource", e.Object)
	if err := d.Parent.Delete(e.Object); err != nil {
		log.Errorf("MongoDB resource removed from Kubernetes, but failed to clean some state in Atlas: %s", err)
		return
	}

	log.Info("Removed MongoDB resource from Kubernetes and Atlas")
}
