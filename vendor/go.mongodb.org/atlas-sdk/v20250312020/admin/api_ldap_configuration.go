// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type LDAPConfigurationApi interface {

	/*
		DeleteLdapUserMapping Remove LDAP User to DN Mapping

		Removes the current LDAP Distinguished Name mapping captured in the ``userToDNMapping`` document from the LDAP configuration for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteLdapUserMappingApiRequest
	*/
	DeleteLdapUserMapping(ctx context.Context, groupId string) DeleteLdapUserMappingApiRequest
	/*
		DeleteLdapUserMapping Remove LDAP User to DN Mapping


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteLdapUserMappingApiParams - Parameters for the request
		@return DeleteLdapUserMappingApiRequest
	*/
	DeleteLdapUserMappingWithParams(ctx context.Context, args *DeleteLdapUserMappingApiParams) DeleteLdapUserMappingApiRequest

	// Method available only for mocking purposes
	DeleteLdapUserMappingExecute(r DeleteLdapUserMappingApiRequest) (*UserSecurity, *http.Response, error)

	/*
		GetUserSecurity Return LDAP or X.509 Configuration

		Returns the current LDAP configuration for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetUserSecurityApiRequest
	*/
	GetUserSecurity(ctx context.Context, groupId string) GetUserSecurityApiRequest
	/*
		GetUserSecurity Return LDAP or X.509 Configuration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetUserSecurityApiParams - Parameters for the request
		@return GetUserSecurityApiRequest
	*/
	GetUserSecurityWithParams(ctx context.Context, args *GetUserSecurityApiParams) GetUserSecurityApiRequest

	// Method available only for mocking purposes
	GetUserSecurityExecute(r GetUserSecurityApiRequest) (*UserSecurity, *http.Response, error)

	/*
		GetUserSecurityVerify Return Status of LDAP Configuration Verification in One Project

		Returns the status of one request to verify one LDAP configuration for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param requestId Unique string that identifies the request to verify an Lightweight Directory Access Protocol (LDAP) configuration.
		@return GetUserSecurityVerifyApiRequest
	*/
	GetUserSecurityVerify(ctx context.Context, groupId string, requestId string) GetUserSecurityVerifyApiRequest
	/*
		GetUserSecurityVerify Return Status of LDAP Configuration Verification in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetUserSecurityVerifyApiParams - Parameters for the request
		@return GetUserSecurityVerifyApiRequest
	*/
	GetUserSecurityVerifyWithParams(ctx context.Context, args *GetUserSecurityVerifyApiParams) GetUserSecurityVerifyApiRequest

	// Method available only for mocking purposes
	GetUserSecurityVerifyExecute(r GetUserSecurityVerifyApiRequest) (*LDAPVerifyConnectivityJobRequest, *http.Response, error)

	/*
			UpdateUserSecurity Update LDAP or X.509 Configuration

			Edits the LDAP configuration for the specified project.

		Updating this configuration triggers a rolling restart of the database.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param userSecurity Updates the LDAP configuration for the specified project.
			@return UpdateUserSecurityApiRequest
	*/
	UpdateUserSecurity(ctx context.Context, groupId string, userSecurity *UserSecurity) UpdateUserSecurityApiRequest
	/*
		UpdateUserSecurity Update LDAP or X.509 Configuration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateUserSecurityApiParams - Parameters for the request
		@return UpdateUserSecurityApiRequest
	*/
	UpdateUserSecurityWithParams(ctx context.Context, args *UpdateUserSecurityApiParams) UpdateUserSecurityApiRequest

	// Method available only for mocking purposes
	UpdateUserSecurityExecute(r UpdateUserSecurityApiRequest) (*UserSecurity, *http.Response, error)

	/*
		VerifyUserSecurityLdap Verify LDAP Configuration in One Project

		Verifies the LDAP configuration for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param lDAPVerifyConnectivityJobRequestParams The LDAP configuration for the specified project that you want to verify.
		@return VerifyUserSecurityLdapApiRequest
	*/
	VerifyUserSecurityLdap(ctx context.Context, groupId string, lDAPVerifyConnectivityJobRequestParams *LDAPVerifyConnectivityJobRequestParams) VerifyUserSecurityLdapApiRequest
	/*
		VerifyUserSecurityLdap Verify LDAP Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param VerifyUserSecurityLdapApiParams - Parameters for the request
		@return VerifyUserSecurityLdapApiRequest
	*/
	VerifyUserSecurityLdapWithParams(ctx context.Context, args *VerifyUserSecurityLdapApiParams) VerifyUserSecurityLdapApiRequest

	// Method available only for mocking purposes
	VerifyUserSecurityLdapExecute(r VerifyUserSecurityLdapApiRequest) (*LDAPVerifyConnectivityJobRequest, *http.Response, error)
}

// LDAPConfigurationApiService LDAPConfigurationApi service
type LDAPConfigurationApiService service

type DeleteLdapUserMappingApiRequest struct {
	ctx        context.Context
	ApiService LDAPConfigurationApi
	groupId    string
}

type DeleteLdapUserMappingApiParams struct {
	GroupId string
}

func (a *LDAPConfigurationApiService) DeleteLdapUserMappingWithParams(ctx context.Context, args *DeleteLdapUserMappingApiParams) DeleteLdapUserMappingApiRequest {
	return DeleteLdapUserMappingApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DeleteLdapUserMappingApiRequest) Execute() (*UserSecurity, *http.Response, error) {
	return r.ApiService.DeleteLdapUserMappingExecute(r)
}

/*
DeleteLdapUserMapping Remove LDAP User to DN Mapping

Removes the current LDAP Distinguished Name mapping captured in the “userToDNMapping“ document from the LDAP configuration for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteLdapUserMappingApiRequest
*/
func (a *LDAPConfigurationApiService) DeleteLdapUserMapping(ctx context.Context, groupId string) DeleteLdapUserMappingApiRequest {
	return DeleteLdapUserMappingApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DeleteLdapUserMappingExecute executes the request
//
//	@return UserSecurity
func (a *LDAPConfigurationApiService) DeleteLdapUserMappingExecute(r DeleteLdapUserMappingApiRequest) (*UserSecurity, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserSecurity
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LDAPConfigurationApiService.DeleteLdapUserMapping")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity/ldap/userToDNMapping"
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

type GetUserSecurityApiRequest struct {
	ctx        context.Context
	ApiService LDAPConfigurationApi
	groupId    string
}

type GetUserSecurityApiParams struct {
	GroupId string
}

func (a *LDAPConfigurationApiService) GetUserSecurityWithParams(ctx context.Context, args *GetUserSecurityApiParams) GetUserSecurityApiRequest {
	return GetUserSecurityApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetUserSecurityApiRequest) Execute() (*UserSecurity, *http.Response, error) {
	return r.ApiService.GetUserSecurityExecute(r)
}

/*
GetUserSecurity Return LDAP or X.509 Configuration

Returns the current LDAP configuration for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetUserSecurityApiRequest
*/
func (a *LDAPConfigurationApiService) GetUserSecurity(ctx context.Context, groupId string) GetUserSecurityApiRequest {
	return GetUserSecurityApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetUserSecurityExecute executes the request
//
//	@return UserSecurity
func (a *LDAPConfigurationApiService) GetUserSecurityExecute(r GetUserSecurityApiRequest) (*UserSecurity, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserSecurity
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LDAPConfigurationApiService.GetUserSecurity")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity"
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

type GetUserSecurityVerifyApiRequest struct {
	ctx        context.Context
	ApiService LDAPConfigurationApi
	groupId    string
	requestId  string
}

type GetUserSecurityVerifyApiParams struct {
	GroupId   string
	RequestId string
}

func (a *LDAPConfigurationApiService) GetUserSecurityVerifyWithParams(ctx context.Context, args *GetUserSecurityVerifyApiParams) GetUserSecurityVerifyApiRequest {
	return GetUserSecurityVerifyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		requestId:  args.RequestId,
	}
}

func (r GetUserSecurityVerifyApiRequest) Execute() (*LDAPVerifyConnectivityJobRequest, *http.Response, error) {
	return r.ApiService.GetUserSecurityVerifyExecute(r)
}

/*
GetUserSecurityVerify Return Status of LDAP Configuration Verification in One Project

Returns the status of one request to verify one LDAP configuration for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param requestId Unique string that identifies the request to verify an Lightweight Directory Access Protocol (LDAP) configuration.
	@return GetUserSecurityVerifyApiRequest
*/
func (a *LDAPConfigurationApiService) GetUserSecurityVerify(ctx context.Context, groupId string, requestId string) GetUserSecurityVerifyApiRequest {
	return GetUserSecurityVerifyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		requestId:  requestId,
	}
}

// GetUserSecurityVerifyExecute executes the request
//
//	@return LDAPVerifyConnectivityJobRequest
func (a *LDAPConfigurationApiService) GetUserSecurityVerifyExecute(r GetUserSecurityVerifyApiRequest) (*LDAPVerifyConnectivityJobRequest, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LDAPVerifyConnectivityJobRequest
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LDAPConfigurationApiService.GetUserSecurityVerify")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity/ldap/verify/{requestId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.requestId == "" {
		return localVarReturnValue, nil, reportError("requestId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"requestId"+"}", url.PathEscape(r.requestId), -1)

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

type UpdateUserSecurityApiRequest struct {
	ctx          context.Context
	ApiService   LDAPConfigurationApi
	groupId      string
	userSecurity *UserSecurity
}

type UpdateUserSecurityApiParams struct {
	GroupId      string
	UserSecurity *UserSecurity
}

func (a *LDAPConfigurationApiService) UpdateUserSecurityWithParams(ctx context.Context, args *UpdateUserSecurityApiParams) UpdateUserSecurityApiRequest {
	return UpdateUserSecurityApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		userSecurity: args.UserSecurity,
	}
}

func (r UpdateUserSecurityApiRequest) Execute() (*UserSecurity, *http.Response, error) {
	return r.ApiService.UpdateUserSecurityExecute(r)
}

/*
UpdateUserSecurity Update LDAP or X.509 Configuration

Edits the LDAP configuration for the specified project.

Updating this configuration triggers a rolling restart of the database.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateUserSecurityApiRequest
*/
func (a *LDAPConfigurationApiService) UpdateUserSecurity(ctx context.Context, groupId string, userSecurity *UserSecurity) UpdateUserSecurityApiRequest {
	return UpdateUserSecurityApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		userSecurity: userSecurity,
	}
}

// UpdateUserSecurityExecute executes the request
//
//	@return UserSecurity
func (a *LDAPConfigurationApiService) UpdateUserSecurityExecute(r UpdateUserSecurityApiRequest) (*UserSecurity, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserSecurity
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LDAPConfigurationApiService.UpdateUserSecurity")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userSecurity == nil {
		return localVarReturnValue, nil, reportError("userSecurity is required and must be specified")
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
	localVarPostBody = r.userSecurity
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

type VerifyUserSecurityLdapApiRequest struct {
	ctx                                    context.Context
	ApiService                             LDAPConfigurationApi
	groupId                                string
	lDAPVerifyConnectivityJobRequestParams *LDAPVerifyConnectivityJobRequestParams
}

type VerifyUserSecurityLdapApiParams struct {
	GroupId                                string
	LDAPVerifyConnectivityJobRequestParams *LDAPVerifyConnectivityJobRequestParams
}

func (a *LDAPConfigurationApiService) VerifyUserSecurityLdapWithParams(ctx context.Context, args *VerifyUserSecurityLdapApiParams) VerifyUserSecurityLdapApiRequest {
	return VerifyUserSecurityLdapApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                args.GroupId,
		lDAPVerifyConnectivityJobRequestParams: args.LDAPVerifyConnectivityJobRequestParams,
	}
}

func (r VerifyUserSecurityLdapApiRequest) Execute() (*LDAPVerifyConnectivityJobRequest, *http.Response, error) {
	return r.ApiService.VerifyUserSecurityLdapExecute(r)
}

/*
VerifyUserSecurityLdap Verify LDAP Configuration in One Project

Verifies the LDAP configuration for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return VerifyUserSecurityLdapApiRequest
*/
func (a *LDAPConfigurationApiService) VerifyUserSecurityLdap(ctx context.Context, groupId string, lDAPVerifyConnectivityJobRequestParams *LDAPVerifyConnectivityJobRequestParams) VerifyUserSecurityLdapApiRequest {
	return VerifyUserSecurityLdapApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                groupId,
		lDAPVerifyConnectivityJobRequestParams: lDAPVerifyConnectivityJobRequestParams,
	}
}

// VerifyUserSecurityLdapExecute executes the request
//
//	@return LDAPVerifyConnectivityJobRequest
func (a *LDAPConfigurationApiService) VerifyUserSecurityLdapExecute(r VerifyUserSecurityLdapApiRequest) (*LDAPVerifyConnectivityJobRequest, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LDAPVerifyConnectivityJobRequest
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LDAPConfigurationApiService.VerifyUserSecurityLdap")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity/ldap/verify"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.lDAPVerifyConnectivityJobRequestParams == nil {
		return localVarReturnValue, nil, reportError("lDAPVerifyConnectivityJobRequestParams is required and must be specified")
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
	localVarPostBody = r.lDAPVerifyConnectivityJobRequestParams
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
