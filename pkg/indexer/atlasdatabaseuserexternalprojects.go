package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDatabaseUserByExternalProjectsRefIndex = "atlasdatabaseuser.spec.externalProjectRef"
)

type AtlasDatabaseUserByExternalProjectsRefIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByExternalProjectsRefIndexer(logger *zap.Logger) *AtlasDatabaseUserByExternalProjectsRefIndexer {
	return &AtlasDatabaseUserByExternalProjectsRefIndexer{
		logger: logger.Named(AtlasDatabaseUserByExternalProjectsRefIndex).Sugar(),
	}
}

func (*AtlasDatabaseUserByExternalProjectsRefIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserByExternalProjectsRefIndexer) Name() string {
	return AtlasDatabaseUserByExternalProjectsRefIndex
}

func (a *AtlasDatabaseUserByExternalProjectsRefIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return []string{user.Spec.ExternalProjectRef.ID}
	}

	return nil
}
