// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type RollingIndexApi interface {

	/*
		CreateRollingIndex Create One Rolling Index

		Creates an index on the cluster identified by its name in a rolling manner. Creating the index in this way allows index builds on one replica set member as a standalone at a time, starting with the secondary members. Creating indexes in this way requires at least one replica set election.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster on which MongoDB Cloud creates an index.
		@param databaseRollingIndexRequest Rolling index to create on the specified cluster.
		@return CreateRollingIndexApiRequest
	*/
	CreateRollingIndex(ctx context.Context, groupId string, clusterName string, databaseRollingIndexRequest *DatabaseRollingIndexRequest) CreateRollingIndexApiRequest
	/*
		CreateRollingIndex Create One Rolling Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateRollingIndexApiParams - Parameters for the request
		@return CreateRollingIndexApiRequest
	*/
	CreateRollingIndexWithParams(ctx context.Context, args *CreateRollingIndexApiParams) CreateRollingIndexApiRequest

	// Method available only for mocking purposes
	CreateRollingIndexExecute(r CreateRollingIndexApiRequest) (*http.Response, error)
}

// RollingIndexApiService RollingIndexApi service
type RollingIndexApiService service

type CreateRollingIndexApiRequest struct {
	ctx                         context.Context
	ApiService                  RollingIndexApi
	groupId                     string
	clusterName                 string
	databaseRollingIndexRequest *DatabaseRollingIndexRequest
}

type CreateRollingIndexApiParams struct {
	GroupId                     string
	ClusterName                 string
	DatabaseRollingIndexRequest *DatabaseRollingIndexRequest
}

func (a *RollingIndexApiService) CreateRollingIndexWithParams(ctx context.Context, args *CreateRollingIndexApiParams) CreateRollingIndexApiRequest {
	return CreateRollingIndexApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		groupId:                     args.GroupId,
		clusterName:                 args.ClusterName,
		databaseRollingIndexRequest: args.DatabaseRollingIndexRequest,
	}
}

func (r CreateRollingIndexApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.CreateRollingIndexExecute(r)
}

/*
CreateRollingIndex Create One Rolling Index

Creates an index on the cluster identified by its name in a rolling manner. Creating the index in this way allows index builds on one replica set member as a standalone at a time, starting with the secondary members. Creating indexes in this way requires at least one replica set election.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster on which MongoDB Cloud creates an index.
	@return CreateRollingIndexApiRequest
*/
func (a *RollingIndexApiService) CreateRollingIndex(ctx context.Context, groupId string, clusterName string, databaseRollingIndexRequest *DatabaseRollingIndexRequest) CreateRollingIndexApiRequest {
	return CreateRollingIndexApiRequest{
		ApiService:                  a,
		ctx:                         ctx,
		groupId:                     groupId,
		clusterName:                 clusterName,
		databaseRollingIndexRequest: databaseRollingIndexRequest,
	}
}

// CreateRollingIndexExecute executes the request
func (a *RollingIndexApiService) CreateRollingIndexExecute(r CreateRollingIndexApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RollingIndexApiService.CreateRollingIndex")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/index"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.databaseRollingIndexRequest == nil {
		return nil, reportError("databaseRollingIndexRequest is required and must be specified")
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
	localVarPostBody = r.databaseRollingIndexRequest
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
