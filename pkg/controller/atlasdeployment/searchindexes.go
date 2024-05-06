package atlasdeployment

import (
	"fmt"
	"maps"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"

	internal "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	IndexTypeVector = "vectorSearch"
	IndexTypeSearch = "search"
)

func getIndexesFromDeploymentStatus(deploymentStatus status.AtlasDeploymentStatus) map[string]string {
	result := map[string]string{}
	if len(deploymentStatus.SearchIndexes) == 0 {
		return nil
	}

	for i := range deploymentStatus.SearchIndexes {
		index := &deploymentStatus.SearchIndexes[i]
		if index.ID == "" {
			continue
		}
		result[index.Name] = index.ID
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

type searchIndexesReconciler struct {
	ctx        *workflow.Context
	deployment *akov2.AtlasDeployment
	k8sClient  client.Client
	projectID  string
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
}

type IndexesErrors map[string][]error

func NewIndexesErrors() IndexesErrors {
	return map[string][]error{}
}

func (i IndexesErrors) Add(indexName string, err error) {
	if _, ok := i[indexName]; !ok {
		i[indexName] = []error{err}
	} else {
		i[indexName] = append(i[indexName], err)
	}
}

func (i IndexesErrors) GetErrors(indexName string) []error {
	val, ok := i[indexName]
	if !ok {
		return nil
	}
	return val
}

func (sr *searchIndexesReconciler) Reconcile() workflow.Result {
	if !verifyAllIndexesNamesAreUnique(sr.deployment.Spec.DeploymentSpec.SearchIndexes) {
		return sr.terminate(status.SearchIndexesNamesAreNotUnique, fmt.Errorf("every index 'Name' must be unique"))
	}
	sr.ctx.Log.Debug("all indexes names are unique")

	previousAKOIndexes := getIndexesFromDeploymentStatus(sr.deployment.Status)
	atlasIndexes := map[string]*internal.SearchIndex{}

	// Map[indexName][]listOfErrors
	indexesErrors := NewIndexesErrors()

	sr.ctx.Log.Debugf("number previous indexes: %d", len(previousAKOIndexes))
	// Fetch existing indices from Atlas
	for prevIndexName, prevIndexID := range previousAKOIndexes {
		if prevIndexID == "" {
			atlasIndexes[prevIndexName] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: prevIndexName}}
			continue
		}
		sr.ctx.Log.Debugf("restoring index %q", prevIndexName)
		resp, httpResp, err := sr.ctx.SdkClient.AtlasSearchApi.GetAtlasSearchIndex(
			sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), prevIndexID).Execute()

		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				sr.removeIndexStatus(prevIndexName)
				continue
			}
			e := fmt.Errorf("couldn't fetch index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
			indexesErrors.Add(prevIndexName, e)
			atlasIndexes[prevIndexName] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: prevIndexName}}
			sr.ctx.Log.Debug(e)
			continue
		}
		if resp == nil {
			e := fmt.Errorf("received an empty index. ID: %s. Status code: %d, E: %w", prevIndexID, httpResp.StatusCode, err)
			indexesErrors.Add(prevIndexName, e)
			atlasIndexes[prevIndexName] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: prevIndexName}}
			sr.ctx.Log.Debug(e)
			continue
		}
		akoIndex, err := internal.NewSearchIndexFromAtlas(*resp)
		if err != nil {
			e := fmt.Errorf("unable to convert index to AKO. Name: %s, ID: %s, E: %w", prevIndexName, prevIndexID, err)
			atlasIndexes[prevIndexName] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: prevIndexName}}
			indexesErrors.Add(prevIndexName, e)
			continue
		}

		atlasIndexes[akoIndex.Name] = akoIndex
	}

	// Build indexes for AKO
	akoIndexes := map[string]*internal.SearchIndex{}
	for i := range sr.deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &sr.deployment.Spec.DeploymentSpec.SearchIndexes[i]
		sr.ctx.Log.Debugf("reading AKO index: '%s'", akoIndex.Name)

		var indexInternal *internal.SearchIndex
		switch akoIndex.Type {
		case IndexTypeSearch:
			if akoIndex.Search == nil {
				e := fmt.Errorf("index '%s' has type '%s' but the spec is missing", akoIndex.Name, IndexTypeSearch)
				indexesErrors.Add(akoIndex.Name, e)
				akoIndexes[akoIndex.Name] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: akoIndex.Name}}
				continue
			}

			var idxConfig akov2.AtlasSearchIndexConfig
			err := sr.k8sClient.Get(sr.ctx.Context, *akoIndex.Search.SearchConfigurationRef.GetObject(sr.deployment.Namespace), &idxConfig)
			if err != nil {
				e := fmt.Errorf("can not get search index configuration for index '%s'. E: %w", akoIndex.Name, err)
				indexesErrors.Add(akoIndex.Name, e)
				akoIndexes[akoIndex.Name] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: akoIndex.Name}}
				continue
			}
			indexInternal = internal.NewSearchIndexFromAKO(akoIndex, &idxConfig.Spec)
		case IndexTypeVector:
			// Vector index doesn't require any external configuration
			indexInternal = internal.NewSearchIndexFromAKO(akoIndex, &akov2.AtlasSearchIndexConfigSpec{})
		default:
			e := fmt.Errorf("index %q has unknown type %q. Can be either %s or %s",
				akoIndex.Name, akoIndex.Type, IndexTypeSearch, IndexTypeVector)
			indexesErrors.Add(akoIndex.Name, e)
			akoIndexes[akoIndex.Name] = &internal.SearchIndex{SearchIndex: akov2.SearchIndex{Name: akoIndex.Name}}
		}
		akoIndexes[akoIndex.Name] = indexInternal
	}

	allIndexes := map[string]*internal.SearchIndex{}
	// note: the order matters! first Atlas, then AKO so we have most up-to-date desired state
	maps.Copy(allIndexes, atlasIndexes)
	maps.Copy(allIndexes, akoIndexes)

	sr.ctx.Log.Debugf("number indexes to process: %d", len(allIndexes))
	if len(allIndexes) == 0 {
		return sr.empty()
	}

	results := make([]workflow.Result, 0, len(allIndexes))
	for i := range allIndexes {
		current := allIndexes[i]

		var akoIdx, atlasIdx *internal.SearchIndex

		if val, ok := akoIndexes[current.Name]; ok {
			akoIdx = val
		}
		if val, ok := atlasIndexes[current.Name]; ok {
			atlasIdx = val
		}

		results = append(results, (&searchIndexReconciler{
			ctx:        sr.ctx,
			deployment: sr.deployment,
			k8sClient:  sr.k8sClient,
			projectID:  sr.projectID,
			indexName:  current.Name,
		}).reconcileInternal(akoIdx, atlasIdx, indexesErrors.GetErrors(current.Name)))
	}

	allDeleted := true
	for i := range results {
		if results[i].IsInProgress() || !results[i].IsOk() {
			return sr.progress()
		}
		allDeleted = allDeleted && results[i].IsDeleted()
	}
	if allDeleted {
		return sr.empty()
	}

	return sr.idle()
}

// This is a special method to curate index status in case index is not in Atlas, not in spec but in status
func (sr *searchIndexesReconciler) removeIndexStatus(prevIndexName string) {
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentUnsetSearchIndexStatus(
		status.NewDeploymentSearchIndexStatus("", status.WithName(prevIndexName))))
}

func (sr *searchIndexesReconciler) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	sr.ctx.Log.Error(err)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	result := workflow.Terminate(reason, errMsg)
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, result)
	return result
}

func (sr *searchIndexesReconciler) progress() workflow.Result {
	result := workflow.InProgress(status.SearchIndexesNotReady, "not all indexes are in READY state")
	sr.ctx.SetConditionFromResult(status.SearchIndexStatusReady, result)
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
