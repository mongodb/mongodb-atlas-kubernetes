// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CloudMigrationServiceApi interface {

	/*
			CreateGroupLiveMigration Create One Migration for One Local Managed Cluster to MongoDB Atlas

			Migrate one cluster that Cloud or Ops Manager manages to MongoDB Atlas.

		 Please make sure to validate your migration before initiating it.

		 You can use this API endpoint for push live migrations only. Your API Key must have the Organization Owner role to successfully call this resource.

		 **NOTE**: Migrating time-series collections is not yet supported on MongoDB 6.0 or higher. Migrations on MongoDB 6.0 or higher will skip any time-series collections on the source cluster. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param liveMigrationRequest20240530 One migration to be created.
			@return CreateGroupLiveMigrationApiRequest
	*/
	CreateGroupLiveMigration(ctx context.Context, groupId string, liveMigrationRequest20240530 *LiveMigrationRequest20240530) CreateGroupLiveMigrationApiRequest
	/*
		CreateGroupLiveMigration Create One Migration for One Local Managed Cluster to MongoDB Atlas


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupLiveMigrationApiParams - Parameters for the request
		@return CreateGroupLiveMigrationApiRequest
	*/
	CreateGroupLiveMigrationWithParams(ctx context.Context, args *CreateGroupLiveMigrationApiParams) CreateGroupLiveMigrationApiRequest

	// Method available only for mocking purposes
	CreateGroupLiveMigrationExecute(r CreateGroupLiveMigrationApiRequest) (*LiveMigrationResponse, *http.Response, error)

	/*
		CreateLinkToken Create One Link-Token

		Create one link-token that contains all the information required to complete the link. MongoDB Atlas uses the link-token for push live migrations only. Live migration (push) allows you to securely push data from Cloud Manager or Ops Manager into MongoDB Atlas. Your API Key must have the Organization Owner role to successfully call this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param targetOrgRequest IP address access list entries associated with the migration.
		@return CreateLinkTokenApiRequest
	*/
	CreateLinkToken(ctx context.Context, orgId string, targetOrgRequest *TargetOrgRequest) CreateLinkTokenApiRequest
	/*
		CreateLinkToken Create One Link-Token


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateLinkTokenApiParams - Parameters for the request
		@return CreateLinkTokenApiRequest
	*/
	CreateLinkTokenWithParams(ctx context.Context, args *CreateLinkTokenApiParams) CreateLinkTokenApiRequest

	// Method available only for mocking purposes
	CreateLinkTokenExecute(r CreateLinkTokenApiRequest) (*TargetOrg, *http.Response, error)

	/*
		CutoverMigration Cut Over One Migrated Cluster

		Cut over the migrated cluster to MongoDB Atlas. Confirm when the cut over completes. When the cut over completes, MongoDB Atlas completes the live migration process and stops synchronizing with the source cluster. Your API Key must have the Organization Owner role to successfully call this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param liveMigrationId Unique 24-hexadecimal digit string that identifies the migration.
		@return CutoverMigrationApiRequest
	*/
	CutoverMigration(ctx context.Context, groupId string, liveMigrationId string) CutoverMigrationApiRequest
	/*
		CutoverMigration Cut Over One Migrated Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CutoverMigrationApiParams - Parameters for the request
		@return CutoverMigrationApiRequest
	*/
	CutoverMigrationWithParams(ctx context.Context, args *CutoverMigrationApiParams) CutoverMigrationApiRequest

	// Method available only for mocking purposes
	CutoverMigrationExecute(r CutoverMigrationApiRequest) (*http.Response, error)

	/*
		DeleteLinkTokens Remove One Link-Token

		Remove one organization link and its associated public API key. MongoDB Atlas uses the link-token for push live migrations only. Live migrations (push) let you securely push data from Cloud Manager or Ops Manager into MongoDB Atlas. Your API Key must have the Organization Owner role to successfully call this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return DeleteLinkTokensApiRequest
	*/
	DeleteLinkTokens(ctx context.Context, orgId string) DeleteLinkTokensApiRequest
	/*
		DeleteLinkTokens Remove One Link-Token


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteLinkTokensApiParams - Parameters for the request
		@return DeleteLinkTokensApiRequest
	*/
	DeleteLinkTokensWithParams(ctx context.Context, args *DeleteLinkTokensApiParams) DeleteLinkTokensApiRequest

	// Method available only for mocking purposes
	DeleteLinkTokensExecute(r DeleteLinkTokensApiRequest) (*http.Response, error)

	/*
		GetGroupLiveMigration Return One Migration Job

		Return details of one cluster migration job. Each push live migration job uses one migration host. Your API Key must have the Organization Member role to successfully call this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param liveMigrationId Unique 24-hexadecimal digit string that identifies the migration.
		@return GetGroupLiveMigrationApiRequest
	*/
	GetGroupLiveMigration(ctx context.Context, groupId string, liveMigrationId string) GetGroupLiveMigrationApiRequest
	/*
		GetGroupLiveMigration Return One Migration Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupLiveMigrationApiParams - Parameters for the request
		@return GetGroupLiveMigrationApiRequest
	*/
	GetGroupLiveMigrationWithParams(ctx context.Context, args *GetGroupLiveMigrationApiParams) GetGroupLiveMigrationApiRequest

	// Method available only for mocking purposes
	GetGroupLiveMigrationExecute(r GetGroupLiveMigrationApiRequest) (*LiveMigrationResponse, *http.Response, error)

	/*
		GetMigrationValidateStatus Return One Migration Validation Job

		Return the status of one migration validation job. Your API Key must have the Organization Owner role to successfully call this resource.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param validationId Unique 24-hexadecimal digit string that identifies the validation job.
		@return GetMigrationValidateStatusApiRequest
	*/
	GetMigrationValidateStatus(ctx context.Context, groupId string, validationId string) GetMigrationValidateStatusApiRequest
	/*
		GetMigrationValidateStatus Return One Migration Validation Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetMigrationValidateStatusApiParams - Parameters for the request
		@return GetMigrationValidateStatusApiRequest
	*/
	GetMigrationValidateStatusWithParams(ctx context.Context, args *GetMigrationValidateStatusApiParams) GetMigrationValidateStatusApiRequest

	// Method available only for mocking purposes
	GetMigrationValidateStatusExecute(r GetMigrationValidateStatusApiRequest) (*LiveImportValidation, *http.Response, error)

	/*
		ListAvailableProjects Return All Projects Available for Migration

		Return all projects that you can migrate to the specified organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListAvailableProjectsApiRequest
	*/
	ListAvailableProjects(ctx context.Context, orgId string) ListAvailableProjectsApiRequest
	/*
		ListAvailableProjects Return All Projects Available for Migration


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListAvailableProjectsApiParams - Parameters for the request
		@return ListAvailableProjectsApiRequest
	*/
	ListAvailableProjectsWithParams(ctx context.Context, args *ListAvailableProjectsApiParams) ListAvailableProjectsApiRequest

	// Method available only for mocking purposes
	ListAvailableProjectsExecute(r ListAvailableProjectsApiRequest) ([]LiveImportAvailableProject, *http.Response, error)

	/*
		ValidateLiveMigrations Validate One Migration Request

		Verifies whether the provided credentials, available disk space, MongoDB versions, and so on meet the requirements of the migration request. If the check passes, the migration can proceed. Your API Key must have the Organization Owner role to successfully call this resource. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param liveMigrationRequest20240530 One migration to be validated.
		@return ValidateLiveMigrationsApiRequest
	*/
	ValidateLiveMigrations(ctx context.Context, groupId string, liveMigrationRequest20240530 *LiveMigrationRequest20240530) ValidateLiveMigrationsApiRequest
	/*
		ValidateLiveMigrations Validate One Migration Request


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ValidateLiveMigrationsApiParams - Parameters for the request
		@return ValidateLiveMigrationsApiRequest
	*/
	ValidateLiveMigrationsWithParams(ctx context.Context, args *ValidateLiveMigrationsApiParams) ValidateLiveMigrationsApiRequest

	// Method available only for mocking purposes
	ValidateLiveMigrationsExecute(r ValidateLiveMigrationsApiRequest) (*LiveImportValidation, *http.Response, error)
}

// CloudMigrationServiceApiService CloudMigrationServiceApi service
type CloudMigrationServiceApiService service

type CreateGroupLiveMigrationApiRequest struct {
	ctx                          context.Context
	ApiService                   CloudMigrationServiceApi
	groupId                      string
	liveMigrationRequest20240530 *LiveMigrationRequest20240530
}

type CreateGroupLiveMigrationApiParams struct {
	GroupId                      string
	LiveMigrationRequest20240530 *LiveMigrationRequest20240530
}

func (a *CloudMigrationServiceApiService) CreateGroupLiveMigrationWithParams(ctx context.Context, args *CreateGroupLiveMigrationApiParams) CreateGroupLiveMigrationApiRequest {
	return CreateGroupLiveMigrationApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		liveMigrationRequest20240530: args.LiveMigrationRequest20240530,
	}
}

func (r CreateGroupLiveMigrationApiRequest) Execute() (*LiveMigrationResponse, *http.Response, error) {
	return r.ApiService.CreateGroupLiveMigrationExecute(r)
}

/*
CreateGroupLiveMigration Create One Migration for One Local Managed Cluster to MongoDB Atlas

Migrate one cluster that Cloud or Ops Manager manages to MongoDB Atlas.

	Please make sure to validate your migration before initiating it.

	You can use this API endpoint for push live migrations only. Your API Key must have the Organization Owner role to successfully call this resource.

	**NOTE**: Migrating time-series collections is not yet supported on MongoDB 6.0 or higher. Migrations on MongoDB 6.0 or higher will skip any time-series collections on the source cluster. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupLiveMigrationApiRequest
*/
func (a *CloudMigrationServiceApiService) CreateGroupLiveMigration(ctx context.Context, groupId string, liveMigrationRequest20240530 *LiveMigrationRequest20240530) CreateGroupLiveMigrationApiRequest {
	return CreateGroupLiveMigrationApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		liveMigrationRequest20240530: liveMigrationRequest20240530,
	}
}

// CreateGroupLiveMigrationExecute executes the request
//
//	@return LiveMigrationResponse
func (a *CloudMigrationServiceApiService) CreateGroupLiveMigrationExecute(r CreateGroupLiveMigrationApiRequest) (*LiveMigrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LiveMigrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.CreateGroupLiveMigration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/liveMigrations"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.liveMigrationRequest20240530 == nil {
		return localVarReturnValue, nil, reportError("liveMigrationRequest20240530 is required and must be specified")
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
	localVarPostBody = r.liveMigrationRequest20240530
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

type CreateLinkTokenApiRequest struct {
	ctx              context.Context
	ApiService       CloudMigrationServiceApi
	orgId            string
	targetOrgRequest *TargetOrgRequest
}

type CreateLinkTokenApiParams struct {
	OrgId            string
	TargetOrgRequest *TargetOrgRequest
}

func (a *CloudMigrationServiceApiService) CreateLinkTokenWithParams(ctx context.Context, args *CreateLinkTokenApiParams) CreateLinkTokenApiRequest {
	return CreateLinkTokenApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            args.OrgId,
		targetOrgRequest: args.TargetOrgRequest,
	}
}

func (r CreateLinkTokenApiRequest) Execute() (*TargetOrg, *http.Response, error) {
	return r.ApiService.CreateLinkTokenExecute(r)
}

/*
CreateLinkToken Create One Link-Token

Create one link-token that contains all the information required to complete the link. MongoDB Atlas uses the link-token for push live migrations only. Live migration (push) allows you to securely push data from Cloud Manager or Ops Manager into MongoDB Atlas. Your API Key must have the Organization Owner role to successfully call this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateLinkTokenApiRequest
*/
func (a *CloudMigrationServiceApiService) CreateLinkToken(ctx context.Context, orgId string, targetOrgRequest *TargetOrgRequest) CreateLinkTokenApiRequest {
	return CreateLinkTokenApiRequest{
		ApiService:       a,
		ctx:              ctx,
		orgId:            orgId,
		targetOrgRequest: targetOrgRequest,
	}
}

// CreateLinkTokenExecute executes the request
//
//	@return TargetOrg
func (a *CloudMigrationServiceApiService) CreateLinkTokenExecute(r CreateLinkTokenApiRequest) (*TargetOrg, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *TargetOrg
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.CreateLinkToken")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/liveMigrations/linkTokens"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.targetOrgRequest == nil {
		return localVarReturnValue, nil, reportError("targetOrgRequest is required and must be specified")
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
	localVarPostBody = r.targetOrgRequest
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

type CutoverMigrationApiRequest struct {
	ctx             context.Context
	ApiService      CloudMigrationServiceApi
	groupId         string
	liveMigrationId string
}

type CutoverMigrationApiParams struct {
	GroupId         string
	LiveMigrationId string
}

func (a *CloudMigrationServiceApiService) CutoverMigrationWithParams(ctx context.Context, args *CutoverMigrationApiParams) CutoverMigrationApiRequest {
	return CutoverMigrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         args.GroupId,
		liveMigrationId: args.LiveMigrationId,
	}
}

func (r CutoverMigrationApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.CutoverMigrationExecute(r)
}

/*
CutoverMigration Cut Over One Migrated Cluster

Cut over the migrated cluster to MongoDB Atlas. Confirm when the cut over completes. When the cut over completes, MongoDB Atlas completes the live migration process and stops synchronizing with the source cluster. Your API Key must have the Organization Owner role to successfully call this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param liveMigrationId Unique 24-hexadecimal digit string that identifies the migration.
	@return CutoverMigrationApiRequest
*/
func (a *CloudMigrationServiceApiService) CutoverMigration(ctx context.Context, groupId string, liveMigrationId string) CutoverMigrationApiRequest {
	return CutoverMigrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         groupId,
		liveMigrationId: liveMigrationId,
	}
}

// CutoverMigrationExecute executes the request
func (a *CloudMigrationServiceApiService) CutoverMigrationExecute(r CutoverMigrationApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPut
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.CutoverMigration")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/liveMigrations/{liveMigrationId}/cutover"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.liveMigrationId == "" {
		return nil, reportError("liveMigrationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"liveMigrationId"+"}", url.PathEscape(r.liveMigrationId), -1)

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

type DeleteLinkTokensApiRequest struct {
	ctx        context.Context
	ApiService CloudMigrationServiceApi
	orgId      string
}

type DeleteLinkTokensApiParams struct {
	OrgId string
}

func (a *CloudMigrationServiceApiService) DeleteLinkTokensWithParams(ctx context.Context, args *DeleteLinkTokensApiParams) DeleteLinkTokensApiRequest {
	return DeleteLinkTokensApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r DeleteLinkTokensApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteLinkTokensExecute(r)
}

/*
DeleteLinkTokens Remove One Link-Token

Remove one organization link and its associated public API key. MongoDB Atlas uses the link-token for push live migrations only. Live migrations (push) let you securely push data from Cloud Manager or Ops Manager into MongoDB Atlas. Your API Key must have the Organization Owner role to successfully call this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return DeleteLinkTokensApiRequest
*/
func (a *CloudMigrationServiceApiService) DeleteLinkTokens(ctx context.Context, orgId string) DeleteLinkTokensApiRequest {
	return DeleteLinkTokensApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// DeleteLinkTokensExecute executes the request
func (a *CloudMigrationServiceApiService) DeleteLinkTokensExecute(r DeleteLinkTokensApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.DeleteLinkTokens")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/liveMigrations/linkTokens"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
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

type GetGroupLiveMigrationApiRequest struct {
	ctx             context.Context
	ApiService      CloudMigrationServiceApi
	groupId         string
	liveMigrationId string
}

type GetGroupLiveMigrationApiParams struct {
	GroupId         string
	LiveMigrationId string
}

func (a *CloudMigrationServiceApiService) GetGroupLiveMigrationWithParams(ctx context.Context, args *GetGroupLiveMigrationApiParams) GetGroupLiveMigrationApiRequest {
	return GetGroupLiveMigrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         args.GroupId,
		liveMigrationId: args.LiveMigrationId,
	}
}

func (r GetGroupLiveMigrationApiRequest) Execute() (*LiveMigrationResponse, *http.Response, error) {
	return r.ApiService.GetGroupLiveMigrationExecute(r)
}

/*
GetGroupLiveMigration Return One Migration Job

Return details of one cluster migration job. Each push live migration job uses one migration host. Your API Key must have the Organization Member role to successfully call this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param liveMigrationId Unique 24-hexadecimal digit string that identifies the migration.
	@return GetGroupLiveMigrationApiRequest
*/
func (a *CloudMigrationServiceApiService) GetGroupLiveMigration(ctx context.Context, groupId string, liveMigrationId string) GetGroupLiveMigrationApiRequest {
	return GetGroupLiveMigrationApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         groupId,
		liveMigrationId: liveMigrationId,
	}
}

// GetGroupLiveMigrationExecute executes the request
//
//	@return LiveMigrationResponse
func (a *CloudMigrationServiceApiService) GetGroupLiveMigrationExecute(r GetGroupLiveMigrationApiRequest) (*LiveMigrationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LiveMigrationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.GetGroupLiveMigration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/liveMigrations/{liveMigrationId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.liveMigrationId == "" {
		return localVarReturnValue, nil, reportError("liveMigrationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"liveMigrationId"+"}", url.PathEscape(r.liveMigrationId), -1)

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

type GetMigrationValidateStatusApiRequest struct {
	ctx          context.Context
	ApiService   CloudMigrationServiceApi
	groupId      string
	validationId string
}

type GetMigrationValidateStatusApiParams struct {
	GroupId      string
	ValidationId string
}

func (a *CloudMigrationServiceApiService) GetMigrationValidateStatusWithParams(ctx context.Context, args *GetMigrationValidateStatusApiParams) GetMigrationValidateStatusApiRequest {
	return GetMigrationValidateStatusApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		validationId: args.ValidationId,
	}
}

func (r GetMigrationValidateStatusApiRequest) Execute() (*LiveImportValidation, *http.Response, error) {
	return r.ApiService.GetMigrationValidateStatusExecute(r)
}

/*
GetMigrationValidateStatus Return One Migration Validation Job

Return the status of one migration validation job. Your API Key must have the Organization Owner role to successfully call this resource.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param validationId Unique 24-hexadecimal digit string that identifies the validation job.
	@return GetMigrationValidateStatusApiRequest
*/
func (a *CloudMigrationServiceApiService) GetMigrationValidateStatus(ctx context.Context, groupId string, validationId string) GetMigrationValidateStatusApiRequest {
	return GetMigrationValidateStatusApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		validationId: validationId,
	}
}

// GetMigrationValidateStatusExecute executes the request
//
//	@return LiveImportValidation
func (a *CloudMigrationServiceApiService) GetMigrationValidateStatusExecute(r GetMigrationValidateStatusApiRequest) (*LiveImportValidation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LiveImportValidation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.GetMigrationValidateStatus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/liveMigrations/validate/{validationId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.validationId == "" {
		return localVarReturnValue, nil, reportError("validationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"validationId"+"}", url.PathEscape(r.validationId), -1)

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

type ListAvailableProjectsApiRequest struct {
	ctx        context.Context
	ApiService CloudMigrationServiceApi
	orgId      string
}

type ListAvailableProjectsApiParams struct {
	OrgId string
}

func (a *CloudMigrationServiceApiService) ListAvailableProjectsWithParams(ctx context.Context, args *ListAvailableProjectsApiParams) ListAvailableProjectsApiRequest {
	return ListAvailableProjectsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r ListAvailableProjectsApiRequest) Execute() ([]LiveImportAvailableProject, *http.Response, error) {
	return r.ApiService.ListAvailableProjectsExecute(r)
}

/*
ListAvailableProjects Return All Projects Available for Migration

Return all projects that you can migrate to the specified organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListAvailableProjectsApiRequest
*/
func (a *CloudMigrationServiceApiService) ListAvailableProjects(ctx context.Context, orgId string) ListAvailableProjectsApiRequest {
	return ListAvailableProjectsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListAvailableProjectsExecute executes the request
//
//	@return []LiveImportAvailableProject
func (a *CloudMigrationServiceApiService) ListAvailableProjectsExecute(r ListAvailableProjectsApiRequest) ([]LiveImportAvailableProject, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []LiveImportAvailableProject
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.ListAvailableProjects")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/liveMigrations/availableProjects"
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

type ValidateLiveMigrationsApiRequest struct {
	ctx                          context.Context
	ApiService                   CloudMigrationServiceApi
	groupId                      string
	liveMigrationRequest20240530 *LiveMigrationRequest20240530
}

type ValidateLiveMigrationsApiParams struct {
	GroupId                      string
	LiveMigrationRequest20240530 *LiveMigrationRequest20240530
}

func (a *CloudMigrationServiceApiService) ValidateLiveMigrationsWithParams(ctx context.Context, args *ValidateLiveMigrationsApiParams) ValidateLiveMigrationsApiRequest {
	return ValidateLiveMigrationsApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		liveMigrationRequest20240530: args.LiveMigrationRequest20240530,
	}
}

func (r ValidateLiveMigrationsApiRequest) Execute() (*LiveImportValidation, *http.Response, error) {
	return r.ApiService.ValidateLiveMigrationsExecute(r)
}

/*
ValidateLiveMigrations Validate One Migration Request

Verifies whether the provided credentials, available disk space, MongoDB versions, and so on meet the requirements of the migration request. If the check passes, the migration can proceed. Your API Key must have the Organization Owner role to successfully call this resource. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ValidateLiveMigrationsApiRequest
*/
func (a *CloudMigrationServiceApiService) ValidateLiveMigrations(ctx context.Context, groupId string, liveMigrationRequest20240530 *LiveMigrationRequest20240530) ValidateLiveMigrationsApiRequest {
	return ValidateLiveMigrationsApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		liveMigrationRequest20240530: liveMigrationRequest20240530,
	}
}

// ValidateLiveMigrationsExecute executes the request
//
//	@return LiveImportValidation
func (a *CloudMigrationServiceApiService) ValidateLiveMigrationsExecute(r ValidateLiveMigrationsApiRequest) (*LiveImportValidation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LiveImportValidation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudMigrationServiceApiService.ValidateLiveMigrations")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/liveMigrations/validate"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.liveMigrationRequest20240530 == nil {
		return localVarReturnValue, nil, reportError("liveMigrationRequest20240530 is required and must be specified")
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
	localVarPostBody = r.liveMigrationRequest20240530
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
