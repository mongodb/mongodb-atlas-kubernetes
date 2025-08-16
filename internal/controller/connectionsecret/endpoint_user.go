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

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

// resolveProjectIDByKey returns the project id from the key
func (r *ConnSecretReconciler) resolveProjectIDByKey(ctx context.Context, key client.ObjectKey) (string, error) {
	proj := &akov2.AtlasProject{}
	if err := r.Client.Get(ctx, key, proj); err != nil {
		return "", err
	}
	if proj.ID() == "" {
		return "", ErrUnresolvedProjectID
	}
	return proj.ID(), nil
}

// resolveProjectNameByKey returns the project name from the key
func (r *ConnSecretReconciler) resolveProjectNameByKey(ctx context.Context, key client.ObjectKey) (string, error) {
	proj := &akov2.AtlasProject{}
	if err := r.Client.Get(ctx, key, proj); err != nil {
		return "", err
	}
	if proj.Spec.Name == "" {
		return "", ErrUnresolvedProjectName
	}
	return kube.NormalizeIdentifier(proj.Spec.Name), nil
}

// GetUserProjectID returns the projectID of the user
func (r *ConnSecretReconciler) getUserProjectID(ctx context.Context, user *akov2.AtlasDatabaseUser) (string, error) {
	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return user.Spec.ExternalProjectRef.ID, nil
	}
	if user.Spec.ProjectRef != nil && user.Spec.ProjectRef.Name != "" {
		return r.resolveProjectIDByKey(ctx, user.AtlasProjectObjectKey())
	}

	return "", fmt.Errorf("missing both external and internal project references")
}

// GetUserProjectName retrives the project name from the AtlasDatabaseUser (either by getting K8s AtlasProject or SDK calls)
func (r *ConnSecretReconciler) getUserProjectName(ctx context.Context, user *akov2.AtlasDatabaseUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil user")
	}
	if user.Spec.ProjectRef != nil && user.Spec.ProjectRef.Name != "" {
		return r.resolveProjectNameByKey(ctx, user.AtlasProjectObjectKey())
	}

	cfg, err := r.ResolveConnectionConfig(ctx, user)
	if err != nil {
		return "", err
	}
	sdk, err := r.AtlasProvider.SdkClientSet(ctx, cfg.Credentials, r.Log)
	if err != nil {
		return "", err
	}
	ap, err := r.ResolveProject(ctx, sdk.SdkClient20250312002, user)
	if err != nil {
		return "", err
	}
	if ap.Name == "" {
		return "", fmt.Errorf("project name not available")
	}

	return kube.NormalizeIdentifier(ap.Name), nil
}
