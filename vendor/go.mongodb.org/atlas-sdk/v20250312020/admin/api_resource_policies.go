// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ResourcePoliciesApi interface {

	/*
		CreateOrgResourcePolicy Create One Atlas Resource Policy

		Create one Atlas Resource Policy for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiAtlasResourcePolicyCreate Atlas Resource Policy to create.
		@return CreateOrgResourcePolicyApiRequest
	*/
	CreateOrgResourcePolicy(ctx context.Context, orgId string, apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate) CreateOrgResourcePolicyApiRequest
	/*
		CreateOrgResourcePolicy Create One Atlas Resource Policy


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgResourcePolicyApiParams - Parameters for the request
		@return CreateOrgResourcePolicyApiRequest
	*/
	CreateOrgResourcePolicyWithParams(ctx context.Context, args *CreateOrgResourcePolicyApiParams) CreateOrgResourcePolicyApiRequest

	// Method available only for mocking purposes
	CreateOrgResourcePolicyExecute(r CreateOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error)

	/*
		DeleteOrgResourcePolicy Delete One Atlas Resource Policy

		Delete one Atlas Resource Policy for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
		@return DeleteOrgResourcePolicyApiRequest
	*/
	DeleteOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string) DeleteOrgResourcePolicyApiRequest
	/*
		DeleteOrgResourcePolicy Delete One Atlas Resource Policy


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgResourcePolicyApiParams - Parameters for the request
		@return DeleteOrgResourcePolicyApiRequest
	*/
	DeleteOrgResourcePolicyWithParams(ctx context.Context, args *DeleteOrgResourcePolicyApiParams) DeleteOrgResourcePolicyApiRequest

	// Method available only for mocking purposes
	DeleteOrgResourcePolicyExecute(r DeleteOrgResourcePolicyApiRequest) (*http.Response, error)

	/*
		GetNonCompliantResources Return All Non-Compliant Resources

		Return all non-compliant resources for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return GetNonCompliantResourcesApiRequest
	*/
	GetNonCompliantResources(ctx context.Context, orgId string) GetNonCompliantResourcesApiRequest
	/*
		GetNonCompliantResources Return All Non-Compliant Resources


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetNonCompliantResourcesApiParams - Parameters for the request
		@return GetNonCompliantResourcesApiRequest
	*/
	GetNonCompliantResourcesWithParams(ctx context.Context, args *GetNonCompliantResourcesApiParams) GetNonCompliantResourcesApiRequest

	// Method available only for mocking purposes
	GetNonCompliantResourcesExecute(r GetNonCompliantResourcesApiRequest) ([]ApiAtlasNonCompliantResource, *http.Response, error)

	/*
		GetOrgResourcePolicy Return One Atlas Resource Policy

		Return one Atlas Resource Policy for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
		@return GetOrgResourcePolicyApiRequest
	*/
	GetOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string) GetOrgResourcePolicyApiRequest
	/*
		GetOrgResourcePolicy Return One Atlas Resource Policy


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgResourcePolicyApiParams - Parameters for the request
		@return GetOrgResourcePolicyApiRequest
	*/
	GetOrgResourcePolicyWithParams(ctx context.Context, args *GetOrgResourcePolicyApiParams) GetOrgResourcePolicyApiRequest

	// Method available only for mocking purposes
	GetOrgResourcePolicyExecute(r GetOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error)

	/*
		ListOrgResourcePolicies Return All Atlas Resource Policies

		Return all Atlas Resource Policies for the organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListOrgResourcePoliciesApiRequest
	*/
	ListOrgResourcePolicies(ctx context.Context, orgId string) ListOrgResourcePoliciesApiRequest
	/*
		ListOrgResourcePolicies Return All Atlas Resource Policies


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgResourcePoliciesApiParams - Parameters for the request
		@return ListOrgResourcePoliciesApiRequest
	*/
	ListOrgResourcePoliciesWithParams(ctx context.Context, args *ListOrgResourcePoliciesApiParams) ListOrgResourcePoliciesApiRequest

	// Method available only for mocking purposes
	ListOrgResourcePoliciesExecute(r ListOrgResourcePoliciesApiRequest) ([]ApiAtlasResourcePolicy, *http.Response, error)

	/*
		UpdateOrgResourcePolicy Update One Atlas Resource Policy

		Update one Atlas Resource Policy for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
		@param apiAtlasResourcePolicyEdit Atlas Resource Policy to update.
		@return UpdateOrgResourcePolicyApiRequest
	*/
	UpdateOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string, apiAtlasResourcePolicyEdit *ApiAtlasResourcePolicyEdit) UpdateOrgResourcePolicyApiRequest
	/*
		UpdateOrgResourcePolicy Update One Atlas Resource Policy


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgResourcePolicyApiParams - Parameters for the request
		@return UpdateOrgResourcePolicyApiRequest
	*/
	UpdateOrgResourcePolicyWithParams(ctx context.Context, args *UpdateOrgResourcePolicyApiParams) UpdateOrgResourcePolicyApiRequest

	// Method available only for mocking purposes
	UpdateOrgResourcePolicyExecute(r UpdateOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error)

	/*
		ValidateResourcePolicies Validate One Atlas Resource Policy

		Validate one Atlas Resource Policy for an organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param apiAtlasResourcePolicyCreate Atlas Resource Policy to create.
		@return ValidateResourcePoliciesApiRequest
	*/
	ValidateResourcePolicies(ctx context.Context, orgId string, apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate) ValidateResourcePoliciesApiRequest
	/*
		ValidateResourcePolicies Validate One Atlas Resource Policy


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ValidateResourcePoliciesApiParams - Parameters for the request
		@return ValidateResourcePoliciesApiRequest
	*/
	ValidateResourcePoliciesWithParams(ctx context.Context, args *ValidateResourcePoliciesApiParams) ValidateResourcePoliciesApiRequest

	// Method available only for mocking purposes
	ValidateResourcePoliciesExecute(r ValidateResourcePoliciesApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error)
}

// ResourcePoliciesApiService ResourcePoliciesApi service
type ResourcePoliciesApiService service

type CreateOrgResourcePolicyApiRequest struct {
	ctx                          context.Context
	ApiService                   ResourcePoliciesApi
	orgId                        string
	apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate
}

type CreateOrgResourcePolicyApiParams struct {
	OrgId                        string
	ApiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate
}

func (a *ResourcePoliciesApiService) CreateOrgResourcePolicyWithParams(ctx context.Context, args *CreateOrgResourcePolicyApiParams) CreateOrgResourcePolicyApiRequest {
	return CreateOrgResourcePolicyApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		orgId:                        args.OrgId,
		apiAtlasResourcePolicyCreate: args.ApiAtlasResourcePolicyCreate,
	}
}

func (r CreateOrgResourcePolicyApiRequest) Execute() (*ApiAtlasResourcePolicy, *http.Response, error) {
	return r.ApiService.CreateOrgResourcePolicyExecute(r)
}

/*
CreateOrgResourcePolicy Create One Atlas Resource Policy

Create one Atlas Resource Policy for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgResourcePolicyApiRequest
*/
func (a *ResourcePoliciesApiService) CreateOrgResourcePolicy(ctx context.Context, orgId string, apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate) CreateOrgResourcePolicyApiRequest {
	return CreateOrgResourcePolicyApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		orgId:                        orgId,
		apiAtlasResourcePolicyCreate: apiAtlasResourcePolicyCreate,
	}
}

// CreateOrgResourcePolicyExecute executes the request
//
//	@return ApiAtlasResourcePolicy
func (a *ResourcePoliciesApiService) CreateOrgResourcePolicyExecute(r CreateOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasResourcePolicy
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.CreateOrgResourcePolicy")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.apiAtlasResourcePolicyCreate == nil {
		return localVarReturnValue, nil, reportError("apiAtlasResourcePolicyCreate is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.apiAtlasResourcePolicyCreate
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

type DeleteOrgResourcePolicyApiRequest struct {
	ctx              context.Context
	ApiService       ResourcePoliciesApi
	orgId            string
	resourcePolicyId string
}

type DeleteOrgResourcePolicyApiParams struct {
	OrgId            string
	ResourcePolicyId string
}

func (a *ResourcePoliciesApiService) DeleteOrgResourcePolicyWithParams(ctx context.Context, args *DeleteOrgResourcePolicyApiParams) DeleteOrgResourcePolicyApiRequest {
	return DeleteOrgResourcePolicyApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            args.OrgId,
		resourcePolicyId: args.ResourcePolicyId,
	}
}

func (r DeleteOrgResourcePolicyApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgResourcePolicyExecute(r)
}

/*
DeleteOrgResourcePolicy Delete One Atlas Resource Policy

Delete one Atlas Resource Policy for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
	@return DeleteOrgResourcePolicyApiRequest
*/
func (a *ResourcePoliciesApiService) DeleteOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string) DeleteOrgResourcePolicyApiRequest {
	return DeleteOrgResourcePolicyApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            orgId,
		resourcePolicyId: resourcePolicyId,
	}
}

// DeleteOrgResourcePolicyExecute executes the request
func (a *ResourcePoliciesApiService) DeleteOrgResourcePolicyExecute(r DeleteOrgResourcePolicyApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.DeleteOrgResourcePolicy")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies/{resourcePolicyId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.resourcePolicyId == "" {
		return nil, reportError("resourcePolicyId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"resourcePolicyId"+"}", url.PathEscape(r.resourcePolicyId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetNonCompliantResourcesApiRequest struct {
	ctx        context.Context
	ApiService ResourcePoliciesApi
	orgId      string
}

type GetNonCompliantResourcesApiParams struct {
	OrgId string
}

func (a *ResourcePoliciesApiService) GetNonCompliantResourcesWithParams(ctx context.Context, args *GetNonCompliantResourcesApiParams) GetNonCompliantResourcesApiRequest {
	return GetNonCompliantResourcesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r GetNonCompliantResourcesApiRequest) Execute() ([]ApiAtlasNonCompliantResource, *http.Response, error) {
	return r.ApiService.GetNonCompliantResourcesExecute(r)
}

/*
GetNonCompliantResources Return All Non-Compliant Resources

Return all non-compliant resources for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return GetNonCompliantResourcesApiRequest
*/
func (a *ResourcePoliciesApiService) GetNonCompliantResources(ctx context.Context, orgId string) GetNonCompliantResourcesApiRequest {
	return GetNonCompliantResourcesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// GetNonCompliantResourcesExecute executes the request
//
//	@return []ApiAtlasNonCompliantResource
func (a *ResourcePoliciesApiService) GetNonCompliantResourcesExecute(r GetNonCompliantResourcesApiRequest) ([]ApiAtlasNonCompliantResource, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []ApiAtlasNonCompliantResource
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.GetNonCompliantResources")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/nonCompliantResources"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetOrgResourcePolicyApiRequest struct {
	ctx              context.Context
	ApiService       ResourcePoliciesApi
	orgId            string
	resourcePolicyId string
}

type GetOrgResourcePolicyApiParams struct {
	OrgId            string
	ResourcePolicyId string
}

func (a *ResourcePoliciesApiService) GetOrgResourcePolicyWithParams(ctx context.Context, args *GetOrgResourcePolicyApiParams) GetOrgResourcePolicyApiRequest {
	return GetOrgResourcePolicyApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            args.OrgId,
		resourcePolicyId: args.ResourcePolicyId,
	}
}

func (r GetOrgResourcePolicyApiRequest) Execute() (*ApiAtlasResourcePolicy, *http.Response, error) {
	return r.ApiService.GetOrgResourcePolicyExecute(r)
}

/*
GetOrgResourcePolicy Return One Atlas Resource Policy

Return one Atlas Resource Policy for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
	@return GetOrgResourcePolicyApiRequest
*/
func (a *ResourcePoliciesApiService) GetOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string) GetOrgResourcePolicyApiRequest {
	return GetOrgResourcePolicyApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            orgId,
		resourcePolicyId: resourcePolicyId,
	}
}

// GetOrgResourcePolicyExecute executes the request
//
//	@return ApiAtlasResourcePolicy
func (a *ResourcePoliciesApiService) GetOrgResourcePolicyExecute(r GetOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasResourcePolicy
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.GetOrgResourcePolicy")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies/{resourcePolicyId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.resourcePolicyId == "" {
		return localVarReturnValue, nil, reportError("resourcePolicyId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"resourcePolicyId"+"}", url.PathEscape(r.resourcePolicyId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type ListOrgResourcePoliciesApiRequest struct {
	ctx        context.Context
	ApiService ResourcePoliciesApi
	orgId      string
}

type ListOrgResourcePoliciesApiParams struct {
	OrgId string
}

func (a *ResourcePoliciesApiService) ListOrgResourcePoliciesWithParams(ctx context.Context, args *ListOrgResourcePoliciesApiParams) ListOrgResourcePoliciesApiRequest {
	return ListOrgResourcePoliciesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r ListOrgResourcePoliciesApiRequest) Execute() ([]ApiAtlasResourcePolicy, *http.Response, error) {
	return r.ApiService.ListOrgResourcePoliciesExecute(r)
}

/*
ListOrgResourcePolicies Return All Atlas Resource Policies

Return all Atlas Resource Policies for the organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgResourcePoliciesApiRequest
*/
func (a *ResourcePoliciesApiService) ListOrgResourcePolicies(ctx context.Context, orgId string) ListOrgResourcePoliciesApiRequest {
	return ListOrgResourcePoliciesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgResourcePoliciesExecute executes the request
//
//	@return []ApiAtlasResourcePolicy
func (a *ResourcePoliciesApiService) ListOrgResourcePoliciesExecute(r ListOrgResourcePoliciesApiRequest) ([]ApiAtlasResourcePolicy, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []ApiAtlasResourcePolicy
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.ListOrgResourcePolicies")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type UpdateOrgResourcePolicyApiRequest struct {
	ctx                        context.Context
	ApiService                 ResourcePoliciesApi
	orgId                      string
	resourcePolicyId           string
	apiAtlasResourcePolicyEdit *ApiAtlasResourcePolicyEdit
}

type UpdateOrgResourcePolicyApiParams struct {
	OrgId                      string
	ResourcePolicyId           string
	ApiAtlasResourcePolicyEdit *ApiAtlasResourcePolicyEdit
}

func (a *ResourcePoliciesApiService) UpdateOrgResourcePolicyWithParams(ctx context.Context, args *UpdateOrgResourcePolicyApiParams) UpdateOrgResourcePolicyApiRequest {
	return UpdateOrgResourcePolicyApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		orgId:                      args.OrgId,
		resourcePolicyId:           args.ResourcePolicyId,
		apiAtlasResourcePolicyEdit: args.ApiAtlasResourcePolicyEdit,
	}
}

func (r UpdateOrgResourcePolicyApiRequest) Execute() (*ApiAtlasResourcePolicy, *http.Response, error) {
	return r.ApiService.UpdateOrgResourcePolicyExecute(r)
}

/*
UpdateOrgResourcePolicy Update One Atlas Resource Policy

Update one Atlas Resource Policy for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param resourcePolicyId Unique 24-hexadecimal digit string that identifies an atlas resource policy.
	@return UpdateOrgResourcePolicyApiRequest
*/
func (a *ResourcePoliciesApiService) UpdateOrgResourcePolicy(ctx context.Context, orgId string, resourcePolicyId string, apiAtlasResourcePolicyEdit *ApiAtlasResourcePolicyEdit) UpdateOrgResourcePolicyApiRequest {
	return UpdateOrgResourcePolicyApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		orgId:                      orgId,
		resourcePolicyId:           resourcePolicyId,
		apiAtlasResourcePolicyEdit: apiAtlasResourcePolicyEdit,
	}
}

// UpdateOrgResourcePolicyExecute executes the request
//
//	@return ApiAtlasResourcePolicy
func (a *ResourcePoliciesApiService) UpdateOrgResourcePolicyExecute(r UpdateOrgResourcePolicyApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasResourcePolicy
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.UpdateOrgResourcePolicy")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies/{resourcePolicyId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.resourcePolicyId == "" {
		return localVarReturnValue, nil, reportError("resourcePolicyId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"resourcePolicyId"+"}", url.PathEscape(r.resourcePolicyId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.apiAtlasResourcePolicyEdit == nil {
		return localVarReturnValue, nil, reportError("apiAtlasResourcePolicyEdit is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.apiAtlasResourcePolicyEdit
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

type ValidateResourcePoliciesApiRequest struct {
	ctx                          context.Context
	ApiService                   ResourcePoliciesApi
	orgId                        string
	apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate
}

type ValidateResourcePoliciesApiParams struct {
	OrgId                        string
	ApiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate
}

func (a *ResourcePoliciesApiService) ValidateResourcePoliciesWithParams(ctx context.Context, args *ValidateResourcePoliciesApiParams) ValidateResourcePoliciesApiRequest {
	return ValidateResourcePoliciesApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		orgId:                        args.OrgId,
		apiAtlasResourcePolicyCreate: args.ApiAtlasResourcePolicyCreate,
	}
}

func (r ValidateResourcePoliciesApiRequest) Execute() (*ApiAtlasResourcePolicy, *http.Response, error) {
	return r.ApiService.ValidateResourcePoliciesExecute(r)
}

/*
ValidateResourcePolicies Validate One Atlas Resource Policy

Validate one Atlas Resource Policy for an organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ValidateResourcePoliciesApiRequest
*/
func (a *ResourcePoliciesApiService) ValidateResourcePolicies(ctx context.Context, orgId string, apiAtlasResourcePolicyCreate *ApiAtlasResourcePolicyCreate) ValidateResourcePoliciesApiRequest {
	return ValidateResourcePoliciesApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		orgId:                        orgId,
		apiAtlasResourcePolicyCreate: apiAtlasResourcePolicyCreate,
	}
}

// ValidateResourcePoliciesExecute executes the request
//
//	@return ApiAtlasResourcePolicy
func (a *ResourcePoliciesApiService) ValidateResourcePoliciesExecute(r ValidateResourcePoliciesApiRequest) (*ApiAtlasResourcePolicy, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasResourcePolicy
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ResourcePoliciesApiService.ValidateResourcePolicies")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/resourcePolicies:validate"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.apiAtlasResourcePolicyCreate == nil {
		return localVarReturnValue, nil, reportError("apiAtlasResourcePolicyCreate is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.apiAtlasResourcePolicyCreate
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
