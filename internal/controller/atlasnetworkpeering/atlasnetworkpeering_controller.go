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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// AtlasNetworkPeeringReconciler reconciles a AtlasNetworkPeering object
type AtlasNetworkPeeringReconciler struct {
	reconciler.AtlasReconciler
	Scheme                   *runtime.Scheme
	EventRecorder            record.EventRecorder
	GlobalPredicates         []predicate.Predicate
	ObjectDeletionProtection bool
	independentSyncPeriod    time.Duration
}

func NewAtlasNetworkPeeringsReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
	independentSyncPeriod time.Duration,
	globalSecretRef client.ObjectKey,
) *AtlasNetworkPeeringReconciler {
	return &AtlasNetworkPeeringReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             logger.Named("controllers").Named("AtlasNetworkPeering").Sugar(),
			GlobalSecretRef: globalSecretRef,
			AtlasProvider:   atlasProvider,
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasPrivateEndpoint"),
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasnetworkpeerings/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasnetworkpeerings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasnetworkpeerings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasnetworkpeerings/finalizers,verbs=update

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
		return result.ReconcileResult(), nil
	}
	return r.handleCustomResource(ctx, &akoNetworkPeering)
}

// For prepares the controller for its target Custom Resource; Network Containers
func (r *AtlasNetworkPeeringReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasNetworkPeering{}, builder.WithPredicates(r.GlobalPredicates...)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasNetworkPeeringReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.For()).
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
		Watches(
			&akov2.AtlasNetworkContainer{},
			handler.EnqueueRequestsFromMapFunc(r.networkPeeringForContainerByIDMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasNetworkPeeringReconciler) networkPeeringForProjectMapFunc() handler.MapFunc {
	return indexer.ProjectsIndexMapperFunc(
		indexer.AtlasNetworkPeeringByProjectIndex,
		func() *akov2.AtlasNetworkPeeringList { return &akov2.AtlasNetworkPeeringList{} },
		indexer.NetworkPeeringRequests,
		r.Client,
		r.Log,
	)
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

func (r *AtlasNetworkPeeringReconciler) networkPeeringForContainerByIDMapFunc() handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		container, ok := obj.(*akov2.AtlasNetworkContainer)
		if !ok {
			r.Log.Warnf("watching AtlasNetworkContainer but got %T", obj)
			return nil
		}
		indexerName := indexer.AtlasNetworkPeeringByContainerIndex
		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasNetworkPeeringByContainerIndex,
				client.ObjectKeyFromObject(container).String(),
			),
		}
		list := &akov2.AtlasNetworkPeeringList{}
		err := r.Client.List(ctx, list, listOpts)
		if err != nil {
			r.Log.Errorf("failed to list from indexer %s: %v", indexerName, err)
			return nil
		}
		requests := make([]reconcile.Request, 0, len(list.Items))
		for _, peering := range list.Items {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      peering.Name,
					Namespace: peering.Namespace,
				},
			})
		}
		return requests
	}
}
