package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasStreamInstancesByConnectionRegistry = ".spec.connectionRegistry"
)

func NewAtlasStreamInstancesByConnectionRegistryIndex(ctx context.Context, logger *zap.SugaredLogger, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasStreamInstance{},
		AtlasStreamInstancesByConnectionRegistry,
		AtlasStreamInstancesByConnectionRegistryIndices(logger.Named("indexers").Named(AtlasStreamInstancesByConnectionRegistry)),
	)
}

func AtlasStreamInstancesByConnectionRegistryIndices(logger *zap.SugaredLogger) client.IndexerFunc {
	return func(object client.Object) []string {
		streamInstance, ok := object.(*akov2.AtlasStreamInstance)
		if !ok {
			logger.Errorf("expected *akov2.AtlasStreamInstance but got %T", object)
			return nil
		}

		if len(streamInstance.Spec.ConnectionRegistry) == 0 {
			return nil
		}

		registry := streamInstance.Spec.ConnectionRegistry
		indices := make([]string, 0, len(registry))
		for i := range registry {
			key := registry[i].GetObject(streamInstance.GetNamespace())
			indices = append(indices, key.String())
		}

		return indices
	}
}
