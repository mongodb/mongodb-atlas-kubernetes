// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PushBasedLogExportApi interface {

	/*
		CreateGroupLogIntegration Create One Log Integration

		Creates a new log integration configuration identified by a unique ID.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param logIntegrationRequest Log integration configuration to create.
		@return CreateGroupLogIntegrationApiRequest
	*/
	CreateGroupLogIntegration(ctx context.Context, groupId string, logIntegrationRequest *LogIntegrationRequest) CreateGroupLogIntegrationApiRequest
	/*
		CreateGroupLogIntegration Create One Log Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupLogIntegrationApiParams - Parameters for the request
		@return CreateGroupLogIntegrationApiRequest
	*/
	CreateGroupLogIntegrationWithParams(ctx context.Context, args *CreateGroupLogIntegrationApiParams) CreateGroupLogIntegrationApiRequest

	// Method available only for mocking purposes
	CreateGroupLogIntegrationExecute(r CreateGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error)

	/*
		CreateLogExport Create One Push-Based Log Export Configuration in One Project

		Configures the project level settings for the push-based log export feature.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param createPushBasedLogExportProjectRequest The project configuration details. The S3 bucket name, IAM role ID, and prefix path fields are required.
		@return CreateLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	CreateLogExport(ctx context.Context, groupId string, createPushBasedLogExportProjectRequest *CreatePushBasedLogExportProjectRequest) CreateLogExportApiRequest
	/*
		CreateLogExport Create One Push-Based Log Export Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateLogExportApiParams - Parameters for the request
		@return CreateLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	CreateLogExportWithParams(ctx context.Context, args *CreateLogExportApiParams) CreateLogExportApiRequest

	// Method available only for mocking purposes
	CreateLogExportExecute(r CreateLogExportApiRequest) (*http.Response, error)

	/*
		DeleteGroupLogIntegration Remove One Log Integration

		Removes the configuration for one log integration identified by its unique ID.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id Unique identifier of the log integration configuration.
		@return DeleteGroupLogIntegrationApiRequest
	*/
	DeleteGroupLogIntegration(ctx context.Context, groupId string, id string) DeleteGroupLogIntegrationApiRequest
	/*
		DeleteGroupLogIntegration Remove One Log Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupLogIntegrationApiParams - Parameters for the request
		@return DeleteGroupLogIntegrationApiRequest
	*/
	DeleteGroupLogIntegrationWithParams(ctx context.Context, args *DeleteGroupLogIntegrationApiParams) DeleteGroupLogIntegrationApiRequest

	// Method available only for mocking purposes
	DeleteGroupLogIntegrationExecute(r DeleteGroupLogIntegrationApiRequest) (*http.Response, error)

	/*
		DeleteLogExport Disable Push-Based Log Export for One Project

		Disables the push-based log export feature by resetting the project level settings to its default configuration.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	DeleteLogExport(ctx context.Context, groupId string) DeleteLogExportApiRequest
	/*
		DeleteLogExport Disable Push-Based Log Export for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteLogExportApiParams - Parameters for the request
		@return DeleteLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	DeleteLogExportWithParams(ctx context.Context, args *DeleteLogExportApiParams) DeleteLogExportApiRequest

	// Method available only for mocking purposes
	DeleteLogExportExecute(r DeleteLogExportApiRequest) (*http.Response, error)

	/*
		GetGroupLogIntegration Return One Log Integration

		Returns the configuration for one log integration identified by its unique ID.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id Unique identifier of the log integration configuration.
		@return GetGroupLogIntegrationApiRequest
	*/
	GetGroupLogIntegration(ctx context.Context, groupId string, id string) GetGroupLogIntegrationApiRequest
	/*
		GetGroupLogIntegration Return One Log Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupLogIntegrationApiParams - Parameters for the request
		@return GetGroupLogIntegrationApiRequest
	*/
	GetGroupLogIntegrationWithParams(ctx context.Context, args *GetGroupLogIntegrationApiParams) GetGroupLogIntegrationApiRequest

	// Method available only for mocking purposes
	GetGroupLogIntegrationExecute(r GetGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error)

	/*
		GetLogExport Return One Push-Based Log Export Configuration in One Project

		Fetches the current project level settings for the push-based log export feature.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	GetLogExport(ctx context.Context, groupId string) GetLogExportApiRequest
	/*
		GetLogExport Return One Push-Based Log Export Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetLogExportApiParams - Parameters for the request
		@return GetLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	GetLogExportWithParams(ctx context.Context, args *GetLogExportApiParams) GetLogExportApiRequest

	// Method available only for mocking purposes
	GetLogExportExecute(r GetLogExportApiRequest) (*PushBasedLogExportProject, *http.Response, error)

	/*
		ListGroupLogIntegrations Return All Active Log Integrations

		Returns all log integration configurations for the project. Optionally filter by integration type.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupLogIntegrationsApiRequest
	*/
	ListGroupLogIntegrations(ctx context.Context, groupId string) ListGroupLogIntegrationsApiRequest
	/*
		ListGroupLogIntegrations Return All Active Log Integrations


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupLogIntegrationsApiParams - Parameters for the request
		@return ListGroupLogIntegrationsApiRequest
	*/
	ListGroupLogIntegrationsWithParams(ctx context.Context, args *ListGroupLogIntegrationsApiParams) ListGroupLogIntegrationsApiRequest

	// Method available only for mocking purposes
	ListGroupLogIntegrationsExecute(r ListGroupLogIntegrationsApiRequest) (*PaginatedLogIntegrationResponse, *http.Response, error)

	/*
		UpdateGroupLogIntegration Update One Log Integration

		Updates the configuration for one log integration identified by its unique ID.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id Unique identifier of the log integration configuration.
		@param logIntegrationRequest Updated log integration configuration.
		@return UpdateGroupLogIntegrationApiRequest
	*/
	UpdateGroupLogIntegration(ctx context.Context, groupId string, id string, logIntegrationRequest *LogIntegrationRequest) UpdateGroupLogIntegrationApiRequest
	/*
		UpdateGroupLogIntegration Update One Log Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupLogIntegrationApiParams - Parameters for the request
		@return UpdateGroupLogIntegrationApiRequest
	*/
	UpdateGroupLogIntegrationWithParams(ctx context.Context, args *UpdateGroupLogIntegrationApiParams) UpdateGroupLogIntegrationApiRequest

	// Method available only for mocking purposes
	UpdateGroupLogIntegrationExecute(r UpdateGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error)

	/*
		UpdateLogExport Update One Push-Based Log Export Configuration in One Project

		Updates the project level settings for the push-based log export feature.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pushBasedLogExportProject The project configuration details. The S3 bucket name, IAM role ID, and prefix path fields are the only fields that may be specified. Fields left unspecified will not be modified.
		@return UpdateLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	UpdateLogExport(ctx context.Context, groupId string, pushBasedLogExportProject *PushBasedLogExportProject) UpdateLogExportApiRequest
	/*
		UpdateLogExport Update One Push-Based Log Export Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateLogExportApiParams - Parameters for the request
		@return UpdateLogExportApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for PushBasedLogExportApi
	*/
	UpdateLogExportWithParams(ctx context.Context, args *UpdateLogExportApiParams) UpdateLogExportApiRequest

	// Method available only for mocking purposes
	UpdateLogExportExecute(r UpdateLogExportApiRequest) (*http.Response, error)
}

// PushBasedLogExportApiService PushBasedLogExportApi service
type PushBasedLogExportApiService service

type CreateGroupLogIntegrationApiRequest struct {
	ctx                   context.Context
	ApiService            PushBasedLogExportApi
	groupId               string
	logIntegrationRequest *LogIntegrationRequest
}

type CreateGroupLogIntegrationApiParams struct {
	GroupId               string
	LogIntegrationRequest *LogIntegrationRequest
}

func (a *PushBasedLogExportApiService) CreateGroupLogIntegrationWithParams(ctx context.Context, args *CreateGroupLogIntegrationApiParams) CreateGroupLogIntegrationApiRequest {
	return CreateGroupLogIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		logIntegrationRequest: args.LogIntegrationRequest,
	}
}

func (r CreateGroupLogIntegrationApiRequest) Execute() (*LogIntegrationResponse, *http.Response, error) {
	return r.ApiService.CreateGroupLogIntegrationExecute(r)
}

/*
CreateGroupLogIntegration Create One Log Integration

Creates a new log integration configuration identified by a unique ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupLogIntegrationApiRequest
*/
func (a *PushBasedLogExportApiService) CreateGroupLogIntegration(ctx context.Context, groupId string, logIntegrationRequest *LogIntegrationRequest) CreateGroupLogIntegrationApiRequest {
	return CreateGroupLogIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               groupId,
		logIntegrationRequest: logIntegrationRequest,
	}
}

// CreateGroupLogIntegrationExecute executes the request
//
//	@return LogIntegrationResponse
func (a *PushBasedLogExportApiService) CreateGroupLogIntegrationExecute(r CreateGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LogIntegrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.CreateGroupLogIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/logIntegrations"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.logIntegrationRequest == nil {
		return localVarReturnValue, nil, reportError("logIntegrationRequest is required and must be specified")
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
	localVarPostBody = r.logIntegrationRequest
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

type CreateLogExportApiRequest struct {
	ctx                                    context.Context
	ApiService                             PushBasedLogExportApi
	groupId                                string
	createPushBasedLogExportProjectRequest *CreatePushBasedLogExportProjectRequest
}

type CreateLogExportApiParams struct {
	GroupId                                string
	CreatePushBasedLogExportProjectRequest *CreatePushBasedLogExportProjectRequest
}

func (a *PushBasedLogExportApiService) CreateLogExportWithParams(ctx context.Context, args *CreateLogExportApiParams) CreateLogExportApiRequest {
	return CreateLogExportApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                args.GroupId,
		createPushBasedLogExportProjectRequest: args.CreatePushBasedLogExportProjectRequest,
	}
}

func (r CreateLogExportApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.CreateLogExportExecute(r)
}

/*
CreateLogExport Create One Push-Based Log Export Configuration in One Project

Configures the project level settings for the push-based log export feature.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateLogExportApiRequest

Deprecated
*/
func (a *PushBasedLogExportApiService) CreateLogExport(ctx context.Context, groupId string, createPushBasedLogExportProjectRequest *CreatePushBasedLogExportProjectRequest) CreateLogExportApiRequest {
	return CreateLogExportApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                groupId,
		createPushBasedLogExportProjectRequest: createPushBasedLogExportProjectRequest,
	}
}

// CreateLogExportExecute executes the request
// Deprecated
func (a *PushBasedLogExportApiService) CreateLogExportExecute(r CreateLogExportApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.CreateLogExport")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pushBasedLogExport"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createPushBasedLogExportProjectRequest == nil {
		return nil, reportError("createPushBasedLogExportProjectRequest is required and must be specified")
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
	localVarPostBody = r.createPushBasedLogExportProjectRequest
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

type DeleteGroupLogIntegrationApiRequest struct {
	ctx        context.Context
	ApiService PushBasedLogExportApi
	groupId    string
	id         string
}

type DeleteGroupLogIntegrationApiParams struct {
	GroupId string
	Id      string
}

func (a *PushBasedLogExportApiService) DeleteGroupLogIntegrationWithParams(ctx context.Context, args *DeleteGroupLogIntegrationApiParams) DeleteGroupLogIntegrationApiRequest {
	return DeleteGroupLogIntegrationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		id:         args.Id,
	}
}

func (r DeleteGroupLogIntegrationApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupLogIntegrationExecute(r)
}

/*
DeleteGroupLogIntegration Remove One Log Integration

Removes the configuration for one log integration identified by its unique ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id Unique identifier of the log integration configuration.
	@return DeleteGroupLogIntegrationApiRequest
*/
func (a *PushBasedLogExportApiService) DeleteGroupLogIntegration(ctx context.Context, groupId string, id string) DeleteGroupLogIntegrationApiRequest {
	return DeleteGroupLogIntegrationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		id:         id,
	}
}

// DeleteGroupLogIntegrationExecute executes the request
func (a *PushBasedLogExportApiService) DeleteGroupLogIntegrationExecute(r DeleteGroupLogIntegrationApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.DeleteGroupLogIntegration")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/logIntegrations/{id}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

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

type DeleteLogExportApiRequest struct {
	ctx        context.Context
	ApiService PushBasedLogExportApi
	groupId    string
}

type DeleteLogExportApiParams struct {
	GroupId string
}

func (a *PushBasedLogExportApiService) DeleteLogExportWithParams(ctx context.Context, args *DeleteLogExportApiParams) DeleteLogExportApiRequest {
	return DeleteLogExportApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DeleteLogExportApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteLogExportExecute(r)
}

/*
DeleteLogExport Disable Push-Based Log Export for One Project

Disables the push-based log export feature by resetting the project level settings to its default configuration.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteLogExportApiRequest

Deprecated
*/
func (a *PushBasedLogExportApiService) DeleteLogExport(ctx context.Context, groupId string) DeleteLogExportApiRequest {
	return DeleteLogExportApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DeleteLogExportExecute executes the request
// Deprecated
func (a *PushBasedLogExportApiService) DeleteLogExportExecute(r DeleteLogExportApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.DeleteLogExport")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pushBasedLogExport"
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

type GetGroupLogIntegrationApiRequest struct {
	ctx        context.Context
	ApiService PushBasedLogExportApi
	groupId    string
	id         string
}

type GetGroupLogIntegrationApiParams struct {
	GroupId string
	Id      string
}

func (a *PushBasedLogExportApiService) GetGroupLogIntegrationWithParams(ctx context.Context, args *GetGroupLogIntegrationApiParams) GetGroupLogIntegrationApiRequest {
	return GetGroupLogIntegrationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		id:         args.Id,
	}
}

func (r GetGroupLogIntegrationApiRequest) Execute() (*LogIntegrationResponse, *http.Response, error) {
	return r.ApiService.GetGroupLogIntegrationExecute(r)
}

/*
GetGroupLogIntegration Return One Log Integration

Returns the configuration for one log integration identified by its unique ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id Unique identifier of the log integration configuration.
	@return GetGroupLogIntegrationApiRequest
*/
func (a *PushBasedLogExportApiService) GetGroupLogIntegration(ctx context.Context, groupId string, id string) GetGroupLogIntegrationApiRequest {
	return GetGroupLogIntegrationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		id:         id,
	}
}

// GetGroupLogIntegrationExecute executes the request
//
//	@return LogIntegrationResponse
func (a *PushBasedLogExportApiService) GetGroupLogIntegrationExecute(r GetGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LogIntegrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.GetGroupLogIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/logIntegrations/{id}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return localVarReturnValue, nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

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

type GetLogExportApiRequest struct {
	ctx        context.Context
	ApiService PushBasedLogExportApi
	groupId    string
}

type GetLogExportApiParams struct {
	GroupId string
}

func (a *PushBasedLogExportApiService) GetLogExportWithParams(ctx context.Context, args *GetLogExportApiParams) GetLogExportApiRequest {
	return GetLogExportApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetLogExportApiRequest) Execute() (*PushBasedLogExportProject, *http.Response, error) {
	return r.ApiService.GetLogExportExecute(r)
}

/*
GetLogExport Return One Push-Based Log Export Configuration in One Project

Fetches the current project level settings for the push-based log export feature.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetLogExportApiRequest

Deprecated
*/
func (a *PushBasedLogExportApiService) GetLogExport(ctx context.Context, groupId string) GetLogExportApiRequest {
	return GetLogExportApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetLogExportExecute executes the request
//
//	@return PushBasedLogExportProject
//
// Deprecated
func (a *PushBasedLogExportApiService) GetLogExportExecute(r GetLogExportApiRequest) (*PushBasedLogExportProject, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PushBasedLogExportProject
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.GetLogExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pushBasedLogExport"
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

type ListGroupLogIntegrationsApiRequest struct {
	ctx             context.Context
	ApiService      PushBasedLogExportApi
	groupId         string
	includeCount    *bool
	itemsPerPage    *int
	pageNum         *int
	integrationType *string
}

type ListGroupLogIntegrationsApiParams struct {
	GroupId         string
	IncludeCount    *bool
	ItemsPerPage    *int
	PageNum         *int
	IntegrationType *string
}

func (a *PushBasedLogExportApiService) ListGroupLogIntegrationsWithParams(ctx context.Context, args *ListGroupLogIntegrationsApiParams) ListGroupLogIntegrationsApiRequest {
	return ListGroupLogIntegrationsApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         args.GroupId,
		includeCount:    args.IncludeCount,
		itemsPerPage:    args.ItemsPerPage,
		pageNum:         args.PageNum,
		integrationType: args.IntegrationType,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupLogIntegrationsApiRequest) IncludeCount(includeCount bool) ListGroupLogIntegrationsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupLogIntegrationsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupLogIntegrationsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupLogIntegrationsApiRequest) PageNum(pageNum int) ListGroupLogIntegrationsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Optional filter by integration type (e.g., &#x60;S3_LOG_EXPORT&#x60;).
func (r ListGroupLogIntegrationsApiRequest) IntegrationType(integrationType string) ListGroupLogIntegrationsApiRequest {
	r.integrationType = &integrationType
	return r
}

func (r ListGroupLogIntegrationsApiRequest) Execute() (*PaginatedLogIntegrationResponse, *http.Response, error) {
	return r.ApiService.ListGroupLogIntegrationsExecute(r)
}

/*
ListGroupLogIntegrations Return All Active Log Integrations

Returns all log integration configurations for the project. Optionally filter by integration type.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupLogIntegrationsApiRequest
*/
func (a *PushBasedLogExportApiService) ListGroupLogIntegrations(ctx context.Context, groupId string) ListGroupLogIntegrationsApiRequest {
	return ListGroupLogIntegrationsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupLogIntegrationsExecute executes the request
//
//	@return PaginatedLogIntegrationResponse
func (a *PushBasedLogExportApiService) ListGroupLogIntegrationsExecute(r ListGroupLogIntegrationsApiRequest) (*PaginatedLogIntegrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedLogIntegrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.ListGroupLogIntegrations")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/logIntegrations"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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
	if r.integrationType != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "integrationType", r.integrationType, "")
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

type UpdateGroupLogIntegrationApiRequest struct {
	ctx                   context.Context
	ApiService            PushBasedLogExportApi
	groupId               string
	id                    string
	logIntegrationRequest *LogIntegrationRequest
}

type UpdateGroupLogIntegrationApiParams struct {
	GroupId               string
	Id                    string
	LogIntegrationRequest *LogIntegrationRequest
}

func (a *PushBasedLogExportApiService) UpdateGroupLogIntegrationWithParams(ctx context.Context, args *UpdateGroupLogIntegrationApiParams) UpdateGroupLogIntegrationApiRequest {
	return UpdateGroupLogIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		id:                    args.Id,
		logIntegrationRequest: args.LogIntegrationRequest,
	}
}

func (r UpdateGroupLogIntegrationApiRequest) Execute() (*LogIntegrationResponse, *http.Response, error) {
	return r.ApiService.UpdateGroupLogIntegrationExecute(r)
}

/*
UpdateGroupLogIntegration Update One Log Integration

Updates the configuration for one log integration identified by its unique ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id Unique identifier of the log integration configuration.
	@return UpdateGroupLogIntegrationApiRequest
*/
func (a *PushBasedLogExportApiService) UpdateGroupLogIntegration(ctx context.Context, groupId string, id string, logIntegrationRequest *LogIntegrationRequest) UpdateGroupLogIntegrationApiRequest {
	return UpdateGroupLogIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               groupId,
		id:                    id,
		logIntegrationRequest: logIntegrationRequest,
	}
}

// UpdateGroupLogIntegrationExecute executes the request
//
//	@return LogIntegrationResponse
func (a *PushBasedLogExportApiService) UpdateGroupLogIntegrationExecute(r UpdateGroupLogIntegrationApiRequest) (*LogIntegrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LogIntegrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.UpdateGroupLogIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/logIntegrations/{id}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return localVarReturnValue, nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.logIntegrationRequest == nil {
		return localVarReturnValue, nil, reportError("logIntegrationRequest is required and must be specified")
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
	localVarPostBody = r.logIntegrationRequest
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

type UpdateLogExportApiRequest struct {
	ctx                       context.Context
	ApiService                PushBasedLogExportApi
	groupId                   string
	pushBasedLogExportProject *PushBasedLogExportProject
}

type UpdateLogExportApiParams struct {
	GroupId                   string
	PushBasedLogExportProject *PushBasedLogExportProject
}

func (a *PushBasedLogExportApiService) UpdateLogExportWithParams(ctx context.Context, args *UpdateLogExportApiParams) UpdateLogExportApiRequest {
	return UpdateLogExportApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   args.GroupId,
		pushBasedLogExportProject: args.PushBasedLogExportProject,
	}
}

func (r UpdateLogExportApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.UpdateLogExportExecute(r)
}

/*
UpdateLogExport Update One Push-Based Log Export Configuration in One Project

Updates the project level settings for the push-based log export feature.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateLogExportApiRequest

Deprecated
*/
func (a *PushBasedLogExportApiService) UpdateLogExport(ctx context.Context, groupId string, pushBasedLogExportProject *PushBasedLogExportProject) UpdateLogExportApiRequest {
	return UpdateLogExportApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   groupId,
		pushBasedLogExportProject: pushBasedLogExportProject,
	}
}

// UpdateLogExportExecute executes the request
// Deprecated
func (a *PushBasedLogExportApiService) UpdateLogExportExecute(r UpdateLogExportApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPatch
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PushBasedLogExportApiService.UpdateLogExport")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pushBasedLogExport"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.pushBasedLogExportProject == nil {
		return nil, reportError("pushBasedLogExportProject is required and must be specified")
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
	localVarPostBody = r.pushBasedLogExportProject
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
