//Copyright 2025 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package atlasipaccesslist

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func (r *AtlasIPAccessListReconciler) handleCustomResource(ctx context.Context, ipAccessList *akov2.AtlasIPAccessList) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(ipAccessList) {
		return r.skip(ctx, ipAccessList)
	}

	conditions := api.InitCondition(ipAccessList, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, ipAccessList)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, ipAccessList)

	isValid := customresource.ValidateResourceVersion(workflowCtx, ipAccessList, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(ipAccessList) {
		return r.unsupport(workflowCtx)
	}

	connectionConfig, err := r.ResolveConnectionConfig(ctx, ipAccessList)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	atlasProject, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20250312011, ipAccessList)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	ipAccessListService := ipaccesslist.NewIPAccessList(sdkClientSet.SdkClient20250312011.ProjectIPAccessListApi)

	return r.handleIPAccessList(workflowCtx, ipAccessListService, atlasProject.ID, ipAccessList)
}

func (r *AtlasIPAccessListReconciler) handleIPAccessList(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	projectID string,
	ipAccessList *akov2.AtlasIPAccessList,
) (ctrl.Result, error) {
	akoIPAccessList, err := ipaccesslist.NewIPAccessListEntries(ipAccessList)
	if err != nil {
		return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.Internal, err)
	}

	atlasIPAccessList, err := ipAccessListService.List(ctx.Context, projectID)
	if err != nil {
		return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.Internal, err)
	}

	existInAtlas := len(atlasIPAccessList) > 0

	if !ipAccessList.GetDeletionTimestamp().IsZero() {
		if existInAtlas {
			r.Log.Infof("deleting ip access list from project %s", projectID)
			return r.deleteAll(ctx, ipAccessListService, ipAccessList, projectID, atlasIPAccessList)
		}

		r.Log.Info("releasing ip access list resource for deletion")
		return r.unmanage(ctx, ipAccessList)
	}

	if toAdd := collection.MapDiff(akoIPAccessList, atlasIPAccessList); len(toAdd) > 0 {
		r.Log.Infof("adding ip access list %v on project %s", toAdd, projectID)
		return r.create(ctx, ipAccessListService, ipAccessList, projectID, toAdd)
	}

	if toDelete := collection.MapDiff(atlasIPAccessList, akoIPAccessList); len(toDelete) > 0 {
		r.Log.Infof("deleting ip access list %v from project %s", toDelete, projectID)
		return r.deletePartial(ctx, ipAccessListService, ipAccessList, projectID, toDelete)
	}

	pending := false
	for _, entry := range akoIPAccessList {
		r.Log.Debugf("retrieving status of ip access list entry %s at project %s", entry.ID(), projectID)
		entryStatus, err := ipAccessListService.Status(ctx.Context, projectID, entry)
		if err != nil {
			return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListFailedToGetState, err)
		}

		r.Log.Debugf("ip access list entry %s status: %s", entry.ID(), entryStatus)
		if entryStatus != "ACTIVE" {
			pending = true
		}

		ctx.EnsureStatusOption(status.AddIPAccessListEntryStatus(entry.ID(), entryStatus))
	}

	if pending {
		return r.inProgress(ctx, ipAccessList, "Atlas has started to add access list entries")
	}

	return r.ready(ctx, ipAccessList)
}
