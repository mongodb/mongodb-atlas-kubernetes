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

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

// nolint:dupl
const ParentBySecretIndex = "parent.secretRef"

type ParentBySecretIndexer struct {
	logger *zap.SugaredLogger
}

func NewParentBySecretIndexer(logger *zap.Logger) *ParentBySecretIndexer {
	return &ParentBySecretIndexer{logger: logger.Named(ParentBySecretIndex).Sugar()}
}
func (*ParentBySecretIndexer) Object() client.Object {
	return &v1.Parent{}
}
func (*ParentBySecretIndexer) Name() string {
	return ParentBySecretIndex
}

// Keys extracts the index key(s) from the given object
func (i *ParentBySecretIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.Parent)
	if !ok {
		i.logger.Errorf("expected *v1.Parent but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.Integrations != nil {
		for _, integration := range *resource.Spec.Integrations {
			if integration.Credentials != nil {
				for _, credential := range *integration.Credentials {
					if credential.SecretRef != nil && credential.SecretRef.Name != "" {
						keys = append(keys, types.NamespacedName{
							Name:      credential.SecretRef.Name,
							Namespace: resource.Namespace,
						}.String())
					}
				}
			}
		}
	}
	return keys
}

//nolint:dupl
func NewParentBySecretMapFunc(kubeClient client.Client) handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		logger := log.FromContext(ctx)

		listOpts := &client.ListOptions{FieldSelector: fields.OneTermEqualSelector(ParentBySecretIndex, types.NamespacedName{
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		}.String())}

		list := &v1.ParentList{}
		err := kubeClient.List(ctx, list, listOpts)
		if err != nil {
			logger.Error(err, "failed to list Parent objects")
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
