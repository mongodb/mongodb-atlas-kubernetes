// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type StreamsApi interface {

	/*
		AcceptVpcPeeringConnection Accept One Incoming VPC Peering Connection

		Requests the acceptance of an incoming VPC Peering connection.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id The VPC Peering Connection id.
		@param vPCPeeringActionChallenge Challenge values for VPC Peering requester account ID, and requester VPC ID.
		@return AcceptVpcPeeringConnectionApiRequest
	*/
	AcceptVpcPeeringConnection(ctx context.Context, groupId string, id string, vPCPeeringActionChallenge *VPCPeeringActionChallenge) AcceptVpcPeeringConnectionApiRequest
	/*
		AcceptVpcPeeringConnection Accept One Incoming VPC Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AcceptVpcPeeringConnectionApiParams - Parameters for the request
		@return AcceptVpcPeeringConnectionApiRequest
	*/
	AcceptVpcPeeringConnectionWithParams(ctx context.Context, args *AcceptVpcPeeringConnectionApiParams) AcceptVpcPeeringConnectionApiRequest

	// Method available only for mocking purposes
	AcceptVpcPeeringConnectionExecute(r AcceptVpcPeeringConnectionApiRequest) (*http.Response, error)

	/*
		CreatePrivateLinkConnection Create One Private Link Connection

		Creates one Private Link in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param streamsPrivateLinkConnection Details to create one Private Link connection for a project. project.
		@return CreatePrivateLinkConnectionApiRequest
	*/
	CreatePrivateLinkConnection(ctx context.Context, groupId string, streamsPrivateLinkConnection *StreamsPrivateLinkConnection) CreatePrivateLinkConnectionApiRequest
	/*
		CreatePrivateLinkConnection Create One Private Link Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreatePrivateLinkConnectionApiParams - Parameters for the request
		@return CreatePrivateLinkConnectionApiRequest
	*/
	CreatePrivateLinkConnectionWithParams(ctx context.Context, args *CreatePrivateLinkConnectionApiParams) CreatePrivateLinkConnectionApiRequest

	// Method available only for mocking purposes
	CreatePrivateLinkConnectionExecute(r CreatePrivateLinkConnectionApiRequest) (*StreamsPrivateLinkConnection, *http.Response, error)

	/*
		CreateStreamConnection Create One Stream Connection

		Creates one connection for a stream workspace in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param streamsConnection Details to create one connection for a streams workspace in the specified project.
		@return CreateStreamConnectionApiRequest
	*/
	CreateStreamConnection(ctx context.Context, groupId string, tenantName string, streamsConnection *StreamsConnection) CreateStreamConnectionApiRequest
	/*
		CreateStreamConnection Create One Stream Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateStreamConnectionApiParams - Parameters for the request
		@return CreateStreamConnectionApiRequest
	*/
	CreateStreamConnectionWithParams(ctx context.Context, args *CreateStreamConnectionApiParams) CreateStreamConnectionApiRequest

	// Method available only for mocking purposes
	CreateStreamConnectionExecute(r CreateStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error)

	/*
		CreateStreamProcessor Create One Stream Processor

		Create one Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param streamsProcessor Details to create an Atlas Streams Processor.
		@return CreateStreamProcessorApiRequest
	*/
	CreateStreamProcessor(ctx context.Context, groupId string, tenantName string, streamsProcessor *StreamsProcessor) CreateStreamProcessorApiRequest
	/*
		CreateStreamProcessor Create One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateStreamProcessorApiParams - Parameters for the request
		@return CreateStreamProcessorApiRequest
	*/
	CreateStreamProcessorWithParams(ctx context.Context, args *CreateStreamProcessorApiParams) CreateStreamProcessorApiRequest

	// Method available only for mocking purposes
	CreateStreamProcessorExecute(r CreateStreamProcessorApiRequest) (*StreamsProcessor, *http.Response, error)

	/*
		CreateStreamWorkspace Create One Stream Workspace

		Creates one stream workspace in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param streamsTenant Details to create one streams workspace in the specified project.
		@return CreateStreamWorkspaceApiRequest
	*/
	CreateStreamWorkspace(ctx context.Context, groupId string, streamsTenant *StreamsTenant) CreateStreamWorkspaceApiRequest
	/*
		CreateStreamWorkspace Create One Stream Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateStreamWorkspaceApiParams - Parameters for the request
		@return CreateStreamWorkspaceApiRequest
	*/
	CreateStreamWorkspaceWithParams(ctx context.Context, args *CreateStreamWorkspaceApiParams) CreateStreamWorkspaceApiRequest

	// Method available only for mocking purposes
	CreateStreamWorkspaceExecute(r CreateStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error)

	/*
		DeletePrivateLinkConnection Delete One Private Link Connection

		Deletes one Private Link in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param connectionId Unique ID that identifies the Private Link connection.
		@return DeletePrivateLinkConnectionApiRequest
	*/
	DeletePrivateLinkConnection(ctx context.Context, groupId string, connectionId string) DeletePrivateLinkConnectionApiRequest
	/*
		DeletePrivateLinkConnection Delete One Private Link Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeletePrivateLinkConnectionApiParams - Parameters for the request
		@return DeletePrivateLinkConnectionApiRequest
	*/
	DeletePrivateLinkConnectionWithParams(ctx context.Context, args *DeletePrivateLinkConnectionApiParams) DeletePrivateLinkConnectionApiRequest

	// Method available only for mocking purposes
	DeletePrivateLinkConnectionExecute(r DeletePrivateLinkConnectionApiRequest) (*http.Response, error)

	/*
		DeleteStreamConnection Delete One Stream Connection

		Delete one connection of the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param connectionName Label that identifies the stream connection.
		@return DeleteStreamConnectionApiRequest
	*/
	DeleteStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string) DeleteStreamConnectionApiRequest
	/*
		DeleteStreamConnection Delete One Stream Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteStreamConnectionApiParams - Parameters for the request
		@return DeleteStreamConnectionApiRequest
	*/
	DeleteStreamConnectionWithParams(ctx context.Context, args *DeleteStreamConnectionApiParams) DeleteStreamConnectionApiRequest

	// Method available only for mocking purposes
	DeleteStreamConnectionExecute(r DeleteStreamConnectionApiRequest) (*http.Response, error)

	/*
		DeleteStreamProcessor Delete One Stream Processor

		Delete a Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@return DeleteStreamProcessorApiRequest
	*/
	DeleteStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) DeleteStreamProcessorApiRequest
	/*
		DeleteStreamProcessor Delete One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteStreamProcessorApiParams - Parameters for the request
		@return DeleteStreamProcessorApiRequest
	*/
	DeleteStreamProcessorWithParams(ctx context.Context, args *DeleteStreamProcessorApiParams) DeleteStreamProcessorApiRequest

	// Method available only for mocking purposes
	DeleteStreamProcessorExecute(r DeleteStreamProcessorApiRequest) (*http.Response, error)

	/*
		DeleteStreamWorkspace Delete One Stream Workspace

		Delete one stream workspace in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace to delete.
		@return DeleteStreamWorkspaceApiRequest
	*/
	DeleteStreamWorkspace(ctx context.Context, groupId string, tenantName string) DeleteStreamWorkspaceApiRequest
	/*
		DeleteStreamWorkspace Delete One Stream Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteStreamWorkspaceApiParams - Parameters for the request
		@return DeleteStreamWorkspaceApiRequest
	*/
	DeleteStreamWorkspaceWithParams(ctx context.Context, args *DeleteStreamWorkspaceApiParams) DeleteStreamWorkspaceApiRequest

	// Method available only for mocking purposes
	DeleteStreamWorkspaceExecute(r DeleteStreamWorkspaceApiRequest) (*http.Response, error)

	/*
		DeleteVpcPeeringConnection Delete One VPC Peering Connection

		Deletes an incoming VPC Peering connection.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id The VPC Peering Connection id.
		@return DeleteVpcPeeringConnectionApiRequest
	*/
	DeleteVpcPeeringConnection(ctx context.Context, groupId string, id string) DeleteVpcPeeringConnectionApiRequest
	/*
		DeleteVpcPeeringConnection Delete One VPC Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteVpcPeeringConnectionApiParams - Parameters for the request
		@return DeleteVpcPeeringConnectionApiRequest
	*/
	DeleteVpcPeeringConnectionWithParams(ctx context.Context, args *DeleteVpcPeeringConnectionApiParams) DeleteVpcPeeringConnectionApiRequest

	// Method available only for mocking purposes
	DeleteVpcPeeringConnectionExecute(r DeleteVpcPeeringConnectionApiRequest) (*http.Response, error)

	/*
		DownloadAuditLogs Download Audit Logs for One Atlas Stream Processing Workspace

		Downloads the audit logs for the specified Atlas Streams Processing workspace or stream processor. By default, logs cover periods of 30 days. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@return DownloadAuditLogsApiRequest
	*/
	DownloadAuditLogs(ctx context.Context, groupId string, tenantName string) DownloadAuditLogsApiRequest
	/*
		DownloadAuditLogs Download Audit Logs for One Atlas Stream Processing Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadAuditLogsApiParams - Parameters for the request
		@return DownloadAuditLogsApiRequest
	*/
	DownloadAuditLogsWithParams(ctx context.Context, args *DownloadAuditLogsApiParams) DownloadAuditLogsApiRequest

	// Method available only for mocking purposes
	DownloadAuditLogsExecute(r DownloadAuditLogsApiRequest) (io.ReadCloser, *http.Response, error)

	/*
		DownloadOperationalLogs Download Operational Logs for One Atlas Stream Processing Workspace

		Downloads the operational logs for the specified Atlas Streams Processing workspace or stream processor. By default, logs cover periods of 30 days. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: "Accept: application/vnd.atlas.2025-03-12+gzip".

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@return DownloadOperationalLogsApiRequest
	*/
	DownloadOperationalLogs(ctx context.Context, groupId string, tenantName string) DownloadOperationalLogsApiRequest
	/*
		DownloadOperationalLogs Download Operational Logs for One Atlas Stream Processing Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DownloadOperationalLogsApiParams - Parameters for the request
		@return DownloadOperationalLogsApiRequest
	*/
	DownloadOperationalLogsWithParams(ctx context.Context, args *DownloadOperationalLogsApiParams) DownloadOperationalLogsApiRequest

	// Method available only for mocking purposes
	DownloadOperationalLogsExecute(r DownloadOperationalLogsApiRequest) (io.ReadCloser, *http.Response, error)

	/*
		GetAccountDetails Return Account ID and VPC ID for One Project and Region

		Returns the Account ID, and the VPC ID for the group and region specified.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return GetAccountDetailsApiRequest
	*/
	GetAccountDetails(ctx context.Context, groupId string) GetAccountDetailsApiRequest
	/*
		GetAccountDetails Return Account ID and VPC ID for One Project and Region


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetAccountDetailsApiParams - Parameters for the request
		@return GetAccountDetailsApiRequest
	*/
	GetAccountDetailsWithParams(ctx context.Context, args *GetAccountDetailsApiParams) GetAccountDetailsApiRequest

	// Method available only for mocking purposes
	GetAccountDetailsExecute(r GetAccountDetailsApiRequest) (*AccountDetails, *http.Response, error)

	/*
		GetPrivateLinkConnection Return One Private Link Connection

		Returns the details of one Private Link connection within the project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param connectionId Unique ID that identifies the Private Link connection.
		@return GetPrivateLinkConnectionApiRequest
	*/
	GetPrivateLinkConnection(ctx context.Context, groupId string, connectionId string) GetPrivateLinkConnectionApiRequest
	/*
		GetPrivateLinkConnection Return One Private Link Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetPrivateLinkConnectionApiParams - Parameters for the request
		@return GetPrivateLinkConnectionApiRequest
	*/
	GetPrivateLinkConnectionWithParams(ctx context.Context, args *GetPrivateLinkConnectionApiParams) GetPrivateLinkConnectionApiRequest

	// Method available only for mocking purposes
	GetPrivateLinkConnectionExecute(r GetPrivateLinkConnectionApiRequest) (*StreamsPrivateLinkConnection, *http.Response, error)

	/*
		GetStreamConnection Return One Stream Connection

		Returns the details of one stream connection within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace to return.
		@param connectionName Label that identifies the stream connection to return.
		@return GetStreamConnectionApiRequest
	*/
	GetStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string) GetStreamConnectionApiRequest
	/*
		GetStreamConnection Return One Stream Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetStreamConnectionApiParams - Parameters for the request
		@return GetStreamConnectionApiRequest
	*/
	GetStreamConnectionWithParams(ctx context.Context, args *GetStreamConnectionApiParams) GetStreamConnectionApiRequest

	// Method available only for mocking purposes
	GetStreamConnectionExecute(r GetStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error)

	/*
		GetStreamProcessor Return One Stream Processor

		Get one Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@return GetStreamProcessorApiRequest
	*/
	GetStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) GetStreamProcessorApiRequest
	/*
		GetStreamProcessor Return One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetStreamProcessorApiParams - Parameters for the request
		@return GetStreamProcessorApiRequest
	*/
	GetStreamProcessorWithParams(ctx context.Context, args *GetStreamProcessorApiParams) GetStreamProcessorApiRequest

	// Method available only for mocking purposes
	GetStreamProcessorExecute(r GetStreamProcessorApiRequest) (*StreamsProcessorWithStats, *http.Response, error)

	/*
		GetStreamProcessors Return All Stream Processors in One Stream Workspace

		Returns all Stream Processors within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@return GetStreamProcessorsApiRequest
	*/
	GetStreamProcessors(ctx context.Context, groupId string, tenantName string) GetStreamProcessorsApiRequest
	/*
		GetStreamProcessors Return All Stream Processors in One Stream Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetStreamProcessorsApiParams - Parameters for the request
		@return GetStreamProcessorsApiRequest
	*/
	GetStreamProcessorsWithParams(ctx context.Context, args *GetStreamProcessorsApiParams) GetStreamProcessorsApiRequest

	// Method available only for mocking purposes
	GetStreamProcessorsExecute(r GetStreamProcessorsApiRequest) (*PaginatedApiStreamsStreamProcessorWithStats, *http.Response, error)

	/*
		GetStreamWorkspace Return One Stream Workspace

		Returns the details of one stream workspace within the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace to return.
		@return GetStreamWorkspaceApiRequest
	*/
	GetStreamWorkspace(ctx context.Context, groupId string, tenantName string) GetStreamWorkspaceApiRequest
	/*
		GetStreamWorkspace Return One Stream Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetStreamWorkspaceApiParams - Parameters for the request
		@return GetStreamWorkspaceApiRequest
	*/
	GetStreamWorkspaceWithParams(ctx context.Context, args *GetStreamWorkspaceApiParams) GetStreamWorkspaceApiRequest

	// Method available only for mocking purposes
	GetStreamWorkspaceExecute(r GetStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error)

	/*
		ListActivePeeringConnections Return All Active Incoming VPC Peering Connections

		Returns a list of active incoming VPC Peering Connections.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListActivePeeringConnectionsApiRequest
	*/
	ListActivePeeringConnections(ctx context.Context, groupId string) ListActivePeeringConnectionsApiRequest
	/*
		ListActivePeeringConnections Return All Active Incoming VPC Peering Connections


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListActivePeeringConnectionsApiParams - Parameters for the request
		@return ListActivePeeringConnectionsApiRequest
	*/
	ListActivePeeringConnectionsWithParams(ctx context.Context, args *ListActivePeeringConnectionsApiParams) ListActivePeeringConnectionsApiRequest

	// Method available only for mocking purposes
	ListActivePeeringConnectionsExecute(r ListActivePeeringConnectionsApiRequest) (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error)

	/*
		ListPrivateLinkConnections Return All Private Link Connections

		Returns all Private Link connections for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListPrivateLinkConnectionsApiRequest
	*/
	ListPrivateLinkConnections(ctx context.Context, groupId string) ListPrivateLinkConnectionsApiRequest
	/*
		ListPrivateLinkConnections Return All Private Link Connections


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListPrivateLinkConnectionsApiParams - Parameters for the request
		@return ListPrivateLinkConnectionsApiRequest
	*/
	ListPrivateLinkConnectionsWithParams(ctx context.Context, args *ListPrivateLinkConnectionsApiParams) ListPrivateLinkConnectionsApiRequest

	// Method available only for mocking purposes
	ListPrivateLinkConnectionsExecute(r ListPrivateLinkConnectionsApiRequest) (*PaginatedApiStreamsPrivateLink, *http.Response, error)

	/*
		ListStreamConnections Return All Connections of the Stream Workspaces

		Returns all connections of the stream workspace for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@return ListStreamConnectionsApiRequest
	*/
	ListStreamConnections(ctx context.Context, groupId string, tenantName string) ListStreamConnectionsApiRequest
	/*
		ListStreamConnections Return All Connections of the Stream Workspaces


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListStreamConnectionsApiParams - Parameters for the request
		@return ListStreamConnectionsApiRequest
	*/
	ListStreamConnectionsWithParams(ctx context.Context, args *ListStreamConnectionsApiParams) ListStreamConnectionsApiRequest

	// Method available only for mocking purposes
	ListStreamConnectionsExecute(r ListStreamConnectionsApiRequest) (*PaginatedApiStreamsConnection, *http.Response, error)

	/*
		ListStreamWorkspaces Return All Stream Workspaces in One Project

		Returns all stream workspaces for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListStreamWorkspacesApiRequest
	*/
	ListStreamWorkspaces(ctx context.Context, groupId string) ListStreamWorkspacesApiRequest
	/*
		ListStreamWorkspaces Return All Stream Workspaces in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListStreamWorkspacesApiParams - Parameters for the request
		@return ListStreamWorkspacesApiRequest
	*/
	ListStreamWorkspacesWithParams(ctx context.Context, args *ListStreamWorkspacesApiParams) ListStreamWorkspacesApiRequest

	// Method available only for mocking purposes
	ListStreamWorkspacesExecute(r ListStreamWorkspacesApiRequest) (*PaginatedApiStreamsTenant, *http.Response, error)

	/*
		ListVpcPeeringConnections Return All VPC Peering Connections

		Returns a list of incoming VPC Peering Connections.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListVpcPeeringConnectionsApiRequest
	*/
	ListVpcPeeringConnections(ctx context.Context, groupId string) ListVpcPeeringConnectionsApiRequest
	/*
		ListVpcPeeringConnections Return All VPC Peering Connections


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListVpcPeeringConnectionsApiParams - Parameters for the request
		@return ListVpcPeeringConnectionsApiRequest
	*/
	ListVpcPeeringConnectionsWithParams(ctx context.Context, args *ListVpcPeeringConnectionsApiParams) ListVpcPeeringConnectionsApiRequest

	// Method available only for mocking purposes
	ListVpcPeeringConnectionsExecute(r ListVpcPeeringConnectionsApiRequest) (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error)

	/*
		RejectVpcPeeringConnection Reject One Incoming VPC Peering Connection

		Requests the rejection of an incoming VPC Peering connection.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param id The VPC Peering Connection id.
		@return RejectVpcPeeringConnectionApiRequest
	*/
	RejectVpcPeeringConnection(ctx context.Context, groupId string, id string) RejectVpcPeeringConnectionApiRequest
	/*
		RejectVpcPeeringConnection Reject One Incoming VPC Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RejectVpcPeeringConnectionApiParams - Parameters for the request
		@return RejectVpcPeeringConnectionApiRequest
	*/
	RejectVpcPeeringConnectionWithParams(ctx context.Context, args *RejectVpcPeeringConnectionApiParams) RejectVpcPeeringConnectionApiRequest

	// Method available only for mocking purposes
	RejectVpcPeeringConnectionExecute(r RejectVpcPeeringConnectionApiRequest) (*http.Response, error)

	/*
		StartStreamProcessor Start One Stream Processor

		Start a Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@return StartStreamProcessorApiRequest
	*/
	StartStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) StartStreamProcessorApiRequest
	/*
		StartStreamProcessor Start One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param StartStreamProcessorApiParams - Parameters for the request
		@return StartStreamProcessorApiRequest
	*/
	StartStreamProcessorWithParams(ctx context.Context, args *StartStreamProcessorApiParams) StartStreamProcessorApiRequest

	// Method available only for mocking purposes
	StartStreamProcessorExecute(r StartStreamProcessorApiRequest) (*http.Response, error)

	/*
		StartStreamProcessorWith Start One Stream Processor With Options

		Start a Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@param streamsStartStreamProcessorWith Options for starting a stream processor.
		@return StartStreamProcessorWithApiRequest
	*/
	StartStreamProcessorWith(ctx context.Context, groupId string, tenantName string, processorName string, streamsStartStreamProcessorWith *StreamsStartStreamProcessorWith) StartStreamProcessorWithApiRequest
	/*
		StartStreamProcessorWith Start One Stream Processor With Options


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param StartStreamProcessorWithApiParams - Parameters for the request
		@return StartStreamProcessorWithApiRequest
	*/
	StartStreamProcessorWithWithParams(ctx context.Context, args *StartStreamProcessorWithApiParams) StartStreamProcessorWithApiRequest

	// Method available only for mocking purposes
	StartStreamProcessorWithExecute(r StartStreamProcessorWithApiRequest) (*http.Response, error)

	/*
		StopStreamProcessor Stop One Stream Processor

		Stop a Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@return StopStreamProcessorApiRequest
	*/
	StopStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) StopStreamProcessorApiRequest
	/*
		StopStreamProcessor Stop One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param StopStreamProcessorApiParams - Parameters for the request
		@return StopStreamProcessorApiRequest
	*/
	StopStreamProcessorWithParams(ctx context.Context, args *StopStreamProcessorApiParams) StopStreamProcessorApiRequest

	// Method available only for mocking purposes
	StopStreamProcessorExecute(r StopStreamProcessorApiRequest) (*http.Response, error)

	/*
		UpdateStreamConnection Update One Stream Connection

		Update one connection for the specified stream workspace in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param connectionName Label that identifies the stream connection.
		@param streamsConnection Details to update one connection for a streams workspace in the specified project.
		@return UpdateStreamConnectionApiRequest
	*/
	UpdateStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string, streamsConnection *StreamsConnection) UpdateStreamConnectionApiRequest
	/*
		UpdateStreamConnection Update One Stream Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateStreamConnectionApiParams - Parameters for the request
		@return UpdateStreamConnectionApiRequest
	*/
	UpdateStreamConnectionWithParams(ctx context.Context, args *UpdateStreamConnectionApiParams) UpdateStreamConnectionApiRequest

	// Method available only for mocking purposes
	UpdateStreamConnectionExecute(r UpdateStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error)

	/*
		UpdateStreamProcessor Update One Stream Processor

		Modify one existing Stream Processor within the specified stream workspace.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace.
		@param processorName Label that identifies the stream processor.
		@param streamsModifyStreamProcessor Modifications to apply to the stream processor.
		@return UpdateStreamProcessorApiRequest
	*/
	UpdateStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string, streamsModifyStreamProcessor *StreamsModifyStreamProcessor) UpdateStreamProcessorApiRequest
	/*
		UpdateStreamProcessor Update One Stream Processor


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateStreamProcessorApiParams - Parameters for the request
		@return UpdateStreamProcessorApiRequest
	*/
	UpdateStreamProcessorWithParams(ctx context.Context, args *UpdateStreamProcessorApiParams) UpdateStreamProcessorApiRequest

	// Method available only for mocking purposes
	UpdateStreamProcessorExecute(r UpdateStreamProcessorApiRequest) (*StreamsProcessorWithStats, *http.Response, error)

	/*
		UpdateStreamWorkspace Update One Stream Workspace

		Update one stream workspace in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param tenantName Label that identifies the stream workspace to update.
		@param streamsTenantUpdateRequest Details to update in the streams workspace.
		@return UpdateStreamWorkspaceApiRequest
	*/
	UpdateStreamWorkspace(ctx context.Context, groupId string, tenantName string, streamsTenantUpdateRequest *StreamsTenantUpdateRequest) UpdateStreamWorkspaceApiRequest
	/*
		UpdateStreamWorkspace Update One Stream Workspace


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateStreamWorkspaceApiParams - Parameters for the request
		@return UpdateStreamWorkspaceApiRequest
	*/
	UpdateStreamWorkspaceWithParams(ctx context.Context, args *UpdateStreamWorkspaceApiParams) UpdateStreamWorkspaceApiRequest

	// Method available only for mocking purposes
	UpdateStreamWorkspaceExecute(r UpdateStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error)

	/*
		WithStreamSampleConnections Create One Stream Workspace with Sample Connections

		Creates one stream workspace in the specified project with sample connections.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param body Details to create one streams workspace in the specified project.
		@return WithStreamSampleConnectionsApiRequest
	*/
	WithStreamSampleConnections(ctx context.Context, groupId string, body *any) WithStreamSampleConnectionsApiRequest
	/*
		WithStreamSampleConnections Create One Stream Workspace with Sample Connections


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param WithStreamSampleConnectionsApiParams - Parameters for the request
		@return WithStreamSampleConnectionsApiRequest
	*/
	WithStreamSampleConnectionsWithParams(ctx context.Context, args *WithStreamSampleConnectionsApiParams) WithStreamSampleConnectionsApiRequest

	// Method available only for mocking purposes
	WithStreamSampleConnectionsExecute(r WithStreamSampleConnectionsApiRequest) (*StreamsTenant, *http.Response, error)
}

// StreamsApiService StreamsApi service
type StreamsApiService service

type AcceptVpcPeeringConnectionApiRequest struct {
	ctx                       context.Context
	ApiService                StreamsApi
	groupId                   string
	id                        string
	vPCPeeringActionChallenge *VPCPeeringActionChallenge
}

type AcceptVpcPeeringConnectionApiParams struct {
	GroupId                   string
	Id                        string
	VPCPeeringActionChallenge *VPCPeeringActionChallenge
}

func (a *StreamsApiService) AcceptVpcPeeringConnectionWithParams(ctx context.Context, args *AcceptVpcPeeringConnectionApiParams) AcceptVpcPeeringConnectionApiRequest {
	return AcceptVpcPeeringConnectionApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   args.GroupId,
		id:                        args.Id,
		vPCPeeringActionChallenge: args.VPCPeeringActionChallenge,
	}
}

func (r AcceptVpcPeeringConnectionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.AcceptVpcPeeringConnectionExecute(r)
}

/*
AcceptVpcPeeringConnection Accept One Incoming VPC Peering Connection

Requests the acceptance of an incoming VPC Peering connection.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id The VPC Peering Connection id.
	@return AcceptVpcPeeringConnectionApiRequest
*/
func (a *StreamsApiService) AcceptVpcPeeringConnection(ctx context.Context, groupId string, id string, vPCPeeringActionChallenge *VPCPeeringActionChallenge) AcceptVpcPeeringConnectionApiRequest {
	return AcceptVpcPeeringConnectionApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		groupId:                   groupId,
		id:                        id,
		vPCPeeringActionChallenge: vPCPeeringActionChallenge,
	}
}

// AcceptVpcPeeringConnectionExecute executes the request
func (a *StreamsApiService) AcceptVpcPeeringConnectionExecute(r AcceptVpcPeeringConnectionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.AcceptVpcPeeringConnection")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/vpcPeeringConnections/{id}:accept"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.vPCPeeringActionChallenge == nil {
		return nil, reportError("vPCPeeringActionChallenge is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.vPCPeeringActionChallenge
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

type CreatePrivateLinkConnectionApiRequest struct {
	ctx                          context.Context
	ApiService                   StreamsApi
	groupId                      string
	streamsPrivateLinkConnection *StreamsPrivateLinkConnection
}

type CreatePrivateLinkConnectionApiParams struct {
	GroupId                      string
	StreamsPrivateLinkConnection *StreamsPrivateLinkConnection
}

func (a *StreamsApiService) CreatePrivateLinkConnectionWithParams(ctx context.Context, args *CreatePrivateLinkConnectionApiParams) CreatePrivateLinkConnectionApiRequest {
	return CreatePrivateLinkConnectionApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		streamsPrivateLinkConnection: args.StreamsPrivateLinkConnection,
	}
}

func (r CreatePrivateLinkConnectionApiRequest) Execute() (*StreamsPrivateLinkConnection, *http.Response, error) {
	return r.ApiService.CreatePrivateLinkConnectionExecute(r)
}

/*
CreatePrivateLinkConnection Create One Private Link Connection

Creates one Private Link in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreatePrivateLinkConnectionApiRequest
*/
func (a *StreamsApiService) CreatePrivateLinkConnection(ctx context.Context, groupId string, streamsPrivateLinkConnection *StreamsPrivateLinkConnection) CreatePrivateLinkConnectionApiRequest {
	return CreatePrivateLinkConnectionApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		streamsPrivateLinkConnection: streamsPrivateLinkConnection,
	}
}

// CreatePrivateLinkConnectionExecute executes the request
//
//	@return StreamsPrivateLinkConnection
func (a *StreamsApiService) CreatePrivateLinkConnectionExecute(r CreatePrivateLinkConnectionApiRequest) (*StreamsPrivateLinkConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsPrivateLinkConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.CreatePrivateLinkConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/privateLinkConnections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsPrivateLinkConnection == nil {
		return localVarReturnValue, nil, reportError("streamsPrivateLinkConnection is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsPrivateLinkConnection
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

type CreateStreamConnectionApiRequest struct {
	ctx               context.Context
	ApiService        StreamsApi
	groupId           string
	tenantName        string
	streamsConnection *StreamsConnection
}

type CreateStreamConnectionApiParams struct {
	GroupId           string
	TenantName        string
	StreamsConnection *StreamsConnection
}

func (a *StreamsApiService) CreateStreamConnectionWithParams(ctx context.Context, args *CreateStreamConnectionApiParams) CreateStreamConnectionApiRequest {
	return CreateStreamConnectionApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		tenantName:        args.TenantName,
		streamsConnection: args.StreamsConnection,
	}
}

func (r CreateStreamConnectionApiRequest) Execute() (*StreamsConnection, *http.Response, error) {
	return r.ApiService.CreateStreamConnectionExecute(r)
}

/*
CreateStreamConnection Create One Stream Connection

Creates one connection for a stream workspace in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return CreateStreamConnectionApiRequest
*/
func (a *StreamsApiService) CreateStreamConnection(ctx context.Context, groupId string, tenantName string, streamsConnection *StreamsConnection) CreateStreamConnectionApiRequest {
	return CreateStreamConnectionApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		tenantName:        tenantName,
		streamsConnection: streamsConnection,
	}
}

// CreateStreamConnectionExecute executes the request
//
//	@return StreamsConnection
func (a *StreamsApiService) CreateStreamConnectionExecute(r CreateStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.CreateStreamConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsConnection == nil {
		return localVarReturnValue, nil, reportError("streamsConnection is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsConnection
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

type CreateStreamProcessorApiRequest struct {
	ctx              context.Context
	ApiService       StreamsApi
	groupId          string
	tenantName       string
	streamsProcessor *StreamsProcessor
}

type CreateStreamProcessorApiParams struct {
	GroupId          string
	TenantName       string
	StreamsProcessor *StreamsProcessor
}

func (a *StreamsApiService) CreateStreamProcessorWithParams(ctx context.Context, args *CreateStreamProcessorApiParams) CreateStreamProcessorApiRequest {
	return CreateStreamProcessorApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          args.GroupId,
		tenantName:       args.TenantName,
		streamsProcessor: args.StreamsProcessor,
	}
}

func (r CreateStreamProcessorApiRequest) Execute() (*StreamsProcessor, *http.Response, error) {
	return r.ApiService.CreateStreamProcessorExecute(r)
}

/*
CreateStreamProcessor Create One Stream Processor

Create one Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return CreateStreamProcessorApiRequest
*/
func (a *StreamsApiService) CreateStreamProcessor(ctx context.Context, groupId string, tenantName string, streamsProcessor *StreamsProcessor) CreateStreamProcessorApiRequest {
	return CreateStreamProcessorApiRequest{
		ApiService:       a,
		ctx:              ctx,
		groupId:          groupId,
		tenantName:       tenantName,
		streamsProcessor: streamsProcessor,
	}
}

// CreateStreamProcessorExecute executes the request
//
//	@return StreamsProcessor
func (a *StreamsApiService) CreateStreamProcessorExecute(r CreateStreamProcessorApiRequest) (*StreamsProcessor, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsProcessor
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.CreateStreamProcessor")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsProcessor == nil {
		return localVarReturnValue, nil, reportError("streamsProcessor is required and must be specified")
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
	localVarPostBody = r.streamsProcessor
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

type CreateStreamWorkspaceApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	streamsTenant *StreamsTenant
}

type CreateStreamWorkspaceApiParams struct {
	GroupId       string
	StreamsTenant *StreamsTenant
}

func (a *StreamsApiService) CreateStreamWorkspaceWithParams(ctx context.Context, args *CreateStreamWorkspaceApiParams) CreateStreamWorkspaceApiRequest {
	return CreateStreamWorkspaceApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		streamsTenant: args.StreamsTenant,
	}
}

func (r CreateStreamWorkspaceApiRequest) Execute() (*StreamsTenant, *http.Response, error) {
	return r.ApiService.CreateStreamWorkspaceExecute(r)
}

/*
CreateStreamWorkspace Create One Stream Workspace

Creates one stream workspace in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateStreamWorkspaceApiRequest
*/
func (a *StreamsApiService) CreateStreamWorkspace(ctx context.Context, groupId string, streamsTenant *StreamsTenant) CreateStreamWorkspaceApiRequest {
	return CreateStreamWorkspaceApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		streamsTenant: streamsTenant,
	}
}

// CreateStreamWorkspaceExecute executes the request
//
//	@return StreamsTenant
func (a *StreamsApiService) CreateStreamWorkspaceExecute(r CreateStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.CreateStreamWorkspace")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsTenant == nil {
		return localVarReturnValue, nil, reportError("streamsTenant is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsTenant
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

type DeletePrivateLinkConnectionApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	connectionId string
}

type DeletePrivateLinkConnectionApiParams struct {
	GroupId      string
	ConnectionId string
}

func (a *StreamsApiService) DeletePrivateLinkConnectionWithParams(ctx context.Context, args *DeletePrivateLinkConnectionApiParams) DeletePrivateLinkConnectionApiRequest {
	return DeletePrivateLinkConnectionApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		connectionId: args.ConnectionId,
	}
}

func (r DeletePrivateLinkConnectionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeletePrivateLinkConnectionExecute(r)
}

/*
DeletePrivateLinkConnection Delete One Private Link Connection

Deletes one Private Link in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param connectionId Unique ID that identifies the Private Link connection.
	@return DeletePrivateLinkConnectionApiRequest
*/
func (a *StreamsApiService) DeletePrivateLinkConnection(ctx context.Context, groupId string, connectionId string) DeletePrivateLinkConnectionApiRequest {
	return DeletePrivateLinkConnectionApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		connectionId: connectionId,
	}
}

// DeletePrivateLinkConnectionExecute executes the request
func (a *StreamsApiService) DeletePrivateLinkConnectionExecute(r DeletePrivateLinkConnectionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DeletePrivateLinkConnection")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/privateLinkConnections/{connectionId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.connectionId == "" {
		return nil, reportError("connectionId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"connectionId"+"}", url.PathEscape(r.connectionId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type DeleteStreamConnectionApiRequest struct {
	ctx            context.Context
	ApiService     StreamsApi
	groupId        string
	tenantName     string
	connectionName string
}

type DeleteStreamConnectionApiParams struct {
	GroupId        string
	TenantName     string
	ConnectionName string
}

func (a *StreamsApiService) DeleteStreamConnectionWithParams(ctx context.Context, args *DeleteStreamConnectionApiParams) DeleteStreamConnectionApiRequest {
	return DeleteStreamConnectionApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		tenantName:     args.TenantName,
		connectionName: args.ConnectionName,
	}
}

func (r DeleteStreamConnectionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteStreamConnectionExecute(r)
}

/*
DeleteStreamConnection Delete One Stream Connection

Delete one connection of the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param connectionName Label that identifies the stream connection.
	@return DeleteStreamConnectionApiRequest
*/
func (a *StreamsApiService) DeleteStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string) DeleteStreamConnectionApiRequest {
	return DeleteStreamConnectionApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		tenantName:     tenantName,
		connectionName: connectionName,
	}
}

// DeleteStreamConnectionExecute executes the request
func (a *StreamsApiService) DeleteStreamConnectionExecute(r DeleteStreamConnectionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DeleteStreamConnection")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.connectionName == "" {
		return nil, reportError("connectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"connectionName"+"}", url.PathEscape(r.connectionName), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type DeleteStreamProcessorApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	tenantName    string
	processorName string
}

type DeleteStreamProcessorApiParams struct {
	GroupId       string
	TenantName    string
	ProcessorName string
}

func (a *StreamsApiService) DeleteStreamProcessorWithParams(ctx context.Context, args *DeleteStreamProcessorApiParams) DeleteStreamProcessorApiRequest {
	return DeleteStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		tenantName:    args.TenantName,
		processorName: args.ProcessorName,
	}
}

func (r DeleteStreamProcessorApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteStreamProcessorExecute(r)
}

/*
DeleteStreamProcessor Delete One Stream Processor

Delete a Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return DeleteStreamProcessorApiRequest
*/
func (a *StreamsApiService) DeleteStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) DeleteStreamProcessorApiRequest {
	return DeleteStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		tenantName:    tenantName,
		processorName: processorName,
	}
}

// DeleteStreamProcessorExecute executes the request
func (a *StreamsApiService) DeleteStreamProcessorExecute(r DeleteStreamProcessorApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DeleteStreamProcessor")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

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

type DeleteStreamWorkspaceApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	tenantName string
}

type DeleteStreamWorkspaceApiParams struct {
	GroupId    string
	TenantName string
}

func (a *StreamsApiService) DeleteStreamWorkspaceWithParams(ctx context.Context, args *DeleteStreamWorkspaceApiParams) DeleteStreamWorkspaceApiRequest {
	return DeleteStreamWorkspaceApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
	}
}

func (r DeleteStreamWorkspaceApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteStreamWorkspaceExecute(r)
}

/*
DeleteStreamWorkspace Delete One Stream Workspace

Delete one stream workspace in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace to delete.
	@return DeleteStreamWorkspaceApiRequest
*/
func (a *StreamsApiService) DeleteStreamWorkspace(ctx context.Context, groupId string, tenantName string) DeleteStreamWorkspaceApiRequest {
	return DeleteStreamWorkspaceApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// DeleteStreamWorkspaceExecute executes the request
func (a *StreamsApiService) DeleteStreamWorkspaceExecute(r DeleteStreamWorkspaceApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DeleteStreamWorkspace")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type DeleteVpcPeeringConnectionApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	id         string
}

type DeleteVpcPeeringConnectionApiParams struct {
	GroupId string
	Id      string
}

func (a *StreamsApiService) DeleteVpcPeeringConnectionWithParams(ctx context.Context, args *DeleteVpcPeeringConnectionApiParams) DeleteVpcPeeringConnectionApiRequest {
	return DeleteVpcPeeringConnectionApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		id:         args.Id,
	}
}

func (r DeleteVpcPeeringConnectionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteVpcPeeringConnectionExecute(r)
}

/*
DeleteVpcPeeringConnection Delete One VPC Peering Connection

Deletes an incoming VPC Peering connection.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id The VPC Peering Connection id.
	@return DeleteVpcPeeringConnectionApiRequest
*/
func (a *StreamsApiService) DeleteVpcPeeringConnection(ctx context.Context, groupId string, id string) DeleteVpcPeeringConnectionApiRequest {
	return DeleteVpcPeeringConnectionApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		id:         id,
	}
}

// DeleteVpcPeeringConnectionExecute executes the request
func (a *StreamsApiService) DeleteVpcPeeringConnectionExecute(r DeleteVpcPeeringConnectionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DeleteVpcPeeringConnection")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/vpcPeeringConnections/{id}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type DownloadAuditLogsApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	tenantName string
	endDate    *int64
	startDate  *int64
	spName     *string
}

type DownloadAuditLogsApiParams struct {
	GroupId    string
	TenantName string
	EndDate    *int64
	StartDate  *int64
	SpName     *string
}

func (a *StreamsApiService) DownloadAuditLogsWithParams(ctx context.Context, args *DownloadAuditLogsApiParams) DownloadAuditLogsApiRequest {
	return DownloadAuditLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
		endDate:    args.EndDate,
		startDate:  args.StartDate,
		spName:     args.SpName,
	}
}

// Timestamp that specifies the end point for the range of log messages to download.  MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadAuditLogsApiRequest) EndDate(endDate int64) DownloadAuditLogsApiRequest {
	r.endDate = &endDate
	return r
}

// Timestamp that specifies the starting point for the range of log messages to download. MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadAuditLogsApiRequest) StartDate(startDate int64) DownloadAuditLogsApiRequest {
	r.startDate = &startDate
	return r
}

// Name of the stream processor to download logs for. An empty string will download logs for all stream processors in the workspace.
func (r DownloadAuditLogsApiRequest) SpName(spName string) DownloadAuditLogsApiRequest {
	r.spName = &spName
	return r
}

func (r DownloadAuditLogsApiRequest) Execute() (io.ReadCloser, *http.Response, error) {
	return r.ApiService.DownloadAuditLogsExecute(r)
}

/*
DownloadAuditLogs Download Audit Logs for One Atlas Stream Processing Workspace

Downloads the audit logs for the specified Atlas Streams Processing workspace or stream processor. By default, logs cover periods of 30 days. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: `Accept: application/vnd.atlas.YYYY-MM-DD+gzip`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return DownloadAuditLogsApiRequest
*/
func (a *StreamsApiService) DownloadAuditLogs(ctx context.Context, groupId string, tenantName string) DownloadAuditLogsApiRequest {
	return DownloadAuditLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// DownloadAuditLogsExecute executes the request
//
//	@return io.ReadCloser
func (a *StreamsApiService) DownloadAuditLogsExecute(r DownloadAuditLogsApiRequest) (io.ReadCloser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue io.ReadCloser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DownloadAuditLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/auditLogs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.endDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endDate", r.endDate, "")
	}
	if r.startDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startDate", r.startDate, "")
	}
	if r.spName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "spName", r.spName, "")
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

type DownloadOperationalLogsApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	tenantName string
	endDate    *int64
	startDate  *int64
	spName     *string
}

type DownloadOperationalLogsApiParams struct {
	GroupId    string
	TenantName string
	EndDate    *int64
	StartDate  *int64
	SpName     *string
}

func (a *StreamsApiService) DownloadOperationalLogsWithParams(ctx context.Context, args *DownloadOperationalLogsApiParams) DownloadOperationalLogsApiRequest {
	return DownloadOperationalLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		tenantName: args.TenantName,
		endDate:    args.EndDate,
		startDate:  args.StartDate,
		spName:     args.SpName,
	}
}

// Timestamp that specifies the end point for the range of log messages to download.  MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadOperationalLogsApiRequest) EndDate(endDate int64) DownloadOperationalLogsApiRequest {
	r.endDate = &endDate
	return r
}

// Timestamp that specifies the starting point for the range of log messages to download. MongoDB Cloud expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch.
func (r DownloadOperationalLogsApiRequest) StartDate(startDate int64) DownloadOperationalLogsApiRequest {
	r.startDate = &startDate
	return r
}

// Name of the stream processor to download logs for. An empty string will download logs for all stream processors in the workspace.
func (r DownloadOperationalLogsApiRequest) SpName(spName string) DownloadOperationalLogsApiRequest {
	r.spName = &spName
	return r
}

func (r DownloadOperationalLogsApiRequest) Execute() (io.ReadCloser, *http.Response, error) {
	return r.ApiService.DownloadOperationalLogsExecute(r)
}

/*
DownloadOperationalLogs Download Operational Logs for One Atlas Stream Processing Workspace

Downloads the operational logs for the specified Atlas Streams Processing workspace or stream processor. By default, logs cover periods of 30 days. The API does not support direct calls with the json response schema. You must request a gzip response schema using an accept header of the format: "Accept: application/vnd.atlas.2025-03-12+gzip".

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return DownloadOperationalLogsApiRequest
*/
func (a *StreamsApiService) DownloadOperationalLogs(ctx context.Context, groupId string, tenantName string) DownloadOperationalLogsApiRequest {
	return DownloadOperationalLogsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// DownloadOperationalLogsExecute executes the request
//
//	@return io.ReadCloser
func (a *StreamsApiService) DownloadOperationalLogsExecute(r DownloadOperationalLogsApiRequest) (io.ReadCloser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue io.ReadCloser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.DownloadOperationalLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}:downloadOperationalLogs"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.endDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endDate", r.endDate, "")
	}
	if r.startDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startDate", r.startDate, "")
	}
	if r.spName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "spName", r.spName, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+gzip"}

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

type GetAccountDetailsApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	cloudProvider *string
	regionName    *string
}

type GetAccountDetailsApiParams struct {
	GroupId       string
	CloudProvider *string
	RegionName    *string
}

func (a *StreamsApiService) GetAccountDetailsWithParams(ctx context.Context, args *GetAccountDetailsApiParams) GetAccountDetailsApiRequest {
	return GetAccountDetailsApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		cloudProvider: args.CloudProvider,
		regionName:    args.RegionName,
	}
}

// One of \&quot;aws\&quot;, \&quot;azure\&quot; or \&quot;gcp\&quot;.
func (r GetAccountDetailsApiRequest) CloudProvider(cloudProvider string) GetAccountDetailsApiRequest {
	r.cloudProvider = &cloudProvider
	return r
}

// The cloud provider specific region name, i.e. \&quot;US_EAST_1\&quot; for cloud provider \&quot;aws\&quot;.
func (r GetAccountDetailsApiRequest) RegionName(regionName string) GetAccountDetailsApiRequest {
	r.regionName = &regionName
	return r
}

func (r GetAccountDetailsApiRequest) Execute() (*AccountDetails, *http.Response, error) {
	return r.ApiService.GetAccountDetailsExecute(r)
}

/*
GetAccountDetails Return Account ID and VPC ID for One Project and Region

Returns the Account ID, and the VPC ID for the group and region specified.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return GetAccountDetailsApiRequest
*/
func (a *StreamsApiService) GetAccountDetails(ctx context.Context, groupId string) GetAccountDetailsApiRequest {
	return GetAccountDetailsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// GetAccountDetailsExecute executes the request
//
//	@return AccountDetails
func (a *StreamsApiService) GetAccountDetailsExecute(r GetAccountDetailsApiRequest) (*AccountDetails, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AccountDetails
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetAccountDetails")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/accountDetails"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudProvider == nil {
		return localVarReturnValue, nil, reportError("cloudProvider is required and must be specified")
	}
	if r.regionName == nil {
		return localVarReturnValue, nil, reportError("regionName is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "cloudProvider", r.cloudProvider, "")
	parameterAddToHeaderOrQuery(localVarQueryParams, "regionName", r.regionName, "")
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

type GetPrivateLinkConnectionApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	connectionId string
}

type GetPrivateLinkConnectionApiParams struct {
	GroupId      string
	ConnectionId string
}

func (a *StreamsApiService) GetPrivateLinkConnectionWithParams(ctx context.Context, args *GetPrivateLinkConnectionApiParams) GetPrivateLinkConnectionApiRequest {
	return GetPrivateLinkConnectionApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		connectionId: args.ConnectionId,
	}
}

func (r GetPrivateLinkConnectionApiRequest) Execute() (*StreamsPrivateLinkConnection, *http.Response, error) {
	return r.ApiService.GetPrivateLinkConnectionExecute(r)
}

/*
GetPrivateLinkConnection Return One Private Link Connection

Returns the details of one Private Link connection within the project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param connectionId Unique ID that identifies the Private Link connection.
	@return GetPrivateLinkConnectionApiRequest
*/
func (a *StreamsApiService) GetPrivateLinkConnection(ctx context.Context, groupId string, connectionId string) GetPrivateLinkConnectionApiRequest {
	return GetPrivateLinkConnectionApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      groupId,
		connectionId: connectionId,
	}
}

// GetPrivateLinkConnectionExecute executes the request
//
//	@return StreamsPrivateLinkConnection
func (a *StreamsApiService) GetPrivateLinkConnectionExecute(r GetPrivateLinkConnectionApiRequest) (*StreamsPrivateLinkConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsPrivateLinkConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetPrivateLinkConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/privateLinkConnections/{connectionId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.connectionId == "" {
		return localVarReturnValue, nil, reportError("connectionId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"connectionId"+"}", url.PathEscape(r.connectionId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type GetStreamConnectionApiRequest struct {
	ctx            context.Context
	ApiService     StreamsApi
	groupId        string
	tenantName     string
	connectionName string
}

type GetStreamConnectionApiParams struct {
	GroupId        string
	TenantName     string
	ConnectionName string
}

func (a *StreamsApiService) GetStreamConnectionWithParams(ctx context.Context, args *GetStreamConnectionApiParams) GetStreamConnectionApiRequest {
	return GetStreamConnectionApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		tenantName:     args.TenantName,
		connectionName: args.ConnectionName,
	}
}

func (r GetStreamConnectionApiRequest) Execute() (*StreamsConnection, *http.Response, error) {
	return r.ApiService.GetStreamConnectionExecute(r)
}

/*
GetStreamConnection Return One Stream Connection

Returns the details of one stream connection within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace to return.
	@param connectionName Label that identifies the stream connection to return.
	@return GetStreamConnectionApiRequest
*/
func (a *StreamsApiService) GetStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string) GetStreamConnectionApiRequest {
	return GetStreamConnectionApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		tenantName:     tenantName,
		connectionName: connectionName,
	}
}

// GetStreamConnectionExecute executes the request
//
//	@return StreamsConnection
func (a *StreamsApiService) GetStreamConnectionExecute(r GetStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetStreamConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.connectionName == "" {
		return localVarReturnValue, nil, reportError("connectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"connectionName"+"}", url.PathEscape(r.connectionName), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type GetStreamProcessorApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	tenantName    string
	processorName string
}

type GetStreamProcessorApiParams struct {
	GroupId       string
	TenantName    string
	ProcessorName string
}

func (a *StreamsApiService) GetStreamProcessorWithParams(ctx context.Context, args *GetStreamProcessorApiParams) GetStreamProcessorApiRequest {
	return GetStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		tenantName:    args.TenantName,
		processorName: args.ProcessorName,
	}
}

func (r GetStreamProcessorApiRequest) Execute() (*StreamsProcessorWithStats, *http.Response, error) {
	return r.ApiService.GetStreamProcessorExecute(r)
}

/*
GetStreamProcessor Return One Stream Processor

Get one Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return GetStreamProcessorApiRequest
*/
func (a *StreamsApiService) GetStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) GetStreamProcessorApiRequest {
	return GetStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		tenantName:    tenantName,
		processorName: processorName,
	}
}

// GetStreamProcessorExecute executes the request
//
//	@return StreamsProcessorWithStats
func (a *StreamsApiService) GetStreamProcessorExecute(r GetStreamProcessorApiRequest) (*StreamsProcessorWithStats, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsProcessorWithStats
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetStreamProcessor")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return localVarReturnValue, nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

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

type GetStreamProcessorsApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	tenantName   string
	itemsPerPage *int
	pageNum      *int
	includeCount *bool
}

type GetStreamProcessorsApiParams struct {
	GroupId      string
	TenantName   string
	ItemsPerPage *int
	PageNum      *int
	IncludeCount *bool
}

func (a *StreamsApiService) GetStreamProcessorsWithParams(ctx context.Context, args *GetStreamProcessorsApiParams) GetStreamProcessorsApiRequest {
	return GetStreamProcessorsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		tenantName:   args.TenantName,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		includeCount: args.IncludeCount,
	}
}

// Number of items that the response returns per page.
func (r GetStreamProcessorsApiRequest) ItemsPerPage(itemsPerPage int) GetStreamProcessorsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetStreamProcessorsApiRequest) PageNum(pageNum int) GetStreamProcessorsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r GetStreamProcessorsApiRequest) IncludeCount(includeCount bool) GetStreamProcessorsApiRequest {
	r.includeCount = &includeCount
	return r
}

func (r GetStreamProcessorsApiRequest) Execute() (*PaginatedApiStreamsStreamProcessorWithStats, *http.Response, error) {
	return r.ApiService.GetStreamProcessorsExecute(r)
}

/*
GetStreamProcessors Return All Stream Processors in One Stream Workspace

Returns all Stream Processors within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return GetStreamProcessorsApiRequest
*/
func (a *StreamsApiService) GetStreamProcessors(ctx context.Context, groupId string, tenantName string) GetStreamProcessorsApiRequest {
	return GetStreamProcessorsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// GetStreamProcessorsExecute executes the request
//
//	@return PaginatedApiStreamsStreamProcessorWithStats
func (a *StreamsApiService) GetStreamProcessorsExecute(r GetStreamProcessorsApiRequest) (*PaginatedApiStreamsStreamProcessorWithStats, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsStreamProcessorWithStats
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetStreamProcessors")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processors"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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
	if r.includeCount != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	} else {
		var defaultValue bool = true
		r.includeCount = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
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

type GetStreamWorkspaceApiRequest struct {
	ctx                context.Context
	ApiService         StreamsApi
	groupId            string
	tenantName         string
	includeConnections *bool
}

type GetStreamWorkspaceApiParams struct {
	GroupId            string
	TenantName         string
	IncludeConnections *bool
}

func (a *StreamsApiService) GetStreamWorkspaceWithParams(ctx context.Context, args *GetStreamWorkspaceApiParams) GetStreamWorkspaceApiRequest {
	return GetStreamWorkspaceApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		tenantName:         args.TenantName,
		includeConnections: args.IncludeConnections,
	}
}

// Flag to indicate whether connections information should be included in the stream workspace.
func (r GetStreamWorkspaceApiRequest) IncludeConnections(includeConnections bool) GetStreamWorkspaceApiRequest {
	r.includeConnections = &includeConnections
	return r
}

func (r GetStreamWorkspaceApiRequest) Execute() (*StreamsTenant, *http.Response, error) {
	return r.ApiService.GetStreamWorkspaceExecute(r)
}

/*
GetStreamWorkspace Return One Stream Workspace

Returns the details of one stream workspace within the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace to return.
	@return GetStreamWorkspaceApiRequest
*/
func (a *StreamsApiService) GetStreamWorkspace(ctx context.Context, groupId string, tenantName string) GetStreamWorkspaceApiRequest {
	return GetStreamWorkspaceApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// GetStreamWorkspaceExecute executes the request
//
//	@return StreamsTenant
func (a *StreamsApiService) GetStreamWorkspaceExecute(r GetStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.GetStreamWorkspace")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeConnections != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeConnections", r.includeConnections, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type ListActivePeeringConnectionsApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	itemsPerPage *int
	pageNum      *int
}

type ListActivePeeringConnectionsApiParams struct {
	GroupId      string
	ItemsPerPage *int
	PageNum      *int
}

func (a *StreamsApiService) ListActivePeeringConnectionsWithParams(ctx context.Context, args *ListActivePeeringConnectionsApiParams) ListActivePeeringConnectionsApiRequest {
	return ListActivePeeringConnectionsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListActivePeeringConnectionsApiRequest) ItemsPerPage(itemsPerPage int) ListActivePeeringConnectionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListActivePeeringConnectionsApiRequest) PageNum(pageNum int) ListActivePeeringConnectionsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListActivePeeringConnectionsApiRequest) Execute() (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error) {
	return r.ApiService.ListActivePeeringConnectionsExecute(r)
}

/*
ListActivePeeringConnections Return All Active Incoming VPC Peering Connections

Returns a list of active incoming VPC Peering Connections.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListActivePeeringConnectionsApiRequest
*/
func (a *StreamsApiService) ListActivePeeringConnections(ctx context.Context, groupId string) ListActivePeeringConnectionsApiRequest {
	return ListActivePeeringConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListActivePeeringConnectionsExecute executes the request
//
//	@return PaginatedApiStreamsVPCPeeringConnection
func (a *StreamsApiService) ListActivePeeringConnectionsExecute(r ListActivePeeringConnectionsApiRequest) (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsVPCPeeringConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.ListActivePeeringConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/activeVpcPeeringConnections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ListPrivateLinkConnectionsApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	itemsPerPage *int
	pageNum      *int
}

type ListPrivateLinkConnectionsApiParams struct {
	GroupId      string
	ItemsPerPage *int
	PageNum      *int
}

func (a *StreamsApiService) ListPrivateLinkConnectionsWithParams(ctx context.Context, args *ListPrivateLinkConnectionsApiParams) ListPrivateLinkConnectionsApiRequest {
	return ListPrivateLinkConnectionsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListPrivateLinkConnectionsApiRequest) ItemsPerPage(itemsPerPage int) ListPrivateLinkConnectionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListPrivateLinkConnectionsApiRequest) PageNum(pageNum int) ListPrivateLinkConnectionsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListPrivateLinkConnectionsApiRequest) Execute() (*PaginatedApiStreamsPrivateLink, *http.Response, error) {
	return r.ApiService.ListPrivateLinkConnectionsExecute(r)
}

/*
ListPrivateLinkConnections Return All Private Link Connections

Returns all Private Link connections for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListPrivateLinkConnectionsApiRequest
*/
func (a *StreamsApiService) ListPrivateLinkConnections(ctx context.Context, groupId string) ListPrivateLinkConnectionsApiRequest {
	return ListPrivateLinkConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListPrivateLinkConnectionsExecute executes the request
//
//	@return PaginatedApiStreamsPrivateLink
func (a *StreamsApiService) ListPrivateLinkConnectionsExecute(r ListPrivateLinkConnectionsApiRequest) (*PaginatedApiStreamsPrivateLink, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsPrivateLink
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.ListPrivateLinkConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/privateLinkConnections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type ListStreamConnectionsApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	tenantName   string
	itemsPerPage *int
	pageNum      *int
}

type ListStreamConnectionsApiParams struct {
	GroupId      string
	TenantName   string
	ItemsPerPage *int
	PageNum      *int
}

func (a *StreamsApiService) ListStreamConnectionsWithParams(ctx context.Context, args *ListStreamConnectionsApiParams) ListStreamConnectionsApiRequest {
	return ListStreamConnectionsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		tenantName:   args.TenantName,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListStreamConnectionsApiRequest) ItemsPerPage(itemsPerPage int) ListStreamConnectionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListStreamConnectionsApiRequest) PageNum(pageNum int) ListStreamConnectionsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListStreamConnectionsApiRequest) Execute() (*PaginatedApiStreamsConnection, *http.Response, error) {
	return r.ApiService.ListStreamConnectionsExecute(r)
}

/*
ListStreamConnections Return All Connections of the Stream Workspaces

Returns all connections of the stream workspace for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@return ListStreamConnectionsApiRequest
*/
func (a *StreamsApiService) ListStreamConnections(ctx context.Context, groupId string, tenantName string) ListStreamConnectionsApiRequest {
	return ListStreamConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		tenantName: tenantName,
	}
}

// ListStreamConnectionsExecute executes the request
//
//	@return PaginatedApiStreamsConnection
func (a *StreamsApiService) ListStreamConnectionsExecute(r ListStreamConnectionsApiRequest) (*PaginatedApiStreamsConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.ListStreamConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type ListStreamWorkspacesApiRequest struct {
	ctx          context.Context
	ApiService   StreamsApi
	groupId      string
	itemsPerPage *int
	pageNum      *int
}

type ListStreamWorkspacesApiParams struct {
	GroupId      string
	ItemsPerPage *int
	PageNum      *int
}

func (a *StreamsApiService) ListStreamWorkspacesWithParams(ctx context.Context, args *ListStreamWorkspacesApiParams) ListStreamWorkspacesApiRequest {
	return ListStreamWorkspacesApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r ListStreamWorkspacesApiRequest) ItemsPerPage(itemsPerPage int) ListStreamWorkspacesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListStreamWorkspacesApiRequest) PageNum(pageNum int) ListStreamWorkspacesApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListStreamWorkspacesApiRequest) Execute() (*PaginatedApiStreamsTenant, *http.Response, error) {
	return r.ApiService.ListStreamWorkspacesExecute(r)
}

/*
ListStreamWorkspaces Return All Stream Workspaces in One Project

Returns all stream workspaces for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListStreamWorkspacesApiRequest
*/
func (a *StreamsApiService) ListStreamWorkspaces(ctx context.Context, groupId string) ListStreamWorkspacesApiRequest {
	return ListStreamWorkspacesApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListStreamWorkspacesExecute executes the request
//
//	@return PaginatedApiStreamsTenant
func (a *StreamsApiService) ListStreamWorkspacesExecute(r ListStreamWorkspacesApiRequest) (*PaginatedApiStreamsTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.ListStreamWorkspaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type ListVpcPeeringConnectionsApiRequest struct {
	ctx                context.Context
	ApiService         StreamsApi
	requesterAccountId *string
	groupId            string
	itemsPerPage       *int
	pageNum            *int
}

type ListVpcPeeringConnectionsApiParams struct {
	RequesterAccountId *string
	GroupId            string
	ItemsPerPage       *int
	PageNum            *int
}

func (a *StreamsApiService) ListVpcPeeringConnectionsWithParams(ctx context.Context, args *ListVpcPeeringConnectionsApiParams) ListVpcPeeringConnectionsApiRequest {
	return ListVpcPeeringConnectionsApiRequest{
		ApiService:         a,
		ctx:                ctx,
		requesterAccountId: args.RequesterAccountId,
		groupId:            args.GroupId,
		itemsPerPage:       args.ItemsPerPage,
		pageNum:            args.PageNum,
	}
}

// The Account ID of the VPC Peering connection/s.
func (r ListVpcPeeringConnectionsApiRequest) RequesterAccountId(requesterAccountId string) ListVpcPeeringConnectionsApiRequest {
	r.requesterAccountId = &requesterAccountId
	return r
}

// Number of items that the response returns per page.
func (r ListVpcPeeringConnectionsApiRequest) ItemsPerPage(itemsPerPage int) ListVpcPeeringConnectionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListVpcPeeringConnectionsApiRequest) PageNum(pageNum int) ListVpcPeeringConnectionsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListVpcPeeringConnectionsApiRequest) Execute() (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error) {
	return r.ApiService.ListVpcPeeringConnectionsExecute(r)
}

/*
ListVpcPeeringConnections Return All VPC Peering Connections

Returns a list of incoming VPC Peering Connections.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListVpcPeeringConnectionsApiRequest
*/
func (a *StreamsApiService) ListVpcPeeringConnections(ctx context.Context, groupId string) ListVpcPeeringConnectionsApiRequest {
	return ListVpcPeeringConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListVpcPeeringConnectionsExecute executes the request
//
//	@return PaginatedApiStreamsVPCPeeringConnection
func (a *StreamsApiService) ListVpcPeeringConnectionsExecute(r ListVpcPeeringConnectionsApiRequest) (*PaginatedApiStreamsVPCPeeringConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiStreamsVPCPeeringConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.ListVpcPeeringConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/vpcPeeringConnections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.requesterAccountId == nil {
		return localVarReturnValue, nil, reportError("requesterAccountId is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "requesterAccountId", r.requesterAccountId, "")
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type RejectVpcPeeringConnectionApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	id         string
}

type RejectVpcPeeringConnectionApiParams struct {
	GroupId string
	Id      string
}

func (a *StreamsApiService) RejectVpcPeeringConnectionWithParams(ctx context.Context, args *RejectVpcPeeringConnectionApiParams) RejectVpcPeeringConnectionApiRequest {
	return RejectVpcPeeringConnectionApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		id:         args.Id,
	}
}

func (r RejectVpcPeeringConnectionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RejectVpcPeeringConnectionExecute(r)
}

/*
RejectVpcPeeringConnection Reject One Incoming VPC Peering Connection

Requests the rejection of an incoming VPC Peering connection.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param id The VPC Peering Connection id.
	@return RejectVpcPeeringConnectionApiRequest
*/
func (a *StreamsApiService) RejectVpcPeeringConnection(ctx context.Context, groupId string, id string) RejectVpcPeeringConnectionApiRequest {
	return RejectVpcPeeringConnectionApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		id:         id,
	}
}

// RejectVpcPeeringConnectionExecute executes the request
func (a *StreamsApiService) RejectVpcPeeringConnectionExecute(r RejectVpcPeeringConnectionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.RejectVpcPeeringConnection")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/vpcPeeringConnections/{id}:reject"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.id == "" {
		return nil, reportError("id is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(r.id), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

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

type StartStreamProcessorApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	tenantName    string
	processorName string
}

type StartStreamProcessorApiParams struct {
	GroupId       string
	TenantName    string
	ProcessorName string
}

func (a *StreamsApiService) StartStreamProcessorWithParams(ctx context.Context, args *StartStreamProcessorApiParams) StartStreamProcessorApiRequest {
	return StartStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		tenantName:    args.TenantName,
		processorName: args.ProcessorName,
	}
}

func (r StartStreamProcessorApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.StartStreamProcessorExecute(r)
}

/*
StartStreamProcessor Start One Stream Processor

Start a Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return StartStreamProcessorApiRequest
*/
func (a *StreamsApiService) StartStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) StartStreamProcessorApiRequest {
	return StartStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		tenantName:    tenantName,
		processorName: processorName,
	}
}

// StartStreamProcessorExecute executes the request
func (a *StreamsApiService) StartStreamProcessorExecute(r StartStreamProcessorApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.StartStreamProcessor")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}:start"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

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

type StartStreamProcessorWithApiRequest struct {
	ctx                             context.Context
	ApiService                      StreamsApi
	groupId                         string
	tenantName                      string
	processorName                   string
	streamsStartStreamProcessorWith *StreamsStartStreamProcessorWith
}

type StartStreamProcessorWithApiParams struct {
	GroupId                         string
	TenantName                      string
	ProcessorName                   string
	StreamsStartStreamProcessorWith *StreamsStartStreamProcessorWith
}

func (a *StreamsApiService) StartStreamProcessorWithWithParams(ctx context.Context, args *StartStreamProcessorWithApiParams) StartStreamProcessorWithApiRequest {
	return StartStreamProcessorWithApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         args.GroupId,
		tenantName:                      args.TenantName,
		processorName:                   args.ProcessorName,
		streamsStartStreamProcessorWith: args.StreamsStartStreamProcessorWith,
	}
}

func (r StartStreamProcessorWithApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.StartStreamProcessorWithExecute(r)
}

/*
StartStreamProcessorWith Start One Stream Processor With Options

Start a Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return StartStreamProcessorWithApiRequest
*/
func (a *StreamsApiService) StartStreamProcessorWith(ctx context.Context, groupId string, tenantName string, processorName string, streamsStartStreamProcessorWith *StreamsStartStreamProcessorWith) StartStreamProcessorWithApiRequest {
	return StartStreamProcessorWithApiRequest{
		ApiService:                      a,
		ctx:                             ctx,
		groupId:                         groupId,
		tenantName:                      tenantName,
		processorName:                   processorName,
		streamsStartStreamProcessorWith: streamsStartStreamProcessorWith,
	}
}

// StartStreamProcessorWithExecute executes the request
func (a *StreamsApiService) StartStreamProcessorWithExecute(r StartStreamProcessorWithApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.StartStreamProcessorWith")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}:startWith"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2025-03-12+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsStartStreamProcessorWith
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

type StopStreamProcessorApiRequest struct {
	ctx           context.Context
	ApiService    StreamsApi
	groupId       string
	tenantName    string
	processorName string
}

type StopStreamProcessorApiParams struct {
	GroupId       string
	TenantName    string
	ProcessorName string
}

func (a *StreamsApiService) StopStreamProcessorWithParams(ctx context.Context, args *StopStreamProcessorApiParams) StopStreamProcessorApiRequest {
	return StopStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		tenantName:    args.TenantName,
		processorName: args.ProcessorName,
	}
}

func (r StopStreamProcessorApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.StopStreamProcessorExecute(r)
}

/*
StopStreamProcessor Stop One Stream Processor

Stop a Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return StopStreamProcessorApiRequest
*/
func (a *StreamsApiService) StopStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string) StopStreamProcessorApiRequest {
	return StopStreamProcessorApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		tenantName:    tenantName,
		processorName: processorName,
	}
}

// StopStreamProcessorExecute executes the request
func (a *StreamsApiService) StopStreamProcessorExecute(r StopStreamProcessorApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.StopStreamProcessor")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}:stop"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

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

type UpdateStreamConnectionApiRequest struct {
	ctx               context.Context
	ApiService        StreamsApi
	groupId           string
	tenantName        string
	connectionName    string
	streamsConnection *StreamsConnection
}

type UpdateStreamConnectionApiParams struct {
	GroupId           string
	TenantName        string
	ConnectionName    string
	StreamsConnection *StreamsConnection
}

func (a *StreamsApiService) UpdateStreamConnectionWithParams(ctx context.Context, args *UpdateStreamConnectionApiParams) UpdateStreamConnectionApiRequest {
	return UpdateStreamConnectionApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           args.GroupId,
		tenantName:        args.TenantName,
		connectionName:    args.ConnectionName,
		streamsConnection: args.StreamsConnection,
	}
}

func (r UpdateStreamConnectionApiRequest) Execute() (*StreamsConnection, *http.Response, error) {
	return r.ApiService.UpdateStreamConnectionExecute(r)
}

/*
UpdateStreamConnection Update One Stream Connection

Update one connection for the specified stream workspace in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param connectionName Label that identifies the stream connection.
	@return UpdateStreamConnectionApiRequest
*/
func (a *StreamsApiService) UpdateStreamConnection(ctx context.Context, groupId string, tenantName string, connectionName string, streamsConnection *StreamsConnection) UpdateStreamConnectionApiRequest {
	return UpdateStreamConnectionApiRequest{
		ApiService:        a,
		ctx:               ctx,
		groupId:           groupId,
		tenantName:        tenantName,
		connectionName:    connectionName,
		streamsConnection: streamsConnection,
	}
}

// UpdateStreamConnectionExecute executes the request
//
//	@return StreamsConnection
func (a *StreamsApiService) UpdateStreamConnectionExecute(r UpdateStreamConnectionApiRequest) (*StreamsConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.UpdateStreamConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.connectionName == "" {
		return localVarReturnValue, nil, reportError("connectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"connectionName"+"}", url.PathEscape(r.connectionName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsConnection == nil {
		return localVarReturnValue, nil, reportError("streamsConnection is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsConnection
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

type UpdateStreamProcessorApiRequest struct {
	ctx                          context.Context
	ApiService                   StreamsApi
	groupId                      string
	tenantName                   string
	processorName                string
	streamsModifyStreamProcessor *StreamsModifyStreamProcessor
}

type UpdateStreamProcessorApiParams struct {
	GroupId                      string
	TenantName                   string
	ProcessorName                string
	StreamsModifyStreamProcessor *StreamsModifyStreamProcessor
}

func (a *StreamsApiService) UpdateStreamProcessorWithParams(ctx context.Context, args *UpdateStreamProcessorApiParams) UpdateStreamProcessorApiRequest {
	return UpdateStreamProcessorApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      args.GroupId,
		tenantName:                   args.TenantName,
		processorName:                args.ProcessorName,
		streamsModifyStreamProcessor: args.StreamsModifyStreamProcessor,
	}
}

func (r UpdateStreamProcessorApiRequest) Execute() (*StreamsProcessorWithStats, *http.Response, error) {
	return r.ApiService.UpdateStreamProcessorExecute(r)
}

/*
UpdateStreamProcessor Update One Stream Processor

Modify one existing Stream Processor within the specified stream workspace.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace.
	@param processorName Label that identifies the stream processor.
	@return UpdateStreamProcessorApiRequest
*/
func (a *StreamsApiService) UpdateStreamProcessor(ctx context.Context, groupId string, tenantName string, processorName string, streamsModifyStreamProcessor *StreamsModifyStreamProcessor) UpdateStreamProcessorApiRequest {
	return UpdateStreamProcessorApiRequest{
		ApiService:                   a,
		ctx:                          ctx,
		groupId:                      groupId,
		tenantName:                   tenantName,
		processorName:                processorName,
		streamsModifyStreamProcessor: streamsModifyStreamProcessor,
	}
}

// UpdateStreamProcessorExecute executes the request
//
//	@return StreamsProcessorWithStats
func (a *StreamsApiService) UpdateStreamProcessorExecute(r UpdateStreamProcessorApiRequest) (*StreamsProcessorWithStats, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsProcessorWithStats
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.UpdateStreamProcessor")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/processor/{processorName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)
	if r.processorName == "" {
		return localVarReturnValue, nil, reportError("processorName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"processorName"+"}", url.PathEscape(r.processorName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsModifyStreamProcessor == nil {
		return localVarReturnValue, nil, reportError("streamsModifyStreamProcessor is required and must be specified")
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
	localVarPostBody = r.streamsModifyStreamProcessor
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

type UpdateStreamWorkspaceApiRequest struct {
	ctx                        context.Context
	ApiService                 StreamsApi
	groupId                    string
	tenantName                 string
	streamsTenantUpdateRequest *StreamsTenantUpdateRequest
}

type UpdateStreamWorkspaceApiParams struct {
	GroupId                    string
	TenantName                 string
	StreamsTenantUpdateRequest *StreamsTenantUpdateRequest
}

func (a *StreamsApiService) UpdateStreamWorkspaceWithParams(ctx context.Context, args *UpdateStreamWorkspaceApiParams) UpdateStreamWorkspaceApiRequest {
	return UpdateStreamWorkspaceApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		tenantName:                 args.TenantName,
		streamsTenantUpdateRequest: args.StreamsTenantUpdateRequest,
	}
}

func (r UpdateStreamWorkspaceApiRequest) Execute() (*StreamsTenant, *http.Response, error) {
	return r.ApiService.UpdateStreamWorkspaceExecute(r)
}

/*
UpdateStreamWorkspace Update One Stream Workspace

Update one stream workspace in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param tenantName Label that identifies the stream workspace to update.
	@return UpdateStreamWorkspaceApiRequest
*/
func (a *StreamsApiService) UpdateStreamWorkspace(ctx context.Context, groupId string, tenantName string, streamsTenantUpdateRequest *StreamsTenantUpdateRequest) UpdateStreamWorkspaceApiRequest {
	return UpdateStreamWorkspaceApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		tenantName:                 tenantName,
		streamsTenantUpdateRequest: streamsTenantUpdateRequest,
	}
}

// UpdateStreamWorkspaceExecute executes the request
//
//	@return StreamsTenant
func (a *StreamsApiService) UpdateStreamWorkspaceExecute(r UpdateStreamWorkspaceApiRequest) (*StreamsTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.UpdateStreamWorkspace")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams/{tenantName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.tenantName == "" {
		return localVarReturnValue, nil, reportError("tenantName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"tenantName"+"}", url.PathEscape(r.tenantName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.streamsTenantUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("streamsTenantUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-02-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.streamsTenantUpdateRequest
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

type WithStreamSampleConnectionsApiRequest struct {
	ctx        context.Context
	ApiService StreamsApi
	groupId    string
	body       *any
}

type WithStreamSampleConnectionsApiParams struct {
	GroupId string
	Body    *any
}

func (a *StreamsApiService) WithStreamSampleConnectionsWithParams(ctx context.Context, args *WithStreamSampleConnectionsApiParams) WithStreamSampleConnectionsApiRequest {
	return WithStreamSampleConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		body:       args.Body,
	}
}

func (r WithStreamSampleConnectionsApiRequest) Execute() (*StreamsTenant, *http.Response, error) {
	return r.ApiService.WithStreamSampleConnectionsExecute(r)
}

/*
WithStreamSampleConnections Create One Stream Workspace with Sample Connections

Creates one stream workspace in the specified project with sample connections.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return WithStreamSampleConnectionsApiRequest
*/
func (a *StreamsApiService) WithStreamSampleConnections(ctx context.Context, groupId string, body *any) WithStreamSampleConnectionsApiRequest {
	return WithStreamSampleConnectionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		body:       body,
	}
}

// WithStreamSampleConnectionsExecute executes the request
//
//	@return StreamsTenant
func (a *StreamsApiService) WithStreamSampleConnectionsExecute(r WithStreamSampleConnectionsApiRequest) (*StreamsTenant, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *StreamsTenant
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StreamsApiService.WithStreamSampleConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/streams:withSampleConnections"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
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
	localVarPostBody = r.body
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
