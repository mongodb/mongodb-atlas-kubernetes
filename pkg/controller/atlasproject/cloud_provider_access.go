package atlasproject

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureProviderAccessStatus(context *workflow.Context, project *v1.AtlasProject, groupID string) workflow.Result {
	roleStatuses := project.Status.DeepCopy().CloudProviderAccessRoles
	roleSpecs := project.Spec.DeepCopy().CloudProviderAccessRoles

	result, condition := syncProviderAccessStatus(context, roleSpecs, roleStatuses, groupID)
	if result != workflow.OK() {
		context.SetConditionFromResult(condition, result)
		return result
	}
	context.SetConditionTrue(status.CloudProviderAccessReadyType)
	return result
}

func syncProviderAccessStatus(customContext *workflow.Context, specs []v1.CloudProviderAccessRole, statuses []status.CloudProviderAccessRole, groupID string) (workflow.Result, status.ConditionType) {
	client := customContext.Client
	logger := customContext.Log
	ctx := context.Background() // TODO: eoaueoaueoa
	specToStatusMap := checkStatuses(logger, specs, statuses)
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
			roleStatus.Update(*role)
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
				roleStatus.Update(role)
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
		roleStatus.Update(*role)
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

func checkStatuses(logger *zap.SugaredLogger, specs []v1.CloudProviderAccessRole, statuses []status.CloudProviderAccessRole) map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole {
	result := make(map[v1.CloudProviderAccessRole]status.CloudProviderAccessRole)
	for _, spec := range specs {
		isCreated := false
		for _, status := range statuses {
			if spec.ProviderName == status.ProviderName && spec.IamAssumedRoleArn == status.IamAssumedRoleArn {
				isCreated = true
				result[spec] = status
				break
			}
		}
		if !isCreated {
			newStatus := status.NewCloudProviderAccessRole(spec.ProviderName, spec.IamAssumedRoleArn)
			result[spec] = newStatus
			statuses = append(statuses, newStatus)
		}
	}
	return result
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

	for spec, status := range expectedRoles {
		if status.RoleID == "" {
			diff.toCreate = append(diff.toCreate, spec)
		}
	}
	return diff, nil
}
