// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package indexer

import (
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	zap "go.uber.org/zap"
	types "k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// nolint:dupl
const FlexClusterByGroupIndex = "flexcluster.groupRef"

type FlexClusterByGroupIndexer struct {
	logger *zap.SugaredLogger
}

func NewFlexClusterByGroupIndexer(logger *zap.Logger) *FlexClusterByGroupIndexer {
	return &FlexClusterByGroupIndexer{logger: logger.Named(FlexClusterByGroupIndex).Sugar()}
}
func (*FlexClusterByGroupIndexer) Object() client.Object {
	return &v1.FlexCluster{}
}
func (*FlexClusterByGroupIndexer) Name() string {
	return FlexClusterByGroupIndex
}

// Keys extracts the index key(s) from the given object
func (i *FlexClusterByGroupIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.FlexCluster)
	if !ok {
		i.logger.Errorf("expected *v1.FlexCluster but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.V20250312.GroupRef != nil && resource.Spec.V20250312.GroupRef.Name != "" {
		keys = append(keys, types.NamespacedName{
			Name:      resource.Spec.V20250312.GroupRef.Name,
			Namespace: resource.Namespace,
		}.String())
	}
	return keys
}
func FlexClusterRequestsFromGroup(list *v1.FlexClusterList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      item.Name,
			Namespace: item.Namespace,
		}})
	}
	return requests
}
