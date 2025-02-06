package indexer

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

type AtlasReferrerByProjectIndexerBase struct {
	name   string
	logger *zap.SugaredLogger
}

func NewAtlasReferrerByProjectIndexer(logger *zap.Logger, name string) *AtlasReferrerByProjectIndexerBase {
	return &AtlasReferrerByProjectIndexerBase{
		name:   name,
		logger: logger.Named(name).Sugar(),
	}
}

func (rb *AtlasReferrerByProjectIndexerBase) Name() string {
	return rb.name
}

func (rb *AtlasReferrerByProjectIndexerBase) Keys(object client.Object) []string {
	pro, ok := object.(project.ProjectReferrerObject)
	if !ok {
		rb.logger.Errorf("expected a project.ProjectReferrerObject but got %T", object)
		return nil
	}

	pdr := pro.ProjectDualRef()
	if pdr == nil || pdr.ProjectRef == nil || pdr.ProjectRef.Name == "" {
		return nil
	}

	return []string{pdr.ProjectRef.GetObject(pro.GetNamespace()).String()}
}

func ProjectsIndexMapperFunc[L client.ObjectList](indexerName string, listGenFn func() L, reqsFn requestsFunc[L], kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		project, ok := obj.(*akov2.AtlasProject)
		if !ok {
			logger.Warnf("watching AtlasProject but got %T", obj)
			return nil
		}

		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexerName,
				client.ObjectKeyFromObject(project).String(),
			),
		}
		list := listGenFn()
		err := kubeClient.List(ctx, list, listOpts)
		if err != nil {
			logger.Errorf("failed to list from indexer %s: %v", indexerName, err)
			return nil
		}
		return reqsFn(list)
	}
}
