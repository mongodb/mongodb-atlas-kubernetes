package atlasproject

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureCustomRoles(ctx *workflow.Context, projectID string, project *v1.AtlasProject) workflow.Result {
	currentCustomRoles, err := fetchCustomRoles(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	ops := calculateChanges(currentCustomRoles, project.Spec.CustomRoles)

	deleteStatus := deleteCustomRoles(ctx, projectID, ops.Delete)
	updateStatus := updateCustomRoles(ctx, projectID, ops.Update)
	createStatus := createCustomRoles(ctx, projectID, ops.Create)

	result := syncCustomRolesStatus(ctx, project.Spec.CustomRoles, createStatus, updateStatus, deleteStatus)

	if !result.IsOk() {
		ctx.SetConditionFromResult(status.ProjectCustomRolesReadyType, result)

		return result
	}

	ctx.SetConditionTrue(status.ProjectCustomRolesReadyType)

	if len(project.Spec.CustomRoles) == 0 {
		ctx.UnsetCondition(status.ProjectCustomRolesReadyType)
	}

	return result
}

func fetchCustomRoles(ctx *workflow.Context, projectID string) ([]v1.CustomRole, error) {
	data, _, err := ctx.Client.CustomDBRoles.List(context.Background(), projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve custom roles from atlas: %w", err)
	}

	if data == nil {
		return []v1.CustomRole{}, nil
	}

	ctx.Log.Debugw("Got Custom Roles", "NumItems", len(*data))

	customRoles := make([]v1.CustomRole, 0, len(*data))

	for _, atlasCustomRole := range *data {
		inheritedRoles := make([]v1.Role, 0, len(atlasCustomRole.InheritedRoles))

		for _, atlasInheritedRole := range atlasCustomRole.InheritedRoles {
			inheritedRoles = append(inheritedRoles, v1.Role{
				Name:     atlasInheritedRole.Role,
				Database: atlasInheritedRole.Db,
			})
		}

		actions := make([]v1.Action, 0, len(atlasCustomRole.Actions))

		for _, atlasAction := range atlasCustomRole.Actions {
			resources := make([]v1.Resource, 0, len(atlasAction.Resources))

			for _, atlasResource := range atlasAction.Resources {
				if atlasResource.Cluster != nil && !*atlasResource.Cluster {
					atlasResource.Cluster = nil
				}
				resources = append(resources, v1.Resource{
					Cluster:    atlasResource.Cluster,
					Database:   atlasResource.DB,
					Collection: atlasResource.Collection,
				})
			}

			actions = append(actions, v1.Action{
				Name:      atlasAction.Action,
				Resources: resources,
			})
		}

		customRoles = append(customRoles, v1.CustomRole{
			Actions:        actions,
			InheritedRoles: inheritedRoles,
			Name:           atlasCustomRole.RoleName,
		})
	}

	return customRoles, nil
}

func deleteCustomRoles(ctx *workflow.Context, projectID string, toDelete map[string]v1.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be deleted", "NumItems", len(toDelete))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toDelete {
		_, err := ctx.Client.CustomDBRoles.Delete(context.Background(), projectID, customRole.Name)

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

func updateCustomRoles(ctx *workflow.Context, projectID string, toUpdate map[string]v1.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be updated", "NumItems", len(toUpdate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toUpdate {
		data := customRole.ToAtlas()
		// Patch fails when sending the role name in the body, needs clarification with cloud team
		data.RoleName = ""
		_, _, err := ctx.Client.CustomDBRoles.Update(context.Background(), projectID, customRole.Name, data)

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

func createCustomRoles(ctx *workflow.Context, projectID string, toCreate map[string]v1.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be added", "NumItems", len(toCreate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toCreate {
		_, _, err := ctx.Client.CustomDBRoles.Create(context.Background(), projectID, customRole.ToAtlas())

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

func mapCustomRolesByName(customRoles []v1.CustomRole) map[string]v1.CustomRole {
	customRolesByName := map[string]v1.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}

type CustomRolesOperations struct {
	Create map[string]v1.CustomRole
	Update map[string]v1.CustomRole
	Delete map[string]v1.CustomRole
}

func calculateChanges(currentCustomRoles []v1.CustomRole, desiredCustomRoles []v1.CustomRole) CustomRolesOperations {
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)
	desiredCustomRolesByName := mapCustomRolesByName(desiredCustomRoles)
	ops := CustomRolesOperations{
		Create: map[string]v1.CustomRole{},
		Update: map[string]v1.CustomRole{},
		Delete: map[string]v1.CustomRole{},
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

func syncCustomRolesStatus(ctx *workflow.Context, desiredCustomRoles []v1.CustomRole, created, updated, deleted map[string]status.CustomRole) workflow.Result {
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
