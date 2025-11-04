// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thirdpartyintegration

import (
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=thirdpartyintegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=thirdpartyintegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=thirdpartyintegrations/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=thirdpartyintegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=thirdpartyintegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=thirdpartyintegrations/finalizers,verbs=update
type ThirdPartyIntegrationHandler struct {
	ctrlstate.StateHandler[v1.ThirdPartyIntegration]
	reconciler.AtlasReconciler
	handlerv20250312 *ThirdPartyIntegrationHandlerv20250312
}

func NewThirdPartyIntegrationReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.ThirdPartyIntegration] {
	// Create version-specific handlers

	handlerv20250312 := NewThirdPartyIntegrationHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("ThirdPartyIntegration-v20250312").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	thirdpartyintegrationHandler := &ThirdPartyIntegrationHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasThirdPartyIntegration").Sugar(),
		},
		handlerv20250312: handlerv20250312,
	}

	return ctrlstate.NewStateReconciler(thirdpartyintegrationHandler, ctrlstate.WithCluster[v1.ThirdPartyIntegration](c), ctrlstate.WithReapplySupport[v1.ThirdPartyIntegration](reapplySupport))
}
