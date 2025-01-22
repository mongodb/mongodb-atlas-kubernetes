package atlasipaccesslist

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func (r *AtlasIPAccessListReconciler) handleIPAccessList(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	projectID string,
	ipAccessList *akov2.AtlasIPAccessList,
) ctrl.Result {
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
			return r.deleteIPAccessList(ctx, ipAccessListService, ipAccessList, projectID, atlasIPAccessList, false)
		}

		r.Log.Info("releasing ip access list resource for deletion")
		return r.unmanage(ctx, ipAccessList)
	}

	if !existInAtlas {
		r.Log.Infof("creating ip access list on project %s", projectID)
		return r.createIPAccessList(ctx, ipAccessListService, ipAccessList, projectID, akoIPAccessList)
	}

	r.Log.Infof("managing ip access list on project %s", projectID)
	return r.manageIPAccessList(ctx, ipAccessListService, ipAccessList, projectID, akoIPAccessList, atlasIPAccessList)
}

func (r *AtlasIPAccessListReconciler) createIPAccessList(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	akoIPAccessList ipaccesslist.IPAccessEntries,
) ctrl.Result {
	err := ipAccessListService.Add(ctx.Context, projectID, akoIPAccessList)
	if err != nil {
		return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListFailedToCreate, err)
	}

	return r.watchState(ctx, ipAccessListService, ipAccessList, projectID, akoIPAccessList)
}

func (r *AtlasIPAccessListReconciler) deleteIPAccessList(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	atlasIPAccessList ipaccesslist.IPAccessEntries,
	partial bool,
) ctrl.Result {
	for _, entry := range atlasIPAccessList {
		err := ipAccessListService.Delete(ctx.Context, projectID, entry)
		if err != nil {
			return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListFailedToDelete, err)
		}

		ctx.EnsureStatusOption(status.RemoveIPAccessListEntryStatus(entry.ID()))
	}

	if partial {
		return r.inProgress(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListPending, "Atlas has started to delete access list entries")
	}

	return r.unmanage(ctx, ipAccessList)
}

func (r *AtlasIPAccessListReconciler) manageIPAccessList(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	akoIPAccessList ipaccesslist.IPAccessEntries,
	atlasIPAccessList ipaccesslist.IPAccessEntries,
) ctrl.Result {
	toAdd := make(ipaccesslist.IPAccessEntries, len(akoIPAccessList))
	for ID, entry := range akoIPAccessList {
		if _, ok := atlasIPAccessList[ID]; !ok {
			r.Log.Debugf("enqueue %s to be added on project %s", entry.ID(), projectID)
			toAdd[ID] = entry
		}
	}

	if len(toAdd) > 0 {
		return r.createIPAccessList(ctx, ipAccessListService, ipAccessList, projectID, toAdd)
	}

	toDelete := make(ipaccesslist.IPAccessEntries, len(atlasIPAccessList))
	for ID, entry := range atlasIPAccessList {
		if _, ok := akoIPAccessList[ID]; !ok {
			r.Log.Debugf("enqueue %s to be deleted on project %s", entry.ID(), projectID)
			toDelete[ID] = entry
		}
	}

	if len(toDelete) > 0 {
		return r.deleteIPAccessList(ctx, ipAccessListService, ipAccessList, projectID, toDelete, true)
	}

	return r.watchState(ctx, ipAccessListService, ipAccessList, projectID, akoIPAccessList)
}

func (r *AtlasIPAccessListReconciler) watchState(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	akoIPAccessList ipaccesslist.IPAccessEntries,
) ctrl.Result {
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
		return r.inProgress(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListPending, "Atlas has started to add access list entries")
	}

	return r.ready(ctx, ipAccessList)
}
