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

package target

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/data"
)

// ConnectionTargetInstance defines the interface for methods on a wrapped connection target instance.
type ConnectionTargetInstance interface {
	GetConnectionTargetType() string
	GetName() string
	IsReady() bool
	GetScopeType() string
	GetProjectID(ctx context.Context) string
	BuildConnectionData(ctx context.Context) *data.ConnectionSecret
}

// ConnectionTarget defines the interface for connection target factories that can
// list and wrap connection target instances.
type ConnectionTarget interface {
	// ListForProject lists all connection target instances of this kind for a given project ID
	ListForProject(ctx context.Context, projectID string) ([]ConnectionTargetInstance, error)
	// GetConnectionTargetInstance returns a wrapped connection target instance if the given object matches
	GetConnectionTargetInstance(obj client.Object) ConnectionTargetInstance
}
