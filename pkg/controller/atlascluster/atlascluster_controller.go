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
	"time"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasClusterReconciler reconciles an AtlasCluster object
type AtlasClusterReconciler struct {
	Client           client.Client
	Log              *zap.SugaredLogger
	Scheme           *runtime.Scheme
	AtlasDomain      string
	GlobalAPISecret  client.ObjectKey
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasClusterReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlascluster", req.NamespacedName)

	cluster := &mdbv1.AtlasCluster{}
	result := customresource.PrepareResource(r.Client, req, cluster, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	ctx := customresource.MarkReconciliationStarted(r.Client, cluster, log)

	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec, "status", cluster.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, cluster)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(cluster, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Client = atlasClient

	c, result := r.ensureClusterState(ctx, project, cluster)
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	if csResult := r.ensureConnectionSecrets(ctx, project, c, cluster); !csResult.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, csResult)
		return csResult.ReconcileResult(), nil
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
	c, err := controller.New("AtlasCluster", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AtlasCluster & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasCluster{}}, &watch.EventHandlerWithDelete{Controller: r}, r.GlobalPredicates...)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements a handler for the Delete event.
func (r *AtlasClusterReconciler) Delete(e event.DeleteEvent) error {
	cluster, ok := e.Object.(*mdbv1.AtlasCluster)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasCluster{}, e.Object)
		return nil
	}

	log := r.Log.With("atlascluster", kube.ObjectKeyFromObject(cluster))

	log.Infow("-> Starting AtlasCluster deletion", "spec", cluster.Spec)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(cluster, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	log = log.With("projectID", project.Status.ID, "clusterName", cluster.Spec.Name)

	if customresource.ResourceShouldBeLeftInAtlas(cluster) {
		log.Infof("Not removing Atlas Cluster from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
	} else if err := r.deleteClusterFromAtlas(cluster, project, log); err != nil {
		log.Error("Failed to remove cluster from Atlas: %s", err)
	}

	// We always remove the connection secrets even if the cluster is not removed from Atlas
	secrets, err := connectionsecret.ListByClusterName(r.Client, cluster.Namespace, project.ID(), cluster.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(context.Background(), &secrets[i]); err != nil {
			log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
	}

	return nil
}

func (r *AtlasClusterReconciler) deleteClusterFromAtlas(cluster *mdbv1.AtlasCluster, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	go func() {
		timeout := time.Now().Add(workflow.DefaultTimeout)

		for time.Now().Before(timeout) {
			_, err = atlasClient.Clusters.Delete(context.Background(), project.Status.ID, cluster.Spec.Name)
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && apiError.ErrorCode == atlas.ClusterNotFound {
				log.Info("Cluster doesn't exist or is already deleted")
				return
			}

			if err != nil {
				log.Errorw("Cannot delete Atlas cluster", "error", err)
				time.Sleep(workflow.DefaultRetry)
				continue
			}

			log.Info("Started Atlas cluster deletion process")
			return
		}

		log.Error("Failed to delete Atlas cluster in time")
	}()
	return nil
}
