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

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
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
	log := r.Log.With("atlascluster", req.NamespacedName)

	cluster := &mdbv1.AtlasCluster{}
	result := customresource.PrepareResource(r.Client, req, cluster, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	ctx := customresource.MarkReconciliationStarted(r.Client, cluster, log)

	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec)
	defer statushandler.Update(ctx, r, cluster)

	project := &mdbv1.AtlasProject{}
	if result := readProjectResource(r, cluster, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, result := atlas.ReadConnection(ctx, r, "TODO!", project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		// merge result into ctx
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	c, result := ensureClusterState(ctx, connection, project, cluster)
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	ctx.
		SetConditionTrue(status.ClusterReadyType).
		EnsureStatusOption(status.AtlasClusterMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasClusterConnectionStringsOption(c.ConnectionStrings)).
		EnsureStatusOption(status.AtlasClusterMongoURIUpdatedOption(c.MongoURIUpdated))

	ctx.SetConditionTrue(status.ReadyType)
	return result.ReconcileResult(), nil
}

func readProjectResource(r *AtlasClusterReconciler, cluster *mdbv1.AtlasCluster, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(context.Background(), cluster.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasCluster{}).
		Complete(r)
}
