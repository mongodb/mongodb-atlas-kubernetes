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

package connsecretsgeneric

import (
	"context"
	"fmt"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func (r *ConnSecretReconciler) GetUserProjectName(ctx context.Context, user *akov2.AtlasDatabaseUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if user.Spec.ProjectRef != nil && user.Spec.ProjectRef.Name != "" {
		proj := &akov2.AtlasProject{}
		key := user.Spec.ProjectRef.GetObject(user.GetNamespace())
		if err := r.Client.Get(ctx, *key, proj); err != nil {
			return "", err
		}
		if proj.Spec.Name != "" {
			return kube.NormalizeIdentifier(proj.Spec.Name), nil
		}
	}

	if r != nil {
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
		return kube.NormalizeIdentifier(ap.Name), nil
	}

	return "", fmt.Errorf("project name not available")
}
