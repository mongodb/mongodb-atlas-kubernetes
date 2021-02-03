package watch

import (
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewResourceWatcher() ResourceWatcher {
	return ResourceWatcher{
		WatchedResources: map[WatchedObject][]client.ObjectKey{},
	}
}

// ResourceWatcher is the object containing the map of watched_resource -> []dependant_resource.
type ResourceWatcher struct {
	WatchedResources map[WatchedObject][]client.ObjectKey
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
		r.WatchedResources[key] = make([]types.NamespacedName, 0)
	}
	found := false
	for _, v := range r.WatchedResources[key] {
		if v == dependentResourceNsName {
			found = true
		}
	}
	if !found {
		r.WatchedResources[key] = append(r.WatchedResources[key], dependentResourceNsName)
		log.Debugf("Watching %s to trigger reconciliation for %s", key, dependentResourceNsName)
	}
}

func (r ResourceWatcher) cleanNonWatchedResources(dependant client.ObjectKey, resourceKind string, watchedKeys []client.ObjectKey) {
	for k, v := range r.WatchedResources {
		if pos(watchedKeys, k.Resource) < 0 || k.ResourceKind != resourceKind {
			var dependantPos int
			if dependantPos = pos(v, dependant); dependantPos >= 0 {
				// we found the old dependency (not watched any more) so we need to remove it
				r.WatchedResources[k] = remove(r.WatchedResources[k], dependantPos)
			}
		}
	}
}

func pos(watchedKeys []client.ObjectKey, key client.ObjectKey) int {
	for i, k := range watchedKeys {
		if k == key {
			return i
		}
	}
	return -1
}

func remove(slice []client.ObjectKey, s int) []client.ObjectKey {
	return append(slice[:s], slice[s+1:]...)
}
