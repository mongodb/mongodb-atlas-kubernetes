/*
Copyright 2020 MongoDB.

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

package atlasproject

import (
	"context"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	akov1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
)

type AtlasProjectExperimentalReconciler struct {
	Client client.Client
	watch.DeprecatedResourceWatcher
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

func (r *AtlasProjectExperimentalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlas_experimental_project", req.NamespacedName)

	projectexperimental := &v1alpha1.AtlasProject{}
	result := customresource.PrepareResource(ctx, r.Client, req, projectexperimental, log)
	if result.IsOk() {
		log.Info("reconciling experimental AtlasProject")
		return result.ReconcileResult(), nil
	}

	return result.ReconcileResult(), nil
}

func (r *AtlasProjectExperimentalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasExperimentalProject").
		For(&akov1alpha1.AtlasProject{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(&corev1.Secret{}, watch.NewSecretHandler(&r.DeprecatedResourceWatcher)).
		Watches(&akov2.AtlasTeam{}, watch.NewAtlasTeamHandler(&r.DeprecatedResourceWatcher)).
		Complete(r)
}
