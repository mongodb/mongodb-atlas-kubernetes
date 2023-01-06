package customresource

import (
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func HaveFinalizer(resource mdbv1.AtlasCustomResource, finalizer string) bool {
	for _, f := range resource.GetFinalizers() {
		if f == finalizer {
			return true
		}
	}

	return false
}

// SetFinalizer Add the given finalizer to the list of resource finalizer
func SetFinalizer(resource mdbv1.AtlasCustomResource, finalizer string) {
	if !HaveFinalizer(resource, finalizer) {
		resource.SetFinalizers(append(resource.GetFinalizers(), finalizer))
	}
}

// UnsetFinalizer remove the given finalizer from the list of resource finalizer
func UnsetFinalizer(resource mdbv1.AtlasCustomResource, finalizer string) {
	finalizers := make([]string, 0, len(resource.GetFinalizers()))

	for _, f := range resource.GetFinalizers() {
		if f != finalizer {
			finalizers = append(finalizers, f)
		}
	}

	resource.SetFinalizers(finalizers)
}
