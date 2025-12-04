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

package atlasproject

import (
	"errors"
	"fmt"
	"sort"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

func ensureCloudProviderIntegration(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	roleStatuses := project.Status.DeepCopy().CloudProviderIntegrations
	roleSpecs := getCloudProviderIntegrations(project.Spec)

	if len(roleSpecs) == 0 && len(roleStatuses) == 0 {
		workflowCtx.UnsetCondition(api.CloudProviderIntegrationReadyType)
		return workflow.OK()
	}

	allAuthorized, err := syncCloudProviderIntegration(workflowCtx, project.ID(), roleSpecs)
	if err != nil {
		result := workflow.Terminate(workflow.ProjectCloudIntegrationsIsNotReadyInAtlas, err)
		workflowCtx.SetConditionFromResult(api.CloudProviderIntegrationReadyType, result)

		return result
	}

	if !allAuthorized {
		workflowCtx.SetConditionFalse(api.CloudProviderIntegrationReadyType)

		return workflow.InProgress(workflow.ProjectCloudIntegrationsIsNotReadyInAtlas, "not all entries are authorized")
	}

	warnDeprecationMsg := ""
	if len(project.Spec.CloudProviderAccessRoles) > 0 {
		warnDeprecationMsg = "The CloudProviderAccessRole has been deprecated, please move your configuration under CloudProviderIntegration."
	}

	workflowCtx.SetConditionTrueMsg(api.CloudProviderIntegrationReadyType, warnDeprecationMsg)

	return workflow.OK()
}

func syncCloudProviderIntegration(workflowCtx *workflow.Context, projectID string, cpaSpecs []akov2.CloudProviderIntegration) (bool, error) {
	// this endpoint does not offer paginated responses
	atlasCPAs, _, err := workflowCtx.SdkClientSet.SdkClient20250312009.CloudProviderAccessApi.
		ListCloudProviderAccess(workflowCtx.Context, projectID).
		Execute()
	if err != nil {
		return false, fmt.Errorf("unable to fetch cloud provider access from Atlas: %w", err)
	}

	AWSRoles := sortAtlasCPAsByRoleID(atlasCPAs.GetAwsIamRoles())
	cpiStatuses := enrichStatuses(initiateStatuses(cpaSpecs), AWSRoles)
	cpiStatusesToUpdate := make([]status.CloudProviderIntegration, 0, len(cpiStatuses))
	withError := false

	for _, cpiStatus := range cpiStatuses {
		switch cpiStatus.Status {
		case status.CloudProviderIntegrationStatusNew, status.CloudProviderIntegrationStatusFailedToCreate:
			createCloudProviderAccess(workflowCtx, projectID, cpiStatus)
			cpiStatusesToUpdate = append(cpiStatusesToUpdate, *cpiStatus)
		case status.CloudProviderIntegrationStatusCreated, status.CloudProviderIntegrationStatusFailedToAuthorize:
			if cpiStatus.IamAssumedRoleArn != "" {
				authorizeCloudProviderAccess(workflowCtx, projectID, cpiStatus)
			}
			cpiStatusesToUpdate = append(cpiStatusesToUpdate, *cpiStatus)
		case status.CloudProviderIntegrationStatusDeAuthorize, status.CloudProviderIntegrationStatusFailedToDeAuthorize:
			deleteCloudProviderAccess(workflowCtx, projectID, cpiStatus)
		case status.CloudProviderIntegrationStatusAuthorized:
			cpiStatusesToUpdate = append(cpiStatusesToUpdate, *cpiStatus)
		}

		if cpiStatus.ErrorMessage != "" {
			withError = true
		}
	}

	workflowCtx.EnsureStatusOption(status.AtlasProjectCloudIntegrationsOption(cpiStatusesToUpdate))

	if withError {
		return false, errors.New("not all items were synchronized successfully")
	}

	for _, cpiStatus := range cpiStatusesToUpdate {
		if cpiStatus.Status != status.CloudProviderIntegrationStatusAuthorized {
			return false, nil
		}
	}

	return true, nil
}

func initiateStatuses(cpiSpecs []akov2.CloudProviderIntegration) []*status.CloudProviderIntegration {
	cpiStatuses := make([]*status.CloudProviderIntegration, 0, len(cpiSpecs))

	for _, cpiSpec := range cpiSpecs {
		newStatus := status.NewCloudProviderIntegration(cpiSpec.ProviderName, cpiSpec.IamAssumedRoleArn)
		cpiStatuses = append(cpiStatuses, &newStatus)
	}

	return cpiStatuses
}

func enrichStatuses(cpiStatuses []*status.CloudProviderIntegration, atlasCPAs []admin.CloudProviderAccessAWSIAMRole) []*status.CloudProviderIntegration {
	// find configured matches: containing IAM Assumed Role ARN
	for _, cpiStatus := range cpiStatuses {
		for _, atlasCPA := range atlasCPAs {
			cpa := atlasCPA

			if isMatch(cpiStatus, &cpa) {
				copyCloudProviderAccessData(cpiStatus, cpa)

				continue
			}
		}
	}

	// Separate created but not authorized entries: when having empty IAM Assumed Role ARN
	noMatch := make([]admin.CloudProviderAccessAWSIAMRole, 0, len(cpiStatuses))
	for _, atlasCPA := range atlasCPAs {
		cpa := atlasCPA

		if _, ok := cpa.GetIamAssumedRoleArnOk(); !ok {
			noMatch = append(noMatch, cpa)
		}
	}

	// find not configured matches: when having empty IAM Assumed Role ARN
	for _, cpiStatus := range cpiStatuses {
		if cpiStatus.IamAssumedRoleArn != "" && cpiStatus.RoleID != "" {
			continue
		}

		if len(noMatch) == 0 {
			break
		}

		copyCloudProviderAccessData(cpiStatus, noMatch[0])
		noMatch = noMatch[1:]
	}

	cpiKey := "%s.%s"
	cpiStatusesMap := map[string]*status.CloudProviderIntegration{}
	for _, cpiStatus := range cpiStatuses {
		if cpiStatus.IamAssumedRoleArn != "" {
			cpiStatusesMap[fmt.Sprintf(cpiKey, cpiStatus.ProviderName, cpiStatus.IamAssumedRoleArn)] = cpiStatus
		}
	}

	// find removals: configured roles matches that are not on spec
	for _, atlasCPA := range atlasCPAs {
		cpa := atlasCPA

		if _, ok := cpa.GetIamAssumedRoleArnOk(); !ok {
			continue
		}

		if _, ok := cpiStatusesMap[fmt.Sprintf(cpiKey, cpa.ProviderName, cpa.GetIamAssumedRoleArn())]; !ok {
			deleteStatus := status.NewCloudProviderIntegration(cpa.ProviderName, cpa.GetIamAssumedRoleArn())
			copyCloudProviderAccessData(&deleteStatus, cpa)
			deleteStatus.Status = status.CloudProviderIntegrationStatusDeAuthorize
			cpiStatuses = append(cpiStatuses, &deleteStatus)
		}
	}

	for _, cpa := range noMatch {
		deleteStatus := status.NewCloudProviderIntegration(cpa.ProviderName, cpa.GetIamAssumedRoleArn())
		copyCloudProviderAccessData(&deleteStatus, cpa)
		deleteStatus.Status = status.CloudProviderIntegrationStatusDeAuthorize
		cpiStatuses = append(cpiStatuses, &deleteStatus)
	}

	return cpiStatuses
}

func sortAtlasCPAsByRoleID(atlasCPAs []admin.CloudProviderAccessAWSIAMRole) []admin.CloudProviderAccessAWSIAMRole {
	sort.Slice(atlasCPAs, func(i, j int) bool {
		return atlasCPAs[i].GetRoleId() < atlasCPAs[j].GetRoleId()
	})

	return atlasCPAs
}

func isMatch(cpaSpec *status.CloudProviderIntegration, atlasCPA *admin.CloudProviderAccessAWSIAMRole) bool {
	return atlasCPA.GetIamAssumedRoleArn() != "" && cpaSpec.IamAssumedRoleArn != "" &&
		atlasCPA.ProviderName == cpaSpec.ProviderName &&
		atlasCPA.GetIamAssumedRoleArn() == cpaSpec.IamAssumedRoleArn
}

func getCloudProviderIntegrations(projectSpec akov2.AtlasProjectSpec) []akov2.CloudProviderIntegration {
	if len(projectSpec.CloudProviderAccessRoles) > 0 {
		cpis := make([]akov2.CloudProviderIntegration, 0, len(projectSpec.CloudProviderIntegrations))

		for _, cpa := range projectSpec.CloudProviderAccessRoles {
			cpis = append(cpis, akov2.CloudProviderIntegration(cpa))
		}

		return cpis
	}

	return projectSpec.CloudProviderIntegrations
}

func copyCloudProviderAccessData(cpiStatus *status.CloudProviderIntegration, atlasCPA admin.CloudProviderAccessAWSIAMRole) {
	cpiStatus.AtlasAWSAccountArn = atlasCPA.GetAtlasAWSAccountArn()
	cpiStatus.AtlasAssumedRoleExternalID = atlasCPA.GetAtlasAssumedRoleExternalId()
	cpiStatus.RoleID = atlasCPA.GetRoleId()
	cpiStatus.CreatedDate = timeutil.FormatISO8601(atlasCPA.GetCreatedDate())

	if authorizedAt, ok := atlasCPA.GetAuthorizedDateOk(); ok {
		cpiStatus.AuthorizedDate = timeutil.FormatISO8601(*authorizedAt)
	}
	cpiStatus.Status = status.CloudProviderIntegrationStatusCreated

	if _, ok := atlasCPA.GetAuthorizedDateOk(); ok {
		cpiStatus.Status = status.CloudProviderIntegrationStatusAuthorized
	}

	if len(atlasCPA.GetFeatureUsages()) > 0 {
		cpiStatus.FeatureUsages = make([]status.FeatureUsage, 0, len(atlasCPA.GetFeatureUsages()))

		for _, feature := range atlasCPA.GetFeatureUsages() {
			id := ""

			if fID, ok := feature.GetFeatureIdOk(); ok {
				id = fmt.Sprintf("%s.%s", fID.GetGroupId(), fID.GetBucketName())
			}

			cpiStatus.FeatureUsages = append(
				cpiStatus.FeatureUsages,
				status.FeatureUsage{
					FeatureID:   id,
					FeatureType: feature.GetFeatureType(),
				},
			)
		}
	}
}

func createCloudProviderAccess(workflowCtx *workflow.Context, projectID string, cpiStatus *status.CloudProviderIntegration) *status.CloudProviderIntegration {
	cpa, _, err := workflowCtx.SdkClientSet.SdkClient20250312009.CloudProviderAccessApi.CreateCloudProviderAccess(
		workflowCtx.Context,
		projectID,
		&admin.CloudProviderAccessRoleRequest{
			ProviderName: cpiStatus.ProviderName,
		},
	).Execute()
	if err != nil {
		workflowCtx.Log.Errorf("failed to start new cloud provider access: %s", err)
		cpiStatus.Status = status.CloudProviderIntegrationStatusFailedToCreate
		cpiStatus.ErrorMessage = err.Error()

		return cpiStatus
	}

	copyCloudProviderAccessData(cpiStatus, convertCloudProviderAccessRole(*cpa))

	return cpiStatus
}

func authorizeCloudProviderAccess(workflowCtx *workflow.Context, projectID string, cpiStatus *status.CloudProviderIntegration) *status.CloudProviderIntegration {
	cpa, _, err := workflowCtx.SdkClientSet.SdkClient20250312009.CloudProviderAccessApi.AuthorizeProviderAccessRole(
		workflowCtx.Context,
		projectID,
		cpiStatus.RoleID,
		&admin.CloudProviderAccessRoleRequestUpdate{
			ProviderName:      cpiStatus.ProviderName,
			IamAssumedRoleArn: &cpiStatus.IamAssumedRoleArn,
		},
	).Execute()
	if err != nil {
		workflowCtx.Log.Errorf(fmt.Sprintf("failed to authorize cloud provider access: %s", err))
		cpiStatus.Status = status.CloudProviderIntegrationStatusFailedToAuthorize
		cpiStatus.ErrorMessage = err.Error()

		return cpiStatus
	}

	copyCloudProviderAccessData(cpiStatus, convertCloudProviderAccessRole(*cpa))

	return cpiStatus
}

func deleteCloudProviderAccess(workflowCtx *workflow.Context, projectID string, cpiStatus *status.CloudProviderIntegration) {
	_, err := workflowCtx.SdkClientSet.SdkClient20250312009.CloudProviderAccessApi.DeauthorizeProviderAccessRole(
		workflowCtx.Context,
		projectID,
		cpiStatus.ProviderName,
		cpiStatus.RoleID,
	).Execute()
	if err != nil {
		workflowCtx.Log.Errorf(fmt.Sprintf("failed to delete cloud provider access: %s", err))
		cpiStatus.Status = status.CloudProviderIntegrationStatusFailedToDeAuthorize
		cpiStatus.ErrorMessage = err.Error()
	}
}

type CloudProviderIntegrationIdentifiable akov2.CloudProviderIntegration

func convertCloudProviderAccessRole(cpa admin.CloudProviderAccessRole) admin.CloudProviderAccessAWSIAMRole {
	return admin.CloudProviderAccessAWSIAMRole{
		Id:            cpa.Id,
		ProviderName:  cpa.ProviderName,
		FeatureUsages: cpa.FeatureUsages,
		CreatedDate:   cpa.CreatedDate,

		AtlasAWSAccountArn:         cpa.AtlasAWSAccountArn,
		AtlasAssumedRoleExternalId: cpa.AtlasAssumedRoleExternalId,
		AuthorizedDate:             cpa.AuthorizedDate,
		IamAssumedRoleArn:          cpa.IamAssumedRoleArn,
		RoleId:                     cpa.RoleId,

		AtlasAzureAppId:    cpa.AtlasAzureAppId,
		ServicePrincipalId: cpa.ServicePrincipalId,
		TenantId:           cpa.TenantId,
		LastUpdatedDate:    cpa.LastUpdatedDate,
	}
}
