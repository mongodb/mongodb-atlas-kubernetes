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
//   <project-id>-<normalized-deployment-name>
//
// Where:
//   - <project-id> is resolved from either ExternalProjectRef.ID or the resolved AtlasProject.Status.ID
//   - <normalized-deployment-name> is produced via kube.NormalizeIdentifier(deployment.Spec.DeploymentSpec.Name)
//
// Purpose:
//   This index allows fast lookup of AtlasDeployment resources by project ID and deployment name,
//   particularly useful for identifying which cluster a user has access to.

const (
	AtlasDeploymentBySpecNameAndProjectID = "atlasdeployment.projectID/spec.name"
)

type AtlasDeploymentBySpecNameIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDeploymentBySpecNameIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDeploymentBySpecNameIndexer {
	return &AtlasDeploymentBySpecNameIndexer{
		ctx:    ctx,
		client: client,
		logger: logger.Named(AtlasDeploymentBySpecNameAndProjectID).Sugar(),
	}
}

func (*AtlasDeploymentBySpecNameIndexer) Object() client.Object {
	return &akov2.AtlasDeployment{}
}

func (*AtlasDeploymentBySpecNameIndexer) Name() string {
	return AtlasDeploymentBySpecNameAndProjectID
}

func (a *AtlasDeploymentBySpecNameIndexer) Keys(object client.Object) []string {
	deployment, ok := object.(*akov2.AtlasDeployment)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDeployment but got %T", object)
		return nil
	}

	name := deployment.GetDeploymentName()
	if name == "" {
		return nil
	}
	name = kube.NormalizeIdentifier(name)

	// First check ExternalProjectRef
	if deployment.Spec.ExternalProjectRef != nil && deployment.Spec.ExternalProjectRef.ID != "" {
		return []string{fmt.Sprintf("%s-%s", deployment.Spec.ExternalProjectRef.ID, name)}
	}

	// Fallback to resolving ProjectRef (name-based)
	if deployment.Spec.ProjectRef != nil && deployment.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *deployment.Spec.ProjectRef.GetObject(deployment.Namespace), project)
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
