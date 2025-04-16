// Copyright 2025 MongoDB Inc
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

package statushandler

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

// Update performs the update (in the form of patch) for the Atlas Custom Resource status.
// It should be a common method for all the controllers
func Update(ctx *workflow.Context, kubeClient client.Client, eventRecorder record.EventRecorder, resource api.AtlasCustomResource) {
	if ctx.LastCondition() != nil {
		logEvent(ctx, eventRecorder, resource)
	}

	if err := patchUpdateStatus(ctx, kubeClient, resource); err != nil {
		if apierrors.IsNotFound(err) {
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
func logEvent(ctx *workflow.Context, eventRecorder record.EventRecorder, resource api.AtlasCustomResource) {
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
