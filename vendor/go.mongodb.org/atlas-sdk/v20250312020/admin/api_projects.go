// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ProjectsApi interface {

	/*
		AddGroupUser Add One MongoDB Cloud User to One Project

		Adds one MongoDB Cloud user to the specified project. If the MongoDB Cloud user is not a member of the project's organization, then the user must accept their invitation to the organization to access information within the specified project. If the MongoDB Cloud User is already a member of the project's organization, then they will be added to the project immediately and an invitation will not be returned by this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupInvitationRequest Adds one MongoDB Cloud user to the specified project.
		@return AddGroupUserApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	AddGroupUser(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) AddGroupUserApiRequest
	/*
		AddGroupUser Add One MongoDB Cloud User to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddGroupUserApiParams - Parameters for the request
		@return AddGroupUserApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	AddGroupUserWithParams(ctx context.Context, args *AddGroupUserApiParams) AddGroupUserApiRequest

	// Method available only for mocking purposes
	AddGroupUserExecute(r AddGroupUserApiRequest) (*OrganizationInvitation, *http.Response, error)

	/*
		CreateGroup Create One Project

		Creates one project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param group Creates one project.
		@return CreateGroupApiRequest
	*/
	CreateGroup(ctx context.Context, group *Group) CreateGroupApiRequest
	/*
		CreateGroup Create One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupApiParams - Parameters for the request
		@return CreateGroupApiRequest
	*/
	CreateGroupWithParams(ctx context.Context, args *CreateGroupApiParams) CreateGroupApiRequest

	// Method available only for mocking purposes
	CreateGroupExecute(r CreateGroupApiRequest) (*Group, *http.Response, error)

	/*
		CreateGroupInvite Create Invitation for One MongoDB Cloud User in One Project

		Invites one MongoDB Cloud user to join the specified project. The MongoDB Cloud user must accept the invitation to access information within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupInvitationRequest Invites one MongoDB Cloud user to join the specified project.
		@return CreateGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	CreateGroupInvite(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) CreateGroupInviteApiRequest
	/*
		CreateGroupInvite Create Invitation for One MongoDB Cloud User in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupInviteApiParams - Parameters for the request
		@return CreateGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	CreateGroupInviteWithParams(ctx context.Context, args *CreateGroupInviteApiParams) CreateGroupInviteApiRequest

	// Method available only for mocking purposes
	CreateGroupInviteExecute(r CreateGroupInviteApiRequest) (*GroupInvitation, *http.Response, error)

	/*
		DeleteGroup Remove One Project

		Removes the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings. You can delete a project only if there are no Online Archives for the clusters in the project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteGroupApiRequest
	*/
	DeleteGroup(ctx context.Context, groupId string) DeleteGroupApiRequest
	/*
		DeleteGroup Remove One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupApiParams - Parameters for the request
		@return DeleteGroupApiRequest
	*/
	DeleteGroupWithParams(ctx context.Context, args *DeleteGroupApiParams) DeleteGroupApiRequest

	// Method available only for mocking purposes
	DeleteGroupExecute(r DeleteGroupApiRequest) (*http.Response, error)

	/*
		DeleteGroupInvite Remove One Invitation from One Project

		Cancels one pending invitation sent to the specified MongoDB Cloud user to join a project. You can't cancel an invitation that the user accepted. Note: deleting a project invitation does not delete an organization invitation even if they were created together.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
		@return DeleteGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	DeleteGroupInvite(ctx context.Context, groupId string, invitationId string) DeleteGroupInviteApiRequest
	/*
		DeleteGroupInvite Remove One Invitation from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupInviteApiParams - Parameters for the request
		@return DeleteGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	DeleteGroupInviteWithParams(ctx context.Context, args *DeleteGroupInviteApiParams) DeleteGroupInviteApiRequest

	// Method available only for mocking purposes
	DeleteGroupInviteExecute(r DeleteGroupInviteApiRequest) (*http.Response, error)

	/*
		DeleteGroupLimit Remove One Project Limit

		Removes the specified project limit. Depending on the limit, Atlas either resets the limit to its default value or removes the limit entirely.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteGroupLimitApiRequest
	*/
	DeleteGroupLimit(ctx context.Context, limitName string, groupId string) DeleteGroupLimitApiRequest
	/*
		DeleteGroupLimit Remove One Project Limit


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupLimitApiParams - Parameters for the request
		@return DeleteGroupLimitApiRequest
	*/
	DeleteGroupLimitWithParams(ctx context.Context, args *DeleteGroupLimitApiParams) DeleteGroupLimitApiRequest

	// Method available only for mocking purposes
	DeleteGroupLimitExecute(r DeleteGroupLimitApiRequest) (*http.Response, error)

	/*
		GetGroup Return One Project

		Returns details about the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupApiRequest
	*/
	GetGroup(ctx context.Context, groupId string) GetGroupApiRequest
	/*
		GetGroup Return One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupApiParams - Parameters for the request
		@return GetGroupApiRequest
	*/
	GetGroupWithParams(ctx context.Context, args *GetGroupApiParams) GetGroupApiRequest

	// Method available only for mocking purposes
	GetGroupExecute(r GetGroupApiRequest) (*Group, *http.Response, error)

	/*
		GetGroupByName Return One Project by Name

		Returns details about the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupName Human-readable label that identifies this project.
		@return GetGroupByNameApiRequest
	*/
	GetGroupByName(ctx context.Context, groupName string) GetGroupByNameApiRequest
	/*
		GetGroupByName Return One Project by Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupByNameApiParams - Parameters for the request
		@return GetGroupByNameApiRequest
	*/
	GetGroupByNameWithParams(ctx context.Context, args *GetGroupByNameApiParams) GetGroupByNameApiRequest

	// Method available only for mocking purposes
	GetGroupByNameExecute(r GetGroupByNameApiRequest) (*Group, *http.Response, error)

	/*
		GetGroupInvite Return One Invitation in One Project by Invitation ID

		Returns the details of one pending invitation to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
		@return GetGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	GetGroupInvite(ctx context.Context, groupId string, invitationId string) GetGroupInviteApiRequest
	/*
		GetGroupInvite Return One Invitation in One Project by Invitation ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupInviteApiParams - Parameters for the request
		@return GetGroupInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	GetGroupInviteWithParams(ctx context.Context, args *GetGroupInviteApiParams) GetGroupInviteApiRequest

	// Method available only for mocking purposes
	GetGroupInviteExecute(r GetGroupInviteApiRequest) (*GroupInvitation, *http.Response, error)

	/*
		GetGroupIpAddresses Return All IP Addresses for One Project

		Returns all IP addresses for this project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupIpAddressesApiRequest
	*/
	GetGroupIpAddresses(ctx context.Context, groupId string) GetGroupIpAddressesApiRequest
	/*
		GetGroupIpAddresses Return All IP Addresses for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupIpAddressesApiParams - Parameters for the request
		@return GetGroupIpAddressesApiRequest
	*/
	GetGroupIpAddressesWithParams(ctx context.Context, args *GetGroupIpAddressesApiParams) GetGroupIpAddressesApiRequest

	// Method available only for mocking purposes
	GetGroupIpAddressesExecute(r GetGroupIpAddressesApiRequest) (*GroupIPAddresses, *http.Response, error)

	/*
		GetGroupLimit Return One Limit for One Project

		Returns the specified limit for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupLimitApiRequest
	*/
	GetGroupLimit(ctx context.Context, limitName string, groupId string) GetGroupLimitApiRequest
	/*
		GetGroupLimit Return One Limit for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupLimitApiParams - Parameters for the request
		@return GetGroupLimitApiRequest
	*/
	GetGroupLimitWithParams(ctx context.Context, args *GetGroupLimitApiParams) GetGroupLimitApiRequest

	// Method available only for mocking purposes
	GetGroupLimitExecute(r GetGroupLimitApiRequest) (*DataFederationLimit, *http.Response, error)

	/*
		GetGroupSettings Return Project Settings

		Returns details about the specified project's settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupSettingsApiRequest
	*/
	GetGroupSettings(ctx context.Context, groupId string) GetGroupSettingsApiRequest
	/*
		GetGroupSettings Return Project Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupSettingsApiParams - Parameters for the request
		@return GetGroupSettingsApiRequest
	*/
	GetGroupSettingsWithParams(ctx context.Context, args *GetGroupSettingsApiParams) GetGroupSettingsApiRequest

	// Method available only for mocking purposes
	GetGroupSettingsExecute(r GetGroupSettingsApiRequest) (*GroupSettings, *http.Response, error)

	/*
		GetMongoDbVersions Return All Available MongoDB LTS Versions for Clusters in One Project

		Returns the MongoDB Long Term Support Major Versions available to new clusters in this project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetMongoDbVersionsApiRequest
	*/
	GetMongoDbVersions(ctx context.Context, groupId string) GetMongoDbVersionsApiRequest
	/*
		GetMongoDbVersions Return All Available MongoDB LTS Versions for Clusters in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetMongoDbVersionsApiParams - Parameters for the request
		@return GetMongoDbVersionsApiRequest
	*/
	GetMongoDbVersionsWithParams(ctx context.Context, args *GetMongoDbVersionsApiParams) GetMongoDbVersionsApiRequest

	// Method available only for mocking purposes
	GetMongoDbVersionsExecute(r GetMongoDbVersionsApiRequest) (*PaginatedAvailableVersion, *http.Response, error)

	/*
		ListGroupInvites Return All Invitations in One Project

		Returns all pending invitations to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	ListGroupInvites(ctx context.Context, groupId string) ListGroupInvitesApiRequest
	/*
		ListGroupInvites Return All Invitations in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupInvitesApiParams - Parameters for the request
		@return ListGroupInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	ListGroupInvitesWithParams(ctx context.Context, args *ListGroupInvitesApiParams) ListGroupInvitesApiRequest

	// Method available only for mocking purposes
	ListGroupInvitesExecute(r ListGroupInvitesApiRequest) ([]GroupInvitation, *http.Response, error)

	/*
		ListGroupLimits Return All Limits for One Project

		Returns all the limits for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupLimitsApiRequest
	*/
	ListGroupLimits(ctx context.Context, groupId string) ListGroupLimitsApiRequest
	/*
		ListGroupLimits Return All Limits for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupLimitsApiParams - Parameters for the request
		@return ListGroupLimitsApiRequest
	*/
	ListGroupLimitsWithParams(ctx context.Context, args *ListGroupLimitsApiParams) ListGroupLimitsApiRequest

	// Method available only for mocking purposes
	ListGroupLimitsExecute(r ListGroupLimitsApiRequest) ([]DataFederationLimit, *http.Response, error)

	/*
		ListGroups Return All Projects

		Returns details about all projects. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListGroupsApiRequest
	*/
	ListGroups(ctx context.Context) ListGroupsApiRequest
	/*
		ListGroups Return All Projects


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupsApiParams - Parameters for the request
		@return ListGroupsApiRequest
	*/
	ListGroupsWithParams(ctx context.Context, args *ListGroupsApiParams) ListGroupsApiRequest

	// Method available only for mocking purposes
	ListGroupsExecute(r ListGroupsApiRequest) (*PaginatedAtlasGroup, *http.Response, error)

	/*
		MigrateGroup Migrate One Project to Another Organization

		Migrates a project from its current organization to another organization. All project users and their roles will be copied to the same project in the destination organization. You must include an organization API key with the Organization Owner role for the destination organization to verify access to the destination organization when you authenticate with Programmatic API Keys. Otherwise, the requesting user must have the Organization Owner role in both organizations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupMigrationRequest Migrates a project from its current organization to another organization.
		@return MigrateGroupApiRequest
	*/
	MigrateGroup(ctx context.Context, groupId string, groupMigrationRequest *GroupMigrationRequest) MigrateGroupApiRequest
	/*
		MigrateGroup Migrate One Project to Another Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param MigrateGroupApiParams - Parameters for the request
		@return MigrateGroupApiRequest
	*/
	MigrateGroupWithParams(ctx context.Context, args *MigrateGroupApiParams) MigrateGroupApiRequest

	// Method available only for mocking purposes
	MigrateGroupExecute(r MigrateGroupApiRequest) (*Group, *http.Response, error)

	/*
			SetGroupLimit Set One Project Limit

			Sets the specified project limit.

		**NOTE**: Increasing the following configuration limits might lead to slower response times in the MongoDB Cloud UI or increased user management overhead leading to authentication or authorization re-architecture. If possible, we recommend that you create additional projects to gain access to more of these resources for a more sustainable growth pattern.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param dataFederationLimit Limit to update.
			@return SetGroupLimitApiRequest
	*/
	SetGroupLimit(ctx context.Context, limitName string, groupId string, dataFederationLimit *DataFederationLimit) SetGroupLimitApiRequest
	/*
		SetGroupLimit Set One Project Limit


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param SetGroupLimitApiParams - Parameters for the request
		@return SetGroupLimitApiRequest
	*/
	SetGroupLimitWithParams(ctx context.Context, args *SetGroupLimitApiParams) SetGroupLimitApiRequest

	// Method available only for mocking purposes
	SetGroupLimitExecute(r SetGroupLimitApiRequest) (*DataFederationLimit, *http.Response, error)

	/*
		UpdateGroup Update One Project

		Updates the human-readable label that identifies the specified project, or the tags associated with the project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupUpdate Project to update.
		@return UpdateGroupApiRequest
	*/
	UpdateGroup(ctx context.Context, groupId string, groupUpdate *GroupUpdate) UpdateGroupApiRequest
	/*
		UpdateGroup Update One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupApiParams - Parameters for the request
		@return UpdateGroupApiRequest
	*/
	UpdateGroupWithParams(ctx context.Context, args *UpdateGroupApiParams) UpdateGroupApiRequest

	// Method available only for mocking purposes
	UpdateGroupExecute(r UpdateGroupApiRequest) (*Group, *http.Response, error)

	/*
		UpdateGroupInvites Update One Invitation in One Project

		Updates the details of one pending invitation to the specified project. To specify which invitation to update, provide the username of the invited user.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupInvitationRequest Updates the details of one pending invitation to the specified project.
		@return UpdateGroupInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateGroupInvites(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) UpdateGroupInvitesApiRequest
	/*
		UpdateGroupInvites Update One Invitation in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupInvitesApiParams - Parameters for the request
		@return UpdateGroupInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateGroupInvitesWithParams(ctx context.Context, args *UpdateGroupInvitesApiParams) UpdateGroupInvitesApiRequest

	// Method available only for mocking purposes
	UpdateGroupInvitesExecute(r UpdateGroupInvitesApiRequest) (*GroupInvitation, *http.Response, error)

	/*
		UpdateGroupSettings Update Project Settings

		Updates the settings of the specified project. You can update any of the options available. MongoDB cloud only updates the options provided in the request.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupSettings Settings to update.
		@return UpdateGroupSettingsApiRequest
	*/
	UpdateGroupSettings(ctx context.Context, groupId string, groupSettings *GroupSettings) UpdateGroupSettingsApiRequest
	/*
		UpdateGroupSettings Update Project Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupSettingsApiParams - Parameters for the request
		@return UpdateGroupSettingsApiRequest
	*/
	UpdateGroupSettingsWithParams(ctx context.Context, args *UpdateGroupSettingsApiParams) UpdateGroupSettingsApiRequest

	// Method available only for mocking purposes
	UpdateGroupSettingsExecute(r UpdateGroupSettingsApiRequest) (*GroupSettings, *http.Response, error)

	/*
		UpdateGroupUserRoles Update Project Roles for One MongoDB Cloud User

		Updates the roles of the specified user in the specified project. To specify the user to update, provide the unique 24-hexadecimal digit string that identifies the user in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param userId Unique 24-hexadecimal digit string that identifies the user to modify.
		@param updateGroupRolesForUser Roles to update for the specified user.
		@return UpdateGroupUserRolesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateGroupUserRoles(ctx context.Context, groupId string, userId string, updateGroupRolesForUser *UpdateGroupRolesForUser) UpdateGroupUserRolesApiRequest
	/*
		UpdateGroupUserRoles Update Project Roles for One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupUserRolesApiParams - Parameters for the request
		@return UpdateGroupUserRolesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateGroupUserRolesWithParams(ctx context.Context, args *UpdateGroupUserRolesApiParams) UpdateGroupUserRolesApiRequest

	// Method available only for mocking purposes
	UpdateGroupUserRolesExecute(r UpdateGroupUserRolesApiRequest) (*UpdateGroupRolesForUser, *http.Response, error)

	/*
		UpdateInviteById Update One Invitation in One Project by Invitation ID

		Updates the details of one pending invitation to the specified project. To specify which invitation to update, provide the unique identification string for that invitation. Use the Return All Project Invitations endpoint to retrieve IDs for all pending project invitations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
		@param groupInvitationUpdateRequest Updates the details of one pending invitation to the specified project.
		@return UpdateInviteByIdApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateInviteById(ctx context.Context, groupId string, invitationId string, groupInvitationUpdateRequest *GroupInvitationUpdateRequest) UpdateInviteByIdApiRequest
	/*
		UpdateInviteById Update One Invitation in One Project by Invitation ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateInviteByIdApiParams - Parameters for the request
		@return UpdateInviteByIdApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ProjectsApi
	*/
	UpdateInviteByIdWithParams(ctx context.Context, args *UpdateInviteByIdApiParams) UpdateInviteByIdApiRequest

	// Method available only for mocking purposes
	UpdateInviteByIdExecute(r UpdateInviteByIdApiRequest) (*GroupInvitation, *http.Response, error)
}

// ProjectsApiService ProjectsApi service
type ProjectsApiService service

type AddGroupUserApiRequest struct {
	ctx                    context.Context
	ApiService             ProjectsApi
	groupId                string
	groupInvitationRequest *GroupInvitationRequest
}

type AddGroupUserApiParams struct {
	GroupId                string
	GroupInvitationRequest *GroupInvitationRequest
}

func (a *ProjectsApiService) AddGroupUserWithParams(ctx context.Context, args *AddGroupUserApiParams) AddGroupUserApiRequest {
	return AddGroupUserApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		groupInvitationRequest: args.GroupInvitationRequest,
	}
}

func (r AddGroupUserApiRequest) Execute() (*OrganizationInvitation, *http.Response, error) {
	return r.ApiService.AddGroupUserExecute(r)
}

/*
AddGroupUser Add One MongoDB Cloud User to One Project

Adds one MongoDB Cloud user to the specified project. If the MongoDB Cloud user is not a member of the project's organization, then the user must accept their invitation to the organization to access information within the specified project. If the MongoDB Cloud User is already a member of the project's organization, then they will be added to the project immediately and an invitation will not be returned by this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return AddGroupUserApiRequest

Deprecated
*/
func (a *ProjectsApiService) AddGroupUser(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) AddGroupUserApiRequest {
	return AddGroupUserApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		groupInvitationRequest: groupInvitationRequest,
	}
}

// AddGroupUserExecute executes the request
//
//	@return OrganizationInvitation
//
// Deprecated
func (a *ProjectsApiService) AddGroupUserExecute(r AddGroupUserApiRequest) (*OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.AddGroupUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/access"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupInvitationRequest == nil {
		return localVarReturnValue, nil, reportError("groupInvitationRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupInvitationRequest
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

type CreateGroupApiRequest struct {
	ctx            context.Context
	ApiService     ProjectsApi
	group          *Group
	projectOwnerId *string
}

type CreateGroupApiParams struct {
	Group          *Group
	ProjectOwnerId *string
}

func (a *ProjectsApiService) CreateGroupWithParams(ctx context.Context, args *CreateGroupApiParams) CreateGroupApiRequest {
	return CreateGroupApiRequest{
		ApiService:     a,
		ctx:            ctx,
		group:          args.Group,
		projectOwnerId: args.ProjectOwnerId,
	}
}

// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.
func (r CreateGroupApiRequest) ProjectOwnerId(projectOwnerId string) CreateGroupApiRequest {
	r.projectOwnerId = &projectOwnerId
	return r
}

func (r CreateGroupApiRequest) Execute() (*Group, *http.Response, error) {
	return r.ApiService.CreateGroupExecute(r)
}

/*
CreateGroup Create One Project

Creates one project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return CreateGroupApiRequest
*/
func (a *ProjectsApiService) CreateGroup(ctx context.Context, group *Group) CreateGroupApiRequest {
	return CreateGroupApiRequest{
		ApiService: a,
		ctx:        ctx,
		group:      group,
	}
}

// CreateGroupExecute executes the request
//
//	@return Group
func (a *ProjectsApiService) CreateGroupExecute(r CreateGroupApiRequest) (*Group, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Group
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.CreateGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.group == nil {
		return localVarReturnValue, nil, reportError("group is required and must be specified")
	}

	if r.projectOwnerId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectOwnerId", r.projectOwnerId, "")
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
	localVarPostBody = r.group
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

type CreateGroupInviteApiRequest struct {
	ctx                    context.Context
	ApiService             ProjectsApi
	groupId                string
	groupInvitationRequest *GroupInvitationRequest
}

type CreateGroupInviteApiParams struct {
	GroupId                string
	GroupInvitationRequest *GroupInvitationRequest
}

func (a *ProjectsApiService) CreateGroupInviteWithParams(ctx context.Context, args *CreateGroupInviteApiParams) CreateGroupInviteApiRequest {
	return CreateGroupInviteApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		groupInvitationRequest: args.GroupInvitationRequest,
	}
}

func (r CreateGroupInviteApiRequest) Execute() (*GroupInvitation, *http.Response, error) {
	return r.ApiService.CreateGroupInviteExecute(r)
}

/*
CreateGroupInvite Create Invitation for One MongoDB Cloud User in One Project

Invites one MongoDB Cloud user to join the specified project. The MongoDB Cloud user must accept the invitation to access information within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupInviteApiRequest

Deprecated
*/
func (a *ProjectsApiService) CreateGroupInvite(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) CreateGroupInviteApiRequest {
	return CreateGroupInviteApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		groupInvitationRequest: groupInvitationRequest,
	}
}

// CreateGroupInviteExecute executes the request
//
//	@return GroupInvitation
//
// Deprecated
func (a *ProjectsApiService) CreateGroupInviteExecute(r CreateGroupInviteApiRequest) (*GroupInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.CreateGroupInvite")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupInvitationRequest == nil {
		return localVarReturnValue, nil, reportError("groupInvitationRequest is required and must be specified")
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
	localVarPostBody = r.groupInvitationRequest
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

type DeleteGroupApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
}

type DeleteGroupApiParams struct {
	GroupId string
}

func (a *ProjectsApiService) DeleteGroupWithParams(ctx context.Context, args *DeleteGroupApiParams) DeleteGroupApiRequest {
	return DeleteGroupApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DeleteGroupApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupExecute(r)
}

/*
DeleteGroup Remove One Project

Removes the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings. You can delete a project only if there are no Online Archives for the clusters in the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteGroupApiRequest
*/
func (a *ProjectsApiService) DeleteGroup(ctx context.Context, groupId string) DeleteGroupApiRequest {
	return DeleteGroupApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DeleteGroupExecute executes the request
func (a *ProjectsApiService) DeleteGroupExecute(r DeleteGroupApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.DeleteGroup")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type DeleteGroupInviteApiRequest struct {
	ctx          context.Context
	ApiService   ProjectsApi
	groupId      string
	invitationId string
}

type DeleteGroupInviteApiParams struct {
	GroupId      string
	InvitationId string
}

func (a *ProjectsApiService) DeleteGroupInviteWithParams(ctx context.Context, args *DeleteGroupInviteApiParams) DeleteGroupInviteApiRequest {
	return DeleteGroupInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		invitationId: args.InvitationId,
	}
}

func (r DeleteGroupInviteApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupInviteExecute(r)
}

/*
DeleteGroupInvite Remove One Invitation from One Project

Cancels one pending invitation sent to the specified MongoDB Cloud user to join a project. You can't cancel an invitation that the user accepted. Note: deleting a project invitation does not delete an organization invitation even if they were created together.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return DeleteGroupInviteApiRequest

Deprecated
*/
func (a *ProjectsApiService) DeleteGroupInvite(ctx context.Context, groupId string, invitationId string) DeleteGroupInviteApiRequest {
	return DeleteGroupInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		invitationId: invitationId,
	}
}

// DeleteGroupInviteExecute executes the request
// Deprecated
func (a *ProjectsApiService) DeleteGroupInviteExecute(r DeleteGroupInviteApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.DeleteGroupInvite")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites/{invitationId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.invitationId == "" {
		return nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

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

type DeleteGroupLimitApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	limitName  string
	groupId    string
}

type DeleteGroupLimitApiParams struct {
	LimitName string
	GroupId   string
}

func (a *ProjectsApiService) DeleteGroupLimitWithParams(ctx context.Context, args *DeleteGroupLimitApiParams) DeleteGroupLimitApiRequest {
	return DeleteGroupLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		limitName:  args.LimitName,
		groupId:    args.GroupId,
	}
}

func (r DeleteGroupLimitApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupLimitExecute(r)
}

/*
DeleteGroupLimit Remove One Project Limit

Removes the specified project limit. Depending on the limit, Atlas either resets the limit to its default value or removes the limit entirely.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteGroupLimitApiRequest
*/
func (a *ProjectsApiService) DeleteGroupLimit(ctx context.Context, limitName string, groupId string) DeleteGroupLimitApiRequest {
	return DeleteGroupLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		limitName:  limitName,
		groupId:    groupId,
	}
}

// DeleteGroupLimitExecute executes the request
func (a *ProjectsApiService) DeleteGroupLimitExecute(r DeleteGroupLimitApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.DeleteGroupLimit")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/limits/{limitName}"
	if r.limitName == "" {
		return nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type GetGroupApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
}

type GetGroupApiParams struct {
	GroupId string
}

func (a *ProjectsApiService) GetGroupWithParams(ctx context.Context, args *GetGroupApiParams) GetGroupApiRequest {
	return GetGroupApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetGroupApiRequest) Execute() (*Group, *http.Response, error) {
	return r.ApiService.GetGroupExecute(r)
}

/*
GetGroup Return One Project

Returns details about the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupApiRequest
*/
func (a *ProjectsApiService) GetGroup(ctx context.Context, groupId string) GetGroupApiRequest {
	return GetGroupApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetGroupExecute executes the request
//
//	@return Group
func (a *ProjectsApiService) GetGroupExecute(r GetGroupApiRequest) (*Group, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Group
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type GetGroupByNameApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupName  string
}

type GetGroupByNameApiParams struct {
	GroupName string
}

func (a *ProjectsApiService) GetGroupByNameWithParams(ctx context.Context, args *GetGroupByNameApiParams) GetGroupByNameApiRequest {
	return GetGroupByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupName:  args.GroupName,
	}
}

func (r GetGroupByNameApiRequest) Execute() (*Group, *http.Response, error) {
	return r.ApiService.GetGroupByNameExecute(r)
}

/*
GetGroupByName Return One Project by Name

Returns details about the specified project. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupName Human-readable label that identifies this project.
	@return GetGroupByNameApiRequest
*/
func (a *ProjectsApiService) GetGroupByName(ctx context.Context, groupName string) GetGroupByNameApiRequest {
	return GetGroupByNameApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupName:  groupName,
	}
}

// GetGroupByNameExecute executes the request
//
//	@return Group
func (a *ProjectsApiService) GetGroupByNameExecute(r GetGroupByNameApiRequest) (*Group, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Group
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroupByName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/byName/{groupName}"
	if r.groupName == "" {
		return localVarReturnValue, nil, reportError("groupName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(r.groupName), -1)

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

type GetGroupInviteApiRequest struct {
	ctx          context.Context
	ApiService   ProjectsApi
	groupId      string
	invitationId string
}

type GetGroupInviteApiParams struct {
	GroupId      string
	InvitationId string
}

func (a *ProjectsApiService) GetGroupInviteWithParams(ctx context.Context, args *GetGroupInviteApiParams) GetGroupInviteApiRequest {
	return GetGroupInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		invitationId: args.InvitationId,
	}
}

func (r GetGroupInviteApiRequest) Execute() (*GroupInvitation, *http.Response, error) {
	return r.ApiService.GetGroupInviteExecute(r)
}

/*
GetGroupInvite Return One Invitation in One Project by Invitation ID

Returns the details of one pending invitation to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return GetGroupInviteApiRequest

Deprecated
*/
func (a *ProjectsApiService) GetGroupInvite(ctx context.Context, groupId string, invitationId string) GetGroupInviteApiRequest {
	return GetGroupInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		invitationId: invitationId,
	}
}

// GetGroupInviteExecute executes the request
//
//	@return GroupInvitation
//
// Deprecated
func (a *ProjectsApiService) GetGroupInviteExecute(r GetGroupInviteApiRequest) (*GroupInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroupInvite")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites/{invitationId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.invitationId == "" {
		return localVarReturnValue, nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

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

type GetGroupIpAddressesApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
}

type GetGroupIpAddressesApiParams struct {
	GroupId string
}

func (a *ProjectsApiService) GetGroupIpAddressesWithParams(ctx context.Context, args *GetGroupIpAddressesApiParams) GetGroupIpAddressesApiRequest {
	return GetGroupIpAddressesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetGroupIpAddressesApiRequest) Execute() (*GroupIPAddresses, *http.Response, error) {
	return r.ApiService.GetGroupIpAddressesExecute(r)
}

/*
GetGroupIpAddresses Return All IP Addresses for One Project

Returns all IP addresses for this project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupIpAddressesApiRequest
*/
func (a *ProjectsApiService) GetGroupIpAddresses(ctx context.Context, groupId string) GetGroupIpAddressesApiRequest {
	return GetGroupIpAddressesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetGroupIpAddressesExecute executes the request
//
//	@return GroupIPAddresses
func (a *ProjectsApiService) GetGroupIpAddressesExecute(r GetGroupIpAddressesApiRequest) (*GroupIPAddresses, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupIPAddresses
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroupIpAddresses")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/ipAddresses"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type GetGroupLimitApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	limitName  string
	groupId    string
}

type GetGroupLimitApiParams struct {
	LimitName string
	GroupId   string
}

func (a *ProjectsApiService) GetGroupLimitWithParams(ctx context.Context, args *GetGroupLimitApiParams) GetGroupLimitApiRequest {
	return GetGroupLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		limitName:  args.LimitName,
		groupId:    args.GroupId,
	}
}

func (r GetGroupLimitApiRequest) Execute() (*DataFederationLimit, *http.Response, error) {
	return r.ApiService.GetGroupLimitExecute(r)
}

/*
GetGroupLimit Return One Limit for One Project

Returns the specified limit for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupLimitApiRequest
*/
func (a *ProjectsApiService) GetGroupLimit(ctx context.Context, limitName string, groupId string) GetGroupLimitApiRequest {
	return GetGroupLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		limitName:  limitName,
		groupId:    groupId,
	}
}

// GetGroupLimitExecute executes the request
//
//	@return DataFederationLimit
func (a *ProjectsApiService) GetGroupLimitExecute(r GetGroupLimitApiRequest) (*DataFederationLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataFederationLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroupLimit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/limits/{limitName}"
	if r.limitName == "" {
		return localVarReturnValue, nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type GetGroupSettingsApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
}

type GetGroupSettingsApiParams struct {
	GroupId string
}

func (a *ProjectsApiService) GetGroupSettingsWithParams(ctx context.Context, args *GetGroupSettingsApiParams) GetGroupSettingsApiRequest {
	return GetGroupSettingsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetGroupSettingsApiRequest) Execute() (*GroupSettings, *http.Response, error) {
	return r.ApiService.GetGroupSettingsExecute(r)
}

/*
GetGroupSettings Return Project Settings

Returns details about the specified project's settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupSettingsApiRequest
*/
func (a *ProjectsApiService) GetGroupSettings(ctx context.Context, groupId string) GetGroupSettingsApiRequest {
	return GetGroupSettingsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetGroupSettingsExecute executes the request
//
//	@return GroupSettings
func (a *ProjectsApiService) GetGroupSettingsExecute(r GetGroupSettingsApiRequest) (*GroupSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetGroupSettings")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/settings"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type GetMongoDbVersionsApiRequest struct {
	ctx           context.Context
	ApiService    ProjectsApi
	groupId       string
	cloudProvider *string
	instanceSize  *string
	defaultStatus *string
	itemsPerPage  *int64
	pageNum       *int
}

type GetMongoDbVersionsApiParams struct {
	GroupId       string
	CloudProvider *string
	InstanceSize  *string
	DefaultStatus *string
	ItemsPerPage  *int64
	PageNum       *int
}

func (a *ProjectsApiService) GetMongoDbVersionsWithParams(ctx context.Context, args *GetMongoDbVersionsApiParams) GetMongoDbVersionsApiRequest {
	return GetMongoDbVersionsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		instanceSize:  args.InstanceSize,
		defaultStatus: args.DefaultStatus,
		itemsPerPage:  args.ItemsPerPage,
		pageNum:       args.PageNum,
	}
}

// Filter results to only one cloud provider.
func (r GetMongoDbVersionsApiRequest) CloudProvider(cloudProvider string) GetMongoDbVersionsApiRequest {
	r.cloudProvider = &cloudProvider
	return r
}

// Filter results to only one instance size.
func (r GetMongoDbVersionsApiRequest) InstanceSize(instanceSize string) GetMongoDbVersionsApiRequest {
	r.instanceSize = &instanceSize
	return r
}

// Filter results to only the default values per tier. This value must be DEFAULT.
func (r GetMongoDbVersionsApiRequest) DefaultStatus(defaultStatus string) GetMongoDbVersionsApiRequest {
	r.defaultStatus = &defaultStatus
	return r
}

// Number of items that the response returns per page.
func (r GetMongoDbVersionsApiRequest) ItemsPerPage(itemsPerPage int64) GetMongoDbVersionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetMongoDbVersionsApiRequest) PageNum(pageNum int) GetMongoDbVersionsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r GetMongoDbVersionsApiRequest) Execute() (*PaginatedAvailableVersion, *http.Response, error) {
	return r.ApiService.GetMongoDbVersionsExecute(r)
}

/*
GetMongoDbVersions Return All Available MongoDB LTS Versions for Clusters in One Project

Returns the MongoDB Long Term Support Major Versions available to new clusters in this project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetMongoDbVersionsApiRequest
*/
func (a *ProjectsApiService) GetMongoDbVersions(ctx context.Context, groupId string) GetMongoDbVersionsApiRequest {
	return GetMongoDbVersionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetMongoDbVersionsExecute executes the request
//
//	@return PaginatedAvailableVersion
func (a *ProjectsApiService) GetMongoDbVersionsExecute(r GetMongoDbVersionsApiRequest) (*PaginatedAvailableVersion, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAvailableVersion
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.GetMongoDbVersions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/mongoDBVersions"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.cloudProvider != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "cloudProvider", r.cloudProvider, "")
	}
	if r.instanceSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "instanceSize", r.instanceSize, "")
	}
	if r.defaultStatus != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "defaultStatus", r.defaultStatus, "")
	}
	if r.itemsPerPage != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	} else {
		var defaultValue int64 = 100
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

type ListGroupInvitesApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
	username   *string
}

type ListGroupInvitesApiParams struct {
	GroupId  string
	Username *string
}

func (a *ProjectsApiService) ListGroupInvitesWithParams(ctx context.Context, args *ListGroupInvitesApiParams) ListGroupInvitesApiRequest {
	return ListGroupInvitesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		username:   args.Username,
	}
}

// Email address of the user account invited to this project.
func (r ListGroupInvitesApiRequest) Username(username string) ListGroupInvitesApiRequest {
	r.username = &username
	return r
}

func (r ListGroupInvitesApiRequest) Execute() ([]GroupInvitation, *http.Response, error) {
	return r.ApiService.ListGroupInvitesExecute(r)
}

/*
ListGroupInvites Return All Invitations in One Project

Returns all pending invitations to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupInvitesApiRequest

Deprecated
*/
func (a *ProjectsApiService) ListGroupInvites(ctx context.Context, groupId string) ListGroupInvitesApiRequest {
	return ListGroupInvitesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupInvitesExecute executes the request
//
//	@return []GroupInvitation
//
// Deprecated
func (a *ProjectsApiService) ListGroupInvitesExecute(r ListGroupInvitesApiRequest) ([]GroupInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []GroupInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.ListGroupInvites")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ListGroupLimitsApiRequest struct {
	ctx        context.Context
	ApiService ProjectsApi
	groupId    string
}

type ListGroupLimitsApiParams struct {
	GroupId string
}

func (a *ProjectsApiService) ListGroupLimitsWithParams(ctx context.Context, args *ListGroupLimitsApiParams) ListGroupLimitsApiRequest {
	return ListGroupLimitsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ListGroupLimitsApiRequest) Execute() ([]DataFederationLimit, *http.Response, error) {
	return r.ApiService.ListGroupLimitsExecute(r)
}

/*
ListGroupLimits Return All Limits for One Project

Returns all the limits for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupLimitsApiRequest
*/
func (a *ProjectsApiService) ListGroupLimits(ctx context.Context, groupId string) ListGroupLimitsApiRequest {
	return ListGroupLimitsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupLimitsExecute executes the request
//
//	@return []DataFederationLimit
func (a *ProjectsApiService) ListGroupLimitsExecute(r ListGroupLimitsApiRequest) ([]DataFederationLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []DataFederationLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.ListGroupLimits")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/limits"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type ListGroupsApiRequest struct {
	ctx          context.Context
	ApiService   ProjectsApi
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupsApiParams struct {
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ProjectsApiService) ListGroupsWithParams(ctx context.Context, args *ListGroupsApiParams) ListGroupsApiRequest {
	return ListGroupsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupsApiRequest) IncludeCount(includeCount bool) ListGroupsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupsApiRequest) PageNum(pageNum int) ListGroupsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupsApiRequest) Execute() (*PaginatedAtlasGroup, *http.Response, error) {
	return r.ApiService.ListGroupsExecute(r)
}

/*
ListGroups Return All Projects

Returns details about all projects. Projects group clusters into logical collections that support an application environment, workload, or both. Each project can have its own users, teams, security, tags, and alert settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListGroupsApiRequest
*/
func (a *ProjectsApiService) ListGroups(ctx context.Context) ListGroupsApiRequest {
	return ListGroupsApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListGroupsExecute executes the request
//
//	@return PaginatedAtlasGroup
func (a *ProjectsApiService) ListGroupsExecute(r ListGroupsApiRequest) (*PaginatedAtlasGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAtlasGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.ListGroups")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups"

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

type MigrateGroupApiRequest struct {
	ctx                   context.Context
	ApiService            ProjectsApi
	groupId               string
	groupMigrationRequest *GroupMigrationRequest
}

type MigrateGroupApiParams struct {
	GroupId               string
	GroupMigrationRequest *GroupMigrationRequest
}

func (a *ProjectsApiService) MigrateGroupWithParams(ctx context.Context, args *MigrateGroupApiParams) MigrateGroupApiRequest {
	return MigrateGroupApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		groupMigrationRequest: args.GroupMigrationRequest,
	}
}

func (r MigrateGroupApiRequest) Execute() (*Group, *http.Response, error) {
	return r.ApiService.MigrateGroupExecute(r)
}

/*
MigrateGroup Migrate One Project to Another Organization

Migrates a project from its current organization to another organization. All project users and their roles will be copied to the same project in the destination organization. You must include an organization API key with the Organization Owner role for the destination organization to verify access to the destination organization when you authenticate with Programmatic API Keys. Otherwise, the requesting user must have the Organization Owner role in both organizations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return MigrateGroupApiRequest
*/
func (a *ProjectsApiService) MigrateGroup(ctx context.Context, groupId string, groupMigrationRequest *GroupMigrationRequest) MigrateGroupApiRequest {
	return MigrateGroupApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               groupId,
		groupMigrationRequest: groupMigrationRequest,
	}
}

// MigrateGroupExecute executes the request
//
//	@return Group
func (a *ProjectsApiService) MigrateGroupExecute(r MigrateGroupApiRequest) (*Group, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Group
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.MigrateGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}:migrate"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupMigrationRequest == nil {
		return localVarReturnValue, nil, reportError("groupMigrationRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-05-30+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-05-30+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupMigrationRequest
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

type SetGroupLimitApiRequest struct {
	ctx                 context.Context
	ApiService          ProjectsApi
	limitName           string
	groupId             string
	dataFederationLimit *DataFederationLimit
}

type SetGroupLimitApiParams struct {
	LimitName           string
	GroupId             string
	DataFederationLimit *DataFederationLimit
}

func (a *ProjectsApiService) SetGroupLimitWithParams(ctx context.Context, args *SetGroupLimitApiParams) SetGroupLimitApiRequest {
	return SetGroupLimitApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		limitName:           args.LimitName,
		groupId:             args.GroupId,
		dataFederationLimit: args.DataFederationLimit,
	}
}

func (r SetGroupLimitApiRequest) Execute() (*DataFederationLimit, *http.Response, error) {
	return r.ApiService.SetGroupLimitExecute(r)
}

/*
SetGroupLimit Set One Project Limit

Sets the specified project limit.

**NOTE**: Increasing the following configuration limits might lead to slower response times in the MongoDB Cloud UI or increased user management overhead leading to authentication or authorization re-architecture. If possible, we recommend that you create additional projects to gain access to more of these resources for a more sustainable growth pattern.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param limitName Human-readable label that identifies this project limit.  | Limit Name | Description | Default | API Override Limit | | --- | --- | --- | --- | | `atlas.project.deployment.clusters` | Limit on the number of clusters in this project | 25 | 100 | | `atlas.project.deployment.nodesPerPrivateLinkRegion` | Limit on the number of nodes per Private Link region in this project | 50 | 90 | | `atlas.project.security.databaseAccess.customRoles` | Limit on the number of custom roles in this project | 100 | 1400 | | `atlas.project.security.databaseAccess.users` | Limit on the number of database users in this project | 100 | 100 | | `atlas.project.security.networkAccess.crossRegionEntries` | Limit on the number of cross-region network access entries in this project | 40 | 220 | | `atlas.project.security.networkAccess.entries` | Limit on the number of network access entries in this project | 200 | 20 | | `dataFederation.bytesProcessed.query` | Limit on the number of bytes processed during a single Data Federation query | N/A | N/A | | `dataFederation.bytesProcessed.daily` | Limit on the number of bytes processed across all Data Federation tenants for the current day | N/A | N/A | | `dataFederation.bytesProcessed.weekly` | Limit on the number of bytes processed across all Data Federation tenants for the current week | N/A | N/A | | `dataFederation.bytesProcessed.monthly` | Limit on the number of bytes processed across all Data Federation tenants for the current month | N/A | N/A | | `atlas.project.deployment.privateServiceConnectionsPerRegionGroup` | Number of Private Service Connections per Region Group | 50 | 100| | `atlas.project.deployment.privateServiceConnectionsSubnetMask` | Subnet mask for GCP PSC Networks. Has lower limit of 20. | 27 | 27|
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return SetGroupLimitApiRequest
*/
func (a *ProjectsApiService) SetGroupLimit(ctx context.Context, limitName string, groupId string, dataFederationLimit *DataFederationLimit) SetGroupLimitApiRequest {
	return SetGroupLimitApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		limitName:           limitName,
		groupId:             groupId,
		dataFederationLimit: dataFederationLimit,
	}
}

// SetGroupLimitExecute executes the request
//
//	@return DataFederationLimit
func (a *ProjectsApiService) SetGroupLimitExecute(r SetGroupLimitApiRequest) (*DataFederationLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataFederationLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.SetGroupLimit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/limits/{limitName}"
	if r.limitName == "" {
		return localVarReturnValue, nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataFederationLimit == nil {
		return localVarReturnValue, nil, reportError("dataFederationLimit is required and must be specified")
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
	localVarPostBody = r.dataFederationLimit
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

type UpdateGroupApiRequest struct {
	ctx         context.Context
	ApiService  ProjectsApi
	groupId     string
	groupUpdate *GroupUpdate
}

type UpdateGroupApiParams struct {
	GroupId     string
	GroupUpdate *GroupUpdate
}

func (a *ProjectsApiService) UpdateGroupWithParams(ctx context.Context, args *UpdateGroupApiParams) UpdateGroupApiRequest {
	return UpdateGroupApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		groupUpdate: args.GroupUpdate,
	}
}

func (r UpdateGroupApiRequest) Execute() (*Group, *http.Response, error) {
	return r.ApiService.UpdateGroupExecute(r)
}

/*
UpdateGroup Update One Project

Updates the human-readable label that identifies the specified project, or the tags associated with the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateGroupApiRequest
*/
func (a *ProjectsApiService) UpdateGroup(ctx context.Context, groupId string, groupUpdate *GroupUpdate) UpdateGroupApiRequest {
	return UpdateGroupApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		groupUpdate: groupUpdate,
	}
}

// UpdateGroupExecute executes the request
//
//	@return Group
func (a *ProjectsApiService) UpdateGroupExecute(r UpdateGroupApiRequest) (*Group, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Group
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.UpdateGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupUpdate == nil {
		return localVarReturnValue, nil, reportError("groupUpdate is required and must be specified")
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
	localVarPostBody = r.groupUpdate
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

type UpdateGroupInvitesApiRequest struct {
	ctx                    context.Context
	ApiService             ProjectsApi
	groupId                string
	groupInvitationRequest *GroupInvitationRequest
}

type UpdateGroupInvitesApiParams struct {
	GroupId                string
	GroupInvitationRequest *GroupInvitationRequest
}

func (a *ProjectsApiService) UpdateGroupInvitesWithParams(ctx context.Context, args *UpdateGroupInvitesApiParams) UpdateGroupInvitesApiRequest {
	return UpdateGroupInvitesApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		groupInvitationRequest: args.GroupInvitationRequest,
	}
}

func (r UpdateGroupInvitesApiRequest) Execute() (*GroupInvitation, *http.Response, error) {
	return r.ApiService.UpdateGroupInvitesExecute(r)
}

/*
UpdateGroupInvites Update One Invitation in One Project

Updates the details of one pending invitation to the specified project. To specify which invitation to update, provide the username of the invited user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateGroupInvitesApiRequest

Deprecated
*/
func (a *ProjectsApiService) UpdateGroupInvites(ctx context.Context, groupId string, groupInvitationRequest *GroupInvitationRequest) UpdateGroupInvitesApiRequest {
	return UpdateGroupInvitesApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		groupInvitationRequest: groupInvitationRequest,
	}
}

// UpdateGroupInvitesExecute executes the request
//
//	@return GroupInvitation
//
// Deprecated
func (a *ProjectsApiService) UpdateGroupInvitesExecute(r UpdateGroupInvitesApiRequest) (*GroupInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.UpdateGroupInvites")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupInvitationRequest == nil {
		return localVarReturnValue, nil, reportError("groupInvitationRequest is required and must be specified")
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
	localVarPostBody = r.groupInvitationRequest
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

type UpdateGroupSettingsApiRequest struct {
	ctx           context.Context
	ApiService    ProjectsApi
	groupId       string
	groupSettings *GroupSettings
}

type UpdateGroupSettingsApiParams struct {
	GroupId       string
	GroupSettings *GroupSettings
}

func (a *ProjectsApiService) UpdateGroupSettingsWithParams(ctx context.Context, args *UpdateGroupSettingsApiParams) UpdateGroupSettingsApiRequest {
	return UpdateGroupSettingsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		groupSettings: args.GroupSettings,
	}
}

func (r UpdateGroupSettingsApiRequest) Execute() (*GroupSettings, *http.Response, error) {
	return r.ApiService.UpdateGroupSettingsExecute(r)
}

/*
UpdateGroupSettings Update Project Settings

Updates the settings of the specified project. You can update any of the options available. MongoDB cloud only updates the options provided in the request.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateGroupSettingsApiRequest
*/
func (a *ProjectsApiService) UpdateGroupSettings(ctx context.Context, groupId string, groupSettings *GroupSettings) UpdateGroupSettingsApiRequest {
	return UpdateGroupSettingsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		groupSettings: groupSettings,
	}
}

// UpdateGroupSettingsExecute executes the request
//
//	@return GroupSettings
func (a *ProjectsApiService) UpdateGroupSettingsExecute(r UpdateGroupSettingsApiRequest) (*GroupSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.UpdateGroupSettings")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/settings"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupSettings == nil {
		return localVarReturnValue, nil, reportError("groupSettings is required and must be specified")
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
	localVarPostBody = r.groupSettings
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

type UpdateGroupUserRolesApiRequest struct {
	ctx                     context.Context
	ApiService              ProjectsApi
	groupId                 string
	userId                  string
	updateGroupRolesForUser *UpdateGroupRolesForUser
}

type UpdateGroupUserRolesApiParams struct {
	GroupId                 string
	UserId                  string
	UpdateGroupRolesForUser *UpdateGroupRolesForUser
}

func (a *ProjectsApiService) UpdateGroupUserRolesWithParams(ctx context.Context, args *UpdateGroupUserRolesApiParams) UpdateGroupUserRolesApiRequest {
	return UpdateGroupUserRolesApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 args.GroupId,
		userId:                  args.UserId,
		updateGroupRolesForUser: args.UpdateGroupRolesForUser,
	}
}

func (r UpdateGroupUserRolesApiRequest) Execute() (*UpdateGroupRolesForUser, *http.Response, error) {
	return r.ApiService.UpdateGroupUserRolesExecute(r)
}

/*
UpdateGroupUserRoles Update Project Roles for One MongoDB Cloud User

Updates the roles of the specified user in the specified project. To specify the user to update, provide the unique 24-hexadecimal digit string that identifies the user in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param userId Unique 24-hexadecimal digit string that identifies the user to modify.
	@return UpdateGroupUserRolesApiRequest

Deprecated
*/
func (a *ProjectsApiService) UpdateGroupUserRoles(ctx context.Context, groupId string, userId string, updateGroupRolesForUser *UpdateGroupRolesForUser) UpdateGroupUserRolesApiRequest {
	return UpdateGroupUserRolesApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 groupId,
		userId:                  userId,
		updateGroupRolesForUser: updateGroupRolesForUser,
	}
}

// UpdateGroupUserRolesExecute executes the request
//
//	@return UpdateGroupRolesForUser
//
// Deprecated
func (a *ProjectsApiService) UpdateGroupUserRolesExecute(r UpdateGroupUserRolesApiRequest) (*UpdateGroupRolesForUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UpdateGroupRolesForUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.UpdateGroupUserRoles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/users/{userId}/roles"
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
	if r.updateGroupRolesForUser == nil {
		return localVarReturnValue, nil, reportError("updateGroupRolesForUser is required and must be specified")
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
	localVarPostBody = r.updateGroupRolesForUser
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

type UpdateInviteByIdApiRequest struct {
	ctx                          context.Context
	ApiService                   ProjectsApi
	groupId                      string
	invitationId                 string
	groupInvitationUpdateRequest *GroupInvitationUpdateRequest
}

type UpdateInviteByIdApiParams struct {
	GroupId                      string
	InvitationId                 string
	GroupInvitationUpdateRequest *GroupInvitationUpdateRequest
}

func (a *ProjectsApiService) UpdateInviteByIdWithParams(ctx context.Context, args *UpdateInviteByIdApiParams) UpdateInviteByIdApiRequest {
	return UpdateInviteByIdApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		invitationId:                 args.InvitationId,
		groupInvitationUpdateRequest: args.GroupInvitationUpdateRequest,
	}
}

func (r UpdateInviteByIdApiRequest) Execute() (*GroupInvitation, *http.Response, error) {
	return r.ApiService.UpdateInviteByIdExecute(r)
}

/*
UpdateInviteById Update One Invitation in One Project by Invitation ID

Updates the details of one pending invitation to the specified project. To specify which invitation to update, provide the unique identification string for that invitation. Use the Return All Project Invitations endpoint to retrieve IDs for all pending project invitations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return UpdateInviteByIdApiRequest

Deprecated
*/
func (a *ProjectsApiService) UpdateInviteById(ctx context.Context, groupId string, invitationId string, groupInvitationUpdateRequest *GroupInvitationUpdateRequest) UpdateInviteByIdApiRequest {
	return UpdateInviteByIdApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		invitationId:                 invitationId,
		groupInvitationUpdateRequest: groupInvitationUpdateRequest,
	}
}

// UpdateInviteByIdExecute executes the request
//
//	@return GroupInvitation
//
// Deprecated
func (a *ProjectsApiService) UpdateInviteByIdExecute(r UpdateInviteByIdApiRequest) (*GroupInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectsApiService.UpdateInviteById")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/invites/{invitationId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.invitationId == "" {
		return localVarReturnValue, nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupInvitationUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("groupInvitationUpdateRequest is required and must be specified")
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
	localVarPostBody = r.groupInvitationUpdateRequest
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
