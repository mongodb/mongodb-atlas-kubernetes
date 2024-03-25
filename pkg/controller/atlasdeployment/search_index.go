package atlasdeployment

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	DeploymentIndicesAnnotation = "mongodb.com/deployment-search-indices"
	DeploymentIndicesSeparator  = ","
)

func handleSearchIndices(ctx *workflow.Context, k8sClient client.Client, deployment *mdbv1.AtlasDeployment, projectID string) error {
	ctx.Log.Debug("starting indexes processing")
	defer ctx.Log.Debug("finished indexes processing")
	// TODO: Verify if all names are unique

	// Step 0
	var previousAKOIndexIDs []string
	if deployment.GetAnnotations() == nil {
		return fmt.Errorf("unable to get deployment annotations")
	}
	if indices, ok := deployment.GetAnnotations()[DeploymentIndicesAnnotation]; ok {
		previousAKOIndexIDs = strings.Split(indices, DeploymentIndicesSeparator)
	}

	var atlasExistingIndices []*admin.ClusterSearchIndex

	// Fetch existing indices from Atlas
	for _, prevIndexID := range previousAKOIndexIDs {
		resp, httpResp, err := ctx.SdkClient.AtlasSearchApi.GetAtlasSearchIndex(context.Background(), projectID, deployment.GetDeploymentName(), prevIndexID).Execute()
		if err != nil {
			ctx.Log.Debug("couldn't fetch index. ID: %s. Status code: %d, Err: %w", prevIndexID, httpResp.StatusCode, err)
			continue
		}
		atlasExistingIndices = append(atlasExistingIndices, resp)
	}

	// Map[indexName][]listOfErrors
	var errors map[string][]error
	// Build indices for AKO
	//var akoIndices = make([]*admin.ClusterSearchIndex, len(deployment.Spec.DeploymentSpec.SearchIndexes))
	for i := range deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &deployment.Spec.DeploymentSpec.SearchIndexes[i]
		var idxConfig mdbv1.AtlasSearchIndexConfig

		err := k8sClient.Get(context.Background(), *akoIndex.IndexConfigRef.GetObject(deployment.Namespace), &idxConfig)
		if err != nil {
			errors[akoIndex.Name] = append(errors[akoIndex.Name], fmt.Errorf("couldn't get search index configuration. E: %w\r\n", err))
			continue
		}
		//akoIndices = append(akoIndices, akoIndex.ToAtlas(&idxConfig.Spec))
	}

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
