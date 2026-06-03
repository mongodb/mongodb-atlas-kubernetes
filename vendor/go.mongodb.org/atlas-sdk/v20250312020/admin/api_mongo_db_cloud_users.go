// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type MongoDBCloudUsersApi interface {

	/*
			AddGroupUserRole Add One Project Role to One MongoDB Cloud User

			Adds one project-level role to the MongoDB Cloud user. You can add a role to an active user or a user that has been invited to join the project.

		**Note**: This resource cannot be used to add a role to users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
			@param addOrRemoveGroupRole Project-level role to assign to the MongoDB Cloud user.
			@return AddGroupUserRoleApiRequest
	*/
	AddGroupUserRole(ctx context.Context, groupId string, userId string, addOrRemoveGroupRole *AddOrRemoveGroupRole) AddGroupUserRoleApiRequest
	/*
		AddGroupUserRole Add One Project Role to One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddGroupUserRoleApiParams - Parameters for the request
		@return AddGroupUserRoleApiRequest
	*/
	AddGroupUserRoleWithParams(ctx context.Context, args *AddGroupUserRoleApiParams) AddGroupUserRoleApiRequest

	// Method available only for mocking purposes
	AddGroupUserRoleExecute(r AddGroupUserRoleApiRequest) (*GroupUserResponse, *http.Response, error)

	/*
			AddGroupUsers Add One MongoDB Cloud User to One Project

			Adds one MongoDB Cloud user to one project.
		- If the user has a pending invitation to join the project's organization, MongoDB Cloud modifies it and grants project access.
		- If the user doesn't have an invitation to join the organization, MongoDB Cloud sends a new invitation that grants the user organization and project access.
		- If the user is already active in the project's organization, MongoDB Cloud grants access to the project.


			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param groupUserRequest The active or pending MongoDB Cloud user that you want to add to the specified project.
			@return AddGroupUsersApiRequest
	*/
	AddGroupUsers(ctx context.Context, groupId string, groupUserRequest *GroupUserRequest) AddGroupUsersApiRequest
	/*
		AddGroupUsers Add One MongoDB Cloud User to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddGroupUsersApiParams - Parameters for the request
		@return AddGroupUsersApiRequest
	*/
	AddGroupUsersWithParams(ctx context.Context, args *AddGroupUsersApiParams) AddGroupUsersApiRequest

	// Method available only for mocking purposes
	AddGroupUsersExecute(r AddGroupUsersApiRequest) (*GroupUserResponse, *http.Response, error)

	/*
			AddOrgRole Add One Organization Role to One MongoDB Cloud User

			Adds one organization-level role to the MongoDB Cloud user. You can add a role to an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: This operation is atomic.

		**Note**: This resource cannot be used to add a role to users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
			@param addOrRemoveOrgRole Organization-level role to assign to the MongoDB Cloud user.
			@return AddOrgRoleApiRequest
	*/
	AddOrgRole(ctx context.Context, orgId string, userId string, addOrRemoveOrgRole *AddOrRemoveOrgRole) AddOrgRoleApiRequest
	/*
		AddOrgRole Add One Organization Role to One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddOrgRoleApiParams - Parameters for the request
		@return AddOrgRoleApiRequest
	*/
	AddOrgRoleWithParams(ctx context.Context, args *AddOrgRoleApiParams) AddOrgRoleApiRequest

	// Method available only for mocking purposes
	AddOrgRoleExecute(r AddOrgRoleApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
			AddOrgTeamUser Add One MongoDB Cloud User to One Team

			Adds one MongoDB Cloud user to one team. You can add an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: This resource cannot be used to add a user invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param teamId Unique 24-hexadecimal digit string that identifies the team to add the MongoDB Cloud user to.
			@param addOrRemoveUserFromTeam The active or pending MongoDB Cloud user that you want to add to the specified team.
			@return AddOrgTeamUserApiRequest
	*/
	AddOrgTeamUser(ctx context.Context, orgId string, teamId string, addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam) AddOrgTeamUserApiRequest
	/*
		AddOrgTeamUser Add One MongoDB Cloud User to One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddOrgTeamUserApiParams - Parameters for the request
		@return AddOrgTeamUserApiRequest
	*/
	AddOrgTeamUserWithParams(ctx context.Context, args *AddOrgTeamUserApiParams) AddOrgTeamUserApiRequest

	// Method available only for mocking purposes
	AddOrgTeamUserExecute(r AddOrgTeamUserApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
			CreateOrgUser Add One MongoDB Cloud User to One Organization

			Invites one new or existing MongoDB Cloud user to join the organization. The invitation to join the organization will be sent to the username provided and must be accepted within 30 days.

		**Note**: If the user does not have an existing MongoDB Cloud account, they will be prompted to finish setting up an account upon accepting the invitation. If the user already has an account, they will still receive an invitation to access the organization.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param orgUserRequest Represents the MongoDB Cloud user to be created within the organization.
			@return CreateOrgUserApiRequest
	*/
	CreateOrgUser(ctx context.Context, orgId string, orgUserRequest *OrgUserRequest) CreateOrgUserApiRequest
	/*
		CreateOrgUser Add One MongoDB Cloud User to One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgUserApiParams - Parameters for the request
		@return CreateOrgUserApiRequest
	*/
	CreateOrgUserWithParams(ctx context.Context, args *CreateOrgUserApiParams) CreateOrgUserApiRequest

	// Method available only for mocking purposes
	CreateOrgUserExecute(r CreateOrgUserApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
			CreateUser Create One MongoDB Cloud User

			Creates one MongoDB Cloud user account. A MongoDB Cloud user account grants access to only the MongoDB Cloud application. To grant database access, create a database user. MongoDB Cloud sends an email to the users you specify, inviting them to join the project. Invited users don't have access to the project until they accept the invitation. Invitations expire after 30 days.

		 MongoDB Cloud limits MongoDB Cloud user membership to a maximum of 250 MongoDB Cloud users per team. MongoDB Cloud limits MongoDB Cloud user membership to 500 MongoDB Cloud users per project and 500 MongoDB Cloud users per organization, which includes the combined membership of all projects in the organization. MongoDB Cloud raises an error if an operation exceeds these limits. For example, if you have an organization with five projects, and each project has 100 MongoDB Cloud users, and each MongoDB Cloud user belongs to only one project, you can't add any MongoDB Cloud users to this organization without first removing existing MongoDB Cloud users from the organization.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param cloudAppUser MongoDB Cloud user account to create.
			@return CreateUserApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	CreateUser(ctx context.Context, cloudAppUser *CloudAppUser) CreateUserApiRequest
	/*
		CreateUser Create One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateUserApiParams - Parameters for the request
		@return CreateUserApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	CreateUserWithParams(ctx context.Context, args *CreateUserApiParams) CreateUserApiRequest

	// Method available only for mocking purposes
	CreateUserExecute(r CreateUserApiRequest) (*CloudAppUser, *http.Response, error)

	/*
			GetGroupUser Return One MongoDB Cloud User in One Project

			Returns information about the specified MongoDB Cloud user within the context of the specified project.

		**Note**: You can only use this resource to fetch information about MongoDB Cloud human users. To return information about an API Key, use the [Return One Organization API Key](#tag/Programmatic-API-Keys/operation/getApiKey) endpoint.

		**Note**: This resource does not return information about pending users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
			@return GetGroupUserApiRequest
	*/
	GetGroupUser(ctx context.Context, groupId string, userId string) GetGroupUserApiRequest
	/*
		GetGroupUser Return One MongoDB Cloud User in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupUserApiParams - Parameters for the request
		@return GetGroupUserApiRequest
	*/
	GetGroupUserWithParams(ctx context.Context, args *GetGroupUserApiParams) GetGroupUserApiRequest

	// Method available only for mocking purposes
	GetGroupUserExecute(r GetGroupUserApiRequest) (*GroupUserResponse, *http.Response, error)

	/*
			GetOrgUser Return One MongoDB Cloud User in One Organization

			Returns information about the specified MongoDB Cloud user within the context of the specified organization.

		**Note**: This resource can only be used to fetch information about MongoDB Cloud human users. To return information about an API Key, use the [Return One Organization API Key](#tag/Programmatic-API-Keys/operation/getApiKey) endpoint.

		**Note**: This resource does not return information about pending users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
			@return GetOrgUserApiRequest
	*/
	GetOrgUser(ctx context.Context, orgId string, userId string) GetOrgUserApiRequest
	/*
		GetOrgUser Return One MongoDB Cloud User in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgUserApiParams - Parameters for the request
		@return GetOrgUserApiRequest
	*/
	GetOrgUserWithParams(ctx context.Context, args *GetOrgUserApiParams) GetOrgUserApiRequest

	// Method available only for mocking purposes
	GetOrgUserExecute(r GetOrgUserApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
		GetUser Return One MongoDB Cloud User by ID

		Returns the details for one MongoDB Cloud user account with the specified unique identifier for the user. You can't use this endpoint to return information on an API Key. To return information about an API Key, use the Return One Organization API Key endpoint. You can always retrieve your own user account. If you are the owner of a MongoDB Cloud organization or project, you can also retrieve the user profile for any user with membership in that organization or project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param userId Unique 24-hexadecimal digit string that identifies this user.
		@return GetUserApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	GetUser(ctx context.Context, userId string) GetUserApiRequest
	/*
		GetUser Return One MongoDB Cloud User by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetUserApiParams - Parameters for the request
		@return GetUserApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	GetUserWithParams(ctx context.Context, args *GetUserApiParams) GetUserApiRequest

	// Method available only for mocking purposes
	GetUserExecute(r GetUserApiRequest) (*CloudAppUser, *http.Response, error)

	/*
		GetUserByName Return One MongoDB Cloud User by Username

		Returns the details for one MongoDB Cloud user account with the specified username. You can't use this endpoint to return information about an API Key. To return information about an API Key, use the Return One Organization API Key endpoint.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param userName Email address that belongs to the MongoDB Cloud user account. You cannot modify this address after creating the user.
		@return GetUserByNameApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	GetUserByName(ctx context.Context, userName string) GetUserByNameApiRequest
	/*
		GetUserByName Return One MongoDB Cloud User by Username


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetUserByNameApiParams - Parameters for the request
		@return GetUserByNameApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for MongoDBCloudUsersApi
	*/
	GetUserByNameWithParams(ctx context.Context, args *GetUserByNameApiParams) GetUserByNameApiRequest

	// Method available only for mocking purposes
	GetUserByNameExecute(r GetUserByNameApiRequest) (*CloudAppUser, *http.Response, error)

	/*
			ListGroupUsers Return All MongoDB Cloud Users in One Project

			Returns details about the pending and active MongoDB Cloud users associated with the specified project.

		**Note**: This resource cannot be used to view details about users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

		**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListGroupUsersApiRequest
	*/
	ListGroupUsers(ctx context.Context, groupId string) ListGroupUsersApiRequest
	/*
		ListGroupUsers Return All MongoDB Cloud Users in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupUsersApiParams - Parameters for the request
		@return ListGroupUsersApiRequest
	*/
	ListGroupUsersWithParams(ctx context.Context, args *ListGroupUsersApiParams) ListGroupUsersApiRequest

	// Method available only for mocking purposes
	ListGroupUsersExecute(r ListGroupUsersApiRequest) (*PaginatedGroupUser, *http.Response, error)

	/*
			ListOrgUsers Return All MongoDB Cloud Users in One Organization

			Returns details about the pending and active MongoDB Cloud users associated with the specified organization.

		**Note**: This resource cannot be used to view details about users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

		**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return ListOrgUsersApiRequest
	*/
	ListOrgUsers(ctx context.Context, orgId string) ListOrgUsersApiRequest
	/*
		ListOrgUsers Return All MongoDB Cloud Users in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgUsersApiParams - Parameters for the request
		@return ListOrgUsersApiRequest
	*/
	ListOrgUsersWithParams(ctx context.Context, args *ListOrgUsersApiParams) ListOrgUsersApiRequest

	// Method available only for mocking purposes
	ListOrgUsersExecute(r ListOrgUsersApiRequest) (*PaginatedOrgUser, *http.Response, error)

	/*
			ListTeamUsers Return All MongoDB Cloud Users Assigned to One Team

			Returns details about the pending and active MongoDB Cloud users associated with the specified team in the organization. Teams enable you to grant project access roles to MongoDB Cloud users.

		**Note**: This resource cannot be used to view details about users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

		**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param teamId Unique 24-hexadecimal digit string that identifies the team whose application users you want to return.
			@return ListTeamUsersApiRequest
	*/
	ListTeamUsers(ctx context.Context, orgId string, teamId string) ListTeamUsersApiRequest
	/*
		ListTeamUsers Return All MongoDB Cloud Users Assigned to One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListTeamUsersApiParams - Parameters for the request
		@return ListTeamUsersApiRequest
	*/
	ListTeamUsersWithParams(ctx context.Context, args *ListTeamUsersApiParams) ListTeamUsersApiRequest

	// Method available only for mocking purposes
	ListTeamUsersExecute(r ListTeamUsersApiRequest) (*PaginatedOrgUser, *http.Response, error)

	/*
			RemoveGroupUser Remove One MongoDB Cloud User from One Project

			Removes one MongoDB Cloud user from the specified project. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: This resource cannot be used to remove pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

		**Note**: To remove pending or active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users can be removed. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the [Return All MongoDB Cloud Users in One Project](#tag/MongoDB-Cloud-Users/operation/listProjectUsers) resource and filter by `username`.
			@return RemoveGroupUserApiRequest
	*/
	RemoveGroupUser(ctx context.Context, groupId string, userId string) RemoveGroupUserApiRequest
	/*
		RemoveGroupUser Remove One MongoDB Cloud User from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveGroupUserApiParams - Parameters for the request
		@return RemoveGroupUserApiRequest
	*/
	RemoveGroupUserWithParams(ctx context.Context, args *RemoveGroupUserApiParams) RemoveGroupUserApiRequest

	// Method available only for mocking purposes
	RemoveGroupUserExecute(r RemoveGroupUserApiRequest) (*http.Response, error)

	/*
			RemoveGroupUserRole Remove One Project Role from One MongoDB Cloud User

			Removes one project-level role from the MongoDB Cloud user. You can remove a role from an active user or a user that has been invited to join the project. To replace a user's only role, add the new role before removing the old role. A user must have at least one role at all times.

		**Note**: This resource cannot be used to remove a role from users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
			@param addOrRemoveGroupRole Project-level role to remove from the MongoDB Cloud user.
			@return RemoveGroupUserRoleApiRequest
	*/
	RemoveGroupUserRole(ctx context.Context, groupId string, userId string, addOrRemoveGroupRole *AddOrRemoveGroupRole) RemoveGroupUserRoleApiRequest
	/*
		RemoveGroupUserRole Remove One Project Role from One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveGroupUserRoleApiParams - Parameters for the request
		@return RemoveGroupUserRoleApiRequest
	*/
	RemoveGroupUserRoleWithParams(ctx context.Context, args *RemoveGroupUserRoleApiParams) RemoveGroupUserRoleApiRequest

	// Method available only for mocking purposes
	RemoveGroupUserRoleExecute(r RemoveGroupUserRoleApiRequest) (*GroupUserResponse, *http.Response, error)

	/*
			RemoveOrgRole Remove One Organization Role from One MongoDB Cloud User

			Removes one organization-level role from the MongoDB Cloud user. You can remove a role from an active user or a user that has not yet accepted the invitation to join the organization. To replace a user's only role, add the new role before removing the old role. A user must have at least one role at all times.

		**Note**: This operation is atomic.

		**Note**: This resource cannot be used to remove a role from users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
			@param addOrRemoveOrgRole Organization-level role to remove from the MongoDB Cloud user.
			@return RemoveOrgRoleApiRequest
	*/
	RemoveOrgRole(ctx context.Context, orgId string, userId string, addOrRemoveOrgRole *AddOrRemoveOrgRole) RemoveOrgRoleApiRequest
	/*
		RemoveOrgRole Remove One Organization Role from One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveOrgRoleApiParams - Parameters for the request
		@return RemoveOrgRoleApiRequest
	*/
	RemoveOrgRoleWithParams(ctx context.Context, args *RemoveOrgRoleApiParams) RemoveOrgRoleApiRequest

	// Method available only for mocking purposes
	RemoveOrgRoleExecute(r RemoveOrgRoleApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
			RemoveOrgTeamUser Remove One MongoDB Cloud User from One Team

			Removes one MongoDB Cloud user from one team. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: This resource cannot be used to remove a user invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param teamId Unique 24-hexadecimal digit string that identifies the team to remove the MongoDB user from.
			@param addOrRemoveUserFromTeam The id of the active or pending MongoDB Cloud user that you want to remove from the specified team.
			@return RemoveOrgTeamUserApiRequest
	*/
	RemoveOrgTeamUser(ctx context.Context, orgId string, teamId string, addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam) RemoveOrgTeamUserApiRequest
	/*
		RemoveOrgTeamUser Remove One MongoDB Cloud User from One Team


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveOrgTeamUserApiParams - Parameters for the request
		@return RemoveOrgTeamUserApiRequest
	*/
	RemoveOrgTeamUserWithParams(ctx context.Context, args *RemoveOrgTeamUserApiParams) RemoveOrgTeamUserApiRequest

	// Method available only for mocking purposes
	RemoveOrgTeamUserExecute(r RemoveOrgTeamUserApiRequest) (*OrgUserResponse, *http.Response, error)

	/*
			RemoveOrgUser Remove One MongoDB Cloud User from One Organization

			Removes one MongoDB Cloud user in the specified organization. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: This resource cannot be used to remove pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

		**Note**: To remove pending or active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users can be removed. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the [Return All MongoDB Cloud Users in One Organization](#tag/MongoDB-Cloud-Users/operation/listOrganizationUsers) resource and filter by `username`.
			@return RemoveOrgUserApiRequest
	*/
	RemoveOrgUser(ctx context.Context, orgId string, userId string) RemoveOrgUserApiRequest
	/*
		RemoveOrgUser Remove One MongoDB Cloud User from One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveOrgUserApiParams - Parameters for the request
		@return RemoveOrgUserApiRequest
	*/
	RemoveOrgUserWithParams(ctx context.Context, args *RemoveOrgUserApiParams) RemoveOrgUserApiRequest

	// Method available only for mocking purposes
	RemoveOrgUserExecute(r RemoveOrgUserApiRequest) (*http.Response, error)

	/*
			UpdateOrgUser Update One MongoDB Cloud User in One Organization

			Updates one MongoDB Cloud user in the specified organization. You can update an active user or a user that has not yet accepted the invitation to join the organization.

		**Note**: Only include the fields you wish to update in the request body. Supplying a field with an empty value will reset that field on the user.

		**Note**: This resource cannot be used to update pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
			@param orgUserUpdateRequest Represents the roles and teams to assign the MongoDB Cloud user.
			@return UpdateOrgUserApiRequest
	*/
	UpdateOrgUser(ctx context.Context, orgId string, userId string, orgUserUpdateRequest *OrgUserUpdateRequest) UpdateOrgUserApiRequest
	/*
		UpdateOrgUser Update One MongoDB Cloud User in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgUserApiParams - Parameters for the request
		@return UpdateOrgUserApiRequest
	*/
	UpdateOrgUserWithParams(ctx context.Context, args *UpdateOrgUserApiParams) UpdateOrgUserApiRequest

	// Method available only for mocking purposes
	UpdateOrgUserExecute(r UpdateOrgUserApiRequest) (*OrgUserResponse, *http.Response, error)
}

// MongoDBCloudUsersApiService MongoDBCloudUsersApi service
type MongoDBCloudUsersApiService service

type AddGroupUserRoleApiRequest struct {
	ctx                  context.Context
	ApiService           MongoDBCloudUsersApi
	groupId              string
	userId               string
	addOrRemoveGroupRole *AddOrRemoveGroupRole
}

type AddGroupUserRoleApiParams struct {
	GroupId              string
	UserId               string
	AddOrRemoveGroupRole *AddOrRemoveGroupRole
}

func (a *MongoDBCloudUsersApiService) AddGroupUserRoleWithParams(ctx context.Context, args *AddGroupUserRoleApiParams) AddGroupUserRoleApiRequest {
	return AddGroupUserRoleApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		groupId:              args.GroupId,
		userId:               args.UserId,
		addOrRemoveGroupRole: args.AddOrRemoveGroupRole,
	}
}

func (r AddGroupUserRoleApiRequest) Execute() (*GroupUserResponse, *http.Response, error) {
	return r.ApiService.AddGroupUserRoleExecute(r)
}

/*
AddGroupUserRole Add One Project Role to One MongoDB Cloud User

Adds one project-level role to the MongoDB Cloud user. You can add a role to an active user or a user that has been invited to join the project.

**Note**: This resource cannot be used to add a role to users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
	@return AddGroupUserRoleApiRequest
*/
func (a *MongoDBCloudUsersApiService) AddGroupUserRole(ctx context.Context, groupId string, userId string, addOrRemoveGroupRole *AddOrRemoveGroupRole) AddGroupUserRoleApiRequest {
	return AddGroupUserRoleApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		groupId:              groupId,
		userId:               userId,
		addOrRemoveGroupRole: addOrRemoveGroupRole,
	}
}

// AddGroupUserRoleExecute executes the request
//
//	@return GroupUserResponse
func (a *MongoDBCloudUsersApiService) AddGroupUserRoleExecute(r AddGroupUserRoleApiRequest) (*GroupUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.AddGroupUserRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users/{userId}:addRole"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.addOrRemoveGroupRole == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveGroupRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveGroupRole
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

type AddGroupUsersApiRequest struct {
	ctx              context.Context
	ApiService       MongoDBCloudUsersApi
	groupId          string
	groupUserRequest *GroupUserRequest
}

type AddGroupUsersApiParams struct {
	GroupId          string
	GroupUserRequest *GroupUserRequest
}

func (a *MongoDBCloudUsersApiService) AddGroupUsersWithParams(ctx context.Context, args *AddGroupUsersApiParams) AddGroupUsersApiRequest {
	return AddGroupUsersApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		groupUserRequest: args.GroupUserRequest,
	}
}

func (r AddGroupUsersApiRequest) Execute() (*GroupUserResponse, *http.Response, error) {
	return r.ApiService.AddGroupUsersExecute(r)
}

/*
AddGroupUsers Add One MongoDB Cloud User to One Project

Adds one MongoDB Cloud user to one project.
- If the user has a pending invitation to join the project's organization, MongoDB Cloud modifies it and grants project access.
- If the user doesn't have an invitation to join the organization, MongoDB Cloud sends a new invitation that grants the user organization and project access.
- If the user is already active in the project's organization, MongoDB Cloud grants access to the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return AddGroupUsersApiRequest
*/
func (a *MongoDBCloudUsersApiService) AddGroupUsers(ctx context.Context, groupId string, groupUserRequest *GroupUserRequest) AddGroupUsersApiRequest {
	return AddGroupUsersApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		groupUserRequest: groupUserRequest,
	}
}

// AddGroupUsersExecute executes the request
//
//	@return GroupUserResponse
func (a *MongoDBCloudUsersApiService) AddGroupUsersExecute(r AddGroupUsersApiRequest) (*GroupUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.AddGroupUsers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupUserRequest == nil {
		return localVarReturnValue, nil, reportError("groupUserRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupUserRequest
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

type AddOrgRoleApiRequest struct {
	ctx                context.Context
	ApiService         MongoDBCloudUsersApi
	orgId              string
	userId             string
	addOrRemoveOrgRole *AddOrRemoveOrgRole
}

type AddOrgRoleApiParams struct {
	OrgId              string
	UserId             string
	AddOrRemoveOrgRole *AddOrRemoveOrgRole
}

func (a *MongoDBCloudUsersApiService) AddOrgRoleWithParams(ctx context.Context, args *AddOrgRoleApiParams) AddOrgRoleApiRequest {
	return AddOrgRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		orgId:              args.OrgId,
		userId:             args.UserId,
		addOrRemoveOrgRole: args.AddOrRemoveOrgRole,
	}
}

func (r AddOrgRoleApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.AddOrgRoleExecute(r)
}

/*
AddOrgRole Add One Organization Role to One MongoDB Cloud User

Adds one organization-level role to the MongoDB Cloud user. You can add a role to an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: This operation is atomic.

**Note**: This resource cannot be used to add a role to users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
	@return AddOrgRoleApiRequest
*/
func (a *MongoDBCloudUsersApiService) AddOrgRole(ctx context.Context, orgId string, userId string, addOrRemoveOrgRole *AddOrRemoveOrgRole) AddOrgRoleApiRequest {
	return AddOrgRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		orgId:              orgId,
		userId:             userId,
		addOrRemoveOrgRole: addOrRemoveOrgRole,
	}
}

// AddOrgRoleExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) AddOrgRoleExecute(r AddOrgRoleApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.AddOrgRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}:addRole"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.addOrRemoveOrgRole == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveOrgRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveOrgRole
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

type AddOrgTeamUserApiRequest struct {
	ctx                     context.Context
	ApiService              MongoDBCloudUsersApi
	orgId                   string
	teamId                  string
	addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam
}

type AddOrgTeamUserApiParams struct {
	OrgId                   string
	TeamId                  string
	AddOrRemoveUserFromTeam *AddOrRemoveUserFromTeam
}

func (a *MongoDBCloudUsersApiService) AddOrgTeamUserWithParams(ctx context.Context, args *AddOrgTeamUserApiParams) AddOrgTeamUserApiRequest {
	return AddOrgTeamUserApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		orgId:                   args.OrgId,
		teamId:                  args.TeamId,
		addOrRemoveUserFromTeam: args.AddOrRemoveUserFromTeam,
	}
}

func (r AddOrgTeamUserApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.AddOrgTeamUserExecute(r)
}

/*
AddOrgTeamUser Add One MongoDB Cloud User to One Team

Adds one MongoDB Cloud user to one team. You can add an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: This resource cannot be used to add a user invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team to add the MongoDB Cloud user to.
	@return AddOrgTeamUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) AddOrgTeamUser(ctx context.Context, orgId string, teamId string, addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam) AddOrgTeamUserApiRequest {
	return AddOrgTeamUserApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		orgId:                   orgId,
		teamId:                  teamId,
		addOrRemoveUserFromTeam: addOrRemoveUserFromTeam,
	}
}

// AddOrgTeamUserExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) AddOrgTeamUserExecute(r AddOrgTeamUserApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.AddOrgTeamUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}:addUser"
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
	if r.addOrRemoveUserFromTeam == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveUserFromTeam is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveUserFromTeam
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

type CreateOrgUserApiRequest struct {
	ctx            context.Context
	ApiService     MongoDBCloudUsersApi
	orgId          string
	orgUserRequest *OrgUserRequest
}

type CreateOrgUserApiParams struct {
	OrgId          string
	OrgUserRequest *OrgUserRequest
}

func (a *MongoDBCloudUsersApiService) CreateOrgUserWithParams(ctx context.Context, args *CreateOrgUserApiParams) CreateOrgUserApiRequest {
	return CreateOrgUserApiRequest{
		ApiService:     a,
		ctx:            ctx,
		orgId:          args.OrgId,
		orgUserRequest: args.OrgUserRequest,
	}
}

func (r CreateOrgUserApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.CreateOrgUserExecute(r)
}

/*
CreateOrgUser Add One MongoDB Cloud User to One Organization

Invites one new or existing MongoDB Cloud user to join the organization. The invitation to join the organization will be sent to the username provided and must be accepted within 30 days.

**Note**: If the user does not have an existing MongoDB Cloud account, they will be prompted to finish setting up an account upon accepting the invitation. If the user already has an account, they will still receive an invitation to access the organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) CreateOrgUser(ctx context.Context, orgId string, orgUserRequest *OrgUserRequest) CreateOrgUserApiRequest {
	return CreateOrgUserApiRequest{
		ApiService:     a,
		ctx:            ctx,
		orgId:          orgId,
		orgUserRequest: orgUserRequest,
	}
}

// CreateOrgUserExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) CreateOrgUserExecute(r CreateOrgUserApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.CreateOrgUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.orgUserRequest == nil {
		return localVarReturnValue, nil, reportError("orgUserRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.orgUserRequest
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

type CreateUserApiRequest struct {
	ctx          context.Context
	ApiService   MongoDBCloudUsersApi
	cloudAppUser *CloudAppUser
}

type CreateUserApiParams struct {
	CloudAppUser *CloudAppUser
}

func (a *MongoDBCloudUsersApiService) CreateUserWithParams(ctx context.Context, args *CreateUserApiParams) CreateUserApiRequest {
	return CreateUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		cloudAppUser: args.CloudAppUser,
	}
}

func (r CreateUserApiRequest) Execute() (*CloudAppUser, *http.Response, error) {
	return r.ApiService.CreateUserExecute(r)
}

/*
CreateUser Create One MongoDB Cloud User

Creates one MongoDB Cloud user account. A MongoDB Cloud user account grants access to only the MongoDB Cloud application. To grant database access, create a database user. MongoDB Cloud sends an email to the users you specify, inviting them to join the project. Invited users don't have access to the project until they accept the invitation. Invitations expire after 30 days.

	MongoDB Cloud limits MongoDB Cloud user membership to a maximum of 250 MongoDB Cloud users per team. MongoDB Cloud limits MongoDB Cloud user membership to 500 MongoDB Cloud users per project and 500 MongoDB Cloud users per organization, which includes the combined membership of all projects in the organization. MongoDB Cloud raises an error if an operation exceeds these limits. For example, if you have an organization with five projects, and each project has 100 MongoDB Cloud users, and each MongoDB Cloud user belongs to only one project, you can't add any MongoDB Cloud users to this organization without first removing existing MongoDB Cloud users from the organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return CreateUserApiRequest

Deprecated
*/
func (a *MongoDBCloudUsersApiService) CreateUser(ctx context.Context, cloudAppUser *CloudAppUser) CreateUserApiRequest {
	return CreateUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		cloudAppUser: cloudAppUser,
	}
}

// CreateUserExecute executes the request
//
//	@return CloudAppUser
//
// Deprecated
func (a *MongoDBCloudUsersApiService) CreateUserExecute(r CreateUserApiRequest) (*CloudAppUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudAppUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.CreateUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/users"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudAppUser == nil {
		return localVarReturnValue, nil, reportError("cloudAppUser is required and must be specified")
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
	localVarPostBody = r.cloudAppUser
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

type GetGroupUserApiRequest struct {
	ctx                   context.Context
	ApiService            MongoDBCloudUsersApi
	groupId               string
	userId                string
	orgMembershipStatuses *[]string
}

type GetGroupUserApiParams struct {
	GroupId               string
	UserId                string
	OrgMembershipStatuses *[]string
}

func (a *MongoDBCloudUsersApiService) GetGroupUserWithParams(ctx context.Context, args *GetGroupUserApiParams) GetGroupUserApiRequest {
	return GetGroupUserApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		userId:                args.UserId,
		orgMembershipStatuses: args.OrgMembershipStatuses,
	}
}

// Organization membership status to filter users by. You can supply this parameter multiple times. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Not supported in deprecated versions.
func (r GetGroupUserApiRequest) OrgMembershipStatuses(orgMembershipStatuses []string) GetGroupUserApiRequest {
	r.orgMembershipStatuses = &orgMembershipStatuses
	return r
}

func (r GetGroupUserApiRequest) Execute() (*GroupUserResponse, *http.Response, error) {
	return r.ApiService.GetGroupUserExecute(r)
}

/*
GetGroupUser Return One MongoDB Cloud User in One Project

Returns information about the specified MongoDB Cloud user within the context of the specified project.

**Note**: You can only use this resource to fetch information about MongoDB Cloud human users. To return information about an API Key, use the [Return One Organization API Key](#tag/Programmatic-API-Keys/operation/getApiKey) endpoint.

**Note**: This resource does not return information about pending users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
	@return GetGroupUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) GetGroupUser(ctx context.Context, groupId string, userId string) GetGroupUserApiRequest {
	return GetGroupUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		userId:     userId,
	}
}

// GetGroupUserExecute executes the request
//
//	@return GroupUserResponse
func (a *MongoDBCloudUsersApiService) GetGroupUserExecute(r GetGroupUserApiRequest) (*GroupUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.GetGroupUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users/{userId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.orgMembershipStatuses != nil {
		t := *r.orgMembershipStatuses
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatuses", t, "multi")

	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type GetOrgUserApiRequest struct {
	ctx                   context.Context
	ApiService            MongoDBCloudUsersApi
	orgId                 string
	userId                string
	orgMembershipStatuses *[]string
}

type GetOrgUserApiParams struct {
	OrgId                 string
	UserId                string
	OrgMembershipStatuses *[]string
}

func (a *MongoDBCloudUsersApiService) GetOrgUserWithParams(ctx context.Context, args *GetOrgUserApiParams) GetOrgUserApiRequest {
	return GetOrgUserApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 args.OrgId,
		userId:                args.UserId,
		orgMembershipStatuses: args.OrgMembershipStatuses,
	}
}

// Organization membership status to filter users by. You can supply this parameter multiple times. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Not supported in deprecated versions.
func (r GetOrgUserApiRequest) OrgMembershipStatuses(orgMembershipStatuses []string) GetOrgUserApiRequest {
	r.orgMembershipStatuses = &orgMembershipStatuses
	return r
}

func (r GetOrgUserApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.GetOrgUserExecute(r)
}

/*
GetOrgUser Return One MongoDB Cloud User in One Organization

Returns information about the specified MongoDB Cloud user within the context of the specified organization.

**Note**: This resource can only be used to fetch information about MongoDB Cloud human users. To return information about an API Key, use the [Return One Organization API Key](#tag/Programmatic-API-Keys/operation/getApiKey) endpoint.

**Note**: This resource does not return information about pending users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
	@return GetOrgUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) GetOrgUser(ctx context.Context, orgId string, userId string) GetOrgUserApiRequest {
	return GetOrgUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		userId:     userId,
	}
}

// GetOrgUserExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) GetOrgUserExecute(r GetOrgUserApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.GetOrgUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.orgMembershipStatuses != nil {
		t := *r.orgMembershipStatuses
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatuses", t, "multi")

	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type GetUserApiRequest struct {
	ctx        context.Context
	ApiService MongoDBCloudUsersApi
	userId     string
}

type GetUserApiParams struct {
	UserId string
}

func (a *MongoDBCloudUsersApiService) GetUserWithParams(ctx context.Context, args *GetUserApiParams) GetUserApiRequest {
	return GetUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		userId:     args.UserId,
	}
}

func (r GetUserApiRequest) Execute() (*CloudAppUser, *http.Response, error) {
	return r.ApiService.GetUserExecute(r)
}

/*
GetUser Return One MongoDB Cloud User by ID

Returns the details for one MongoDB Cloud user account with the specified unique identifier for the user. You can't use this endpoint to return information on an API Key. To return information about an API Key, use the Return One Organization API Key endpoint. You can always retrieve your own user account. If you are the owner of a MongoDB Cloud organization or project, you can also retrieve the user profile for any user with membership in that organization or project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param userId Unique 24-hexadecimal digit string that identifies this user.
	@return GetUserApiRequest

Deprecated
*/
func (a *MongoDBCloudUsersApiService) GetUser(ctx context.Context, userId string) GetUserApiRequest {
	return GetUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		userId:     userId,
	}
}

// GetUserExecute executes the request
//
//	@return CloudAppUser
//
// Deprecated
func (a *MongoDBCloudUsersApiService) GetUserExecute(r GetUserApiRequest) (*CloudAppUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudAppUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.GetUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/users/{userId}"
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
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

type GetUserByNameApiRequest struct {
	ctx        context.Context
	ApiService MongoDBCloudUsersApi
	userName   string
}

type GetUserByNameApiParams struct {
	UserName string
}

func (a *MongoDBCloudUsersApiService) GetUserByNameWithParams(ctx context.Context, args *GetUserByNameApiParams) GetUserByNameApiRequest {
	return GetUserByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		userName:   args.UserName,
	}
}

func (r GetUserByNameApiRequest) Execute() (*CloudAppUser, *http.Response, error) {
	return r.ApiService.GetUserByNameExecute(r)
}

/*
GetUserByName Return One MongoDB Cloud User by Username

Returns the details for one MongoDB Cloud user account with the specified username. You can't use this endpoint to return information about an API Key. To return information about an API Key, use the Return One Organization API Key endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param userName Email address that belongs to the MongoDB Cloud user account. You cannot modify this address after creating the user.
	@return GetUserByNameApiRequest

Deprecated
*/
func (a *MongoDBCloudUsersApiService) GetUserByName(ctx context.Context, userName string) GetUserByNameApiRequest {
	return GetUserByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		userName:   userName,
	}
}

// GetUserByNameExecute executes the request
//
//	@return CloudAppUser
//
// Deprecated
func (a *MongoDBCloudUsersApiService) GetUserByNameExecute(r GetUserByNameApiRequest) (*CloudAppUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudAppUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.GetUserByName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/users/byName/{userName}"
	if r.userName == "" {
		return localVarReturnValue, nil, reportError("userName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userName"+"}", url.PathEscape(r.userName), -1)

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

type ListGroupUsersApiRequest struct {
	ctx                   context.Context
	ApiService            MongoDBCloudUsersApi
	groupId               string
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
	flattenTeams          *bool
	includeOrgUsers       *bool
	orgMembershipStatus   *string
	orgMembershipStatuses *[]string
	username              *string
}

type ListGroupUsersApiParams struct {
	GroupId               string
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
	FlattenTeams          *bool
	IncludeOrgUsers       *bool
	OrgMembershipStatus   *string
	OrgMembershipStatuses *[]string
	Username              *string
}

func (a *MongoDBCloudUsersApiService) ListGroupUsersWithParams(ctx context.Context, args *ListGroupUsersApiParams) ListGroupUsersApiRequest {
	return ListGroupUsersApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
		flattenTeams:          args.FlattenTeams,
		includeOrgUsers:       args.IncludeOrgUsers,
		orgMembershipStatus:   args.OrgMembershipStatus,
		orgMembershipStatuses: args.OrgMembershipStatuses,
		username:              args.Username,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupUsersApiRequest) IncludeCount(includeCount bool) ListGroupUsersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupUsersApiRequest) ItemsPerPage(itemsPerPage int) ListGroupUsersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupUsersApiRequest) PageNum(pageNum int) ListGroupUsersApiRequest {
	r.pageNum = &pageNum
	return r
}

// Flag that indicates whether the returned list should include users who belong to a team with a role in this project. You might not have assigned the individual users a role in this project. If &#x60;\&quot;flattenTeams\&quot; : false&#x60;, this resource returns only users with a role in the project.  If &#x60;\&quot;flattenTeams\&quot; : true&#x60;, this resource returns both users with roles in the project and users who belong to teams with roles in the project.
func (r ListGroupUsersApiRequest) FlattenTeams(flattenTeams bool) ListGroupUsersApiRequest {
	r.flattenTeams = &flattenTeams
	return r
}

// Flag that indicates whether the returned list should include users with implicit access to the project, the Organization Owner or Organization Read Only role. You might not have assigned the individual users a role in this project. If &#x60;\&quot;includeOrgUsers\&quot;: false&#x60;, this resource returns only users with a role in the project. If &#x60;\&quot;includeOrgUsers\&quot;: true&#x60;, this resource returns both users with roles in the project and users who have implicit access to the project through their organization role.
func (r ListGroupUsersApiRequest) IncludeOrgUsers(includeOrgUsers bool) ListGroupUsersApiRequest {
	r.includeOrgUsers = &includeOrgUsers
	return r
}

// Deprecated: Use &#x60;orgMembershipStatuses&#x60; instead. Organization membership status to filter users by. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Not supported in deprecated versions.
// Deprecated
func (r ListGroupUsersApiRequest) OrgMembershipStatus(orgMembershipStatus string) ListGroupUsersApiRequest {
	r.orgMembershipStatus = &orgMembershipStatus
	return r
}

// Organization membership status to filter users by. You can supply this parameter multiple times. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. Replaces the deprecated &#x60;orgMembershipStatus&#x60; parameter. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Cannot be combined with &#x60;orgMembershipStatus&#x60;. Not supported in deprecated versions.
func (r ListGroupUsersApiRequest) OrgMembershipStatuses(orgMembershipStatuses []string) ListGroupUsersApiRequest {
	r.orgMembershipStatuses = &orgMembershipStatuses
	return r
}

// Email address to filter users by. Not supported in deprecated versions.
func (r ListGroupUsersApiRequest) Username(username string) ListGroupUsersApiRequest {
	r.username = &username
	return r
}

func (r ListGroupUsersApiRequest) Execute() (*PaginatedGroupUser, *http.Response, error) {
	return r.ApiService.ListGroupUsersExecute(r)
}

/*
ListGroupUsers Return All MongoDB Cloud Users in One Project

Returns details about the pending and active MongoDB Cloud users associated with the specified project.

**Note**: This resource cannot be used to view details about users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupUsersApiRequest
*/
func (a *MongoDBCloudUsersApiService) ListGroupUsers(ctx context.Context, groupId string) ListGroupUsersApiRequest {
	return ListGroupUsersApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupUsersExecute executes the request
//
//	@return PaginatedGroupUser
func (a *MongoDBCloudUsersApiService) ListGroupUsersExecute(r ListGroupUsersApiRequest) (*PaginatedGroupUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedGroupUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.ListGroupUsers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users"
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
	if r.flattenTeams != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "flattenTeams", r.flattenTeams, "")
	} else {
		var defaultValue bool = false
		r.flattenTeams = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "flattenTeams", r.flattenTeams, "")
	}
	if r.includeOrgUsers != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeOrgUsers", r.includeOrgUsers, "")
	} else {
		var defaultValue bool = false
		r.includeOrgUsers = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeOrgUsers", r.includeOrgUsers, "")
	}
	if r.orgMembershipStatus != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatus", r.orgMembershipStatus, "")
	}
	if r.orgMembershipStatuses != nil {
		t := *r.orgMembershipStatuses
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatuses", t, "multi")

	}
	if r.username != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "username", r.username, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type ListOrgUsersApiRequest struct {
	ctx                   context.Context
	ApiService            MongoDBCloudUsersApi
	orgId                 string
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
	username              *string
	orgMembershipStatus   *string
	orgMembershipStatuses *[]string
}

type ListOrgUsersApiParams struct {
	OrgId                 string
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
	Username              *string
	OrgMembershipStatus   *string
	OrgMembershipStatuses *[]string
}

func (a *MongoDBCloudUsersApiService) ListOrgUsersWithParams(ctx context.Context, args *ListOrgUsersApiParams) ListOrgUsersApiRequest {
	return ListOrgUsersApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 args.OrgId,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
		username:              args.Username,
		orgMembershipStatus:   args.OrgMembershipStatus,
		orgMembershipStatuses: args.OrgMembershipStatuses,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgUsersApiRequest) IncludeCount(includeCount bool) ListOrgUsersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgUsersApiRequest) ItemsPerPage(itemsPerPage int) ListOrgUsersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgUsersApiRequest) PageNum(pageNum int) ListOrgUsersApiRequest {
	r.pageNum = &pageNum
	return r
}

// Email address to filter users by. Not supported in deprecated versions.
func (r ListOrgUsersApiRequest) Username(username string) ListOrgUsersApiRequest {
	r.username = &username
	return r
}

// Deprecated: Use &#x60;orgMembershipStatuses&#x60; instead. Organization membership status to filter users by. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Not supported in deprecated versions.
// Deprecated
func (r ListOrgUsersApiRequest) OrgMembershipStatus(orgMembershipStatus string) ListOrgUsersApiRequest {
	r.orgMembershipStatus = &orgMembershipStatus
	return r
}

// Organization membership status to filter users by. You can supply this parameter multiple times. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. Replaces the deprecated &#x60;orgMembershipStatus&#x60; parameter. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Cannot be combined with &#x60;orgMembershipStatus&#x60;. Not supported in deprecated versions.
func (r ListOrgUsersApiRequest) OrgMembershipStatuses(orgMembershipStatuses []string) ListOrgUsersApiRequest {
	r.orgMembershipStatuses = &orgMembershipStatuses
	return r
}

func (r ListOrgUsersApiRequest) Execute() (*PaginatedOrgUser, *http.Response, error) {
	return r.ApiService.ListOrgUsersExecute(r)
}

/*
ListOrgUsers Return All MongoDB Cloud Users in One Organization

Returns details about the pending and active MongoDB Cloud users associated with the specified organization.

**Note**: This resource cannot be used to view details about users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgUsersApiRequest
*/
func (a *MongoDBCloudUsersApiService) ListOrgUsers(ctx context.Context, orgId string) ListOrgUsersApiRequest {
	return ListOrgUsersApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgUsersExecute executes the request
//
//	@return PaginatedOrgUser
func (a *MongoDBCloudUsersApiService) ListOrgUsersExecute(r ListOrgUsersApiRequest) (*PaginatedOrgUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOrgUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.ListOrgUsers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

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
	if r.username != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "username", r.username, "")
	}
	if r.orgMembershipStatus != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatus", r.orgMembershipStatus, "")
	}
	if r.orgMembershipStatuses != nil {
		t := *r.orgMembershipStatuses
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatuses", t, "multi")

	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type ListTeamUsersApiRequest struct {
	ctx                   context.Context
	ApiService            MongoDBCloudUsersApi
	orgId                 string
	teamId                string
	itemsPerPage          *int
	pageNum               *int
	username              *string
	orgMembershipStatus   *string
	orgMembershipStatuses *[]string
	userId                *string
}

type ListTeamUsersApiParams struct {
	OrgId                 string
	TeamId                string
	ItemsPerPage          *int
	PageNum               *int
	Username              *string
	OrgMembershipStatus   *string
	OrgMembershipStatuses *[]string
	UserId                *string
}

func (a *MongoDBCloudUsersApiService) ListTeamUsersWithParams(ctx context.Context, args *ListTeamUsersApiParams) ListTeamUsersApiRequest {
	return ListTeamUsersApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 args.OrgId,
		teamId:                args.TeamId,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
		username:              args.Username,
		orgMembershipStatus:   args.OrgMembershipStatus,
		orgMembershipStatuses: args.OrgMembershipStatuses,
		userId:                args.UserId,
	}
}

// Number of items that the response returns per page.
func (r ListTeamUsersApiRequest) ItemsPerPage(itemsPerPage int) ListTeamUsersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListTeamUsersApiRequest) PageNum(pageNum int) ListTeamUsersApiRequest {
	r.pageNum = &pageNum
	return r
}

// Email address to filter users by. Not supported in deprecated versions.
func (r ListTeamUsersApiRequest) Username(username string) ListTeamUsersApiRequest {
	r.username = &username
	return r
}

// Deprecated: Use &#x60;orgMembershipStatuses&#x60; instead. Organization membership status to filter users by. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Not supported in deprecated versions.
// Deprecated
func (r ListTeamUsersApiRequest) OrgMembershipStatus(orgMembershipStatus string) ListTeamUsersApiRequest {
	r.orgMembershipStatus = &orgMembershipStatus
	return r
}

// Organization membership status to filter users by. You can supply this parameter multiple times. Allowed values: &#x60;ACTIVE&#x60;, &#x60;PENDING&#x60;, &#x60;INVITATION_EXPIRED&#x60;, &#x60;INVITATION_REJECTED&#x60;. Replaces the deprecated &#x60;orgMembershipStatus&#x60; parameter. If you exclude this parameter, this resource returns ACTIVE and PENDING users. Cannot be combined with &#x60;orgMembershipStatus&#x60;. Not supported in deprecated versions.
func (r ListTeamUsersApiRequest) OrgMembershipStatuses(orgMembershipStatuses []string) ListTeamUsersApiRequest {
	r.orgMembershipStatuses = &orgMembershipStatuses
	return r
}

// Unique 24-hexadecimal digit string to filter users by. Not supported in deprecated versions.
func (r ListTeamUsersApiRequest) UserId(userId string) ListTeamUsersApiRequest {
	r.userId = &userId
	return r
}

func (r ListTeamUsersApiRequest) Execute() (*PaginatedOrgUser, *http.Response, error) {
	return r.ApiService.ListTeamUsersExecute(r)
}

/*
ListTeamUsers Return All MongoDB Cloud Users Assigned to One Team

Returns details about the pending and active MongoDB Cloud users associated with the specified team in the organization. Teams enable you to grant project access roles to MongoDB Cloud users.

**Note**: This resource cannot be used to view details about users invited via the deprecated [Invite One MongoDB Cloud User to Join One Project](#tag/Projects/operation/createProjectInvitation) endpoint.

**Note**: To return both pending and active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users will be returned. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team whose application users you want to return.
	@return ListTeamUsersApiRequest
*/
func (a *MongoDBCloudUsersApiService) ListTeamUsers(ctx context.Context, orgId string, teamId string) ListTeamUsersApiRequest {
	return ListTeamUsersApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		teamId:     teamId,
	}
}

// ListTeamUsersExecute executes the request
//
//	@return PaginatedOrgUser
func (a *MongoDBCloudUsersApiService) ListTeamUsersExecute(r ListTeamUsersApiRequest) (*PaginatedOrgUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOrgUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.ListTeamUsers")
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
	if r.username != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "username", r.username, "")
	}
	if r.orgMembershipStatus != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatus", r.orgMembershipStatus, "")
	}
	if r.orgMembershipStatuses != nil {
		t := *r.orgMembershipStatuses
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgMembershipStatuses", t, "multi")

	}
	if r.userId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "userId", r.userId, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type RemoveGroupUserApiRequest struct {
	ctx        context.Context
	ApiService MongoDBCloudUsersApi
	groupId    string
	userId     string
}

type RemoveGroupUserApiParams struct {
	GroupId string
	UserId  string
}

func (a *MongoDBCloudUsersApiService) RemoveGroupUserWithParams(ctx context.Context, args *RemoveGroupUserApiParams) RemoveGroupUserApiRequest {
	return RemoveGroupUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		userId:     args.UserId,
	}
}

func (r RemoveGroupUserApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RemoveGroupUserExecute(r)
}

/*
RemoveGroupUser Remove One MongoDB Cloud User from One Project

Removes one MongoDB Cloud user from the specified project. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: This resource cannot be used to remove pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

**Note**: To remove pending or active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users can be removed. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the [Return All MongoDB Cloud Users in One Project](#tag/MongoDB-Cloud-Users/operation/listProjectUsers) resource and filter by `username`.
	@return RemoveGroupUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) RemoveGroupUser(ctx context.Context, groupId string, userId string) RemoveGroupUserApiRequest {
	return RemoveGroupUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		userId:     userId,
	}
}

// RemoveGroupUserExecute executes the request
func (a *MongoDBCloudUsersApiService) RemoveGroupUserExecute(r RemoveGroupUserApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.RemoveGroupUser")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users/{userId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type RemoveGroupUserRoleApiRequest struct {
	ctx                  context.Context
	ApiService           MongoDBCloudUsersApi
	groupId              string
	userId               string
	addOrRemoveGroupRole *AddOrRemoveGroupRole
}

type RemoveGroupUserRoleApiParams struct {
	GroupId              string
	UserId               string
	AddOrRemoveGroupRole *AddOrRemoveGroupRole
}

func (a *MongoDBCloudUsersApiService) RemoveGroupUserRoleWithParams(ctx context.Context, args *RemoveGroupUserRoleApiParams) RemoveGroupUserRoleApiRequest {
	return RemoveGroupUserRoleApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		groupId:              args.GroupId,
		userId:               args.UserId,
		addOrRemoveGroupRole: args.AddOrRemoveGroupRole,
	}
}

func (r RemoveGroupUserRoleApiRequest) Execute() (*GroupUserResponse, *http.Response, error) {
	return r.ApiService.RemoveGroupUserRoleExecute(r)
}

/*
RemoveGroupUserRole Remove One Project Role from One MongoDB Cloud User

Removes one project-level role from the MongoDB Cloud user. You can remove a role from an active user or a user that has been invited to join the project. To replace a user's only role, add the new role before removing the old role. A user must have at least one role at all times.

**Note**: This resource cannot be used to remove a role from users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the project. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Project resource and filter by `username`.
	@return RemoveGroupUserRoleApiRequest
*/
func (a *MongoDBCloudUsersApiService) RemoveGroupUserRole(ctx context.Context, groupId string, userId string, addOrRemoveGroupRole *AddOrRemoveGroupRole) RemoveGroupUserRoleApiRequest {
	return RemoveGroupUserRoleApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		groupId:              groupId,
		userId:               userId,
		addOrRemoveGroupRole: addOrRemoveGroupRole,
	}
}

// RemoveGroupUserRoleExecute executes the request
//
//	@return GroupUserResponse
func (a *MongoDBCloudUsersApiService) RemoveGroupUserRoleExecute(r RemoveGroupUserRoleApiRequest) (*GroupUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.RemoveGroupUserRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users/{userId}:removeRole"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.addOrRemoveGroupRole == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveGroupRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveGroupRole
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

type RemoveOrgRoleApiRequest struct {
	ctx                context.Context
	ApiService         MongoDBCloudUsersApi
	orgId              string
	userId             string
	addOrRemoveOrgRole *AddOrRemoveOrgRole
}

type RemoveOrgRoleApiParams struct {
	OrgId              string
	UserId             string
	AddOrRemoveOrgRole *AddOrRemoveOrgRole
}

func (a *MongoDBCloudUsersApiService) RemoveOrgRoleWithParams(ctx context.Context, args *RemoveOrgRoleApiParams) RemoveOrgRoleApiRequest {
	return RemoveOrgRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		orgId:              args.OrgId,
		userId:             args.UserId,
		addOrRemoveOrgRole: args.AddOrRemoveOrgRole,
	}
}

func (r RemoveOrgRoleApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.RemoveOrgRoleExecute(r)
}

/*
RemoveOrgRole Remove One Organization Role from One MongoDB Cloud User

Removes one organization-level role from the MongoDB Cloud user. You can remove a role from an active user or a user that has not yet accepted the invitation to join the organization. To replace a user's only role, add the new role before removing the old role. A user must have at least one role at all times.

**Note**: This operation is atomic.

**Note**: This resource cannot be used to remove a role from users invited using the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
	@return RemoveOrgRoleApiRequest
*/
func (a *MongoDBCloudUsersApiService) RemoveOrgRole(ctx context.Context, orgId string, userId string, addOrRemoveOrgRole *AddOrRemoveOrgRole) RemoveOrgRoleApiRequest {
	return RemoveOrgRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		orgId:              orgId,
		userId:             userId,
		addOrRemoveOrgRole: addOrRemoveOrgRole,
	}
}

// RemoveOrgRoleExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) RemoveOrgRoleExecute(r RemoveOrgRoleApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.RemoveOrgRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}:removeRole"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.addOrRemoveOrgRole == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveOrgRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveOrgRole
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

type RemoveOrgTeamUserApiRequest struct {
	ctx                     context.Context
	ApiService              MongoDBCloudUsersApi
	orgId                   string
	teamId                  string
	addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam
}

type RemoveOrgTeamUserApiParams struct {
	OrgId                   string
	TeamId                  string
	AddOrRemoveUserFromTeam *AddOrRemoveUserFromTeam
}

func (a *MongoDBCloudUsersApiService) RemoveOrgTeamUserWithParams(ctx context.Context, args *RemoveOrgTeamUserApiParams) RemoveOrgTeamUserApiRequest {
	return RemoveOrgTeamUserApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		orgId:                   args.OrgId,
		teamId:                  args.TeamId,
		addOrRemoveUserFromTeam: args.AddOrRemoveUserFromTeam,
	}
}

func (r RemoveOrgTeamUserApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.RemoveOrgTeamUserExecute(r)
}

/*
RemoveOrgTeamUser Remove One MongoDB Cloud User from One Team

Removes one MongoDB Cloud user from one team. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: This resource cannot be used to remove a user invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param teamId Unique 24-hexadecimal digit string that identifies the team to remove the MongoDB user from.
	@return RemoveOrgTeamUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) RemoveOrgTeamUser(ctx context.Context, orgId string, teamId string, addOrRemoveUserFromTeam *AddOrRemoveUserFromTeam) RemoveOrgTeamUserApiRequest {
	return RemoveOrgTeamUserApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		orgId:                   orgId,
		teamId:                  teamId,
		addOrRemoveUserFromTeam: addOrRemoveUserFromTeam,
	}
}

// RemoveOrgTeamUserExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) RemoveOrgTeamUserExecute(r RemoveOrgTeamUserApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.RemoveOrgTeamUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/teams/{teamId}:removeUser"
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
	if r.addOrRemoveUserFromTeam == nil {
		return localVarReturnValue, nil, reportError("addOrRemoveUserFromTeam is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.addOrRemoveUserFromTeam
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

type RemoveOrgUserApiRequest struct {
	ctx        context.Context
	ApiService MongoDBCloudUsersApi
	orgId      string
	userId     string
}

type RemoveOrgUserApiParams struct {
	OrgId  string
	UserId string
}

func (a *MongoDBCloudUsersApiService) RemoveOrgUserWithParams(ctx context.Context, args *RemoveOrgUserApiParams) RemoveOrgUserApiRequest {
	return RemoveOrgUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		userId:     args.UserId,
	}
}

func (r RemoveOrgUserApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RemoveOrgUserExecute(r)
}

/*
RemoveOrgUser Remove One MongoDB Cloud User from One Organization

Removes one MongoDB Cloud user in the specified organization. You can remove an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: This resource cannot be used to remove pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

**Note**: To remove pending or active users, use v2-{2025-02-19} or later. If using a deprecated version, only active users can be removed. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the [Return All MongoDB Cloud Users in One Organization](#tag/MongoDB-Cloud-Users/operation/listOrganizationUsers) resource and filter by `username`.
	@return RemoveOrgUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) RemoveOrgUser(ctx context.Context, orgId string, userId string) RemoveOrgUserApiRequest {
	return RemoveOrgUserApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		userId:     userId,
	}
}

// RemoveOrgUserExecute executes the request
func (a *MongoDBCloudUsersApiService) RemoveOrgUserExecute(r RemoveOrgUserApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.RemoveOrgUser")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

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

type UpdateOrgUserApiRequest struct {
	ctx                  context.Context
	ApiService           MongoDBCloudUsersApi
	orgId                string
	userId               string
	orgUserUpdateRequest *OrgUserUpdateRequest
}

type UpdateOrgUserApiParams struct {
	OrgId                string
	UserId               string
	OrgUserUpdateRequest *OrgUserUpdateRequest
}

func (a *MongoDBCloudUsersApiService) UpdateOrgUserWithParams(ctx context.Context, args *UpdateOrgUserApiParams) UpdateOrgUserApiRequest {
	return UpdateOrgUserApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		orgId:                args.OrgId,
		userId:               args.UserId,
		orgUserUpdateRequest: args.OrgUserUpdateRequest,
	}
}

func (r UpdateOrgUserApiRequest) Execute() (*OrgUserResponse, *http.Response, error) {
	return r.ApiService.UpdateOrgUserExecute(r)
}

/*
UpdateOrgUser Update One MongoDB Cloud User in One Organization

Updates one MongoDB Cloud user in the specified organization. You can update an active user or a user that has not yet accepted the invitation to join the organization.

**Note**: Only include the fields you wish to update in the request body. Supplying a field with an empty value will reset that field on the user.

**Note**: This resource cannot be used to update pending users invited via the deprecated Invite One MongoDB Cloud User to Join One Project endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the pending or active user in the organization. If you need to lookup a user's `userId` or verify a user's status in the organization, use the Return All MongoDB Cloud Users in One Organization resource and filter by `username`.
	@return UpdateOrgUserApiRequest
*/
func (a *MongoDBCloudUsersApiService) UpdateOrgUser(ctx context.Context, orgId string, userId string, orgUserUpdateRequest *OrgUserUpdateRequest) UpdateOrgUserApiRequest {
	return UpdateOrgUserApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		orgId:                orgId,
		userId:               userId,
		orgUserUpdateRequest: orgUserUpdateRequest,
	}
}

// UpdateOrgUserExecute executes the request
//
//	@return OrgUserResponse
func (a *MongoDBCloudUsersApiService) UpdateOrgUserExecute(r UpdateOrgUserApiRequest) (*OrgUserResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgUserResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MongoDBCloudUsersApiService.UpdateOrgUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.orgUserUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("orgUserUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-02-19+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.orgUserUpdateRequest
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
