//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasNetworkPeeringByProjectIndex = "atlasnetworkpeering.spec.projectRef"
)

type AtlasNetworkPeeringByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasNetworkPeeringByProjectIndexer(logger *zap.Logger) *AtlasNetworkPeeringByProjectIndexer {
	return &AtlasNetworkPeeringByProjectIndexer{
		logger: logger.Named(AtlasNetworkPeeringByProjectIndex).Sugar(),
	}
}

func (*AtlasNetworkPeeringByProjectIndexer) Object() client.Object {
	return &akov2.AtlasNetworkPeering{}
}

func (*AtlasNetworkPeeringByProjectIndexer) Name() string {
	return AtlasNetworkPeeringByProjectIndex
}

func (a *AtlasNetworkPeeringByProjectIndexer) Keys(object client.Object) []string {
	pe, ok := object.(*akov2.AtlasNetworkPeering)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasNetworkPeering but got %T", object)
		return nil
	}

	if pe.Spec.ProjectRef == nil || pe.Spec.ProjectRef.Name == "" {
		return nil
	}

	return []string{pe.Spec.ProjectRef.GetObject(pe.GetNamespace()).String()}
}
