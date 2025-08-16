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
	"fmt"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

// Index Format:
//   <project-id>-<normalized-federation-name>
//
// Where:
//   - <project-id> is resolved from either ExternalProjectRef.ID or the resolved AtlasProject.Status.ID
//   - <normalized-federation-name> is produced via kube.NormalizeIdentifier(dataFederation.Spec.Name)
//
// Purpose:
//   This index allows fast lookup of AtlasDataFederation resources by project ID and federation name,
//   particularly useful for identifying which cluster a user has access to.

const (
	AtlasDataFederationBySpecNameAndProjectID = "atlasdatafederation.projectID/spec.name"
)

type AtlasDataFederationBySpecNameIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDataFederationBySpecNameIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDataFederationBySpecNameIndexer {
	return &AtlasDataFederationBySpecNameIndexer{
		ctx:    ctx,
		client: client,
		logger: logger.Named(AtlasDataFederationBySpecNameAndProjectID).Sugar(),
	}
}

func (*AtlasDataFederationBySpecNameIndexer) Object() client.Object {
	return &akov2.AtlasDataFederation{}
}

func (*AtlasDataFederationBySpecNameIndexer) Name() string {
	return AtlasDataFederationBySpecNameAndProjectID
}

func (a *AtlasDataFederationBySpecNameIndexer) Keys(object client.Object) []string {
	df, ok := object.(*akov2.AtlasDataFederation)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDataFederation but got %T", object)
		return nil
	}

	name := df.Spec.Name
	if name == "" {
		return nil
	}
	name = kube.NormalizeIdentifier(name)

	if df.Spec.Project.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *df.Spec.Project.GetObject(df.Namespace), project)
		if err != nil {
			a.logger.Errorf("unable to find project to index: %s", err)
			return nil
		}

		if project.ID() != "" {
			return []string{fmt.Sprintf("%s-%s", project.ID(), name)}
		}
	}

	return nil
}
