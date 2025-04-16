// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasdeployment

import (
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
)

const (
	IndexStatusActive = "READY"
)

type searchIndexReconcileRequest struct {
	ctx           *workflow.Context
	deployment    *akov2.AtlasDeployment
	k8sClient     client.Client
	projectID     string
	indexName     string
	searchService searchindex.AtlasSearchIdxService
}

func (sr *searchIndexReconcileRequest) Handle(spec *akov2.SearchIndex, previous *status.DeploymentSearchIndexStatus) workflow.Result {
	var stateInAtlas, stateInAKO *searchindex.SearchIndex
	name := ""

	if previous != nil {
		name = previous.Name
		var err error
		sr.ctx.Log.Debugf("restoring index %q from status", previous.Name)
		stateInAtlas, err = sr.searchService.GetIndex(
			sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), previous.Name, previous.ID)
		if err != nil {
			if !errors.Is(err, searchindex.ErrNotFound) {
				return sr.terminate(stateInAtlas, err)
			}
			stateInAtlas = nil // Not Found = not in Atlas
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

	return sr.reconcileInternal(name, stateInAKO, stateInAtlas)
}

func (sr *searchIndexReconcileRequest) reconcileInternal(indexName string, stateInAKO, stateInAtlas *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("starting reconciliation for index '%s'", sr.indexName)
	defer sr.ctx.Log.Debugf("finished reconciliation for index '%s'", sr.indexName)

	inAtlas := stateInAtlas != nil
	inSpec := stateInAKO != nil

	var currentStatus string
	if stateInAtlas != nil && stateInAtlas.Status != nil {
		currentStatus = *stateInAtlas.Status
	}
	// Atlas is still processing the index, nothing can't be done
	if currentStatus != IndexStatusActive && currentStatus != "" {
		return sr.progress(stateInAtlas)
	}

	switch {
	case !inAtlas && inSpec:
		return sr.create(stateInAKO)
	case inAtlas && inSpec:
		return sr.compare(stateInAKO, stateInAtlas)
	case inAtlas && !inSpec:
		return sr.delete(stateInAtlas)
	default: // not in Atlas or K8s, only if removed from Atlas from elsewhere
		return sr.deleted(indexName)
	}
}

func (sr *searchIndexReconcileRequest) idle(index *searchindex.SearchIndex) workflow.Result {
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
func (sr *searchIndexReconcileRequest) terminate(index *searchindex.SearchIndex, err error) workflow.Result {
	msg := fmt.Errorf("error with processing index %s. err: %w", index.Name, err)
	sr.ctx.Log.Debug(msg)
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentSetSearchIndexStatus(status.NewDeploymentSearchIndexStatus(
		status.SearchIndexStatusError,
		status.WithMsg(msg.Error()),
		status.WithName(index.Name),
	)))
	terminate := workflow.Terminate(status.SearchIndexStatusError, msg)
	sr.ctx.SetConditionFromResult(api.SearchIndexesReadyType, terminate)
	return terminate
}

func (sr *searchIndexReconcileRequest) create(index *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[creating] index %s", index.Name)
	defer sr.ctx.Log.Debugf("[creation finished] for index %s", index.Name)
	akoIdx, err := sr.searchService.CreateIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), index)
	if err != nil {
		return sr.terminate(index, err)
	}
	return sr.progress(akoIdx)
}

func (sr *searchIndexReconcileRequest) progress(index *searchindex.SearchIndex) workflow.Result {
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

func (sr *searchIndexReconcileRequest) delete(index *searchindex.SearchIndex) workflow.Result {
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

func (sr *searchIndexReconcileRequest) deleted(indexName string) workflow.Result {
	sr.ctx.EnsureStatusOption(status.AtlasDeploymentUnsetSearchIndexStatus(status.NewDeploymentSearchIndexStatus("",
		status.WithName(indexName))))
	return workflow.Deleted()
}

func (sr *searchIndexReconcileRequest) compare(akoIdx, atlasIdx *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("[syncing] index %s", akoIdx.Name)
	defer sr.ctx.Log.Debugf("[update finished] for index %s", akoIdx.Name)

	isEqual, err := akoIdx.EqualTo(atlasIdx)
	if err != nil {
		sr.terminate(atlasIdx, err)
	}
	if isEqual {
		sr.ctx.Log.Debugf("index %s is already up to date", akoIdx.Name)
		return sr.idle(atlasIdx)
	}
	return sr.update(akoIdx, atlasIdx)
}

func (sr *searchIndexReconcileRequest) update(akoIdx, atlasIdx *searchindex.SearchIndex) workflow.Result {
	sr.ctx.Log.Debugf("updating index %s...", akoIdx.Name)
	convertedIdx, err := sr.searchService.UpdateIndex(sr.ctx.Context, sr.projectID, sr.deployment.GetDeploymentName(), akoIdx)
	if err != nil {
		return sr.terminate(atlasIdx, err)
	}
	return sr.progress(convertedIdx)
}
