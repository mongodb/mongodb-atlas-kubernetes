package indexer

import (
	"context"

	"go.uber.org/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasSearchIndexToDeploymentRegistry = ".spec.deploymentSpec.searchIndexes"
)

func NewAtlasSearchIndexConfigsToDeploymentIndex(ctx context.Context, logger *zap.SugaredLogger, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasDeployment{},
		AtlasSearchIndexToDeploymentRegistry,
		AtlasSearchIndexKeysToDeployment(logger.Named("indexers").Named(AtlasSearchIndexToDeploymentRegistry)),
	)
}

func AtlasSearchIndexKeysToDeployment(logger *zap.SugaredLogger) client.IndexerFunc {
	return func(object client.Object) []string {
		deployment, ok := object.(*akov2.AtlasDeployment)
		if !ok {
			logger.Errorf("expected *akov2.AtlasSearchIndexConfig but got %T", object)
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
}
