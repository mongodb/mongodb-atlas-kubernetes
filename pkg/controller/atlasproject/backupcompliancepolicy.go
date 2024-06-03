/*
Copyright 2023 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasproject

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const bcpNotMet = "BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS"

type backupComplianceController struct {
	ctx     *workflow.Context
	client  client.Client
	project *akov2.AtlasProject
}

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	ctx.Log.Debug("starting backup compliance policy processing")
	defer ctx.Log.Debug("finished backup compliance policy processing")

	b := backupComplianceController{
		ctx:     ctx,
		client:  r.Client,
		project: project,
	}

	return b.handlePending()
}

func (b *backupComplianceController) handlePending() workflow.Result {
	bcp, found, err := b.getAtlasBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}

	akoEmpty := true
	if bcpRef := b.project.Spec.BackupCompliancePolicyRef; bcpRef != nil {
		akoEmpty = bcpRef.IsEmpty()
	}

	atlasEmpty := !found

	switch {
	case !akoEmpty && atlasEmpty:
		return b.upsert(bcp)
	case !akoEmpty && !atlasEmpty:
		return b.upsert(bcp)
	case akoEmpty && !atlasEmpty:
		return b.delete()
	default:
		return b.unmanage()
	}
}

// upsert updates the backup compliance settings for a project. These settings can only be updated (not created), so
// we also use this for creation too.
func (b *backupComplianceController) upsert(atlasBCP *admin.DataProtectionSettings20231001) workflow.Result {
	b.ctx.Log.Debug("updating backup compliance policy")
	akoBCP, err := b.getAKOBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}
	equal, err := cmp.SemanticEqual(akoBCP, akov2.NewBCPFromAtlas(atlasBCP))
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}
	if !equal {
		_, _, err = b.ctx.SdkClient.CloudBackupsApi.UpdateDataProtectionSettings(b.ctx.Context, b.project.ID(), akoBCP.ToAtlas(b.project.ID())).OverwriteBackupPolicies(akoBCP.Spec.OverwriteBackupPolicies).Execute()
		if err != nil {
			if admin.IsErrorCode(err, bcpNotMet) {
				return b.terminate(workflow.ProjectBackupCompliancePolicyNotMet, err)
			}
			return b.terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, err)
		}
	}
	return b.idle()
}

// delete begins the deletion process for a backup compliance policy. However, the admin API cannot delete BCPs (which must be
// done via support), so we notify the user of next steps via status.
func (b *backupComplianceController) delete() workflow.Result {
	b.ctx.Log.Debug("deleting backup compliance policy")
	return b.terminate(
		workflow.ProjectBackupCompliancePolicyCannotDelete,
		errors.New("cannot delete backup compliance policy through the API - to delete a backup compliance policy, please contact support"),
	).WithoutRetry()
}

// terminate transitions to pending state if an error occurred.
func (b *backupComplianceController) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	b.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err.Error())
	b.ctx.SetConditionFromResult(api.BackupComplianceReadyType, result)
	return result
}

// unmanage transitions to pending state if there is no managed BCP.
func (b *backupComplianceController) unmanage() workflow.Result {
	b.ctx.UnsetCondition(api.BackupComplianceReadyType)
	return workflow.OK()
}

// idle transitions BCP to idle state when ready and idle.
func (b *backupComplianceController) idle() workflow.Result {
	b.ctx.SetConditionTrue(api.BackupComplianceReadyType)
	return workflow.OK()
}

func (b *backupComplianceController) getAtlasBackupCompliancePolicy() (*admin.DataProtectionSettings20231001, bool, error) {
	bcp, _, err := b.ctx.SdkClient.CloudBackupsApi.GetDataProtectionSettings(b.ctx.Context, b.project.ID()).Execute()
	if err != nil {
		// Note: getting backup compliance policies never yields a 404
		return nil, false, fmt.Errorf("error finding backup compliance policy: %w", err)
	}

	var emptyPolicy admin.DataProtectionSettings20231001
	if bcp == nil || reflect.DeepEqual(&emptyPolicy, bcp) {
		return nil, false, nil
	}

	return bcp, true, nil
}

func (b *backupComplianceController) getAKOBackupCompliancePolicy() (*akov2.AtlasBackupCompliancePolicy, error) {
	bcp := &akov2.AtlasBackupCompliancePolicy{}
	err := b.client.Get(b.ctx.Context, *b.project.BackupCompliancePolicyObjectKey(), bcp)
	if err != nil {
		return nil, err
	}
	return bcp, nil
}
