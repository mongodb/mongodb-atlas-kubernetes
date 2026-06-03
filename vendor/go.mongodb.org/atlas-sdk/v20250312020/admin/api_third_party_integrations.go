// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ThirdPartyIntegrationsApi interface {

	/*
		CreateGroupIntegration Create One Third-Party Service Integration

		Adds the settings for configuring one third-party service integration. These settings apply to all databases managed in the specified MongoDB Cloud project. Each project can have only one configuration per `{INTEGRATION-TYPE}`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param thirdPartyIntegration Third-party integration that you want to configure for your project.
		@return CreateGroupIntegrationApiRequest
	*/
	CreateGroupIntegration(ctx context.Context, integrationType string, groupId string, thirdPartyIntegration *ThirdPartyIntegration) CreateGroupIntegrationApiRequest
	/*
		CreateGroupIntegration Create One Third-Party Service Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupIntegrationApiParams - Parameters for the request
		@return CreateGroupIntegrationApiRequest
	*/
	CreateGroupIntegrationWithParams(ctx context.Context, args *CreateGroupIntegrationApiParams) CreateGroupIntegrationApiRequest

	// Method available only for mocking purposes
	CreateGroupIntegrationExecute(r CreateGroupIntegrationApiRequest) (*PaginatedIntegration, *http.Response, error)

	/*
		DeleteGroupIntegration Remove One Third-Party Service Integration

		Removes the settings that permit configuring one third-party service integration. These settings apply to all databases managed in one MongoDB Cloud project. If you delete an integration from a project, you remove that integration configuration only for that project. This action doesn't affect any other project or organization's configured `{INTEGRATION-TYPE}` integrations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DeleteGroupIntegrationApiRequest
	*/
	DeleteGroupIntegration(ctx context.Context, integrationType string, groupId string) DeleteGroupIntegrationApiRequest
	/*
		DeleteGroupIntegration Remove One Third-Party Service Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupIntegrationApiParams - Parameters for the request
		@return DeleteGroupIntegrationApiRequest
	*/
	DeleteGroupIntegrationWithParams(ctx context.Context, args *DeleteGroupIntegrationApiParams) DeleteGroupIntegrationApiRequest

	// Method available only for mocking purposes
	DeleteGroupIntegrationExecute(r DeleteGroupIntegrationApiRequest) (*http.Response, error)

	/*
		GetGroupIntegration Return One Third-Party Service Integration

		Returns the settings for configuring integration with one third-party service. These settings apply to all databases managed in one MongoDB Cloud project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
		@return GetGroupIntegrationApiRequest
	*/
	GetGroupIntegration(ctx context.Context, groupId string, integrationType string) GetGroupIntegrationApiRequest
	/*
		GetGroupIntegration Return One Third-Party Service Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupIntegrationApiParams - Parameters for the request
		@return GetGroupIntegrationApiRequest
	*/
	GetGroupIntegrationWithParams(ctx context.Context, args *GetGroupIntegrationApiParams) GetGroupIntegrationApiRequest

	// Method available only for mocking purposes
	GetGroupIntegrationExecute(r GetGroupIntegrationApiRequest) (*ThirdPartyIntegration, *http.Response, error)

	/*
		ListGroupIntegrations Return All Active Third-Party Service Integrations

		Returns the settings that permit integrations with all configured third-party services. These settings apply to all databases managed in one MongoDB Cloud project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupIntegrationsApiRequest
	*/
	ListGroupIntegrations(ctx context.Context, groupId string) ListGroupIntegrationsApiRequest
	/*
		ListGroupIntegrations Return All Active Third-Party Service Integrations


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupIntegrationsApiParams - Parameters for the request
		@return ListGroupIntegrationsApiRequest
	*/
	ListGroupIntegrationsWithParams(ctx context.Context, args *ListGroupIntegrationsApiParams) ListGroupIntegrationsApiRequest

	// Method available only for mocking purposes
	ListGroupIntegrationsExecute(r ListGroupIntegrationsApiRequest) (*PaginatedIntegration, *http.Response, error)

	/*
		UpdateGroupIntegration Update One Third-Party Service Integration

		Updates the settings for configuring integration with one third-party service. These settings apply to all databases managed in one MongoDB Cloud project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param thirdPartyIntegration Third-party integration that you want to configure for your project.
		@return UpdateGroupIntegrationApiRequest
	*/
	UpdateGroupIntegration(ctx context.Context, integrationType string, groupId string, thirdPartyIntegration *ThirdPartyIntegration) UpdateGroupIntegrationApiRequest
	/*
		UpdateGroupIntegration Update One Third-Party Service Integration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupIntegrationApiParams - Parameters for the request
		@return UpdateGroupIntegrationApiRequest
	*/
	UpdateGroupIntegrationWithParams(ctx context.Context, args *UpdateGroupIntegrationApiParams) UpdateGroupIntegrationApiRequest

	// Method available only for mocking purposes
	UpdateGroupIntegrationExecute(r UpdateGroupIntegrationApiRequest) (*PaginatedIntegration, *http.Response, error)
}

// ThirdPartyIntegrationsApiService ThirdPartyIntegrationsApi service
type ThirdPartyIntegrationsApiService service

type CreateGroupIntegrationApiRequest struct {
	ctx                   context.Context
	ApiService            ThirdPartyIntegrationsApi
	integrationType       string
	groupId               string
	thirdPartyIntegration *ThirdPartyIntegration
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
}

type CreateGroupIntegrationApiParams struct {
	IntegrationType       string
	GroupId               string
	ThirdPartyIntegration *ThirdPartyIntegration
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
}

func (a *ThirdPartyIntegrationsApiService) CreateGroupIntegrationWithParams(ctx context.Context, args *CreateGroupIntegrationApiParams) CreateGroupIntegrationApiRequest {
	return CreateGroupIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		integrationType:       args.IntegrationType,
		groupId:               args.GroupId,
		thirdPartyIntegration: args.ThirdPartyIntegration,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r CreateGroupIntegrationApiRequest) IncludeCount(includeCount bool) CreateGroupIntegrationApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r CreateGroupIntegrationApiRequest) ItemsPerPage(itemsPerPage int) CreateGroupIntegrationApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r CreateGroupIntegrationApiRequest) PageNum(pageNum int) CreateGroupIntegrationApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r CreateGroupIntegrationApiRequest) Execute() (*PaginatedIntegration, *http.Response, error) {
	return r.ApiService.CreateGroupIntegrationExecute(r)
}

/*
CreateGroupIntegration Create One Third-Party Service Integration

Adds the settings for configuring one third-party service integration. These settings apply to all databases managed in the specified MongoDB Cloud project. Each project can have only one configuration per `{INTEGRATION-TYPE}`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupIntegrationApiRequest
*/
func (a *ThirdPartyIntegrationsApiService) CreateGroupIntegration(ctx context.Context, integrationType string, groupId string, thirdPartyIntegration *ThirdPartyIntegration) CreateGroupIntegrationApiRequest {
	return CreateGroupIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		integrationType:       integrationType,
		groupId:               groupId,
		thirdPartyIntegration: thirdPartyIntegration,
	}
}

// CreateGroupIntegrationExecute executes the request
//
//	@return PaginatedIntegration
func (a *ThirdPartyIntegrationsApiService) CreateGroupIntegrationExecute(r CreateGroupIntegrationApiRequest) (*PaginatedIntegration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedIntegration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ThirdPartyIntegrationsApiService.CreateGroupIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/integrations/{integrationType}"
	if r.integrationType == "" {
		return localVarReturnValue, nil, reportError("integrationType is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"integrationType"+"}", url.PathEscape(r.integrationType), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.thirdPartyIntegration == nil {
		return localVarReturnValue, nil, reportError("thirdPartyIntegration is required and must be specified")
	}

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
	localVarPostBody = r.thirdPartyIntegration
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

type DeleteGroupIntegrationApiRequest struct {
	ctx             context.Context
	ApiService      ThirdPartyIntegrationsApi
	integrationType string
	groupId         string
}

type DeleteGroupIntegrationApiParams struct {
	IntegrationType string
	GroupId         string
}

func (a *ThirdPartyIntegrationsApiService) DeleteGroupIntegrationWithParams(ctx context.Context, args *DeleteGroupIntegrationApiParams) DeleteGroupIntegrationApiRequest {
	return DeleteGroupIntegrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		integrationType: args.IntegrationType,
		groupId:         args.GroupId,
	}
}

func (r DeleteGroupIntegrationApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupIntegrationExecute(r)
}

/*
DeleteGroupIntegration Remove One Third-Party Service Integration

Removes the settings that permit configuring one third-party service integration. These settings apply to all databases managed in one MongoDB Cloud project. If you delete an integration from a project, you remove that integration configuration only for that project. This action doesn't affect any other project or organization's configured `{INTEGRATION-TYPE}` integrations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DeleteGroupIntegrationApiRequest
*/
func (a *ThirdPartyIntegrationsApiService) DeleteGroupIntegration(ctx context.Context, integrationType string, groupId string) DeleteGroupIntegrationApiRequest {
	return DeleteGroupIntegrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		integrationType: integrationType,
		groupId:         groupId,
	}
}

// DeleteGroupIntegrationExecute executes the request
func (a *ThirdPartyIntegrationsApiService) DeleteGroupIntegrationExecute(r DeleteGroupIntegrationApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ThirdPartyIntegrationsApiService.DeleteGroupIntegration")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/integrations/{integrationType}"
	if r.integrationType == "" {
		return nil, reportError("integrationType is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"integrationType"+"}", url.PathEscape(r.integrationType), -1)
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

type GetGroupIntegrationApiRequest struct {
	ctx             context.Context
	ApiService      ThirdPartyIntegrationsApi
	groupId         string
	integrationType string
}

type GetGroupIntegrationApiParams struct {
	GroupId         string
	IntegrationType string
}

func (a *ThirdPartyIntegrationsApiService) GetGroupIntegrationWithParams(ctx context.Context, args *GetGroupIntegrationApiParams) GetGroupIntegrationApiRequest {
	return GetGroupIntegrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         args.GroupId,
		integrationType: args.IntegrationType,
	}
}

func (r GetGroupIntegrationApiRequest) Execute() (*ThirdPartyIntegration, *http.Response, error) {
	return r.ApiService.GetGroupIntegrationExecute(r)
}

/*
GetGroupIntegration Return One Third-Party Service Integration

Returns the settings for configuring integration with one third-party service. These settings apply to all databases managed in one MongoDB Cloud project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
	@return GetGroupIntegrationApiRequest
*/
func (a *ThirdPartyIntegrationsApiService) GetGroupIntegration(ctx context.Context, groupId string, integrationType string) GetGroupIntegrationApiRequest {
	return GetGroupIntegrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         groupId,
		integrationType: integrationType,
	}
}

// GetGroupIntegrationExecute executes the request
//
//	@return ThirdPartyIntegration
func (a *ThirdPartyIntegrationsApiService) GetGroupIntegrationExecute(r GetGroupIntegrationApiRequest) (*ThirdPartyIntegration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ThirdPartyIntegration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ThirdPartyIntegrationsApiService.GetGroupIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/integrations/{integrationType}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.integrationType == "" {
		return localVarReturnValue, nil, reportError("integrationType is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"integrationType"+"}", url.PathEscape(r.integrationType), -1)

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

type ListGroupIntegrationsApiRequest struct {
	ctx          context.Context
	ApiService   ThirdPartyIntegrationsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupIntegrationsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ThirdPartyIntegrationsApiService) ListGroupIntegrationsWithParams(ctx context.Context, args *ListGroupIntegrationsApiParams) ListGroupIntegrationsApiRequest {
	return ListGroupIntegrationsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupIntegrationsApiRequest) IncludeCount(includeCount bool) ListGroupIntegrationsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupIntegrationsApiRequest) ItemsPerPage(itemsPerPage int) ListGroupIntegrationsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupIntegrationsApiRequest) PageNum(pageNum int) ListGroupIntegrationsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupIntegrationsApiRequest) Execute() (*PaginatedIntegration, *http.Response, error) {
	return r.ApiService.ListGroupIntegrationsExecute(r)
}

/*
ListGroupIntegrations Return All Active Third-Party Service Integrations

Returns the settings that permit integrations with all configured third-party services. These settings apply to all databases managed in one MongoDB Cloud project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupIntegrationsApiRequest
*/
func (a *ThirdPartyIntegrationsApiService) ListGroupIntegrations(ctx context.Context, groupId string) ListGroupIntegrationsApiRequest {
	return ListGroupIntegrationsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupIntegrationsExecute executes the request
//
//	@return PaginatedIntegration
func (a *ThirdPartyIntegrationsApiService) ListGroupIntegrationsExecute(r ListGroupIntegrationsApiRequest) (*PaginatedIntegration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedIntegration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ThirdPartyIntegrationsApiService.ListGroupIntegrations")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/integrations"
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

type UpdateGroupIntegrationApiRequest struct {
	ctx                   context.Context
	ApiService            ThirdPartyIntegrationsApi
	integrationType       string
	groupId               string
	thirdPartyIntegration *ThirdPartyIntegration
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
}

type UpdateGroupIntegrationApiParams struct {
	IntegrationType       string
	GroupId               string
	ThirdPartyIntegration *ThirdPartyIntegration
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
}

func (a *ThirdPartyIntegrationsApiService) UpdateGroupIntegrationWithParams(ctx context.Context, args *UpdateGroupIntegrationApiParams) UpdateGroupIntegrationApiRequest {
	return UpdateGroupIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		integrationType:       args.IntegrationType,
		groupId:               args.GroupId,
		thirdPartyIntegration: args.ThirdPartyIntegration,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r UpdateGroupIntegrationApiRequest) IncludeCount(includeCount bool) UpdateGroupIntegrationApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r UpdateGroupIntegrationApiRequest) ItemsPerPage(itemsPerPage int) UpdateGroupIntegrationApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r UpdateGroupIntegrationApiRequest) PageNum(pageNum int) UpdateGroupIntegrationApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r UpdateGroupIntegrationApiRequest) Execute() (*PaginatedIntegration, *http.Response, error) {
	return r.ApiService.UpdateGroupIntegrationExecute(r)
}

/*
UpdateGroupIntegration Update One Third-Party Service Integration

Updates the settings for configuring integration with one third-party service. These settings apply to all databases managed in one MongoDB Cloud project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param integrationType Human-readable label that identifies the service which you want to integrate with MongoDB Cloud.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateGroupIntegrationApiRequest
*/
func (a *ThirdPartyIntegrationsApiService) UpdateGroupIntegration(ctx context.Context, integrationType string, groupId string, thirdPartyIntegration *ThirdPartyIntegration) UpdateGroupIntegrationApiRequest {
	return UpdateGroupIntegrationApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		integrationType:       integrationType,
		groupId:               groupId,
		thirdPartyIntegration: thirdPartyIntegration,
	}
}

// UpdateGroupIntegrationExecute executes the request
//
//	@return PaginatedIntegration
func (a *ThirdPartyIntegrationsApiService) UpdateGroupIntegrationExecute(r UpdateGroupIntegrationApiRequest) (*PaginatedIntegration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedIntegration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ThirdPartyIntegrationsApiService.UpdateGroupIntegration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/integrations/{integrationType}"
	if r.integrationType == "" {
		return localVarReturnValue, nil, reportError("integrationType is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"integrationType"+"}", url.PathEscape(r.integrationType), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.thirdPartyIntegration == nil {
		return localVarReturnValue, nil, reportError("thirdPartyIntegration is required and must be specified")
	}

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
	localVarPostBody = r.thirdPartyIntegration
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
