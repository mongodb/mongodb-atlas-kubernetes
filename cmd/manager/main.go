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

	"go.uber.org/zap/zapcore"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

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

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"
)

const (
	objectDeletionProtectionDefault    = false
	subobjectDeletionProtectionDefault = false

	objectDeletionProtectionEnvVar    = "UNSUPPORTED_OBJECT_DELETION_PROTECTION"
	subobjectDeletionProtectionEnvVar = "UNSUPPORTED_SUBOBJECT_DELETION_PROTECTION"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))
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
		cacheFunc = cache.MultiNamespacedCacheBuilder(namespaces)
	} else {
		cacheFunc = cache.BuilderWithOptions(cache.Options{
			SelectorsByObject: cache.SelectorsByObject{
				&corev1.Secret{}: {
					Label: labels.SelectorFromSet(labels.Set{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					}),
				},
			},
		})
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsAddr,
		Port:                   9443,
		Namespace:              config.Namespace,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "06d035fb.mongodb.com",
		SyncPeriod:             &syncPeriod,
		NewCache:               cacheFunc,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// globalPredicates should be used for general controller Predicates
	// that should be applied to all controllers in order to limit the
	// resources they receive events for.
	globalPredicates := []predicate.Predicate{
		watch.CommonPredicates(),                                  // ignore spurious changes. status changes etc.
		watch.SelectNamespacesPredicate(config.WatchedNamespaces), // select only desired namespaces
	}

	if err = (&atlasdeployment.AtlasDeploymentReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		Scheme:                      mgr.GetScheme(),
		AtlasDomain:                 config.AtlasDomain,
		GlobalAPISecret:             config.GlobalAPISecret,
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDeployment"),
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDeployment")
		os.Exit(1)
	}

	if err = (&atlasproject.AtlasProjectReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:                      mgr.GetScheme(),
		AtlasDomain:                 config.AtlasDomain,
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalAPISecret:             config.GlobalAPISecret,
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasProject"),
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasProject")
		os.Exit(1)
	}

	if err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		ResourceWatcher:             watch.NewResourceWatcher(),
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:                      mgr.GetScheme(),
		AtlasDomain:                 config.AtlasDomain,
		GlobalAPISecret:             config.GlobalAPISecret,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDatabaseUser"),
		GlobalPredicates:            globalPredicates,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDatabaseUser")
		os.Exit(1)
	}

	if err = (&atlasdatafederation.AtlasDataFederationReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		Scheme:                      mgr.GetScheme(),
		AtlasDomain:                 config.AtlasDomain,
		GlobalAPISecret:             config.GlobalAPISecret,
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDataFederation"),
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDataFederation")
		os.Exit(1)
	}

	if err = (&atlasfederatedauth.AtlasFederatedAuthReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
		Scheme:                      mgr.GetScheme(),
		AtlasDomain:                 config.AtlasDomain,
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasFederatedAuth"),
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasFederatedAuth")
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

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
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
	config.ObjectDeletionProtection = objectDeletionProtectionDefault
	config.SubObjectDeletionProtection = subobjectDeletionProtectionDefault

	// TODO: replace with the CLI flags at feature completion
	enableDeletionProtectionFromEnvVars(config, version.Version)
}

func enableDeletionProtectionFromEnvVars(config *Config, v string) {
	if version.IsRelease(v) {
		if isOn(os.Getenv(objectDeletionProtectionEnvVar)) ||
			isOn(os.Getenv(subobjectDeletionProtectionEnvVar)) {
			log.Printf("Deletion Protection feature is not available yet in production releases")
		}
		return
	}

	if isOn(os.Getenv(objectDeletionProtectionEnvVar)) {
		config.ObjectDeletionProtection = true
	}
	if isOn(os.Getenv(subobjectDeletionProtectionEnvVar)) {
		config.SubObjectDeletionProtection = true
	}
}

func isOn(value string) bool {
	return strings.ToLower(value) == "on"
}
