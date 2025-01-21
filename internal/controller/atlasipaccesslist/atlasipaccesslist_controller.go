/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package atlasipaccesslist

import (
	"context"
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
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

// AtlasIPAccessListReconciler reconciles a AtlasIPAccessList object
type AtlasIPAccessListReconciler struct {
	reconciler.AtlasReconciler

	Scheme                   *runtime.Scheme
	EventRecorder            record.EventRecorder
	AtlasProvider            atlas.Provider
	GlobalPredicates         []predicate.Predicate
	ObjectDeletionProtection bool
	independentSyncPeriod    time.Duration
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasipaccesslists,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasipaccesslists/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasipaccesslists/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasipaccesslists,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasipaccesslists/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasipaccesslists/finalizers,verbs=update

func (r *AtlasIPAccessListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infow("-> Starting AtlasIPAccessList reconciliation")

	ipAccessList := akov2.AtlasIPAccessList{}
	result := customresource.PrepareResource(ctx, r.Client, req, &ipAccessList, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	return r.ensureCustomResource(ctx, &ipAccessList), nil
}

func (r *AtlasIPAccessListReconciler) ensureCustomResource(ctx context.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	if customresource.ReconciliationShouldBeSkipped(ipAccessList) {
		return r.skip(ctx, ipAccessList)
	}

	conditions := api.InitCondition(ipAccessList, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, ipAccessList)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, ipAccessList)

	isValid := customresource.ValidateResourceVersion(workflowCtx, ipAccessList, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(ipAccessList) {
		return r.unsupport(workflowCtx)
	}

	credentials, err := r.ResolveCredentials(ctx, ipAccessList)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClient, orgID, err := r.AtlasProvider.SdkClient(ctx, credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	atlasProject, err := r.ResolveProject(ctx, sdkClient, ipAccessList, orgID)
	if err != nil {
		return r.terminate(workflowCtx, ipAccessList, api.ReadyType, workflow.AtlasAPIAccessNotConfigured, err)
	}
	ipAccessListService := ipaccesslist.NewIPAccessList(sdkClient.ProjectIPAccessListApi)

	return r.handleIPAccessList(workflowCtx, ipAccessListService, atlasProject.ID, ipAccessList)
}

func (r *AtlasIPAccessListReconciler) skip(ctx context.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasIPAccessList reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", ipAccessList.Spec)
	if !ipAccessList.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, ipAccessList, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err.Error())
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) invalidate(invalid workflow.Result) ctrl.Result {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasIPAccessList is invalid: %v", invalid)
	return invalid.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) unsupport(ctx *workflow.Context) ctrl.Result {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, "the AtlasIPAccessList is not supported by Atlas for government").
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) inProgress(
	ctx *workflow.Context,
	ipAccessList *akov2.AtlasIPAccessList,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	msg string,
) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	result := workflow.InProgress(reason, msg)
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) unmanage(ctx *workflow.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.Deleted().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) ready(ctx *workflow.Context, ipAccessList *akov2.AtlasIPAccessList) ctrl.Result {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, ipAccessList, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, ipAccessList, api.ReadyType, workflow.AtlasFinalizerNotSet, err)
	}

	ctx.SetConditionTrue(api.ReadyType).
		SetConditionTrue(api.IPAccessListReady)

	if ipAccessList.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult()
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasIPAccessListReconciler) terminate(
	ctx *workflow.Context,
	ipAccessList *akov2.AtlasIPAccessList,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	err error,
) ctrl.Result {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", ipAccessList, ipAccessList.GetNamespace(), ipAccessList.GetName(), condition, err)
	result := workflow.Terminate(reason, err.Error())
	ctx.SetConditionFalse(api.ReadyType).
		SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasIPAccessListReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasIPAccessList").
		For(&akov2.AtlasIPAccessList{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.ipAccessListForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.ipAccessListForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasIPAccessListReconciler) ipAccessListForProjectMapFunc() handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		atlasProject, ok := obj.(*akov2.AtlasProject)
		if !ok {
			r.Log.Warnf("watching Project but got %T", obj)

			return nil
		}

		list := &akov2.AtlasIPAccessListList{}
		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasIPAccessListCredentialsIndex,
				client.ObjectKeyFromObject(atlasProject).String(),
			),
		}
		err := r.Client.List(ctx, list, listOpts)
		if err != nil {
			r.Log.Errorf("failed to list AtlasIPAccessList: %s", err)

			return []reconcile.Request{}
		}

		requests := make([]reconcile.Request, 0, len(list.Items))
		for _, item := range list.Items {
			requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace}})
		}

		return requests
	}
}

func (r *AtlasIPAccessListReconciler) ipAccessListForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasIPAccessListCredentialsIndex,
		func() *akov2.AtlasIPAccessListList { return &akov2.AtlasIPAccessListList{} },
		indexer.IPAccessListRequests,
		r.Client,
		r.Log,
	)
}

func NewAtlasIPAccessListReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	independentSyncPeriod time.Duration,
	logger *zap.Logger,
) *AtlasIPAccessListReconciler {
	return &AtlasIPAccessListReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: mgr.GetClient(),
			Log:    logger.Named("controllers").Named("AtlasIPAccessList").Sugar(),
		},
		Scheme:                   mgr.GetScheme(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasIPAccessList"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}
