package watch

import (
	"sync"

	"golang.org/x/exp/maps"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewDeprecatedResourceWatcher creates a new resource watcher.
//
// NOTE: DeprecatedResourceWatcher is DEPRECATED and DISCOURAGED to be used in new implementations.
// Use controller-runtime intrinsics instead, see https://book.kubebuilder.io/reference/watching-resources/externally-managed.
func NewDeprecatedResourceWatcher() DeprecatedResourceWatcher {
	return DeprecatedResourceWatcher{
		mtx:              &sync.RWMutex{},
		watchedResources: map[WatchedObject]map[client.ObjectKey]bool{},
	}
}

// DeprecatedResourceWatcher is the object containing the map of watched_resource -> []dependant_resource.
type DeprecatedResourceWatcher struct {
	mtx              *sync.RWMutex
	watchedResources map[WatchedObject]map[client.ObjectKey]bool
}

// WatchedResourcesSnapshot returns the most recent snapshot of watched resources.
// Note that entries here can change concurrently.
func (r *DeprecatedResourceWatcher) WatchedResourcesSnapshot() map[WatchedObject]map[client.ObjectKey]bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if len(r.watchedResources) == 0 {
		return nil
	}

	result := make(map[WatchedObject]map[client.ObjectKey]bool, len(r.watchedResources))
	for watchedKey, dependentResources := range r.watchedResources {
		dependentResourcesCopy := make(map[client.ObjectKey]bool, len(dependentResources))
		maps.Copy(dependentResourcesCopy, dependentResources)
		result[watchedKey] = dependentResourcesCopy
	}
	return result
}

// EnsureResourcesAreWatched registers a dependant for the watched objects.
// This will let the controller to react on the events for the watched objects and trigger reconciliation for dependants.
func (r *DeprecatedResourceWatcher) EnsureResourcesAreWatched(dependant client.ObjectKey, resourceKind string, log *zap.SugaredLogger, watchedObjectsKeys ...client.ObjectKey) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, watchedObjectKey := range watchedObjectsKeys {
		r.unsafeAddWatchedResourceIfNotAdded(watchedObjectKey, resourceKind, dependant, log)
	}
	// Next we need to clean any watched resources that are not referenced any more. This could happen if the SecretRef
	// has been updated to reference another Secret, for example
	r.unsafeCleanNonWatchedResources(dependant, resourceKind, watchedObjectsKeys)
}

func (r *DeprecatedResourceWatcher) EnsureMultiplesResourcesAreWatched(dependant client.ObjectKey, log *zap.SugaredLogger, resources ...WatchedObject) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, res := range resources {
		r.unsafeAddWatchedResourceIfNotAdded(res.Resource, res.ResourceKind, dependant, log)
		log.Debugf("resource watcher: watching %v to trigger reconciliation for %v", res.Resource, dependant)
	}

	r.unsafeCleanNonWatchedResourcesExceptMultiple(dependant, resources...)
}

func (r *DeprecatedResourceWatcher) unsafeAddWatchedResourceIfNotAdded(watchedObjectKey client.ObjectKey, resourceKind string, dependentResourceNsName client.ObjectKey, log *zap.SugaredLogger) {
	key := WatchedObject{ResourceKind: resourceKind, Resource: watchedObjectKey}
	if _, ok := r.watchedResources[key]; !ok {
		r.watchedResources[key] = make(map[client.ObjectKey]bool)
	}
	if _, ok := r.watchedResources[key][dependentResourceNsName]; !ok {
		log.Debugf("resource watcher: watching %s to trigger reconciliation for %s", key, dependentResourceNsName)
	}
	r.watchedResources[key][dependentResourceNsName] = true
}

func (r *DeprecatedResourceWatcher) unsafeCleanNonWatchedResources(dependant client.ObjectKey, resourceKind string, watchedKeys []client.ObjectKey) {
	for k, v := range r.watchedResources {
		if !contains(watchedKeys, k.Resource) || k.ResourceKind != resourceKind {
			delete(v, dependant)
		}
	}
}

func (r *DeprecatedResourceWatcher) unsafeCleanNonWatchedResourcesExceptMultiple(dependant client.ObjectKey, resources ...WatchedObject) {
	for k, v := range r.watchedResources {
		toRemove := true
		for _, res := range resources {
			if res.Resource == k.Resource {
				toRemove = false
			}
		}
		if toRemove {
			delete(v, dependant)
		}
	}
}

func contains(watchedKeys []client.ObjectKey, key client.ObjectKey) bool {
	for _, k := range watchedKeys {
		if k == key {
			return true
		}
	}
	return false
}
