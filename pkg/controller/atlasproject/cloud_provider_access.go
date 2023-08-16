package atlasproject

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"

	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureProviderAccessStatus(ctx context.Context, workflowCtx *workflow.Context, project *v1.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canCloudProviderAccessReconcile(ctx, workflowCtx.Client, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.CloudProviderAccessReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.CloudProviderAccessReadyType, result)

		return result
	}

	roleStatuses := project.Status.DeepCopy().CloudProviderAccessRoles
	roleSpecs := project.Spec.DeepCopy().CloudProviderAccessRoles

	if len(roleSpecs) == 0 && len(roleStatuses) == 0 {
		workflowCtx.UnsetCondition(status.CloudProviderAccessReadyType)
		return workflow.OK()
	}

	allAuthorized, err := syncCloudProviderAccess(ctx, workflowCtx, project.ID(), project.Spec.CloudProviderAccessRoles)
	if err != nil {
		result := workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, err.Error())
		workflowCtx.SetConditionFromResult(status.CloudProviderAccessReadyType, result)

		return result
	}

	if !allAuthorized {
		workflowCtx.SetConditionFalse(status.CloudProviderAccessReadyType)

		return workflow.InProgress(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "not all entries are authorized")
	}

	workflowCtx.SetConditionTrue(status.CloudProviderAccessReadyType)
	return workflow.OK()
}

func syncCloudProviderAccess(ctx context.Context, workflowCtx *workflow.Context, projectID string, cpaSpecs []v1.CloudProviderAccessRole) (bool, error) {
	atlasCPAs, _, err := workflowCtx.Client.CloudProviderAccess.ListRoles(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("unable to fetch cloud provider access from Atlas: %w", err)
	}

	AWSRoles := sortAtlasCPAsByRoleID(atlasCPAs.AWSIAMRoles)
	cpaStatuses := enrichStatuses(initiateStatuses(cpaSpecs), AWSRoles)
	cpaStatusesToUpdate := make([]status.CloudProviderAccessRole, 0, len(cpaStatuses))
	withError := false

	for _, cpaStatus := range cpaStatuses {
		switch cpaStatus.Status {
		case status.CloudProviderAccessStatusNew, status.CloudProviderAccessStatusFailedToCreate:
			createCloudProviderAccess(ctx, workflowCtx, projectID, cpaStatus)
			cpaStatusesToUpdate = append(cpaStatusesToUpdate, *cpaStatus)
		case status.CloudProviderAccessStatusCreated, status.CloudProviderAccessStatusFailedToAuthorize:
			if cpaStatus.IamAssumedRoleArn != "" {
				authorizeCloudProviderAccess(ctx, workflowCtx, projectID, cpaStatus)
			}
			cpaStatusesToUpdate = append(cpaStatusesToUpdate, *cpaStatus)
		case status.CloudProviderAccessStatusDeAuthorize, status.CloudProviderAccessStatusFailedToDeAuthorize:
			deleteCloudProviderAccess(ctx, workflowCtx, projectID, cpaStatus)
		case status.CloudProviderAccessStatusAuthorized:
			cpaStatusesToUpdate = append(cpaStatusesToUpdate, *cpaStatus)
		}

		if cpaStatus.ErrorMessage != "" {
			withError = true
		}
	}

	workflowCtx.EnsureStatusOption(status.AtlasProjectCloudAccessRolesOption(cpaStatusesToUpdate))

	if withError {
		return false, errors.New("not all items were synchronized successfully")
	}

	for _, capStatus := range cpaStatusesToUpdate {
		if capStatus.Status != status.CloudProviderAccessStatusAuthorized {
			return false, nil
		}
	}

	return true, nil
}

func initiateStatuses(cpaSpecs []v1.CloudProviderAccessRole) []*status.CloudProviderAccessRole {
	cpaStatuses := make([]*status.CloudProviderAccessRole, 0, len(cpaSpecs))

	for _, cpaSpec := range cpaSpecs {
		newStatus := status.NewCloudProviderAccessRole(cpaSpec.ProviderName, cpaSpec.IamAssumedRoleArn)
		cpaStatuses = append(cpaStatuses, &newStatus)
	}

	return cpaStatuses
}

func enrichStatuses(cpaStatuses []*status.CloudProviderAccessRole, atlasCPAs []mongodbatlas.CloudProviderAccessRole) []*status.CloudProviderAccessRole {
	// find configured matches: containing IAM Assumed Role ARN
	for _, cpaStatus := range cpaStatuses {
		for _, atlasCPA := range atlasCPAs {
			cpa := atlasCPA

			if isMatch(cpaStatus, &cpa) {
				copyCloudProviderAccessData(cpaStatus, &cpa)

				continue
			}
		}
	}

	// Separate created but not authorized entries: when having empty IAM Assumed Role ARN
	noMatch := make([]*mongodbatlas.CloudProviderAccessRole, 0, len(cpaStatuses))
	for _, atlasCPA := range atlasCPAs {
		cpa := atlasCPA

		if cpa.IAMAssumedRoleARN == "" {
			noMatch = append(noMatch, &cpa)
		}
	}

	// find not configured matches: when having empty IAM Assumed Role ARN
	for _, cpaStatus := range cpaStatuses {
		if cpaStatus.IamAssumedRoleArn != "" && cpaStatus.RoleID != "" {
			continue
		}

		if len(noMatch) == 0 {
			break
		}

		copyCloudProviderAccessData(cpaStatus, noMatch[0])
		noMatch = noMatch[1:]
	}

	cpaKey := "%s.%s"
	cpaStatusesMap := map[string]*status.CloudProviderAccessRole{}
	for _, cpaStatus := range cpaStatuses {
		if cpaStatus.IamAssumedRoleArn != "" {
			cpaStatusesMap[fmt.Sprintf(cpaKey, cpaStatus.ProviderName, cpaStatus.IamAssumedRoleArn)] = cpaStatus
		}
	}

	// find removals: configured roles matches that are not on spec
	for _, atlasCPA := range atlasCPAs {
		cpa := atlasCPA

		if cpa.IAMAssumedRoleARN == "" {
			continue
		}

		if _, ok := cpaStatusesMap[fmt.Sprintf(cpaKey, cpa.ProviderName, cpa.IAMAssumedRoleARN)]; !ok {
			deleteStatus := status.NewCloudProviderAccessRole(cpa.ProviderName, cpa.IAMAssumedRoleARN)
			copyCloudProviderAccessData(&deleteStatus, &cpa)
			deleteStatus.Status = status.CloudProviderAccessStatusDeAuthorize
			cpaStatuses = append(cpaStatuses, &deleteStatus)
		}
	}

	for _, cpa := range noMatch {
		deleteStatus := status.NewCloudProviderAccessRole(cpa.ProviderName, cpa.IAMAssumedRoleARN)
		copyCloudProviderAccessData(&deleteStatus, cpa)
		deleteStatus.Status = status.CloudProviderAccessStatusDeAuthorize
		cpaStatuses = append(cpaStatuses, &deleteStatus)
	}

	return cpaStatuses
}

func sortAtlasCPAsByRoleID(atlasCPAs []mongodbatlas.CloudProviderAccessRole) []mongodbatlas.CloudProviderAccessRole {
	fmt.Println(atlasCPAs)
	sort.Slice(atlasCPAs, func(i, j int) bool {
		return atlasCPAs[i].RoleID < atlasCPAs[j].RoleID
	})
	fmt.Println(atlasCPAs)
	return atlasCPAs
}

func isMatch(cpaSpec *status.CloudProviderAccessRole, atlasCPA *mongodbatlas.CloudProviderAccessRole) bool {
	return atlasCPA.IAMAssumedRoleARN != "" && cpaSpec.IamAssumedRoleArn != "" &&
		atlasCPA.ProviderName == cpaSpec.ProviderName &&
		atlasCPA.IAMAssumedRoleARN == cpaSpec.IamAssumedRoleArn
}

func copyCloudProviderAccessData(cpaStatus *status.CloudProviderAccessRole, atlasCPA *mongodbatlas.CloudProviderAccessRole) {
	cpaStatus.AtlasAWSAccountArn = atlasCPA.AtlasAWSAccountARN
	cpaStatus.AtlasAssumedRoleExternalID = atlasCPA.AtlasAssumedRoleExternalID
	cpaStatus.RoleID = atlasCPA.RoleID
	cpaStatus.CreatedDate = atlasCPA.CreatedDate
	cpaStatus.AuthorizedDate = atlasCPA.AuthorizedDate
	cpaStatus.Status = status.CloudProviderAccessStatusCreated

	if atlasCPA.AuthorizedDate != "" {
		cpaStatus.Status = status.CloudProviderAccessStatusAuthorized
	}

	if len(atlasCPA.FeatureUsages) > 0 {
		cpaStatus.FeatureUsages = make([]status.FeatureUsage, 0, len(atlasCPA.FeatureUsages))

		for _, feature := range atlasCPA.FeatureUsages {
			if feature == nil {
				continue
			}

			id := ""

			if feature.FeatureID != nil {
				id = feature.FeatureID.(string)
			}

			cpaStatus.FeatureUsages = append(
				cpaStatus.FeatureUsages,
				status.FeatureUsage{
					FeatureID:   id,
					FeatureType: feature.FeatureType,
				},
			)
		}
	}
}

func createCloudProviderAccess(ctx context.Context, workflowCtx *workflow.Context, projectID string, cpaStatus *status.CloudProviderAccessRole) *status.CloudProviderAccessRole {
	cpa, _, err := workflowCtx.Client.CloudProviderAccess.CreateRole(
		ctx,
		projectID,
		&mongodbatlas.CloudProviderAccessRoleRequest{
			ProviderName: cpaStatus.ProviderName,
		},
	)
	if err != nil {
		workflowCtx.Log.Errorf("failed to start new cloud provider access: %s", err)
		cpaStatus.Status = status.CloudProviderAccessStatusFailedToCreate
		cpaStatus.ErrorMessage = err.Error()

		return cpaStatus
	}

	copyCloudProviderAccessData(cpaStatus, cpa)

	return cpaStatus
}

func authorizeCloudProviderAccess(ctx context.Context, workflowCtx *workflow.Context, projectID string, cpaStatus *status.CloudProviderAccessRole) *status.CloudProviderAccessRole {
	cpa, _, err := workflowCtx.Client.CloudProviderAccess.AuthorizeRole(
		ctx,
		projectID,
		cpaStatus.RoleID,
		&mongodbatlas.CloudProviderAccessRoleRequest{
			ProviderName:      cpaStatus.ProviderName,
			IAMAssumedRoleARN: &cpaStatus.IamAssumedRoleArn,
		},
	)
	if err != nil {
		workflowCtx.Log.Errorf(fmt.Sprintf("failed to authorize cloud provider access: %s", err))
		cpaStatus.Status = status.CloudProviderAccessStatusFailedToAuthorize
		cpaStatus.ErrorMessage = err.Error()

		return cpaStatus
	}

	copyCloudProviderAccessData(cpaStatus, cpa)

	return cpaStatus
}

func deleteCloudProviderAccess(ctx context.Context, workflowCtx *workflow.Context, projectID string, cpaStatus *status.CloudProviderAccessRole) {
	_, err := workflowCtx.Client.CloudProviderAccess.DeauthorizeRole(
		ctx,
		&mongodbatlas.CloudProviderDeauthorizationRequest{
			ProviderName: cpaStatus.ProviderName,
			GroupID:      projectID,
			RoleID:       cpaStatus.RoleID,
		},
	)
	if err != nil {
		workflowCtx.Log.Errorf(fmt.Sprintf("failed to delete cloud provider access: %s", err))
		cpaStatus.Status = status.CloudProviderAccessStatusFailedToDeAuthorize
		cpaStatus.ErrorMessage = err.Error()
	}
}

func canCloudProviderAccessReconcile(ctx context.Context, atlasClient mongodbatlas.Client, protected bool, akoProject *v1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &v1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	list, _, err := atlasClient.CloudProviderAccess.ListRoles(ctx, akoProject.ID())
	if err != nil {
		return false, err
	}

	atlasList := make([]CloudProviderAccessIdentifiable, 0, len(list.AWSIAMRoles))
	for _, r := range list.AWSIAMRoles {
		if r.IAMAssumedRoleARN != "" {
			atlasList = append(atlasList,
				CloudProviderAccessIdentifiable{
					ProviderName:      r.ProviderName,
					IamAssumedRoleArn: r.IAMAssumedRoleARN,
				},
			)
		}
	}

	if len(atlasList) == 0 {
		return true, nil
	}

	akoLastList := make([]CloudProviderAccessIdentifiable, len(latestConfig.CloudProviderAccessRoles))
	for i, v := range latestConfig.CloudProviderAccessRoles {
		akoLastList[i] = CloudProviderAccessIdentifiable(v)
	}

	diff := set.Difference(atlasList, akoLastList)

	if len(diff) == 0 {
		return true, nil
	}

	akoCurrentList := make([]CloudProviderAccessIdentifiable, len(akoProject.Spec.CloudProviderAccessRoles))
	for i, v := range akoProject.Spec.CloudProviderAccessRoles {
		akoCurrentList[i] = CloudProviderAccessIdentifiable(v)
	}

	diff = set.Difference(akoCurrentList, atlasList)

	return len(diff) == 0, nil
}

type CloudProviderAccessIdentifiable v1.CloudProviderAccessRole

func (cpa CloudProviderAccessIdentifiable) Identifier() interface{} {
	return fmt.Sprintf("%s.%s", cpa.ProviderName, cpa.IamAssumedRoleArn)
}
