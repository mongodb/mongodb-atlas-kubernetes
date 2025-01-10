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
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

// AtlasNetworkPeeringReconciler reconciles a AtlasNetworkPeering object
type AtlasNetworkPeeringReconciler struct {
	reconciler.AtlasReconciler
	AtlasProvider            atlas.Provider
	Scheme                   *runtime.Scheme
	EventRecorder            record.EventRecorder
	GlobalPredicates         []predicate.Predicate
	ObjectDeletionProtection bool
}

func NewAtlasNetworkPeeringsReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasNetworkPeeringReconciler {
	return &AtlasNetworkPeeringReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: mgr.GetClient(),
			Log:    logger.Named("controllers").Named("AtlasNetworkPeering").Sugar(),
		},
		Scheme:                   mgr.GetScheme(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasPrivateEndpoint"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
	}
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

	typeName := reflect.TypeOf(akoNetworkPeering).Name()
	if customresource.ReconciliationShouldBeSkipped(&akoNetworkPeering) {
		return r.Skip(ctx, typeName, &akoNetworkPeering, &akoNetworkPeering.Spec)
	}

	conditions := api.InitCondition(&akoNetworkPeering, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, &akoNetworkPeering)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, &akoNetworkPeering)

	isValid := customresource.ValidateResourceVersion(workflowCtx, &akoNetworkPeering, r.Log)
	if !isValid.IsOk() {
		return r.Invalidate(typeName, isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(&akoNetworkPeering) {
		return r.Unsupport(workflowCtx, typeName)
	}

	credentials, err := r.ResolveCredentials(ctx, &akoNetworkPeering)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClient, orgID, err := r.AtlasProvider.SdkClient(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err)
	}
	project, err := r.ResolveProject(ctx, sdkClient, &akoNetworkPeering, orgID)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err)
	}
	return r.handle(&reconcileRequest{
		workflowCtx:    workflowCtx,
		service:        networkpeering.NewNetworkPeeringService(sdkClient.NetworkPeeringApi),
		projectID:      project.ID,
		networkPeering: &akoNetworkPeering,
	})
}

func (r *AtlasNetworkPeeringReconciler) handle(req *reconcileRequest) (ctrl.Result, error) {
	r.Log.Infow("handling network peering reconcile request",
		"service set", (req.service != nil), "projectID", req.projectID, "networkPeering", req.networkPeering)
	atlasPeer, err := req.service.GetPeer(req.workflowCtx.Context, req.projectID, req.networkPeering.Spec.ContainerID)
	if err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	inAtlas := atlasPeer != nil
	deleted := req.networkPeering.DeletionTimestamp != nil

	switch {
	case !deleted && !inAtlas:
		return r.create(req)
	case !deleted && inAtlas:
		return r.sync(req)
	case deleted && inAtlas:
		return r.delete(req)
	default: // deleted && !inAtlas so nothing to do
		return r.noop(req)
	}
}

func (r *AtlasNetworkPeeringReconciler) create(req *reconcileRequest) (ctrl.Result, error) {
	spec := req.networkPeering.Spec
	ctx := req.workflowCtx.Context
	containerID := containerID(req.networkPeering)
	if containerID == "" {
		requestedContainer := &networkpeering.ProviderContainer{
			AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
				AtlasCIDRBlock: spec.AtlasCIDRBlock,
			},
			Provider: spec.Provider,
		}
		if spec.Provider != string(provider.ProviderGCP) {
			requestedContainer.ContainerRegion = spec.ContainerRegion
		} // TODO: else for GCP regions
		createdContainer, err := req.service.CreateContainer(ctx, req.projectID, requestedContainer)
		if err != nil {
			return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
		}
		containerID = createdContainer.ID
		req.networkPeering.Status.ContainerID = containerID
	}
	requestedPeer := networkpeering.NetworkPeer{
		AtlasNetworkPeeringConfig: spec.AtlasNetworkPeeringConfig,
		ID:                        containerID,
	}
	createdPeer, err := req.service.CreatePeer(ctx, req.projectID, &requestedPeer)
	if err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	req.networkPeering.Status.ID = createdPeer.ID
	if err := r.saveStatus(&req.networkPeering.Status); err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasNetworkPeeringReconciler) sync(_ *reconcileRequest) (ctrl.Result, error) {
	panic("TBD")
}

func (r *AtlasNetworkPeeringReconciler) delete(_ *reconcileRequest) (ctrl.Result, error) {
	panic("TBD")
}

func (r *AtlasNetworkPeeringReconciler) noop(_ *reconcileRequest) (ctrl.Result, error) {
	panic("TBD")
}

func (r *AtlasNetworkPeeringReconciler) saveStatus(_ *status.AtlasNetworkPeeringStatus) error {
	panic("TBD")
}

func (r *AtlasNetworkPeeringReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) (ctrl.Result, error) {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err.Error())
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	return result.ReconcileResult(), nil
}

func containerID(peer *akov2.AtlasNetworkPeering) string {
	if peer.Spec.ContainerID != "" {
		return peer.Spec.ContainerID
	}
	return peer.Status.ContainerID
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasNetworkPeeringReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&akov2.AtlasNetworkPeering{}).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.networkPeeringForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.networkPeeringForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasNetworkPeeringReconciler) networkPeeringForProjectMapFunc() handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		atlasProject, ok := obj.(*akov2.AtlasProject)
		if !ok {
			r.Log.Warnf("watching Project but got %T", obj)

			return nil
		}

		npList := &akov2.AtlasNetworkPeeringList{}
		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasNetworkPeeringByProjectIndex,
				client.ObjectKeyFromObject(atlasProject).String(),
			),
		}
		err := r.Client.List(ctx, npList, listOpts)
		if err != nil {
			r.Log.Errorf("failed to list AtlasPrivateEndpoint: %s", err)

			return []reconcile.Request{}
		}

		requests := make([]reconcile.Request, 0, len(npList.Items))
		for _, item := range npList.Items {
			requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace}})
		}

		return requests
	}
}

func (r *AtlasNetworkPeeringReconciler) networkPeeringForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasNetworkPeeringCredentialsIndex,
		func() *akov2.AtlasNetworkPeeringList { return &akov2.AtlasNetworkPeeringList{} },
		indexer.NetworkPeeringRequests,
		r.Client,
		r.Log,
	)
}
