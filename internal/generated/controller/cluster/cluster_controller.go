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

package cluster

import (
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=clusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=clusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=clusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=clusters/finalizers,verbs=update
type ClusterHandler struct {
	ctrlstate.StateHandler[v1.Cluster]
	reconciler.AtlasReconciler
	handlerv20250312 *ClusterHandlerv20250312
	handlerv20250313 *ClusterHandlerv20250313
}

func NewClusterReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.Cluster] {
	// Create version-specific handlers

	handlerv20250312 := NewClusterHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("Cluster-v20250312").Sugar(), globalSecretRef)
	handlerv20250313 := NewClusterHandlerv20250313(atlasProvider, c.GetClient(), logger.Named("controllers").Named("Cluster-v20250313").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	clusterHandler := &ClusterHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasCluster").Sugar(),
		},
		handlerv20250312: handlerv20250312,
		handlerv20250313: handlerv20250313,
	}

	return ctrlstate.NewStateReconciler(clusterHandler, ctrlstate.WithCluster[v1.Cluster](c), ctrlstate.WithReapplySupport[v1.Cluster](reapplySupport))
}
