// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AlertConfigurationsApi interface {

	/*
			CreateAlertConfig Create One Alert Configuration in One Project

			Creates one alert configuration for the specified project. Alert configurations define the triggers and notification methods for alerts.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param groupAlertsConfig Creates one alert configuration for the specified project.
			@return CreateAlertConfigApiRequest
	*/
	CreateAlertConfig(ctx context.Context, groupId string, groupAlertsConfig *GroupAlertsConfig) CreateAlertConfigApiRequest
	/*
		CreateAlertConfig Create One Alert Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateAlertConfigApiParams - Parameters for the request
		@return CreateAlertConfigApiRequest
	*/
	CreateAlertConfigWithParams(ctx context.Context, args *CreateAlertConfigApiParams) CreateAlertConfigApiRequest

	// Method available only for mocking purposes
	CreateAlertConfigExecute(r CreateAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error)

	/*
			DeleteAlertConfig Remove One Alert Configuration from One Project

			Removes one alert configuration from the specified project. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
			@return DeleteAlertConfigApiRequest
	*/
	DeleteAlertConfig(ctx context.Context, groupId string, alertConfigId string) DeleteAlertConfigApiRequest
	/*
		DeleteAlertConfig Remove One Alert Configuration from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteAlertConfigApiParams - Parameters for the request
		@return DeleteAlertConfigApiRequest
	*/
	DeleteAlertConfigWithParams(ctx context.Context, args *DeleteAlertConfigApiParams) DeleteAlertConfigApiRequest

	// Method available only for mocking purposes
	DeleteAlertConfigExecute(r DeleteAlertConfigApiRequest) (*http.Response, error)

	/*
			GetAlertConfig Return One Alert Configuration from One Project

			Returns the specified alert configuration from the specified project. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
			@return GetAlertConfigApiRequest
	*/
	GetAlertConfig(ctx context.Context, groupId string, alertConfigId string) GetAlertConfigApiRequest
	/*
		GetAlertConfig Return One Alert Configuration from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAlertConfigApiParams - Parameters for the request
		@return GetAlertConfigApiRequest
	*/
	GetAlertConfigWithParams(ctx context.Context, args *GetAlertConfigApiParams) GetAlertConfigApiRequest

	// Method available only for mocking purposes
	GetAlertConfigExecute(r GetAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error)

	/*
			GetAlertConfigs Return All Alert Configurations Set for One Alert

			Returns all alert configurations set for the specified alert. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertId Unique 24-hexadecimal digit string that identifies the alert.
			@return GetAlertConfigsApiRequest
	*/
	GetAlertConfigs(ctx context.Context, groupId string, alertId string) GetAlertConfigsApiRequest
	/*
		GetAlertConfigs Return All Alert Configurations Set for One Alert


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAlertConfigsApiParams - Parameters for the request
		@return GetAlertConfigsApiRequest
	*/
	GetAlertConfigsWithParams(ctx context.Context, args *GetAlertConfigsApiParams) GetAlertConfigsApiRequest

	// Method available only for mocking purposes
	GetAlertConfigsExecute(r GetAlertConfigsApiRequest) (*PaginatedAlertConfig, *http.Response, error)

	/*
			ListAlertConfigs Return All Alert Configurations in One Project

			Returns all alert configurations for one project. These alert configurations apply to any component in the project. Alert configurations define the triggers and notification methods for alerts.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListAlertConfigsApiRequest
	*/
	ListAlertConfigs(ctx context.Context, groupId string) ListAlertConfigsApiRequest
	/*
		ListAlertConfigs Return All Alert Configurations in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListAlertConfigsApiParams - Parameters for the request
		@return ListAlertConfigsApiRequest
	*/
	ListAlertConfigsWithParams(ctx context.Context, args *ListAlertConfigsApiParams) ListAlertConfigsApiRequest

	// Method available only for mocking purposes
	ListAlertConfigsExecute(r ListAlertConfigsApiRequest) (*PaginatedAlertConfig, *http.Response, error)

	/*
		ListMatcherFieldNames Return All Alert Configuration Matchers Field Names

		Get all field names that the `matchers.fieldName` parameter accepts when you create or update an Alert Configuration. You can successfully call this endpoint with any assigned role.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListMatcherFieldNamesApiRequest
	*/
	ListMatcherFieldNames(ctx context.Context) ListMatcherFieldNamesApiRequest
	/*
		ListMatcherFieldNames Return All Alert Configuration Matchers Field Names


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListMatcherFieldNamesApiParams - Parameters for the request
		@return ListMatcherFieldNamesApiRequest
	*/
	ListMatcherFieldNamesWithParams(ctx context.Context, args *ListMatcherFieldNamesApiParams) ListMatcherFieldNamesApiRequest

	// Method available only for mocking purposes
	ListMatcherFieldNamesExecute(r ListMatcherFieldNamesApiRequest) ([]string, *http.Response, error)

	/*
			ToggleAlertConfig Toggle State of One Alert Configuration in One Project

			Enables or disables the specified alert configuration in the specified project. The resource enables the specified alert configuration if currently enabled. The resource disables the specified alert configuration if currently disabled.

		**NOTE**: This endpoint updates only the enabled/disabled state for the alert configuration. To update more than just this configuration, see Update One Alert Configuration.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration that triggered this alert.
			@param alertsToggle Enables or disables the specified alert configuration in the specified project.
			@return ToggleAlertConfigApiRequest
	*/
	ToggleAlertConfig(ctx context.Context, groupId string, alertConfigId string, alertsToggle *AlertsToggle) ToggleAlertConfigApiRequest
	/*
		ToggleAlertConfig Toggle State of One Alert Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ToggleAlertConfigApiParams - Parameters for the request
		@return ToggleAlertConfigApiRequest
	*/
	ToggleAlertConfigWithParams(ctx context.Context, args *ToggleAlertConfigApiParams) ToggleAlertConfigApiRequest

	// Method available only for mocking purposes
	ToggleAlertConfigExecute(r ToggleAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error)

	/*
			UpdateAlertConfig Update One Alert Configuration in One Project

			Updates one alert configuration in the specified project. Alert configurations define the triggers and notification methods for alerts.

		**NOTE**: To enable or disable the alert configuration, see Toggle One State of One Alert Configuration in One Project.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
			@param groupAlertsConfig Updates one alert configuration in the specified project.
			@return UpdateAlertConfigApiRequest
	*/
	UpdateAlertConfig(ctx context.Context, groupId string, alertConfigId string, groupAlertsConfig *GroupAlertsConfig) UpdateAlertConfigApiRequest
	/*
		UpdateAlertConfig Update One Alert Configuration in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateAlertConfigApiParams - Parameters for the request
		@return UpdateAlertConfigApiRequest
	*/
	UpdateAlertConfigWithParams(ctx context.Context, args *UpdateAlertConfigApiParams) UpdateAlertConfigApiRequest

	// Method available only for mocking purposes
	UpdateAlertConfigExecute(r UpdateAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error)
}

// AlertConfigurationsApiService AlertConfigurationsApi service
type AlertConfigurationsApiService service

type CreateAlertConfigApiRequest struct {
	ctx               context.Context
	ApiService        AlertConfigurationsApi
	groupId           string
	groupAlertsConfig *GroupAlertsConfig
}

type CreateAlertConfigApiParams struct {
	GroupId           string
	GroupAlertsConfig *GroupAlertsConfig
}

func (a *AlertConfigurationsApiService) CreateAlertConfigWithParams(ctx context.Context, args *CreateAlertConfigApiParams) CreateAlertConfigApiRequest {
	return CreateAlertConfigApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		groupAlertsConfig: args.GroupAlertsConfig,
	}
}

func (r CreateAlertConfigApiRequest) Execute() (*GroupAlertsConfig, *http.Response, error) {
	return r.ApiService.CreateAlertConfigExecute(r)
}

/*
CreateAlertConfig Create One Alert Configuration in One Project

Creates one alert configuration for the specified project. Alert configurations define the triggers and notification methods for alerts.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateAlertConfigApiRequest
*/
func (a *AlertConfigurationsApiService) CreateAlertConfig(ctx context.Context, groupId string, groupAlertsConfig *GroupAlertsConfig) CreateAlertConfigApiRequest {
	return CreateAlertConfigApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		groupAlertsConfig: groupAlertsConfig,
	}
}

// CreateAlertConfigExecute executes the request
//
//	@return GroupAlertsConfig
func (a *AlertConfigurationsApiService) CreateAlertConfigExecute(r CreateAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupAlertsConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.CreateAlertConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupAlertsConfig == nil {
		return localVarReturnValue, nil, reportError("groupAlertsConfig is required and must be specified")
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
	localVarPostBody = r.groupAlertsConfig
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

type DeleteAlertConfigApiRequest struct {
	ctx           context.Context
	ApiService    AlertConfigurationsApi
	groupId       string
	alertConfigId string
}

type DeleteAlertConfigApiParams struct {
	GroupId       string
	AlertConfigId string
}

func (a *AlertConfigurationsApiService) DeleteAlertConfigWithParams(ctx context.Context, args *DeleteAlertConfigApiParams) DeleteAlertConfigApiRequest {
	return DeleteAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		alertConfigId: args.AlertConfigId,
	}
}

func (r DeleteAlertConfigApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteAlertConfigExecute(r)
}

/*
DeleteAlertConfig Remove One Alert Configuration from One Project

Removes one alert configuration from the specified project. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
	@return DeleteAlertConfigApiRequest
*/
func (a *AlertConfigurationsApiService) DeleteAlertConfig(ctx context.Context, groupId string, alertConfigId string) DeleteAlertConfigApiRequest {
	return DeleteAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		alertConfigId: alertConfigId,
	}
}

// DeleteAlertConfigExecute executes the request
func (a *AlertConfigurationsApiService) DeleteAlertConfigExecute(r DeleteAlertConfigApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.DeleteAlertConfig")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.alertConfigId == "" {
		return nil, reportError("alertConfigId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"alertConfigId"+"}", url.PathEscape(r.alertConfigId), -1)

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

type GetAlertConfigApiRequest struct {
	ctx           context.Context
	ApiService    AlertConfigurationsApi
	groupId       string
	alertConfigId string
}

type GetAlertConfigApiParams struct {
	GroupId       string
	AlertConfigId string
}

func (a *AlertConfigurationsApiService) GetAlertConfigWithParams(ctx context.Context, args *GetAlertConfigApiParams) GetAlertConfigApiRequest {
	return GetAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		alertConfigId: args.AlertConfigId,
	}
}

func (r GetAlertConfigApiRequest) Execute() (*GroupAlertsConfig, *http.Response, error) {
	return r.ApiService.GetAlertConfigExecute(r)
}

/*
GetAlertConfig Return One Alert Configuration from One Project

Returns the specified alert configuration from the specified project. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
	@return GetAlertConfigApiRequest
*/
func (a *AlertConfigurationsApiService) GetAlertConfig(ctx context.Context, groupId string, alertConfigId string) GetAlertConfigApiRequest {
	return GetAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		alertConfigId: alertConfigId,
	}
}

// GetAlertConfigExecute executes the request
//
//	@return GroupAlertsConfig
func (a *AlertConfigurationsApiService) GetAlertConfigExecute(r GetAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupAlertsConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.GetAlertConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.alertConfigId == "" {
		return localVarReturnValue, nil, reportError("alertConfigId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"alertConfigId"+"}", url.PathEscape(r.alertConfigId), -1)

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

type GetAlertConfigsApiRequest struct {
	ctx          context.Context
	ApiService   AlertConfigurationsApi
	groupId      string
	alertId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type GetAlertConfigsApiParams struct {
	GroupId      string
	AlertId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *AlertConfigurationsApiService) GetAlertConfigsWithParams(ctx context.Context, args *GetAlertConfigsApiParams) GetAlertConfigsApiRequest {
	return GetAlertConfigsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		alertId:      args.AlertId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r GetAlertConfigsApiRequest) IncludeCount(includeCount bool) GetAlertConfigsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r GetAlertConfigsApiRequest) ItemsPerPage(itemsPerPage int) GetAlertConfigsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetAlertConfigsApiRequest) PageNum(pageNum int) GetAlertConfigsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r GetAlertConfigsApiRequest) Execute() (*PaginatedAlertConfig, *http.Response, error) {
	return r.ApiService.GetAlertConfigsExecute(r)
}

/*
GetAlertConfigs Return All Alert Configurations Set for One Alert

Returns all alert configurations set for the specified alert. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertId Unique 24-hexadecimal digit string that identifies the alert.
	@return GetAlertConfigsApiRequest
*/
func (a *AlertConfigurationsApiService) GetAlertConfigs(ctx context.Context, groupId string, alertId string) GetAlertConfigsApiRequest {
	return GetAlertConfigsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		alertId:    alertId,
	}
}

// GetAlertConfigsExecute executes the request
//
//	@return PaginatedAlertConfig
func (a *AlertConfigurationsApiService) GetAlertConfigsExecute(r GetAlertConfigsApiRequest) (*PaginatedAlertConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAlertConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.GetAlertConfigs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alerts/{alertId}/alertConfigs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.alertId == "" {
		return localVarReturnValue, nil, reportError("alertId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"alertId"+"}", url.PathEscape(r.alertId), -1)

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

type ListAlertConfigsApiRequest struct {
	ctx          context.Context
	ApiService   AlertConfigurationsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListAlertConfigsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *AlertConfigurationsApiService) ListAlertConfigsWithParams(ctx context.Context, args *ListAlertConfigsApiParams) ListAlertConfigsApiRequest {
	return ListAlertConfigsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListAlertConfigsApiRequest) IncludeCount(includeCount bool) ListAlertConfigsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListAlertConfigsApiRequest) ItemsPerPage(itemsPerPage int) ListAlertConfigsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListAlertConfigsApiRequest) PageNum(pageNum int) ListAlertConfigsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListAlertConfigsApiRequest) Execute() (*PaginatedAlertConfig, *http.Response, error) {
	return r.ApiService.ListAlertConfigsExecute(r)
}

/*
ListAlertConfigs Return All Alert Configurations in One Project

Returns all alert configurations for one project. These alert configurations apply to any component in the project. Alert configurations define the triggers and notification methods for alerts.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListAlertConfigsApiRequest
*/
func (a *AlertConfigurationsApiService) ListAlertConfigs(ctx context.Context, groupId string) ListAlertConfigsApiRequest {
	return ListAlertConfigsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListAlertConfigsExecute executes the request
//
//	@return PaginatedAlertConfig
func (a *AlertConfigurationsApiService) ListAlertConfigsExecute(r ListAlertConfigsApiRequest) (*PaginatedAlertConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAlertConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.ListAlertConfigs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs"
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

type ListMatcherFieldNamesApiRequest struct {
	ctx        context.Context
	ApiService AlertConfigurationsApi
}

type ListMatcherFieldNamesApiParams struct {
}

func (a *AlertConfigurationsApiService) ListMatcherFieldNamesWithParams(ctx context.Context, args *ListMatcherFieldNamesApiParams) ListMatcherFieldNamesApiRequest {
	return ListMatcherFieldNamesApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

func (r ListMatcherFieldNamesApiRequest) Execute() ([]string, *http.Response, error) {
	return r.ApiService.ListMatcherFieldNamesExecute(r)
}

/*
ListMatcherFieldNames Return All Alert Configuration Matchers Field Names

Get all field names that the `matchers.fieldName` parameter accepts when you create or update an Alert Configuration. You can successfully call this endpoint with any assigned role.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListMatcherFieldNamesApiRequest
*/
func (a *AlertConfigurationsApiService) ListMatcherFieldNames(ctx context.Context) ListMatcherFieldNamesApiRequest {
	return ListMatcherFieldNamesApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListMatcherFieldNamesExecute executes the request
//
//	@return []string
func (a *AlertConfigurationsApiService) ListMatcherFieldNamesExecute(r ListMatcherFieldNamesApiRequest) ([]string, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.ListMatcherFieldNames")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/alertConfigs/matchers/fieldNames"

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

type ToggleAlertConfigApiRequest struct {
	ctx           context.Context
	ApiService    AlertConfigurationsApi
	groupId       string
	alertConfigId string
	alertsToggle  *AlertsToggle
}

type ToggleAlertConfigApiParams struct {
	GroupId       string
	AlertConfigId string
	AlertsToggle  *AlertsToggle
}

func (a *AlertConfigurationsApiService) ToggleAlertConfigWithParams(ctx context.Context, args *ToggleAlertConfigApiParams) ToggleAlertConfigApiRequest {
	return ToggleAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		alertConfigId: args.AlertConfigId,
		alertsToggle:  args.AlertsToggle,
	}
}

func (r ToggleAlertConfigApiRequest) Execute() (*GroupAlertsConfig, *http.Response, error) {
	return r.ApiService.ToggleAlertConfigExecute(r)
}

/*
ToggleAlertConfig Toggle State of One Alert Configuration in One Project

Enables or disables the specified alert configuration in the specified project. The resource enables the specified alert configuration if currently enabled. The resource disables the specified alert configuration if currently disabled.

**NOTE**: This endpoint updates only the enabled/disabled state for the alert configuration. To update more than just this configuration, see Update One Alert Configuration.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration that triggered this alert.
	@return ToggleAlertConfigApiRequest
*/
func (a *AlertConfigurationsApiService) ToggleAlertConfig(ctx context.Context, groupId string, alertConfigId string, alertsToggle *AlertsToggle) ToggleAlertConfigApiRequest {
	return ToggleAlertConfigApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		alertConfigId: alertConfigId,
		alertsToggle:  alertsToggle,
	}
}

// ToggleAlertConfigExecute executes the request
//
//	@return GroupAlertsConfig
func (a *AlertConfigurationsApiService) ToggleAlertConfigExecute(r ToggleAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupAlertsConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.ToggleAlertConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.alertConfigId == "" {
		return localVarReturnValue, nil, reportError("alertConfigId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"alertConfigId"+"}", url.PathEscape(r.alertConfigId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.alertsToggle == nil {
		return localVarReturnValue, nil, reportError("alertsToggle is required and must be specified")
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
	localVarPostBody = r.alertsToggle
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

type UpdateAlertConfigApiRequest struct {
	ctx               context.Context
	ApiService        AlertConfigurationsApi
	groupId           string
	alertConfigId     string
	groupAlertsConfig *GroupAlertsConfig
}

type UpdateAlertConfigApiParams struct {
	GroupId           string
	AlertConfigId     string
	GroupAlertsConfig *GroupAlertsConfig
}

func (a *AlertConfigurationsApiService) UpdateAlertConfigWithParams(ctx context.Context, args *UpdateAlertConfigApiParams) UpdateAlertConfigApiRequest {
	return UpdateAlertConfigApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		alertConfigId:     args.AlertConfigId,
		groupAlertsConfig: args.GroupAlertsConfig,
	}
}

func (r UpdateAlertConfigApiRequest) Execute() (*GroupAlertsConfig, *http.Response, error) {
	return r.ApiService.UpdateAlertConfigExecute(r)
}

/*
UpdateAlertConfig Update One Alert Configuration in One Project

Updates one alert configuration in the specified project. Alert configurations define the triggers and notification methods for alerts.

**NOTE**: To enable or disable the alert configuration, see Toggle One State of One Alert Configuration in One Project.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
	@return UpdateAlertConfigApiRequest
*/
func (a *AlertConfigurationsApiService) UpdateAlertConfig(ctx context.Context, groupId string, alertConfigId string, groupAlertsConfig *GroupAlertsConfig) UpdateAlertConfigApiRequest {
	return UpdateAlertConfigApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		alertConfigId:     alertConfigId,
		groupAlertsConfig: groupAlertsConfig,
	}
}

// UpdateAlertConfigExecute executes the request
//
//	@return GroupAlertsConfig
func (a *AlertConfigurationsApiService) UpdateAlertConfigExecute(r UpdateAlertConfigApiRequest) (*GroupAlertsConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *GroupAlertsConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertConfigurationsApiService.UpdateAlertConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.alertConfigId == "" {
		return localVarReturnValue, nil, reportError("alertConfigId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"alertConfigId"+"}", url.PathEscape(r.alertConfigId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.groupAlertsConfig == nil {
		return localVarReturnValue, nil, reportError("groupAlertsConfig is required and must be specified")
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
	localVarPostBody = r.groupAlertsConfig
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
