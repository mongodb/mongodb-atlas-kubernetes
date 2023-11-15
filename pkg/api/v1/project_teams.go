package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

// +kubebuilder:validation:Enum=GROUP_OWNER;GROUP_CLUSTER_MANAGER;GROUP_DATA_ACCESS_ADMIN;GROUP_DATA_ACCESS_READ_WRITE;GROUP_DATA_ACCESS_READ_ONLY;GROUP_READ_ONLY

type TeamRole string

const (
	TeamRoleOwner               TeamRole = "GROUP_OWNER"
	TeamRoleClusterManager      TeamRole = "GROUP_CLUSTER_MANAGER"
	TeamRoleDataAccessAdmin     TeamRole = "GROUP_DATA_ACCESS_ADMIN"
	TeamRoleDataAccessReadWrite TeamRole = "GROUP_DATA_ACCESS_READ_WRITE"
	TeamRoleDataAccessReadOnly  TeamRole = "GROUP_DATA_ACCESS_READ_ONLY"
	TeamRoleReadOnly            TeamRole = "GROUP_READ_ONLY"
)

type Team struct {
	// Reference to the team which will assigned to the project
	TeamRef common.ResourceRefNamespaced `json:"teamRef"`
	// +kubebuilder:validation:MinItems=1
	// Roles the users of the team has over the project
	Roles []TeamRole `json:"roles"`
}

func (in *Team) ToAtlas(teamID string) *mongodbatlas.ProjectTeam {
	result := &mongodbatlas.ProjectTeam{
		TeamID:    teamID,
		RoleNames: make([]string, 0, len(in.Roles)),
	}

	for _, role := range in.Roles {
		result.RoleNames = append(result.RoleNames, string(role))
	}

	return result
}
