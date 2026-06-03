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

type CollectionLevelMetricsApi interface {

	/*
		GetClusterNamespaces Return Ranked Namespaces from One Cluster

		Return the subset of namespaces from the given cluster sorted by highest total execution time (descending) within the given time window.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
		@param clusterView Human-readable label that identifies the cluster topology to retrieve metrics for.
		@return GetClusterNamespacesApiRequest
	*/
	GetClusterNamespaces(ctx context.Context, groupId string, clusterName string, clusterView string) GetClusterNamespacesApiRequest
	/*
		GetClusterNamespaces Return Ranked Namespaces from One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterNamespacesApiParams - Parameters for the request
		@return GetClusterNamespacesApiRequest
	*/
	GetClusterNamespacesWithParams(ctx context.Context, args *GetClusterNamespacesApiParams) GetClusterNamespacesApiRequest

	// Method available only for mocking purposes
	GetClusterNamespacesExecute(r GetClusterNamespacesApiRequest) (*CollStatsRankedNamespaces, *http.Response, error)

	/*
		GetProcessNamespaces Return Ranked Namespaces from One Host

		Return the subset of namespaces from the given process ranked by highest total execution time (descending) within the given time window.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@return GetProcessNamespacesApiRequest
	*/
	GetProcessNamespaces(ctx context.Context, groupId string, processId string) GetProcessNamespacesApiRequest
	/*
		GetProcessNamespaces Return Ranked Namespaces from One Host


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetProcessNamespacesApiParams - Parameters for the request
		@return GetProcessNamespacesApiRequest
	*/
	GetProcessNamespacesWithParams(ctx context.Context, args *GetProcessNamespacesApiParams) GetProcessNamespacesApiRequest

	// Method available only for mocking purposes
	GetProcessNamespacesExecute(r GetProcessNamespacesApiRequest) (*CollStatsRankedNamespaces, *http.Response, error)

	/*
		ListCollStatMeasurements Return Cluster-Level Query Latency

		Get a list of the Coll Stats Latency cluster-level measurements for the given namespace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to retrieve metrics for.
		@param clusterView Human-readable label that identifies the cluster topology to retrieve metrics for.
		@param databaseName Human-readable label that identifies the database.
		@param collectionName Human-readable label that identifies the collection.
		@return ListCollStatMeasurementsApiRequest
	*/
	ListCollStatMeasurements(ctx context.Context, groupId string, clusterName string, clusterView string, databaseName string, collectionName string) ListCollStatMeasurementsApiRequest
	/*
		ListCollStatMeasurements Return Cluster-Level Query Latency


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListCollStatMeasurementsApiParams - Parameters for the request
		@return ListCollStatMeasurementsApiRequest
	*/
	ListCollStatMeasurementsWithParams(ctx context.Context, args *ListCollStatMeasurementsApiParams) ListCollStatMeasurementsApiRequest

	// Method available only for mocking purposes
	ListCollStatMeasurementsExecute(r ListCollStatMeasurementsApiRequest) (*MeasurementsCollStatsLatencyCluster, *http.Response, error)

	/*
		ListCollStatMetrics Return All Metric Names

		Returns all available Coll Stats Latency metric names and their respective units for the specified project at the time of request.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListCollStatMetricsApiRequest
	*/
	ListCollStatMetrics(ctx context.Context, groupId string) ListCollStatMetricsApiRequest
	/*
		ListCollStatMetrics Return All Metric Names


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListCollStatMetricsApiParams - Parameters for the request
		@return ListCollStatMetricsApiRequest
	*/
	ListCollStatMetricsWithParams(ctx context.Context, args *ListCollStatMetricsApiParams) ListCollStatMetricsApiRequest

	// Method available only for mocking purposes
	ListCollStatMetricsExecute(r ListCollStatMetricsApiRequest) (*CollStatsLatencyNamespaceMetrics, *http.Response, error)

	/*
		ListPinnedNamespaces Return Pinned Namespaces

		Returns a list of given cluster's pinned namespaces, a set of namespaces manually selected by users to collect query latency metrics on.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to retrieve pinned namespaces for.
		@return ListPinnedNamespacesApiRequest
	*/
	ListPinnedNamespaces(ctx context.Context, groupId string, clusterName string) ListPinnedNamespacesApiRequest
	/*
		ListPinnedNamespaces Return Pinned Namespaces


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPinnedNamespacesApiParams - Parameters for the request
		@return ListPinnedNamespacesApiRequest
	*/
	ListPinnedNamespacesWithParams(ctx context.Context, args *ListPinnedNamespacesApiParams) ListPinnedNamespacesApiRequest

	// Method available only for mocking purposes
	ListPinnedNamespacesExecute(r ListPinnedNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error)

	/*
		ListProcessMeasurements Return Host-Level Query Latency

		Get a list of the Coll Stats Latency process-level measurements for the given namespace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@param databaseName Human-readable label that identifies the database.
		@param collectionName Human-readable label that identifies the collection.
		@return ListProcessMeasurementsApiRequest
	*/
	ListProcessMeasurements(ctx context.Context, groupId string, processId string, databaseName string, collectionName string) ListProcessMeasurementsApiRequest
	/*
		ListProcessMeasurements Return Host-Level Query Latency


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListProcessMeasurementsApiParams - Parameters for the request
		@return ListProcessMeasurementsApiRequest
	*/
	ListProcessMeasurementsWithParams(ctx context.Context, args *ListProcessMeasurementsApiParams) ListProcessMeasurementsApiRequest

	// Method available only for mocking purposes
	ListProcessMeasurementsExecute(r ListProcessMeasurementsApiRequest) (*MeasurementsCollStatsLatencyHost, *http.Response, error)

	/*
		PinNamespaces Pin Namespaces

		Pin provided list of namespaces for collection-level latency metrics collection for the given Group and Cluster. This initializes a pinned namespaces list or replaces any existing pinned namespaces list for the Group and Cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
		@param namespacesRequest List of namespace strings (combination of database and collection name) to pin for query latency metric collection.
		@return PinNamespacesApiRequest
	*/
	PinNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) PinNamespacesApiRequest
	/*
		PinNamespaces Pin Namespaces


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param PinNamespacesApiParams - Parameters for the request
		@return PinNamespacesApiRequest
	*/
	PinNamespacesWithParams(ctx context.Context, args *PinNamespacesApiParams) PinNamespacesApiRequest

	// Method available only for mocking purposes
	PinNamespacesExecute(r PinNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error)

	/*
		UnpinNamespaces Unpin Namespaces

		Unpin provided list of namespaces for collection-level latency metrics collection for the given Group and Cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to unpin namespaces from.
		@param namespacesRequest List of namespace strings (combination of database and collection name) to pin for query latency metric collection.
		@return UnpinNamespacesApiRequest
	*/
	UnpinNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) UnpinNamespacesApiRequest
	/*
		UnpinNamespaces Unpin Namespaces


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UnpinNamespacesApiParams - Parameters for the request
		@return UnpinNamespacesApiRequest
	*/
	UnpinNamespacesWithParams(ctx context.Context, args *UnpinNamespacesApiParams) UnpinNamespacesApiRequest

	// Method available only for mocking purposes
	UnpinNamespacesExecute(r UnpinNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error)

	/*
		UpdatePinnedNamespaces Add Pinned Namespaces

		Add provided list of namespaces to existing pinned namespaces list for collection-level latency metrics collection for the given Group and Cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
		@param namespacesRequest List of namespace strings (combination of database and collection name) to pin for query latency metric collection.
		@return UpdatePinnedNamespacesApiRequest
	*/
	UpdatePinnedNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) UpdatePinnedNamespacesApiRequest
	/*
		UpdatePinnedNamespaces Add Pinned Namespaces


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdatePinnedNamespacesApiParams - Parameters for the request
		@return UpdatePinnedNamespacesApiRequest
	*/
	UpdatePinnedNamespacesWithParams(ctx context.Context, args *UpdatePinnedNamespacesApiParams) UpdatePinnedNamespacesApiRequest

	// Method available only for mocking purposes
	UpdatePinnedNamespacesExecute(r UpdatePinnedNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error)
}

// CollectionLevelMetricsApiService CollectionLevelMetricsApi service
type CollectionLevelMetricsApiService service

type GetClusterNamespacesApiRequest struct {
	ctx         context.Context
	ApiService  CollectionLevelMetricsApi
	groupId     string
	clusterName string
	clusterView string
	start       *time.Time
	end         *time.Time
	period      *string
}

type GetClusterNamespacesApiParams struct {
	GroupId     string
	ClusterName string
	ClusterView string
	Start       *time.Time
	End         *time.Time
	Period      *string
}

func (a *CollectionLevelMetricsApiService) GetClusterNamespacesWithParams(ctx context.Context, args *GetClusterNamespacesApiParams) GetClusterNamespacesApiRequest {
	return GetClusterNamespacesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		clusterView: args.ClusterView,
		start:       args.Start,
		end:         args.End,
		period:      args.Period,
	}
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetClusterNamespacesApiRequest) Start(start time.Time) GetClusterNamespacesApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetClusterNamespacesApiRequest) End(end time.Time) GetClusterNamespacesApiRequest {
	r.end = &end
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetClusterNamespacesApiRequest) Period(period string) GetClusterNamespacesApiRequest {
	r.period = &period
	return r
}

func (r GetClusterNamespacesApiRequest) Execute() (*CollStatsRankedNamespaces, *http.Response, error) {
	return r.ApiService.GetClusterNamespacesExecute(r)
}

/*
GetClusterNamespaces Return Ranked Namespaces from One Cluster

Return the subset of namespaces from the given cluster sorted by highest total execution time (descending) within the given time window.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
	@param clusterView Human-readable label that identifies the cluster topology to retrieve metrics for.
	@return GetClusterNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) GetClusterNamespaces(ctx context.Context, groupId string, clusterName string, clusterView string) GetClusterNamespacesApiRequest {
	return GetClusterNamespacesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		clusterView: clusterView,
	}
}

// GetClusterNamespacesExecute executes the request
//
//	@return CollStatsRankedNamespaces
func (a *CollectionLevelMetricsApiService) GetClusterNamespacesExecute(r GetClusterNamespacesApiRequest) (*CollStatsRankedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CollStatsRankedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.GetClusterNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/{clusterView}/collStats/namespaces"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.clusterView == "" {
		return localVarReturnValue, nil, reportError("clusterView is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterView"+"}", url.PathEscape(r.clusterView), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type GetProcessNamespacesApiRequest struct {
	ctx        context.Context
	ApiService CollectionLevelMetricsApi
	groupId    string
	processId  string
	start      *time.Time
	end        *time.Time
	period     *string
}

type GetProcessNamespacesApiParams struct {
	GroupId   string
	ProcessId string
	Start     *time.Time
	End       *time.Time
	Period    *string
}

func (a *CollectionLevelMetricsApiService) GetProcessNamespacesWithParams(ctx context.Context, args *GetProcessNamespacesApiParams) GetProcessNamespacesApiRequest {
	return GetProcessNamespacesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		processId:  args.ProcessId,
		start:      args.Start,
		end:        args.End,
		period:     args.Period,
	}
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessNamespacesApiRequest) Start(start time.Time) GetProcessNamespacesApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessNamespacesApiRequest) End(end time.Time) GetProcessNamespacesApiRequest {
	r.end = &end
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetProcessNamespacesApiRequest) Period(period string) GetProcessNamespacesApiRequest {
	r.period = &period
	return r
}

func (r GetProcessNamespacesApiRequest) Execute() (*CollStatsRankedNamespaces, *http.Response, error) {
	return r.ApiService.GetProcessNamespacesExecute(r)
}

/*
GetProcessNamespaces Return Ranked Namespaces from One Host

Return the subset of namespaces from the given process ranked by highest total execution time (descending) within the given time window.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetProcessNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) GetProcessNamespaces(ctx context.Context, groupId string, processId string) GetProcessNamespacesApiRequest {
	return GetProcessNamespacesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// GetProcessNamespacesExecute executes the request
//
//	@return CollStatsRankedNamespaces
func (a *CollectionLevelMetricsApiService) GetProcessNamespacesExecute(r GetProcessNamespacesApiRequest) (*CollStatsRankedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CollStatsRankedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.GetProcessNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/collStats/namespaces"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type ListCollStatMeasurementsApiRequest struct {
	ctx            context.Context
	ApiService     CollectionLevelMetricsApi
	groupId        string
	clusterName    string
	clusterView    string
	databaseName   string
	collectionName string
	metrics        *[]string
	start          *time.Time
	end            *time.Time
	period         *string
}

type ListCollStatMeasurementsApiParams struct {
	GroupId        string
	ClusterName    string
	ClusterView    string
	DatabaseName   string
	CollectionName string
	Metrics        *[]string
	Start          *time.Time
	End            *time.Time
	Period         *string
}

func (a *CollectionLevelMetricsApiService) ListCollStatMeasurementsWithParams(ctx context.Context, args *ListCollStatMeasurementsApiParams) ListCollStatMeasurementsApiRequest {
	return ListCollStatMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		clusterView:    args.ClusterView,
		databaseName:   args.DatabaseName,
		collectionName: args.CollectionName,
		metrics:        args.Metrics,
		start:          args.Start,
		end:            args.End,
		period:         args.Period,
	}
}

// List that contains the metrics that you want to retrieve for the associated data series. If you don&#39;t set this parameter, this resource returns data series for all Coll Stats Latency metrics.
func (r ListCollStatMeasurementsApiRequest) Metrics(metrics []string) ListCollStatMeasurementsApiRequest {
	r.metrics = &metrics
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListCollStatMeasurementsApiRequest) Start(start time.Time) ListCollStatMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListCollStatMeasurementsApiRequest) End(end time.Time) ListCollStatMeasurementsApiRequest {
	r.end = &end
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r ListCollStatMeasurementsApiRequest) Period(period string) ListCollStatMeasurementsApiRequest {
	r.period = &period
	return r
}

func (r ListCollStatMeasurementsApiRequest) Execute() (*MeasurementsCollStatsLatencyCluster, *http.Response, error) {
	return r.ApiService.ListCollStatMeasurementsExecute(r)
}

/*
ListCollStatMeasurements Return Cluster-Level Query Latency

Get a list of the Coll Stats Latency cluster-level measurements for the given namespace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to retrieve metrics for.
	@param clusterView Human-readable label that identifies the cluster topology to retrieve metrics for.
	@param databaseName Human-readable label that identifies the database.
	@param collectionName Human-readable label that identifies the collection.
	@return ListCollStatMeasurementsApiRequest
*/
func (a *CollectionLevelMetricsApiService) ListCollStatMeasurements(ctx context.Context, groupId string, clusterName string, clusterView string, databaseName string, collectionName string) ListCollStatMeasurementsApiRequest {
	return ListCollStatMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		clusterView:    clusterView,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}

// ListCollStatMeasurementsExecute executes the request
//
//	@return MeasurementsCollStatsLatencyCluster
func (a *CollectionLevelMetricsApiService) ListCollStatMeasurementsExecute(r ListCollStatMeasurementsApiRequest) (*MeasurementsCollStatsLatencyCluster, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementsCollStatsLatencyCluster
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.ListCollStatMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/{clusterView}/{databaseName}/{collectionName}/collStats/measurements"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.clusterView == "" {
		return localVarReturnValue, nil, reportError("clusterView is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterView"+"}", url.PathEscape(r.clusterView), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.metrics != nil {
		t := *r.metrics
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "metrics", t, "multi")

	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type ListCollStatMetricsApiRequest struct {
	ctx        context.Context
	ApiService CollectionLevelMetricsApi
	groupId    string
}

type ListCollStatMetricsApiParams struct {
	GroupId string
}

func (a *CollectionLevelMetricsApiService) ListCollStatMetricsWithParams(ctx context.Context, args *ListCollStatMetricsApiParams) ListCollStatMetricsApiRequest {
	return ListCollStatMetricsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ListCollStatMetricsApiRequest) Execute() (*CollStatsLatencyNamespaceMetrics, *http.Response, error) {
	return r.ApiService.ListCollStatMetricsExecute(r)
}

/*
ListCollStatMetrics Return All Metric Names

Returns all available Coll Stats Latency metric names and their respective units for the specified project at the time of request.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListCollStatMetricsApiRequest
*/
func (a *CollectionLevelMetricsApiService) ListCollStatMetrics(ctx context.Context, groupId string) ListCollStatMetricsApiRequest {
	return ListCollStatMetricsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListCollStatMetricsExecute executes the request
//
//	@return CollStatsLatencyNamespaceMetrics
func (a *CollectionLevelMetricsApiService) ListCollStatMetricsExecute(r ListCollStatMetricsApiRequest) (*CollStatsLatencyNamespaceMetrics, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CollStatsLatencyNamespaceMetrics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.ListCollStatMetrics")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/collStats/metrics"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type ListPinnedNamespacesApiRequest struct {
	ctx         context.Context
	ApiService  CollectionLevelMetricsApi
	groupId     string
	clusterName string
}

type ListPinnedNamespacesApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *CollectionLevelMetricsApiService) ListPinnedNamespacesWithParams(ctx context.Context, args *ListPinnedNamespacesApiParams) ListPinnedNamespacesApiRequest {
	return ListPinnedNamespacesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r ListPinnedNamespacesApiRequest) Execute() (*PinnedNamespaces, *http.Response, error) {
	return r.ApiService.ListPinnedNamespacesExecute(r)
}

/*
ListPinnedNamespaces Return Pinned Namespaces

Returns a list of given cluster's pinned namespaces, a set of namespaces manually selected by users to collect query latency metrics on.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to retrieve pinned namespaces for.
	@return ListPinnedNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) ListPinnedNamespaces(ctx context.Context, groupId string, clusterName string) ListPinnedNamespacesApiRequest {
	return ListPinnedNamespacesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListPinnedNamespacesExecute executes the request
//
//	@return PinnedNamespaces
func (a *CollectionLevelMetricsApiService) ListPinnedNamespacesExecute(r ListPinnedNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PinnedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.ListPinnedNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/collStats/pinned"
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

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type ListProcessMeasurementsApiRequest struct {
	ctx            context.Context
	ApiService     CollectionLevelMetricsApi
	groupId        string
	processId      string
	databaseName   string
	collectionName string
	metrics        *[]string
	start          *time.Time
	end            *time.Time
	period         *string
}

type ListProcessMeasurementsApiParams struct {
	GroupId        string
	ProcessId      string
	DatabaseName   string
	CollectionName string
	Metrics        *[]string
	Start          *time.Time
	End            *time.Time
	Period         *string
}

func (a *CollectionLevelMetricsApiService) ListProcessMeasurementsWithParams(ctx context.Context, args *ListProcessMeasurementsApiParams) ListProcessMeasurementsApiRequest {
	return ListProcessMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		processId:      args.ProcessId,
		databaseName:   args.DatabaseName,
		collectionName: args.CollectionName,
		metrics:        args.Metrics,
		start:          args.Start,
		end:            args.End,
		period:         args.Period,
	}
}

// List that contains the metrics that you want to retrieve for the associated data series. If you don&#39;t set this parameter, this resource returns data series for all Coll Stats Latency metrics.
func (r ListProcessMeasurementsApiRequest) Metrics(metrics []string) ListProcessMeasurementsApiRequest {
	r.metrics = &metrics
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListProcessMeasurementsApiRequest) Start(start time.Time) ListProcessMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListProcessMeasurementsApiRequest) End(end time.Time) ListProcessMeasurementsApiRequest {
	r.end = &end
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r ListProcessMeasurementsApiRequest) Period(period string) ListProcessMeasurementsApiRequest {
	r.period = &period
	return r
}

func (r ListProcessMeasurementsApiRequest) Execute() (*MeasurementsCollStatsLatencyHost, *http.Response, error) {
	return r.ApiService.ListProcessMeasurementsExecute(r)
}

/*
ListProcessMeasurements Return Host-Level Query Latency

Get a list of the Coll Stats Latency process-level measurements for the given namespace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@param databaseName Human-readable label that identifies the database.
	@param collectionName Human-readable label that identifies the collection.
	@return ListProcessMeasurementsApiRequest
*/
func (a *CollectionLevelMetricsApiService) ListProcessMeasurements(ctx context.Context, groupId string, processId string, databaseName string, collectionName string) ListProcessMeasurementsApiRequest {
	return ListProcessMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		processId:      processId,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}

// ListProcessMeasurementsExecute executes the request
//
//	@return MeasurementsCollStatsLatencyHost
func (a *CollectionLevelMetricsApiService) ListProcessMeasurementsExecute(r ListProcessMeasurementsApiRequest) (*MeasurementsCollStatsLatencyHost, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementsCollStatsLatencyHost
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.ListProcessMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/{databaseName}/{collectionName}/collStats/measurements"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.metrics != nil {
		t := *r.metrics
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "metrics", t, "multi")

	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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

type PinNamespacesApiRequest struct {
	ctx               context.Context
	ApiService        CollectionLevelMetricsApi
	groupId           string
	clusterName       string
	namespacesRequest *NamespacesRequest
}

type PinNamespacesApiParams struct {
	GroupId           string
	ClusterName       string
	NamespacesRequest *NamespacesRequest
}

func (a *CollectionLevelMetricsApiService) PinNamespacesWithParams(ctx context.Context, args *PinNamespacesApiParams) PinNamespacesApiRequest {
	return PinNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		clusterName:       args.ClusterName,
		namespacesRequest: args.NamespacesRequest,
	}
}

func (r PinNamespacesApiRequest) Execute() (*PinnedNamespaces, *http.Response, error) {
	return r.ApiService.PinNamespacesExecute(r)
}

/*
PinNamespaces Pin Namespaces

Pin provided list of namespaces for collection-level latency metrics collection for the given Group and Cluster. This initializes a pinned namespaces list or replaces any existing pinned namespaces list for the Group and Cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
	@return PinNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) PinNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) PinNamespacesApiRequest {
	return PinNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		clusterName:       clusterName,
		namespacesRequest: namespacesRequest,
	}
}

// PinNamespacesExecute executes the request
//
//	@return PinnedNamespaces
func (a *CollectionLevelMetricsApiService) PinNamespacesExecute(r PinNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PinnedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.PinNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/collStats/pinned"
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
	if r.namespacesRequest == nil {
		return localVarReturnValue, nil, reportError("namespacesRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.namespacesRequest
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

type UnpinNamespacesApiRequest struct {
	ctx               context.Context
	ApiService        CollectionLevelMetricsApi
	groupId           string
	clusterName       string
	namespacesRequest *NamespacesRequest
}

type UnpinNamespacesApiParams struct {
	GroupId           string
	ClusterName       string
	NamespacesRequest *NamespacesRequest
}

func (a *CollectionLevelMetricsApiService) UnpinNamespacesWithParams(ctx context.Context, args *UnpinNamespacesApiParams) UnpinNamespacesApiRequest {
	return UnpinNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		clusterName:       args.ClusterName,
		namespacesRequest: args.NamespacesRequest,
	}
}

func (r UnpinNamespacesApiRequest) Execute() (*PinnedNamespaces, *http.Response, error) {
	return r.ApiService.UnpinNamespacesExecute(r)
}

/*
UnpinNamespaces Unpin Namespaces

Unpin provided list of namespaces for collection-level latency metrics collection for the given Group and Cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to unpin namespaces from.
	@return UnpinNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) UnpinNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) UnpinNamespacesApiRequest {
	return UnpinNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		clusterName:       clusterName,
		namespacesRequest: namespacesRequest,
	}
}

// UnpinNamespacesExecute executes the request
//
//	@return PinnedNamespaces
func (a *CollectionLevelMetricsApiService) UnpinNamespacesExecute(r UnpinNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PinnedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.UnpinNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/collStats/unpin"
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
	if r.namespacesRequest == nil {
		return localVarReturnValue, nil, reportError("namespacesRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.namespacesRequest
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

type UpdatePinnedNamespacesApiRequest struct {
	ctx               context.Context
	ApiService        CollectionLevelMetricsApi
	groupId           string
	clusterName       string
	namespacesRequest *NamespacesRequest
}

type UpdatePinnedNamespacesApiParams struct {
	GroupId           string
	ClusterName       string
	NamespacesRequest *NamespacesRequest
}

func (a *CollectionLevelMetricsApiService) UpdatePinnedNamespacesWithParams(ctx context.Context, args *UpdatePinnedNamespacesApiParams) UpdatePinnedNamespacesApiRequest {
	return UpdatePinnedNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		clusterName:       args.ClusterName,
		namespacesRequest: args.NamespacesRequest,
	}
}

func (r UpdatePinnedNamespacesApiRequest) Execute() (*PinnedNamespaces, *http.Response, error) {
	return r.ApiService.UpdatePinnedNamespacesExecute(r)
}

/*
UpdatePinnedNamespaces Add Pinned Namespaces

Add provided list of namespaces to existing pinned namespaces list for collection-level latency metrics collection for the given Group and Cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to pin namespaces to.
	@return UpdatePinnedNamespacesApiRequest
*/
func (a *CollectionLevelMetricsApiService) UpdatePinnedNamespaces(ctx context.Context, groupId string, clusterName string, namespacesRequest *NamespacesRequest) UpdatePinnedNamespacesApiRequest {
	return UpdatePinnedNamespacesApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		clusterName:       clusterName,
		namespacesRequest: namespacesRequest,
	}
}

// UpdatePinnedNamespacesExecute executes the request
//
//	@return PinnedNamespaces
func (a *CollectionLevelMetricsApiService) UpdatePinnedNamespacesExecute(r UpdatePinnedNamespacesApiRequest) (*PinnedNamespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PinnedNamespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CollectionLevelMetricsApiService.UpdatePinnedNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/collStats/pinned"
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
	if r.namespacesRequest == nil {
		return localVarReturnValue, nil, reportError("namespacesRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.namespacesRequest
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
