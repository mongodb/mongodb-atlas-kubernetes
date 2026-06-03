// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type RootApi interface {

	/*
		GetSystemStatus Return the Status of This MongoDB Application

		This resource returns information about the MongoDB application along with API key meta data.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return GetSystemStatusApiRequest
	*/
	GetSystemStatus(ctx context.Context) GetSystemStatusApiRequest
	/*
		GetSystemStatus Return the Status of This MongoDB Application


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetSystemStatusApiParams - Parameters for the request
		@return GetSystemStatusApiRequest
	*/
	GetSystemStatusWithParams(ctx context.Context, args *GetSystemStatusApiParams) GetSystemStatusApiRequest

	// Method available only for mocking purposes
	GetSystemStatusExecute(r GetSystemStatusApiRequest) (*SystemStatus, *http.Response, error)

	/*
		ListControlPlaneAddresses Return All Control Plane IP Addresses

		Returns all control plane IP addresses.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListControlPlaneAddressesApiRequest
	*/
	ListControlPlaneAddresses(ctx context.Context) ListControlPlaneAddressesApiRequest
	/*
		ListControlPlaneAddresses Return All Control Plane IP Addresses


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListControlPlaneAddressesApiParams - Parameters for the request
		@return ListControlPlaneAddressesApiRequest
	*/
	ListControlPlaneAddressesWithParams(ctx context.Context, args *ListControlPlaneAddressesApiParams) ListControlPlaneAddressesApiRequest

	// Method available only for mocking purposes
	ListControlPlaneAddressesExecute(r ListControlPlaneAddressesApiRequest) (*ControlPlaneIPAddresses, *http.Response, error)
}

// RootApiService RootApi service
type RootApiService service

type GetSystemStatusApiRequest struct {
	ctx        context.Context
	ApiService RootApi
}

type GetSystemStatusApiParams struct {
}

func (a *RootApiService) GetSystemStatusWithParams(ctx context.Context, args *GetSystemStatusApiParams) GetSystemStatusApiRequest {
	return GetSystemStatusApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

func (r GetSystemStatusApiRequest) Execute() (*SystemStatus, *http.Response, error) {
	return r.ApiService.GetSystemStatusExecute(r)
}

/*
GetSystemStatus Return the Status of This MongoDB Application

This resource returns information about the MongoDB application along with API key meta data.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return GetSystemStatusApiRequest
*/
func (a *RootApiService) GetSystemStatus(ctx context.Context) GetSystemStatusApiRequest {
	return GetSystemStatusApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// GetSystemStatusExecute executes the request
//
//	@return SystemStatus
func (a *RootApiService) GetSystemStatusExecute(r GetSystemStatusApiRequest) (*SystemStatus, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SystemStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RootApiService.GetSystemStatus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2"

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

type ListControlPlaneAddressesApiRequest struct {
	ctx        context.Context
	ApiService RootApi
}

type ListControlPlaneAddressesApiParams struct {
}

func (a *RootApiService) ListControlPlaneAddressesWithParams(ctx context.Context, args *ListControlPlaneAddressesApiParams) ListControlPlaneAddressesApiRequest {
	return ListControlPlaneAddressesApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

func (r ListControlPlaneAddressesApiRequest) Execute() (*ControlPlaneIPAddresses, *http.Response, error) {
	return r.ApiService.ListControlPlaneAddressesExecute(r)
}

/*
ListControlPlaneAddresses Return All Control Plane IP Addresses

Returns all control plane IP addresses.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListControlPlaneAddressesApiRequest
*/
func (a *RootApiService) ListControlPlaneAddresses(ctx context.Context) ListControlPlaneAddressesApiRequest {
	return ListControlPlaneAddressesApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListControlPlaneAddressesExecute executes the request
//
//	@return ControlPlaneIPAddresses
func (a *RootApiService) ListControlPlaneAddressesExecute(r ListControlPlaneAddressesApiRequest) (*ControlPlaneIPAddresses, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ControlPlaneIPAddresses
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RootApiService.ListControlPlaneAddresses")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/unauth/controlPlaneIPAddresses"

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-11-15+json"}

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
