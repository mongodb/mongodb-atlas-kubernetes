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
	"flag"
	"log"
	"os"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlascluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(mdbv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	// controller-runtime/pkg/log/zap is a wrapper over zap that implements logr
	// logr looks quite limited in functionality so we better use Zap directly.
	// Though we still need the controller-runtime library and go-logr/zapr as they are used in controller-runtime
	// logging
	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.StacktraceLevel(zap.ErrorLevel))

	config := parseConfiguration(logger.Sugar())
	ctrl.SetLogger(zapr.NewLogger(logger))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "06d035fb.mongodb.com",
		Namespace:              config.WatchedNamespaces,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	operatorPod := operatorPodObjectKey()

	if err = (&atlascluster.AtlasClusterReconciler{
		Client:      mgr.GetClient(),
		Log:         logger.Named("controllers").Named("AtlasCluster").Sugar(),
		Scheme:      mgr.GetScheme(),
		AtlasDomain: config.AtlasDomain,
		OperatorPod: operatorPod,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasCluster")
		os.Exit(1)
	}

	if err = (&atlasproject.AtlasProjectReconciler{
		Client:          mgr.GetClient(),
		Log:             logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:          mgr.GetScheme(),
		AtlasDomain:     config.AtlasDomain,
		ResourceWatcher: watch.NewResourceWatcher(),
		OperatorPod:     operatorPod,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasProject")
		os.Exit(1)
	}

	if err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		Client:      mgr.GetClient(),
		Log:         logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:      mgr.GetScheme(),
		AtlasDomain: config.AtlasDomain,
		OperatorPod: operatorPod,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDatabaseUser")
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
	AtlasDomain          string
	EnableLeaderElection bool
	MetricsAddr          string
	WatchedNamespaces    string
	ProbeAddr            string
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func parseConfiguration(log *zap.SugaredLogger) Config {
	config := Config{}
	flag.StringVar(&config.AtlasDomain, "atlas-domain", "https://cloud.mongodb.com", "the Atlas URL domain name (no slash in the end).")
	flag.StringVar(&config.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&config.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&config.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	flag.Parse()

	// dev note: we pass the watched namespace as the env variable to use the Kubernetes Downward API. Unfortunately
	// there is no way to use it for container arguments
	watchedNamespace := os.Getenv("WATCHED_NAMESPACE")
	if watchedNamespace != "" {
		log.Infof("The Operator is watching the namespace %s", watchedNamespace)
	}
	config.WatchedNamespaces = watchedNamespace
	return config
}

func operatorPodObjectKey() client.ObjectKey {
	operatorName := os.Getenv("OPERATOR_NAME")
	if operatorName == "" {
		log.Fatal(`"OPERATOR_NAME" environment variable must be set!`)
	}
	operatorNamespace := os.Getenv("OPERATOR_NAMESPACE")
	if operatorNamespace == "" {
		log.Fatal(`"OPERATOR_NAMESPACE" environment variable must be set!`)
	}

	return client.ObjectKey{Namespace: operatorNamespace, Name: operatorName}
}
