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

//nolint:dupl
package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasDatabaseUserByProject = "atlasdatabaseuser.spec.projectRef,externalProjectID"
)

type AtlasDatabaseUserByProjectIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByProjectIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDatabaseUserByProjectIndexer {
	return &AtlasDatabaseUserByProjectIndexer{
		ctx:    ctx,
		client: client,
		logger: logger.Named(AtlasDatabaseUserByProject).Sugar(),
	}
}

func (*AtlasDatabaseUserByProjectIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserByProjectIndexer) Name() string {
	return AtlasDatabaseUserByProject
}

func (a *AtlasDatabaseUserByProjectIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return []string{user.Spec.ExternalProjectRef.ID}
	}

	if user.Spec.ProjectRef != nil && user.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *user.Spec.ProjectRef.GetObject(user.Namespace), project)
		if err != nil {
			a.logger.Errorf("unable to find project to index: %s", err)

			return nil
		}

		if project.ID() != "" {
			return []string{project.ID()}
		}
	}

	return nil
}
