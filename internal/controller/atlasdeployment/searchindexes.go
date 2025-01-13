package atlasdeployment

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
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
	ctx           *workflow.Context
	deployment    *akov2.AtlasDeployment
	k8sClient     client.Client
	projectID     string
	searchService searchindex.AtlasSearchIdxService
}

func handleSearchIndexes(ctx *workflow.Context, k8sClient client.Client, searchService searchindex.AtlasSearchIdxService, deployment *akov2.AtlasDeployment, projectID string) workflow.Result {
	ctx.Log.Debug("starting indexes processing")
	defer ctx.Log.Debug("finished indexes processing")

	reconciler := &searchIndexesReconciler{
		ctx:           ctx,
		k8sClient:     k8sClient,
		deployment:    deployment,
		projectID:     projectID,
		searchService: searchService,
	}

	return reconciler.Reconcile()
}

func (sr *searchIndexesReconciler) Reconcile() workflow.Result {
	if !verifyAllIndexesNamesAreUnique(sr.deployment.Spec.DeploymentSpec.SearchIndexes) {
		return sr.terminate(api.SearchIndexesNamesAreNotUnique, fmt.Errorf("every index 'Name' must be unique"))
	}
	sr.ctx.Log.Debug("all indexes names are unique")

	type tuple struct {
		previous *status.DeploymentSearchIndexStatus
		spec     *akov2.SearchIndex
	}
	allIndexes := map[string]tuple{}

	sr.ctx.Log.Debugf("number previous indexes: %d", len(sr.deployment.Status.SearchIndexes))
	// Build indexes based on previously reconciled indexes
	for i := range sr.deployment.Status.SearchIndexes {
		searchIndexStatus := sr.deployment.Status.SearchIndexes[i]
		if searchIndexStatus.ID == "" {
			continue
		}
		allIndexes[searchIndexStatus.Name] = tuple{previous: &searchIndexStatus}
	}

	// Build indexes based on the spec
	for i := range sr.deployment.Spec.DeploymentSpec.SearchIndexes {
		akoIndex := &sr.deployment.Spec.DeploymentSpec.SearchIndexes[i]
		sr.ctx.Log.Debugf("reading AKO index: '%s'", akoIndex.Name)
		var entry tuple
		if _, ok := allIndexes[akoIndex.Name]; ok {
			entry = allIndexes[akoIndex.Name]
		}
		entry.spec = akoIndex
		allIndexes[akoIndex.Name] = entry
	}

	sr.ctx.Log.Debugf("number indexes to process: %d", len(allIndexes))
	if len(allIndexes) == 0 {
		return sr.empty()
	}

	results := make([]workflow.Result, 0, len(allIndexes))
	for indexName, val := range allIndexes {
		results = append(results, (&searchIndexReconcileRequest{
			ctx:           sr.ctx,
			deployment:    sr.deployment,
			k8sClient:     sr.k8sClient,
			projectID:     sr.projectID,
			indexName:     indexName,
			searchService: sr.searchService,
		}).Handle(val.spec, val.previous))
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

func (sr *searchIndexesReconciler) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	sr.ctx.Log.Error(err)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	result := workflow.Terminate(reason, errMsg)
	sr.ctx.SetConditionFromResult(api.SearchIndexesReadyType, result)
	return result
}

func (sr *searchIndexesReconciler) progress() workflow.Result {
	result := workflow.InProgress(api.SearchIndexesNotReady, "not all indexes are in READY state")
	sr.ctx.SetConditionFromResult(status.SearchIndexStatusReady, result)
	return result
}

func (sr *searchIndexesReconciler) empty() workflow.Result {
	sr.ctx.UnsetCondition(api.SearchIndexesReadyType)
	return workflow.OK()
}

func (sr *searchIndexesReconciler) idle() workflow.Result {
	sr.ctx.SetConditionTrue(api.SearchIndexesReadyType)
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentRemoveStatusesWithEmptyIDs())
	return workflow.OK()
}
