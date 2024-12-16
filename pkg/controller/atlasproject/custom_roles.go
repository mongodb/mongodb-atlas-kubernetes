package atlasproject

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type roleController struct {
	ctx     *workflow.Context
	project *akov2.AtlasProject
	service customroles.CustomRoleService
}

func hasSkippedCustomRoles(atlasProject *akov2.AtlasProject) (bool, error) {
	lastSkippedSpec := akov2.AtlasProjectSpec{}
	lastSkippedSpecString, ok := atlasProject.Annotations[customresource.AnnotationLastSkippedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(lastSkippedSpecString), &lastSkippedSpec); err != nil {
			return false, fmt.Errorf("failed to parse last skipped configuration: %w", err)
		}

		return len(lastSkippedSpec.CustomRoles) != 0, nil
	}

	return false, nil
}

func hasLastAppliedCustomRoles(atlasProject *akov2.AtlasProject) (bool, error) {
	lastAppliedSpec := akov2.AtlasProjectSpec{}
	lastAppliedSpecStr, ok := atlasProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if !ok {
		return false, nil
	}

	if err := json.Unmarshal([]byte(lastAppliedSpecStr), &lastAppliedSpec); err != nil {
		return false, fmt.Errorf("failed to parse last applied configuration: %w", err)
	}

	return len(lastAppliedSpec.CustomRoles) != 0, nil
}

func ensureCustomRoles(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	skipped, err := hasSkippedCustomRoles(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if skipped {
		workflowCtx.UnsetCondition(api.ProjectCustomRolesReadyType)

		return workflow.OK()
	}

	hadPreviousCustomRoles, err := hasLastAppliedCustomRoles(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	r := roleController{
		ctx:     workflowCtx,
		project: project,
		service: customroles.NewCustomRoles(workflowCtx.SdkClient.CustomDatabaseRolesApi),
	}

	currentCustomRoles, err := r.service.List(r.ctx.Context, r.project.ID())
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	akoRoles := make([]customroles.CustomRole, len(project.Spec.CustomRoles))
	for i := range project.Spec.CustomRoles {
		akoRoles[i] = customroles.NewCustomRole(&project.Spec.CustomRoles[i])
	}

	ops := calculateChanges(currentCustomRoles, akoRoles)

	var deleteStatus map[string]status.CustomRole
	if hadPreviousCustomRoles {
		deleteStatus = r.deleteCustomRoles(workflowCtx, project.ID(), ops.Delete)
	}
	updateStatus := r.updateCustomRoles(workflowCtx, project.ID(), ops.Update)
	createStatus := r.createCustomRoles(workflowCtx, project.ID(), ops.Create)

	result := syncCustomRolesStatus(workflowCtx, akoRoles, createStatus, updateStatus, deleteStatus)

	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.ProjectCustomRolesReadyType, result)

		return result
	}

	workflowCtx.SetConditionTrue(api.ProjectCustomRolesReadyType)

	if len(project.Spec.CustomRoles) == 0 {
		workflowCtx.UnsetCondition(api.ProjectCustomRolesReadyType)
	}

	return result
}

func (r *roleController) deleteCustomRoles(ctx *workflow.Context, projectID string, toDelete map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be deleted", "NumItems", len(toDelete))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toDelete {
		err := r.service.Delete(ctx.Context, projectID, customRole.Name)

		opStatus, errorMsg := evaluateOperation(err)
		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to delete custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Removed Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func (r *roleController) updateCustomRoles(ctx *workflow.Context, projectID string, toUpdate map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be updated", "NumItems", len(toUpdate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toUpdate {
		err := r.service.Update(ctx.Context, projectID, customRole.Name, customRole)

		opStatus, errorMsg := evaluateOperation(err)

		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to update custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Updated Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func (r *roleController) createCustomRoles(ctx *workflow.Context, projectID string, toCreate map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be added", "NumItems", len(toCreate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toCreate {
		err := r.service.Create(ctx.Context, projectID, customRole)

		opStatus, errorMsg := evaluateOperation(err)

		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to create custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Created Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func mapCustomRolesByName(customRoles []customroles.CustomRole) map[string]customroles.CustomRole {
	customRolesByName := map[string]customroles.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}

type CustomRolesOperations struct {
	Create map[string]customroles.CustomRole
	Update map[string]customroles.CustomRole
	Delete map[string]customroles.CustomRole
}

func calculateChanges(currentCustomRoles []customroles.CustomRole, desiredCustomRoles []customroles.CustomRole) CustomRolesOperations {
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)
	desiredCustomRolesByName := mapCustomRolesByName(desiredCustomRoles)
	ops := CustomRolesOperations{
		Create: map[string]customroles.CustomRole{},
		Update: map[string]customroles.CustomRole{},
		Delete: map[string]customroles.CustomRole{},
	}

	for _, currentCustomRole := range currentCustomRoles {
		if _, ok := desiredCustomRolesByName[currentCustomRole.Name]; !ok {
			ops.Delete[currentCustomRole.Name] = currentCustomRole
		}
	}

	for _, desiredCustomRole := range desiredCustomRoles {
		customRole, ok := currentCustomRolesByName[desiredCustomRole.Name]

		if !ok {
			ops.Create[desiredCustomRole.Name] = desiredCustomRole

			continue
		}

		if d := cmp.Diff(desiredCustomRole, customRole, cmpopts.EquateEmpty()); d != "" {
			ops.Update[desiredCustomRole.Name] = desiredCustomRole
		}
	}

	return ops
}

func evaluateOperation(err error) (status.CustomRoleStatus, string) {
	if err != nil {
		return status.CustomRoleStatusFailed, err.Error()
	}

	return status.CustomRoleStatusOK, ""
}

func syncCustomRolesStatus(ctx *workflow.Context, desiredCustomRoles []customroles.CustomRole, created, updated, deleted map[string]status.CustomRole) workflow.Result {
	statuses := make([]status.CustomRole, 0, len(desiredCustomRoles))
	var err error

	for _, customRole := range desiredCustomRoles {
		stat, ok := deleted[customRole.Name]

		if ok {
			if stat.Status == status.CustomRoleStatusFailed {
				statuses = append(statuses, stat)
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		if stat, ok = updated[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		if stat, ok = created[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		statuses = append(statuses, status.CustomRole{
			Name:   customRole.Name,
			Status: status.CustomRoleStatusOK,
		})
	}

	ctx.EnsureStatusOption(status.AtlasProjectSetCustomRolesOption(&statuses))

	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, fmt.Sprintf("failed to apply changes to custom roles: %s", err.Error()))
	}

	return workflow.OK()
}
