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
	AtlasDataFederationByProjectID = "atlasdatafederation.spec.projectID"
)

type AtlasDataFederationByProjectIDIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDataFederationByProjectIDIndexer(ctx context.Context, c client.Client, logger *zap.Logger) *AtlasDataFederationByProjectIDIndexer {
	return &AtlasDataFederationByProjectIDIndexer{
		ctx:    ctx,
		client: c,
		logger: logger.Named(AtlasDataFederationByProjectID).Sugar(),
	}
}

func (*AtlasDataFederationByProjectIDIndexer) Object() client.Object {
	return &akov2.AtlasDataFederation{}
}

func (*AtlasDataFederationByProjectIDIndexer) Name() string {
	return AtlasDataFederationByProjectID
}

func (a *AtlasDataFederationByProjectIDIndexer) Keys(object client.Object) []string {
	df, ok := object.(*akov2.AtlasDataFederation)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDataFederation but got %T", object)
		return nil
	}

	if df.Spec.Project.Name != "" {
		project := &akov2.AtlasProject{}
		if err := a.client.Get(a.ctx, *df.Spec.Project.GetObject(df.Namespace), project); err != nil {
			a.logger.Errorf("unable to find project to index: %s", err)
			return nil
		}
		if project.ID() != "" {
			return []string{project.ID()}
		}
	}

	return nil
}
