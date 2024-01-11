package controller

import (
	"fmt"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"

	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func RegisterControllers(mgr manager.Manager, provider atlas.Provider, cfg Config, globalPredicates []predicate.Predicate, logger *zap.Logger) error {
	if err := (&atlasproject.AtlasProjectReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasProject"),
		AtlasProvider:               provider,
		ObjectDeletionProtection:    cfg.DeletionProtection.Object,
		SubObjectDeletionProtection: cfg.DeletionProtection.SubObject,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create AtlasProject controller: %w", err)
	}

	if err := (&atlasdeployment.AtlasDeploymentReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDeployment"),
		AtlasProvider:               provider,
		ObjectDeletionProtection:    cfg.DeletionProtection.Object,
		SubObjectDeletionProtection: cfg.DeletionProtection.SubObject,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create AtlasDeployment controller: %w", err)
	}

	if err := (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		Client:                        mgr.GetClient(),
		Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:                        mgr.GetScheme(),
		ResourceWatcher:               watch.NewResourceWatcher(),
		GlobalPredicates:              globalPredicates,
		EventRecorder:                 mgr.GetEventRecorderFor("AtlasDatabaseUser"),
		AtlasProvider:                 provider,
		ObjectDeletionProtection:      cfg.DeletionProtection.Object,
		SubObjectDeletionProtection:   cfg.DeletionProtection.SubObject,
		FeaturePreviewOIDCAuthEnabled: cfg.FeatureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create AtlasDatabaseUser controller: %w", err)
	}

	if err := (&atlasdatafederation.AtlasDataFederationReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDataFederation"),
		AtlasProvider:               provider,
		ObjectDeletionProtection:    cfg.DeletionProtection.Object,
		SubObjectDeletionProtection: cfg.DeletionProtection.SubObject,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create AtlasDataFederation controller: %w", err)
	}

	if err := (&atlasfederatedauth.AtlasFederatedAuthReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasFederatedAuth"),
		AtlasProvider:               provider,
		ObjectDeletionProtection:    cfg.DeletionProtection.Object,
		SubObjectDeletionProtection: cfg.DeletionProtection.SubObject,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create AtlasFederatedAuth controller: %w", err)
	}

	return nil
}
