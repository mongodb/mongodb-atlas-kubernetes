// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RateLimitingApi interface {

	/*
		GetRateLimit Return One Rate Limit

		Get one rate limit endpoint set.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param endpointSetId The ID of the rate limit endpoint set.
		@return GetRateLimitApiRequest
	*/
	GetRateLimit(ctx context.Context, endpointSetId string) GetRateLimitApiRequest
	/*
		GetRateLimit Return One Rate Limit


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetRateLimitApiParams - Parameters for the request
		@return GetRateLimitApiRequest
	*/
	GetRateLimitWithParams(ctx context.Context, args *GetRateLimitApiParams) GetRateLimitApiRequest

	// Method available only for mocking purposes
	GetRateLimitExecute(r GetRateLimitApiRequest) (*RateLimitEndpointSetResponse, *http.Response, error)

	/*
		ListRateLimits Return All Rate Limits

		Get all rate limits for all v2 Atlas Administration API endpoint sets.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListRateLimitsApiRequest
	*/
	ListRateLimits(ctx context.Context) ListRateLimitsApiRequest
	/*
		ListRateLimits Return All Rate Limits


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListRateLimitsApiParams - Parameters for the request
		@return ListRateLimitsApiRequest
	*/
	ListRateLimitsWithParams(ctx context.Context, args *ListRateLimitsApiParams) ListRateLimitsApiRequest

	// Method available only for mocking purposes
	ListRateLimitsExecute(r ListRateLimitsApiRequest) (*PaginatedRateLimitEndpointSets, *http.Response, error)
}

// RateLimitingApiService RateLimitingApi service
type RateLimitingApiService service

type GetRateLimitApiRequest struct {
	ctx           context.Context
	ApiService    RateLimitingApi
	endpointSetId string
	groupId       *string
	orgId         *string
	userId        *string
	ipAddress     *string
}

type GetRateLimitApiParams struct {
	EndpointSetId string
	GroupId       *string
	OrgId         *string
	UserId        *string
	IpAddress     *string
}

func (a *RateLimitingApiService) GetRateLimitWithParams(ctx context.Context, args *GetRateLimitApiParams) GetRateLimitApiRequest {
	return GetRateLimitApiRequest{
		ApiService:    a,
		ctx:           ctx,
		endpointSetId: args.EndpointSetId,
		groupId:       args.GroupId,
		orgId:         args.OrgId,
		userId:        args.UserId,
		ipAddress:     args.IpAddress,
	}
}

// Unique 24-hexadecimal digit string that identifies the Atlas Project to request rate limits for. When this parameter is provided, the limits returned are applicable to the specified project. The requesting user must have the Project Read Only role for the specified project.
func (r GetRateLimitApiRequest) GroupId(groupId string) GetRateLimitApiRequest {
	r.groupId = &groupId
	return r
}

// Unique 24-hexadecimal digit string that identifies the Atlas Organization to request rate limits for. When this parameter is provided, the limits returned are applicable to the specified organization. The requesting user must have the Organization Read Only role for the specified organization.
func (r GetRateLimitApiRequest) OrgId(orgId string) GetRateLimitApiRequest {
	r.orgId = &orgId
	return r
}

// A string that identifies the Atlas user to request rate limits for. The ID can for example be the Service Account Client ID or the API Public Key. When this parameter is provided, the limits returned are applicable to the specified  user. The requesting user must be the same as the specified user.
func (r GetRateLimitApiRequest) UserId(userId string) GetRateLimitApiRequest {
	r.userId = &userId
	return r
}

// An IP address to request rate limits for. When this parameter is provided, the limits returned are applicable to the specified IP address. The requesting user must have the same IP address as the one provided in the request.
func (r GetRateLimitApiRequest) IpAddress(ipAddress string) GetRateLimitApiRequest {
	r.ipAddress = &ipAddress
	return r
}

func (r GetRateLimitApiRequest) Execute() (*RateLimitEndpointSetResponse, *http.Response, error) {
	return r.ApiService.GetRateLimitExecute(r)
}

/*
GetRateLimit Return One Rate Limit

Get one rate limit endpoint set.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param endpointSetId The ID of the rate limit endpoint set.
	@return GetRateLimitApiRequest
*/
func (a *RateLimitingApiService) GetRateLimit(ctx context.Context, endpointSetId string) GetRateLimitApiRequest {
	return GetRateLimitApiRequest{
		ApiService:    a,
		ctx:           ctx,
		endpointSetId: endpointSetId,
	}
}

// GetRateLimitExecute executes the request
//
//	@return RateLimitEndpointSetResponse
func (a *RateLimitingApiService) GetRateLimitExecute(r GetRateLimitApiRequest) (*RateLimitEndpointSetResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *RateLimitEndpointSetResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RateLimitingApiService.GetRateLimit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/rateLimits/{endpointSetId}"
	if r.endpointSetId == "" {
		return localVarReturnValue, nil, reportError("endpointSetId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointSetId"+"}", url.PathEscape(r.endpointSetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.groupId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "groupId", r.groupId, "")
	}
	if r.orgId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgId", r.orgId, "")
	}
	if r.userId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "userId", r.userId, "")
	}
	if r.ipAddress != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "ipAddress", r.ipAddress, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

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

type ListRateLimitsApiRequest struct {
	ctx          context.Context
	ApiService   RateLimitingApi
	itemsPerPage *int
	pageNum      *int
	groupId      *string
	orgId        *string
	userId       *string
	ipAddress    *string
	name         *string
	endpointPath *string
}

type ListRateLimitsApiParams struct {
	ItemsPerPage *int
	PageNum      *int
	GroupId      *string
	OrgId        *string
	UserId       *string
	IpAddress    *string
	Name         *string
	EndpointPath *string
}

func (a *RateLimitingApiService) ListRateLimitsWithParams(ctx context.Context, args *ListRateLimitsApiParams) ListRateLimitsApiRequest {
	return ListRateLimitsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		groupId:      args.GroupId,
		orgId:        args.OrgId,
		userId:       args.UserId,
		ipAddress:    args.IpAddress,
		name:         args.Name,
		endpointPath: args.EndpointPath,
	}
}

// Number of items that the response returns per page.
func (r ListRateLimitsApiRequest) ItemsPerPage(itemsPerPage int) ListRateLimitsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListRateLimitsApiRequest) PageNum(pageNum int) ListRateLimitsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Unique 24-hexadecimal digit string that identifies the Atlas Project to request rate limits for. When this parameter is provided, only group scoped endpoint sets are returned and the limits returned are applicable to the specified project. The requesting user must have the Project Read Only role for the specified project.
func (r ListRateLimitsApiRequest) GroupId(groupId string) ListRateLimitsApiRequest {
	r.groupId = &groupId
	return r
}

// Unique 24-hexadecimal digit string that identifies the Atlas Organization to request rate limits for. When this parameter is provided, only organization scoped endpoint sets are returned and the limits returned are applicable to the specified organization. The requesting user must have the Organization Read Only role for the specified organization.
func (r ListRateLimitsApiRequest) OrgId(orgId string) ListRateLimitsApiRequest {
	r.orgId = &orgId
	return r
}

// A string that identifies the Atlas user to request rate limits for. The ID can for example be the Service Account Client ID or the API Public Key. When this parameter is provided, only user scoped endpoint sets are returned and the limits returned are applicable to the specified user. The requesting user must be the same as the specified user.
func (r ListRateLimitsApiRequest) UserId(userId string) ListRateLimitsApiRequest {
	r.userId = &userId
	return r
}

// An IP address to request rate limits for. When this parameter is provided, only IP scoped endpoint sets are returned and the limits returned are applicable to the specified IP address. The requesting user must have the same IP address as the one provided in the request.
func (r ListRateLimitsApiRequest) IpAddress(ipAddress string) ListRateLimitsApiRequest {
	r.ipAddress = &ipAddress
	return r
}

// Filters the returned endpoint sets by the provided endpoint set name. Multiple names may be provided, for example &#x60;/rateLimits?name&#x3D;Name1&amp;name&#x3D;Name2&#x60;. For names that use spaces, replace the space with its URL-encoded value (&#x60;%20&#x60;).
func (r ListRateLimitsApiRequest) Name(name string) ListRateLimitsApiRequest {
	r.name = &name
	return r
}

// Filters the returned endpoint sets by the provided endpoint path. Multiple paths may be provided, for example &#x60;/rateLimits?endpointPath&#x3D;%2Fapi%2Fatlas%2Fv2%2Fclusters&amp;endpointPath&#x3D;%2Fapi%2Fatlas%2Fv2%2Fgroups%2F%7BgroupId%7D%2F&#x60;. Replace &#x60;/&#x60;, &#x60;{&#x60; and &#x60;}&#x60; with their URL-encoded values (&#x60;%2F&#x60;, &#x60;%7B&#x60; and &#x60;%7D&#x60; respectively).
func (r ListRateLimitsApiRequest) EndpointPath(endpointPath string) ListRateLimitsApiRequest {
	r.endpointPath = &endpointPath
	return r
}

func (r ListRateLimitsApiRequest) Execute() (*PaginatedRateLimitEndpointSets, *http.Response, error) {
	return r.ApiService.ListRateLimitsExecute(r)
}

/*
ListRateLimits Return All Rate Limits

Get all rate limits for all v2 Atlas Administration API endpoint sets.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListRateLimitsApiRequest
*/
func (a *RateLimitingApiService) ListRateLimits(ctx context.Context) ListRateLimitsApiRequest {
	return ListRateLimitsApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListRateLimitsExecute executes the request
//
//	@return PaginatedRateLimitEndpointSets
func (a *RateLimitingApiService) ListRateLimitsExecute(r ListRateLimitsApiRequest) (*PaginatedRateLimitEndpointSets, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedRateLimitEndpointSets
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RateLimitingApiService.ListRateLimits")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/rateLimits"

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
	if r.groupId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "groupId", r.groupId, "")
	}
	if r.orgId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orgId", r.orgId, "")
	}
	if r.userId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "userId", r.userId, "")
	}
	if r.ipAddress != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "ipAddress", r.ipAddress, "")
	}
	if r.name != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name", r.name, "")
	}
	if r.endpointPath != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endpointPath", r.endpointPath, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

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
