package indexer

import (
	"context"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AtlasStreamInstancesByConnectionRegistry = ".spec.connectionRegistry"
)

func NewAtlasStreamInstancesByConnectionRegistryIndex(ctx context.Context, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasStreamInstance{},
		AtlasStreamInstancesByConnectionRegistry,
		AtlasStreamInstancesByConnectionRegistryIndices,
	)
}

func AtlasStreamInstancesByConnectionRegistryIndices(object client.Object) []string {
	streamInstance := object.(*akov2.AtlasStreamInstance)
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
