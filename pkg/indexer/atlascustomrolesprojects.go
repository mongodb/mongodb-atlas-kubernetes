//nolint:dupl
package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasCustomRoleByProject = "atlascustomrole.spec.projectRef,externalProjectID"
)

type AtlasCustomRoleByProjectIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasCustomRoleByProjectIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasCustomRoleByProjectIndexer {
	return &AtlasCustomRoleByProjectIndexer{
		ctx:    ctx,
		client: client,
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

	if role.Spec.ProjectRef != nil && role.Spec.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *role.Spec.ProjectRef.GetObject(role.Namespace), project)
		if err != nil {
			a.logger.Errorf("unable to find project to index: %s", err)

			return nil
		}

		if project.ID() != "" {
			return []string{project.ID()}
		}
	}

	return nil
}
