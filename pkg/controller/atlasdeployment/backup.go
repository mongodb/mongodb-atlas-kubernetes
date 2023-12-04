package atlasdeployment

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/kube"

	"go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/sync/errgroup"

	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func (r *AtlasDeploymentReconciler) ensureBackupScheduleAndPolicy(
	service *workflow.Context,
	projectID string,
	deployment *mdbv1.AtlasDeployment,
	isEnabled bool,
) error {
	if deployment.Spec.BackupScheduleRef.Name == "" {
		r.Log.Debug("no backup schedule configured for the deployment")

		err := r.garbageCollectBackupResource(service.Context, deployment.GetDeploymentName())
		if err != nil {
			return err
		}
		return nil
	}

	if !isEnabled {
		return fmt.Errorf("can not proceed with backup configuration. Backups are not enabled for cluster %s", deployment.GetDeploymentName())
	}

	resourcesToWatch := []watch.WatchedObject{}
	defer func() {
		service.AddResourcesToWatch(resourcesToWatch...)
		r.Log.Debugf("watched backup schedule and policy resources: %v\r\n", r.WatchedResources)
	}()

	bSchedule, err := r.ensureBackupSchedule(service, deployment, &resourcesToWatch)
	if err != nil {
		return err
	}

	bPolicy, err := r.ensureBackupPolicy(service, bSchedule, &resourcesToWatch)
	if err != nil {
		return err
	}

	return r.updateBackupScheduleAndPolicy(service.Context, service, projectID, deployment, bSchedule, bPolicy)
}

func (r *AtlasDeploymentReconciler) ensureBackupSchedule(
	service *workflow.Context,
	deployment *mdbv1.AtlasDeployment,
	resourcesToWatch *[]watch.WatchedObject,
) (*mdbv1.AtlasBackupSchedule, error) {
	backupScheduleRef := deployment.Spec.BackupScheduleRef.GetObject(deployment.Namespace)
	bSchedule := &mdbv1.AtlasBackupSchedule{}
	err := r.Client.Get(service.Context, *backupScheduleRef, bSchedule)
	if err != nil {
		return nil, fmt.Errorf("%v AtlasBackupSchedule resource is not found. e: %w", *backupScheduleRef, err)
	}

	resourceVersionIsValid := customresource.ValidateResourceVersion(service, bSchedule, r.Log)
	if !resourceVersionIsValid.IsOk() {
		errText := fmt.Sprintf("backup schedule validation result: %v", resourceVersionIsValid)
		r.Log.Debug(errText)
		return nil, errors.New(errText)
	}

	if err = validate.BackupSchedule(bSchedule, deployment); err != nil {
		return nil, err
	}

	if !r.AtlasProvider.IsResourceSupported(bSchedule) {
		return nil, errors.New("the AtlasBackupSchedule is not supported by Atlas for government")
	}

	bSchedule.UpdateStatus([]status.Condition{}, status.AtlasBackupScheduleSetDeploymentID(deployment.GetDeploymentName()))

	if err = r.Client.Status().Update(service.Context, bSchedule); err != nil {
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

	if err = r.Client.Update(service.Context, bSchedule); err != nil {
		r.Log.Errorw("failed to update BackupSchedule object", "error", err)
		return nil, err
	}

	*resourcesToWatch = append(*resourcesToWatch, watch.WatchedObject{ResourceKind: bSchedule.Kind, Resource: *backupScheduleRef})

	return bSchedule, nil
}

func (r *AtlasDeploymentReconciler) ensureBackupPolicy(
	service *workflow.Context,
	bSchedule *mdbv1.AtlasBackupSchedule,
	resourcesToWatch *[]watch.WatchedObject,
) (*mdbv1.AtlasBackupPolicy, error) {
	bPolicyRef := *bSchedule.Spec.PolicyRef.GetObject(bSchedule.Namespace)
	bPolicy := &mdbv1.AtlasBackupPolicy{}
	err := r.Client.Get(service.Context, bPolicyRef, bPolicy)
	if err != nil {
		return nil, fmt.Errorf("unable to get AtlasBackupPolicy resource %s. e: %w", bPolicyRef.String(), err)
	}

	resourceVersionIsValid := customresource.ValidateResourceVersion(service, bPolicy, r.Log)
	if !resourceVersionIsValid.IsOk() {
		errText := fmt.Sprintf("backup policy validation result: %v", resourceVersionIsValid)
		r.Log.Debug(errText)
		return nil, errors.New(errText)
	}

	if !r.AtlasProvider.IsResourceSupported(bPolicy) {
		return nil, errors.New("the AtlasBackupPolicy is not supported by Atlas for government")
	}

	scheduleRef := kube.ObjectKeyFromObject(bSchedule).String()
	bPolicy.UpdateStatus([]status.Condition{}, status.AtlasBackupPolicySetScheduleID(scheduleRef))

	if err = r.Client.Status().Update(service.Context, bPolicy); err != nil {
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

	if err = r.Client.Update(service.Context, bPolicy); err != nil {
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
	deployment *mdbv1.AtlasDeployment,
	bSchedule *mdbv1.AtlasBackupSchedule,
	bPolicy *mdbv1.AtlasBackupPolicy,
) error {
	clusterName := deployment.GetDeploymentName()
	currentSchedule, response, err := service.Client.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		errMessage := "unable to get current backup configuration for project"
		r.Log.Debugf("%s: %s:%s, %v", errMessage, projectID, clusterName, err)
		return fmt.Errorf("%s: %s:%s, %w", errMessage, projectID, clusterName, err)
	}

	if currentSchedule == nil && response != nil {
		return fmt.Errorf("can not get сurrent backup configuration. response status: %s", response.Status)
	}

	r.Log.Debugf("successfully received backup configuration: %v", currentSchedule)

	r.Log.Debugf("updating backup configuration for the atlas deployment: %v", clusterName)

	apiScheduleReq := bSchedule.ToAtlas(currentSchedule.ClusterID, clusterName, deployment.GetReplicationSetID(), bPolicy)

	// There is only one policy, always
	apiScheduleReq.Policies[0].ID = currentSchedule.Policies[0].ID

	equal, err := backupSchedulesAreEqual(currentSchedule, apiScheduleReq)
	if err != nil {
		return fmt.Errorf("can not compare BackupSchedule resources: %w", err)
	}

	if equal {
		r.Log.Debug("backup schedules are equal, nothing to change")
		return nil
	}

	r.Log.Debugf("applying backup configuration: %v", *bSchedule)
	if _, _, err := service.Client.CloudProviderSnapshotBackupPolicies.Update(ctx, projectID, clusterName, apiScheduleReq); err != nil {
		return fmt.Errorf("unable to create backup schedule %s. e: %w", client.ObjectKeyFromObject(bSchedule).String(), err)
	}
	r.Log.Infof("successfully updated backup configuration for deployment %v", clusterName)
	return nil
}

func backupSchedulesAreEqual(currentSchedule *mongodbatlas.CloudProviderSnapshotBackupPolicy, newSchedule *mongodbatlas.CloudProviderSnapshotBackupPolicy) (bool, error) {
	currentCopy := mongodbatlas.CloudProviderSnapshotBackupPolicy{}
	err := compat.JSONCopy(&currentCopy, currentSchedule)
	if err != nil {
		return false, err
	}

	newCopy := mongodbatlas.CloudProviderSnapshotBackupPolicy{}
	err = compat.JSONCopy(&newCopy, newSchedule)
	if err != nil {
		return false, err
	}

	normalizeBackupSchedule(&currentCopy)
	normalizeBackupSchedule(&newCopy)
	d := cmp.Diff(&currentCopy, &newCopy, cmpopts.EquateEmpty())
	if d != "" {
		return false, nil
	}
	return true, nil
}

func normalizeBackupSchedule(s *mongodbatlas.CloudProviderSnapshotBackupPolicy) {
	s.Links = nil
	s.NextSnapshot = ""
	if len(s.Policies) > 0 && len(s.Policies[0].PolicyItems) > 0 {
		for i := range s.Policies[0].PolicyItems {
			s.Policies[0].PolicyItems[i].ID = ""
		}
	}
	s.UpdateSnapshots = nil

	if len(s.CopySettings) > 0 {
		for i := range s.CopySettings {
			s.CopySettings[i].ReplicationSpecID = nil
		}
	}
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

				if err = r.Client.Status().Update(ctx, &backupSchedule); err != nil {
					r.Log.Errorw("failed to update BackupSchedule status", "error", err)
					return err
				}

				lastScheduleRef := false
				if len(backupSchedule.Status.DeploymentIDs) == 0 &&
					customresource.HaveFinalizer(&backupSchedule, customresource.FinalizerLabel) {
					customresource.UnsetFinalizer(&backupSchedule, customresource.FinalizerLabel)
					lastScheduleRef = true
				}

				if err = r.Client.Update(ctx, &backupSchedule); err != nil {
					r.Log.Errorw("failed to update BackupSchedule object", "error", err)
					return err
				}

				if !lastScheduleRef {
					continue
				}

				bPolicy := &mdbv1.AtlasBackupPolicy{}
				bPolicyRef := *backupSchedule.Spec.PolicyRef.GetObject(backupSchedule.Namespace)
				err = r.Client.Get(ctx, bPolicyRef, bPolicy)
				if err != nil {
					return fmt.Errorf("failed to retrieve list of backup schedules: %w", err)
				}

				scheduleRef := kube.ObjectKeyFromObject(&backupSchedule).String()
				bPolicy.UpdateStatus([]status.Condition{}, status.AtlasBackupPolicyUnsetScheduleID(scheduleRef))

				if err = r.Client.Status().Update(ctx, bPolicy); err != nil {
					r.Log.Errorw("failed to update BackupPolicy status", "error", err)
					return err
				}

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
