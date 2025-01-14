package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
