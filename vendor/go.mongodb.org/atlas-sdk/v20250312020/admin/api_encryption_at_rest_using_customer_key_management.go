// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type EncryptionAtRestUsingCustomerKeyManagementApi interface {

	/*
		CreateRestPrivateEndpoint Create One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

		Creates a private endpoint in the specified region for encryption at rest using customer key management.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoint to create.
		@param eARPrivateEndpoint Creates a private endpoint in the specified region.
		@return CreateRestPrivateEndpointApiRequest
	*/
	CreateRestPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, eARPrivateEndpoint *EARPrivateEndpoint) CreateRestPrivateEndpointApiRequest
	/*
		CreateRestPrivateEndpoint Create One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateRestPrivateEndpointApiParams - Parameters for the request
		@return CreateRestPrivateEndpointApiRequest
	*/
	CreateRestPrivateEndpointWithParams(ctx context.Context, args *CreateRestPrivateEndpointApiParams) CreateRestPrivateEndpointApiRequest

	// Method available only for mocking purposes
	CreateRestPrivateEndpointExecute(r CreateRestPrivateEndpointApiRequest) (*EARPrivateEndpoint, *http.Response, error)

	/*
			GetEncryptionAtRest Return One Configuration for Encryption at Rest Using Customer-Managed Keys for One Project

			Returns the configuration for encryption at rest using the keys you manage through your cloud provider. MongoDB Cloud encrypts all storage even if you don't use your own key management.

		**LIMITED TO M10 OR GREATER:** MongoDB Cloud limits this feature to dedicated cluster tiers of M10 and greater.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return GetEncryptionAtRestApiRequest
	*/
	GetEncryptionAtRest(ctx context.Context, groupId string) GetEncryptionAtRestApiRequest
	/*
		GetEncryptionAtRest Return One Configuration for Encryption at Rest Using Customer-Managed Keys for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetEncryptionAtRestApiParams - Parameters for the request
		@return GetEncryptionAtRestApiRequest
	*/
	GetEncryptionAtRestWithParams(ctx context.Context, args *GetEncryptionAtRestApiParams) GetEncryptionAtRestApiRequest

	// Method available only for mocking purposes
	GetEncryptionAtRestExecute(r GetEncryptionAtRestApiRequest) (*EncryptionAtRest, *http.Response, error)

	/*
		GetRestPrivateEndpoint Return One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

		Returns one private endpoint, identified by its ID, for encryption at rest using Customer Key Management.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint.
		@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint.
		@return GetRestPrivateEndpointApiRequest
	*/
	GetRestPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) GetRestPrivateEndpointApiRequest
	/*
		GetRestPrivateEndpoint Return One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetRestPrivateEndpointApiParams - Parameters for the request
		@return GetRestPrivateEndpointApiRequest
	*/
	GetRestPrivateEndpointWithParams(ctx context.Context, args *GetRestPrivateEndpointApiParams) GetRestPrivateEndpointApiRequest

	// Method available only for mocking purposes
	GetRestPrivateEndpointExecute(r GetRestPrivateEndpointApiRequest) (*EARPrivateEndpoint, *http.Response, error)

	/*
		ListRestPrivateEndpoints Return Private Endpoints for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

		Returns the private endpoints of the specified cloud provider for encryption at rest using customer key management.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoints to return.
		@return ListRestPrivateEndpointsApiRequest
	*/
	ListRestPrivateEndpoints(ctx context.Context, groupId string, cloudProvider string) ListRestPrivateEndpointsApiRequest
	/*
		ListRestPrivateEndpoints Return Private Endpoints for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListRestPrivateEndpointsApiParams - Parameters for the request
		@return ListRestPrivateEndpointsApiRequest
	*/
	ListRestPrivateEndpointsWithParams(ctx context.Context, args *ListRestPrivateEndpointsApiParams) ListRestPrivateEndpointsApiRequest

	// Method available only for mocking purposes
	ListRestPrivateEndpointsExecute(r ListRestPrivateEndpointsApiRequest) (*PaginatedApiAtlasEARPrivateEndpoint, *http.Response, error)

	/*
		RequestPrivateEndpointDeletion Delete One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider from One Project

		Deletes one private endpoint, identified by its ID, for encryption at rest using Customer Key Management.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint to delete.
		@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint to delete.
		@return RequestPrivateEndpointDeletionApiRequest
	*/
	RequestPrivateEndpointDeletion(ctx context.Context, groupId string, cloudProvider string, endpointId string) RequestPrivateEndpointDeletionApiRequest
	/*
		RequestPrivateEndpointDeletion Delete One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RequestPrivateEndpointDeletionApiParams - Parameters for the request
		@return RequestPrivateEndpointDeletionApiRequest
	*/
	RequestPrivateEndpointDeletionWithParams(ctx context.Context, args *RequestPrivateEndpointDeletionApiParams) RequestPrivateEndpointDeletionApiRequest

	// Method available only for mocking purposes
	RequestPrivateEndpointDeletionExecute(r RequestPrivateEndpointDeletionApiRequest) (*http.Response, error)

	/*
			UpdateEncryptionAtRest Update Encryption at Rest Configuration in One Project

			Updates the configuration for encryption at rest using the keys you manage through your cloud provider. MongoDB Cloud encrypts all storage even if you don't use your own key management. This feature isn't available for `M0` free clusters, `M2`, `M5`, or serverless clusters.

		 After you configure at least one Encryption at Rest using a Customer Key Management provider for the MongoDB Cloud project, Project Owners can enable Encryption at Rest using Customer Key Management for each MongoDB Cloud cluster for which they require encryption. The Encryption at Rest using Customer Key Management provider doesn't have to match the cluster cloud service provider. MongoDB Cloud doesn't automatically rotate user-managed encryption keys. Defer to your preferred Encryption at Rest using Customer Key Management provider's documentation and guidance for best practices on key rotation. MongoDB Cloud automatically creates a 90-day key rotation alert when you configure Encryption at Rest using Customer Key Management using your Key Management in an MongoDB Cloud project. MongoDB Cloud encrypts all storage whether or not you use your own key management.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param encryptionAtRest Required parameters depend on whether someone has enabled Encryption at Rest using Customer Key Management:  If you have enabled Encryption at Rest using Customer Key Management (CMK), Atlas requires all of the parameters for the desired encryption provider.  - To use AWS Key Management Service (KMS), MongoDB Cloud requires all the fields in the `awsKms` object. - To use Azure Key Vault, MongoDB Cloud requires all the fields in the `azureKeyVault` object. - To use Google Cloud Key Management Service (KMS), MongoDB Cloud requires all the fields in the `googleCloudKms` object For authentication, you must provide either `serviceAccountKey` (static credentials) or `roleId` (service-account–based authentication) Once `roleId` is configured, `serviceAccountKey` is no longer supported.  If you enabled Encryption at Rest using Customer Key Management, administrators can pass only the changed fields for the `awsKms`, `azureKeyVault`, or `googleCloudKms` object to update the configuration to this endpoint.
			@return UpdateEncryptionAtRestApiRequest
	*/
	UpdateEncryptionAtRest(ctx context.Context, groupId string, encryptionAtRest *EncryptionAtRest) UpdateEncryptionAtRestApiRequest
	/*
		UpdateEncryptionAtRest Update Encryption at Rest Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateEncryptionAtRestApiParams - Parameters for the request
		@return UpdateEncryptionAtRestApiRequest
	*/
	UpdateEncryptionAtRestWithParams(ctx context.Context, args *UpdateEncryptionAtRestApiParams) UpdateEncryptionAtRestApiRequest

	// Method available only for mocking purposes
	UpdateEncryptionAtRestExecute(r UpdateEncryptionAtRestApiRequest) (*EncryptionAtRest, *http.Response, error)
}

// EncryptionAtRestUsingCustomerKeyManagementApiService EncryptionAtRestUsingCustomerKeyManagementApi service
type EncryptionAtRestUsingCustomerKeyManagementApiService service

type CreateRestPrivateEndpointApiRequest struct {
	ctx                context.Context
	ApiService         EncryptionAtRestUsingCustomerKeyManagementApi
	groupId            string
	cloudProvider      string
	eARPrivateEndpoint *EARPrivateEndpoint
}

type CreateRestPrivateEndpointApiParams struct {
	GroupId            string
	CloudProvider      string
	EARPrivateEndpoint *EARPrivateEndpoint
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) CreateRestPrivateEndpointWithParams(ctx context.Context, args *CreateRestPrivateEndpointApiParams) CreateRestPrivateEndpointApiRequest {
	return CreateRestPrivateEndpointApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		cloudProvider:      args.CloudProvider,
		eARPrivateEndpoint: args.EARPrivateEndpoint,
	}
}

func (r CreateRestPrivateEndpointApiRequest) Execute() (*EARPrivateEndpoint, *http.Response, error) {
	return r.ApiService.CreateRestPrivateEndpointExecute(r)
}

/*
CreateRestPrivateEndpoint Create One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

Creates a private endpoint in the specified region for encryption at rest using customer key management.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoint to create.
	@return CreateRestPrivateEndpointApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) CreateRestPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, eARPrivateEndpoint *EARPrivateEndpoint) CreateRestPrivateEndpointApiRequest {
	return CreateRestPrivateEndpointApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            groupId,
		cloudProvider:      cloudProvider,
		eARPrivateEndpoint: eARPrivateEndpoint,
	}
}

// CreateRestPrivateEndpointExecute executes the request
//
//	@return EARPrivateEndpoint
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) CreateRestPrivateEndpointExecute(r CreateRestPrivateEndpointApiRequest) (*EARPrivateEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EARPrivateEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.CreateRestPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest/{cloudProvider}/privateEndpoints"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.eARPrivateEndpoint == nil {
		return localVarReturnValue, nil, reportError("eARPrivateEndpoint is required and must be specified")
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
	localVarPostBody = r.eARPrivateEndpoint
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

type GetEncryptionAtRestApiRequest struct {
	ctx        context.Context
	ApiService EncryptionAtRestUsingCustomerKeyManagementApi
	groupId    string
}

type GetEncryptionAtRestApiParams struct {
	GroupId string
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetEncryptionAtRestWithParams(ctx context.Context, args *GetEncryptionAtRestApiParams) GetEncryptionAtRestApiRequest {
	return GetEncryptionAtRestApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetEncryptionAtRestApiRequest) Execute() (*EncryptionAtRest, *http.Response, error) {
	return r.ApiService.GetEncryptionAtRestExecute(r)
}

/*
GetEncryptionAtRest Return One Configuration for Encryption at Rest Using Customer-Managed Keys for One Project

Returns the configuration for encryption at rest using the keys you manage through your cloud provider. MongoDB Cloud encrypts all storage even if you don't use your own key management.

**LIMITED TO M10 OR GREATER:** MongoDB Cloud limits this feature to dedicated cluster tiers of M10 and greater.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetEncryptionAtRestApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetEncryptionAtRest(ctx context.Context, groupId string) GetEncryptionAtRestApiRequest {
	return GetEncryptionAtRestApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetEncryptionAtRestExecute executes the request
//
//	@return EncryptionAtRest
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetEncryptionAtRestExecute(r GetEncryptionAtRestApiRequest) (*EncryptionAtRest, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EncryptionAtRest
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.GetEncryptionAtRest")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest"
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

type GetRestPrivateEndpointApiRequest struct {
	ctx           context.Context
	ApiService    EncryptionAtRestUsingCustomerKeyManagementApi
	groupId       string
	cloudProvider string
	endpointId    string
}

type GetRestPrivateEndpointApiParams struct {
	GroupId       string
	CloudProvider string
	EndpointId    string
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetRestPrivateEndpointWithParams(ctx context.Context, args *GetRestPrivateEndpointApiParams) GetRestPrivateEndpointApiRequest {
	return GetRestPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		endpointId:    args.EndpointId,
	}
}

func (r GetRestPrivateEndpointApiRequest) Execute() (*EARPrivateEndpoint, *http.Response, error) {
	return r.ApiService.GetRestPrivateEndpointExecute(r)
}

/*
GetRestPrivateEndpoint Return One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

Returns one private endpoint, identified by its ID, for encryption at rest using Customer Key Management.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint.
	@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint.
	@return GetRestPrivateEndpointApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetRestPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) GetRestPrivateEndpointApiRequest {
	return GetRestPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
		endpointId:    endpointId,
	}
}

// GetRestPrivateEndpointExecute executes the request
//
//	@return EARPrivateEndpoint
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) GetRestPrivateEndpointExecute(r GetRestPrivateEndpointApiRequest) (*EARPrivateEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EARPrivateEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.GetRestPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest/{cloudProvider}/privateEndpoints/{endpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)
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

type ListRestPrivateEndpointsApiRequest struct {
	ctx           context.Context
	ApiService    EncryptionAtRestUsingCustomerKeyManagementApi
	groupId       string
	cloudProvider string
	includeCount  *bool
	itemsPerPage  *int
	pageNum       *int
}

type ListRestPrivateEndpointsApiParams struct {
	GroupId       string
	CloudProvider string
	IncludeCount  *bool
	ItemsPerPage  *int
	PageNum       *int
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) ListRestPrivateEndpointsWithParams(ctx context.Context, args *ListRestPrivateEndpointsApiParams) ListRestPrivateEndpointsApiRequest {
	return ListRestPrivateEndpointsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		includeCount:  args.IncludeCount,
		itemsPerPage:  args.ItemsPerPage,
		pageNum:       args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListRestPrivateEndpointsApiRequest) IncludeCount(includeCount bool) ListRestPrivateEndpointsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListRestPrivateEndpointsApiRequest) ItemsPerPage(itemsPerPage int) ListRestPrivateEndpointsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListRestPrivateEndpointsApiRequest) PageNum(pageNum int) ListRestPrivateEndpointsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListRestPrivateEndpointsApiRequest) Execute() (*PaginatedApiAtlasEARPrivateEndpoint, *http.Response, error) {
	return r.ApiService.ListRestPrivateEndpointsExecute(r)
}

/*
ListRestPrivateEndpoints Return Private Endpoints for Encryption at Rest Using Customer Key Management for One Cloud Provider in One Project

Returns the private endpoints of the specified cloud provider for encryption at rest using customer key management.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoints to return.
	@return ListRestPrivateEndpointsApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) ListRestPrivateEndpoints(ctx context.Context, groupId string, cloudProvider string) ListRestPrivateEndpointsApiRequest {
	return ListRestPrivateEndpointsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
	}
}

// ListRestPrivateEndpointsExecute executes the request
//
//	@return PaginatedApiAtlasEARPrivateEndpoint
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) ListRestPrivateEndpointsExecute(r ListRestPrivateEndpointsApiRequest) (*PaginatedApiAtlasEARPrivateEndpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasEARPrivateEndpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.ListRestPrivateEndpoints")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest/{cloudProvider}/privateEndpoints"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)

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

type RequestPrivateEndpointDeletionApiRequest struct {
	ctx           context.Context
	ApiService    EncryptionAtRestUsingCustomerKeyManagementApi
	groupId       string
	cloudProvider string
	endpointId    string
}

type RequestPrivateEndpointDeletionApiParams struct {
	GroupId       string
	CloudProvider string
	EndpointId    string
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) RequestPrivateEndpointDeletionWithParams(ctx context.Context, args *RequestPrivateEndpointDeletionApiParams) RequestPrivateEndpointDeletionApiRequest {
	return RequestPrivateEndpointDeletionApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		endpointId:    args.EndpointId,
	}
}

func (r RequestPrivateEndpointDeletionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RequestPrivateEndpointDeletionExecute(r)
}

/*
RequestPrivateEndpointDeletion Delete One Private Endpoint for Encryption at Rest Using Customer Key Management for One Cloud Provider from One Project

Deletes one private endpoint, identified by its ID, for encryption at rest using Customer Key Management.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint to delete.
	@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint to delete.
	@return RequestPrivateEndpointDeletionApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) RequestPrivateEndpointDeletion(ctx context.Context, groupId string, cloudProvider string, endpointId string) RequestPrivateEndpointDeletionApiRequest {
	return RequestPrivateEndpointDeletionApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
		endpointId:    endpointId,
	}
}

// RequestPrivateEndpointDeletionExecute executes the request
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) RequestPrivateEndpointDeletionExecute(r RequestPrivateEndpointDeletionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.RequestPrivateEndpointDeletion")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest/{cloudProvider}/privateEndpoints/{endpointId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)
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

type UpdateEncryptionAtRestApiRequest struct {
	ctx              context.Context
	ApiService       EncryptionAtRestUsingCustomerKeyManagementApi
	groupId          string
	encryptionAtRest *EncryptionAtRest
}

type UpdateEncryptionAtRestApiParams struct {
	GroupId          string
	EncryptionAtRest *EncryptionAtRest
}

func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) UpdateEncryptionAtRestWithParams(ctx context.Context, args *UpdateEncryptionAtRestApiParams) UpdateEncryptionAtRestApiRequest {
	return UpdateEncryptionAtRestApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		encryptionAtRest: args.EncryptionAtRest,
	}
}

func (r UpdateEncryptionAtRestApiRequest) Execute() (*EncryptionAtRest, *http.Response, error) {
	return r.ApiService.UpdateEncryptionAtRestExecute(r)
}

/*
UpdateEncryptionAtRest Update Encryption at Rest Configuration in One Project

Updates the configuration for encryption at rest using the keys you manage through your cloud provider. MongoDB Cloud encrypts all storage even if you don't use your own key management. This feature isn't available for `M0` free clusters, `M2`, `M5`, or serverless clusters.

	After you configure at least one Encryption at Rest using a Customer Key Management provider for the MongoDB Cloud project, Project Owners can enable Encryption at Rest using Customer Key Management for each MongoDB Cloud cluster for which they require encryption. The Encryption at Rest using Customer Key Management provider doesn't have to match the cluster cloud service provider. MongoDB Cloud doesn't automatically rotate user-managed encryption keys. Defer to your preferred Encryption at Rest using Customer Key Management provider's documentation and guidance for best practices on key rotation. MongoDB Cloud automatically creates a 90-day key rotation alert when you configure Encryption at Rest using Customer Key Management using your Key Management in an MongoDB Cloud project. MongoDB Cloud encrypts all storage whether or not you use your own key management.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateEncryptionAtRestApiRequest
*/
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) UpdateEncryptionAtRest(ctx context.Context, groupId string, encryptionAtRest *EncryptionAtRest) UpdateEncryptionAtRestApiRequest {
	return UpdateEncryptionAtRestApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		encryptionAtRest: encryptionAtRest,
	}
}

// UpdateEncryptionAtRestExecute executes the request
//
//	@return EncryptionAtRest
func (a *EncryptionAtRestUsingCustomerKeyManagementApiService) UpdateEncryptionAtRestExecute(r UpdateEncryptionAtRestApiRequest) (*EncryptionAtRest, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *EncryptionAtRest
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "EncryptionAtRestUsingCustomerKeyManagementApiService.UpdateEncryptionAtRest")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/encryptionAtRest"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.encryptionAtRest == nil {
		return localVarReturnValue, nil, reportError("encryptionAtRest is required and must be specified")
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
	localVarPostBody = r.encryptionAtRest
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
