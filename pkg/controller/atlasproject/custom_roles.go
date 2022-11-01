package atlasproject

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureCustomRoles(ctx *workflow.Context, projectID string, project *v1.AtlasProject) (result workflow.Result) {
	isFailure := func(err error) bool {
		if err != nil {
			result = workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
			ctx.SetConditionFromResult(status.ProjectCustomRolesReadyType, result)

			return true
		}

		return false
	}

	currentCustomRoles, err := fetchCustomRoles(ctx, projectID)
	if isFailure(err) {
		return
	}

	err = deleteCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if isFailure(err) {
		return
	}

	err = updateCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if isFailure(err) {
		return
	}

	err = createCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if isFailure(err) {
		return
	}

	ctx.SetConditionTrue(status.ProjectCustomRolesReadyType)

	if len(project.Spec.CustomRoles) == 0 {
		ctx.UnsetCondition(status.ProjectCustomRolesReadyType)
	}

	return workflow.OK()
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

func deleteCustomRoles(ctx *workflow.Context, projectID string, currentCustomRoles []v1.CustomRole, desiredCustomRoles []v1.CustomRole) error {
	toDelete := make([]v1.CustomRole, 0, len(currentCustomRoles))
	desiredCustomRolesByName := mapCustomRolesByName(desiredCustomRoles)

	for _, currentCustomRole := range currentCustomRoles {
		if _, ok := desiredCustomRolesByName[currentCustomRole.Name]; !ok {
			toDelete = append(toDelete, currentCustomRole)
		}
	}

	ctx.Log.Debugw("Custom Roles to be deleted", "NumItems", len(toDelete))

	var err error
	for _, customRole := range toDelete {
		_, delErr := ctx.Client.CustomDBRoles.Delete(context.Background(), projectID, customRole.Name)

		if delErr != nil {
			err = multierror.Append(err, fmt.Errorf("failed to delete custom role \"%s\": %w", customRole.Name, delErr))
		}

		ctx.Log.Debugw("Removed Custom Role from current AtlasProject", "custom role", customRole.Name)
	}

	return err
}

func updateCustomRoles(ctx *workflow.Context, projectID string, currentCustomRoles []v1.CustomRole, desiredCustomRoles []v1.CustomRole) error {
	toUpdate := make([]v1.CustomRole, 0, len(desiredCustomRoles))
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)

	for _, desiredCustomRole := range desiredCustomRoles {
		customRole, ok := currentCustomRolesByName[desiredCustomRole.Name]

		if !ok {
			continue
		}

		if d := cmp.Diff(desiredCustomRole, customRole, cmpopts.EquateEmpty()); d != "" {
			toUpdate = append(toUpdate, desiredCustomRole)
		}
	}

	ctx.Log.Debugw("Custom Roles to be updated", "NumItems", len(toUpdate))

	var err error
	for _, customRole := range toUpdate {
		data := customRole.ToAtlas()
		// Patch fails when sending the role name in the body, needs clarification with cloud team
		data.RoleName = ""
		_, _, upErr := ctx.Client.CustomDBRoles.Update(context.Background(), projectID, customRole.Name, data)

		if upErr != nil {
			err = multierror.Append(err, fmt.Errorf("failed to update custom role \"%s\": %w", customRole.Name, upErr))
		}

		ctx.Log.Debugw("Updated Custom Role in current AtlasProject", "custom role", customRole.Name)
	}

	return err
}

func createCustomRoles(ctx *workflow.Context, projectID string, currentCustomRoles []v1.CustomRole, desiredCustomRoles []v1.CustomRole) error {
	toCreate := make([]v1.CustomRole, 0, len(desiredCustomRoles))
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)

	for _, desiredCustomRole := range desiredCustomRoles {
		if _, ok := currentCustomRolesByName[desiredCustomRole.Name]; !ok {
			toCreate = append(toCreate, desiredCustomRole)
		}
	}

	ctx.Log.Debugw("Custom Roles to be added", "NumItems", len(toCreate))

	var err error
	for _, customRole := range toCreate {
		_, _, creErr := ctx.Client.CustomDBRoles.Create(context.Background(), projectID, customRole.ToAtlas())

		if creErr != nil {
			err = multierror.Append(err, fmt.Errorf("failed to create custom role \"%s\": %w", customRole.Name, creErr))
		}

		ctx.Log.Debugw("Created Custom Role in current AtlasProject", "custom role", customRole.Name)
	}

	return err
}

func mapCustomRolesByName(customRoles []v1.CustomRole) map[string]v1.CustomRole {
	customRolesByName := map[string]v1.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}
