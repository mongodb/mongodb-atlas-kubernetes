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

package connectionsecret

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func HasReadyCondition(conditions []api.Condition) bool {
	for _, c := range conditions {
		if c.Type == api.ReadyType && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsDeploymentReady(d *akov2.AtlasDeployment) bool {
	return HasReadyCondition(d.Status.Conditions)
}

func IsDatabaseUserReady(u *akov2.AtlasDatabaseUser) bool {
	return HasReadyCondition(u.Status.Conditions)
}

func ResolveProjectIDFromDeployment(ctx context.Context, c client.Client, d *akov2.AtlasDeployment) (string, error) {
	if d.Spec.ExternalProjectRef != nil && d.Spec.ExternalProjectRef.ID != "" {
		return d.Spec.ExternalProjectRef.ID, nil
	}
	if d.Spec.ProjectRef != nil && d.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		if err := c.Get(ctx, *d.Spec.ProjectRef.GetObject(d.Namespace), project); err != nil {
			return "", fmt.Errorf("failed to resolve projectRef from deployment: %w", err)
		}
		return project.ID(), nil
	}
	return "", fmt.Errorf("missing both external and internal project references")
}

func ResolveProjectIDFromDatabaseUser(ctx context.Context, c client.Client, u *akov2.AtlasDatabaseUser) (string, error) {
	if u.Spec.ExternalProjectRef != nil && u.Spec.ExternalProjectRef.ID != "" {
		return u.Spec.ExternalProjectRef.ID, nil
	}
	if u.Spec.ProjectRef != nil && u.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		if err := c.Get(ctx, *u.Spec.ProjectRef.GetObject(u.Namespace), project); err != nil {
			return "", fmt.Errorf("failed to resolve projectRef from user: %w", err)
		}
		return project.ID(), nil
	}
	return "", fmt.Errorf("missing both external and internal project references")
}
