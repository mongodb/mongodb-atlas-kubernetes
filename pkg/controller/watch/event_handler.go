package watch

import (
	"go.uber.org/zap"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

type ResourceEventHandler struct {
	Controller interface{}
}

type Deleter interface {
	Delete(e event.DeleteEvent) error
}

type Creator interface {
	Create(e event.CreateEvent) error
}

type Updater interface {
	Update(e event.UpdateEvent) error
}

type Genericer interface {
	Generic(e event.GenericEvent) error
}

var _ handler.EventHandler = &ResourceEventHandler{}

func (d *ResourceEventHandler) Delete(e event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	ctrl, ok := d.Controller.(Deleter)
	if !ok {
		return
	}

	objectKey := kube.ObjectKeyFromObject(e.Meta)
	log := zap.S().With("resource", objectKey)

	if err := ctrl.Delete(e); err != nil {
		log.Errorf("Object (%s) removed from Kubernetes, but controller could not delete it: %s", e.Object.GetObjectKind(), err)
		return
	}
}

func (d *ResourceEventHandler) Create(e event.CreateEvent, _ workqueue.RateLimitingInterface) {
	ctrl, ok := d.Controller.(Creator)
	if !ok {
		return
	}

	objectKey := kube.ObjectKeyFromObject(e.Meta)
	log := zap.S().With("resource", objectKey)

	if err := ctrl.Create(e); err != nil {
		log.Errorf("Object (%s) created in Kubernetes, but controller could not create it: %s", e.Object.GetObjectKind(), err)
		return
	}
}

func (d *ResourceEventHandler) Update(e event.UpdateEvent, _ workqueue.RateLimitingInterface) {
	ctrl, ok := d.Controller.(Updater)
	if !ok {
		return
	}

	log := zap.S().With(
		"resourceOld", kube.ObjectKeyFromObject(e.MetaOld),
		"resourceNew", kube.ObjectKeyFromObject(e.MetaNew),
	)

	if err := ctrl.Update(e); err != nil {
		log.Errorf("Object (%s -> %s) updated in Kubernetes, but controller could not update it: %s", e.ObjectOld.GetObjectKind(), e.ObjectNew.GetObjectKind(), err)
		return
	}
}

func (d *ResourceEventHandler) Generic(e event.GenericEvent, _ workqueue.RateLimitingInterface) {
	ctrl, ok := d.Controller.(Genericer)
	if !ok {
		return
	}

	objectKey := kube.ObjectKeyFromObject(e.Meta)
	log := zap.S().With("resource", objectKey)

	if err := ctrl.Generic(e); err != nil {
		log.Errorf("Object (%s) received a generic event from Kubernetes, but controller could not handle it: %s", e.Object.GetObjectKind(), err)
		return
	}
}
