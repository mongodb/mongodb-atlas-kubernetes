// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AlertsApi interface {

	/*
			AcknowledgeAlert Acknowledge One Alert from One Project

			Confirms receipt of one existing alert. This alert applies to any component in one project. Acknowledging an alert prevents successive notifications. You receive an alert when a monitored component meets or exceeds a value you set until you acknowledge the alert. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

		This resource remains under revision and may change. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertId Unique 24-hexadecimal digit string that identifies the alert.
			@param acknowledgeAlert Acknowledges or unacknowledges one alert.
			@return AcknowledgeAlertApiRequest
	*/
	AcknowledgeAlert(ctx context.Context, groupId string, alertId string, acknowledgeAlert *AcknowledgeAlert) AcknowledgeAlertApiRequest
	/*
		AcknowledgeAlert Acknowledge One Alert from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AcknowledgeAlertApiParams - Parameters for the request
		@return AcknowledgeAlertApiRequest
	*/
	AcknowledgeAlertWithParams(ctx context.Context, args *AcknowledgeAlertApiParams) AcknowledgeAlertApiRequest

	// Method available only for mocking purposes
	AcknowledgeAlertExecute(r AcknowledgeAlertApiRequest) (*AlertViewForNdsGroup, *http.Response, error)

	/*
			GetAlert Return One Alert from One Project

			Returns one alert. This alert applies to any component in one project. You receive an alert when a monitored component meets or exceeds a value you set. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertId Unique 24-hexadecimal digit string that identifies the alert.
			@return GetAlertApiRequest
	*/
	GetAlert(ctx context.Context, groupId string, alertId string) GetAlertApiRequest
	/*
		GetAlert Return One Alert from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAlertApiParams - Parameters for the request
		@return GetAlertApiRequest
	*/
	GetAlertWithParams(ctx context.Context, args *GetAlertApiParams) GetAlertApiRequest

	// Method available only for mocking purposes
	GetAlertExecute(r GetAlertApiRequest) (*AlertViewForNdsGroup, *http.Response, error)

	/*
			GetAlertConfigAlerts Return All Open Alerts for One Alert Configuration

			Returns all open alerts that the specified alert configuration triggers. These alert configurations apply to the specified project only. Alert configurations define the triggers and notification methods for alerts. Open alerts have been triggered but remain unacknowledged. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
			@return GetAlertConfigAlertsApiRequest
	*/
	GetAlertConfigAlerts(ctx context.Context, groupId string, alertConfigId string) GetAlertConfigAlertsApiRequest
	/*
		GetAlertConfigAlerts Return All Open Alerts for One Alert Configuration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAlertConfigAlertsApiParams - Parameters for the request
		@return GetAlertConfigAlertsApiRequest
	*/
	GetAlertConfigAlertsWithParams(ctx context.Context, args *GetAlertConfigAlertsApiParams) GetAlertConfigAlertsApiRequest

	// Method available only for mocking purposes
	GetAlertConfigAlertsExecute(r GetAlertConfigAlertsApiRequest) (*PaginatedAlert, *http.Response, error)

	/*
			ListAlerts Return All Alerts from One Project

			Returns all alerts. These alerts apply to all components in one project. You receive an alert when a monitored component meets or exceeds a value you set.

		This resource remains under revision and may change.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListAlertsApiRequest
	*/
	ListAlerts(ctx context.Context, groupId string) ListAlertsApiRequest
	/*
		ListAlerts Return All Alerts from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListAlertsApiParams - Parameters for the request
		@return ListAlertsApiRequest
	*/
	ListAlertsWithParams(ctx context.Context, args *ListAlertsApiParams) ListAlertsApiRequest

	// Method available only for mocking purposes
	ListAlertsExecute(r ListAlertsApiRequest) (*PaginatedAlert, *http.Response, error)
}

// AlertsApiService AlertsApi service
type AlertsApiService service

type AcknowledgeAlertApiRequest struct {
	ctx              context.Context
	ApiService       AlertsApi
	groupId          string
	alertId          string
	acknowledgeAlert *AcknowledgeAlert
}

type AcknowledgeAlertApiParams struct {
	GroupId          string
	AlertId          string
	AcknowledgeAlert *AcknowledgeAlert
}

func (a *AlertsApiService) AcknowledgeAlertWithParams(ctx context.Context, args *AcknowledgeAlertApiParams) AcknowledgeAlertApiRequest {
	return AcknowledgeAlertApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		alertId:          args.AlertId,
		acknowledgeAlert: args.AcknowledgeAlert,
	}
}

func (r AcknowledgeAlertApiRequest) Execute() (*AlertViewForNdsGroup, *http.Response, error) {
	return r.ApiService.AcknowledgeAlertExecute(r)
}

/*
AcknowledgeAlert Acknowledge One Alert from One Project

Confirms receipt of one existing alert. This alert applies to any component in one project. Acknowledging an alert prevents successive notifications. You receive an alert when a monitored component meets or exceeds a value you set until you acknowledge the alert. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

This resource remains under revision and may change. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertId Unique 24-hexadecimal digit string that identifies the alert.
	@return AcknowledgeAlertApiRequest
*/
func (a *AlertsApiService) AcknowledgeAlert(ctx context.Context, groupId string, alertId string, acknowledgeAlert *AcknowledgeAlert) AcknowledgeAlertApiRequest {
	return AcknowledgeAlertApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		alertId:          alertId,
		acknowledgeAlert: acknowledgeAlert,
	}
}

// AcknowledgeAlertExecute executes the request
//
//	@return AlertViewForNdsGroup
func (a *AlertsApiService) AcknowledgeAlertExecute(r AcknowledgeAlertApiRequest) (*AlertViewForNdsGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AlertViewForNdsGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertsApiService.AcknowledgeAlert")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alerts/{alertId}"
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
	if r.acknowledgeAlert == nil {
		return localVarReturnValue, nil, reportError("acknowledgeAlert is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-05-30+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-05-30+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.acknowledgeAlert
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

type GetAlertApiRequest struct {
	ctx        context.Context
	ApiService AlertsApi
	groupId    string
	alertId    string
}

type GetAlertApiParams struct {
	GroupId string
	AlertId string
}

func (a *AlertsApiService) GetAlertWithParams(ctx context.Context, args *GetAlertApiParams) GetAlertApiRequest {
	return GetAlertApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		alertId:    args.AlertId,
	}
}

func (r GetAlertApiRequest) Execute() (*AlertViewForNdsGroup, *http.Response, error) {
	return r.ApiService.GetAlertExecute(r)
}

/*
GetAlert Return One Alert from One Project

Returns one alert. This alert applies to any component in one project. You receive an alert when a monitored component meets or exceeds a value you set. Use the Return All Alerts from One Project endpoint to retrieve all alerts to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertId Unique 24-hexadecimal digit string that identifies the alert.
	@return GetAlertApiRequest
*/
func (a *AlertsApiService) GetAlert(ctx context.Context, groupId string, alertId string) GetAlertApiRequest {
	return GetAlertApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		alertId:    alertId,
	}
}

// GetAlertExecute executes the request
//
//	@return AlertViewForNdsGroup
func (a *AlertsApiService) GetAlertExecute(r GetAlertApiRequest) (*AlertViewForNdsGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AlertViewForNdsGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertsApiService.GetAlert")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alerts/{alertId}"
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

type GetAlertConfigAlertsApiRequest struct {
	ctx           context.Context
	ApiService    AlertsApi
	groupId       string
	alertConfigId string
	includeCount  *bool
	itemsPerPage  *int
	pageNum       *int
}

type GetAlertConfigAlertsApiParams struct {
	GroupId       string
	AlertConfigId string
	IncludeCount  *bool
	ItemsPerPage  *int
	PageNum       *int
}

func (a *AlertsApiService) GetAlertConfigAlertsWithParams(ctx context.Context, args *GetAlertConfigAlertsApiParams) GetAlertConfigAlertsApiRequest {
	return GetAlertConfigAlertsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		alertConfigId: args.AlertConfigId,
		includeCount:  args.IncludeCount,
		itemsPerPage:  args.ItemsPerPage,
		pageNum:       args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r GetAlertConfigAlertsApiRequest) IncludeCount(includeCount bool) GetAlertConfigAlertsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r GetAlertConfigAlertsApiRequest) ItemsPerPage(itemsPerPage int) GetAlertConfigAlertsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetAlertConfigAlertsApiRequest) PageNum(pageNum int) GetAlertConfigAlertsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r GetAlertConfigAlertsApiRequest) Execute() (*PaginatedAlert, *http.Response, error) {
	return r.ApiService.GetAlertConfigAlertsExecute(r)
}

/*
GetAlertConfigAlerts Return All Open Alerts for One Alert Configuration

Returns all open alerts that the specified alert configuration triggers. These alert configurations apply to the specified project only. Alert configurations define the triggers and notification methods for alerts. Open alerts have been triggered but remain unacknowledged. Use the Return All Alert Configurations for One Project endpoint to retrieve all alert configurations to which the authenticated user has access.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param alertConfigId Unique 24-hexadecimal digit string that identifies the alert configuration.
	@return GetAlertConfigAlertsApiRequest
*/
func (a *AlertsApiService) GetAlertConfigAlerts(ctx context.Context, groupId string, alertConfigId string) GetAlertConfigAlertsApiRequest {
	return GetAlertConfigAlertsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		alertConfigId: alertConfigId,
	}
}

// GetAlertConfigAlertsExecute executes the request
//
//	@return PaginatedAlert
func (a *AlertsApiService) GetAlertConfigAlertsExecute(r GetAlertConfigAlertsApiRequest) (*PaginatedAlert, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAlert
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertsApiService.GetAlertConfigAlerts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}/alerts"
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

type ListAlertsApiRequest struct {
	ctx          context.Context
	ApiService   AlertsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	status       *string
}

type ListAlertsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Status       *string
}

func (a *AlertsApiService) ListAlertsWithParams(ctx context.Context, args *ListAlertsApiParams) ListAlertsApiRequest {
	return ListAlertsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		status:       args.Status,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListAlertsApiRequest) IncludeCount(includeCount bool) ListAlertsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListAlertsApiRequest) ItemsPerPage(itemsPerPage int) ListAlertsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListAlertsApiRequest) PageNum(pageNum int) ListAlertsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Status of the alerts to return. Omit this parameter to return all alerts in all statuses. TRACKING indicates the alert condition exists but has not persisted for the minimum notification delay. OPEN indicates the alert condition currently exists. CLOSED indicates the alert condition has been resolved.
func (r ListAlertsApiRequest) Status(status string) ListAlertsApiRequest {
	r.status = &status
	return r
}

func (r ListAlertsApiRequest) Execute() (*PaginatedAlert, *http.Response, error) {
	return r.ApiService.ListAlertsExecute(r)
}

/*
ListAlerts Return All Alerts from One Project

Returns all alerts. These alerts apply to all components in one project. You receive an alert when a monitored component meets or exceeds a value you set.

This resource remains under revision and may change.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListAlertsApiRequest
*/
func (a *AlertsApiService) ListAlerts(ctx context.Context, groupId string) ListAlertsApiRequest {
	return ListAlertsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListAlertsExecute executes the request
//
//	@return PaginatedAlert
func (a *AlertsApiService) ListAlertsExecute(r ListAlertsApiRequest) (*PaginatedAlert, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAlert
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AlertsApiService.ListAlerts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/alerts"
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
	if r.status != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "status", r.status, "")
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
