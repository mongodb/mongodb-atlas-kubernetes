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

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

// nolint:dupl
const DatabaseUserByGroupId = "databaseuser.groupId"

type DatabaseUserByGroupIdIndexer struct {
	logger *zap.SugaredLogger
	client client.Client
	ctx    context.Context
}

func NewDatabaseUserBySecretIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *DatabaseUserByGroupIdIndexer {
	return &DatabaseUserByGroupIdIndexer{logger: logger.Named(DatabaseUserByGroupId).Sugar(), client: client, ctx: ctx}
}

func (*DatabaseUserByGroupIdIndexer) Object() client.Object {
	return &generatedv1.DatabaseUser{}
}

func (*DatabaseUserByGroupIdIndexer) Name() string {
	return DatabaseUserByGroupId
}

// Keys extracts the index key(s) from the given object
func (i *DatabaseUserByGroupIdIndexer) Keys(object client.Object) []string {
	user, ok := object.(*generatedv1.DatabaseUser)
	if !ok {
		i.logger.Errorf("expected *v1.DatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.V20250312 == nil {
		return nil
	}

	if user.Spec.V20250312.GroupId != nil {
		return []string{*user.Spec.V20250312.GroupId}
	}

	if user.Spec.V20250312.GroupRef == nil {
		return nil
	}

	groupRef := *user.Spec.V20250312.GroupRef
	group := &generatedv1.Group{}

	err := i.client.Get(i.ctx, client.ObjectKey{
		Namespace: user.GetNamespace(),
		Name:      groupRef.Name,
	}, group)

	if err != nil {
		i.logger.Errorf("error getting group: %v", err)
		return nil
	}

	if group.Status.V20250312 == nil || group.Status.V20250312.Id == nil {
		return nil
	}

	return []string{*group.Status.V20250312.Id}
}
