// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CustomDatabaseRolesApi interface {

	/*
		CreateCustomDbRole Create One Custom Role

		Creates one custom role in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param userCustomDBRole Creates one custom role in the specified project.
		@return CreateCustomDbRoleApiRequest
	*/
	CreateCustomDbRole(ctx context.Context, groupId string, userCustomDBRole *UserCustomDBRole) CreateCustomDbRoleApiRequest
	/*
		CreateCustomDbRole Create One Custom Role


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateCustomDbRoleApiParams - Parameters for the request
		@return CreateCustomDbRoleApiRequest
	*/
	CreateCustomDbRoleWithParams(ctx context.Context, args *CreateCustomDbRoleApiParams) CreateCustomDbRoleApiRequest

	// Method available only for mocking purposes
	CreateCustomDbRoleExecute(r CreateCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error)

	/*
		DeleteCustomDbRole Remove One Custom Role from One Project

		Removes one custom role from the specified project. You can't remove a custom role that would leave one or more child roles with no parent roles or actions. You also can't remove a custom role that would leave one or more database users without roles.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
		@return DeleteCustomDbRoleApiRequest
	*/
	DeleteCustomDbRole(ctx context.Context, groupId string, roleName string) DeleteCustomDbRoleApiRequest
	/*
		DeleteCustomDbRole Remove One Custom Role from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteCustomDbRoleApiParams - Parameters for the request
		@return DeleteCustomDbRoleApiRequest
	*/
	DeleteCustomDbRoleWithParams(ctx context.Context, args *DeleteCustomDbRoleApiParams) DeleteCustomDbRoleApiRequest

	// Method available only for mocking purposes
	DeleteCustomDbRoleExecute(r DeleteCustomDbRoleApiRequest) (*http.Response, error)

	/*
		GetCustomDbRole Return One Custom Role in One Project

		Returns one custom role for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
		@return GetCustomDbRoleApiRequest
	*/
	GetCustomDbRole(ctx context.Context, groupId string, roleName string) GetCustomDbRoleApiRequest
	/*
		GetCustomDbRole Return One Custom Role in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetCustomDbRoleApiParams - Parameters for the request
		@return GetCustomDbRoleApiRequest
	*/
	GetCustomDbRoleWithParams(ctx context.Context, args *GetCustomDbRoleApiParams) GetCustomDbRoleApiRequest

	// Method available only for mocking purposes
	GetCustomDbRoleExecute(r GetCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error)

	/*
		ListCustomDbRoles Return All Custom Roles in One Project

		Returns all custom roles for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListCustomDbRolesApiRequest
	*/
	ListCustomDbRoles(ctx context.Context, groupId string) ListCustomDbRolesApiRequest
	/*
		ListCustomDbRoles Return All Custom Roles in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListCustomDbRolesApiParams - Parameters for the request
		@return ListCustomDbRolesApiRequest
	*/
	ListCustomDbRolesWithParams(ctx context.Context, args *ListCustomDbRolesApiParams) ListCustomDbRolesApiRequest

	// Method available only for mocking purposes
	ListCustomDbRolesExecute(r ListCustomDbRolesApiRequest) ([]UserCustomDBRole, *http.Response, error)

	/*
		UpdateCustomDbRole Update One Custom Role in One Project

		Updates one custom role in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
		@param updateCustomDBRole Updates one custom role in the specified project.
		@return UpdateCustomDbRoleApiRequest
	*/
	UpdateCustomDbRole(ctx context.Context, groupId string, roleName string, updateCustomDBRole *UpdateCustomDBRole) UpdateCustomDbRoleApiRequest
	/*
		UpdateCustomDbRole Update One Custom Role in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateCustomDbRoleApiParams - Parameters for the request
		@return UpdateCustomDbRoleApiRequest
	*/
	UpdateCustomDbRoleWithParams(ctx context.Context, args *UpdateCustomDbRoleApiParams) UpdateCustomDbRoleApiRequest

	// Method available only for mocking purposes
	UpdateCustomDbRoleExecute(r UpdateCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error)
}

// CustomDatabaseRolesApiService CustomDatabaseRolesApi service
type CustomDatabaseRolesApiService service

type CreateCustomDbRoleApiRequest struct {
	ctx              context.Context
	ApiService       CustomDatabaseRolesApi
	groupId          string
	userCustomDBRole *UserCustomDBRole
}

type CreateCustomDbRoleApiParams struct {
	GroupId          string
	UserCustomDBRole *UserCustomDBRole
}

func (a *CustomDatabaseRolesApiService) CreateCustomDbRoleWithParams(ctx context.Context, args *CreateCustomDbRoleApiParams) CreateCustomDbRoleApiRequest {
	return CreateCustomDbRoleApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		userCustomDBRole: args.UserCustomDBRole,
	}
}

func (r CreateCustomDbRoleApiRequest) Execute() (*UserCustomDBRole, *http.Response, error) {
	return r.ApiService.CreateCustomDbRoleExecute(r)
}

/*
CreateCustomDbRole Create One Custom Role

Creates one custom role in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateCustomDbRoleApiRequest
*/
func (a *CustomDatabaseRolesApiService) CreateCustomDbRole(ctx context.Context, groupId string, userCustomDBRole *UserCustomDBRole) CreateCustomDbRoleApiRequest {
	return CreateCustomDbRoleApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		userCustomDBRole: userCustomDBRole,
	}
}

// CreateCustomDbRoleExecute executes the request
//
//	@return UserCustomDBRole
func (a *CustomDatabaseRolesApiService) CreateCustomDbRoleExecute(r CreateCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserCustomDBRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CustomDatabaseRolesApiService.CreateCustomDbRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/customDBRoles/roles"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userCustomDBRole == nil {
		return localVarReturnValue, nil, reportError("userCustomDBRole is required and must be specified")
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
	localVarPostBody = r.userCustomDBRole
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

type DeleteCustomDbRoleApiRequest struct {
	ctx        context.Context
	ApiService CustomDatabaseRolesApi
	groupId    string
	roleName   string
}

type DeleteCustomDbRoleApiParams struct {
	GroupId  string
	RoleName string
}

func (a *CustomDatabaseRolesApiService) DeleteCustomDbRoleWithParams(ctx context.Context, args *DeleteCustomDbRoleApiParams) DeleteCustomDbRoleApiRequest {
	return DeleteCustomDbRoleApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		roleName:   args.RoleName,
	}
}

func (r DeleteCustomDbRoleApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteCustomDbRoleExecute(r)
}

/*
DeleteCustomDbRole Remove One Custom Role from One Project

Removes one custom role from the specified project. You can't remove a custom role that would leave one or more child roles with no parent roles or actions. You also can't remove a custom role that would leave one or more database users without roles.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
	@return DeleteCustomDbRoleApiRequest
*/
func (a *CustomDatabaseRolesApiService) DeleteCustomDbRole(ctx context.Context, groupId string, roleName string) DeleteCustomDbRoleApiRequest {
	return DeleteCustomDbRoleApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		roleName:   roleName,
	}
}

// DeleteCustomDbRoleExecute executes the request
func (a *CustomDatabaseRolesApiService) DeleteCustomDbRoleExecute(r DeleteCustomDbRoleApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CustomDatabaseRolesApiService.DeleteCustomDbRole")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/customDBRoles/roles/{roleName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.roleName == "" {
		return nil, reportError("roleName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleName"+"}", url.PathEscape(r.roleName), -1)

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

type GetCustomDbRoleApiRequest struct {
	ctx        context.Context
	ApiService CustomDatabaseRolesApi
	groupId    string
	roleName   string
}

type GetCustomDbRoleApiParams struct {
	GroupId  string
	RoleName string
}

func (a *CustomDatabaseRolesApiService) GetCustomDbRoleWithParams(ctx context.Context, args *GetCustomDbRoleApiParams) GetCustomDbRoleApiRequest {
	return GetCustomDbRoleApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		roleName:   args.RoleName,
	}
}

func (r GetCustomDbRoleApiRequest) Execute() (*UserCustomDBRole, *http.Response, error) {
	return r.ApiService.GetCustomDbRoleExecute(r)
}

/*
GetCustomDbRole Return One Custom Role in One Project

Returns one custom role for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
	@return GetCustomDbRoleApiRequest
*/
func (a *CustomDatabaseRolesApiService) GetCustomDbRole(ctx context.Context, groupId string, roleName string) GetCustomDbRoleApiRequest {
	return GetCustomDbRoleApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		roleName:   roleName,
	}
}

// GetCustomDbRoleExecute executes the request
//
//	@return UserCustomDBRole
func (a *CustomDatabaseRolesApiService) GetCustomDbRoleExecute(r GetCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserCustomDBRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CustomDatabaseRolesApiService.GetCustomDbRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/customDBRoles/roles/{roleName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.roleName == "" {
		return localVarReturnValue, nil, reportError("roleName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleName"+"}", url.PathEscape(r.roleName), -1)

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

type ListCustomDbRolesApiRequest struct {
	ctx        context.Context
	ApiService CustomDatabaseRolesApi
	groupId    string
}

type ListCustomDbRolesApiParams struct {
	GroupId string
}

func (a *CustomDatabaseRolesApiService) ListCustomDbRolesWithParams(ctx context.Context, args *ListCustomDbRolesApiParams) ListCustomDbRolesApiRequest {
	return ListCustomDbRolesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ListCustomDbRolesApiRequest) Execute() ([]UserCustomDBRole, *http.Response, error) {
	return r.ApiService.ListCustomDbRolesExecute(r)
}

/*
ListCustomDbRoles Return All Custom Roles in One Project

Returns all custom roles for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListCustomDbRolesApiRequest
*/
func (a *CustomDatabaseRolesApiService) ListCustomDbRoles(ctx context.Context, groupId string) ListCustomDbRolesApiRequest {
	return ListCustomDbRolesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListCustomDbRolesExecute executes the request
//
//	@return []UserCustomDBRole
func (a *CustomDatabaseRolesApiService) ListCustomDbRolesExecute(r ListCustomDbRolesApiRequest) ([]UserCustomDBRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []UserCustomDBRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CustomDatabaseRolesApiService.ListCustomDbRoles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/customDBRoles/roles"
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

type UpdateCustomDbRoleApiRequest struct {
	ctx                context.Context
	ApiService         CustomDatabaseRolesApi
	groupId            string
	roleName           string
	updateCustomDBRole *UpdateCustomDBRole
}

type UpdateCustomDbRoleApiParams struct {
	GroupId            string
	RoleName           string
	UpdateCustomDBRole *UpdateCustomDBRole
}

func (a *CustomDatabaseRolesApiService) UpdateCustomDbRoleWithParams(ctx context.Context, args *UpdateCustomDbRoleApiParams) UpdateCustomDbRoleApiRequest {
	return UpdateCustomDbRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		roleName:           args.RoleName,
		updateCustomDBRole: args.UpdateCustomDBRole,
	}
}

func (r UpdateCustomDbRoleApiRequest) Execute() (*UserCustomDBRole, *http.Response, error) {
	return r.ApiService.UpdateCustomDbRoleExecute(r)
}

/*
UpdateCustomDbRole Update One Custom Role in One Project

Updates one custom role in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param roleName Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
	@return UpdateCustomDbRoleApiRequest
*/
func (a *CustomDatabaseRolesApiService) UpdateCustomDbRole(ctx context.Context, groupId string, roleName string, updateCustomDBRole *UpdateCustomDBRole) UpdateCustomDbRoleApiRequest {
	return UpdateCustomDbRoleApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            groupId,
		roleName:           roleName,
		updateCustomDBRole: updateCustomDBRole,
	}
}

// UpdateCustomDbRoleExecute executes the request
//
//	@return UserCustomDBRole
func (a *CustomDatabaseRolesApiService) UpdateCustomDbRoleExecute(r UpdateCustomDbRoleApiRequest) (*UserCustomDBRole, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserCustomDBRole
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CustomDatabaseRolesApiService.UpdateCustomDbRole")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/customDBRoles/roles/{roleName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.roleName == "" {
		return localVarReturnValue, nil, reportError("roleName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"roleName"+"}", url.PathEscape(r.roleName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateCustomDBRole == nil {
		return localVarReturnValue, nil, reportError("updateCustomDBRole is required and must be specified")
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
	localVarPostBody = r.updateCustomDBRole
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
