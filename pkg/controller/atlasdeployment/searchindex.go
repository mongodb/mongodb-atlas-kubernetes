package atlasdeployment

import (
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	IndexStatusActive = "STEADY"
)

type searchIndexReconciler struct {
	ctx           *workflow.Context
	deployment    *akov2.AtlasDeployment
	k8sClient     client.Client
	projectID     string
	indexName     string
	searchService searchindex.AtlasSearchIdxService
}

func (sr *searchIndexReconciler) Reconcile(spec *akov2.SearchIndex, previous *status.DeploymentSearchIndexStatus) workflow.Result {
	var stateInAtlas, stateInAKO *searchindex.SearchIndex

	if previous != nil {
		var err error
		sr.ctx.Log.Debugf("restoring index %q from status", previous.Name)
		stateInAtlas, err = sr.searchService.GetIndex(
			sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), previous.Name, previous.ID)
		if err != nil {
			if errors.Is(err, searchindex.ErrNotFound) {
				// transition into deleted state to clear out any previous search indexes
				sr.deleted(previous.Name)
				// then immediately transition into terminated state to requeue
				return sr.terminate(stateInAtlas, err)
			}
			return sr.terminate(stateInAtlas, err)
		}
	}

	if spec != nil {
		sr.ctx.Log.Debugf("restoring index %q from spec", spec.Name)
		internalState := &searchindex.SearchIndex{
			SearchIndex:                akov2.SearchIndex{Name: spec.Name},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         nil,
			Status:                     nil,
		}
		switch spec.Type {
		case IndexTypeSearch:
			if spec.Search == nil {
				err := fmt.Errorf("index '%s' has type '%s' but the spec is missing", spec.Name, IndexTypeSearch)
				return sr.terminate(internalState, err)
			}

			var idxConfig akov2.AtlasSearchIndexConfig
			err := sr.k8sClient.Get(sr.ctx.Context, *spec.Search.SearchConfigurationRef.GetObject(sr.deployment.Namespace), &idxConfig)
			if err != nil {
				err := fmt.Errorf("can not get search index configuration for index '%s'. E: %w", spec.Name, err)
				return sr.terminate(internalState, err)
			}
			stateInAKO = searchindex.NewSearchIndex(spec, &idxConfig.Spec)
		case IndexTypeVector:
			// Vector index doesn't require any external configuration
			stateInAKO = searchindex.NewSearchIndex(spec, &akov2.AtlasSearchIndexConfigSpec{})
		default:
			err := fmt.Errorf("index %q has unknown type %q. Can be either %s or %s",
				spec.Name, spec.Type, IndexTypeSearch, IndexTypeVector)
			return sr.terminate(internalState, err)
		}
	}

	return sr.reconcileInternal(stateInAKO, stateInAtlas)
}

func (sr *searchIndexReconciler) reconcileInternal(stateInAKO, stateInAtlas *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("starting reconciliation for index '%s'", sr.indexName)
	defer sr.ctx.Log.Debugf("finished reconciliation for index '%s'", sr.indexName)

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
	sr.ctx.SetConditionFromResult(api.SearchIndexesReadyType, ok)
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
	sr.ctx.SetConditionFromResult(api.SearchIndexesReadyType, terminate)
	return terminate
}

func (sr *searchIndexReconciler) create(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[creating] index %s", index.Name)
	defer sr.ctx.Log.Debugf("[creation finished] for index %s", index.Name)
	akoIdx, err := sr.searchService.CreateIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), index)
	if err != nil {
		return sr.terminate(index, err)
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
	sr.ctx.SetConditionFromResult(api.SearchIndexesReadyType, inProgress)
	return inProgress
}

func (sr *searchIndexReconciler) delete(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[deleting] index %s", index.Name)
	defer sr.ctx.Log.Debugf("[deletion finished] for index %s", index.Name)
	if index.ID == nil {
		return workflow.OK()
	}

	if err := sr.searchService.DeleteIndex(
		sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), index.GetID()); err != nil {
		return sr.terminate(index, err)
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
	convertedIdx, err := sr.searchService.UpdateIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), akoIdx)
	if err != nil {
		return sr.terminate(atlasIdx, err)
	}
	return sr.progress(convertedIdx)
}
