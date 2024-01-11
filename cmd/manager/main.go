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
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	logging "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/log"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/manager"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

var (
	scheme      = runtime.NewScheme()
	setupLog    = ctrl.Log.WithName("setup")
	managerCfg  manager.Config
	operatorCfg controller.Config
	logCfg      logging.Config
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = mdbv1.AddToScheme(scheme)
}

func run(fs *flag.FlagSet) int {
	cfg, err := controller.DefaultConfig()
	if err != nil {
		setupLog.Error(err, "unable to create operator configuration")

		return 1
	}

	operatorCfg = cfg
	managerCfg = manager.DefaultConfig()
	logCfg = logging.DefaultConfig()

	parseConfiguration(fs)

	setupLog.Info("starting with configuration", "config", operatorCfg, "version", version.Version)

	err = manager.Setup(scheme, ctrl.GetConfigOrDie(), managerCfg, operatorCfg, logCfg, nil)
	if err != nil {
		setupLog.Error(err, "failed to setup the manager")

		return 1
	}

	return 0
}

func parseConfiguration(fs *flag.FlagSet) {
	manager.RegisterFlags(&managerCfg, fs)
	controller.RegisterFlags(&operatorCfg, fs)
	logging.RegisterFlags(&logCfg, fs)

	appVersion := fs.Bool("v", false, "prints application version")

	// No need to check for errors because Parse would exit on error.
	_ = fs.Parse(os.Args[1:])

	if *appVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	managerCfg.ParseWatchedNamespaces()
	operatorCfg.ParseFlagsFromEnv(fs)
}

func main() {
	os.Exit(run(flag.CommandLine))
}
