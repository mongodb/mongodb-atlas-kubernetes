package atlasdeployment

import (
	"errors"
	"fmt"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	IndexStatusActive = "STEADY"
)

type searchIndexReconciler struct {
	ctx        *workflow.Context
	deployment *akov2.AtlasDeployment
	k8sClient  client.Client
	projectID  string
	indexName  string
}

func (sr *searchIndexReconciler) reconcileInternal(stateInAKO, stateInAtlas *searchindex.SearchIndex, errs []error) workflow.Result {
	sr.ctx.Log.Debugf("starting reconciliation for index '%s'", sr.indexName)
	defer sr.ctx.Log.Debugf("finished reconciliation for index '%s'", sr.indexName)

	if len(errs) != 0 {
		return sr.terminate(&searchindex.SearchIndex{
			SearchIndex:                akov2.SearchIndex{Name: sr.indexName},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         nil,
			Status:                     nil,
		}, errors.Join(errs...))
	}

	emptyInAtlas := stateInAtlas == nil
	emptyInAKO := stateInAKO == nil

	var currentStatus string
	if stateInAtlas != nil && stateInAtlas.Status != nil {
		currentStatus = *stateInAtlas.Status
	}
	// Atlas is still processing the index, nothing can't be done
	if currentStatus != IndexStatusActive && currentStatus != "" {
		return sr.progress(stateInAtlas)
	}

	switch {
	case emptyInAtlas && !emptyInAKO:
		return sr.create(stateInAKO)
	case !emptyInAtlas && !emptyInAKO:
		return sr.update(stateInAKO, stateInAtlas)
	case !emptyInAtlas && emptyInAKO:
		return sr.delete(stateInAtlas)
	default:
		return workflow.OK()
	}
}

func (sr *searchIndexReconciler) idle(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[idle] index '%s'", index.Name)
	msg := fmt.Sprintf("Atlas search index status: %s", index.GetStatus())
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(status.NewDeploymentSearchIndexStatus(
		status.SearchIndexStatusReady,
		status.WithID(index.GetID()),
		status.WithName(index.Name),
		status.WithMsg(msg))))
	ok := workflow.OK()
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, ok)
	return ok
}

// Never set the ID (status.WithID()) on terminate. It may be empty and the AKO will lose track on this index
func (sr *searchIndexReconciler) terminate(index *searchindex.SearchIndex, err error) workflow.Result {
	msg := fmt.Errorf("error with processing index %s. err: %w", index.Name, err)
	sr.ctx.Log.Debug(msg)
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(status.NewDeploymentSearchIndexStatus(
		status.SearchIndexStatusError,
		status.WithMsg(msg.Error()),
		status.WithName(index.Name),
	)))
	terminate := workflow.Terminate(status.SearchIndexStatusError, msg.Error())
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, terminate)
	return terminate
}

func (sr *searchIndexReconciler) create(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[creating] index %s", index.Name)
	defer sr.ctx.Log.Debugf("[creation finished] for index %s", index.Name)
	atlasIdx, err := index.ToAtlas()
	if err != nil {
		return sr.terminate(index, err)
	}
	respIdx, resp, err := sr.ctx.SdkClient.AtlasSearchApi.CreateAtlasSearchIndex(sr.ctx.Context,
		sr.projectID, sr.deployment.GetDeploymentName(), atlasIdx).Execute()
	if err != nil || resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return sr.terminate(index, fmt.Errorf("failed to create index: %w, status: %d", err, resp.StatusCode))
	}
	if respIdx == nil {
		return sr.terminate(index, fmt.Errorf("returned an empty index as a result of creation"))
	}
	akoIdx, err := searchindex.NewSearchIndexFromAtlas(*respIdx)
	if err != nil {
		return sr.terminate(index, fmt.Errorf("unable to convert index to AKO: %w", err))
	}
	return sr.progress(akoIdx)
}

func (sr *searchIndexReconciler) progress(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("index %s is progress: %s", index.Name, index.GetStatus())
	msg := fmt.Sprintf("Atlas search index status: %s", index.GetStatus())
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(
		status.NewDeploymentSearchIndexStatus(status.SearchIndexStatusInProgress,
			status.WithMsg(msg),
			status.WithID(index.GetID()),
			status.WithName(index.Name))))
	inProgress := workflow.InProgress(status.SearchIndexStatusInProgress, msg)
	sr.ctx.SetConditionFromResult(status.SearchIndexesReadyType, inProgress)
	return inProgress
}

func (sr *searchIndexReconciler) delete(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[deleting] index %s", index.Name)
	defer sr.ctx.Log.Debugf("[deletion finished] for index %s", index.Name)
	if index.ID == nil {
		return workflow.OK()
	}

	_, resp, err := sr.ctx.SdkClient.AtlasSearchApi.DeleteAtlasSearchIndex(sr.ctx.Context, sr.projectID,
		sr.deployment.GetDeploymentName(), *index.ID).Execute()
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNotFound || err != nil {
		return sr.terminate(index, fmt.Errorf("failed to delete index: %w, status: %d", err, resp.StatusCode))
	}

	return sr.deleted(index.Name)
}

func (sr *searchIndexReconciler) deleted(indexName string) workflow.Result {
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentUnsetSearchIndexStatus(status.NewDeploymentSearchIndexStatus("",
		status.WithName(indexName))))
	return workflow.Deleted()
}

func (sr *searchIndexReconciler) update(akoIdx, atlasIdx *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[updating] index %s", akoIdx.Name)
	defer sr.ctx.Log.Debugf("[update finished] for index %s", akoIdx.Name)

	isEqual, err := akoIdx.EqualTo(atlasIdx)
	if err != nil {
		sr.terminate(atlasIdx, err)
	}
	if isEqual {
		sr.ctx.Log.Debugf("index %s is already updated", akoIdx.Name)
		return sr.idle(atlasIdx)
	}

	sr.ctx.Log.Debugf("updating index %s...", akoIdx.Name)
	toUpdateIdx, err := akoIdx.ToAtlas()
	if err != nil {
		return sr.terminate(akoIdx, fmt.Errorf("unable to convert index to AKO: %w", err))
	}
	respIdx, resp, err := sr.ctx.SdkClient.AtlasSearchApi.UpdateAtlasSearchIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(),
		atlasIdx.GetID(), toUpdateIdx).Execute()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK || err != nil {
		return sr.terminate(akoIdx, fmt.Errorf("failed to update index: %w, status: %d", err, resp.StatusCode))
	}
	if respIdx == nil {
		return sr.terminate(akoIdx, fmt.Errorf("update returned an empty index: %w", err))
	}
	convertedIdx, err := searchindex.NewSearchIndexFromAtlas(*respIdx)
	if err != nil {
		return sr.terminate(convertedIdx, fmt.Errorf("failed to convert updated index to AKO: %w", err))
	}
	return sr.progress(convertedIdx)
}
