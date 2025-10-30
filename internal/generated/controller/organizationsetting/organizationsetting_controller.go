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

package organizationsetting

import (
	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=organizationsettings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=organizationsettings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=organizationsettings/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=organizationsettings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=organizationsettings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=organizationsettings/finalizers,verbs=update
type OrganizationSettingHandler struct {
	ctrlstate.StateHandler[v1.OrganizationSetting]
	reconciler.AtlasReconciler
	handlerv20250312 *OrganizationSettingHandlerv20250312
}

func NewOrganizationSettingReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.OrganizationSetting] {
	// Create version-specific handlers

	handlerv20250312 := NewOrganizationSettingHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("OrganizationSetting-v20250312").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	organizationsettingHandler := &OrganizationSettingHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasOrganizationSetting").Sugar(),
		},
		handlerv20250312: handlerv20250312,
	}

	return ctrlstate.NewStateReconciler(organizationsettingHandler, ctrlstate.WithCluster[v1.OrganizationSetting](c), ctrlstate.WithReapplySupport[v1.OrganizationSetting](reapplySupport))
}
