package statushandler

import (
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// Update performs the update (in the form of patch) for the Atlas Custom Resource status.
// It should be a common method for all the controllers
func Update(ctx *workflow.Context, kubeClient client.Client, eventRecorder record.EventRecorder, resource mdbv1.AtlasCustomResource) {
	if ctx.LastCondition() != nil {
		logEvent(ctx, eventRecorder, resource)
	}

	resource.UpdateStatus(ctx.Conditions(), ctx.StatusOptions()...)

	if err := patchUpdateStatus(kubeClient, resource); err != nil {
		if apiErrors.IsNotFound(err) {
			ctx.Log.Infof("The resource %s no longer exists, not updating the status", kube.ObjectKey(resource.GetNamespace(), resource.GetName()))
			return
		}
		// Implementation logic: we deliberately don't return the 'error' to avoid cumbersome handling logic as the
		// failed update of the status is not something that should block reconciliation
		ctx.Log.Errorf("Failed to update status: %s", err)
	}
}

// logEvent logs the last condition to the output and also creates the Event for it in Kubernetes.
// Some tradeoffs about event submission: the Event always requires the 'reason' and 'message' though our Status
// conditions may lack that in case the condition is successful ("true"). In this case we leave the message empty
// and use the ConditionType as the reason.
func logEvent(ctx *workflow.Context, eventRecorder record.EventRecorder, resource mdbv1.AtlasCustomResource) {
	ctx.Log.Infow("Status update", "lastCondition", ctx.LastCondition())

	eventType := "Normal"
	reason := string(ctx.LastCondition().Type)
	msg := ""
	if ctx.LastCondition().Reason != "" {
		reason = ctx.LastCondition().Reason
		msg = ctx.LastCondition().Message
	}
	if ctx.LastConditionWarn() {
		eventType = "Warning"
	}
	eventRecorder.Event(resource, eventType, reason, msg)
}
