// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type ActivityFeedApi interface {

	/*
		GetGroupActivityFeed Return Pre-Filtered Activity Feed Link for One Project

		Returns a pre-filtered activity feed link for the specified project based on the provided date range and event types. The returned link can be shared and opened to view the activity feed with the same filters applied.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupActivityFeedApiRequest
	*/
	GetGroupActivityFeed(ctx context.Context, groupId string) GetGroupActivityFeedApiRequest
	/*
		GetGroupActivityFeed Return Pre-Filtered Activity Feed Link for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupActivityFeedApiParams - Parameters for the request
		@return GetGroupActivityFeedApiRequest
	*/
	GetGroupActivityFeedWithParams(ctx context.Context, args *GetGroupActivityFeedApiParams) GetGroupActivityFeedApiRequest

	// Method available only for mocking purposes
	GetGroupActivityFeedExecute(r GetGroupActivityFeedApiRequest) (*ActivityFeedLinkResponse, *http.Response, error)

	/*
		GetOrgActivityFeed Return Pre-Filtered Activity Feed Link for One Organization

		Returns a pre-filtered activity feed link for the specified organization based on the provided date range and event types. The returned link can be shared and opened to view the activity feed with the same filters applied.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return GetOrgActivityFeedApiRequest
	*/
	GetOrgActivityFeed(ctx context.Context, orgId string) GetOrgActivityFeedApiRequest
	/*
		GetOrgActivityFeed Return Pre-Filtered Activity Feed Link for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgActivityFeedApiParams - Parameters for the request
		@return GetOrgActivityFeedApiRequest
	*/
	GetOrgActivityFeedWithParams(ctx context.Context, args *GetOrgActivityFeedApiParams) GetOrgActivityFeedApiRequest

	// Method available only for mocking purposes
	GetOrgActivityFeedExecute(r GetOrgActivityFeedApiRequest) (*ActivityFeedLinkResponse, *http.Response, error)
}

// ActivityFeedApiService ActivityFeedApi service
type ActivityFeedApiService service

type GetGroupActivityFeedApiRequest struct {
	ctx         context.Context
	ApiService  ActivityFeedApi
	groupId     string
	eventType   *[]string
	maxDate     *time.Time
	minDate     *time.Time
	clusterName *[]string
}

type GetGroupActivityFeedApiParams struct {
	GroupId     string
	EventType   *[]string
	MaxDate     *time.Time
	MinDate     *time.Time
	ClusterName *[]string
}

func (a *ActivityFeedApiService) GetGroupActivityFeedWithParams(ctx context.Context, args *GetGroupActivityFeedApiParams) GetGroupActivityFeedApiRequest {
	return GetGroupActivityFeedApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		eventType:   args.EventType,
		maxDate:     args.MaxDate,
		minDate:     args.MinDate,
		clusterName: args.ClusterName,
	}
}

// Category of incident recorded at this moment in time.  **IMPORTANT**: The complete list of event type values changes frequently.
func (r GetGroupActivityFeedApiRequest) EventType(eventType []string) GetGroupActivityFeedApiRequest {
	r.eventType = &eventType
	return r
}

// End date and time for events to include in the activity feed link. ISO 8601 timestamp format in UTC.
func (r GetGroupActivityFeedApiRequest) MaxDate(maxDate time.Time) GetGroupActivityFeedApiRequest {
	r.maxDate = &maxDate
	return r
}

// Start date and time for events to include in the activity feed link. ISO 8601 timestamp format in UTC.
func (r GetGroupActivityFeedApiRequest) MinDate(minDate time.Time) GetGroupActivityFeedApiRequest {
	r.minDate = &minDate
	return r
}

// Human-readable label that identifies the cluster.
func (r GetGroupActivityFeedApiRequest) ClusterName(clusterName []string) GetGroupActivityFeedApiRequest {
	r.clusterName = &clusterName
	return r
}

func (r GetGroupActivityFeedApiRequest) Execute() (*ActivityFeedLinkResponse, *http.Response, error) {
	return r.ApiService.GetGroupActivityFeedExecute(r)
}

/*
GetGroupActivityFeed Return Pre-Filtered Activity Feed Link for One Project

Returns a pre-filtered activity feed link for the specified project based on the provided date range and event types. The returned link can be shared and opened to view the activity feed with the same filters applied.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupActivityFeedApiRequest
*/
func (a *ActivityFeedApiService) GetGroupActivityFeed(ctx context.Context, groupId string) GetGroupActivityFeedApiRequest {
	return GetGroupActivityFeedApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetGroupActivityFeedExecute executes the request
//
//	@return ActivityFeedLinkResponse
func (a *ActivityFeedApiService) GetGroupActivityFeedExecute(r GetGroupActivityFeedApiRequest) (*ActivityFeedLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ActivityFeedLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ActivityFeedApiService.GetGroupActivityFeed")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/activityFeed"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.eventType != nil {
		t := *r.eventType
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "eventType", t, "multi")

	}
	if r.maxDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "maxDate", r.maxDate, "")
	}
	if r.minDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "minDate", r.minDate, "")
	}
	if r.clusterName != nil {
		t := *r.clusterName
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "clusterName", t, "multi")

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

type GetOrgActivityFeedApiRequest struct {
	ctx        context.Context
	ApiService ActivityFeedApi
	orgId      string
	eventType  *[]string
	maxDate    *time.Time
	minDate    *time.Time
}

type GetOrgActivityFeedApiParams struct {
	OrgId     string
	EventType *[]string
	MaxDate   *time.Time
	MinDate   *time.Time
}

func (a *ActivityFeedApiService) GetOrgActivityFeedWithParams(ctx context.Context, args *GetOrgActivityFeedApiParams) GetOrgActivityFeedApiRequest {
	return GetOrgActivityFeedApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		eventType:  args.EventType,
		maxDate:    args.MaxDate,
		minDate:    args.MinDate,
	}
}

// Category of incident recorded at this moment in time.  **IMPORTANT**: The complete list of event type values changes frequently.
func (r GetOrgActivityFeedApiRequest) EventType(eventType []string) GetOrgActivityFeedApiRequest {
	r.eventType = &eventType
	return r
}

// End date and time for events to include in the activity feed link. ISO 8601 timestamp format in UTC.
func (r GetOrgActivityFeedApiRequest) MaxDate(maxDate time.Time) GetOrgActivityFeedApiRequest {
	r.maxDate = &maxDate
	return r
}

// Start date and time for events to include in the activity feed link. ISO 8601 timestamp format in UTC.
func (r GetOrgActivityFeedApiRequest) MinDate(minDate time.Time) GetOrgActivityFeedApiRequest {
	r.minDate = &minDate
	return r
}

func (r GetOrgActivityFeedApiRequest) Execute() (*ActivityFeedLinkResponse, *http.Response, error) {
	return r.ApiService.GetOrgActivityFeedExecute(r)
}

/*
GetOrgActivityFeed Return Pre-Filtered Activity Feed Link for One Organization

Returns a pre-filtered activity feed link for the specified organization based on the provided date range and event types. The returned link can be shared and opened to view the activity feed with the same filters applied.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return GetOrgActivityFeedApiRequest
*/
func (a *ActivityFeedApiService) GetOrgActivityFeed(ctx context.Context, orgId string) GetOrgActivityFeedApiRequest {
	return GetOrgActivityFeedApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// GetOrgActivityFeedExecute executes the request
//
//	@return ActivityFeedLinkResponse
func (a *ActivityFeedApiService) GetOrgActivityFeedExecute(r GetOrgActivityFeedApiRequest) (*ActivityFeedLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ActivityFeedLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ActivityFeedApiService.GetOrgActivityFeed")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/activityFeed"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.eventType != nil {
		t := *r.eventType
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "eventType", t, "multi")

	}
	if r.maxDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "maxDate", r.maxDate, "")
	}
	if r.minDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "minDate", r.minDate, "")
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
