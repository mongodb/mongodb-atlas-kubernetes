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

type ProjectsClientMock struct {
	GetAllProjectsFunc  func() (*mongodbatlas.Projects, *mongodbatlas.Response, error)
	GetAllProjectsCalls int

	GetOneProjectFunc     func(projectID string) (*mongodbatlas.Project, *mongodbatlas.Response, error)
	GetOneProjectRequests map[string]struct{}

	GetOneProjectByNameFunc     func(projectName string) (*mongodbatlas.Project, *mongodbatlas.Response, error)
	GetOneProjectByNameRequests map[string]struct{}

	CreateFunc     func(project *mongodbatlas.Project) (*mongodbatlas.Project, *mongodbatlas.Response, error)
	CreateRequests []*mongodbatlas.Project

	UpdateFunc     func(projectID string, project *mongodbatlas.ProjectUpdateRequest) (*mongodbatlas.Project, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.ProjectUpdateRequest

	DeleteFunc     func(projectID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}

	GetProjectTeamsAssignedFunc     func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error)
	GetProjectTeamsAssignedRequests map[string]struct{}

	AddTeamsToProjectFunc     func(projectId string, teams []*mongodbatlas.ProjectTeam) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error)
	AddTeamsToProjectRequests map[string][]*mongodbatlas.ProjectTeam

	RemoveUserFromProjectFunc     func(projectID string, userID string) (*mongodbatlas.Response, error)
	RemoveUserFromProjectRequests map[string]struct{}

	InvitationsFunc     func(projectID string, invitations *mongodbatlas.InvitationOptions) ([]*mongodbatlas.Invitation, *mongodbatlas.Response, error)
	InvitationsRequests map[string]*mongodbatlas.InvitationOptions

	InvitationFunc     func(projectID string, invitationID string) (*mongodbatlas.Invitation, *mongodbatlas.Response, error)
	InvitationRequests map[string]struct{}

	InviteUserFunc     func(projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error)
	InviteUserRequests map[string]*mongodbatlas.Invitation

	UpdateInvitationFunc     func(projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error)
	UpdateInvitationRequests map[string]*mongodbatlas.Invitation

	UpdateInvitationByIDFunc     func(projectID string, invitationID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error)
	UpdateInvitationByIDRequests map[string]*mongodbatlas.Invitation

	DeleteInvitationFunc     func(projectID string, invitationID string) (*mongodbatlas.Response, error)
	DeleteInvitationRequests map[string]struct{}

	GetProjectSettingsFunc     func(projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error)
	GetProjectSettingsRequests map[string]struct{}

	UpdateProjectSettingsFunc     func(projectID string, settings *mongodbatlas.ProjectSettings) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error)
	UpdateProjectSettingsRequests map[string]*mongodbatlas.ProjectSettings
}

func (c *ProjectsClientMock) GetAllProjects(_ context.Context, _ *mongodbatlas.ListOptions) (*mongodbatlas.Projects, *mongodbatlas.Response, error) {
	c.GetAllProjectsCalls++

	return c.GetAllProjectsFunc()
}

func (c *ProjectsClientMock) GetOneProject(_ context.Context, projectID string) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.GetOneProjectRequests == nil {
		c.GetOneProjectRequests = map[string]struct{}{}
	}

	c.GetOneProjectRequests[projectID] = struct{}{}

	return c.GetOneProjectFunc(projectID)
}

func (c *ProjectsClientMock) GetOneProjectByName(_ context.Context, projectName string) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.GetOneProjectByNameRequests == nil {
		c.GetOneProjectByNameRequests = map[string]struct{}{}
	}

	c.GetOneProjectByNameRequests[projectName] = struct{}{}

	return c.GetOneProjectFunc(projectName)
}

func (c *ProjectsClientMock) Create(_ context.Context, project *mongodbatlas.Project, _ *mongodbatlas.CreateProjectOptions) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = []*mongodbatlas.Project{}
	}

	c.CreateRequests = append(c.CreateRequests, project)

	return c.CreateFunc(project)
}

func (c *ProjectsClientMock) Update(_ context.Context, projectID string, project *mongodbatlas.ProjectUpdateRequest) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.ProjectUpdateRequest{}
	}

	c.UpdateRequests[projectID] = project

	return c.UpdateFunc(projectID, project)
}

func (c *ProjectsClientMock) Delete(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[projectID] = struct{}{}

	return c.DeleteFunc(projectID)
}

func (c *ProjectsClientMock) GetProjectTeamsAssigned(_ context.Context, projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
	if c.GetProjectTeamsAssignedRequests == nil {
		c.GetProjectTeamsAssignedRequests = map[string]struct{}{}
	}

	c.GetProjectTeamsAssignedRequests[projectID] = struct{}{}

	return c.GetProjectTeamsAssignedFunc(projectID)
}

func (c *ProjectsClientMock) AddTeamsToProject(_ context.Context, projectID string, teams []*mongodbatlas.ProjectTeam) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
	if c.AddTeamsToProjectRequests == nil {
		c.AddTeamsToProjectRequests = map[string][]*mongodbatlas.ProjectTeam{}
	}

	c.AddTeamsToProjectRequests[projectID] = teams

	return c.AddTeamsToProjectFunc(projectID, teams)
}

func (c *ProjectsClientMock) RemoveUserFromProject(_ context.Context, projectID string, userID string) (*mongodbatlas.Response, error) {
	if c.RemoveUserFromProjectRequests == nil {
		c.RemoveUserFromProjectRequests = map[string]struct{}{}
	}

	c.RemoveUserFromProjectRequests[fmt.Sprintf("%s.%s", projectID, userID)] = struct{}{}

	return c.RemoveUserFromProjectFunc(projectID, userID)
}

func (c *ProjectsClientMock) Invitations(_ context.Context, projectID string, invitation *mongodbatlas.InvitationOptions) ([]*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InvitationsRequests == nil {
		c.InvitationsRequests = map[string]*mongodbatlas.InvitationOptions{}
	}

	c.InvitationsRequests[projectID] = invitation

	return c.InvitationsFunc(projectID, invitation)
}

func (c *ProjectsClientMock) Invitation(_ context.Context, projectID string, invitationID string) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InvitationRequests == nil {
		c.InvitationRequests = map[string]struct{}{}
	}

	c.InvitationRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = struct{}{}

	return c.InvitationFunc(projectID, invitationID)
}

func (c *ProjectsClientMock) InviteUser(_ context.Context, projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InviteUserRequests == nil {
		c.InviteUserRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.InviteUserRequests[projectID] = invitation

	return c.InviteUserFunc(projectID, invitation)
}

func (c *ProjectsClientMock) UpdateInvitation(_ context.Context, projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.UpdateInvitationRequests == nil {
		c.UpdateInvitationRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.UpdateInvitationRequests[projectID] = invitation

	return c.UpdateInvitationFunc(projectID, invitation)
}

func (c *ProjectsClientMock) UpdateInvitationByID(_ context.Context, projectID string, invitationID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.UpdateInvitationByIDRequests == nil {
		c.UpdateInvitationByIDRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.UpdateInvitationByIDRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = invitation

	return c.UpdateInvitationByIDFunc(projectID, invitationID, invitation)
}

func (c *ProjectsClientMock) DeleteInvitation(_ context.Context, projectID string, invitationID string) (*mongodbatlas.Response, error) {
	if c.DeleteInvitationRequests == nil {
		c.DeleteInvitationRequests = map[string]struct{}{}
	}

	c.DeleteInvitationRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = struct{}{}

	return c.DeleteInvitationFunc(projectID, invitationID)
}

func (c *ProjectsClientMock) GetProjectSettings(_ context.Context, projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
	if c.GetProjectSettingsRequests == nil {
		c.GetProjectSettingsRequests = map[string]struct{}{}
	}

	c.GetProjectSettingsRequests[projectID] = struct{}{}

	return c.GetProjectSettingsFunc(projectID)
}

func (c *ProjectsClientMock) UpdateProjectSettings(_ context.Context, projectID string, settings *mongodbatlas.ProjectSettings) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
	if c.UpdateProjectSettingsRequests == nil {
		c.UpdateProjectSettingsRequests = map[string]*mongodbatlas.ProjectSettings{}
	}

	c.UpdateProjectSettingsRequests[projectID] = settings

	return c.UpdateProjectSettingsFunc(projectID, settings)
}
