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
//   <project-id>-<normalized-username>
//
// Where:
//   - <project-id> is resolved from either ExternalProjectRef.ID or the resolved AtlasProject.Status.ID
//   - <normalized-username> is produced via kube.NormalizeIdentifier(user.Spec.Username)
//
// Purpose:
//   This index enables fast lookup of AtlasDatabaseUser objects by a combination of project ID and username,
//   such as when reconciling resources linked to a connection secret.

const (
	AtlasDatabaseUserBySpecUsernameAndProjectID = "atlasdatabaseuser.projectID/spec.username"
)

type AtlasDatabaseUserBySpecUsernameIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserBySpecUsernameIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDatabaseUserBySpecUsernameIndexer {
	return &AtlasDatabaseUserBySpecUsernameIndexer{
		ctx:    ctx,
		client: client,
		logger: logger.Named(AtlasDatabaseUserBySpecUsernameAndProjectID).Sugar(),
	}
}

func (*AtlasDatabaseUserBySpecUsernameIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserBySpecUsernameIndexer) Name() string {
	return AtlasDatabaseUserBySpecUsernameAndProjectID
}

func (a *AtlasDatabaseUserBySpecUsernameIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDatabaseUser but got %T", object)
		return nil
	}

	username := user.Spec.Username
	if username == "" {
		return nil
	}

	username = kube.NormalizeIdentifier(username)
	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return []string{fmt.Sprintf("%s-%s", user.Spec.ExternalProjectRef.ID, username)}
	}

	if user.Spec.ProjectRef != nil && user.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *user.Spec.ProjectRef.GetObject(user.Namespace), project)
		if err != nil {
			a.logger.Errorf("unable to find project to index: %s", err)
			return nil
		}

		if project.ID() != "" {
			return []string{fmt.Sprintf("%s-%s", project.ID(), username)}
		}
	}

	return nil
}
