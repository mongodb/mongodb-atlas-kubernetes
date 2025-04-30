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

package integrations

import (
	"context"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi"
)

// AtlasThirdPartyIntegrationsReconciler reconciles a AtlasNetworkPeering object
type AtlasThirdPartyIntegrationsReconciler struct {
	reconciler.AtlasReconciler
	AtlasProvider            atlas.Provider
	Scheme                   *runtime.Scheme
	EventRecorder            record.EventRecorder
	GlobalPredicates         []predicate.Predicate
	ObjectDeletionProtection bool
	independentSyncPeriod    time.Duration
}

func NewAtlas3rdPartyIntegrationsReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
	independentSyncPeriod time.Duration,
	globalSecretRef client.ObjectKey,
) *AtlasThirdPartyIntegrationsReconciler {
	return &AtlasThirdPartyIntegrationsReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             logger.Named("controllers").Named("Atlas3rdPartyIntegrationsReconciler").Sugar(),
			GlobalSecretRef: globalSecretRef,
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasThirdPartyIntegration"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AtlasThirdPartyIntegrationsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infow("-> Starting AtlasThirdPartyIntegration reconciliation")

	akoIntegrations := nextapi.AtlasThirdPartyIntegration{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoIntegrations, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	return r.handleCustomResource(ctx, &akoIntegrations)
}

func (r *AtlasThirdPartyIntegrationsReconciler) handleCustomResource(ctx context.Context, integration *nextapi.AtlasThirdPartyIntegration) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(integration) {
		return r.Skip(ctx, "AtlasThirdPartyIntegration", integration, &integration.Spec)
	}

	panic("unimplemented")
}

// For prepares the controller for its target Custom Resource; Network Containers
func (r *AtlasThirdPartyIntegrationsReconciler) For() (client.Object, builder.Predicates) {
	return &nextapi.AtlasThirdPartyIntegration{}, builder.WithPredicates(r.GlobalPredicates...)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasThirdPartyIntegrationsReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.For()).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.integrationForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.integrationForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasThirdPartyIntegrationsReconciler) integrationForProjectMapFunc() handler.MapFunc {
	return indexer.ProjectsIndexMapperFunc(
		indexer.AtlasThirdPartyIntegrationByProjectIndex,
		func() *nextapi.AtlasThirdPartyIntegrationList { return &nextapi.AtlasThirdPartyIntegrationList{} },
		indexer.AtlasThirdPartyIntegrationRequests,
		r.Client,
		r.Log,
	)
}

func (r *AtlasThirdPartyIntegrationsReconciler) integrationForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasThirdPartyIntegrationCredentialsIndex,
		func() *nextapi.AtlasThirdPartyIntegrationList { return &nextapi.AtlasThirdPartyIntegrationList{} },
		indexer.AtlasThirdPartyIntegrationRequests,
		r.Client,
		r.Log,
	)
}
