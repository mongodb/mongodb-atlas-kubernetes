package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDatabaseUserByProjectsRefIndex = "atlasdatabaseuser.spec.projectRef"
)

type AtlasDatabaseUserByProjectsRefIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByProjectsRefIndexer(logger *zap.Logger) *AtlasDatabaseUserByProjectsRefIndexer {
	return &AtlasDatabaseUserByProjectsRefIndexer{
		logger: logger.Named(AtlasDatabaseUserByProjectsRefIndex).Sugar(),
	}
}

func (*AtlasDatabaseUserByProjectsRefIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserByProjectsRefIndexer) Name() string {
	return AtlasDatabaseUserByProjectsRefIndex
}

func (a *AtlasDatabaseUserByProjectsRefIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.Project != nil && user.Spec.Project.Name != "" {
		return []string{user.Spec.Project.GetObject(user.Namespace).String()}
	}

	return nil
}
