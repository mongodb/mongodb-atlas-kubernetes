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
	"reflect"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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
		atlasClient, _, err := r.AtlasProvider.Client(ctx, project.ConnectionSecretObjectKey(), log)
		if err != nil {
			return r.terminate(err, nil)
		}
		return r.reconcile(ctx, atlasClient, project)
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
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *AtlasProjectExperimentalReconciler) findExperimentalProjectForStableProject(ctx context.Context, obj client.Object) []reconcile.Request {
	stableProject, ok := obj.(*akov2.AtlasProject)
	if !ok {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      stableProject.Name,
			Namespace: stableProject.Namespace,
		},
	}}
}

func (r *AtlasProjectExperimentalReconciler) reconcile(ctx context.Context, atlasClient *mongodbatlas.Client, experimentalProject *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	var stableProject akov2.AtlasProject
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: experimentalProject.GetNamespace(), Name: experimentalProject.GetName()}, &stableProject)

	switch {
	case errors.IsNotFound(err):
		// no stable project is present, create one
		return r.create(ctx, experimentalProject)
	case err != nil:
		// some error occurred
		return r.terminate(err, experimentalProject)
	}

	var idle ctrl.Result

	// one stable project is present, update it
	result, err := r.update(ctx, experimentalProject, &stableProject)
	if err != nil {
		return r.terminate(err, experimentalProject)
	}
	if result != idle {
		return result, nil
	}

	result, err = r.reconcileAudit(ctx, atlasClient, experimentalProject)
	if err != nil {
		return r.terminate(err, experimentalProject)
	}
	if result != idle {
		return result, nil
	}

	if err := r.Client.SubResource("status").Update(ctx, experimentalProject, &client.SubResourceUpdateOptions{}); err != nil {
		return r.terminate(err, experimentalProject)
	}

	return r.idle(experimentalProject)
}

func (r *AtlasProjectExperimentalReconciler) create(ctx context.Context, experimentalProject *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	var stableProject akov2.AtlasProject

	stableProject.Spec = experimentalProject.Spec.AtlasProjectSpec
	stableProject.Namespace = experimentalProject.Namespace
	stableProject.Name = experimentalProject.Name
	stableProject.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			Kind:       experimentalProject.Kind,
			APIVersion: experimentalProject.APIVersion,
			Name:       experimentalProject.Name,
			UID:        experimentalProject.UID,
		},
	}

	if err := r.Client.Create(ctx, &stableProject); err != nil {
		return r.terminate(err, experimentalProject)
	}

	return r.progress(experimentalProject)
}

func (r *AtlasProjectExperimentalReconciler) update(ctx context.Context, experimentalProject *akov1alpha1.AtlasProject, stableProject *akov2.AtlasProject) (ctrl.Result, error) {
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: stableProject.GetNamespace(), Name: stableProject.GetName()}, stableProject); err != nil {
		return r.terminate(err, experimentalProject)
	}

	// first, update the stable project spec, if needed
	if !reflect.DeepEqual(stableProject.Spec, experimentalProject.Spec.AtlasProjectSpec) {
		stableProject.Spec = experimentalProject.Spec.AtlasProjectSpec
		if err := r.Client.Update(ctx, stableProject); err != nil {
			return r.terminate(err, experimentalProject)
		}
	}

	// then, always update the experimental project status based on the stable project status
	experimentalProject.Status = stableProject.Status
	return r.idle(experimentalProject)
}

func (r *AtlasProjectExperimentalReconciler) reconcileAudit(ctx context.Context, atlasClient *mongodbatlas.Client, project *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	if project.Spec.AuditRef.Name == "" {
		return r.empty(project)
	}

	if project.Status.ID == "" {
		return r.progress(project)
	}

	var audit akov1alpha1.AtlasAuditing
	if err := r.Client.Get(ctx, *project.Spec.AuditRef.GetObject(project.Namespace), &audit); err != nil {
		return r.terminate(err, project)
	}

	atlasAuditing := &mongodbatlas.Auditing{
		AuditAuthorizationSuccess: pointer.MakePtr(audit.Spec.AuditAuthorizationSuccess),
		Enabled:                   pointer.MakePtr(audit.Spec.Enabled),
	}
	if len(audit.Spec.AuditFilter.Raw) > 0 {
		atlasAuditing.AuditFilter = string(audit.Spec.AuditFilter.Raw)
	}

	_, _, err := atlasClient.Auditing.Configure(ctx, project.Status.ID, atlasAuditing)
	if err != nil {
		return r.terminate(err, project)
	}

	return r.idle(project)
}

func (r *AtlasProjectExperimentalReconciler) progress(project *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	if project != nil {
		project.Status.Conditions = replaceStatusCondition(project.Status.Conditions, api.Condition{
			Type:               "AtlasAuditingReady",
			Status:             "False",
			Reason:             "InProgress",
			LastTransitionTime: metav1.Now(),
		})
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *AtlasProjectExperimentalReconciler) empty(project *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	if project != nil {
		project.Status.Conditions = replaceStatusCondition(project.Status.Conditions)
	}
	return ctrl.Result{}, nil
}

func (r *AtlasProjectExperimentalReconciler) idle(project *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	if project != nil {
		project.Status.Conditions = replaceStatusCondition(project.Status.Conditions, api.Condition{
			Type:               "AtlasAuditingReady",
			Status:             "True",
			LastTransitionTime: metav1.Now(),
		})
	}

	return ctrl.Result{}, nil
}

func (r *AtlasProjectExperimentalReconciler) terminate(err error, project *akov1alpha1.AtlasProject) (ctrl.Result, error) {
	if project != nil {
		project.Status.Conditions = replaceStatusCondition(project.Status.Conditions, api.Condition{
			Type:               "AtlasAuditingReady",
			Status:             "False",
			Reason:             "Error",
			LastTransitionTime: metav1.Now(),
			Message:            err.Error(),
		})
	}

	return ctrl.Result{}, err
}

func replaceStatusCondition(conditions []api.Condition, condition ...api.Condition) []api.Condition {
	result := make([]api.Condition, 0, len(conditions))
	for _, c := range conditions {
		if c.Type != "AtlasAuditingReady" {
			result = append(result, c)
		}
	}
	return append(result, condition...)
}
