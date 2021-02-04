package watch

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// WatchedObject is the  object watched by controller. Includes its type and namespace+name
type WatchedObject struct {
	ResourceKind string
	Resource     types.NamespacedName
}

func (w WatchedObject) String() string {
	return fmt.Sprintf("%s (%s)", w.Resource, w.ResourceKind)
}

// ResourcesHandler is a special implementation of 'handler.EventHandler' that checks if the event for
// WatchedObject must trigger reconciliation for any Operator managed Resource (AtlasProject, AtlasCluster etc). This is
// done via consulting the 'TrackedResources' map. The map is stored in relevant Reconciler which ensures it's up-to-date
// on each reconciliation
type ResourcesHandler struct {
	TrackedResources map[WatchedObject][]types.NamespacedName
}

// Create handles the Create event for the resource.
// Note that we implement Create in addition to Update to be able to handle cases when config map or secret is deleted
// and then created again.
func (c *ResourcesHandler) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	c.doHandle(e.Meta.GetNamespace(), e.Meta.GetName(), e.Object.GetObjectKind().GroupVersionKind().Kind, q)
}

func (c *ResourcesHandler) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if !shouldHandleUpdate(e) {
		return
	}
	c.doHandle(e.MetaOld.GetNamespace(), e.MetaOld.GetName(), e.ObjectOld.GetObjectKind().GroupVersionKind().Kind, q)
}

// shouldHandleUpdate return true if the update event must be handled. This should happen only if the real data has
// changed, not status etc.
func shouldHandleUpdate(e event.UpdateEvent) bool {
	switch v := e.ObjectOld.(type) {
	case *corev1.ConfigMap:
		return !reflect.DeepEqual(v.Data, e.ObjectNew.(*corev1.ConfigMap).Data)
	case *corev1.Secret:
		return !reflect.DeepEqual(v.Data, e.ObjectNew.(*corev1.Secret).Data)
	}
	return true
}

func (c *ResourcesHandler) doHandle(namespace, name, kind string, q workqueue.RateLimitingInterface) {
	watchedResource := WatchedObject{
		ResourceKind: kind,
		Resource:     types.NamespacedName{Name: name, Namespace: namespace},
	}
	for _, v := range c.TrackedResources[watchedResource] {
		zap.S().Infof("%s has been modified -> triggering reconciliation for the %s", watchedResource, v)
		q.Add(reconcile.Request{NamespacedName: v})
	}
}

// Delete (Seems we don't need to react on watched resources removal..)
func (c *ResourcesHandler) Delete(event.DeleteEvent, workqueue.RateLimitingInterface) {}

func (c *ResourcesHandler) Generic(event.GenericEvent, workqueue.RateLimitingInterface) {}
