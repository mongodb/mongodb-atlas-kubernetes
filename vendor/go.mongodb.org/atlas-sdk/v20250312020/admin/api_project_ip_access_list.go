// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ProjectIPAccessListApi interface {

	/*
		CreateAccessListEntry Add Entries to Project IP Access List

		Adds one or more access list entries to the specified project. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. Write each entry as either one IP address or one CIDR-notated block of IP addresses. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations. This endpoint doesn't support concurrent `POST` requests. You must submit multiple `POST` requests synchronously.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param networkPermissionEntry One or more access list entries to add to the specified project.
		@return CreateAccessListEntryApiRequest
	*/
	CreateAccessListEntry(ctx context.Context, groupId string, networkPermissionEntry *[]NetworkPermissionEntry) CreateAccessListEntryApiRequest
	/*
		CreateAccessListEntry Add Entries to Project IP Access List


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateAccessListEntryApiParams - Parameters for the request
		@return CreateAccessListEntryApiRequest
	*/
	CreateAccessListEntryWithParams(ctx context.Context, args *CreateAccessListEntryApiParams) CreateAccessListEntryApiRequest

	// Method available only for mocking purposes
	CreateAccessListEntryExecute(r CreateAccessListEntryApiRequest) (*PaginatedNetworkAccess, *http.Response, error)

	/*
		DeleteAccessListEntry Remove One Entry from One Project IP Access List

		Removes one access list entry from the specified project's IP access list. Each entry in the project's IP access list contains one IP address, one CIDR-notated block of IP addresses, or one AWS Security Group ID. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param entryValue Access list entry that you want to remove from the project's IP access list. This value can use one of the following: one AWS security group ID, one IP address, or one CIDR block of addresses. For CIDR blocks that use a subnet mask, replace the forward slash (`/`) with its URL-encoded value (`%2F`). When you remove an entry from the IP access list, existing connections from the removed address or addresses may remain open for a variable amount of time. The amount of time it takes MongoDB Cloud to close the connection depends upon several factors, including:  - how your application established the connection, - how MongoDB Cloud or the driver using the address behaves, and - which protocol (like TCP or UDP) the connection uses.
		@return DeleteAccessListEntryApiRequest
	*/
	DeleteAccessListEntry(ctx context.Context, groupId string, entryValue string) DeleteAccessListEntryApiRequest
	/*
		DeleteAccessListEntry Remove One Entry from One Project IP Access List


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteAccessListEntryApiParams - Parameters for the request
		@return DeleteAccessListEntryApiRequest
	*/
	DeleteAccessListEntryWithParams(ctx context.Context, args *DeleteAccessListEntryApiParams) DeleteAccessListEntryApiRequest

	// Method available only for mocking purposes
	DeleteAccessListEntryExecute(r DeleteAccessListEntryApiRequest) (*http.Response, error)

	/*
		GetAccessListEntry Return One Project IP Access List Entry

		Returns one access list entry from the specified project's IP access list. Each entry in the project's IP access list contains either one IP address or one CIDR-notated block of IP addresses. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. This endpoint (`/groups/{GROUP-ID}/accessList`) manages the Project IP Access List. It doesn't manage the access list for MongoDB Cloud organizations. The Programmatic API Keys endpoint (`/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist`) manages those access lists.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param entryValue Access list entry that you want to return from the project's IP access list. This value can use one of the following: one AWS security group ID, one IP address, or one CIDR block of addresses. For CIDR blocks that use a subnet mask, replace the forward slash (`/`) with its URL-encoded value (`%2F`).
		@return GetAccessListEntryApiRequest
	*/
	GetAccessListEntry(ctx context.Context, groupId string, entryValue string) GetAccessListEntryApiRequest
	/*
		GetAccessListEntry Return One Project IP Access List Entry


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAccessListEntryApiParams - Parameters for the request
		@return GetAccessListEntryApiRequest
	*/
	GetAccessListEntryWithParams(ctx context.Context, args *GetAccessListEntryApiParams) GetAccessListEntryApiRequest

	// Method available only for mocking purposes
	GetAccessListEntryExecute(r GetAccessListEntryApiRequest) (*NetworkPermissionEntry, *http.Response, error)

	/*
		GetAccessListStatus Return Status of One Project IP Access List Entry

		Returns the status of one project IP access list entry. This resource checks if the provided project IP access list entry applies to all cloud providers serving clusters from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param entryValue Network address or cloud provider security construct that identifies which project access list entry to be verified.
		@return GetAccessListStatusApiRequest
	*/
	GetAccessListStatus(ctx context.Context, groupId string, entryValue string) GetAccessListStatusApiRequest
	/*
		GetAccessListStatus Return Status of One Project IP Access List Entry


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAccessListStatusApiParams - Parameters for the request
		@return GetAccessListStatusApiRequest
	*/
	GetAccessListStatusWithParams(ctx context.Context, args *GetAccessListStatusApiParams) GetAccessListStatusApiRequest

	// Method available only for mocking purposes
	GetAccessListStatusExecute(r GetAccessListStatusApiRequest) (*NetworkPermissionEntryStatus, *http.Response, error)

	/*
		ListAccessListEntries Return All Project IP Access List Entries

		Returns all access list entries from the specified project's IP access list. Each entry in the project's IP access list contains either one IP address or one CIDR-notated block of IP addresses. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListAccessListEntriesApiRequest
	*/
	ListAccessListEntries(ctx context.Context, groupId string) ListAccessListEntriesApiRequest
	/*
		ListAccessListEntries Return All Project IP Access List Entries


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListAccessListEntriesApiParams - Parameters for the request
		@return ListAccessListEntriesApiRequest
	*/
	ListAccessListEntriesWithParams(ctx context.Context, args *ListAccessListEntriesApiParams) ListAccessListEntriesApiRequest

	// Method available only for mocking purposes
	ListAccessListEntriesExecute(r ListAccessListEntriesApiRequest) (*PaginatedNetworkAccess, *http.Response, error)
}

// ProjectIPAccessListApiService ProjectIPAccessListApi service
type ProjectIPAccessListApiService service

type CreateAccessListEntryApiRequest struct {
	ctx                    context.Context
	ApiService             ProjectIPAccessListApi
	groupId                string
	networkPermissionEntry *[]NetworkPermissionEntry
	includeCount           *bool
	itemsPerPage           *int
	pageNum                *int
}

type CreateAccessListEntryApiParams struct {
	GroupId                string
	NetworkPermissionEntry *[]NetworkPermissionEntry
	IncludeCount           *bool
	ItemsPerPage           *int
	PageNum                *int
}

func (a *ProjectIPAccessListApiService) CreateAccessListEntryWithParams(ctx context.Context, args *CreateAccessListEntryApiParams) CreateAccessListEntryApiRequest {
	return CreateAccessListEntryApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		networkPermissionEntry: args.NetworkPermissionEntry,
		includeCount:           args.IncludeCount,
		itemsPerPage:           args.ItemsPerPage,
		pageNum:                args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r CreateAccessListEntryApiRequest) IncludeCount(includeCount bool) CreateAccessListEntryApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r CreateAccessListEntryApiRequest) ItemsPerPage(itemsPerPage int) CreateAccessListEntryApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r CreateAccessListEntryApiRequest) PageNum(pageNum int) CreateAccessListEntryApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r CreateAccessListEntryApiRequest) Execute() (*PaginatedNetworkAccess, *http.Response, error) {
	return r.ApiService.CreateAccessListEntryExecute(r)
}

/*
CreateAccessListEntry Add Entries to Project IP Access List

Adds one or more access list entries to the specified project. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. Write each entry as either one IP address or one CIDR-notated block of IP addresses. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations. This endpoint doesn't support concurrent `POST` requests. You must submit multiple `POST` requests synchronously.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateAccessListEntryApiRequest
*/
func (a *ProjectIPAccessListApiService) CreateAccessListEntry(ctx context.Context, groupId string, networkPermissionEntry *[]NetworkPermissionEntry) CreateAccessListEntryApiRequest {
	return CreateAccessListEntryApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		networkPermissionEntry: networkPermissionEntry,
	}
}

// CreateAccessListEntryExecute executes the request
//
//	@return PaginatedNetworkAccess
func (a *ProjectIPAccessListApiService) CreateAccessListEntryExecute(r CreateAccessListEntryApiRequest) (*PaginatedNetworkAccess, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedNetworkAccess
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectIPAccessListApiService.CreateAccessListEntry")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/accessList"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.networkPermissionEntry == nil {
		return localVarReturnValue, nil, reportError("networkPermissionEntry is required and must be specified")
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
	localVarPostBody = r.networkPermissionEntry
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

type DeleteAccessListEntryApiRequest struct {
	ctx        context.Context
	ApiService ProjectIPAccessListApi
	groupId    string
	entryValue string
}

type DeleteAccessListEntryApiParams struct {
	GroupId    string
	EntryValue string
}

func (a *ProjectIPAccessListApiService) DeleteAccessListEntryWithParams(ctx context.Context, args *DeleteAccessListEntryApiParams) DeleteAccessListEntryApiRequest {
	return DeleteAccessListEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		entryValue: args.EntryValue,
	}
}

func (r DeleteAccessListEntryApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteAccessListEntryExecute(r)
}

/*
DeleteAccessListEntry Remove One Entry from One Project IP Access List

Removes one access list entry from the specified project's IP access list. Each entry in the project's IP access list contains one IP address, one CIDR-notated block of IP addresses, or one AWS Security Group ID. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param entryValue Access list entry that you want to remove from the project's IP access list. This value can use one of the following: one AWS security group ID, one IP address, or one CIDR block of addresses. For CIDR blocks that use a subnet mask, replace the forward slash (`/`) with its URL-encoded value (`%2F`). When you remove an entry from the IP access list, existing connections from the removed address or addresses may remain open for a variable amount of time. The amount of time it takes MongoDB Cloud to close the connection depends upon several factors, including:  - how your application established the connection, - how MongoDB Cloud or the driver using the address behaves, and - which protocol (like TCP or UDP) the connection uses.
	@return DeleteAccessListEntryApiRequest
*/
func (a *ProjectIPAccessListApiService) DeleteAccessListEntry(ctx context.Context, groupId string, entryValue string) DeleteAccessListEntryApiRequest {
	return DeleteAccessListEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		entryValue: entryValue,
	}
}

// DeleteAccessListEntryExecute executes the request
func (a *ProjectIPAccessListApiService) DeleteAccessListEntryExecute(r DeleteAccessListEntryApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectIPAccessListApiService.DeleteAccessListEntry")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/accessList/{entryValue}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.entryValue == "" {
		return nil, reportError("entryValue is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"entryValue"+"}", url.PathEscape(r.entryValue), -1)

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

type GetAccessListEntryApiRequest struct {
	ctx        context.Context
	ApiService ProjectIPAccessListApi
	groupId    string
	entryValue string
}

type GetAccessListEntryApiParams struct {
	GroupId    string
	EntryValue string
}

func (a *ProjectIPAccessListApiService) GetAccessListEntryWithParams(ctx context.Context, args *GetAccessListEntryApiParams) GetAccessListEntryApiRequest {
	return GetAccessListEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		entryValue: args.EntryValue,
	}
}

func (r GetAccessListEntryApiRequest) Execute() (*NetworkPermissionEntry, *http.Response, error) {
	return r.ApiService.GetAccessListEntryExecute(r)
}

/*
GetAccessListEntry Return One Project IP Access List Entry

Returns one access list entry from the specified project's IP access list. Each entry in the project's IP access list contains either one IP address or one CIDR-notated block of IP addresses. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. This endpoint (`/groups/{GROUP-ID}/accessList`) manages the Project IP Access List. It doesn't manage the access list for MongoDB Cloud organizations. The Programmatic API Keys endpoint (`/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist`) manages those access lists.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param entryValue Access list entry that you want to return from the project's IP access list. This value can use one of the following: one AWS security group ID, one IP address, or one CIDR block of addresses. For CIDR blocks that use a subnet mask, replace the forward slash (`/`) with its URL-encoded value (`%2F`).
	@return GetAccessListEntryApiRequest
*/
func (a *ProjectIPAccessListApiService) GetAccessListEntry(ctx context.Context, groupId string, entryValue string) GetAccessListEntryApiRequest {
	return GetAccessListEntryApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		entryValue: entryValue,
	}
}

// GetAccessListEntryExecute executes the request
//
//	@return NetworkPermissionEntry
func (a *ProjectIPAccessListApiService) GetAccessListEntryExecute(r GetAccessListEntryApiRequest) (*NetworkPermissionEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *NetworkPermissionEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectIPAccessListApiService.GetAccessListEntry")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/accessList/{entryValue}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.entryValue == "" {
		return localVarReturnValue, nil, reportError("entryValue is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"entryValue"+"}", url.PathEscape(r.entryValue), -1)

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

type GetAccessListStatusApiRequest struct {
	ctx        context.Context
	ApiService ProjectIPAccessListApi
	groupId    string
	entryValue string
}

type GetAccessListStatusApiParams struct {
	GroupId    string
	EntryValue string
}

func (a *ProjectIPAccessListApiService) GetAccessListStatusWithParams(ctx context.Context, args *GetAccessListStatusApiParams) GetAccessListStatusApiRequest {
	return GetAccessListStatusApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		entryValue: args.EntryValue,
	}
}

func (r GetAccessListStatusApiRequest) Execute() (*NetworkPermissionEntryStatus, *http.Response, error) {
	return r.ApiService.GetAccessListStatusExecute(r)
}

/*
GetAccessListStatus Return Status of One Project IP Access List Entry

Returns the status of one project IP access list entry. This resource checks if the provided project IP access list entry applies to all cloud providers serving clusters from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param entryValue Network address or cloud provider security construct that identifies which project access list entry to be verified.
	@return GetAccessListStatusApiRequest
*/
func (a *ProjectIPAccessListApiService) GetAccessListStatus(ctx context.Context, groupId string, entryValue string) GetAccessListStatusApiRequest {
	return GetAccessListStatusApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		entryValue: entryValue,
	}
}

// GetAccessListStatusExecute executes the request
//
//	@return NetworkPermissionEntryStatus
func (a *ProjectIPAccessListApiService) GetAccessListStatusExecute(r GetAccessListStatusApiRequest) (*NetworkPermissionEntryStatus, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *NetworkPermissionEntryStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectIPAccessListApiService.GetAccessListStatus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/accessList/{entryValue}/status"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.entryValue == "" {
		return localVarReturnValue, nil, reportError("entryValue is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"entryValue"+"}", url.PathEscape(r.entryValue), -1)

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

type ListAccessListEntriesApiRequest struct {
	ctx          context.Context
	ApiService   ProjectIPAccessListApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListAccessListEntriesApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ProjectIPAccessListApiService) ListAccessListEntriesWithParams(ctx context.Context, args *ListAccessListEntriesApiParams) ListAccessListEntriesApiRequest {
	return ListAccessListEntriesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListAccessListEntriesApiRequest) IncludeCount(includeCount bool) ListAccessListEntriesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListAccessListEntriesApiRequest) ItemsPerPage(itemsPerPage int) ListAccessListEntriesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListAccessListEntriesApiRequest) PageNum(pageNum int) ListAccessListEntriesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListAccessListEntriesApiRequest) Execute() (*PaginatedNetworkAccess, *http.Response, error) {
	return r.ApiService.ListAccessListEntriesExecute(r)
}

/*
ListAccessListEntries Return All Project IP Access List Entries

Returns all access list entries from the specified project's IP access list. Each entry in the project's IP access list contains either one IP address or one CIDR-notated block of IP addresses. MongoDB Cloud only allows client connections to the cluster from entries in the project's IP access list. This resource replaces the whitelist resource. MongoDB Cloud removed whitelists in July 2021. Update your applications to use this new resource. The `/groups/{GROUP-ID}/accessList` endpoint manages the database IP access list. This endpoint is distinct from the `orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/accesslist` endpoint, which manages the access list for MongoDB Cloud organizations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListAccessListEntriesApiRequest
*/
func (a *ProjectIPAccessListApiService) ListAccessListEntries(ctx context.Context, groupId string) ListAccessListEntriesApiRequest {
	return ListAccessListEntriesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListAccessListEntriesExecute executes the request
//
//	@return PaginatedNetworkAccess
func (a *ProjectIPAccessListApiService) ListAccessListEntriesExecute(r ListAccessListEntriesApiRequest) (*PaginatedNetworkAccess, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedNetworkAccess
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ProjectIPAccessListApiService.ListAccessListEntries")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/accessList"
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
