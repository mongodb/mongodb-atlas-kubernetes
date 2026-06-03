// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OnlineArchiveApi interface {

	/*
		CreateOnlineArchive Create One Online Archive

		Creates one online archive. This archive stores data from one cluster within one project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to create one online archive.
		@param backupOnlineArchiveCreate Creates one online archive.
		@return CreateOnlineArchiveApiRequest
	*/
	CreateOnlineArchive(ctx context.Context, groupId string, clusterName string, backupOnlineArchiveCreate *BackupOnlineArchiveCreate) CreateOnlineArchiveApiRequest
	/*
		CreateOnlineArchive Create One Online Archive


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOnlineArchiveApiParams - Parameters for the request
		@return CreateOnlineArchiveApiRequest
	*/
	CreateOnlineArchiveWithParams(ctx context.Context, args *CreateOnlineArchiveApiParams) CreateOnlineArchiveApiRequest

	// Method available only for mocking purposes
	CreateOnlineArchiveExecute(r CreateOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error)

	/*
		DeleteOnlineArchive Remove One Online Archive

		Removes one online archive. This archive stores data from one cluster within one project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to delete.
		@param clusterName Human-readable label that identifies the cluster that contains the collection from which you want to remove an online archive.
		@return DeleteOnlineArchiveApiRequest
	*/
	DeleteOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string) DeleteOnlineArchiveApiRequest
	/*
		DeleteOnlineArchive Remove One Online Archive


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOnlineArchiveApiParams - Parameters for the request
		@return DeleteOnlineArchiveApiRequest
	*/
	DeleteOnlineArchiveWithParams(ctx context.Context, args *DeleteOnlineArchiveApiParams) DeleteOnlineArchiveApiRequest

	// Method available only for mocking purposes
	DeleteOnlineArchiveExecute(r DeleteOnlineArchiveApiRequest) (*http.Response, error)

	/*
		DownloadQueryLogs Download Online Archive Query Logs

		Downloads query logs for the specified online archive. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to return the query logs from one online archive.
		@return DownloadQueryLogsApiRequest
	*/
	DownloadQueryLogs(ctx context.Context, groupId string, clusterName string) DownloadQueryLogsApiRequest
	/*
		DownloadQueryLogs Download Online Archive Query Logs


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadQueryLogsApiParams - Parameters for the request
		@return DownloadQueryLogsApiRequest
	*/
	DownloadQueryLogsWithParams(ctx context.Context, args *DownloadQueryLogsApiParams) DownloadQueryLogsApiRequest

	// Method available only for mocking purposes
	DownloadQueryLogsExecute(r DownloadQueryLogsApiRequest) (io.ReadCloser, *http.Response, error)

	/*
		GetOnlineArchive Return One Online Archive

		Returns one online archive for one cluster. This archive stores data from one cluster within one project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to return.
		@param clusterName Human-readable label that identifies the cluster that contains the specified collection from which Application created the online archive.
		@return GetOnlineArchiveApiRequest
	*/
	GetOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string) GetOnlineArchiveApiRequest
	/*
		GetOnlineArchive Return One Online Archive


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOnlineArchiveApiParams - Parameters for the request
		@return GetOnlineArchiveApiRequest
	*/
	GetOnlineArchiveWithParams(ctx context.Context, args *GetOnlineArchiveApiParams) GetOnlineArchiveApiRequest

	// Method available only for mocking purposes
	GetOnlineArchiveExecute(r GetOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error)

	/*
		ListOnlineArchives Return All Online Archives for One Cluster

		Returns details of all online archives. This archive stores data from one cluster within one project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to return the online archives.
		@return ListOnlineArchivesApiRequest
	*/
	ListOnlineArchives(ctx context.Context, groupId string, clusterName string) ListOnlineArchivesApiRequest
	/*
		ListOnlineArchives Return All Online Archives for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOnlineArchivesApiParams - Parameters for the request
		@return ListOnlineArchivesApiRequest
	*/
	ListOnlineArchivesWithParams(ctx context.Context, args *ListOnlineArchivesApiParams) ListOnlineArchivesApiRequest

	// Method available only for mocking purposes
	ListOnlineArchivesExecute(r ListOnlineArchivesApiRequest) (*PaginatedOnlineArchive, *http.Response, error)

	/*
		UpdateOnlineArchive Update One Online Archive

		Updates, pauses, or resumes one online archive. This archive stores data from one cluster within one project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to update.
		@param clusterName Human-readable label that identifies the cluster that contains the specified collection from which Application created the online archive.
		@param backupOnlineArchive Updates, pauses, or resumes one online archive.
		@return UpdateOnlineArchiveApiRequest
	*/
	UpdateOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string, backupOnlineArchive *BackupOnlineArchive) UpdateOnlineArchiveApiRequest
	/*
		UpdateOnlineArchive Update One Online Archive


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOnlineArchiveApiParams - Parameters for the request
		@return UpdateOnlineArchiveApiRequest
	*/
	UpdateOnlineArchiveWithParams(ctx context.Context, args *UpdateOnlineArchiveApiParams) UpdateOnlineArchiveApiRequest

	// Method available only for mocking purposes
	UpdateOnlineArchiveExecute(r UpdateOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error)
}

// OnlineArchiveApiService OnlineArchiveApi service
type OnlineArchiveApiService service

type CreateOnlineArchiveApiRequest struct {
	ctx                       context.Context
	ApiService                OnlineArchiveApi
	groupId                   string
	clusterName               string
	backupOnlineArchiveCreate *BackupOnlineArchiveCreate
}

type CreateOnlineArchiveApiParams struct {
	GroupId                   string
	ClusterName               string
	BackupOnlineArchiveCreate *BackupOnlineArchiveCreate
}

func (a *OnlineArchiveApiService) CreateOnlineArchiveWithParams(ctx context.Context, args *CreateOnlineArchiveApiParams) CreateOnlineArchiveApiRequest {
	return CreateOnlineArchiveApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   args.GroupId,
		clusterName:               args.ClusterName,
		backupOnlineArchiveCreate: args.BackupOnlineArchiveCreate,
	}
}

func (r CreateOnlineArchiveApiRequest) Execute() (*BackupOnlineArchive, *http.Response, error) {
	return r.ApiService.CreateOnlineArchiveExecute(r)
}

/*
CreateOnlineArchive Create One Online Archive

Creates one online archive. This archive stores data from one cluster within one project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to create one online archive.
	@return CreateOnlineArchiveApiRequest
*/
func (a *OnlineArchiveApiService) CreateOnlineArchive(ctx context.Context, groupId string, clusterName string, backupOnlineArchiveCreate *BackupOnlineArchiveCreate) CreateOnlineArchiveApiRequest {
	return CreateOnlineArchiveApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   groupId,
		clusterName:               clusterName,
		backupOnlineArchiveCreate: backupOnlineArchiveCreate,
	}
}

// CreateOnlineArchiveExecute executes the request
//
//	@return BackupOnlineArchive
func (a *OnlineArchiveApiService) CreateOnlineArchiveExecute(r CreateOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupOnlineArchive
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.CreateOnlineArchive")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.backupOnlineArchiveCreate == nil {
		return localVarReturnValue, nil, reportError("backupOnlineArchiveCreate is required and must be specified")
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
	localVarPostBody = r.backupOnlineArchiveCreate
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

type DeleteOnlineArchiveApiRequest struct {
	ctx         context.Context
	ApiService  OnlineArchiveApi
	groupId     string
	archiveId   string
	clusterName string
}

type DeleteOnlineArchiveApiParams struct {
	GroupId     string
	ArchiveId   string
	ClusterName string
}

func (a *OnlineArchiveApiService) DeleteOnlineArchiveWithParams(ctx context.Context, args *DeleteOnlineArchiveApiParams) DeleteOnlineArchiveApiRequest {
	return DeleteOnlineArchiveApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		archiveId:   args.ArchiveId,
		clusterName: args.ClusterName,
	}
}

func (r DeleteOnlineArchiveApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOnlineArchiveExecute(r)
}

/*
DeleteOnlineArchive Remove One Online Archive

Removes one online archive. This archive stores data from one cluster within one project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to delete.
	@param clusterName Human-readable label that identifies the cluster that contains the collection from which you want to remove an online archive.
	@return DeleteOnlineArchiveApiRequest
*/
func (a *OnlineArchiveApiService) DeleteOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string) DeleteOnlineArchiveApiRequest {
	return DeleteOnlineArchiveApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		archiveId:   archiveId,
		clusterName: clusterName,
	}
}

// DeleteOnlineArchiveExecute executes the request
func (a *OnlineArchiveApiService) DeleteOnlineArchiveExecute(r DeleteOnlineArchiveApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.DeleteOnlineArchive")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives/{archiveId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.archiveId == "" {
		return nil, reportError("archiveId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"archiveId"+"}", url.PathEscape(r.archiveId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

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

type DownloadQueryLogsApiRequest struct {
	ctx         context.Context
	ApiService  OnlineArchiveApi
	groupId     string
	clusterName string
	startDate   *int64
	endDate     *int64
	archiveOnly *bool
}

type DownloadQueryLogsApiParams struct {
	GroupId     string
	ClusterName string
	StartDate   *int64
	EndDate     *int64
	ArchiveOnly *bool
}

func (a *OnlineArchiveApiService) DownloadQueryLogsWithParams(ctx context.Context, args *DownloadQueryLogsApiParams) DownloadQueryLogsApiRequest {
	return DownloadQueryLogsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		startDate:   args.StartDate,
		endDate:     args.EndDate,
		archiveOnly: args.ArchiveOnly,
	}
}

// Date and time that specifies the starting point for the range of log messages to return. This resource expresses this value in the number of seconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).
func (r DownloadQueryLogsApiRequest) StartDate(startDate int64) DownloadQueryLogsApiRequest {
	r.startDate = &startDate
	return r
}

// Date and time that specifies the end point for the range of log messages to return. This resource expresses this value in the number of seconds that have elapsed since the [UNIX epoch](https://en.wikipedia.org/wiki/Unix_time).
func (r DownloadQueryLogsApiRequest) EndDate(endDate int64) DownloadQueryLogsApiRequest {
	r.endDate = &endDate
	return r
}

// Flag that indicates whether to download logs for queries against your online archive only or both your online archive and cluster.
func (r DownloadQueryLogsApiRequest) ArchiveOnly(archiveOnly bool) DownloadQueryLogsApiRequest {
	r.archiveOnly = &archiveOnly
	return r
}

func (r DownloadQueryLogsApiRequest) Execute() (io.ReadCloser, *http.Response, error) {
	return r.ApiService.DownloadQueryLogsExecute(r)
}

/*
DownloadQueryLogs Download Online Archive Query Logs

Downloads query logs for the specified online archive. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to return the query logs from one online archive.
	@return DownloadQueryLogsApiRequest
*/
func (a *OnlineArchiveApiService) DownloadQueryLogs(ctx context.Context, groupId string, clusterName string) DownloadQueryLogsApiRequest {
	return DownloadQueryLogsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// DownloadQueryLogsExecute executes the request
//
//	@return io.ReadCloser
func (a *OnlineArchiveApiService) DownloadQueryLogsExecute(r DownloadQueryLogsApiRequest) (io.ReadCloser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue io.ReadCloser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.DownloadQueryLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives/queryLogs.gz"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.startDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startDate", r.startDate, "")
	}
	if r.endDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endDate", r.endDate, "")
	}
	if r.archiveOnly != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "archiveOnly", r.archiveOnly, "")
	} else {
		var defaultValue bool = false
		r.archiveOnly = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "archiveOnly", r.archiveOnly, "")
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

type GetOnlineArchiveApiRequest struct {
	ctx         context.Context
	ApiService  OnlineArchiveApi
	groupId     string
	archiveId   string
	clusterName string
}

type GetOnlineArchiveApiParams struct {
	GroupId     string
	ArchiveId   string
	ClusterName string
}

func (a *OnlineArchiveApiService) GetOnlineArchiveWithParams(ctx context.Context, args *GetOnlineArchiveApiParams) GetOnlineArchiveApiRequest {
	return GetOnlineArchiveApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		archiveId:   args.ArchiveId,
		clusterName: args.ClusterName,
	}
}

func (r GetOnlineArchiveApiRequest) Execute() (*BackupOnlineArchive, *http.Response, error) {
	return r.ApiService.GetOnlineArchiveExecute(r)
}

/*
GetOnlineArchive Return One Online Archive

Returns one online archive for one cluster. This archive stores data from one cluster within one project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to return.
	@param clusterName Human-readable label that identifies the cluster that contains the specified collection from which Application created the online archive.
	@return GetOnlineArchiveApiRequest
*/
func (a *OnlineArchiveApiService) GetOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string) GetOnlineArchiveApiRequest {
	return GetOnlineArchiveApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		archiveId:   archiveId,
		clusterName: clusterName,
	}
}

// GetOnlineArchiveExecute executes the request
//
//	@return BackupOnlineArchive
func (a *OnlineArchiveApiService) GetOnlineArchiveExecute(r GetOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupOnlineArchive
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.GetOnlineArchive")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives/{archiveId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.archiveId == "" {
		return localVarReturnValue, nil, reportError("archiveId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"archiveId"+"}", url.PathEscape(r.archiveId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

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

type ListOnlineArchivesApiRequest struct {
	ctx          context.Context
	ApiService   OnlineArchiveApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListOnlineArchivesApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *OnlineArchiveApiService) ListOnlineArchivesWithParams(ctx context.Context, args *ListOnlineArchivesApiParams) ListOnlineArchivesApiRequest {
	return ListOnlineArchivesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOnlineArchivesApiRequest) IncludeCount(includeCount bool) ListOnlineArchivesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOnlineArchivesApiRequest) ItemsPerPage(itemsPerPage int) ListOnlineArchivesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOnlineArchivesApiRequest) PageNum(pageNum int) ListOnlineArchivesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListOnlineArchivesApiRequest) Execute() (*PaginatedOnlineArchive, *http.Response, error) {
	return r.ApiService.ListOnlineArchivesExecute(r)
}

/*
ListOnlineArchives Return All Online Archives for One Cluster

Returns details of all online archives. This archive stores data from one cluster within one project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that contains the collection for which you want to return the online archives.
	@return ListOnlineArchivesApiRequest
*/
func (a *OnlineArchiveApiService) ListOnlineArchives(ctx context.Context, groupId string, clusterName string) ListOnlineArchivesApiRequest {
	return ListOnlineArchivesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListOnlineArchivesExecute executes the request
//
//	@return PaginatedOnlineArchive
func (a *OnlineArchiveApiService) ListOnlineArchivesExecute(r ListOnlineArchivesApiRequest) (*PaginatedOnlineArchive, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOnlineArchive
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.ListOnlineArchives")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

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

type UpdateOnlineArchiveApiRequest struct {
	ctx                 context.Context
	ApiService          OnlineArchiveApi
	groupId             string
	archiveId           string
	clusterName         string
	backupOnlineArchive *BackupOnlineArchive
}

type UpdateOnlineArchiveApiParams struct {
	GroupId             string
	ArchiveId           string
	ClusterName         string
	BackupOnlineArchive *BackupOnlineArchive
}

func (a *OnlineArchiveApiService) UpdateOnlineArchiveWithParams(ctx context.Context, args *UpdateOnlineArchiveApiParams) UpdateOnlineArchiveApiRequest {
	return UpdateOnlineArchiveApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             args.GroupId,
		archiveId:           args.ArchiveId,
		clusterName:         args.ClusterName,
		backupOnlineArchive: args.BackupOnlineArchive,
	}
}

func (r UpdateOnlineArchiveApiRequest) Execute() (*BackupOnlineArchive, *http.Response, error) {
	return r.ApiService.UpdateOnlineArchiveExecute(r)
}

/*
UpdateOnlineArchive Update One Online Archive

Updates, pauses, or resumes one online archive. This archive stores data from one cluster within one project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param archiveId Unique 24-hexadecimal digit string that identifies the online archive to update.
	@param clusterName Human-readable label that identifies the cluster that contains the specified collection from which Application created the online archive.
	@return UpdateOnlineArchiveApiRequest
*/
func (a *OnlineArchiveApiService) UpdateOnlineArchive(ctx context.Context, groupId string, archiveId string, clusterName string, backupOnlineArchive *BackupOnlineArchive) UpdateOnlineArchiveApiRequest {
	return UpdateOnlineArchiveApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             groupId,
		archiveId:           archiveId,
		clusterName:         clusterName,
		backupOnlineArchive: backupOnlineArchive,
	}
}

// UpdateOnlineArchiveExecute executes the request
//
//	@return BackupOnlineArchive
func (a *OnlineArchiveApiService) UpdateOnlineArchiveExecute(r UpdateOnlineArchiveApiRequest) (*BackupOnlineArchive, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupOnlineArchive
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OnlineArchiveApiService.UpdateOnlineArchive")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/onlineArchives/{archiveId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.archiveId == "" {
		return localVarReturnValue, nil, reportError("archiveId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"archiveId"+"}", url.PathEscape(r.archiveId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.backupOnlineArchive == nil {
		return localVarReturnValue, nil, reportError("backupOnlineArchive is required and must be specified")
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
	localVarPostBody = r.backupOnlineArchive
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
