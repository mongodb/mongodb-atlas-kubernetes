package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasStreamInstanceByConnectionIndex = ".spec.connectionRegistry"
)

type AtlasStreamInstanceByConnectionIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasStreamInstanceByConnectionIndexer(logger *zap.Logger) *AtlasStreamInstanceByConnectionIndexer {
	return &AtlasStreamInstanceByConnectionIndexer{
		logger: logger.Named(AtlasStreamInstanceByConnectionIndex).Sugar(),
	}
}

func (*AtlasStreamInstanceByConnectionIndexer) Object() client.Object {
	return &akov2.AtlasStreamInstance{}
}
func (*AtlasStreamInstanceByConnectionIndexer) Name() string {
	return AtlasStreamInstanceByConnectionIndex
}

func (a *AtlasStreamInstanceByConnectionIndexer) Keys(object client.Object) []string {
	streamInstance, ok := object.(*akov2.AtlasStreamInstance)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasStreamInstance but got %T", object)
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
