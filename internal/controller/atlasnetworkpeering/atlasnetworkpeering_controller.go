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
	"fmt"
	"reflect"
	"time"

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
	independentSyncPeriod    time.Duration
}

func NewAtlasNetworkPeeringsReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
	independentSyncPeriod time.Duration,
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
		independentSyncPeriod:    independentSyncPeriod,
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
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err), nil
	}
	sdkClient, orgID, err := r.AtlasProvider.SdkClient(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err), nil
	}
	project, err := r.ResolveProject(ctx, sdkClient, &akoNetworkPeering, orgID)
	if err != nil {
		return r.terminate(workflowCtx, &akoNetworkPeering, workflow.AtlasAPIAccessNotConfigured, err), nil
	}
	return r.handle(&reconcileRequest{
		workflowCtx:    workflowCtx,
		service:        networkpeering.NewNetworkPeeringService(sdkClient.NetworkPeeringApi),
		projectID:      project.ID,
		networkPeering: &akoNetworkPeering,
	}), nil
}

func (r *AtlasNetworkPeeringReconciler) handle(req *reconcileRequest) ctrl.Result {
	r.Log.Infow("handling network peering reconcile request",
		"service set", (req.service != nil), "projectID", req.projectID, "networkPeering", req.networkPeering)
	var atlasPeer *networkpeering.NetworkPeer
	if req.networkPeering.Status.ID != "" {
		peer, err := req.service.GetPeer(req.workflowCtx.Context, req.projectID, req.networkPeering.Status.ID)
		if err != nil && !errors.Is(err, networkpeering.ErrNotFound) {
			return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
		}
		atlasPeer = peer
	}
	inAtlas := atlasPeer != nil
	deleted := req.networkPeering.DeletionTimestamp != nil

	switch {
	case !deleted && !inAtlas:
		return r.create(req)
	case !deleted && inAtlas:
		return r.sync(req, atlasPeer)
	case deleted && inAtlas:
		return r.delete(req, atlasPeer)
	default:
		return r.unmanage(req) // deleted && !inAtlas
	}
}

func (r *AtlasNetworkPeeringReconciler) create(req *reconcileRequest) ctrl.Result {
	container, err := r.handleContainer(req)
	if err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	req.workflowCtx.EnsureStatusOption(updateContainerStatusOption(container))
	specPeer := networkpeering.NewNetworkPeeringSpec(&req.networkPeering.Spec.AtlasNetworkPeeringConfig)
	specPeer.ContainerID = container.ID
	newPeer, err := req.service.CreatePeer(req.workflowCtx.Context, req.projectID, specPeer)
	if err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	req.workflowCtx.EnsureStatusOption(updatePeeringStatusOption(newPeer))
	return workflow.InProgress(
		workflow.NetworkPeeringConnectionCreating,
		fmt.Sprintf("Network Peering Connection %s is %s",
			req.networkPeering.Status.ID, req.networkPeering.Status.Status),
	).ReconcileResult()
}

func (r *AtlasNetworkPeeringReconciler) sync(req *reconcileRequest, atlasPeer *networkpeering.NetworkPeer) ctrl.Result {
	container, err := r.handleContainer(req)
	if err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	}
	req.workflowCtx.EnsureStatusOption(updateContainerStatusOption(container))
	switch {
	case atlasPeer.Failed():
		err := fmt.Errorf("peer connection failed: %s", atlasPeer.ErrorMessage)
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, err)
	case !atlasPeer.Available():
		return r.inProgress(req, atlasPeer)
	}
	return r.ready(req, atlasPeer)
}

func (r *AtlasNetworkPeeringReconciler) delete(req *reconcileRequest, atlasPeer *networkpeering.NetworkPeer) ctrl.Result {
	id := req.networkPeering.Status.ID
	peer := atlasPeer
	if id != "" && !atlasPeer.Closing() {
		if err := req.service.DeletePeer(req.workflowCtx.Context, req.projectID, id); err != nil {
			wrappedErr := fmt.Errorf("failed to delete peer connection %s: %w", id, err)
			return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
		}
		closingPeer, err := req.service.GetPeer(req.workflowCtx.Context, req.projectID, id)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get closing peer connection %s: %w", id, err)
			return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
		}
		peer = closingPeer
	}
	return r.inProgress(req, peer)
}

func (r *AtlasNetworkPeeringReconciler) unmanage(req *reconcileRequest) ctrl.Result {
	req.workflowCtx.EnsureStatusOption(clearPeeringStatusOption())
	if _, err := r.handleContainer(req); err != nil {
		containerID := containerID(req.networkPeering)
		if !errors.Is(err, networkpeering.ErrContainerInUse) {
			wrappedErr := fmt.Errorf("failed to clear container %s: %w", containerID, err)
			return r.terminate(req.workflowCtx, req.networkPeering, workflow.Internal, wrappedErr)
		}
	}
	req.workflowCtx.EnsureStatusOption(clearContainerStatusOption())
	if err := customresource.ManageFinalizer(req.workflowCtx.Context, r.Client, req.networkPeering, customresource.UnsetFinalizer); err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.Deleted().ReconcileResult()
}

func (r *AtlasNetworkPeeringReconciler) inProgress(req *reconcileRequest, peer *networkpeering.NetworkPeer) ctrl.Result {
	req.workflowCtx.EnsureStatusOption(updatePeeringStatusOption(peer))

	return workflow.InProgress(
		workflow.NetworkPeeringConnectionPending,
		fmt.Sprintf("Network Peering Connection %s is %s", peer.ID, peer.Status),
	).ReconcileResult()
}

func (r *AtlasNetworkPeeringReconciler) ready(req *reconcileRequest, peer *networkpeering.NetworkPeer) ctrl.Result {
	if err := customresource.ManageFinalizer(req.workflowCtx.Context, r.Client, req.networkPeering, customresource.SetFinalizer); err != nil {
		return r.terminate(req.workflowCtx, req.networkPeering, workflow.AtlasFinalizerNotSet, err)
	}

	req.workflowCtx.EnsureStatusOption(updatePeeringStatusOption(peer))
	req.workflowCtx.SetConditionTrue(api.ReadyType)

	if req.networkPeering.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult()
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasNetworkPeeringReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) ctrl.Result {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}

func updateContainerStatusOption(container *networkpeering.ProviderContainer) status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		applyContainerStatus(peeringStatus, container)
	}
}

func applyContainerStatus(peeringStatus *status.AtlasNetworkPeeringStatus, container *networkpeering.ProviderContainer) {
	peeringStatus.ContainerID = container.ID
	peeringStatus.ContainerProvisioned = container.Provisioned
	providerName := provider.ProviderName(container.Provider)
	switch {
	case providerName == provider.ProviderAWS && container.AWSStatus != nil:
		if peeringStatus.AWSStatus == nil {
			peeringStatus.AWSStatus = &status.AWSStatus{}
		}
		peeringStatus.AWSStatus.ContainerVpcID = container.AWSStatus.VpcID
	case providerName == provider.ProviderAzure && container.AzureStatus != nil:
		if peeringStatus.AzureStatus == nil {
			peeringStatus.AzureStatus = &status.AzureStatus{}
		}
		peeringStatus.AzureStatus.AzureSubscriptionID = container.AzureStatus.AzureSubscriptionID
		peeringStatus.AzureStatus.VnetName = container.AzureStatus.VnetName
	case providerName == provider.ProviderGCP && container.GoogleStatus != nil:
		if peeringStatus.GoogleStatus == nil {
			peeringStatus.GoogleStatus = &status.GoogleStatus{}
		}
		peeringStatus.GoogleStatus.GCPProjectID = container.GoogleStatus.GCPProjectID
		peeringStatus.GoogleStatus.NetworkName = container.GoogleStatus.NetworkName
	}
}

func updatePeeringStatusOption(peer *networkpeering.NetworkPeer) status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		applyPeeringStatus(peeringStatus, peer)
	}
}

func applyPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus, peer *networkpeering.NetworkPeer) {
	peeringStatus.ID = peer.ID
	peeringStatus.Status = peer.Status
	peeringStatus.Error = peer.ErrorMessage
	providerName := provider.ProviderName(peer.Provider)
	if providerName == provider.ProviderAWS && peer.AWSStatus != nil {
		if peeringStatus.AWSStatus == nil {
			peeringStatus.AWSStatus = &status.AWSStatus{}
		}
		peeringStatus.AWSStatus.ConnectionID = peer.AWSStatus.ConnectionID
	}
}

func clearPeeringStatusOption() status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		clearPeeringStatus(peeringStatus)
	}
}

func clearPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus) {
	peeringStatus.ID = ""
	peeringStatus.Status = ""
	peeringStatus.Error = ""
	if peeringStatus.AWSStatus != nil {
		peeringStatus.AWSStatus.ConnectionID = ""
	}
}

func clearContainerStatusOption() status.AtlasNetworkPeeringStatusOption {
	return func(peeringStatus *status.AtlasNetworkPeeringStatus) {
		clearContainerStatus(peeringStatus)
	}
}

func clearContainerStatus(peeringStatus *status.AtlasNetworkPeeringStatus) {
	peeringStatus.ContainerID = ""
	peeringStatus.ContainerProvisioned = false
	peeringStatus.AWSStatus = nil
	peeringStatus.AzureStatus = nil
	peeringStatus.GoogleStatus = nil
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
