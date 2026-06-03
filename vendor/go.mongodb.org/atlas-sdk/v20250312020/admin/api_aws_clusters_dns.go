// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AWSClustersDNSApi interface {

	/*
		GetAwsCustomDns Return One Custom DNS Configuration for Atlas Clusters on AWS

		Returns the custom DNS configuration for AWS clusters in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetAwsCustomDnsApiRequest
	*/
	GetAwsCustomDns(ctx context.Context, groupId string) GetAwsCustomDnsApiRequest
	/*
		GetAwsCustomDns Return One Custom DNS Configuration for Atlas Clusters on AWS


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAwsCustomDnsApiParams - Parameters for the request
		@return GetAwsCustomDnsApiRequest
	*/
	GetAwsCustomDnsWithParams(ctx context.Context, args *GetAwsCustomDnsApiParams) GetAwsCustomDnsApiRequest

	// Method available only for mocking purposes
	GetAwsCustomDnsExecute(r GetAwsCustomDnsApiRequest) (*AWSCustomDNSEnabled, *http.Response, error)

	/*
		ToggleAwsCustomDns Update State of One Custom DNS Configuration for Atlas Clusters on AWS

		Enables or disables the custom DNS configuration for AWS clusters in the specified project. Enable custom DNS if you use AWS VPC peering and use your own DNS servers.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param aWSCustomDNSEnabled Enables or disables the custom DNS configuration for AWS clusters in the specified project.
		@return ToggleAwsCustomDnsApiRequest
	*/
	ToggleAwsCustomDns(ctx context.Context, groupId string, aWSCustomDNSEnabled *AWSCustomDNSEnabled) ToggleAwsCustomDnsApiRequest
	/*
		ToggleAwsCustomDns Update State of One Custom DNS Configuration for Atlas Clusters on AWS


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ToggleAwsCustomDnsApiParams - Parameters for the request
		@return ToggleAwsCustomDnsApiRequest
	*/
	ToggleAwsCustomDnsWithParams(ctx context.Context, args *ToggleAwsCustomDnsApiParams) ToggleAwsCustomDnsApiRequest

	// Method available only for mocking purposes
	ToggleAwsCustomDnsExecute(r ToggleAwsCustomDnsApiRequest) (*AWSCustomDNSEnabled, *http.Response, error)
}

// AWSClustersDNSApiService AWSClustersDNSApi service
type AWSClustersDNSApiService service

type GetAwsCustomDnsApiRequest struct {
	ctx        context.Context
	ApiService AWSClustersDNSApi
	groupId    string
}

type GetAwsCustomDnsApiParams struct {
	GroupId string
}

func (a *AWSClustersDNSApiService) GetAwsCustomDnsWithParams(ctx context.Context, args *GetAwsCustomDnsApiParams) GetAwsCustomDnsApiRequest {
	return GetAwsCustomDnsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetAwsCustomDnsApiRequest) Execute() (*AWSCustomDNSEnabled, *http.Response, error) {
	return r.ApiService.GetAwsCustomDnsExecute(r)
}

/*
GetAwsCustomDns Return One Custom DNS Configuration for Atlas Clusters on AWS

Returns the custom DNS configuration for AWS clusters in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetAwsCustomDnsApiRequest
*/
func (a *AWSClustersDNSApiService) GetAwsCustomDns(ctx context.Context, groupId string) GetAwsCustomDnsApiRequest {
	return GetAwsCustomDnsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetAwsCustomDnsExecute executes the request
//
//	@return AWSCustomDNSEnabled
func (a *AWSClustersDNSApiService) GetAwsCustomDnsExecute(r GetAwsCustomDnsApiRequest) (*AWSCustomDNSEnabled, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AWSCustomDNSEnabled
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AWSClustersDNSApiService.GetAwsCustomDns")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/awsCustomDNS"
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

type ToggleAwsCustomDnsApiRequest struct {
	ctx                 context.Context
	ApiService          AWSClustersDNSApi
	groupId             string
	aWSCustomDNSEnabled *AWSCustomDNSEnabled
}

type ToggleAwsCustomDnsApiParams struct {
	GroupId             string
	AWSCustomDNSEnabled *AWSCustomDNSEnabled
}

func (a *AWSClustersDNSApiService) ToggleAwsCustomDnsWithParams(ctx context.Context, args *ToggleAwsCustomDnsApiParams) ToggleAwsCustomDnsApiRequest {
	return ToggleAwsCustomDnsApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             args.GroupId,
		aWSCustomDNSEnabled: args.AWSCustomDNSEnabled,
	}
}

func (r ToggleAwsCustomDnsApiRequest) Execute() (*AWSCustomDNSEnabled, *http.Response, error) {
	return r.ApiService.ToggleAwsCustomDnsExecute(r)
}

/*
ToggleAwsCustomDns Update State of One Custom DNS Configuration for Atlas Clusters on AWS

Enables or disables the custom DNS configuration for AWS clusters in the specified project. Enable custom DNS if you use AWS VPC peering and use your own DNS servers.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ToggleAwsCustomDnsApiRequest
*/
func (a *AWSClustersDNSApiService) ToggleAwsCustomDns(ctx context.Context, groupId string, aWSCustomDNSEnabled *AWSCustomDNSEnabled) ToggleAwsCustomDnsApiRequest {
	return ToggleAwsCustomDnsApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             groupId,
		aWSCustomDNSEnabled: aWSCustomDNSEnabled,
	}
}

// ToggleAwsCustomDnsExecute executes the request
//
//	@return AWSCustomDNSEnabled
func (a *AWSClustersDNSApiService) ToggleAwsCustomDnsExecute(r ToggleAwsCustomDnsApiRequest) (*AWSCustomDNSEnabled, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AWSCustomDNSEnabled
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AWSClustersDNSApiService.ToggleAwsCustomDns")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/awsCustomDNS"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.aWSCustomDNSEnabled == nil {
		return localVarReturnValue, nil, reportError("aWSCustomDNSEnabled is required and must be specified")
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
	localVarPostBody = r.aWSCustomDNSEnabled
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
