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

package reconciler

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func (r *AtlasReconciler) ResolveSDKClientSet(ctx context.Context, referrer project.ProjectReferrerObject) (*atlas.ClientSet, error) {
	connectionConfig, err := r.ResolveConnectionConfig(ctx, referrer)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve connection config: %w", err)
	}
	sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, r.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate client set: %w", err)
	}
	return sdkClientSet, nil
}
