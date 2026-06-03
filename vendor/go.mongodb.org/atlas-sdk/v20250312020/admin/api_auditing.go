// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AuditingApi interface {

	/*
		GetGroupAuditLog Return Auditing Configuration for One Project

		Returns the auditing configuration for the specified project. The auditing configuration defines the events that MongoDB Cloud records in the audit log. This feature isn't available for `M0`, `M2`, `M5`, or serverless clusters.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetGroupAuditLogApiRequest
	*/
	GetGroupAuditLog(ctx context.Context, groupId string) GetGroupAuditLogApiRequest
	/*
		GetGroupAuditLog Return Auditing Configuration for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupAuditLogApiParams - Parameters for the request
		@return GetGroupAuditLogApiRequest
	*/
	GetGroupAuditLogWithParams(ctx context.Context, args *GetGroupAuditLogApiParams) GetGroupAuditLogApiRequest

	// Method available only for mocking purposes
	GetGroupAuditLogExecute(r GetGroupAuditLogApiRequest) (*AuditLog, *http.Response, error)

	/*
		UpdateAuditLog Update Auditing Configuration for One Project

		Updates the auditing configuration for the specified project. The auditing configuration defines the events that MongoDB Cloud records in the audit log. This feature isn't available for `M0`, `M2`, `M5`, or serverless clusters.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param auditLog Updated auditing configuration for the specified project.
		@return UpdateAuditLogApiRequest
	*/
	UpdateAuditLog(ctx context.Context, groupId string, auditLog *AuditLog) UpdateAuditLogApiRequest
	/*
		UpdateAuditLog Update Auditing Configuration for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateAuditLogApiParams - Parameters for the request
		@return UpdateAuditLogApiRequest
	*/
	UpdateAuditLogWithParams(ctx context.Context, args *UpdateAuditLogApiParams) UpdateAuditLogApiRequest

	// Method available only for mocking purposes
	UpdateAuditLogExecute(r UpdateAuditLogApiRequest) (*AuditLog, *http.Response, error)
}

// AuditingApiService AuditingApi service
type AuditingApiService service

type GetGroupAuditLogApiRequest struct {
	ctx        context.Context
	ApiService AuditingApi
	groupId    string
}

type GetGroupAuditLogApiParams struct {
	GroupId string
}

func (a *AuditingApiService) GetGroupAuditLogWithParams(ctx context.Context, args *GetGroupAuditLogApiParams) GetGroupAuditLogApiRequest {
	return GetGroupAuditLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetGroupAuditLogApiRequest) Execute() (*AuditLog, *http.Response, error) {
	return r.ApiService.GetGroupAuditLogExecute(r)
}

/*
GetGroupAuditLog Return Auditing Configuration for One Project

Returns the auditing configuration for the specified project. The auditing configuration defines the events that MongoDB Cloud records in the audit log. This feature isn't available for `M0`, `M2`, `M5`, or serverless clusters.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetGroupAuditLogApiRequest
*/
func (a *AuditingApiService) GetGroupAuditLog(ctx context.Context, groupId string) GetGroupAuditLogApiRequest {
	return GetGroupAuditLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetGroupAuditLogExecute executes the request
//
//	@return AuditLog
func (a *AuditingApiService) GetGroupAuditLogExecute(r GetGroupAuditLogApiRequest) (*AuditLog, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AuditLog
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuditingApiService.GetGroupAuditLog")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/auditLog"
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

type UpdateAuditLogApiRequest struct {
	ctx        context.Context
	ApiService AuditingApi
	groupId    string
	auditLog   *AuditLog
}

type UpdateAuditLogApiParams struct {
	GroupId  string
	AuditLog *AuditLog
}

func (a *AuditingApiService) UpdateAuditLogWithParams(ctx context.Context, args *UpdateAuditLogApiParams) UpdateAuditLogApiRequest {
	return UpdateAuditLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		auditLog:   args.AuditLog,
	}
}

func (r UpdateAuditLogApiRequest) Execute() (*AuditLog, *http.Response, error) {
	return r.ApiService.UpdateAuditLogExecute(r)
}

/*
UpdateAuditLog Update Auditing Configuration for One Project

Updates the auditing configuration for the specified project. The auditing configuration defines the events that MongoDB Cloud records in the audit log. This feature isn't available for `M0`, `M2`, `M5`, or serverless clusters.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateAuditLogApiRequest
*/
func (a *AuditingApiService) UpdateAuditLog(ctx context.Context, groupId string, auditLog *AuditLog) UpdateAuditLogApiRequest {
	return UpdateAuditLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		auditLog:   auditLog,
	}
}

// UpdateAuditLogExecute executes the request
//
//	@return AuditLog
func (a *AuditingApiService) UpdateAuditLogExecute(r UpdateAuditLogApiRequest) (*AuditLog, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AuditLog
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuditingApiService.UpdateAuditLog")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/auditLog"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.auditLog == nil {
		return localVarReturnValue, nil, reportError("auditLog is required and must be specified")
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
	localVarPostBody = r.auditLog
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
