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

package references

import (
	"context"
	"errors"
	"fmt"

	"github.com/crd2go/crd2go/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// GetGroupID returns the group IP (if present) from a given Kubernetes group reference
// TODO: Autogenerate with Scaffolder
func GetGroupID(ctx context.Context, c client.Client, groupRef *k8s.LocalReference, namespace string) (string, error) {
	if groupRef == nil {
		return "", errors.New("group reference is nil")
	}

	group := &akov2generated.Group{}
	err := c.Get(ctx, client.ObjectKey{
		Namespace: groupRef.Name,
		Name:      namespace,
	}, group)

	if err != nil {
		return "", fmt.Errorf("failed to get object: %w", err)
	}

	// for each suported Group version...
	if group.Status.V20250312 != nil && group.Status.V20250312.Id != nil {
		return *group.Status.V20250312.Id, nil
	}

	return "", errors.New("group ID is not available")
}
