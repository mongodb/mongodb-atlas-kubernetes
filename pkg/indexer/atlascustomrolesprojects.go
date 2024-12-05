//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasCustomRoleByProject = "atlascustomrole.spec.projectRef,externalProjectID"
)

type AtlasCustomRoleByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasCustomRoleByProjectIndexer(logger *zap.Logger) *AtlasCustomRoleByProjectIndexer {
	return &AtlasCustomRoleByProjectIndexer{
		logger: logger.Named(AtlasCustomRoleByProject).Sugar(),
	}
}

func (*AtlasCustomRoleByProjectIndexer) Object() client.Object {
	return &akov2.AtlasCustomRole{}
}

func (*AtlasCustomRoleByProjectIndexer) Name() string {
	return AtlasCustomRoleByProject
}

func (a *AtlasCustomRoleByProjectIndexer) Keys(object client.Object) []string {
	role, ok := object.(*akov2.AtlasCustomRole)
	if !ok {
		a.logger.Errorf("expected *v1.AtlasCustomRole but got %T", object)
		return nil
	}

	if role.Spec.Project != nil && role.Spec.Project.Name != "" {
		return []string{role.Spec.Project.GetObject(role.Namespace).String()}
	}

	return nil
}
