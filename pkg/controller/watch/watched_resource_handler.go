package watch

import (
	"context"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
// WatchedObject must trigger reconciliation for any Operator managed Resource (AtlasProject, AtlasDeployment etc). This is
// done via consulting the 'TrackedResources' map. The map is stored in relevant Reconciler which ensures it's up-to-date
// on each reconciliation
type ResourcesHandler struct {
	ResourceKind    string
	ResourceWatcher *DeprecatedResourceWatcher
}

// NewSecretHandler TODO Igor: refactor this to create generic constructor
func NewSecretHandler(r *DeprecatedResourceWatcher) *ResourcesHandler {
	return &ResourcesHandler{ResourceKind: "Secret", ResourceWatcher: r}
}

func NewBackupScheduleHandler(r *DeprecatedResourceWatcher) *ResourcesHandler {
	return &ResourcesHandler{ResourceKind: "AtlasBackupSchedule", ResourceWatcher: r}
}

func NewBackupPolicyHandler(r *DeprecatedResourceWatcher) *ResourcesHandler {
	return &ResourcesHandler{ResourceKind: "AtlasBackupPolicy", ResourceWatcher: r}
}

func NewAtlasTeamHandler(r *DeprecatedResourceWatcher) *ResourcesHandler {
	return &ResourcesHandler{ResourceKind: "AtlasTeam", ResourceWatcher: r}
}

// Create handles the Create event for the resource.
// Note that we implement Create in addition to Update to be able to handle cases when config map or secret is deleted
// and then created again.
func (c *ResourcesHandler) Create(ctx context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
	c.doHandle(ctx, e.Object.GetNamespace(), e.Object.GetName(), e.Object.GetObjectKind().GroupVersionKind().Kind, q)
}

func (c *ResourcesHandler) Update(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if !shouldHandleUpdate(e) {
		zap.S().Debugf("resource watcher: skipping update for resource %v", e.ObjectNew.GetName())
		return
	}
	// For some reasons e.ObjectOld.GetObjectKind().GroupVersionKind().Kind is empty... that's why we have to keep the
	// resourceKind as a field in the handler...
	c.doHandle(ctx, e.ObjectOld.GetNamespace(), e.ObjectNew.GetName(), c.ResourceKind, q)
}

// shouldHandleUpdate return true if the update event must be handled. This should happen only if the real data has
// changed, not status etc.
func shouldHandleUpdate(e event.UpdateEvent) bool {
	switch v := e.ObjectOld.(type) {
	case *corev1.ConfigMap:
		return !reflect.DeepEqual(v.Data, e.ObjectNew.(*corev1.ConfigMap).Data)
	case *corev1.Secret:
		return !reflect.DeepEqual(v.Data, e.ObjectNew.(*corev1.Secret).Data)
	case *akov2.AtlasTeam:
		return !reflect.DeepEqual(v.Spec, e.ObjectNew.(*akov2.AtlasTeam).Spec)
	case *akov2.AtlasBackupSchedule:
		return !reflect.DeepEqual(v.Spec, e.ObjectNew.(*akov2.AtlasBackupSchedule).Spec)
	case *akov2.AtlasBackupPolicy:
		return !reflect.DeepEqual(v.Spec, e.ObjectNew.(*akov2.AtlasBackupPolicy).Spec)
	}
	return true
}

func (c *ResourcesHandler) doHandle(_ context.Context, namespace, name, kind string, q workqueue.RateLimitingInterface) {
	watchedResource := WatchedObject{
		ResourceKind: kind,
		Resource:     types.NamespacedName{Name: name, Namespace: namespace},
	}
	trackedResources := c.ResourceWatcher.WatchedResourcesSnapshot()
	for k := range trackedResources[watchedResource] {
		zap.S().Infof("%s has been modified -> triggering reconciliation for the %s", watchedResource, k)
		q.Add(reconcile.Request{NamespacedName: k})
	}
}

// Delete (Seems we don't need to react on watched resources removal..)
func (c *ResourcesHandler) Delete(context.Context, event.DeleteEvent, workqueue.RateLimitingInterface) {
}

func (c *ResourcesHandler) Generic(context.Context, event.GenericEvent, workqueue.RateLimitingInterface) {
}
