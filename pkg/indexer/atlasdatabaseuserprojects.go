//nolint:dupl
package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDatabaseUserByProject = "atlasdatabaseuser.spec.projectRef,externalProjectID"
)

type AtlasDatabaseUserByProjectIndexer struct {
	ctx    context.Context
	client client.Client
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserByProjectIndexer(ctx context.Context, client client.Client, logger *zap.Logger) *AtlasDatabaseUserByProjectIndexer {
	return &AtlasDatabaseUserByProjectIndexer{
		ctx:    ctx,
		client: client,
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

	// TODO: remove in a separate PR
	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return []string{user.Spec.ExternalProjectRef.ID}
	}
	// TODO: end

	if user.Spec.Project != nil && user.Spec.Project.Name != "" {
		project := &akov2.AtlasProject{}
		err := a.client.Get(a.ctx, *user.Spec.Project.GetObject(user.Namespace), project)
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
