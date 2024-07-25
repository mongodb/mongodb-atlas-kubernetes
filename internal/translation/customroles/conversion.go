package customroles

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type CustomRole struct {
	*akov2.CustomRole
}

func NewCustomRole(role *akov2.CustomRole) CustomRole {
	return CustomRole{
		CustomRole: role,
	}
}

func toAtlas(role *CustomRole) *admin.UserCustomDBRole {
	atlas := admin.NewUserCustomDBRoleWithDefaults()
	atlas.SetRoleName(role.Name)
	atlas.SetActions(toAtlasActions(role))
	atlas.SetInheritedRoles(toAtlasInheritedRoles(role))

	return atlas
}

func toAtlasUpdate(role *CustomRole) *admin.UpdateCustomDBRole {
	atlas := admin.NewUpdateCustomDBRoleWithDefaults()
	atlas.SetActions(toAtlasActions(role))
	atlas.SetInheritedRoles(toAtlasInheritedRoles(role))

	return atlas
}

func toAtlasActions(role *CustomRole) []admin.DatabasePrivilegeAction {
	actions := make([]admin.DatabasePrivilegeAction, 0, len(role.Actions))
	for _, action := range role.Actions {
		resources := make([]admin.DatabasePermittedNamespaceResource, 0, len(action.Resources))
		for _, resource := range action.Resources {
			if resource.Cluster != nil && !*resource.Cluster {
				resource.Cluster = nil
			}
			resources = append(resources, admin.DatabasePermittedNamespaceResource{
				Collection: admin.GetOrDefault(resource.Collection, ""),
				Db:         admin.GetOrDefault(resource.Database, ""),
				Cluster:    admin.GetOrDefault(resource.Cluster, false),
			})
		}
		actions = append(actions, admin.DatabasePrivilegeAction{
			Action:    action.Name,
			Resources: &resources,
		})
	}

	return actions
}

func toAtlasInheritedRoles(role *CustomRole) []admin.DatabaseInheritedRole {
	inheritedRoles := make([]admin.DatabaseInheritedRole, 0, len(role.InheritedRoles))
	for _, inheritedRole := range role.InheritedRoles {
		inheritedRoles = append(inheritedRoles, admin.DatabaseInheritedRole{
			Db:   inheritedRole.Database,
			Role: inheritedRole.Name,
		})
	}

	return inheritedRoles
}

func fromAtlas(role *admin.UserCustomDBRole) CustomRole {
	inheritedRoles := make([]akov2.Role, 0, len(*role.InheritedRoles))

	for _, atlasInheritedRole := range *role.InheritedRoles {
		inheritedRoles = append(inheritedRoles, akov2.Role{
			Name:     atlasInheritedRole.Role,
			Database: atlasInheritedRole.Db,
		})
	}

	actions := make([]akov2.Action, 0, len(*role.Actions))

	for _, atlasAction := range *role.Actions {
		resources := make([]akov2.Resource, 0, len(*atlasAction.Resources))

		for _, atlasResource := range *atlasAction.Resources {
			resources = append(resources, akov2.Resource{
				Cluster:    pointer.MakePtr(atlasResource.Cluster),
				Database:   pointer.MakePtr(atlasResource.Db),
				Collection: pointer.MakePtr(atlasResource.Collection),
			})
		}

		actions = append(actions, akov2.Action{
			Name:      atlasAction.Action,
			Resources: resources,
		})
	}

	return CustomRole{
		CustomRole: &akov2.CustomRole{
			Name:           role.RoleName,
			InheritedRoles: inheritedRoles,
			Actions:        actions,
		},
	}
}
