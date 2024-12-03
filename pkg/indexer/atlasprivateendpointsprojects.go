package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasPrivateEndpointByProjectIndex = "atlasprivateendpoint.spec.projectRef"
)

type AtlasPrivateEndpointByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasPrivateEndpointByProjectIndexer(logger *zap.Logger) *AtlasPrivateEndpointByProjectIndexer {
	return &AtlasPrivateEndpointByProjectIndexer{
		logger: logger.Named(AtlasPrivateEndpointByProjectIndex).Sugar(),
	}
}

func (*AtlasPrivateEndpointByProjectIndexer) Object() client.Object {
	return &akov2.AtlasPrivateEndpoint{}
}

func (*AtlasPrivateEndpointByProjectIndexer) Name() string {
	return AtlasPrivateEndpointByProjectIndex
}

func (a *AtlasPrivateEndpointByProjectIndexer) Keys(object client.Object) []string {
	pe, ok := object.(*akov2.AtlasPrivateEndpoint)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasPrivateEndpoint but got %T", object)
		return nil
	}

	if pe.Spec.Project == nil || pe.Spec.Project.Name == "" {
		return nil
	}

	return []string{pe.Spec.Project.GetObject(pe.GetNamespace()).String()}
}
