package watch

import (
	"k8s.io/apimachinery/pkg/types"
)

func NewResourceWatcher() ResourceWatcher {
	return ResourceWatcher{
		WatchedResources: map[WatchedObject][]types.NamespacedName{},
	}
}

type ResourceWatcher struct {
	WatchedResources map[WatchedObject][]types.NamespacedName
}

// Methods for updating the watched resources are TODO
