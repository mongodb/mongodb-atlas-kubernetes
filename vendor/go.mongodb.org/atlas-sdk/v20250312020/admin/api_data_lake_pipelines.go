// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DataLakePipelinesApi interface {

	/*
		CreatePipeline Create One Data Lake Pipeline

		Creates one Data Lake Pipeline.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param dataLakeIngestionPipeline Creates one Data Lake Pipeline.
		@return CreatePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	CreatePipeline(ctx context.Context, groupId string, dataLakeIngestionPipeline *DataLakeIngestionPipeline) CreatePipelineApiRequest
	/*
		CreatePipeline Create One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreatePipelineApiParams - Parameters for the request
		@return CreatePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	CreatePipelineWithParams(ctx context.Context, args *CreatePipelineApiParams) CreatePipelineApiRequest

	// Method available only for mocking purposes
	CreatePipelineExecute(r CreatePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error)

	/*
		DeletePipeline Remove One Data Lake Pipeline

		Removes one Data Lake Pipeline.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return DeletePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	DeletePipeline(ctx context.Context, groupId string, pipelineName string) DeletePipelineApiRequest
	/*
		DeletePipeline Remove One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeletePipelineApiParams - Parameters for the request
		@return DeletePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	DeletePipelineWithParams(ctx context.Context, args *DeletePipelineApiParams) DeletePipelineApiRequest

	// Method available only for mocking purposes
	DeletePipelineExecute(r DeletePipelineApiRequest) (*http.Response, error)

	/*
		DeletePipelineRun Delete One Pipeline Run Dataset

		Deletes dataset that Atlas generated during the specified pipeline run.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@param pipelineRunId Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
		@return DeletePipelineRunApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	DeletePipelineRun(ctx context.Context, groupId string, pipelineName string, pipelineRunId string) DeletePipelineRunApiRequest
	/*
		DeletePipelineRun Delete One Pipeline Run Dataset


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeletePipelineRunApiParams - Parameters for the request
		@return DeletePipelineRunApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	DeletePipelineRunWithParams(ctx context.Context, args *DeletePipelineRunApiParams) DeletePipelineRunApiRequest

	// Method available only for mocking purposes
	DeletePipelineRunExecute(r DeletePipelineRunApiRequest) (any, *http.Response, error)

	/*
		GetAvailablePipelineSchedules Return All Ingestion Schedules for One Data Lake Pipeline

		Returns a list of backup schedule policy items that you can use as a Data Lake Pipeline source.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return GetAvailablePipelineSchedulesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetAvailablePipelineSchedules(ctx context.Context, groupId string, pipelineName string) GetAvailablePipelineSchedulesApiRequest
	/*
		GetAvailablePipelineSchedules Return All Ingestion Schedules for One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAvailablePipelineSchedulesApiParams - Parameters for the request
		@return GetAvailablePipelineSchedulesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetAvailablePipelineSchedulesWithParams(ctx context.Context, args *GetAvailablePipelineSchedulesApiParams) GetAvailablePipelineSchedulesApiRequest

	// Method available only for mocking purposes
	GetAvailablePipelineSchedulesExecute(r GetAvailablePipelineSchedulesApiRequest) ([]DiskBackupApiPolicyItem, *http.Response, error)

	/*
		GetAvailablePipelineSnapshots Return All Backup Snapshots for One Data Lake Pipeline

		Returns a list of backup snapshots that you can use to trigger an on demand pipeline run.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return GetAvailablePipelineSnapshotsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetAvailablePipelineSnapshots(ctx context.Context, groupId string, pipelineName string) GetAvailablePipelineSnapshotsApiRequest
	/*
		GetAvailablePipelineSnapshots Return All Backup Snapshots for One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAvailablePipelineSnapshotsApiParams - Parameters for the request
		@return GetAvailablePipelineSnapshotsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetAvailablePipelineSnapshotsWithParams(ctx context.Context, args *GetAvailablePipelineSnapshotsApiParams) GetAvailablePipelineSnapshotsApiRequest

	// Method available only for mocking purposes
	GetAvailablePipelineSnapshotsExecute(r GetAvailablePipelineSnapshotsApiRequest) (*PaginatedBackupSnapshot, *http.Response, error)

	/*
		GetPipeline Return One Data Lake Pipeline

		Returns the details of one Data Lake Pipeline within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return GetPipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetPipeline(ctx context.Context, groupId string, pipelineName string) GetPipelineApiRequest
	/*
		GetPipeline Return One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetPipelineApiParams - Parameters for the request
		@return GetPipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetPipelineWithParams(ctx context.Context, args *GetPipelineApiParams) GetPipelineApiRequest

	// Method available only for mocking purposes
	GetPipelineExecute(r GetPipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error)

	/*
		GetPipelineRun Return One Data Lake Pipeline Run

		Returns the details of one Data Lake Pipeline run within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@param pipelineRunId Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
		@return GetPipelineRunApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetPipelineRun(ctx context.Context, groupId string, pipelineName string, pipelineRunId string) GetPipelineRunApiRequest
	/*
		GetPipelineRun Return One Data Lake Pipeline Run


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetPipelineRunApiParams - Parameters for the request
		@return GetPipelineRunApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	GetPipelineRunWithParams(ctx context.Context, args *GetPipelineRunApiParams) GetPipelineRunApiRequest

	// Method available only for mocking purposes
	GetPipelineRunExecute(r GetPipelineRunApiRequest) (*IngestionPipelineRun, *http.Response, error)

	/*
		ListPipelineRuns Return All Data Lake Pipeline Runs in One Project

		Returns a list of past Data Lake Pipeline runs.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return ListPipelineRunsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ListPipelineRuns(ctx context.Context, groupId string, pipelineName string) ListPipelineRunsApiRequest
	/*
		ListPipelineRuns Return All Data Lake Pipeline Runs in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPipelineRunsApiParams - Parameters for the request
		@return ListPipelineRunsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ListPipelineRunsWithParams(ctx context.Context, args *ListPipelineRunsApiParams) ListPipelineRunsApiRequest

	// Method available only for mocking purposes
	ListPipelineRunsExecute(r ListPipelineRunsApiRequest) (*PaginatedPipelineRun, *http.Response, error)

	/*
		ListPipelines Return All Data Lake Pipelines in One Project

		Returns a list of Data Lake Pipelines.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListPipelinesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ListPipelines(ctx context.Context, groupId string) ListPipelinesApiRequest
	/*
		ListPipelines Return All Data Lake Pipelines in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPipelinesApiParams - Parameters for the request
		@return ListPipelinesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ListPipelinesWithParams(ctx context.Context, args *ListPipelinesApiParams) ListPipelinesApiRequest

	// Method available only for mocking purposes
	ListPipelinesExecute(r ListPipelinesApiRequest) ([]DataLakeIngestionPipeline, *http.Response, error)

	/*
		PausePipeline Pause One Data Lake Pipeline

		Pauses ingestion for a Data Lake Pipeline within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return PausePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	PausePipeline(ctx context.Context, groupId string, pipelineName string) PausePipelineApiRequest
	/*
		PausePipeline Pause One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param PausePipelineApiParams - Parameters for the request
		@return PausePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	PausePipelineWithParams(ctx context.Context, args *PausePipelineApiParams) PausePipelineApiRequest

	// Method available only for mocking purposes
	PausePipelineExecute(r PausePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error)

	/*
		ResumePipeline Resume One Data Lake Pipeline

		Resumes ingestion for a Data Lake Pipeline within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@return ResumePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ResumePipeline(ctx context.Context, groupId string, pipelineName string) ResumePipelineApiRequest
	/*
		ResumePipeline Resume One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ResumePipelineApiParams - Parameters for the request
		@return ResumePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	ResumePipelineWithParams(ctx context.Context, args *ResumePipelineApiParams) ResumePipelineApiRequest

	// Method available only for mocking purposes
	ResumePipelineExecute(r ResumePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error)

	/*
		TriggerPipeline Trigger On-Demand Snapshot Ingestion

		Triggers a Data Lake Pipeline ingestion of a specified snapshot.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@param triggerIngestionPipelineRequest Triggers a single ingestion run of a snapshot.
		@return TriggerPipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	TriggerPipeline(ctx context.Context, groupId string, pipelineName string, triggerIngestionPipelineRequest *TriggerIngestionPipelineRequest) TriggerPipelineApiRequest
	/*
		TriggerPipeline Trigger On-Demand Snapshot Ingestion


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param TriggerPipelineApiParams - Parameters for the request
		@return TriggerPipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	TriggerPipelineWithParams(ctx context.Context, args *TriggerPipelineApiParams) TriggerPipelineApiRequest

	// Method available only for mocking purposes
	TriggerPipelineExecute(r TriggerPipelineApiRequest) (*IngestionPipelineRun, *http.Response, error)

	/*
		UpdatePipeline Update One Data Lake Pipeline

		Updates one Data Lake Pipeline.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
		@param dataLakeIngestionPipeline Updates one Data Lake Pipeline.
		@return UpdatePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	UpdatePipeline(ctx context.Context, groupId string, pipelineName string, dataLakeIngestionPipeline *DataLakeIngestionPipeline) UpdatePipelineApiRequest
	/*
		UpdatePipeline Update One Data Lake Pipeline


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdatePipelineApiParams - Parameters for the request
		@return UpdatePipelineApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for DataLakePipelinesApi
	*/
	UpdatePipelineWithParams(ctx context.Context, args *UpdatePipelineApiParams) UpdatePipelineApiRequest

	// Method available only for mocking purposes
	UpdatePipelineExecute(r UpdatePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error)
}

// DataLakePipelinesApiService DataLakePipelinesApi service
type DataLakePipelinesApiService service

type CreatePipelineApiRequest struct {
	ctx                       context.Context
	ApiService                DataLakePipelinesApi
	groupId                   string
	dataLakeIngestionPipeline *DataLakeIngestionPipeline
}

type CreatePipelineApiParams struct {
	GroupId                   string
	DataLakeIngestionPipeline *DataLakeIngestionPipeline
}

func (a *DataLakePipelinesApiService) CreatePipelineWithParams(ctx context.Context, args *CreatePipelineApiParams) CreatePipelineApiRequest {
	return CreatePipelineApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   args.GroupId,
		dataLakeIngestionPipeline: args.DataLakeIngestionPipeline,
	}
}

func (r CreatePipelineApiRequest) Execute() (*DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.CreatePipelineExecute(r)
}

/*
CreatePipeline Create One Data Lake Pipeline

Creates one Data Lake Pipeline.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreatePipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) CreatePipeline(ctx context.Context, groupId string, dataLakeIngestionPipeline *DataLakeIngestionPipeline) CreatePipelineApiRequest {
	return CreatePipelineApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   groupId,
		dataLakeIngestionPipeline: dataLakeIngestionPipeline,
	}
}

// CreatePipelineExecute executes the request
//
//	@return DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) CreatePipelineExecute(r CreatePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.CreatePipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataLakeIngestionPipeline == nil {
		return localVarReturnValue, nil, reportError("dataLakeIngestionPipeline is required and must be specified")
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
	localVarPostBody = r.dataLakeIngestionPipeline
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

type DeletePipelineApiRequest struct {
	ctx          context.Context
	ApiService   DataLakePipelinesApi
	groupId      string
	pipelineName string
}

type DeletePipelineApiParams struct {
	GroupId      string
	PipelineName string
}

func (a *DataLakePipelinesApiService) DeletePipelineWithParams(ctx context.Context, args *DeletePipelineApiParams) DeletePipelineApiRequest {
	return DeletePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		pipelineName: args.PipelineName,
	}
}

func (r DeletePipelineApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeletePipelineExecute(r)
}

/*
DeletePipeline Remove One Data Lake Pipeline

Removes one Data Lake Pipeline.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return DeletePipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) DeletePipeline(ctx context.Context, groupId string, pipelineName string) DeletePipelineApiRequest {
	return DeletePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// DeletePipelineExecute executes the request
// Deprecated
func (a *DataLakePipelinesApiService) DeletePipelineExecute(r DeletePipelineApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.DeletePipeline")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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

type DeletePipelineRunApiRequest struct {
	ctx           context.Context
	ApiService    DataLakePipelinesApi
	groupId       string
	pipelineName  string
	pipelineRunId string
}

type DeletePipelineRunApiParams struct {
	GroupId       string
	PipelineName  string
	PipelineRunId string
}

func (a *DataLakePipelinesApiService) DeletePipelineRunWithParams(ctx context.Context, args *DeletePipelineRunApiParams) DeletePipelineRunApiRequest {
	return DeletePipelineRunApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		pipelineName:  args.PipelineName,
		pipelineRunId: args.PipelineRunId,
	}
}

func (r DeletePipelineRunApiRequest) Execute() (any, *http.Response, error) {
	return r.ApiService.DeletePipelineRunExecute(r)
}

/*
DeletePipelineRun Delete One Pipeline Run Dataset

Deletes dataset that Atlas generated during the specified pipeline run.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@param pipelineRunId Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
	@return DeletePipelineRunApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) DeletePipelineRun(ctx context.Context, groupId string, pipelineName string, pipelineRunId string) DeletePipelineRunApiRequest {
	return DeletePipelineRunApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		pipelineName:  pipelineName,
		pipelineRunId: pipelineRunId,
	}
}

// DeletePipelineRunExecute executes the request
//
//	@return any
//
// Deprecated
func (a *DataLakePipelinesApiService) DeletePipelineRunExecute(r DeletePipelineRunApiRequest) (any, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue any
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.DeletePipelineRun")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/runs/{pipelineRunId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)
	if r.pipelineRunId == "" {
		return localVarReturnValue, nil, reportError("pipelineRunId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineRunId"+"}", url.PathEscape(r.pipelineRunId), -1)

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

type GetAvailablePipelineSchedulesApiRequest struct {
	ctx          context.Context
	ApiService   DataLakePipelinesApi
	groupId      string
	pipelineName string
}

type GetAvailablePipelineSchedulesApiParams struct {
	GroupId      string
	PipelineName string
}

func (a *DataLakePipelinesApiService) GetAvailablePipelineSchedulesWithParams(ctx context.Context, args *GetAvailablePipelineSchedulesApiParams) GetAvailablePipelineSchedulesApiRequest {
	return GetAvailablePipelineSchedulesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		pipelineName: args.PipelineName,
	}
}

func (r GetAvailablePipelineSchedulesApiRequest) Execute() ([]DiskBackupApiPolicyItem, *http.Response, error) {
	return r.ApiService.GetAvailablePipelineSchedulesExecute(r)
}

/*
GetAvailablePipelineSchedules Return All Ingestion Schedules for One Data Lake Pipeline

Returns a list of backup schedule policy items that you can use as a Data Lake Pipeline source.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return GetAvailablePipelineSchedulesApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) GetAvailablePipelineSchedules(ctx context.Context, groupId string, pipelineName string) GetAvailablePipelineSchedulesApiRequest {
	return GetAvailablePipelineSchedulesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// GetAvailablePipelineSchedulesExecute executes the request
//
//	@return []DiskBackupApiPolicyItem
//
// Deprecated
func (a *DataLakePipelinesApiService) GetAvailablePipelineSchedulesExecute(r GetAvailablePipelineSchedulesApiRequest) ([]DiskBackupApiPolicyItem, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []DiskBackupApiPolicyItem
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.GetAvailablePipelineSchedules")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/availableSchedules"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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

type GetAvailablePipelineSnapshotsApiRequest struct {
	ctx            context.Context
	ApiService     DataLakePipelinesApi
	groupId        string
	pipelineName   string
	includeCount   *bool
	itemsPerPage   *int
	pageNum        *int
	completedAfter *time.Time
}

type GetAvailablePipelineSnapshotsApiParams struct {
	GroupId        string
	PipelineName   string
	IncludeCount   *bool
	ItemsPerPage   *int
	PageNum        *int
	CompletedAfter *time.Time
}

func (a *DataLakePipelinesApiService) GetAvailablePipelineSnapshotsWithParams(ctx context.Context, args *GetAvailablePipelineSnapshotsApiParams) GetAvailablePipelineSnapshotsApiRequest {
	return GetAvailablePipelineSnapshotsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		pipelineName:   args.PipelineName,
		includeCount:   args.IncludeCount,
		itemsPerPage:   args.ItemsPerPage,
		pageNum:        args.PageNum,
		completedAfter: args.CompletedAfter,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r GetAvailablePipelineSnapshotsApiRequest) IncludeCount(includeCount bool) GetAvailablePipelineSnapshotsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r GetAvailablePipelineSnapshotsApiRequest) ItemsPerPage(itemsPerPage int) GetAvailablePipelineSnapshotsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetAvailablePipelineSnapshotsApiRequest) PageNum(pageNum int) GetAvailablePipelineSnapshotsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Date and time after which MongoDB Cloud created the snapshot. If specified, MongoDB Cloud returns available backup snapshots created after this time and date only. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
func (r GetAvailablePipelineSnapshotsApiRequest) CompletedAfter(completedAfter time.Time) GetAvailablePipelineSnapshotsApiRequest {
	r.completedAfter = &completedAfter
	return r
}

func (r GetAvailablePipelineSnapshotsApiRequest) Execute() (*PaginatedBackupSnapshot, *http.Response, error) {
	return r.ApiService.GetAvailablePipelineSnapshotsExecute(r)
}

/*
GetAvailablePipelineSnapshots Return All Backup Snapshots for One Data Lake Pipeline

Returns a list of backup snapshots that you can use to trigger an on demand pipeline run.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return GetAvailablePipelineSnapshotsApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) GetAvailablePipelineSnapshots(ctx context.Context, groupId string, pipelineName string) GetAvailablePipelineSnapshotsApiRequest {
	return GetAvailablePipelineSnapshotsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// GetAvailablePipelineSnapshotsExecute executes the request
//
//	@return PaginatedBackupSnapshot
//
// Deprecated
func (a *DataLakePipelinesApiService) GetAvailablePipelineSnapshotsExecute(r GetAvailablePipelineSnapshotsApiRequest) (*PaginatedBackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedBackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.GetAvailablePipelineSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/availableSnapshots"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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
	if r.completedAfter != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "completedAfter", r.completedAfter, "")
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

type GetPipelineApiRequest struct {
	ctx          context.Context
	ApiService   DataLakePipelinesApi
	groupId      string
	pipelineName string
}

type GetPipelineApiParams struct {
	GroupId      string
	PipelineName string
}

func (a *DataLakePipelinesApiService) GetPipelineWithParams(ctx context.Context, args *GetPipelineApiParams) GetPipelineApiRequest {
	return GetPipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		pipelineName: args.PipelineName,
	}
}

func (r GetPipelineApiRequest) Execute() (*DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.GetPipelineExecute(r)
}

/*
GetPipeline Return One Data Lake Pipeline

Returns the details of one Data Lake Pipeline within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return GetPipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) GetPipeline(ctx context.Context, groupId string, pipelineName string) GetPipelineApiRequest {
	return GetPipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// GetPipelineExecute executes the request
//
//	@return DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) GetPipelineExecute(r GetPipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.GetPipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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

type GetPipelineRunApiRequest struct {
	ctx           context.Context
	ApiService    DataLakePipelinesApi
	groupId       string
	pipelineName  string
	pipelineRunId string
}

type GetPipelineRunApiParams struct {
	GroupId       string
	PipelineName  string
	PipelineRunId string
}

func (a *DataLakePipelinesApiService) GetPipelineRunWithParams(ctx context.Context, args *GetPipelineRunApiParams) GetPipelineRunApiRequest {
	return GetPipelineRunApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		pipelineName:  args.PipelineName,
		pipelineRunId: args.PipelineRunId,
	}
}

func (r GetPipelineRunApiRequest) Execute() (*IngestionPipelineRun, *http.Response, error) {
	return r.ApiService.GetPipelineRunExecute(r)
}

/*
GetPipelineRun Return One Data Lake Pipeline Run

Returns the details of one Data Lake Pipeline run within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@param pipelineRunId Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
	@return GetPipelineRunApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) GetPipelineRun(ctx context.Context, groupId string, pipelineName string, pipelineRunId string) GetPipelineRunApiRequest {
	return GetPipelineRunApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		pipelineName:  pipelineName,
		pipelineRunId: pipelineRunId,
	}
}

// GetPipelineRunExecute executes the request
//
//	@return IngestionPipelineRun
//
// Deprecated
func (a *DataLakePipelinesApiService) GetPipelineRunExecute(r GetPipelineRunApiRequest) (*IngestionPipelineRun, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *IngestionPipelineRun
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.GetPipelineRun")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/runs/{pipelineRunId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)
	if r.pipelineRunId == "" {
		return localVarReturnValue, nil, reportError("pipelineRunId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineRunId"+"}", url.PathEscape(r.pipelineRunId), -1)

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

type ListPipelineRunsApiRequest struct {
	ctx           context.Context
	ApiService    DataLakePipelinesApi
	groupId       string
	pipelineName  string
	includeCount  *bool
	itemsPerPage  *int
	pageNum       *int
	createdBefore *time.Time
}

type ListPipelineRunsApiParams struct {
	GroupId       string
	PipelineName  string
	IncludeCount  *bool
	ItemsPerPage  *int
	PageNum       *int
	CreatedBefore *time.Time
}

func (a *DataLakePipelinesApiService) ListPipelineRunsWithParams(ctx context.Context, args *ListPipelineRunsApiParams) ListPipelineRunsApiRequest {
	return ListPipelineRunsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		pipelineName:  args.PipelineName,
		includeCount:  args.IncludeCount,
		itemsPerPage:  args.ItemsPerPage,
		pageNum:       args.PageNum,
		createdBefore: args.CreatedBefore,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListPipelineRunsApiRequest) IncludeCount(includeCount bool) ListPipelineRunsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListPipelineRunsApiRequest) ItemsPerPage(itemsPerPage int) ListPipelineRunsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListPipelineRunsApiRequest) PageNum(pageNum int) ListPipelineRunsApiRequest {
	r.pageNum = &pageNum
	return r
}

// If specified, Atlas returns only Data Lake Pipeline runs initiated before this time and date.
func (r ListPipelineRunsApiRequest) CreatedBefore(createdBefore time.Time) ListPipelineRunsApiRequest {
	r.createdBefore = &createdBefore
	return r
}

func (r ListPipelineRunsApiRequest) Execute() (*PaginatedPipelineRun, *http.Response, error) {
	return r.ApiService.ListPipelineRunsExecute(r)
}

/*
ListPipelineRuns Return All Data Lake Pipeline Runs in One Project

Returns a list of past Data Lake Pipeline runs.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return ListPipelineRunsApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) ListPipelineRuns(ctx context.Context, groupId string, pipelineName string) ListPipelineRunsApiRequest {
	return ListPipelineRunsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// ListPipelineRunsExecute executes the request
//
//	@return PaginatedPipelineRun
//
// Deprecated
func (a *DataLakePipelinesApiService) ListPipelineRunsExecute(r ListPipelineRunsApiRequest) (*PaginatedPipelineRun, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedPipelineRun
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.ListPipelineRuns")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/runs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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
	if r.createdBefore != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "createdBefore", r.createdBefore, "")
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

type ListPipelinesApiRequest struct {
	ctx        context.Context
	ApiService DataLakePipelinesApi
	groupId    string
}

type ListPipelinesApiParams struct {
	GroupId string
}

func (a *DataLakePipelinesApiService) ListPipelinesWithParams(ctx context.Context, args *ListPipelinesApiParams) ListPipelinesApiRequest {
	return ListPipelinesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r ListPipelinesApiRequest) Execute() ([]DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.ListPipelinesExecute(r)
}

/*
ListPipelines Return All Data Lake Pipelines in One Project

Returns a list of Data Lake Pipelines.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListPipelinesApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) ListPipelines(ctx context.Context, groupId string) ListPipelinesApiRequest {
	return ListPipelinesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListPipelinesExecute executes the request
//
//	@return []DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) ListPipelinesExecute(r ListPipelinesApiRequest) ([]DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.ListPipelines")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines"
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

type PausePipelineApiRequest struct {
	ctx          context.Context
	ApiService   DataLakePipelinesApi
	groupId      string
	pipelineName string
}

type PausePipelineApiParams struct {
	GroupId      string
	PipelineName string
}

func (a *DataLakePipelinesApiService) PausePipelineWithParams(ctx context.Context, args *PausePipelineApiParams) PausePipelineApiRequest {
	return PausePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		pipelineName: args.PipelineName,
	}
}

func (r PausePipelineApiRequest) Execute() (*DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.PausePipelineExecute(r)
}

/*
PausePipeline Pause One Data Lake Pipeline

Pauses ingestion for a Data Lake Pipeline within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return PausePipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) PausePipeline(ctx context.Context, groupId string, pipelineName string) PausePipelineApiRequest {
	return PausePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// PausePipelineExecute executes the request
//
//	@return DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) PausePipelineExecute(r PausePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.PausePipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/pause"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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

type ResumePipelineApiRequest struct {
	ctx          context.Context
	ApiService   DataLakePipelinesApi
	groupId      string
	pipelineName string
}

type ResumePipelineApiParams struct {
	GroupId      string
	PipelineName string
}

func (a *DataLakePipelinesApiService) ResumePipelineWithParams(ctx context.Context, args *ResumePipelineApiParams) ResumePipelineApiRequest {
	return ResumePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		pipelineName: args.PipelineName,
	}
}

func (r ResumePipelineApiRequest) Execute() (*DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.ResumePipelineExecute(r)
}

/*
ResumePipeline Resume One Data Lake Pipeline

Resumes ingestion for a Data Lake Pipeline within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return ResumePipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) ResumePipeline(ctx context.Context, groupId string, pipelineName string) ResumePipelineApiRequest {
	return ResumePipelineApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		pipelineName: pipelineName,
	}
}

// ResumePipelineExecute executes the request
//
//	@return DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) ResumePipelineExecute(r ResumePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.ResumePipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/resume"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

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

type TriggerPipelineApiRequest struct {
	ctx                             context.Context
	ApiService                      DataLakePipelinesApi
	groupId                         string
	pipelineName                    string
	triggerIngestionPipelineRequest *TriggerIngestionPipelineRequest
}

type TriggerPipelineApiParams struct {
	GroupId                         string
	PipelineName                    string
	TriggerIngestionPipelineRequest *TriggerIngestionPipelineRequest
}

func (a *DataLakePipelinesApiService) TriggerPipelineWithParams(ctx context.Context, args *TriggerPipelineApiParams) TriggerPipelineApiRequest {
	return TriggerPipelineApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         args.GroupId,
		pipelineName:                    args.PipelineName,
		triggerIngestionPipelineRequest: args.TriggerIngestionPipelineRequest,
	}
}

func (r TriggerPipelineApiRequest) Execute() (*IngestionPipelineRun, *http.Response, error) {
	return r.ApiService.TriggerPipelineExecute(r)
}

/*
TriggerPipeline Trigger On-Demand Snapshot Ingestion

Triggers a Data Lake Pipeline ingestion of a specified snapshot.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return TriggerPipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) TriggerPipeline(ctx context.Context, groupId string, pipelineName string, triggerIngestionPipelineRequest *TriggerIngestionPipelineRequest) TriggerPipelineApiRequest {
	return TriggerPipelineApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         groupId,
		pipelineName:                    pipelineName,
		triggerIngestionPipelineRequest: triggerIngestionPipelineRequest,
	}
}

// TriggerPipelineExecute executes the request
//
//	@return IngestionPipelineRun
//
// Deprecated
func (a *DataLakePipelinesApiService) TriggerPipelineExecute(r TriggerPipelineApiRequest) (*IngestionPipelineRun, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *IngestionPipelineRun
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.TriggerPipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}/trigger"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.triggerIngestionPipelineRequest == nil {
		return localVarReturnValue, nil, reportError("triggerIngestionPipelineRequest is required and must be specified")
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
	localVarPostBody = r.triggerIngestionPipelineRequest
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

type UpdatePipelineApiRequest struct {
	ctx                       context.Context
	ApiService                DataLakePipelinesApi
	groupId                   string
	pipelineName              string
	dataLakeIngestionPipeline *DataLakeIngestionPipeline
}

type UpdatePipelineApiParams struct {
	GroupId                   string
	PipelineName              string
	DataLakeIngestionPipeline *DataLakeIngestionPipeline
}

func (a *DataLakePipelinesApiService) UpdatePipelineWithParams(ctx context.Context, args *UpdatePipelineApiParams) UpdatePipelineApiRequest {
	return UpdatePipelineApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   args.GroupId,
		pipelineName:              args.PipelineName,
		dataLakeIngestionPipeline: args.DataLakeIngestionPipeline,
	}
}

func (r UpdatePipelineApiRequest) Execute() (*DataLakeIngestionPipeline, *http.Response, error) {
	return r.ApiService.UpdatePipelineExecute(r)
}

/*
UpdatePipeline Update One Data Lake Pipeline

Updates one Data Lake Pipeline.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param pipelineName Human-readable label that identifies the Data Lake Pipeline.
	@return UpdatePipelineApiRequest

Deprecated
*/
func (a *DataLakePipelinesApiService) UpdatePipeline(ctx context.Context, groupId string, pipelineName string, dataLakeIngestionPipeline *DataLakeIngestionPipeline) UpdatePipelineApiRequest {
	return UpdatePipelineApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   groupId,
		pipelineName:              pipelineName,
		dataLakeIngestionPipeline: dataLakeIngestionPipeline,
	}
}

// UpdatePipelineExecute executes the request
//
//	@return DataLakeIngestionPipeline
//
// Deprecated
func (a *DataLakePipelinesApiService) UpdatePipelineExecute(r UpdatePipelineApiRequest) (*DataLakeIngestionPipeline, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataLakeIngestionPipeline
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DataLakePipelinesApiService.UpdatePipeline")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/pipelines/{pipelineName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.pipelineName == "" {
		return localVarReturnValue, nil, reportError("pipelineName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"pipelineName"+"}", url.PathEscape(r.pipelineName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataLakeIngestionPipeline == nil {
		return localVarReturnValue, nil, reportError("dataLakeIngestionPipeline is required and must be specified")
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
	localVarPostBody = r.dataLakeIngestionPipeline
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
