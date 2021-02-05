package watch

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewResourceWatcher() ResourceWatcher {
	return ResourceWatcher{
		WatchedResources: map[WatchedObject]map[client.ObjectKey]bool{},
	}
}

// ResourceWatcher is the object containing the map of watched_resource -> []dependant_resource.
type ResourceWatcher struct {
	WatchedResources map[WatchedObject]map[client.ObjectKey]bool
}

// EnsureResourcesAreWatched registers a dependant for the watched objects.
// This will let the controller to react on the events for the watched objects and trigger reconciliation for dependants.
func (r ResourceWatcher) EnsureResourcesAreWatched(dependant client.ObjectKey, resourceKind string, log *zap.SugaredLogger, watchedObjectsKeys ...client.ObjectKey) {
	for _, watchedObjectKey := range watchedObjectsKeys {
		r.addWatchedResourceIfNotAdded(watchedObjectKey, resourceKind, dependant, log)
	}

	// Next we need to clean any watched resources that are not referenced any more. This could happen if the SecretRef
	// has been updated to reference another Secret, for example
	r.cleanNonWatchedResources(dependant, resourceKind, watchedObjectsKeys)
}

func (r *ResourceWatcher) addWatchedResourceIfNotAdded(watchedObjectKey client.ObjectKey, resourceKind string, dependentResourceNsName client.ObjectKey, log *zap.SugaredLogger) {
	key := WatchedObject{ResourceKind: resourceKind, Resource: watchedObjectKey}
	if _, ok := r.WatchedResources[key]; !ok {
		r.WatchedResources[key] = make(map[client.ObjectKey]bool)
	}
	if _, ok := r.WatchedResources[key][dependentResourceNsName]; !ok {
		log.Debugf("Watching %s to trigger reconciliation for %s", key, dependentResourceNsName)
	}
	r.WatchedResources[key][dependentResourceNsName] = true
}

func (r ResourceWatcher) cleanNonWatchedResources(dependant client.ObjectKey, resourceKind string, watchedKeys []client.ObjectKey) {
	for k, v := range r.WatchedResources {
		if !contains(watchedKeys, k.Resource) || k.ResourceKind != resourceKind {
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
