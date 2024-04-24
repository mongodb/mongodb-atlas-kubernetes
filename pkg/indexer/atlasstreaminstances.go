package indexer

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
	streamInstance, ok := object.(*akov2.AtlasStreamInstance)
	if !ok {
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
