package v1

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
