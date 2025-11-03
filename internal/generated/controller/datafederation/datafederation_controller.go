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

package datafederation

import (
	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"
)

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=datafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=datafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=datafederations/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=datafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=datafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=datafederations/finalizers,verbs=update
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

type DataFederationHandler struct {
	ctrlstate.StateHandler[v1.DataFederation]
	reconciler.AtlasReconciler
	handlerv20250312 *DataFederationHandlerv20250312
}

func NewDataFederationReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.DataFederation] {
	// Create version-specific handlers

	handlerv20250312 := NewDataFederationHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("DataFederation-v20250312").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	datafederationHandler := &DataFederationHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		},
		handlerv20250312: handlerv20250312,
	}

	return ctrlstate.NewStateReconciler(datafederationHandler, ctrlstate.WithCluster[v1.DataFederation](c), ctrlstate.WithReapplySupport[v1.DataFederation](reapplySupport))
}
