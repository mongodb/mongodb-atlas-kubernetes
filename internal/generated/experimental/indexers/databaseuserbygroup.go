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
	"context"

	zap "go.uber.org/zap"
	fields "k8s.io/apimachinery/pkg/fields"
	types "k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	handler "sigs.k8s.io/controller-runtime/pkg/handler"
	log "sigs.k8s.io/controller-runtime/pkg/log"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// nolint:dupl
const DatabaseUserByGroupIndex = "databaseuser.groupRef"

type DatabaseUserByGroupIndexer struct {
	logger *zap.SugaredLogger
}

func NewDatabaseUserByGroupIndexer(logger *zap.Logger) *DatabaseUserByGroupIndexer {
	return &DatabaseUserByGroupIndexer{logger: logger.Named(DatabaseUserByGroupIndex).Sugar()}
}
func (*DatabaseUserByGroupIndexer) Object() client.Object {
	return &v1.DatabaseUser{}
}
func (*DatabaseUserByGroupIndexer) Name() string {
	return DatabaseUserByGroupIndex
}

// Keys extracts the index key(s) from the given object
func (i *DatabaseUserByGroupIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.DatabaseUser)
	if !ok {
		i.logger.Errorf("expected *v1.DatabaseUser but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.V20250312 != nil && resource.Spec.V20250312.GroupRef != nil && resource.Spec.V20250312.GroupRef.Name != "" {
		keys = append(keys, types.NamespacedName{
			Name:      resource.Spec.V20250312.GroupRef.Name,
			Namespace: resource.Namespace,
		}.String())
	}
	return keys
}

func NewDatabaseUserByGroupMapFunc(kubeClient client.Client) handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		logger := log.FromContext(ctx)

		listOpts := &client.ListOptions{FieldSelector: fields.OneTermEqualSelector(DatabaseUserByGroupIndex, types.NamespacedName{
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		}.String())}

		list := &v1.DatabaseUserList{}
		err := kubeClient.List(ctx, list, listOpts)
		if err != nil {
			logger.Error(err, "failed to list DatabaseUser objects")
			return nil
		}

		requests := make([]reconcile.Request, 0, len(list.Items))
		for _, item := range list.Items {
			requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      item.Name,
				Namespace: item.Namespace,
			}})
		}

		return requests
	}
}
