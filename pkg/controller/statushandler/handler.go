package statushandler

import (
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Update performs the update (in the form of patch) for the Atlas Custom Resource status.
// It should be a common method for all the controllers
func Update(ctx *workflow.Context, kubeClient client.Client, resource mdbv1.AtlasCustomResource) {
	if ctx.LastCondition() != nil {
		ctx.Log.Infow("Status update", "lastCondition", ctx.LastCondition())
	}

	resource.UpdateStatus(ctx.Conditions(), ctx.StatusOptions()...)

	if err := patchUpdateStatus(kubeClient, resource); err != nil {
		ctx.Log.Errorf("Failed to update status: %s", err)
	}
}
