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
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	client.Client
	Log    zap.SugaredLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *AtlasProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasproject", req.NamespacedName)

	project := &mdbv1.AtlasProject{}
	if result := customresource.GetResource(r.Client, req, project, log); !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	ctx := workflow.NewContext(log)

	log.Infow("-> Starting AtlasProject reconciliation", "spec", project.Spec)

	if project.Spec.ConnectionSecret == nil {
		log.Error("So far the Connection Secret in AtlasProject is mandatory!")
		return reconcile.Result{}, nil
	}

	connection, result := atlas.ReadConnection(ctx, r.Client, "TODO!", project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		// merge result into ctx
		statushandler.Update(ctx.SetConditionFromResult(status.ProjectReadyType, result), r.Client, project)
		log.Debugf("returning %+v", result.ReconcileResult())
		return result.ReconcileResult(), nil
	}

	var projectID string
	if projectID, result = ensureProjectExists(ctx, connection, project); !result.IsOk() {
		statushandler.Update(ctx.SetConditionFromResult(status.ProjectReadyType, result), r.Client, project)
		return result.ReconcileResult(), nil
	}
	ctx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))

	// Updating the status with "projectReady = true" and "IPAccessListReady = false" (not as separate updates!)
	ctx.SetConditionTrue(status.ProjectReadyType)
	statushandler.Update(ctx.SetConditionFalse(status.IPAccessListReadyType), r.Client, project)

	// TODO projectAccessList

	return ctrl.Result{}, nil
}

func (r *AtlasProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasProject{}).
		WithEventFilter(watch.CommonPredicates()).
		Complete(r)
}
