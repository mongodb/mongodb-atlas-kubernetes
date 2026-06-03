// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CloudProviderAccessApi interface {

	/*
		AuthorizeProviderAccessRole Authorize One Cloud Provider Access Role

		Grants access to the specified project for the specified access role. This API endpoint is one step in a procedure to create unified access for MongoDB Cloud services. This is not required for GCP service account access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param roleId Unique 24-hexadecimal digit string that identifies the role.
		@param cloudProviderAccessRoleRequestUpdate Grants access to the specified project for the specified access role.
		@return AuthorizeProviderAccessRoleApiRequest
	*/
	AuthorizeProviderAccessRole(ctx context.Context, groupId string, roleId string, cloudProviderAccessRoleRequestUpdate *CloudProviderAccessRoleRequestUpdate) AuthorizeProviderAccessRoleApiRequest
	/*
		AuthorizeProviderAccessRole Authorize One Cloud Provider Access Role


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AuthorizeProviderAccessRoleApiParams - Parameters for the request
		@return AuthorizeProviderAccessRoleApiRequest
	*/
	AuthorizeProviderAccessRoleWithParams(ctx context.Context, args *AuthorizeProviderAccessRoleApiParams) AuthorizeProviderAccessRoleApiRequest

	// Method available only for mocking purposes
	AuthorizeProviderAccessRoleExecute(r AuthorizeProviderAccessRoleApiRequest) (*CloudProviderAccessRole, *http.Response, error)

	/*
		CreateCloudProviderAccess Create One Cloud Provider Access Role

		Creates one access role for the specified cloud provider. Some MongoDB Cloud features use these cloud provider access roles for authentication. For the GCP provider, if the project folder is not yet provisioned, Atlas will now create the role asynchronously. An intermediate role with status `IN_PROGRESS` will be returned, and the final service account will be provisioned. Once the GCP project is set up, subsequent requests will create the service account synchronously.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProviderAccessRoleRequest Creates one role for the specified cloud provider.
		@return CreateCloudProviderAccessApiRequest
	*/
	CreateCloudProviderAccess(ctx context.Context, groupId string, cloudProviderAccessRoleRequest *CloudProviderAccessRoleRequest) CreateCloudProviderAccessApiRequest
	/*
		CreateCloudProviderAccess Create One Cloud Provider Access Role


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateCloudProviderAccessApiParams - Parameters for the request
		@return CreateCloudProviderAccessApiRequest
	*/
	CreateCloudProviderAccessWithParams(ctx context.Context, args *CreateCloudProviderAccessApiParams) CreateCloudProviderAccessApiRequest

	// Method available only for mocking purposes
	CreateCloudProviderAccessExecute(r CreateCloudProviderAccessApiRequest) (*CloudProviderAccessRole, *http.Response, error)

	/*
		DeauthorizeProviderAccessRole Deauthorize One Cloud Provider Access Role

		Revokes access to the specified project for the specified access role.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider of the role to deauthorize.
		@param roleId Unique 24-hexadecimal digit string that identifies the role.
		@return DeauthorizeProviderAccessRoleApiRequest
	*/
	DeauthorizeProviderAccessRole(ctx context.Context, groupId string, cloudProvider string, roleId string) DeauthorizeProviderAccessRoleApiRequest
	/*
		DeauthorizeProviderAccessRole Deauthorize One Cloud Provider Access Role


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeauthorizeProviderAccessRoleApiParams - Parameters for the request
		@return DeauthorizeProviderAccessRoleApiRequest
	*/
	DeauthorizeProviderAccessRoleWithParams(ctx context.Context, args *DeauthorizeProviderAccessRoleApiParams) DeauthorizeProviderAccessRoleApiRequest

	// Method available only for mocking purposes
	DeauthorizeProviderAccessRoleExecute(r DeauthorizeProviderAccessRoleApiRequest) (*http.Response, error)

	/*
		GetCloudProviderAccess Return One Cloud Provider Access Role

		Returns the access role with the specified id and with access to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param roleId Unique 24-hexadecimal digit string that identifies the role.
		@return GetCloudProviderAccessApiRequest
	*/
	GetCloudProviderAccess(ctx context.Context, groupId string, roleId string) GetCloudProviderAccessApiRequest
	/*
		GetCloudProviderAccess Return One Cloud Provider Access Role


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetCloudProviderAccessApiParams - Parameters for the request
		@return GetCloudProviderAccessApiRequest
	*/
	GetCloudProviderAccessWithParams(ctx context.Context, args *GetCloudProviderAccessApiParams) GetCloudProviderAccessApiRequest

	// Method available only for mocking purposes
	GetCloudProviderAccessExecute(r GetCloudProviderAccessApiRequest) (*CloudProviderAccessRole, *http.Response, error)

	/*
		ListCloudProviderAccess Return All Cloud Provider Access Roles

		Returns all cloud provider access roles with access to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListCloudProviderAccessApiRequest
	*/
	ListCloudProviderAccess(ctx context.Context, groupId string) ListCloudProviderAccessApiRequest
	/*
		ListCloudProviderAccess Return All Cloud Provider Access Roles


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListCloudProviderAccessApiParams - Parameters for the request
		@return ListCloudProviderAccessApiRequest
	*/
	ListCloudProviderAccessWithParams(ctx context.Context, args *ListCloudProviderAccessApiParams) ListCloudProviderAccessApiRequest

	// Method available only for mocking purposes
	ListCloudProviderAccessExecute(r ListCloudProviderAccessApiRequest) (*CloudProviderAccessRoles, *http.Response, error)
}

// CloudProviderAccessApiService CloudProviderAccessApi service
type CloudProviderAccessApiService service

type AuthorizeProviderAccessRoleApiRequest struct {
	ctx                                  context.Context
	ApiService                           CloudProviderAccessApi
	groupId                              string
	roleId                               string
	cloudProviderAccessRoleRequestUpdate *CloudProviderAccessRoleRequestUpdate
}

type AuthorizeProviderAccessRoleApiParams struct {
	GroupId                              string
	RoleId                               string
	CloudProviderAccessRoleRequestUpdate *CloudProviderAccessRoleRequestUpdate
}

func (a *CloudProviderAccessApiService) AuthorizeProviderAccessRoleWithParams(ctx context.Context, args *AuthorizeProviderAccessRoleApiParams) AuthorizeProviderAccessRoleApiRequest {
	return AuthorizeProviderAccessRoleApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              args.GroupId,
		roleId:                               args.RoleId,
		cloudProviderAccessRoleRequestUpdate: args.CloudProviderAccessRoleRequestUpdate,
	}
}

func (r AuthorizeProviderAccessRoleApiRequest) Execute() (*CloudProviderAccessRole, *http.Response, error) {
	return r.ApiService.AuthorizeProviderAccessRoleExecute(r)
}

/*
AuthorizeProviderAccessRole Authorize One Cloud Provider Access Role

Grants access to the specified project for the specified access role. This API endpoint is one step in a procedure to create unified access for MongoDB Cloud services. This is not required for GCP service account access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param roleId Unique 24-hexadecimal digit string that identifies the role.
	@return AuthorizeProviderAccessRoleApiRequest
*/
func (a *CloudProviderAccessApiService) AuthorizeProviderAccessRole(ctx context.Context, groupId string, roleId string, cloudProviderAccessRoleRequestUpdate *CloudProviderAccessRoleRequestUpdate) AuthorizeProviderAccessRoleApiRequest {
	return AuthorizeProviderAccessRoleApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              groupId,
		roleId:                               roleId,
		cloudProviderAccessRoleRequestUpdate: cloudProviderAccessRoleRequestUpdate,
	}
}

// AuthorizeProviderAccessRoleExecute executes the request
//
//	@return CloudProviderAccessRole
func (a *CloudProviderAccessApiService) AuthorizeProviderAccessRoleExecute(r AuthorizeProviderAccessRoleApiRequest) (*CloudProviderAccessRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderAccessRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudProviderAccessApiService.AuthorizeProviderAccessRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/cloudProviderAccess/{roleId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.roleId == "" {
		return localVarReturnValue, nil, reportError("roleId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleId"+"}", url.PathEscape(r.roleId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudProviderAccessRoleRequestUpdate == nil {
		return localVarReturnValue, nil, reportError("cloudProviderAccessRoleRequestUpdate is required and must be specified")
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
	localVarPostBody = r.cloudProviderAccessRoleRequestUpdate
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

type CreateCloudProviderAccessApiRequest struct {
	ctx                            context.Context
	ApiService                     CloudProviderAccessApi
	groupId                        string
	cloudProviderAccessRoleRequest *CloudProviderAccessRoleRequest
}

type CreateCloudProviderAccessApiParams struct {
	GroupId                        string
	CloudProviderAccessRoleRequest *CloudProviderAccessRoleRequest
}

func (a *CloudProviderAccessApiService) CreateCloudProviderAccessWithParams(ctx context.Context, args *CreateCloudProviderAccessApiParams) CreateCloudProviderAccessApiRequest {
	return CreateCloudProviderAccessApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        args.GroupId,
		cloudProviderAccessRoleRequest: args.CloudProviderAccessRoleRequest,
	}
}

func (r CreateCloudProviderAccessApiRequest) Execute() (*CloudProviderAccessRole, *http.Response, error) {
	return r.ApiService.CreateCloudProviderAccessExecute(r)
}

/*
CreateCloudProviderAccess Create One Cloud Provider Access Role

Creates one access role for the specified cloud provider. Some MongoDB Cloud features use these cloud provider access roles for authentication. For the GCP provider, if the project folder is not yet provisioned, Atlas will now create the role asynchronously. An intermediate role with status `IN_PROGRESS` will be returned, and the final service account will be provisioned. Once the GCP project is set up, subsequent requests will create the service account synchronously.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateCloudProviderAccessApiRequest
*/
func (a *CloudProviderAccessApiService) CreateCloudProviderAccess(ctx context.Context, groupId string, cloudProviderAccessRoleRequest *CloudProviderAccessRoleRequest) CreateCloudProviderAccessApiRequest {
	return CreateCloudProviderAccessApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        groupId,
		cloudProviderAccessRoleRequest: cloudProviderAccessRoleRequest,
	}
}

// CreateCloudProviderAccessExecute executes the request
//
//	@return CloudProviderAccessRole
func (a *CloudProviderAccessApiService) CreateCloudProviderAccessExecute(r CreateCloudProviderAccessApiRequest) (*CloudProviderAccessRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderAccessRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudProviderAccessApiService.CreateCloudProviderAccess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/cloudProviderAccess"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudProviderAccessRoleRequest == nil {
		return localVarReturnValue, nil, reportError("cloudProviderAccessRoleRequest is required and must be specified")
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
	localVarPostBody = r.cloudProviderAccessRoleRequest
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

type DeauthorizeProviderAccessRoleApiRequest struct {
	ctx           context.Context
	ApiService    CloudProviderAccessApi
	groupId       string
	cloudProvider string
	roleId        string
}

type DeauthorizeProviderAccessRoleApiParams struct {
	GroupId       string
	CloudProvider string
	RoleId        string
}

func (a *CloudProviderAccessApiService) DeauthorizeProviderAccessRoleWithParams(ctx context.Context, args *DeauthorizeProviderAccessRoleApiParams) DeauthorizeProviderAccessRoleApiRequest {
	return DeauthorizeProviderAccessRoleApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		roleId:        args.RoleId,
	}
}

func (r DeauthorizeProviderAccessRoleApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeauthorizeProviderAccessRoleExecute(r)
}

/*
DeauthorizeProviderAccessRole Deauthorize One Cloud Provider Access Role

Revokes access to the specified project for the specified access role.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider of the role to deauthorize.
	@param roleId Unique 24-hexadecimal digit string that identifies the role.
	@return DeauthorizeProviderAccessRoleApiRequest
*/
func (a *CloudProviderAccessApiService) DeauthorizeProviderAccessRole(ctx context.Context, groupId string, cloudProvider string, roleId string) DeauthorizeProviderAccessRoleApiRequest {
	return DeauthorizeProviderAccessRoleApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
		roleId:        roleId,
	}
}

// DeauthorizeProviderAccessRoleExecute executes the request
func (a *CloudProviderAccessApiService) DeauthorizeProviderAccessRoleExecute(r DeauthorizeProviderAccessRoleApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudProviderAccessApiService.DeauthorizeProviderAccessRole")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/cloudProviderAccess/{cloudProvider}/{roleId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)
	if r.roleId == "" {
		return nil, reportError("roleId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleId"+"}", url.PathEscape(r.roleId), -1)

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

type GetCloudProviderAccessApiRequest struct {
	ctx        context.Context
	ApiService CloudProviderAccessApi
	groupId    string
	roleId     string
}

type GetCloudProviderAccessApiParams struct {
	GroupId string
	RoleId  string
}

func (a *CloudProviderAccessApiService) GetCloudProviderAccessWithParams(ctx context.Context, args *GetCloudProviderAccessApiParams) GetCloudProviderAccessApiRequest {
	return GetCloudProviderAccessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		roleId:     args.RoleId,
	}
}

func (r GetCloudProviderAccessApiRequest) Execute() (*CloudProviderAccessRole, *http.Response, error) {
	return r.ApiService.GetCloudProviderAccessExecute(r)
}

/*
GetCloudProviderAccess Return One Cloud Provider Access Role

Returns the access role with the specified id and with access to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param roleId Unique 24-hexadecimal digit string that identifies the role.
	@return GetCloudProviderAccessApiRequest
*/
func (a *CloudProviderAccessApiService) GetCloudProviderAccess(ctx context.Context, groupId string, roleId string) GetCloudProviderAccessApiRequest {
	return GetCloudProviderAccessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		roleId:     roleId,
	}
}

// GetCloudProviderAccessExecute executes the request
//
//	@return CloudProviderAccessRole
func (a *CloudProviderAccessApiService) GetCloudProviderAccessExecute(r GetCloudProviderAccessApiRequest) (*CloudProviderAccessRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderAccessRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudProviderAccessApiService.GetCloudProviderAccess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/cloudProviderAccess/{roleId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.roleId == "" {
		return localVarReturnValue, nil, reportError("roleId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleId"+"}", url.PathEscape(r.roleId), -1)

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

type ListCloudProviderAccessApiRequest struct {
	ctx        context.Context
	ApiService CloudProviderAccessApi
	groupId    string
}

type ListCloudProviderAccessApiParams struct {
	GroupId string
}

func (a *CloudProviderAccessApiService) ListCloudProviderAccessWithParams(ctx context.Context, args *ListCloudProviderAccessApiParams) ListCloudProviderAccessApiRequest {
	return ListCloudProviderAccessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ListCloudProviderAccessApiRequest) Execute() (*CloudProviderAccessRoles, *http.Response, error) {
	return r.ApiService.ListCloudProviderAccessExecute(r)
}

/*
ListCloudProviderAccess Return All Cloud Provider Access Roles

Returns all cloud provider access roles with access to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListCloudProviderAccessApiRequest
*/
func (a *CloudProviderAccessApiService) ListCloudProviderAccess(ctx context.Context, groupId string) ListCloudProviderAccessApiRequest {
	return ListCloudProviderAccessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListCloudProviderAccessExecute executes the request
//
//	@return CloudProviderAccessRoles
func (a *CloudProviderAccessApiService) ListCloudProviderAccessExecute(r ListCloudProviderAccessApiRequest) (*CloudProviderAccessRoles, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderAccessRoles
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudProviderAccessApiService.ListCloudProviderAccess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/cloudProviderAccess"
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
