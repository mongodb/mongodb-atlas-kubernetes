package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasIPAccessListByProjectIndex = "atlasipaccesslist.spec.projectRef"
)

type AtlasIPAccessListByProjectIndexer struct {
	AtlasReferrerByProjectIndexerBase
}

func NewAtlasIPAccessListByProjectIndexer(logger *zap.Logger) *AtlasIPAccessListByProjectIndexer {
	return &AtlasIPAccessListByProjectIndexer{
		AtlasReferrerByProjectIndexerBase: *NewAtlasReferrerByProjectIndexer(
			logger,
			AtlasIPAccessListByProjectIndex,
		),
	}
}

func (*AtlasIPAccessListByProjectIndexer) Object() client.Object {
	return &akov2.AtlasIPAccessList{}
}
