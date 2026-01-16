// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasdeployment

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
)

func (r *AtlasDeploymentReconciler) ensureBackupScheduleAndPolicy(service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, projectID string, deployment *akov2.AtlasDeployment, zoneID string) transitionFn {
	if deployment.Spec.BackupScheduleRef.Name == "" {
		r.Log.Debug("no backup schedule configured for the deployment")

		err := r.garbageCollectBackupResource(service.Context, deployment.GetDeploymentName())
		if err != nil {
			return r.transitionFromLegacy(service, deploymentService, projectID, deployment, err)
		}
		return nil
	}

	if deployment.Spec.DeploymentSpec.BackupEnabled == nil || !*deployment.Spec.DeploymentSpec.BackupEnabled {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, fmt.Errorf("can not proceed with backup configuration. Backups are not enabled for cluster %s", deployment.GetDeploymentName()))
	}

	bSchedule, err := r.ensureBackupSchedule(service, deployment)
	if err != nil {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, err)
	}

	bPolicy, err := r.ensureBackupPolicy(service, bSchedule)
	if err != nil {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, err)
	}

	return r.updateBackupScheduleAndPolicy(service.Context, service, deploymentService, projectID, deployment, bSchedule, bPolicy, zoneID)
}

func (r *AtlasDeploymentReconciler) ensureBackupSchedule(
	service *workflow.Context,
	deployment *akov2.AtlasDeployment,
) (*akov2.AtlasBackupSchedule, error) {
	backupScheduleRef := deployment.Spec.BackupScheduleRef.GetObject(deployment.Namespace)
	bSchedule := &akov2.AtlasBackupSchedule{}
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

	bSchedule.UpdateStatus([]api.Condition{}, status.AtlasBackupScheduleSetDeploymentID(deployment.GetDeploymentName()))

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

	return bSchedule, nil
}

func (r *AtlasDeploymentReconciler) ensureBackupPolicy(service *workflow.Context, bSchedule *akov2.AtlasBackupSchedule) (*akov2.AtlasBackupPolicy, error) {
	bPolicyRef := *bSchedule.Spec.PolicyRef.GetObject(bSchedule.Namespace)
	bPolicy := &akov2.AtlasBackupPolicy{}
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
	bPolicy.UpdateStatus([]api.Condition{}, status.AtlasBackupPolicySetScheduleID(scheduleRef))

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

	return bPolicy, nil
}

func (r *AtlasDeploymentReconciler) updateBackupScheduleAndPolicy(ctx context.Context, service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, projectID string, deployment *akov2.AtlasDeployment, bSchedule *akov2.AtlasBackupSchedule, bPolicy *akov2.AtlasBackupPolicy, zoneID string) transitionFn {
	clusterName := deployment.GetDeploymentName()
	currentSchedule, response, err := service.SdkClientSet.SdkClient20250312012.CloudBackupsApi.GetBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		errMessage := "unable to get current backup configuration for project"
		r.Log.Debugf("%s: %s:%s, %v", errMessage, projectID, clusterName, err)
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, fmt.Errorf("%s: %s:%s, %w", errMessage, projectID, clusterName, err))
	}

	if currentSchedule == nil && response != nil {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, fmt.Errorf("can not get сurrent backup configuration. response status: %s", response.Status))
	}

	r.Log.Debugf("successfully received backup configuration: %v", currentSchedule)

	r.Log.Debugf("updating backup configuration for the atlas deployment: %v", clusterName)

	apiScheduleReq := bSchedule.ToAtlas(currentSchedule.GetClusterId(), clusterName, zoneID, bPolicy)

	// There is only one policy, always
	apiScheduleReq.GetPolicies()[0].SetId(currentSchedule.GetPolicies()[0].GetId())

	equal, err := backupSchedulesAreEqual(currentSchedule, apiScheduleReq)
	if err != nil {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, fmt.Errorf("can not compare BackupSchedule resources: %w", err))
	}

	if equal {
		r.Log.Debug("backup schedules are equal, nothing to change")
		return nil
	}

	r.Log.Debugf("applying backup configuration: %v", *bSchedule)
	if _, _, err := service.SdkClientSet.SdkClient20250312012.CloudBackupsApi.UpdateBackupSchedule(ctx, projectID, clusterName, apiScheduleReq).Execute(); err != nil {
		return r.transitionFromLegacy(service, deploymentService, projectID, deployment, fmt.Errorf("unable to create backup schedule %s. e: %w", client.ObjectKeyFromObject(bSchedule).String(), err))
	}
	r.Log.Infof("successfully updated backup configuration for deployment %v", clusterName)
	return r.transitionFromLegacy(service, deploymentService, projectID, deployment, nil)
}

func backupSchedulesAreEqual(currentSchedule, newSchedule *admin.DiskBackupSnapshotSchedule20240805) (bool, error) {
	currentCopy, err := deepCopy(currentSchedule)
	if err != nil {
		return false, err
	}

	newCopy, err := deepCopy(newSchedule)
	if err != nil {
		return false, err
	}

	normalizeBackupSchedule(currentCopy)
	normalizeBackupSchedule(newCopy)

	return reflect.DeepEqual(currentCopy, newCopy), nil
}

func deepCopy[T any](src *T) (*T, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}

	dst := new(T)
	err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func normalizeBackupSchedule(s *admin.DiskBackupSnapshotSchedule20240805) {
	policies := s.GetPolicies()
	if len(policies) > 0 && len(policies[0].GetPolicyItems()) > 0 {
		items := policies[0].GetPolicyItems()
		for i := range items {
			items[i].Id = nil
		}

		policies[0].SetPolicyItems(items)
	}

	copySettings := s.GetCopySettings()
	if len(copySettings) > 0 {
		for i := range copySettings {
			copySettings[i].ZoneId = ""
		}
	}

	s.Links = nil
	s.NextSnapshot = nil
	s.UpdateSnapshots = nil
	s.SetPolicies(policies)
	s.SetCopySettings(copySettings)
}

func (r *AtlasDeploymentReconciler) garbageCollectBackupResource(ctx context.Context, clusterName string) error {
	schedules := &akov2.AtlasBackupScheduleList{}

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

				backupSchedule.UpdateStatus([]api.Condition{}, status.AtlasBackupScheduleUnsetDeploymentID(clusterName))

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

				bPolicy := &akov2.AtlasBackupPolicy{}
				bPolicyRef := *backupSchedule.Spec.PolicyRef.GetObject(backupSchedule.Namespace)
				err = r.Client.Get(ctx, bPolicyRef, bPolicy)
				if err != nil {
					return fmt.Errorf("failed to retrieve list of backup schedules: %w", err)
				}

				scheduleRef := kube.ObjectKeyFromObject(&backupSchedule).String()
				bPolicy.UpdateStatus([]api.Condition{}, status.AtlasBackupPolicyUnsetScheduleID(scheduleRef))

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
