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
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func (r *AtlasIPAccessListReconciler) create(
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

	return r.inProgress(ctx, ipAccessList, "Atlas has started to add access list entries")
}

func (r *AtlasIPAccessListReconciler) deleteAll(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	atlasIPAccessList ipaccesslist.IPAccessEntries,
) ctrl.Result {
	err := r.delete(ctx, ipAccessListService, projectID, atlasIPAccessList)
	if err != nil {
		return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListFailedToDelete, err)
	}

	return r.unmanage(ctx, ipAccessList)
}

func (r *AtlasIPAccessListReconciler) deletePartial(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	ipAccessList *akov2.AtlasIPAccessList,
	projectID string,
	atlasIPAccessList ipaccesslist.IPAccessEntries,
) ctrl.Result {
	err := r.delete(ctx, ipAccessListService, projectID, atlasIPAccessList)
	if err != nil {
		return r.terminate(ctx, ipAccessList, api.IPAccessListReady, workflow.IPAccessListFailedToDelete, err)
	}

	return r.inProgress(ctx, ipAccessList, "Atlas has started to delete access list entries")
}

func (r *AtlasIPAccessListReconciler) delete(
	ctx *workflow.Context,
	ipAccessListService ipaccesslist.IPAccessListService,
	projectID string,
	atlasIPAccessList ipaccesslist.IPAccessEntries,
) error {
	for _, entry := range atlasIPAccessList {
		err := ipAccessListService.Delete(ctx.Context, projectID, entry)
		if err != nil {
			return err
		}

		ctx.EnsureStatusOption(status.RemoveIPAccessListEntryStatus(entry.ID()))
	}

	return nil
}

func (r *AtlasIPAccessListReconciler) skip(ctx context.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasIPAccessList reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", ipAccessList.Spec)
	if !ipAccessList.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, ipAccessList, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) invalidate(invalid workflow.Result) ctrl.Result {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasIPAccessList is invalid: %v", invalid)
	return invalid.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) unsupport(ctx *workflow.Context) ctrl.Result {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, errors.New("the AtlasIPAccessList is not supported by Atlas for government")).
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) inProgress(
	ctx *workflow.Context,
	ipAccessList *akov2.AtlasIPAccessList,
	msg string,
) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	result := workflow.InProgress(workflow.IPAccessListPending, msg)
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(api.IPAccessListReady, result)

	return result.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) unmanage(ctx *workflow.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.Deleted().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) ready(ctx *workflow.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	ctx.SetConditionTrue(api.ReadyType).
		SetConditionTrue(api.IPAccessListReady)

	if ipAccessList.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult()
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) terminate(
	ctx *workflow.Context,
	ipAccessList *akov2.AtlasIPAccessList,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	err error,
) ctrl.Result {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", ipAccessList, ipAccessList.GetNamespace(), ipAccessList.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}
