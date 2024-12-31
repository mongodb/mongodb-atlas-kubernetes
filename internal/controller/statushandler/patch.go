package statushandler

import (
	"encoding/json"
	"reflect"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

// patchUpdateStatus performs the JSONPatch patch update to the Atlas Custom Resource.
// The "jsonPatch" merge allows to update only status field so is more
func patchUpdateStatus(ctx *workflow.Context, kubeClient client.Client, resource api.AtlasCustomResource) error {
	// we just copied an api.AtlasCustomResource so it must be one
	resourceCopy := resource.DeepCopyObject().(api.AtlasCustomResource)
	resourceCopy.UpdateStatus(ctx.Conditions(), ctx.StatusOptions()...)

	if reflect.DeepEqual(resource.GetStatus(), resourceCopy.GetStatus()) {
		return nil
	}

	data, err := json.Marshal([]map[string]interface{}{{
		"op":    "replace",
		"path":  "/status",
		"value": resourceCopy.GetStatus(),
	}})
	if err != nil {
		return err
	}

	return kubeClient.Status().Patch(ctx.Context, resource, client.RawPatch(types.JSONPatchType, data))
}
