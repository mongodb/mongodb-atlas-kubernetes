package controller

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasbackupcompliancepolicy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlascustomrole"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasipaccesslist"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasprivateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlassearchindexconfig"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasstream"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
)

type ManagerAware interface {
	SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error
}

type Reconciler interface {
	dryrun.Reconciler
	ManagerAware
}

type Registry struct {
	predicates            []predicate.Predicate
	deletionProtection    bool
	independentSyncPeriod time.Duration
	featureFlags          *featureflags.FeatureFlags

	logger      *zap.Logger
	reconcilers []Reconciler
}

func NewRegistry(predicates []predicate.Predicate, deletionProtection bool, logger *zap.Logger, independentSyncPeriod time.Duration, featureFlags *featureflags.FeatureFlags) *Registry {
	return &Registry{
		predicates:            predicates,
		deletionProtection:    deletionProtection,
		logger:                logger,
		independentSyncPeriod: independentSyncPeriod,
		featureFlags:          featureFlags,
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
	reconcilers = append(reconcilers, atlasproject.NewAtlasProjectReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasdeployment.NewAtlasDeploymentReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasdatabaseuser.NewAtlasDatabaseUserReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.featureFlags, r.logger))
	reconcilers = append(reconcilers, atlasdatafederation.NewAtlasDataFederationReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasfederatedauth.NewAtlasFederatedAuthReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsInstanceReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsConnectionReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlassearchindexconfig.NewAtlasSearchIndexConfigReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasbackupcompliancepolicy.NewAtlasBackupCompliancePolicyReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlascustomrole.NewAtlasCustomRoleReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasprivateendpoint.NewAtlasPrivateEndpointReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasipaccesslist.NewAtlasIPAccessListReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	r.reconcilers = reconcilers
}
