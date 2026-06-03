// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ServerlessPrivateEndpointsApi interface {

	/*
			CreateServerlessPrivateEndpoint Create One Private Endpoint for One Serverless Instance

			Creates one private endpoint for one serverless instance.

		 A new endpoint won't be immediately available after creation.  Read the steps in the linked tutorial for detailed guidance.

		This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param instanceName Human-readable label that identifies the serverless instance for which the tenant endpoint will be created.
			@param serverlessTenantCreateRequest Information about the Private Endpoint to create for the Serverless Instance.
			@return CreateServerlessPrivateEndpointApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	CreateServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, serverlessTenantCreateRequest *ServerlessTenantCreateRequest) CreateServerlessPrivateEndpointApiRequest
	/*
		CreateServerlessPrivateEndpoint Create One Private Endpoint for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateServerlessPrivateEndpointApiParams - Parameters for the request
		@return CreateServerlessPrivateEndpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	CreateServerlessPrivateEndpointWithParams(ctx context.Context, args *CreateServerlessPrivateEndpointApiParams) CreateServerlessPrivateEndpointApiRequest

	// Method available only for mocking purposes
	CreateServerlessPrivateEndpointExecute(r CreateServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error)

	/*
			DeleteServerlessPrivateEndpoint Remove One Private Endpoint for One Serverless Instance

			Remove one private endpoint from one serverless instance.

		This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param instanceName Human-readable label that identifies the serverless instance from which the tenant endpoint will be removed.
			@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint which will be removed.
			@return DeleteServerlessPrivateEndpointApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	DeleteServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string) DeleteServerlessPrivateEndpointApiRequest
	/*
		DeleteServerlessPrivateEndpoint Remove One Private Endpoint for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteServerlessPrivateEndpointApiParams - Parameters for the request
		@return DeleteServerlessPrivateEndpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	DeleteServerlessPrivateEndpointWithParams(ctx context.Context, args *DeleteServerlessPrivateEndpointApiParams) DeleteServerlessPrivateEndpointApiRequest

	// Method available only for mocking purposes
	DeleteServerlessPrivateEndpointExecute(r DeleteServerlessPrivateEndpointApiRequest) (*http.Response, error)

	/*
			GetServerlessPrivateEndpoint Return One Private Endpoint for One Serverless Instance

			Return one private endpoint for one serverless instance. Identify this endpoint using its unique ID. You must have at least the Project Read Only role for the project to successfully call this resource.

		This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint.
			@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint.
			@return GetServerlessPrivateEndpointApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	GetServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string) GetServerlessPrivateEndpointApiRequest
	/*
		GetServerlessPrivateEndpoint Return One Private Endpoint for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServerlessPrivateEndpointApiParams - Parameters for the request
		@return GetServerlessPrivateEndpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	GetServerlessPrivateEndpointWithParams(ctx context.Context, args *GetServerlessPrivateEndpointApiParams) GetServerlessPrivateEndpointApiRequest

	// Method available only for mocking purposes
	GetServerlessPrivateEndpointExecute(r GetServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error)

	/*
			ListServerlessPrivateEndpoint Return All Private Endpoints for One Serverless Instance

			Returns all private endpoints for one serverless instance. You must have at least the Project Read Only role for the project to successfully call this resource.

		This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint.
			@return ListServerlessPrivateEndpointApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	ListServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string) ListServerlessPrivateEndpointApiRequest
	/*
		ListServerlessPrivateEndpoint Return All Private Endpoints for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListServerlessPrivateEndpointApiParams - Parameters for the request
		@return ListServerlessPrivateEndpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	ListServerlessPrivateEndpointWithParams(ctx context.Context, args *ListServerlessPrivateEndpointApiParams) ListServerlessPrivateEndpointApiRequest

	// Method available only for mocking purposes
	ListServerlessPrivateEndpointExecute(r ListServerlessPrivateEndpointApiRequest) ([]ServerlessTenantEndpoint, *http.Response, error)

	/*
			UpdateServerlessPrivateEndpoint Update One Private Endpoint for One Serverless Instance

			Updates one private endpoint for one serverless instance.

		This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint that will be updated.
			@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint which will be updated.
			@param serverlessTenantEndpointUpdate Object used for update.
			@return UpdateServerlessPrivateEndpointApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	UpdateServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string, serverlessTenantEndpointUpdate *ServerlessTenantEndpointUpdate) UpdateServerlessPrivateEndpointApiRequest
	/*
		UpdateServerlessPrivateEndpoint Update One Private Endpoint for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateServerlessPrivateEndpointApiParams - Parameters for the request
		@return UpdateServerlessPrivateEndpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ServerlessPrivateEndpointsApi
	*/
	UpdateServerlessPrivateEndpointWithParams(ctx context.Context, args *UpdateServerlessPrivateEndpointApiParams) UpdateServerlessPrivateEndpointApiRequest

	// Method available only for mocking purposes
	UpdateServerlessPrivateEndpointExecute(r UpdateServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error)
}

// ServerlessPrivateEndpointsApiService ServerlessPrivateEndpointsApi service
type ServerlessPrivateEndpointsApiService service

type CreateServerlessPrivateEndpointApiRequest struct {
	ctx                           context.Context
	ApiService                    ServerlessPrivateEndpointsApi
	groupId                       string
	instanceName                  string
	serverlessTenantCreateRequest *ServerlessTenantCreateRequest
}

type CreateServerlessPrivateEndpointApiParams struct {
	GroupId                       string
	InstanceName                  string
	ServerlessTenantCreateRequest *ServerlessTenantCreateRequest
}

func (a *ServerlessPrivateEndpointsApiService) CreateServerlessPrivateEndpointWithParams(ctx context.Context, args *CreateServerlessPrivateEndpointApiParams) CreateServerlessPrivateEndpointApiRequest {
	return CreateServerlessPrivateEndpointApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       args.GroupId,
		instanceName:                  args.InstanceName,
		serverlessTenantCreateRequest: args.ServerlessTenantCreateRequest,
	}
}

func (r CreateServerlessPrivateEndpointApiRequest) Execute() (*ServerlessTenantEndpoint, *http.Response, error) {
	return r.ApiService.CreateServerlessPrivateEndpointExecute(r)
}

/*
CreateServerlessPrivateEndpoint Create One Private Endpoint for One Serverless Instance

Creates one private endpoint for one serverless instance.

	A new endpoint won't be immediately available after creation.  Read the steps in the linked tutorial for detailed guidance.

This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param instanceName Human-readable label that identifies the serverless instance for which the tenant endpoint will be created.
	@return CreateServerlessPrivateEndpointApiRequest

Deprecated
*/
func (a *ServerlessPrivateEndpointsApiService) CreateServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, serverlessTenantCreateRequest *ServerlessTenantCreateRequest) CreateServerlessPrivateEndpointApiRequest {
	return CreateServerlessPrivateEndpointApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       groupId,
		instanceName:                  instanceName,
		serverlessTenantCreateRequest: serverlessTenantCreateRequest,
	}
}

// CreateServerlessPrivateEndpointExecute executes the request
//
//	@return ServerlessTenantEndpoint
//
// Deprecated
func (a *ServerlessPrivateEndpointsApiService) CreateServerlessPrivateEndpointExecute(r CreateServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessTenantEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessPrivateEndpointsApiService.CreateServerlessPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateEndpoint/serverless/instance/{instanceName}/endpoint"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.instanceName == "" {
		return localVarReturnValue, nil, reportError("instanceName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(r.instanceName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serverlessTenantCreateRequest == nil {
		return localVarReturnValue, nil, reportError("serverlessTenantCreateRequest is required and must be specified")
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
	localVarPostBody = r.serverlessTenantCreateRequest
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

type DeleteServerlessPrivateEndpointApiRequest struct {
	ctx          context.Context
	ApiService   ServerlessPrivateEndpointsApi
	groupId      string
	instanceName string
	endpointId   string
}

type DeleteServerlessPrivateEndpointApiParams struct {
	GroupId      string
	InstanceName string
	EndpointId   string
}

func (a *ServerlessPrivateEndpointsApiService) DeleteServerlessPrivateEndpointWithParams(ctx context.Context, args *DeleteServerlessPrivateEndpointApiParams) DeleteServerlessPrivateEndpointApiRequest {
	return DeleteServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		instanceName: args.InstanceName,
		endpointId:   args.EndpointId,
	}
}

func (r DeleteServerlessPrivateEndpointApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteServerlessPrivateEndpointExecute(r)
}

/*
DeleteServerlessPrivateEndpoint Remove One Private Endpoint for One Serverless Instance

Remove one private endpoint from one serverless instance.

This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param instanceName Human-readable label that identifies the serverless instance from which the tenant endpoint will be removed.
	@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint which will be removed.
	@return DeleteServerlessPrivateEndpointApiRequest

Deprecated
*/
func (a *ServerlessPrivateEndpointsApiService) DeleteServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string) DeleteServerlessPrivateEndpointApiRequest {
	return DeleteServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		instanceName: instanceName,
		endpointId:   endpointId,
	}
}

// DeleteServerlessPrivateEndpointExecute executes the request
// Deprecated
func (a *ServerlessPrivateEndpointsApiService) DeleteServerlessPrivateEndpointExecute(r DeleteServerlessPrivateEndpointApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessPrivateEndpointsApiService.DeleteServerlessPrivateEndpoint")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateEndpoint/serverless/instance/{instanceName}/endpoint/{endpointId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.instanceName == "" {
		return nil, reportError("instanceName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(r.instanceName), -1)
	if r.endpointId == "" {
		return nil, reportError("endpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointId"+"}", url.PathEscape(r.endpointId), -1)

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

type GetServerlessPrivateEndpointApiRequest struct {
	ctx          context.Context
	ApiService   ServerlessPrivateEndpointsApi
	groupId      string
	instanceName string
	endpointId   string
}

type GetServerlessPrivateEndpointApiParams struct {
	GroupId      string
	InstanceName string
	EndpointId   string
}

func (a *ServerlessPrivateEndpointsApiService) GetServerlessPrivateEndpointWithParams(ctx context.Context, args *GetServerlessPrivateEndpointApiParams) GetServerlessPrivateEndpointApiRequest {
	return GetServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		instanceName: args.InstanceName,
		endpointId:   args.EndpointId,
	}
}

func (r GetServerlessPrivateEndpointApiRequest) Execute() (*ServerlessTenantEndpoint, *http.Response, error) {
	return r.ApiService.GetServerlessPrivateEndpointExecute(r)
}

/*
GetServerlessPrivateEndpoint Return One Private Endpoint for One Serverless Instance

Return one private endpoint for one serverless instance. Identify this endpoint using its unique ID. You must have at least the Project Read Only role for the project to successfully call this resource.

This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint.
	@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint.
	@return GetServerlessPrivateEndpointApiRequest

Deprecated
*/
func (a *ServerlessPrivateEndpointsApiService) GetServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string) GetServerlessPrivateEndpointApiRequest {
	return GetServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		instanceName: instanceName,
		endpointId:   endpointId,
	}
}

// GetServerlessPrivateEndpointExecute executes the request
//
//	@return ServerlessTenantEndpoint
//
// Deprecated
func (a *ServerlessPrivateEndpointsApiService) GetServerlessPrivateEndpointExecute(r GetServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessTenantEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessPrivateEndpointsApiService.GetServerlessPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateEndpoint/serverless/instance/{instanceName}/endpoint/{endpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.instanceName == "" {
		return localVarReturnValue, nil, reportError("instanceName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(r.instanceName), -1)
	if r.endpointId == "" {
		return localVarReturnValue, nil, reportError("endpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointId"+"}", url.PathEscape(r.endpointId), -1)

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

type ListServerlessPrivateEndpointApiRequest struct {
	ctx          context.Context
	ApiService   ServerlessPrivateEndpointsApi
	groupId      string
	instanceName string
}

type ListServerlessPrivateEndpointApiParams struct {
	GroupId      string
	InstanceName string
}

func (a *ServerlessPrivateEndpointsApiService) ListServerlessPrivateEndpointWithParams(ctx context.Context, args *ListServerlessPrivateEndpointApiParams) ListServerlessPrivateEndpointApiRequest {
	return ListServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		instanceName: args.InstanceName,
	}
}

func (r ListServerlessPrivateEndpointApiRequest) Execute() ([]ServerlessTenantEndpoint, *http.Response, error) {
	return r.ApiService.ListServerlessPrivateEndpointExecute(r)
}

/*
ListServerlessPrivateEndpoint Return All Private Endpoints for One Serverless Instance

Returns all private endpoints for one serverless instance. You must have at least the Project Read Only role for the project to successfully call this resource.

This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint.
	@return ListServerlessPrivateEndpointApiRequest

Deprecated
*/
func (a *ServerlessPrivateEndpointsApiService) ListServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string) ListServerlessPrivateEndpointApiRequest {
	return ListServerlessPrivateEndpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		instanceName: instanceName,
	}
}

// ListServerlessPrivateEndpointExecute executes the request
//
//	@return []ServerlessTenantEndpoint
//
// Deprecated
func (a *ServerlessPrivateEndpointsApiService) ListServerlessPrivateEndpointExecute(r ListServerlessPrivateEndpointApiRequest) ([]ServerlessTenantEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []ServerlessTenantEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessPrivateEndpointsApiService.ListServerlessPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateEndpoint/serverless/instance/{instanceName}/endpoint"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.instanceName == "" {
		return localVarReturnValue, nil, reportError("instanceName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(r.instanceName), -1)

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

type UpdateServerlessPrivateEndpointApiRequest struct {
	ctx                            context.Context
	ApiService                     ServerlessPrivateEndpointsApi
	groupId                        string
	instanceName                   string
	endpointId                     string
	serverlessTenantEndpointUpdate *ServerlessTenantEndpointUpdate
}

type UpdateServerlessPrivateEndpointApiParams struct {
	GroupId                        string
	InstanceName                   string
	EndpointId                     string
	ServerlessTenantEndpointUpdate *ServerlessTenantEndpointUpdate
}

func (a *ServerlessPrivateEndpointsApiService) UpdateServerlessPrivateEndpointWithParams(ctx context.Context, args *UpdateServerlessPrivateEndpointApiParams) UpdateServerlessPrivateEndpointApiRequest {
	return UpdateServerlessPrivateEndpointApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        args.GroupId,
		instanceName:                   args.InstanceName,
		endpointId:                     args.EndpointId,
		serverlessTenantEndpointUpdate: args.ServerlessTenantEndpointUpdate,
	}
}

func (r UpdateServerlessPrivateEndpointApiRequest) Execute() (*ServerlessTenantEndpoint, *http.Response, error) {
	return r.ApiService.UpdateServerlessPrivateEndpointExecute(r)
}

/*
UpdateServerlessPrivateEndpoint Update One Private Endpoint for One Serverless Instance

Updates one private endpoint for one serverless instance.

This feature does not work for Flex clusters. To continue using Private Endpoints once Serverless is replaced by Flex, please use a Dedicated cluster instead. This endpoint will be sunset on January 22, 2026.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param instanceName Human-readable label that identifies the serverless instance associated with the tenant endpoint that will be updated.
	@param endpointId Unique 24-hexadecimal digit string that identifies the tenant endpoint which will be updated.
	@return UpdateServerlessPrivateEndpointApiRequest

Deprecated
*/
func (a *ServerlessPrivateEndpointsApiService) UpdateServerlessPrivateEndpoint(ctx context.Context, groupId string, instanceName string, endpointId string, serverlessTenantEndpointUpdate *ServerlessTenantEndpointUpdate) UpdateServerlessPrivateEndpointApiRequest {
	return UpdateServerlessPrivateEndpointApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        groupId,
		instanceName:                   instanceName,
		endpointId:                     endpointId,
		serverlessTenantEndpointUpdate: serverlessTenantEndpointUpdate,
	}
}

// UpdateServerlessPrivateEndpointExecute executes the request
//
//	@return ServerlessTenantEndpoint
//
// Deprecated
func (a *ServerlessPrivateEndpointsApiService) UpdateServerlessPrivateEndpointExecute(r UpdateServerlessPrivateEndpointApiRequest) (*ServerlessTenantEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessTenantEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServerlessPrivateEndpointsApiService.UpdateServerlessPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateEndpoint/serverless/instance/{instanceName}/endpoint/{endpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.instanceName == "" {
		return localVarReturnValue, nil, reportError("instanceName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(r.instanceName), -1)
	if r.endpointId == "" {
		return localVarReturnValue, nil, reportError("endpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointId"+"}", url.PathEscape(r.endpointId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serverlessTenantEndpointUpdate == nil {
		return localVarReturnValue, nil, reportError("serverlessTenantEndpointUpdate is required and must be specified")
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
	localVarPostBody = r.serverlessTenantEndpointUpdate
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
