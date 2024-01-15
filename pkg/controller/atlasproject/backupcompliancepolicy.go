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
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *mdbv1.AtlasProject) workflow.Result {
	defer func() { r.garbageCollectBackupResource(ctx.Context, project.GetName()) }()

	if IsBackupComplianceEmpty(project.Spec.BackupCompliancePolicyRef) {
		// check if it is actually enabled in Atlas
		// if it is enabled in Atlas, we still have to signal here via the status condition
		// that there is an not-deleted-yet backup compliance policy in Atlas
		ctx.UnsetCondition(status.BackupComplianceReadyType)
		return workflow.OK()
	}

	// reference set
	// TODO start watching backup compliance CR
	// check reference points to existing compliance policy CR
	compliancePolicy := &mdbv1.AtlasBackupCompliancePolicy{}
	err := r.Client.Get(ctx.Context, *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace), compliancePolicy)
	if err != nil {
		ctx.Log.Errorf("failed to get backup compliance policy: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}
	// check if compliance policy exists in atlas (and matches)
	// if match, return workflow.OK()
	atlasCompliancePolicy, _, err := ctx.Client.BackupCompliancePolicy.Get(ctx.Context, project.ID())
	if err != nil {
		ctx.Log.Errorf("failed to get backup compliance policy from atlas: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}

	// otherwise, create/update compliance policy...

	// check existing backups meet requirements
	// TODO POTENTIAL RACE WITH DEPLOYMENT CONTROLLER
	// if dont meet, set status, return workflow.Terminate()
	backups, err := r.getBackupPoliciesInProject(ctx.Context, project)
	if err != nil {
		ctx.Log.Errorf("failed to get backup policies: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}

	if err = currentBackupPoliciesMatchCompliance(backups, *compliancePolicy); err != nil {
		// TODO figure out appropriate status names/messages
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyNotMet, err.Error())
	}
	// otherwise, continue...

	// enable finalizer on compliance policy CR (if doesn't already exist)
	err = customresource.ManageFinalizer(context.Background(), r.Client, compliancePolicy, customresource.SetFinalizer)
	if err != nil {
		ctx.SetConditionFalse(status.BackupComplianceReadyType)
		return workflow.Terminate(workflow.AtlasFinalizerNotSet, err.Error())
	}
	// finalizer blocks deletion while there are references and/or compliance policy exists in atlas

	// add annotation to compliance policy for associated atlas project
	compliancePolicy.SetAnnotations(map[string]string{
		// TODO pick better name for project annotation
		"mongodbatlas/project": project.ID(),
	})
	r.Client.Update(context.Background(), compliancePolicy)

	// create compliance policy in atlas
	result := syncBackupCompliancePolicy(ctx, project.ID(), *compliancePolicy, *atlasCompliancePolicy)
	if !result.IsOk() {
		ctx.SetConditionFromResult(status.BackupComplianceReadyType, result)
		return result
	}

	return workflow.OK()
}

// TODO: there is certainly a better way of doing this
// can we annotate these seperate resources to attribute them to projects/deployments?
func (r *AtlasProjectReconciler) getBackupPoliciesInProject(ctx context.Context, project *mdbv1.AtlasProject) ([]mdbv1.AtlasBackupPolicyItem, error) {
	policies := []mdbv1.AtlasBackupPolicyItem{}
	deployments := &mdbv1.AtlasDeploymentList{}

	// Get all deployments
	err := r.Client.List(ctx, deployments)
	if err != nil {
		return policies, fmt.Errorf("failed to retrieve list of deployments: %w", err)
	}
	// We only want deployments in this project
	for _, d := range deployments.Items {
		if d.Spec.Project.Name != project.Name || d.Spec.Project.Namespace != project.Namespace {
			continue
		}
		// Get backup schedule for deployment
		schedule := &mdbv1.AtlasBackupSchedule{}
		err = r.Client.Get(ctx, *d.Spec.BackupScheduleRef.GetObject(d.Namespace), schedule)
		if err != nil {
			return policies, fmt.Errorf("failed to retrieve backup schedule: %w", err)
		}
		// Get backup policy from schedule
		policy := &mdbv1.AtlasBackupPolicy{}
		err = r.Client.Get(ctx, *schedule.Spec.PolicyRef.GetObject(d.Namespace), policy)
		if err != nil {
			return policies, fmt.Errorf("failed to retreieve backup policy: %w", err)
		}
		// apparently there is only ever 1 item, but better safe than sorry
		policies = append(policies, policy.Spec.Items...)
	}
	return policies, nil
}

// syncBackupCompliancePolicy compares the compliance policy specified in Kubernetes to the one currently present in Atlas, updating Atlas should they differ.
func syncBackupCompliancePolicy(ctx *workflow.Context, groupID string, kubeCompliancePolicy mdbv1.AtlasBackupCompliancePolicy, atlasCompliancePolicy mongodbatlas.BackupCompliancePolicy) workflow.Result {
	// TODO do we need to check this, or can we just always update?
	localCompliancePolicy := kubeCompliancePolicy.ToAtlas()
	if cmp.Diff(localCompliancePolicy, atlasCompliancePolicy, cmpopts.EquateEmpty()) != "" {
		_, _, err := ctx.Client.BackupCompliancePolicy.Update(context.Background(), groupID, localCompliancePolicy)
		if err != nil {
			ctx.Log.Errorf("failed to update backup compliance policy in atlas: %v", err)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, err.Error())
		}
	}

	return workflow.OK()
}

// currentBackupPoliciesMatchCompliance checks all backup policies present in the project, assessing if they meet the requirements specified in the backup compliance policy.
func currentBackupPoliciesMatchCompliance(backups []mdbv1.AtlasBackupPolicyItem, compliance mdbv1.AtlasBackupCompliancePolicy) error {
	// error rather than bool means we can accumulate errors and report all insufficient backup policies, rather than just the first we encounter
	var err error
	for _, complianceScheduledPolicy := range compliance.Spec.ScheduledPolicyItems {
		for _, backup := range backups {
			if backup.FrequencyType != complianceScheduledPolicy.FrequencyType {
				continue
			}
			if !compareBackupPolicyItem(backup, complianceScheduledPolicy) {
				// TODO: ideally have some identifying information here, but currently no way to tell backup policies apart (pass in policy rather than policyitem?)
				errors.Join(err, errors.New("existing backup policy does not satisfy backup compliance policy"))
			}
		}
	}
	return err
}

// TODO: likely needs renaming because we're actually checking >= rather than ==
// compareBackupPolicyItem checks that policy item x satisfies the minimums set in y
func compareBackupPolicyItem(x, y mdbv1.AtlasBackupPolicyItem) bool {
	return x.FrequencyType == y.FrequencyType &&
		x.FrequencyInterval >= y.FrequencyInterval &&
		normalizeRetention(x) >= normalizeRetention(y)
}

// normalizeRetention 'normalizes' the retention, which is otherwise defined by both its value and its units.
func normalizeRetention(policy mdbv1.AtlasBackupPolicyItem) int {
	switch policy.RetentionUnit {
	case "days":
		return policy.RetentionValue
	case "weeks":
		return policy.RetentionValue * 7
	case "months":
		return policy.RetentionValue * 31
	}
	return -1
}

func IsBackupComplianceEmpty(backupCompliancePolicyRef *common.ResourceRefNamespaced) bool {
	return (backupCompliancePolicyRef == nil) || (backupCompliancePolicyRef.Name == "")
}
