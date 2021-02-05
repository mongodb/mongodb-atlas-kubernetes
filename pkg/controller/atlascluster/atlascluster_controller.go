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
	"errors"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasClusterReconciler reconciles an AtlasCluster object
type AtlasClusterReconciler struct {
	Client      client.Client
	Log         *zap.SugaredLogger
	Scheme      *runtime.Scheme
	AtlasDomain string
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

	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec, "status", cluster.Status)
	defer statushandler.Update(ctx, r.Client, cluster)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(cluster, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, result := atlas.ReadConnection(log, r.Client, "TODO!", project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		// merge result into ctx
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	c, result := r.ensureClusterState(log, connection, project, cluster)
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

func (r *AtlasClusterReconciler) readProjectResource(cluster *mdbv1.AtlasCluster, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(context.Background(), cluster.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasCluster{}).
		WithEventFilter(watch.CommonPredicates()).
		Watches(
			&source.Kind{Type: &mdbv1.AtlasCluster{}},
			&watch.DeleteEventHandler{Controller: r},
			builder.WithPredicates(watch.DeleteOnly()),
		).
		Complete(r)
}

// Delete implements a handler for the Delete event.
func (r *AtlasClusterReconciler) Delete(obj runtime.Object) error {
	cluster, ok := obj.(*mdbv1.AtlasCluster)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasCluster{}, obj)
		return nil
	}

	log := r.Log.With("atlascluster", kube.ObjectKeyFromObject(cluster))

	log.Infow("-> Starting AtlasCluster deletion", "spec", cluster.Spec)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(cluster, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	connection, result := atlas.ReadConnection(log, r.Client, "TODO!", project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		return errors.New("cannot read Atlas connection")
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	_, err = atlasClient.Clusters.Delete(context.Background(), project.Status.ID, cluster.Spec.Name)
	if err != nil {
		return fmt.Errorf("cannot delete Atlas cluster: %w", err)
	}

	log.Infow("Started Atlas cluster deletion process", "projectID", project.Status.ID, "clusterName", cluster.Name)

	return nil
}
