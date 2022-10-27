package atlasproject

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureCustomRoles(ctx *workflow.Context, projectID string, project *v1.AtlasProject) (result workflow.Result) {
	currentCustomRoles, err := fetchCustomRoles(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	err = deleteCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	err = updateCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	err = createCustomRoles(ctx, projectID, currentCustomRoles, project.Spec.CustomRoles)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	return workflow.OK()
}

func fetchCustomRoles(ctx *workflow.Context, projectID string) ([]v1.CustomRole, error) {
	data, _, err := ctx.Client.CustomDBRoles.List(context.Background(), projectID, nil)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugw("Got Project Settings", "data", data)

	if data == nil {
		return []v1.CustomRole{}, nil
	}

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
		if customRole, ok := desiredCustomRolesByName[currentCustomRole.Name]; !ok {
			toDelete = append(toDelete, customRole)
		}
	}

	for _, customRole := range toDelete {
		_, err := ctx.Client.CustomDBRoles.Delete(context.Background(), projectID, customRole.Name)

		if err != nil {
			return err
		}

		ctx.Log.Debugw("Removed custom role from current AtlasProject", "custom role", customRole.Name)
	}

	return nil
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
			toUpdate = append(toUpdate, customRole)
		}
	}

	for _, customRole := range toUpdate {
		_, _, err := ctx.Client.CustomDBRoles.Update(context.Background(), projectID, customRole.Name, customRole.ToAtlas())

		if err != nil {
			return err
		}

		ctx.Log.Debugw("Updated custom role in current AtlasProject", "custom role", customRole.Name)
	}

	return nil
}

func createCustomRoles(ctx *workflow.Context, projectID string, currentCustomRoles []v1.CustomRole, desiredCustomRoles []v1.CustomRole) error {
	toCreate := make([]v1.CustomRole, 0, len(desiredCustomRoles))
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)

	for _, desiredCustomRole := range desiredCustomRoles {
		if customRole, ok := currentCustomRolesByName[desiredCustomRole.Name]; !ok {
			toCreate = append(toCreate, customRole)
		}
	}

	for _, customRole := range toCreate {
		_, _, err := ctx.Client.CustomDBRoles.Create(context.Background(), projectID, customRole.ToAtlas())

		if err != nil {
			return err
		}

		ctx.Log.Debugw("Created custom role in current AtlasProject", "custom role", customRole.Name)
	}

	return nil
}

func mapCustomRolesByName(customRoles []v1.CustomRole) map[string]v1.CustomRole {
	customRolesByName := map[string]v1.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}
