//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDatabaseUserByProject = "atlasdatabaseuser.spec.projectRef,externalProjectID"
)

type AtlasDatabaseUserByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByProjectIndexer(logger *zap.Logger) *AtlasDatabaseUserByProjectIndexer {
	return &AtlasDatabaseUserByProjectIndexer{
		logger: logger.Named(AtlasDatabaseUserByProject).Sugar(),
	}
}

func (*AtlasDatabaseUserByProjectIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserByProjectIndexer) Name() string {
	return AtlasDatabaseUserByProject
}

func (a *AtlasDatabaseUserByProjectIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.Project != nil && user.Spec.Project.Name != "" {
		return []string{user.Spec.Project.GetObject(user.Namespace).String()}
	}

	return nil
}
