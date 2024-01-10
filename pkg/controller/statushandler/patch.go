package statushandler

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type patchValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// patchUpdateStatus performs the JSONPatch patch update to the Atlas Custom Resource.
// The "jsonPatch" merge allows to update only status field so is more
func patchUpdateStatus(ctx context.Context, kubeClient client.Client, resource mdbv1.AtlasCustomResource) error {
	return doPatch(ctx, kubeClient, resource, resource.GetStatus())
}

func doPatch(ctx context.Context, kubeClient client.Client, resource client.Object, statusValue interface{}) error {
	payload := []patchValue{{
		Op:    "replace",
		Path:  "/status",
		Value: statusValue,
	}}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	patch := client.RawPatch(types.JSONPatchType, data)
	return kubeClient.Status().Patch(ctx, resource, patch)
}
