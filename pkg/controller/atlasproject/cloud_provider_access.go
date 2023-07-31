package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureProviderAccessStatus(ctx context.Context, customContext *workflow.Context, project *v1.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canCloudProviderAccessReconcile(ctx, customContext.Client, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		customContext.SetConditionFromResult(status.CloudProviderAccessReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		customContext.SetConditionFromResult(status.CloudProviderAccessReadyType, result)

		return result
	}

	roleStatuses := project.Status.DeepCopy().CloudProviderAccessRoles
	roleSpecs := project.Spec.DeepCopy().CloudProviderAccessRoles

	if len(roleSpecs) == 0 && len(roleStatuses) == 0 {
		customContext.UnsetCondition(status.CloudProviderAccessReadyType)
		return workflow.OK()
	}

	result, condition := syncProviderAccessStatus(ctx, customContext, roleSpecs, roleStatuses, project.ID())
	if result != workflow.OK() {
		customContext.SetConditionFromResult(condition, result)
		return result
	}
	customContext.SetConditionTrue(status.CloudProviderAccessReadyType)
	return result
}

func syncProviderAccessStatus(ctx context.Context, customContext *workflow.Context, specs []v1.CloudProviderAccessRole, statuses []status.CloudProviderAccessRole, groupID string) (workflow.Result, status.ConditionType) {
	client := customContext.Client
	logger := customContext.Log
	specToStatusMap, haveDuplicate, cantMatch := checkStatuses(specs, statuses)
	if haveDuplicate {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "some roles contains same ARN value"), status.CloudProviderAccessReadyType
	}
	if cantMatch {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "More than one new role"+
			" with ARN may correspond to an existing empty role. Keep only one new role containing data from the status "+
			"field and delete all other roles. You can add them again after authorization is complete."), status.CloudProviderAccessReadyType
	}
	defer func() {
		SetNewStatuses(customContext, specToStatusMap)
	}()

	diff, err := sortAccessRoles(ctx, client.CloudProviderAccess, logger, specToStatusMap, groupID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, fmt.Sprintf("failed to sort access roles: %s", err)),
			status.CloudProviderAccessReadyType
	}
	err = deleteAccessRoles(ctx, client.CloudProviderAccess, logger, diff.toDelete, groupID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, fmt.Sprintf("failed to delete access roles: %s", err)),
			status.CloudProviderAccessReadyType
	}
	err = createAccessRoles(ctx, client.CloudProviderAccess, logger, diff.toCreate, specToStatusMap, groupID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, fmt.Sprintf("failed to create access roles: %s", err)),
			status.CloudProviderAccessReadyType
	}

	tryToAuthorize(ctx, client.CloudProviderAccess, logger, specToStatusMap, groupID)
	updateAccessRoles(diff.toUpdate, specToStatusMap)
	return ensureCloudProviderAccessStatus(specToStatusMap)
}

func tryToAuthorize(ctx context.Context, access mongodbatlas.CloudProviderAccessService, logger *zap.SugaredLogger, statusMap map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole, groupID string) {
	for spec, roleStatus := range statusMap {
		if roleStatus.Status == status.StatusCreated {
			request := mongodbatlas.CloudProviderAuthorizationRequest{
				ProviderName:      spec.ProviderName,
				IAMAssumedRoleARN: spec.IamAssumedRoleArn,
			}
			role, _, err := access.AuthorizeRole(ctx, groupID, roleStatus.RoleID, &request)
			if err != nil {
				roleStatus.FailedToAuthorise(fmt.Sprintf("cant authorize role. %s", err))
				logger.Errorw("cant authorize role", "role", roleStatus.RoleID, "error", err)
				statusMap[spec] = roleStatus
				continue
			}
			roleStatus.Update(*role, roleStatus.IsEmptyARN())
			statusMap[spec] = roleStatus
		}
	}
}

func SetNewStatuses(customContext *workflow.Context, specToStatus map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole) {
	newRoleStatuses := make([]status.CloudProviderAccessRole, 0, len(specToStatus))
	for _, roleStatus := range specToStatus {
		newRoleStatuses = append(newRoleStatuses, roleStatus)
	}
	customContext.EnsureStatusOption(status.AtlasProjectCloudAccessRolesOption(newRoleStatuses))
}

func ensureCloudProviderAccessStatus(statusMap map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole) (workflow.Result, status.ConditionType) {
	ok := true
	for _, roleStatus := range statusMap {
		if roleStatus.Status != status.StatusReady {
			ok = false
		}
	}
	if !ok {
		return workflow.Terminate(workflow.ProjectCloudAccessRolesIsNotReadyInAtlas, "not all roles are ready"),
			status.CloudProviderAccessReadyType
	}
	return workflow.OK(), status.CloudProviderAccessReadyType
}

func updateAccessRoles(toUpdate []mongodbatlas.AWSIAMRole, specToStatus map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole) {
	for _, role := range toUpdate {
		for spec, roleStatus := range specToStatus {
			if role.RoleID == roleStatus.RoleID {
				roleStatus.Update(role, roleStatus.IsEmptyARN())
				specToStatus[spec] = roleStatus
			}
		}
	}
}

func createAccessRoles(ctx context.Context, accessClient mongodbatlas.CloudProviderAccessService, logger *zap.SugaredLogger,
	toCreate []v1.CloudProviderAccessRole, specToStatus map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole, groupID string) error {
	for _, spec := range toCreate {
		role, _, err := accessClient.CreateRole(ctx, groupID, &mongodbatlas.CloudProviderAccessRoleRequest{
			ProviderName: spec.ProviderName,
		})
		if err != nil {
			logger.Error("failed to create access role", zap.Error(err))
			roleStatus, ok := specToStatus[spec]
			if !ok {
				logger.Error("failed to find status for access role")
			}
			roleStatus.Failed(err.Error())
			specToStatus[spec] = roleStatus
			return err
		}
		roleStatus, ok := specToStatus[spec]
		if !ok {
			logger.Error("failed to find status for access role")
			roleStatus.Failed("failed to find status for access role")
			specToStatus[spec] = roleStatus
			continue
		}
		roleStatus.Update(*role, roleStatus.IsEmptyARN())
		specToStatus[spec] = roleStatus
	}
	return nil
}

func deleteAccessRoles(ctx context.Context, accessClient mongodbatlas.CloudProviderAccessService, logger *zap.SugaredLogger, toDelete map[string]string, groupID string) error {
	for roleID, providerName := range toDelete {
		request := mongodbatlas.CloudProviderDeauthorizationRequest{
			ProviderName: providerName,
			GroupID:      groupID,
			RoleID:       roleID,
		}
		_, err := accessClient.DeauthorizeRole(ctx, &request)
		if err != nil {
			logger.Error("failed to deauthorize role", zap.Error(err))
			return err
		}
	}
	return nil
}

func checkStatuses(specs []v1.CloudProviderAccessRole, statuses []status.CloudProviderAccessRole) (map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole, bool, bool) {
	result := make(map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole)
	existStatusWithEmptyARN := false
	emptyRoleIsAssign := false
	var emptyArnRoleStatus status.CloudProviderAccessRole
	for _, spec := range specs {
		isCreated := false
		for _, existedStatus := range statuses {
			if spec.ProviderName == existedStatus.ProviderName && spec.IamAssumedRoleArn == existedStatus.IamAssumedRoleArn {
				isCreated = true
				if _, ok := result[spec]; !ok {
					result[spec] = existedStatus
				} else {
					return nil, true, false
				}
				break
			}
			if existedStatus.IsEmptyARN() {
				existStatusWithEmptyARN = true
				emptyArnRoleStatus = existedStatus
			}
		}
		if !isCreated {
			if emptyRoleIsAssign {
				return nil, false, true
			}
			if existStatusWithEmptyARN {
				emptyRoleIsAssign = true
				if spec.IamAssumedRoleArn != "" {
					emptyArnRoleStatus.Status = status.StatusCreated
					emptyArnRoleStatus.IamAssumedRoleArn = spec.IamAssumedRoleArn
					result[spec] = emptyArnRoleStatus
				} else {
					result[spec] = emptyArnRoleStatus
				}
			} else {
				newStatus := status.NewCloudProviderAccessRole(spec.ProviderName, spec.IamAssumedRoleArn)
				result[spec] = newStatus
				statuses = append(statuses, newStatus)
			}
		}
	}
	return result, false, false
}

type accessRoleDiff struct {
	toCreate []v1.CloudProviderAccessRole
	toUpdate []mongodbatlas.AWSIAMRole
	toDelete map[string]string // roleId -> providerName
}

func sortAccessRoles(ctx context.Context, accessClient mongodbatlas.CloudProviderAccessService, logger *zap.SugaredLogger, expectedRoles map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole, groupID string) (accessRoleDiff, error) {
	roleList, _, err := accessClient.ListRoles(ctx, groupID)
	if err != nil {
		logger.Error("failed to list access roles", zap.Error(err))
		return accessRoleDiff{}, err
	}
	logger.Debugf("found %d access roles", len(roleList.AWSIAMRoles))
	existedRoles := roleList.AWSIAMRoles
	diff := accessRoleDiff{}
	diff.toDelete = make(map[string]string)
	for _, existedRole := range existedRoles {
		toDelete := true
		for _, status := range expectedRoles {
			if status.RoleID == existedRole.RoleID {
				toDelete = false
				diff.toUpdate = append(diff.toUpdate, existedRole)
				break
			}
		}
		if toDelete {
			diff.toDelete[existedRole.RoleID] = existedRole.ProviderName
		}
	}

	for spec, existedStatus := range expectedRoles {
		if existedStatus.RoleID == "" {
			diff.toCreate = append(diff.toCreate, spec)
		}
	}
	return diff, nil
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
