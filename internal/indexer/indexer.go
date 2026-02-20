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

package indexer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	connectionsecretindexer "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/indexer"
	indexer "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/experimental/indexers"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

type Indexer interface {
	Object() client.Object
	Name() string
	Keys(client.Object) []string
}

// RegisterAll registers all known indexers to the given manager.
// It uses the given logger to create a new named "indexer" logger,
// passing that to each indexer.
func RegisterAll(ctx context.Context, c cluster.Cluster, logger *zap.Logger) error {
	logger = logger.Named("indexer")
	indexers := []Indexer{}
	indexers = append(indexers,
		NewAtlasBackupScheduleByBackupPolicyIndexer(logger),
		NewAtlasDeploymentByBackupScheduleIndexer(logger),
		NewAtlasDeploymentBySearchIndexIndexer(logger),
		NewAtlasStreamConnectionBySecretIndexer(logger),
		NewAtlasStreamInstanceByProjectIndexer(logger),
		NewAtlasStreamInstanceByConnectionIndexer(logger),
		NewAtlasProjectByBackupCompliancePolicyIndexer(logger),
		NewAtlasProjectByConnectionSecretIndexer(logger),
		NewAtlasProjectByTeamIndexer(logger),
		NewAtlasFederatedAuthBySecretsIndexer(logger),
		NewAtlasDatabaseUserBySecretsIndexer(logger),
		NewAtlasDatabaseUserByCredentialIndexer(logger),
		NewAtlasDeploymentByCredentialIndexer(logger),
		NewAtlasDatabaseUserByProjectIndexer(ctx, c.GetClient(), logger),
		NewAtlasDataFederationByProjectIndexer(logger),
		NewAtlasCustomRoleByCredentialIndexer(logger),
		NewAtlasCustomRoleByProjectIndexer(logger),
		NewAtlasPrivateEndpointByCredentialIndexer(logger),
		NewAtlasPrivateEndpointByProjectIndexer(logger),
		NewAtlasIPAccessListCredentialsByCredentialIndexer(logger),
		NewAtlasIPAccessListByProjectIndexer(logger),
		NewAtlasNetworkPeeringByCredentialIndexer(logger),
		NewAtlasNetworkPeeringByProjectIndexer(logger),
		NewAtlasNetworkContainerByCredentialIndexer(logger),
		NewAtlasNetworkContainerByProjectIndexer(logger),
		NewAtlasNetworkPeeringByContainerIndexer(logger),
		NewAtlasThirdPartyIntegrationByProjectIndexer(logger),
		NewAtlasThirdPartyIntegrationByCredentialIndexer(logger),
		NewAtlasThirdPartyIntegrationBySecretsIndexer(logger),
		NewAtlasOrgSettingsByConnectionSecretIndexer(logger),
	)
	if version.IsExperimental() {
		// add experimental indexers here
		indexers = append(indexers,
			connectionsecretindexer.NewFlexClusterByGroupIdIndexer(logger),
			connectionsecretindexer.NewClusterByGroupIdIndexer(logger),
			connectionsecretindexer.NewDatabaseUserBySecretIndexer(ctx, c.GetClient(), logger),

			indexer.NewFlexClusterByGroupIndexer(logger),
			indexer.NewClusterByGroupIndexer(logger),
			indexer.NewDatabaseUserBySecretIndexer(logger),
			indexer.NewDatabaseUserByGroupIndexer(logger),
		)
	}
	return Register(ctx, c, indexers...)
}

// Register registers the given indexers to the given manager's field indexer.
func Register(ctx context.Context, c cluster.Cluster, indexers ...Indexer) error {
	for _, indexer := range indexers {
		err := c.GetFieldIndexer().IndexField(ctx, indexer.Object(), indexer.Name(), indexer.Keys)
		if err != nil {
			return fmt.Errorf("error registering indexer %q: %w", indexer.Name(), err)
		}
	}

	return nil
}
