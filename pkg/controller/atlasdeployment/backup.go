package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"

	"go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/sync/errgroup"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func (r *AtlasDeploymentReconciler) ensureBackupScheduleAndPolicy(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	backupScheduleRef types.NamespacedName,
	enabled bool,
	requestNamespacedName client.ObjectKey,
) error {
	if backupScheduleRef.Name == "" || backupScheduleRef.Namespace == "" {
		r.Log.Debug("no backup schedule configured for the deployment")

		err := r.garbageCollectBackupResource(ctx, clusterName)
		if err != nil {
			return err
		}

		return nil
	}

	if !enabled {
		return fmt.Errorf("can not proceed with backup configuration. Backups are not enabled for cluster %s", clusterName)
	}

	resourcesToWatch := []watch.WatchedObject{}
	defer func() {
		r.EnsureMultiplesResourcesAreWatched(requestNamespacedName, r.Log, resourcesToWatch...)
		r.Log.Debugf("watched backup schedule and policy resources: %v\r\n", r.WatchedResources)
	}()

	bSchedule, err := r.ensureBackupSchedule(ctx, service, backupScheduleRef, clusterName, &resourcesToWatch)
	if err != nil {
		return err
	}

	bPolicy, err := r.ensureBackupPolicy(ctx, service, *bSchedule.Spec.PolicyRef.GetObject(bSchedule.Namespace), &resourcesToWatch)
	if err != nil {
		return err
	}

	return r.updateBackupScheduleAndPolicy(ctx, service, projectID, clusterName, bSchedule, bPolicy)
}

func (r *AtlasDeploymentReconciler) ensureBackupSchedule(
	ctx context.Context,
	service *workflow.Context,
	backupScheduleRef types.NamespacedName,
	clusterName string,
	resourcesToWatch *[]watch.WatchedObject,
) (*mdbv1.AtlasBackupSchedule, error) {
	bSchedule := &mdbv1.AtlasBackupSchedule{}
	err := r.Client.Get(ctx, backupScheduleRef, bSchedule)
	if err != nil {
		return nil, fmt.Errorf("%v backupschedule resource is not found. e: %w", backupScheduleRef, err)
	}

	resourceVersionIsValid := customresource.ValidateResourceVersion(service, bSchedule, r.Log)
	if !resourceVersionIsValid.IsOk() {
		errText := fmt.Sprintf("backup schedule validation result: %v", resourceVersionIsValid)
		r.Log.Debug(errText)
		return nil, errors.New(errText)
	}

	bSchedule.UpdateStatus([]status.Condition{}, status.AtlasBackupScheduleSetDeploymentID(clusterName))

	if err = r.Client.Status().Update(ctx, bSchedule); err != nil {
		r.Log.Errorw("failed to update BackupSchedule status", "error", err)
		return nil, err
	}

	if bSchedule.GetDeletionTimestamp().IsZero() {
		if len(bSchedule.Status.DeploymentIDs) > 0 {
			r.Log.Debugw("adding deletion finalizer", "name", customresource.FinalizerLabel)
			customresource.SetFinalizer(bSchedule, customresource.FinalizerLabel)
		} else {
			r.Log.Debugw("removing deletion finalizer", "name", customresource.FinalizerLabel)
			customresource.UnsetFinalizer(bSchedule, customresource.FinalizerLabel)
		}
	}

	if !bSchedule.GetDeletionTimestamp().IsZero() && customresource.HaveFinalizer(bSchedule, customresource.FinalizerLabel) {
		r.Log.Warnf("backupSchedule %s is assigned to at least one deployment. Remove it from all deployment before delete", bSchedule.Name)
	}

	if err = r.Client.Update(ctx, bSchedule); err != nil {
		r.Log.Errorw("failed to update BackupSchedule object", "error", err)
		return nil, err
	}

	*resourcesToWatch = append(*resourcesToWatch, watch.WatchedObject{ResourceKind: bSchedule.Kind, Resource: backupScheduleRef})

	return bSchedule, nil
}

func (r *AtlasDeploymentReconciler) ensureBackupPolicy(
	ctx context.Context,
	service *workflow.Context,
	bPolicyRef types.NamespacedName,
	resourcesToWatch *[]watch.WatchedObject,
) (*mdbv1.AtlasBackupPolicy, error) {
	bPolicy := &mdbv1.AtlasBackupPolicy{}
	err := r.Client.Get(ctx, bPolicyRef, bPolicy)
	if err != nil {
		return nil, fmt.Errorf("unable to get backuppolicy resource %s. e: %w", bPolicyRef.String(), err)
	}

	resourceVersionIsValid := customresource.ValidateResourceVersion(service, bPolicy, r.Log)
	if !resourceVersionIsValid.IsOk() {
		errText := fmt.Sprintf("backup policy validation result: %v", resourceVersionIsValid)
		r.Log.Debug(errText)
		return nil, errors.New(errText)
	}

	bPolicy.UpdateStatus([]status.Condition{}, status.AtlasBackupPolicySetScheduleID(bPolicyRef.String()))

	if err = r.Client.Status().Update(ctx, bPolicy); err != nil {
		r.Log.Errorw("failed to update BackupPolicy status", "error", err)
		return nil, err
	}

	if bPolicy.GetDeletionTimestamp().IsZero() {
		if len(bPolicy.Status.BackupScheduleIDs) > 0 {
			r.Log.Debugw("adding deletion finalizer", "name", customresource.FinalizerLabel)
			customresource.SetFinalizer(bPolicy, customresource.FinalizerLabel)
		} else {
			r.Log.Debugw("removing deletion finalizer", "name", customresource.FinalizerLabel)
			customresource.UnsetFinalizer(bPolicy, customresource.FinalizerLabel)
		}
	}

	if !bPolicy.GetDeletionTimestamp().IsZero() && customresource.HaveFinalizer(bPolicy, customresource.FinalizerLabel) {
		r.Log.Warnf("backupPolicy %s is assigned to at least one BackupSchedule. Remove it from all BackupSchedules before delete", bPolicy.Name)
	}

	if err = r.Client.Update(ctx, bPolicy); err != nil {
		r.Log.Errorw("failed to update BackupPolicy object", "error", err)
		return nil, err
	}

	*resourcesToWatch = append(*resourcesToWatch, watch.WatchedObject{ResourceKind: bPolicy.Kind, Resource: bPolicyRef})

	return bPolicy, nil
}

func (r *AtlasDeploymentReconciler) updateBackupScheduleAndPolicy(
	ctx context.Context,
	service *workflow.Context,
	projectID string,
	clusterName string,
	bSchedule *mdbv1.AtlasBackupSchedule,
	bPolicy *mdbv1.AtlasBackupPolicy,
) error {
	// Create new backup configuration
	r.Log.Debugf("updating backup configuration for the atlas deployment: %v", clusterName)

	apiPolicy := mongodbatlas.Policy{}

	for _, bpItem := range bPolicy.Spec.Items {
		apiPolicy.PolicyItems = append(apiPolicy.PolicyItems, mongodbatlas.PolicyItem{
			FrequencyInterval: bpItem.FrequencyInterval,
			FrequencyType:     strings.ToLower(bpItem.FrequencyType),
			RetentionValue:    bpItem.RetentionValue,
			RetentionUnit:     strings.ToLower(bpItem.RetentionUnit),
		})
	}

	r.Log.Debugf("updating backup configuration for the atlas deployment: %v", clusterName)
	apiScheduleReq := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
		ClusterName:           clusterName,
		ReferenceHourOfDay:    &bSchedule.Spec.ReferenceHourOfDay,
		ReferenceMinuteOfHour: &bSchedule.Spec.ReferenceMinuteOfHour,
		RestoreWindowDays:     &bSchedule.Spec.RestoreWindowDays,
		UpdateSnapshots:       &bSchedule.Spec.UpdateSnapshots,
		Policies:              []mongodbatlas.Policy{apiPolicy},
	}

	currentSchedule, response, err := service.Client.CloudProviderSnapshotBackupPolicies.Delete(ctx, projectID, clusterName)
	if err != nil {
		errMessage := "unable to delete current backup configuration for project"
		r.Log.Debugf("%s: %s:%s, %v", errMessage, projectID, clusterName, err)
		return fmt.Errorf("%s: %s:%s, %w", errMessage, projectID, clusterName, err)
	}

	if currentSchedule == nil && response != nil {
		return fmt.Errorf("can't delete сurrent backup configuration. response status: %s", response.Status)
	}

	r.Log.Debugf("successfully deleted backup configuration. Default schedule received: %v", currentSchedule)

	apiScheduleReq.ClusterID = currentSchedule.ClusterID
	// There is only one policy, always
	apiScheduleReq.Policies[0].ID = currentSchedule.Policies[0].ID

	r.Log.Debugf("applying backup configuration: %v", *bSchedule)
	if _, _, err := service.Client.CloudProviderSnapshotBackupPolicies.Update(ctx, projectID, clusterName, apiScheduleReq); err != nil {
		return fmt.Errorf("unable to create backupschedule %s. e: %w", client.ObjectKeyFromObject(bSchedule).String(), err)
	}
	r.Log.Infof("successfully updated backup configuration for deployment %v", clusterName)
	return nil
}

func (r *AtlasDeploymentReconciler) garbageCollectBackupResource(ctx context.Context, clusterName string) error {
	schedules := &mdbv1.AtlasBackupScheduleList{}

	err := r.Client.List(ctx, schedules)
	if err != nil {
		return fmt.Errorf("failed to retrieve list of backup schedules: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, bSchedule := range schedules.Items {
		backupSchedule := bSchedule
		g.Go(func() error {
			for _, id := range backupSchedule.Status.DeploymentIDs {
				if id != clusterName {
					continue
				}

				backupSchedule.UpdateStatus([]status.Condition{}, status.AtlasBackupScheduleUnsetDeploymentID(clusterName))

				if len(backupSchedule.Status.DeploymentIDs) == 0 &&
					customresource.HaveFinalizer(&backupSchedule, customresource.FinalizerLabel) {
					customresource.UnsetFinalizer(&backupSchedule, customresource.FinalizerLabel)
				}

				if err = r.Client.Update(ctx, &backupSchedule); err != nil {
					r.Log.Errorw("failed to update BackupSchedule object", "error", err)
					return err
				}

				if backupSchedule.DeletionTimestamp.IsZero() {
					continue
				}

				bPolicy := &mdbv1.AtlasBackupPolicy{}
				bPolicyRef := *backupSchedule.Spec.PolicyRef.GetObject(backupSchedule.Namespace)
				err = r.Client.Get(ctx, bPolicyRef, bPolicy)
				if err != nil {
					return fmt.Errorf("failed to retrieve list of backup schedules: %w", err)
				}

				bPolicy.UpdateStatus([]status.Condition{}, status.AtlasBackupPolicyUnsetScheduleID(bPolicyRef.String()))

				if len(bPolicy.Status.BackupScheduleIDs) == 0 &&
					customresource.HaveFinalizer(bPolicy, customresource.FinalizerLabel) {
					customresource.UnsetFinalizer(bPolicy, customresource.FinalizerLabel)
				}

				if err = r.Client.Update(ctx, bPolicy); err != nil {
					r.Log.Errorw("failed to update BackupPolicy object", "error", err)
					return err
				}
			}

			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}
