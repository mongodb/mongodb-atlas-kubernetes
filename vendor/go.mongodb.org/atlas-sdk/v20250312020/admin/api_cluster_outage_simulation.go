// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ClusterOutageSimulationApi interface {

	/*
		EndOutageSimulation End One Outage Simulation

		Ends a cluster outage simulation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that is undergoing outage simulation.
		@return EndOutageSimulationApiRequest
	*/
	EndOutageSimulation(ctx context.Context, groupId string, clusterName string) EndOutageSimulationApiRequest
	/*
		EndOutageSimulation End One Outage Simulation


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param EndOutageSimulationApiParams - Parameters for the request
		@return EndOutageSimulationApiRequest
	*/
	EndOutageSimulationWithParams(ctx context.Context, args *EndOutageSimulationApiParams) EndOutageSimulationApiRequest

	// Method available only for mocking purposes
	EndOutageSimulationExecute(r EndOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error)

	/*
		GetOutageSimulation Return One Outage Simulation

		Returns one outage simulation for one cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that is undergoing outage simulation.
		@return GetOutageSimulationApiRequest
	*/
	GetOutageSimulation(ctx context.Context, groupId string, clusterName string) GetOutageSimulationApiRequest
	/*
		GetOutageSimulation Return One Outage Simulation


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOutageSimulationApiParams - Parameters for the request
		@return GetOutageSimulationApiRequest
	*/
	GetOutageSimulationWithParams(ctx context.Context, args *GetOutageSimulationApiParams) GetOutageSimulationApiRequest

	// Method available only for mocking purposes
	GetOutageSimulationExecute(r GetOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error)

	/*
		StartOutageSimulation Start One Outage Simulation

		Starts a cluster outage simulation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster to undergo an outage simulation.
		@param clusterOutageSimulation Describes the outage simulation.
		@return StartOutageSimulationApiRequest
	*/
	StartOutageSimulation(ctx context.Context, groupId string, clusterName string, clusterOutageSimulation *ClusterOutageSimulation) StartOutageSimulationApiRequest
	/*
		StartOutageSimulation Start One Outage Simulation


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param StartOutageSimulationApiParams - Parameters for the request
		@return StartOutageSimulationApiRequest
	*/
	StartOutageSimulationWithParams(ctx context.Context, args *StartOutageSimulationApiParams) StartOutageSimulationApiRequest

	// Method available only for mocking purposes
	StartOutageSimulationExecute(r StartOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error)
}

// ClusterOutageSimulationApiService ClusterOutageSimulationApi service
type ClusterOutageSimulationApiService service

type EndOutageSimulationApiRequest struct {
	ctx         context.Context
	ApiService  ClusterOutageSimulationApi
	groupId     string
	clusterName string
}

type EndOutageSimulationApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClusterOutageSimulationApiService) EndOutageSimulationWithParams(ctx context.Context, args *EndOutageSimulationApiParams) EndOutageSimulationApiRequest {
	return EndOutageSimulationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r EndOutageSimulationApiRequest) Execute() (*ClusterOutageSimulation, *http.Response, error) {
	return r.ApiService.EndOutageSimulationExecute(r)
}

/*
EndOutageSimulation End One Outage Simulation

Ends a cluster outage simulation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that is undergoing outage simulation.
	@return EndOutageSimulationApiRequest
*/
func (a *ClusterOutageSimulationApiService) EndOutageSimulation(ctx context.Context, groupId string, clusterName string) EndOutageSimulationApiRequest {
	return EndOutageSimulationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// EndOutageSimulationExecute executes the request
//
//	@return ClusterOutageSimulation
func (a *ClusterOutageSimulationApiService) EndOutageSimulationExecute(r EndOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterOutageSimulation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClusterOutageSimulationApiService.EndOutageSimulation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/outageSimulation"
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

type GetOutageSimulationApiRequest struct {
	ctx         context.Context
	ApiService  ClusterOutageSimulationApi
	groupId     string
	clusterName string
}

type GetOutageSimulationApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClusterOutageSimulationApiService) GetOutageSimulationWithParams(ctx context.Context, args *GetOutageSimulationApiParams) GetOutageSimulationApiRequest {
	return GetOutageSimulationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetOutageSimulationApiRequest) Execute() (*ClusterOutageSimulation, *http.Response, error) {
	return r.ApiService.GetOutageSimulationExecute(r)
}

/*
GetOutageSimulation Return One Outage Simulation

Returns one outage simulation for one cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that is undergoing outage simulation.
	@return GetOutageSimulationApiRequest
*/
func (a *ClusterOutageSimulationApiService) GetOutageSimulation(ctx context.Context, groupId string, clusterName string) GetOutageSimulationApiRequest {
	return GetOutageSimulationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetOutageSimulationExecute executes the request
//
//	@return ClusterOutageSimulation
func (a *ClusterOutageSimulationApiService) GetOutageSimulationExecute(r GetOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterOutageSimulation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClusterOutageSimulationApiService.GetOutageSimulation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/outageSimulation"
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

type StartOutageSimulationApiRequest struct {
	ctx                     context.Context
	ApiService              ClusterOutageSimulationApi
	groupId                 string
	clusterName             string
	clusterOutageSimulation *ClusterOutageSimulation
}

type StartOutageSimulationApiParams struct {
	GroupId                 string
	ClusterName             string
	ClusterOutageSimulation *ClusterOutageSimulation
}

func (a *ClusterOutageSimulationApiService) StartOutageSimulationWithParams(ctx context.Context, args *StartOutageSimulationApiParams) StartOutageSimulationApiRequest {
	return StartOutageSimulationApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 args.GroupId,
		clusterName:             args.ClusterName,
		clusterOutageSimulation: args.ClusterOutageSimulation,
	}
}

func (r StartOutageSimulationApiRequest) Execute() (*ClusterOutageSimulation, *http.Response, error) {
	return r.ApiService.StartOutageSimulationExecute(r)
}

/*
StartOutageSimulation Start One Outage Simulation

Starts a cluster outage simulation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster to undergo an outage simulation.
	@return StartOutageSimulationApiRequest
*/
func (a *ClusterOutageSimulationApiService) StartOutageSimulation(ctx context.Context, groupId string, clusterName string, clusterOutageSimulation *ClusterOutageSimulation) StartOutageSimulationApiRequest {
	return StartOutageSimulationApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 groupId,
		clusterName:             clusterName,
		clusterOutageSimulation: clusterOutageSimulation,
	}
}

// StartOutageSimulationExecute executes the request
//
//	@return ClusterOutageSimulation
func (a *ClusterOutageSimulationApiService) StartOutageSimulationExecute(r StartOutageSimulationApiRequest) (*ClusterOutageSimulation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterOutageSimulation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClusterOutageSimulationApiService.StartOutageSimulation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/outageSimulation"
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
	if r.clusterOutageSimulation == nil {
		return localVarReturnValue, nil, reportError("clusterOutageSimulation is required and must be specified")
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
	localVarPostBody = r.clusterOutageSimulation
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
