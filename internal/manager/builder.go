package manager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/zapr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	logging "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/log"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func Setup(scheme *runtime.Scheme, kubeCfg *rest.Config, managerCfg Config, operatorCfg controller.Config, logCfg logging.Config, withCtx context.Context) error {
	logger, err := logging.NewLogger(logCfg, os.Stdout)
	if err != nil {
		return fmt.Errorf("unable to create logger: %w", err)
	}

	ctrl.SetLogger(zapr.NewLogger(logger))

	mgr, err := NewManager(scheme, kubeCfg, managerCfg, controller.GetOperatorNamespace())
	if err != nil {
		return fmt.Errorf("unable to create controller manager: %w", err)
	}

	atlasProvider := atlas.NewProductionProvider(operatorCfg.Domain, operatorCfg.APISecret, mgr.GetClient())

	err = controller.RegisterControllers(mgr, atlasProvider, operatorCfg, globalPredicates(managerCfg.Namespaces), logger)
	if err != nil {
		return fmt.Errorf("unable to create controllers: %w", err)
	}

	if err = mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}
	if err = mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	if withCtx == nil {
		withCtx = ctrl.SetupSignalHandler()
	}

	if err = mgr.Start(withCtx); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}

func NewManager(scheme *runtime.Scheme, kubeCfg *rest.Config, cfg Config, operatorNamespace string) (manager.Manager, error) {
	mgr, err := ctrl.NewManager(kubeCfg, ctrl.Options{
		Scheme:                  scheme,
		Cache:                   cacheOptions(cfg.Namespaces, cfg.SyncPeriod),
		LeaderElection:          cfg.EnableLeaderElection,
		LeaderElectionID:        "ako-cloud.mongodb.com",
		LeaderElectionNamespace: operatorNamespace,
		MetricsBindAddress:      cfg.MetricsBindAddress,
		HealthProbeBindAddress:  cfg.HealthProbeBindAddress,
		WebhookServer: webhook.NewServer(
			webhook.Options{
				Port: 9443,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to setup the manager: %w", err)
	}

	return mgr, nil
}

func cacheOptions(managedNamespaces []string, syncPeriod time.Duration) cache.Options {
	opts := cache.Options{
		SyncPeriod: &syncPeriod,
	}

	switch len(managedNamespaces) {
	case 0:
		if opts.ByObject == nil {
			opts.ByObject = map[client.Object]cache.ByObject{}
		}
		opts.ByObject[&corev1.Secret{}] = cache.ByObject{
			Label: labels.SelectorFromSet(labels.Set{connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal}),
		}
	default:
		opts.Namespaces = managedNamespaces
	}

	return opts
}

func globalPredicates(namespaces []string) []predicate.Predicate {
	managedNamespace := map[string]bool{}
	for _, ns := range namespaces {
		managedNamespace[ns] = true
	}

	if len(managedNamespace) == 0 {
		managedNamespace[""] = true
	}

	// globalPredicates should be used for general controller Predicates
	// that should be applied to all controllers in order to limit the
	// resources they receive events for.
	return []predicate.Predicate{
		watch.CommonPredicates(),                          // ignore spurious changes. status changes etc.
		watch.SelectNamespacesPredicate(managedNamespace), // select only desired namespaces
	}
}
