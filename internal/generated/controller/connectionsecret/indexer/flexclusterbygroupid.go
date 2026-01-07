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

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// nolint:dupl
const FlexClusterByGroupIdIndex = "flexcluster.groupId"

type FlexClusterByGroupIdIndexer struct {
	logger *zap.SugaredLogger
}

func NewFlexClusterByGroupIdIndexer(logger *zap.Logger) *FlexClusterByGroupIdIndexer {
	return &FlexClusterByGroupIdIndexer{logger: logger.Named(FlexClusterByGroupIdIndex).Sugar()}
}

func (*FlexClusterByGroupIdIndexer) Object() client.Object {
	return &v1.FlexCluster{}
}

func (*FlexClusterByGroupIdIndexer) Name() string {
	return FlexClusterByGroupIdIndex
}

// Keys extracts the index key(s) from the given object
func (i *FlexClusterByGroupIdIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.FlexCluster)
	if !ok {
		i.logger.Errorf("expected *v1.FlexCluster but got %T", object)
		return nil
	}

	if resource.Status.V20250312 == nil || resource.Status.V20250312.GroupId == nil {
		return nil
	}

	return []string{*resource.Status.V20250312.GroupId}
}
