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
	AtlasDeploymentByProject = "atlasdeployment.spec.projectRef,externalProjectID"
)

type AtlasDeploymentByProjectIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDeploymentByProjectIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDeploymentByProjectIndexer {
	return &AtlasDeploymentByProjectIndexer{
		ctx:    ctx,
		client: client,
		logger: logger.Named(AtlasDeploymentByProject).Sugar(),
	}
}

func (*AtlasDeploymentByProjectIndexer) Object() client.Object {
	return &akov2.AtlasDeployment{}
}

func (*AtlasDeploymentByProjectIndexer) Name() string {
	return AtlasDeploymentByProject
}

func (a *AtlasDeploymentByProjectIndexer) Keys(object client.Object) []string {
	deployment, ok := object.(*akov2.AtlasDeployment)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDeployment but got %T", object)
		return nil
	}

	if deployment.Spec.ExternalProjectRef != nil && deployment.Spec.ExternalProjectRef.ID != "" {
		return []string{deployment.Spec.ExternalProjectRef.ID}
	}

	if deployment.Spec.ProjectRef != nil && deployment.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *deployment.Spec.ProjectRef.GetObject(deployment.Namespace), project)
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
