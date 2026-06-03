// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DatabaseUsersApi interface {

	/*
		CreateDatabaseUser Create One Database User in One Project

		Creates one database user in the specified project. This MongoDB Cloud supports a maximum of 100 database users per project. If you require more than 100 database users on a project, contact Support.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudDatabaseUser Creates one database user in the specified project.
		@return CreateDatabaseUserApiRequest
	*/
	CreateDatabaseUser(ctx context.Context, groupId string, cloudDatabaseUser *CloudDatabaseUser) CreateDatabaseUserApiRequest
	/*
		CreateDatabaseUser Create One Database User in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateDatabaseUserApiParams - Parameters for the request
		@return CreateDatabaseUserApiRequest
	*/
	CreateDatabaseUserWithParams(ctx context.Context, args *CreateDatabaseUserApiParams) CreateDatabaseUserApiRequest

	// Method available only for mocking purposes
	CreateDatabaseUserExecute(r CreateDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error)

	/*
		DeleteDatabaseUser Remove One Database User from One Project

		Removes one database user from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
		@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
		@return DeleteDatabaseUserApiRequest
	*/
	DeleteDatabaseUser(ctx context.Context, groupId string, databaseName string, username string) DeleteDatabaseUserApiRequest
	/*
		DeleteDatabaseUser Remove One Database User from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteDatabaseUserApiParams - Parameters for the request
		@return DeleteDatabaseUserApiRequest
	*/
	DeleteDatabaseUserWithParams(ctx context.Context, args *DeleteDatabaseUserApiParams) DeleteDatabaseUserApiRequest

	// Method available only for mocking purposes
	DeleteDatabaseUserExecute(r DeleteDatabaseUserApiRequest) (*http.Response, error)

	/*
		GetDatabaseUser Return One Database User from One Project

		Returns one database user that belong to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
		@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
		@return GetDatabaseUserApiRequest
	*/
	GetDatabaseUser(ctx context.Context, groupId string, databaseName string, username string) GetDatabaseUserApiRequest
	/*
		GetDatabaseUser Return One Database User from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetDatabaseUserApiParams - Parameters for the request
		@return GetDatabaseUserApiRequest
	*/
	GetDatabaseUserWithParams(ctx context.Context, args *GetDatabaseUserApiParams) GetDatabaseUserApiRequest

	// Method available only for mocking purposes
	GetDatabaseUserExecute(r GetDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error)

	/*
		ListDatabaseUsers Return All Database Users in One Project

		Returns all database users that belong to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListDatabaseUsersApiRequest
	*/
	ListDatabaseUsers(ctx context.Context, groupId string) ListDatabaseUsersApiRequest
	/*
		ListDatabaseUsers Return All Database Users in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDatabaseUsersApiParams - Parameters for the request
		@return ListDatabaseUsersApiRequest
	*/
	ListDatabaseUsersWithParams(ctx context.Context, args *ListDatabaseUsersApiParams) ListDatabaseUsersApiRequest

	// Method available only for mocking purposes
	ListDatabaseUsersExecute(r ListDatabaseUsersApiRequest) (*PaginatedApiAtlasDatabaseUser, *http.Response, error)

	/*
		UpdateDatabaseUser Update One Database User in One Project

		Updates one database user that belongs to the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
		@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
		@param cloudDatabaseUser Updates one database user that belongs to the specified project.
		@return UpdateDatabaseUserApiRequest
	*/
	UpdateDatabaseUser(ctx context.Context, groupId string, databaseName string, username string, cloudDatabaseUser *CloudDatabaseUser) UpdateDatabaseUserApiRequest
	/*
		UpdateDatabaseUser Update One Database User in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateDatabaseUserApiParams - Parameters for the request
		@return UpdateDatabaseUserApiRequest
	*/
	UpdateDatabaseUserWithParams(ctx context.Context, args *UpdateDatabaseUserApiParams) UpdateDatabaseUserApiRequest

	// Method available only for mocking purposes
	UpdateDatabaseUserExecute(r UpdateDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error)
}

// DatabaseUsersApiService DatabaseUsersApi service
type DatabaseUsersApiService service

type CreateDatabaseUserApiRequest struct {
	ctx               context.Context
	ApiService        DatabaseUsersApi
	groupId           string
	cloudDatabaseUser *CloudDatabaseUser
}

type CreateDatabaseUserApiParams struct {
	GroupId           string
	CloudDatabaseUser *CloudDatabaseUser
}

func (a *DatabaseUsersApiService) CreateDatabaseUserWithParams(ctx context.Context, args *CreateDatabaseUserApiParams) CreateDatabaseUserApiRequest {
	return CreateDatabaseUserApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		cloudDatabaseUser: args.CloudDatabaseUser,
	}
}

func (r CreateDatabaseUserApiRequest) Execute() (*CloudDatabaseUser, *http.Response, error) {
	return r.ApiService.CreateDatabaseUserExecute(r)
}

/*
CreateDatabaseUser Create One Database User in One Project

Creates one database user in the specified project. This MongoDB Cloud supports a maximum of 100 database users per project. If you require more than 100 database users on a project, contact Support.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateDatabaseUserApiRequest
*/
func (a *DatabaseUsersApiService) CreateDatabaseUser(ctx context.Context, groupId string, cloudDatabaseUser *CloudDatabaseUser) CreateDatabaseUserApiRequest {
	return CreateDatabaseUserApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		cloudDatabaseUser: cloudDatabaseUser,
	}
}

// CreateDatabaseUserExecute executes the request
//
//	@return CloudDatabaseUser
func (a *DatabaseUsersApiService) CreateDatabaseUserExecute(r CreateDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudDatabaseUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseUsersApiService.CreateDatabaseUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudDatabaseUser == nil {
		return localVarReturnValue, nil, reportError("cloudDatabaseUser is required and must be specified")
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
	localVarPostBody = r.cloudDatabaseUser
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

type DeleteDatabaseUserApiRequest struct {
	ctx          context.Context
	ApiService   DatabaseUsersApi
	groupId      string
	databaseName string
	username     string
}

type DeleteDatabaseUserApiParams struct {
	GroupId      string
	DatabaseName string
	Username     string
}

func (a *DatabaseUsersApiService) DeleteDatabaseUserWithParams(ctx context.Context, args *DeleteDatabaseUserApiParams) DeleteDatabaseUserApiRequest {
	return DeleteDatabaseUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		databaseName: args.DatabaseName,
		username:     args.Username,
	}
}

func (r DeleteDatabaseUserApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteDatabaseUserExecute(r)
}

/*
DeleteDatabaseUser Remove One Database User from One Project

Removes one database user from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
	@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
	@return DeleteDatabaseUserApiRequest
*/
func (a *DatabaseUsersApiService) DeleteDatabaseUser(ctx context.Context, groupId string, databaseName string, username string) DeleteDatabaseUserApiRequest {
	return DeleteDatabaseUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		databaseName: databaseName,
		username:     username,
	}
}

// DeleteDatabaseUserExecute executes the request
func (a *DatabaseUsersApiService) DeleteDatabaseUserExecute(r DeleteDatabaseUserApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseUsersApiService.DeleteDatabaseUser")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers/{databaseName}/{username}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.databaseName == "" {
		return nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.username == "" {
		return nil, reportError("username is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(r.username), -1)

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

type GetDatabaseUserApiRequest struct {
	ctx          context.Context
	ApiService   DatabaseUsersApi
	groupId      string
	databaseName string
	username     string
}

type GetDatabaseUserApiParams struct {
	GroupId      string
	DatabaseName string
	Username     string
}

func (a *DatabaseUsersApiService) GetDatabaseUserWithParams(ctx context.Context, args *GetDatabaseUserApiParams) GetDatabaseUserApiRequest {
	return GetDatabaseUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		databaseName: args.DatabaseName,
		username:     args.Username,
	}
}

func (r GetDatabaseUserApiRequest) Execute() (*CloudDatabaseUser, *http.Response, error) {
	return r.ApiService.GetDatabaseUserExecute(r)
}

/*
GetDatabaseUser Return One Database User from One Project

Returns one database user that belong to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
	@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
	@return GetDatabaseUserApiRequest
*/
func (a *DatabaseUsersApiService) GetDatabaseUser(ctx context.Context, groupId string, databaseName string, username string) GetDatabaseUserApiRequest {
	return GetDatabaseUserApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		databaseName: databaseName,
		username:     username,
	}
}

// GetDatabaseUserExecute executes the request
//
//	@return CloudDatabaseUser
func (a *DatabaseUsersApiService) GetDatabaseUserExecute(r GetDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudDatabaseUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseUsersApiService.GetDatabaseUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers/{databaseName}/{username}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.username == "" {
		return localVarReturnValue, nil, reportError("username is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(r.username), -1)

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

type ListDatabaseUsersApiRequest struct {
	ctx          context.Context
	ApiService   DatabaseUsersApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListDatabaseUsersApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *DatabaseUsersApiService) ListDatabaseUsersWithParams(ctx context.Context, args *ListDatabaseUsersApiParams) ListDatabaseUsersApiRequest {
	return ListDatabaseUsersApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListDatabaseUsersApiRequest) IncludeCount(includeCount bool) ListDatabaseUsersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListDatabaseUsersApiRequest) ItemsPerPage(itemsPerPage int) ListDatabaseUsersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListDatabaseUsersApiRequest) PageNum(pageNum int) ListDatabaseUsersApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListDatabaseUsersApiRequest) Execute() (*PaginatedApiAtlasDatabaseUser, *http.Response, error) {
	return r.ApiService.ListDatabaseUsersExecute(r)
}

/*
ListDatabaseUsers Return All Database Users in One Project

Returns all database users that belong to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListDatabaseUsersApiRequest
*/
func (a *DatabaseUsersApiService) ListDatabaseUsers(ctx context.Context, groupId string) ListDatabaseUsersApiRequest {
	return ListDatabaseUsersApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListDatabaseUsersExecute executes the request
//
//	@return PaginatedApiAtlasDatabaseUser
func (a *DatabaseUsersApiService) ListDatabaseUsersExecute(r ListDatabaseUsersApiRequest) (*PaginatedApiAtlasDatabaseUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasDatabaseUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseUsersApiService.ListDatabaseUsers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers"
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

type UpdateDatabaseUserApiRequest struct {
	ctx               context.Context
	ApiService        DatabaseUsersApi
	groupId           string
	databaseName      string
	username          string
	cloudDatabaseUser *CloudDatabaseUser
}

type UpdateDatabaseUserApiParams struct {
	GroupId           string
	DatabaseName      string
	Username          string
	CloudDatabaseUser *CloudDatabaseUser
}

func (a *DatabaseUsersApiService) UpdateDatabaseUserWithParams(ctx context.Context, args *UpdateDatabaseUserApiParams) UpdateDatabaseUserApiRequest {
	return UpdateDatabaseUserApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		databaseName:      args.DatabaseName,
		username:          args.Username,
		cloudDatabaseUser: args.CloudDatabaseUser,
	}
}

func (r UpdateDatabaseUserApiRequest) Execute() (*CloudDatabaseUser, *http.Response, error) {
	return r.ApiService.UpdateDatabaseUserExecute(r)
}

/*
UpdateDatabaseUser Update One Database User in One Project

Updates one database user that belongs to the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param databaseName The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
	@param username Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
	@return UpdateDatabaseUserApiRequest
*/
func (a *DatabaseUsersApiService) UpdateDatabaseUser(ctx context.Context, groupId string, databaseName string, username string, cloudDatabaseUser *CloudDatabaseUser) UpdateDatabaseUserApiRequest {
	return UpdateDatabaseUserApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		databaseName:      databaseName,
		username:          username,
		cloudDatabaseUser: cloudDatabaseUser,
	}
}

// UpdateDatabaseUserExecute executes the request
//
//	@return CloudDatabaseUser
func (a *DatabaseUsersApiService) UpdateDatabaseUserExecute(r UpdateDatabaseUserApiRequest) (*CloudDatabaseUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudDatabaseUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseUsersApiService.UpdateDatabaseUser")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers/{databaseName}/{username}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.username == "" {
		return localVarReturnValue, nil, reportError("username is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(r.username), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudDatabaseUser == nil {
		return localVarReturnValue, nil, reportError("cloudDatabaseUser is required and must be specified")
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
	localVarPostBody = r.cloudDatabaseUser
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
