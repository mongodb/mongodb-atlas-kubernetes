package atlasdeployment

import (
	"fmt"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	IndexStatusActive = "ACTIVE"
)

type searchIndexReconciler struct {
	ctx        *workflow.Context
	deployment *akov2.AtlasDeployment
	k8sClient  client.Client
	projectID  string
	indexName  string
}

func newStatusEntry(index *searchindex.SearchIndex, newStatus status.IndexStatus, msg string) status.DeploymentSearchIndexStatus {
	var sts string
	if newStatus == "" {
		sts = pointer.GetOrDefault(index.Status, status.SearchIndexStatusUnknown)
	}
	return status.DeploymentSearchIndexStatus{
		Name:   index.Name,
		ID:     pointer.GetOrDefault(index.ID, ""),
		Status: status.IndexStatus(sts),
		// TODO: propagate to index DTO!
		ConfigRef: common.ResourceRefNamespaced{},
		Message:   msg,
	}
}

func (sr *searchIndexReconciler) Reconcile(stateInAKO, stateInAtlas *searchindex.SearchIndex, errs []error) workflow.Result {
	sr.ctx.Log.Debugf("starting reconciliation for index '%s'", sr.indexName)
	defer sr.ctx.Log.Debugf("finished reconciliation for index '%s'", sr.indexName)

	// TODO terminate if there are errors
	//if len(errs) != 0 {
	//	return sr.terminate()
	//}

	emptyInAtlas := stateInAKO == nil
	emptyInAKO := stateInAtlas == nil

	var currentStatus string
	if stateInAtlas != nil && stateInAtlas.Status != nil {
		currentStatus = *stateInAtlas.Status
	}

	switch currentStatus {
	case IndexStatusActive, "":
		break
	default:
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
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(status.NewDeploymentSearchIndexStatus(
		status.SearchIndexStatusReady, status.WithID(index.GetID()))))
	return workflow.OK()
}
func (sr *searchIndexReconciler) terminate(index *searchindex.SearchIndex, err error) workflow.Result {
	msg := fmt.Errorf("error with processing index '%s'. err: %w", index.Name, err)
	sr.ctx.Log.Debug(msg)
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(status.NewDeploymentSearchIndexStatus(
		status.SearchIndexStatusError,
		status.WithMsg(msg.Error()),
		status.WithID(index.GetID()),
		status.WithName(index.Name),
	)))
	return workflow.Terminate(status.SearchIndexStatusError, msg.Error())
}

func (sr *searchIndexReconciler) create(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[creating] index '%s'", index.Name)
	defer sr.ctx.Log.Debugf("[creation finished] for index '%s'", index.Name)
	atlasIdx, err := index.ToAtlas()
	if err != nil {
		return sr.terminate(index, err)
	}
	respIdx, resp, err := sr.ctx.SdkClient.AtlasSearchApi.CreateAtlasSearchIndex(sr.ctx.Context,
		sr.projectID, sr.deployment.GetDeploymentName(), atlasIdx).Execute()
	if err != nil || resp.StatusCode != http.StatusCreated {
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
	sr.ctx.Log.Debugf("index '%s' is progress: %s", index.Name, index.GetStatus())
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(
		status.NewDeploymentSearchIndexStatus(status.SearchIndexStatusInProgress,
			status.WithMsg(*index.Status), status.WithID(index.GetID()), status.WithName(index.Name))))
	return workflow.InProgress(status.SearchIndexStatusInProgress, index.GetStatus())
}

func (sr *searchIndexReconciler) delete(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[deleting] index '%s'", index.Name)
	defer sr.ctx.Log.Debugf("[deletion finished] for index '%s'", index.Name)
	if index.ID == nil {
		return workflow.OK()
	}
	// TODO: get the index first and check the status. if it's not ACTIVE, return terminate

	_, resp, err := sr.ctx.SdkClient.AtlasSearchApi.DeleteAtlasSearchIndex(sr.ctx.Context, sr.projectID,
		sr.deployment.GetDeploymentName(), *index.ID).Execute()
	if resp.StatusCode != http.StatusNoContent || err != nil {
		return sr.terminate(index, fmt.Errorf("failed to delete index: %w, status: %d", err, resp.StatusCode))
	}
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentUnsetSearchIndexStatus(status.NewDeploymentSearchIndexStatus("")))
	return sr.progress(index)
}

func (sr *searchIndexReconciler) update(akoIdx, atlasIdx *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[updating] index '%s'", akoIdx.Name)
	defer sr.ctx.Log.Debugf("[updating] for index '%s'", akoIdx.Name)

	// TODO: find diff
	if akoIdx.EqualTo(atlasIdx) {
		return sr.idle(atlasIdx)
	}

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
