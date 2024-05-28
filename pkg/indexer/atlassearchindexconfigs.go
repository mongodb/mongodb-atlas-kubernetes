package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDeploymentBySearchIndexIndex = "atlasdeployment.spec.deploymentSpec.searchIndexes"
)

type AtlasDeploymentBySearchIndexIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDeploymentBySearchIndexIndexer(logger *zap.Logger) *AtlasDeploymentBySearchIndexIndexer {
	return &AtlasDeploymentBySearchIndexIndexer{
		logger: logger.Named(AtlasDeploymentBySearchIndexIndex).Sugar(),
	}
}

func (*AtlasDeploymentBySearchIndexIndexer) Object() client.Object {
	return &akov2.AtlasDeployment{}
}

func (*AtlasDeploymentBySearchIndexIndexer) Name() string {
	return AtlasDeploymentBySearchIndexIndex
}

func (a *AtlasDeploymentBySearchIndexIndexer) Keys(object client.Object) []string {
	deployment, ok := object.(*akov2.AtlasDeployment)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDeployment but got %T", object)
		return nil
	}

	if deployment.Spec.DeploymentSpec == nil {
		return nil
	}

	if len(deployment.Spec.DeploymentSpec.SearchIndexes) == 0 {
		return nil
	}

	searchIndexes := deployment.Spec.DeploymentSpec.SearchIndexes

	result := make([]string, 0, len(searchIndexes))
	for i := range searchIndexes {
		idx := &searchIndexes[i]
		if idx.Search == nil {
			continue
		}

		// searchIndexConfigKey -> deploymentName
		searchIndexKey := idx.Search.SearchConfigurationRef.GetObject(deployment.GetNamespace())
		result = append(result, searchIndexKey.String())
	}

	return result
}
