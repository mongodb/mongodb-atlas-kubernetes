package customresource

import (
	"context"
	"fmt"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func AddFinalizer(ctx context.Context, client client.Client, resource mdbv1.AtlasCustomResource, finalizer string) error {
	err := client.Get(ctx, kube.ObjectKeyFromObject(resource), resource)
	if err != nil {
		return fmt.Errorf("cannot get resource to add finalizer: %w", err)
	}

	SetFinalizer(resource, finalizer)

	if err = client.Update(ctx, resource); err != nil {
		return fmt.Errorf("failed to add deletion finalizer from %s: %w", resource, err)
	}
	return nil
}

func DeleteFinalizer(ctx context.Context, client client.Client, resource mdbv1.AtlasCustomResource, finalizer string) error {
	err := client.Get(ctx, kube.ObjectKeyFromObject(resource), resource)
	if err != nil {
		return fmt.Errorf("cannot get resource to delete finalizer: %w", err)
	}

	UnsetFinalizer(resource, finalizer)

	if err = client.Update(ctx, resource); err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", resource, err)
	}
	return nil
}
