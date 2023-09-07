package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type MockProjectsClient struct {
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

func (c *MockProjectsClient) GetAllProjects(_ context.Context, _ *mongodbatlas.ListOptions) (*mongodbatlas.Projects, *mongodbatlas.Response, error) {
	c.GetAllProjectsCalls++

	return c.GetAllProjectsFunc()
}

func (c *MockProjectsClient) GetOneProject(_ context.Context, projectID string) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.GetOneProjectRequests == nil {
		c.GetOneProjectRequests = map[string]struct{}{}
	}

	c.GetOneProjectRequests[projectID] = struct{}{}

	return c.GetOneProjectFunc(projectID)
}

func (c *MockProjectsClient) GetOneProjectByName(_ context.Context, projectName string) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.GetOneProjectByNameRequests == nil {
		c.GetOneProjectByNameRequests = map[string]struct{}{}
	}

	c.GetOneProjectByNameRequests[projectName] = struct{}{}

	return c.GetOneProjectFunc(projectName)
}

func (c *MockProjectsClient) Create(_ context.Context, project *mongodbatlas.Project, _ *mongodbatlas.CreateProjectOptions) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = []*mongodbatlas.Project{}
	}

	c.CreateRequests = append(c.CreateRequests, project)

	return c.CreateFunc(project)
}

func (c *MockProjectsClient) Update(_ context.Context, projectID string, project *mongodbatlas.ProjectUpdateRequest) (*mongodbatlas.Project, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.ProjectUpdateRequest{}
	}

	c.UpdateRequests[projectID] = project

	return c.UpdateFunc(projectID, project)
}

func (c *MockProjectsClient) Delete(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[projectID] = struct{}{}

	return c.DeleteFunc(projectID)
}

func (c *MockProjectsClient) GetProjectTeamsAssigned(_ context.Context, projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
	if c.GetProjectTeamsAssignedRequests == nil {
		c.GetProjectTeamsAssignedRequests = map[string]struct{}{}
	}

	c.GetProjectTeamsAssignedRequests[projectID] = struct{}{}

	return c.GetProjectTeamsAssignedFunc(projectID)
}

func (c *MockProjectsClient) AddTeamsToProject(_ context.Context, projectID string, teams []*mongodbatlas.ProjectTeam) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
	if c.AddTeamsToProjectRequests == nil {
		c.AddTeamsToProjectRequests = map[string][]*mongodbatlas.ProjectTeam{}
	}

	c.AddTeamsToProjectRequests[projectID] = teams

	return c.AddTeamsToProjectFunc(projectID, teams)
}

func (c *MockProjectsClient) RemoveUserFromProject(_ context.Context, projectID string, userID string) (*mongodbatlas.Response, error) {
	if c.RemoveUserFromProjectRequests == nil {
		c.RemoveUserFromProjectRequests = map[string]struct{}{}
	}

	c.RemoveUserFromProjectRequests[fmt.Sprintf("%s.%s", projectID, userID)] = struct{}{}

	return c.RemoveUserFromProjectFunc(projectID, userID)
}

func (c *MockProjectsClient) Invitations(_ context.Context, projectID string, invitation *mongodbatlas.InvitationOptions) ([]*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InvitationsRequests == nil {
		c.InvitationsRequests = map[string]*mongodbatlas.InvitationOptions{}
	}

	c.InvitationsRequests[projectID] = invitation

	return c.InvitationsFunc(projectID, invitation)
}

func (c *MockProjectsClient) Invitation(_ context.Context, projectID string, invitationID string) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InvitationRequests == nil {
		c.InvitationRequests = map[string]struct{}{}
	}

	c.InvitationRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = struct{}{}

	return c.InvitationFunc(projectID, invitationID)
}

func (c *MockProjectsClient) InviteUser(_ context.Context, projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.InviteUserRequests == nil {
		c.InviteUserRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.InviteUserRequests[projectID] = invitation

	return c.InviteUserFunc(projectID, invitation)
}

func (c *MockProjectsClient) UpdateInvitation(_ context.Context, projectID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.UpdateInvitationRequests == nil {
		c.UpdateInvitationRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.UpdateInvitationRequests[projectID] = invitation

	return c.UpdateInvitationFunc(projectID, invitation)
}

func (c *MockProjectsClient) UpdateInvitationByID(_ context.Context, projectID string, invitationID string, invitation *mongodbatlas.Invitation) (*mongodbatlas.Invitation, *mongodbatlas.Response, error) {
	if c.UpdateInvitationByIDRequests == nil {
		c.UpdateInvitationByIDRequests = map[string]*mongodbatlas.Invitation{}
	}

	c.UpdateInvitationByIDRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = invitation

	return c.UpdateInvitationByIDFunc(projectID, invitationID, invitation)
}

func (c *MockProjectsClient) DeleteInvitation(_ context.Context, projectID string, invitationID string) (*mongodbatlas.Response, error) {
	if c.DeleteInvitationRequests == nil {
		c.DeleteInvitationRequests = map[string]struct{}{}
	}

	c.DeleteInvitationRequests[fmt.Sprintf("%s.%s", projectID, invitationID)] = struct{}{}

	return c.DeleteInvitationFunc(projectID, invitationID)
}

func (c *MockProjectsClient) GetProjectSettings(_ context.Context, projectID string) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
	if c.GetProjectSettingsRequests == nil {
		c.GetProjectSettingsRequests = map[string]struct{}{}
	}

	c.GetProjectSettingsRequests[projectID] = struct{}{}

	return c.GetProjectSettingsFunc(projectID)
}

func (c *MockProjectsClient) UpdateProjectSettings(_ context.Context, projectID string, settings *mongodbatlas.ProjectSettings) (*mongodbatlas.ProjectSettings, *mongodbatlas.Response, error) {
	if c.UpdateProjectSettingsRequests == nil {
		c.UpdateProjectSettingsRequests = map[string]*mongodbatlas.ProjectSettings{}
	}

	c.UpdateProjectSettingsRequests[projectID] = settings

	return c.UpdateProjectSettingsFunc(projectID, settings)
}
