package atlasdeployment

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-multierror"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func ensureAtlasSearch(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	deployment *mdbv1.AtlasDeployment,
) workflow.Result {
	service.Log.Info("starting reconciliation of AtlasSearch")
	atlasSearch := getAtlasSearch(deployment)

	service.Log.Info("syncing custom-analyzers...")
	analyzersStatus, err := syncCustomAnalyzers(ctx, service, projectID, deployment.GetDeploymentName(), atlasSearch.CustomAnalyzers)
	service.Log.Debugf("%d custom-analyzer(s) synchronized: %v", len(analyzersStatus), analyzersStatus)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintf("failed to sync custom-analyzers: %s", err))
	}

	service.Log.Debugf("collecting all indexes information...")
	desiredIndexes, existingIndexes, err := collectIndexesInfo(ctx, service, projectID, deployment.GetDeploymentName(), atlasSearch.Databases)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintf("failed to collect information about the indexes: %s", err.Error()))
	}

	indexesStatus := make([]*status.AtlasIndex, 0)
	if deployment.Status.AtlasSearch != nil {
		indexesStatus = deployment.Status.AtlasSearch.Indexes
	}

	service.Log.Debugf("syncing indexes...")
	newStatuses, err := syncIndexes(ctx, service, projectID, deployment.GetDeploymentName(), desiredIndexes, existingIndexes, indexesStatus)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	service.Log.Debugf("%d index(es) synchronized", len(newStatuses))

	return ensureAtlasSearchStatus(service, newStatuses, analyzersStatus)
}

func ensureAtlasSearchStatus(
	service *workflow.Context,
	indexesStatus []*status.AtlasIndex,
	analyzersStatus []string,
) workflow.Result {
	if len(indexesStatus) == 0 && len(analyzersStatus) == 0 {
		service.EnsureStatusOption(status.AtlasDeploymentAtlasSearch(nil))
		service.UnsetCondition(status.AtlasSearchReadyType)

		return workflow.OK()
	}

	service.EnsureStatusOption(status.AtlasDeploymentAtlasSearch(&status.AtlasSearch{
		CustomAnalyzers: analyzersStatus,
		Indexes:         indexesStatus,
	}))

	inProgress := false
	hasFailed := false
	for _, indexStatus := range indexesStatus {
		if indexStatus.Status == status.IndexStatusInProgress {
			inProgress = true
		}

		if indexStatus.Status == status.IndexStatusFailed {
			hasFailed = true
		}
	}

	if inProgress {
		service.SetConditionFalse(status.AtlasSearchReadyType)
		return workflow.InProgress(workflow.AtlasSearchInProgress, "indexes are being created")
	}

	if hasFailed {
		service.SetConditionFalse(status.AtlasSearchReadyType)
		return workflow.Terminate(workflow.AtlasSearchIncomplete, "some of the indexes were not successfully created")
	}

	service.SetConditionTrue(status.AtlasSearchReadyType)
	return workflow.OK()
}

func syncIndexes(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	desiredIndexes map[string]*mongodbatlas.SearchIndex,
	existingIndexes map[string]*mongodbatlas.SearchIndex,
	statuses []*status.AtlasIndex,
) ([]*status.AtlasIndex, error) {
	var errors error
	newStatuses := make([]*status.AtlasIndex, 0, len(desiredIndexes))
	toCreate, toUpdate, toDelete, statusUpdate := sortIndexes(desiredIndexes, existingIndexes, statuses)

	createdStatuses, err := createIndexes(ctx, service, projectID, clusterName, toCreate)
	if err != nil {
		errors = multierror.Append(err)
	}
	newStatuses = append(newStatuses, createdStatuses...)

	updateStatuses, err := updateIndexes(ctx, service, projectID, clusterName, toUpdate)
	if err != nil {
		errors = multierror.Append(err)
	}
	newStatuses = append(newStatuses, updateStatuses...)

	err = deleteIndexes(ctx, service, projectID, clusterName, toDelete)
	if err != nil {
		errors = multierror.Append(err)
	}

	newStatuses = append(newStatuses, updateIndexesStatuses(service, statusUpdate)...)

	return newStatuses, errors
}

func collectIndexesInfo(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	databases []mdbv1.AtlasSearchDatabase,
) (desired, existing map[string]*mongodbatlas.SearchIndex, err error) {
	desired = map[string]*mongodbatlas.SearchIndex{}
	existing = map[string]*mongodbatlas.SearchIndex{}

	for _, database := range databases {
		for _, collection := range database.Collections {
			existingIndexes, _, err := service.Client.Search.ListIndexes(ctx, projectID, clusterName, database.Database, collection.CollectionName, nil)
			if err != nil {
				return nil, nil, err
			}

			for _, existingIndex := range existingIndexes {
				existing[fmt.Sprintf("%s.%s.%s", existingIndex.Database, existingIndex.CollectionName, existingIndex.Name)] = existingIndex
			}

			for _, index := range collection.Indexes {
				desired[fmt.Sprintf("%s.%s.%s", database.Database, collection.CollectionName, index.Name)] = index.ToAtlas(database.Database, collection.CollectionName)
			}
		}
	}

	return desired, existing, err
}

func sortIndexes(
	desiredIndexes map[string]*mongodbatlas.SearchIndex,
	existingIndexes map[string]*mongodbatlas.SearchIndex,
	statuses []*status.AtlasIndex,
) (toCreate []*mongodbatlas.SearchIndex, toUpdate map[string]*mongodbatlas.SearchIndex, toDelete []string, statusUpdate []*mongodbatlas.SearchIndex) {
	toDelete = make([]string, 0, len(existingIndexes))
	for identifier, existingIndex := range existingIndexes {
		if _, ok := desiredIndexes[identifier]; !ok {
			toDelete = append(toDelete, existingIndex.IndexID)
		}
	}

	for _, indexStatus := range statuses {
		_, desiredOk := desiredIndexes[fmt.Sprintf("%s.%s.%s", indexStatus.Database, indexStatus.CollectionName, indexStatus.Name)]

		if !desiredOk {
			toDelete = append(toDelete, indexStatus.ID)
		}
	}

	toCreate = make([]*mongodbatlas.SearchIndex, 0, len(desiredIndexes))
	toUpdate = map[string]*mongodbatlas.SearchIndex{}
	for identifier, desiredIndex := range desiredIndexes {
		existingIndex, ok := existingIndexes[identifier]

		if !ok {
			toCreate = append(toCreate, desiredIndex)

			continue
		}

		if !hasIndexChanged(desiredIndex, existingIndex) {
			toUpdate[existingIndex.IndexID] = desiredIndex

			continue
		}

		statusUpdate = append(statusUpdate, existingIndex)
	}

	return toCreate, toUpdate, toDelete, statusUpdate
}

func createIndexes(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	indexes []*mongodbatlas.SearchIndex,
) ([]*status.AtlasIndex, error) {
	var allErrors error
	statuses := make([]*status.AtlasIndex, 0, len(indexes))

	for _, index := range indexes {
		service.Log.Debugf("creating index %s at %s.%s", index.Name, index.Database, index.CollectionName)
		atlasIndex, _, err := service.Client.Search.CreateIndex(
			ctx,
			projectID,
			clusterName,
			index,
		)
		if err != nil {
			allErrors = multierror.Append(err)
		}

		statuses = append(statuses, status.NewStatusFromAtlas(atlasIndex, err))
	}

	return statuses, allErrors
}

func updateIndexes(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	indexes map[string]*mongodbatlas.SearchIndex,
) ([]*status.AtlasIndex, error) {
	var allErrors error
	statuses := make([]*status.AtlasIndex, 0, len(indexes))

	for indexID, index := range indexes {
		service.Log.Debugf("updating index %s at %s.%s", index.Name, index.Database, index.CollectionName)
		atlasIndex, _, err := service.Client.Search.UpdateIndex(
			ctx,
			projectID,
			clusterName,
			indexID,
			index,
		)
		if err != nil {
			allErrors = multierror.Append(err)
		}

		statuses = append(statuses, status.NewStatusFromAtlas(atlasIndex, err))
	}

	return statuses, allErrors
}

func deleteIndexes(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	indexes []string,
) error {
	var allErrors error

	for _, id := range indexes {
		indexID := id

		service.Log.Debugf("deleting index %s", indexID)
		_, err := service.Client.Search.DeleteIndex(
			ctx,
			projectID,
			clusterName,
			indexID,
		)
		if err != nil {
			allErrors = multierror.Append(err)
		}
	}

	return allErrors
}

func updateIndexesStatuses(
	service *workflow.Context,
	indexes []*mongodbatlas.SearchIndex,
) []*status.AtlasIndex {
	statuses := make([]*status.AtlasIndex, 0, len(indexes))

	for _, index := range indexes {
		existingIndex := index

		service.Log.Debugf("updating status of the index %s", index.Name)
		statuses = append(statuses, status.NewStatusFromAtlas(existingIndex, nil))
	}

	return statuses
}

func hasIndexChanged(desired, existing *mongodbatlas.SearchIndex) bool {
	if desired.Name != existing.Name {
		return false
	}

	if desired.Analyzer != existing.Analyzer {
		return false
	}

	if desired.SearchAnalyzer != existing.SearchAnalyzer {
		return false
	}

	if desired.Mappings.Dynamic != existing.Mappings.Dynamic {
		return false
	}

	if desired.Mappings.Fields == nil && existing.Mappings.Fields != nil {
		return false
	}

	if desired.Mappings.Fields != nil && existing.Mappings.Fields == nil {
		return false
	}

	if desired.Mappings.Fields != nil && existing.Mappings.Fields != nil {
		for key, value := range *desired.Mappings.Fields {
			atlasValue, ok := (*existing.Mappings.Fields)[key]

			if !ok {
				return false
			}

			if !reflect.DeepEqual(value, atlasValue) {
				return false
			}
		}
	}

	if len(desired.Synonyms) != len(existing.Synonyms) {
		return false
	}

	return true
}

func getAtlasSearch(deployment *mdbv1.AtlasDeployment) *mdbv1.AtlasSearch {
	if deployment.Spec.DeploymentSpec != nil &&
		deployment.Spec.DeploymentSpec.AtlasSearch != nil {
		return deployment.Spec.DeploymentSpec.AtlasSearch
	}

	if deployment.Spec.AdvancedDeploymentSpec != nil &&
		deployment.Spec.AdvancedDeploymentSpec.AtlasSearch != nil {
		return deployment.Spec.AdvancedDeploymentSpec.AtlasSearch
	}

	return &mdbv1.AtlasSearch{}
}
