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

package controller

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasbackupcompliancepolicy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlascustomrole"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasipaccesslist"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasnetworkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasnetworkpeering"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasorgsettings"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasprivateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlassearchindexconfig"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasstream"
	integrations "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasthirdpartyintegrations"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/group"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

const DefaultReapplySupport = true

type Reconciler interface {
	reconcile.Reconciler
	For() (client.Object, builder.Predicates)
	SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error
}

type Registry struct {
	sharedPredicates      []predicate.Predicate
	deletionProtection    bool
	independentSyncPeriod time.Duration
	featureFlags          *featureflags.FeatureFlags

	logger          *zap.Logger
	reconcilers     []Reconciler
	globalSecretRef client.ObjectKey

	reapplySupport          bool
	maxConcurrentReconciles int
}

func NewRegistry(predicates []predicate.Predicate, deletionProtection bool, logger *zap.Logger, independentSyncPeriod time.Duration, featureFlags *featureflags.FeatureFlags, globalSecretRef client.ObjectKey, maxConcurrentReconciles int) *Registry {
	return &Registry{
		sharedPredicates:        predicates,
		deletionProtection:      deletionProtection,
		logger:                  logger,
		independentSyncPeriod:   independentSyncPeriod,
		featureFlags:            featureFlags,
		globalSecretRef:         globalSecretRef,
		reapplySupport:          DefaultReapplySupport,
		maxConcurrentReconciles: maxConcurrentReconciles,
	}
}

func (r *Registry) RegisterWithDryRunManager(mgr *dryrun.Manager, ap atlas.Provider) error {
	r.registerControllers(mgr, ap)

	for _, reconciler := range r.reconcilers {
		mgr.SetupReconciler(reconciler)
	}

	return nil
}

func (r *Registry) RegisterWithManager(mgr ctrl.Manager, skipNameValidation bool, ap atlas.Provider) error {
	r.registerControllers(mgr, ap)

	for _, reconciler := range r.reconcilers {
		if err := reconciler.SetupWithManager(mgr, skipNameValidation); err != nil {
			return fmt.Errorf("failed to set up with manager: %w", err)
		}
	}
	return nil
}

func (r *Registry) registerControllers(c cluster.Cluster, ap atlas.Provider) {
	if len(r.reconcilers) > 0 {
		return
	}

	var reconcilers []Reconciler
	reconcilers = append(reconcilers, atlasproject.NewAtlasProjectReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasdeployment.NewAtlasDeploymentReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasdatabaseuser.NewAtlasDatabaseUserReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.featureFlags, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasdatafederation.NewAtlasDataFederationReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasfederatedauth.NewAtlasFederatedAuthReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsInstanceReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsConnectionReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlassearchindexconfig.NewAtlasSearchIndexConfigReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasbackupcompliancepolicy.NewAtlasBackupCompliancePolicyReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlascustomrole.NewAtlasCustomRoleReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasprivateendpoint.NewAtlasPrivateEndpointReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasipaccesslist.NewAtlasIPAccessListReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasnetworkcontainer.NewAtlasNetworkContainerReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.logger, r.independentSyncPeriod, r.globalSecretRef, r.maxConcurrentReconciles))
	reconcilers = append(reconcilers, atlasnetworkpeering.NewAtlasNetworkPeeringsReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.logger, r.independentSyncPeriod, r.globalSecretRef, r.maxConcurrentReconciles))

	orgSettingsReconciler := atlasorgsettings.NewAtlasOrgSettingsReconciler(c, ap, r.logger, r.globalSecretRef, r.reapplySupport)
	reconcilers = append(reconcilers, newCtrlStateReconciler(orgSettingsReconciler, r.maxConcurrentReconciles))
	integrationsReconciler := integrations.NewAtlasThirdPartyIntegrationsReconciler(c, ap, r.deletionProtection, r.logger, r.globalSecretRef, r.reapplySupport)
	reconcilers = append(reconcilers, newCtrlStateReconciler(integrationsReconciler, r.maxConcurrentReconciles))

	if version.IsExperimental() {
		// Add experimental controllers here
		reconcilers = append(reconcilers, connectionsecret.NewConnectionSecretReconciler(c, r.defaultPredicates(), ap, r.logger, r.globalSecretRef))
		groupReconciler := group.NewGroupReconciler(c, ap, r.logger, r.globalSecretRef, r.deletionProtection, true, r.defaultPredicates())
		reconcilers = append(reconcilers, newCtrlStateReconciler(groupReconciler, r.maxConcurrentReconciles))
	}

	r.reconcilers = reconcilers
}

// deprecatedPredicates are to be phased out in favor of defaultPredicates
func (r *Registry) deprecatedPredicates() []predicate.Predicate {
	return append(r.sharedPredicates, watch.DeprecatedCommonPredicates[client.Object]())
}

// defaultPredicates minimize the reconciliations controllers actually do, avoiding
// spurious after delete handling and acting on finalizers setting or unsetting
func (r *Registry) defaultPredicates() []predicate.Predicate {
	return append(r.sharedPredicates, watch.DefaultPredicates[client.Object]())
}

type ctrlStateReconciler[T any] struct {
	*ctrlstate.Reconciler[T]
	maxConcurrentReconciles int
}

func newCtrlStateReconciler[T any](r *ctrlstate.Reconciler[T], maxConcurrentReconciles int) *ctrlStateReconciler[T] {
	return &ctrlStateReconciler[T]{Reconciler: r, maxConcurrentReconciles: maxConcurrentReconciles}
}

func (nr *ctrlStateReconciler[T]) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	defaultReconcilerOptions := controller.TypedOptions[reconcile.Request]{
		RateLimiter:             ratelimit.NewRateLimiter[reconcile.Request](),
		SkipNameValidation:      pointer.MakePtr(skipNameValidation),
		MaxConcurrentReconciles: nr.maxConcurrentReconciles,
	}
	return nr.Reconciler.SetupWithManager(mgr, defaultReconcilerOptions)
}
