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

type EventsApi interface {

	/*
			GetGroupEvent Return One Event from One Project

			Returns one event for the specified project. Events identify significant database, billing, or security activities or status changes. Use the Return Events from One Project endpoint to retrieve all events to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param eventId Unique 24-hexadecimal digit string that identifies the event that you want to return.
			@return GetGroupEventApiRequest
	*/
	GetGroupEvent(ctx context.Context, groupId string, eventId string) GetGroupEventApiRequest
	/*
		GetGroupEvent Return One Event from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupEventApiParams - Parameters for the request
		@return GetGroupEventApiRequest
	*/
	GetGroupEventWithParams(ctx context.Context, args *GetGroupEventApiParams) GetGroupEventApiRequest

	// Method available only for mocking purposes
	GetGroupEventExecute(r GetGroupEventApiRequest) (*EventViewForNdsGroup, *http.Response, error)

	/*
			GetOrgEvent Return One Event from One Organization

			Returns one event for the specified organization. Events identify significant database, billing, or security activities or status changes. Use the Return Events from One Organization endpoint to retrieve all events to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param eventId Unique 24-hexadecimal digit string that identifies the event that you want to return.
			@return GetOrgEventApiRequest
	*/
	GetOrgEvent(ctx context.Context, orgId string, eventId string) GetOrgEventApiRequest
	/*
		GetOrgEvent Return One Event from One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgEventApiParams - Parameters for the request
		@return GetOrgEventApiRequest
	*/
	GetOrgEventWithParams(ctx context.Context, args *GetOrgEventApiParams) GetOrgEventApiRequest

	// Method available only for mocking purposes
	GetOrgEventExecute(r GetOrgEventApiRequest) (*EventViewForOrg, *http.Response, error)

	/*
		ListEventTypes Return All Event Types

		Returns a list of all event types, along with a description and additional metadata about each event.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListEventTypesApiRequest
	*/
	ListEventTypes(ctx context.Context) ListEventTypesApiRequest
	/*
		ListEventTypes Return All Event Types


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListEventTypesApiParams - Parameters for the request
		@return ListEventTypesApiRequest
	*/
	ListEventTypesWithParams(ctx context.Context, args *ListEventTypesApiParams) ListEventTypesApiRequest

	// Method available only for mocking purposes
	ListEventTypesExecute(r ListEventTypesApiRequest) (*PaginatedEventTypeDetailsResponse, *http.Response, error)

	/*
			ListGroupEvents Return Events from One Project

			Returns events for the specified project. Events identify significant database, billing, or security activities or status changes.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListGroupEventsApiRequest
	*/
	ListGroupEvents(ctx context.Context, groupId string) ListGroupEventsApiRequest
	/*
		ListGroupEvents Return Events from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupEventsApiParams - Parameters for the request
		@return ListGroupEventsApiRequest
	*/
	ListGroupEventsWithParams(ctx context.Context, args *ListGroupEventsApiParams) ListGroupEventsApiRequest

	// Method available only for mocking purposes
	ListGroupEventsExecute(r ListGroupEventsApiRequest) (*GroupPaginatedEvent, *http.Response, error)

	/*
			ListOrgEvents Return Events from One Organization

			Returns events for the specified organization. Events identify significant database, billing, or security activities or status changes.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return ListOrgEventsApiRequest
	*/
	ListOrgEvents(ctx context.Context, orgId string) ListOrgEventsApiRequest
	/*
		ListOrgEvents Return Events from One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgEventsApiParams - Parameters for the request
		@return ListOrgEventsApiRequest
	*/
	ListOrgEventsWithParams(ctx context.Context, args *ListOrgEventsApiParams) ListOrgEventsApiRequest

	// Method available only for mocking purposes
	ListOrgEventsExecute(r ListOrgEventsApiRequest) (*OrgPaginatedEvent, *http.Response, error)
}

// EventsApiService EventsApi service
type EventsApiService service

type GetGroupEventApiRequest struct {
	ctx        context.Context
	ApiService EventsApi
	groupId    string
	eventId    string
	includeRaw *bool
}

type GetGroupEventApiParams struct {
	GroupId    string
	EventId    string
	IncludeRaw *bool
}

func (a *EventsApiService) GetGroupEventWithParams(ctx context.Context, args *GetGroupEventApiParams) GetGroupEventApiRequest {
	return GetGroupEventApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		eventId:    args.EventId,
		includeRaw: args.IncludeRaw,
	}
}

// Flag that indicates whether to include the raw document in the output. The raw document contains additional meta information about the event.
func (r GetGroupEventApiRequest) IncludeRaw(includeRaw bool) GetGroupEventApiRequest {
	r.includeRaw = &includeRaw
	return r
}

func (r GetGroupEventApiRequest) Execute() (*EventViewForNdsGroup, *http.Response, error) {
	return r.ApiService.GetGroupEventExecute(r)
}

/*
GetGroupEvent Return One Event from One Project

Returns one event for the specified project. Events identify significant database, billing, or security activities or status changes. Use the Return Events from One Project endpoint to retrieve all events to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param eventId Unique 24-hexadecimal digit string that identifies the event that you want to return.
	@return GetGroupEventApiRequest
*/
func (a *EventsApiService) GetGroupEvent(ctx context.Context, groupId string, eventId string) GetGroupEventApiRequest {
	return GetGroupEventApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		eventId:    eventId,
	}
}

// GetGroupEventExecute executes the request
//
//	@return EventViewForNdsGroup
func (a *EventsApiService) GetGroupEventExecute(r GetGroupEventApiRequest) (*EventViewForNdsGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EventViewForNdsGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventsApiService.GetGroupEvent")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/events/{eventId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.eventId == "" {
		return localVarReturnValue, nil, reportError("eventId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"eventId"+"}", url.PathEscape(r.eventId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeRaw != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
	} else {
		var defaultValue bool = false
		r.includeRaw = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
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

type GetOrgEventApiRequest struct {
	ctx        context.Context
	ApiService EventsApi
	orgId      string
	eventId    string
	includeRaw *bool
}

type GetOrgEventApiParams struct {
	OrgId      string
	EventId    string
	IncludeRaw *bool
}

func (a *EventsApiService) GetOrgEventWithParams(ctx context.Context, args *GetOrgEventApiParams) GetOrgEventApiRequest {
	return GetOrgEventApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		eventId:    args.EventId,
		includeRaw: args.IncludeRaw,
	}
}

// Flag that indicates whether to include the raw document in the output. The raw document contains additional meta information about the event.
func (r GetOrgEventApiRequest) IncludeRaw(includeRaw bool) GetOrgEventApiRequest {
	r.includeRaw = &includeRaw
	return r
}

func (r GetOrgEventApiRequest) Execute() (*EventViewForOrg, *http.Response, error) {
	return r.ApiService.GetOrgEventExecute(r)
}

/*
GetOrgEvent Return One Event from One Organization

Returns one event for the specified organization. Events identify significant database, billing, or security activities or status changes. Use the Return Events from One Organization endpoint to retrieve all events to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param eventId Unique 24-hexadecimal digit string that identifies the event that you want to return.
	@return GetOrgEventApiRequest
*/
func (a *EventsApiService) GetOrgEvent(ctx context.Context, orgId string, eventId string) GetOrgEventApiRequest {
	return GetOrgEventApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		eventId:    eventId,
	}
}

// GetOrgEventExecute executes the request
//
//	@return EventViewForOrg
func (a *EventsApiService) GetOrgEventExecute(r GetOrgEventApiRequest) (*EventViewForOrg, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EventViewForOrg
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventsApiService.GetOrgEvent")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/events/{eventId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.eventId == "" {
		return localVarReturnValue, nil, reportError("eventId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"eventId"+"}", url.PathEscape(r.eventId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeRaw != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
	} else {
		var defaultValue bool = false
		r.includeRaw = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
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

type ListEventTypesApiRequest struct {
	ctx          context.Context
	ApiService   EventsApi
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListEventTypesApiParams struct {
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *EventsApiService) ListEventTypesWithParams(ctx context.Context, args *ListEventTypesApiParams) ListEventTypesApiRequest {
	return ListEventTypesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListEventTypesApiRequest) IncludeCount(includeCount bool) ListEventTypesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListEventTypesApiRequest) ItemsPerPage(itemsPerPage int) ListEventTypesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListEventTypesApiRequest) PageNum(pageNum int) ListEventTypesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListEventTypesApiRequest) Execute() (*PaginatedEventTypeDetailsResponse, *http.Response, error) {
	return r.ApiService.ListEventTypesExecute(r)
}

/*
ListEventTypes Return All Event Types

Returns a list of all event types, along with a description and additional metadata about each event.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListEventTypesApiRequest
*/
func (a *EventsApiService) ListEventTypes(ctx context.Context) ListEventTypesApiRequest {
	return ListEventTypesApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListEventTypesExecute executes the request
//
//	@return PaginatedEventTypeDetailsResponse
func (a *EventsApiService) ListEventTypesExecute(r ListEventTypesApiRequest) (*PaginatedEventTypeDetailsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedEventTypeDetailsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventsApiService.ListEventTypes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/eventTypes"

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

type ListGroupEventsApiRequest struct {
	ctx               context.Context
	ApiService        EventsApi
	groupId           string
	includeCount      *bool
	itemsPerPage      *int
	pageNum           *int
	clusterNames      *[]string
	eventType         *[]string
	excludedEventType *[]string
	includeRaw        *bool
	maxDate           *time.Time
	minDate           *time.Time
}

type ListGroupEventsApiParams struct {
	GroupId           string
	IncludeCount      *bool
	ItemsPerPage      *int
	PageNum           *int
	ClusterNames      *[]string
	EventType         *[]string
	ExcludedEventType *[]string
	IncludeRaw        *bool
	MaxDate           *time.Time
	MinDate           *time.Time
}

func (a *EventsApiService) ListGroupEventsWithParams(ctx context.Context, args *ListGroupEventsApiParams) ListGroupEventsApiRequest {
	return ListGroupEventsApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		includeCount:      args.IncludeCount,
		itemsPerPage:      args.ItemsPerPage,
		pageNum:           args.PageNum,
		clusterNames:      args.ClusterNames,
		eventType:         args.EventType,
		excludedEventType: args.ExcludedEventType,
		includeRaw:        args.IncludeRaw,
		maxDate:           args.MaxDate,
		minDate:           args.MinDate,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupEventsApiRequest) IncludeCount(includeCount bool) ListGroupEventsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupEventsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupEventsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupEventsApiRequest) PageNum(pageNum int) ListGroupEventsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Human-readable label that identifies the cluster.
func (r ListGroupEventsApiRequest) ClusterNames(clusterNames []string) ListGroupEventsApiRequest {
	r.clusterNames = &clusterNames
	return r
}

// Category of incident recorded at this moment in time.  **IMPORTANT**: The complete list of event type values changes frequently.
func (r ListGroupEventsApiRequest) EventType(eventType []string) ListGroupEventsApiRequest {
	r.eventType = &eventType
	return r
}

// Category of event that you would like to exclude from query results, such as &#x60;CLUSTER_CREATED&#x60;.  **IMPORTANT**: Event type names change frequently. Verify that you specify the event type correctly by checking the complete list of event types.
func (r ListGroupEventsApiRequest) ExcludedEventType(excludedEventType []string) ListGroupEventsApiRequest {
	r.excludedEventType = &excludedEventType
	return r
}

// Flag that indicates whether to include the raw document in the output. The raw document contains additional meta information about the event.
func (r ListGroupEventsApiRequest) IncludeRaw(includeRaw bool) ListGroupEventsApiRequest {
	r.includeRaw = &includeRaw
	return r
}

// Date and time from when MongoDB Cloud stops returning events. This parameter uses the ISO 8601 timestamp format in UTC.
func (r ListGroupEventsApiRequest) MaxDate(maxDate time.Time) ListGroupEventsApiRequest {
	r.maxDate = &maxDate
	return r
}

// Date and time from when MongoDB Cloud starts returning events. This parameter uses the ISO 8601 timestamp format in UTC.
func (r ListGroupEventsApiRequest) MinDate(minDate time.Time) ListGroupEventsApiRequest {
	r.minDate = &minDate
	return r
}

func (r ListGroupEventsApiRequest) Execute() (*GroupPaginatedEvent, *http.Response, error) {
	return r.ApiService.ListGroupEventsExecute(r)
}

/*
ListGroupEvents Return Events from One Project

Returns events for the specified project. Events identify significant database, billing, or security activities or status changes.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupEventsApiRequest
*/
func (a *EventsApiService) ListGroupEvents(ctx context.Context, groupId string) ListGroupEventsApiRequest {
	return ListGroupEventsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupEventsExecute executes the request
//
//	@return GroupPaginatedEvent
func (a *EventsApiService) ListGroupEventsExecute(r ListGroupEventsApiRequest) (*GroupPaginatedEvent, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupPaginatedEvent
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventsApiService.ListGroupEvents")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/events"
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
	if r.clusterNames != nil {
		t := *r.clusterNames
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "clusterNames", t, "multi")

	}
	if r.eventType != nil {
		t := *r.eventType
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "eventType", t, "multi")

	}
	if r.excludedEventType != nil {
		t := *r.excludedEventType
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "excludedEventType", t, "multi")

	}
	if r.includeRaw != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
	} else {
		var defaultValue bool = false
		r.includeRaw = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
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

type ListOrgEventsApiRequest struct {
	ctx          context.Context
	ApiService   EventsApi
	orgId        string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	eventType    *[]string
	includeRaw   *bool
	maxDate      *time.Time
	minDate      *time.Time
}

type ListOrgEventsApiParams struct {
	OrgId        string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	EventType    *[]string
	IncludeRaw   *bool
	MaxDate      *time.Time
	MinDate      *time.Time
}

func (a *EventsApiService) ListOrgEventsWithParams(ctx context.Context, args *ListOrgEventsApiParams) ListOrgEventsApiRequest {
	return ListOrgEventsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		eventType:    args.EventType,
		includeRaw:   args.IncludeRaw,
		maxDate:      args.MaxDate,
		minDate:      args.MinDate,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgEventsApiRequest) IncludeCount(includeCount bool) ListOrgEventsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgEventsApiRequest) ItemsPerPage(itemsPerPage int) ListOrgEventsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgEventsApiRequest) PageNum(pageNum int) ListOrgEventsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Category of incident recorded at this moment in time.  **IMPORTANT**: The complete list of event type values changes frequently.
func (r ListOrgEventsApiRequest) EventType(eventType []string) ListOrgEventsApiRequest {
	r.eventType = &eventType
	return r
}

// Flag that indicates whether to include the raw document in the output. The raw document contains additional meta information about the event.
func (r ListOrgEventsApiRequest) IncludeRaw(includeRaw bool) ListOrgEventsApiRequest {
	r.includeRaw = &includeRaw
	return r
}

// Date and time from when MongoDB Cloud stops returning events. This parameter uses the ISO 8601 timestamp format in UTC.
func (r ListOrgEventsApiRequest) MaxDate(maxDate time.Time) ListOrgEventsApiRequest {
	r.maxDate = &maxDate
	return r
}

// Date and time from when MongoDB Cloud starts returning events. This parameter uses the ISO 8601 timestamp format in UTC.
func (r ListOrgEventsApiRequest) MinDate(minDate time.Time) ListOrgEventsApiRequest {
	r.minDate = &minDate
	return r
}

func (r ListOrgEventsApiRequest) Execute() (*OrgPaginatedEvent, *http.Response, error) {
	return r.ApiService.ListOrgEventsExecute(r)
}

/*
ListOrgEvents Return Events from One Organization

Returns events for the specified organization. Events identify significant database, billing, or security activities or status changes.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgEventsApiRequest
*/
func (a *EventsApiService) ListOrgEvents(ctx context.Context, orgId string) ListOrgEventsApiRequest {
	return ListOrgEventsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgEventsExecute executes the request
//
//	@return OrgPaginatedEvent
func (a *EventsApiService) ListOrgEventsExecute(r ListOrgEventsApiRequest) (*OrgPaginatedEvent, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrgPaginatedEvent
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EventsApiService.ListOrgEvents")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/events"
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
	if r.eventType != nil {
		t := *r.eventType
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "eventType", t, "multi")

	}
	if r.includeRaw != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
	} else {
		var defaultValue bool = false
		r.includeRaw = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeRaw", r.includeRaw, "")
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
