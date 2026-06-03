// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type X509AuthenticationApi interface {

	/*
			CreateDatabaseUserCert Create One X.509 Certificate for One Database User

			Generates one X.509 certificate for the specified MongoDB user. Atlas manages the            certificate and MongoDB user that belong to one project.


		To get MongoDB Cloud to generate a managed certificate for a database user, set `"x509Type" : "MANAGED"` on the desired MongoDB Database User.

		If you are managing your own Certificate Authority (CA) in Self-Managed X.509 mode, you must generate certificates for database users using your own CA.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param username Human-readable label that represents the MongoDB database user account for whom to create a certificate.
			@param userCert Generates one X.509 certificate for the specified MongoDB user.
			@return CreateDatabaseUserCertApiRequest
	*/
	CreateDatabaseUserCert(ctx context.Context, groupId string, username string, userCert *UserCert) CreateDatabaseUserCertApiRequest
	/*
		CreateDatabaseUserCert Create One X.509 Certificate for One Database User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateDatabaseUserCertApiParams - Parameters for the request
		@return CreateDatabaseUserCertApiRequest
	*/
	CreateDatabaseUserCertWithParams(ctx context.Context, args *CreateDatabaseUserCertApiParams) CreateDatabaseUserCertApiRequest

	// Method available only for mocking purposes
	CreateDatabaseUserCertExecute(r CreateDatabaseUserCertApiRequest) (string, *http.Response, error)

	/*
			DisableSecurityCustomerX509 Disable Customer-Managed X.509

			Clears the customer-managed X.509 settings on a project, including the uploaded Certificate Authority, which disables self-managed X.509.

		 Updating this configuration triggers a rolling restart of the database. You must have the Project Owner role to use this endpoint.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return DisableSecurityCustomerX509ApiRequest
	*/
	DisableSecurityCustomerX509(ctx context.Context, groupId string) DisableSecurityCustomerX509ApiRequest
	/*
		DisableSecurityCustomerX509 Disable Customer-Managed X.509


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DisableSecurityCustomerX509ApiParams - Parameters for the request
		@return DisableSecurityCustomerX509ApiRequest
	*/
	DisableSecurityCustomerX509WithParams(ctx context.Context, args *DisableSecurityCustomerX509ApiParams) DisableSecurityCustomerX509ApiRequest

	// Method available only for mocking purposes
	DisableSecurityCustomerX509Execute(r DisableSecurityCustomerX509ApiRequest) (*UserSecurity, *http.Response, error)

	/*
		ListDatabaseUserCerts Return All X.509 Certificates Assigned to One Database User

		Returns all unexpired X.509 certificates for the specified MongoDB user. This MongoDB user belongs to one project. Atlas manages these certificates and the MongoDB user.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param username Human-readable label that represents the MongoDB database user account whose certificates you want to return.
		@return ListDatabaseUserCertsApiRequest
	*/
	ListDatabaseUserCerts(ctx context.Context, groupId string, username string) ListDatabaseUserCertsApiRequest
	/*
		ListDatabaseUserCerts Return All X.509 Certificates Assigned to One Database User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDatabaseUserCertsApiParams - Parameters for the request
		@return ListDatabaseUserCertsApiRequest
	*/
	ListDatabaseUserCertsWithParams(ctx context.Context, args *ListDatabaseUserCertsApiParams) ListDatabaseUserCertsApiRequest

	// Method available only for mocking purposes
	ListDatabaseUserCertsExecute(r ListDatabaseUserCertsApiRequest) (*PaginatedUserCert, *http.Response, error)
}

// X509AuthenticationApiService X509AuthenticationApi service
type X509AuthenticationApiService service

type CreateDatabaseUserCertApiRequest struct {
	ctx        context.Context
	ApiService X509AuthenticationApi
	groupId    string
	username   string
	userCert   *UserCert
}

type CreateDatabaseUserCertApiParams struct {
	GroupId  string
	Username string
	UserCert *UserCert
}

func (a *X509AuthenticationApiService) CreateDatabaseUserCertWithParams(ctx context.Context, args *CreateDatabaseUserCertApiParams) CreateDatabaseUserCertApiRequest {
	return CreateDatabaseUserCertApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		username:   args.Username,
		userCert:   args.UserCert,
	}
}

func (r CreateDatabaseUserCertApiRequest) Execute() (string, *http.Response, error) {
	return r.ApiService.CreateDatabaseUserCertExecute(r)
}

/*
CreateDatabaseUserCert Create One X.509 Certificate for One Database User

Generates one X.509 certificate for the specified MongoDB user. Atlas manages the            certificate and MongoDB user that belong to one project.

To get MongoDB Cloud to generate a managed certificate for a database user, set `"x509Type" : "MANAGED"` on the desired MongoDB Database User.

If you are managing your own Certificate Authority (CA) in Self-Managed X.509 mode, you must generate certificates for database users using your own CA.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param username Human-readable label that represents the MongoDB database user account for whom to create a certificate.
	@return CreateDatabaseUserCertApiRequest
*/
func (a *X509AuthenticationApiService) CreateDatabaseUserCert(ctx context.Context, groupId string, username string, userCert *UserCert) CreateDatabaseUserCertApiRequest {
	return CreateDatabaseUserCertApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		username:   username,
		userCert:   userCert,
	}
}

// CreateDatabaseUserCertExecute executes the request
//
//	@return string
func (a *X509AuthenticationApiService) CreateDatabaseUserCertExecute(r CreateDatabaseUserCertApiRequest) (string, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "X509AuthenticationApiService.CreateDatabaseUserCert")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers/{username}/certs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.username == "" {
		return localVarReturnValue, nil, reportError("username is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(r.username), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userCert == nil {
		return localVarReturnValue, nil, reportError("userCert is required and must be specified")
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
	localVarPostBody = r.userCert
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

type DisableSecurityCustomerX509ApiRequest struct {
	ctx        context.Context
	ApiService X509AuthenticationApi
	groupId    string
}

type DisableSecurityCustomerX509ApiParams struct {
	GroupId string
}

func (a *X509AuthenticationApiService) DisableSecurityCustomerX509WithParams(ctx context.Context, args *DisableSecurityCustomerX509ApiParams) DisableSecurityCustomerX509ApiRequest {
	return DisableSecurityCustomerX509ApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DisableSecurityCustomerX509ApiRequest) Execute() (*UserSecurity, *http.Response, error) {
	return r.ApiService.DisableSecurityCustomerX509Execute(r)
}

/*
DisableSecurityCustomerX509 Disable Customer-Managed X.509

Clears the customer-managed X.509 settings on a project, including the uploaded Certificate Authority, which disables self-managed X.509.

	Updating this configuration triggers a rolling restart of the database. You must have the Project Owner role to use this endpoint.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DisableSecurityCustomerX509ApiRequest
*/
func (a *X509AuthenticationApiService) DisableSecurityCustomerX509(ctx context.Context, groupId string) DisableSecurityCustomerX509ApiRequest {
	return DisableSecurityCustomerX509ApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DisableSecurityCustomerX509Execute executes the request
//
//	@return UserSecurity
func (a *X509AuthenticationApiService) DisableSecurityCustomerX509Execute(r DisableSecurityCustomerX509ApiRequest) (*UserSecurity, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UserSecurity
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "X509AuthenticationApiService.DisableSecurityCustomerX509")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/userSecurity/customerX509"
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

type ListDatabaseUserCertsApiRequest struct {
	ctx          context.Context
	ApiService   X509AuthenticationApi
	groupId      string
	username     string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListDatabaseUserCertsApiParams struct {
	GroupId      string
	Username     string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *X509AuthenticationApiService) ListDatabaseUserCertsWithParams(ctx context.Context, args *ListDatabaseUserCertsApiParams) ListDatabaseUserCertsApiRequest {
	return ListDatabaseUserCertsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		username:     args.Username,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListDatabaseUserCertsApiRequest) IncludeCount(includeCount bool) ListDatabaseUserCertsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListDatabaseUserCertsApiRequest) ItemsPerPage(itemsPerPage int) ListDatabaseUserCertsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListDatabaseUserCertsApiRequest) PageNum(pageNum int) ListDatabaseUserCertsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListDatabaseUserCertsApiRequest) Execute() (*PaginatedUserCert, *http.Response, error) {
	return r.ApiService.ListDatabaseUserCertsExecute(r)
}

/*
ListDatabaseUserCerts Return All X.509 Certificates Assigned to One Database User

Returns all unexpired X.509 certificates for the specified MongoDB user. This MongoDB user belongs to one project. Atlas manages these certificates and the MongoDB user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param username Human-readable label that represents the MongoDB database user account whose certificates you want to return.
	@return ListDatabaseUserCertsApiRequest
*/
func (a *X509AuthenticationApiService) ListDatabaseUserCerts(ctx context.Context, groupId string, username string) ListDatabaseUserCertsApiRequest {
	return ListDatabaseUserCertsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		username:   username,
	}
}

// ListDatabaseUserCertsExecute executes the request
//
//	@return PaginatedUserCert
func (a *X509AuthenticationApiService) ListDatabaseUserCertsExecute(r ListDatabaseUserCertsApiRequest) (*PaginatedUserCert, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedUserCert
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "X509AuthenticationApiService.ListDatabaseUserCerts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/databaseUsers/{username}/certs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.username == "" {
		return localVarReturnValue, nil, reportError("username is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"username"+"}", url.PathEscape(r.username), -1)

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
