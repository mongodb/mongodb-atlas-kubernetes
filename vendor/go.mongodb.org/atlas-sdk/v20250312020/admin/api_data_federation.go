// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DataFederationApi interface {

	/*
		CreateDataFederation Create One Federated Database Instance in One Project

		Creates one federated database instance in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param dataLakeTenant Details to create one federated database instance in the specified project.
		@return CreateDataFederationApiRequest
	*/
	CreateDataFederation(ctx context.Context, groupId string, dataLakeTenant *DataLakeTenant) CreateDataFederationApiRequest
	/*
		CreateDataFederation Create One Federated Database Instance in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateDataFederationApiParams - Parameters for the request
		@return CreateDataFederationApiRequest
	*/
	CreateDataFederationWithParams(ctx context.Context, args *CreateDataFederationApiParams) CreateDataFederationApiRequest

	// Method available only for mocking purposes
	CreateDataFederationExecute(r CreateDataFederationApiRequest) (*DataLakeTenant, *http.Response, error)

	/*
		CreatePrivateEndpointId Create One Federated Database Instance and Online Archive Private Endpoint for One Project

		Adds one private endpoint for Federated Database Instances and Online Archives to the specified projects. If the endpoint ID already exists and the associated comment is unchanged, Atlas Data Federation makes no change to the endpoint ID list. If the endpoint ID already exists and the associated comment is changed, Atlas Data Federation updates the comment value only in the endpoint ID list. If the endpoint ID doesn't exist, Atlas Data Federation appends the new endpoint to the list of endpoints in the endpoint ID list. Each region has an associated service name for the various endpoints. For the latest list of supported regions and their service names, see the external documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param privateNetworkEndpointIdEntry Private endpoint for Federated Database Instances and Online Archives to add to the specified project.
		@return CreatePrivateEndpointIdApiRequest
	*/
	CreatePrivateEndpointId(ctx context.Context, groupId string, privateNetworkEndpointIdEntry *PrivateNetworkEndpointIdEntry) CreatePrivateEndpointIdApiRequest
	/*
		CreatePrivateEndpointId Create One Federated Database Instance and Online Archive Private Endpoint for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreatePrivateEndpointIdApiParams - Parameters for the request
		@return CreatePrivateEndpointIdApiRequest
	*/
	CreatePrivateEndpointIdWithParams(ctx context.Context, args *CreatePrivateEndpointIdApiParams) CreatePrivateEndpointIdApiRequest

	// Method available only for mocking purposes
	CreatePrivateEndpointIdExecute(r CreatePrivateEndpointIdApiRequest) (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error)

	/*
		DeleteDataFederation Remove One Federated Database Instance from One Project

		Removes one federated database instance from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance to remove.
		@return DeleteDataFederationApiRequest
	*/
	DeleteDataFederation(ctx context.Context, groupId string, tenantName string) DeleteDataFederationApiRequest
	/*
		DeleteDataFederation Remove One Federated Database Instance from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteDataFederationApiParams - Parameters for the request
		@return DeleteDataFederationApiRequest
	*/
	DeleteDataFederationWithParams(ctx context.Context, args *DeleteDataFederationApiParams) DeleteDataFederationApiRequest

	// Method available only for mocking purposes
	DeleteDataFederationExecute(r DeleteDataFederationApiRequest) (*http.Response, error)

	/*
		DeleteDataFederationLimit Delete One Query Limit for One Federated Database Instance

		Deletes one query limit for one federated database instance.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
		@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
		@return DeleteDataFederationLimitApiRequest
	*/
	DeleteDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string) DeleteDataFederationLimitApiRequest
	/*
		DeleteDataFederationLimit Delete One Query Limit for One Federated Database Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteDataFederationLimitApiParams - Parameters for the request
		@return DeleteDataFederationLimitApiRequest
	*/
	DeleteDataFederationLimitWithParams(ctx context.Context, args *DeleteDataFederationLimitApiParams) DeleteDataFederationLimitApiRequest

	// Method available only for mocking purposes
	DeleteDataFederationLimitExecute(r DeleteDataFederationLimitApiRequest) (*http.Response, error)

	/*
		DeletePrivateEndpointId Remove One Federated Database Instance and Online Archive Private Endpoint from One Project

		Removes one private endpoint for Federated Database Instances and Online Archives in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param endpointId Unique 22-character alphanumeric string that identifies the private endpoint to remove. Atlas Data Federation supports AWS private endpoints using the AWS PrivateLink feature.
		@return DeletePrivateEndpointIdApiRequest
	*/
	DeletePrivateEndpointId(ctx context.Context, groupId string, endpointId string) DeletePrivateEndpointIdApiRequest
	/*
		DeletePrivateEndpointId Remove One Federated Database Instance and Online Archive Private Endpoint from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeletePrivateEndpointIdApiParams - Parameters for the request
		@return DeletePrivateEndpointIdApiRequest
	*/
	DeletePrivateEndpointIdWithParams(ctx context.Context, args *DeletePrivateEndpointIdApiParams) DeletePrivateEndpointIdApiRequest

	// Method available only for mocking purposes
	DeletePrivateEndpointIdExecute(r DeletePrivateEndpointIdApiRequest) (*http.Response, error)

	/*
		DownloadFederationQueryLogs Download Query Logs for One Federated Database Instance

		Downloads the query logs for the specified federated database instance. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance for which you want to download query logs.
		@return DownloadFederationQueryLogsApiRequest
	*/
	DownloadFederationQueryLogs(ctx context.Context, groupId string, tenantName string) DownloadFederationQueryLogsApiRequest
	/*
		DownloadFederationQueryLogs Download Query Logs for One Federated Database Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadFederationQueryLogsApiParams - Parameters for the request
		@return DownloadFederationQueryLogsApiRequest
	*/
	DownloadFederationQueryLogsWithParams(ctx context.Context, args *DownloadFederationQueryLogsApiParams) DownloadFederationQueryLogsApiRequest

	// Method available only for mocking purposes
	DownloadFederationQueryLogsExecute(r DownloadFederationQueryLogsApiRequest) (io.ReadCloser, *http.Response, error)

	/*
		GetDataFederation Return One Federated Database Instance in One Project

		Returns the details of one federated database instance within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the Federated Database to return.
		@return GetDataFederationApiRequest
	*/
	GetDataFederation(ctx context.Context, groupId string, tenantName string) GetDataFederationApiRequest
	/*
		GetDataFederation Return One Federated Database Instance in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetDataFederationApiParams - Parameters for the request
		@return GetDataFederationApiRequest
	*/
	GetDataFederationWithParams(ctx context.Context, args *GetDataFederationApiParams) GetDataFederationApiRequest

	// Method available only for mocking purposes
	GetDataFederationExecute(r GetDataFederationApiRequest) (*DataLakeTenant, *http.Response, error)

	/*
		GetDataFederationLimit Return One Federated Database Instance Query Limit for One Project

		Returns the details of one query limit for the specified federated database instance in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
		@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
		@return GetDataFederationLimitApiRequest
	*/
	GetDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string) GetDataFederationLimitApiRequest
	/*
		GetDataFederationLimit Return One Federated Database Instance Query Limit for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetDataFederationLimitApiParams - Parameters for the request
		@return GetDataFederationLimitApiRequest
	*/
	GetDataFederationLimitWithParams(ctx context.Context, args *GetDataFederationLimitApiParams) GetDataFederationLimitApiRequest

	// Method available only for mocking purposes
	GetDataFederationLimitExecute(r GetDataFederationLimitApiRequest) (*DataFederationTenantQueryLimit, *http.Response, error)

	/*
		GetPrivateEndpointId Return One Federated Database Instance and Online Archive Private Endpoint in One Project

		Returns the specified private endpoint for Federated Database Instances or Online Archives in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param endpointId Unique 22-character alphanumeric string that identifies the private endpoint to return. Atlas Data Federation supports AWS private endpoints using the AWS PrivateLink feature.
		@return GetPrivateEndpointIdApiRequest
	*/
	GetPrivateEndpointId(ctx context.Context, groupId string, endpointId string) GetPrivateEndpointIdApiRequest
	/*
		GetPrivateEndpointId Return One Federated Database Instance and Online Archive Private Endpoint in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetPrivateEndpointIdApiParams - Parameters for the request
		@return GetPrivateEndpointIdApiRequest
	*/
	GetPrivateEndpointIdWithParams(ctx context.Context, args *GetPrivateEndpointIdApiParams) GetPrivateEndpointIdApiRequest

	// Method available only for mocking purposes
	GetPrivateEndpointIdExecute(r GetPrivateEndpointIdApiRequest) (*PrivateNetworkEndpointIdEntry, *http.Response, error)

	/*
		ListDataFederation Return All Federated Database Instances in One Project

		Returns the details of all federated database instances in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListDataFederationApiRequest
	*/
	ListDataFederation(ctx context.Context, groupId string) ListDataFederationApiRequest
	/*
		ListDataFederation Return All Federated Database Instances in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDataFederationApiParams - Parameters for the request
		@return ListDataFederationApiRequest
	*/
	ListDataFederationWithParams(ctx context.Context, args *ListDataFederationApiParams) ListDataFederationApiRequest

	// Method available only for mocking purposes
	ListDataFederationExecute(r ListDataFederationApiRequest) ([]DataLakeTenant, *http.Response, error)

	/*
		ListDataFederationLimits Return All Query Limits for One Federated Database Instance

		Returns query limits for a federated databases instance in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance for which you want to retrieve query limits.
		@return ListDataFederationLimitsApiRequest
	*/
	ListDataFederationLimits(ctx context.Context, groupId string, tenantName string) ListDataFederationLimitsApiRequest
	/*
		ListDataFederationLimits Return All Query Limits for One Federated Database Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDataFederationLimitsApiParams - Parameters for the request
		@return ListDataFederationLimitsApiRequest
	*/
	ListDataFederationLimitsWithParams(ctx context.Context, args *ListDataFederationLimitsApiParams) ListDataFederationLimitsApiRequest

	// Method available only for mocking purposes
	ListDataFederationLimitsExecute(r ListDataFederationLimitsApiRequest) ([]DataFederationTenantQueryLimit, *http.Response, error)

	/*
		ListPrivateEndpointIds Return All Federated Database Instance and Online Archive Private Endpoints in One Project

		Returns all private endpoints for Federated Database Instances and Online Archives in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListPrivateEndpointIdsApiRequest
	*/
	ListPrivateEndpointIds(ctx context.Context, groupId string) ListPrivateEndpointIdsApiRequest
	/*
		ListPrivateEndpointIds Return All Federated Database Instance and Online Archive Private Endpoints in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPrivateEndpointIdsApiParams - Parameters for the request
		@return ListPrivateEndpointIdsApiRequest
	*/
	ListPrivateEndpointIdsWithParams(ctx context.Context, args *ListPrivateEndpointIdsApiParams) ListPrivateEndpointIdsApiRequest

	// Method available only for mocking purposes
	ListPrivateEndpointIdsExecute(r ListPrivateEndpointIdsApiRequest) (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error)

	/*
		SetDataFederationLimit Configure One Query Limit for One Federated Database Instance

		Creates or updates one query limit for one federated database instance.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
		@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
		@param dataFederationTenantQueryLimit Creates or updates one query limit for one federated database instance.
		@return SetDataFederationLimitApiRequest
	*/
	SetDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string, dataFederationTenantQueryLimit *DataFederationTenantQueryLimit) SetDataFederationLimitApiRequest
	/*
		SetDataFederationLimit Configure One Query Limit for One Federated Database Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param SetDataFederationLimitApiParams - Parameters for the request
		@return SetDataFederationLimitApiRequest
	*/
	SetDataFederationLimitWithParams(ctx context.Context, args *SetDataFederationLimitApiParams) SetDataFederationLimitApiRequest

	// Method available only for mocking purposes
	SetDataFederationLimitExecute(r SetDataFederationLimitApiRequest) (*DataFederationTenantQueryLimit, *http.Response, error)

	/*
		UpdateDataFederation Update One Federated Database Instance in One Project

		Updates the details of one federated database instance in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Human-readable label that identifies the federated database instance to update.
		@param dataLakeTenant Details of one Federated Database to update in the specified project.
		@return UpdateDataFederationApiRequest
	*/
	UpdateDataFederation(ctx context.Context, groupId string, tenantName string, dataLakeTenant *DataLakeTenant) UpdateDataFederationApiRequest
	/*
		UpdateDataFederation Update One Federated Database Instance in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateDataFederationApiParams - Parameters for the request
		@return UpdateDataFederationApiRequest
	*/
	UpdateDataFederationWithParams(ctx context.Context, args *UpdateDataFederationApiParams) UpdateDataFederationApiRequest

	// Method available only for mocking purposes
	UpdateDataFederationExecute(r UpdateDataFederationApiRequest) (*DataLakeTenant, *http.Response, error)
}

// DataFederationApiService DataFederationApi service
type DataFederationApiService service

type CreateDataFederationApiRequest struct {
	ctx                context.Context
	ApiService         DataFederationApi
	groupId            string
	dataLakeTenant     *DataLakeTenant
	skipRoleValidation *bool
}

type CreateDataFederationApiParams struct {
	GroupId            string
	DataLakeTenant     *DataLakeTenant
	SkipRoleValidation *bool
}

func (a *DataFederationApiService) CreateDataFederationWithParams(ctx context.Context, args *CreateDataFederationApiParams) CreateDataFederationApiRequest {
	return CreateDataFederationApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		dataLakeTenant:     args.DataLakeTenant,
		skipRoleValidation: args.SkipRoleValidation,
	}
}

// Flag that indicates whether this request should check if the requesting IAM role can read from the S3 bucket. AWS checks if the role can list the objects in the bucket before writing to it. Some IAM roles only need write permissions. This flag allows you to skip that check.
func (r CreateDataFederationApiRequest) SkipRoleValidation(skipRoleValidation bool) CreateDataFederationApiRequest {
	r.skipRoleValidation = &skipRoleValidation
	return r
}

func (r CreateDataFederationApiRequest) Execute() (*DataLakeTenant, *http.Response, error) {
	return r.ApiService.CreateDataFederationExecute(r)
}

/*
CreateDataFederation Create One Federated Database Instance in One Project

Creates one federated database instance in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateDataFederationApiRequest
*/
func (a *DataFederationApiService) CreateDataFederation(ctx context.Context, groupId string, dataLakeTenant *DataLakeTenant) CreateDataFederationApiRequest {
	return CreateDataFederationApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		dataLakeTenant: dataLakeTenant,
	}
}

// CreateDataFederationExecute executes the request
//
//	@return DataLakeTenant
func (a *DataFederationApiService) CreateDataFederationExecute(r CreateDataFederationApiRequest) (*DataLakeTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.CreateDataFederation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataLakeTenant == nil {
		return localVarReturnValue, nil, reportError("dataLakeTenant is required and must be specified")
	}

	if r.skipRoleValidation != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "skipRoleValidation", r.skipRoleValidation, "")
	} else {
		var defaultValue bool = false
		r.skipRoleValidation = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "skipRoleValidation", r.skipRoleValidation, "")
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
	localVarPostBody = r.dataLakeTenant
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

type CreatePrivateEndpointIdApiRequest struct {
	ctx                           context.Context
	ApiService                    DataFederationApi
	groupId                       string
	privateNetworkEndpointIdEntry *PrivateNetworkEndpointIdEntry
}

type CreatePrivateEndpointIdApiParams struct {
	GroupId                       string
	PrivateNetworkEndpointIdEntry *PrivateNetworkEndpointIdEntry
}

func (a *DataFederationApiService) CreatePrivateEndpointIdWithParams(ctx context.Context, args *CreatePrivateEndpointIdApiParams) CreatePrivateEndpointIdApiRequest {
	return CreatePrivateEndpointIdApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       args.GroupId,
		privateNetworkEndpointIdEntry: args.PrivateNetworkEndpointIdEntry,
	}
}

func (r CreatePrivateEndpointIdApiRequest) Execute() (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error) {
	return r.ApiService.CreatePrivateEndpointIdExecute(r)
}

/*
CreatePrivateEndpointId Create One Federated Database Instance and Online Archive Private Endpoint for One Project

Adds one private endpoint for Federated Database Instances and Online Archives to the specified projects. If the endpoint ID already exists and the associated comment is unchanged, Atlas Data Federation makes no change to the endpoint ID list. If the endpoint ID already exists and the associated comment is changed, Atlas Data Federation updates the comment value only in the endpoint ID list. If the endpoint ID doesn't exist, Atlas Data Federation appends the new endpoint to the list of endpoints in the endpoint ID list. Each region has an associated service name for the various endpoints. For the latest list of supported regions and their service names, see the external documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreatePrivateEndpointIdApiRequest
*/
func (a *DataFederationApiService) CreatePrivateEndpointId(ctx context.Context, groupId string, privateNetworkEndpointIdEntry *PrivateNetworkEndpointIdEntry) CreatePrivateEndpointIdApiRequest {
	return CreatePrivateEndpointIdApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       groupId,
		privateNetworkEndpointIdEntry: privateNetworkEndpointIdEntry,
	}
}

// CreatePrivateEndpointIdExecute executes the request
//
//	@return PaginatedPrivateNetworkEndpointIdEntry
func (a *DataFederationApiService) CreatePrivateEndpointIdExecute(r CreatePrivateEndpointIdApiRequest) (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedPrivateNetworkEndpointIdEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.CreatePrivateEndpointId")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateNetworkSettings/endpointIds"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.privateNetworkEndpointIdEntry == nil {
		return localVarReturnValue, nil, reportError("privateNetworkEndpointIdEntry is required and must be specified")
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
	localVarPostBody = r.privateNetworkEndpointIdEntry
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

type DeleteDataFederationApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
}

type DeleteDataFederationApiParams struct {
	GroupId    string
	TenantName string
}

func (a *DataFederationApiService) DeleteDataFederationWithParams(ctx context.Context, args *DeleteDataFederationApiParams) DeleteDataFederationApiRequest {
	return DeleteDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
	}
}

func (r DeleteDataFederationApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteDataFederationExecute(r)
}

/*
DeleteDataFederation Remove One Federated Database Instance from One Project

Removes one federated database instance from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance to remove.
	@return DeleteDataFederationApiRequest
*/
func (a *DataFederationApiService) DeleteDataFederation(ctx context.Context, groupId string, tenantName string) DeleteDataFederationApiRequest {
	return DeleteDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// DeleteDataFederationExecute executes the request
func (a *DataFederationApiService) DeleteDataFederationExecute(r DeleteDataFederationApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.DeleteDataFederation")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

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

type DeleteDataFederationLimitApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
	limitName  string
}

type DeleteDataFederationLimitApiParams struct {
	GroupId    string
	TenantName string
	LimitName  string
}

func (a *DataFederationApiService) DeleteDataFederationLimitWithParams(ctx context.Context, args *DeleteDataFederationLimitApiParams) DeleteDataFederationLimitApiRequest {
	return DeleteDataFederationLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
		limitName:  args.LimitName,
	}
}

func (r DeleteDataFederationLimitApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteDataFederationLimitExecute(r)
}

/*
DeleteDataFederationLimit Delete One Query Limit for One Federated Database Instance

Deletes one query limit for one federated database instance.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
	@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
	@return DeleteDataFederationLimitApiRequest
*/
func (a *DataFederationApiService) DeleteDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string) DeleteDataFederationLimitApiRequest {
	return DeleteDataFederationLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
		limitName:  limitName,
	}
}

// DeleteDataFederationLimitExecute executes the request
func (a *DataFederationApiService) DeleteDataFederationLimitExecute(r DeleteDataFederationLimitApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.DeleteDataFederationLimit")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}/limits/{limitName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.limitName == "" {
		return nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)

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

type DeletePrivateEndpointIdApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	endpointId string
}

type DeletePrivateEndpointIdApiParams struct {
	GroupId    string
	EndpointId string
}

func (a *DataFederationApiService) DeletePrivateEndpointIdWithParams(ctx context.Context, args *DeletePrivateEndpointIdApiParams) DeletePrivateEndpointIdApiRequest {
	return DeletePrivateEndpointIdApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		endpointId: args.EndpointId,
	}
}

func (r DeletePrivateEndpointIdApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeletePrivateEndpointIdExecute(r)
}

/*
DeletePrivateEndpointId Remove One Federated Database Instance and Online Archive Private Endpoint from One Project

Removes one private endpoint for Federated Database Instances and Online Archives in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param endpointId Unique 22-character alphanumeric string that identifies the private endpoint to remove. Atlas Data Federation supports AWS private endpoints using the AWS PrivateLink feature.
	@return DeletePrivateEndpointIdApiRequest
*/
func (a *DataFederationApiService) DeletePrivateEndpointId(ctx context.Context, groupId string, endpointId string) DeletePrivateEndpointIdApiRequest {
	return DeletePrivateEndpointIdApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		endpointId: endpointId,
	}
}

// DeletePrivateEndpointIdExecute executes the request
func (a *DataFederationApiService) DeletePrivateEndpointIdExecute(r DeletePrivateEndpointIdApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.DeletePrivateEndpointId")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateNetworkSettings/endpointIds/{endpointId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.endpointId == "" {
		return nil, reportError("endpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointId"+"}", url.PathEscape(r.endpointId), -1)

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

type DownloadFederationQueryLogsApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
	endDate    *int64
	startDate  *int64
}

type DownloadFederationQueryLogsApiParams struct {
	GroupId    string
	TenantName string
	EndDate    *int64
	StartDate  *int64
}

func (a *DataFederationApiService) DownloadFederationQueryLogsWithParams(ctx context.Context, args *DownloadFederationQueryLogsApiParams) DownloadFederationQueryLogsApiRequest {
	return DownloadFederationQueryLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
		endDate:    args.EndDate,
		startDate:  args.StartDate,
	}
}

// Timestamp that specifies the end point for the range of log messages to download.  MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadFederationQueryLogsApiRequest) EndDate(endDate int64) DownloadFederationQueryLogsApiRequest {
	r.endDate = &endDate
	return r
}

// Timestamp that specifies the starting point for the range of log messages to download. MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadFederationQueryLogsApiRequest) StartDate(startDate int64) DownloadFederationQueryLogsApiRequest {
	r.startDate = &startDate
	return r
}

func (r DownloadFederationQueryLogsApiRequest) Execute() (io.ReadCloser, *http.Response, error) {
	return r.ApiService.DownloadFederationQueryLogsExecute(r)
}

/*
DownloadFederationQueryLogs Download Query Logs for One Federated Database Instance

Downloads the query logs for the specified federated database instance. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance for which you want to download query logs.
	@return DownloadFederationQueryLogsApiRequest
*/
func (a *DataFederationApiService) DownloadFederationQueryLogs(ctx context.Context, groupId string, tenantName string) DownloadFederationQueryLogsApiRequest {
	return DownloadFederationQueryLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// DownloadFederationQueryLogsExecute executes the request
//
//	@return io.ReadCloser
func (a *DataFederationApiService) DownloadFederationQueryLogsExecute(r DownloadFederationQueryLogsApiRequest) (io.ReadCloser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue io.ReadCloser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.DownloadFederationQueryLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}/queryLogs.gz"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.endDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endDate", r.endDate, "")
	}
	if r.startDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startDate", r.startDate, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+gzip"}

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

type GetDataFederationApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
}

type GetDataFederationApiParams struct {
	GroupId    string
	TenantName string
}

func (a *DataFederationApiService) GetDataFederationWithParams(ctx context.Context, args *GetDataFederationApiParams) GetDataFederationApiRequest {
	return GetDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
	}
}

func (r GetDataFederationApiRequest) Execute() (*DataLakeTenant, *http.Response, error) {
	return r.ApiService.GetDataFederationExecute(r)
}

/*
GetDataFederation Return One Federated Database Instance in One Project

Returns the details of one federated database instance within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the Federated Database to return.
	@return GetDataFederationApiRequest
*/
func (a *DataFederationApiService) GetDataFederation(ctx context.Context, groupId string, tenantName string) GetDataFederationApiRequest {
	return GetDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// GetDataFederationExecute executes the request
//
//	@return DataLakeTenant
func (a *DataFederationApiService) GetDataFederationExecute(r GetDataFederationApiRequest) (*DataLakeTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.GetDataFederation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

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

type GetDataFederationLimitApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
	limitName  string
}

type GetDataFederationLimitApiParams struct {
	GroupId    string
	TenantName string
	LimitName  string
}

func (a *DataFederationApiService) GetDataFederationLimitWithParams(ctx context.Context, args *GetDataFederationLimitApiParams) GetDataFederationLimitApiRequest {
	return GetDataFederationLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
		limitName:  args.LimitName,
	}
}

func (r GetDataFederationLimitApiRequest) Execute() (*DataFederationTenantQueryLimit, *http.Response, error) {
	return r.ApiService.GetDataFederationLimitExecute(r)
}

/*
GetDataFederationLimit Return One Federated Database Instance Query Limit for One Project

Returns the details of one query limit for the specified federated database instance in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
	@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
	@return GetDataFederationLimitApiRequest
*/
func (a *DataFederationApiService) GetDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string) GetDataFederationLimitApiRequest {
	return GetDataFederationLimitApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
		limitName:  limitName,
	}
}

// GetDataFederationLimitExecute executes the request
//
//	@return DataFederationTenantQueryLimit
func (a *DataFederationApiService) GetDataFederationLimitExecute(r GetDataFederationLimitApiRequest) (*DataFederationTenantQueryLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataFederationTenantQueryLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.GetDataFederationLimit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}/limits/{limitName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.limitName == "" {
		return localVarReturnValue, nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)

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

type GetPrivateEndpointIdApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	endpointId string
}

type GetPrivateEndpointIdApiParams struct {
	GroupId    string
	EndpointId string
}

func (a *DataFederationApiService) GetPrivateEndpointIdWithParams(ctx context.Context, args *GetPrivateEndpointIdApiParams) GetPrivateEndpointIdApiRequest {
	return GetPrivateEndpointIdApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		endpointId: args.EndpointId,
	}
}

func (r GetPrivateEndpointIdApiRequest) Execute() (*PrivateNetworkEndpointIdEntry, *http.Response, error) {
	return r.ApiService.GetPrivateEndpointIdExecute(r)
}

/*
GetPrivateEndpointId Return One Federated Database Instance and Online Archive Private Endpoint in One Project

Returns the specified private endpoint for Federated Database Instances or Online Archives in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param endpointId Unique 22-character alphanumeric string that identifies the private endpoint to return. Atlas Data Federation supports AWS private endpoints using the AWS PrivateLink feature.
	@return GetPrivateEndpointIdApiRequest
*/
func (a *DataFederationApiService) GetPrivateEndpointId(ctx context.Context, groupId string, endpointId string) GetPrivateEndpointIdApiRequest {
	return GetPrivateEndpointIdApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		endpointId: endpointId,
	}
}

// GetPrivateEndpointIdExecute executes the request
//
//	@return PrivateNetworkEndpointIdEntry
func (a *DataFederationApiService) GetPrivateEndpointIdExecute(r GetPrivateEndpointIdApiRequest) (*PrivateNetworkEndpointIdEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PrivateNetworkEndpointIdEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.GetPrivateEndpointId")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateNetworkSettings/endpointIds/{endpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.endpointId == "" {
		return localVarReturnValue, nil, reportError("endpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"endpointId"+"}", url.PathEscape(r.endpointId), -1)

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

type ListDataFederationApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	type_      *string
}

type ListDataFederationApiParams struct {
	GroupId string
	Type_   *string
}

func (a *DataFederationApiService) ListDataFederationWithParams(ctx context.Context, args *ListDataFederationApiParams) ListDataFederationApiRequest {
	return ListDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		type_:      args.Type_,
	}
}

// Type of Federated Database Instances to return.
func (r ListDataFederationApiRequest) Type_(type_ string) ListDataFederationApiRequest {
	r.type_ = &type_
	return r
}

func (r ListDataFederationApiRequest) Execute() ([]DataLakeTenant, *http.Response, error) {
	return r.ApiService.ListDataFederationExecute(r)
}

/*
ListDataFederation Return All Federated Database Instances in One Project

Returns the details of all federated database instances in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListDataFederationApiRequest
*/
func (a *DataFederationApiService) ListDataFederation(ctx context.Context, groupId string) ListDataFederationApiRequest {
	return ListDataFederationApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListDataFederationExecute executes the request
//
//	@return []DataLakeTenant
func (a *DataFederationApiService) ListDataFederationExecute(r ListDataFederationApiRequest) ([]DataLakeTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []DataLakeTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.ListDataFederation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.type_ != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "type", r.type_, "")
	} else {
		var defaultValue string = "USER"
		r.type_ = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "type", r.type_, "")
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

type ListDataFederationLimitsApiRequest struct {
	ctx        context.Context
	ApiService DataFederationApi
	groupId    string
	tenantName string
}

type ListDataFederationLimitsApiParams struct {
	GroupId    string
	TenantName string
}

func (a *DataFederationApiService) ListDataFederationLimitsWithParams(ctx context.Context, args *ListDataFederationLimitsApiParams) ListDataFederationLimitsApiRequest {
	return ListDataFederationLimitsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
	}
}

func (r ListDataFederationLimitsApiRequest) Execute() ([]DataFederationTenantQueryLimit, *http.Response, error) {
	return r.ApiService.ListDataFederationLimitsExecute(r)
}

/*
ListDataFederationLimits Return All Query Limits for One Federated Database Instance

Returns query limits for a federated databases instance in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance for which you want to retrieve query limits.
	@return ListDataFederationLimitsApiRequest
*/
func (a *DataFederationApiService) ListDataFederationLimits(ctx context.Context, groupId string, tenantName string) ListDataFederationLimitsApiRequest {
	return ListDataFederationLimitsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// ListDataFederationLimitsExecute executes the request
//
//	@return []DataFederationTenantQueryLimit
func (a *DataFederationApiService) ListDataFederationLimitsExecute(r ListDataFederationLimitsApiRequest) ([]DataFederationTenantQueryLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []DataFederationTenantQueryLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.ListDataFederationLimits")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}/limits"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

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

type ListPrivateEndpointIdsApiRequest struct {
	ctx          context.Context
	ApiService   DataFederationApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListPrivateEndpointIdsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *DataFederationApiService) ListPrivateEndpointIdsWithParams(ctx context.Context, args *ListPrivateEndpointIdsApiParams) ListPrivateEndpointIdsApiRequest {
	return ListPrivateEndpointIdsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListPrivateEndpointIdsApiRequest) IncludeCount(includeCount bool) ListPrivateEndpointIdsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListPrivateEndpointIdsApiRequest) ItemsPerPage(itemsPerPage int) ListPrivateEndpointIdsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListPrivateEndpointIdsApiRequest) PageNum(pageNum int) ListPrivateEndpointIdsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListPrivateEndpointIdsApiRequest) Execute() (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error) {
	return r.ApiService.ListPrivateEndpointIdsExecute(r)
}

/*
ListPrivateEndpointIds Return All Federated Database Instance and Online Archive Private Endpoints in One Project

Returns all private endpoints for Federated Database Instances and Online Archives in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListPrivateEndpointIdsApiRequest
*/
func (a *DataFederationApiService) ListPrivateEndpointIds(ctx context.Context, groupId string) ListPrivateEndpointIdsApiRequest {
	return ListPrivateEndpointIdsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListPrivateEndpointIdsExecute executes the request
//
//	@return PaginatedPrivateNetworkEndpointIdEntry
func (a *DataFederationApiService) ListPrivateEndpointIdsExecute(r ListPrivateEndpointIdsApiRequest) (*PaginatedPrivateNetworkEndpointIdEntry, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedPrivateNetworkEndpointIdEntry
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.ListPrivateEndpointIds")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateNetworkSettings/endpointIds"
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

type SetDataFederationLimitApiRequest struct {
	ctx                            context.Context
	ApiService                     DataFederationApi
	groupId                        string
	tenantName                     string
	limitName                      string
	dataFederationTenantQueryLimit *DataFederationTenantQueryLimit
}

type SetDataFederationLimitApiParams struct {
	GroupId                        string
	TenantName                     string
	LimitName                      string
	DataFederationTenantQueryLimit *DataFederationTenantQueryLimit
}

func (a *DataFederationApiService) SetDataFederationLimitWithParams(ctx context.Context, args *SetDataFederationLimitApiParams) SetDataFederationLimitApiRequest {
	return SetDataFederationLimitApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        args.GroupId,
		tenantName:                     args.TenantName,
		limitName:                      args.LimitName,
		dataFederationTenantQueryLimit: args.DataFederationTenantQueryLimit,
	}
}

func (r SetDataFederationLimitApiRequest) Execute() (*DataFederationTenantQueryLimit, *http.Response, error) {
	return r.ApiService.SetDataFederationLimitExecute(r)
}

/*
SetDataFederationLimit Configure One Query Limit for One Federated Database Instance

Creates or updates one query limit for one federated database instance.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance to which the query limit applies.
	@param limitName Human-readable label that identifies this data federation instance limit.  | Limit Name | Description | Default | | --- | --- | --- | | `bytesProcessed.query` | Limit on the number of bytes processed during a single data federation query | N/A | | `bytesProcessed.daily` | Limit on the number of bytes processed for the data federation instance for the current day | N/A | | `bytesProcessed.weekly` | Limit on the number of bytes processed for the data federation instance for the current week | N/A | | `bytesProcessed.monthly` | Limit on the number of bytes processed for the data federation instance for the current month | N/A |
	@return SetDataFederationLimitApiRequest
*/
func (a *DataFederationApiService) SetDataFederationLimit(ctx context.Context, groupId string, tenantName string, limitName string, dataFederationTenantQueryLimit *DataFederationTenantQueryLimit) SetDataFederationLimitApiRequest {
	return SetDataFederationLimitApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        groupId,
		tenantName:                     tenantName,
		limitName:                      limitName,
		dataFederationTenantQueryLimit: dataFederationTenantQueryLimit,
	}
}

// SetDataFederationLimitExecute executes the request
//
//	@return DataFederationTenantQueryLimit
func (a *DataFederationApiService) SetDataFederationLimitExecute(r SetDataFederationLimitApiRequest) (*DataFederationTenantQueryLimit, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataFederationTenantQueryLimit
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.SetDataFederationLimit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}/limits/{limitName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.limitName == "" {
		return localVarReturnValue, nil, reportError("limitName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"limitName"+"}", url.PathEscape(r.limitName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataFederationTenantQueryLimit == nil {
		return localVarReturnValue, nil, reportError("dataFederationTenantQueryLimit is required and must be specified")
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
	localVarPostBody = r.dataFederationTenantQueryLimit
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

type UpdateDataFederationApiRequest struct {
	ctx                context.Context
	ApiService         DataFederationApi
	groupId            string
	tenantName         string
	skipRoleValidation *bool
	dataLakeTenant     *DataLakeTenant
}

type UpdateDataFederationApiParams struct {
	GroupId            string
	TenantName         string
	SkipRoleValidation *bool
	DataLakeTenant     *DataLakeTenant
}

func (a *DataFederationApiService) UpdateDataFederationWithParams(ctx context.Context, args *UpdateDataFederationApiParams) UpdateDataFederationApiRequest {
	return UpdateDataFederationApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		tenantName:         args.TenantName,
		skipRoleValidation: args.SkipRoleValidation,
		dataLakeTenant:     args.DataLakeTenant,
	}
}

// Flag that indicates whether this request should check if the requesting IAM role can read from the S3 bucket. AWS checks if the role can list the objects in the bucket before writing to it. Some IAM roles only need write permissions. This flag allows you to skip that check.
func (r UpdateDataFederationApiRequest) SkipRoleValidation(skipRoleValidation bool) UpdateDataFederationApiRequest {
	r.skipRoleValidation = &skipRoleValidation
	return r
}

func (r UpdateDataFederationApiRequest) Execute() (*DataLakeTenant, *http.Response, error) {
	return r.ApiService.UpdateDataFederationExecute(r)
}

/*
UpdateDataFederation Update One Federated Database Instance in One Project

Updates the details of one federated database instance in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Human-readable label that identifies the federated database instance to update.
	@return UpdateDataFederationApiRequest
*/
func (a *DataFederationApiService) UpdateDataFederation(ctx context.Context, groupId string, tenantName string, dataLakeTenant *DataLakeTenant) UpdateDataFederationApiRequest {
	return UpdateDataFederationApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		tenantName:     tenantName,
		dataLakeTenant: dataLakeTenant,
	}
}

// UpdateDataFederationExecute executes the request
//
//	@return DataLakeTenant
func (a *DataFederationApiService) UpdateDataFederationExecute(r UpdateDataFederationApiRequest) (*DataLakeTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataFederationApiService.UpdateDataFederation")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/dataFederation/{tenantName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.skipRoleValidation == nil {
		return localVarReturnValue, nil, reportError("skipRoleValidation is required and must be specified")
	}
	if r.dataLakeTenant == nil {
		return localVarReturnValue, nil, reportError("dataLakeTenant is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "skipRoleValidation", r.skipRoleValidation, "")
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
	localVarPostBody = r.dataLakeTenant
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
