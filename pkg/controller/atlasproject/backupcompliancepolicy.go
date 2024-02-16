/*
Copyright 2023 MongoDB.

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
	"slices"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const ProjectAnnotation = "mongodbatlas/project"

func (r *AtlasProjectReconciler) ensureBackupCompliance(ctx *workflow.Context, project *mdbv1.AtlasProject) workflow.Result {
	defer func() {
		// TODO set condition from error here
		r.garbageCollectBackupResource(ctx.Context, project)
	}()

	if IsBackupComplianceEmpty(project.Spec.BackupCompliancePolicyRef) {
		// check if it is actually enabled in Atlas
		atlasCompliancePolicy, _, err := ctx.SdkClient.CloudBackupsApi.GetDataProtectionSettings(ctx.Context, project.ID()).Execute()
		if err != nil {
			ctx.Log.Errorf("failed to get backup compliance policy from atlas: %v", err)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
		}
		// if not in atlas, we can return OK
		// atlas returns an empty object rather than an error
		if (*atlasCompliancePolicy == admin.DataProtectionSettings20231001{}) {
			return workflow.OK()
		}
		// if it is enabled in Atlas, we still have to signal here via the status condition
		// that there is an not-deleted-yet backup compliance policy in Atlas
		ctx.SetConditionFalseMsg(status.BackupComplianceReadyType, "Backup Compliance Policy must be deleted via Support")
		return workflow.OK()
	}

	compliancePolicy := &mdbv1.AtlasBackupCompliancePolicy{}
	defer func() {
		ctx.AddResourcesToWatch(watch.WatchedObject{ResourceKind: compliancePolicy.Kind, Resource: *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace)})
		r.Log.Debugf("watched backup compliance policy resource: %v\r\n", project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace))
	}()

	// reference set
	// TODO start watching backup compliance CR
	// check reference points to existing compliance policy CR
	err := r.Client.Get(ctx.Context, *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace), compliancePolicy)
	if err != nil {
		ctx.Log.Errorf("failed to get backup compliance policy: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}

	if compliancePolicy.Annotations == nil {
		compliancePolicy.Annotations = map[string]string{}
	}
	if projectIds, ok := compliancePolicy.Annotations[ProjectAnnotation]; !ok {
		compliancePolicy.Annotations[ProjectAnnotation] = project.ID()
	} else {
		if !slices.Contains(strings.Split(projectIds, ","), project.ID()) {
			compliancePolicy.Annotations[ProjectAnnotation] = projectIds + "," + project.ID()
		}
	}

	err = r.Client.Update(ctx.Context, compliancePolicy)
	if err != nil {
		ctx.Log.Errorf("failed to update backup compliance policy: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}

	// check if compliance policy exists in atlas (and matches)
	// if match, return workflow.OK()
	atlasCompliancePolicy, _, err := ctx.SdkClient.CloudBackupsApi.GetDataProtectionSettings(ctx.Context, project.ID()).Execute()

	if err != nil {
		ctx.Log.Errorf("failed to get backup compliance policy from atlas: %v", err)
		return workflow.Terminate(workflow.ProjectBackupCompliancePolicyUnavailable, err.Error())
	}

	if cmp.Equal(atlasCompliancePolicy, compliancePolicy.ToAtlas(project.ID())) {
		return workflow.OK()
	}

	// otherwise, create/update compliance policy...

	// create compliance policy in atlas
	result := syncBackupCompliancePolicy(ctx, project.ID(), *compliancePolicy, *atlasCompliancePolicy)
	ctx.SetConditionFromResult(status.BackupComplianceReadyType, result)
	return result
}

func IsBackupComplianceEmpty(backupCompliancePolicyRef *common.ResourceRefNamespaced) bool {
	return (backupCompliancePolicyRef == nil) || (backupCompliancePolicyRef.Name == "")
}

// syncBackupCompliancePolicy compares the compliance policy specified in Kubernetes to the one currently present in Atlas, updating Atlas should they differ.
func syncBackupCompliancePolicy(ctx *workflow.Context, groupID string, kubeCompliancePolicy mdbv1.AtlasBackupCompliancePolicy, atlasCompliancePolicy admin.DataProtectionSettings20231001) workflow.Result {
	localCompliancePolicy := kubeCompliancePolicy.ToAtlas(groupID)
	// TODO diff will not work here - we need to implement a comparison ourselves
	// also note that the API returns "unknown" for the first & last names, making it hard to compare
	if cmp.Diff(localCompliancePolicy, atlasCompliancePolicy, cmpopts.EquateEmpty()) != "" {
		_, _, err := ctx.SdkClient.CloudBackupsApi.UpdateDataProtectionSettings(ctx.Context, groupID, localCompliancePolicy).Execute()
		if err != nil {
			ctx.Log.Errorf("failed to update backup compliance policy in atlas: %v", err)
			return workflow.Terminate(workflow.ProjectBackupCompliancePolicyNotCreatedInAtlas, err.Error())
		}
	}

	return workflow.OK()
}
