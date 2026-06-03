// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MaintenanceWindowsApi interface {

	/*
		DeferMaintenanceWindow Defer One Maintenance Window for One Project

		Defers the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeferMaintenanceWindowApiRequest
	*/
	DeferMaintenanceWindow(ctx context.Context, groupId string) DeferMaintenanceWindowApiRequest
	/*
		DeferMaintenanceWindow Defer One Maintenance Window for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeferMaintenanceWindowApiParams - Parameters for the request
		@return DeferMaintenanceWindowApiRequest
	*/
	DeferMaintenanceWindowWithParams(ctx context.Context, args *DeferMaintenanceWindowApiParams) DeferMaintenanceWindowApiRequest

	// Method available only for mocking purposes
	DeferMaintenanceWindowExecute(r DeferMaintenanceWindowApiRequest) (*http.Response, error)

	/*
		GetMaintenanceWindow Return One Maintenance Window for One Project

		Returns the maintenance window for the specified project. MongoDB Cloud starts those maintenance activities when needed. You can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetMaintenanceWindowApiRequest
	*/
	GetMaintenanceWindow(ctx context.Context, groupId string) GetMaintenanceWindowApiRequest
	/*
		GetMaintenanceWindow Return One Maintenance Window for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetMaintenanceWindowApiParams - Parameters for the request
		@return GetMaintenanceWindowApiRequest
	*/
	GetMaintenanceWindowWithParams(ctx context.Context, args *GetMaintenanceWindowApiParams) GetMaintenanceWindowApiRequest

	// Method available only for mocking purposes
	GetMaintenanceWindowExecute(r GetMaintenanceWindowApiRequest) (*GroupMaintenanceWindow, *http.Response, error)

	/*
		ResetMaintenanceWindow Reset One Maintenance Window for One Project

		Resets the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ResetMaintenanceWindowApiRequest
	*/
	ResetMaintenanceWindow(ctx context.Context, groupId string) ResetMaintenanceWindowApiRequest
	/*
		ResetMaintenanceWindow Reset One Maintenance Window for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ResetMaintenanceWindowApiParams - Parameters for the request
		@return ResetMaintenanceWindowApiRequest
	*/
	ResetMaintenanceWindowWithParams(ctx context.Context, args *ResetMaintenanceWindowApiParams) ResetMaintenanceWindowApiRequest

	// Method available only for mocking purposes
	ResetMaintenanceWindowExecute(r ResetMaintenanceWindowApiRequest) (*http.Response, error)

	/*
		ToggleMaintenanceAutoDefer Toggle Automatic Deferral of Maintenance for One Project

		Toggles automatic deferral of the maintenance window for the specified project. When automatic deferral is enabled, all maintenance windows are deferred for one week. This endpoint controls the same underlying feature as the `autoDeferOnceEnabled` field in the PATCH `/maintenanceWindow` endpoint. The difference is that this endpoint toggles the current value (switches from enabled to disabled or vice versa), while the `autoDeferOnceEnabled` field allows you to set a specific value. For most use cases, the PATCH endpoint with `autoDeferOnceEnabled` is recommended because it allows setting an explicit value rather than toggling.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ToggleMaintenanceAutoDeferApiRequest
	*/
	ToggleMaintenanceAutoDefer(ctx context.Context, groupId string) ToggleMaintenanceAutoDeferApiRequest
	/*
		ToggleMaintenanceAutoDefer Toggle Automatic Deferral of Maintenance for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ToggleMaintenanceAutoDeferApiParams - Parameters for the request
		@return ToggleMaintenanceAutoDeferApiRequest
	*/
	ToggleMaintenanceAutoDeferWithParams(ctx context.Context, args *ToggleMaintenanceAutoDeferApiParams) ToggleMaintenanceAutoDeferApiRequest

	// Method available only for mocking purposes
	ToggleMaintenanceAutoDeferExecute(r ToggleMaintenanceAutoDeferApiRequest) (*http.Response, error)

	/*
		UpdateMaintenanceWindow Update One Maintenance Window for One Project

		Updates the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time. Updating the maintenance window will reset any maintenance deferrals for this project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param groupMaintenanceWindow Updates the maintenance window for the specified project.
		@return UpdateMaintenanceWindowApiRequest
	*/
	UpdateMaintenanceWindow(ctx context.Context, groupId string, groupMaintenanceWindow *GroupMaintenanceWindow) UpdateMaintenanceWindowApiRequest
	/*
		UpdateMaintenanceWindow Update One Maintenance Window for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateMaintenanceWindowApiParams - Parameters for the request
		@return UpdateMaintenanceWindowApiRequest
	*/
	UpdateMaintenanceWindowWithParams(ctx context.Context, args *UpdateMaintenanceWindowApiParams) UpdateMaintenanceWindowApiRequest

	// Method available only for mocking purposes
	UpdateMaintenanceWindowExecute(r UpdateMaintenanceWindowApiRequest) (*http.Response, error)
}

// MaintenanceWindowsApiService MaintenanceWindowsApi service
type MaintenanceWindowsApiService service

type DeferMaintenanceWindowApiRequest struct {
	ctx        context.Context
	ApiService MaintenanceWindowsApi
	groupId    string
}

type DeferMaintenanceWindowApiParams struct {
	GroupId string
}

func (a *MaintenanceWindowsApiService) DeferMaintenanceWindowWithParams(ctx context.Context, args *DeferMaintenanceWindowApiParams) DeferMaintenanceWindowApiRequest {
	return DeferMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DeferMaintenanceWindowApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeferMaintenanceWindowExecute(r)
}

/*
DeferMaintenanceWindow Defer One Maintenance Window for One Project

Defers the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeferMaintenanceWindowApiRequest
*/
func (a *MaintenanceWindowsApiService) DeferMaintenanceWindow(ctx context.Context, groupId string) DeferMaintenanceWindowApiRequest {
	return DeferMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DeferMaintenanceWindowExecute executes the request
func (a *MaintenanceWindowsApiService) DeferMaintenanceWindowExecute(r DeferMaintenanceWindowApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MaintenanceWindowsApiService.DeferMaintenanceWindow")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/maintenanceWindow/defer"
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

type GetMaintenanceWindowApiRequest struct {
	ctx        context.Context
	ApiService MaintenanceWindowsApi
	groupId    string
}

type GetMaintenanceWindowApiParams struct {
	GroupId string
}

func (a *MaintenanceWindowsApiService) GetMaintenanceWindowWithParams(ctx context.Context, args *GetMaintenanceWindowApiParams) GetMaintenanceWindowApiRequest {
	return GetMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetMaintenanceWindowApiRequest) Execute() (*GroupMaintenanceWindow, *http.Response, error) {
	return r.ApiService.GetMaintenanceWindowExecute(r)
}

/*
GetMaintenanceWindow Return One Maintenance Window for One Project

Returns the maintenance window for the specified project. MongoDB Cloud starts those maintenance activities when needed. You can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetMaintenanceWindowApiRequest
*/
func (a *MaintenanceWindowsApiService) GetMaintenanceWindow(ctx context.Context, groupId string) GetMaintenanceWindowApiRequest {
	return GetMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetMaintenanceWindowExecute executes the request
//
//	@return GroupMaintenanceWindow
func (a *MaintenanceWindowsApiService) GetMaintenanceWindowExecute(r GetMaintenanceWindowApiRequest) (*GroupMaintenanceWindow, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupMaintenanceWindow
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MaintenanceWindowsApiService.GetMaintenanceWindow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/maintenanceWindow"
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

type ResetMaintenanceWindowApiRequest struct {
	ctx        context.Context
	ApiService MaintenanceWindowsApi
	groupId    string
}

type ResetMaintenanceWindowApiParams struct {
	GroupId string
}

func (a *MaintenanceWindowsApiService) ResetMaintenanceWindowWithParams(ctx context.Context, args *ResetMaintenanceWindowApiParams) ResetMaintenanceWindowApiRequest {
	return ResetMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ResetMaintenanceWindowApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.ResetMaintenanceWindowExecute(r)
}

/*
ResetMaintenanceWindow Reset One Maintenance Window for One Project

Resets the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ResetMaintenanceWindowApiRequest
*/
func (a *MaintenanceWindowsApiService) ResetMaintenanceWindow(ctx context.Context, groupId string) ResetMaintenanceWindowApiRequest {
	return ResetMaintenanceWindowApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ResetMaintenanceWindowExecute executes the request
func (a *MaintenanceWindowsApiService) ResetMaintenanceWindowExecute(r ResetMaintenanceWindowApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MaintenanceWindowsApiService.ResetMaintenanceWindow")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/maintenanceWindow"
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

type ToggleMaintenanceAutoDeferApiRequest struct {
	ctx        context.Context
	ApiService MaintenanceWindowsApi
	groupId    string
}

type ToggleMaintenanceAutoDeferApiParams struct {
	GroupId string
}

func (a *MaintenanceWindowsApiService) ToggleMaintenanceAutoDeferWithParams(ctx context.Context, args *ToggleMaintenanceAutoDeferApiParams) ToggleMaintenanceAutoDeferApiRequest {
	return ToggleMaintenanceAutoDeferApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ToggleMaintenanceAutoDeferApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.ToggleMaintenanceAutoDeferExecute(r)
}

/*
ToggleMaintenanceAutoDefer Toggle Automatic Deferral of Maintenance for One Project

Toggles automatic deferral of the maintenance window for the specified project. When automatic deferral is enabled, all maintenance windows are deferred for one week. This endpoint controls the same underlying feature as the `autoDeferOnceEnabled` field in the PATCH `/maintenanceWindow` endpoint. The difference is that this endpoint toggles the current value (switches from enabled to disabled or vice versa), while the `autoDeferOnceEnabled` field allows you to set a specific value. For most use cases, the PATCH endpoint with `autoDeferOnceEnabled` is recommended because it allows setting an explicit value rather than toggling.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ToggleMaintenanceAutoDeferApiRequest
*/
func (a *MaintenanceWindowsApiService) ToggleMaintenanceAutoDefer(ctx context.Context, groupId string) ToggleMaintenanceAutoDeferApiRequest {
	return ToggleMaintenanceAutoDeferApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ToggleMaintenanceAutoDeferExecute executes the request
func (a *MaintenanceWindowsApiService) ToggleMaintenanceAutoDeferExecute(r ToggleMaintenanceAutoDeferApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MaintenanceWindowsApiService.ToggleMaintenanceAutoDefer")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/maintenanceWindow/autoDefer"
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

type UpdateMaintenanceWindowApiRequest struct {
	ctx                    context.Context
	ApiService             MaintenanceWindowsApi
	groupId                string
	groupMaintenanceWindow *GroupMaintenanceWindow
}

type UpdateMaintenanceWindowApiParams struct {
	GroupId                string
	GroupMaintenanceWindow *GroupMaintenanceWindow
}

func (a *MaintenanceWindowsApiService) UpdateMaintenanceWindowWithParams(ctx context.Context, args *UpdateMaintenanceWindowApiParams) UpdateMaintenanceWindowApiRequest {
	return UpdateMaintenanceWindowApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		groupMaintenanceWindow: args.GroupMaintenanceWindow,
	}
}

func (r UpdateMaintenanceWindowApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.UpdateMaintenanceWindowExecute(r)
}

/*
UpdateMaintenanceWindow Update One Maintenance Window for One Project

Updates the maintenance window for the specified project. Urgent maintenance activities such as security patches can't wait for your chosen window. MongoDB Cloud starts those maintenance activities when needed. After you schedule maintenance for your cluster, you can't change your maintenance window until the current maintenance efforts complete. The maintenance procedure that MongoDB Cloud performs requires at least one replica set election during the maintenance window per replica set. Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or unexpected system issues could delay the start time. Updating the maintenance window will reset any maintenance deferrals for this project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateMaintenanceWindowApiRequest
*/
func (a *MaintenanceWindowsApiService) UpdateMaintenanceWindow(ctx context.Context, groupId string, groupMaintenanceWindow *GroupMaintenanceWindow) UpdateMaintenanceWindowApiRequest {
	return UpdateMaintenanceWindowApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		groupMaintenanceWindow: groupMaintenanceWindow,
	}
}

// UpdateMaintenanceWindowExecute executes the request
func (a *MaintenanceWindowsApiService) UpdateMaintenanceWindowExecute(r UpdateMaintenanceWindowApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPatch
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MaintenanceWindowsApiService.UpdateMaintenanceWindow")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/maintenanceWindow"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupMaintenanceWindow == nil {
		return nil, reportError("groupMaintenanceWindow is required and must be specified")
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
	localVarPostBody = r.groupMaintenanceWindow
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
