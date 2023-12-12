package atlasproject

import (
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *mdbv1.AtlasProject) workflow.Result {
// Reference set

	// check reference points to existing compliance policy CR

	// check if compliance policy exists in atlas (and matches)
	// if match, return workflow.OK()
	// otherwise, create/update compliance policy...

	// check existing backups meet requirements
		// TODO POTENTIAL RACE WITH DEPLOYMENT CONTROLLER
	// if dont meet, set status, return workflow.Terminate()
	// otherwise, continue...

	// enable finalizer on compliance policy CR (if doesn't already exist)
	// finalizer blocks deletion while there are references and/or compliance policy exists in atlas

	// add annotation to compliance policy for associated atlas project

	// create compliance policy in atlas

// Reference unset

	// Check finalizer
	// Check if compliance policy CR is referenced in k8s
	// Check if compliance policy is in use in Atlas itself

	return workflow.OK()
}
