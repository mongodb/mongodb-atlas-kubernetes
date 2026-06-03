// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TeamsApi interface {

	/*
		AddGroupTeams Add Multiple Teams to One Project

		Adds multiple teams to the specified project. All members of a team share the same project access. MongoDB Cloud limits the number of users to a maximum of 100 teams per project and a maximum of 250 teams per organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param teamRole Teams and their roles to be added to the specified project.
		@return AddGroupTeamsApiRequest
	*/
	AddGroupTeams(ctx context.Context, groupId string, teamRole *[]TeamRole) AddGroupTeamsApiRequest
	/*
		AddGroupTeams Add Multiple Teams to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddGroupTeamsApiParams - Parameters for the request
		@return AddGroupTeamsApiRequest
	*/
	AddGroupTeamsWithParams(ctx context.Context, args *AddGroupTeamsApiParams) AddGroupTeamsApiRequest

	// Method available only for mocking purposes
	AddGroupTeamsExecute(r AddGroupTeamsApiRequest) (*PaginatedTeamRole, *http.Response, error)

	/*
			AddTeamUsers Assign MongoDB Cloud Users in One Organization to One Team

			Adds one or more MongoDB Cloud users from the specified organization to the specified team. Teams enable you to grant project access roles to MongoDB Cloud users. You can assign up to 250 MongoDB Cloud users from one organization to one team.

		**Note**: This endpoint is deprecated. Use [Add One MongoDB Cloud User to One Team](#tag/MongoDB-Cloud-Users/operation/addUserToTeam) to add an active or pending user to a team.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param teamId Unique 24-hexadecimal character string that identifies the team to which you want to add MongoDB Cloud users.
			@param addUserToTeam One or more MongoDB Cloud users that you want to add to the specified team.
			@return AddTeamUsersApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for TeamsApi
	*/
	AddTeamUsers(ctx context.Context, orgId string, teamId string, addUserToTeam *[]AddUserToTeam) AddTeamUsersApiRequest
	/*
		AddTeamUsers Assign MongoDB Cloud Users in One Organization to One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddTeamUsersApiParams - Parameters for the request
		@return AddTeamUsersApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for TeamsApi
	*/
	AddTeamUsersWithParams(ctx context.Context, args *AddTeamUsersApiParams) AddTeamUsersApiRequest

	// Method available only for mocking purposes
	AddTeamUsersExecute(r AddTeamUsersApiRequest) (*PaginatedApiAppUser, *http.Response, error)

	/*
		CreateOrgTeam Create One Team in One Organization

		Creates one team in the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users. MongoDB Cloud limits the number of teams to a maximum of 250 teams per organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param team Team that you want to create in the specified organization.
		@return CreateOrgTeamApiRequest
	*/
	CreateOrgTeam(ctx context.Context, orgId string, team *Team) CreateOrgTeamApiRequest
	/*
		CreateOrgTeam Create One Team in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgTeamApiParams - Parameters for the request
		@return CreateOrgTeamApiRequest
	*/
	CreateOrgTeamWithParams(ctx context.Context, args *CreateOrgTeamApiParams) CreateOrgTeamApiRequest

	// Method available only for mocking purposes
	CreateOrgTeamExecute(r CreateOrgTeamApiRequest) (*Team, *http.Response, error)

	/*
		DeleteOrgTeam Remove One Team from One Organization

		Removes one team specified using its unique 24-hexadecimal digit identifier from the organization specified using its unique 24-hexadecimal digit identifier.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to delete.
		@return DeleteOrgTeamApiRequest
	*/
	DeleteOrgTeam(ctx context.Context, orgId string, teamId string) DeleteOrgTeamApiRequest
	/*
		DeleteOrgTeam Remove One Team from One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgTeamApiParams - Parameters for the request
		@return DeleteOrgTeamApiRequest
	*/
	DeleteOrgTeamWithParams(ctx context.Context, args *DeleteOrgTeamApiParams) DeleteOrgTeamApiRequest

	// Method available only for mocking purposes
	DeleteOrgTeamExecute(r DeleteOrgTeamApiRequest) (*http.Response, error)

	/*
		GetGroupTeam Return One Team in One Project

		Returns one team to which the authenticated user has access in the project specified using its unique 24-hexadecimal digit identifier. All members of the team share the same project access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param teamId Unique 24-hexadecimal digit string that identifies the team for which you want to get.
		@return GetGroupTeamApiRequest
	*/
	GetGroupTeam(ctx context.Context, groupId string, teamId string) GetGroupTeamApiRequest
	/*
		GetGroupTeam Return One Team in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupTeamApiParams - Parameters for the request
		@return GetGroupTeamApiRequest
	*/
	GetGroupTeamWithParams(ctx context.Context, args *GetGroupTeamApiParams) GetGroupTeamApiRequest

	// Method available only for mocking purposes
	GetGroupTeamExecute(r GetGroupTeamApiRequest) (*TeamRole, *http.Response, error)

	/*
		GetOrgTeam Return One Team by ID

		Returns one team that you identified using its unique 24-hexadecimal digit ID. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param teamId Unique 24-hexadecimal digit string that identifies the team whose information you want to return.
		@return GetOrgTeamApiRequest
	*/
	GetOrgTeam(ctx context.Context, orgId string, teamId string) GetOrgTeamApiRequest
	/*
		GetOrgTeam Return One Team by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgTeamApiParams - Parameters for the request
		@return GetOrgTeamApiRequest
	*/
	GetOrgTeamWithParams(ctx context.Context, args *GetOrgTeamApiParams) GetOrgTeamApiRequest

	// Method available only for mocking purposes
	GetOrgTeamExecute(r GetOrgTeamApiRequest) (*TeamResponse, *http.Response, error)

	/*
		GetTeamByName Return One Team by Name

		Returns one team that you identified using its human-readable name. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param teamName Name of the team whose information you want to return.
		@return GetTeamByNameApiRequest
	*/
	GetTeamByName(ctx context.Context, orgId string, teamName string) GetTeamByNameApiRequest
	/*
		GetTeamByName Return One Team by Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetTeamByNameApiParams - Parameters for the request
		@return GetTeamByNameApiRequest
	*/
	GetTeamByNameWithParams(ctx context.Context, args *GetTeamByNameApiParams) GetTeamByNameApiRequest

	// Method available only for mocking purposes
	GetTeamByNameExecute(r GetTeamByNameApiRequest) (*TeamResponse, *http.Response, error)

	/*
		ListGroupTeams Return All Teams in One Project

		Returns all teams to which the authenticated user has access in the project specified using its unique 24-hexadecimal digit identifier. All members of the team share the same project access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupTeamsApiRequest
	*/
	ListGroupTeams(ctx context.Context, groupId string) ListGroupTeamsApiRequest
	/*
		ListGroupTeams Return All Teams in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupTeamsApiParams - Parameters for the request
		@return ListGroupTeamsApiRequest
	*/
	ListGroupTeamsWithParams(ctx context.Context, args *ListGroupTeamsApiParams) ListGroupTeamsApiRequest

	// Method available only for mocking purposes
	ListGroupTeamsExecute(r ListGroupTeamsApiRequest) (*PaginatedTeamRole, *http.Response, error)

	/*
		ListOrgTeams Return All Teams in One Organization

		Returns all teams that belong to the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users. MongoDB Cloud only returns teams for which you have access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListOrgTeamsApiRequest
	*/
	ListOrgTeams(ctx context.Context, orgId string) ListOrgTeamsApiRequest
	/*
		ListOrgTeams Return All Teams in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgTeamsApiParams - Parameters for the request
		@return ListOrgTeamsApiRequest
	*/
	ListOrgTeamsWithParams(ctx context.Context, args *ListOrgTeamsApiParams) ListOrgTeamsApiRequest

	// Method available only for mocking purposes
	ListOrgTeamsExecute(r ListOrgTeamsApiRequest) (*PaginatedTeam, *http.Response, error)

	/*
		RemoveGroupTeam Remove One Team from One Project

		Removes one team specified using its unique 24-hexadecimal digit identifier from the project specified using its unique 24-hexadecimal digit identifier.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to remove from the specified project.
		@return RemoveGroupTeamApiRequest
	*/
	RemoveGroupTeam(ctx context.Context, groupId string, teamId string) RemoveGroupTeamApiRequest
	/*
		RemoveGroupTeam Remove One Team from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveGroupTeamApiParams - Parameters for the request
		@return RemoveGroupTeamApiRequest
	*/
	RemoveGroupTeamWithParams(ctx context.Context, args *RemoveGroupTeamApiParams) RemoveGroupTeamApiRequest

	// Method available only for mocking purposes
	RemoveGroupTeamExecute(r RemoveGroupTeamApiRequest) (*http.Response, error)

	/*
			RemoveUserFromTeam Remove One MongoDB Cloud User from One Team

			Removes one MongoDB Cloud user from the specified team. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

		**Note**: This endpoint is deprecated. Use [Remove One MongoDB Cloud User from One Team](#tag/MongoDB-Cloud-Users/operation/removeUserFromTeam) to remove an active or pending user from a team.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param teamId Unique 24-hexadecimal digit string that identifies the team from which you want to remove one database application user.
			@param userId Unique 24-hexadecimal digit string that identifies MongoDB Cloud user that you want to remove from the specified team.
			@return RemoveUserFromTeamApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for TeamsApi
	*/
	RemoveUserFromTeam(ctx context.Context, orgId string, teamId string, userId string) RemoveUserFromTeamApiRequest
	/*
		RemoveUserFromTeam Remove One MongoDB Cloud User from One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveUserFromTeamApiParams - Parameters for the request
		@return RemoveUserFromTeamApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for TeamsApi
	*/
	RemoveUserFromTeamWithParams(ctx context.Context, args *RemoveUserFromTeamApiParams) RemoveUserFromTeamApiRequest

	// Method available only for mocking purposes
	RemoveUserFromTeamExecute(r RemoveUserFromTeamApiRequest) (*http.Response, error)

	/*
		RenameOrgTeam Rename One Team

		Renames one team in the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to rename.
		@param teamUpdate Details to update on the specified team.
		@return RenameOrgTeamApiRequest
	*/
	RenameOrgTeam(ctx context.Context, orgId string, teamId string, teamUpdate *TeamUpdate) RenameOrgTeamApiRequest
	/*
		RenameOrgTeam Rename One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RenameOrgTeamApiParams - Parameters for the request
		@return RenameOrgTeamApiRequest
	*/
	RenameOrgTeamWithParams(ctx context.Context, args *RenameOrgTeamApiParams) RenameOrgTeamApiRequest

	// Method available only for mocking purposes
	RenameOrgTeamExecute(r RenameOrgTeamApiRequest) (*TeamResponse, *http.Response, error)

	/*
		UpdateGroupTeam Update Team Roles in One Project

		Updates the project roles assigned to the specified team. You can grant team roles for specific projects and grant project access roles to users in the team. All members of the team share the same project access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param teamId Unique 24-hexadecimal digit string that identifies the team for which you want to update roles.
		@param teamRole The project roles assigned to the specified team.
		@return UpdateGroupTeamApiRequest
	*/
	UpdateGroupTeam(ctx context.Context, groupId string, teamId string, teamRole *TeamRole) UpdateGroupTeamApiRequest
	/*
		UpdateGroupTeam Update Team Roles in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupTeamApiParams - Parameters for the request
		@return UpdateGroupTeamApiRequest
	*/
	UpdateGroupTeamWithParams(ctx context.Context, args *UpdateGroupTeamApiParams) UpdateGroupTeamApiRequest

	// Method available only for mocking purposes
	UpdateGroupTeamExecute(r UpdateGroupTeamApiRequest) (*PaginatedTeamRole, *http.Response, error)
}

// TeamsApiService TeamsApi service
type TeamsApiService service

type AddGroupTeamsApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	groupId    string
	teamRole   *[]TeamRole
}

type AddGroupTeamsApiParams struct {
	GroupId  string
	TeamRole *[]TeamRole
}

func (a *TeamsApiService) AddGroupTeamsWithParams(ctx context.Context, args *AddGroupTeamsApiParams) AddGroupTeamsApiRequest {
	return AddGroupTeamsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		teamRole:   args.TeamRole,
	}
}

func (r AddGroupTeamsApiRequest) Execute() (*PaginatedTeamRole, *http.Response, error) {
	return r.ApiService.AddGroupTeamsExecute(r)
}

/*
AddGroupTeams Add Multiple Teams to One Project

Adds multiple teams to the specified project. All members of a team share the same project access. MongoDB Cloud limits the number of users to a maximum of 100 teams per project and a maximum of 250 teams per organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return AddGroupTeamsApiRequest
*/
func (a *TeamsApiService) AddGroupTeams(ctx context.Context, groupId string, teamRole *[]TeamRole) AddGroupTeamsApiRequest {
	return AddGroupTeamsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		teamRole:   teamRole,
	}
}

// AddGroupTeamsExecute executes the request
//
//	@return PaginatedTeamRole
func (a *TeamsApiService) AddGroupTeamsExecute(r AddGroupTeamsApiRequest) (*PaginatedTeamRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedTeamRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.AddGroupTeams")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/teams"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.teamRole == nil {
		return localVarReturnValue, nil, reportError("teamRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.teamRole
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type AddTeamUsersApiRequest struct {
	ctx           context.Context
	ApiService    TeamsApi
	orgId         string
	teamId        string
	addUserToTeam *[]AddUserToTeam
}

type AddTeamUsersApiParams struct {
	OrgId         string
	TeamId        string
	AddUserToTeam *[]AddUserToTeam
}

func (a *TeamsApiService) AddTeamUsersWithParams(ctx context.Context, args *AddTeamUsersApiParams) AddTeamUsersApiRequest {
	return AddTeamUsersApiRequest{
		ApiService:    a,
		ctx:           ctx,
		orgId:         args.OrgId,
		teamId:        args.TeamId,
		addUserToTeam: args.AddUserToTeam,
	}
}

func (r AddTeamUsersApiRequest) Execute() (*PaginatedApiAppUser, *http.Response, error) {
	return r.ApiService.AddTeamUsersExecute(r)
}

/*
AddTeamUsers Assign MongoDB Cloud Users in One Organization to One Team

Adds one or more MongoDB Cloud users from the specified organization to the specified team. Teams enable you to grant project access roles to MongoDB Cloud users. You can assign up to 250 MongoDB Cloud users from one organization to one team.

**Note**: This endpoint is deprecated. Use [Add One MongoDB Cloud User to One Team](#tag/MongoDB-Cloud-Users/operation/addUserToTeam) to add an active or pending user to a team.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal character string that identifies the team to which you want to add MongoDB Cloud users.
	@return AddTeamUsersApiRequest

Deprecated
*/
func (a *TeamsApiService) AddTeamUsers(ctx context.Context, orgId string, teamId string, addUserToTeam *[]AddUserToTeam) AddTeamUsersApiRequest {
	return AddTeamUsersApiRequest{
		ApiService:    a,
		ctx:           ctx,
		orgId:         orgId,
		teamId:        teamId,
		addUserToTeam: addUserToTeam,
	}
}

// AddTeamUsersExecute executes the request
//
//	@return PaginatedApiAppUser
//
// Deprecated
func (a *TeamsApiService) AddTeamUsersExecute(r AddTeamUsersApiRequest) (*PaginatedApiAppUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAppUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.AddTeamUsers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}/users"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamId == "" {
		return localVarReturnValue, nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.addUserToTeam == nil {
		return localVarReturnValue, nil, reportError("addUserToTeam is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addUserToTeam
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type CreateOrgTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	team       *Team
}

type CreateOrgTeamApiParams struct {
	OrgId string
	Team  *Team
}

func (a *TeamsApiService) CreateOrgTeamWithParams(ctx context.Context, args *CreateOrgTeamApiParams) CreateOrgTeamApiRequest {
	return CreateOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		team:       args.Team,
	}
}

func (r CreateOrgTeamApiRequest) Execute() (*Team, *http.Response, error) {
	return r.ApiService.CreateOrgTeamExecute(r)
}

/*
CreateOrgTeam Create One Team in One Organization

Creates one team in the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users. MongoDB Cloud limits the number of teams to a maximum of 250 teams per organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgTeamApiRequest
*/
func (a *TeamsApiService) CreateOrgTeam(ctx context.Context, orgId string, team *Team) CreateOrgTeamApiRequest {
	return CreateOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		team:       team,
	}
}

// CreateOrgTeamExecute executes the request
//
//	@return Team
func (a *TeamsApiService) CreateOrgTeamExecute(r CreateOrgTeamApiRequest) (*Team, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Team
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.CreateOrgTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.team == nil {
		return localVarReturnValue, nil, reportError("team is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.team
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type DeleteOrgTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	teamId     string
}

type DeleteOrgTeamApiParams struct {
	OrgId  string
	TeamId string
}

func (a *TeamsApiService) DeleteOrgTeamWithParams(ctx context.Context, args *DeleteOrgTeamApiParams) DeleteOrgTeamApiRequest {
	return DeleteOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		teamId:     args.TeamId,
	}
}

func (r DeleteOrgTeamApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgTeamExecute(r)
}

/*
DeleteOrgTeam Remove One Team from One Organization

Removes one team specified using its unique 24-hexadecimal digit identifier from the organization specified using its unique 24-hexadecimal digit identifier.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to delete.
	@return DeleteOrgTeamApiRequest
*/
func (a *TeamsApiService) DeleteOrgTeam(ctx context.Context, orgId string, teamId string) DeleteOrgTeamApiRequest {
	return DeleteOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamId:     teamId,
	}
}

// DeleteOrgTeamExecute executes the request
func (a *TeamsApiService) DeleteOrgTeamExecute(r DeleteOrgTeamApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.DeleteOrgTeam")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamId == "" {
		return nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type GetGroupTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	groupId    string
	teamId     string
}

type GetGroupTeamApiParams struct {
	GroupId string
	TeamId  string
}

func (a *TeamsApiService) GetGroupTeamWithParams(ctx context.Context, args *GetGroupTeamApiParams) GetGroupTeamApiRequest {
	return GetGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		teamId:     args.TeamId,
	}
}

func (r GetGroupTeamApiRequest) Execute() (*TeamRole, *http.Response, error) {
	return r.ApiService.GetGroupTeamExecute(r)
}

/*
GetGroupTeam Return One Team in One Project

Returns one team to which the authenticated user has access in the project specified using its unique 24-hexadecimal digit identifier. All members of the team share the same project access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param teamId Unique 24-hexadecimal digit string that identifies the team for which you want to get.
	@return GetGroupTeamApiRequest
*/
func (a *TeamsApiService) GetGroupTeam(ctx context.Context, groupId string, teamId string) GetGroupTeamApiRequest {
	return GetGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		teamId:     teamId,
	}
}

// GetGroupTeamExecute executes the request
//
//	@return TeamRole
func (a *TeamsApiService) GetGroupTeamExecute(r GetGroupTeamApiRequest) (*TeamRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *TeamRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.GetGroupTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/teams/{teamId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.teamId == "" {
		return localVarReturnValue, nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type GetOrgTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	teamId     string
}

type GetOrgTeamApiParams struct {
	OrgId  string
	TeamId string
}

func (a *TeamsApiService) GetOrgTeamWithParams(ctx context.Context, args *GetOrgTeamApiParams) GetOrgTeamApiRequest {
	return GetOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		teamId:     args.TeamId,
	}
}

func (r GetOrgTeamApiRequest) Execute() (*TeamResponse, *http.Response, error) {
	return r.ApiService.GetOrgTeamExecute(r)
}

/*
GetOrgTeam Return One Team by ID

Returns one team that you identified using its unique 24-hexadecimal digit ID. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team whose information you want to return.
	@return GetOrgTeamApiRequest
*/
func (a *TeamsApiService) GetOrgTeam(ctx context.Context, orgId string, teamId string) GetOrgTeamApiRequest {
	return GetOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamId:     teamId,
	}
}

// GetOrgTeamExecute executes the request
//
//	@return TeamResponse
func (a *TeamsApiService) GetOrgTeamExecute(r GetOrgTeamApiRequest) (*TeamResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *TeamResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.GetOrgTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamId == "" {
		return localVarReturnValue, nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type GetTeamByNameApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	teamName   string
}

type GetTeamByNameApiParams struct {
	OrgId    string
	TeamName string
}

func (a *TeamsApiService) GetTeamByNameWithParams(ctx context.Context, args *GetTeamByNameApiParams) GetTeamByNameApiRequest {
	return GetTeamByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		teamName:   args.TeamName,
	}
}

func (r GetTeamByNameApiRequest) Execute() (*TeamResponse, *http.Response, error) {
	return r.ApiService.GetTeamByNameExecute(r)
}

/*
GetTeamByName Return One Team by Name

Returns one team that you identified using its human-readable name. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamName Name of the team whose information you want to return.
	@return GetTeamByNameApiRequest
*/
func (a *TeamsApiService) GetTeamByName(ctx context.Context, orgId string, teamName string) GetTeamByNameApiRequest {
	return GetTeamByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamName:   teamName,
	}
}

// GetTeamByNameExecute executes the request
//
//	@return TeamResponse
func (a *TeamsApiService) GetTeamByNameExecute(r GetTeamByNameApiRequest) (*TeamResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *TeamResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.GetTeamByName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/byName/{teamName}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamName == "" {
		return localVarReturnValue, nil, reportError("teamName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamName"+"}", url.PathEscape(r.teamName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ListGroupTeamsApiRequest struct {
	ctx          context.Context
	ApiService   TeamsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupTeamsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *TeamsApiService) ListGroupTeamsWithParams(ctx context.Context, args *ListGroupTeamsApiParams) ListGroupTeamsApiRequest {
	return ListGroupTeamsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupTeamsApiRequest) IncludeCount(includeCount bool) ListGroupTeamsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupTeamsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupTeamsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupTeamsApiRequest) PageNum(pageNum int) ListGroupTeamsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupTeamsApiRequest) Execute() (*PaginatedTeamRole, *http.Response, error) {
	return r.ApiService.ListGroupTeamsExecute(r)
}

/*
ListGroupTeams Return All Teams in One Project

Returns all teams to which the authenticated user has access in the project specified using its unique 24-hexadecimal digit identifier. All members of the team share the same project access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupTeamsApiRequest
*/
func (a *TeamsApiService) ListGroupTeams(ctx context.Context, groupId string) ListGroupTeamsApiRequest {
	return ListGroupTeamsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupTeamsExecute executes the request
//
//	@return PaginatedTeamRole
func (a *TeamsApiService) ListGroupTeamsExecute(r ListGroupTeamsApiRequest) (*PaginatedTeamRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedTeamRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.ListGroupTeams")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/teams"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeCount != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	} else {
		var defaultValue bool = true
		r.includeCount = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	}
	if r.itemsPerPage != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	} else {
		var defaultValue int = 100
		r.itemsPerPage = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	}
	if r.pageNum != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	} else {
		var defaultValue int = 1
		r.pageNum = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ListOrgTeamsApiRequest struct {
	ctx          context.Context
	ApiService   TeamsApi
	orgId        string
	itemsPerPage *int
	includeCount *bool
	pageNum      *int
}

type ListOrgTeamsApiParams struct {
	OrgId        string
	ItemsPerPage *int
	IncludeCount *bool
	PageNum      *int
}

func (a *TeamsApiService) ListOrgTeamsWithParams(ctx context.Context, args *ListOrgTeamsApiParams) ListOrgTeamsApiRequest {
	return ListOrgTeamsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		itemsPerPage: args.ItemsPerPage,
		includeCount: args.IncludeCount,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListOrgTeamsApiRequest) ItemsPerPage(itemsPerPage int) ListOrgTeamsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgTeamsApiRequest) IncludeCount(includeCount bool) ListOrgTeamsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgTeamsApiRequest) PageNum(pageNum int) ListOrgTeamsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOrgTeamsApiRequest) Execute() (*PaginatedTeam, *http.Response, error) {
	return r.ApiService.ListOrgTeamsExecute(r)
}

/*
ListOrgTeams Return All Teams in One Organization

Returns all teams that belong to the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users. MongoDB Cloud only returns teams for which you have access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgTeamsApiRequest
*/
func (a *TeamsApiService) ListOrgTeams(ctx context.Context, orgId string) ListOrgTeamsApiRequest {
	return ListOrgTeamsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgTeamsExecute executes the request
//
//	@return PaginatedTeam
func (a *TeamsApiService) ListOrgTeamsExecute(r ListOrgTeamsApiRequest) (*PaginatedTeam, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedTeam
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.ListOrgTeams")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.itemsPerPage != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	} else {
		var defaultValue int = 100
		r.itemsPerPage = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	}
	if r.includeCount != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	} else {
		var defaultValue bool = true
		r.includeCount = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	}
	if r.pageNum != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	} else {
		var defaultValue int = 1
		r.pageNum = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type RemoveGroupTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	groupId    string
	teamId     string
}

type RemoveGroupTeamApiParams struct {
	GroupId string
	TeamId  string
}

func (a *TeamsApiService) RemoveGroupTeamWithParams(ctx context.Context, args *RemoveGroupTeamApiParams) RemoveGroupTeamApiRequest {
	return RemoveGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		teamId:     args.TeamId,
	}
}

func (r RemoveGroupTeamApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RemoveGroupTeamExecute(r)
}

/*
RemoveGroupTeam Remove One Team from One Project

Removes one team specified using its unique 24-hexadecimal digit identifier from the project specified using its unique 24-hexadecimal digit identifier.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to remove from the specified project.
	@return RemoveGroupTeamApiRequest
*/
func (a *TeamsApiService) RemoveGroupTeam(ctx context.Context, groupId string, teamId string) RemoveGroupTeamApiRequest {
	return RemoveGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		teamId:     teamId,
	}
}

// RemoveGroupTeamExecute executes the request
func (a *TeamsApiService) RemoveGroupTeamExecute(r RemoveGroupTeamApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.RemoveGroupTeam")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/teams/{teamId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.teamId == "" {
		return nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type RemoveUserFromTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	teamId     string
	userId     string
}

type RemoveUserFromTeamApiParams struct {
	OrgId  string
	TeamId string
	UserId string
}

func (a *TeamsApiService) RemoveUserFromTeamWithParams(ctx context.Context, args *RemoveUserFromTeamApiParams) RemoveUserFromTeamApiRequest {
	return RemoveUserFromTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		teamId:     args.TeamId,
		userId:     args.UserId,
	}
}

func (r RemoveUserFromTeamApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RemoveUserFromTeamExecute(r)
}

/*
RemoveUserFromTeam Remove One MongoDB Cloud User from One Team

Removes one MongoDB Cloud user from the specified team. This team belongs to one organization. Teams enable you to grant project access roles to MongoDB Cloud users.

**Note**: This endpoint is deprecated. Use [Remove One MongoDB Cloud User from One Team](#tag/MongoDB-Cloud-Users/operation/removeUserFromTeam) to remove an active or pending user from a team.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team from which you want to remove one database application user.
	@param userId Unique 24-hexadecimal digit string that identifies MongoDB Cloud user that you want to remove from the specified team.
	@return RemoveUserFromTeamApiRequest

Deprecated
*/
func (a *TeamsApiService) RemoveUserFromTeam(ctx context.Context, orgId string, teamId string, userId string) RemoveUserFromTeamApiRequest {
	return RemoveUserFromTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamId:     teamId,
		userId:     userId,
	}
}

// RemoveUserFromTeamExecute executes the request
// Deprecated
func (a *TeamsApiService) RemoveUserFromTeamExecute(r RemoveUserFromTeamApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.RemoveUserFromTeam")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}/users/{userId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamId == "" {
		return nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)
	if r.userId == "" {
		return nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type RenameOrgTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	orgId      string
	teamId     string
	teamUpdate *TeamUpdate
}

type RenameOrgTeamApiParams struct {
	OrgId      string
	TeamId     string
	TeamUpdate *TeamUpdate
}

func (a *TeamsApiService) RenameOrgTeamWithParams(ctx context.Context, args *RenameOrgTeamApiParams) RenameOrgTeamApiRequest {
	return RenameOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		teamId:     args.TeamId,
		teamUpdate: args.TeamUpdate,
	}
}

func (r RenameOrgTeamApiRequest) Execute() (*TeamResponse, *http.Response, error) {
	return r.ApiService.RenameOrgTeamExecute(r)
}

/*
RenameOrgTeam Rename One Team

Renames one team in the specified organization. Teams enable you to grant project access roles to MongoDB Cloud users.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team that you want to rename.
	@return RenameOrgTeamApiRequest
*/
func (a *TeamsApiService) RenameOrgTeam(ctx context.Context, orgId string, teamId string, teamUpdate *TeamUpdate) RenameOrgTeamApiRequest {
	return RenameOrgTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamId:     teamId,
		teamUpdate: teamUpdate,
	}
}

// RenameOrgTeamExecute executes the request
//
//	@return TeamResponse
func (a *TeamsApiService) RenameOrgTeamExecute(r RenameOrgTeamApiRequest) (*TeamResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *TeamResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.RenameOrgTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.teamId == "" {
		return localVarReturnValue, nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.teamUpdate == nil {
		return localVarReturnValue, nil, reportError("teamUpdate is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.teamUpdate
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateGroupTeamApiRequest struct {
	ctx        context.Context
	ApiService TeamsApi
	groupId    string
	teamId     string
	teamRole   *TeamRole
}

type UpdateGroupTeamApiParams struct {
	GroupId  string
	TeamId   string
	TeamRole *TeamRole
}

func (a *TeamsApiService) UpdateGroupTeamWithParams(ctx context.Context, args *UpdateGroupTeamApiParams) UpdateGroupTeamApiRequest {
	return UpdateGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		teamId:     args.TeamId,
		teamRole:   args.TeamRole,
	}
}

func (r UpdateGroupTeamApiRequest) Execute() (*PaginatedTeamRole, *http.Response, error) {
	return r.ApiService.UpdateGroupTeamExecute(r)
}

/*
UpdateGroupTeam Update Team Roles in One Project

Updates the project roles assigned to the specified team. You can grant team roles for specific projects and grant project access roles to users in the team. All members of the team share the same project access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param teamId Unique 24-hexadecimal digit string that identifies the team for which you want to update roles.
	@return UpdateGroupTeamApiRequest
*/
func (a *TeamsApiService) UpdateGroupTeam(ctx context.Context, groupId string, teamId string, teamRole *TeamRole) UpdateGroupTeamApiRequest {
	return UpdateGroupTeamApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		teamId:     teamId,
		teamRole:   teamRole,
	}
}

// UpdateGroupTeamExecute executes the request
//
//	@return PaginatedTeamRole
func (a *TeamsApiService) UpdateGroupTeamExecute(r UpdateGroupTeamApiRequest) (*PaginatedTeamRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedTeamRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TeamsApiService.UpdateGroupTeam")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/teams/{teamId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.teamId == "" {
		return localVarReturnValue, nil, reportError("teamId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"teamId"+"}", url.PathEscape(r.teamId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.teamRole == nil {
		return localVarReturnValue, nil, reportError("teamRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.teamRole
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
