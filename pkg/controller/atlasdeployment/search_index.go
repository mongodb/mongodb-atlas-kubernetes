package atlasdeployment

import (
	"context"
	"fmt"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	DeploymentIndicesAnnotation = "mongodb.com/deployment-search-indices"
	DeploymentIndicesSeparator  = ","
	IndexToIDSeparator          = ":"
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
func handleSearchIndices(ctx *workflow.Context, k8sClient client.Client, deployment *akov2.AtlasDeployment, projectID string) error {
	ctx.Log.Debug("starting indexes processing")
	defer ctx.Log.Debug("finished indexes processing")
	// TODO: Verify if all names are unique

	if deployment.GetAnnotations() == nil {
		return fmt.Errorf("unable to get deployment annotations")
	}

	if !verifyAllIndicesNamesAreUnique(deployment.Spec.DeploymentSpec.SearchIndexes) {
		return fmt.Errorf("every index 'Name' must be unique")
	}

	previousAKOIndices := getIndicesFromAnnotations(deployment.GetAnnotations())

	var atlasExistingIndices []*searchindex.SearchIndex
	// Map[indexName][]listOfErrors
	var errors map[string][]error

	// Fetch existing indices from Atlas
	for prevIndexName, prevIndexID := range previousAKOIndices {
		resp, httpResp, err := ctx.SdkClient.AtlasSearchApi.GetAtlasSearchIndex(
			context.Background(), projectID, deployment.GetDeploymentName(), prevIndexID).Execute()
		if err != nil {
			ctx.Log.Debug("couldn't fetch index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
			continue
		}
		if resp == nil {
			ctx.Log.Debug("received an empty index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
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
	var akoIndices = make([]*searchindex.SearchIndex, len(deployment.Spec.DeploymentSpec.SearchIndexes))
	for i := range deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &deployment.Spec.DeploymentSpec.SearchIndexes[i]
		var idxConfig akov2.AtlasSearchIndexConfig

		err := k8sClient.Get(context.Background(), *akoIndex.IndexConfigRef.GetObject(deployment.Namespace), &idxConfig)
		if err != nil {
			errors[akoIndex.Name] = append(errors[akoIndex.Name], fmt.Errorf("couldn't get search index configuration. E: %w", err))
			continue
		}
		akoIndices = append(akoIndices, searchindex.NewSearchIndexFromAKO(akoIndex, &idxConfig.Spec))
	}

	// indicesToCreate := findIntersection(akoIndices, atlasIndices, leftJoin)
	// indicesToDelete := findIntersection(akoIndices, atlasIndices, rightJoin)
	// indicesToUpdate := findIntersection(akoIndices, atlasIndices, innerJoin)
	// setAnnotationsForAKO

	// Action plan:
	// 0. List all existing indexes IDs from the Annotations
	//    0.1 Get all indexes by IDs
	// 1. For each index in Atlas
	//    1.1 Get the search index configuration
	//    1.2 Build the search index
	//    1.3 Save to a list of configured indexes using Atlas DTO (!!!)
	// 2. Get the current indexes for a deployment
	// 3. Compare configured indexes, find diffs
	//    3.1 Indexes to create
	//    3.2 Indexes to update
	//    3.3 Indexes to delete
	// 4. Store new indexes IDs in the annotations
	// 5. Update status for each Index
	// indexInAKOEmpty := len(deployment.Spec.DeploymentSpec.SearchIndexes) == 0

	// callctx, cancelF := context.WithTimeout(context.Background(), 5 * time.Second)
	// defer cancelF()
	// currentIndexesInAtlas, httpResp, err := ctx.SdkClient.AtlasSearchApi.ListAtlasSearchIndexes(callctx, projectID,
	return nil
}
