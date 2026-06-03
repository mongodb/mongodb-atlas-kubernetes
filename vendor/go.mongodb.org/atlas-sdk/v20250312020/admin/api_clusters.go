// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type ClustersApi interface {

	/*
		AutoScalingConfiguration Return Auto Scaling Configuration for One Sharded Cluster

		Returns the internal configuration of AutoScaling for sharded clusters. This endpoint can be used for diagnostic purposes to ensure that sharded clusters updated from older APIs have gained support for AutoScaling each shard independently.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies this cluster.
		@return AutoScalingConfigurationApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ClustersApi
	*/
	AutoScalingConfiguration(ctx context.Context, groupId string, clusterName string) AutoScalingConfigurationApiRequest
	/*
		AutoScalingConfiguration Return Auto Scaling Configuration for One Sharded Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param AutoScalingConfigurationApiParams - Parameters for the request
		@return AutoScalingConfigurationApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ClustersApi
	*/
	AutoScalingConfigurationWithParams(ctx context.Context, args *AutoScalingConfigurationApiParams) AutoScalingConfigurationApiRequest

	// Method available only for mocking purposes
	AutoScalingConfigurationExecute(r AutoScalingConfigurationApiRequest) (*ClusterDescriptionAutoScalingModeConfiguration, *http.Response, error)

	/*
			CreateCluster Create One Cluster in One Project

			Creates one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. This resource can create clusters with asymmetrically-sized shards. Each project supports up to 25 database deployments. This feature is not available for serverless clusters.

		Please note that using an `instanceSize` of M2 or M5 will create a Flex cluster instead. Support for the `instanceSize` of M2 or M5 will be discontinued in January 2026. We recommend using the Create Flex Cluster API for such configurations moving forward. Deprecated versions: v2-{2024-08-05}, v2-{2023-02-01}, v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterDescription20240805 Cluster to create in this project.
			@return CreateClusterApiRequest
	*/
	CreateCluster(ctx context.Context, groupId string, clusterDescription20240805 *ClusterDescription20240805) CreateClusterApiRequest
	/*
		CreateCluster Create One Cluster in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateClusterApiParams - Parameters for the request
		@return CreateClusterApiRequest
	*/
	CreateClusterWithParams(ctx context.Context, args *CreateClusterApiParams) CreateClusterApiRequest

	// Method available only for mocking purposes
	CreateClusterExecute(r CreateClusterApiRequest) (*ClusterDescription20240805, *http.Response, error)

	/*
			DeleteCluster Remove One Cluster from One Project

			Removes one cluster from the specified project. The cluster must have termination protection disabled in order to be deleted. This feature is not available for serverless clusters.

		This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the Delete Flex Cluster endpoint for Flex clusters instead. Deprecated versions: v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies the cluster.
			@return DeleteClusterApiRequest
	*/
	DeleteCluster(ctx context.Context, groupId string, clusterName string) DeleteClusterApiRequest
	/*
		DeleteCluster Remove One Cluster from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterApiParams - Parameters for the request
		@return DeleteClusterApiRequest
	*/
	DeleteClusterWithParams(ctx context.Context, args *DeleteClusterApiParams) DeleteClusterApiRequest

	// Method available only for mocking purposes
	DeleteClusterExecute(r DeleteClusterApiRequest) (*http.Response, error)

	/*
			GetCluster Return One Cluster from One Project

			Returns the details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. The response includes clusters with asymmetrically-sized shards. This feature is not available for serverless clusters.

		This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the Get Flex Cluster endpoint for Flex clusters instead. Deprecated versions: v2-{2023-02-01}, v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param clusterName Human-readable label that identifies this cluster.
			@return GetClusterApiRequest
	*/
	GetCluster(ctx context.Context, groupId string, clusterName string) GetClusterApiRequest
	/*
		GetCluster Return One Cluster from One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterApiParams - Parameters for the request
		@return GetClusterApiRequest
	*/
	GetClusterWithParams(ctx context.Context, args *GetClusterApiParams) GetClusterApiRequest

	// Method available only for mocking purposes
	GetClusterExecute(r GetClusterApiRequest) (*ClusterDescription20240805, *http.Response, error)

	/*
		GetClusterStatus Return Status of All Cluster Operations

		Returns the status of all changes that you made to the specified cluster in the specified project. Use this resource to check the progress MongoDB Cloud has made in processing your changes. The response does not include the deployment of new dedicated clusters.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return GetClusterStatusApiRequest
	*/
	GetClusterStatus(ctx context.Context, groupId string, clusterName string) GetClusterStatusApiRequest
	/*
		GetClusterStatus Return Status of All Cluster Operations


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterStatusApiParams - Parameters for the request
		@return GetClusterStatusApiRequest
	*/
	GetClusterStatusWithParams(ctx context.Context, args *GetClusterStatusApiParams) GetClusterStatusApiRequest

	// Method available only for mocking purposes
	GetClusterStatusExecute(r GetClusterStatusApiRequest) (*ClusterStatus, *http.Response, error)

	/*
		GetProcessArgs Return Advanced Configuration Options for One Cluster

		Returns the advanced configuration details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. Advanced configuration details include the read/write concern, index and oplog limits, and other database settings. This feature isn't available for `M0` free clusters, `M2` and `M5` shared-tier clusters, flex clusters, or serverless clusters. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return GetProcessArgsApiRequest
	*/
	GetProcessArgs(ctx context.Context, groupId string, clusterName string) GetProcessArgsApiRequest
	/*
		GetProcessArgs Return Advanced Configuration Options for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetProcessArgsApiParams - Parameters for the request
		@return GetProcessArgsApiRequest
	*/
	GetProcessArgsWithParams(ctx context.Context, args *GetProcessArgsApiParams) GetProcessArgsApiRequest

	// Method available only for mocking purposes
	GetProcessArgsExecute(r GetProcessArgsApiRequest) (*ClusterDescriptionProcessArgs20240805, *http.Response, error)

	/*
		GetSampleDatasetLoad Return Status of Sample Dataset Load for One Cluster

		Checks the progress of loading the sample dataset into one cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param sampleDatasetId Unique 24-hexadecimal digit string that identifies the loaded sample dataset.
		@return GetSampleDatasetLoadApiRequest
	*/
	GetSampleDatasetLoad(ctx context.Context, groupId string, sampleDatasetId string) GetSampleDatasetLoadApiRequest
	/*
		GetSampleDatasetLoad Return Status of Sample Dataset Load for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetSampleDatasetLoadApiParams - Parameters for the request
		@return GetSampleDatasetLoadApiRequest
	*/
	GetSampleDatasetLoadWithParams(ctx context.Context, args *GetSampleDatasetLoadApiParams) GetSampleDatasetLoadApiRequest

	// Method available only for mocking purposes
	GetSampleDatasetLoadExecute(r GetSampleDatasetLoadApiRequest) (*SampleDatasetStatus, *http.Response, error)

	/*
		GrantMongoEmployeeAccess Grant MongoDB Employee Cluster Access for One Cluster

		Grants MongoDB employee cluster access for the given duration and at the specified level for one cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies this cluster.
		@param employeeAccessGrant Grant access level and expiration.
		@return GrantMongoEmployeeAccessApiRequest
	*/
	GrantMongoEmployeeAccess(ctx context.Context, groupId string, clusterName string, employeeAccessGrant *EmployeeAccessGrant) GrantMongoEmployeeAccessApiRequest
	/*
		GrantMongoEmployeeAccess Grant MongoDB Employee Cluster Access for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GrantMongoEmployeeAccessApiParams - Parameters for the request
		@return GrantMongoEmployeeAccessApiRequest
	*/
	GrantMongoEmployeeAccessWithParams(ctx context.Context, args *GrantMongoEmployeeAccessApiParams) GrantMongoEmployeeAccessApiRequest

	// Method available only for mocking purposes
	GrantMongoEmployeeAccessExecute(r GrantMongoEmployeeAccessApiRequest) (*http.Response, error)

	/*
		ListClusterDetails Return All Authorized Clusters in All Projects

		Returns the details for all clusters in all projects to which you have access. Clusters contain a group of hosts that maintain the same data set. The response does not include multi-cloud clusters.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListClusterDetailsApiRequest
	*/
	ListClusterDetails(ctx context.Context) ListClusterDetailsApiRequest
	/*
		ListClusterDetails Return All Authorized Clusters in All Projects


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterDetailsApiParams - Parameters for the request
		@return ListClusterDetailsApiRequest
	*/
	ListClusterDetailsWithParams(ctx context.Context, args *ListClusterDetailsApiParams) ListClusterDetailsApiRequest

	// Method available only for mocking purposes
	ListClusterDetailsExecute(r ListClusterDetailsApiRequest) (*PaginatedOrgGroup, *http.Response, error)

	/*
		ListClusterProviderRegions Return All Cloud Provider Regions

		Returns the list of regions available for the specified cloud provider at the specified tier.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListClusterProviderRegionsApiRequest
	*/
	ListClusterProviderRegions(ctx context.Context, groupId string) ListClusterProviderRegionsApiRequest
	/*
		ListClusterProviderRegions Return All Cloud Provider Regions


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterProviderRegionsApiParams - Parameters for the request
		@return ListClusterProviderRegionsApiRequest
	*/
	ListClusterProviderRegionsWithParams(ctx context.Context, args *ListClusterProviderRegionsApiParams) ListClusterProviderRegionsApiRequest

	// Method available only for mocking purposes
	ListClusterProviderRegionsExecute(r ListClusterProviderRegionsApiRequest) (*PaginatedApiAtlasProviderRegions, *http.Response, error)

	/*
			ListClusters Return All Clusters in One Project

			Returns the details for all clusters in the specific project to which you have access. Clusters contain a group of hosts that maintain the same data set. The response includes clusters with asymmetrically-sized shards. This feature is not  available for serverless clusters.

		This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the List Flex Clusters endpoint for Flex clusters instead. Deprecated versions: v2-{2023-02-01}, v2-{2023-01-01}

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@return ListClustersApiRequest
	*/
	ListClusters(ctx context.Context, groupId string) ListClustersApiRequest
	/*
		ListClusters Return All Clusters in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClustersApiParams - Parameters for the request
		@return ListClustersApiRequest
	*/
	ListClustersWithParams(ctx context.Context, args *ListClustersApiParams) ListClustersApiRequest

	// Method available only for mocking purposes
	ListClustersExecute(r ListClustersApiRequest) (*PaginatedClusterDescription20240805, *http.Response, error)

	/*
		PinFeatureCompatibilityVersion Pin Feature Compatibility Version for One Cluster in One Project

		Pins the Feature Compatibility Version (FCV) to the current MongoDB version and sets the pin expiration date. If an FCV pin already exists for the cluster, calling this method will only update the expiration date of the existing pin and will not re-pin the FCV.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies this cluster.
		@param pinFCV Optional request parameters for tuning FCV pinning configuration.
		@return PinFeatureCompatibilityVersionApiRequest
	*/
	PinFeatureCompatibilityVersion(ctx context.Context, groupId string, clusterName string, pinFCV *PinFCV) PinFeatureCompatibilityVersionApiRequest
	/*
		PinFeatureCompatibilityVersion Pin Feature Compatibility Version for One Cluster in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param PinFeatureCompatibilityVersionApiParams - Parameters for the request
		@return PinFeatureCompatibilityVersionApiRequest
	*/
	PinFeatureCompatibilityVersionWithParams(ctx context.Context, args *PinFeatureCompatibilityVersionApiParams) PinFeatureCompatibilityVersionApiRequest

	// Method available only for mocking purposes
	PinFeatureCompatibilityVersionExecute(r PinFeatureCompatibilityVersionApiRequest) (*http.Response, error)

	/*
		RequestSampleDatasetLoad Load Sample Dataset into One Cluster

		Requests loading the MongoDB sample dataset into the specified cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param name Human-readable label that identifies the cluster into which you load the sample dataset.
		@return RequestSampleDatasetLoadApiRequest
	*/
	RequestSampleDatasetLoad(ctx context.Context, groupId string, name string) RequestSampleDatasetLoadApiRequest
	/*
		RequestSampleDatasetLoad Load Sample Dataset into One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RequestSampleDatasetLoadApiParams - Parameters for the request
		@return RequestSampleDatasetLoadApiRequest
	*/
	RequestSampleDatasetLoadWithParams(ctx context.Context, args *RequestSampleDatasetLoadApiParams) RequestSampleDatasetLoadApiRequest

	// Method available only for mocking purposes
	RequestSampleDatasetLoadExecute(r RequestSampleDatasetLoadApiRequest) (*SampleDatasetStatus, *http.Response, error)

	/*
		RestartPrimaries Test Failover for One Cluster

		Starts a failover test for the specified cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. A failover test checks how MongoDB Cloud handles the failure of the cluster's primary node. During the test, MongoDB Cloud shuts down the primary node and elects a new primary. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@return RestartPrimariesApiRequest
	*/
	RestartPrimaries(ctx context.Context, groupId string, clusterName string) RestartPrimariesApiRequest
	/*
		RestartPrimaries Test Failover for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RestartPrimariesApiParams - Parameters for the request
		@return RestartPrimariesApiRequest
	*/
	RestartPrimariesWithParams(ctx context.Context, args *RestartPrimariesApiParams) RestartPrimariesApiRequest

	// Method available only for mocking purposes
	RestartPrimariesExecute(r RestartPrimariesApiRequest) (*http.Response, error)

	/*
		RevokeMongoEmployeeAccess Revoke MongoDB Employee Cluster Access for One Cluster

		Revokes a previously granted MongoDB employee cluster access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies this cluster.
		@return RevokeMongoEmployeeAccessApiRequest
	*/
	RevokeMongoEmployeeAccess(ctx context.Context, groupId string, clusterName string) RevokeMongoEmployeeAccessApiRequest
	/*
		RevokeMongoEmployeeAccess Revoke MongoDB Employee Cluster Access for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param RevokeMongoEmployeeAccessApiParams - Parameters for the request
		@return RevokeMongoEmployeeAccessApiRequest
	*/
	RevokeMongoEmployeeAccessWithParams(ctx context.Context, args *RevokeMongoEmployeeAccessApiParams) RevokeMongoEmployeeAccessApiRequest

	// Method available only for mocking purposes
	RevokeMongoEmployeeAccessExecute(r RevokeMongoEmployeeAccessApiRequest) (*http.Response, error)

	/*
		UnpinFeatureCompatibilityVersion Unpin Feature Compatibility Version for One Cluster in One Project

		Unpins the current fixed Feature Compatibility Version (FCV). This feature is not available for clusters on rapid release.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies this cluster.
		@return UnpinFeatureCompatibilityVersionApiRequest
	*/
	UnpinFeatureCompatibilityVersion(ctx context.Context, groupId string, clusterName string) UnpinFeatureCompatibilityVersionApiRequest
	/*
		UnpinFeatureCompatibilityVersion Unpin Feature Compatibility Version for One Cluster in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UnpinFeatureCompatibilityVersionApiParams - Parameters for the request
		@return UnpinFeatureCompatibilityVersionApiRequest
	*/
	UnpinFeatureCompatibilityVersionWithParams(ctx context.Context, args *UnpinFeatureCompatibilityVersionApiParams) UnpinFeatureCompatibilityVersionApiRequest

	// Method available only for mocking purposes
	UnpinFeatureCompatibilityVersionExecute(r UnpinFeatureCompatibilityVersionApiRequest) (*http.Response, error)

	/*
		UpdateCluster Update One Cluster in One Project

		Updates the details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. This resource can update clusters with asymmetrically-sized shards. To update a cluster's termination protection, the requesting Service Account or API Key must have the Project Owner role. For all other updates, the requesting Service Account or API Key must have the Project Cluster Manager role or the Project Replica Set Manager role. You can't modify a paused cluster (`paused : true`). You must call this endpoint to set `paused : false`. After this endpoint responds with `paused : false`, you can call it again with the changes you want to make to the cluster. This feature is not available for serverless clusters. Deprecated versions: v2-{2024-08-05}, v2-{2023-02-01}, v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param clusterDescription20240805 Cluster to update in the specified project.
		@return UpdateClusterApiRequest
	*/
	UpdateCluster(ctx context.Context, groupId string, clusterName string, clusterDescription20240805 *ClusterDescription20240805) UpdateClusterApiRequest
	/*
		UpdateCluster Update One Cluster in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterApiParams - Parameters for the request
		@return UpdateClusterApiRequest
	*/
	UpdateClusterWithParams(ctx context.Context, args *UpdateClusterApiParams) UpdateClusterApiRequest

	// Method available only for mocking purposes
	UpdateClusterExecute(r UpdateClusterApiRequest) (*ClusterDescription20240805, *http.Response, error)

	/*
		UpdateProcessArgs Update Advanced Configuration Options for One Cluster

		Updates the advanced configuration details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. Advanced configuration details include the read/write concern, index and oplog limits, and other database settings. This feature isn't available for `M0` free clusters, `M2` and `M5` shared-tier clusters, flex clusters, or serverless clusters. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Human-readable label that identifies the cluster.
		@param clusterDescriptionProcessArgs20240805 Advanced configuration details to add for one cluster in the specified project.
		@return UpdateProcessArgsApiRequest
	*/
	UpdateProcessArgs(ctx context.Context, groupId string, clusterName string, clusterDescriptionProcessArgs20240805 *ClusterDescriptionProcessArgs20240805) UpdateProcessArgsApiRequest
	/*
		UpdateProcessArgs Update Advanced Configuration Options for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateProcessArgsApiParams - Parameters for the request
		@return UpdateProcessArgsApiRequest
	*/
	UpdateProcessArgsWithParams(ctx context.Context, args *UpdateProcessArgsApiParams) UpdateProcessArgsApiRequest

	// Method available only for mocking purposes
	UpdateProcessArgsExecute(r UpdateProcessArgsApiRequest) (*ClusterDescriptionProcessArgs20240805, *http.Response, error)

	/*
			UpgradeClusterToServerless Upgrade One Shared-Tier Cluster to One Serverless Instance

			This endpoint has been deprecated as of February 2025 as we no longer support the creation of new serverless instances. Please use the Upgrade Flex Cluster endpoint to upgrade Flex clusters.

		 Upgrades a shared-tier cluster to a serverless instance in the specified project.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param serverlessInstanceDescription Details of the shared-tier cluster upgrade in the specified project.
			@return UpgradeClusterToServerlessApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for ClustersApi
	*/
	UpgradeClusterToServerless(ctx context.Context, groupId string, serverlessInstanceDescription *ServerlessInstanceDescription) UpgradeClusterToServerlessApiRequest
	/*
		UpgradeClusterToServerless Upgrade One Shared-Tier Cluster to One Serverless Instance


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpgradeClusterToServerlessApiParams - Parameters for the request
		@return UpgradeClusterToServerlessApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for ClustersApi
	*/
	UpgradeClusterToServerlessWithParams(ctx context.Context, args *UpgradeClusterToServerlessApiParams) UpgradeClusterToServerlessApiRequest

	// Method available only for mocking purposes
	UpgradeClusterToServerlessExecute(r UpgradeClusterToServerlessApiRequest) (*ServerlessInstanceDescription, *http.Response, error)

	/*
			UpgradeTenantUpgrade Upgrade One Shared-Tier Cluster

			Upgrades a shared-tier cluster to a Flex or Dedicated (M10+) cluster in the specified project. Each project supports up to 25 clusters.

		This endpoint can also be used to upgrade Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) API or former M2/M5 clusters that have been migrated to Flex clusters, using `instanceSizeName` to “M2” or “M5” until January 2026. This functionality will be available until January 22, 2026, after which it will only be available for M0 clusters. Please use the Upgrade Flex Cluster endpoint instead.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
			@param legacyAtlasTenantClusterUpgradeRequest Details of the shared-tier cluster upgrade in the specified project.
			@return UpgradeTenantUpgradeApiRequest
	*/
	UpgradeTenantUpgrade(ctx context.Context, groupId string, legacyAtlasTenantClusterUpgradeRequest *LegacyAtlasTenantClusterUpgradeRequest) UpgradeTenantUpgradeApiRequest
	/*
		UpgradeTenantUpgrade Upgrade One Shared-Tier Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpgradeTenantUpgradeApiParams - Parameters for the request
		@return UpgradeTenantUpgradeApiRequest
	*/
	UpgradeTenantUpgradeWithParams(ctx context.Context, args *UpgradeTenantUpgradeApiParams) UpgradeTenantUpgradeApiRequest

	// Method available only for mocking purposes
	UpgradeTenantUpgradeExecute(r UpgradeTenantUpgradeApiRequest) (*LegacyAtlasCluster, *http.Response, error)
}

// ClustersApiService ClustersApi service
type ClustersApiService service

type AutoScalingConfigurationApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type AutoScalingConfigurationApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) AutoScalingConfigurationWithParams(ctx context.Context, args *AutoScalingConfigurationApiParams) AutoScalingConfigurationApiRequest {
	return AutoScalingConfigurationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r AutoScalingConfigurationApiRequest) Execute() (*ClusterDescriptionAutoScalingModeConfiguration, *http.Response, error) {
	return r.ApiService.AutoScalingConfigurationExecute(r)
}

/*
AutoScalingConfiguration Return Auto Scaling Configuration for One Sharded Cluster

Returns the internal configuration of AutoScaling for sharded clusters. This endpoint can be used for diagnostic purposes to ensure that sharded clusters updated from older APIs have gained support for AutoScaling each shard independently.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return AutoScalingConfigurationApiRequest

Deprecated
*/
func (a *ClustersApiService) AutoScalingConfiguration(ctx context.Context, groupId string, clusterName string) AutoScalingConfigurationApiRequest {
	return AutoScalingConfigurationApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// AutoScalingConfigurationExecute executes the request
//
//	@return ClusterDescriptionAutoScalingModeConfiguration
//
// Deprecated
func (a *ClustersApiService) AutoScalingConfigurationExecute(r AutoScalingConfigurationApiRequest) (*ClusterDescriptionAutoScalingModeConfiguration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescriptionAutoScalingModeConfiguration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.AutoScalingConfiguration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/autoScalingConfiguration"
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

type CreateClusterApiRequest struct {
	ctx                                context.Context
	ApiService                         ClustersApi
	groupId                            string
	clusterDescription20240805         *ClusterDescription20240805
	useEffectiveInstanceFields         *bool
	useEffectiveFieldsReplicationSpecs *bool
}

type CreateClusterApiParams struct {
	GroupId                            string
	ClusterDescription20240805         *ClusterDescription20240805
	UseEffectiveInstanceFields         *bool
	UseEffectiveFieldsReplicationSpecs *bool
}

func (a *ClustersApiService) CreateClusterWithParams(ctx context.Context, args *CreateClusterApiParams) CreateClusterApiRequest {
	return CreateClusterApiRequest{
		ApiService:                         a,
		ctx:                                ctx,
		groupId:                            args.GroupId,
		clusterDescription20240805:         args.ClusterDescription20240805,
		useEffectiveInstanceFields:         args.UseEffectiveInstanceFields,
		useEffectiveFieldsReplicationSpecs: args.UseEffectiveFieldsReplicationSpecs,
	}
}

// Controls how hardware specification fields are returned in the response after cluster creation. When set to true, returns the original client-specified values and provides separate effective fields showing current operational values. When false (default), hardware specification fields show current operational values directly. Primarily used for autoscaling compatibility.
func (r CreateClusterApiRequest) UseEffectiveInstanceFields(useEffectiveInstanceFields bool) CreateClusterApiRequest {
	r.useEffectiveInstanceFields = &useEffectiveInstanceFields
	return r
}

// Controls how &#x60;replicationSpecs&#x60; fields are returned in the response. When set to &#x60;true&#x60;, stores the client&#39;s view of &#x60;replicationSpecs&#x60; and returns it in &#x60;replicationSpecs&#x60;, while the actual cluster state (including auto-scaled hardware and auto-added shards) is returned in &#x60;effectiveReplicationSpecs&#x60;. When &#x60;false&#x60; (default), &#x60;replicationSpecs&#x60; contains the actual cluster state.
func (r CreateClusterApiRequest) UseEffectiveFieldsReplicationSpecs(useEffectiveFieldsReplicationSpecs bool) CreateClusterApiRequest {
	r.useEffectiveFieldsReplicationSpecs = &useEffectiveFieldsReplicationSpecs
	return r
}

func (r CreateClusterApiRequest) Execute() (*ClusterDescription20240805, *http.Response, error) {
	return r.ApiService.CreateClusterExecute(r)
}

/*
CreateCluster Create One Cluster in One Project

Creates one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. This resource can create clusters with asymmetrically-sized shards. Each project supports up to 25 database deployments. This feature is not available for serverless clusters.

Please note that using an `instanceSize` of M2 or M5 will create a Flex cluster instead. Support for the `instanceSize` of M2 or M5 will be discontinued in January 2026. We recommend using the Create Flex Cluster API for such configurations moving forward. Deprecated versions: v2-{2024-08-05}, v2-{2023-02-01}, v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateClusterApiRequest
*/
func (a *ClustersApiService) CreateCluster(ctx context.Context, groupId string, clusterDescription20240805 *ClusterDescription20240805) CreateClusterApiRequest {
	return CreateClusterApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterDescription20240805: clusterDescription20240805,
	}
}

// CreateClusterExecute executes the request
//
//	@return ClusterDescription20240805
func (a *ClustersApiService) CreateClusterExecute(r CreateClusterApiRequest) (*ClusterDescription20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescription20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.CreateCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.clusterDescription20240805 == nil {
		return localVarReturnValue, nil, reportError("clusterDescription20240805 is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-10-23+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-10-23+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.useEffectiveInstanceFields != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Instance-Fields", r.useEffectiveInstanceFields, "")
	}
	if r.useEffectiveFieldsReplicationSpecs != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Fields-Replication-Specs", r.useEffectiveFieldsReplicationSpecs, "")
	}
	// body params
	localVarPostBody = r.clusterDescription20240805
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

type DeleteClusterApiRequest struct {
	ctx           context.Context
	ApiService    ClustersApi
	groupId       string
	clusterName   string
	retainBackups *bool
}

type DeleteClusterApiParams struct {
	GroupId       string
	ClusterName   string
	RetainBackups *bool
}

func (a *ClustersApiService) DeleteClusterWithParams(ctx context.Context, args *DeleteClusterApiParams) DeleteClusterApiRequest {
	return DeleteClusterApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		clusterName:   args.ClusterName,
		retainBackups: args.RetainBackups,
	}
}

// Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster.
func (r DeleteClusterApiRequest) RetainBackups(retainBackups bool) DeleteClusterApiRequest {
	r.retainBackups = &retainBackups
	return r
}

func (r DeleteClusterApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterExecute(r)
}

/*
DeleteCluster Remove One Cluster from One Project

Removes one cluster from the specified project. The cluster must have termination protection disabled in order to be deleted. This feature is not available for serverless clusters.

This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the Delete Flex Cluster endpoint for Flex clusters instead. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return DeleteClusterApiRequest
*/
func (a *ClustersApiService) DeleteCluster(ctx context.Context, groupId string, clusterName string) DeleteClusterApiRequest {
	return DeleteClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// DeleteClusterExecute executes the request
func (a *ClustersApiService) DeleteClusterExecute(r DeleteClusterApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.DeleteCluster")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}"
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

	if r.retainBackups != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "retainBackups", r.retainBackups, "")
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

type GetClusterApiRequest struct {
	ctx                                context.Context
	ApiService                         ClustersApi
	groupId                            string
	clusterName                        string
	useEffectiveInstanceFields         *bool
	useEffectiveFieldsReplicationSpecs *bool
}

type GetClusterApiParams struct {
	GroupId                            string
	ClusterName                        string
	UseEffectiveInstanceFields         *bool
	UseEffectiveFieldsReplicationSpecs *bool
}

func (a *ClustersApiService) GetClusterWithParams(ctx context.Context, args *GetClusterApiParams) GetClusterApiRequest {
	return GetClusterApiRequest{
		ApiService:                         a,
		ctx:                                ctx,
		groupId:                            args.GroupId,
		clusterName:                        args.ClusterName,
		useEffectiveInstanceFields:         args.UseEffectiveInstanceFields,
		useEffectiveFieldsReplicationSpecs: args.UseEffectiveFieldsReplicationSpecs,
	}
}

// Controls how hardware specification fields are returned in the response. When set to true, returns the original client-specified values and provides separate effective fields showing current operational values. When false (default), hardware specification fields show current operational values directly. Primarily used for autoscaling compatibility.
func (r GetClusterApiRequest) UseEffectiveInstanceFields(useEffectiveInstanceFields bool) GetClusterApiRequest {
	r.useEffectiveInstanceFields = &useEffectiveInstanceFields
	return r
}

// Controls how &#x60;replicationSpecs&#x60; are returned in the response. When set to &#x60;true&#x60;, returns the client-specified view in &#x60;replicationSpecs&#x60; and the actual cluster state in &#x60;effectiveReplicationSpecs&#x60;. When &#x60;false&#x60; (default), &#x60;replicationSpecs&#x60; contains the actual cluster state.
func (r GetClusterApiRequest) UseEffectiveFieldsReplicationSpecs(useEffectiveFieldsReplicationSpecs bool) GetClusterApiRequest {
	r.useEffectiveFieldsReplicationSpecs = &useEffectiveFieldsReplicationSpecs
	return r
}

func (r GetClusterApiRequest) Execute() (*ClusterDescription20240805, *http.Response, error) {
	return r.ApiService.GetClusterExecute(r)
}

/*
GetCluster Return One Cluster from One Project

Returns the details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. The response includes clusters with asymmetrically-sized shards. This feature is not available for serverless clusters.

This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the Get Flex Cluster endpoint for Flex clusters instead. Deprecated versions: v2-{2023-02-01}, v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return GetClusterApiRequest
*/
func (a *ClustersApiService) GetCluster(ctx context.Context, groupId string, clusterName string) GetClusterApiRequest {
	return GetClusterApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetClusterExecute executes the request
//
//	@return ClusterDescription20240805
func (a *ClustersApiService) GetClusterExecute(r GetClusterApiRequest) (*ClusterDescription20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescription20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.GetCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}"
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
	if r.useEffectiveInstanceFields != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Instance-Fields", r.useEffectiveInstanceFields, "")
	}
	if r.useEffectiveFieldsReplicationSpecs != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Fields-Replication-Specs", r.useEffectiveFieldsReplicationSpecs, "")
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

type GetClusterStatusApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type GetClusterStatusApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) GetClusterStatusWithParams(ctx context.Context, args *GetClusterStatusApiParams) GetClusterStatusApiRequest {
	return GetClusterStatusApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetClusterStatusApiRequest) Execute() (*ClusterStatus, *http.Response, error) {
	return r.ApiService.GetClusterStatusExecute(r)
}

/*
GetClusterStatus Return Status of All Cluster Operations

Returns the status of all changes that you made to the specified cluster in the specified project. Use this resource to check the progress MongoDB Cloud has made in processing your changes. The response does not include the deployment of new dedicated clusters.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return GetClusterStatusApiRequest
*/
func (a *ClustersApiService) GetClusterStatus(ctx context.Context, groupId string, clusterName string) GetClusterStatusApiRequest {
	return GetClusterStatusApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetClusterStatusExecute executes the request
//
//	@return ClusterStatus
func (a *ClustersApiService) GetClusterStatusExecute(r GetClusterStatusApiRequest) (*ClusterStatus, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.GetClusterStatus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/status"
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

type GetProcessArgsApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type GetProcessArgsApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) GetProcessArgsWithParams(ctx context.Context, args *GetProcessArgsApiParams) GetProcessArgsApiRequest {
	return GetProcessArgsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetProcessArgsApiRequest) Execute() (*ClusterDescriptionProcessArgs20240805, *http.Response, error) {
	return r.ApiService.GetProcessArgsExecute(r)
}

/*
GetProcessArgs Return Advanced Configuration Options for One Cluster

Returns the advanced configuration details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. Advanced configuration details include the read/write concern, index and oplog limits, and other database settings. This feature isn't available for `M0` free clusters, `M2` and `M5` shared-tier clusters, flex clusters, or serverless clusters. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return GetProcessArgsApiRequest
*/
func (a *ClustersApiService) GetProcessArgs(ctx context.Context, groupId string, clusterName string) GetProcessArgsApiRequest {
	return GetProcessArgsApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetProcessArgsExecute executes the request
//
//	@return ClusterDescriptionProcessArgs20240805
func (a *ClustersApiService) GetProcessArgsExecute(r GetProcessArgsApiRequest) (*ClusterDescriptionProcessArgs20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescriptionProcessArgs20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.GetProcessArgs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/processArgs"
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

type GetSampleDatasetLoadApiRequest struct {
	ctx             context.Context
	ApiService      ClustersApi
	groupId         string
	sampleDatasetId string
}

type GetSampleDatasetLoadApiParams struct {
	GroupId         string
	SampleDatasetId string
}

func (a *ClustersApiService) GetSampleDatasetLoadWithParams(ctx context.Context, args *GetSampleDatasetLoadApiParams) GetSampleDatasetLoadApiRequest {
	return GetSampleDatasetLoadApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         args.GroupId,
		sampleDatasetId: args.SampleDatasetId,
	}
}

func (r GetSampleDatasetLoadApiRequest) Execute() (*SampleDatasetStatus, *http.Response, error) {
	return r.ApiService.GetSampleDatasetLoadExecute(r)
}

/*
GetSampleDatasetLoad Return Status of Sample Dataset Load for One Cluster

Checks the progress of loading the sample dataset into one cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param sampleDatasetId Unique 24-hexadecimal digit string that identifies the loaded sample dataset.
	@return GetSampleDatasetLoadApiRequest
*/
func (a *ClustersApiService) GetSampleDatasetLoad(ctx context.Context, groupId string, sampleDatasetId string) GetSampleDatasetLoadApiRequest {
	return GetSampleDatasetLoadApiRequest{
		ApiService:      a,
		ctx:             ctx,
		groupId:         groupId,
		sampleDatasetId: sampleDatasetId,
	}
}

// GetSampleDatasetLoadExecute executes the request
//
//	@return SampleDatasetStatus
func (a *ClustersApiService) GetSampleDatasetLoadExecute(r GetSampleDatasetLoadApiRequest) (*SampleDatasetStatus, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SampleDatasetStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.GetSampleDatasetLoad")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/sampleDatasetLoad/{sampleDatasetId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.sampleDatasetId == "" {
		return localVarReturnValue, nil, reportError("sampleDatasetId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"sampleDatasetId"+"}", url.PathEscape(r.sampleDatasetId), -1)

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

type GrantMongoEmployeeAccessApiRequest struct {
	ctx                 context.Context
	ApiService          ClustersApi
	groupId             string
	clusterName         string
	employeeAccessGrant *EmployeeAccessGrant
}

type GrantMongoEmployeeAccessApiParams struct {
	GroupId             string
	ClusterName         string
	EmployeeAccessGrant *EmployeeAccessGrant
}

func (a *ClustersApiService) GrantMongoEmployeeAccessWithParams(ctx context.Context, args *GrantMongoEmployeeAccessApiParams) GrantMongoEmployeeAccessApiRequest {
	return GrantMongoEmployeeAccessApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             args.GroupId,
		clusterName:         args.ClusterName,
		employeeAccessGrant: args.EmployeeAccessGrant,
	}
}

func (r GrantMongoEmployeeAccessApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.GrantMongoEmployeeAccessExecute(r)
}

/*
GrantMongoEmployeeAccess Grant MongoDB Employee Cluster Access for One Cluster

Grants MongoDB employee cluster access for the given duration and at the specified level for one cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return GrantMongoEmployeeAccessApiRequest
*/
func (a *ClustersApiService) GrantMongoEmployeeAccess(ctx context.Context, groupId string, clusterName string, employeeAccessGrant *EmployeeAccessGrant) GrantMongoEmployeeAccessApiRequest {
	return GrantMongoEmployeeAccessApiRequest{
		ApiService:          a,
		ctx:                 ctx,
		groupId:             groupId,
		clusterName:         clusterName,
		employeeAccessGrant: employeeAccessGrant,
	}
}

// GrantMongoEmployeeAccessExecute executes the request
func (a *ClustersApiService) GrantMongoEmployeeAccessExecute(r GrantMongoEmployeeAccessApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.GrantMongoEmployeeAccess")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}:grantMongoDBEmployeeAccess"
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
	if r.employeeAccessGrant == nil {
		return nil, reportError("employeeAccessGrant is required and must be specified")
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
	localVarPostBody = r.employeeAccessGrant
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

type ListClusterDetailsApiRequest struct {
	ctx          context.Context
	ApiService   ClustersApi
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListClusterDetailsApiParams struct {
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *ClustersApiService) ListClusterDetailsWithParams(ctx context.Context, args *ListClusterDetailsApiParams) ListClusterDetailsApiRequest {
	return ListClusterDetailsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClusterDetailsApiRequest) IncludeCount(includeCount bool) ListClusterDetailsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterDetailsApiRequest) ItemsPerPage(itemsPerPage int) ListClusterDetailsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterDetailsApiRequest) PageNum(pageNum int) ListClusterDetailsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListClusterDetailsApiRequest) Execute() (*PaginatedOrgGroup, *http.Response, error) {
	return r.ApiService.ListClusterDetailsExecute(r)
}

/*
ListClusterDetails Return All Authorized Clusters in All Projects

Returns the details for all clusters in all projects to which you have access. Clusters contain a group of hosts that maintain the same data set. The response does not include multi-cloud clusters.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListClusterDetailsApiRequest
*/
func (a *ClustersApiService) ListClusterDetails(ctx context.Context) ListClusterDetailsApiRequest {
	return ListClusterDetailsApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListClusterDetailsExecute executes the request
//
//	@return PaginatedOrgGroup
func (a *ClustersApiService) ListClusterDetailsExecute(r ListClusterDetailsApiRequest) (*PaginatedOrgGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOrgGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.ListClusterDetails")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/clusters"

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

type ListClusterProviderRegionsApiRequest struct {
	ctx          context.Context
	ApiService   ClustersApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	providers    *[]string
	tier         *string
}

type ListClusterProviderRegionsApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Providers    *[]string
	Tier         *string
}

func (a *ClustersApiService) ListClusterProviderRegionsWithParams(ctx context.Context, args *ListClusterProviderRegionsApiParams) ListClusterProviderRegionsApiRequest {
	return ListClusterProviderRegionsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		providers:    args.Providers,
		tier:         args.Tier,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClusterProviderRegionsApiRequest) IncludeCount(includeCount bool) ListClusterProviderRegionsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClusterProviderRegionsApiRequest) ItemsPerPage(itemsPerPage int) ListClusterProviderRegionsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClusterProviderRegionsApiRequest) PageNum(pageNum int) ListClusterProviderRegionsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Cloud providers whose regions to retrieve. When you specify multiple providers, the response can return only tiers and regions that support multi-cloud clusters.
func (r ListClusterProviderRegionsApiRequest) Providers(providers []string) ListClusterProviderRegionsApiRequest {
	r.providers = &providers
	return r
}

// Cluster tier for which to retrieve the regions.
func (r ListClusterProviderRegionsApiRequest) Tier(tier string) ListClusterProviderRegionsApiRequest {
	r.tier = &tier
	return r
}

func (r ListClusterProviderRegionsApiRequest) Execute() (*PaginatedApiAtlasProviderRegions, *http.Response, error) {
	return r.ApiService.ListClusterProviderRegionsExecute(r)
}

/*
ListClusterProviderRegions Return All Cloud Provider Regions

Returns the list of regions available for the specified cloud provider at the specified tier.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListClusterProviderRegionsApiRequest
*/
func (a *ClustersApiService) ListClusterProviderRegions(ctx context.Context, groupId string) ListClusterProviderRegionsApiRequest {
	return ListClusterProviderRegionsApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListClusterProviderRegionsExecute executes the request
//
//	@return PaginatedApiAtlasProviderRegions
func (a *ClustersApiService) ListClusterProviderRegionsExecute(r ListClusterProviderRegionsApiRequest) (*PaginatedApiAtlasProviderRegions, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiAtlasProviderRegions
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.ListClusterProviderRegions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/provider/regions"
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
	if r.providers != nil {
		t := *r.providers
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "providers", t, "multi")

	}
	if r.tier != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "tier", r.tier, "")
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

type ListClustersApiRequest struct {
	ctx                               context.Context
	ApiService                        ClustersApi
	groupId                           string
	includeCount                      *bool
	itemsPerPage                      *int
	pageNum                           *int
	includeDeletedWithRetainedBackups *bool
	useEffectiveInstanceFields        *bool
}

type ListClustersApiParams struct {
	GroupId                           string
	IncludeCount                      *bool
	ItemsPerPage                      *int
	PageNum                           *int
	IncludeDeletedWithRetainedBackups *bool
	UseEffectiveInstanceFields        *bool
}

func (a *ClustersApiService) ListClustersWithParams(ctx context.Context, args *ListClustersApiParams) ListClustersApiRequest {
	return ListClustersApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		groupId:                           args.GroupId,
		includeCount:                      args.IncludeCount,
		itemsPerPage:                      args.ItemsPerPage,
		pageNum:                           args.PageNum,
		includeDeletedWithRetainedBackups: args.IncludeDeletedWithRetainedBackups,
		useEffectiveInstanceFields:        args.UseEffectiveInstanceFields,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListClustersApiRequest) IncludeCount(includeCount bool) ListClustersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListClustersApiRequest) ItemsPerPage(itemsPerPage int) ListClustersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListClustersApiRequest) PageNum(pageNum int) ListClustersApiRequest {
	r.pageNum = &pageNum
	return r
}

// Flag that indicates whether to return Clusters with retain backups.
func (r ListClustersApiRequest) IncludeDeletedWithRetainedBackups(includeDeletedWithRetainedBackups bool) ListClustersApiRequest {
	r.includeDeletedWithRetainedBackups = &includeDeletedWithRetainedBackups
	return r
}

// Controls how hardware specification fields are returned in the response. When set to true, returns the original client-specified values and provides separate effective fields showing current operational values. When false (default), hardware specification fields show current operational values directly. Primarily used for autoscaling compatibility.
func (r ListClustersApiRequest) UseEffectiveInstanceFields(useEffectiveInstanceFields bool) ListClustersApiRequest {
	r.useEffectiveInstanceFields = &useEffectiveInstanceFields
	return r
}

func (r ListClustersApiRequest) Execute() (*PaginatedClusterDescription20240805, *http.Response, error) {
	return r.ApiService.ListClustersExecute(r)
}

/*
ListClusters Return All Clusters in One Project

Returns the details for all clusters in the specific project to which you have access. Clusters contain a group of hosts that maintain the same data set. The response includes clusters with asymmetrically-sized shards. This feature is not  available for serverless clusters.

This endpoint can also be used on Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) endpoint or former M2/M5 clusters that have been migrated to Flex clusters until January 2026. Please use the List Flex Clusters endpoint for Flex clusters instead. Deprecated versions: v2-{2023-02-01}, v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListClustersApiRequest
*/
func (a *ClustersApiService) ListClusters(ctx context.Context, groupId string) ListClustersApiRequest {
	return ListClustersApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListClustersExecute executes the request
//
//	@return PaginatedClusterDescription20240805
func (a *ClustersApiService) ListClustersExecute(r ListClustersApiRequest) (*PaginatedClusterDescription20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedClusterDescription20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.ListClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters"
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
	if r.includeDeletedWithRetainedBackups != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeDeletedWithRetainedBackups", r.includeDeletedWithRetainedBackups, "")
	} else {
		var defaultValue bool = false
		r.includeDeletedWithRetainedBackups = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeDeletedWithRetainedBackups", r.includeDeletedWithRetainedBackups, "")
	}
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
	if r.useEffectiveInstanceFields != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Instance-Fields", r.useEffectiveInstanceFields, "")
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

type PinFeatureCompatibilityVersionApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
	pinFCV      *PinFCV
}

type PinFeatureCompatibilityVersionApiParams struct {
	GroupId     string
	ClusterName string
	PinFCV      *PinFCV
}

func (a *ClustersApiService) PinFeatureCompatibilityVersionWithParams(ctx context.Context, args *PinFeatureCompatibilityVersionApiParams) PinFeatureCompatibilityVersionApiRequest {
	return PinFeatureCompatibilityVersionApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		pinFCV:      args.PinFCV,
	}
}

func (r PinFeatureCompatibilityVersionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.PinFeatureCompatibilityVersionExecute(r)
}

/*
PinFeatureCompatibilityVersion Pin Feature Compatibility Version for One Cluster in One Project

Pins the Feature Compatibility Version (FCV) to the current MongoDB version and sets the pin expiration date. If an FCV pin already exists for the cluster, calling this method will only update the expiration date of the existing pin and will not re-pin the FCV.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return PinFeatureCompatibilityVersionApiRequest
*/
func (a *ClustersApiService) PinFeatureCompatibilityVersion(ctx context.Context, groupId string, clusterName string, pinFCV *PinFCV) PinFeatureCompatibilityVersionApiRequest {
	return PinFeatureCompatibilityVersionApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		pinFCV:      pinFCV,
	}
}

// PinFeatureCompatibilityVersionExecute executes the request
func (a *ClustersApiService) PinFeatureCompatibilityVersionExecute(r PinFeatureCompatibilityVersionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.PinFeatureCompatibilityVersion")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}:pinFeatureCompatibilityVersion"
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
	localVarPostBody = r.pinFCV
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

type RequestSampleDatasetLoadApiRequest struct {
	ctx        context.Context
	ApiService ClustersApi
	groupId    string
	name       string
}

type RequestSampleDatasetLoadApiParams struct {
	GroupId string
	Name    string
}

func (a *ClustersApiService) RequestSampleDatasetLoadWithParams(ctx context.Context, args *RequestSampleDatasetLoadApiParams) RequestSampleDatasetLoadApiRequest {
	return RequestSampleDatasetLoadApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		name:       args.Name,
	}
}

func (r RequestSampleDatasetLoadApiRequest) Execute() (*SampleDatasetStatus, *http.Response, error) {
	return r.ApiService.RequestSampleDatasetLoadExecute(r)
}

/*
RequestSampleDatasetLoad Load Sample Dataset into One Cluster

Requests loading the MongoDB sample dataset into the specified cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param name Human-readable label that identifies the cluster into which you load the sample dataset.
	@return RequestSampleDatasetLoadApiRequest
*/
func (a *ClustersApiService) RequestSampleDatasetLoad(ctx context.Context, groupId string, name string) RequestSampleDatasetLoadApiRequest {
	return RequestSampleDatasetLoadApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		name:       name,
	}
}

// RequestSampleDatasetLoadExecute executes the request
//
//	@return SampleDatasetStatus
func (a *ClustersApiService) RequestSampleDatasetLoadExecute(r RequestSampleDatasetLoadApiRequest) (*SampleDatasetStatus, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SampleDatasetStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.RequestSampleDatasetLoad")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/sampleDatasetLoad/{name}"
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

type RestartPrimariesApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type RestartPrimariesApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) RestartPrimariesWithParams(ctx context.Context, args *RestartPrimariesApiParams) RestartPrimariesApiRequest {
	return RestartPrimariesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r RestartPrimariesApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RestartPrimariesExecute(r)
}

/*
RestartPrimaries Test Failover for One Cluster

Starts a failover test for the specified cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. A failover test checks how MongoDB Cloud handles the failure of the cluster's primary node. During the test, MongoDB Cloud shuts down the primary node and elects a new primary. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return RestartPrimariesApiRequest
*/
func (a *ClustersApiService) RestartPrimaries(ctx context.Context, groupId string, clusterName string) RestartPrimariesApiRequest {
	return RestartPrimariesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// RestartPrimariesExecute executes the request
func (a *ClustersApiService) RestartPrimariesExecute(r RestartPrimariesApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.RestartPrimaries")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/restartPrimaries"
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

type RevokeMongoEmployeeAccessApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type RevokeMongoEmployeeAccessApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) RevokeMongoEmployeeAccessWithParams(ctx context.Context, args *RevokeMongoEmployeeAccessApiParams) RevokeMongoEmployeeAccessApiRequest {
	return RevokeMongoEmployeeAccessApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r RevokeMongoEmployeeAccessApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.RevokeMongoEmployeeAccessExecute(r)
}

/*
RevokeMongoEmployeeAccess Revoke MongoDB Employee Cluster Access for One Cluster

Revokes a previously granted MongoDB employee cluster access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return RevokeMongoEmployeeAccessApiRequest
*/
func (a *ClustersApiService) RevokeMongoEmployeeAccess(ctx context.Context, groupId string, clusterName string) RevokeMongoEmployeeAccessApiRequest {
	return RevokeMongoEmployeeAccessApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// RevokeMongoEmployeeAccessExecute executes the request
func (a *ClustersApiService) RevokeMongoEmployeeAccessExecute(r RevokeMongoEmployeeAccessApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.RevokeMongoEmployeeAccess")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}:revokeMongoDBEmployeeAccess"
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

type UnpinFeatureCompatibilityVersionApiRequest struct {
	ctx         context.Context
	ApiService  ClustersApi
	groupId     string
	clusterName string
}

type UnpinFeatureCompatibilityVersionApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *ClustersApiService) UnpinFeatureCompatibilityVersionWithParams(ctx context.Context, args *UnpinFeatureCompatibilityVersionApiParams) UnpinFeatureCompatibilityVersionApiRequest {
	return UnpinFeatureCompatibilityVersionApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r UnpinFeatureCompatibilityVersionApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.UnpinFeatureCompatibilityVersionExecute(r)
}

/*
UnpinFeatureCompatibilityVersion Unpin Feature Compatibility Version for One Cluster in One Project

Unpins the current fixed Feature Compatibility Version (FCV). This feature is not available for clusters on rapid release.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies this cluster.
	@return UnpinFeatureCompatibilityVersionApiRequest
*/
func (a *ClustersApiService) UnpinFeatureCompatibilityVersion(ctx context.Context, groupId string, clusterName string) UnpinFeatureCompatibilityVersionApiRequest {
	return UnpinFeatureCompatibilityVersionApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// UnpinFeatureCompatibilityVersionExecute executes the request
func (a *ClustersApiService) UnpinFeatureCompatibilityVersionExecute(r UnpinFeatureCompatibilityVersionApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.UnpinFeatureCompatibilityVersion")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}:unpinFeatureCompatibilityVersion"
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

type UpdateClusterApiRequest struct {
	ctx                                context.Context
	ApiService                         ClustersApi
	groupId                            string
	clusterName                        string
	clusterDescription20240805         *ClusterDescription20240805
	useEffectiveInstanceFields         *bool
	useEffectiveFieldsReplicationSpecs *bool
}

type UpdateClusterApiParams struct {
	GroupId                            string
	ClusterName                        string
	ClusterDescription20240805         *ClusterDescription20240805
	UseEffectiveInstanceFields         *bool
	UseEffectiveFieldsReplicationSpecs *bool
}

func (a *ClustersApiService) UpdateClusterWithParams(ctx context.Context, args *UpdateClusterApiParams) UpdateClusterApiRequest {
	return UpdateClusterApiRequest{
		ApiService:                         a,
		ctx:                                ctx,
		groupId:                            args.GroupId,
		clusterName:                        args.ClusterName,
		clusterDescription20240805:         args.ClusterDescription20240805,
		useEffectiveInstanceFields:         args.UseEffectiveInstanceFields,
		useEffectiveFieldsReplicationSpecs: args.UseEffectiveFieldsReplicationSpecs,
	}
}

// Controls how hardware specification fields are returned in the response after cluster updates. When set to true, returns the original client-specified values and provides separate effective fields showing current operational values. When false (default), hardware specification fields show current operational values directly. Note: When using this header with autoscaling enabled, MongoDB ignores &#x60;replicationSpecs&#x60; changes during updates. To intentionally override the &#x60;replicationSpecs&#x60;, disable this header.
func (r UpdateClusterApiRequest) UseEffectiveInstanceFields(useEffectiveInstanceFields bool) UpdateClusterApiRequest {
	r.useEffectiveInstanceFields = &useEffectiveInstanceFields
	return r
}

// Controls how &#x60;replicationSpecs&#x60; fields are returned in the response. When set to &#x60;true&#x60;, stores the client&#39;s view of &#x60;replicationSpecs&#x60; and returns it in &#x60;replicationSpecs&#x60;, while the actual cluster state (including auto-scaled hardware and auto-added shards) is returned in &#x60;effectiveReplicationSpecs&#x60;. When &#x60;false&#x60; (default), &#x60;replicationSpecs&#x60; contains the actual cluster state.
func (r UpdateClusterApiRequest) UseEffectiveFieldsReplicationSpecs(useEffectiveFieldsReplicationSpecs bool) UpdateClusterApiRequest {
	r.useEffectiveFieldsReplicationSpecs = &useEffectiveFieldsReplicationSpecs
	return r
}

func (r UpdateClusterApiRequest) Execute() (*ClusterDescription20240805, *http.Response, error) {
	return r.ApiService.UpdateClusterExecute(r)
}

/*
UpdateCluster Update One Cluster in One Project

Updates the details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. This resource can update clusters with asymmetrically-sized shards. To update a cluster's termination protection, the requesting Service Account or API Key must have the Project Owner role. For all other updates, the requesting Service Account or API Key must have the Project Cluster Manager role or the Project Replica Set Manager role. You can't modify a paused cluster (`paused : true`). You must call this endpoint to set `paused : false`. After this endpoint responds with `paused : false`, you can call it again with the changes you want to make to the cluster. This feature is not available for serverless clusters. Deprecated versions: v2-{2024-08-05}, v2-{2023-02-01}, v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return UpdateClusterApiRequest
*/
func (a *ClustersApiService) UpdateCluster(ctx context.Context, groupId string, clusterName string, clusterDescription20240805 *ClusterDescription20240805) UpdateClusterApiRequest {
	return UpdateClusterApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterName:                clusterName,
		clusterDescription20240805: clusterDescription20240805,
	}
}

// UpdateClusterExecute executes the request
//
//	@return ClusterDescription20240805
func (a *ClustersApiService) UpdateClusterExecute(r UpdateClusterApiRequest) (*ClusterDescription20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescription20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.UpdateCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}"
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
	if r.clusterDescription20240805 == nil {
		return localVarReturnValue, nil, reportError("clusterDescription20240805 is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-10-23+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-10-23+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.useEffectiveInstanceFields != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Instance-Fields", r.useEffectiveInstanceFields, "")
	}
	if r.useEffectiveFieldsReplicationSpecs != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Use-Effective-Fields-Replication-Specs", r.useEffectiveFieldsReplicationSpecs, "")
	}
	// body params
	localVarPostBody = r.clusterDescription20240805
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

type UpdateProcessArgsApiRequest struct {
	ctx                                   context.Context
	ApiService                            ClustersApi
	groupId                               string
	clusterName                           string
	clusterDescriptionProcessArgs20240805 *ClusterDescriptionProcessArgs20240805
}

type UpdateProcessArgsApiParams struct {
	GroupId                               string
	ClusterName                           string
	ClusterDescriptionProcessArgs20240805 *ClusterDescriptionProcessArgs20240805
}

func (a *ClustersApiService) UpdateProcessArgsWithParams(ctx context.Context, args *UpdateProcessArgsApiParams) UpdateProcessArgsApiRequest {
	return UpdateProcessArgsApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               args.GroupId,
		clusterName:                           args.ClusterName,
		clusterDescriptionProcessArgs20240805: args.ClusterDescriptionProcessArgs20240805,
	}
}

func (r UpdateProcessArgsApiRequest) Execute() (*ClusterDescriptionProcessArgs20240805, *http.Response, error) {
	return r.ApiService.UpdateProcessArgsExecute(r)
}

/*
UpdateProcessArgs Update Advanced Configuration Options for One Cluster

Updates the advanced configuration details for one cluster in the specified project. Clusters contain a group of hosts that maintain the same data set. Advanced configuration details include the read/write concern, index and oplog limits, and other database settings. This feature isn't available for `M0` free clusters, `M2` and `M5` shared-tier clusters, flex clusters, or serverless clusters. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Human-readable label that identifies the cluster.
	@return UpdateProcessArgsApiRequest
*/
func (a *ClustersApiService) UpdateProcessArgs(ctx context.Context, groupId string, clusterName string, clusterDescriptionProcessArgs20240805 *ClusterDescriptionProcessArgs20240805) UpdateProcessArgsApiRequest {
	return UpdateProcessArgsApiRequest{
		ApiService:                            a,
		ctx:                                   ctx,
		groupId:                               groupId,
		clusterName:                           clusterName,
		clusterDescriptionProcessArgs20240805: clusterDescriptionProcessArgs20240805,
	}
}

// UpdateProcessArgsExecute executes the request
//
//	@return ClusterDescriptionProcessArgs20240805
func (a *ClustersApiService) UpdateProcessArgsExecute(r UpdateProcessArgsApiRequest) (*ClusterDescriptionProcessArgs20240805, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterDescriptionProcessArgs20240805
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.UpdateProcessArgs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/processArgs"
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
	if r.clusterDescriptionProcessArgs20240805 == nil {
		return localVarReturnValue, nil, reportError("clusterDescriptionProcessArgs20240805 is required and must be specified")
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
	localVarPostBody = r.clusterDescriptionProcessArgs20240805
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

type UpgradeClusterToServerlessApiRequest struct {
	ctx                           context.Context
	ApiService                    ClustersApi
	groupId                       string
	serverlessInstanceDescription *ServerlessInstanceDescription
}

type UpgradeClusterToServerlessApiParams struct {
	GroupId                       string
	ServerlessInstanceDescription *ServerlessInstanceDescription
}

func (a *ClustersApiService) UpgradeClusterToServerlessWithParams(ctx context.Context, args *UpgradeClusterToServerlessApiParams) UpgradeClusterToServerlessApiRequest {
	return UpgradeClusterToServerlessApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       args.GroupId,
		serverlessInstanceDescription: args.ServerlessInstanceDescription,
	}
}

func (r UpgradeClusterToServerlessApiRequest) Execute() (*ServerlessInstanceDescription, *http.Response, error) {
	return r.ApiService.UpgradeClusterToServerlessExecute(r)
}

/*
UpgradeClusterToServerless Upgrade One Shared-Tier Cluster to One Serverless Instance

This endpoint has been deprecated as of February 2025 as we no longer support the creation of new serverless instances. Please use the Upgrade Flex Cluster endpoint to upgrade Flex clusters.

	Upgrades a shared-tier cluster to a serverless instance in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpgradeClusterToServerlessApiRequest

Deprecated
*/
func (a *ClustersApiService) UpgradeClusterToServerless(ctx context.Context, groupId string, serverlessInstanceDescription *ServerlessInstanceDescription) UpgradeClusterToServerlessApiRequest {
	return UpgradeClusterToServerlessApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		groupId:                       groupId,
		serverlessInstanceDescription: serverlessInstanceDescription,
	}
}

// UpgradeClusterToServerlessExecute executes the request
//
//	@return ServerlessInstanceDescription
//
// Deprecated
func (a *ClustersApiService) UpgradeClusterToServerlessExecute(r UpgradeClusterToServerlessApiRequest) (*ServerlessInstanceDescription, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ServerlessInstanceDescription
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.UpgradeClusterToServerless")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/tenantUpgradeToServerless"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.serverlessInstanceDescription == nil {
		return localVarReturnValue, nil, reportError("serverlessInstanceDescription is required and must be specified")
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
	localVarPostBody = r.serverlessInstanceDescription
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

type UpgradeTenantUpgradeApiRequest struct {
	ctx                                    context.Context
	ApiService                             ClustersApi
	groupId                                string
	legacyAtlasTenantClusterUpgradeRequest *LegacyAtlasTenantClusterUpgradeRequest
}

type UpgradeTenantUpgradeApiParams struct {
	GroupId                                string
	LegacyAtlasTenantClusterUpgradeRequest *LegacyAtlasTenantClusterUpgradeRequest
}

func (a *ClustersApiService) UpgradeTenantUpgradeWithParams(ctx context.Context, args *UpgradeTenantUpgradeApiParams) UpgradeTenantUpgradeApiRequest {
	return UpgradeTenantUpgradeApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                args.GroupId,
		legacyAtlasTenantClusterUpgradeRequest: args.LegacyAtlasTenantClusterUpgradeRequest,
	}
}

func (r UpgradeTenantUpgradeApiRequest) Execute() (*LegacyAtlasCluster, *http.Response, error) {
	return r.ApiService.UpgradeTenantUpgradeExecute(r)
}

/*
UpgradeTenantUpgrade Upgrade One Shared-Tier Cluster

Upgrades a shared-tier cluster to a Flex or Dedicated (M10+) cluster in the specified project. Each project supports up to 25 clusters.

This endpoint can also be used to upgrade Flex clusters that were created using the [Create Cluster](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createCluster) API or former M2/M5 clusters that have been migrated to Flex clusters, using `instanceSizeName` to “M2” or “M5” until January 2026. This functionality will be available until January 22, 2026, after which it will only be available for M0 clusters. Please use the Upgrade Flex Cluster endpoint instead.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return UpgradeTenantUpgradeApiRequest
*/
func (a *ClustersApiService) UpgradeTenantUpgrade(ctx context.Context, groupId string, legacyAtlasTenantClusterUpgradeRequest *LegacyAtlasTenantClusterUpgradeRequest) UpgradeTenantUpgradeApiRequest {
	return UpgradeTenantUpgradeApiRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		groupId:                                groupId,
		legacyAtlasTenantClusterUpgradeRequest: legacyAtlasTenantClusterUpgradeRequest,
	}
}

// UpgradeTenantUpgradeExecute executes the request
//
//	@return LegacyAtlasCluster
func (a *ClustersApiService) UpgradeTenantUpgradeExecute(r UpgradeTenantUpgradeApiRequest) (*LegacyAtlasCluster, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *LegacyAtlasCluster
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClustersApiService.UpgradeTenantUpgrade")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/tenantUpgrade"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.legacyAtlasTenantClusterUpgradeRequest == nil {
		return localVarReturnValue, nil, reportError("legacyAtlasTenantClusterUpgradeRequest is required and must be specified")
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
	localVarPostBody = r.legacyAtlasTenantClusterUpgradeRequest
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
