// Copyright 2024 MongoDB Inc
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

package atlasproject

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

type backupComplianceController struct {
	ctx     *workflow.Context
	client  client.Client
	project *akov2.AtlasProject
}

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	ctx.Log.Debug("starting backup compliance policy processing")
	defer ctx.Log.Debug("finished backup compliance policy processing")

	b := backupComplianceController{
		ctx:     ctx,
		client:  r.Client,
		project: project,
	}

	c, ok := ctx.GetCondition(api.BackupComplianceReadyType)
	if ok {
		switch reason := workflow.ConditionReason(c.Reason); reason {
		case workflow.ProjectBackupCompliancePolicyUpdating:
			return b.handleUpserting()
		}
	}

	return b.handlePending()
}

func (b *backupComplianceController) handlePending() workflow.DeprecatedResult {
	bcp, inAtlas, err := b.getAtlasBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}

	inAKO := false
	if bcpRef := b.project.Spec.BackupCompliancePolicyRef; bcpRef != nil {
		inAKO = !bcpRef.IsEmpty()
	}

	switch {
	case inAKO && !inAtlas:
		return b.upsert(bcp)
	case inAKO && inAtlas:
		return b.upsert(bcp)
	case !inAKO && inAtlas:
		return b.delete()
	default:
		return b.unmanage()
	}
}

func (b *backupComplianceController) handleUpserting() workflow.DeprecatedResult {
	atlasBCP, found, err := b.getAtlasBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, err)
	}
	if !found {
		return b.terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, errors.New("bcp not found in Atlas"))
	}
	akoBCP, err := b.getAKOBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}
	equal, err := cmp.SemanticEqual(&akoBCP.Spec, akov2.NewBCPFromAtlas(atlasBCP))
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}
	if equal {
		lastApplied, ok := akoBCP.GetAnnotations()[customresource.AnnotationLastAppliedConfiguration]
		if ok {
			temp := &akov2.AtlasBackupCompliancePolicy{}
			err = json.Unmarshal([]byte(lastApplied), temp)
			if err != nil {
				return b.terminate(workflow.Internal, err)
			}
			equal = (akoBCP.Spec.OverwriteBackupPolicies == temp.Spec.OverwriteBackupPolicies)
		}
	}

	switch {
	case !equal:
		return b.terminate(workflow.ProjectBackupCompliancePolicyUpdating, errors.New("aborting update: spec has changed"))
	case atlasBCP.GetState() != "ACTIVE":
		return b.progress(
			workflow.ProjectBackupCompliancePolicyUpdating,
			fmt.Sprintf("backup compliance policy not ready yet: %q", atlasBCP.GetState()),
			"updating backup compliance policy",
		)
	default:
		return b.idle()
	}
}

// upsert updates the backup compliance settings for a project. These settings can only be updated (not created), so
// we also use this for creation too.
func (b *backupComplianceController) upsert(atlasBCP *admin.DataProtectionSettings20231001) workflow.DeprecatedResult {
	b.ctx.Log.Debug("updating backup compliance policy")
	akoBCP, err := b.getAKOBackupCompliancePolicy()
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}

	equal, err := cmp.SemanticEqual(akoBCP.Spec.DeepCopy(), akov2.NewBCPFromAtlas(atlasBCP))
	if err != nil {
		return b.terminate(workflow.Internal, err)
	}
	if equal {
		lastApplied, ok := akoBCP.GetAnnotations()[customresource.AnnotationLastAppliedConfiguration]
		if ok {
			temp := &akov2.AtlasBackupCompliancePolicy{}
			err = json.Unmarshal([]byte(lastApplied), temp)
			if err != nil {
				return b.terminate(workflow.Internal, err)
			}
			equal = (akoBCP.Spec.OverwriteBackupPolicies == temp.Spec.OverwriteBackupPolicies)
		}
	}

	if !equal {
		atlasBCP, _, err = b.ctx.SdkClientSet.SdkClient20250312009.CloudBackupsApi.UpdateCompliancePolicy(b.ctx.Context, b.project.ID(), akoBCP.ToAtlas(b.project.ID())).OverwriteBackupPolicies(akoBCP.Spec.OverwriteBackupPolicies).Execute()
		if err != nil {
			if admin.IsErrorCode(err, atlas.BackupComplianceNotMet) {
				return b.terminate(workflow.ProjectBackupCompliancePolicyNotMet, err)
			}
			return b.terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, err)
		}
	}

	if atlasBCP.GetState() != "ACTIVE" {
		return b.progress(
			workflow.ProjectBackupCompliancePolicyUpdating,
			fmt.Sprintf("backup compliance policy not ready yet: %q", atlasBCP.GetState()),
			"updating backup compliance policy",
		)
	}

	return b.idle()
}

// delete begins the deletion process for a backup compliance policy. However, the admin API cannot delete BCPs (which must be
// done via support), so we notify the user of next steps via status.
func (b *backupComplianceController) delete() workflow.DeprecatedResult {
	b.ctx.Log.Debug("deleting backup compliance policy")
	return b.terminate(
		workflow.ProjectBackupCompliancePolicyCannotDelete,
		errors.New("deletion of backup compliance policies is not supported in Atlas: please contact support"),
	).WithoutRetry()
}

func (b *backupComplianceController) progress(state workflow.ConditionReason, fineMsg, coarseMsg string) workflow.DeprecatedResult {
	var (
		fineProgress   = workflow.InProgress(state, fineMsg)
		coarseProgress = workflow.InProgress(state, coarseMsg)
	)

	b.ctx.SetConditionFromResult(api.BackupComplianceReadyType, fineProgress)
	return coarseProgress
}

// terminate transitions to pending state if an error occurred.
func (b *backupComplianceController) terminate(reason workflow.ConditionReason, err error) workflow.DeprecatedResult {
	b.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	b.ctx.SetConditionFromResult(api.BackupComplianceReadyType, result)
	return result
}

// unmanage transitions to pending state if there is no managed BCP.
func (b *backupComplianceController) unmanage() workflow.DeprecatedResult {
	b.ctx.UnsetCondition(api.BackupComplianceReadyType)
	return workflow.OK()
}

// idle transitions BCP to idle state when ready and idle.
func (b *backupComplianceController) idle() workflow.DeprecatedResult {
	b.ctx.SetConditionTrue(api.BackupComplianceReadyType)
	return workflow.OK()
}

func (b *backupComplianceController) getAtlasBackupCompliancePolicy() (*admin.DataProtectionSettings20231001, bool, error) {
	bcp, _, err := b.ctx.SdkClientSet.SdkClient20250312009.CloudBackupsApi.GetCompliancePolicy(b.ctx.Context, b.project.ID()).Execute()
	if err != nil {
		// NOTE: getting backup compliance policies never yields a 404
		return nil, false, fmt.Errorf("error finding backup compliance policy: %w", err)
	}

	var emptyPolicy admin.DataProtectionSettings20231001
	if bcp == nil || reflect.DeepEqual(&emptyPolicy, bcp) {
		return nil, false, nil
	}

	return bcp, true, nil
}

func (b *backupComplianceController) getAKOBackupCompliancePolicy() (*akov2.AtlasBackupCompliancePolicy, error) {
	if b.project.Spec.BackupCompliancePolicyRef == nil {
		return nil, errors.New("bcp not found in Kubernetes")
	}
	bcp := &akov2.AtlasBackupCompliancePolicy{}
	err := b.client.Get(b.ctx.Context, *b.project.Spec.BackupCompliancePolicyRef.GetObject(b.project.Namespace), bcp)
	if err != nil {
		return nil, err
	}
	return bcp, nil
}
