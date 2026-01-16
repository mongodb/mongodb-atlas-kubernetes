// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasorgsettings

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlrtbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/atlasorgsettings"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	mckpredicate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasorgsettings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasorgsettings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasorgsettings/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasorgsettings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasorgsettings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasorgsettings/finalizers,verbs=update

type serviceBuilderFunc func(*atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService

type AtlasOrgSettingsHandler struct {
	ctrlstate.StateHandler[akov2.AtlasOrgSettings]
	reconciler.AtlasReconciler
	serviceBuilder serviceBuilderFunc
}

func NewAtlasOrgSettingsReconciler(
	c cluster.Cluster,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef client.ObjectKey,
	reapplySupport bool,
) *ctrlstate.Reconciler[akov2.AtlasOrgSettings] {
	orgSettingsHandler := &AtlasOrgSettingsHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			AtlasProvider:   atlasProvider,
			Log:             logger.Named("controllers").Named("AtlasOrgSettings").Sugar(),
			GlobalSecretRef: globalSecretRef,
		},
		serviceBuilder: func(clientSet *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
			return atlasorgsettings.NewAtlasOrgSettingsService(clientSet.SdkClient20250312012.OrganizationsApi)
		},
	}
	return ctrlstate.NewStateReconciler(
		orgSettingsHandler,
		ctrlstate.WithCluster[akov2.AtlasOrgSettings](c),
		ctrlstate.WithReapplySupport[akov2.AtlasOrgSettings](reapplySupport),
	)
}

func (h *AtlasOrgSettingsHandler) For() (client.Object, builder.Predicates) {
	obj := &akov2.AtlasOrgSettings{}
	return obj, ctrlrtbuilder.WithPredicates(
		predicate.Or(
			mckpredicate.AnnotationChanged("mongodb.com/reapply-period"),
			predicate.GenerationChangedPredicate{},
		),
		mckpredicate.IgnoreDeletedPredicate[client.Object](),
	)
}

func (h *AtlasOrgSettingsHandler) SetupWithManager(mgr ctrl.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).
		Named("AtlasOrgSettings").
		For(h.For()).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(h.findSecretsForOrgSettings()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(defaultOptions).Complete(rec)
}

func (h *AtlasOrgSettingsHandler) findSecretsForOrgSettings() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasOrgSettingsBySecretsIndex,
		func() *akov2.AtlasOrgSettingsList { return &akov2.AtlasOrgSettingsList{} },
		indexer.AtlasOrgSettingsRequest, h.Client, h.Log)
}
