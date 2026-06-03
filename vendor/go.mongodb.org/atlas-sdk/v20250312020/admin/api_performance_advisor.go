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

type PerformanceAdvisorApi interface {

	/*
		DisableManagedSlowMs Disable Managed Slow Operation Threshold

		Disables the slow operation threshold that MongoDB Cloud calculated for the specified project. The threshold determines which operations the Performance Advisor and Query Profiler considers slow. When enabled, MongoDB Cloud uses the average execution time for operations on your cluster to determine slow-running queries. As a result, the threshold is more pertinent to your cluster workload. The slow operation threshold is enabled by default for dedicated clusters (M10+). When disabled, MongoDB Cloud considers any operation that takes longer than 100 milliseconds to be slow.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DisableManagedSlowMsApiRequest
	*/
	DisableManagedSlowMs(ctx context.Context, groupId string) DisableManagedSlowMsApiRequest
	/*
		DisableManagedSlowMs Disable Managed Slow Operation Threshold


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DisableManagedSlowMsApiParams - Parameters for the request
		@return DisableManagedSlowMsApiRequest
	*/
	DisableManagedSlowMsWithParams(ctx context.Context, args *DisableManagedSlowMsApiParams) DisableManagedSlowMsApiRequest

	// Method available only for mocking purposes
	DisableManagedSlowMsExecute(r DisableManagedSlowMsApiRequest) (*http.Response, error)

	/*
		EnableManagedSlowMs Enable Managed Slow Operation Threshold

		Enables MongoDB Cloud to use its slow operation threshold for the specified project. The threshold determines which operations the Performance Advisor and Query Profiler considers slow. When enabled, MongoDB Cloud uses the average execution time for operations on your cluster to determine slow-running queries. As a result, the threshold is more pertinent to your cluster workload. The slow operation threshold is enabled by default for dedicated clusters (M10+). When disabled, MongoDB Cloud considers any operation that takes longer than 100 milliseconds to be slow.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return EnableManagedSlowMsApiRequest
	*/
	EnableManagedSlowMs(ctx context.Context, groupId string) EnableManagedSlowMsApiRequest
	/*
		EnableManagedSlowMs Enable Managed Slow Operation Threshold


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param EnableManagedSlowMsApiParams - Parameters for the request
		@return EnableManagedSlowMsApiRequest
	*/
	EnableManagedSlowMsWithParams(ctx context.Context, args *EnableManagedSlowMsApiParams) EnableManagedSlowMsApiRequest

	// Method available only for mocking purposes
	EnableManagedSlowMsExecute(r EnableManagedSlowMsApiRequest) (*http.Response, error)

	/*
		GetManagedSlowMs Return Managed Slow Operation Threshold Status

		Get whether the Managed Slow MS feature is enabled.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetManagedSlowMsApiRequest
	*/
	GetManagedSlowMs(ctx context.Context, groupId string) GetManagedSlowMsApiRequest
	/*
		GetManagedSlowMs Return Managed Slow Operation Threshold Status


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetManagedSlowMsApiParams - Parameters for the request
		@return GetManagedSlowMsApiRequest
	*/
	GetManagedSlowMsWithParams(ctx context.Context, args *GetManagedSlowMsApiParams) GetManagedSlowMsApiRequest

	// Method available only for mocking purposes
	GetManagedSlowMsExecute(r GetManagedSlowMsApiRequest) (bool, *http.Response, error)

	/*
		GetServerlessAutoIndexing Return Serverless Auto-Indexing Status

		Get whether the Serverless Auto Indexing feature is enabled. This endpoint returns a value for Flex clusters that were created with the Create Serverless Instance endpoint or Flex clusters that were migrated from Serverless instances. However, the value returned is not indicative of the Auto Indexing state as Auto Indexing is unavailable for Flex clusters. This endpoint will be sunset on January 22, 2026.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return GetServerlessAutoIndexingApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PerformanceAdvisorApi
	*/
	GetServerlessAutoIndexing(ctx context.Context, groupId string, clusterName string) GetServerlessAutoIndexingApiRequest
	/*
		GetServerlessAutoIndexing Return Serverless Auto-Indexing Status


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServerlessAutoIndexingApiParams - Parameters for the request
		@return GetServerlessAutoIndexingApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PerformanceAdvisorApi
	*/
	GetServerlessAutoIndexingWithParams(ctx context.Context, args *GetServerlessAutoIndexingApiParams) GetServerlessAutoIndexingApiRequest

	// Method available only for mocking purposes
	GetServerlessAutoIndexingExecute(r GetServerlessAutoIndexingApiRequest) (bool, *http.Response, error)

	/*
		ListClusterSuggestedIndexes Return All Suggested Indexes

		Returns the indexes that the Performance Advisor suggests. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListClusterSuggestedIndexesApiRequest
	*/
	ListClusterSuggestedIndexes(ctx context.Context, groupId string, clusterName string) ListClusterSuggestedIndexesApiRequest
	/*
		ListClusterSuggestedIndexes Return All Suggested Indexes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterSuggestedIndexesApiParams - Parameters for the request
		@return ListClusterSuggestedIndexesApiRequest
	*/
	ListClusterSuggestedIndexesWithParams(ctx context.Context, args *ListClusterSuggestedIndexesApiParams) ListClusterSuggestedIndexesApiRequest

	// Method available only for mocking purposes
	ListClusterSuggestedIndexesExecute(r ListClusterSuggestedIndexesApiRequest) (*PerformanceAdvisorResponse, *http.Response, error)

	/*
		ListDropIndexSuggestions Return All Suggested Indexes to Drop

		Returns the indexes that the Performance Advisor suggests to drop. The Performance Advisor suggests dropping unused, redundant, and hidden indexes to improve write performance and increase storage space.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListDropIndexSuggestionsApiRequest
	*/
	ListDropIndexSuggestions(ctx context.Context, groupId string, clusterName string) ListDropIndexSuggestionsApiRequest
	/*
		ListDropIndexSuggestions Return All Suggested Indexes to Drop


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDropIndexSuggestionsApiParams - Parameters for the request
		@return ListDropIndexSuggestionsApiRequest
	*/
	ListDropIndexSuggestionsWithParams(ctx context.Context, args *ListDropIndexSuggestionsApiParams) ListDropIndexSuggestionsApiRequest

	// Method available only for mocking purposes
	ListDropIndexSuggestionsExecute(r ListDropIndexSuggestionsApiRequest) (*DropIndexSuggestionsResponse, *http.Response, error)

	/*
		ListPerformanceAdvisorNamespaces Return All Namespaces for One Host

		Returns up to 20 namespaces for collections experiencing slow queries on the specified host. If you specify a secondary member of a replica set that hasn't received any database read operations, the endpoint doesn't return any namespaces.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return ListPerformanceAdvisorNamespacesApiRequest
	*/
	ListPerformanceAdvisorNamespaces(ctx context.Context, groupId string, processId string) ListPerformanceAdvisorNamespacesApiRequest
	/*
		ListPerformanceAdvisorNamespaces Return All Namespaces for One Host


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPerformanceAdvisorNamespacesApiParams - Parameters for the request
		@return ListPerformanceAdvisorNamespacesApiRequest
	*/
	ListPerformanceAdvisorNamespacesWithParams(ctx context.Context, args *ListPerformanceAdvisorNamespacesApiParams) ListPerformanceAdvisorNamespacesApiRequest

	// Method available only for mocking purposes
	ListPerformanceAdvisorNamespacesExecute(r ListPerformanceAdvisorNamespacesApiRequest) (*Namespaces, *http.Response, error)

	/*
		ListSchemaAdvice Return Schema Advice

		Returns the schema suggestions that the Performance Advisor detects. The Performance Advisor provides holistic schema recommendations for your cluster by sampling documents in your most active collections and collections with slow-running queries.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListSchemaAdviceApiRequest
	*/
	ListSchemaAdvice(ctx context.Context, groupId string, clusterName string) ListSchemaAdviceApiRequest
	/*
		ListSchemaAdvice Return Schema Advice


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListSchemaAdviceApiParams - Parameters for the request
		@return ListSchemaAdviceApiRequest
	*/
	ListSchemaAdviceWithParams(ctx context.Context, args *ListSchemaAdviceApiParams) ListSchemaAdviceApiRequest

	// Method available only for mocking purposes
	ListSchemaAdviceExecute(r ListSchemaAdviceApiRequest) (*SchemaAdvisorResponse, *http.Response, error)

	/*
		ListSlowQueryLogs Return Slow Queries

		Returns log lines for slow queries that the Performance Advisor and Query Profiler identified. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance. MongoDB Cloud bases the threshold for slow queries on the average time of operations on your cluster. This enables workload-relevant recommendations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return ListSlowQueryLogsApiRequest
	*/
	ListSlowQueryLogs(ctx context.Context, groupId string, processId string) ListSlowQueryLogsApiRequest
	/*
		ListSlowQueryLogs Return Slow Queries


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListSlowQueryLogsApiParams - Parameters for the request
		@return ListSlowQueryLogsApiRequest
	*/
	ListSlowQueryLogsWithParams(ctx context.Context, args *ListSlowQueryLogsApiParams) ListSlowQueryLogsApiRequest

	// Method available only for mocking purposes
	ListSlowQueryLogsExecute(r ListSlowQueryLogsApiRequest) (*PerformanceAdvisorSlowQueryList, *http.Response, error)

	/*
		ListSuggestedIndexes Return All Suggested Indexes

		Returns the indexes that the Performance Advisor suggests. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return ListSuggestedIndexesApiRequest
	*/
	ListSuggestedIndexes(ctx context.Context, groupId string, processId string) ListSuggestedIndexesApiRequest
	/*
		ListSuggestedIndexes Return All Suggested Indexes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListSuggestedIndexesApiParams - Parameters for the request
		@return ListSuggestedIndexesApiRequest
	*/
	ListSuggestedIndexesWithParams(ctx context.Context, args *ListSuggestedIndexesApiParams) ListSuggestedIndexesApiRequest

	// Method available only for mocking purposes
	ListSuggestedIndexesExecute(r ListSuggestedIndexesApiRequest) (*PerformanceAdvisorResponse, *http.Response, error)

	/*
		SetServerlessAutoIndexing Set Serverless Auto-Indexing Status

		Set whether the Serverless Auto Indexing feature is enabled. This endpoint sets a value for Flex clusters that were created with the Create Serverless Instance endpoint or Flex clusters that were migrated from Serverless instances. However, the value returned is not indicative of the Auto Indexing state as Auto Indexing is unavailable for Flex clusters. This endpoint will be sunset on January 22, 2026.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return SetServerlessAutoIndexingApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PerformanceAdvisorApi
	*/
	SetServerlessAutoIndexing(ctx context.Context, groupId string, clusterName string) SetServerlessAutoIndexingApiRequest
	/*
		SetServerlessAutoIndexing Set Serverless Auto-Indexing Status


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param SetServerlessAutoIndexingApiParams - Parameters for the request
		@return SetServerlessAutoIndexingApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PerformanceAdvisorApi
	*/
	SetServerlessAutoIndexingWithParams(ctx context.Context, args *SetServerlessAutoIndexingApiParams) SetServerlessAutoIndexingApiRequest

	// Method available only for mocking purposes
	SetServerlessAutoIndexingExecute(r SetServerlessAutoIndexingApiRequest) (*http.Response, error)
}

// PerformanceAdvisorApiService PerformanceAdvisorApi service
type PerformanceAdvisorApiService service

type DisableManagedSlowMsApiRequest struct {
	ctx        context.Context
	ApiService PerformanceAdvisorApi
	groupId    string
}

type DisableManagedSlowMsApiParams struct {
	GroupId string
}

func (a *PerformanceAdvisorApiService) DisableManagedSlowMsWithParams(ctx context.Context, args *DisableManagedSlowMsApiParams) DisableManagedSlowMsApiRequest {
	return DisableManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DisableManagedSlowMsApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DisableManagedSlowMsExecute(r)
}

/*
DisableManagedSlowMs Disable Managed Slow Operation Threshold

Disables the slow operation threshold that MongoDB Cloud calculated for the specified project. The threshold determines which operations the Performance Advisor and Query Profiler considers slow. When enabled, MongoDB Cloud uses the average execution time for operations on your cluster to determine slow-running queries. As a result, the threshold is more pertinent to your cluster workload. The slow operation threshold is enabled by default for dedicated clusters (M10+). When disabled, MongoDB Cloud considers any operation that takes longer than 100 milliseconds to be slow.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DisableManagedSlowMsApiRequest
*/
func (a *PerformanceAdvisorApiService) DisableManagedSlowMs(ctx context.Context, groupId string) DisableManagedSlowMsApiRequest {
	return DisableManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DisableManagedSlowMsExecute executes the request
func (a *PerformanceAdvisorApiService) DisableManagedSlowMsExecute(r DisableManagedSlowMsApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.DisableManagedSlowMs")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/managedSlowMs/disable"
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

type EnableManagedSlowMsApiRequest struct {
	ctx        context.Context
	ApiService PerformanceAdvisorApi
	groupId    string
}

type EnableManagedSlowMsApiParams struct {
	GroupId string
}

func (a *PerformanceAdvisorApiService) EnableManagedSlowMsWithParams(ctx context.Context, args *EnableManagedSlowMsApiParams) EnableManagedSlowMsApiRequest {
	return EnableManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r EnableManagedSlowMsApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.EnableManagedSlowMsExecute(r)
}

/*
EnableManagedSlowMs Enable Managed Slow Operation Threshold

Enables MongoDB Cloud to use its slow operation threshold for the specified project. The threshold determines which operations the Performance Advisor and Query Profiler considers slow. When enabled, MongoDB Cloud uses the average execution time for operations on your cluster to determine slow-running queries. As a result, the threshold is more pertinent to your cluster workload. The slow operation threshold is enabled by default for dedicated clusters (M10+). When disabled, MongoDB Cloud considers any operation that takes longer than 100 milliseconds to be slow.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return EnableManagedSlowMsApiRequest
*/
func (a *PerformanceAdvisorApiService) EnableManagedSlowMs(ctx context.Context, groupId string) EnableManagedSlowMsApiRequest {
	return EnableManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// EnableManagedSlowMsExecute executes the request
func (a *PerformanceAdvisorApiService) EnableManagedSlowMsExecute(r EnableManagedSlowMsApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.EnableManagedSlowMs")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/managedSlowMs/enable"
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

type GetManagedSlowMsApiRequest struct {
	ctx        context.Context
	ApiService PerformanceAdvisorApi
	groupId    string
}

type GetManagedSlowMsApiParams struct {
	GroupId string
}

func (a *PerformanceAdvisorApiService) GetManagedSlowMsWithParams(ctx context.Context, args *GetManagedSlowMsApiParams) GetManagedSlowMsApiRequest {
	return GetManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetManagedSlowMsApiRequest) Execute() (bool, *http.Response, error) {
	return r.ApiService.GetManagedSlowMsExecute(r)
}

/*
GetManagedSlowMs Return Managed Slow Operation Threshold Status

Get whether the Managed Slow MS feature is enabled.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetManagedSlowMsApiRequest
*/
func (a *PerformanceAdvisorApiService) GetManagedSlowMs(ctx context.Context, groupId string) GetManagedSlowMsApiRequest {
	return GetManagedSlowMsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetManagedSlowMsExecute executes the request
//
//	@return bool
func (a *PerformanceAdvisorApiService) GetManagedSlowMsExecute(r GetManagedSlowMsApiRequest) (bool, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue bool
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.GetManagedSlowMs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/managedSlowMs"
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

type GetServerlessAutoIndexingApiRequest struct {
	ctx         context.Context
	ApiService  PerformanceAdvisorApi
	groupId     string
	clusterName string
}

type GetServerlessAutoIndexingApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *PerformanceAdvisorApiService) GetServerlessAutoIndexingWithParams(ctx context.Context, args *GetServerlessAutoIndexingApiParams) GetServerlessAutoIndexingApiRequest {
	return GetServerlessAutoIndexingApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetServerlessAutoIndexingApiRequest) Execute() (bool, *http.Response, error) {
	return r.ApiService.GetServerlessAutoIndexingExecute(r)
}

/*
GetServerlessAutoIndexing Return Serverless Auto-Indexing Status

Get whether the Serverless Auto Indexing feature is enabled. This endpoint returns a value for Flex clusters that were created with the Create Serverless Instance endpoint or Flex clusters that were migrated from Serverless instances. However, the value returned is not indicative of the Auto Indexing state as Auto Indexing is unavailable for Flex clusters. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return GetServerlessAutoIndexingApiRequest

Deprecated
*/
func (a *PerformanceAdvisorApiService) GetServerlessAutoIndexing(ctx context.Context, groupId string, clusterName string) GetServerlessAutoIndexingApiRequest {
	return GetServerlessAutoIndexingApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetServerlessAutoIndexingExecute executes the request
//
//	@return bool
//
// Deprecated
func (a *PerformanceAdvisorApiService) GetServerlessAutoIndexingExecute(r GetServerlessAutoIndexingApiRequest) (bool, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue bool
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.GetServerlessAutoIndexing")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/performanceAdvisor/autoIndexing"
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

type ListClusterSuggestedIndexesApiRequest struct {
	ctx         context.Context
	ApiService  PerformanceAdvisorApi
	groupId     string
	clusterName string
	processIds  *[]string
	namespaces  *[]string
	since       *int64
	until       *int64
}

type ListClusterSuggestedIndexesApiParams struct {
	GroupId     string
	ClusterName string
	ProcessIds  *[]string
	Namespaces  *[]string
	Since       *int64
	Until       *int64
}

func (a *PerformanceAdvisorApiService) ListClusterSuggestedIndexesWithParams(ctx context.Context, args *ListClusterSuggestedIndexesApiParams) ListClusterSuggestedIndexesApiRequest {
	return ListClusterSuggestedIndexesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		processIds:  args.ProcessIds,
		namespaces:  args.Namespaces,
		since:       args.Since,
		until:       args.Until,
	}
}

// Process IDs from which to retrieve suggested indexes. A &#x60;processId&#x60; is a combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (&#x60;mongod&#x60; or &#x60;mongos&#x60;). The port must be the IANA port on which the MongoDB process listens for requests. To include multiple &#x60;processIds&#x60;, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each &#x60;processId&#x60;.
func (r ListClusterSuggestedIndexesApiRequest) ProcessIds(processIds []string) ListClusterSuggestedIndexesApiRequest {
	r.processIds = &processIds
	return r
}

// Namespaces from which to retrieve suggested indexes. A namespace consists of one database and one collection resource written as &#x60;.&#x60;: &#x60;&lt;database&gt;.&lt;collection&gt;&#x60;. To include multiple namespaces, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each namespace. Omit this parameter to return results for all namespaces.
func (r ListClusterSuggestedIndexesApiRequest) Namespaces(namespaces []string) ListClusterSuggestedIndexesApiRequest {
	r.namespaces = &namespaces
	return r
}

// Date and time from which the query retrieves the suggested indexes. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **until** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r ListClusterSuggestedIndexesApiRequest) Since(since int64) ListClusterSuggestedIndexesApiRequest {
	r.since = &since
	return r
}

// Date and time up until which the query retrieves the suggested indexes. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you specify the **until** parameter, you must specify the **since** parameter. - If you specify neither the **since** nor the **until** parameters, the endpoint returns data from the previous 24 hours.
func (r ListClusterSuggestedIndexesApiRequest) Until(until int64) ListClusterSuggestedIndexesApiRequest {
	r.until = &until
	return r
}

func (r ListClusterSuggestedIndexesApiRequest) Execute() (*PerformanceAdvisorResponse, *http.Response, error) {
	return r.ApiService.ListClusterSuggestedIndexesExecute(r)
}

/*
ListClusterSuggestedIndexes Return All Suggested Indexes

Returns the indexes that the Performance Advisor suggests. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListClusterSuggestedIndexesApiRequest
*/
func (a *PerformanceAdvisorApiService) ListClusterSuggestedIndexes(ctx context.Context, groupId string, clusterName string) ListClusterSuggestedIndexesApiRequest {
	return ListClusterSuggestedIndexesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterSuggestedIndexesExecute executes the request
//
//	@return PerformanceAdvisorResponse
func (a *PerformanceAdvisorApiService) ListClusterSuggestedIndexesExecute(r ListClusterSuggestedIndexesApiRequest) (*PerformanceAdvisorResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PerformanceAdvisorResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListClusterSuggestedIndexes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/performanceAdvisor/suggestedIndexes"
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
	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
	}
	if r.until != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "until", r.until, "")
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

type ListDropIndexSuggestionsApiRequest struct {
	ctx         context.Context
	ApiService  PerformanceAdvisorApi
	groupId     string
	clusterName string
}

type ListDropIndexSuggestionsApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *PerformanceAdvisorApiService) ListDropIndexSuggestionsWithParams(ctx context.Context, args *ListDropIndexSuggestionsApiParams) ListDropIndexSuggestionsApiRequest {
	return ListDropIndexSuggestionsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r ListDropIndexSuggestionsApiRequest) Execute() (*DropIndexSuggestionsResponse, *http.Response, error) {
	return r.ApiService.ListDropIndexSuggestionsExecute(r)
}

/*
ListDropIndexSuggestions Return All Suggested Indexes to Drop

Returns the indexes that the Performance Advisor suggests to drop. The Performance Advisor suggests dropping unused, redundant, and hidden indexes to improve write performance and increase storage space.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListDropIndexSuggestionsApiRequest
*/
func (a *PerformanceAdvisorApiService) ListDropIndexSuggestions(ctx context.Context, groupId string, clusterName string) ListDropIndexSuggestionsApiRequest {
	return ListDropIndexSuggestionsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListDropIndexSuggestionsExecute executes the request
//
//	@return DropIndexSuggestionsResponse
func (a *PerformanceAdvisorApiService) ListDropIndexSuggestionsExecute(r ListDropIndexSuggestionsApiRequest) (*DropIndexSuggestionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DropIndexSuggestionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListDropIndexSuggestions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/performanceAdvisor/dropIndexSuggestions"
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

type ListPerformanceAdvisorNamespacesApiRequest struct {
	ctx        context.Context
	ApiService PerformanceAdvisorApi
	groupId    string
	processId  string
	duration   *int64
	since      *int64
}

type ListPerformanceAdvisorNamespacesApiParams struct {
	GroupId   string
	ProcessId string
	Duration  *int64
	Since     *int64
}

func (a *PerformanceAdvisorApiService) ListPerformanceAdvisorNamespacesWithParams(ctx context.Context, args *ListPerformanceAdvisorNamespacesApiParams) ListPerformanceAdvisorNamespacesApiRequest {
	return ListPerformanceAdvisorNamespacesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		processId:  args.ProcessId,
		duration:   args.Duration,
		since:      args.Since,
	}
}

// Length of time expressed during which the query finds suggested indexes among the managed namespaces in the cluster. This parameter expresses its value in milliseconds.  - If you don&#39;t specify the **since** parameter, the endpoint returns data covering the duration before the current time. - If you specify neither the **duration** nor **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListPerformanceAdvisorNamespacesApiRequest) Duration(duration int64) ListPerformanceAdvisorNamespacesApiRequest {
	r.duration = &duration
	return r
}

// Date and time from which the query retrieves the suggested indexes. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **duration** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **duration** nor the **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListPerformanceAdvisorNamespacesApiRequest) Since(since int64) ListPerformanceAdvisorNamespacesApiRequest {
	r.since = &since
	return r
}

func (r ListPerformanceAdvisorNamespacesApiRequest) Execute() (*Namespaces, *http.Response, error) {
	return r.ApiService.ListPerformanceAdvisorNamespacesExecute(r)
}

/*
ListPerformanceAdvisorNamespaces Return All Namespaces for One Host

Returns up to 20 namespaces for collections experiencing slow queries on the specified host. If you specify a secondary member of a replica set that hasn't received any database read operations, the endpoint doesn't return any namespaces.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return ListPerformanceAdvisorNamespacesApiRequest
*/
func (a *PerformanceAdvisorApiService) ListPerformanceAdvisorNamespaces(ctx context.Context, groupId string, processId string) ListPerformanceAdvisorNamespacesApiRequest {
	return ListPerformanceAdvisorNamespacesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// ListPerformanceAdvisorNamespacesExecute executes the request
//
//	@return Namespaces
func (a *PerformanceAdvisorApiService) ListPerformanceAdvisorNamespacesExecute(r ListPerformanceAdvisorNamespacesApiRequest) (*Namespaces, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *Namespaces
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListPerformanceAdvisorNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/performanceAdvisor/namespaces"
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

	if r.duration != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "duration", r.duration, "")
	}
	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
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

type ListSchemaAdviceApiRequest struct {
	ctx         context.Context
	ApiService  PerformanceAdvisorApi
	groupId     string
	clusterName string
}

type ListSchemaAdviceApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *PerformanceAdvisorApiService) ListSchemaAdviceWithParams(ctx context.Context, args *ListSchemaAdviceApiParams) ListSchemaAdviceApiRequest {
	return ListSchemaAdviceApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r ListSchemaAdviceApiRequest) Execute() (*SchemaAdvisorResponse, *http.Response, error) {
	return r.ApiService.ListSchemaAdviceExecute(r)
}

/*
ListSchemaAdvice Return Schema Advice

Returns the schema suggestions that the Performance Advisor detects. The Performance Advisor provides holistic schema recommendations for your cluster by sampling documents in your most active collections and collections with slow-running queries.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListSchemaAdviceApiRequest
*/
func (a *PerformanceAdvisorApiService) ListSchemaAdvice(ctx context.Context, groupId string, clusterName string) ListSchemaAdviceApiRequest {
	return ListSchemaAdviceApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListSchemaAdviceExecute executes the request
//
//	@return SchemaAdvisorResponse
func (a *PerformanceAdvisorApiService) ListSchemaAdviceExecute(r ListSchemaAdviceApiRequest) (*SchemaAdvisorResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SchemaAdvisorResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListSchemaAdvice")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/performanceAdvisor/schemaAdvice"
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

type ListSlowQueryLogsApiRequest struct {
	ctx                 context.Context
	ApiService          PerformanceAdvisorApi
	groupId             string
	processId           string
	duration            *int64
	namespaces          *[]string
	nLogs               *int64
	since               *int64
	includeMetrics      *bool
	includeReplicaState *bool
	includeOpType       *bool
}

type ListSlowQueryLogsApiParams struct {
	GroupId             string
	ProcessId           string
	Duration            *int64
	Namespaces          *[]string
	NLogs               *int64
	Since               *int64
	IncludeMetrics      *bool
	IncludeReplicaState *bool
	IncludeOpType       *bool
}

func (a *PerformanceAdvisorApiService) ListSlowQueryLogsWithParams(ctx context.Context, args *ListSlowQueryLogsApiParams) ListSlowQueryLogsApiRequest {
	return ListSlowQueryLogsApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             args.GroupId,
		processId:           args.ProcessId,
		duration:            args.Duration,
		namespaces:          args.Namespaces,
		nLogs:               args.NLogs,
		since:               args.Since,
		includeMetrics:      args.IncludeMetrics,
		includeReplicaState: args.IncludeReplicaState,
		includeOpType:       args.IncludeOpType,
	}
}

// Length of time expressed during which the query finds slow queries among the managed namespaces in the cluster. This parameter expresses its value in milliseconds.  - If you don&#39;t specify the **since** parameter, the endpoint returns data covering the duration before the current time. - If you specify neither the **duration** nor **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListSlowQueryLogsApiRequest) Duration(duration int64) ListSlowQueryLogsApiRequest {
	r.duration = &duration
	return r
}

// Namespaces from which to retrieve slow queries. A namespace consists of one database and one collection resource written as &#x60;.&#x60;: &#x60;&lt;database&gt;.&lt;collection&gt;&#x60;. To include multiple namespaces, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each namespace. Omit this parameter to return results for all namespaces.
func (r ListSlowQueryLogsApiRequest) Namespaces(namespaces []string) ListSlowQueryLogsApiRequest {
	r.namespaces = &namespaces
	return r
}

// Maximum number of lines from the log to return.
func (r ListSlowQueryLogsApiRequest) NLogs(nLogs int64) ListSlowQueryLogsApiRequest {
	r.nLogs = &nLogs
	return r
}

// Date and time from which the query retrieves the slow queries. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **duration** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **duration** nor the **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListSlowQueryLogsApiRequest) Since(since int64) ListSlowQueryLogsApiRequest {
	r.since = &since
	return r
}

// Whether or not to include metrics extracted from the slow query log as separate fields.
func (r ListSlowQueryLogsApiRequest) IncludeMetrics(includeMetrics bool) ListSlowQueryLogsApiRequest {
	r.includeMetrics = &includeMetrics
	return r
}

// Whether or not to include the replica state of the host when the slow query log was generated as a separate field.
func (r ListSlowQueryLogsApiRequest) IncludeReplicaState(includeReplicaState bool) ListSlowQueryLogsApiRequest {
	r.includeReplicaState = &includeReplicaState
	return r
}

// Whether or not to include the operation type (read/write/command) extracted from the slow query log as a separate field.
func (r ListSlowQueryLogsApiRequest) IncludeOpType(includeOpType bool) ListSlowQueryLogsApiRequest {
	r.includeOpType = &includeOpType
	return r
}

func (r ListSlowQueryLogsApiRequest) Execute() (*PerformanceAdvisorSlowQueryList, *http.Response, error) {
	return r.ApiService.ListSlowQueryLogsExecute(r)
}

/*
ListSlowQueryLogs Return Slow Queries

Returns log lines for slow queries that the Performance Advisor and Query Profiler identified. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance. MongoDB Cloud bases the threshold for slow queries on the average time of operations on your cluster. This enables workload-relevant recommendations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return ListSlowQueryLogsApiRequest
*/
func (a *PerformanceAdvisorApiService) ListSlowQueryLogs(ctx context.Context, groupId string, processId string) ListSlowQueryLogsApiRequest {
	return ListSlowQueryLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// ListSlowQueryLogsExecute executes the request
//
//	@return PerformanceAdvisorSlowQueryList
func (a *PerformanceAdvisorApiService) ListSlowQueryLogsExecute(r ListSlowQueryLogsApiRequest) (*PerformanceAdvisorSlowQueryList, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PerformanceAdvisorSlowQueryList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListSlowQueryLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/performanceAdvisor/slowQueryLogs"
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

	if r.duration != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "duration", r.duration, "")
	}
	if r.namespaces != nil {
		t := *r.namespaces
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "namespaces", t, "multi")

	}
	if r.nLogs != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	} else {
		var defaultValue int64 = 20000
		r.nLogs = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	}
	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
	}
	if r.includeMetrics != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeMetrics", r.includeMetrics, "")
	} else {
		var defaultValue bool = false
		r.includeMetrics = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeMetrics", r.includeMetrics, "")
	}
	if r.includeReplicaState != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeReplicaState", r.includeReplicaState, "")
	} else {
		var defaultValue bool = false
		r.includeReplicaState = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeReplicaState", r.includeReplicaState, "")
	}
	if r.includeOpType != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeOpType", r.includeOpType, "")
	} else {
		var defaultValue bool = false
		r.includeOpType = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeOpType", r.includeOpType, "")
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

type ListSuggestedIndexesApiRequest struct {
	ctx          context.Context
	ApiService   PerformanceAdvisorApi
	groupId      string
	processId    string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	duration     *int64
	namespaces   *[]string
	nExamples    *int64
	nIndexes     *int64
	since        *int64
}

type ListSuggestedIndexesApiParams struct {
	GroupId      string
	ProcessId    string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Duration     *int64
	Namespaces   *[]string
	NExamples    *int64
	NIndexes     *int64
	Since        *int64
}

func (a *PerformanceAdvisorApiService) ListSuggestedIndexesWithParams(ctx context.Context, args *ListSuggestedIndexesApiParams) ListSuggestedIndexesApiRequest {
	return ListSuggestedIndexesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		processId:    args.ProcessId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		duration:     args.Duration,
		namespaces:   args.Namespaces,
		nExamples:    args.NExamples,
		nIndexes:     args.NIndexes,
		since:        args.Since,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListSuggestedIndexesApiRequest) IncludeCount(includeCount bool) ListSuggestedIndexesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListSuggestedIndexesApiRequest) ItemsPerPage(itemsPerPage int) ListSuggestedIndexesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListSuggestedIndexesApiRequest) PageNum(pageNum int) ListSuggestedIndexesApiRequest {
	r.pageNum = &pageNum
	return r
}

// Length of time expressed during which the query finds suggested indexes among the managed namespaces in the cluster. This parameter expresses its value in milliseconds.  - If you don&#39;t specify the **since** parameter, the endpoint returns data covering the duration before the current time. - If you specify neither the **duration** nor **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListSuggestedIndexesApiRequest) Duration(duration int64) ListSuggestedIndexesApiRequest {
	r.duration = &duration
	return r
}

// Namespaces from which to retrieve suggested indexes. A namespace consists of one database and one collection resource written as &#x60;.&#x60;: &#x60;&lt;database&gt;.&lt;collection&gt;&#x60;. To include multiple namespaces, pass the parameter multiple times delimited with an ampersand (&#x60;&amp;&#x60;) between each namespace. Omit this parameter to return results for all namespaces.
func (r ListSuggestedIndexesApiRequest) Namespaces(namespaces []string) ListSuggestedIndexesApiRequest {
	r.namespaces = &namespaces
	return r
}

// Maximum number of example queries that benefit from the suggested index.
func (r ListSuggestedIndexesApiRequest) NExamples(nExamples int64) ListSuggestedIndexesApiRequest {
	r.nExamples = &nExamples
	return r
}

// Number that indicates the maximum indexes to suggest.
func (r ListSuggestedIndexesApiRequest) NIndexes(nIndexes int64) ListSuggestedIndexesApiRequest {
	r.nIndexes = &nIndexes
	return r
}

// Date and time from which the query retrieves the suggested indexes. This parameter expresses its value in the number of milliseconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).  - If you don&#39;t specify the **duration** parameter, the endpoint returns data covering from the **since** value and the current time. - If you specify neither the **duration** nor the **since** parameters, the endpoint returns data from the previous 24 hours.
func (r ListSuggestedIndexesApiRequest) Since(since int64) ListSuggestedIndexesApiRequest {
	r.since = &since
	return r
}

func (r ListSuggestedIndexesApiRequest) Execute() (*PerformanceAdvisorResponse, *http.Response, error) {
	return r.ApiService.ListSuggestedIndexesExecute(r)
}

/*
ListSuggestedIndexes Return All Suggested Indexes

Returns the indexes that the Performance Advisor suggests. The Performance Advisor monitors queries that MongoDB considers slow and suggests new indexes to improve query performance.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of host and port that serves the MongoDB process. The host must be the hostname, FQDN, IPv4 address, or IPv6 address of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return ListSuggestedIndexesApiRequest
*/
func (a *PerformanceAdvisorApiService) ListSuggestedIndexes(ctx context.Context, groupId string, processId string) ListSuggestedIndexesApiRequest {
	return ListSuggestedIndexesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// ListSuggestedIndexesExecute executes the request
//
//	@return PerformanceAdvisorResponse
func (a *PerformanceAdvisorApiService) ListSuggestedIndexesExecute(r ListSuggestedIndexesApiRequest) (*PerformanceAdvisorResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PerformanceAdvisorResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.ListSuggestedIndexes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/performanceAdvisor/suggestedIndexes"
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
	if r.duration != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "duration", r.duration, "")
	}
	if r.namespaces != nil {
		t := *r.namespaces
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "namespaces", t, "multi")

	}
	if r.nExamples != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nExamples", r.nExamples, "")
	} else {
		var defaultValue int64 = 5
		r.nExamples = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "nExamples", r.nExamples, "")
	}
	if r.nIndexes != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nIndexes", r.nIndexes, "")
	}
	if r.since != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "since", r.since, "")
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

type SetServerlessAutoIndexingApiRequest struct {
	ctx         context.Context
	ApiService  PerformanceAdvisorApi
	groupId     string
	clusterName string
	enable      *bool
}

type SetServerlessAutoIndexingApiParams struct {
	GroupId     string
	ClusterName string
	Enable      *bool
}

func (a *PerformanceAdvisorApiService) SetServerlessAutoIndexingWithParams(ctx context.Context, args *SetServerlessAutoIndexingApiParams) SetServerlessAutoIndexingApiRequest {
	return SetServerlessAutoIndexingApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		enable:      args.Enable,
	}
}

// Value that we want to set for the Serverless Auto Indexing toggle.
func (r SetServerlessAutoIndexingApiRequest) Enable(enable bool) SetServerlessAutoIndexingApiRequest {
	r.enable = &enable
	return r
}

func (r SetServerlessAutoIndexingApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.SetServerlessAutoIndexingExecute(r)
}

/*
SetServerlessAutoIndexing Set Serverless Auto-Indexing Status

Set whether the Serverless Auto Indexing feature is enabled. This endpoint sets a value for Flex clusters that were created with the Create Serverless Instance endpoint or Flex clusters that were migrated from Serverless instances. However, the value returned is not indicative of the Auto Indexing state as Auto Indexing is unavailable for Flex clusters. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return SetServerlessAutoIndexingApiRequest

Deprecated
*/
func (a *PerformanceAdvisorApiService) SetServerlessAutoIndexing(ctx context.Context, groupId string, clusterName string) SetServerlessAutoIndexingApiRequest {
	return SetServerlessAutoIndexingApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// SetServerlessAutoIndexingExecute executes the request
// Deprecated
func (a *PerformanceAdvisorApiService) SetServerlessAutoIndexingExecute(r SetServerlessAutoIndexingApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PerformanceAdvisorApiService.SetServerlessAutoIndexing")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/performanceAdvisor/autoIndexing"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.enable == nil {
		return nil, reportError("enable is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "enable", r.enable, "")
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
