package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasNetworkPeeringByContainerIndex = "atlasnetworkpeering.spec.container"
)

type NetworkPeeringByContainerIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasNetworkPeeringByContainerIndexer(logger *zap.Logger) *NetworkPeeringByContainerIndexer {
	return &NetworkPeeringByContainerIndexer{
		logger: logger.Named(AtlasNetworkPeeringByContainerIndex).Sugar(),
	}
}

func (*NetworkPeeringByContainerIndexer) Object() client.Object {
	return &akov2.AtlasNetworkPeering{}
}

func (*NetworkPeeringByContainerIndexer) Name() string {
	return AtlasNetworkPeeringByContainerIndex
}

func (p *NetworkPeeringByContainerIndexer) Keys(object client.Object) []string {
	peering, ok := object.(*akov2.AtlasNetworkPeering)
	if !ok {
		p.logger.Errorf("expected *akov2.AtlasNetworkPeering but got %T", object)
		return nil
	}

	containerName := peering.Spec.ContainerRef.Name
	if containerName != "" {
		return []string{containerName}
	}

	return nil
}
