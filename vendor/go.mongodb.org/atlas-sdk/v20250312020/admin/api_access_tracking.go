// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AccessTrackingApi interface {

	/*
		GetAccessHistoryCluster Return Database Access History for One Cluster by Cluster Name

		Returns the access logs of one cluster identified by the cluster's name. Access logs contain a list of authentication requests made against your cluster. You can't use this feature on tenant-tier clusters (M0, M2, M5).

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return GetAccessHistoryClusterApiRequest
	*/
	GetAccessHistoryCluster(ctx context.Context, groupId string, clusterName string) GetAccessHistoryClusterApiRequest
	/*
		GetAccessHistoryCluster Return Database Access History for One Cluster by Cluster Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAccessHistoryClusterApiParams - Parameters for the request
		@return GetAccessHistoryClusterApiRequest
	*/
	GetAccessHistoryClusterWithParams(ctx context.Context, args *GetAccessHistoryClusterApiParams) GetAccessHistoryClusterApiRequest

	// Method available only for mocking purposes
	GetAccessHistoryClusterExecute(r GetAccessHistoryClusterApiRequest) (*MongoDBAccessLogsList, *http.Response, error)

	/*
		GetAccessHistoryProcess Return Database Access History for One Cluster by Hostname

		Returns the access logs of one cluster identified by the cluster's hostname. Access logs contain a list of authentication requests made against your clusters. You can't use this feature on tenant-tier clusters (M0, M2, M5).

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param hostname Fully qualified domain name or IP address of the MongoDB host that stores the log files that you want to download.
		@return GetAccessHistoryProcessApiRequest
	*/
	GetAccessHistoryProcess(ctx context.Context, groupId string, hostname string) GetAccessHistoryProcessApiRequest
	/*
		GetAccessHistoryProcess Return Database Access History for One Cluster by Hostname


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAccessHistoryProcessApiParams - Parameters for the request
		@return GetAccessHistoryProcessApiRequest
	*/
	GetAccessHistoryProcessWithParams(ctx context.Context, args *GetAccessHistoryProcessApiParams) GetAccessHistoryProcessApiRequest

	// Method available only for mocking purposes
	GetAccessHistoryProcessExecute(r GetAccessHistoryProcessApiRequest) (*MongoDBAccessLogsList, *http.Response, error)
}

// AccessTrackingApiService AccessTrackingApi service
type AccessTrackingApiService service

type GetAccessHistoryClusterApiRequest struct {
	ctx         context.Context
	ApiService  AccessTrackingApi
	groupId     string
	clusterName string
	authResult  *bool
	end         *int64
	ipAddress   *string
	nLogs       *int
	start       *int64
}

type GetAccessHistoryClusterApiParams struct {
	GroupId     string
	ClusterName string
	AuthResult  *bool
	End         *int64
	IpAddress   *string
	NLogs       *int
	Start       *int64
}

func (a *AccessTrackingApiService) GetAccessHistoryClusterWithParams(ctx context.Context, args *GetAccessHistoryClusterApiParams) GetAccessHistoryClusterApiRequest {
	return GetAccessHistoryClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		authResult:  args.AuthResult,
		end:         args.End,
		ipAddress:   args.IpAddress,
		nLogs:       args.NLogs,
		start:       args.Start,
	}
}

// Flag that indicates whether the response returns the successful authentication attempts only.
func (r GetAccessHistoryClusterApiRequest) AuthResult(authResult bool) GetAccessHistoryClusterApiRequest {
	r.authResult = &authResult
	return r
}

// Date and time when to stop retrieving database history. If you specify **end**, you must also specify **start**. This parameter uses UNIX epoch time in milliseconds.
func (r GetAccessHistoryClusterApiRequest) End(end int64) GetAccessHistoryClusterApiRequest {
	r.end = &end
	return r
}

// One Internet Protocol address that attempted to authenticate with the database.
func (r GetAccessHistoryClusterApiRequest) IpAddress(ipAddress string) GetAccessHistoryClusterApiRequest {
	r.ipAddress = &ipAddress
	return r
}

// Maximum number of lines from the log to return.
func (r GetAccessHistoryClusterApiRequest) NLogs(nLogs int) GetAccessHistoryClusterApiRequest {
	r.nLogs = &nLogs
	return r
}

// Date and time when MongoDB Cloud begins retrieving database history. If you specify **start**, you must also specify **end**. This parameter uses UNIX epoch time in milliseconds.
func (r GetAccessHistoryClusterApiRequest) Start(start int64) GetAccessHistoryClusterApiRequest {
	r.start = &start
	return r
}

func (r GetAccessHistoryClusterApiRequest) Execute() (*MongoDBAccessLogsList, *http.Response, error) {
	return r.ApiService.GetAccessHistoryClusterExecute(r)
}

/*
GetAccessHistoryCluster Return Database Access History for One Cluster by Cluster Name

Returns the access logs of one cluster identified by the cluster's name. Access logs contain a list of authentication requests made against your cluster. You can't use this feature on tenant-tier clusters (M0, M2, M5).

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return GetAccessHistoryClusterApiRequest
*/
func (a *AccessTrackingApiService) GetAccessHistoryCluster(ctx context.Context, groupId string, clusterName string) GetAccessHistoryClusterApiRequest {
	return GetAccessHistoryClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetAccessHistoryClusterExecute executes the request
//
//	@return MongoDBAccessLogsList
func (a *AccessTrackingApiService) GetAccessHistoryClusterExecute(r GetAccessHistoryClusterApiRequest) (*MongoDBAccessLogsList, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MongoDBAccessLogsList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AccessTrackingApiService.GetAccessHistoryCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dbAccessHistory/clusters/{clusterName}"
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

	if r.authResult != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "authResult", r.authResult, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.ipAddress != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "ipAddress", r.ipAddress, "")
	}
	if r.nLogs != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	} else {
		var defaultValue int = 20000
		r.nLogs = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
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

type GetAccessHistoryProcessApiRequest struct {
	ctx        context.Context
	ApiService AccessTrackingApi
	groupId    string
	hostname   string
	authResult *bool
	end        *int64
	ipAddress  *string
	nLogs      *int
	start      *int64
}

type GetAccessHistoryProcessApiParams struct {
	GroupId    string
	Hostname   string
	AuthResult *bool
	End        *int64
	IpAddress  *string
	NLogs      *int
	Start      *int64
}

func (a *AccessTrackingApiService) GetAccessHistoryProcessWithParams(ctx context.Context, args *GetAccessHistoryProcessApiParams) GetAccessHistoryProcessApiRequest {
	return GetAccessHistoryProcessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		hostname:   args.Hostname,
		authResult: args.AuthResult,
		end:        args.End,
		ipAddress:  args.IpAddress,
		nLogs:      args.NLogs,
		start:      args.Start,
	}
}

// Flag that indicates whether the response returns the successful authentication attempts only.
func (r GetAccessHistoryProcessApiRequest) AuthResult(authResult bool) GetAccessHistoryProcessApiRequest {
	r.authResult = &authResult
	return r
}

// Date and time when to stop retrieving database history. If you specify **end**, you must also specify **start**. This parameter uses UNIX epoch time in milliseconds.
func (r GetAccessHistoryProcessApiRequest) End(end int64) GetAccessHistoryProcessApiRequest {
	r.end = &end
	return r
}

// One Internet Protocol address that attempted to authenticate with the database.
func (r GetAccessHistoryProcessApiRequest) IpAddress(ipAddress string) GetAccessHistoryProcessApiRequest {
	r.ipAddress = &ipAddress
	return r
}

// Maximum number of lines from the log to return.
func (r GetAccessHistoryProcessApiRequest) NLogs(nLogs int) GetAccessHistoryProcessApiRequest {
	r.nLogs = &nLogs
	return r
}

// Date and time when MongoDB Cloud begins retrieving database history. If you specify **start**, you must also specify **end**. This parameter uses UNIX epoch time in milliseconds.
func (r GetAccessHistoryProcessApiRequest) Start(start int64) GetAccessHistoryProcessApiRequest {
	r.start = &start
	return r
}

func (r GetAccessHistoryProcessApiRequest) Execute() (*MongoDBAccessLogsList, *http.Response, error) {
	return r.ApiService.GetAccessHistoryProcessExecute(r)
}

/*
GetAccessHistoryProcess Return Database Access History for One Cluster by Hostname

Returns the access logs of one cluster identified by the cluster's hostname. Access logs contain a list of authentication requests made against your clusters. You can't use this feature on tenant-tier clusters (M0, M2, M5).

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param hostname Fully qualified domain name or IP address of the MongoDB host that stores the log files that you want to download.
	@return GetAccessHistoryProcessApiRequest
*/
func (a *AccessTrackingApiService) GetAccessHistoryProcess(ctx context.Context, groupId string, hostname string) GetAccessHistoryProcessApiRequest {
	return GetAccessHistoryProcessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		hostname:   hostname,
	}
}

// GetAccessHistoryProcessExecute executes the request
//
//	@return MongoDBAccessLogsList
func (a *AccessTrackingApiService) GetAccessHistoryProcessExecute(r GetAccessHistoryProcessApiRequest) (*MongoDBAccessLogsList, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MongoDBAccessLogsList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AccessTrackingApiService.GetAccessHistoryProcess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dbAccessHistory/processes/{hostname}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.hostname == "" {
		return localVarReturnValue, nil, reportError("hostname is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"hostname"+"}", url.PathEscape(r.hostname), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.authResult != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "authResult", r.authResult, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	if r.ipAddress != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "ipAddress", r.ipAddress, "")
	}
	if r.nLogs != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	} else {
		var defaultValue int = 20000
		r.nLogs = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "nLogs", r.nLogs, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
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
