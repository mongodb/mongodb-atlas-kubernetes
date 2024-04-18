package atlasdeployment

import (
	"context"
	"fmt"
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
	DeploymentIndicesAnnotation = "mongodb.com/deployment-search-indices"
	DeploymentIndicesSeparator  = ","
	IndexToIDSeparator          = ":"
	IndexStatusFormat           = "SearchIndex-%s"
)

// getIndicesFromAnnotations returns a map IndexName -> IndexID
func getIndicesFromAnnotations(in map[string]string) map[string]string {
	result := map[string]string{}
	indices, ok := in[DeploymentIndicesAnnotation]
	if !ok {
		return nil
	}
	indexNameIDPairs := strings.Split(indices, DeploymentIndicesSeparator)
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

func verifyAllIndicesNamesAreUnique(indices []akov2.SearchIndex) bool {
	buff := make(map[string]bool, len(indices))
	for i := range indices {
		if _, ok := buff[indices[i].Name]; ok {
			return false
		}
		buff[indices[i].Name] = true
	}
	return true
}

func findIndicesIntersection(akoIndices, atlasIndices []*searchindex.SearchIndex, intersection IntersectionType) []searchindex.SearchIndex {
	var result []searchindex.SearchIndex
	switch intersection {
	case ToCreate:
		for i := range akoIndices {
			found := false
			for j := range atlasIndices {
				if akoIndices[i].Name == atlasIndices[j].Name {
					found = true
					continue
				}
			}
			if !found {
				if akoIndices[i] != nil {
					result = append(result, *(akoIndices[i]))
				}
			}
		}

	case ToUpdate:
		for i := range akoIndices {
			for j := range atlasIndices {
				if akoIndices[i].Name == atlasIndices[j].Name {
					if akoIndices[i] != nil {
						result = append(result, *(akoIndices[i]))
					}
				}
			}
		}

	case ToDelete:
		for i := range atlasIndices {
			found := false
			for j := range akoIndices {
				if akoIndices[j].Name == atlasIndices[i].Name {
					found = true
					continue
				}
			}
			if !found {
				if atlasIndices[i] != nil {
					result = append(result, *(atlasIndices[i]))
				}
			}
		}
	}

	return result
}

type searchIndexReconciler struct {
	ctx          *workflow.Context
	deployment   *akov2.AtlasDeployment
	k8sClient    client.Client
	projectID    string
	akoIndices   map[string]*searchindex.SearchIndex
	atlasIndices map[string]*searchindex.SearchIndex
	atlasErrors  map[string]error
}

func handleSearchIndices(ctx *workflow.Context, k8sClient client.Client, deployment *akov2.AtlasDeployment, projectID string) workflow.Result {
	ctx.Log.Debug("starting indexes processing")
	defer ctx.Log.Debug("finished indexes processing")

	reconciler := &searchIndexReconciler{
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
	// indexInAKOEmpty := len(deployment.Spec.DeploymentSpec.SearchIndexes) == 0

	// callctx, cancelF := context.WithTimeout(context.Background(), 5 * time.Second)
	// defer cancelF()
	// currentIndexesInAtlas, httpResp, err := ctx.SdkClient.AtlasSearchApi.ListAtlasSearchIndexes(callctx, projectID,
}

func (sr *searchIndexReconciler) ensureSearchIndex(index *searchindex.SearchIndex) workflow.Result {

	// !isPresentInAtlas(index) && isPresentInAKO(index) -> create
	// isPresentInAtlas(index) && isPresentInAKO(index) -> update
	// isPresentInAtlas(index) && !isPresentInAKO(index) -> delete
	return workflow.OK()
}
func (sr *searchIndexReconciler) Reconcile() workflow.Result {

	if !verifyAllIndicesNamesAreUnique(sr.deployment.Spec.DeploymentSpec.SearchIndexes) {
		return sr.terminate(status.SearchIndexesNamesAreNotUnique, fmt.Errorf("every index 'Name' must be unique"))
	}
	// TODO: do a per index state machine

	previousAKOIndices := getIndicesFromAnnotations(sr.deployment.GetAnnotations())

	var atlasExistingIndices []*searchindex.SearchIndex
	// Map[indexName][]listOfErrors
	var errors map[string][]error

	// Fetch existing indices from Atlas
	for prevIndexName, prevIndexID := range previousAKOIndices {
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
			errors[prevIndexName] = append(errors[prevIndexName],
				fmt.Errorf("unable to convert index to AKO. Name: %s, ID: %s, E: %w", prevIndexName, prevIndexID, err))
		}
		atlasExistingIndices = append(atlasExistingIndices, akoIndex)
	}

	// Build indices for AKO
	var akoIndices = make([]*searchindex.SearchIndex, len(sr.deployment.Spec.DeploymentSpec.SearchIndexes))
	for i := range sr.deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &sr.deployment.Spec.DeploymentSpec.SearchIndexes[i]
		var idxConfig akov2.AtlasSearchIndexConfig

		err := sr.k8sClient.Get(context.Background(), *akoIndex.IndexConfigRef.GetObject(sr.deployment.Namespace), &idxConfig)
		if err != nil {
			errors[akoIndex.Name] = append(errors[akoIndex.Name], fmt.Errorf("couldn't get search index configuration. E: %w", err))
			continue
		}
		akoIndices = append(akoIndices, searchindex.NewSearchIndexFromAKO(akoIndex, &idxConfig.Spec))
	}

	//var workflowResults map[string]workflow.Result
	//
	//for i := range

	//indicesToCreate := findIndicesIntersection(akoIndices, atlasExistingIndices, ToCreate)
	//indicesToUpdate := findIndicesIntersection(akoIndices, atlasExistingIndices, ToUpdate)
	//indicesToDelete := findIndicesIntersection(akoIndices, atlasExistingIndices, ToDelete)
	return workflow.OK()
}

func (sr *searchIndexReconciler) handleCreating(indices []*searchindex.SearchIndex) {
	sr.ctx.Log.Debug("starting handling indices creating")
	defer sr.ctx.Log.Debug("starting handling indices creation")
	//for _, index := range indices {
	//_, _, err := sr.ctx.SdkClient.AtlasSearchApi.CreateAtlasSearchIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), index.ToAtlas()).Execute()
	//if err != nil {
	//	sr.indexStatus[index.Name] = workflow.SearchIndexesReady
	//}

	//return nil
}

func (sr *searchIndexReconciler) handleUpdating(indices []*searchindex.SearchIndex) {
	sr.ctx.Log.Debug("starting handling indices updates")
	defer sr.ctx.Log.Debug("starting handling indices updates")
}

func (sr *searchIndexReconciler) handleDeleting(indices []*searchindex.SearchIndex) {
	sr.ctx.Log.Debug("starting handling indices deletion")
	defer sr.ctx.Log.Debug("starting handling indices deletion")
}

func (sr *searchIndexReconciler) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	sr.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err.Error())
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, result)
	return result
}

func (sr *searchIndexReconciler) empty() workflow.Result {
	sr.ctx.UnsetCondition(status.SearchIndexesReadyType)
	return workflow.OK()
}

func (sr *searchIndexReconciler) idle() workflow.Result {
	sr.ctx.SetConditionTrue(status.SearchIndexesReadyType)
	return workflow.OK()
}
