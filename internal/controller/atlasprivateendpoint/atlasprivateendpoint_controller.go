/*
Copyright 2024 MongoDB.

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

package atlasprivateendpoint

import (
	"context"
	"errors"
	"fmt"
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
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

// AtlasPrivateEndpointReconciler reconciles a AtlasPrivateEndpoint object
type AtlasPrivateEndpointReconciler struct {
	reconciler.AtlasReconciler
	Scheme           *runtime.Scheme
	EventRecorder    record.EventRecorder
	AtlasProvider    atlas.Provider
	GlobalPredicates []predicate.Predicate

	ObjectDeletionProtection bool
	independentSyncPeriod    time.Duration
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprivateendpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprivateendpoints/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprivateendpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprivateendpoints/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *AtlasPrivateEndpointReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infow("-> Starting AtlasPrivateEndpoint reconciliation")

	akoPrivateEndpoint := akov2.AtlasPrivateEndpoint{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoPrivateEndpoint, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), errors.New(result.GetMessage())
	}

	return r.ensureCustomResource(ctx, &akoPrivateEndpoint)
}

func (r *AtlasPrivateEndpointReconciler) ensureCustomResource(ctx context.Context, akoPrivateEndpoint *akov2.AtlasPrivateEndpoint) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(akoPrivateEndpoint) {
		return r.skip(ctx, akoPrivateEndpoint), nil
	}

	conditions := api.InitCondition(akoPrivateEndpoint, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, akoPrivateEndpoint)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, akoPrivateEndpoint)

	isValid := customresource.ValidateResourceVersion(workflowCtx, akoPrivateEndpoint, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid), nil
	}

	if !r.AtlasProvider.IsResourceSupported(akoPrivateEndpoint) {
		return r.unsupport(workflowCtx), nil
	}

	credentials, err := r.ResolveCredentials(ctx, akoPrivateEndpoint)
	if err != nil {
		return r.terminate(workflowCtx, akoPrivateEndpoint, nil, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClient, orgID, err := r.AtlasProvider.SdkClient(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, akoPrivateEndpoint, nil, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	atlasProject, err := r.ResolveProject(ctx, sdkClient, akoPrivateEndpoint, orgID)
	if err != nil {
		return r.terminate(workflowCtx, akoPrivateEndpoint, nil, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	privateEndpointService := privateendpoint.NewPrivateEndpointAPI(sdkClient.PrivateEndpointServicesApi)

	return r.handlePrivateEndpointService(workflowCtx, privateEndpointService, atlasProject.ID, akoPrivateEndpoint)
}

func (r *AtlasPrivateEndpointReconciler) skip(ctx context.Context, akoPrivateEndpoint *akov2.AtlasPrivateEndpoint) ctrl.Result {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasPrivateEndpoint reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", akoPrivateEndpoint.Spec)
	if !akoPrivateEndpoint.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, akoPrivateEndpoint, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasPrivateEndpointReconciler) invalidate(invalid workflow.Result) ctrl.Result {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasPrivateEndpoint is invalid: %v", invalid)
	return invalid.ReconcileResult()
}

func (r *AtlasPrivateEndpointReconciler) unsupport(ctx *workflow.Context) ctrl.Result {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, errors.New("the AtlasPrivateEndpoint is not supported by Atlas for government")).
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult()
}

func (r *AtlasPrivateEndpointReconciler) terminate(
	ctx *workflow.Context,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	atlasPEService privateendpoint.EndpointService,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	err error,
) (ctrl.Result, error) {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", akoPrivateEndpoint, akoPrivateEndpoint.GetNamespace(), akoPrivateEndpoint.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	if atlasPEService != nil {
		ctx.EnsureStatusOption(privateendpoint.NewPrivateEndpointStatus(atlasPEService))
	}

	return result.ReconcileResult(), nil
}

func (r *AtlasPrivateEndpointReconciler) inProgress(
	ctx *workflow.Context,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	atlasPEService privateendpoint.EndpointService,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	msg string,
) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, akoPrivateEndpoint, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	result := workflow.InProgress(reason, msg)
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result).
		EnsureStatusOption(privateendpoint.NewPrivateEndpointStatus(atlasPEService))

	return result.ReconcileResult(), nil
}

func (r *AtlasPrivateEndpointReconciler) ready(ctx *workflow.Context, akoPrivateEndpoint *akov2.AtlasPrivateEndpoint, atlasPEService privateendpoint.EndpointService) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, akoPrivateEndpoint, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	ctx.SetConditionTrue(api.PrivateEndpointServiceReady).
		SetConditionTrue(api.PrivateEndpointReady).
		SetConditionTrue(api.ReadyType).
		EnsureStatusOption(privateendpoint.NewPrivateEndpointStatus(atlasPEService))

	if akoPrivateEndpoint.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult(), nil
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasPrivateEndpointReconciler) waitForConfiguration(ctx *workflow.Context, akoPrivateEndpoint *akov2.AtlasPrivateEndpoint, atlasPEService privateendpoint.EndpointService) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, akoPrivateEndpoint, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	result := workflow.InProgress(workflow.PrivateEndpointConfigurationPending, "waiting for private endpoint configuration from customer side").
		WithoutRetry()
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionTrue(api.PrivateEndpointServiceReady).
		SetConditionFromResult(api.PrivateEndpointReady, result).
		EnsureStatusOption(privateendpoint.NewPrivateEndpointStatus(atlasPEService))

	return result.ReconcileResult(), nil
}

func (r *AtlasPrivateEndpointReconciler) unmanage(ctx *workflow.Context, akoPrivateEndpoint *akov2.AtlasPrivateEndpoint) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, akoPrivateEndpoint, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, nil, api.ReadyType, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.Deleted().ReconcileResult(), nil
}

func (r *AtlasPrivateEndpointReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasPrivateEndpoint{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasPrivateEndpointReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasPrivateEndpoint").
		For(r.For()).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.privateEndpointForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.privateEndpointForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasPrivateEndpointReconciler) privateEndpointForProjectMapFunc() handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		atlasProject, ok := obj.(*akov2.AtlasProject)
		if !ok {
			r.Log.Warnf("watching Project but got %T", obj)

			return nil
		}

		peList := &akov2.AtlasPrivateEndpointList{}
		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasPrivateEndpointByProjectIndex,
				client.ObjectKeyFromObject(atlasProject).String(),
			),
		}
		err := r.Client.List(ctx, peList, listOpts)
		if err != nil {
			r.Log.Errorf("failed to list AtlasPrivateEndpoint: %s", err)

			return []reconcile.Request{}
		}

		requests := make([]reconcile.Request, 0, len(peList.Items))
		for _, item := range peList.Items {
			requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace}})
		}

		return requests
	}
}

func (r *AtlasPrivateEndpointReconciler) privateEndpointForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasPrivateEndpointCredentialsIndex,
		func() *akov2.AtlasPrivateEndpointList { return &akov2.AtlasPrivateEndpointList{} },
		indexer.PrivateEndpointRequests,
		r.Client,
		r.Log,
	)
}

func NewAtlasPrivateEndpointReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	independentSyncPeriod time.Duration,
	logger *zap.Logger,
) *AtlasPrivateEndpointReconciler {
	return &AtlasPrivateEndpointReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: c.GetClient(),
			Log:    logger.Named("controllers").Named("AtlasPrivateEndpoint").Sugar(),
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasPrivateEndpoint"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}
