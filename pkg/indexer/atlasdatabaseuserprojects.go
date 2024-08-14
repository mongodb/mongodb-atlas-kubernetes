package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDatabaseUserByProjectsIndex = "atlasdatabaseuser.spec.project"
)

type AtlasDatabaseUserByProjectsIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByProjectsIndexer(logger *zap.Logger) *AtlasDatabaseUserByProjectsIndexer {
	return &AtlasDatabaseUserByProjectsIndexer{
		logger: logger.Named(AtlasDatabaseUserByProjectsIndex).Sugar(),
	}
}

func (*AtlasDatabaseUserByProjectsIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserByProjectsIndexer) Name() string {
	return AtlasDatabaseUserByProjectsIndex
}

func (a *AtlasDatabaseUserByProjectsIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.Project != nil && user.Spec.Project.Name != "" {
		return []string{user.Spec.Project.GetObject(user.Namespace).String()}
	}

	if user.Spec.AtlasRef != nil && user.Spec.AtlasRef.ID != "" {
		return []string{user.Spec.AtlasRef.ID}
	}

	return nil
}
