// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
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
	// Reference to the AtlasTeam custom resource which will be assigned to the project.
	TeamRef common.ResourceRefNamespaced `json:"teamRef"`
	// +kubebuilder:validation:MinItems=1
	// Roles the users in the team has within the project.
	Roles []TeamRole `json:"roles"`
}

func (in *Team) ToAtlas(teamID string) admin.TeamRole {
	roleNames := make([]string, 0, len(in.Roles))
	result := admin.TeamRole{
		TeamId:    teamID,
		RoleNames: roleNames,
	}

	for _, role := range in.Roles {
		roleNames = append(roleNames, string(role))
	}
	result.SetRoleNames(roleNames)

	return result
}
