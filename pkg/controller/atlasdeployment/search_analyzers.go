package atlasdeployment

import (
	"context"
	"fmt"
	"reflect"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"

	"go.mongodb.org/atlas/mongodbatlas"
)

func syncCustomAnalyzers(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	analyzers []mdbv1.CustomAnalyzer,
) ([]string, error) {
	existingAnalyzers, _, err := service.Client.Search.ListAnalyzers(ctx, projectID, clusterName, nil)
	if err != nil {
		return []string{}, fmt.Errorf("failed to retrieve list of custom analyzers: %w", err)
	}

	if len(analyzers) == 0 {
		if len(existingAnalyzers) == 0 {
			service.Log.Debug("no custom-analyzers to sync")
			return []string{}, nil
		}

		service.Log.Debug("removing all custom-analyzers")
		_, _, err = service.Client.Search.UpdateAllAnalyzers(ctx, projectID, clusterName, []*mongodbatlas.SearchAnalyzer{})
		if err != nil {
			return []string{}, fmt.Errorf("failed to remove custom analyzers: %w", err)
		}

		return []string{}, nil
	}

	if !hasCustomAnalyzersChanged(analyzers, existingAnalyzers) {
		return []string{}, nil
	}

	data := make([]*mongodbatlas.SearchAnalyzer, 0, len(analyzers))

	for _, analyzer := range analyzers {
		data = append(data, analyzer.ToAtlas())
	}

	service.Log.Debug("updating all custom-analyzers")
	result, _, err := service.Client.Search.UpdateAllAnalyzers(ctx, projectID, clusterName, data)
	if err != nil {
		return []string{}, fmt.Errorf("failed to create/update custom analyzers: %w", err)
	}

	statuses := make([]string, 0, len(result))

	for _, v := range result {
		statuses = append(statuses, v.Name)
	}

	return statuses, nil
}

func hasCustomAnalyzersChanged(
	desiredAnalyzers []mdbv1.CustomAnalyzer,
	existingAnalyzers []*mongodbatlas.SearchAnalyzer,
) bool {
	if len(desiredAnalyzers) != len(existingAnalyzers) {
		return true
	}

	existingAnalyzersMap := set.FromSlice(existingAnalyzers, func(item *mongodbatlas.SearchAnalyzer) string {
		return item.Name
	})

	for _, analyzer := range desiredAnalyzers {
		existingAnalyzer, ok := existingAnalyzersMap[analyzer.Name]

		if !ok {
			return true
		}

		if !reflect.DeepEqual(analyzer, existingAnalyzer) {
			return true
		}
	}

	return false
}
