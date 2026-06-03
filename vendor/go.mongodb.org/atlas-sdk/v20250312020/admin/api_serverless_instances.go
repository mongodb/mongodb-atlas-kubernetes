// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ServerlessInstancesApi interface {

	/*
			GetServerlessInstance Return One Serverless Instance from One Project

			Returns details for one serverless instance in the specified project.

		This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. Continuous backups are not supported and `serverlessContinuousBackupEnabled` will not take effect on these clusters. This endpoint will be sunset on January 15, 2027. Please use the Get Flex Cluster endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param name Human-readable label that identifies the serverless instance.
			@return GetServerlessInstanceApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessInstancesApi
	*/
	GetServerlessInstance(ctx context.Context, groupId string, name string) GetServerlessInstanceApiRequest
	/*
		GetServerlessInstance Return One Serverless Instance from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServerlessInstanceApiParams - Parameters for the request
		@return GetServerlessInstanceApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessInstancesApi
	*/
	GetServerlessInstanceWithParams(ctx context.Context, args *GetServerlessInstanceApiParams) GetServerlessInstanceApiRequest

	// Method available only for mocking purposes
	GetServerlessInstanceExecute(r GetServerlessInstanceApiRequest) (*ServerlessInstanceDescription, *http.Response, error)

	/*
			ListServerlessInstances Return All Serverless Instances in One Project

			Returns details for all serverless instances in the specified project.

		This endpoint also lists Flex clusters that were created using the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or former Serverless instances that have been migrated to Flex clusters, until January 15, 2026 after which this endpoint will begin returning an empty list. The endpoint will be removed entirely on January 15, 2027. Continuous backups are not supported and `serverlessContinuousBackupEnabled` will not take effect on these clusters. Please use the List Flex Clusters endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListServerlessInstancesApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessInstancesApi
	*/
	ListServerlessInstances(ctx context.Context, groupId string) ListServerlessInstancesApiRequest
	/*
		ListServerlessInstances Return All Serverless Instances in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListServerlessInstancesApiParams - Parameters for the request
		@return ListServerlessInstancesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessInstancesApi
	*/
	ListServerlessInstancesWithParams(ctx context.Context, args *ListServerlessInstancesApiParams) ListServerlessInstancesApiRequest

	// Method available only for mocking purposes
	ListServerlessInstancesExecute(r ListServerlessInstancesApiRequest) (*PaginatedServerlessInstanceDescription, *http.Response, error)
}

// ServerlessInstancesApiService ServerlessInstancesApi service
type ServerlessInstancesApiService service

type GetServerlessInstanceApiRequest struct {
	ctx        context.Context
	ApiService ServerlessInstancesApi
	groupId    string
	name       string
}

type GetServerlessInstanceApiParams struct {
	GroupId string
	Name    string
}

func (a *ServerlessInstancesApiService) GetServerlessInstanceWithParams(ctx context.Context, args *GetServerlessInstanceApiParams) GetServerlessInstanceApiRequest {
	return GetServerlessInstanceApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		name:       args.Name,
	}
}

func (r GetServerlessInstanceApiRequest) Execute() (*ServerlessInstanceDescription, *http.Response, error) {
	return r.ApiService.GetServerlessInstanceExecute(r)
}

/*
GetServerlessInstance Return One Serverless Instance from One Project

Returns details for one serverless instance in the specified project.

This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. Continuous backups are not supported and `serverlessContinuousBackupEnabled` will not take effect on these clusters. This endpoint will be sunset on January 15, 2027. Please use the Get Flex Cluster endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param name Human-readable label that identifies the serverless instance.
	@return GetServerlessInstanceApiRequest

Deprecated
*/
func (a *ServerlessInstancesApiService) GetServerlessInstance(ctx context.Context, groupId string, name string) GetServerlessInstanceApiRequest {
	return GetServerlessInstanceApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		name:       name,
	}
}

// GetServerlessInstanceExecute executes the request
//
//	@return ServerlessInstanceDescription
//
// Deprecated
func (a *ServerlessInstancesApiService) GetServerlessInstanceExecute(r GetServerlessInstanceApiRequest) (*ServerlessInstanceDescription, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessInstanceDescription
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessInstancesApiService.GetServerlessInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{name}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.name == "" {
		return localVarReturnValue, nil, reportError("name is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", url.PathEscape(r.name), -1)

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

type ListServerlessInstancesApiRequest struct {
	ctx          context.Context
	ApiService   ServerlessInstancesApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListServerlessInstancesApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ServerlessInstancesApiService) ListServerlessInstancesWithParams(ctx context.Context, args *ListServerlessInstancesApiParams) ListServerlessInstancesApiRequest {
	return ListServerlessInstancesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListServerlessInstancesApiRequest) IncludeCount(includeCount bool) ListServerlessInstancesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListServerlessInstancesApiRequest) ItemsPerPage(itemsPerPage int) ListServerlessInstancesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListServerlessInstancesApiRequest) PageNum(pageNum int) ListServerlessInstancesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListServerlessInstancesApiRequest) Execute() (*PaginatedServerlessInstanceDescription, *http.Response, error) {
	return r.ApiService.ListServerlessInstancesExecute(r)
}

/*
ListServerlessInstances Return All Serverless Instances in One Project

Returns details for all serverless instances in the specified project.

This endpoint also lists Flex clusters that were created using the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or former Serverless instances that have been migrated to Flex clusters, until January 15, 2026 after which this endpoint will begin returning an empty list. The endpoint will be removed entirely on January 15, 2027. Continuous backups are not supported and `serverlessContinuousBackupEnabled` will not take effect on these clusters. Please use the List Flex Clusters endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListServerlessInstancesApiRequest

Deprecated
*/
func (a *ServerlessInstancesApiService) ListServerlessInstances(ctx context.Context, groupId string) ListServerlessInstancesApiRequest {
	return ListServerlessInstancesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListServerlessInstancesExecute executes the request
//
//	@return PaginatedServerlessInstanceDescription
//
// Deprecated
func (a *ServerlessInstancesApiService) ListServerlessInstancesExecute(r ListServerlessInstancesApiRequest) (*PaginatedServerlessInstanceDescription, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedServerlessInstanceDescription
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessInstancesApiService.ListServerlessInstances")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless"
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
