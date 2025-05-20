package status

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/json"
)

type Resource struct {
	client.Object `json:"-"`
	Status        Status `json:"status,omitempty"`
}

type Status struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func GetStatus(obj client.Object) *Resource {
	result := &Resource{
		Object: obj,
	}
	json.MustUnmarshal(json.MustMarshal(obj), result)
	return result
}

func PatchStatus(ctx context.Context, c client.Client, obj client.Object, status any) error {
	patchErr := c.Status().Patch(ctx, obj, client.RawPatch(types.MergePatchType, json.MustMarshal(status)))
	if patchErr != nil {
		return fmt.Errorf("failed to patch status: %w", patchErr)
	}
	return nil
}
