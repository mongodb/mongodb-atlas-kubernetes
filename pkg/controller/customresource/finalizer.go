package customresource

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

const FinalizerLabel = "mongodbatlas/finalizer"

type FinalizerOperator func(resource mdbv1.AtlasCustomResource, finalizer string)

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

func ManageFinalizer(
	ctx context.Context,
	client client.Client,
	resource mdbv1.AtlasCustomResource,
	op FinalizerOperator,
) error {
	err := client.Get(ctx, kube.ObjectKeyFromObject(resource), resource)
	if err != nil {
		return fmt.Errorf("failed to get %t before removing deletion finalizer: %w", resource, err)
	}

	op(resource, FinalizerLabel)

	err = client.Update(ctx, resource)
	if err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", resource.GetName(), err)
	}

	return nil
}
