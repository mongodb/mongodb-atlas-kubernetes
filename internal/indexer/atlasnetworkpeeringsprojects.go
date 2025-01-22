//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasNetworkPeeringByProjectIndex = "atlasnetworkpeering.spec.projectRef"
)

type AtlasNetworkPeeringByProjectIndexer struct {
	AtlasReferrerByProjectIndexerBase
}

func NewAtlasNetworkPeeringByProjectIndexer(logger *zap.Logger) *AtlasNetworkPeeringByProjectIndexer {
	return &AtlasNetworkPeeringByProjectIndexer{
		AtlasReferrerByProjectIndexerBase: *NewAtlasReferrerByProjectIndexer(
			logger,
			AtlasNetworkPeeringByProjectIndex,
		),
	}
}

func (*AtlasNetworkPeeringByProjectIndexer) Object() client.Object {
	return &akov2.AtlasNetworkPeering{}
}
