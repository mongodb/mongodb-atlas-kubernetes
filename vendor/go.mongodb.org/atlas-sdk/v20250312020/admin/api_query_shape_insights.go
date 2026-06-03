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

type QueryShapeInsightsApi interface {

	/*
		GetClusterQueryShape Return One Query Shape

		Returns the details for a single query shape. This endpoint only returns query shapes with REJECTED status. If the specified query shape hash does not correspond to a rejected query shape, a 404 Not Found error is returned.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
		@return GetClusterQueryShapeApiRequest
	*/
	GetClusterQueryShape(ctx context.Context, groupId string, clusterName string, queryShapeHash string) GetClusterQueryShapeApiRequest
	/*
		GetClusterQueryShape Return One Query Shape


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterQueryShapeApiParams - Parameters for the request
		@return GetClusterQueryShapeApiRequest
	*/
	GetClusterQueryShapeWithParams(ctx context.Context, args *GetClusterQueryShapeApiParams) GetClusterQueryShapeApiRequest

	// Method available only for mocking purposes
	GetClusterQueryShapeExecute(r GetClusterQueryShapeApiRequest) (*QueryShapeResponse, *http.Response, error)

	/*
		GetQueryShapeDetails Return Query Shape Details

		Returns the metadata and statistics summary for a given query shape hash.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
		@return GetQueryShapeDetailsApiRequest
	*/
	GetQueryShapeDetails(ctx context.Context, groupId string, clusterName string, queryShapeHash string) GetQueryShapeDetailsApiRequest
	/*
		GetQueryShapeDetails Return Query Shape Details


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetQueryShapeDetailsApiParams - Parameters for the request
		@return GetQueryShapeDetailsApiRequest
	*/
	GetQueryShapeDetailsWithParams(ctx context.Context, args *GetQueryShapeDetailsApiParams) GetQueryShapeDetailsApiRequest

	// Method available only for mocking purposes
	GetQueryShapeDetailsExecute(r GetQueryShapeDetailsApiRequest) (*QueryStatsDetailsResponse, *http.Response, error)

	/*
		ListClusterQueryShapes Return All Query Shapes

		Returns a list of query shapes for one cluster. Query shapes may be filtered by their status; at present, this endpoint supports only the REJECTED status.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListClusterQueryShapesApiRequest
	*/
	ListClusterQueryShapes(ctx context.Context, groupId string, clusterName string) ListClusterQueryShapesApiRequest
	/*
		ListClusterQueryShapes Return All Query Shapes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterQueryShapesApiParams - Parameters for the request
		@return ListClusterQueryShapesApiRequest
	*/
	ListClusterQueryShapesWithParams(ctx context.Context, args *ListClusterQueryShapesApiParams) ListClusterQueryShapesApiRequest

	// Method available only for mocking purposes
	ListClusterQueryShapesExecute(r ListClusterQueryShapesApiRequest) (*PaginatedQueryShapes, *http.Response, error)

	/*
		ListQueryShapeSummaries Return Query Statistic Summaries

		Returns a list of query shape statistics summaries for a given cluster. Query shape statistics provide performance insights about MongoDB queries, helping users identify problematic query patterns and potential optimizations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListQueryShapeSummariesApiRequest
	*/
	ListQueryShapeSummaries(ctx context.Context, groupId string, clusterName string) ListQueryShapeSummariesApiRequest
	/*
		ListQueryShapeSummaries Return Query Statistic Summaries


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListQueryShapeSummariesApiParams - Parameters for the request
		@return ListQueryShapeSummariesApiRequest
	*/
	ListQueryShapeSummariesWithParams(ctx context.Context, args *ListQueryShapeSummariesApiParams) ListQueryShapeSummariesApiRequest

	// Method available only for mocking purposes
	ListQueryShapeSummariesExecute(r ListQueryShapeSummariesApiRequest) (*QueryStatsSummaryListResponse, *http.Response, error)

	/*
		UpdateClusterQueryShape Update Query Shape Rejection Status

		Updates the rejection status of a query shape. Use this endpoint to reject a query shape (preventing it from executing on the cluster) or to unreject a previously rejected query shape (allowing it to execute again). This operation is idempotent: rejecting an already rejected query shape or unrejecting an already unrejected query shape will return 200 OK.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
		@param queryShapeUpdateRequest The desired rejection status for the query shape. Provide REJECTED to block the query shape from executing, or UNREJECTED to allow it to execute.
		@return UpdateClusterQueryShapeApiRequest
	*/
	UpdateClusterQueryShape(ctx context.Context, groupId string, clusterName string, queryShapeHash string, queryShapeUpdateRequest *QueryShapeUpdateRequest) UpdateClusterQueryShapeApiRequest
	/*
		UpdateClusterQueryShape Update Query Shape Rejection Status


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterQueryShapeApiParams - Parameters for the request
		@return UpdateClusterQueryShapeApiRequest
	*/
	UpdateClusterQueryShapeWithParams(ctx context.Context, args *UpdateClusterQueryShapeApiParams) UpdateClusterQueryShapeApiRequest

	// Method available only for mocking purposes
	UpdateClusterQueryShapeExecute(r UpdateClusterQueryShapeApiRequest) (*QueryShapeResponse, *http.Response, error)
}

// QueryShapeInsightsApiService QueryShapeInsightsApi service
type QueryShapeInsightsApiService service

type GetClusterQueryShapeApiRequest struct {
	ctx            context.Context
	ApiService     QueryShapeInsightsApi
	groupId        string
	clusterName    string
	queryShapeHash string
}

type GetClusterQueryShapeApiParams struct {
	GroupId        string
	ClusterName    string
	QueryShapeHash string
}

func (a *QueryShapeInsightsApiService) GetClusterQueryShapeWithParams(ctx context.Context, args *GetClusterQueryShapeApiParams) GetClusterQueryShapeApiRequest {
	return GetClusterQueryShapeApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		queryShapeHash: args.QueryShapeHash,
	}
}

func (r GetClusterQueryShapeApiRequest) Execute() (*QueryShapeResponse, *http.Response, error) {
	return r.ApiService.GetClusterQueryShapeExecute(r)
}

/*
GetClusterQueryShape Return One Query Shape

Returns the details for a single query shape. This endpoint only returns query shapes with REJECTED status. If the specified query shape hash does not correspond to a rejected query shape, a 404 Not Found error is returned.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
	@return GetClusterQueryShapeApiRequest
*/
func (a *QueryShapeInsightsApiService) GetClusterQueryShape(ctx context.Context, groupId string, clusterName string, queryShapeHash string) GetClusterQueryShapeApiRequest {
	return GetClusterQueryShapeApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		queryShapeHash: queryShapeHash,
	}
}

// GetClusterQueryShapeExecute executes the request
//
//	@return QueryShapeResponse
func (a *QueryShapeInsightsApiService) GetClusterQueryShapeExecute(r GetClusterQueryShapeApiRequest) (*QueryShapeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *QueryShapeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QueryShapeInsightsApiService.GetClusterQueryShape")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/queryShapes/{queryShapeHash}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.queryShapeHash == "" {
		return localVarReturnValue, nil, reportError("queryShapeHash is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"queryShapeHash"+"}", url.PathEscape(r.queryShapeHash), -1)

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

type GetQueryShapeDetailsApiRequest struct {
	ctx            context.Context
	ApiService     QueryShapeInsightsApi
	groupId        string
	clusterName    string
	queryShapeHash string
	since          *int64
	until          *int64
	processIds     *[]string
}

type GetQueryShapeDetailsApiParams struct {
	GroupId        string
	ClusterName    string
	QueryShapeHash string
	Since          *int64
	Until          *int64
	ProcessIds     *[]string
}

func (a *QueryShapeInsightsApiService) GetQueryShapeDetailsWithParams(ctx context.Context, args *GetQueryShapeDetailsApiParams) GetQueryShapeDetailsApiRequest {
	return GetQueryShapeDetailsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		queryShapeHash: args.QueryShapeHash,
		since:          args.Since,
		until:          args.Until,
		processIds:     args.ProcessIds,
	}
}

// Date and time from which to retrieve query shape statistics. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **until** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r GetQueryShapeDetailsApiRequest) Since(since int64) GetQueryShapeDetailsApiRequest {
	r.since = &since
	return r
}

// Date and time up until which to retrieve query shape statistics. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you specify the **until** parameter, you must specify the **since** parameter. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r GetQueryShapeDetailsApiRequest) Until(until int64) GetQueryShapeDetailsApiRequest {
	r.until = &until
	return r
}

// Process IDs from which to retrieve query shape statistics. A &#x60;processId&#x60; is a combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (&#x60;mongod&#x60; or &#x60;mongos&#x60;). The port must be the IANA port on which the MongoDB process listens for requests. To include multiple &#x60;processIds&#x60;, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each &#x60;processId&#x60;.
func (r GetQueryShapeDetailsApiRequest) ProcessIds(processIds []string) GetQueryShapeDetailsApiRequest {
	r.processIds = &processIds
	return r
}

func (r GetQueryShapeDetailsApiRequest) Execute() (*QueryStatsDetailsResponse, *http.Response, error) {
	return r.ApiService.GetQueryShapeDetailsExecute(r)
}

/*
GetQueryShapeDetails Return Query Shape Details

Returns the metadata and statistics summary for a given query shape hash.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
	@return GetQueryShapeDetailsApiRequest
*/
func (a *QueryShapeInsightsApiService) GetQueryShapeDetails(ctx context.Context, groupId string, clusterName string, queryShapeHash string) GetQueryShapeDetailsApiRequest {
	return GetQueryShapeDetailsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		queryShapeHash: queryShapeHash,
	}
}

// GetQueryShapeDetailsExecute executes the request
//
//	@return QueryStatsDetailsResponse
func (a *QueryShapeInsightsApiService) GetQueryShapeDetailsExecute(r GetQueryShapeDetailsApiRequest) (*QueryStatsDetailsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *QueryStatsDetailsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QueryShapeInsightsApiService.GetQueryShapeDetails")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/queryShapeInsights/{queryShapeHash}/details"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.queryShapeHash == "" {
		return localVarReturnValue, nil, reportError("queryShapeHash is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"queryShapeHash"+"}", url.PathEscape(r.queryShapeHash), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
	}
	if r.until != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "until", r.until, "")
	}
	if r.processIds != nil {
		t := *r.processIds
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "processIds", t, "multi")

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

type ListClusterQueryShapesApiRequest struct {
	ctx          context.Context
	ApiService   QueryShapeInsightsApi
	groupId      string
	clusterName  string
	status       *string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListClusterQueryShapesApiParams struct {
	GroupId      string
	ClusterName  string
	Status       *string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *QueryShapeInsightsApiService) ListClusterQueryShapesWithParams(ctx context.Context, args *ListClusterQueryShapesApiParams) ListClusterQueryShapesApiRequest {
	return ListClusterQueryShapesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		status:       args.Status,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// The status of query shapes to retrieve. Only REJECTED status is supported. If omitted, defaults to REJECTED.
func (r ListClusterQueryShapesApiRequest) Status(status string) ListClusterQueryShapesApiRequest {
	r.status = &status
	return r
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClusterQueryShapesApiRequest) IncludeCount(includeCount bool) ListClusterQueryShapesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterQueryShapesApiRequest) ItemsPerPage(itemsPerPage int) ListClusterQueryShapesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterQueryShapesApiRequest) PageNum(pageNum int) ListClusterQueryShapesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListClusterQueryShapesApiRequest) Execute() (*PaginatedQueryShapes, *http.Response, error) {
	return r.ApiService.ListClusterQueryShapesExecute(r)
}

/*
ListClusterQueryShapes Return All Query Shapes

Returns a list of query shapes for one cluster. Query shapes may be filtered by their status; at present, this endpoint supports only the REJECTED status.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListClusterQueryShapesApiRequest
*/
func (a *QueryShapeInsightsApiService) ListClusterQueryShapes(ctx context.Context, groupId string, clusterName string) ListClusterQueryShapesApiRequest {
	return ListClusterQueryShapesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterQueryShapesExecute executes the request
//
//	@return PaginatedQueryShapes
func (a *QueryShapeInsightsApiService) ListClusterQueryShapesExecute(r ListClusterQueryShapesApiRequest) (*PaginatedQueryShapes, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedQueryShapes
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QueryShapeInsightsApiService.ListClusterQueryShapes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/queryShapes"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.status != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "status", r.status, "")
	} else {
		var defaultValue string = "REJECTED"
		r.status = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "status", r.status, "")
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

type ListQueryShapeSummariesApiRequest struct {
	ctx              context.Context
	ApiService       QueryShapeInsightsApi
	groupId          string
	clusterName      string
	since            *int64
	until            *int64
	processIds       *[]string
	namespaces       *[]string
	commands         *[]string
	nSummaries       *int64
	series           *[]string
	queryShapeHashes *[]string
}

type ListQueryShapeSummariesApiParams struct {
	GroupId          string
	ClusterName      string
	Since            *int64
	Until            *int64
	ProcessIds       *[]string
	Namespaces       *[]string
	Commands         *[]string
	NSummaries       *int64
	Series           *[]string
	QueryShapeHashes *[]string
}

func (a *QueryShapeInsightsApiService) ListQueryShapeSummariesWithParams(ctx context.Context, args *ListQueryShapeSummariesApiParams) ListQueryShapeSummariesApiRequest {
	return ListQueryShapeSummariesApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		clusterName:      args.ClusterName,
		since:            args.Since,
		until:            args.Until,
		processIds:       args.ProcessIds,
		namespaces:       args.Namespaces,
		commands:         args.Commands,
		nSummaries:       args.NSummaries,
		series:           args.Series,
		queryShapeHashes: args.QueryShapeHashes,
	}
}

// Date and time from which to retrieve query shape statistics. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **until** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r ListQueryShapeSummariesApiRequest) Since(since int64) ListQueryShapeSummariesApiRequest {
	r.since = &since
	return r
}

// Date and time up until which to retrieve query shape statistics. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you specify the **until** parameter, you must specify the **since** parameter. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r ListQueryShapeSummariesApiRequest) Until(until int64) ListQueryShapeSummariesApiRequest {
	r.until = &until
	return r
}

// Process IDs from which to retrieve query shape statistics. A &#x60;processId&#x60; is a combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (&#x60;mongod&#x60; or &#x60;mongos&#x60;). The port must be the IANA port on which the MongoDB process listens for requests. To include multiple &#x60;processId&#x60;, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each &#x60;processId&#x60;.
func (r ListQueryShapeSummariesApiRequest) ProcessIds(processIds []string) ListQueryShapeSummariesApiRequest {
	r.processIds = &processIds
	return r
}

// Namespaces from which to retrieve query shape statistics. A namespace consists of one database and one collection resource written as &#x60;.&#x60;: &#x60;&lt;database&gt;.&lt;collection&gt;&#x60;. To include multiple namespaces, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each namespace. Omit this parameter to return results for all namespaces.
func (r ListQueryShapeSummariesApiRequest) Namespaces(namespaces []string) ListQueryShapeSummariesApiRequest {
	r.namespaces = &namespaces
	return r
}

// Retrieve query shape statistics matching specified MongoDB commands. To include multiple commands, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each command. The currently supported parameters are find, distinct, and aggregate. Omit this parameter to return results for all supported commands.
func (r ListQueryShapeSummariesApiRequest) Commands(commands []string) ListQueryShapeSummariesApiRequest {
	r.commands = &commands
	return r
}

// Maximum number of query statistic summaries to return.
func (r ListQueryShapeSummariesApiRequest) NSummaries(nSummaries int64) ListQueryShapeSummariesApiRequest {
	r.nSummaries = &nSummaries
	return r
}

// Query shape statistics data series to retrieve. A series represents a specific metric about query execution. To include multiple series, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each series. Omit this parameter to return results for all available series.
func (r ListQueryShapeSummariesApiRequest) Series(series []string) ListQueryShapeSummariesApiRequest {
	r.series = &series
	return r
}

// A list of SHA256 hashes of desired query shapes, output by MongoDB commands like &#x60;$queryStats&#x60; and $explain or slow query logs. To include multiple series, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each series. Omit this parameter to return results for all available series.
func (r ListQueryShapeSummariesApiRequest) QueryShapeHashes(queryShapeHashes []string) ListQueryShapeSummariesApiRequest {
	r.queryShapeHashes = &queryShapeHashes
	return r
}

func (r ListQueryShapeSummariesApiRequest) Execute() (*QueryStatsSummaryListResponse, *http.Response, error) {
	return r.ApiService.ListQueryShapeSummariesExecute(r)
}

/*
ListQueryShapeSummaries Return Query Statistic Summaries

Returns a list of query shape statistics summaries for a given cluster. Query shape statistics provide performance insights about MongoDB queries, helping users identify problematic query patterns and potential optimizations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListQueryShapeSummariesApiRequest
*/
func (a *QueryShapeInsightsApiService) ListQueryShapeSummaries(ctx context.Context, groupId string, clusterName string) ListQueryShapeSummariesApiRequest {
	return ListQueryShapeSummariesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListQueryShapeSummariesExecute executes the request
//
//	@return QueryStatsSummaryListResponse
func (a *QueryShapeInsightsApiService) ListQueryShapeSummariesExecute(r ListQueryShapeSummariesApiRequest) (*QueryStatsSummaryListResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *QueryStatsSummaryListResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QueryShapeInsightsApiService.ListQueryShapeSummaries")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/queryShapeInsights/summaries"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
	}
	if r.until != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "until", r.until, "")
	}
	if r.processIds != nil {
		t := *r.processIds
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "processIds", t, "multi")

	}
	if r.namespaces != nil {
		t := *r.namespaces
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "namespaces", t, "multi")

	}
	if r.commands != nil {
		t := *r.commands
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "commands", t, "multi")

	}
	if r.nSummaries != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nSummaries", r.nSummaries, "")
	} else {
		var defaultValue int64 = 100
		r.nSummaries = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "nSummaries", r.nSummaries, "")
	}
	if r.series != nil {
		t := *r.series
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "series", t, "multi")

	}
	if r.queryShapeHashes != nil {
		t := *r.queryShapeHashes
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "queryShapeHashes", t, "multi")

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

type UpdateClusterQueryShapeApiRequest struct {
	ctx                     context.Context
	ApiService              QueryShapeInsightsApi
	groupId                 string
	clusterName             string
	queryShapeHash          string
	queryShapeUpdateRequest *QueryShapeUpdateRequest
}

type UpdateClusterQueryShapeApiParams struct {
	GroupId                 string
	ClusterName             string
	QueryShapeHash          string
	QueryShapeUpdateRequest *QueryShapeUpdateRequest
}

func (a *QueryShapeInsightsApiService) UpdateClusterQueryShapeWithParams(ctx context.Context, args *UpdateClusterQueryShapeApiParams) UpdateClusterQueryShapeApiRequest {
	return UpdateClusterQueryShapeApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 args.GroupId,
		clusterName:             args.ClusterName,
		queryShapeHash:          args.QueryShapeHash,
		queryShapeUpdateRequest: args.QueryShapeUpdateRequest,
	}
}

func (r UpdateClusterQueryShapeApiRequest) Execute() (*QueryShapeResponse, *http.Response, error) {
	return r.ApiService.UpdateClusterQueryShapeExecute(r)
}

/*
UpdateClusterQueryShape Update Query Shape Rejection Status

Updates the rejection status of a query shape. Use this endpoint to reject a query shape (preventing it from executing on the cluster) or to unreject a previously rejected query shape (allowing it to execute again). This operation is idempotent: rejecting an already rejected query shape or unrejecting an already unrejected query shape will return 200 OK.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param queryShapeHash A SHA256 hash of a query shape, output by MongoDB commands like `$queryStats` and `$explain` or slow query logs.
	@return UpdateClusterQueryShapeApiRequest
*/
func (a *QueryShapeInsightsApiService) UpdateClusterQueryShape(ctx context.Context, groupId string, clusterName string, queryShapeHash string, queryShapeUpdateRequest *QueryShapeUpdateRequest) UpdateClusterQueryShapeApiRequest {
	return UpdateClusterQueryShapeApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 groupId,
		clusterName:             clusterName,
		queryShapeHash:          queryShapeHash,
		queryShapeUpdateRequest: queryShapeUpdateRequest,
	}
}

// UpdateClusterQueryShapeExecute executes the request
//
//	@return QueryShapeResponse
func (a *QueryShapeInsightsApiService) UpdateClusterQueryShapeExecute(r UpdateClusterQueryShapeApiRequest) (*QueryShapeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *QueryShapeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QueryShapeInsightsApiService.UpdateClusterQueryShape")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/queryShapes/{queryShapeHash}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.queryShapeHash == "" {
		return localVarReturnValue, nil, reportError("queryShapeHash is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"queryShapeHash"+"}", url.PathEscape(r.queryShapeHash), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.queryShapeUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("queryShapeUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-03-12+json"}

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
	// body params
	localVarPostBody = r.queryShapeUpdateRequest
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
