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

package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// ClusterTarget is the factory type that implements ConnectionTarget
type ClusterTarget struct {
	Client client.Client
}

// NewClusterTarget creates a new ClusterTarget
func NewClusterTarget(c client.Client) *ClusterTarget {
	return &ClusterTarget{Client: c}
}

// ListForProject lists all FlexCluster connection targets for a given project ID
func (e *ClusterTarget) ListForProject(ctx context.Context, projectID string) ([]target.ConnectionTargetInstance, error) {
	list := &generatedv1.ClusterList{}

	if err := e.Client.List(ctx, list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.ClusterByGroupIdIndex, projectID),
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
func (e *ClusterTarget) GetConnectionTargetInstance(obj client.Object) target.ConnectionTargetInstance {
	c, ok := obj.(*generatedv1.Cluster)
	if !ok {
		return nil
	}

	if c.Spec.V20250312 != nil {
		return &ClusterInstance_v20250312{c}
	}

	// add additional if statements and return dedicated versioned instances when new versions are supported

	return nil
}
