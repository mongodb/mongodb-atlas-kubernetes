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

package databaseuser

import (
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312006/admin"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	translate "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=databaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=databaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=databaseusers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=databaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=databaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=databaseusers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

type Handler struct {
	ctrlstate.StateHandler[akov2generated.DatabaseUser]
	reconciler.AtlasReconciler
	deletionProtection bool
	predicates         []predicate.Predicate
	handlerv20250312   ctrlstate.VersionedHandlerFunc[v20250312sdk.APIClient, akov2generated.DatabaseUser]
}

func NewDatabaseUserReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, deletionProtection bool, reapplySupport bool, predicates []predicate.Predicate) *ctrlstate.Reconciler[akov2generated.DatabaseUser] {
	// Create main handler dispatcher
	databaseuserHandler := &Handler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		},
		deletionProtection: deletionProtection,
		handlerv20250312:   handlerv20250312Func,
		predicates:         predicates,
	}

	return ctrlstate.NewStateReconciler(databaseuserHandler, ctrlstate.WithCluster[akov2generated.DatabaseUser](c), ctrlstate.WithReapplySupport[akov2generated.DatabaseUser](reapplySupport))
}
func handlerv20250312Func(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translatorRequest *translate.Request, deletionProtection bool) ctrlstate.StateHandler[akov2generated.DatabaseUser] {
	return NewHandlerv20250312(kubeClient, atlasClient, translatorRequest, deletionProtection)
}
