// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CloudBackupsApi interface {

	/*
		CancelBackupRestoreJob Cancel One Restore Job for One Cluster

		Cancels one cloud backup restore job of one cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to remove.
		@return CancelBackupRestoreJobApiRequest
	*/
	CancelBackupRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) CancelBackupRestoreJobApiRequest
	/*
		CancelBackupRestoreJob Cancel One Restore Job for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CancelBackupRestoreJobApiParams - Parameters for the request
		@return CancelBackupRestoreJobApiRequest
	*/
	CancelBackupRestoreJobWithParams(ctx context.Context, args *CancelBackupRestoreJobApiParams) CancelBackupRestoreJobApiRequest

	// Method available only for mocking purposes
	CancelBackupRestoreJobExecute(r CancelBackupRestoreJobApiRequest) (*http.Response, error)

	/*
		CreateBackupExport Create One Snapshot Export Job

		Exports one backup Snapshot for dedicated Atlas cluster using Cloud Backups to an Export Bucket.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param diskBackupExportJobRequest Information about the Cloud Backup Snapshot Export Job to create.
		@return CreateBackupExportApiRequest
	*/
	CreateBackupExport(ctx context.Context, groupId string, clusterName string, diskBackupExportJobRequest *DiskBackupExportJobRequest) CreateBackupExportApiRequest
	/*
		CreateBackupExport Create One Snapshot Export Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateBackupExportApiParams - Parameters for the request
		@return CreateBackupExportApiRequest
	*/
	CreateBackupExportWithParams(ctx context.Context, args *CreateBackupExportApiParams) CreateBackupExportApiRequest

	// Method available only for mocking purposes
	CreateBackupExportExecute(r CreateBackupExportApiRequest) (*DiskBackupExportJob, *http.Response, error)

	/*
		CreateBackupPrivateEndpoint Create One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project

		Creates a private endpoint in the specified region for secure, private connectivity between Atlas and cloud provider object storage services for backup operations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoint to create.
		@param objectStoragePrivateEndpointRequest Creates a private endpoint in the specified region for object storage backup operations.
		@return CreateBackupPrivateEndpointApiRequest
	*/
	CreateBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, objectStoragePrivateEndpointRequest *ObjectStoragePrivateEndpointRequest) CreateBackupPrivateEndpointApiRequest
	/*
		CreateBackupPrivateEndpoint Create One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateBackupPrivateEndpointApiParams - Parameters for the request
		@return CreateBackupPrivateEndpointApiRequest
	*/
	CreateBackupPrivateEndpointWithParams(ctx context.Context, args *CreateBackupPrivateEndpointApiParams) CreateBackupPrivateEndpointApiRequest

	// Method available only for mocking purposes
	CreateBackupPrivateEndpointExecute(r CreateBackupPrivateEndpointApiRequest) (*ObjectStoragePrivateEndpointResponse, *http.Response, error)

	/*
		CreateBackupRestoreJob Create One Restore Job of One Cluster

		Restores one snapshot of one cluster from the specified project. Atlas takes on-demand snapshots immediately and scheduled snapshots at regular intervals. If an on-demand snapshot with a status of `queued` or `inProgress` exists, before taking another snapshot, wait until Atlas completes processing the previously taken on-demand snapshot.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param diskBackupSnapshotRestoreJob Restores one snapshot of one cluster from the specified project.
		@return CreateBackupRestoreJobApiRequest
	*/
	CreateBackupRestoreJob(ctx context.Context, groupId string, clusterName string, diskBackupSnapshotRestoreJob *DiskBackupSnapshotRestoreJob) CreateBackupRestoreJobApiRequest
	/*
		CreateBackupRestoreJob Create One Restore Job of One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateBackupRestoreJobApiParams - Parameters for the request
		@return CreateBackupRestoreJobApiRequest
	*/
	CreateBackupRestoreJobWithParams(ctx context.Context, args *CreateBackupRestoreJobApiParams) CreateBackupRestoreJobApiRequest

	// Method available only for mocking purposes
	CreateBackupRestoreJobExecute(r CreateBackupRestoreJobApiRequest) (*DiskBackupSnapshotRestoreJob, *http.Response, error)

	/*
		CreateExportBucket Create One Snapshot Export Bucket

		Creates a Snapshot Export Bucket for an AWS S3 Bucket, Azure Blob Storage Container, or Google Cloud Storage Bucket. Once created, an snapshots can be exported to the Export Bucket and its referenced AWS S3 Bucket, Azure Blob Storage Container, or Google Cloud Storage Bucket. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param diskBackupSnapshotExportBucketRequest Specifies the role and AWS S3 Bucket, Azure Blob Storage Container, or Google Cloud Storage Bucket that the Export Bucket should reference.
		@return CreateExportBucketApiRequest
	*/
	CreateExportBucket(ctx context.Context, groupId string, diskBackupSnapshotExportBucketRequest *DiskBackupSnapshotExportBucketRequest) CreateExportBucketApiRequest
	/*
		CreateExportBucket Create One Snapshot Export Bucket


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateExportBucketApiParams - Parameters for the request
		@return CreateExportBucketApiRequest
	*/
	CreateExportBucketWithParams(ctx context.Context, args *CreateExportBucketApiParams) CreateExportBucketApiRequest

	// Method available only for mocking purposes
	CreateExportBucketExecute(r CreateExportBucketApiRequest) (*DiskBackupSnapshotExportBucketResponse, *http.Response, error)

	/*
			CreateServerlessRestoreJob Create One Restore Job for One Serverless Instance

			Restores one snapshot of one serverless instance from the specified project.

		This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Create Flex Backup Restore Job endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the serverless instance whose snapshot you want to restore.
			@param serverlessBackupRestoreJob Restores one snapshot of one serverless instance from the specified project.
			@return CreateServerlessRestoreJobApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	CreateServerlessRestoreJob(ctx context.Context, groupId string, clusterName string, serverlessBackupRestoreJob *ServerlessBackupRestoreJob) CreateServerlessRestoreJobApiRequest
	/*
		CreateServerlessRestoreJob Create One Restore Job for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateServerlessRestoreJobApiParams - Parameters for the request
		@return CreateServerlessRestoreJobApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	CreateServerlessRestoreJobWithParams(ctx context.Context, args *CreateServerlessRestoreJobApiParams) CreateServerlessRestoreJobApiRequest

	// Method available only for mocking purposes
	CreateServerlessRestoreJobExecute(r CreateServerlessRestoreJobApiRequest) (*ServerlessBackupRestoreJob, *http.Response, error)

	/*
		DeleteBackupPrivateEndpoint Delete One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider from One Project

		Deletes one private endpoint, identified by its ID, for object storage backup operations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint to delete.
		@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint to delete.
		@return DeleteBackupPrivateEndpointApiRequest
	*/
	DeleteBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) DeleteBackupPrivateEndpointApiRequest
	/*
		DeleteBackupPrivateEndpoint Delete One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteBackupPrivateEndpointApiParams - Parameters for the request
		@return DeleteBackupPrivateEndpointApiRequest
	*/
	DeleteBackupPrivateEndpointWithParams(ctx context.Context, args *DeleteBackupPrivateEndpointApiParams) DeleteBackupPrivateEndpointApiRequest

	// Method available only for mocking purposes
	DeleteBackupPrivateEndpointExecute(r DeleteBackupPrivateEndpointApiRequest) (*http.Response, error)

	/*
		DeleteBackupShardedCluster Remove One Sharded Cluster Cloud Backup

		Removes one snapshot of one sharded cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return DeleteBackupShardedClusterApiRequest
	*/
	DeleteBackupShardedCluster(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteBackupShardedClusterApiRequest
	/*
		DeleteBackupShardedCluster Remove One Sharded Cluster Cloud Backup


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteBackupShardedClusterApiParams - Parameters for the request
		@return DeleteBackupShardedClusterApiRequest
	*/
	DeleteBackupShardedClusterWithParams(ctx context.Context, args *DeleteBackupShardedClusterApiParams) DeleteBackupShardedClusterApiRequest

	// Method available only for mocking purposes
	DeleteBackupShardedClusterExecute(r DeleteBackupShardedClusterApiRequest) (*http.Response, error)

	/*
		DeleteClusterBackupSchedule Remove All Cloud Backup Schedules

		Removes all cloud backup schedules for the specified cluster. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return DeleteClusterBackupScheduleApiRequest
	*/
	DeleteClusterBackupSchedule(ctx context.Context, groupId string, clusterName string) DeleteClusterBackupScheduleApiRequest
	/*
		DeleteClusterBackupSchedule Remove All Cloud Backup Schedules


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterBackupScheduleApiParams - Parameters for the request
		@return DeleteClusterBackupScheduleApiRequest
	*/
	DeleteClusterBackupScheduleWithParams(ctx context.Context, args *DeleteClusterBackupScheduleApiParams) DeleteClusterBackupScheduleApiRequest

	// Method available only for mocking purposes
	DeleteClusterBackupScheduleExecute(r DeleteClusterBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error)

	/*
		DeleteClusterBackupSnapshot Remove One Replica Set Cloud Backup

		Removes the specified snapshot.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return DeleteClusterBackupSnapshotApiRequest
	*/
	DeleteClusterBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteClusterBackupSnapshotApiRequest
	/*
		DeleteClusterBackupSnapshot Remove One Replica Set Cloud Backup


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterBackupSnapshotApiParams - Parameters for the request
		@return DeleteClusterBackupSnapshotApiRequest
	*/
	DeleteClusterBackupSnapshotWithParams(ctx context.Context, args *DeleteClusterBackupSnapshotApiParams) DeleteClusterBackupSnapshotApiRequest

	// Method available only for mocking purposes
	DeleteClusterBackupSnapshotExecute(r DeleteClusterBackupSnapshotApiRequest) (*http.Response, error)

	/*
		DeleteExportBucket Delete One Snapshot Export Bucket

		Deletes an Export Bucket. Auto export must be disabled on all clusters in this Project exporting to this Export Bucket before revoking access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param exportBucketId Unique 24-hexadecimal character string that identifies the Export Bucket.
		@return DeleteExportBucketApiRequest
	*/
	DeleteExportBucket(ctx context.Context, groupId string, exportBucketId string) DeleteExportBucketApiRequest
	/*
		DeleteExportBucket Delete One Snapshot Export Bucket


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteExportBucketApiParams - Parameters for the request
		@return DeleteExportBucketApiRequest
	*/
	DeleteExportBucketWithParams(ctx context.Context, args *DeleteExportBucketApiParams) DeleteExportBucketApiRequest

	// Method available only for mocking purposes
	DeleteExportBucketExecute(r DeleteExportBucketApiRequest) (*http.Response, error)

	/*
		DisableCompliancePolicy Disable Backup Compliance Policy Settings

		Disables the Backup Compliance Policy settings with the specified project. As a prerequisite, a support ticket needs to be file first, instructions in https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return DisableCompliancePolicyApiRequest
	*/
	DisableCompliancePolicy(ctx context.Context, groupId string) DisableCompliancePolicyApiRequest
	/*
		DisableCompliancePolicy Disable Backup Compliance Policy Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DisableCompliancePolicyApiParams - Parameters for the request
		@return DisableCompliancePolicyApiRequest
	*/
	DisableCompliancePolicyWithParams(ctx context.Context, args *DisableCompliancePolicyApiParams) DisableCompliancePolicyApiRequest

	// Method available only for mocking purposes
	DisableCompliancePolicyExecute(r DisableCompliancePolicyApiRequest) (*http.Response, error)

	/*
		GetBackupExport Return One Snapshot Export Job

		Returns one Cloud Backup Snapshot Export Job associated with the specified Atlas cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param exportId Unique 24-hexadecimal character string that identifies the Export Job.
		@return GetBackupExportApiRequest
	*/
	GetBackupExport(ctx context.Context, groupId string, clusterName string, exportId string) GetBackupExportApiRequest
	/*
		GetBackupExport Return One Snapshot Export Job


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetBackupExportApiParams - Parameters for the request
		@return GetBackupExportApiRequest
	*/
	GetBackupExportWithParams(ctx context.Context, args *GetBackupExportApiParams) GetBackupExportApiRequest

	// Method available only for mocking purposes
	GetBackupExportExecute(r GetBackupExportApiRequest) (*DiskBackupExportJob, *http.Response, error)

	/*
		GetBackupPrivateEndpoint Return One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project

		Returns one private endpoint, identified by its ID, for object storage backup operations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint.
		@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint.
		@return GetBackupPrivateEndpointApiRequest
	*/
	GetBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) GetBackupPrivateEndpointApiRequest
	/*
		GetBackupPrivateEndpoint Return One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetBackupPrivateEndpointApiParams - Parameters for the request
		@return GetBackupPrivateEndpointApiRequest
	*/
	GetBackupPrivateEndpointWithParams(ctx context.Context, args *GetBackupPrivateEndpointApiParams) GetBackupPrivateEndpointApiRequest

	// Method available only for mocking purposes
	GetBackupPrivateEndpointExecute(r GetBackupPrivateEndpointApiRequest) (*ObjectStoragePrivateEndpointResponse, *http.Response, error)

	/*
		GetBackupRestoreJob Return One Restore Job for One Cluster

		Returns one cloud backup restore job for one cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster with the restore jobs you want to return.
		@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to return.
		@return GetBackupRestoreJobApiRequest
	*/
	GetBackupRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) GetBackupRestoreJobApiRequest
	/*
		GetBackupRestoreJob Return One Restore Job for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetBackupRestoreJobApiParams - Parameters for the request
		@return GetBackupRestoreJobApiRequest
	*/
	GetBackupRestoreJobWithParams(ctx context.Context, args *GetBackupRestoreJobApiParams) GetBackupRestoreJobApiRequest

	// Method available only for mocking purposes
	GetBackupRestoreJobExecute(r GetBackupRestoreJobApiRequest) (*DiskBackupSnapshotRestoreJob, *http.Response, error)

	/*
		GetBackupSchedule Return One Cloud Backup Schedule

		Returns the cloud backup schedule for the specified cluster within the specified project. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return GetBackupScheduleApiRequest
	*/
	GetBackupSchedule(ctx context.Context, groupId string, clusterName string) GetBackupScheduleApiRequest
	/*
		GetBackupSchedule Return One Cloud Backup Schedule


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetBackupScheduleApiParams - Parameters for the request
		@return GetBackupScheduleApiRequest
	*/
	GetBackupScheduleWithParams(ctx context.Context, args *GetBackupScheduleApiParams) GetBackupScheduleApiRequest

	// Method available only for mocking purposes
	GetBackupScheduleExecute(r GetBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error)

	/*
		GetBackupShardedCluster Return One Sharded Cluster Cloud Backup

		Returns one snapshot of one sharded cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return GetBackupShardedClusterApiRequest
	*/
	GetBackupShardedCluster(ctx context.Context, groupId string, clusterName string, snapshotId string) GetBackupShardedClusterApiRequest
	/*
		GetBackupShardedCluster Return One Sharded Cluster Cloud Backup


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetBackupShardedClusterApiParams - Parameters for the request
		@return GetBackupShardedClusterApiRequest
	*/
	GetBackupShardedClusterWithParams(ctx context.Context, args *GetBackupShardedClusterApiParams) GetBackupShardedClusterApiRequest

	// Method available only for mocking purposes
	GetBackupShardedClusterExecute(r GetBackupShardedClusterApiRequest) (*DiskBackupShardedClusterSnapshot, *http.Response, error)

	/*
		GetClusterBackupSnapshot Return One Replica Set Cloud Backup

		Returns one snapshot from the specified cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@return GetClusterBackupSnapshotApiRequest
	*/
	GetClusterBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetClusterBackupSnapshotApiRequest
	/*
		GetClusterBackupSnapshot Return One Replica Set Cloud Backup


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterBackupSnapshotApiParams - Parameters for the request
		@return GetClusterBackupSnapshotApiRequest
	*/
	GetClusterBackupSnapshotWithParams(ctx context.Context, args *GetClusterBackupSnapshotApiParams) GetClusterBackupSnapshotApiRequest

	// Method available only for mocking purposes
	GetClusterBackupSnapshotExecute(r GetClusterBackupSnapshotApiRequest) (*DiskBackupReplicaSet, *http.Response, error)

	/*
		GetCompliancePolicy Return Backup Compliance Policy Settings

		Returns the Backup Compliance Policy settings with the specified project. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetCompliancePolicyApiRequest
	*/
	GetCompliancePolicy(ctx context.Context, groupId string) GetCompliancePolicyApiRequest
	/*
		GetCompliancePolicy Return Backup Compliance Policy Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetCompliancePolicyApiParams - Parameters for the request
		@return GetCompliancePolicyApiRequest
	*/
	GetCompliancePolicyWithParams(ctx context.Context, args *GetCompliancePolicyApiParams) GetCompliancePolicyApiRequest

	// Method available only for mocking purposes
	GetCompliancePolicyExecute(r GetCompliancePolicyApiRequest) (*DataProtectionSettings20231001, *http.Response, error)

	/*
		GetExportBucket Return One Snapshot Export Bucket

		Returns one Export Bucket associated with the specified Project. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param exportBucketId Unique 24-hexadecimal character string that identifies the Export Bucket.
		@return GetExportBucketApiRequest
	*/
	GetExportBucket(ctx context.Context, groupId string, exportBucketId string) GetExportBucketApiRequest
	/*
		GetExportBucket Return One Snapshot Export Bucket


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetExportBucketApiParams - Parameters for the request
		@return GetExportBucketApiRequest
	*/
	GetExportBucketWithParams(ctx context.Context, args *GetExportBucketApiParams) GetExportBucketApiRequest

	// Method available only for mocking purposes
	GetExportBucketExecute(r GetExportBucketApiRequest) (*DiskBackupSnapshotExportBucketResponse, *http.Response, error)

	/*
			GetServerlessBackupSnapshot Return One Snapshot of One Serverless Instance

			Returns one snapshot of one serverless instance from the specified project.

		This endpoint can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) API or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Get Flex Backup endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the serverless instance.
			@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
			@return GetServerlessBackupSnapshotApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	GetServerlessBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetServerlessBackupSnapshotApiRequest
	/*
		GetServerlessBackupSnapshot Return One Snapshot of One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServerlessBackupSnapshotApiParams - Parameters for the request
		@return GetServerlessBackupSnapshotApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	GetServerlessBackupSnapshotWithParams(ctx context.Context, args *GetServerlessBackupSnapshotApiParams) GetServerlessBackupSnapshotApiRequest

	// Method available only for mocking purposes
	GetServerlessBackupSnapshotExecute(r GetServerlessBackupSnapshotApiRequest) (*ServerlessBackupSnapshot, *http.Response, error)

	/*
			GetServerlessRestoreJob Return One Restore Job for One Serverless Instance

			Returns one restore job for one serverless instance from the specified project.

		This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Get Flex Backup Restore Job endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the serverless instance.
			@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to return.
			@return GetServerlessRestoreJobApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	GetServerlessRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) GetServerlessRestoreJobApiRequest
	/*
		GetServerlessRestoreJob Return One Restore Job for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetServerlessRestoreJobApiParams - Parameters for the request
		@return GetServerlessRestoreJobApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	GetServerlessRestoreJobWithParams(ctx context.Context, args *GetServerlessRestoreJobApiParams) GetServerlessRestoreJobApiRequest

	// Method available only for mocking purposes
	GetServerlessRestoreJobExecute(r GetServerlessRestoreJobApiRequest) (*ServerlessBackupRestoreJob, *http.Response, error)

	/*
		ListBackupExports Return All Snapshot Export Jobs

		Returns all Cloud Backup Snapshot Export Jobs associated with the specified Atlas cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListBackupExportsApiRequest
	*/
	ListBackupExports(ctx context.Context, groupId string, clusterName string) ListBackupExportsApiRequest
	/*
		ListBackupExports Return All Snapshot Export Jobs


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListBackupExportsApiParams - Parameters for the request
		@return ListBackupExportsApiRequest
	*/
	ListBackupExportsWithParams(ctx context.Context, args *ListBackupExportsApiParams) ListBackupExportsApiRequest

	// Method available only for mocking purposes
	ListBackupExportsExecute(r ListBackupExportsApiRequest) (*PaginatedApiAtlasDiskBackupExportJob, *http.Response, error)

	/*
		ListBackupPrivateEndpoints Return Object Storage Private Endpoints for Cloud Backups for One Cloud Provider in One Project

		Returns the private endpoints of the specified cloud provider for object storage backup operations.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoints to return.
		@return ListBackupPrivateEndpointsApiRequest
	*/
	ListBackupPrivateEndpoints(ctx context.Context, groupId string, cloudProvider string) ListBackupPrivateEndpointsApiRequest
	/*
		ListBackupPrivateEndpoints Return Object Storage Private Endpoints for Cloud Backups for One Cloud Provider in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListBackupPrivateEndpointsApiParams - Parameters for the request
		@return ListBackupPrivateEndpointsApiRequest
	*/
	ListBackupPrivateEndpointsWithParams(ctx context.Context, args *ListBackupPrivateEndpointsApiParams) ListBackupPrivateEndpointsApiRequest

	// Method available only for mocking purposes
	ListBackupPrivateEndpointsExecute(r ListBackupPrivateEndpointsApiRequest) (*PaginatedApiAtlasObjectStoragePrivateEndpointResponse, *http.Response, error)

	/*
		ListBackupRestoreJobs Return All Restore Jobs for One Cluster

		Returns all cloud backup restore jobs for one cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster with the restore jobs you want to return.
		@return ListBackupRestoreJobsApiRequest
	*/
	ListBackupRestoreJobs(ctx context.Context, groupId string, clusterName string) ListBackupRestoreJobsApiRequest
	/*
		ListBackupRestoreJobs Return All Restore Jobs for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListBackupRestoreJobsApiParams - Parameters for the request
		@return ListBackupRestoreJobsApiRequest
	*/
	ListBackupRestoreJobsWithParams(ctx context.Context, args *ListBackupRestoreJobsApiParams) ListBackupRestoreJobsApiRequest

	// Method available only for mocking purposes
	ListBackupRestoreJobsExecute(r ListBackupRestoreJobsApiRequest) (*PaginatedCloudBackupRestoreJob, *http.Response, error)

	/*
		ListBackupShardedClusters Return All Sharded Cluster Cloud Backups

		Returns all snapshots of one sharded cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListBackupShardedClustersApiRequest
	*/
	ListBackupShardedClusters(ctx context.Context, groupId string, clusterName string) ListBackupShardedClustersApiRequest
	/*
		ListBackupShardedClusters Return All Sharded Cluster Cloud Backups


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListBackupShardedClustersApiParams - Parameters for the request
		@return ListBackupShardedClustersApiRequest
	*/
	ListBackupShardedClustersWithParams(ctx context.Context, args *ListBackupShardedClustersApiParams) ListBackupShardedClustersApiRequest

	// Method available only for mocking purposes
	ListBackupShardedClustersExecute(r ListBackupShardedClustersApiRequest) (*PaginatedCloudBackupShardedClusterSnapshot, *http.Response, error)

	/*
		ListBackupSnapshots Return All Replica Set Cloud Backups

		Returns all snapshots of one cluster from the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return ListBackupSnapshotsApiRequest
	*/
	ListBackupSnapshots(ctx context.Context, groupId string, clusterName string) ListBackupSnapshotsApiRequest
	/*
		ListBackupSnapshots Return All Replica Set Cloud Backups


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListBackupSnapshotsApiParams - Parameters for the request
		@return ListBackupSnapshotsApiRequest
	*/
	ListBackupSnapshotsWithParams(ctx context.Context, args *ListBackupSnapshotsApiParams) ListBackupSnapshotsApiRequest

	// Method available only for mocking purposes
	ListBackupSnapshotsExecute(r ListBackupSnapshotsApiRequest) (*PaginatedCloudBackupReplicaSet, *http.Response, error)

	/*
		ListExportBuckets Return All Snapshot Export Buckets

		Returns all Export Buckets associated with the specified Project. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListExportBucketsApiRequest
	*/
	ListExportBuckets(ctx context.Context, groupId string) ListExportBucketsApiRequest
	/*
		ListExportBuckets Return All Snapshot Export Buckets


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListExportBucketsApiParams - Parameters for the request
		@return ListExportBucketsApiRequest
	*/
	ListExportBucketsWithParams(ctx context.Context, args *ListExportBucketsApiParams) ListExportBucketsApiRequest

	// Method available only for mocking purposes
	ListExportBucketsExecute(r ListExportBucketsApiRequest) (*PaginatedBackupSnapshotExportBuckets, *http.Response, error)

	/*
			ListServerlessBackupSnapshots Return All Snapshots of One Serverless Instance

			Returns all snapshots of one serverless instance from the specified project.

		This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the List Flex Backups endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the serverless instance.
			@return ListServerlessBackupSnapshotsApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	ListServerlessBackupSnapshots(ctx context.Context, groupId string, clusterName string) ListServerlessBackupSnapshotsApiRequest
	/*
		ListServerlessBackupSnapshots Return All Snapshots of One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListServerlessBackupSnapshotsApiParams - Parameters for the request
		@return ListServerlessBackupSnapshotsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	ListServerlessBackupSnapshotsWithParams(ctx context.Context, args *ListServerlessBackupSnapshotsApiParams) ListServerlessBackupSnapshotsApiRequest

	// Method available only for mocking purposes
	ListServerlessBackupSnapshotsExecute(r ListServerlessBackupSnapshotsApiRequest) (*PaginatedApiAtlasServerlessBackupSnapshot, *http.Response, error)

	/*
			ListServerlessRestoreJobs Return All Restore Jobs for One Serverless Instance

			Returns all restore jobs for one serverless instance from the specified project.

		This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the List Flex Backup Restore Jobs endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the serverless instance.
			@return ListServerlessRestoreJobsApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	ListServerlessRestoreJobs(ctx context.Context, groupId string, clusterName string) ListServerlessRestoreJobsApiRequest
	/*
		ListServerlessRestoreJobs Return All Restore Jobs for One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListServerlessRestoreJobsApiParams - Parameters for the request
		@return ListServerlessRestoreJobsApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for CloudBackupsApi
	*/
	ListServerlessRestoreJobsWithParams(ctx context.Context, args *ListServerlessRestoreJobsApiParams) ListServerlessRestoreJobsApiRequest

	// Method available only for mocking purposes
	ListServerlessRestoreJobsExecute(r ListServerlessRestoreJobsApiRequest) (*PaginatedApiAtlasServerlessBackupRestoreJob, *http.Response, error)

	/*
		TakeSnapshots Take One On-Demand Snapshot

		Takes one on-demand snapshot for the specified cluster. Atlas takes on-demand snapshots immediately and scheduled snapshots at regular intervals. If an on-demand snapshot with a status of `queued` or `inProgress` exists, before taking another snapshot, wait until Atlas completes processing the previously taken on-demand snapshot.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param diskBackupOnDemandSnapshotRequest Takes one on-demand snapshot.
		@return TakeSnapshotsApiRequest
	*/
	TakeSnapshots(ctx context.Context, groupId string, clusterName string, diskBackupOnDemandSnapshotRequest *DiskBackupOnDemandSnapshotRequest) TakeSnapshotsApiRequest
	/*
		TakeSnapshots Take One On-Demand Snapshot


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param TakeSnapshotsApiParams - Parameters for the request
		@return TakeSnapshotsApiRequest
	*/
	TakeSnapshotsWithParams(ctx context.Context, args *TakeSnapshotsApiParams) TakeSnapshotsApiRequest

	// Method available only for mocking purposes
	TakeSnapshotsExecute(r TakeSnapshotsApiRequest) (*DiskBackupSnapshot, *http.Response, error)

	/*
		UpdateBackupExportBucket Update One Export Bucket Private Networking Settings

		Updates the private networking settings for one snapshot export bucket in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param exportBucketId Unique 24-hexadecimal character string that identifies the snapshot export bucket.
		@param updateRequirePrivateNetworkingRequest Updates the private networking requirement for the snapshot export bucket.
		@return UpdateBackupExportBucketApiRequest
	*/
	UpdateBackupExportBucket(ctx context.Context, groupId string, exportBucketId string, updateRequirePrivateNetworkingRequest *UpdateRequirePrivateNetworkingRequest) UpdateBackupExportBucketApiRequest
	/*
		UpdateBackupExportBucket Update One Export Bucket Private Networking Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateBackupExportBucketApiParams - Parameters for the request
		@return UpdateBackupExportBucketApiRequest
	*/
	UpdateBackupExportBucketWithParams(ctx context.Context, args *UpdateBackupExportBucketApiParams) UpdateBackupExportBucketApiRequest

	// Method available only for mocking purposes
	UpdateBackupExportBucketExecute(r UpdateBackupExportBucketApiRequest) (*DiskBackupSnapshotAWSExportBucketResponse, *http.Response, error)

	/*
		UpdateBackupSchedule Update Cloud Backup Schedule for One Cluster

		Updates the cloud backup schedule for one cluster within the specified project. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param diskBackupSnapshotSchedule20240805 Updates the cloud backup schedule for one cluster within the specified project.  **Note**: In the request body, provide only the fields that you want to update.
		@return UpdateBackupScheduleApiRequest
	*/
	UpdateBackupSchedule(ctx context.Context, groupId string, clusterName string, diskBackupSnapshotSchedule20240805 *DiskBackupSnapshotSchedule20240805) UpdateBackupScheduleApiRequest
	/*
		UpdateBackupSchedule Update Cloud Backup Schedule for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateBackupScheduleApiParams - Parameters for the request
		@return UpdateBackupScheduleApiRequest
	*/
	UpdateBackupScheduleWithParams(ctx context.Context, args *UpdateBackupScheduleApiParams) UpdateBackupScheduleApiRequest

	// Method available only for mocking purposes
	UpdateBackupScheduleExecute(r UpdateBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error)

	/*
		UpdateBackupSnapshot Update Expiration Date for One Cloud Backup

		Changes the expiration date for one cloud backup snapshot for one cluster in the specified project, the requesting Service Account or API Key must have the Project Backup Manager role.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
		@param backupSnapshotRetention Changes the expiration date for one cloud backup snapshot for one cluster in the specified project.
		@return UpdateBackupSnapshotApiRequest
	*/
	UpdateBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string, backupSnapshotRetention *BackupSnapshotRetention) UpdateBackupSnapshotApiRequest
	/*
		UpdateBackupSnapshot Update Expiration Date for One Cloud Backup


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateBackupSnapshotApiParams - Parameters for the request
		@return UpdateBackupSnapshotApiRequest
	*/
	UpdateBackupSnapshotWithParams(ctx context.Context, args *UpdateBackupSnapshotApiParams) UpdateBackupSnapshotApiRequest

	// Method available only for mocking purposes
	UpdateBackupSnapshotExecute(r UpdateBackupSnapshotApiRequest) (*DiskBackupReplicaSet, *http.Response, error)

	/*
		UpdateCompliancePolicy Update Backup Compliance Policy Settings

		Updates the Backup Compliance Policy settings for the specified project. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param dataProtectionSettings20231001 The new Backup Compliance Policy settings.
		@return UpdateCompliancePolicyApiRequest
	*/
	UpdateCompliancePolicy(ctx context.Context, groupId string, dataProtectionSettings20231001 *DataProtectionSettings20231001) UpdateCompliancePolicyApiRequest
	/*
		UpdateCompliancePolicy Update Backup Compliance Policy Settings


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateCompliancePolicyApiParams - Parameters for the request
		@return UpdateCompliancePolicyApiRequest
	*/
	UpdateCompliancePolicyWithParams(ctx context.Context, args *UpdateCompliancePolicyApiParams) UpdateCompliancePolicyApiRequest

	// Method available only for mocking purposes
	UpdateCompliancePolicyExecute(r UpdateCompliancePolicyApiRequest) (*DataProtectionSettings20231001, *http.Response, error)
}

// CloudBackupsApiService CloudBackupsApi service
type CloudBackupsApiService service

type CancelBackupRestoreJobApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	restoreJobId string
}

type CancelBackupRestoreJobApiParams struct {
	GroupId      string
	ClusterName  string
	RestoreJobId string
}

func (a *CloudBackupsApiService) CancelBackupRestoreJobWithParams(ctx context.Context, args *CancelBackupRestoreJobApiParams) CancelBackupRestoreJobApiRequest {
	return CancelBackupRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		restoreJobId: args.RestoreJobId,
	}
}

func (r CancelBackupRestoreJobApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.CancelBackupRestoreJobExecute(r)
}

/*
CancelBackupRestoreJob Cancel One Restore Job for One Cluster

Cancels one cloud backup restore job of one cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to remove.
	@return CancelBackupRestoreJobApiRequest
*/
func (a *CloudBackupsApiService) CancelBackupRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) CancelBackupRestoreJobApiRequest {
	return CancelBackupRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		clusterName:  clusterName,
		restoreJobId: restoreJobId,
	}
}

// CancelBackupRestoreJobExecute executes the request
func (a *CloudBackupsApiService) CancelBackupRestoreJobExecute(r CancelBackupRestoreJobApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CancelBackupRestoreJob")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/restoreJobs/{restoreJobId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.restoreJobId == "" {
		return nil, reportError("restoreJobId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"restoreJobId"+"}", url.PathEscape(r.restoreJobId), -1)

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

type CreateBackupExportApiRequest struct {
	ctx                        context.Context
	ApiService                 CloudBackupsApi
	groupId                    string
	clusterName                string
	diskBackupExportJobRequest *DiskBackupExportJobRequest
}

type CreateBackupExportApiParams struct {
	GroupId                    string
	ClusterName                string
	DiskBackupExportJobRequest *DiskBackupExportJobRequest
}

func (a *CloudBackupsApiService) CreateBackupExportWithParams(ctx context.Context, args *CreateBackupExportApiParams) CreateBackupExportApiRequest {
	return CreateBackupExportApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		clusterName:                args.ClusterName,
		diskBackupExportJobRequest: args.DiskBackupExportJobRequest,
	}
}

func (r CreateBackupExportApiRequest) Execute() (*DiskBackupExportJob, *http.Response, error) {
	return r.ApiService.CreateBackupExportExecute(r)
}

/*
CreateBackupExport Create One Snapshot Export Job

Exports one backup Snapshot for dedicated Atlas cluster using Cloud Backups to an Export Bucket.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return CreateBackupExportApiRequest
*/
func (a *CloudBackupsApiService) CreateBackupExport(ctx context.Context, groupId string, clusterName string, diskBackupExportJobRequest *DiskBackupExportJobRequest) CreateBackupExportApiRequest {
	return CreateBackupExportApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterName:                clusterName,
		diskBackupExportJobRequest: diskBackupExportJobRequest,
	}
}

// CreateBackupExportExecute executes the request
//
//	@return DiskBackupExportJob
func (a *CloudBackupsApiService) CreateBackupExportExecute(r CreateBackupExportApiRequest) (*DiskBackupExportJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupExportJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CreateBackupExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/exports"
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
	if r.diskBackupExportJobRequest == nil {
		return localVarReturnValue, nil, reportError("diskBackupExportJobRequest is required and must be specified")
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
	localVarPostBody = r.diskBackupExportJobRequest
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

type CreateBackupPrivateEndpointApiRequest struct {
	ctx                                 context.Context
	ApiService                          CloudBackupsApi
	groupId                             string
	cloudProvider                       string
	objectStoragePrivateEndpointRequest *ObjectStoragePrivateEndpointRequest
}

type CreateBackupPrivateEndpointApiParams struct {
	GroupId                             string
	CloudProvider                       string
	ObjectStoragePrivateEndpointRequest *ObjectStoragePrivateEndpointRequest
}

func (a *CloudBackupsApiService) CreateBackupPrivateEndpointWithParams(ctx context.Context, args *CreateBackupPrivateEndpointApiParams) CreateBackupPrivateEndpointApiRequest {
	return CreateBackupPrivateEndpointApiRequest{
		ApiService:                          a,
		ctx:                                 ctx,
		groupId:                             args.GroupId,
		cloudProvider:                       args.CloudProvider,
		objectStoragePrivateEndpointRequest: args.ObjectStoragePrivateEndpointRequest,
	}
}

func (r CreateBackupPrivateEndpointApiRequest) Execute() (*ObjectStoragePrivateEndpointResponse, *http.Response, error) {
	return r.ApiService.CreateBackupPrivateEndpointExecute(r)
}

/*
CreateBackupPrivateEndpoint Create One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project

Creates a private endpoint in the specified region for secure, private connectivity between Atlas and cloud provider object storage services for backup operations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoint to create.
	@return CreateBackupPrivateEndpointApiRequest
*/
func (a *CloudBackupsApiService) CreateBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, objectStoragePrivateEndpointRequest *ObjectStoragePrivateEndpointRequest) CreateBackupPrivateEndpointApiRequest {
	return CreateBackupPrivateEndpointApiRequest{
		ApiService:                          a,
		ctx:                                 ctx,
		groupId:                             groupId,
		cloudProvider:                       cloudProvider,
		objectStoragePrivateEndpointRequest: objectStoragePrivateEndpointRequest,
	}
}

// CreateBackupPrivateEndpointExecute executes the request
//
//	@return ObjectStoragePrivateEndpointResponse
func (a *CloudBackupsApiService) CreateBackupPrivateEndpointExecute(r CreateBackupPrivateEndpointApiRequest) (*ObjectStoragePrivateEndpointResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ObjectStoragePrivateEndpointResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CreateBackupPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/{cloudProvider}/privateEndpoints"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.objectStoragePrivateEndpointRequest == nil {
		return localVarReturnValue, nil, reportError("objectStoragePrivateEndpointRequest is required and must be specified")
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
	localVarPostBody = r.objectStoragePrivateEndpointRequest
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

type CreateBackupRestoreJobApiRequest struct {
	ctx                          context.Context
	ApiService                   CloudBackupsApi
	groupId                      string
	clusterName                  string
	diskBackupSnapshotRestoreJob *DiskBackupSnapshotRestoreJob
}

type CreateBackupRestoreJobApiParams struct {
	GroupId                      string
	ClusterName                  string
	DiskBackupSnapshotRestoreJob *DiskBackupSnapshotRestoreJob
}

func (a *CloudBackupsApiService) CreateBackupRestoreJobWithParams(ctx context.Context, args *CreateBackupRestoreJobApiParams) CreateBackupRestoreJobApiRequest {
	return CreateBackupRestoreJobApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		clusterName:                  args.ClusterName,
		diskBackupSnapshotRestoreJob: args.DiskBackupSnapshotRestoreJob,
	}
}

func (r CreateBackupRestoreJobApiRequest) Execute() (*DiskBackupSnapshotRestoreJob, *http.Response, error) {
	return r.ApiService.CreateBackupRestoreJobExecute(r)
}

/*
CreateBackupRestoreJob Create One Restore Job of One Cluster

Restores one snapshot of one cluster from the specified project. Atlas takes on-demand snapshots immediately and scheduled snapshots at regular intervals. If an on-demand snapshot with a status of `queued` or `inProgress` exists, before taking another snapshot, wait until Atlas completes processing the previously taken on-demand snapshot.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return CreateBackupRestoreJobApiRequest
*/
func (a *CloudBackupsApiService) CreateBackupRestoreJob(ctx context.Context, groupId string, clusterName string, diskBackupSnapshotRestoreJob *DiskBackupSnapshotRestoreJob) CreateBackupRestoreJobApiRequest {
	return CreateBackupRestoreJobApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		clusterName:                  clusterName,
		diskBackupSnapshotRestoreJob: diskBackupSnapshotRestoreJob,
	}
}

// CreateBackupRestoreJobExecute executes the request
//
//	@return DiskBackupSnapshotRestoreJob
func (a *CloudBackupsApiService) CreateBackupRestoreJobExecute(r CreateBackupRestoreJobApiRequest) (*DiskBackupSnapshotRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CreateBackupRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/restoreJobs"
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
	if r.diskBackupSnapshotRestoreJob == nil {
		return localVarReturnValue, nil, reportError("diskBackupSnapshotRestoreJob is required and must be specified")
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
	localVarPostBody = r.diskBackupSnapshotRestoreJob
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

type CreateExportBucketApiRequest struct {
	ctx                                   context.Context
	ApiService                            CloudBackupsApi
	groupId                               string
	diskBackupSnapshotExportBucketRequest *DiskBackupSnapshotExportBucketRequest
}

type CreateExportBucketApiParams struct {
	GroupId                               string
	DiskBackupSnapshotExportBucketRequest *DiskBackupSnapshotExportBucketRequest
}

func (a *CloudBackupsApiService) CreateExportBucketWithParams(ctx context.Context, args *CreateExportBucketApiParams) CreateExportBucketApiRequest {
	return CreateExportBucketApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               args.GroupId,
		diskBackupSnapshotExportBucketRequest: args.DiskBackupSnapshotExportBucketRequest,
	}
}

func (r CreateExportBucketApiRequest) Execute() (*DiskBackupSnapshotExportBucketResponse, *http.Response, error) {
	return r.ApiService.CreateExportBucketExecute(r)
}

/*
CreateExportBucket Create One Snapshot Export Bucket

Creates a Snapshot Export Bucket for an AWS S3 Bucket, Azure Blob Storage Container, or Google Cloud Storage Bucket. Once created, an snapshots can be exported to the Export Bucket and its referenced AWS S3 Bucket, Azure Blob Storage Container, or Google Cloud Storage Bucket. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateExportBucketApiRequest
*/
func (a *CloudBackupsApiService) CreateExportBucket(ctx context.Context, groupId string, diskBackupSnapshotExportBucketRequest *DiskBackupSnapshotExportBucketRequest) CreateExportBucketApiRequest {
	return CreateExportBucketApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               groupId,
		diskBackupSnapshotExportBucketRequest: diskBackupSnapshotExportBucketRequest,
	}
}

// CreateExportBucketExecute executes the request
//
//	@return DiskBackupSnapshotExportBucketResponse
func (a *CloudBackupsApiService) CreateExportBucketExecute(r CreateExportBucketApiRequest) (*DiskBackupSnapshotExportBucketResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotExportBucketResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CreateExportBucket")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/exportBuckets"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.diskBackupSnapshotExportBucketRequest == nil {
		return localVarReturnValue, nil, reportError("diskBackupSnapshotExportBucketRequest is required and must be specified")
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
	localVarPostBody = r.diskBackupSnapshotExportBucketRequest
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

type CreateServerlessRestoreJobApiRequest struct {
	ctx                        context.Context
	ApiService                 CloudBackupsApi
	groupId                    string
	clusterName                string
	serverlessBackupRestoreJob *ServerlessBackupRestoreJob
}

type CreateServerlessRestoreJobApiParams struct {
	GroupId                    string
	ClusterName                string
	ServerlessBackupRestoreJob *ServerlessBackupRestoreJob
}

func (a *CloudBackupsApiService) CreateServerlessRestoreJobWithParams(ctx context.Context, args *CreateServerlessRestoreJobApiParams) CreateServerlessRestoreJobApiRequest {
	return CreateServerlessRestoreJobApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		clusterName:                args.ClusterName,
		serverlessBackupRestoreJob: args.ServerlessBackupRestoreJob,
	}
}

func (r CreateServerlessRestoreJobApiRequest) Execute() (*ServerlessBackupRestoreJob, *http.Response, error) {
	return r.ApiService.CreateServerlessRestoreJobExecute(r)
}

/*
CreateServerlessRestoreJob Create One Restore Job for One Serverless Instance

Restores one snapshot of one serverless instance from the specified project.

This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Create Flex Backup Restore Job endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the serverless instance whose snapshot you want to restore.
	@return CreateServerlessRestoreJobApiRequest

Deprecated
*/
func (a *CloudBackupsApiService) CreateServerlessRestoreJob(ctx context.Context, groupId string, clusterName string, serverlessBackupRestoreJob *ServerlessBackupRestoreJob) CreateServerlessRestoreJobApiRequest {
	return CreateServerlessRestoreJobApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterName:                clusterName,
		serverlessBackupRestoreJob: serverlessBackupRestoreJob,
	}
}

// CreateServerlessRestoreJobExecute executes the request
//
//	@return ServerlessBackupRestoreJob
//
// Deprecated
func (a *CloudBackupsApiService) CreateServerlessRestoreJobExecute(r CreateServerlessRestoreJobApiRequest) (*ServerlessBackupRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessBackupRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.CreateServerlessRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/backup/restoreJobs"
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
	if r.serverlessBackupRestoreJob == nil {
		return localVarReturnValue, nil, reportError("serverlessBackupRestoreJob is required and must be specified")
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
	localVarPostBody = r.serverlessBackupRestoreJob
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

type DeleteBackupPrivateEndpointApiRequest struct {
	ctx           context.Context
	ApiService    CloudBackupsApi
	groupId       string
	cloudProvider string
	endpointId    string
}

type DeleteBackupPrivateEndpointApiParams struct {
	GroupId       string
	CloudProvider string
	EndpointId    string
}

func (a *CloudBackupsApiService) DeleteBackupPrivateEndpointWithParams(ctx context.Context, args *DeleteBackupPrivateEndpointApiParams) DeleteBackupPrivateEndpointApiRequest {
	return DeleteBackupPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		endpointId:    args.EndpointId,
	}
}

func (r DeleteBackupPrivateEndpointApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteBackupPrivateEndpointExecute(r)
}

/*
DeleteBackupPrivateEndpoint Delete One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider from One Project

Deletes one private endpoint, identified by its ID, for object storage backup operations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint to delete.
	@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint to delete.
	@return DeleteBackupPrivateEndpointApiRequest
*/
func (a *CloudBackupsApiService) DeleteBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) DeleteBackupPrivateEndpointApiRequest {
	return DeleteBackupPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
		endpointId:    endpointId,
	}
}

// DeleteBackupPrivateEndpointExecute executes the request
func (a *CloudBackupsApiService) DeleteBackupPrivateEndpointExecute(r DeleteBackupPrivateEndpointApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DeleteBackupPrivateEndpoint")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/{cloudProvider}/privateEndpoints/{endpointId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-11-13+json"}

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

type DeleteBackupShardedClusterApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	snapshotId  string
}

type DeleteBackupShardedClusterApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *CloudBackupsApiService) DeleteBackupShardedClusterWithParams(ctx context.Context, args *DeleteBackupShardedClusterApiParams) DeleteBackupShardedClusterApiRequest {
	return DeleteBackupShardedClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r DeleteBackupShardedClusterApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteBackupShardedClusterExecute(r)
}

/*
DeleteBackupShardedCluster Remove One Sharded Cluster Cloud Backup

Removes one snapshot of one sharded cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return DeleteBackupShardedClusterApiRequest
*/
func (a *CloudBackupsApiService) DeleteBackupShardedCluster(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteBackupShardedClusterApiRequest {
	return DeleteBackupShardedClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// DeleteBackupShardedClusterExecute executes the request
func (a *CloudBackupsApiService) DeleteBackupShardedClusterExecute(r DeleteBackupShardedClusterApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DeleteBackupShardedCluster")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/shardedCluster/{snapshotId}"
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

type DeleteClusterBackupScheduleApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
}

type DeleteClusterBackupScheduleApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *CloudBackupsApiService) DeleteClusterBackupScheduleWithParams(ctx context.Context, args *DeleteClusterBackupScheduleApiParams) DeleteClusterBackupScheduleApiRequest {
	return DeleteClusterBackupScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r DeleteClusterBackupScheduleApiRequest) Execute() (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	return r.ApiService.DeleteClusterBackupScheduleExecute(r)
}

/*
DeleteClusterBackupSchedule Remove All Cloud Backup Schedules

Removes all cloud backup schedules for the specified cluster. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return DeleteClusterBackupScheduleApiRequest
*/
func (a *CloudBackupsApiService) DeleteClusterBackupSchedule(ctx context.Context, groupId string, clusterName string) DeleteClusterBackupScheduleApiRequest {
	return DeleteClusterBackupScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// DeleteClusterBackupScheduleExecute executes the request
//
//	@return DiskBackupSnapshotSchedule20240805
func (a *CloudBackupsApiService) DeleteClusterBackupScheduleExecute(r DeleteClusterBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotSchedule20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DeleteClusterBackupSchedule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/schedule"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type DeleteClusterBackupSnapshotApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	snapshotId  string
}

type DeleteClusterBackupSnapshotApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *CloudBackupsApiService) DeleteClusterBackupSnapshotWithParams(ctx context.Context, args *DeleteClusterBackupSnapshotApiParams) DeleteClusterBackupSnapshotApiRequest {
	return DeleteClusterBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r DeleteClusterBackupSnapshotApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterBackupSnapshotExecute(r)
}

/*
DeleteClusterBackupSnapshot Remove One Replica Set Cloud Backup

Removes the specified snapshot.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return DeleteClusterBackupSnapshotApiRequest
*/
func (a *CloudBackupsApiService) DeleteClusterBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) DeleteClusterBackupSnapshotApiRequest {
	return DeleteClusterBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// DeleteClusterBackupSnapshotExecute executes the request
func (a *CloudBackupsApiService) DeleteClusterBackupSnapshotExecute(r DeleteClusterBackupSnapshotApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DeleteClusterBackupSnapshot")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/{snapshotId}"
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

type DeleteExportBucketApiRequest struct {
	ctx            context.Context
	ApiService     CloudBackupsApi
	groupId        string
	exportBucketId string
}

type DeleteExportBucketApiParams struct {
	GroupId        string
	ExportBucketId string
}

func (a *CloudBackupsApiService) DeleteExportBucketWithParams(ctx context.Context, args *DeleteExportBucketApiParams) DeleteExportBucketApiRequest {
	return DeleteExportBucketApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		exportBucketId: args.ExportBucketId,
	}
}

func (r DeleteExportBucketApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteExportBucketExecute(r)
}

/*
DeleteExportBucket Delete One Snapshot Export Bucket

Deletes an Export Bucket. Auto export must be disabled on all clusters in this Project exporting to this Export Bucket before revoking access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param exportBucketId Unique 24-hexadecimal character string that identifies the Export Bucket.
	@return DeleteExportBucketApiRequest
*/
func (a *CloudBackupsApiService) DeleteExportBucket(ctx context.Context, groupId string, exportBucketId string) DeleteExportBucketApiRequest {
	return DeleteExportBucketApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		exportBucketId: exportBucketId,
	}
}

// DeleteExportBucketExecute executes the request
func (a *CloudBackupsApiService) DeleteExportBucketExecute(r DeleteExportBucketApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DeleteExportBucket")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/exportBuckets/{exportBucketId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.exportBucketId == "" {
		return nil, reportError("exportBucketId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"exportBucketId"+"}", url.PathEscape(r.exportBucketId), -1)

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

type DisableCompliancePolicyApiRequest struct {
	ctx        context.Context
	ApiService CloudBackupsApi
	groupId    string
}

type DisableCompliancePolicyApiParams struct {
	GroupId string
}

func (a *CloudBackupsApiService) DisableCompliancePolicyWithParams(ctx context.Context, args *DisableCompliancePolicyApiParams) DisableCompliancePolicyApiRequest {
	return DisableCompliancePolicyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r DisableCompliancePolicyApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DisableCompliancePolicyExecute(r)
}

/*
DisableCompliancePolicy Disable Backup Compliance Policy Settings

Disables the Backup Compliance Policy settings with the specified project. As a prerequisite, a support ticket needs to be file first, instructions in https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DisableCompliancePolicyApiRequest
*/
func (a *CloudBackupsApiService) DisableCompliancePolicy(ctx context.Context, groupId string) DisableCompliancePolicyApiRequest {
	return DisableCompliancePolicyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// DisableCompliancePolicyExecute executes the request
func (a *CloudBackupsApiService) DisableCompliancePolicyExecute(r DisableCompliancePolicyApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.DisableCompliancePolicy")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backupCompliancePolicy"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-11-13+json"}

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

type GetBackupExportApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	exportId    string
}

type GetBackupExportApiParams struct {
	GroupId     string
	ClusterName string
	ExportId    string
}

func (a *CloudBackupsApiService) GetBackupExportWithParams(ctx context.Context, args *GetBackupExportApiParams) GetBackupExportApiRequest {
	return GetBackupExportApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		exportId:    args.ExportId,
	}
}

func (r GetBackupExportApiRequest) Execute() (*DiskBackupExportJob, *http.Response, error) {
	return r.ApiService.GetBackupExportExecute(r)
}

/*
GetBackupExport Return One Snapshot Export Job

Returns one Cloud Backup Snapshot Export Job associated with the specified Atlas cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param exportId Unique 24-hexadecimal character string that identifies the Export Job.
	@return GetBackupExportApiRequest
*/
func (a *CloudBackupsApiService) GetBackupExport(ctx context.Context, groupId string, clusterName string, exportId string) GetBackupExportApiRequest {
	return GetBackupExportApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		exportId:    exportId,
	}
}

// GetBackupExportExecute executes the request
//
//	@return DiskBackupExportJob
func (a *CloudBackupsApiService) GetBackupExportExecute(r GetBackupExportApiRequest) (*DiskBackupExportJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupExportJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetBackupExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/exports/{exportId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.exportId == "" {
		return localVarReturnValue, nil, reportError("exportId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(r.exportId), -1)

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

type GetBackupPrivateEndpointApiRequest struct {
	ctx           context.Context
	ApiService    CloudBackupsApi
	groupId       string
	cloudProvider string
	endpointId    string
}

type GetBackupPrivateEndpointApiParams struct {
	GroupId       string
	CloudProvider string
	EndpointId    string
}

func (a *CloudBackupsApiService) GetBackupPrivateEndpointWithParams(ctx context.Context, args *GetBackupPrivateEndpointApiParams) GetBackupPrivateEndpointApiRequest {
	return GetBackupPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		endpointId:    args.EndpointId,
	}
}

func (r GetBackupPrivateEndpointApiRequest) Execute() (*ObjectStoragePrivateEndpointResponse, *http.Response, error) {
	return r.ApiService.GetBackupPrivateEndpointExecute(r)
}

/*
GetBackupPrivateEndpoint Return One Object Storage Private Endpoint for Cloud Backups for One Cloud Provider in One Project

Returns one private endpoint, identified by its ID, for object storage backup operations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider of the private endpoint.
	@param endpointId Unique 24-hexadecimal digit string that identifies the private endpoint.
	@return GetBackupPrivateEndpointApiRequest
*/
func (a *CloudBackupsApiService) GetBackupPrivateEndpoint(ctx context.Context, groupId string, cloudProvider string, endpointId string) GetBackupPrivateEndpointApiRequest {
	return GetBackupPrivateEndpointApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
		endpointId:    endpointId,
	}
}

// GetBackupPrivateEndpointExecute executes the request
//
//	@return ObjectStoragePrivateEndpointResponse
func (a *CloudBackupsApiService) GetBackupPrivateEndpointExecute(r GetBackupPrivateEndpointApiRequest) (*ObjectStoragePrivateEndpointResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ObjectStoragePrivateEndpointResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetBackupPrivateEndpoint")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/{cloudProvider}/privateEndpoints/{endpointId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)
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

type GetBackupRestoreJobApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	restoreJobId string
}

type GetBackupRestoreJobApiParams struct {
	GroupId      string
	ClusterName  string
	RestoreJobId string
}

func (a *CloudBackupsApiService) GetBackupRestoreJobWithParams(ctx context.Context, args *GetBackupRestoreJobApiParams) GetBackupRestoreJobApiRequest {
	return GetBackupRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		restoreJobId: args.RestoreJobId,
	}
}

func (r GetBackupRestoreJobApiRequest) Execute() (*DiskBackupSnapshotRestoreJob, *http.Response, error) {
	return r.ApiService.GetBackupRestoreJobExecute(r)
}

/*
GetBackupRestoreJob Return One Restore Job for One Cluster

Returns one cloud backup restore job for one cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the restore jobs you want to return.
	@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to return.
	@return GetBackupRestoreJobApiRequest
*/
func (a *CloudBackupsApiService) GetBackupRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) GetBackupRestoreJobApiRequest {
	return GetBackupRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		clusterName:  clusterName,
		restoreJobId: restoreJobId,
	}
}

// GetBackupRestoreJobExecute executes the request
//
//	@return DiskBackupSnapshotRestoreJob
func (a *CloudBackupsApiService) GetBackupRestoreJobExecute(r GetBackupRestoreJobApiRequest) (*DiskBackupSnapshotRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetBackupRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/restoreJobs/{restoreJobId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.restoreJobId == "" {
		return localVarReturnValue, nil, reportError("restoreJobId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"restoreJobId"+"}", url.PathEscape(r.restoreJobId), -1)

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

type GetBackupScheduleApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
}

type GetBackupScheduleApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *CloudBackupsApiService) GetBackupScheduleWithParams(ctx context.Context, args *GetBackupScheduleApiParams) GetBackupScheduleApiRequest {
	return GetBackupScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetBackupScheduleApiRequest) Execute() (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	return r.ApiService.GetBackupScheduleExecute(r)
}

/*
GetBackupSchedule Return One Cloud Backup Schedule

Returns the cloud backup schedule for the specified cluster within the specified project. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return GetBackupScheduleApiRequest
*/
func (a *CloudBackupsApiService) GetBackupSchedule(ctx context.Context, groupId string, clusterName string) GetBackupScheduleApiRequest {
	return GetBackupScheduleApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetBackupScheduleExecute executes the request
//
//	@return DiskBackupSnapshotSchedule20240805
func (a *CloudBackupsApiService) GetBackupScheduleExecute(r GetBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotSchedule20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetBackupSchedule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/schedule"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

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

type GetBackupShardedClusterApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	snapshotId  string
}

type GetBackupShardedClusterApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *CloudBackupsApiService) GetBackupShardedClusterWithParams(ctx context.Context, args *GetBackupShardedClusterApiParams) GetBackupShardedClusterApiRequest {
	return GetBackupShardedClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r GetBackupShardedClusterApiRequest) Execute() (*DiskBackupShardedClusterSnapshot, *http.Response, error) {
	return r.ApiService.GetBackupShardedClusterExecute(r)
}

/*
GetBackupShardedCluster Return One Sharded Cluster Cloud Backup

Returns one snapshot of one sharded cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return GetBackupShardedClusterApiRequest
*/
func (a *CloudBackupsApiService) GetBackupShardedCluster(ctx context.Context, groupId string, clusterName string, snapshotId string) GetBackupShardedClusterApiRequest {
	return GetBackupShardedClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// GetBackupShardedClusterExecute executes the request
//
//	@return DiskBackupShardedClusterSnapshot
func (a *CloudBackupsApiService) GetBackupShardedClusterExecute(r GetBackupShardedClusterApiRequest) (*DiskBackupShardedClusterSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupShardedClusterSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetBackupShardedCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/shardedCluster/{snapshotId}"
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

type GetClusterBackupSnapshotApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	snapshotId  string
}

type GetClusterBackupSnapshotApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *CloudBackupsApiService) GetClusterBackupSnapshotWithParams(ctx context.Context, args *GetClusterBackupSnapshotApiParams) GetClusterBackupSnapshotApiRequest {
	return GetClusterBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r GetClusterBackupSnapshotApiRequest) Execute() (*DiskBackupReplicaSet, *http.Response, error) {
	return r.ApiService.GetClusterBackupSnapshotExecute(r)
}

/*
GetClusterBackupSnapshot Return One Replica Set Cloud Backup

Returns one snapshot from the specified cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return GetClusterBackupSnapshotApiRequest
*/
func (a *CloudBackupsApiService) GetClusterBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetClusterBackupSnapshotApiRequest {
	return GetClusterBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// GetClusterBackupSnapshotExecute executes the request
//
//	@return DiskBackupReplicaSet
func (a *CloudBackupsApiService) GetClusterBackupSnapshotExecute(r GetClusterBackupSnapshotApiRequest) (*DiskBackupReplicaSet, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupReplicaSet
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetClusterBackupSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/{snapshotId}"
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

type GetCompliancePolicyApiRequest struct {
	ctx        context.Context
	ApiService CloudBackupsApi
	groupId    string
}

type GetCompliancePolicyApiParams struct {
	GroupId string
}

func (a *CloudBackupsApiService) GetCompliancePolicyWithParams(ctx context.Context, args *GetCompliancePolicyApiParams) GetCompliancePolicyApiRequest {
	return GetCompliancePolicyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r GetCompliancePolicyApiRequest) Execute() (*DataProtectionSettings20231001, *http.Response, error) {
	return r.ApiService.GetCompliancePolicyExecute(r)
}

/*
GetCompliancePolicy Return Backup Compliance Policy Settings

Returns the Backup Compliance Policy settings with the specified project. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetCompliancePolicyApiRequest
*/
func (a *CloudBackupsApiService) GetCompliancePolicy(ctx context.Context, groupId string) GetCompliancePolicyApiRequest {
	return GetCompliancePolicyApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetCompliancePolicyExecute executes the request
//
//	@return DataProtectionSettings20231001
func (a *CloudBackupsApiService) GetCompliancePolicyExecute(r GetCompliancePolicyApiRequest) (*DataProtectionSettings20231001, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataProtectionSettings20231001
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetCompliancePolicy")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backupCompliancePolicy"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-10-01+json"}

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

type GetExportBucketApiRequest struct {
	ctx            context.Context
	ApiService     CloudBackupsApi
	groupId        string
	exportBucketId string
}

type GetExportBucketApiParams struct {
	GroupId        string
	ExportBucketId string
}

func (a *CloudBackupsApiService) GetExportBucketWithParams(ctx context.Context, args *GetExportBucketApiParams) GetExportBucketApiRequest {
	return GetExportBucketApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		exportBucketId: args.ExportBucketId,
	}
}

func (r GetExportBucketApiRequest) Execute() (*DiskBackupSnapshotExportBucketResponse, *http.Response, error) {
	return r.ApiService.GetExportBucketExecute(r)
}

/*
GetExportBucket Return One Snapshot Export Bucket

Returns one Export Bucket associated with the specified Project. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param exportBucketId Unique 24-hexadecimal character string that identifies the Export Bucket.
	@return GetExportBucketApiRequest
*/
func (a *CloudBackupsApiService) GetExportBucket(ctx context.Context, groupId string, exportBucketId string) GetExportBucketApiRequest {
	return GetExportBucketApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		exportBucketId: exportBucketId,
	}
}

// GetExportBucketExecute executes the request
//
//	@return DiskBackupSnapshotExportBucketResponse
func (a *CloudBackupsApiService) GetExportBucketExecute(r GetExportBucketApiRequest) (*DiskBackupSnapshotExportBucketResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotExportBucketResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetExportBucket")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/exportBuckets/{exportBucketId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.exportBucketId == "" {
		return localVarReturnValue, nil, reportError("exportBucketId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"exportBucketId"+"}", url.PathEscape(r.exportBucketId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-05-30+json"}

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

type GetServerlessBackupSnapshotApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
	snapshotId  string
}

type GetServerlessBackupSnapshotApiParams struct {
	GroupId     string
	ClusterName string
	SnapshotId  string
}

func (a *CloudBackupsApiService) GetServerlessBackupSnapshotWithParams(ctx context.Context, args *GetServerlessBackupSnapshotApiParams) GetServerlessBackupSnapshotApiRequest {
	return GetServerlessBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		snapshotId:  args.SnapshotId,
	}
}

func (r GetServerlessBackupSnapshotApiRequest) Execute() (*ServerlessBackupSnapshot, *http.Response, error) {
	return r.ApiService.GetServerlessBackupSnapshotExecute(r)
}

/*
GetServerlessBackupSnapshot Return One Snapshot of One Serverless Instance

Returns one snapshot of one serverless instance from the specified project.

This endpoint can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) API or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Get Flex Backup endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the serverless instance.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return GetServerlessBackupSnapshotApiRequest

Deprecated
*/
func (a *CloudBackupsApiService) GetServerlessBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string) GetServerlessBackupSnapshotApiRequest {
	return GetServerlessBackupSnapshotApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		snapshotId:  snapshotId,
	}
}

// GetServerlessBackupSnapshotExecute executes the request
//
//	@return ServerlessBackupSnapshot
//
// Deprecated
func (a *CloudBackupsApiService) GetServerlessBackupSnapshotExecute(r GetServerlessBackupSnapshotApiRequest) (*ServerlessBackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessBackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetServerlessBackupSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/backup/snapshots/{snapshotId}"
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

type GetServerlessRestoreJobApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	restoreJobId string
}

type GetServerlessRestoreJobApiParams struct {
	GroupId      string
	ClusterName  string
	RestoreJobId string
}

func (a *CloudBackupsApiService) GetServerlessRestoreJobWithParams(ctx context.Context, args *GetServerlessRestoreJobApiParams) GetServerlessRestoreJobApiRequest {
	return GetServerlessRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		clusterName:  args.ClusterName,
		restoreJobId: args.RestoreJobId,
	}
}

func (r GetServerlessRestoreJobApiRequest) Execute() (*ServerlessBackupRestoreJob, *http.Response, error) {
	return r.ApiService.GetServerlessRestoreJobExecute(r)
}

/*
GetServerlessRestoreJob Return One Restore Job for One Serverless Instance

Returns one restore job for one serverless instance from the specified project.

This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the Get Flex Backup Restore Job endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the serverless instance.
	@param restoreJobId Unique 24-hexadecimal digit string that identifies the restore job to return.
	@return GetServerlessRestoreJobApiRequest

Deprecated
*/
func (a *CloudBackupsApiService) GetServerlessRestoreJob(ctx context.Context, groupId string, clusterName string, restoreJobId string) GetServerlessRestoreJobApiRequest {
	return GetServerlessRestoreJobApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		clusterName:  clusterName,
		restoreJobId: restoreJobId,
	}
}

// GetServerlessRestoreJobExecute executes the request
//
//	@return ServerlessBackupRestoreJob
//
// Deprecated
func (a *CloudBackupsApiService) GetServerlessRestoreJobExecute(r GetServerlessRestoreJobApiRequest) (*ServerlessBackupRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessBackupRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.GetServerlessRestoreJob")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/backup/restoreJobs/{restoreJobId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.restoreJobId == "" {
		return localVarReturnValue, nil, reportError("restoreJobId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"restoreJobId"+"}", url.PathEscape(r.restoreJobId), -1)

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

type ListBackupExportsApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListBackupExportsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *CloudBackupsApiService) ListBackupExportsWithParams(ctx context.Context, args *ListBackupExportsApiParams) ListBackupExportsApiRequest {
	return ListBackupExportsApiRequest{
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
func (r ListBackupExportsApiRequest) IncludeCount(includeCount bool) ListBackupExportsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListBackupExportsApiRequest) ItemsPerPage(itemsPerPage int) ListBackupExportsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListBackupExportsApiRequest) PageNum(pageNum int) ListBackupExportsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListBackupExportsApiRequest) Execute() (*PaginatedApiAtlasDiskBackupExportJob, *http.Response, error) {
	return r.ApiService.ListBackupExportsExecute(r)
}

/*
ListBackupExports Return All Snapshot Export Jobs

Returns all Cloud Backup Snapshot Export Jobs associated with the specified Atlas cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListBackupExportsApiRequest
*/
func (a *CloudBackupsApiService) ListBackupExports(ctx context.Context, groupId string, clusterName string) ListBackupExportsApiRequest {
	return ListBackupExportsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListBackupExportsExecute executes the request
//
//	@return PaginatedApiAtlasDiskBackupExportJob
func (a *CloudBackupsApiService) ListBackupExportsExecute(r ListBackupExportsApiRequest) (*PaginatedApiAtlasDiskBackupExportJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasDiskBackupExportJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListBackupExports")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/exports"
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

type ListBackupPrivateEndpointsApiRequest struct {
	ctx           context.Context
	ApiService    CloudBackupsApi
	groupId       string
	cloudProvider string
	includeCount  *bool
	itemsPerPage  *int
	pageNum       *int
}

type ListBackupPrivateEndpointsApiParams struct {
	GroupId       string
	CloudProvider string
	IncludeCount  *bool
	ItemsPerPage  *int
	PageNum       *int
}

func (a *CloudBackupsApiService) ListBackupPrivateEndpointsWithParams(ctx context.Context, args *ListBackupPrivateEndpointsApiParams) ListBackupPrivateEndpointsApiRequest {
	return ListBackupPrivateEndpointsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		includeCount:  args.IncludeCount,
		itemsPerPage:  args.ItemsPerPage,
		pageNum:       args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListBackupPrivateEndpointsApiRequest) IncludeCount(includeCount bool) ListBackupPrivateEndpointsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListBackupPrivateEndpointsApiRequest) ItemsPerPage(itemsPerPage int) ListBackupPrivateEndpointsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListBackupPrivateEndpointsApiRequest) PageNum(pageNum int) ListBackupPrivateEndpointsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListBackupPrivateEndpointsApiRequest) Execute() (*PaginatedApiAtlasObjectStoragePrivateEndpointResponse, *http.Response, error) {
	return r.ApiService.ListBackupPrivateEndpointsExecute(r)
}

/*
ListBackupPrivateEndpoints Return Object Storage Private Endpoints for Cloud Backups for One Cloud Provider in One Project

Returns the private endpoints of the specified cloud provider for object storage backup operations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param cloudProvider Human-readable label that identifies the cloud provider for the private endpoints to return.
	@return ListBackupPrivateEndpointsApiRequest
*/
func (a *CloudBackupsApiService) ListBackupPrivateEndpoints(ctx context.Context, groupId string, cloudProvider string) ListBackupPrivateEndpointsApiRequest {
	return ListBackupPrivateEndpointsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		cloudProvider: cloudProvider,
	}
}

// ListBackupPrivateEndpointsExecute executes the request
//
//	@return PaginatedApiAtlasObjectStoragePrivateEndpointResponse
func (a *CloudBackupsApiService) ListBackupPrivateEndpointsExecute(r ListBackupPrivateEndpointsApiRequest) (*PaginatedApiAtlasObjectStoragePrivateEndpointResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasObjectStoragePrivateEndpointResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListBackupPrivateEndpoints")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/{cloudProvider}/privateEndpoints"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.cloudProvider == "" {
		return localVarReturnValue, nil, reportError("cloudProvider is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"cloudProvider"+"}", url.PathEscape(r.cloudProvider), -1)

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

type ListBackupRestoreJobsApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListBackupRestoreJobsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *CloudBackupsApiService) ListBackupRestoreJobsWithParams(ctx context.Context, args *ListBackupRestoreJobsApiParams) ListBackupRestoreJobsApiRequest {
	return ListBackupRestoreJobsApiRequest{
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
func (r ListBackupRestoreJobsApiRequest) IncludeCount(includeCount bool) ListBackupRestoreJobsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListBackupRestoreJobsApiRequest) ItemsPerPage(itemsPerPage int) ListBackupRestoreJobsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListBackupRestoreJobsApiRequest) PageNum(pageNum int) ListBackupRestoreJobsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListBackupRestoreJobsApiRequest) Execute() (*PaginatedCloudBackupRestoreJob, *http.Response, error) {
	return r.ApiService.ListBackupRestoreJobsExecute(r)
}

/*
ListBackupRestoreJobs Return All Restore Jobs for One Cluster

Returns all cloud backup restore jobs for one cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster with the restore jobs you want to return.
	@return ListBackupRestoreJobsApiRequest
*/
func (a *CloudBackupsApiService) ListBackupRestoreJobs(ctx context.Context, groupId string, clusterName string) ListBackupRestoreJobsApiRequest {
	return ListBackupRestoreJobsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListBackupRestoreJobsExecute executes the request
//
//	@return PaginatedCloudBackupRestoreJob
func (a *CloudBackupsApiService) ListBackupRestoreJobsExecute(r ListBackupRestoreJobsApiRequest) (*PaginatedCloudBackupRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedCloudBackupRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListBackupRestoreJobs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/restoreJobs"
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

type ListBackupShardedClustersApiRequest struct {
	ctx         context.Context
	ApiService  CloudBackupsApi
	groupId     string
	clusterName string
}

type ListBackupShardedClustersApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *CloudBackupsApiService) ListBackupShardedClustersWithParams(ctx context.Context, args *ListBackupShardedClustersApiParams) ListBackupShardedClustersApiRequest {
	return ListBackupShardedClustersApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r ListBackupShardedClustersApiRequest) Execute() (*PaginatedCloudBackupShardedClusterSnapshot, *http.Response, error) {
	return r.ApiService.ListBackupShardedClustersExecute(r)
}

/*
ListBackupShardedClusters Return All Sharded Cluster Cloud Backups

Returns all snapshots of one sharded cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListBackupShardedClustersApiRequest
*/
func (a *CloudBackupsApiService) ListBackupShardedClusters(ctx context.Context, groupId string, clusterName string) ListBackupShardedClustersApiRequest {
	return ListBackupShardedClustersApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListBackupShardedClustersExecute executes the request
//
//	@return PaginatedCloudBackupShardedClusterSnapshot
func (a *CloudBackupsApiService) ListBackupShardedClustersExecute(r ListBackupShardedClustersApiRequest) (*PaginatedCloudBackupShardedClusterSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedCloudBackupShardedClusterSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListBackupShardedClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/shardedClusters"
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

type ListBackupSnapshotsApiRequest struct {
	ctx                   context.Context
	ApiService            CloudBackupsApi
	groupId               string
	clusterName           string
	includeCount          *bool
	itemsPerPage          *int
	pageNum               *int
	pointInTimeUtcSeconds *int64
	oplogTs               *int64
	oplogInc              *int64
}

type ListBackupSnapshotsApiParams struct {
	GroupId               string
	ClusterName           string
	IncludeCount          *bool
	ItemsPerPage          *int
	PageNum               *int
	PointInTimeUtcSeconds *int64
	OplogTs               *int64
	OplogInc              *int64
}

func (a *CloudBackupsApiService) ListBackupSnapshotsWithParams(ctx context.Context, args *ListBackupSnapshotsApiParams) ListBackupSnapshotsApiRequest {
	return ListBackupSnapshotsApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		groupId:               args.GroupId,
		clusterName:           args.ClusterName,
		includeCount:          args.IncludeCount,
		itemsPerPage:          args.ItemsPerPage,
		pageNum:               args.PageNum,
		pointInTimeUtcSeconds: args.PointInTimeUtcSeconds,
		oplogTs:               args.OplogTs,
		oplogInc:              args.OplogInc,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListBackupSnapshotsApiRequest) IncludeCount(includeCount bool) ListBackupSnapshotsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListBackupSnapshotsApiRequest) ItemsPerPage(itemsPerPage int) ListBackupSnapshotsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListBackupSnapshotsApiRequest) PageNum(pageNum int) ListBackupSnapshotsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Desired point in time, expressed as the number of seconds that have elapsed since the UNIX epoch. If specified, returns the closest snapshot created before that point in time. Mutually exclusive with &#x60;oplogTs&#x60; and &#x60;oplogInc&#x60;.
func (r ListBackupSnapshotsApiRequest) PointInTimeUtcSeconds(pointInTimeUtcSeconds int64) ListBackupSnapshotsApiRequest {
	r.pointInTimeUtcSeconds = &pointInTimeUtcSeconds
	return r
}

// Oplog timestamp that represents the desired point in time. This is the first part of an Oplog timestamp. Must be used with &#x60;oplogInc&#x60;. Mutually exclusive with &#x60;pointInTimeUtcSeconds&#x60;.
func (r ListBackupSnapshotsApiRequest) OplogTs(oplogTs int64) ListBackupSnapshotsApiRequest {
	r.oplogTs = &oplogTs
	return r
}

// Oplog operation number that represents the desired point in time. This is the second part of an Oplog timestamp. Must be used with &#x60;oplogTs&#x60;. Mutually exclusive with &#x60;pointInTimeUtcSeconds&#x60;.
func (r ListBackupSnapshotsApiRequest) OplogInc(oplogInc int64) ListBackupSnapshotsApiRequest {
	r.oplogInc = &oplogInc
	return r
}

func (r ListBackupSnapshotsApiRequest) Execute() (*PaginatedCloudBackupReplicaSet, *http.Response, error) {
	return r.ApiService.ListBackupSnapshotsExecute(r)
}

/*
ListBackupSnapshots Return All Replica Set Cloud Backups

Returns all snapshots of one cluster from the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return ListBackupSnapshotsApiRequest
*/
func (a *CloudBackupsApiService) ListBackupSnapshots(ctx context.Context, groupId string, clusterName string) ListBackupSnapshotsApiRequest {
	return ListBackupSnapshotsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListBackupSnapshotsExecute executes the request
//
//	@return PaginatedCloudBackupReplicaSet
func (a *CloudBackupsApiService) ListBackupSnapshotsExecute(r ListBackupSnapshotsApiRequest) (*PaginatedCloudBackupReplicaSet, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedCloudBackupReplicaSet
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListBackupSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots"
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
	if r.pointInTimeUtcSeconds != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pointInTimeUtcSeconds", r.pointInTimeUtcSeconds, "")
	}
	if r.oplogTs != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "oplogTs", r.oplogTs, "")
	}
	if r.oplogInc != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "oplogInc", r.oplogInc, "")
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

type ListExportBucketsApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListExportBucketsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *CloudBackupsApiService) ListExportBucketsWithParams(ctx context.Context, args *ListExportBucketsApiParams) ListExportBucketsApiRequest {
	return ListExportBucketsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListExportBucketsApiRequest) IncludeCount(includeCount bool) ListExportBucketsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListExportBucketsApiRequest) ItemsPerPage(itemsPerPage int) ListExportBucketsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListExportBucketsApiRequest) PageNum(pageNum int) ListExportBucketsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListExportBucketsApiRequest) Execute() (*PaginatedBackupSnapshotExportBuckets, *http.Response, error) {
	return r.ApiService.ListExportBucketsExecute(r)
}

/*
ListExportBuckets Return All Snapshot Export Buckets

Returns all Export Buckets associated with the specified Project. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListExportBucketsApiRequest
*/
func (a *CloudBackupsApiService) ListExportBuckets(ctx context.Context, groupId string) ListExportBucketsApiRequest {
	return ListExportBucketsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListExportBucketsExecute executes the request
//
//	@return PaginatedBackupSnapshotExportBuckets
func (a *CloudBackupsApiService) ListExportBucketsExecute(r ListExportBucketsApiRequest) (*PaginatedBackupSnapshotExportBuckets, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedBackupSnapshotExportBuckets
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListExportBuckets")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/exportBuckets"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-05-30+json"}

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

type ListServerlessBackupSnapshotsApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListServerlessBackupSnapshotsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *CloudBackupsApiService) ListServerlessBackupSnapshotsWithParams(ctx context.Context, args *ListServerlessBackupSnapshotsApiParams) ListServerlessBackupSnapshotsApiRequest {
	return ListServerlessBackupSnapshotsApiRequest{
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
func (r ListServerlessBackupSnapshotsApiRequest) IncludeCount(includeCount bool) ListServerlessBackupSnapshotsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListServerlessBackupSnapshotsApiRequest) ItemsPerPage(itemsPerPage int) ListServerlessBackupSnapshotsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListServerlessBackupSnapshotsApiRequest) PageNum(pageNum int) ListServerlessBackupSnapshotsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListServerlessBackupSnapshotsApiRequest) Execute() (*PaginatedApiAtlasServerlessBackupSnapshot, *http.Response, error) {
	return r.ApiService.ListServerlessBackupSnapshotsExecute(r)
}

/*
ListServerlessBackupSnapshots Return All Snapshots of One Serverless Instance

Returns all snapshots of one serverless instance from the specified project.

This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the List Flex Backups endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the serverless instance.
	@return ListServerlessBackupSnapshotsApiRequest

Deprecated
*/
func (a *CloudBackupsApiService) ListServerlessBackupSnapshots(ctx context.Context, groupId string, clusterName string) ListServerlessBackupSnapshotsApiRequest {
	return ListServerlessBackupSnapshotsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListServerlessBackupSnapshotsExecute executes the request
//
//	@return PaginatedApiAtlasServerlessBackupSnapshot
//
// Deprecated
func (a *CloudBackupsApiService) ListServerlessBackupSnapshotsExecute(r ListServerlessBackupSnapshotsApiRequest) (*PaginatedApiAtlasServerlessBackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasServerlessBackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListServerlessBackupSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/backup/snapshots"
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

type ListServerlessRestoreJobsApiRequest struct {
	ctx          context.Context
	ApiService   CloudBackupsApi
	groupId      string
	clusterName  string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListServerlessRestoreJobsApiParams struct {
	GroupId      string
	ClusterName  string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *CloudBackupsApiService) ListServerlessRestoreJobsWithParams(ctx context.Context, args *ListServerlessRestoreJobsApiParams) ListServerlessRestoreJobsApiRequest {
	return ListServerlessRestoreJobsApiRequest{
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
func (r ListServerlessRestoreJobsApiRequest) IncludeCount(includeCount bool) ListServerlessRestoreJobsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListServerlessRestoreJobsApiRequest) ItemsPerPage(itemsPerPage int) ListServerlessRestoreJobsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListServerlessRestoreJobsApiRequest) PageNum(pageNum int) ListServerlessRestoreJobsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListServerlessRestoreJobsApiRequest) Execute() (*PaginatedApiAtlasServerlessBackupRestoreJob, *http.Response, error) {
	return r.ApiService.ListServerlessRestoreJobsExecute(r)
}

/*
ListServerlessRestoreJobs Return All Restore Jobs for One Serverless Instance

Returns all restore jobs for one serverless instance from the specified project.

This API can also be used on Flex clusters that were created with the [Create Serverless Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Serverless-Instances/operation/createServerlessInstance) endpoint or Flex clusters that were migrated from Serverless instances. This endpoint will be sunset on January 22, 2026. Please use the List Flex Backup Restore Jobs endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the serverless instance.
	@return ListServerlessRestoreJobsApiRequest

Deprecated
*/
func (a *CloudBackupsApiService) ListServerlessRestoreJobs(ctx context.Context, groupId string, clusterName string) ListServerlessRestoreJobsApiRequest {
	return ListServerlessRestoreJobsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListServerlessRestoreJobsExecute executes the request
//
//	@return PaginatedApiAtlasServerlessBackupRestoreJob
//
// Deprecated
func (a *CloudBackupsApiService) ListServerlessRestoreJobsExecute(r ListServerlessRestoreJobsApiRequest) (*PaginatedApiAtlasServerlessBackupRestoreJob, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasServerlessBackupRestoreJob
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.ListServerlessRestoreJobs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/serverless/{clusterName}/backup/restoreJobs"
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

type TakeSnapshotsApiRequest struct {
	ctx                               context.Context
	ApiService                        CloudBackupsApi
	groupId                           string
	clusterName                       string
	diskBackupOnDemandSnapshotRequest *DiskBackupOnDemandSnapshotRequest
}

type TakeSnapshotsApiParams struct {
	GroupId                           string
	ClusterName                       string
	DiskBackupOnDemandSnapshotRequest *DiskBackupOnDemandSnapshotRequest
}

func (a *CloudBackupsApiService) TakeSnapshotsWithParams(ctx context.Context, args *TakeSnapshotsApiParams) TakeSnapshotsApiRequest {
	return TakeSnapshotsApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		groupId:                           args.GroupId,
		clusterName:                       args.ClusterName,
		diskBackupOnDemandSnapshotRequest: args.DiskBackupOnDemandSnapshotRequest,
	}
}

func (r TakeSnapshotsApiRequest) Execute() (*DiskBackupSnapshot, *http.Response, error) {
	return r.ApiService.TakeSnapshotsExecute(r)
}

/*
TakeSnapshots Take One On-Demand Snapshot

Takes one on-demand snapshot for the specified cluster. Atlas takes on-demand snapshots immediately and scheduled snapshots at regular intervals. If an on-demand snapshot with a status of `queued` or `inProgress` exists, before taking another snapshot, wait until Atlas completes processing the previously taken on-demand snapshot.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return TakeSnapshotsApiRequest
*/
func (a *CloudBackupsApiService) TakeSnapshots(ctx context.Context, groupId string, clusterName string, diskBackupOnDemandSnapshotRequest *DiskBackupOnDemandSnapshotRequest) TakeSnapshotsApiRequest {
	return TakeSnapshotsApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		groupId:                           groupId,
		clusterName:                       clusterName,
		diskBackupOnDemandSnapshotRequest: diskBackupOnDemandSnapshotRequest,
	}
}

// TakeSnapshotsExecute executes the request
//
//	@return DiskBackupSnapshot
func (a *CloudBackupsApiService) TakeSnapshotsExecute(r TakeSnapshotsApiRequest) (*DiskBackupSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.TakeSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots"
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
	if r.diskBackupOnDemandSnapshotRequest == nil {
		return localVarReturnValue, nil, reportError("diskBackupOnDemandSnapshotRequest is required and must be specified")
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
	localVarPostBody = r.diskBackupOnDemandSnapshotRequest
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

type UpdateBackupExportBucketApiRequest struct {
	ctx                                   context.Context
	ApiService                            CloudBackupsApi
	groupId                               string
	exportBucketId                        string
	updateRequirePrivateNetworkingRequest *UpdateRequirePrivateNetworkingRequest
}

type UpdateBackupExportBucketApiParams struct {
	GroupId                               string
	ExportBucketId                        string
	UpdateRequirePrivateNetworkingRequest *UpdateRequirePrivateNetworkingRequest
}

func (a *CloudBackupsApiService) UpdateBackupExportBucketWithParams(ctx context.Context, args *UpdateBackupExportBucketApiParams) UpdateBackupExportBucketApiRequest {
	return UpdateBackupExportBucketApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               args.GroupId,
		exportBucketId:                        args.ExportBucketId,
		updateRequirePrivateNetworkingRequest: args.UpdateRequirePrivateNetworkingRequest,
	}
}

func (r UpdateBackupExportBucketApiRequest) Execute() (*DiskBackupSnapshotAWSExportBucketResponse, *http.Response, error) {
	return r.ApiService.UpdateBackupExportBucketExecute(r)
}

/*
UpdateBackupExportBucket Update One Export Bucket Private Networking Settings

Updates the private networking settings for one snapshot export bucket in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param exportBucketId Unique 24-hexadecimal character string that identifies the snapshot export bucket.
	@return UpdateBackupExportBucketApiRequest
*/
func (a *CloudBackupsApiService) UpdateBackupExportBucket(ctx context.Context, groupId string, exportBucketId string, updateRequirePrivateNetworkingRequest *UpdateRequirePrivateNetworkingRequest) UpdateBackupExportBucketApiRequest {
	return UpdateBackupExportBucketApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               groupId,
		exportBucketId:                        exportBucketId,
		updateRequirePrivateNetworkingRequest: updateRequirePrivateNetworkingRequest,
	}
}

// UpdateBackupExportBucketExecute executes the request
//
//	@return DiskBackupSnapshotAWSExportBucketResponse
func (a *CloudBackupsApiService) UpdateBackupExportBucketExecute(r UpdateBackupExportBucketApiRequest) (*DiskBackupSnapshotAWSExportBucketResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotAWSExportBucketResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.UpdateBackupExportBucket")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backup/exportBuckets/{exportBucketId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.exportBucketId == "" {
		return localVarReturnValue, nil, reportError("exportBucketId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"exportBucketId"+"}", url.PathEscape(r.exportBucketId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateRequirePrivateNetworkingRequest == nil {
		return localVarReturnValue, nil, reportError("updateRequirePrivateNetworkingRequest is required and must be specified")
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
	localVarPostBody = r.updateRequirePrivateNetworkingRequest
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

type UpdateBackupScheduleApiRequest struct {
	ctx                                context.Context
	ApiService                         CloudBackupsApi
	groupId                            string
	clusterName                        string
	diskBackupSnapshotSchedule20240805 *DiskBackupSnapshotSchedule20240805
}

type UpdateBackupScheduleApiParams struct {
	GroupId                            string
	ClusterName                        string
	DiskBackupSnapshotSchedule20240805 *DiskBackupSnapshotSchedule20240805
}

func (a *CloudBackupsApiService) UpdateBackupScheduleWithParams(ctx context.Context, args *UpdateBackupScheduleApiParams) UpdateBackupScheduleApiRequest {
	return UpdateBackupScheduleApiRequest{
		ApiService:                         a,
		ctx:                                ctx,
		groupId:                            args.GroupId,
		clusterName:                        args.ClusterName,
		diskBackupSnapshotSchedule20240805: args.DiskBackupSnapshotSchedule20240805,
	}
}

func (r UpdateBackupScheduleApiRequest) Execute() (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	return r.ApiService.UpdateBackupScheduleExecute(r)
}

/*
UpdateBackupSchedule Update Cloud Backup Schedule for One Cluster

Updates the cloud backup schedule for one cluster within the specified project. This schedule defines when MongoDB Cloud takes scheduled snapshots and how long it stores those snapshots. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return UpdateBackupScheduleApiRequest
*/
func (a *CloudBackupsApiService) UpdateBackupSchedule(ctx context.Context, groupId string, clusterName string, diskBackupSnapshotSchedule20240805 *DiskBackupSnapshotSchedule20240805) UpdateBackupScheduleApiRequest {
	return UpdateBackupScheduleApiRequest{
		ApiService:                         a,
		ctx:                                ctx,
		groupId:                            groupId,
		clusterName:                        clusterName,
		diskBackupSnapshotSchedule20240805: diskBackupSnapshotSchedule20240805,
	}
}

// UpdateBackupScheduleExecute executes the request
//
//	@return DiskBackupSnapshotSchedule20240805
func (a *CloudBackupsApiService) UpdateBackupScheduleExecute(r UpdateBackupScheduleApiRequest) (*DiskBackupSnapshotSchedule20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupSnapshotSchedule20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.UpdateBackupSchedule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/schedule"
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
	if r.diskBackupSnapshotSchedule20240805 == nil {
		return localVarReturnValue, nil, reportError("diskBackupSnapshotSchedule20240805 is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.diskBackupSnapshotSchedule20240805
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

type UpdateBackupSnapshotApiRequest struct {
	ctx                     context.Context
	ApiService              CloudBackupsApi
	groupId                 string
	clusterName             string
	snapshotId              string
	backupSnapshotRetention *BackupSnapshotRetention
}

type UpdateBackupSnapshotApiParams struct {
	GroupId                 string
	ClusterName             string
	SnapshotId              string
	BackupSnapshotRetention *BackupSnapshotRetention
}

func (a *CloudBackupsApiService) UpdateBackupSnapshotWithParams(ctx context.Context, args *UpdateBackupSnapshotApiParams) UpdateBackupSnapshotApiRequest {
	return UpdateBackupSnapshotApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 args.GroupId,
		clusterName:             args.ClusterName,
		snapshotId:              args.SnapshotId,
		backupSnapshotRetention: args.BackupSnapshotRetention,
	}
}

func (r UpdateBackupSnapshotApiRequest) Execute() (*DiskBackupReplicaSet, *http.Response, error) {
	return r.ApiService.UpdateBackupSnapshotExecute(r)
}

/*
UpdateBackupSnapshot Update Expiration Date for One Cloud Backup

Changes the expiration date for one cloud backup snapshot for one cluster in the specified project, the requesting Service Account or API Key must have the Project Backup Manager role.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@param snapshotId Unique 24-hexadecimal digit string that identifies the desired snapshot.
	@return UpdateBackupSnapshotApiRequest
*/
func (a *CloudBackupsApiService) UpdateBackupSnapshot(ctx context.Context, groupId string, clusterName string, snapshotId string, backupSnapshotRetention *BackupSnapshotRetention) UpdateBackupSnapshotApiRequest {
	return UpdateBackupSnapshotApiRequest{
		ApiService:              a,
		ctx:                     ctx,
		groupId:                 groupId,
		clusterName:             clusterName,
		snapshotId:              snapshotId,
		backupSnapshotRetention: backupSnapshotRetention,
	}
}

// UpdateBackupSnapshotExecute executes the request
//
//	@return DiskBackupReplicaSet
func (a *CloudBackupsApiService) UpdateBackupSnapshotExecute(r UpdateBackupSnapshotApiRequest) (*DiskBackupReplicaSet, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DiskBackupReplicaSet
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.UpdateBackupSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/backup/snapshots/{snapshotId}"
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
	if r.backupSnapshotRetention == nil {
		return localVarReturnValue, nil, reportError("backupSnapshotRetention is required and must be specified")
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
	localVarPostBody = r.backupSnapshotRetention
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

type UpdateCompliancePolicyApiRequest struct {
	ctx                            context.Context
	ApiService                     CloudBackupsApi
	groupId                        string
	dataProtectionSettings20231001 *DataProtectionSettings20231001
	overwriteBackupPolicies        *bool
}

type UpdateCompliancePolicyApiParams struct {
	GroupId                        string
	DataProtectionSettings20231001 *DataProtectionSettings20231001
	OverwriteBackupPolicies        *bool
}

func (a *CloudBackupsApiService) UpdateCompliancePolicyWithParams(ctx context.Context, args *UpdateCompliancePolicyApiParams) UpdateCompliancePolicyApiRequest {
	return UpdateCompliancePolicyApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        args.GroupId,
		dataProtectionSettings20231001: args.DataProtectionSettings20231001,
		overwriteBackupPolicies:        args.OverwriteBackupPolicies,
	}
}

// Flag that indicates whether to overwrite non complying backup policies with the new data protection settings or not.
func (r UpdateCompliancePolicyApiRequest) OverwriteBackupPolicies(overwriteBackupPolicies bool) UpdateCompliancePolicyApiRequest {
	r.overwriteBackupPolicies = &overwriteBackupPolicies
	return r
}

func (r UpdateCompliancePolicyApiRequest) Execute() (*DataProtectionSettings20231001, *http.Response, error) {
	return r.ApiService.UpdateCompliancePolicyExecute(r)
}

/*
UpdateCompliancePolicy Update Backup Compliance Policy Settings

Updates the Backup Compliance Policy settings for the specified project. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpdateCompliancePolicyApiRequest
*/
func (a *CloudBackupsApiService) UpdateCompliancePolicy(ctx context.Context, groupId string, dataProtectionSettings20231001 *DataProtectionSettings20231001) UpdateCompliancePolicyApiRequest {
	return UpdateCompliancePolicyApiRequest{
		ApiService:                     a,
		ctx:                            ctx,
		groupId:                        groupId,
		dataProtectionSettings20231001: dataProtectionSettings20231001,
	}
}

// UpdateCompliancePolicyExecute executes the request
//
//	@return DataProtectionSettings20231001
func (a *CloudBackupsApiService) UpdateCompliancePolicyExecute(r UpdateCompliancePolicyApiRequest) (*DataProtectionSettings20231001, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *DataProtectionSettings20231001
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CloudBackupsApiService.UpdateCompliancePolicy")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/backupCompliancePolicy"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.dataProtectionSettings20231001 == nil {
		return localVarReturnValue, nil, reportError("dataProtectionSettings20231001 is required and must be specified")
	}

	if r.overwriteBackupPolicies != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "overwriteBackupPolicies", r.overwriteBackupPolicies, "")
	} else {
		var defaultValue bool = true
		r.overwriteBackupPolicies = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "overwriteBackupPolicies", r.overwriteBackupPolicies, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-10-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-10-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.dataProtectionSettings20231001
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
