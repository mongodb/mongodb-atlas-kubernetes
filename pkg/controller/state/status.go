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

package state

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusObject interface {
	GetConditions() []metav1.Condition
}

type resource struct {
	Status statusResource `json:"status,omitempty"`
}

type statusResource struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func newStatusObject(conditions []metav1.Condition) *resource {
	return &resource{Status: statusResource{Conditions: conditions}}
}

func patchStatus(ctx context.Context, c client.Client, obj client.Object, status any) error {
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}
	patchErr := c.Status().Patch(ctx, obj, client.RawPatch(types.MergePatchType, statusJSON))
	if patchErr != nil {
		return fmt.Errorf("failed to patch status: %w", patchErr)
	}
	return nil
}
