package atlasproject

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func ensureCustomRoles(workflowCtx *workflow.Context, project *akov2.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canCustomRolesReconcile(workflowCtx, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.ProjectCustomRolesReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.ProjectCustomRolesReadyType, result)

		return result
	}

	currentCustomRoles, err := fetchCustomRoles(workflowCtx, project.ID())
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	ops := calculateChanges(currentCustomRoles, project.Spec.CustomRoles)

	deleteStatus := deleteCustomRoles(workflowCtx, project.ID(), ops.Delete)
	updateStatus := updateCustomRoles(workflowCtx, project.ID(), ops.Update)
	createStatus := createCustomRoles(workflowCtx, project.ID(), ops.Create)

	result := syncCustomRolesStatus(workflowCtx, project.Spec.CustomRoles, createStatus, updateStatus, deleteStatus)

	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.ProjectCustomRolesReadyType, result)

		return result
	}

	workflowCtx.SetConditionTrue(status.ProjectCustomRolesReadyType)

	if len(project.Spec.CustomRoles) == 0 {
		workflowCtx.UnsetCondition(status.ProjectCustomRolesReadyType)
	}

	return result
}

func fetchCustomRoles(ctx *workflow.Context, projectID string) ([]akov2.CustomRole, error) {
	data, _, err := ctx.Client.CustomDBRoles.List(ctx.Context, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve custom roles from atlas: %w", err)
	}

	if data == nil {
		return []akov2.CustomRole{}, nil
	}

	ctx.Log.Debugw("Got Custom Roles", "NumItems", len(*data))

	return mapToOperator(data), nil
}

func mapToOperator(atlasCustomRoles *[]mongodbatlas.CustomDBRole) []akov2.CustomRole {
	customRoles := make([]akov2.CustomRole, 0, len(*atlasCustomRoles))

	for _, atlasCustomRole := range *atlasCustomRoles {
		inheritedRoles := make([]akov2.Role, 0, len(atlasCustomRole.InheritedRoles))

		for _, atlasInheritedRole := range atlasCustomRole.InheritedRoles {
			inheritedRoles = append(inheritedRoles, akov2.Role{
				Name:     atlasInheritedRole.Role,
				Database: atlasInheritedRole.Db,
			})
		}

		actions := make([]akov2.Action, 0, len(atlasCustomRole.Actions))

		for _, atlasAction := range atlasCustomRole.Actions {
			resources := make([]akov2.Resource, 0, len(atlasAction.Resources))

			for _, atlasResource := range atlasAction.Resources {
				resources = append(resources, akov2.Resource{
					Cluster:    atlasResource.Cluster,
					Database:   atlasResource.DB,
					Collection: atlasResource.Collection,
				})
			}

			actions = append(actions, akov2.Action{
				Name:      atlasAction.Action,
				Resources: resources,
			})
		}

		customRoles = append(customRoles, akov2.CustomRole{
			Actions:        actions,
			InheritedRoles: inheritedRoles,
			Name:           atlasCustomRole.RoleName,
		})
	}

	return customRoles
}

func deleteCustomRoles(ctx *workflow.Context, projectID string, toDelete map[string]akov2.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be deleted", "NumItems", len(toDelete))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toDelete {
		_, err := ctx.Client.CustomDBRoles.Delete(ctx.Context, projectID, customRole.Name)

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

func updateCustomRoles(ctx *workflow.Context, projectID string, toUpdate map[string]akov2.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be updated", "NumItems", len(toUpdate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toUpdate {
		data := customRole.ToAtlas()
		// Patch fails when sending the role name in the body, needs clarification with cloud team
		data.RoleName = ""
		_, _, err := ctx.Client.CustomDBRoles.Update(ctx.Context, projectID, customRole.Name, data)

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

func createCustomRoles(ctx *workflow.Context, projectID string, toCreate map[string]akov2.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be added", "NumItems", len(toCreate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toCreate {
		_, _, err := ctx.Client.CustomDBRoles.Create(ctx.Context, projectID, customRole.ToAtlas())

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

func mapCustomRolesByName(customRoles []akov2.CustomRole) map[string]akov2.CustomRole {
	customRolesByName := map[string]akov2.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}

type CustomRolesOperations struct {
	Create map[string]akov2.CustomRole
	Update map[string]akov2.CustomRole
	Delete map[string]akov2.CustomRole
}

func calculateChanges(currentCustomRoles []akov2.CustomRole, desiredCustomRoles []akov2.CustomRole) CustomRolesOperations {
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)
	desiredCustomRolesByName := mapCustomRolesByName(desiredCustomRoles)
	ops := CustomRolesOperations{
		Create: map[string]akov2.CustomRole{},
		Update: map[string]akov2.CustomRole{},
		Delete: map[string]akov2.CustomRole{},
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

func syncCustomRolesStatus(ctx *workflow.Context, desiredCustomRoles []akov2.CustomRole, created, updated, deleted map[string]status.CustomRole) workflow.Result {
	statuses := make([]status.CustomRole, 0, len(desiredCustomRoles))
	var err error

	for _, customRole := range desiredCustomRoles {
		stat, ok := deleted[customRole.Name]

		if ok {
			if stat.Status == status.CustomRoleStatusFailed {
				statuses = append(statuses, stat)
				err = errors.Join(err, fmt.Errorf(stat.Error))
			}

			continue
		}

		if stat, ok = updated[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf(stat.Error))
			}

			continue
		}

		if stat, ok = created[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf(stat.Error))
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

func canCustomRolesReconcile(workflowCtx *workflow.Context, protected bool, akoProject *akov2.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &akov2.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	atlasData, _, err := workflowCtx.Client.CustomDBRoles.List(workflowCtx.Context, akoProject.ID(), nil)
	if err != nil {
		return false, err
	}

	if atlasData == nil || len(*atlasData) == 0 {
		return true, nil
	}

	atlasCustomRoles := mapToOperator(atlasData)

	if cmp.Diff(latestConfig.CustomRoles, atlasCustomRoles, cmpopts.EquateEmpty()) == "" {
		return true, nil
	}

	return cmp.Diff(akoProject.Spec.CustomRoles, atlasCustomRoles, cmpopts.EquateEmpty()) == "", nil
}
