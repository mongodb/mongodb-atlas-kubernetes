// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrations

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	mckpredicate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasthirdpartyintegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasthirdpartyintegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasthirdpartyintegrations/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasthirdpartyintegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasthirdpartyintegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasthirdpartyintegrations/finalizers,verbs=update

type serviceBuilderFunc func(*atlas.ClientSet) thirdpartyintegration.ThirdPartyIntegrationService

type AtlasThirdPartyIntegrationHandler struct {
	ctrlstate.StateHandler[akov2.AtlasThirdPartyIntegration]
	reconciler.AtlasReconciler
	deletionProtection bool
	serviceBuilder     serviceBuilderFunc
}

func NewAtlasThirdPartyIntegrationsReconciler(
	c cluster.Cluster,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
	globalSecretRef client.ObjectKey,
	reapplySupport bool,
) *ctrlstate.Reconciler[akov2.AtlasThirdPartyIntegration] {
	intHandler := &AtlasThirdPartyIntegrationHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			AtlasProvider:   atlasProvider,
			Log:             logger.Named("controllers").Named("AtlasThirdPartyIntegration").Sugar(),
			GlobalSecretRef: globalSecretRef,
		},
		deletionProtection: deletionProtection,
		serviceBuilder:     thirdpartyintegration.NewThirdPartyIntegrationServiceFromClientSet,
	}
	return ctrlstate.NewStateReconciler(
		intHandler,
		ctrlstate.WithCluster[akov2.AtlasThirdPartyIntegration](c),
		ctrlstate.WithReapplySupport[akov2.AtlasThirdPartyIntegration](reapplySupport),
	)
}

// For prepares the controller for its target Custom Resource; AtlasThirdPartyIntegration
func (r *AtlasThirdPartyIntegrationHandler) For() (client.Object, builder.Predicates) {
	obj := &akov2.AtlasThirdPartyIntegration{}
	return obj, ctrlrtbuilder.WithPredicates(
		predicate.Or(
			mckpredicate.AnnotationChanged("mongodb.com/reapply-period"),
			predicate.GenerationChangedPredicate{},
		),
		mckpredicate.IgnoreDeletedPredicate[client.Object](),
	)
}

func (h *AtlasThirdPartyIntegrationHandler) SetupWithManager(mgr ctrl.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).
		For(h.For()).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(h.integrationForProjectMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(h.integrationForSecretMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(defaultOptions).Complete(rec)
}

func (h *AtlasThirdPartyIntegrationHandler) integrationForProjectMapFunc() handler.MapFunc {
	return indexer.ProjectsIndexMapperFunc(
		string(indexer.AtlasThirdPartyIntegrationByProjectIndex),
		func() *akov2.AtlasThirdPartyIntegrationList { return &akov2.AtlasThirdPartyIntegrationList{} },
		indexer.AtlasThirdPartyIntegrationRequests,
		h.Client,
		h.Log,
	)
}

func (h *AtlasThirdPartyIntegrationHandler) integrationForSecretMapFunc() handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		secret, ok := obj.(*corev1.Secret)
		if !ok {
			h.Log.Warnf("watching Secret but got %T", obj)
			return nil
		}

		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasThirdPartyIntegrationBySecretsIndex,
				client.ObjectKeyFromObject(secret).String(),
			),
		}
		list1 := &akov2.AtlasThirdPartyIntegrationList{}
		err := h.Client.List(ctx, list1, listOpts)
		if err != nil {
			h.Log.Errorf("failed to list from indexer %s: %v",
				indexer.AtlasThirdPartyIntegrationBySecretsIndex, err)
			return nil
		}
		requests1 := indexer.AtlasThirdPartyIntegrationRequests(list1)

		listOpts = &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexer.AtlasThirdPartyIntegrationCredentialsIndex,
				client.ObjectKeyFromObject(secret).String(),
			),
		}
		list2 := &akov2.AtlasThirdPartyIntegrationList{}
		err = h.Client.List(ctx, list2, listOpts)
		if err != nil {
			h.Log.Errorf("failed to list from indexer %s: %v",
				indexer.AtlasThirdPartyIntegrationCredentialsIndex, err)
			return nil
		}
		requests2 := indexer.AtlasThirdPartyIntegrationRequests(list2)

		return append(requests1, requests2...)
	}
}

type reconcileRequest struct {
	ClientSet   *atlas.ClientSet
	Project     *project.Project
	Service     thirdpartyintegration.ThirdPartyIntegrationService
	integration *akov2.AtlasThirdPartyIntegration
}

func (h *AtlasThirdPartyIntegrationHandler) newReconcileRequest(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (*reconcileRequest, error) {
	req := reconcileRequest{}
	sdkClientSet, err := h.ResolveSDKClientSet(ctx, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve connection config: %w", err)
	}
	req.ClientSet = sdkClientSet
	req.Service = h.serviceBuilder(sdkClientSet)
	resolvedProject, err := h.ResolveProject(ctx, sdkClientSet.SdkClient20250312011, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch referenced project: %w", err)
	}
	req.Project = resolvedProject
	req.integration = integration
	return &req, nil
}
