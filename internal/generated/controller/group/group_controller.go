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

package group

import (
	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"
)

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=groups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=groups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=groups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=groups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=groups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=groups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

type GroupHandler struct {
	ctrlstate.StateHandler[v1.Group]
	reconciler.AtlasReconciler
	handlerv20250312 *GroupHandlerv20250312
}

func NewGroupReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.Group] {
	// Create version-specific handlers

	handlerv20250312 := NewGroupHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("Group-v20250312").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	groupHandler := &GroupHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasGroup").Sugar(),
		},
		handlerv20250312: handlerv20250312,
	}

	return ctrlstate.NewStateReconciler(groupHandler, ctrlstate.WithCluster[v1.Group](c), ctrlstate.WithReapplySupport[v1.Group](reapplySupport))
}
