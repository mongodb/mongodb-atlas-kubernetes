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
	"fmt"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

const (
	// crdVersion of the controler
	crdVersion = "v1"
)

var (
	// sdkVersions supported by this controller
	sdkVersions = []string{"v20250312"}
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

type Handler struct {
	ctrlstate.StateHandler[akov2generated.Group]
	reconciler.AtlasReconciler
	deletionProtection bool
	predicates         []predicate.Predicate
	translators        map[string]crapi.Translator
	handlerv20250312   ctrlstate.VersionedHandlerFunc[v20250312sdk.APIClient, akov2generated.Group]
}

func NewGroupReconciler(
	c cluster.Cluster,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef client.ObjectKey,
	deletionProtection bool,
	reapplySupport bool,
	predicates []predicate.Predicate) (*ctrlstate.Reconciler[akov2generated.Group], error) {
	crd, err := crds.EmbeddedCRD("Group")
	if err != nil {
		return nil, fmt.Errorf("failed to read CRD for Group: %w", err)
	}
	translators, err := crapi.NewPerVersionTranslators(c.GetScheme(), crd, crdVersion, sdkVersions...)
	if err != nil {
		return nil, fmt.Errorf("failed to get translator set for Group: %w", err)
	}
	// Create main handler dispatcher
	groupHandler := &Handler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasGroup").Sugar(),
		},
		deletionProtection: deletionProtection,
		handlerv20250312:   handlerv20250312Func,
		predicates:         predicates,
		translators:        translators,
	}

	return ctrlstate.NewStateReconciler(groupHandler, ctrlstate.WithCluster[akov2generated.Group](c), ctrlstate.WithReapplySupport[akov2generated.Group](reapplySupport)), nil
}
func handlerv20250312Func(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translator crapi.Translator, deletionProtection bool) ctrlstate.StateHandler[akov2generated.Group] {
	return NewHandlerv20250312(kubeClient, atlasClient, translator, deletionProtection)
}
