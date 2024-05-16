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
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/labels"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	akov1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
)

const (
	ExperimentalProject = "atlas.mongodb.com/experimentalproject"
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

	project := &akov1alpha1.AtlasProject{}
	result := customresource.PrepareResource(ctx, r.Client, req, project, log)
	if result.IsOk() {
		return r.reconcile(ctx, project)
	}

	return result.ReconcileResult(), nil
}

func (r *AtlasProjectExperimentalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasExperimentalProject").
		For(&akov1alpha1.AtlasProject{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.findExperimentalProjectForStableProject),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Complete(r)
}

func (r *AtlasProjectExperimentalReconciler) findExperimentalProjectForStableProject(ctx context.Context, obj client.Object) []reconcile.Request {
	stableProject, ok := obj.(*akov2.AtlasProject)
	if !ok {
		return nil
	}
	value, ok := stableProject.Labels[ExperimentalProject]
	if !ok {
		return nil
	}

	name := strings.TrimPrefix(strings.TrimPrefix(value, stableProject.Namespace), "-")
	if len(name) == 0 {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: stableProject.Namespace,
		},
	}}
}

func (r *AtlasProjectExperimentalReconciler) reconcile(ctx context.Context, experimentalProject *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	var stableProjects akov2.AtlasProjectList
	if err := r.Client.List(ctx, &stableProjects, &client.ListOptions{
		LabelSelector: mustParseLabelSelector(ExperimentalProject + "=" + experimentalProject.Namespace + "-" + experimentalProject.Name),
	}); err != nil {
		return r.terminate(err)
	}

	switch {
	// no stable project is present, create one
	case len(stableProjects.Items) == 0:
		return r.create(ctx, experimentalProject)
	// one stable project is present, update it
	case len(stableProjects.Items) == 1:
		return r.update(ctx, experimentalProject, stableProjects.Items[0].DeepCopy())
	}

	return r.terminate(fmt.Errorf("expected none or one AtlasProject but found %v", len(stableProjects.Items)))
}

func (r *AtlasProjectExperimentalReconciler) create(ctx context.Context, experimentalProject *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	var stableProject akov2.AtlasProject

	// overwrite spec
	stableProject.Spec = experimentalProject.Spec.AtlasProjectSpec
	stableProject.Namespace = experimentalProject.Namespace
	stableProject.GenerateName = experimentalProject.Name
	stableProject.Labels = map[string]string{
		ExperimentalProject: experimentalProject.Namespace + "-" + experimentalProject.Name,
	}

	if err := r.Client.Create(ctx, &stableProject); err != nil {
		return r.terminate(err)
	}

	return r.progress()
}

func (r *AtlasProjectExperimentalReconciler) update(ctx context.Context, experimentalProject *akov1alpha1.AtlasProject, stableProject *akov2.AtlasProject) (ctrl.Result, error) {
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: stableProject.GetNamespace(), Name: stableProject.GetName()}, stableProject); err != nil {
		return r.terminate(err)
	}

	// first, update the stable project spec, if needed
	if !reflect.DeepEqual(stableProject.Spec, experimentalProject.Spec.AtlasProjectSpec) {
		stableProject.Spec = experimentalProject.Spec.AtlasProjectSpec
		if err := r.Client.Update(ctx, stableProject); err != nil {
			return r.terminate(err)
		}
	}

	// then, always update the experimental project status based on the stable project status
	experimentalProject.Status = stableProject.Status
	if err := r.Client.SubResource("status").Update(ctx, experimentalProject, &client.SubResourceUpdateOptions{}); err != nil {
		return r.terminate(err)
	}

	return r.idle()
}

func (r *AtlasProjectExperimentalReconciler) progress() (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *AtlasProjectExperimentalReconciler) idle() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *AtlasProjectExperimentalReconciler) terminate(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func mustParseLabelSelector(label string) labels.Selector {
	sel, err := labels.Parse(label)
	if err != nil {
		panic(err)
	}
	return sel
}
