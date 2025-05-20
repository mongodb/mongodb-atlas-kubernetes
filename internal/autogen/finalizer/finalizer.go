package finalizer

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func UnsetFinalizers(ctx context.Context, c client.Client, o client.Object, finalizer ...string) error {
	for _, f := range finalizer {
		controllerutil.RemoveFinalizer(o, f)
	}

	data, err := json.Marshal([]map[string]interface{}{{
		"op":    "replace",
		"path":  "/metadata/finalizers",
		"value": o.GetFinalizers(),
	}})
	if err != nil {
		return err
	}

	return c.Patch(ctx, o, client.RawPatch(types.JSONPatchType, data))
}

func EnsureFinalizers(ctx context.Context, c client.Client, o client.Object, finalizer ...string) error {
	hasAllFinalizers := true
	for _, f := range finalizer {
		if !controllerutil.ContainsFinalizer(o, f) {
			hasAllFinalizers = false
		}
	}
	if hasAllFinalizers {
		return nil
	}

	for _, f := range finalizer {
		controllerutil.AddFinalizer(o, f)
	}

	data, err := json.Marshal([]map[string]interface{}{{
		"op":    "replace",
		"path":  "/metadata/finalizers",
		"value": o.GetFinalizers(),
	}})
	if err != nil {
		return err
	}

	return c.Patch(ctx, o, client.RawPatch(types.JSONPatchType, data))
}
