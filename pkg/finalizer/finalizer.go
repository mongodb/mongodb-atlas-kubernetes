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

	data, err := json.Marshal([]map[string]any{{
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

	data, err := json.Marshal([]map[string]any{{
		"op":    "replace",
		"path":  "/metadata/finalizers",
		"value": o.GetFinalizers(),
	}})
	if err != nil {
		return err
	}

	return c.Patch(ctx, o, client.RawPatch(types.JSONPatchType, data))
}
