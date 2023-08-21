package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"
)

type CustomRole struct {
	// Human-readable label that identifies the role. This name must be unique for this custom role in this project.
	Name string `json:"name"`
	// List of the built-in roles that this custom role inherits.
	// +optional
	InheritedRoles []Role `json:"inheritedRoles,omitempty"`
	// List of the individual privilege actions that the role grants.
	// +optional
	Actions []Action `json:"actions,omitempty"`
}

type Role struct {
	// Human-readable label that identifies the role inherited.
	Name string `json:"name"`
	// Human-readable label that identifies the database on which someone grants the action to one MongoDB user.
	Database string `json:"database"`
}

type Action struct {
	// Human-readable label that identifies the privilege action.
	Name string `json:"name"`
	// List of resources on which you grant the action.
	Resources []Resource `json:"resources"`
}

type Resource struct {
	// Flag that indicates whether to grant the action on the cluster resource. If true, MongoDB Cloud ignores Database and Collection parameters.
	Cluster *bool `json:"cluster,omitempty"`
	// Human-readable label that identifies the database on which you grant the action to one MongoDB user.
	Database *string `json:"database,omitempty"`
	// Human-readable label that identifies the collection on which you grant the action to one MongoDB user.
	Collection *string `json:"collection,omitempty"`
}

func (in *CustomRole) ToAtlas() *mongodbatlas.CustomDBRole {
	actions := make([]mongodbatlas.Action, 0, len(in.Actions))

	for _, action := range in.Actions {
		resources := make([]mongodbatlas.Resource, 0, len(action.Resources))

		for _, resource := range action.Resources {
			if resource.Cluster != nil {
				if !*resource.Cluster {
					resource.Cluster = nil
				}
			}
			resources = append(resources, mongodbatlas.Resource{
				Collection: resource.Collection,
				DB:         resource.Database,
				Cluster:    resource.Cluster,
			})
		}

		actions = append(actions, mongodbatlas.Action{
			Action:    action.Name,
			Resources: resources,
		})
	}

	inheritedRoles := make([]mongodbatlas.InheritedRole, 0, len(in.InheritedRoles))

	for _, inheritedRole := range in.InheritedRoles {
		inheritedRoles = append(inheritedRoles, mongodbatlas.InheritedRole{
			Db:   inheritedRole.Database,
			Role: inheritedRole.Name,
		})
	}

	return &mongodbatlas.CustomDBRole{
		Actions:        actions,
		InheritedRoles: inheritedRoles,
		RoleName:       in.Name,
	}
}
