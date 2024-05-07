/*
Copyright 2020 The Kubernetes authors.

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

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlassearchindexconfig"

	"go.uber.org/zap/zapcore"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller"

	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasstream"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

const (
	objectDeletionProtectionFlag       = "object-deletion-protection"
	subobjectDeletionProtectionFlag    = "subobject-deletion-protection"
	objectDeletionProtectionEnvVar     = "OBJECT_DELETION_PROTECTION"
	subobjectDeletionProtectionEnvVar  = "SUBOBJECT_DELETION_PROTECTION"
	objectDeletionProtectionDefault    = true
	subobjectDeletionProtectionDefault = false
	subobjectDeletionProtectionMessage = "Note: sub-object deletion protection is IGNORED because it does not work deterministically."
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))
}

func main() {
	// controller-runtime/pkg/log/zap is a wrapper over zap that implements logr
	// logr looks quite limited in functionality so we better use Zap directly.
	// Though we still need the controller-runtime library and go-logr/zapr as they are used in controller-runtime
	// logging
	ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.StacktraceLevel(zap.ErrorLevel))
	config := parseConfiguration()
	logger, err := initCustomZapLogger(config.LogLevel, config.LogEncoder)
	if err != nil {
		fmt.Printf("error instantiating logger: %v\r\n", err)
		os.Exit(1)
	}

	logger.Info("starting with configuration", zap.Any("config", config), zap.Any("version", version.Version))

	ctrl.SetLogger(zapr.NewLogger(logger))

	syncPeriod := time.Hour * 3

	var cacheFunc cache.NewCacheFunc
	if len(config.WatchedNamespaces) > 1 {
		var namespaces []string
		for ns := range config.WatchedNamespaces {
			namespaces = append(namespaces, ns)
		}
		cacheFunc = controller.MultiNamespacedCacheBuilder(namespaces)
	} else {
		cacheFunc = controller.CustomLabelSelectorCacheBuilder(
			&corev1.Secret{},
			labels.SelectorFromSet(labels.Set{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			}),
		)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:  scheme,
		Metrics: metricsserver.Options{BindAddress: config.MetricsAddr},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{config.Namespace: {}},
			SyncPeriod:        &syncPeriod,
		},
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "06d035fb.mongodb.com",
		NewCache:               cacheFunc,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()

	// globalPredicates should be used for general controller Predicates
	// that should be applied to all controllers in order to limit the
	// resources they receive events for.
	globalPredicates := []predicate.Predicate{
		watch.CommonPredicates(),                                  // ignore spurious changes. status changes etc.
		watch.SelectNamespacesPredicate(config.WatchedNamespaces), // select only desired namespaces
	}

	atlasProvider := atlas.NewProductionProvider(config.AtlasDomain, config.GlobalAPISecret, mgr.GetClient())

	if err = (&atlasdeployment.AtlasDeploymentReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		Scheme:                      mgr.GetScheme(),
		DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDeployment"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDeployment")
		os.Exit(1)
	}

	if err = (&atlasproject.AtlasProjectReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:                      mgr.GetScheme(),
		DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasProject"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasProject")
		os.Exit(1)
	}

	if err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		DeprecatedResourceWatcher:     watch.NewDeprecatedResourceWatcher(),
		Client:                        mgr.GetClient(),
		Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:                        mgr.GetScheme(),
		EventRecorder:                 mgr.GetEventRecorderFor("AtlasDatabaseUser"),
		AtlasProvider:                 atlasProvider,
		GlobalPredicates:              globalPredicates,
		ObjectDeletionProtection:      config.ObjectDeletionProtection,
		SubObjectDeletionProtection:   false,
		FeaturePreviewOIDCAuthEnabled: config.FeatureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDatabaseUser")
		os.Exit(1)
	}

	if err = (&atlasdatafederation.AtlasDataFederationReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		Scheme:                      mgr.GetScheme(),
		DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDataFederation"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDataFederation")
		os.Exit(1)
	}

	if err = (&atlasfederatedauth.AtlasFederatedAuthReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
		Scheme:                      mgr.GetScheme(),
		DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasFederatedAuth"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasFederatedAuth")
		os.Exit(1)
	}

	if err = (&atlasstream.AtlasStreamsInstanceReconciler{
		Scheme:                      mgr.GetScheme(),
		Client:                      mgr.GetClient(),
		EventRecorder:               mgr.GetEventRecorderFor("AtlasStreamsInstance"),
		GlobalPredicates:            globalPredicates,
		Log:                         logger.Named("controllers").Named("AtlasStreamsInstance").Sugar(),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasStreamsInstance")
		os.Exit(1)
	}

	if err = (&atlasstream.AtlasStreamsConnectionReconciler{
		Scheme:                      mgr.GetScheme(),
		Client:                      mgr.GetClient(),
		EventRecorder:               mgr.GetEventRecorderFor("AtlasStreamsConnection"),
		GlobalPredicates:            globalPredicates,
		Log:                         logger.Named("controllers").Named("AtlasStreamsConnection").Sugar(),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasStreamsConnection")
		os.Exit(1)
	}

	if err = (&atlassearchindexconfig.AtlasSearchIndexConfigReconciler{
		Scheme:                      mgr.GetScheme(),
		Client:                      mgr.GetClient(),
		EventRecorder:               mgr.GetEventRecorderFor("AtlasSearchIndexConfig"),
		GlobalPredicates:            globalPredicates,
		Log:                         logger.Named("controllers").Named("AtlasSearchIndexConfig").Sugar(),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: false,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasSearchIndexConfig")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info(subobjectDeletionProtectionMessage)
	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

type Config struct {
	AtlasDomain                 string
	EnableLeaderElection        bool
	MetricsAddr                 string
	Namespace                   string
	WatchedNamespaces           map[string]bool
	ProbeAddr                   string
	GlobalAPISecret             client.ObjectKey
	LogLevel                    string
	LogEncoder                  string
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	FeatureFlags                *featureflags.FeatureFlags
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func parseConfiguration() Config {
	var globalAPISecretName string
	config := Config{}
	flag.StringVar(&config.AtlasDomain, "atlas-domain", "https://cloud.mongodb.com/", "the Atlas URL domain name (with slash in the end).")
	flag.StringVar(&config.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&config.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&globalAPISecretName, "global-api-secret-name", "", "The name of the Secret that contains Atlas API keys. "+
		"It is used by the Operator if AtlasProject configuration doesn't contain API key reference. Defaults to <deployment_name>-api-key.")
	flag.BoolVar(&config.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&config.LogLevel, "log-level", "info", "Log level. Available values: debug | info | warn | error | dpanic | panic | fatal")
	flag.StringVar(&config.LogEncoder, "log-encoder", "json", "Log encoder. Available values: json | console")
	flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "Defines if the operator deletes Atlas resource "+
		"when a Custom Resource is deleted")
	flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "Defines if the operator overwrites "+
		"(and consequently delete) subresources that were not previously created by the operator. "+subobjectDeletionProtectionMessage)
	appVersion := flag.Bool("v", false, "prints application version")
	flag.Parse()

	if *appVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	config.GlobalAPISecret = operatorGlobalKeySecretOrDefault(globalAPISecretName)

	// dev note: we pass the watched namespace as the env variable to use the Kubernetes Downward API. Unfortunately
	// there is no way to use it for container arguments
	watchedNamespace := os.Getenv("WATCH_NAMESPACE")
	config.WatchedNamespaces = make(map[string]bool)
	for _, namespace := range strings.Split(watchedNamespace, ",") {
		namespace = strings.TrimSpace(namespace)
		config.WatchedNamespaces[namespace] = true
	}

	if len(config.WatchedNamespaces) == 1 {
		config.Namespace = watchedNamespace
	}

	configureDeletionProtection(&config)

	config.FeatureFlags = featureflags.NewFeatureFlags(os.Environ)
	return config
}

func operatorGlobalKeySecretOrDefault(secretNameOverride string) client.ObjectKey {
	secretName := secretNameOverride
	if secretName == "" {
		operatorPodName := os.Getenv("OPERATOR_POD_NAME")
		if operatorPodName == "" {
			log.Fatal(`"OPERATOR_POD_NAME" environment variable must be set!`)
		}
		deploymentName, err := kube.ParseDeploymentNameFromPodName(operatorPodName)
		if err != nil {
			log.Fatalf(`Failed to get Operator Deployment name from "OPERATOR_POD_NAME" environment variable: %s`, err.Error())
		}
		secretName = deploymentName + "-api-key"
	}
	operatorNamespace := os.Getenv("OPERATOR_NAMESPACE")
	if operatorNamespace == "" {
		log.Fatal(`"OPERATOR_NAMESPACE" environment variable must be set!`)
	}

	return client.ObjectKey{Namespace: operatorNamespace, Name: secretName}
}

func initCustomZapLogger(level, encoding string) (*zap.Logger, error) {
	lv := zap.AtomicLevel{}
	err := lv.UnmarshalText([]byte(strings.ToLower(level)))
	if err != nil {
		return nil, err
	}

	enc := strings.ToLower(encoding)
	if enc != "json" && enc != "console" {
		return nil, errors.New("'encoding' parameter can only by either 'json' or 'console'")
	}

	cfg := zap.Config{
		Level:             lv,
		OutputPaths:       []string{"stdout"},
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          enc,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	return cfg.Build()
}

func configureDeletionProtection(config *Config) {
	if config == nil {
		return
	}

	objectDeletionSet := false
	subObjectDeletionSet := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == objectDeletionProtectionFlag {
			objectDeletionSet = true
		}

		if f.Name == subobjectDeletionProtectionFlag {
			subObjectDeletionSet = true
		}
	})

	if !objectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(objectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			config.ObjectDeletionProtection = true
		case "false":
			config.ObjectDeletionProtection = false
		default:
			config.ObjectDeletionProtection = objectDeletionProtectionDefault
		}
	}

	if !subObjectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(subobjectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			config.SubObjectDeletionProtection = true
		case "false":
			config.SubObjectDeletionProtection = false
		default:
			config.SubObjectDeletionProtection = subobjectDeletionProtectionDefault
		}
	}
}
