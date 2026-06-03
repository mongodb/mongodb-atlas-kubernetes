// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type FlexSnapshotsApi interface {

	/*
		DownloadFlexBackup Download One Flex Cluster Snapshot

		Requests one snapshot for the specified flex cluster. This resource returns a `snapshotURL` that you can use to download the snapshot. This `snapshotURL` remains active for four hours after you make the request.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param name Human-readable label that identifies the flex cluster.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param flexBackupSnapshotDownloadCreate20241113 Snapshot to be downloaded.
		@return DownloadFlexBackupApiRequest
	*/
	DownloadFlexBackup(ctx context.Context, name string, groupId string, flexBackupSnapshotDownloadCreate20241113 *FlexBackupSnapshotDownloadCreate20241113) DownloadFlexBackupApiRequest
	/*
		DownloadFlexBackup Download One Flex Cluster Snapshot


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadFlexBackupApiParams - Parameters for the request
		@return DownloadFlexBackupApiRequest
	*/
	DownloadFlexBackupWithParams(ctx context.Context, args *DownloadFlexBackupApiParams) DownloadFlexBackupApiRequest

	// Method available only for mocking purposes
	DownloadFlexBackupExecute(r DownloadFlexBackupApiRequest) (*FlexBackupRestoreJob20241113, *http.Response, error)

	/*
		GetFlexBackupSnapshot Return One Snapshot for One Flex Cluster

		Returns one snapshot of one flex cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param name Human-readable label that identifies the flex cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return GetFlexBackupSnapshotApiRequest
	*/
	GetFlexBackupSnapshot(ctx context.Context, groupId string, name string, snapshotId string) GetFlexBackupSnapshotApiRequest
	/*
		GetFlexBackupSnapshot Return One Snapshot for One Flex Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetFlexBackupSnapshotApiParams - Parameters for the request
		@return GetFlexBackupSnapshotApiRequest
	*/
	GetFlexBackupSnapshotWithParams(ctx context.Context, args *GetFlexBackupSnapshotApiParams) GetFlexBackupSnapshotApiRequest

	// Method available only for mocking purposes
	GetFlexBackupSnapshotExecute(r GetFlexBackupSnapshotApiRequest) (*FlexBackupSnapshot20241113, *http.Response, error)

	/*
		ListFlexBackupSnapshots Return All Snapshots for One Flex Cluster

		Returns all snapshots of one flex cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param name Human-readable label that identifies the flex cluster.
		@return ListFlexBackupSnapshotsApiRequest
	*/
	ListFlexBackupSnapshots(ctx context.Context, groupId string, name string) ListFlexBackupSnapshotsApiRequest
	/*
		ListFlexBackupSnapshots Return All Snapshots for One Flex Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListFlexBackupSnapshotsApiParams - Parameters for the request
		@return ListFlexBackupSnapshotsApiRequest
	*/
	ListFlexBackupSnapshotsWithParams(ctx context.Context, args *ListFlexBackupSnapshotsApiParams) ListFlexBackupSnapshotsApiRequest

	// Method available only for mocking purposes
	ListFlexBackupSnapshotsExecute(r ListFlexBackupSnapshotsApiRequest) (*PaginatedApiAtlasFlexBackupSnapshot20241113, *http.Response, error)
}

// FlexSnapshotsApiService FlexSnapshotsApi service
type FlexSnapshotsApiService service

type DownloadFlexBackupApiRequest struct {
	ctx                                      context.Context
	ApiService                               FlexSnapshotsApi
	name                                     string
	groupId                                  string
	flexBackupSnapshotDownloadCreate20241113 *FlexBackupSnapshotDownloadCreate20241113
}

type DownloadFlexBackupApiParams struct {
	Name                                     string
	GroupId                                  string
	FlexBackupSnapshotDownloadCreate20241113 *FlexBackupSnapshotDownloadCreate20241113
}

func (a *FlexSnapshotsApiService) DownloadFlexBackupWithParams(ctx context.Context, args *DownloadFlexBackupApiParams) DownloadFlexBackupApiRequest {
	return DownloadFlexBackupApiRequest{
		ApiService:                               a,
		ctx:                                      ctx,
		name:                                     args.Name,
		groupId:                                  args.GroupId,
		flexBackupSnapshotDownloadCreate20241113: args.FlexBackupSnapshotDownloadCreate20241113,
	}
}

func (r DownloadFlexBackupApiRequest) Execute() (*FlexBackupRestoreJob20241113, *http.Response, error) {
	return r.ApiService.DownloadFlexBackupExecute(r)
}

/*
DownloadFlexBackup Download One Flex Cluster Snapshot

Requests one snapshot for the specified flex cluster. This resource returns a `snapshotURL` that you can use to download the snapshot. This `snapshotURL` remains active for four hours after you make the request.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param name Human-readable label that identifies the flex cluster.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DownloadFlexBackupApiRequest
*/
func (a *FlexSnapshotsApiService) DownloadFlexBackup(ctx context.Context, name string, groupId string, flexBackupSnapshotDownloadCreate20241113 *FlexBackupSnapshotDownloadCreate20241113) DownloadFlexBackupApiRequest {
	return DownloadFlexBackupApiRequest{
		ApiService:                               a,
		ctx:                                      ctx,
		name:                                     name,
		groupId:                                  groupId,
		flexBackupSnapshotDownloadCreate20241113: flexBackupSnapshotDownloadCreate20241113,
	}
}

// DownloadFlexBackupExecute executes the request
//
//	@return FlexBackupRestoreJob20241113
func (a *FlexSnapshotsApiService) DownloadFlexBackupExecute(r DownloadFlexBackupApiRequest) (*FlexBackupRestoreJob20241113, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *FlexBackupRestoreJob20241113
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "FlexSnapshotsApiService.DownloadFlexBackup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/flexClusters/{name}/backup/download"
	if r.name == "" {
		return localVarReturnValue, nil, reportError("name is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", url.PathEscape(r.name), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.flexBackupSnapshotDownloadCreate20241113 == nil {
		return localVarReturnValue, nil, reportError("flexBackupSnapshotDownloadCreate20241113 is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-11-13+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-11-13+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.flexBackupSnapshotDownloadCreate20241113
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

type GetFlexBackupSnapshotApiRequest struct {
	ctx        context.Context
	ApiService FlexSnapshotsApi
	groupId    string
	name       string
	snapshotId string
}

type GetFlexBackupSnapshotApiParams struct {
	GroupId    string
	Name       string
	SnapshotId string
}

func (a *FlexSnapshotsApiService) GetFlexBackupSnapshotWithParams(ctx context.Context, args *GetFlexBackupSnapshotApiParams) GetFlexBackupSnapshotApiRequest {
	return GetFlexBackupSnapshotApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		name:       args.Name,
		snapshotId: args.SnapshotId,
	}
}

func (r GetFlexBackupSnapshotApiRequest) Execute() (*FlexBackupSnapshot20241113, *http.Response, error) {
	return r.ApiService.GetFlexBackupSnapshotExecute(r)
}

/*
GetFlexBackupSnapshot Return One Snapshot for One Flex Cluster

Returns one snapshot of one flex cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param name Human-readable label that identifies the flex cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return GetFlexBackupSnapshotApiRequest
*/
func (a *FlexSnapshotsApiService) GetFlexBackupSnapshot(ctx context.Context, groupId string, name string, snapshotId string) GetFlexBackupSnapshotApiRequest {
	return GetFlexBackupSnapshotApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		name:       name,
		snapshotId: snapshotId,
	}
}

// GetFlexBackupSnapshotExecute executes the request
//
//	@return FlexBackupSnapshot20241113
func (a *FlexSnapshotsApiService) GetFlexBackupSnapshotExecute(r GetFlexBackupSnapshotApiRequest) (*FlexBackupSnapshot20241113, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *FlexBackupSnapshot20241113
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "FlexSnapshotsApiService.GetFlexBackupSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/flexClusters/{name}/backup/snapshots/{snapshotId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.name == "" {
		return localVarReturnValue, nil, reportError("name is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", url.PathEscape(r.name), -1)
	if r.snapshotId == "" {
		return localVarReturnValue, nil, reportError("snapshotId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"snapshotId"+"}", url.PathEscape(r.snapshotId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-11-13+json"}

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

type ListFlexBackupSnapshotsApiRequest struct {
	ctx          context.Context
	ApiService   FlexSnapshotsApi
	groupId      string
	name         string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListFlexBackupSnapshotsApiParams struct {
	GroupId      string
	Name         string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *FlexSnapshotsApiService) ListFlexBackupSnapshotsWithParams(ctx context.Context, args *ListFlexBackupSnapshotsApiParams) ListFlexBackupSnapshotsApiRequest {
	return ListFlexBackupSnapshotsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		name:         args.Name,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListFlexBackupSnapshotsApiRequest) IncludeCount(includeCount bool) ListFlexBackupSnapshotsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListFlexBackupSnapshotsApiRequest) ItemsPerPage(itemsPerPage int) ListFlexBackupSnapshotsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListFlexBackupSnapshotsApiRequest) PageNum(pageNum int) ListFlexBackupSnapshotsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListFlexBackupSnapshotsApiRequest) Execute() (*PaginatedApiAtlasFlexBackupSnapshot20241113, *http.Response, error) {
	return r.ApiService.ListFlexBackupSnapshotsExecute(r)
}

/*
ListFlexBackupSnapshots Return All Snapshots for One Flex Cluster

Returns all snapshots of one flex cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param name Human-readable label that identifies the flex cluster.
	@return ListFlexBackupSnapshotsApiRequest
*/
func (a *FlexSnapshotsApiService) ListFlexBackupSnapshots(ctx context.Context, groupId string, name string) ListFlexBackupSnapshotsApiRequest {
	return ListFlexBackupSnapshotsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		name:       name,
	}
}

// ListFlexBackupSnapshotsExecute executes the request
//
//	@return PaginatedApiAtlasFlexBackupSnapshot20241113
func (a *FlexSnapshotsApiService) ListFlexBackupSnapshotsExecute(r ListFlexBackupSnapshotsApiRequest) (*PaginatedApiAtlasFlexBackupSnapshot20241113, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasFlexBackupSnapshot20241113
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "FlexSnapshotsApiService.ListFlexBackupSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/flexClusters/{name}/backup/snapshots"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.name == "" {
		return localVarReturnValue, nil, reportError("name is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", url.PathEscape(r.name), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-11-13+json"}

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
