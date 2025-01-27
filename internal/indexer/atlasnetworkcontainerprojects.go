//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasNetworkContainerByProjectIndex = "atlasnetworkcontainer.spec.projectRef"
)

type AtlasNetworkContainerByProjectIndexer struct {
	AtlasReferrerByProjectIndexerBase
}

func NewAtlasNetworkContainerByProjectIndexer(logger *zap.Logger) *AtlasNetworkContainerByProjectIndexer {
	return &AtlasNetworkContainerByProjectIndexer{
		AtlasReferrerByProjectIndexerBase: *NewAtlasReferrerByProjectIndexer(
			logger,
			AtlasNetworkContainerByProjectIndex,
		),
	}
}

func (*AtlasNetworkContainerByProjectIndexer) Object() client.Object {
	return &akov2.AtlasNetworkContainer{}
}
