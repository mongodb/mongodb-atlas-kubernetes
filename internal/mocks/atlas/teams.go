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

package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type TeamsClientMock struct {
	ListFunc     func(orgID string) ([]mongodbatlas.Team, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	GetOneTeamByNameFunc     func(orgID string, name string) (*mongodbatlas.Team, *mongodbatlas.Response, error)
	GetOneTeamByNameRequests map[string]struct{}

	GetTeamUsersAssignedFunc     func(orgID string, teamID string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error)
	GetTeamUsersAssignedRequests map[string]struct{}

	CreateFunc     func(orgID string, team *mongodbatlas.Team) (*mongodbatlas.Team, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.Team

	RenameFunc     func(orgID string, teamID string, name string) (*mongodbatlas.Team, *mongodbatlas.Response, error)
	RenameRequests map[string]struct{}

	UpdateTeamRolesFunc     func(orgID string, teamID string, roles *mongodbatlas.TeamUpdateRoles) ([]mongodbatlas.TeamRoles, *mongodbatlas.Response, error)
	UpdateTeamRolesRequests map[string]*mongodbatlas.TeamUpdateRoles

	AddUsersToTeamFunc     func(orgID string, teamID string, userIDs []string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error)
	AddUsersToTeamRequests map[string][]string

	RemoveUserToTeamFunc     func(orgID string, teamID string, userID string) (*mongodbatlas.Response, error)
	RemoveUserToTeamRequests map[string]struct{}

	RemoveTeamFromOrganizationFunc     func(orgID string, teamID string) (*mongodbatlas.Response, error)
	RemoveTeamFromOrganizationRequests map[string]struct{}

	RemoveTeamFromProjectFunc     func(projectID string, teamID string) (*mongodbatlas.Response, error)
	RemoveTeamFromProjectRequests map[string]struct{}
}

func (c *TeamsClientMock) List(_ context.Context, orgID string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.Team, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[orgID] = struct{}{}

	return c.ListFunc(orgID)
}

func (c *TeamsClientMock) Get(_ context.Context, orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", orgID, teamID)] = struct{}{}

	return c.GetFunc(orgID, teamID)
}

func (c *TeamsClientMock) GetOneTeamByName(_ context.Context, orgID string, name string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", orgID, name)] = struct{}{}

	return c.GetFunc(orgID, name)
}

func (c *TeamsClientMock) GetTeamUsersAssigned(_ context.Context, orgID string, teamID string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error) {
	if c.GetTeamUsersAssignedRequests == nil {
		c.GetTeamUsersAssignedRequests = map[string]struct{}{}
	}

	c.GetTeamUsersAssignedRequests[fmt.Sprintf("%s.%s", orgID, teamID)] = struct{}{}

	return c.GetTeamUsersAssignedFunc(orgID, teamID)
}

func (c *TeamsClientMock) Create(_ context.Context, orgID string, team *mongodbatlas.Team) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.Team{}
	}

	c.CreateRequests[orgID] = team

	return c.CreateFunc(orgID, team)
}

func (c *TeamsClientMock) Rename(_ context.Context, orgID string, teamID string, name string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	if c.RenameRequests == nil {
		c.RenameRequests = map[string]struct{}{}
	}

	c.RenameRequests[fmt.Sprintf("%s.%s.%s", orgID, teamID, name)] = struct{}{}

	return c.RenameFunc(orgID, teamID, name)
}

func (c *TeamsClientMock) UpdateTeamRoles(_ context.Context, orgID string, teamID string, roles *mongodbatlas.TeamUpdateRoles) ([]mongodbatlas.TeamRoles, *mongodbatlas.Response, error) {
	if c.UpdateTeamRolesRequests == nil {
		c.UpdateTeamRolesRequests = map[string]*mongodbatlas.TeamUpdateRoles{}
	}

	c.UpdateTeamRolesRequests[fmt.Sprintf("%s.%s", orgID, teamID)] = roles

	return c.UpdateTeamRolesFunc(orgID, teamID, roles)
}

func (c *TeamsClientMock) AddUsersToTeam(_ context.Context, orgID string, teamID string, userIDs []string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error) {
	if c.AddUsersToTeamRequests == nil {
		c.AddUsersToTeamRequests = map[string][]string{}
	}

	c.AddUsersToTeamRequests[fmt.Sprintf("%s.%s", orgID, teamID)] = userIDs

	return c.AddUsersToTeamFunc(orgID, teamID, userIDs)
}

func (c *TeamsClientMock) RemoveUserToTeam(_ context.Context, orgID string, teamID string, userID string) (*mongodbatlas.Response, error) {
	if c.RemoveUserToTeamRequests == nil {
		c.RemoveUserToTeamRequests = map[string]struct{}{}
	}

	c.RemoveUserToTeamRequests[fmt.Sprintf("%s.%s.%s", orgID, teamID, userID)] = struct{}{}

	return c.RemoveUserToTeamFunc(orgID, teamID, userID)
}

func (c *TeamsClientMock) RemoveTeamFromOrganization(_ context.Context, orgID string, teamID string) (*mongodbatlas.Response, error) {
	if c.RemoveTeamFromOrganizationRequests == nil {
		c.RemoveTeamFromOrganizationRequests = map[string]struct{}{}
	}

	c.RemoveTeamFromOrganizationRequests[fmt.Sprintf("%s.%s", orgID, teamID)] = struct{}{}

	return c.RemoveTeamFromOrganizationFunc(orgID, teamID)
}

func (c *TeamsClientMock) RemoveTeamFromProject(_ context.Context, projectID string, teamID string) (*mongodbatlas.Response, error) {
	if c.RemoveTeamFromProjectRequests == nil {
		c.RemoveTeamFromProjectRequests = map[string]struct{}{}
	}

	c.RemoveTeamFromProjectRequests[fmt.Sprintf("%s.%s", projectID, teamID)] = struct{}{}

	return c.RemoveTeamFromProjectFunc(projectID, teamID)
}
