// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ProgrammaticAPIKeysApi interface {

	/*
		AddGroupApiKey Assign One Organization API Key to One Project

		Assigns the specified organization API key to the specified project. Users with the Project Owner role in the project associated with the API key can then use the organization API key to access the resources.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to assign to one project.
		@param userAccessRoleAssignment Organization API key to be assigned to the specified project.
		@return AddGroupApiKeyApiRequest
	*/
	AddGroupApiKey(ctx context.Context, groupId string, apiUserId string, userAccessRoleAssignment *[]UserAccessRoleAssignment) AddGroupApiKeyApiRequest
	/*
		AddGroupApiKey Assign One Organization API Key to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AddGroupApiKeyApiParams - Parameters for the request
		@return AddGroupApiKeyApiRequest
	*/
	AddGroupApiKeyWithParams(ctx context.Context, args *AddGroupApiKeyApiParams) AddGroupApiKeyApiRequest

	// Method available only for mocking purposes
	AddGroupApiKeyExecute(r AddGroupApiKeyApiRequest) (*http.Response, error)

	/*
		CreateGroupApiKey Create and Assign One Organization API Key to One Project

		Creates and assigns the specified organization API key to the specified project. Users with the Project Owner role in the project associated with the API key can use the organization API key to access the resources.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param createAtlasProjectApiKey Organization API key to be created and assigned to the specified project.
		@return CreateGroupApiKeyApiRequest
	*/
	CreateGroupApiKey(ctx context.Context, groupId string, createAtlasProjectApiKey *CreateAtlasProjectApiKey) CreateGroupApiKeyApiRequest
	/*
		CreateGroupApiKey Create and Assign One Organization API Key to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupApiKeyApiParams - Parameters for the request
		@return CreateGroupApiKeyApiRequest
	*/
	CreateGroupApiKeyWithParams(ctx context.Context, args *CreateGroupApiKeyApiParams) CreateGroupApiKeyApiRequest

	// Method available only for mocking purposes
	CreateGroupApiKeyExecute(r CreateGroupApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error)

	/*
		CreateOrgAccessEntry Create One Access List Entry for One Organization API Key

		Creates the access list entries for the specified organization API key. Resources require all API requests originate from IP addresses on the API access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to create a new access list entry.
		@param userAccessListRequest Access list entries to be created for the specified organization API key.
		@return CreateOrgAccessEntryApiRequest
	*/
	CreateOrgAccessEntry(ctx context.Context, orgId string, apiUserId string, userAccessListRequest *[]UserAccessListRequest) CreateOrgAccessEntryApiRequest
	/*
		CreateOrgAccessEntry Create One Access List Entry for One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgAccessEntryApiParams - Parameters for the request
		@return CreateOrgAccessEntryApiRequest
	*/
	CreateOrgAccessEntryWithParams(ctx context.Context, args *CreateOrgAccessEntryApiParams) CreateOrgAccessEntryApiRequest

	// Method available only for mocking purposes
	CreateOrgAccessEntryExecute(r CreateOrgAccessEntryApiRequest) (*PaginatedApiUserAccessListResponse, *http.Response, error)

	/*
		CreateOrgApiKey Create One Organization API Key

		Creates one API key for the specified organization. An organization API key grants programmatic access to an organization. You can't use the API key to log into the console.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param createAtlasOrganizationApiKey Organization API Key to be created.
		@return CreateOrgApiKeyApiRequest
	*/
	CreateOrgApiKey(ctx context.Context, orgId string, createAtlasOrganizationApiKey *CreateAtlasOrganizationApiKey) CreateOrgApiKeyApiRequest
	/*
		CreateOrgApiKey Create One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgApiKeyApiParams - Parameters for the request
		@return CreateOrgApiKeyApiRequest
	*/
	CreateOrgApiKeyWithParams(ctx context.Context, args *CreateOrgApiKeyApiParams) CreateOrgApiKeyApiRequest

	// Method available only for mocking purposes
	CreateOrgApiKeyExecute(r CreateOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error)

	/*
		DeleteAccessEntry Remove One Access List Entry for One Organization API Key

		Removes the specified access list entry from the specified organization API key. Resources require all API requests originate from the IP addresses on the API access list. In addition, you cannot remove the requesting IP address from the requesting organization API key.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to remove access list entries.
		@param ipAddress One IP address or multiple IP addresses represented as one CIDR block to limit requests to API resources in the specified organization. When adding a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
		@return DeleteAccessEntryApiRequest
	*/
	DeleteAccessEntry(ctx context.Context, orgId string, apiUserId string, ipAddress string) DeleteAccessEntryApiRequest
	/*
		DeleteAccessEntry Remove One Access List Entry for One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteAccessEntryApiParams - Parameters for the request
		@return DeleteAccessEntryApiRequest
	*/
	DeleteAccessEntryWithParams(ctx context.Context, args *DeleteAccessEntryApiParams) DeleteAccessEntryApiRequest

	// Method available only for mocking purposes
	DeleteAccessEntryExecute(r DeleteAccessEntryApiRequest) (*http.Response, error)

	/*
		DeleteOrgApiKey Remove One Organization API Key

		Removes one organization API key from the specified organization. When you remove an API key from an organization, MongoDB Cloud also removes that key from any projects that use that key.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key.
		@return DeleteOrgApiKeyApiRequest
	*/
	DeleteOrgApiKey(ctx context.Context, orgId string, apiUserId string) DeleteOrgApiKeyApiRequest
	/*
		DeleteOrgApiKey Remove One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgApiKeyApiParams - Parameters for the request
		@return DeleteOrgApiKeyApiRequest
	*/
	DeleteOrgApiKeyWithParams(ctx context.Context, args *DeleteOrgApiKeyApiParams) DeleteOrgApiKeyApiRequest

	// Method available only for mocking purposes
	DeleteOrgApiKeyExecute(r DeleteOrgApiKeyApiRequest) (*http.Response, error)

	/*
		GetOrgAccessEntry Return One Access List Entry for One Organization API Key

		Returns one access list entry for the specified organization API key. Resources require  all API requests originate from IP addresses on the API access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param ipAddress One IP address or multiple IP addresses represented as one CIDR block to limit  requests to API resources in the specified organization. When adding a CIDR block with a subnet mask, such as  192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for  which you want to return access list entries.
		@return GetOrgAccessEntryApiRequest
	*/
	GetOrgAccessEntry(ctx context.Context, orgId string, ipAddress string, apiUserId string) GetOrgAccessEntryApiRequest
	/*
		GetOrgAccessEntry Return One Access List Entry for One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgAccessEntryApiParams - Parameters for the request
		@return GetOrgAccessEntryApiRequest
	*/
	GetOrgAccessEntryWithParams(ctx context.Context, args *GetOrgAccessEntryApiParams) GetOrgAccessEntryApiRequest

	// Method available only for mocking purposes
	GetOrgAccessEntryExecute(r GetOrgAccessEntryApiRequest) (*UserAccessListResponse, *http.Response, error)

	/*
		GetOrgApiKey Return One Organization API Key

		Returns one organization API key. The organization API keys grant programmatic access to an organization. You can't use the API key to log into MongoDB Cloud through the user interface.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that  you want to update.
		@return GetOrgApiKeyApiRequest
	*/
	GetOrgApiKey(ctx context.Context, orgId string, apiUserId string) GetOrgApiKeyApiRequest
	/*
		GetOrgApiKey Return One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgApiKeyApiParams - Parameters for the request
		@return GetOrgApiKeyApiRequest
	*/
	GetOrgApiKeyWithParams(ctx context.Context, args *GetOrgApiKeyApiParams) GetOrgApiKeyApiRequest

	// Method available only for mocking purposes
	GetOrgApiKeyExecute(r GetOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error)

	/*
		ListGroupApiKeys Return All Organization API Keys Assigned to One Project

		Returns all organization API keys that you assigned to the specified project. Users with the Project Owner role in the project associated with the API key can use the organization API key to access the resources.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupApiKeysApiRequest
	*/
	ListGroupApiKeys(ctx context.Context, groupId string) ListGroupApiKeysApiRequest
	/*
		ListGroupApiKeys Return All Organization API Keys Assigned to One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupApiKeysApiParams - Parameters for the request
		@return ListGroupApiKeysApiRequest
	*/
	ListGroupApiKeysWithParams(ctx context.Context, args *ListGroupApiKeysApiParams) ListGroupApiKeysApiRequest

	// Method available only for mocking purposes
	ListGroupApiKeysExecute(r ListGroupApiKeysApiRequest) (*PaginatedApiApiUser, *http.Response, error)

	/*
		ListOrgAccessEntries Return All Access List Entries for One Organization API Key

		Returns all access list entries that you configured for the specified organization API key.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to return access list entries.
		@return ListOrgAccessEntriesApiRequest
	*/
	ListOrgAccessEntries(ctx context.Context, orgId string, apiUserId string) ListOrgAccessEntriesApiRequest
	/*
		ListOrgAccessEntries Return All Access List Entries for One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgAccessEntriesApiParams - Parameters for the request
		@return ListOrgAccessEntriesApiRequest
	*/
	ListOrgAccessEntriesWithParams(ctx context.Context, args *ListOrgAccessEntriesApiParams) ListOrgAccessEntriesApiRequest

	// Method available only for mocking purposes
	ListOrgAccessEntriesExecute(r ListOrgAccessEntriesApiRequest) (*PaginatedApiUserAccessListResponse, *http.Response, error)

	/*
		ListOrgApiKeys Return All Organization API Keys

		Returns all organization API keys for the specified organization. The organization API keys grant programmatic access to an organization. You can't use the API key to log into MongoDB Cloud through the console.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListOrgApiKeysApiRequest
	*/
	ListOrgApiKeys(ctx context.Context, orgId string) ListOrgApiKeysApiRequest
	/*
		ListOrgApiKeys Return All Organization API Keys


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgApiKeysApiParams - Parameters for the request
		@return ListOrgApiKeysApiRequest
	*/
	ListOrgApiKeysWithParams(ctx context.Context, args *ListOrgApiKeysApiParams) ListOrgApiKeysApiRequest

	// Method available only for mocking purposes
	ListOrgApiKeysExecute(r ListOrgApiKeysApiRequest) (*PaginatedApiApiUser, *http.Response, error)

	/*
		RemoveGroupApiKey Unassign One Organization API Key from One Project

		Removes one organization API key from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to unassign from one project.
		@return RemoveGroupApiKeyApiRequest
	*/
	RemoveGroupApiKey(ctx context.Context, groupId string, apiUserId string) RemoveGroupApiKeyApiRequest
	/*
		RemoveGroupApiKey Unassign One Organization API Key from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RemoveGroupApiKeyApiParams - Parameters for the request
		@return RemoveGroupApiKeyApiRequest
	*/
	RemoveGroupApiKeyWithParams(ctx context.Context, args *RemoveGroupApiKeyApiParams) RemoveGroupApiKeyApiRequest

	// Method available only for mocking purposes
	RemoveGroupApiKeyExecute(r RemoveGroupApiKeyApiRequest) (*http.Response, error)

	/*
		UpdateApiKeyRoles Update Organization API Key Roles for One Project

		Updates the roles of the organization API key that you specify for the project that you specify. You must specify at least one valid role for the project. The application removes any roles that you do not include in this request if they were previously set in the organization API key that you specify for the project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to unassign from one project.
		@param updateAtlasProjectApiKey Organization API Key to be updated. This request requires a minimum of one of the two body parameters.
		@return UpdateApiKeyRolesApiRequest
	*/
	UpdateApiKeyRoles(ctx context.Context, groupId string, apiUserId string, updateAtlasProjectApiKey *UpdateAtlasProjectApiKey) UpdateApiKeyRolesApiRequest
	/*
		UpdateApiKeyRoles Update Organization API Key Roles for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateApiKeyRolesApiParams - Parameters for the request
		@return UpdateApiKeyRolesApiRequest
	*/
	UpdateApiKeyRolesWithParams(ctx context.Context, args *UpdateApiKeyRolesApiParams) UpdateApiKeyRolesApiRequest

	// Method available only for mocking purposes
	UpdateApiKeyRolesExecute(r UpdateApiKeyRolesApiRequest) (*ApiKeyUserDetails, *http.Response, error)

	/*
		UpdateOrgApiKey Update One Organization API Key

		Updates one organization API key in the specified organization. The organization API keys  grant programmatic access to an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key you  want to update.
		@param updateAtlasOrganizationApiKey Organization API key to be updated. This request requires a minimum of one of the two body parameters.
		@return UpdateOrgApiKeyApiRequest
	*/
	UpdateOrgApiKey(ctx context.Context, orgId string, apiUserId string, updateAtlasOrganizationApiKey *UpdateAtlasOrganizationApiKey) UpdateOrgApiKeyApiRequest
	/*
		UpdateOrgApiKey Update One Organization API Key


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgApiKeyApiParams - Parameters for the request
		@return UpdateOrgApiKeyApiRequest
	*/
	UpdateOrgApiKeyWithParams(ctx context.Context, args *UpdateOrgApiKeyApiParams) UpdateOrgApiKeyApiRequest

	// Method available only for mocking purposes
	UpdateOrgApiKeyExecute(r UpdateOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error)
}

// ProgrammaticAPIKeysApiService ProgrammaticAPIKeysApi service
type ProgrammaticAPIKeysApiService service

type AddGroupApiKeyApiRequest struct {
	ctx                      context.Context
	ApiService               ProgrammaticAPIKeysApi
	groupId                  string
	apiUserId                string
	userAccessRoleAssignment *[]UserAccessRoleAssignment
}

type AddGroupApiKeyApiParams struct {
	GroupId                  string
	ApiUserId                string
	UserAccessRoleAssignment *[]UserAccessRoleAssignment
}

func (a *ProgrammaticAPIKeysApiService) AddGroupApiKeyWithParams(ctx context.Context, args *AddGroupApiKeyApiParams) AddGroupApiKeyApiRequest {
	return AddGroupApiKeyApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		apiUserId:                args.ApiUserId,
		userAccessRoleAssignment: args.UserAccessRoleAssignment,
	}
}

func (r AddGroupApiKeyApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.AddGroupApiKeyExecute(r)
}

/*
AddGroupApiKey Assign One Organization API Key to One Project

Assigns the specified organization API key to the specified project. Users with the Project Owner role in the project associated with the API key can then use the organization API key to access the resources.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to assign to one project.
	@return AddGroupApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) AddGroupApiKey(ctx context.Context, groupId string, apiUserId string, userAccessRoleAssignment *[]UserAccessRoleAssignment) AddGroupApiKeyApiRequest {
	return AddGroupApiKeyApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		apiUserId:                apiUserId,
		userAccessRoleAssignment: userAccessRoleAssignment,
	}
}

// AddGroupApiKeyExecute executes the request
func (a *ProgrammaticAPIKeysApiService) AddGroupApiKeyExecute(r AddGroupApiKeyApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.AddGroupApiKey")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/apiKeys/{apiUserId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.apiUserId == "" {
		return nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userAccessRoleAssignment == nil {
		return nil, reportError("userAccessRoleAssignment is required and must be specified")
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
	localVarPostBody = r.userAccessRoleAssignment
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

type CreateGroupApiKeyApiRequest struct {
	ctx                      context.Context
	ApiService               ProgrammaticAPIKeysApi
	groupId                  string
	createAtlasProjectApiKey *CreateAtlasProjectApiKey
}

type CreateGroupApiKeyApiParams struct {
	GroupId                  string
	CreateAtlasProjectApiKey *CreateAtlasProjectApiKey
}

func (a *ProgrammaticAPIKeysApiService) CreateGroupApiKeyWithParams(ctx context.Context, args *CreateGroupApiKeyApiParams) CreateGroupApiKeyApiRequest {
	return CreateGroupApiKeyApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		createAtlasProjectApiKey: args.CreateAtlasProjectApiKey,
	}
}

func (r CreateGroupApiKeyApiRequest) Execute() (*ApiKeyUserDetails, *http.Response, error) {
	return r.ApiService.CreateGroupApiKeyExecute(r)
}

/*
CreateGroupApiKey Create and Assign One Organization API Key to One Project

Creates and assigns the specified organization API key to the specified project. Users with the Project Owner role in the project associated with the API key can use the organization API key to access the resources.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) CreateGroupApiKey(ctx context.Context, groupId string, createAtlasProjectApiKey *CreateAtlasProjectApiKey) CreateGroupApiKeyApiRequest {
	return CreateGroupApiKeyApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		createAtlasProjectApiKey: createAtlasProjectApiKey,
	}
}

// CreateGroupApiKeyExecute executes the request
//
//	@return ApiKeyUserDetails
func (a *ProgrammaticAPIKeysApiService) CreateGroupApiKeyExecute(r CreateGroupApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiKeyUserDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.CreateGroupApiKey")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/apiKeys"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createAtlasProjectApiKey == nil {
		return localVarReturnValue, nil, reportError("createAtlasProjectApiKey is required and must be specified")
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
	localVarPostBody = r.createAtlasProjectApiKey
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

type CreateOrgAccessEntryApiRequest struct {
	ctx                   context.Context
	ApiService            ProgrammaticAPIKeysApi
	orgId                 string
	apiUserId             string
	userAccessListRequest *[]UserAccessListRequest
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
}

type CreateOrgAccessEntryApiParams struct {
	OrgId                 string
	ApiUserId             string
	UserAccessListRequest *[]UserAccessListRequest
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
}

func (a *ProgrammaticAPIKeysApiService) CreateOrgAccessEntryWithParams(ctx context.Context, args *CreateOrgAccessEntryApiParams) CreateOrgAccessEntryApiRequest {
	return CreateOrgAccessEntryApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 args.OrgId,
		apiUserId:             args.ApiUserId,
		userAccessListRequest: args.UserAccessListRequest,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r CreateOrgAccessEntryApiRequest) IncludeCount(includeCount bool) CreateOrgAccessEntryApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r CreateOrgAccessEntryApiRequest) ItemsPerPage(itemsPerPage int) CreateOrgAccessEntryApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r CreateOrgAccessEntryApiRequest) PageNum(pageNum int) CreateOrgAccessEntryApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r CreateOrgAccessEntryApiRequest) Execute() (*PaginatedApiUserAccessListResponse, *http.Response, error) {
	return r.ApiService.CreateOrgAccessEntryExecute(r)
}

/*
CreateOrgAccessEntry Create One Access List Entry for One Organization API Key

Creates the access list entries for the specified organization API key. Resources require all API requests originate from IP addresses on the API access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to create a new access list entry.
	@return CreateOrgAccessEntryApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) CreateOrgAccessEntry(ctx context.Context, orgId string, apiUserId string, userAccessListRequest *[]UserAccessListRequest) CreateOrgAccessEntryApiRequest {
	return CreateOrgAccessEntryApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 orgId,
		apiUserId:             apiUserId,
		userAccessListRequest: userAccessListRequest,
	}
}

// CreateOrgAccessEntryExecute executes the request
//
//	@return PaginatedApiUserAccessListResponse
func (a *ProgrammaticAPIKeysApiService) CreateOrgAccessEntryExecute(r CreateOrgAccessEntryApiRequest) (*PaginatedApiUserAccessListResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiUserAccessListResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.CreateOrgAccessEntry")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}/accessList"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userAccessListRequest == nil {
		return localVarReturnValue, nil, reportError("userAccessListRequest is required and must be specified")
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
	localVarPostBody = r.userAccessListRequest
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

type CreateOrgApiKeyApiRequest struct {
	ctx                           context.Context
	ApiService                    ProgrammaticAPIKeysApi
	orgId                         string
	createAtlasOrganizationApiKey *CreateAtlasOrganizationApiKey
}

type CreateOrgApiKeyApiParams struct {
	OrgId                         string
	CreateAtlasOrganizationApiKey *CreateAtlasOrganizationApiKey
}

func (a *ProgrammaticAPIKeysApiService) CreateOrgApiKeyWithParams(ctx context.Context, args *CreateOrgApiKeyApiParams) CreateOrgApiKeyApiRequest {
	return CreateOrgApiKeyApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         args.OrgId,
		createAtlasOrganizationApiKey: args.CreateAtlasOrganizationApiKey,
	}
}

func (r CreateOrgApiKeyApiRequest) Execute() (*ApiKeyUserDetails, *http.Response, error) {
	return r.ApiService.CreateOrgApiKeyExecute(r)
}

/*
CreateOrgApiKey Create One Organization API Key

Creates one API key for the specified organization. An organization API key grants programmatic access to an organization. You can't use the API key to log into the console.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) CreateOrgApiKey(ctx context.Context, orgId string, createAtlasOrganizationApiKey *CreateAtlasOrganizationApiKey) CreateOrgApiKeyApiRequest {
	return CreateOrgApiKeyApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         orgId,
		createAtlasOrganizationApiKey: createAtlasOrganizationApiKey,
	}
}

// CreateOrgApiKeyExecute executes the request
//
//	@return ApiKeyUserDetails
func (a *ProgrammaticAPIKeysApiService) CreateOrgApiKeyExecute(r CreateOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiKeyUserDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.CreateOrgApiKey")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createAtlasOrganizationApiKey == nil {
		return localVarReturnValue, nil, reportError("createAtlasOrganizationApiKey is required and must be specified")
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
	localVarPostBody = r.createAtlasOrganizationApiKey
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

type DeleteAccessEntryApiRequest struct {
	ctx        context.Context
	ApiService ProgrammaticAPIKeysApi
	orgId      string
	apiUserId  string
	ipAddress  string
}

type DeleteAccessEntryApiParams struct {
	OrgId     string
	ApiUserId string
	IpAddress string
}

func (a *ProgrammaticAPIKeysApiService) DeleteAccessEntryWithParams(ctx context.Context, args *DeleteAccessEntryApiParams) DeleteAccessEntryApiRequest {
	return DeleteAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		apiUserId:  args.ApiUserId,
		ipAddress:  args.IpAddress,
	}
}

func (r DeleteAccessEntryApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteAccessEntryExecute(r)
}

/*
DeleteAccessEntry Remove One Access List Entry for One Organization API Key

Removes the specified access list entry from the specified organization API key. Resources require all API requests originate from the IP addresses on the API access list. In addition, you cannot remove the requesting IP address from the requesting organization API key.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to remove access list entries.
	@param ipAddress One IP address or multiple IP addresses represented as one CIDR block to limit requests to API resources in the specified organization. When adding a CIDR block with a subnet mask, such as 192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
	@return DeleteAccessEntryApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) DeleteAccessEntry(ctx context.Context, orgId string, apiUserId string, ipAddress string) DeleteAccessEntryApiRequest {
	return DeleteAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		apiUserId:  apiUserId,
		ipAddress:  ipAddress,
	}
}

// DeleteAccessEntryExecute executes the request
func (a *ProgrammaticAPIKeysApiService) DeleteAccessEntryExecute(r DeleteAccessEntryApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.DeleteAccessEntry")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}/accessList/{ipAddress}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)
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

type DeleteOrgApiKeyApiRequest struct {
	ctx        context.Context
	ApiService ProgrammaticAPIKeysApi
	orgId      string
	apiUserId  string
}

type DeleteOrgApiKeyApiParams struct {
	OrgId     string
	ApiUserId string
}

func (a *ProgrammaticAPIKeysApiService) DeleteOrgApiKeyWithParams(ctx context.Context, args *DeleteOrgApiKeyApiParams) DeleteOrgApiKeyApiRequest {
	return DeleteOrgApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		apiUserId:  args.ApiUserId,
	}
}

func (r DeleteOrgApiKeyApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgApiKeyExecute(r)
}

/*
DeleteOrgApiKey Remove One Organization API Key

Removes one organization API key from the specified organization. When you remove an API key from an organization, MongoDB Cloud also removes that key from any projects that use that key.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key.
	@return DeleteOrgApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) DeleteOrgApiKey(ctx context.Context, orgId string, apiUserId string) DeleteOrgApiKeyApiRequest {
	return DeleteOrgApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		apiUserId:  apiUserId,
	}
}

// DeleteOrgApiKeyExecute executes the request
func (a *ProgrammaticAPIKeysApiService) DeleteOrgApiKeyExecute(r DeleteOrgApiKeyApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.DeleteOrgApiKey")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

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

type GetOrgAccessEntryApiRequest struct {
	ctx        context.Context
	ApiService ProgrammaticAPIKeysApi
	orgId      string
	ipAddress  string
	apiUserId  string
}

type GetOrgAccessEntryApiParams struct {
	OrgId     string
	IpAddress string
	ApiUserId string
}

func (a *ProgrammaticAPIKeysApiService) GetOrgAccessEntryWithParams(ctx context.Context, args *GetOrgAccessEntryApiParams) GetOrgAccessEntryApiRequest {
	return GetOrgAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		ipAddress:  args.IpAddress,
		apiUserId:  args.ApiUserId,
	}
}

func (r GetOrgAccessEntryApiRequest) Execute() (*UserAccessListResponse, *http.Response, error) {
	return r.ApiService.GetOrgAccessEntryExecute(r)
}

/*
GetOrgAccessEntry Return One Access List Entry for One Organization API Key

Returns one access list entry for the specified organization API key. Resources require  all API requests originate from IP addresses on the API access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param ipAddress One IP address or multiple IP addresses represented as one CIDR block to limit  requests to API resources in the specified organization. When adding a CIDR block with a subnet mask, such as  192.0.2.0/24, use the URL-encoded value %2F for the forward slash /.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for  which you want to return access list entries.
	@return GetOrgAccessEntryApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) GetOrgAccessEntry(ctx context.Context, orgId string, ipAddress string, apiUserId string) GetOrgAccessEntryApiRequest {
	return GetOrgAccessEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		ipAddress:  ipAddress,
		apiUserId:  apiUserId,
	}
}

// GetOrgAccessEntryExecute executes the request
//
//	@return UserAccessListResponse
func (a *ProgrammaticAPIKeysApiService) GetOrgAccessEntryExecute(r GetOrgAccessEntryApiRequest) (*UserAccessListResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserAccessListResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.GetOrgAccessEntry")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}/accessList/{ipAddress}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.ipAddress == "" {
		return localVarReturnValue, nil, reportError("ipAddress is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"ipAddress"+"}", url.PathEscape(r.ipAddress), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

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

type GetOrgApiKeyApiRequest struct {
	ctx        context.Context
	ApiService ProgrammaticAPIKeysApi
	orgId      string
	apiUserId  string
}

type GetOrgApiKeyApiParams struct {
	OrgId     string
	ApiUserId string
}

func (a *ProgrammaticAPIKeysApiService) GetOrgApiKeyWithParams(ctx context.Context, args *GetOrgApiKeyApiParams) GetOrgApiKeyApiRequest {
	return GetOrgApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		apiUserId:  args.ApiUserId,
	}
}

func (r GetOrgApiKeyApiRequest) Execute() (*ApiKeyUserDetails, *http.Response, error) {
	return r.ApiService.GetOrgApiKeyExecute(r)
}

/*
GetOrgApiKey Return One Organization API Key

Returns one organization API key. The organization API keys grant programmatic access to an organization. You can't use the API key to log into MongoDB Cloud through the user interface.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that  you want to update.
	@return GetOrgApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) GetOrgApiKey(ctx context.Context, orgId string, apiUserId string) GetOrgApiKeyApiRequest {
	return GetOrgApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		apiUserId:  apiUserId,
	}
}

// GetOrgApiKeyExecute executes the request
//
//	@return ApiKeyUserDetails
func (a *ProgrammaticAPIKeysApiService) GetOrgApiKeyExecute(r GetOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiKeyUserDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.GetOrgApiKey")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

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

type ListGroupApiKeysApiRequest struct {
	ctx          context.Context
	ApiService   ProgrammaticAPIKeysApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupApiKeysApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ProgrammaticAPIKeysApiService) ListGroupApiKeysWithParams(ctx context.Context, args *ListGroupApiKeysApiParams) ListGroupApiKeysApiRequest {
	return ListGroupApiKeysApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupApiKeysApiRequest) IncludeCount(includeCount bool) ListGroupApiKeysApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupApiKeysApiRequest) ItemsPerPage(itemsPerPage int) ListGroupApiKeysApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupApiKeysApiRequest) PageNum(pageNum int) ListGroupApiKeysApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupApiKeysApiRequest) Execute() (*PaginatedApiApiUser, *http.Response, error) {
	return r.ApiService.ListGroupApiKeysExecute(r)
}

/*
ListGroupApiKeys Return All Organization API Keys Assigned to One Project

Returns all organization API keys that you assigned to the specified project. Users with the Project Owner role in the project associated with the API key can use the organization API key to access the resources.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupApiKeysApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) ListGroupApiKeys(ctx context.Context, groupId string) ListGroupApiKeysApiRequest {
	return ListGroupApiKeysApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupApiKeysExecute executes the request
//
//	@return PaginatedApiApiUser
func (a *ProgrammaticAPIKeysApiService) ListGroupApiKeysExecute(r ListGroupApiKeysApiRequest) (*PaginatedApiApiUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiApiUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.ListGroupApiKeys")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/apiKeys"
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

type ListOrgAccessEntriesApiRequest struct {
	ctx          context.Context
	ApiService   ProgrammaticAPIKeysApi
	orgId        string
	apiUserId    string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListOrgAccessEntriesApiParams struct {
	OrgId        string
	ApiUserId    string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ProgrammaticAPIKeysApiService) ListOrgAccessEntriesWithParams(ctx context.Context, args *ListOrgAccessEntriesApiParams) ListOrgAccessEntriesApiRequest {
	return ListOrgAccessEntriesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		apiUserId:    args.ApiUserId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgAccessEntriesApiRequest) IncludeCount(includeCount bool) ListOrgAccessEntriesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgAccessEntriesApiRequest) ItemsPerPage(itemsPerPage int) ListOrgAccessEntriesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgAccessEntriesApiRequest) PageNum(pageNum int) ListOrgAccessEntriesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOrgAccessEntriesApiRequest) Execute() (*PaginatedApiUserAccessListResponse, *http.Response, error) {
	return r.ApiService.ListOrgAccessEntriesExecute(r)
}

/*
ListOrgAccessEntries Return All Access List Entries for One Organization API Key

Returns all access list entries that you configured for the specified organization API key.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key for which you want to return access list entries.
	@return ListOrgAccessEntriesApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) ListOrgAccessEntries(ctx context.Context, orgId string, apiUserId string) ListOrgAccessEntriesApiRequest {
	return ListOrgAccessEntriesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		apiUserId:  apiUserId,
	}
}

// ListOrgAccessEntriesExecute executes the request
//
//	@return PaginatedApiUserAccessListResponse
func (a *ProgrammaticAPIKeysApiService) ListOrgAccessEntriesExecute(r ListOrgAccessEntriesApiRequest) (*PaginatedApiUserAccessListResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiUserAccessListResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.ListOrgAccessEntries")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}/accessList"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

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

type ListOrgApiKeysApiRequest struct {
	ctx          context.Context
	ApiService   ProgrammaticAPIKeysApi
	orgId        string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListOrgApiKeysApiParams struct {
	OrgId        string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ProgrammaticAPIKeysApiService) ListOrgApiKeysWithParams(ctx context.Context, args *ListOrgApiKeysApiParams) ListOrgApiKeysApiRequest {
	return ListOrgApiKeysApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgApiKeysApiRequest) IncludeCount(includeCount bool) ListOrgApiKeysApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgApiKeysApiRequest) ItemsPerPage(itemsPerPage int) ListOrgApiKeysApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgApiKeysApiRequest) PageNum(pageNum int) ListOrgApiKeysApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOrgApiKeysApiRequest) Execute() (*PaginatedApiApiUser, *http.Response, error) {
	return r.ApiService.ListOrgApiKeysExecute(r)
}

/*
ListOrgApiKeys Return All Organization API Keys

Returns all organization API keys for the specified organization. The organization API keys grant programmatic access to an organization. You can't use the API key to log into MongoDB Cloud through the console.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgApiKeysApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) ListOrgApiKeys(ctx context.Context, orgId string) ListOrgApiKeysApiRequest {
	return ListOrgApiKeysApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgApiKeysExecute executes the request
//
//	@return PaginatedApiApiUser
func (a *ProgrammaticAPIKeysApiService) ListOrgApiKeysExecute(r ListOrgApiKeysApiRequest) (*PaginatedApiApiUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiApiUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.ListOrgApiKeys")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys"
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

type RemoveGroupApiKeyApiRequest struct {
	ctx        context.Context
	ApiService ProgrammaticAPIKeysApi
	groupId    string
	apiUserId  string
}

type RemoveGroupApiKeyApiParams struct {
	GroupId   string
	ApiUserId string
}

func (a *ProgrammaticAPIKeysApiService) RemoveGroupApiKeyWithParams(ctx context.Context, args *RemoveGroupApiKeyApiParams) RemoveGroupApiKeyApiRequest {
	return RemoveGroupApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		apiUserId:  args.ApiUserId,
	}
}

func (r RemoveGroupApiKeyApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RemoveGroupApiKeyExecute(r)
}

/*
RemoveGroupApiKey Unassign One Organization API Key from One Project

Removes one organization API key from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to unassign from one project.
	@return RemoveGroupApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) RemoveGroupApiKey(ctx context.Context, groupId string, apiUserId string) RemoveGroupApiKeyApiRequest {
	return RemoveGroupApiKeyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		apiUserId:  apiUserId,
	}
}

// RemoveGroupApiKeyExecute executes the request
func (a *ProgrammaticAPIKeysApiService) RemoveGroupApiKeyExecute(r RemoveGroupApiKeyApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.RemoveGroupApiKey")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/apiKeys/{apiUserId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.apiUserId == "" {
		return nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

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

type UpdateApiKeyRolesApiRequest struct {
	ctx                      context.Context
	ApiService               ProgrammaticAPIKeysApi
	groupId                  string
	apiUserId                string
	updateAtlasProjectApiKey *UpdateAtlasProjectApiKey
	pageNum                  *int
	itemsPerPage             *int
	includeCount             *bool
}

type UpdateApiKeyRolesApiParams struct {
	GroupId                  string
	ApiUserId                string
	UpdateAtlasProjectApiKey *UpdateAtlasProjectApiKey
	PageNum                  *int
	ItemsPerPage             *int
	IncludeCount             *bool
}

func (a *ProgrammaticAPIKeysApiService) UpdateApiKeyRolesWithParams(ctx context.Context, args *UpdateApiKeyRolesApiParams) UpdateApiKeyRolesApiRequest {
	return UpdateApiKeyRolesApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		apiUserId:                args.ApiUserId,
		updateAtlasProjectApiKey: args.UpdateAtlasProjectApiKey,
		pageNum:                  args.PageNum,
		itemsPerPage:             args.ItemsPerPage,
		includeCount:             args.IncludeCount,
	}
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r UpdateApiKeyRolesApiRequest) PageNum(pageNum int) UpdateApiKeyRolesApiRequest {
	r.pageNum = &pageNum
	return r
}

// Number of items that the response returns per page.
func (r UpdateApiKeyRolesApiRequest) ItemsPerPage(itemsPerPage int) UpdateApiKeyRolesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r UpdateApiKeyRolesApiRequest) IncludeCount(includeCount bool) UpdateApiKeyRolesApiRequest {
	r.includeCount = &includeCount
	return r
}

func (r UpdateApiKeyRolesApiRequest) Execute() (*ApiKeyUserDetails, *http.Response, error) {
	return r.ApiService.UpdateApiKeyRolesExecute(r)
}

/*
UpdateApiKeyRoles Update Organization API Key Roles for One Project

Updates the roles of the organization API key that you specify for the project that you specify. You must specify at least one valid role for the project. The application removes any roles that you do not include in this request if they were previously set in the organization API key that you specify for the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key that you want to unassign from one project.
	@return UpdateApiKeyRolesApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) UpdateApiKeyRoles(ctx context.Context, groupId string, apiUserId string, updateAtlasProjectApiKey *UpdateAtlasProjectApiKey) UpdateApiKeyRolesApiRequest {
	return UpdateApiKeyRolesApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		apiUserId:                apiUserId,
		updateAtlasProjectApiKey: updateAtlasProjectApiKey,
	}
}

// UpdateApiKeyRolesExecute executes the request
//
//	@return ApiKeyUserDetails
func (a *ProgrammaticAPIKeysApiService) UpdateApiKeyRolesExecute(r UpdateApiKeyRolesApiRequest) (*ApiKeyUserDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiKeyUserDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.UpdateApiKeyRoles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/apiKeys/{apiUserId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateAtlasProjectApiKey == nil {
		return localVarReturnValue, nil, reportError("updateAtlasProjectApiKey is required and must be specified")
	}

	if r.pageNum != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	} else {
		var defaultValue int = 1
		r.pageNum = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	}
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
	localVarPostBody = r.updateAtlasProjectApiKey
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

type UpdateOrgApiKeyApiRequest struct {
	ctx                           context.Context
	ApiService                    ProgrammaticAPIKeysApi
	orgId                         string
	apiUserId                     string
	updateAtlasOrganizationApiKey *UpdateAtlasOrganizationApiKey
}

type UpdateOrgApiKeyApiParams struct {
	OrgId                         string
	ApiUserId                     string
	UpdateAtlasOrganizationApiKey *UpdateAtlasOrganizationApiKey
}

func (a *ProgrammaticAPIKeysApiService) UpdateOrgApiKeyWithParams(ctx context.Context, args *UpdateOrgApiKeyApiParams) UpdateOrgApiKeyApiRequest {
	return UpdateOrgApiKeyApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         args.OrgId,
		apiUserId:                     args.ApiUserId,
		updateAtlasOrganizationApiKey: args.UpdateAtlasOrganizationApiKey,
	}
}

func (r UpdateOrgApiKeyApiRequest) Execute() (*ApiKeyUserDetails, *http.Response, error) {
	return r.ApiService.UpdateOrgApiKeyExecute(r)
}

/*
UpdateOrgApiKey Update One Organization API Key

Updates one organization API key in the specified organization. The organization API keys  grant programmatic access to an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param apiUserId Unique 24-hexadecimal digit string that identifies this organization API key you  want to update.
	@return UpdateOrgApiKeyApiRequest
*/
func (a *ProgrammaticAPIKeysApiService) UpdateOrgApiKey(ctx context.Context, orgId string, apiUserId string, updateAtlasOrganizationApiKey *UpdateAtlasOrganizationApiKey) UpdateOrgApiKeyApiRequest {
	return UpdateOrgApiKeyApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         orgId,
		apiUserId:                     apiUserId,
		updateAtlasOrganizationApiKey: updateAtlasOrganizationApiKey,
	}
}

// UpdateOrgApiKeyExecute executes the request
//
//	@return ApiKeyUserDetails
func (a *ProgrammaticAPIKeysApiService) UpdateOrgApiKeyExecute(r UpdateOrgApiKeyApiRequest) (*ApiKeyUserDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiKeyUserDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProgrammaticAPIKeysApiService.UpdateOrgApiKey")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/apiKeys/{apiUserId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.apiUserId == "" {
		return localVarReturnValue, nil, reportError("apiUserId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"apiUserId"+"}", url.PathEscape(r.apiUserId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateAtlasOrganizationApiKey == nil {
		return localVarReturnValue, nil, reportError("updateAtlasOrganizationApiKey is required and must be specified")
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
	localVarPostBody = r.updateAtlasOrganizationApiKey
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
