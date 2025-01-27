/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package atlasnetworkcontainer

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// AtlasNetworkContainerReconciler reconciles a AtlasNetworkContainer object
type AtlasNetworkContainerReconciler struct {
	reconciler.AtlasReconciler
	AtlasProvider            atlas.Provider
	Scheme                   *runtime.Scheme
	EventRecorder            record.EventRecorder
	GlobalPredicates         []predicate.Predicate
	ObjectDeletionProtection bool
	independentSyncPeriod    time.Duration
}

// type reconcileRequest struct {
// 	workflowCtx      *workflow.Context
// 	service          networkcontainer.NetworkContainerService
// 	projectID        string
// 	networkContainer *akov2.AtlasNetworkContainer
// }

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkcontainers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkcontainers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkcontainers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AtlasNetworkContainer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *AtlasNetworkContainerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infow("-> Starting AtlasNetworkContainer reconciliation")

	akoNetworkContainer := akov2.AtlasNetworkContainer{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoNetworkContainer, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), errors.New(result.GetMessage())
	}
	return r.handle(ctx, &akoNetworkContainer)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasNetworkContainerReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&akov2.AtlasNetworkContainer{}).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.networkContainerForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.networkContainerForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasNetworkContainerReconciler) networkContainerForProjectMapFunc() handler.MapFunc {
	return indexer.ProjectsIndexMapperFunc(
		indexer.AtlasNetworkContainerByProjectIndex,
		func() *akov2.AtlasNetworkContainerList { return &akov2.AtlasNetworkContainerList{} },
		indexer.NetworkContainerRequests,
		r.Client,
		r.Log,
	)
}

func (r *AtlasNetworkContainerReconciler) networkContainerForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasNetworkContainerCredentialsIndex,
		func() *akov2.AtlasNetworkContainerList { return &akov2.AtlasNetworkContainerList{} },
		indexer.NetworkContainerRequests,
		r.Client,
		r.Log,
	)
}

func NewAtlasNetworkContainerReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
	independentSyncPeriod time.Duration,
) *AtlasNetworkContainerReconciler {
	return &AtlasNetworkContainerReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: c.GetClient(),
			Log:    logger.Named("controllers").Named("AtlasNetworkContainer").Sugar(),
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasNetworkContainer"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}
