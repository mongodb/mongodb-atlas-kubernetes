package status

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	mustUnmarshal(mustMarshal(obj), result)
	return result
}

func PatchStatus(ctx context.Context, c client.Client, obj client.Object, status any) error {
	patchErr := c.Status().Patch(ctx, obj, client.RawPatch(types.MergePatchType, mustMarshal(status)))
	if patchErr != nil {
		return fmt.Errorf("failed to patch status: %w", patchErr)
	}
	return nil
}

func mustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
