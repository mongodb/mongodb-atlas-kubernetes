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

package experimentalconnectionsecret

import (
	"context"
	"fmt"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

// GetUserProjectID returns the projectID of the user
func (r *ConnSecretReconciler) getUserProjectID(ctx context.Context, user *akov2.AtlasDatabaseUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil user")
	}
	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return user.Spec.ExternalProjectRef.ID, nil
	}
	return resolveProjectIDByKey(ctx, r.Client, user.AtlasProjectObjectKey())
}
