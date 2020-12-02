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
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	client.Client
	Log    zap.SugaredLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mongodb.com.mongodb.com,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mongodb.com.mongodb.com,resources=atlasprojects/status,verbs=get;update;patch

func (r *AtlasProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.With("atlasproject", req.NamespacedName)

	project := &mdbv1.AtlasProject{}
	if err := r.Client.Get(ctx, kube.ObjectKey(req.Namespace, req.Name), project); err != nil {
		// TODO make generic (update status, log message)
		log.Error(err, "Failed to read the AtlasProject")
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	if project.Spec.ConnectionSecret == nil {
		log.Error("So far the Connection Secret in AtlasProject is mandatory!")
		return reconcile.Result{}, nil
	}

	connection, err := atlas.ReadConnection(r.Client, "TODO!", project.ConnectionSecretObjectKey(), log)
	if err != nil {
		log.Error(err, "Failed to read Atlas Connection details")
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	if err := ensureProjectExists(connection, project, log); err != nil {
		log.Error(err, "Failed to read the AtlasProject")
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	// TODO projectAccessList

	return ctrl.Result{}, nil
}

func (r *AtlasProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasProject{}).
		Complete(r)
}
