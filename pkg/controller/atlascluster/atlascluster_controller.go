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
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/validate"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
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
	watch.ResourceWatcher
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
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupschedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupschedules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackupschedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackupschedules/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackuppolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackuppolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackuppolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackuppolicies/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasClusterReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlascluster", req.NamespacedName)

	cluster := &mdbv1.AtlasCluster{}
	result := customresource.PrepareResource(r.Client, req, cluster, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if shouldSkip := customresource.ReconciliationShouldBeSkipped(cluster); shouldSkip {
		log.Infow(fmt.Sprintf("-> Skipping AtlasCluster reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", cluster.Spec)
		return workflow.OK().ReconcileResult(), nil
	}

	ctx := customresource.MarkReconciliationStarted(r.Client, cluster, log)
	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec, "status", cluster.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, cluster)

	if err := validate.ClusterSpec(cluster.Spec); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	ctx.SetConditionTrue(status.ValidationSucceeded)

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

	handleCluster := r.selectClusterHandler(cluster)
	if result, _ := handleCluster(ctx, project, cluster, req); !result.IsOk() {
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}

	if !cluster.IsServerless() {
		if result := r.handleAdvancedOptions(ctx, project, cluster); !result.IsOk() {
			ctx.SetConditionFromResult(status.ClusterReadyType, result)
			return result.ReconcileResult(), nil
		}
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasClusterReconciler) selectClusterHandler(cluster *mdbv1.AtlasCluster) clusterHandlerFunc {
	if cluster.IsAdvancedCluster() {
		return r.handleAdvancedCluster
	}
	if cluster.IsServerless() {
		return r.handleServerlessInstance
	}
	return r.handleRegularCluster
}

func (r *AtlasClusterReconciler) handleClusterBackupSchedule(ctx *workflow.Context, c *mdbv1.AtlasCluster, projectID, cName string, backupEnabled bool, req ctrl.Request) error {
	if c.Spec.BackupScheduleRef.Name == "" && c.Spec.BackupScheduleRef.Namespace == "" {
		r.Log.Debug("no backup schedule configured for the cluster")
		return nil
	}

	if !backupEnabled {
		return fmt.Errorf("can not proceed with backup schedule. Backups are not enabled for cluster %v", c.ClusterName)
	}

	resourcesToWatch := []watch.WatchedObject{}

	// Process backup schedule
	bSchedule := &mdbv1.AtlasBackupSchedule{}
	bKey := types.NamespacedName{Namespace: c.Spec.BackupScheduleRef.Namespace, Name: c.Spec.BackupScheduleRef.Name}
	err := r.Client.Get(context.Background(), bKey, bSchedule)
	if err != nil {
		return fmt.Errorf("%v backupschedule resource is not found. e: %w", c.Spec.BackupScheduleRef, err)
	}
	resourcesToWatch = append(resourcesToWatch, watch.WatchedObject{ResourceKind: bSchedule.Kind, Resource: bKey})

	// Process backup policy for the schedule
	bPolicy := &mdbv1.AtlasBackupPolicy{}
	pKey := types.NamespacedName{Namespace: bSchedule.Spec.PolicyRef.Namespace, Name: bSchedule.Spec.PolicyRef.Name}
	err = r.Client.Get(context.Background(), pKey, bPolicy)
	if err != nil {
		return fmt.Errorf("unable to get backuppolicy resource %s/%s. e: %w", bSchedule.Spec.PolicyRef.Namespace, bSchedule.Spec.PolicyRef.Name, err)
	}
	resourcesToWatch = append(resourcesToWatch, watch.WatchedObject{ResourceKind: bPolicy.Kind, Resource: pKey})

	// Create new backup schedule
	r.Log.Infof("updating backupschedule for the atlas cluster: %v", cName)
	apiScheduleRes := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
		ClusterName:           cName,
		ReferenceHourOfDay:    &bSchedule.Spec.ReferenceHourOfDay,
		ReferenceMinuteOfHour: &bSchedule.Spec.ReferenceMinuteOfHour,
		RestoreWindowDays:     &bSchedule.Spec.RestoreWindowDays,
		UpdateSnapshots:       &bSchedule.Spec.UpdateSnapshots,
		Policies:              nil,
	}

	// No matter what happens we should add watchers to both schedule and policy
	defer func() {
		r.EnsureMultiplesResourcesAreWatched(req.NamespacedName, r.Log, resourcesToWatch...)
		r.Log.Debugf("watched backup schedule resources: %v\r\n", r.WatchedResources)
	}()

	apiPolicy := mongodbatlas.Policy{}

	for _, bpItem := range bPolicy.Spec.Items {
		apiPolicy.PolicyItems = append(apiPolicy.PolicyItems, mongodbatlas.PolicyItem{
			FrequencyInterval: bpItem.FrequencyInterval,
			FrequencyType:     strings.ToLower(bpItem.FrequencyType),
			RetentionValue:    bpItem.RetentionValue,
			RetentionUnit:     strings.ToLower(bpItem.RetentionUnit),
		})
	}
	apiScheduleRes.Policies = []mongodbatlas.Policy{apiPolicy}

	currentSchedule, _, err := ctx.Client.CloudProviderSnapshotBackupPolicies.Delete(context.Background(), projectID, cName)
	if err != nil {
		r.Log.Debugf("unable to delete current backup policy for project: %v:%v, %v", projectID, cName, err)
	}
	r.Log.Debugf("successfully deleted backup policy. Default schedule received: %v", currentSchedule)

	apiScheduleRes.ClusterID = currentSchedule.ClusterID
	// There is only one policy, always
	apiScheduleRes.Policies[0].ID = currentSchedule.Policies[0].ID

	r.Log.Debugf("applying backupschedule policy: %v", *apiScheduleRes)
	if _, _, err := ctx.Client.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, cName, apiScheduleRes); err != nil {
		return fmt.Errorf("unable to create backupschedule %v. e: %w", bKey, err)
	}
	r.Log.Infof("successfully updated backupschedule for cluster %v", cName)
	return nil
}

// handleAdvancedCluster ensures the state of the cluster using the Advanced Cluster API
func (r *AtlasClusterReconciler) handleAdvancedCluster(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, req reconcile.Request) (workflow.Result, error) {
	c, result := r.ensureAdvancedClusterState(ctx, project, cluster)
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	if err := r.handleClusterBackupSchedule(ctx, cluster, project.ID(), c.Name, *c.BackupEnabled, req); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(ctx, project, c.Name, c.ConnectionStrings, cluster); !csResult.IsOk() {
		return csResult, nil
	}

	ctx.
		SetConditionTrue(status.ClusterReadyType).
		EnsureStatusOption(status.AtlasClusterMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasClusterConnectionStringsOption(c.ConnectionStrings))

	ctx.SetConditionTrue(status.ReadyType)
	return result, nil
}

// handleServerlessInstance ensures the state of the serverless instance using the serverless API
func (r *AtlasClusterReconciler) handleServerlessInstance(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, req reconcile.Request) (workflow.Result, error) {
	c, result := ensureServerlessInstanceState(ctx, project, cluster.Spec.ServerlessSpec)
	return r.ensureConnectionSecretsAndSetStatusOptions(ctx, project, cluster, result, c)
}

// handleRegularCluster ensures the state of the cluster using the Regular Cluster API
func (r *AtlasClusterReconciler) handleRegularCluster(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, req reconcile.Request) (workflow.Result, error) {
	c, result := ensureClusterState(ctx, project, cluster)
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	if err := r.handleClusterBackupSchedule(ctx, cluster, project.ID(), c.Name, *c.ProviderBackupEnabled || *c.BackupEnabled, req); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result, nil
	}
	return r.ensureConnectionSecretsAndSetStatusOptions(ctx, project, cluster, result, c)
}

// ensureConnectionSecretsAndSetStatusOptions creates the relevant connection secrets and sets
// status options to the given context. This function can be used for regular clusters and serverless instances
func (r *AtlasClusterReconciler) ensureConnectionSecretsAndSetStatusOptions(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, result workflow.Result, c *mongodbatlas.Cluster) (workflow.Result, error) {
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasClusterStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(ctx, project, c.Name, c.ConnectionStrings, cluster); !csResult.IsOk() {
		return csResult, nil
	}

	ctx.
		SetConditionTrue(status.ClusterReadyType).
		EnsureStatusOption(status.AtlasClusterMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasClusterConnectionStringsOption(c.ConnectionStrings)).
		EnsureStatusOption(status.AtlasClusterMongoURIUpdatedOption(c.MongoURIUpdated))

	ctx.SetConditionTrue(status.ReadyType)
	return result, nil
}

func (r *AtlasClusterReconciler) handleAdvancedOptions(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster) workflow.Result {
	clusterName := cluster.GetClusterName()
	atlasArgs, _, err := ctx.Client.Clusters.GetProcessArgs(context.Background(), project.Status.ID, clusterName)
	if err != nil {
		return workflow.Terminate(workflow.Internal, "cannot get process args")
	}

	if cluster.Spec.ProcessArgs == nil {
		return workflow.OK()
	}

	if !cluster.Spec.ProcessArgs.IsEqual(atlasArgs) {
		options := mongodbatlas.ProcessArgs(*cluster.Spec.ProcessArgs)
		args, resp, err := ctx.Client.Clusters.UpdateProcessArgs(context.Background(), project.Status.ID, clusterName, &options)
		ctx.Log.Debugw("ProcessArgs Update", "args", args, "resp", resp.Body, "err", err)
		if err != nil {
			return workflow.Terminate(workflow.Internal, "cannot update process args")
		}

		workflow.InProgress(workflow.ClusterAdvancedOptionsAreNotReady, "cluster Advanced Configuration Options are being updated")
	}

	return workflow.OK()
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

	// Watch for Backup schedules
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupSchedule{}}, watch.NewBackupScheduleHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	// Watch for Backup policies
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupPolicy{}}, watch.NewBackupPolicyHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	return nil
}

// Delete implements a handler for the Delete event.
func (r *AtlasClusterReconciler) Delete(e event.DeleteEvent) error {
	// TODO: Add deletion for AtlasBackupSchedule and AtlasBackupPolicy
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

	log = log.With("projectID", project.Status.ID, "clusterName", cluster.GetClusterName())

	if customresource.ResourceShouldBeLeftInAtlas(cluster) {
		log.Infof("Not removing Atlas Cluster from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
	} else if err := r.deleteClusterFromAtlas(cluster, project, log); err != nil {
		log.Error("Failed to remove cluster from Atlas: %s", err)
	}

	// We always remove the connection secrets even if the cluster is not removed from Atlas
	secrets, err := connectionsecret.ListByClusterName(r.Client, cluster.Namespace, project.ID(), cluster.GetClusterName())
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(context.Background(), &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
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
			deleteClusterFunc := atlasClient.Clusters.Delete
			if cluster.Spec.AdvancedClusterSpec != nil {
				deleteClusterFunc = atlasClient.AdvancedClusters.Delete
			}
			if cluster.IsServerless() {
				deleteClusterFunc = atlasClient.ServerlessInstances.Delete
			}

			_, err = deleteClusterFunc(context.Background(), project.Status.ID, cluster.GetClusterName())

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

type clusterHandlerFunc func(ctx *workflow.Context, project *mdbv1.AtlasProject, cluster *mdbv1.AtlasCluster, req reconcile.Request) (workflow.Result, error)
