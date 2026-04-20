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

package customresource

import (
	"context"
	"encoding/json"
	"reflect"
	"slices"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

const FinalizerLabel = "mongodbatlas/finalizer"

type FinalizerOperator func(resource api.AtlasCustomResource, finalizer string)

func HaveFinalizer(resource api.AtlasCustomResource, finalizer string) bool {
	return slices.Contains(resource.GetFinalizers(), finalizer)
}

// SetFinalizer Add the given finalizer to the list of resource finalizer
func SetFinalizer(resource api.AtlasCustomResource, finalizer string) {
	if !HaveFinalizer(resource, finalizer) {
		resource.SetFinalizers(append(resource.GetFinalizers(), finalizer))
	}
}

// UnsetFinalizer remove the given finalizer from the list of resource finalizer
func UnsetFinalizer(resource api.AtlasCustomResource, finalizer string) {
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
	c client.Client,
	resource api.AtlasCustomResource,
	op FinalizerOperator,
) error {
	// we just copied an api.AtlasCustomResource so it must be one
	resourceCopy := resource.DeepCopyObject().(api.AtlasCustomResource)
	op(resourceCopy, FinalizerLabel)

	if reflect.DeepEqual(resource.GetFinalizers(), resourceCopy.GetFinalizers()) {
		return nil
	}

	data, err := json.Marshal([]map[string]any{{
		"op":    "replace",
		"path":  "/metadata/finalizers",
		"value": resourceCopy.GetFinalizers(),
	}})
	if err != nil {
		return err
	}

	return c.Patch(ctx, resource, client.RawPatch(types.JSONPatchType, data))
}
