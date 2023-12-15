package atlasproject

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *mdbv1.AtlasProject) workflow.Result {
	// Reference set
	if project.Spec.BackupCompliancePolicyRef.Name != "" {
		// check reference points to existing compliance policy CR
		compliancePolicy := &mdbv1.AtlasBackupCompliancePolicy{}
		r.Client.Get(context.Background(), *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace), compliancePolicy)
		// check if compliance policy exists in atlas (and matches)
		atlasCompliancePolicy, _, err := ctx.Client.BackupCompliancePolicy.Get(context.Background(), project.ID())
		if err != nil {
			// TODO Does this error or return empty when no compliance policy set?
			ctx.SetConditionFalse(status.BackupComplianceReadyType)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
		}
		// if match, return workflow.OK()
		if compliancePolicyMatch(compliancePolicy, atlasCompliancePolicy) {
			return workflow.OK()
		}
		// otherwise, create/update compliance policy...

		// check existing backups meet requirements
		// TODO POTENTIAL RACE WITH DEPLOYMENT CONTROLLER
		// if dont meet, set status, return workflow.Terminate()
		if !currentBackupPoliciesMatchCompliance() {
			// TODO figure out appropriate status names
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
		_, _, err = ctx.Client.BackupCompliancePolicy.Update(context.Background(), project.ID(), compliancePolicy.ToAtlas())
		if err != nil {
			ctx.SetConditionFalse(status.BackupComplianceReadyType)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
		}
		return workflow.OK()
	}
	// Reference unset

	// Check finalizer
	// Check if compliance policy CR is referenced in k8s
	// Check if compliance policy is in use in Atlas itself

	return workflow.OK()
}

func compliancePolicyMatch(k8s *mdbv1.AtlasBackupCompliancePolicy, atlas *mongodbatlas.BackupCompliancePolicy) bool {
	// TODO diff k8s object & atlas object
	return true
}

func currentBackupPoliciesMatchCompliance() bool {
	// TODO check current backup policies in project against compliance policy
	// TODO should this return a bool, or an error?
	return true
}
