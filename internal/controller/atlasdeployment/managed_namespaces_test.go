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

package atlasdeployment

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
)

func TestEnsureManagedNamespaces(t *testing.T) {
	projectID := "test-project"
	deploymentName := "test-deployment"
	exampleManagedNamespace := akov2.ManagedNamespace{
		Db:                     "test-db",
		Collection:             "test-collection",
		CustomShardKey:         "test-shard-key",
		NumInitialChunks:       10,
		PresplitHashedZones:    pointer.MakePtr(false),
		IsCustomShardKeyHashed: pointer.MakePtr(true),
		IsShardKeyUnique:       pointer.MakePtr(true),
	}

	for _, tc := range []struct {
		name              string
		managedNamespaces []akov2.ManagedNamespace
		deploymentAPI     deployment.AtlasDeploymentsService
		isOK              bool
	}{
		{
			name: "No managed namespace in AKO or Atlas (no op)",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return(nil, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Get errors",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return(nil, errors.New("test GET error"))

				return service
			}(),
			isOK: false,
		},
		{
			name: "Managed namespace in AKO but not Atlas (create)",
			managedNamespaces: []akov2.ManagedNamespace{
				exampleManagedNamespace,
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return(nil, nil)
				service.EXPECT().CreateManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Create errors",
			managedNamespaces: []akov2.ManagedNamespace{
				exampleManagedNamespace,
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return(nil, nil)
				service.EXPECT().CreateManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(errors.New("test create error"))

				return service
			}(),
			isOK: false,
		},
		{
			name: "Managed namespace in Atlas but not AKO (delete)",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return([]akov2.ManagedNamespace{exampleManagedNamespace}, nil)
				service.EXPECT().DeleteManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Delete errors",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return([]akov2.ManagedNamespace{exampleManagedNamespace}, nil)
				service.EXPECT().DeleteManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(errors.New("test delete error"))

				return service
			}(),
			isOK: false,
		},
		{
			name:              "Managed namespace the same in both AKO and Atlas (no op)",
			managedNamespaces: []akov2.ManagedNamespace{exampleManagedNamespace},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return([]akov2.ManagedNamespace{exampleManagedNamespace}, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Managed namespace different in AKO and Atlas (update)",
			managedNamespaces: []akov2.ManagedNamespace{
				{
					Db:                     "new-test-db",
					Collection:             "new-test-collection",
					CustomShardKey:         "new-test-shard-key",
					NumInitialChunks:       12,
					PresplitHashedZones:    pointer.MakePtr(false),
					IsCustomShardKeyHashed: pointer.MakePtr(true),
					IsShardKeyUnique:       pointer.MakePtr(true),
				},
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetManagedNamespaces(context.Background(), projectID, deploymentName).Return([]akov2.ManagedNamespace{exampleManagedNamespace}, nil)
				service.EXPECT().DeleteManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(nil)
				service.EXPECT().CreateManagedNamespace(context.Background(), projectID, deploymentName, mock.AnythingOfType("*v1.ManagedNamespace")).Return(nil)

				return service
			}(),
			isOK: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &AtlasDeploymentReconciler{}
			ctx := &workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
			}

			result := r.ensureManagedNamespaces(
				ctx,
				tc.deploymentAPI,
				projectID,
				string(akov2.TypeGeoSharded),
				tc.managedNamespaces,
				deploymentName,
			)

			equal := (result.IsOk() == tc.isOK)
			if !equal {
				t.Log(result.GetMessage())
			}
			require.True(t, equal)
		})
	}
}
