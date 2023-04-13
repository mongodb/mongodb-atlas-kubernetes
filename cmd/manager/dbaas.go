/*
Copyright 2023 The Kubernetes authors.

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
	"fmt"

	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	dbaasv1beta1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1beta1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasconnection"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasinstance"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasinventory"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/dbaasprovider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
)

// checkAndEnableDBaaS check if DBaaSProvider CRD is deployed and start the DBaaS related controllers
func checkAndEnableDBaaS(mgr manager.Manager, logger *zap.Logger, config Config) (bool, error) {
	cfg := mgr.GetConfig()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return false, fmt.Errorf("unable to create clientset: %w", err)
	}

	dbaasInstalled, err := checkDBaaSCRDInstalled(clientset)
	if err != nil {
		return false, fmt.Errorf("unable to check OpenShift Database Access (DBaaS) Provider CRD: %w", err)
	}
	if dbaasInstalled {
		if err = (&dbaasprovider.DBaaSProviderReconciler{
			Client:    mgr.GetClient(),
			Scheme:    mgr.GetScheme(),
			Log:       logger.Named("controllers").Named("DBaaSProvider").Sugar(),
			Clientset: clientset,
		}).SetupWithManager(mgr); err != nil {
			return false, fmt.Errorf("unable to create DBaaSProvider controller: %w", err)
		}

		if err = (&atlasinventory.MongoDBAtlasInventoryReconciler{
			Client:          mgr.GetClient(),
			Log:             logger.Named("controllers").Named("MongoDBAtlasInventory").Sugar(),
			Scheme:          mgr.GetScheme(),
			AtlasDomain:     config.AtlasDomain,
			ResourceWatcher: watch.NewResourceWatcher(),
			GlobalAPISecret: config.GlobalAPISecret,
			EventRecorder:   mgr.GetEventRecorderFor("MongoDBAtlasInventory"),
		}).SetupWithManager(mgr); err != nil {
			return false, fmt.Errorf("unable to create MongoDBAtlasInventory controller: %w", err)
		}

		if err = (&atlasconnection.MongoDBAtlasConnectionReconciler{
			Client:          mgr.GetClient(),
			Clientset:       clientset,
			Log:             logger.Named("controllers").Named("MongoDBAtlasConnection").Sugar(),
			Scheme:          mgr.GetScheme(),
			AtlasDomain:     config.AtlasDomain,
			ResourceWatcher: watch.NewResourceWatcher(),
			GlobalAPISecret: config.GlobalAPISecret,
			EventRecorder:   mgr.GetEventRecorderFor("MongoDBAtlasConnection"),
		}).SetupWithManager(mgr); err != nil {
			return false, fmt.Errorf("unable to create MongoDBAtlasConnection controller: %w", err)
		}

		if err = (&atlasinstance.MongoDBAtlasInstanceReconciler{
			Client:          mgr.GetClient(),
			Clientset:       clientset,
			Log:             logger.Named("controllers").Named("MongoDBAtlasInstance").Sugar(),
			Scheme:          mgr.GetScheme(),
			AtlasDomain:     config.AtlasDomain,
			ResourceWatcher: watch.NewResourceWatcher(),
			GlobalAPISecret: config.GlobalAPISecret,
			EventRecorder:   mgr.GetEventRecorderFor("MongoDBAtlasInstance"),
		}).SetupWithManager(mgr); err != nil {
			return false, fmt.Errorf("unable to create MongoDBAtlasInstance controller: %w", err)
		}
	} else {
		return false, nil
	}
	return true, nil
}

// checkDBaaSCRDInstalled checks whether dbaas provider CRD, has been created yet
func checkDBaaSCRDInstalled(clientset kubernetes.Interface) (bool, error) {
	resources, err := clientset.Discovery().ServerResourcesForGroupVersion(dbaasv1beta1.GroupVersion.String())
	if err != nil {
		if apiErrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check DBaaSProvider CRD:%w", err)
	}
	for _, r := range resources.APIResources {
		if r.Kind == dbaasProviderKind {
			return true, nil
		}
	}
	return false, nil
}
