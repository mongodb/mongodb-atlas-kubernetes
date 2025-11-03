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
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// nolint:dupl
const DataFederationByGroupIndex = "datafederation.groupRef"

type DataFederationByGroupIndexer struct {
	logger *zap.SugaredLogger
}

func NewDataFederationByGroupIndexer(logger *zap.Logger) *DataFederationByGroupIndexer {
	return &DataFederationByGroupIndexer{logger: logger.Named(DataFederationByGroupIndex).Sugar()}
}
func (*DataFederationByGroupIndexer) Object() client.Object {
	return &v1.DataFederation{}
}
func (*DataFederationByGroupIndexer) Name() string {
	return DataFederationByGroupIndex
}

// Keys extracts the index key(s) from the given object
func (i *DataFederationByGroupIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.DataFederation)
	if !ok {
		i.logger.Errorf("expected *v1.DataFederation but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.V20250312.GroupRef != nil && resource.Spec.V20250312.GroupRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.GroupRef.Name)
	}
	return keys
}
