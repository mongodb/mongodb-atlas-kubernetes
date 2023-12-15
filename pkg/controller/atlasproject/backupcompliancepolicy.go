package atlasproject

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *mdbv1.AtlasProject) workflow.Result {
	// Reference set
	if project.Spec.BackupCompliancePolicyRef.Name != "" {
		// TODO start watching backup compliance CR
		// check reference points to existing compliance policy CR
		compliancePolicy := &mdbv1.AtlasBackupCompliancePolicy{}
		err := r.Client.Get(context.Background(), *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace), compliancePolicy)
		if err != nil {
			ctx.Log.Errorf("failed to get backup compliance policy ")
		}
		// check if compliance policy exists in atlas (and matches)
		// if match, return workflow.OK()
		atlasCompliancePolicy, _, err := ctx.Client.BackupCompliancePolicy.Get(context.Background(), project.ID())
		if err != nil {
			ctx.Log.Errorf("failed to get backup compliance policy from atlas: %v", err)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
		}

		// otherwise, create/update compliance policy...

		// check existing backups meet requirements
		// TODO POTENTIAL RACE WITH DEPLOYMENT CONTROLLER
		// if dont meet, set status, return workflow.Terminate()

		if !currentBackupPoliciesMatchCompliance() {
			// TODO figure out appropriate status names/messages
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyNotMet, "current backup policies do not satisfy this compliance policy")
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
	// Reference unset
	// TODO should this be an else or can we just let the flow fall through?

	// Check finalizer
	// Check if compliance policy CR is referenced in k8s
	// Check if compliance policy is in use in Atlas itself

	return workflow.OK()
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
func currentBackupPoliciesMatchCompliance() bool {
	// TODO should this return a bool, or an error?
	return true
}
