//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasDataFederationByProject = "atlasdatafederation.spec.project"
)

type AtlasDataFederationByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDataFederationByProjectIndexer(logger *zap.Logger) *AtlasDataFederationByProjectIndexer {
	return &AtlasDataFederationByProjectIndexer{
		logger: logger.Named(AtlasDatabaseUserByProject).Sugar(),
	}
}

func (*AtlasDataFederationByProjectIndexer) Object() client.Object {
	return &akov2.AtlasDataFederation{}
}

func (*AtlasDataFederationByProjectIndexer) Name() string {
	return AtlasDataFederationByProject
}

func (a *AtlasDataFederationByProjectIndexer) Keys(object client.Object) []string {
	datafederation, ok := object.(*akov2.AtlasDataFederation)
	if !ok {
		a.logger.Errorf("expected *AtlasDataFederation but got %T", object)
		return nil
	}

	if datafederation.Spec.Project.IsEmpty() {
		return nil
	}

	return []string{datafederation.Spec.Project.GetObject(datafederation.Namespace).String()}
}
