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
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// nolint:dupl
const ClusterByGroupIdIndex = "cluster.groupId"

type ClusterByGroupIdIndexer struct {
	logger *zap.SugaredLogger
}

func NewClusterByGroupIdIndexer(logger *zap.Logger) *ClusterByGroupIdIndexer {
	return &ClusterByGroupIdIndexer{logger: logger.Named(ClusterByGroupIdIndex).Sugar()}
}

func (*ClusterByGroupIdIndexer) Object() client.Object {
	return &generatedv1.Cluster{}
}

func (*ClusterByGroupIdIndexer) Name() string {
	return ClusterByGroupIdIndex
}

// Keys extracts the index key(s) from the given object
func (i *ClusterByGroupIdIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*generatedv1.Cluster)
	if !ok {
		i.logger.Errorf("expected *generatedv1.Cluster but got %T", object)
		return nil
	}

	if resource.Status.V20250312 == nil || resource.Status.V20250312.GroupId == nil {
		return nil
	}

	return []string{*resource.Status.V20250312.GroupId}
}
