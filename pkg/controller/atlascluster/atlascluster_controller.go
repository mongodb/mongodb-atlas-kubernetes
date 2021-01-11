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

package atlascluster

import (
	"context"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// AtlasClusterReconciler reconciles a AtlasCluster object
type AtlasClusterReconciler struct {
	client.Client
	Log    *zap.SugaredLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters/status,verbs=get;update;patch

func (r *AtlasClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.With("atlascluster", req.NamespacedName)

	cluster := &mdbv1.AtlasCluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		log.Error(err, "Failed to read AtlasCluster")
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	log = log.With("clusterName", cluster.Spec.Name)

	project := &mdbv1.AtlasProject{}
	if err := r.Get(ctx, req.NamespacedName, project); err != nil {
		log.Error(err, "Failed to read Project from AtlasCluster")
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec)

	wctx := workflow.NewContext(log)
	defer statushandler.Update(wctx, r, cluster)

	connection, result := atlas.ReadConnection(wctx, r, "TODO!", project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		// merge result into ctx
		wctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	c, result := ensureClusterState(wctx, connection, project, cluster)
	if c.StateName != "" {
		wctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		wctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	wctx.
		SetConditionTrue(status.ClusterReadyType).
		EnsureStatusOption(status.AtlasClusterMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasClusterConnectionStringsOption(c.ConnectionStrings)).
		EnsureStatusOption(status.AtlasClusterMongoURIUpdatedOption(c.MongoURIUpdated))

	return result.ReconcileResult(), nil
}

func (r *AtlasClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasCluster{}).
		Complete(r)
}
