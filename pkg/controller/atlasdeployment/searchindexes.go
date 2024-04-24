package atlasdeployment

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type IntersectionType byte

const (
	ToCreate IntersectionType = iota
	ToUpdate
	ToDelete
)

const (
	DeploymentIndexesAnnotation = "mongodb.com/deployment-search-indices"
	DeploymentIndexesSeparator  = ","
	IndexToIDSeparator          = ":"
	IndexStatusFormat           = "SearchIndex-%s"
)

var (
	ErrNoIndexConfig = "index configuration is not available"
)

// getIndexesFromAnnotations returns a map IndexName -> IndexID
func getIndexesFromAnnotations(in map[string]string) map[string]string {
	result := map[string]string{}
	indexes, ok := in[DeploymentIndexesAnnotation]
	if !ok {
		return nil
	}
	indexNameIDPairs := strings.Split(indexes, DeploymentIndexesSeparator)
	for _, pair := range indexNameIDPairs {
		res := strings.Split(pair, IndexToIDSeparator)
		if len(res) != 2 {
			continue
		}
		if res[1] == "" {
			continue
		}
		result[res[0]] = res[1]
	}
	return result
}

func verifyAllIndexesNamesAreUnique(indexes []akov2.SearchIndex) bool {
	buff := make(map[string]bool, len(indexes))
	for i := range indexes {
		if _, ok := buff[indexes[i].Name]; ok {
			return false
		}
		buff[indexes[i].Name] = true
	}
	return true
}

func findIndexesIntersection(akoIndexes, atlasIndexes []*searchindex.SearchIndex, intersection IntersectionType) []searchindex.SearchIndex {
	var result []searchindex.SearchIndex
	switch intersection {
	case ToCreate:
		for i := range akoIndexes {
			found := false
			for j := range atlasIndexes {
				if akoIndexes[i].Name == atlasIndexes[j].Name {
					found = true
					continue
				}
			}
			if !found {
				if akoIndexes[i] != nil {
					result = append(result, *(akoIndexes[i]))
				}
			}
		}

	case ToUpdate:
		for i := range akoIndexes {
			for j := range atlasIndexes {
				if akoIndexes[i].Name == atlasIndexes[j].Name {
					if akoIndexes[i] != nil {
						result = append(result, *(akoIndexes[i]))
					}
				}
			}
		}

	case ToDelete:
		for i := range atlasIndexes {
			found := false
			for j := range akoIndexes {
				if akoIndexes[j].Name == atlasIndexes[i].Name {
					found = true
					continue
				}
			}
			if !found {
				if atlasIndexes[i] != nil {
					result = append(result, *(atlasIndexes[i]))
				}
			}
		}
	}

	return result
}

type searchIndexesReconciler struct {
	ctx         *workflow.Context
	deployment  *akov2.AtlasDeployment
	k8sClient   client.Client
	projectID   string
	atlasErrors map[string]error
}

func handleSearchIndexes(ctx *workflow.Context, k8sClient client.Client, deployment *akov2.AtlasDeployment, projectID string) workflow.Result {
	ctx.Log.Debug("starting indexes processing")
	defer ctx.Log.Debug("finished indexes processing")

	reconciler := &searchIndexesReconciler{
		ctx:        ctx,
		k8sClient:  k8sClient,
		deployment: deployment,
		projectID:  projectID,
	}

	return reconciler.Reconcile()
	// Action plan:
	// 0. List all existing indexes IDs from the Annotations
	//    0.1 Get all indexes by IDs
	// 1. For each index in Atlas
	//    1.1 Get the search index configuration
	//    1.2 Build the search index
	//    1.3 Save to a list of configured indexes using Atlas DTO (!!!)
	// 2. Get the current indexes for a deployment
	// 3. Compare configured indexes
	//    3.1 Indexes to create
	//    3.2 Indexes to update
	//      3.2.1 Find diffs. Update only those that are different
	//    3.3 Indexes to delete
	// 4. Store new indexes IDs in the annotations
	// 5. Update status for each Index
}

func (sr *searchIndexesReconciler) Reconcile() workflow.Result {
	if !verifyAllIndexesNamesAreUnique(sr.deployment.Spec.DeploymentSpec.SearchIndexes) {
		return sr.terminate(status.SearchIndexesNamesAreNotUnique, fmt.Errorf("every index 'Name' must be unique"))
	}

	// TODO: refactor that to be taken from the Status
	previousAKOIndexes := getIndexesFromAnnotations(sr.deployment.GetAnnotations())
	atlasIndexes := map[string]*searchindex.SearchIndex{}

	// Map[indexName][]listOfErrors
	var indexesErrors map[string][]error

	// Fetch existing indices from Atlas
	for prevIndexName, prevIndexID := range previousAKOIndexes {
		resp, httpResp, err := sr.ctx.SdkClient.AtlasSearchApi.GetAtlasSearchIndex(
			context.Background(), sr.projectID, sr.deployment.GetDeploymentName(), prevIndexID).Execute()
		// TODO: store the errors in the sr.errors
		if err != nil {
			sr.ctx.Log.Debug("couldn't fetch index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
			continue
		}
		if resp == nil {
			sr.ctx.Log.Debug("received an empty index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
			continue
		}

		akoIndex, err := searchindex.NewSearchIndexFromAtlas(*resp)
		if err != nil {
			indexesErrors[prevIndexName] = append(indexesErrors[prevIndexName],
				fmt.Errorf("unable to convert index to AKO. Name: %s, ID: %s, E: %w", prevIndexName, prevIndexID, err))
		}

		atlasIndexes[akoIndex.Name] = akoIndex
	}

	// Build indexes for AKO
	akoIndexes := map[string]*searchindex.SearchIndex{}
	for i := range sr.deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &sr.deployment.Spec.DeploymentSpec.SearchIndexes[i]
		var idxConfig akov2.AtlasSearchIndexConfig

		err := sr.k8sClient.Get(context.Background(), *akoIndex.IndexConfigRef.GetObject(sr.deployment.Namespace), &idxConfig)
		if err != nil {
			indexesErrors[akoIndex.Name] = append(indexesErrors[akoIndex.Name], fmt.Errorf("couldn't get search index configuration. E: %w", err))
			continue
		}
		akoIndexes[akoIndex.Name] = searchindex.NewSearchIndexFromAKO(akoIndex, &idxConfig.Spec)
	}
	var allIndexes map[string]*searchindex.SearchIndex
	// note: order matters! first Atlas, then AKO so we have most up-to-date desired state
	maps.Copy(allIndexes, atlasIndexes)
	maps.Copy(allIndexes, akoIndexes)

	if len(allIndexes) == 0 {
		return sr.empty()
	}

	results := make([]workflow.Result, 0, len(allIndexes))
	for i := range allIndexes {
		current := allIndexes[i]

		var akoIdx, atlasIdx *searchindex.SearchIndex
		var errors []error

		if _, ok := akoIndexes[current.Name]; ok {
			akoIdx = current
		}
		if _, ok := atlasIndexes[current.Name]; ok {
			atlasIdx = current
		}
		if _, ok := indexesErrors[current.Name]; ok {
			errors = indexesErrors[current.Name]
		}

		results = append(results, (&searchIndexReconciler{
			ctx:         sr.ctx,
			deployment:  sr.deployment,
			k8sClient:   sr.k8sClient,
			projectID:   sr.projectID,
			indexErrors: errors,
		}).Reconcile(akoIdx, atlasIdx))
	}

	for i := range results {
		if !results[i].IsOk() {
			return sr.terminate(status.SearchIndexesNotAllReady, nil)
		}
	}

	return sr.idle()
}

func (sr *searchIndexesReconciler) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	sr.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err.Error())
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, result)
	return result
}

func (sr *searchIndexesReconciler) empty() workflow.Result {
	sr.ctx.UnsetCondition(status.SearchIndexesReadyType)
	return workflow.OK()
}

func (sr *searchIndexesReconciler) idle() workflow.Result {
	sr.ctx.SetConditionTrue(status.SearchIndexesReadyType)
	return workflow.OK()
}
