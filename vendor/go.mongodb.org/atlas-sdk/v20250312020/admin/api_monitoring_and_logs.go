// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type MonitoringAndLogsApi interface {

	/*
		DownloadClusterLog Download Logs for One Cluster Host in One Project

		Returns a compressed (.gz) log file that contains a range of log messages for the specified host for the specified project. MongoDB updates process and audit logs from the cluster backend infrastructure every five minutes. Logs are stored in chunks approximately five minutes in length, but this duration may vary. If you poll the API for log files, we recommend polling every five minutes even though consecutive polls could contain some overlapping logs. This feature isn't available for `M0` free clusters, `M2`, `M5`, flex, or serverless clusters. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param hostName Human-readable label that identifies the host that stores the log files that you want to download.
		@param logName Human-readable label that identifies the log file that you want to return. To return audit logs, enable *Database Auditing* for the specified project.
		@return DownloadClusterLogApiRequest
	*/
	DownloadClusterLog(ctx context.Context, groupId string, hostName string, logName string) DownloadClusterLogApiRequest
	/*
		DownloadClusterLog Download Logs for One Cluster Host in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadClusterLogApiParams - Parameters for the request
		@return DownloadClusterLogApiRequest
	*/
	DownloadClusterLogWithParams(ctx context.Context, args *DownloadClusterLogApiParams) DownloadClusterLogApiRequest

	// Method available only for mocking purposes
	DownloadClusterLogExecute(r DownloadClusterLogApiRequest) (io.ReadCloser, *http.Response, error)

	/*
		GetDatabase Return One Database for One MongoDB Process

		Returns one database running on the specified host for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param databaseName Human-readable label that identifies the database that the specified MongoDB process serves.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return GetDatabaseApiRequest
	*/
	GetDatabase(ctx context.Context, groupId string, databaseName string, processId string) GetDatabaseApiRequest
	/*
		GetDatabase Return One Database for One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetDatabaseApiParams - Parameters for the request
		@return GetDatabaseApiRequest
	*/
	GetDatabaseWithParams(ctx context.Context, args *GetDatabaseApiParams) GetDatabaseApiRequest

	// Method available only for mocking purposes
	GetDatabaseExecute(r GetDatabaseApiRequest) (*MesurementsDatabase, *http.Response, error)

	/*
		GetDatabaseMeasurements Return Measurements for One Database in One MongoDB Process

		Returns the measurements of one database for the specified host for the specified project. Returns the database's on-disk storage space based on the MongoDB `dbStats` command output. To calculate some metric series, Atlas takes the rate between every two adjacent points. For these metric series, the first data point has a null value because Atlas can't calculate a rate for the first data point given the query time range. Atlas retrieves database metrics every 20 minutes but reduces frequency when necessary to optimize database performance.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param databaseName Human-readable label that identifies the database that the specified MongoDB process serves.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return GetDatabaseMeasurementsApiRequest
	*/
	GetDatabaseMeasurements(ctx context.Context, groupId string, databaseName string, processId string) GetDatabaseMeasurementsApiRequest
	/*
		GetDatabaseMeasurements Return Measurements for One Database in One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetDatabaseMeasurementsApiParams - Parameters for the request
		@return GetDatabaseMeasurementsApiRequest
	*/
	GetDatabaseMeasurementsWithParams(ctx context.Context, args *GetDatabaseMeasurementsApiParams) GetDatabaseMeasurementsApiRequest

	// Method available only for mocking purposes
	GetDatabaseMeasurementsExecute(r GetDatabaseMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error)

	/*
		GetGroupProcess Return One MongoDB Process by ID

		Returns the processes for the specified host for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return GetGroupProcessApiRequest
	*/
	GetGroupProcess(ctx context.Context, groupId string, processId string) GetGroupProcessApiRequest
	/*
		GetGroupProcess Return One MongoDB Process by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupProcessApiParams - Parameters for the request
		@return GetGroupProcessApiRequest
	*/
	GetGroupProcessWithParams(ctx context.Context, args *GetGroupProcessApiParams) GetGroupProcessApiRequest

	// Method available only for mocking purposes
	GetGroupProcessExecute(r GetGroupProcessApiRequest) (*ApiHostViewAtlas, *http.Response, error)

	/*
		GetIndexMeasurements Return Atlas Search Metrics for One Index in One Namespace

		Returns the Atlas Search metrics data series within the provided time range for one namespace and index name on the specified process. You must have the Project Read Only or higher role to view the Atlas Search metric types.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@param indexName Human-readable label that identifies the index.
		@param databaseName Human-readable label that identifies the database.
		@param collectionName Human-readable label that identifies the collection.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetIndexMeasurementsApiRequest
	*/
	GetIndexMeasurements(ctx context.Context, processId string, indexName string, databaseName string, collectionName string, groupId string) GetIndexMeasurementsApiRequest
	/*
		GetIndexMeasurements Return Atlas Search Metrics for One Index in One Namespace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetIndexMeasurementsApiParams - Parameters for the request
		@return GetIndexMeasurementsApiRequest
	*/
	GetIndexMeasurementsWithParams(ctx context.Context, args *GetIndexMeasurementsApiParams) GetIndexMeasurementsApiRequest

	// Method available only for mocking purposes
	GetIndexMeasurementsExecute(r GetIndexMeasurementsApiRequest) (*MeasurementsIndexes, *http.Response, error)

	/*
		GetProcessDisk Return Measurements for One Disk

		Returns measurement details for one disk or partition for the specified host for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param partitionName Human-readable label of the disk or partition to which the measurements apply.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return GetProcessDiskApiRequest
	*/
	GetProcessDisk(ctx context.Context, partitionName string, groupId string, processId string) GetProcessDiskApiRequest
	/*
		GetProcessDisk Return Measurements for One Disk


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetProcessDiskApiParams - Parameters for the request
		@return GetProcessDiskApiRequest
	*/
	GetProcessDiskWithParams(ctx context.Context, args *GetProcessDiskApiParams) GetProcessDiskApiRequest

	// Method available only for mocking purposes
	GetProcessDiskExecute(r GetProcessDiskApiRequest) (*MeasurementDiskPartition, *http.Response, error)

	/*
			GetProcessDiskMeasurements Return Measurements of One Disk for One MongoDB Process

			Returns the measurements of one disk or partition for the specified host for the specified project. Returned value can be one of the following:
		- Throughput of I/O operations for the disk partition used for the MongoDB process
		- Percentage of time during which requests the partition issued and serviced
		- Latency per operation type of the disk partition used for the MongoDB process
		- Amount of free and used disk space on the disk partition used for the MongoDB process.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param partitionName Human-readable label of the disk or partition to which the measurements apply.
			@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
			@return GetProcessDiskMeasurementsApiRequest
	*/
	GetProcessDiskMeasurements(ctx context.Context, groupId string, partitionName string, processId string) GetProcessDiskMeasurementsApiRequest
	/*
		GetProcessDiskMeasurements Return Measurements of One Disk for One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetProcessDiskMeasurementsApiParams - Parameters for the request
		@return GetProcessDiskMeasurementsApiRequest
	*/
	GetProcessDiskMeasurementsWithParams(ctx context.Context, args *GetProcessDiskMeasurementsApiParams) GetProcessDiskMeasurementsApiRequest

	// Method available only for mocking purposes
	GetProcessDiskMeasurementsExecute(r GetProcessDiskMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error)

	/*
			GetProcessMeasurements Return Measurements for One MongoDB Process

			Returns disk, partition, or host measurements per process for the specified host for the specified project. Returned value can be one of the following:
		- Throughput of I/O operations for the disk partition used for the MongoDB process
		- Percentage of time during which requests the partition issued and serviced
		- Latency per operation type of the disk partition used for the MongoDB process
		- Amount of free and used disk space on the disk partition used for the MongoDB process
		- Measurements for the host, such as CPU usage or number of I/O operations.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
			@return GetProcessMeasurementsApiRequest
	*/
	GetProcessMeasurements(ctx context.Context, groupId string, processId string) GetProcessMeasurementsApiRequest
	/*
		GetProcessMeasurements Return Measurements for One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetProcessMeasurementsApiParams - Parameters for the request
		@return GetProcessMeasurementsApiRequest
	*/
	GetProcessMeasurementsWithParams(ctx context.Context, args *GetProcessMeasurementsApiParams) GetProcessMeasurementsApiRequest

	// Method available only for mocking purposes
	GetProcessMeasurementsExecute(r GetProcessMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error)

	/*
		ListDatabases Return Available Databases for One MongoDB Process

		Returns the list of databases running on the specified host for the specified project. `M0` free clusters, `M2`, `M5`, serverless, and Flex clusters have some operational limits. The MongoDB Cloud process must be a `mongod`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return ListDatabasesApiRequest
	*/
	ListDatabases(ctx context.Context, groupId string, processId string) ListDatabasesApiRequest
	/*
		ListDatabases Return Available Databases for One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListDatabasesApiParams - Parameters for the request
		@return ListDatabasesApiRequest
	*/
	ListDatabasesWithParams(ctx context.Context, args *ListDatabasesApiParams) ListDatabasesApiRequest

	// Method available only for mocking purposes
	ListDatabasesExecute(r ListDatabasesApiRequest) (*PaginatedDatabase, *http.Response, error)

	/*
		ListGroupProcesses Return All MongoDB Processes in One Project

		Returns details of all processes for the specified project. A MongoDB process can be either a `mongod` or `mongos`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupProcessesApiRequest
	*/
	ListGroupProcesses(ctx context.Context, groupId string) ListGroupProcessesApiRequest
	/*
		ListGroupProcesses Return All MongoDB Processes in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupProcessesApiParams - Parameters for the request
		@return ListGroupProcessesApiRequest
	*/
	ListGroupProcessesWithParams(ctx context.Context, args *ListGroupProcessesApiParams) ListGroupProcessesApiRequest

	// Method available only for mocking purposes
	ListGroupProcessesExecute(r ListGroupProcessesApiRequest) (*PaginatedHostViewAtlas, *http.Response, error)

	/*
		ListHostFtsMetrics Return All Atlas Search Metric Types for One Process

		Returns all Atlas Search metric types available for one process in the specified project. You must have the Project Read Only or higher role to view the Atlas Search metric types.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListHostFtsMetricsApiRequest
	*/
	ListHostFtsMetrics(ctx context.Context, processId string, groupId string) ListHostFtsMetricsApiRequest
	/*
		ListHostFtsMetrics Return All Atlas Search Metric Types for One Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListHostFtsMetricsApiParams - Parameters for the request
		@return ListHostFtsMetricsApiRequest
	*/
	ListHostFtsMetricsWithParams(ctx context.Context, args *ListHostFtsMetricsApiParams) ListHostFtsMetricsApiRequest

	// Method available only for mocking purposes
	ListHostFtsMetricsExecute(r ListHostFtsMetricsApiRequest) (*CloudSearchMetrics, *http.Response, error)

	/*
		ListIndexMeasurements Return All Atlas Search Index Metrics for One Namespace

		Returns the Atlas Search index metrics within the specified time range for one namespace in the specified process.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@param databaseName Human-readable label that identifies the database.
		@param collectionName Human-readable label that identifies the collection.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListIndexMeasurementsApiRequest
	*/
	ListIndexMeasurements(ctx context.Context, processId string, databaseName string, collectionName string, groupId string) ListIndexMeasurementsApiRequest
	/*
		ListIndexMeasurements Return All Atlas Search Index Metrics for One Namespace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListIndexMeasurementsApiParams - Parameters for the request
		@return ListIndexMeasurementsApiRequest
	*/
	ListIndexMeasurementsWithParams(ctx context.Context, args *ListIndexMeasurementsApiParams) ListIndexMeasurementsApiRequest

	// Method available only for mocking purposes
	ListIndexMeasurementsExecute(r ListIndexMeasurementsApiRequest) (*MeasurementsIndexes, *http.Response, error)

	/*
		ListMeasurements Return Atlas Search Hardware and Status Metrics

		Returns the Atlas Search hardware and status data series within the provided time range for one process in the specified project. You must have the Project Read Only or higher role to view the Atlas Search metric types.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListMeasurementsApiRequest
	*/
	ListMeasurements(ctx context.Context, processId string, groupId string) ListMeasurementsApiRequest
	/*
		ListMeasurements Return Atlas Search Hardware and Status Metrics


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListMeasurementsApiParams - Parameters for the request
		@return ListMeasurementsApiRequest
	*/
	ListMeasurementsWithParams(ctx context.Context, args *ListMeasurementsApiParams) ListMeasurementsApiRequest

	// Method available only for mocking purposes
	ListMeasurementsExecute(r ListMeasurementsApiRequest) (*MeasurementsNonIndex, *http.Response, error)

	/*
		ListProcessDisks Return Available Disks for One MongoDB Process

		Returns the list of disks or partitions for the specified host for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
		@return ListProcessDisksApiRequest
	*/
	ListProcessDisks(ctx context.Context, groupId string, processId string) ListProcessDisksApiRequest
	/*
		ListProcessDisks Return Available Disks for One MongoDB Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListProcessDisksApiParams - Parameters for the request
		@return ListProcessDisksApiRequest
	*/
	ListProcessDisksWithParams(ctx context.Context, args *ListProcessDisksApiParams) ListProcessDisksApiRequest

	// Method available only for mocking purposes
	ListProcessDisksExecute(r ListProcessDisksApiRequest) (*PaginatedDiskPartition, *http.Response, error)
}

// MonitoringAndLogsApiService MonitoringAndLogsApi service
type MonitoringAndLogsApiService service

type DownloadClusterLogApiRequest struct {
	ctx        context.Context
	ApiService MonitoringAndLogsApi
	groupId    string
	hostName   string
	logName    string
	endDate    *int64
	startDate  *int64
}

type DownloadClusterLogApiParams struct {
	GroupId   string
	HostName  string
	LogName   string
	EndDate   *int64
	StartDate *int64
}

func (a *MonitoringAndLogsApiService) DownloadClusterLogWithParams(ctx context.Context, args *DownloadClusterLogApiParams) DownloadClusterLogApiRequest {
	return DownloadClusterLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		hostName:   args.HostName,
		logName:    args.LogName,
		endDate:    args.EndDate,
		startDate:  args.StartDate,
	}
}

// Specifies the date and time for the ending point of the range of log messages to retrieve, in the number of seconds that have elapsed since the UNIX epoch. This value will default to 24 hours after the start date. If the start date is also unspecified, the value will default to the time of the request.
func (r DownloadClusterLogApiRequest) EndDate(endDate int64) DownloadClusterLogApiRequest {
	r.endDate = &endDate
	return r
}

// Specifies the date and time for the starting point of the range of log messages to retrieve, in the number of seconds that have elapsed since the UNIX epoch. This value will default to 24 hours prior to the end date. If the end date is also unspecified, the value will default to 24 hours prior to the time of the request.
func (r DownloadClusterLogApiRequest) StartDate(startDate int64) DownloadClusterLogApiRequest {
	r.startDate = &startDate
	return r
}

func (r DownloadClusterLogApiRequest) Execute() (io.ReadCloser, *http.Response, error) {
	return r.ApiService.DownloadClusterLogExecute(r)
}

/*
DownloadClusterLog Download Logs for One Cluster Host in One Project

Returns a compressed (.gz) log file that contains a range of log messages for the specified host for the specified project. MongoDB updates process and audit logs from the cluster backend infrastructure every five minutes. Logs are stored in chunks approximately five minutes in length, but this duration may vary. If you poll the API for log files, we recommend polling every five minutes even though consecutive polls could contain some overlapping logs. This feature isn't available for `M0` free clusters, `M2`, `M5`, flex, or serverless clusters. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param hostName Human-readable label that identifies the host that stores the log files that you want to download.
	@param logName Human-readable label that identifies the log file that you want to return. To return audit logs, enable *Database Auditing* for the specified project.
	@return DownloadClusterLogApiRequest
*/
func (a *MonitoringAndLogsApiService) DownloadClusterLog(ctx context.Context, groupId string, hostName string, logName string) DownloadClusterLogApiRequest {
	return DownloadClusterLogApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		hostName:   hostName,
		logName:    logName,
	}
}

// DownloadClusterLogExecute executes the request
//
//	@return io.ReadCloser
func (a *MonitoringAndLogsApiService) DownloadClusterLogExecute(r DownloadClusterLogApiRequest) (io.ReadCloser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue io.ReadCloser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.DownloadClusterLog")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{hostName}/logs/{logName}.gz"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.hostName == "" {
		return localVarReturnValue, nil, reportError("hostName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"hostName"+"}", url.PathEscape(r.hostName), -1)
	if r.logName == "" {
		return localVarReturnValue, nil, reportError("logName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"logName"+"}", url.PathEscape(r.logName), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+gzip"}

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

type GetDatabaseApiRequest struct {
	ctx          context.Context
	ApiService   MonitoringAndLogsApi
	groupId      string
	databaseName string
	processId    string
}

type GetDatabaseApiParams struct {
	GroupId      string
	DatabaseName string
	ProcessId    string
}

func (a *MonitoringAndLogsApiService) GetDatabaseWithParams(ctx context.Context, args *GetDatabaseApiParams) GetDatabaseApiRequest {
	return GetDatabaseApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		databaseName: args.DatabaseName,
		processId:    args.ProcessId,
	}
}

func (r GetDatabaseApiRequest) Execute() (*MesurementsDatabase, *http.Response, error) {
	return r.ApiService.GetDatabaseExecute(r)
}

/*
GetDatabase Return One Database for One MongoDB Process

Returns one database running on the specified host for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param databaseName Human-readable label that identifies the database that the specified MongoDB process serves.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetDatabaseApiRequest
*/
func (a *MonitoringAndLogsApiService) GetDatabase(ctx context.Context, groupId string, databaseName string, processId string) GetDatabaseApiRequest {
	return GetDatabaseApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		databaseName: databaseName,
		processId:    processId,
	}
}

// GetDatabaseExecute executes the request
//
//	@return MesurementsDatabase
func (a *MonitoringAndLogsApiService) GetDatabaseExecute(r GetDatabaseApiRequest) (*MesurementsDatabase, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MesurementsDatabase
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetDatabase")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/databases/{databaseName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

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

type GetDatabaseMeasurementsApiRequest struct {
	ctx          context.Context
	ApiService   MonitoringAndLogsApi
	groupId      string
	databaseName string
	processId    string
	granularity  *string
	m            *[]string
	period       *string
	start        *time.Time
	end          *time.Time
}

type GetDatabaseMeasurementsApiParams struct {
	GroupId      string
	DatabaseName string
	ProcessId    string
	Granularity  *string
	M            *[]string
	Period       *string
	Start        *time.Time
	End          *time.Time
}

func (a *MonitoringAndLogsApiService) GetDatabaseMeasurementsWithParams(ctx context.Context, args *GetDatabaseMeasurementsApiParams) GetDatabaseMeasurementsApiRequest {
	return GetDatabaseMeasurementsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		databaseName: args.DatabaseName,
		processId:    args.ProcessId,
		granularity:  args.Granularity,
		m:            args.M,
		period:       args.Period,
		start:        args.Start,
		end:          args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r GetDatabaseMeasurementsApiRequest) Granularity(granularity string) GetDatabaseMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// One or more types of measurement to request for this MongoDB process. If omitted, the resource returns all measurements. To specify multiple values for &#x60;m&#x60;, repeat the &#x60;m&#x60; parameter for each value. Specify measurements that apply to the specified host. MongoDB Cloud returns an error if you specified any invalid measurements.
func (r GetDatabaseMeasurementsApiRequest) M(m []string) GetDatabaseMeasurementsApiRequest {
	r.m = &m
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetDatabaseMeasurementsApiRequest) Period(period string) GetDatabaseMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetDatabaseMeasurementsApiRequest) Start(start time.Time) GetDatabaseMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetDatabaseMeasurementsApiRequest) End(end time.Time) GetDatabaseMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r GetDatabaseMeasurementsApiRequest) Execute() (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	return r.ApiService.GetDatabaseMeasurementsExecute(r)
}

/*
GetDatabaseMeasurements Return Measurements for One Database in One MongoDB Process

Returns the measurements of one database for the specified host for the specified project. Returns the database's on-disk storage space based on the MongoDB `dbStats` command output. To calculate some metric series, Atlas takes the rate between every two adjacent points. For these metric series, the first data point has a null value because Atlas can't calculate a rate for the first data point given the query time range. Atlas retrieves database metrics every 20 minutes but reduces frequency when necessary to optimize database performance.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param databaseName Human-readable label that identifies the database that the specified MongoDB process serves.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetDatabaseMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) GetDatabaseMeasurements(ctx context.Context, groupId string, databaseName string, processId string) GetDatabaseMeasurementsApiRequest {
	return GetDatabaseMeasurementsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		databaseName: databaseName,
		processId:    processId,
	}
}

// GetDatabaseMeasurementsExecute executes the request
//
//	@return ApiMeasurementsGeneralViewAtlas
func (a *MonitoringAndLogsApiService) GetDatabaseMeasurementsExecute(r GetDatabaseMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiMeasurementsGeneralViewAtlas
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetDatabaseMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/databases/{databaseName}/measurements"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}

	if r.m != nil {
		t := *r.m
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "m", t, "multi")

	}
	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
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

type GetGroupProcessApiRequest struct {
	ctx        context.Context
	ApiService MonitoringAndLogsApi
	groupId    string
	processId  string
}

type GetGroupProcessApiParams struct {
	GroupId   string
	ProcessId string
}

func (a *MonitoringAndLogsApiService) GetGroupProcessWithParams(ctx context.Context, args *GetGroupProcessApiParams) GetGroupProcessApiRequest {
	return GetGroupProcessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		processId:  args.ProcessId,
	}
}

func (r GetGroupProcessApiRequest) Execute() (*ApiHostViewAtlas, *http.Response, error) {
	return r.ApiService.GetGroupProcessExecute(r)
}

/*
GetGroupProcess Return One MongoDB Process by ID

Returns the processes for the specified host for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetGroupProcessApiRequest
*/
func (a *MonitoringAndLogsApiService) GetGroupProcess(ctx context.Context, groupId string, processId string) GetGroupProcessApiRequest {
	return GetGroupProcessApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// GetGroupProcessExecute executes the request
//
//	@return ApiHostViewAtlas
func (a *MonitoringAndLogsApiService) GetGroupProcessExecute(r GetGroupProcessApiRequest) (*ApiHostViewAtlas, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiHostViewAtlas
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetGroupProcess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

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

type GetIndexMeasurementsApiRequest struct {
	ctx            context.Context
	ApiService     MonitoringAndLogsApi
	processId      string
	indexName      string
	databaseName   string
	collectionName string
	groupId        string
	granularity    *string
	metrics        *[]string
	period         *string
	start          *time.Time
	end            *time.Time
}

type GetIndexMeasurementsApiParams struct {
	ProcessId      string
	IndexName      string
	DatabaseName   string
	CollectionName string
	GroupId        string
	Granularity    *string
	Metrics        *[]string
	Period         *string
	Start          *time.Time
	End            *time.Time
}

func (a *MonitoringAndLogsApiService) GetIndexMeasurementsWithParams(ctx context.Context, args *GetIndexMeasurementsApiParams) GetIndexMeasurementsApiRequest {
	return GetIndexMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		processId:      args.ProcessId,
		indexName:      args.IndexName,
		databaseName:   args.DatabaseName,
		collectionName: args.CollectionName,
		groupId:        args.GroupId,
		granularity:    args.Granularity,
		metrics:        args.Metrics,
		period:         args.Period,
		start:          args.Start,
		end:            args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r GetIndexMeasurementsApiRequest) Granularity(granularity string) GetIndexMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// List that contains the measurements that MongoDB Atlas reports for the associated data series.
func (r GetIndexMeasurementsApiRequest) Metrics(metrics []string) GetIndexMeasurementsApiRequest {
	r.metrics = &metrics
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetIndexMeasurementsApiRequest) Period(period string) GetIndexMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetIndexMeasurementsApiRequest) Start(start time.Time) GetIndexMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetIndexMeasurementsApiRequest) End(end time.Time) GetIndexMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r GetIndexMeasurementsApiRequest) Execute() (*MeasurementsIndexes, *http.Response, error) {
	return r.ApiService.GetIndexMeasurementsExecute(r)
}

/*
GetIndexMeasurements Return Atlas Search Metrics for One Index in One Namespace

Returns the Atlas Search metrics data series within the provided time range for one namespace and index name on the specified process. You must have the Project Read Only or higher role to view the Atlas Search metric types.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@param indexName Human-readable label that identifies the index.
	@param databaseName Human-readable label that identifies the database.
	@param collectionName Human-readable label that identifies the collection.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetIndexMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) GetIndexMeasurements(ctx context.Context, processId string, indexName string, databaseName string, collectionName string, groupId string) GetIndexMeasurementsApiRequest {
	return GetIndexMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		processId:      processId,
		indexName:      indexName,
		databaseName:   databaseName,
		collectionName: collectionName,
		groupId:        groupId,
	}
}

// GetIndexMeasurementsExecute executes the request
//
//	@return MeasurementsIndexes
func (a *MonitoringAndLogsApiService) GetIndexMeasurementsExecute(r GetIndexMeasurementsApiRequest) (*MeasurementsIndexes, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementsIndexes
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetIndexMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/hosts/{processId}/fts/metrics/indexes/{databaseName}/{collectionName}/{indexName}/measurements"
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)
	if r.indexName == "" {
		return localVarReturnValue, nil, reportError("indexName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexName"+"}", url.PathEscape(r.indexName), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}
	if r.metrics == nil {
		return localVarReturnValue, nil, reportError("metrics is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	{
		t := *r.metrics
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "metrics", t, "multi")
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

type GetProcessDiskApiRequest struct {
	ctx           context.Context
	ApiService    MonitoringAndLogsApi
	partitionName string
	groupId       string
	processId     string
}

type GetProcessDiskApiParams struct {
	PartitionName string
	GroupId       string
	ProcessId     string
}

func (a *MonitoringAndLogsApiService) GetProcessDiskWithParams(ctx context.Context, args *GetProcessDiskApiParams) GetProcessDiskApiRequest {
	return GetProcessDiskApiRequest{
		ApiService:    a,
		ctx:           ctx,
		partitionName: args.PartitionName,
		groupId:       args.GroupId,
		processId:     args.ProcessId,
	}
}

func (r GetProcessDiskApiRequest) Execute() (*MeasurementDiskPartition, *http.Response, error) {
	return r.ApiService.GetProcessDiskExecute(r)
}

/*
GetProcessDisk Return Measurements for One Disk

Returns measurement details for one disk or partition for the specified host for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param partitionName Human-readable label of the disk or partition to which the measurements apply.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetProcessDiskApiRequest
*/
func (a *MonitoringAndLogsApiService) GetProcessDisk(ctx context.Context, partitionName string, groupId string, processId string) GetProcessDiskApiRequest {
	return GetProcessDiskApiRequest{
		ApiService:    a,
		ctx:           ctx,
		partitionName: partitionName,
		groupId:       groupId,
		processId:     processId,
	}
}

// GetProcessDiskExecute executes the request
//
//	@return MeasurementDiskPartition
func (a *MonitoringAndLogsApiService) GetProcessDiskExecute(r GetProcessDiskApiRequest) (*MeasurementDiskPartition, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementDiskPartition
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetProcessDisk")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/disks/{partitionName}"
	if r.partitionName == "" {
		return localVarReturnValue, nil, reportError("partitionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"partitionName"+"}", url.PathEscape(r.partitionName), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

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

type GetProcessDiskMeasurementsApiRequest struct {
	ctx           context.Context
	ApiService    MonitoringAndLogsApi
	groupId       string
	partitionName string
	processId     string
	granularity   *string
	m             *[]string
	period        *string
	start         *time.Time
	end           *time.Time
}

type GetProcessDiskMeasurementsApiParams struct {
	GroupId       string
	PartitionName string
	ProcessId     string
	Granularity   *string
	M             *[]string
	Period        *string
	Start         *time.Time
	End           *time.Time
}

func (a *MonitoringAndLogsApiService) GetProcessDiskMeasurementsWithParams(ctx context.Context, args *GetProcessDiskMeasurementsApiParams) GetProcessDiskMeasurementsApiRequest {
	return GetProcessDiskMeasurementsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		partitionName: args.PartitionName,
		processId:     args.ProcessId,
		granularity:   args.Granularity,
		m:             args.M,
		period:        args.Period,
		start:         args.Start,
		end:           args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r GetProcessDiskMeasurementsApiRequest) Granularity(granularity string) GetProcessDiskMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// One or more types of measurement to request for this MongoDB process. If omitted, the resource returns all measurements. To specify multiple values for &#x60;m&#x60;, repeat the &#x60;m&#x60; parameter for each value. Specify measurements that apply to the specified host. MongoDB Cloud returns an error if you specified any invalid measurements.
func (r GetProcessDiskMeasurementsApiRequest) M(m []string) GetProcessDiskMeasurementsApiRequest {
	r.m = &m
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetProcessDiskMeasurementsApiRequest) Period(period string) GetProcessDiskMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessDiskMeasurementsApiRequest) Start(start time.Time) GetProcessDiskMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessDiskMeasurementsApiRequest) End(end time.Time) GetProcessDiskMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r GetProcessDiskMeasurementsApiRequest) Execute() (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	return r.ApiService.GetProcessDiskMeasurementsExecute(r)
}

/*
GetProcessDiskMeasurements Return Measurements of One Disk for One MongoDB Process

Returns the measurements of one disk or partition for the specified host for the specified project. Returned value can be one of the following:
- Throughput of I/O operations for the disk partition used for the MongoDB process
- Percentage of time during which requests the partition issued and serviced
- Latency per operation type of the disk partition used for the MongoDB process
- Amount of free and used disk space on the disk partition used for the MongoDB process.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param partitionName Human-readable label of the disk or partition to which the measurements apply.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetProcessDiskMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) GetProcessDiskMeasurements(ctx context.Context, groupId string, partitionName string, processId string) GetProcessDiskMeasurementsApiRequest {
	return GetProcessDiskMeasurementsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		partitionName: partitionName,
		processId:     processId,
	}
}

// GetProcessDiskMeasurementsExecute executes the request
//
//	@return ApiMeasurementsGeneralViewAtlas
func (a *MonitoringAndLogsApiService) GetProcessDiskMeasurementsExecute(r GetProcessDiskMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiMeasurementsGeneralViewAtlas
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetProcessDiskMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/disks/{partitionName}/measurements"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.partitionName == "" {
		return localVarReturnValue, nil, reportError("partitionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"partitionName"+"}", url.PathEscape(r.partitionName), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}

	if r.m != nil {
		t := *r.m
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "m", t, "multi")

	}
	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
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

type GetProcessMeasurementsApiRequest struct {
	ctx         context.Context
	ApiService  MonitoringAndLogsApi
	groupId     string
	processId   string
	granularity *string
	m           *[]string
	period      *string
	start       *time.Time
	end         *time.Time
}

type GetProcessMeasurementsApiParams struct {
	GroupId     string
	ProcessId   string
	Granularity *string
	M           *[]string
	Period      *string
	Start       *time.Time
	End         *time.Time
}

func (a *MonitoringAndLogsApiService) GetProcessMeasurementsWithParams(ctx context.Context, args *GetProcessMeasurementsApiParams) GetProcessMeasurementsApiRequest {
	return GetProcessMeasurementsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		processId:   args.ProcessId,
		granularity: args.Granularity,
		m:           args.M,
		period:      args.Period,
		start:       args.Start,
		end:         args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r GetProcessMeasurementsApiRequest) Granularity(granularity string) GetProcessMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// One or more types of measurement to request for this MongoDB process. If omitted, the resource returns all measurements. To specify multiple values for &#x60;m&#x60;, repeat the &#x60;m&#x60; parameter for each value. Specify measurements that apply to the specified host. MongoDB Cloud returns an error if you specified any invalid measurements.
func (r GetProcessMeasurementsApiRequest) M(m []string) GetProcessMeasurementsApiRequest {
	r.m = &m
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r GetProcessMeasurementsApiRequest) Period(period string) GetProcessMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessMeasurementsApiRequest) Start(start time.Time) GetProcessMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r GetProcessMeasurementsApiRequest) End(end time.Time) GetProcessMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r GetProcessMeasurementsApiRequest) Execute() (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	return r.ApiService.GetProcessMeasurementsExecute(r)
}

/*
GetProcessMeasurements Return Measurements for One MongoDB Process

Returns disk, partition, or host measurements per process for the specified host for the specified project. Returned value can be one of the following:
- Throughput of I/O operations for the disk partition used for the MongoDB process
- Percentage of time during which requests the partition issued and serviced
- Latency per operation type of the disk partition used for the MongoDB process
- Amount of free and used disk space on the disk partition used for the MongoDB process
- Measurements for the host, such as CPU usage or number of I/O operations.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return GetProcessMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) GetProcessMeasurements(ctx context.Context, groupId string, processId string) GetProcessMeasurementsApiRequest {
	return GetProcessMeasurementsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// GetProcessMeasurementsExecute executes the request
//
//	@return ApiMeasurementsGeneralViewAtlas
func (a *MonitoringAndLogsApiService) GetProcessMeasurementsExecute(r GetProcessMeasurementsApiRequest) (*ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiMeasurementsGeneralViewAtlas
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.GetProcessMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/measurements"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}

	if r.m != nil {
		t := *r.m
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "m", t, "multi")

	}
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
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

type ListDatabasesApiRequest struct {
	ctx          context.Context
	ApiService   MonitoringAndLogsApi
	groupId      string
	processId    string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListDatabasesApiParams struct {
	GroupId      string
	ProcessId    string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *MonitoringAndLogsApiService) ListDatabasesWithParams(ctx context.Context, args *ListDatabasesApiParams) ListDatabasesApiRequest {
	return ListDatabasesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		processId:    args.ProcessId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListDatabasesApiRequest) IncludeCount(includeCount bool) ListDatabasesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListDatabasesApiRequest) ItemsPerPage(itemsPerPage int) ListDatabasesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListDatabasesApiRequest) PageNum(pageNum int) ListDatabasesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListDatabasesApiRequest) Execute() (*PaginatedDatabase, *http.Response, error) {
	return r.ApiService.ListDatabasesExecute(r)
}

/*
ListDatabases Return Available Databases for One MongoDB Process

Returns the list of databases running on the specified host for the specified project. `M0` free clusters, `M2`, `M5`, serverless, and Flex clusters have some operational limits. The MongoDB Cloud process must be a `mongod`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return ListDatabasesApiRequest
*/
func (a *MonitoringAndLogsApiService) ListDatabases(ctx context.Context, groupId string, processId string) ListDatabasesApiRequest {
	return ListDatabasesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// ListDatabasesExecute executes the request
//
//	@return PaginatedDatabase
func (a *MonitoringAndLogsApiService) ListDatabasesExecute(r ListDatabasesApiRequest) (*PaginatedDatabase, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedDatabase
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListDatabases")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/databases"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

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

type ListGroupProcessesApiRequest struct {
	ctx          context.Context
	ApiService   MonitoringAndLogsApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupProcessesApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *MonitoringAndLogsApiService) ListGroupProcessesWithParams(ctx context.Context, args *ListGroupProcessesApiParams) ListGroupProcessesApiRequest {
	return ListGroupProcessesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupProcessesApiRequest) IncludeCount(includeCount bool) ListGroupProcessesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupProcessesApiRequest) ItemsPerPage(itemsPerPage int) ListGroupProcessesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupProcessesApiRequest) PageNum(pageNum int) ListGroupProcessesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupProcessesApiRequest) Execute() (*PaginatedHostViewAtlas, *http.Response, error) {
	return r.ApiService.ListGroupProcessesExecute(r)
}

/*
ListGroupProcesses Return All MongoDB Processes in One Project

Returns details of all processes for the specified project. A MongoDB process can be either a `mongod` or `mongos`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupProcessesApiRequest
*/
func (a *MonitoringAndLogsApiService) ListGroupProcesses(ctx context.Context, groupId string) ListGroupProcessesApiRequest {
	return ListGroupProcessesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupProcessesExecute executes the request
//
//	@return PaginatedHostViewAtlas
func (a *MonitoringAndLogsApiService) ListGroupProcessesExecute(r ListGroupProcessesApiRequest) (*PaginatedHostViewAtlas, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedHostViewAtlas
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListGroupProcesses")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes"
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

type ListHostFtsMetricsApiRequest struct {
	ctx        context.Context
	ApiService MonitoringAndLogsApi
	processId  string
	groupId    string
}

type ListHostFtsMetricsApiParams struct {
	ProcessId string
	GroupId   string
}

func (a *MonitoringAndLogsApiService) ListHostFtsMetricsWithParams(ctx context.Context, args *ListHostFtsMetricsApiParams) ListHostFtsMetricsApiRequest {
	return ListHostFtsMetricsApiRequest{
		ApiService: a,
		ctx:        ctx,
		processId:  args.ProcessId,
		groupId:    args.GroupId,
	}
}

func (r ListHostFtsMetricsApiRequest) Execute() (*CloudSearchMetrics, *http.Response, error) {
	return r.ApiService.ListHostFtsMetricsExecute(r)
}

/*
ListHostFtsMetrics Return All Atlas Search Metric Types for One Process

Returns all Atlas Search metric types available for one process in the specified project. You must have the Project Read Only or higher role to view the Atlas Search metric types.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListHostFtsMetricsApiRequest
*/
func (a *MonitoringAndLogsApiService) ListHostFtsMetrics(ctx context.Context, processId string, groupId string) ListHostFtsMetricsApiRequest {
	return ListHostFtsMetricsApiRequest{
		ApiService: a,
		ctx:        ctx,
		processId:  processId,
		groupId:    groupId,
	}
}

// ListHostFtsMetricsExecute executes the request
//
//	@return CloudSearchMetrics
func (a *MonitoringAndLogsApiService) ListHostFtsMetricsExecute(r ListHostFtsMetricsApiRequest) (*CloudSearchMetrics, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudSearchMetrics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListHostFtsMetrics")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/hosts/{processId}/fts/metrics"
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)
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

type ListIndexMeasurementsApiRequest struct {
	ctx            context.Context
	ApiService     MonitoringAndLogsApi
	processId      string
	databaseName   string
	collectionName string
	groupId        string
	granularity    *string
	metrics        *[]string
	period         *string
	start          *time.Time
	end            *time.Time
}

type ListIndexMeasurementsApiParams struct {
	ProcessId      string
	DatabaseName   string
	CollectionName string
	GroupId        string
	Granularity    *string
	Metrics        *[]string
	Period         *string
	Start          *time.Time
	End            *time.Time
}

func (a *MonitoringAndLogsApiService) ListIndexMeasurementsWithParams(ctx context.Context, args *ListIndexMeasurementsApiParams) ListIndexMeasurementsApiRequest {
	return ListIndexMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		processId:      args.ProcessId,
		databaseName:   args.DatabaseName,
		collectionName: args.CollectionName,
		groupId:        args.GroupId,
		granularity:    args.Granularity,
		metrics:        args.Metrics,
		period:         args.Period,
		start:          args.Start,
		end:            args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r ListIndexMeasurementsApiRequest) Granularity(granularity string) ListIndexMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// List that contains the measurements that MongoDB Atlas reports for the associated data series.
func (r ListIndexMeasurementsApiRequest) Metrics(metrics []string) ListIndexMeasurementsApiRequest {
	r.metrics = &metrics
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r ListIndexMeasurementsApiRequest) Period(period string) ListIndexMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListIndexMeasurementsApiRequest) Start(start time.Time) ListIndexMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListIndexMeasurementsApiRequest) End(end time.Time) ListIndexMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r ListIndexMeasurementsApiRequest) Execute() (*MeasurementsIndexes, *http.Response, error) {
	return r.ApiService.ListIndexMeasurementsExecute(r)
}

/*
ListIndexMeasurements Return All Atlas Search Index Metrics for One Namespace

Returns the Atlas Search index metrics within the specified time range for one namespace in the specified process.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@param databaseName Human-readable label that identifies the database.
	@param collectionName Human-readable label that identifies the collection.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListIndexMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) ListIndexMeasurements(ctx context.Context, processId string, databaseName string, collectionName string, groupId string) ListIndexMeasurementsApiRequest {
	return ListIndexMeasurementsApiRequest{
		ApiService:     a,
		ctx:            ctx,
		processId:      processId,
		databaseName:   databaseName,
		collectionName: collectionName,
		groupId:        groupId,
	}
}

// ListIndexMeasurementsExecute executes the request
//
//	@return MeasurementsIndexes
func (a *MonitoringAndLogsApiService) ListIndexMeasurementsExecute(r ListIndexMeasurementsApiRequest) (*MeasurementsIndexes, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementsIndexes
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListIndexMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/hosts/{processId}/fts/metrics/indexes/{databaseName}/{collectionName}/measurements"
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}
	if r.metrics == nil {
		return localVarReturnValue, nil, reportError("metrics is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	{
		t := *r.metrics
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "metrics", t, "multi")
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

type ListMeasurementsApiRequest struct {
	ctx         context.Context
	ApiService  MonitoringAndLogsApi
	processId   string
	groupId     string
	granularity *string
	metrics     *[]string
	period      *string
	start       *time.Time
	end         *time.Time
}

type ListMeasurementsApiParams struct {
	ProcessId   string
	GroupId     string
	Granularity *string
	Metrics     *[]string
	Period      *string
	Start       *time.Time
	End         *time.Time
}

func (a *MonitoringAndLogsApiService) ListMeasurementsWithParams(ctx context.Context, args *ListMeasurementsApiParams) ListMeasurementsApiRequest {
	return ListMeasurementsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		processId:   args.ProcessId,
		groupId:     args.GroupId,
		granularity: args.Granularity,
		metrics:     args.Metrics,
		period:      args.Period,
		start:       args.Start,
		end:         args.End,
	}
}

// Duration that specifies the interval at which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC.
func (r ListMeasurementsApiRequest) Granularity(granularity string) ListMeasurementsApiRequest {
	r.granularity = &granularity
	return r
}

// List that contains the metrics that you want MongoDB Atlas to report for the associated data series. If you don&#39;t set this parameter, this resource returns all hardware and status metrics for the associated data series.
func (r ListMeasurementsApiRequest) Metrics(metrics []string) ListMeasurementsApiRequest {
	r.metrics = &metrics
	return r
}

// Duration over which Atlas reports the metrics. This parameter expresses its value in the ISO 8601 duration format in UTC. Include this parameter when you do not set **start** and **end**.
func (r ListMeasurementsApiRequest) Period(period string) ListMeasurementsApiRequest {
	r.period = &period
	return r
}

// Date and time when MongoDB Cloud begins reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListMeasurementsApiRequest) Start(start time.Time) ListMeasurementsApiRequest {
	r.start = &start
	return r
}

// Date and time when MongoDB Cloud stops reporting the metrics. This parameter expresses its value in the ISO 8601 timestamp format in UTC. Include this parameter when you do not set **period**.
func (r ListMeasurementsApiRequest) End(end time.Time) ListMeasurementsApiRequest {
	r.end = &end
	return r
}

func (r ListMeasurementsApiRequest) Execute() (*MeasurementsNonIndex, *http.Response, error) {
	return r.ApiService.ListMeasurementsExecute(r)
}

/*
ListMeasurements Return Atlas Search Hardware and Status Metrics

Returns the Atlas Search hardware and status data series within the provided time range for one process in the specified project. You must have the Project Read Only or higher role to view the Atlas Search metric types.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param processId Combination of hostname and IANA port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (mongod or mongos). The port must be the IANA port on which the MongoDB process listens for requests.
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListMeasurementsApiRequest
*/
func (a *MonitoringAndLogsApiService) ListMeasurements(ctx context.Context, processId string, groupId string) ListMeasurementsApiRequest {
	return ListMeasurementsApiRequest{
		ApiService: a,
		ctx:        ctx,
		processId:  processId,
		groupId:    groupId,
	}
}

// ListMeasurementsExecute executes the request
//
//	@return MeasurementsNonIndex
func (a *MonitoringAndLogsApiService) ListMeasurementsExecute(r ListMeasurementsApiRequest) (*MeasurementsNonIndex, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *MeasurementsNonIndex
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListMeasurements")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/hosts/{processId}/fts/metrics/measurements"
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.granularity == nil {
		return localVarReturnValue, nil, reportError("granularity is required and must be specified")
	}
	if r.metrics == nil {
		return localVarReturnValue, nil, reportError("metrics is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "granularity", r.granularity, "")
	if r.period != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "period", r.period, "")
	}
	if r.start != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "start", r.start, "")
	}
	if r.end != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "end", r.end, "")
	}
	{
		t := *r.metrics
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "metrics", t, "multi")
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

type ListProcessDisksApiRequest struct {
	ctx          context.Context
	ApiService   MonitoringAndLogsApi
	groupId      string
	processId    string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListProcessDisksApiParams struct {
	GroupId      string
	ProcessId    string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *MonitoringAndLogsApiService) ListProcessDisksWithParams(ctx context.Context, args *ListProcessDisksApiParams) ListProcessDisksApiRequest {
	return ListProcessDisksApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		processId:    args.ProcessId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListProcessDisksApiRequest) IncludeCount(includeCount bool) ListProcessDisksApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListProcessDisksApiRequest) ItemsPerPage(itemsPerPage int) ListProcessDisksApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListProcessDisksApiRequest) PageNum(pageNum int) ListProcessDisksApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListProcessDisksApiRequest) Execute() (*PaginatedDiskPartition, *http.Response, error) {
	return r.ApiService.ListProcessDisksExecute(r)
}

/*
ListProcessDisks Return Available Disks for One MongoDB Process

Returns the list of disks or partitions for the specified host for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param processId Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	@return ListProcessDisksApiRequest
*/
func (a *MonitoringAndLogsApiService) ListProcessDisks(ctx context.Context, groupId string, processId string) ListProcessDisksApiRequest {
	return ListProcessDisksApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		processId:  processId,
	}
}

// ListProcessDisksExecute executes the request
//
//	@return PaginatedDiskPartition
func (a *MonitoringAndLogsApiService) ListProcessDisksExecute(r ListProcessDisksApiRequest) (*PaginatedDiskPartition, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedDiskPartition
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MonitoringAndLogsApiService.ListProcessDisks")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/processes/{processId}/disks"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.processId == "" {
		return localVarReturnValue, nil, reportError("processId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processId"+"}", url.PathEscape(r.processId), -1)

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
