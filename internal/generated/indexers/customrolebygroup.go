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
	types "k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// nolint:dupl
const CustomRoleByGroupIndex = "customrole.groupRef"

type CustomRoleByGroupIndexer struct {
	logger *zap.SugaredLogger
}

func NewCustomRoleByGroupIndexer(logger *zap.Logger) *CustomRoleByGroupIndexer {
	return &CustomRoleByGroupIndexer{logger: logger.Named(CustomRoleByGroupIndex).Sugar()}
}
func (*CustomRoleByGroupIndexer) Object() client.Object {
	return &v1.CustomRole{}
}
func (*CustomRoleByGroupIndexer) Name() string {
	return CustomRoleByGroupIndex
}

// Keys extracts the index key(s) from the given object
func (i *CustomRoleByGroupIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.CustomRole)
	if !ok {
		i.logger.Errorf("expected *v1.CustomRole but got %T", object)
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
func CustomRoleRequestsFromGroup(list *v1.CustomRoleList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      item.Name,
			Namespace: item.Namespace,
		}})
	}
	return requests
}
