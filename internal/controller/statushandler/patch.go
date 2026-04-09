// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	data, err := json.Marshal([]map[string]any{{
		"op":    "replace",
		"path":  "/status",
		"value": resourceCopy.GetStatus(),
	}})
	if err != nil {
		return err
	}

	return kubeClient.Status().Patch(ctx.Context, resource, client.RawPatch(types.JSONPatchType, data))
}
