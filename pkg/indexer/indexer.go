package indexer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Indexer interface {
	Object() client.Object
	Name() string
	Keys(client.Object) []string
}

// RegisterAll registers all known indexers to the given manager.
// It uses the given logger to create a new named "indexer" logger,
// passing that to each indexer.
func RegisterAll(ctx context.Context, mgr manager.Manager, logger *zap.Logger) error {
	logger = logger.Named("indexer")
	return Register(ctx, mgr,
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
		NewAtlasDatabaseUserByProjectIndexer(ctx, mgr.GetClient(), logger),
	)
}

// Register registers the given indexers to the given manager's field indexer.
func Register(ctx context.Context, mgr manager.Manager, indexers ...Indexer) error {
	for _, indexer := range indexers {
		err := mgr.GetFieldIndexer().IndexField(ctx, indexer.Object(), indexer.Name(), indexer.Keys)
		if err != nil {
			return fmt.Errorf("error registering indexer %q: %w", indexer.Name(), err)
		}
	}

	return nil
}
