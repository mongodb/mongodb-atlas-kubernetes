package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasPrivateEndpointByProjectIndex = "atlasprivateendpoint.spec.projectRef"
)

type AtlasPrivateEndpointByProjectIndexer struct {
	AtlasReferrerByProjectIndexerBase
}

func NewAtlasPrivateEndpointByProjectIndexer(logger *zap.Logger) *AtlasPrivateEndpointByProjectIndexer {
	return &AtlasPrivateEndpointByProjectIndexer{
		AtlasReferrerByProjectIndexerBase: *NewAtlasReferrerByProjectIndexer(
			logger,
			AtlasPrivateEndpointByProjectIndex,
		),
	}
}

func (*AtlasPrivateEndpointByProjectIndexer) Object() client.Object {
	return &akov2.AtlasPrivateEndpoint{}
}
