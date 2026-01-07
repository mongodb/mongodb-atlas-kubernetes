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

package flexcluster

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// FlexClusterTarget is the factory type that implements ConnectionTarget
type FlexClusterTarget struct {
	Client client.Client
}

// NewFlexClusterTarget creates a new FlexClusterTarget
func NewFlexClusterTarget(c client.Client) *FlexClusterTarget {
	return &FlexClusterTarget{Client: c}
}

// ListForProject lists all FlexCluster connection targets for a given project ID
func (e *FlexClusterTarget) ListForProject(ctx context.Context, projectID string) ([]target.ConnectionTargetInstance, error) {
	list := &generatedv1.FlexClusterList{}

	if err := e.Client.List(ctx, list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.FlexClusterByGroupIdIndex, projectID),
	}); err != nil {
		return nil, err
	}

	out := make([]target.ConnectionTargetInstance, 0, len(list.Items))
	for i := range list.Items {
		if instance := e.GetConnectionTargetInstance(&list.Items[i]); instance != nil {
			out = append(out, instance)
		}
	}

	return out, nil
}

// WrapObject wraps a client.Object if it's a FlexCluster, returning the wrapped target or nil if not matching
func (e *FlexClusterTarget) GetConnectionTargetInstance(obj client.Object) target.ConnectionTargetInstance {
	fc, ok := obj.(*generatedv1.FlexCluster)
	if !ok {
		return nil
	}

	if fc.Spec.V20250312 != nil {
		return &FlexClusterInstance_v20250312{fc}
	}

	// add additional if statements and return dedicated versioned instances when new versions are supported

	return nil
}
