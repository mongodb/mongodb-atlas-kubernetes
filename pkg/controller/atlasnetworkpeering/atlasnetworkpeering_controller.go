// Package atlasnetworkpeering holds the network peering controller
package atlasnetworkpeering

/*
Copyright 2024.

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

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	typeName = "AtlasNetworkPeering"
)

// AtlasNetworkPeeringReconciler reconciles a AtlasNetworkPeering object
type AtlasNetworkPeeringReconciler struct {
	reconciler.Reconciler
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

type reconcileRequest struct {
	workflowCtx    *workflow.Context
	service        networkpeering.NetworkPeeringService
	projectID      string
	networkPeering *akov2.AtlasNetworkPeering
}

//+kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AtlasNetworkPeering object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AtlasNetworkPeeringReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infow("-> Starting AtlasNetworkPeering reconciliation")

	akoNetworkPeering := akov2.AtlasNetworkPeering{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoNetworkPeering, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), errors.New(result.GetMessage())
	}

	if customresource.ReconciliationShouldBeSkipped(&akoNetworkPeering) {
		return r.Skip(ctx, typeName, &akoNetworkPeering, &akoNetworkPeering.Spec)
	}

	conditions := api.InitCondition(&akoNetworkPeering, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, &akoNetworkPeering)

	isValid := customresource.ValidateResourceVersion(workflowCtx, &akoNetworkPeering, r.Log)
	if !isValid.IsOk() {
		return r.Invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(&akoNetworkPeering) {
		return r.Unsupport(workflowCtx, typeName)
	}

	projectRefs := &akoNetworkPeering.Spec.ProjectReferences
	credentials, err := r.SelectCredentials(ctx, projectRefs, &akoNetworkPeering)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClient, _, err := r.AtlasProvider.SdkClient(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	projectID, err := r.GetProjectID(ctx, projectRefs, sdkClient, akoNetworkPeering.Namespace)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	return r.handle(&reconcileRequest{
		workflowCtx:    workflowCtx,
		service:        networkpeering.NewNetworkPeeringService(sdkClient.NetworkPeeringApi),
		projectID:      projectID,
		networkPeering: &akoNetworkPeering,
	})
}

func (r *AtlasNetworkPeeringReconciler) handle(req *reconcileRequest) (ctrl.Result, error) {
	r.Log.Infow("handling network peering reconcile request",
		"service set", (req.service != nil), "projectID", req.projectID, "networkPeering", req.networkPeering)
	// TODO: state machine goes here
	return ctrl.Result{}, nil
}

func (r *AtlasNetworkPeeringReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	err error,
) (ctrl.Result, error) {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err.Error())
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	return result.ReconcileResult(), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasNetworkPeeringReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&akov2.AtlasNetworkPeering{}).
		Complete(r)
}
