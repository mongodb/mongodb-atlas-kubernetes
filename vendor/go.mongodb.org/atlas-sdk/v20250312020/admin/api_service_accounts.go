// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ServiceAccountsApi interface {

	/*
		CreateAccessList Add Access List Entries for One Project Service Account

		Add Access List Entries for the specified Service Account for the project. Resources require all API requests to originate from IP addresses on the API access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clientId The Client ID of the Service Account.
		@param serviceAccountIPAccessListEntry A list of access list entries to add to the access list of the specified Service Account for the project.
		@return CreateAccessListApiRequest
	*/
	CreateAccessList(ctx context.Context, groupId string, clientId string, serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry) CreateAccessListApiRequest
	/*
		CreateAccessList Add Access List Entries for One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateAccessListApiParams - Parameters for the request
		@return CreateAccessListApiRequest
	*/
	CreateAccessListWithParams(ctx context.Context, args *CreateAccessListApiParams) CreateAccessListApiRequest

	// Method available only for mocking purposes
	CreateAccessListExecute(r CreateAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error)

	/*
		CreateGroupSecret Create One Project Service Account Secret

		Create a secret for the specified Service Account in the specified Project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clientId The Client ID of the Service Account.
		@param serviceAccountSecretRequest Details for the new secret.
		@return CreateGroupSecretApiRequest
	*/
	CreateGroupSecret(ctx context.Context, groupId string, clientId string, serviceAccountSecretRequest *ServiceAccountSecretRequest) CreateGroupSecretApiRequest
	/*
		CreateGroupSecret Create One Project Service Account Secret


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupSecretApiParams - Parameters for the request
		@return CreateGroupSecretApiRequest
	*/
	CreateGroupSecretWithParams(ctx context.Context, args *CreateGroupSecretApiParams) CreateGroupSecretApiRequest

	// Method available only for mocking purposes
	CreateGroupSecretExecute(r CreateGroupSecretApiRequest) (*ServiceAccountSecret, *http.Response, error)

	/*
		CreateGroupServiceAccount Create One Project Service Account

		Creates one Service Account for the specified Project. The Service Account will automatically be added as an Organization Member to the Organization that the specified Project is a part of.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupServiceAccountRequest Details of the new Service Account.
		@return CreateGroupServiceAccountApiRequest
	*/
	CreateGroupServiceAccount(ctx context.Context, groupId string, groupServiceAccountRequest *GroupServiceAccountRequest) CreateGroupServiceAccountApiRequest
	/*
		CreateGroupServiceAccount Create One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupServiceAccountApiParams - Parameters for the request
		@return CreateGroupServiceAccountApiRequest
	*/
	CreateGroupServiceAccountWithParams(ctx context.Context, args *CreateGroupServiceAccountApiParams) CreateGroupServiceAccountApiRequest

	// Method available only for mocking purposes
	CreateGroupServiceAccountExecute(r CreateGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error)

	/*
		CreateOrgAccessList Add Access List Entries for One Organization Service Account

		Add Access List Entries for the specified Service Account for the organization. Resources require all API requests to originate from IP addresses on the API access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@param serviceAccountIPAccessListEntry A list of access list entries to add to the access list of the specified Service Account for the organization.
		@return CreateOrgAccessListApiRequest
	*/
	CreateOrgAccessList(ctx context.Context, orgId string, clientId string, serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry) CreateOrgAccessListApiRequest
	/*
		CreateOrgAccessList Add Access List Entries for One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgAccessListApiParams - Parameters for the request
		@return CreateOrgAccessListApiRequest
	*/
	CreateOrgAccessListWithParams(ctx context.Context, args *CreateOrgAccessListApiParams) CreateOrgAccessListApiRequest

	// Method available only for mocking purposes
	CreateOrgAccessListExecute(r CreateOrgAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error)

	/*
		CreateOrgSecret Create One Organization Service Account Secret

		Create a secret for the specified Service Account.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@param serviceAccountSecretRequest Details for the new secret.
		@return CreateOrgSecretApiRequest
	*/
	CreateOrgSecret(ctx context.Context, orgId string, clientId string, serviceAccountSecretRequest *ServiceAccountSecretRequest) CreateOrgSecretApiRequest
	/*
		CreateOrgSecret Create One Organization Service Account Secret


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgSecretApiParams - Parameters for the request
		@return CreateOrgSecretApiRequest
	*/
	CreateOrgSecretWithParams(ctx context.Context, args *CreateOrgSecretApiParams) CreateOrgSecretApiRequest

	// Method available only for mocking purposes
	CreateOrgSecretExecute(r CreateOrgSecretApiRequest) (*ServiceAccountSecret, *http.Response, error)

	/*
		CreateOrgServiceAccount Create One Organization Service Account

		Creates one Service Account for the specified Organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param orgServiceAccountRequest Details of the new Service Account.
		@return CreateOrgServiceAccountApiRequest
	*/
	CreateOrgServiceAccount(ctx context.Context, orgId string, orgServiceAccountRequest *OrgServiceAccountRequest) CreateOrgServiceAccountApiRequest
	/*
		CreateOrgServiceAccount Create One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgServiceAccountApiParams - Parameters for the request
		@return CreateOrgServiceAccountApiRequest
	*/
	CreateOrgServiceAccountWithParams(ctx context.Context, args *CreateOrgServiceAccountApiParams) CreateOrgServiceAccountApiRequest

	// Method available only for mocking purposes
	CreateOrgServiceAccountExecute(r CreateOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error)

	/*
		DeleteGroupAccessEntry Remove One Access List Entry from One Project Service Account

		Removes the specified access list entry from the specified Service Account for the project. You can't remove the requesting IP address from the access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clientId The Client ID of the Service Account.
		@param ipAddress One IP address or multiple IP addresses represented as one CIDR block. When specifying a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
		@return DeleteGroupAccessEntryApiRequest
	*/
	DeleteGroupAccessEntry(ctx context.Context, groupId string, clientId string, ipAddress string) DeleteGroupAccessEntryApiRequest
	/*
		DeleteGroupAccessEntry Remove One Access List Entry from One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupAccessEntryApiParams - Parameters for the request
		@return DeleteGroupAccessEntryApiRequest
	*/
	DeleteGroupAccessEntryWithParams(ctx context.Context, args *DeleteGroupAccessEntryApiParams) DeleteGroupAccessEntryApiRequest

	// Method available only for mocking purposes
	DeleteGroupAccessEntryExecute(r DeleteGroupAccessEntryApiRequest) (*http.Response, error)

	/*
		DeleteGroupSecret Delete One Project Service Account Secret

		Deletes the specified Service Account secret.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param secretId Unique 24-hexadecimal digit string that identifies the secret.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteGroupSecretApiRequest
	*/
	DeleteGroupSecret(ctx context.Context, clientId string, secretId string, groupId string) DeleteGroupSecretApiRequest
	/*
		DeleteGroupSecret Delete One Project Service Account Secret


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupSecretApiParams - Parameters for the request
		@return DeleteGroupSecretApiRequest
	*/
	DeleteGroupSecretWithParams(ctx context.Context, args *DeleteGroupSecretApiParams) DeleteGroupSecretApiRequest

	// Method available only for mocking purposes
	DeleteGroupSecretExecute(r DeleteGroupSecretApiRequest) (*http.Response, error)

	/*
		DeleteGroupServiceAccount Remove One Project Service Account

		Removes the specified Service Account from the specified project. The Service Account will still be a part of the Organization it was created in, and the credentials will remain active until expired or manually revoked.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteGroupServiceAccountApiRequest
	*/
	DeleteGroupServiceAccount(ctx context.Context, clientId string, groupId string) DeleteGroupServiceAccountApiRequest
	/*
		DeleteGroupServiceAccount Remove One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupServiceAccountApiParams - Parameters for the request
		@return DeleteGroupServiceAccountApiRequest
	*/
	DeleteGroupServiceAccountWithParams(ctx context.Context, args *DeleteGroupServiceAccountApiParams) DeleteGroupServiceAccountApiRequest

	// Method available only for mocking purposes
	DeleteGroupServiceAccountExecute(r DeleteGroupServiceAccountApiRequest) (*http.Response, error)

	/*
		DeleteOrgAccessEntry Remove One Access List Entry from One Organization Service Account

		Removes the specified access list entry from the specified Service Account for the organization. You can't remove the requesting IP address from the access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@param ipAddress One IP address or multiple IP addresses represented as one CIDR block. When specifying a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
		@return DeleteOrgAccessEntryApiRequest
	*/
	DeleteOrgAccessEntry(ctx context.Context, orgId string, clientId string, ipAddress string) DeleteOrgAccessEntryApiRequest
	/*
		DeleteOrgAccessEntry Remove One Access List Entry from One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgAccessEntryApiParams - Parameters for the request
		@return DeleteOrgAccessEntryApiRequest
	*/
	DeleteOrgAccessEntryWithParams(ctx context.Context, args *DeleteOrgAccessEntryApiParams) DeleteOrgAccessEntryApiRequest

	// Method available only for mocking purposes
	DeleteOrgAccessEntryExecute(r DeleteOrgAccessEntryApiRequest) (*http.Response, error)

	/*
		DeleteOrgSecret Delete One Organization Service Account Secret

		Deletes the specified Service Account secret.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param secretId Unique 24-hexadecimal digit string that identifies the secret.
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return DeleteOrgSecretApiRequest
	*/
	DeleteOrgSecret(ctx context.Context, clientId string, secretId string, orgId string) DeleteOrgSecretApiRequest
	/*
		DeleteOrgSecret Delete One Organization Service Account Secret


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgSecretApiParams - Parameters for the request
		@return DeleteOrgSecretApiRequest
	*/
	DeleteOrgSecretWithParams(ctx context.Context, args *DeleteOrgSecretApiParams) DeleteOrgSecretApiRequest

	// Method available only for mocking purposes
	DeleteOrgSecretExecute(r DeleteOrgSecretApiRequest) (*http.Response, error)

	/*
		DeleteOrgServiceAccount Delete One Organization Service Account

		Deletes the specified Service Account.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return DeleteOrgServiceAccountApiRequest
	*/
	DeleteOrgServiceAccount(ctx context.Context, clientId string, orgId string) DeleteOrgServiceAccountApiRequest
	/*
		DeleteOrgServiceAccount Delete One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgServiceAccountApiParams - Parameters for the request
		@return DeleteOrgServiceAccountApiRequest
	*/
	DeleteOrgServiceAccountWithParams(ctx context.Context, args *DeleteOrgServiceAccountApiParams) DeleteOrgServiceAccountApiRequest

	// Method available only for mocking purposes
	DeleteOrgServiceAccountExecute(r DeleteOrgServiceAccountApiRequest) (*http.Response, error)

	/*
		GetGroupServiceAccount Return One Project Service Account

		Returns one Service Account in the specified Project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clientId The Client ID of the Service Account.
		@return GetGroupServiceAccountApiRequest
	*/
	GetGroupServiceAccount(ctx context.Context, groupId string, clientId string) GetGroupServiceAccountApiRequest
	/*
		GetGroupServiceAccount Return One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupServiceAccountApiParams - Parameters for the request
		@return GetGroupServiceAccountApiRequest
	*/
	GetGroupServiceAccountWithParams(ctx context.Context, args *GetGroupServiceAccountApiParams) GetGroupServiceAccountApiRequest

	// Method available only for mocking purposes
	GetGroupServiceAccountExecute(r GetGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error)

	/*
		GetOrgServiceAccount Return One Organization Service Account

		Returns the specified Service Account.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@return GetOrgServiceAccountApiRequest
	*/
	GetOrgServiceAccount(ctx context.Context, orgId string, clientId string) GetOrgServiceAccountApiRequest
	/*
		GetOrgServiceAccount Return One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgServiceAccountApiParams - Parameters for the request
		@return GetOrgServiceAccountApiRequest
	*/
	GetOrgServiceAccountWithParams(ctx context.Context, args *GetOrgServiceAccountApiParams) GetOrgServiceAccountApiRequest

	// Method available only for mocking purposes
	GetOrgServiceAccountExecute(r GetOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error)

	/*
		GetServiceAccountGroups Return All Service Account Project Assignments

		Returns a list of all projects the specified Service Account is a part of.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@return GetServiceAccountGroupsApiRequest
	*/
	GetServiceAccountGroups(ctx context.Context, orgId string, clientId string) GetServiceAccountGroupsApiRequest
	/*
		GetServiceAccountGroups Return All Service Account Project Assignments


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServiceAccountGroupsApiParams - Parameters for the request
		@return GetServiceAccountGroupsApiRequest
	*/
	GetServiceAccountGroupsWithParams(ctx context.Context, args *GetServiceAccountGroupsApiParams) GetServiceAccountGroupsApiRequest

	// Method available only for mocking purposes
	GetServiceAccountGroupsExecute(r GetServiceAccountGroupsApiRequest) (*PaginatedServiceAccountGroup, *http.Response, error)

	/*
		InviteGroupServiceAccount Assign One Service Account to One Project

		Assigns the specified Service Account to the specified Project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupServiceAccountRoleAssignment The Project permissions for the Service Account in the specified Project.
		@return InviteGroupServiceAccountApiRequest
	*/
	InviteGroupServiceAccount(ctx context.Context, clientId string, groupId string, groupServiceAccountRoleAssignment *GroupServiceAccountRoleAssignment) InviteGroupServiceAccountApiRequest
	/*
		InviteGroupServiceAccount Assign One Service Account to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param InviteGroupServiceAccountApiParams - Parameters for the request
		@return InviteGroupServiceAccountApiRequest
	*/
	InviteGroupServiceAccountWithParams(ctx context.Context, args *InviteGroupServiceAccountApiParams) InviteGroupServiceAccountApiRequest

	// Method available only for mocking purposes
	InviteGroupServiceAccountExecute(r InviteGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error)

	/*
		ListAccessList Return All Access List Entries for One Project Service Account

		Returns all access list entries that you configured for the specified Service Account for the project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clientId The Client ID of the Service Account.
		@return ListAccessListApiRequest
	*/
	ListAccessList(ctx context.Context, groupId string, clientId string) ListAccessListApiRequest
	/*
		ListAccessList Return All Access List Entries for One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListAccessListApiParams - Parameters for the request
		@return ListAccessListApiRequest
	*/
	ListAccessListWithParams(ctx context.Context, args *ListAccessListApiParams) ListAccessListApiRequest

	// Method available only for mocking purposes
	ListAccessListExecute(r ListAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error)

	/*
		ListGroupServiceAccounts Return All Project Service Accounts

		Returns all Service Accounts for the specified Project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupServiceAccountsApiRequest
	*/
	ListGroupServiceAccounts(ctx context.Context, groupId string) ListGroupServiceAccountsApiRequest
	/*
		ListGroupServiceAccounts Return All Project Service Accounts


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupServiceAccountsApiParams - Parameters for the request
		@return ListGroupServiceAccountsApiRequest
	*/
	ListGroupServiceAccountsWithParams(ctx context.Context, args *ListGroupServiceAccountsApiParams) ListGroupServiceAccountsApiRequest

	// Method available only for mocking purposes
	ListGroupServiceAccountsExecute(r ListGroupServiceAccountsApiRequest) (*PaginatedGroupServiceAccounts, *http.Response, error)

	/*
		ListOrgAccessList Return All Access List Entries for One Organization Service Account

		Returns all access list entries that you configured for the specified Service Account for the organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param clientId The Client ID of the Service Account.
		@return ListOrgAccessListApiRequest
	*/
	ListOrgAccessList(ctx context.Context, orgId string, clientId string) ListOrgAccessListApiRequest
	/*
		ListOrgAccessList Return All Access List Entries for One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgAccessListApiParams - Parameters for the request
		@return ListOrgAccessListApiRequest
	*/
	ListOrgAccessListWithParams(ctx context.Context, args *ListOrgAccessListApiParams) ListOrgAccessListApiRequest

	// Method available only for mocking purposes
	ListOrgAccessListExecute(r ListOrgAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error)

	/*
		ListOrgServiceAccounts Return All Organization Service Accounts

		Returns all Service Accounts for the specified Organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListOrgServiceAccountsApiRequest
	*/
	ListOrgServiceAccounts(ctx context.Context, orgId string) ListOrgServiceAccountsApiRequest
	/*
		ListOrgServiceAccounts Return All Organization Service Accounts


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgServiceAccountsApiParams - Parameters for the request
		@return ListOrgServiceAccountsApiRequest
	*/
	ListOrgServiceAccountsWithParams(ctx context.Context, args *ListOrgServiceAccountsApiParams) ListOrgServiceAccountsApiRequest

	// Method available only for mocking purposes
	ListOrgServiceAccountsExecute(r ListOrgServiceAccountsApiRequest) (*PaginatedOrgServiceAccounts, *http.Response, error)

	/*
		UpdateGroupServiceAccount Update One Project Service Account

		Updates one Service Account in the specified Project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupServiceAccountUpdateRequest The new details for the Service Account.
		@return UpdateGroupServiceAccountApiRequest
	*/
	UpdateGroupServiceAccount(ctx context.Context, clientId string, groupId string, groupServiceAccountUpdateRequest *GroupServiceAccountUpdateRequest) UpdateGroupServiceAccountApiRequest
	/*
		UpdateGroupServiceAccount Update One Project Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupServiceAccountApiParams - Parameters for the request
		@return UpdateGroupServiceAccountApiRequest
	*/
	UpdateGroupServiceAccountWithParams(ctx context.Context, args *UpdateGroupServiceAccountApiParams) UpdateGroupServiceAccountApiRequest

	// Method available only for mocking purposes
	UpdateGroupServiceAccountExecute(r UpdateGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error)

	/*
		UpdateOrgServiceAccount Update One Organization Service Account

		Updates the specified Service Account in the specified Organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param clientId The Client ID of the Service Account.
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param orgServiceAccountUpdateRequest The new details for the Service Account.
		@return UpdateOrgServiceAccountApiRequest
	*/
	UpdateOrgServiceAccount(ctx context.Context, clientId string, orgId string, orgServiceAccountUpdateRequest *OrgServiceAccountUpdateRequest) UpdateOrgServiceAccountApiRequest
	/*
		UpdateOrgServiceAccount Update One Organization Service Account


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgServiceAccountApiParams - Parameters for the request
		@return UpdateOrgServiceAccountApiRequest
	*/
	UpdateOrgServiceAccountWithParams(ctx context.Context, args *UpdateOrgServiceAccountApiParams) UpdateOrgServiceAccountApiRequest

	// Method available only for mocking purposes
	UpdateOrgServiceAccountExecute(r UpdateOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error)
}

// ServiceAccountsApiService ServiceAccountsApi service
type ServiceAccountsApiService service

type CreateAccessListApiRequest struct {
	ctx                             context.Context
	ApiService                      ServiceAccountsApi
	groupId                         string
	clientId                        string
	serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry
	includeCount                    *bool
	itemsPerPage                    *int
	pageNum                         *int
}

type CreateAccessListApiParams struct {
	GroupId                         string
	ClientId                        string
	ServiceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry
	IncludeCount                    *bool
	ItemsPerPage                    *int
	PageNum                         *int
}

func (a *ServiceAccountsApiService) CreateAccessListWithParams(ctx context.Context, args *CreateAccessListApiParams) CreateAccessListApiRequest {
	return CreateAccessListApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         args.GroupId,
		clientId:                        args.ClientId,
		serviceAccountIPAccessListEntry: args.ServiceAccountIPAccessListEntry,
		includeCount:                    args.IncludeCount,
		itemsPerPage:                    args.ItemsPerPage,
		pageNum:                         args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r CreateAccessListApiRequest) IncludeCount(includeCount bool) CreateAccessListApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r CreateAccessListApiRequest) ItemsPerPage(itemsPerPage int) CreateAccessListApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r CreateAccessListApiRequest) PageNum(pageNum int) CreateAccessListApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r CreateAccessListApiRequest) Execute() (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	return r.ApiService.CreateAccessListExecute(r)
}

/*
CreateAccessList Add Access List Entries for One Project Service Account

Add Access List Entries for the specified Service Account for the project. Resources require all API requests to originate from IP addresses on the API access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clientId The Client ID of the Service Account.
	@return CreateAccessListApiRequest
*/
func (a *ServiceAccountsApiService) CreateAccessList(ctx context.Context, groupId string, clientId string, serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry) CreateAccessListApiRequest {
	return CreateAccessListApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         groupId,
		clientId:                        clientId,
		serviceAccountIPAccessListEntry: serviceAccountIPAccessListEntry,
	}
}

// CreateAccessListExecute executes the request
//
//	@return PaginatedServiceAccountIPAccessEntry
func (a *ServiceAccountsApiService) CreateAccessListExecute(r CreateAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServiceAccountIPAccessEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateAccessList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}/accessList"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serviceAccountIPAccessListEntry == nil {
		return localVarReturnValue, nil, reportError("serviceAccountIPAccessListEntry is required and must be specified")
	}

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
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.serviceAccountIPAccessListEntry
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

type CreateGroupSecretApiRequest struct {
	ctx                         context.Context
	ApiService                  ServiceAccountsApi
	groupId                     string
	clientId                    string
	serviceAccountSecretRequest *ServiceAccountSecretRequest
}

type CreateGroupSecretApiParams struct {
	GroupId                     string
	ClientId                    string
	ServiceAccountSecretRequest *ServiceAccountSecretRequest
}

func (a *ServiceAccountsApiService) CreateGroupSecretWithParams(ctx context.Context, args *CreateGroupSecretApiParams) CreateGroupSecretApiRequest {
	return CreateGroupSecretApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		groupId:                     args.GroupId,
		clientId:                    args.ClientId,
		serviceAccountSecretRequest: args.ServiceAccountSecretRequest,
	}
}

func (r CreateGroupSecretApiRequest) Execute() (*ServiceAccountSecret, *http.Response, error) {
	return r.ApiService.CreateGroupSecretExecute(r)
}

/*
CreateGroupSecret Create One Project Service Account Secret

Create a secret for the specified Service Account in the specified Project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clientId The Client ID of the Service Account.
	@return CreateGroupSecretApiRequest
*/
func (a *ServiceAccountsApiService) CreateGroupSecret(ctx context.Context, groupId string, clientId string, serviceAccountSecretRequest *ServiceAccountSecretRequest) CreateGroupSecretApiRequest {
	return CreateGroupSecretApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		groupId:                     groupId,
		clientId:                    clientId,
		serviceAccountSecretRequest: serviceAccountSecretRequest,
	}
}

// CreateGroupSecretExecute executes the request
//
//	@return ServiceAccountSecret
func (a *ServiceAccountsApiService) CreateGroupSecretExecute(r CreateGroupSecretApiRequest) (*ServiceAccountSecret, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServiceAccountSecret
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateGroupSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}/secrets"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serviceAccountSecretRequest == nil {
		return localVarReturnValue, nil, reportError("serviceAccountSecretRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.serviceAccountSecretRequest
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

type CreateGroupServiceAccountApiRequest struct {
	ctx                        context.Context
	ApiService                 ServiceAccountsApi
	groupId                    string
	groupServiceAccountRequest *GroupServiceAccountRequest
}

type CreateGroupServiceAccountApiParams struct {
	GroupId                    string
	GroupServiceAccountRequest *GroupServiceAccountRequest
}

func (a *ServiceAccountsApiService) CreateGroupServiceAccountWithParams(ctx context.Context, args *CreateGroupServiceAccountApiParams) CreateGroupServiceAccountApiRequest {
	return CreateGroupServiceAccountApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		groupServiceAccountRequest: args.GroupServiceAccountRequest,
	}
}

func (r CreateGroupServiceAccountApiRequest) Execute() (*GroupServiceAccount, *http.Response, error) {
	return r.ApiService.CreateGroupServiceAccountExecute(r)
}

/*
CreateGroupServiceAccount Create One Project Service Account

Creates one Service Account for the specified Project. The Service Account will automatically be added as an Organization Member to the Organization that the specified Project is a part of.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) CreateGroupServiceAccount(ctx context.Context, groupId string, groupServiceAccountRequest *GroupServiceAccountRequest) CreateGroupServiceAccountApiRequest {
	return CreateGroupServiceAccountApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		groupServiceAccountRequest: groupServiceAccountRequest,
	}
}

// CreateGroupServiceAccountExecute executes the request
//
//	@return GroupServiceAccount
func (a *ServiceAccountsApiService) CreateGroupServiceAccountExecute(r CreateGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateGroupServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupServiceAccountRequest == nil {
		return localVarReturnValue, nil, reportError("groupServiceAccountRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupServiceAccountRequest
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

type CreateOrgAccessListApiRequest struct {
	ctx                             context.Context
	ApiService                      ServiceAccountsApi
	orgId                           string
	clientId                        string
	serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry
	includeCount                    *bool
	itemsPerPage                    *int
	pageNum                         *int
}

type CreateOrgAccessListApiParams struct {
	OrgId                           string
	ClientId                        string
	ServiceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry
	IncludeCount                    *bool
	ItemsPerPage                    *int
	PageNum                         *int
}

func (a *ServiceAccountsApiService) CreateOrgAccessListWithParams(ctx context.Context, args *CreateOrgAccessListApiParams) CreateOrgAccessListApiRequest {
	return CreateOrgAccessListApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		orgId:                           args.OrgId,
		clientId:                        args.ClientId,
		serviceAccountIPAccessListEntry: args.ServiceAccountIPAccessListEntry,
		includeCount:                    args.IncludeCount,
		itemsPerPage:                    args.ItemsPerPage,
		pageNum:                         args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r CreateOrgAccessListApiRequest) IncludeCount(includeCount bool) CreateOrgAccessListApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r CreateOrgAccessListApiRequest) ItemsPerPage(itemsPerPage int) CreateOrgAccessListApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r CreateOrgAccessListApiRequest) PageNum(pageNum int) CreateOrgAccessListApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r CreateOrgAccessListApiRequest) Execute() (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	return r.ApiService.CreateOrgAccessListExecute(r)
}

/*
CreateOrgAccessList Add Access List Entries for One Organization Service Account

Add Access List Entries for the specified Service Account for the organization. Resources require all API requests to originate from IP addresses on the API access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@return CreateOrgAccessListApiRequest
*/
func (a *ServiceAccountsApiService) CreateOrgAccessList(ctx context.Context, orgId string, clientId string, serviceAccountIPAccessListEntry *[]ServiceAccountIPAccessListEntry) CreateOrgAccessListApiRequest {
	return CreateOrgAccessListApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		orgId:                           orgId,
		clientId:                        clientId,
		serviceAccountIPAccessListEntry: serviceAccountIPAccessListEntry,
	}
}

// CreateOrgAccessListExecute executes the request
//
//	@return PaginatedServiceAccountIPAccessEntry
func (a *ServiceAccountsApiService) CreateOrgAccessListExecute(r CreateOrgAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServiceAccountIPAccessEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateOrgAccessList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/accessList"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serviceAccountIPAccessListEntry == nil {
		return localVarReturnValue, nil, reportError("serviceAccountIPAccessListEntry is required and must be specified")
	}

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
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.serviceAccountIPAccessListEntry
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

type CreateOrgSecretApiRequest struct {
	ctx                         context.Context
	ApiService                  ServiceAccountsApi
	orgId                       string
	clientId                    string
	serviceAccountSecretRequest *ServiceAccountSecretRequest
}

type CreateOrgSecretApiParams struct {
	OrgId                       string
	ClientId                    string
	ServiceAccountSecretRequest *ServiceAccountSecretRequest
}

func (a *ServiceAccountsApiService) CreateOrgSecretWithParams(ctx context.Context, args *CreateOrgSecretApiParams) CreateOrgSecretApiRequest {
	return CreateOrgSecretApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		orgId:                       args.OrgId,
		clientId:                    args.ClientId,
		serviceAccountSecretRequest: args.ServiceAccountSecretRequest,
	}
}

func (r CreateOrgSecretApiRequest) Execute() (*ServiceAccountSecret, *http.Response, error) {
	return r.ApiService.CreateOrgSecretExecute(r)
}

/*
CreateOrgSecret Create One Organization Service Account Secret

Create a secret for the specified Service Account.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@return CreateOrgSecretApiRequest
*/
func (a *ServiceAccountsApiService) CreateOrgSecret(ctx context.Context, orgId string, clientId string, serviceAccountSecretRequest *ServiceAccountSecretRequest) CreateOrgSecretApiRequest {
	return CreateOrgSecretApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		orgId:                       orgId,
		clientId:                    clientId,
		serviceAccountSecretRequest: serviceAccountSecretRequest,
	}
}

// CreateOrgSecretExecute executes the request
//
//	@return ServiceAccountSecret
func (a *ServiceAccountsApiService) CreateOrgSecretExecute(r CreateOrgSecretApiRequest) (*ServiceAccountSecret, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServiceAccountSecret
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateOrgSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/secrets"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serviceAccountSecretRequest == nil {
		return localVarReturnValue, nil, reportError("serviceAccountSecretRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.serviceAccountSecretRequest
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

type CreateOrgServiceAccountApiRequest struct {
	ctx                      context.Context
	ApiService               ServiceAccountsApi
	orgId                    string
	orgServiceAccountRequest *OrgServiceAccountRequest
}

type CreateOrgServiceAccountApiParams struct {
	OrgId                    string
	OrgServiceAccountRequest *OrgServiceAccountRequest
}

func (a *ServiceAccountsApiService) CreateOrgServiceAccountWithParams(ctx context.Context, args *CreateOrgServiceAccountApiParams) CreateOrgServiceAccountApiRequest {
	return CreateOrgServiceAccountApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		orgId:                    args.OrgId,
		orgServiceAccountRequest: args.OrgServiceAccountRequest,
	}
}

func (r CreateOrgServiceAccountApiRequest) Execute() (*OrgServiceAccount, *http.Response, error) {
	return r.ApiService.CreateOrgServiceAccountExecute(r)
}

/*
CreateOrgServiceAccount Create One Organization Service Account

Creates one Service Account for the specified Organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) CreateOrgServiceAccount(ctx context.Context, orgId string, orgServiceAccountRequest *OrgServiceAccountRequest) CreateOrgServiceAccountApiRequest {
	return CreateOrgServiceAccountApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		orgId:                    orgId,
		orgServiceAccountRequest: orgServiceAccountRequest,
	}
}

// CreateOrgServiceAccountExecute executes the request
//
//	@return OrgServiceAccount
func (a *ServiceAccountsApiService) CreateOrgServiceAccountExecute(r CreateOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.CreateOrgServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.orgServiceAccountRequest == nil {
		return localVarReturnValue, nil, reportError("orgServiceAccountRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.orgServiceAccountRequest
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

type DeleteGroupAccessEntryApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	groupId    string
	clientId   string
	ipAddress  string
}

type DeleteGroupAccessEntryApiParams struct {
	GroupId   string
	ClientId  string
	IpAddress string
}

func (a *ServiceAccountsApiService) DeleteGroupAccessEntryWithParams(ctx context.Context, args *DeleteGroupAccessEntryApiParams) DeleteGroupAccessEntryApiRequest {
	return DeleteGroupAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		clientId:   args.ClientId,
		ipAddress:  args.IpAddress,
	}
}

func (r DeleteGroupAccessEntryApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupAccessEntryExecute(r)
}

/*
DeleteGroupAccessEntry Remove One Access List Entry from One Project Service Account

Removes the specified access list entry from the specified Service Account for the project. You can't remove the requesting IP address from the access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clientId The Client ID of the Service Account.
	@param ipAddress One IP address or multiple IP addresses represented as one CIDR block. When specifying a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
	@return DeleteGroupAccessEntryApiRequest
*/
func (a *ServiceAccountsApiService) DeleteGroupAccessEntry(ctx context.Context, groupId string, clientId string, ipAddress string) DeleteGroupAccessEntryApiRequest {
	return DeleteGroupAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		clientId:   clientId,
		ipAddress:  ipAddress,
	}
}

// DeleteGroupAccessEntryExecute executes the request
func (a *ServiceAccountsApiService) DeleteGroupAccessEntryExecute(r DeleteGroupAccessEntryApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteGroupAccessEntry")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}/accessList/{ipAddress}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.ipAddress == "" {
		return nil, reportError("ipAddress is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"ipAddress"+"}", url.PathEscape(r.ipAddress), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteGroupSecretApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	clientId   string
	secretId   string
	groupId    string
}

type DeleteGroupSecretApiParams struct {
	ClientId string
	SecretId string
	GroupId  string
}

func (a *ServiceAccountsApiService) DeleteGroupSecretWithParams(ctx context.Context, args *DeleteGroupSecretApiParams) DeleteGroupSecretApiRequest {
	return DeleteGroupSecretApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   args.ClientId,
		secretId:   args.SecretId,
		groupId:    args.GroupId,
	}
}

func (r DeleteGroupSecretApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupSecretExecute(r)
}

/*
DeleteGroupSecret Delete One Project Service Account Secret

Deletes the specified Service Account secret.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param secretId Unique 24-hexadecimal digit string that identifies the secret.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteGroupSecretApiRequest
*/
func (a *ServiceAccountsApiService) DeleteGroupSecret(ctx context.Context, clientId string, secretId string, groupId string) DeleteGroupSecretApiRequest {
	return DeleteGroupSecretApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   clientId,
		secretId:   secretId,
		groupId:    groupId,
	}
}

// DeleteGroupSecretExecute executes the request
func (a *ServiceAccountsApiService) DeleteGroupSecretExecute(r DeleteGroupSecretApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteGroupSecret")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}/secrets/{secretId}"
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.secretId == "" {
		return nil, reportError("secretId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"secretId"+"}", url.PathEscape(r.secretId), -1)
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteGroupServiceAccountApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	clientId   string
	groupId    string
}

type DeleteGroupServiceAccountApiParams struct {
	ClientId string
	GroupId  string
}

func (a *ServiceAccountsApiService) DeleteGroupServiceAccountWithParams(ctx context.Context, args *DeleteGroupServiceAccountApiParams) DeleteGroupServiceAccountApiRequest {
	return DeleteGroupServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   args.ClientId,
		groupId:    args.GroupId,
	}
}

func (r DeleteGroupServiceAccountApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupServiceAccountExecute(r)
}

/*
DeleteGroupServiceAccount Remove One Project Service Account

Removes the specified Service Account from the specified project. The Service Account will still be a part of the Organization it was created in, and the credentials will remain active until expired or manually revoked.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteGroupServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) DeleteGroupServiceAccount(ctx context.Context, clientId string, groupId string) DeleteGroupServiceAccountApiRequest {
	return DeleteGroupServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   clientId,
		groupId:    groupId,
	}
}

// DeleteGroupServiceAccountExecute executes the request
func (a *ServiceAccountsApiService) DeleteGroupServiceAccountExecute(r DeleteGroupServiceAccountApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteGroupServiceAccount")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}"
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteOrgAccessEntryApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	orgId      string
	clientId   string
	ipAddress  string
}

type DeleteOrgAccessEntryApiParams struct {
	OrgId     string
	ClientId  string
	IpAddress string
}

func (a *ServiceAccountsApiService) DeleteOrgAccessEntryWithParams(ctx context.Context, args *DeleteOrgAccessEntryApiParams) DeleteOrgAccessEntryApiRequest {
	return DeleteOrgAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		clientId:   args.ClientId,
		ipAddress:  args.IpAddress,
	}
}

func (r DeleteOrgAccessEntryApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgAccessEntryExecute(r)
}

/*
DeleteOrgAccessEntry Remove One Access List Entry from One Organization Service Account

Removes the specified access list entry from the specified Service Account for the organization. You can't remove the requesting IP address from the access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@param ipAddress One IP address or multiple IP addresses represented as one CIDR block. When specifying a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
	@return DeleteOrgAccessEntryApiRequest
*/
func (a *ServiceAccountsApiService) DeleteOrgAccessEntry(ctx context.Context, orgId string, clientId string, ipAddress string) DeleteOrgAccessEntryApiRequest {
	return DeleteOrgAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		clientId:   clientId,
		ipAddress:  ipAddress,
	}
}

// DeleteOrgAccessEntryExecute executes the request
func (a *ServiceAccountsApiService) DeleteOrgAccessEntryExecute(r DeleteOrgAccessEntryApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteOrgAccessEntry")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/accessList/{ipAddress}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.ipAddress == "" {
		return nil, reportError("ipAddress is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"ipAddress"+"}", url.PathEscape(r.ipAddress), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteOrgSecretApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	clientId   string
	secretId   string
	orgId      string
}

type DeleteOrgSecretApiParams struct {
	ClientId string
	SecretId string
	OrgId    string
}

func (a *ServiceAccountsApiService) DeleteOrgSecretWithParams(ctx context.Context, args *DeleteOrgSecretApiParams) DeleteOrgSecretApiRequest {
	return DeleteOrgSecretApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   args.ClientId,
		secretId:   args.SecretId,
		orgId:      args.OrgId,
	}
}

func (r DeleteOrgSecretApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgSecretExecute(r)
}

/*
DeleteOrgSecret Delete One Organization Service Account Secret

Deletes the specified Service Account secret.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param secretId Unique 24-hexadecimal digit string that identifies the secret.
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return DeleteOrgSecretApiRequest
*/
func (a *ServiceAccountsApiService) DeleteOrgSecret(ctx context.Context, clientId string, secretId string, orgId string) DeleteOrgSecretApiRequest {
	return DeleteOrgSecretApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   clientId,
		secretId:   secretId,
		orgId:      orgId,
	}
}

// DeleteOrgSecretExecute executes the request
func (a *ServiceAccountsApiService) DeleteOrgSecretExecute(r DeleteOrgSecretApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteOrgSecret")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/secrets/{secretId}"
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.secretId == "" {
		return nil, reportError("secretId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"secretId"+"}", url.PathEscape(r.secretId), -1)
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteOrgServiceAccountApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	clientId   string
	orgId      string
}

type DeleteOrgServiceAccountApiParams struct {
	ClientId string
	OrgId    string
}

func (a *ServiceAccountsApiService) DeleteOrgServiceAccountWithParams(ctx context.Context, args *DeleteOrgServiceAccountApiParams) DeleteOrgServiceAccountApiRequest {
	return DeleteOrgServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   args.ClientId,
		orgId:      args.OrgId,
	}
}

func (r DeleteOrgServiceAccountApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgServiceAccountExecute(r)
}

/*
DeleteOrgServiceAccount Delete One Organization Service Account

Deletes the specified Service Account.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return DeleteOrgServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) DeleteOrgServiceAccount(ctx context.Context, clientId string, orgId string) DeleteOrgServiceAccountApiRequest {
	return DeleteOrgServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		clientId:   clientId,
		orgId:      orgId,
	}
}

// DeleteOrgServiceAccountExecute executes the request
func (a *ServiceAccountsApiService) DeleteOrgServiceAccountExecute(r DeleteOrgServiceAccountApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.DeleteOrgServiceAccount")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}"
	if r.clientId == "" {
		return nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetGroupServiceAccountApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	groupId    string
	clientId   string
}

type GetGroupServiceAccountApiParams struct {
	GroupId  string
	ClientId string
}

func (a *ServiceAccountsApiService) GetGroupServiceAccountWithParams(ctx context.Context, args *GetGroupServiceAccountApiParams) GetGroupServiceAccountApiRequest {
	return GetGroupServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		clientId:   args.ClientId,
	}
}

func (r GetGroupServiceAccountApiRequest) Execute() (*GroupServiceAccount, *http.Response, error) {
	return r.ApiService.GetGroupServiceAccountExecute(r)
}

/*
GetGroupServiceAccount Return One Project Service Account

Returns one Service Account in the specified Project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clientId The Client ID of the Service Account.
	@return GetGroupServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) GetGroupServiceAccount(ctx context.Context, groupId string, clientId string) GetGroupServiceAccountApiRequest {
	return GetGroupServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		clientId:   clientId,
	}
}

// GetGroupServiceAccountExecute executes the request
//
//	@return GroupServiceAccount
func (a *ServiceAccountsApiService) GetGroupServiceAccountExecute(r GetGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.GetGroupServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetOrgServiceAccountApiRequest struct {
	ctx        context.Context
	ApiService ServiceAccountsApi
	orgId      string
	clientId   string
}

type GetOrgServiceAccountApiParams struct {
	OrgId    string
	ClientId string
}

func (a *ServiceAccountsApiService) GetOrgServiceAccountWithParams(ctx context.Context, args *GetOrgServiceAccountApiParams) GetOrgServiceAccountApiRequest {
	return GetOrgServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		clientId:   args.ClientId,
	}
}

func (r GetOrgServiceAccountApiRequest) Execute() (*OrgServiceAccount, *http.Response, error) {
	return r.ApiService.GetOrgServiceAccountExecute(r)
}

/*
GetOrgServiceAccount Return One Organization Service Account

Returns the specified Service Account.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@return GetOrgServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) GetOrgServiceAccount(ctx context.Context, orgId string, clientId string) GetOrgServiceAccountApiRequest {
	return GetOrgServiceAccountApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		clientId:   clientId,
	}
}

// GetOrgServiceAccountExecute executes the request
//
//	@return OrgServiceAccount
func (a *ServiceAccountsApiService) GetOrgServiceAccountExecute(r GetOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.GetOrgServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetServiceAccountGroupsApiRequest struct {
	ctx          context.Context
	ApiService   ServiceAccountsApi
	orgId        string
	clientId     string
	itemsPerPage *int
	pageNum      *int
}

type GetServiceAccountGroupsApiParams struct {
	OrgId        string
	ClientId     string
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServiceAccountsApiService) GetServiceAccountGroupsWithParams(ctx context.Context, args *GetServiceAccountGroupsApiParams) GetServiceAccountGroupsApiRequest {
	return GetServiceAccountGroupsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		clientId:     args.ClientId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r GetServiceAccountGroupsApiRequest) ItemsPerPage(itemsPerPage int) GetServiceAccountGroupsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetServiceAccountGroupsApiRequest) PageNum(pageNum int) GetServiceAccountGroupsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r GetServiceAccountGroupsApiRequest) Execute() (*PaginatedServiceAccountGroup, *http.Response, error) {
	return r.ApiService.GetServiceAccountGroupsExecute(r)
}

/*
GetServiceAccountGroups Return All Service Account Project Assignments

Returns a list of all projects the specified Service Account is a part of.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@return GetServiceAccountGroupsApiRequest
*/
func (a *ServiceAccountsApiService) GetServiceAccountGroups(ctx context.Context, orgId string, clientId string) GetServiceAccountGroupsApiRequest {
	return GetServiceAccountGroupsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		clientId:   clientId,
	}
}

// GetServiceAccountGroupsExecute executes the request
//
//	@return PaginatedServiceAccountGroup
func (a *ServiceAccountsApiService) GetServiceAccountGroupsExecute(r GetServiceAccountGroupsApiRequest) (*PaginatedServiceAccountGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServiceAccountGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.GetServiceAccountGroups")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/groups"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

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
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type InviteGroupServiceAccountApiRequest struct {
	ctx                               context.Context
	ApiService                        ServiceAccountsApi
	clientId                          string
	groupId                           string
	groupServiceAccountRoleAssignment *GroupServiceAccountRoleAssignment
}

type InviteGroupServiceAccountApiParams struct {
	ClientId                          string
	GroupId                           string
	GroupServiceAccountRoleAssignment *GroupServiceAccountRoleAssignment
}

func (a *ServiceAccountsApiService) InviteGroupServiceAccountWithParams(ctx context.Context, args *InviteGroupServiceAccountApiParams) InviteGroupServiceAccountApiRequest {
	return InviteGroupServiceAccountApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		clientId:                          args.ClientId,
		groupId:                           args.GroupId,
		groupServiceAccountRoleAssignment: args.GroupServiceAccountRoleAssignment,
	}
}

func (r InviteGroupServiceAccountApiRequest) Execute() (*GroupServiceAccount, *http.Response, error) {
	return r.ApiService.InviteGroupServiceAccountExecute(r)
}

/*
InviteGroupServiceAccount Assign One Service Account to One Project

Assigns the specified Service Account to the specified Project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return InviteGroupServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) InviteGroupServiceAccount(ctx context.Context, clientId string, groupId string, groupServiceAccountRoleAssignment *GroupServiceAccountRoleAssignment) InviteGroupServiceAccountApiRequest {
	return InviteGroupServiceAccountApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		clientId:                          clientId,
		groupId:                           groupId,
		groupServiceAccountRoleAssignment: groupServiceAccountRoleAssignment,
	}
}

// InviteGroupServiceAccountExecute executes the request
//
//	@return GroupServiceAccount
func (a *ServiceAccountsApiService) InviteGroupServiceAccountExecute(r InviteGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.InviteGroupServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}:invite"
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupServiceAccountRoleAssignment == nil {
		return localVarReturnValue, nil, reportError("groupServiceAccountRoleAssignment is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupServiceAccountRoleAssignment
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

type ListAccessListApiRequest struct {
	ctx          context.Context
	ApiService   ServiceAccountsApi
	groupId      string
	clientId     string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListAccessListApiParams struct {
	GroupId      string
	ClientId     string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServiceAccountsApiService) ListAccessListWithParams(ctx context.Context, args *ListAccessListApiParams) ListAccessListApiRequest {
	return ListAccessListApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clientId:     args.ClientId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListAccessListApiRequest) IncludeCount(includeCount bool) ListAccessListApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListAccessListApiRequest) ItemsPerPage(itemsPerPage int) ListAccessListApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListAccessListApiRequest) PageNum(pageNum int) ListAccessListApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListAccessListApiRequest) Execute() (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	return r.ApiService.ListAccessListExecute(r)
}

/*
ListAccessList Return All Access List Entries for One Project Service Account

Returns all access list entries that you configured for the specified Service Account for the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clientId The Client ID of the Service Account.
	@return ListAccessListApiRequest
*/
func (a *ServiceAccountsApiService) ListAccessList(ctx context.Context, groupId string, clientId string) ListAccessListApiRequest {
	return ListAccessListApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		clientId:   clientId,
	}
}

// ListAccessListExecute executes the request
//
//	@return PaginatedServiceAccountIPAccessEntry
func (a *ServiceAccountsApiService) ListAccessListExecute(r ListAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServiceAccountIPAccessEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.ListAccessList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}/accessList"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type ListGroupServiceAccountsApiRequest struct {
	ctx          context.Context
	ApiService   ServiceAccountsApi
	groupId      string
	itemsPerPage *int
	pageNum      *int
}

type ListGroupServiceAccountsApiParams struct {
	GroupId      string
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServiceAccountsApiService) ListGroupServiceAccountsWithParams(ctx context.Context, args *ListGroupServiceAccountsApiParams) ListGroupServiceAccountsApiRequest {
	return ListGroupServiceAccountsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListGroupServiceAccountsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupServiceAccountsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupServiceAccountsApiRequest) PageNum(pageNum int) ListGroupServiceAccountsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupServiceAccountsApiRequest) Execute() (*PaginatedGroupServiceAccounts, *http.Response, error) {
	return r.ApiService.ListGroupServiceAccountsExecute(r)
}

/*
ListGroupServiceAccounts Return All Project Service Accounts

Returns all Service Accounts for the specified Project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupServiceAccountsApiRequest
*/
func (a *ServiceAccountsApiService) ListGroupServiceAccounts(ctx context.Context, groupId string) ListGroupServiceAccountsApiRequest {
	return ListGroupServiceAccountsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupServiceAccountsExecute executes the request
//
//	@return PaginatedGroupServiceAccounts
func (a *ServiceAccountsApiService) ListGroupServiceAccountsExecute(r ListGroupServiceAccountsApiRequest) (*PaginatedGroupServiceAccounts, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedGroupServiceAccounts
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.ListGroupServiceAccounts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type ListOrgAccessListApiRequest struct {
	ctx          context.Context
	ApiService   ServiceAccountsApi
	orgId        string
	clientId     string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListOrgAccessListApiParams struct {
	OrgId        string
	ClientId     string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServiceAccountsApiService) ListOrgAccessListWithParams(ctx context.Context, args *ListOrgAccessListApiParams) ListOrgAccessListApiRequest {
	return ListOrgAccessListApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		clientId:     args.ClientId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgAccessListApiRequest) IncludeCount(includeCount bool) ListOrgAccessListApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgAccessListApiRequest) ItemsPerPage(itemsPerPage int) ListOrgAccessListApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgAccessListApiRequest) PageNum(pageNum int) ListOrgAccessListApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOrgAccessListApiRequest) Execute() (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	return r.ApiService.ListOrgAccessListExecute(r)
}

/*
ListOrgAccessList Return All Access List Entries for One Organization Service Account

Returns all access list entries that you configured for the specified Service Account for the organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param clientId The Client ID of the Service Account.
	@return ListOrgAccessListApiRequest
*/
func (a *ServiceAccountsApiService) ListOrgAccessList(ctx context.Context, orgId string, clientId string) ListOrgAccessListApiRequest {
	return ListOrgAccessListApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		clientId:   clientId,
	}
}

// ListOrgAccessListExecute executes the request
//
//	@return PaginatedServiceAccountIPAccessEntry
func (a *ServiceAccountsApiService) ListOrgAccessListExecute(r ListOrgAccessListApiRequest) (*PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServiceAccountIPAccessEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.ListOrgAccessList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}/accessList"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type ListOrgServiceAccountsApiRequest struct {
	ctx          context.Context
	ApiService   ServiceAccountsApi
	orgId        string
	itemsPerPage *int
	pageNum      *int
}

type ListOrgServiceAccountsApiParams struct {
	OrgId        string
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServiceAccountsApiService) ListOrgServiceAccountsWithParams(ctx context.Context, args *ListOrgServiceAccountsApiParams) ListOrgServiceAccountsApiRequest {
	return ListOrgServiceAccountsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListOrgServiceAccountsApiRequest) ItemsPerPage(itemsPerPage int) ListOrgServiceAccountsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgServiceAccountsApiRequest) PageNum(pageNum int) ListOrgServiceAccountsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOrgServiceAccountsApiRequest) Execute() (*PaginatedOrgServiceAccounts, *http.Response, error) {
	return r.ApiService.ListOrgServiceAccountsExecute(r)
}

/*
ListOrgServiceAccounts Return All Organization Service Accounts

Returns all Service Accounts for the specified Organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgServiceAccountsApiRequest
*/
func (a *ServiceAccountsApiService) ListOrgServiceAccounts(ctx context.Context, orgId string) ListOrgServiceAccountsApiRequest {
	return ListOrgServiceAccountsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgServiceAccountsExecute executes the request
//
//	@return PaginatedOrgServiceAccounts
func (a *ServiceAccountsApiService) ListOrgServiceAccountsExecute(r ListOrgServiceAccountsApiRequest) (*PaginatedOrgServiceAccounts, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOrgServiceAccounts
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.ListOrgServiceAccounts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type UpdateGroupServiceAccountApiRequest struct {
	ctx                              context.Context
	ApiService                       ServiceAccountsApi
	clientId                         string
	groupId                          string
	groupServiceAccountUpdateRequest *GroupServiceAccountUpdateRequest
}

type UpdateGroupServiceAccountApiParams struct {
	ClientId                         string
	GroupId                          string
	GroupServiceAccountUpdateRequest *GroupServiceAccountUpdateRequest
}

func (a *ServiceAccountsApiService) UpdateGroupServiceAccountWithParams(ctx context.Context, args *UpdateGroupServiceAccountApiParams) UpdateGroupServiceAccountApiRequest {
	return UpdateGroupServiceAccountApiRequest{
		ApiService:                       a,
		ctx:                              ctx,
		clientId:                         args.ClientId,
		groupId:                          args.GroupId,
		groupServiceAccountUpdateRequest: args.GroupServiceAccountUpdateRequest,
	}
}

func (r UpdateGroupServiceAccountApiRequest) Execute() (*GroupServiceAccount, *http.Response, error) {
	return r.ApiService.UpdateGroupServiceAccountExecute(r)
}

/*
UpdateGroupServiceAccount Update One Project Service Account

Updates one Service Account in the specified Project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateGroupServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) UpdateGroupServiceAccount(ctx context.Context, clientId string, groupId string, groupServiceAccountUpdateRequest *GroupServiceAccountUpdateRequest) UpdateGroupServiceAccountApiRequest {
	return UpdateGroupServiceAccountApiRequest{
		ApiService:                       a,
		ctx:                              ctx,
		clientId:                         clientId,
		groupId:                          groupId,
		groupServiceAccountUpdateRequest: groupServiceAccountUpdateRequest,
	}
}

// UpdateGroupServiceAccountExecute executes the request
//
//	@return GroupServiceAccount
func (a *ServiceAccountsApiService) UpdateGroupServiceAccountExecute(r UpdateGroupServiceAccountApiRequest) (*GroupServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.UpdateGroupServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serviceAccounts/{clientId}"
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupServiceAccountUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("groupServiceAccountUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.groupServiceAccountUpdateRequest
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

type UpdateOrgServiceAccountApiRequest struct {
	ctx                            context.Context
	ApiService                     ServiceAccountsApi
	clientId                       string
	orgId                          string
	orgServiceAccountUpdateRequest *OrgServiceAccountUpdateRequest
}

type UpdateOrgServiceAccountApiParams struct {
	ClientId                       string
	OrgId                          string
	OrgServiceAccountUpdateRequest *OrgServiceAccountUpdateRequest
}

func (a *ServiceAccountsApiService) UpdateOrgServiceAccountWithParams(ctx context.Context, args *UpdateOrgServiceAccountApiParams) UpdateOrgServiceAccountApiRequest {
	return UpdateOrgServiceAccountApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		clientId:                       args.ClientId,
		orgId:                          args.OrgId,
		orgServiceAccountUpdateRequest: args.OrgServiceAccountUpdateRequest,
	}
}

func (r UpdateOrgServiceAccountApiRequest) Execute() (*OrgServiceAccount, *http.Response, error) {
	return r.ApiService.UpdateOrgServiceAccountExecute(r)
}

/*
UpdateOrgServiceAccount Update One Organization Service Account

Updates the specified Service Account in the specified Organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId The Client ID of the Service Account.
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return UpdateOrgServiceAccountApiRequest
*/
func (a *ServiceAccountsApiService) UpdateOrgServiceAccount(ctx context.Context, clientId string, orgId string, orgServiceAccountUpdateRequest *OrgServiceAccountUpdateRequest) UpdateOrgServiceAccountApiRequest {
	return UpdateOrgServiceAccountApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		clientId:                       clientId,
		orgId:                          orgId,
		orgServiceAccountUpdateRequest: orgServiceAccountUpdateRequest,
	}
}

// UpdateOrgServiceAccountExecute executes the request
//
//	@return OrgServiceAccount
func (a *ServiceAccountsApiService) UpdateOrgServiceAccountExecute(r UpdateOrgServiceAccountApiRequest) (*OrgServiceAccount, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgServiceAccount
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceAccountsApiService.UpdateOrgServiceAccount")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}"
	if r.clientId == "" {
		return localVarReturnValue, nil, reportError("clientId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(r.clientId), -1)
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.orgServiceAccountUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("orgServiceAccountUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.orgServiceAccountUpdateRequest
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
