// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type LegacyBackupApi interface {

	/*
		CreateClusterRestoreJob Create One Legacy Backup Restore Job

		Restores one legacy backup for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation. This endpoint doesn't support creating checkpoint restore jobs for sharded clusters, or creating restore jobs for queryable backup snapshots. If you create an automated restore job by specifying `delivery.methodName` of `AUTOMATED_RESTORE` in your request body, MongoDB Cloud removes all existing data on the target cluster prior to the restore.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
		@param backupRestoreJob Legacy backup to restore to one cluster in the specified project.
		@return CreateClusterRestoreJobApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	CreateClusterRestoreJob(ctx context.Context, groupId string, clusterName string, backupRestoreJob *BackupRestoreJob) CreateClusterRestoreJobApiRequest
	/*
		CreateClusterRestoreJob Create One Legacy Backup Restore Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateClusterRestoreJobApiParams - Parameters for the request
		@return CreateClusterRestoreJobApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	CreateClusterRestoreJobWithParams(ctx context.Context, args *CreateClusterRestoreJobApiParams) CreateClusterRestoreJobApiRequest

	// Method available only for mocking purposes
	CreateClusterRestoreJobExecute(r CreateClusterRestoreJobApiRequest) (*PaginatedRestoreJob, *http.Response, error)

	/*
		DeleteClusterSnapshot Remove One Legacy Backup Snapshot

		Removes one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return DeleteClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	DeleteClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteClusterSnapshotApiRequest
	/*
		DeleteClusterSnapshot Remove One Legacy Backup Snapshot


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterSnapshotApiParams - Parameters for the request
		@return DeleteClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	DeleteClusterSnapshotWithParams(ctx context.Context, args *DeleteClusterSnapshotApiParams) DeleteClusterSnapshotApiRequest

	// Method available only for mocking purposes
	DeleteClusterSnapshotExecute(r DeleteClusterSnapshotApiRequest) (*http.Response, error)

	/*
		GetClusterBackupCheckpoint Return One Legacy Backup Checkpoint

		Returns one legacy backup checkpoint for one cluster in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param checkpointId Unique 24-hexadecimal digit string that identifies the checkpoint.
		@param clusterName Human-readable label that identifies the cluster that contains the checkpoints that you want to return.
		@return GetClusterBackupCheckpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterBackupCheckpoint(ctx context.Context, groupId string, checkpointId string, clusterName string) GetClusterBackupCheckpointApiRequest
	/*
		GetClusterBackupCheckpoint Return One Legacy Backup Checkpoint


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterBackupCheckpointApiParams - Parameters for the request
		@return GetClusterBackupCheckpointApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterBackupCheckpointWithParams(ctx context.Context, args *GetClusterBackupCheckpointApiParams) GetClusterBackupCheckpointApiRequest

	// Method available only for mocking purposes
	GetClusterBackupCheckpointExecute(r GetClusterBackupCheckpointApiRequest) (*ApiAtlasCheckpoint, *http.Response, error)

	/*
			GetClusterRestoreJob Return One Legacy Backup Restore Job

			Returns one legacy backup restore job for one cluster in the specified project.

		 Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
			@param jobId Unique 24-hexadecimal digit string that identifies the restore job.
			@return GetClusterRestoreJobApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterRestoreJob(ctx context.Context, groupId string, clusterName string, jobId string) GetClusterRestoreJobApiRequest
	/*
		GetClusterRestoreJob Return One Legacy Backup Restore Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterRestoreJobApiParams - Parameters for the request
		@return GetClusterRestoreJobApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterRestoreJobWithParams(ctx context.Context, args *GetClusterRestoreJobApiParams) GetClusterRestoreJobApiRequest

	// Method available only for mocking purposes
	GetClusterRestoreJobExecute(r GetClusterRestoreJobApiRequest) (*BackupRestoreJob, *http.Response, error)

	/*
		GetClusterSnapshot Return One Legacy Backup Snapshot

		Returns one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return GetClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetClusterSnapshotApiRequest
	/*
		GetClusterSnapshot Return One Legacy Backup Snapshot


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterSnapshotApiParams - Parameters for the request
		@return GetClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterSnapshotWithParams(ctx context.Context, args *GetClusterSnapshotApiParams) GetClusterSnapshotApiRequest

	// Method available only for mocking purposes
	GetClusterSnapshotExecute(r GetClusterSnapshotApiRequest) (*BackupSnapshot, *http.Response, error)

	/*
			GetClusterSnapshotSchedule Return One Snapshot Schedule

			Returns the snapshot schedule for one cluster in the specified project.

		 Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
			@return GetClusterSnapshotScheduleApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterSnapshotSchedule(ctx context.Context, groupId string, clusterName string) GetClusterSnapshotScheduleApiRequest
	/*
		GetClusterSnapshotSchedule Return One Snapshot Schedule


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterSnapshotScheduleApiParams - Parameters for the request
		@return GetClusterSnapshotScheduleApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	GetClusterSnapshotScheduleWithParams(ctx context.Context, args *GetClusterSnapshotScheduleApiParams) GetClusterSnapshotScheduleApiRequest

	// Method available only for mocking purposes
	GetClusterSnapshotScheduleExecute(r GetClusterSnapshotScheduleApiRequest) (*ApiAtlasSnapshotSchedule, *http.Response, error)

	/*
		ListClusterBackupCheckpoints Return All Legacy Backup Checkpoints

		Returns all legacy backup checkpoints for one cluster in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster that contains the checkpoints that you want to return.
		@return ListClusterBackupCheckpointsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterBackupCheckpoints(ctx context.Context, groupId string, clusterName string) ListClusterBackupCheckpointsApiRequest
	/*
		ListClusterBackupCheckpoints Return All Legacy Backup Checkpoints


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterBackupCheckpointsApiParams - Parameters for the request
		@return ListClusterBackupCheckpointsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterBackupCheckpointsWithParams(ctx context.Context, args *ListClusterBackupCheckpointsApiParams) ListClusterBackupCheckpointsApiRequest

	// Method available only for mocking purposes
	ListClusterBackupCheckpointsExecute(r ListClusterBackupCheckpointsApiRequest) (*PaginatedApiAtlasCheckpoint, *http.Response, error)

	/*
			ListClusterRestoreJobs Return All Legacy Backup Restore Jobs

			Returns all legacy backup restore jobs for one cluster in the specified project.

		 Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation. If you use the `BATCH-ID` query parameter, you can retrieve all restore jobs in the specified batch. When creating a restore job for a sharded cluster, MongoDB Cloud creates a separate job for each shard, plus another for the config server. Each of those jobs are part of a batch. However, a batch can't include a restore job for a replica set.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
			@return ListClusterRestoreJobsApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterRestoreJobs(ctx context.Context, groupId string, clusterName string) ListClusterRestoreJobsApiRequest
	/*
		ListClusterRestoreJobs Return All Legacy Backup Restore Jobs


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterRestoreJobsApiParams - Parameters for the request
		@return ListClusterRestoreJobsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterRestoreJobsWithParams(ctx context.Context, args *ListClusterRestoreJobsApiParams) ListClusterRestoreJobsApiRequest

	// Method available only for mocking purposes
	ListClusterRestoreJobsExecute(r ListClusterRestoreJobsApiRequest) (*PaginatedRestoreJob, *http.Response, error)

	/*
		ListClusterSnapshots Return All Legacy Backup Snapshots

		Returns all legacy backup snapshots for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListClusterSnapshotsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterSnapshots(ctx context.Context, groupId string, clusterName string) ListClusterSnapshotsApiRequest
	/*
		ListClusterSnapshots Return All Legacy Backup Snapshots


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterSnapshotsApiParams - Parameters for the request
		@return ListClusterSnapshotsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	ListClusterSnapshotsWithParams(ctx context.Context, args *ListClusterSnapshotsApiParams) ListClusterSnapshotsApiRequest

	// Method available only for mocking purposes
	ListClusterSnapshotsExecute(r ListClusterSnapshotsApiRequest) (*PaginatedSnapshot, *http.Response, error)

	/*
		UpdateClusterSnapshot Update Expiration Date for One Legacy Backup Snapshot

		Changes the expiration date for one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@param backupSnapshot Changes One Legacy Backup Snapshot Expiration.
		@return UpdateClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	UpdateClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string, backupSnapshot *BackupSnapshot) UpdateClusterSnapshotApiRequest
	/*
		UpdateClusterSnapshot Update Expiration Date for One Legacy Backup Snapshot


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterSnapshotApiParams - Parameters for the request
		@return UpdateClusterSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	UpdateClusterSnapshotWithParams(ctx context.Context, args *UpdateClusterSnapshotApiParams) UpdateClusterSnapshotApiRequest

	// Method available only for mocking purposes
	UpdateClusterSnapshotExecute(r UpdateClusterSnapshotApiRequest) (*BackupSnapshot, *http.Response, error)

	/*
			UpdateClusterSnapshotSchedule Update Snapshot Schedule for One Cluster

			Updates the snapshot schedule for one cluster in the specified project.

		 Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
			@param apiAtlasSnapshotSchedule Update the snapshot schedule for one cluster in the specified project.
			@return UpdateClusterSnapshotScheduleApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	UpdateClusterSnapshotSchedule(ctx context.Context, groupId string, clusterName string, apiAtlasSnapshotSchedule *ApiAtlasSnapshotSchedule) UpdateClusterSnapshotScheduleApiRequest
	/*
		UpdateClusterSnapshotSchedule Update Snapshot Schedule for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterSnapshotScheduleApiParams - Parameters for the request
		@return UpdateClusterSnapshotScheduleApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for LegacyBackupApi
	*/
	UpdateClusterSnapshotScheduleWithParams(ctx context.Context, args *UpdateClusterSnapshotScheduleApiParams) UpdateClusterSnapshotScheduleApiRequest

	// Method available only for mocking purposes
	UpdateClusterSnapshotScheduleExecute(r UpdateClusterSnapshotScheduleApiRequest) (*ApiAtlasSnapshotSchedule, *http.Response, error)
}

// LegacyBackupApiService LegacyBackupApi service
type LegacyBackupApiService service

type CreateClusterRestoreJobApiRequest struct {
	ctx              context.Context
	ApiService       LegacyBackupApi
	groupId          string
	clusterName      string
	backupRestoreJob *BackupRestoreJob
}

type CreateClusterRestoreJobApiParams struct {
	GroupId          string
	ClusterName      string
	BackupRestoreJob *BackupRestoreJob
}

func (a *LegacyBackupApiService) CreateClusterRestoreJobWithParams(ctx context.Context, args *CreateClusterRestoreJobApiParams) CreateClusterRestoreJobApiRequest {
	return CreateClusterRestoreJobApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		clusterName:      args.ClusterName,
		backupRestoreJob: args.BackupRestoreJob,
	}
}

func (r CreateClusterRestoreJobApiRequest) Execute() (*PaginatedRestoreJob, *http.Response, error) {
	return r.ApiService.CreateClusterRestoreJobExecute(r)
}

/*
CreateClusterRestoreJob Create One Legacy Backup Restore Job

Restores one legacy backup for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation. This endpoint doesn't support creating checkpoint restore jobs for sharded clusters, or creating restore jobs for queryable backup snapshots. If you create an automated restore job by specifying `delivery.methodName` of `AUTOMATED_RESTORE` in your request body, MongoDB Cloud removes all existing data on the target cluster prior to the restore.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
	@return CreateClusterRestoreJobApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) CreateClusterRestoreJob(ctx context.Context, groupId string, clusterName string, backupRestoreJob *BackupRestoreJob) CreateClusterRestoreJobApiRequest {
	return CreateClusterRestoreJobApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		clusterName:      clusterName,
		backupRestoreJob: backupRestoreJob,
	}
}

// CreateClusterRestoreJobExecute executes the request
//
//	@return PaginatedRestoreJob
//
// Deprecated
func (a *LegacyBackupApiService) CreateClusterRestoreJobExecute(r CreateClusterRestoreJobApiRequest) (*PaginatedRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.CreateClusterRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/restoreJobs"
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
	if r.backupRestoreJob == nil {
		return localVarReturnValue, nil, reportError("backupRestoreJob is required and must be specified")
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
	localVarPostBody = r.backupRestoreJob
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

type DeleteClusterSnapshotApiRequest struct {
	ctx         context.Context
	ApiService  LegacyBackupApi
	groupId     string
	clusterName string
	snapshotId  string
}

type DeleteClusterSnapshotApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *LegacyBackupApiService) DeleteClusterSnapshotWithParams(ctx context.Context, args *DeleteClusterSnapshotApiParams) DeleteClusterSnapshotApiRequest {
	return DeleteClusterSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r DeleteClusterSnapshotApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterSnapshotExecute(r)
}

/*
DeleteClusterSnapshot Remove One Legacy Backup Snapshot

Removes one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return DeleteClusterSnapshotApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) DeleteClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteClusterSnapshotApiRequest {
	return DeleteClusterSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// DeleteClusterSnapshotExecute executes the request
// Deprecated
func (a *LegacyBackupApiService) DeleteClusterSnapshotExecute(r DeleteClusterSnapshotApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.DeleteClusterSnapshot")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshots/{snapshotId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.snapshotId == "" {
		return nil, reportError("snapshotId is empty and must be specified")
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

type GetClusterBackupCheckpointApiRequest struct {
	ctx          context.Context
	ApiService   LegacyBackupApi
	groupId      string
	checkpointId string
	clusterName  string
}

type GetClusterBackupCheckpointApiParams struct {
	GroupId      string
	CheckpointId string
	ClusterName  string
}

func (a *LegacyBackupApiService) GetClusterBackupCheckpointWithParams(ctx context.Context, args *GetClusterBackupCheckpointApiParams) GetClusterBackupCheckpointApiRequest {
	return GetClusterBackupCheckpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		checkpointId: args.CheckpointId,
		clusterName:  args.ClusterName,
	}
}

func (r GetClusterBackupCheckpointApiRequest) Execute() (*ApiAtlasCheckpoint, *http.Response, error) {
	return r.ApiService.GetClusterBackupCheckpointExecute(r)
}

/*
GetClusterBackupCheckpoint Return One Legacy Backup Checkpoint

Returns one legacy backup checkpoint for one cluster in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param checkpointId Unique 24-hexadecimal digit string that identifies the checkpoint.
	@param clusterName Human-readable label that identifies the cluster that contains the checkpoints that you want to return.
	@return GetClusterBackupCheckpointApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) GetClusterBackupCheckpoint(ctx context.Context, groupId string, checkpointId string, clusterName string) GetClusterBackupCheckpointApiRequest {
	return GetClusterBackupCheckpointApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		checkpointId: checkpointId,
		clusterName:  clusterName,
	}
}

// GetClusterBackupCheckpointExecute executes the request
//
//	@return ApiAtlasCheckpoint
//
// Deprecated
func (a *LegacyBackupApiService) GetClusterBackupCheckpointExecute(r GetClusterBackupCheckpointApiRequest) (*ApiAtlasCheckpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasCheckpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.GetClusterBackupCheckpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backupCheckpoints/{checkpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.checkpointId == "" {
		return localVarReturnValue, nil, reportError("checkpointId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"checkpointId"+"}", url.PathEscape(r.checkpointId), -1)
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

type GetClusterRestoreJobApiRequest struct {
	ctx         context.Context
	ApiService  LegacyBackupApi
	groupId     string
	clusterName string
	jobId       string
}

type GetClusterRestoreJobApiParams struct {
	GroupId     string
	ClusterName string
	JobId       string
}

func (a *LegacyBackupApiService) GetClusterRestoreJobWithParams(ctx context.Context, args *GetClusterRestoreJobApiParams) GetClusterRestoreJobApiRequest {
	return GetClusterRestoreJobApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		jobId:       args.JobId,
	}
}

func (r GetClusterRestoreJobApiRequest) Execute() (*BackupRestoreJob, *http.Response, error) {
	return r.ApiService.GetClusterRestoreJobExecute(r)
}

/*
GetClusterRestoreJob Return One Legacy Backup Restore Job

Returns one legacy backup restore job for one cluster in the specified project.

	Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
	@param jobId Unique 24-hexadecimal digit string that identifies the restore job.
	@return GetClusterRestoreJobApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) GetClusterRestoreJob(ctx context.Context, groupId string, clusterName string, jobId string) GetClusterRestoreJobApiRequest {
	return GetClusterRestoreJobApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		jobId:       jobId,
	}
}

// GetClusterRestoreJobExecute executes the request
//
//	@return BackupRestoreJob
//
// Deprecated
func (a *LegacyBackupApiService) GetClusterRestoreJobExecute(r GetClusterRestoreJobApiRequest) (*BackupRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.GetClusterRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/restoreJobs/{jobId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.jobId == "" {
		return localVarReturnValue, nil, reportError("jobId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"jobId"+"}", url.PathEscape(r.jobId), -1)

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

type GetClusterSnapshotApiRequest struct {
	ctx         context.Context
	ApiService  LegacyBackupApi
	groupId     string
	clusterName string
	snapshotId  string
}

type GetClusterSnapshotApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *LegacyBackupApiService) GetClusterSnapshotWithParams(ctx context.Context, args *GetClusterSnapshotApiParams) GetClusterSnapshotApiRequest {
	return GetClusterSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r GetClusterSnapshotApiRequest) Execute() (*BackupSnapshot, *http.Response, error) {
	return r.ApiService.GetClusterSnapshotExecute(r)
}

/*
GetClusterSnapshot Return One Legacy Backup Snapshot

Returns one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return GetClusterSnapshotApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) GetClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetClusterSnapshotApiRequest {
	return GetClusterSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// GetClusterSnapshotExecute executes the request
//
//	@return BackupSnapshot
//
// Deprecated
func (a *LegacyBackupApiService) GetClusterSnapshotExecute(r GetClusterSnapshotApiRequest) (*BackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.GetClusterSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshots/{snapshotId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
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

type GetClusterSnapshotScheduleApiRequest struct {
	ctx         context.Context
	ApiService  LegacyBackupApi
	groupId     string
	clusterName string
}

type GetClusterSnapshotScheduleApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *LegacyBackupApiService) GetClusterSnapshotScheduleWithParams(ctx context.Context, args *GetClusterSnapshotScheduleApiParams) GetClusterSnapshotScheduleApiRequest {
	return GetClusterSnapshotScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetClusterSnapshotScheduleApiRequest) Execute() (*ApiAtlasSnapshotSchedule, *http.Response, error) {
	return r.ApiService.GetClusterSnapshotScheduleExecute(r)
}

/*
GetClusterSnapshotSchedule Return One Snapshot Schedule

Returns the snapshot schedule for one cluster in the specified project.

	Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
	@return GetClusterSnapshotScheduleApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) GetClusterSnapshotSchedule(ctx context.Context, groupId string, clusterName string) GetClusterSnapshotScheduleApiRequest {
	return GetClusterSnapshotScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetClusterSnapshotScheduleExecute executes the request
//
//	@return ApiAtlasSnapshotSchedule
//
// Deprecated
func (a *LegacyBackupApiService) GetClusterSnapshotScheduleExecute(r GetClusterSnapshotScheduleApiRequest) (*ApiAtlasSnapshotSchedule, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasSnapshotSchedule
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.GetClusterSnapshotSchedule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshotSchedule"
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

type ListClusterBackupCheckpointsApiRequest struct {
	ctx          context.Context
	ApiService   LegacyBackupApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListClusterBackupCheckpointsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *LegacyBackupApiService) ListClusterBackupCheckpointsWithParams(ctx context.Context, args *ListClusterBackupCheckpointsApiParams) ListClusterBackupCheckpointsApiRequest {
	return ListClusterBackupCheckpointsApiRequest{
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
func (r ListClusterBackupCheckpointsApiRequest) IncludeCount(includeCount bool) ListClusterBackupCheckpointsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterBackupCheckpointsApiRequest) ItemsPerPage(itemsPerPage int) ListClusterBackupCheckpointsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterBackupCheckpointsApiRequest) PageNum(pageNum int) ListClusterBackupCheckpointsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListClusterBackupCheckpointsApiRequest) Execute() (*PaginatedApiAtlasCheckpoint, *http.Response, error) {
	return r.ApiService.ListClusterBackupCheckpointsExecute(r)
}

/*
ListClusterBackupCheckpoints Return All Legacy Backup Checkpoints

Returns all legacy backup checkpoints for one cluster in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster that contains the checkpoints that you want to return.
	@return ListClusterBackupCheckpointsApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) ListClusterBackupCheckpoints(ctx context.Context, groupId string, clusterName string) ListClusterBackupCheckpointsApiRequest {
	return ListClusterBackupCheckpointsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterBackupCheckpointsExecute executes the request
//
//	@return PaginatedApiAtlasCheckpoint
//
// Deprecated
func (a *LegacyBackupApiService) ListClusterBackupCheckpointsExecute(r ListClusterBackupCheckpointsApiRequest) (*PaginatedApiAtlasCheckpoint, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasCheckpoint
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.ListClusterBackupCheckpoints")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backupCheckpoints"
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

type ListClusterRestoreJobsApiRequest struct {
	ctx          context.Context
	ApiService   LegacyBackupApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	batchId      *string
}

type ListClusterRestoreJobsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	BatchId      *string
}

func (a *LegacyBackupApiService) ListClusterRestoreJobsWithParams(ctx context.Context, args *ListClusterRestoreJobsApiParams) ListClusterRestoreJobsApiRequest {
	return ListClusterRestoreJobsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		batchId:      args.BatchId,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClusterRestoreJobsApiRequest) IncludeCount(includeCount bool) ListClusterRestoreJobsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterRestoreJobsApiRequest) ItemsPerPage(itemsPerPage int) ListClusterRestoreJobsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterRestoreJobsApiRequest) PageNum(pageNum int) ListClusterRestoreJobsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Unique 24-hexadecimal digit string that identifies the batch of restore jobs to return. Timestamp in ISO 8601 date and time format in UTC when creating a restore job for a sharded cluster, Application creates a separate job for each shard, plus another for the config host. Each of these jobs comprise one batch. A restore job for a replica set can&#39;t be part of a batch.
func (r ListClusterRestoreJobsApiRequest) BatchId(batchId string) ListClusterRestoreJobsApiRequest {
	r.batchId = &batchId
	return r
}

func (r ListClusterRestoreJobsApiRequest) Execute() (*PaginatedRestoreJob, *http.Response, error) {
	return r.ApiService.ListClusterRestoreJobsExecute(r)
}

/*
ListClusterRestoreJobs Return All Legacy Backup Restore Jobs

Returns all legacy backup restore jobs for one cluster in the specified project.

	Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation. If you use the `BATCH-ID` query parameter, you can retrieve all restore jobs in the specified batch. When creating a restore job for a sharded cluster, MongoDB Cloud creates a separate job for each shard, plus another for the config server. Each of those jobs are part of a batch. However, a batch can't include a restore job for a replica set.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
	@return ListClusterRestoreJobsApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) ListClusterRestoreJobs(ctx context.Context, groupId string, clusterName string) ListClusterRestoreJobsApiRequest {
	return ListClusterRestoreJobsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterRestoreJobsExecute executes the request
//
//	@return PaginatedRestoreJob
//
// Deprecated
func (a *LegacyBackupApiService) ListClusterRestoreJobsExecute(r ListClusterRestoreJobsApiRequest) (*PaginatedRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.ListClusterRestoreJobs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/restoreJobs"
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
	if r.batchId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "batchId", r.batchId, "")
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

type ListClusterSnapshotsApiRequest struct {
	ctx          context.Context
	ApiService   LegacyBackupApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	completed    *string
}

type ListClusterSnapshotsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Completed    *string
}

func (a *LegacyBackupApiService) ListClusterSnapshotsWithParams(ctx context.Context, args *ListClusterSnapshotsApiParams) ListClusterSnapshotsApiRequest {
	return ListClusterSnapshotsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		completed:    args.Completed,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClusterSnapshotsApiRequest) IncludeCount(includeCount bool) ListClusterSnapshotsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterSnapshotsApiRequest) ItemsPerPage(itemsPerPage int) ListClusterSnapshotsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterSnapshotsApiRequest) PageNum(pageNum int) ListClusterSnapshotsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Human-readable label that specifies whether to return only completed, incomplete, or all snapshots. By default, MongoDB Cloud only returns completed snapshots.
func (r ListClusterSnapshotsApiRequest) Completed(completed string) ListClusterSnapshotsApiRequest {
	r.completed = &completed
	return r
}

func (r ListClusterSnapshotsApiRequest) Execute() (*PaginatedSnapshot, *http.Response, error) {
	return r.ApiService.ListClusterSnapshotsExecute(r)
}

/*
ListClusterSnapshots Return All Legacy Backup Snapshots

Returns all legacy backup snapshots for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListClusterSnapshotsApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) ListClusterSnapshots(ctx context.Context, groupId string, clusterName string) ListClusterSnapshotsApiRequest {
	return ListClusterSnapshotsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterSnapshotsExecute executes the request
//
//	@return PaginatedSnapshot
//
// Deprecated
func (a *LegacyBackupApiService) ListClusterSnapshotsExecute(r ListClusterSnapshotsApiRequest) (*PaginatedSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.ListClusterSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshots"
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
	if r.completed != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "completed", r.completed, "")
	} else {
		var defaultValue string = "true"
		r.completed = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "completed", r.completed, "")
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

type UpdateClusterSnapshotApiRequest struct {
	ctx            context.Context
	ApiService     LegacyBackupApi
	groupId        string
	clusterName    string
	snapshotId     string
	backupSnapshot *BackupSnapshot
}

type UpdateClusterSnapshotApiParams struct {
	GroupId        string
	ClusterName    string
	SnapshotId     string
	BackupSnapshot *BackupSnapshot
}

func (a *LegacyBackupApiService) UpdateClusterSnapshotWithParams(ctx context.Context, args *UpdateClusterSnapshotApiParams) UpdateClusterSnapshotApiRequest {
	return UpdateClusterSnapshotApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		snapshotId:     args.SnapshotId,
		backupSnapshot: args.BackupSnapshot,
	}
}

func (r UpdateClusterSnapshotApiRequest) Execute() (*BackupSnapshot, *http.Response, error) {
	return r.ApiService.UpdateClusterSnapshotExecute(r)
}

/*
UpdateClusterSnapshot Update Expiration Date for One Legacy Backup Snapshot

Changes the expiration date for one legacy backup snapshot for one cluster in the specified project. Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return UpdateClusterSnapshotApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) UpdateClusterSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string, backupSnapshot *BackupSnapshot) UpdateClusterSnapshotApiRequest {
	return UpdateClusterSnapshotApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		snapshotId:     snapshotId,
		backupSnapshot: backupSnapshot,
	}
}

// UpdateClusterSnapshotExecute executes the request
//
//	@return BackupSnapshot
//
// Deprecated
func (a *LegacyBackupApiService) UpdateClusterSnapshotExecute(r UpdateClusterSnapshotApiRequest) (*BackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.UpdateClusterSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshots/{snapshotId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.snapshotId == "" {
		return localVarReturnValue, nil, reportError("snapshotId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"snapshotId"+"}", url.PathEscape(r.snapshotId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.backupSnapshot == nil {
		return localVarReturnValue, nil, reportError("backupSnapshot is required and must be specified")
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
	localVarPostBody = r.backupSnapshot
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

type UpdateClusterSnapshotScheduleApiRequest struct {
	ctx                      context.Context
	ApiService               LegacyBackupApi
	groupId                  string
	clusterName              string
	apiAtlasSnapshotSchedule *ApiAtlasSnapshotSchedule
}

type UpdateClusterSnapshotScheduleApiParams struct {
	GroupId                  string
	ClusterName              string
	ApiAtlasSnapshotSchedule *ApiAtlasSnapshotSchedule
}

func (a *LegacyBackupApiService) UpdateClusterSnapshotScheduleWithParams(ctx context.Context, args *UpdateClusterSnapshotScheduleApiParams) UpdateClusterSnapshotScheduleApiRequest {
	return UpdateClusterSnapshotScheduleApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		clusterName:              args.ClusterName,
		apiAtlasSnapshotSchedule: args.ApiAtlasSnapshotSchedule,
	}
}

func (r UpdateClusterSnapshotScheduleApiRequest) Execute() (*ApiAtlasSnapshotSchedule, *http.Response, error) {
	return r.ApiService.UpdateClusterSnapshotScheduleExecute(r)
}

/*
UpdateClusterSnapshotSchedule Update Snapshot Schedule for One Cluster

Updates the snapshot schedule for one cluster in the specified project.

	Effective 23 March 2020, all new clusters can use only Cloud Backups. When you upgrade to 4.2, your backup system upgrades to cloud backup if it is currently set to legacy backup. After this upgrade, all your existing legacy backup snapshots remain available. They expire over time in accordance with your retention policy. Your backup policy resets to the default schedule. If you had a custom backup policy in place with legacy backups, you must re-create it with the procedure outlined in the Cloud Backup documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the snapshot you want to return.
	@return UpdateClusterSnapshotScheduleApiRequest

Deprecated
*/
func (a *LegacyBackupApiService) UpdateClusterSnapshotSchedule(ctx context.Context, groupId string, clusterName string, apiAtlasSnapshotSchedule *ApiAtlasSnapshotSchedule) UpdateClusterSnapshotScheduleApiRequest {
	return UpdateClusterSnapshotScheduleApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		clusterName:              clusterName,
		apiAtlasSnapshotSchedule: apiAtlasSnapshotSchedule,
	}
}

// UpdateClusterSnapshotScheduleExecute executes the request
//
//	@return ApiAtlasSnapshotSchedule
//
// Deprecated
func (a *LegacyBackupApiService) UpdateClusterSnapshotScheduleExecute(r UpdateClusterSnapshotScheduleApiRequest) (*ApiAtlasSnapshotSchedule, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiAtlasSnapshotSchedule
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LegacyBackupApiService.UpdateClusterSnapshotSchedule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/snapshotSchedule"
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
	if r.apiAtlasSnapshotSchedule == nil {
		return localVarReturnValue, nil, reportError("apiAtlasSnapshotSchedule is required and must be specified")
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
	localVarPostBody = r.apiAtlasSnapshotSchedule
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
