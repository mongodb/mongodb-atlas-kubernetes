// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AtlasSearchApi interface {

	/*
		CreateClusterFtsIndex Create One Atlas Search Index

		Creates one Atlas Search index on the specified collection. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index. Only clusters running MongoDB v4.2 or later can use Atlas Search.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection on which to create an Atlas Search index.
		@param clusterSearchIndex Creates one Atlas Search index on the specified collection.
		@return CreateClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	CreateClusterFtsIndex(ctx context.Context, groupId string, clusterName string, clusterSearchIndex *ClusterSearchIndex) CreateClusterFtsIndexApiRequest
	/*
		CreateClusterFtsIndex Create One Atlas Search Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateClusterFtsIndexApiParams - Parameters for the request
		@return CreateClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	CreateClusterFtsIndexWithParams(ctx context.Context, args *CreateClusterFtsIndexApiParams) CreateClusterFtsIndexApiRequest

	// Method available only for mocking purposes
	CreateClusterFtsIndexExecute(r CreateClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error)

	/*
		CreateClusterSearchDeployment Create Search Nodes

		Creates Search Nodes for the specified cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Label that identifies the cluster to create Search Nodes for.
		@param apiSearchDeploymentRequest Creates Search Nodes for the specified cluster.
		@return CreateClusterSearchDeploymentApiRequest
	*/
	CreateClusterSearchDeployment(ctx context.Context, groupId string, clusterName string, apiSearchDeploymentRequest *ApiSearchDeploymentRequest) CreateClusterSearchDeploymentApiRequest
	/*
		CreateClusterSearchDeployment Create Search Nodes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateClusterSearchDeploymentApiParams - Parameters for the request
		@return CreateClusterSearchDeploymentApiRequest
	*/
	CreateClusterSearchDeploymentWithParams(ctx context.Context, args *CreateClusterSearchDeploymentApiParams) CreateClusterSearchDeploymentApiRequest

	// Method available only for mocking purposes
	CreateClusterSearchDeploymentExecute(r CreateClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error)

	/*
		CreateClusterSearchIndex Create One Atlas Search Index

		Creates one Atlas Search index on the specified collection. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index. Only clusters running MongoDB v4.2 or later can use Atlas Search.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection on which to create an Atlas Search index.
		@param searchIndexCreateRequest Creates one Atlas Search index on the specified collection.
		@return CreateClusterSearchIndexApiRequest
	*/
	CreateClusterSearchIndex(ctx context.Context, groupId string, clusterName string, searchIndexCreateRequest *SearchIndexCreateRequest) CreateClusterSearchIndexApiRequest
	/*
		CreateClusterSearchIndex Create One Atlas Search Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateClusterSearchIndexApiParams - Parameters for the request
		@return CreateClusterSearchIndexApiRequest
	*/
	CreateClusterSearchIndexWithParams(ctx context.Context, args *CreateClusterSearchIndexApiParams) CreateClusterSearchIndexApiRequest

	// Method available only for mocking purposes
	CreateClusterSearchIndexExecute(r CreateClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error)

	/*
		DeleteClusterFtsIndex Remove One Atlas Search Index

		Removes one Atlas Search index that you identified with its unique ID.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
		@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search index. Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
		@return DeleteClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	DeleteClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string) DeleteClusterFtsIndexApiRequest
	/*
		DeleteClusterFtsIndex Remove One Atlas Search Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterFtsIndexApiParams - Parameters for the request
		@return DeleteClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	DeleteClusterFtsIndexWithParams(ctx context.Context, args *DeleteClusterFtsIndexApiParams) DeleteClusterFtsIndexApiRequest

	// Method available only for mocking purposes
	DeleteClusterFtsIndexExecute(r DeleteClusterFtsIndexApiRequest) (*http.Response, error)

	/*
		DeleteClusterSearchDeployment Delete Search Nodes

		Deletes the Search Nodes for the specified cluster.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Label that identifies the cluster to delete.
		@return DeleteClusterSearchDeploymentApiRequest
	*/
	DeleteClusterSearchDeployment(ctx context.Context, groupId string, clusterName string) DeleteClusterSearchDeploymentApiRequest
	/*
		DeleteClusterSearchDeployment Delete Search Nodes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterSearchDeploymentApiParams - Parameters for the request
		@return DeleteClusterSearchDeploymentApiRequest
	*/
	DeleteClusterSearchDeploymentWithParams(ctx context.Context, args *DeleteClusterSearchDeploymentApiParams) DeleteClusterSearchDeploymentApiRequest

	// Method available only for mocking purposes
	DeleteClusterSearchDeploymentExecute(r DeleteClusterSearchDeploymentApiRequest) (*http.Response, error)

	/*
		DeleteClusterSearchIndex Remove One Atlas Search Index by ID

		Removes one Atlas Search index that you identified with its unique ID. This deletion is eventually consistent.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
		@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search index. Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
		@return DeleteClusterSearchIndexApiRequest
	*/
	DeleteClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string) DeleteClusterSearchIndexApiRequest
	/*
		DeleteClusterSearchIndex Remove One Atlas Search Index by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteClusterSearchIndexApiParams - Parameters for the request
		@return DeleteClusterSearchIndexApiRequest
	*/
	DeleteClusterSearchIndexWithParams(ctx context.Context, args *DeleteClusterSearchIndexApiParams) DeleteClusterSearchIndexApiRequest

	// Method available only for mocking purposes
	DeleteClusterSearchIndexExecute(r DeleteClusterSearchIndexApiRequest) (*http.Response, error)

	/*
		DeleteIndexByName Remove One Atlas Search Index by Name

		Removes one Atlas Search index that you identified with its database, collection, and name. This deletion is eventually consistent.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
		@param collectionName Name of the collection that contains one or more Atlas Search indexes.
		@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
		@param indexName Name of the Atlas Search index to delete.
		@return DeleteIndexByNameApiRequest
	*/
	DeleteIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string) DeleteIndexByNameApiRequest
	/*
		DeleteIndexByName Remove One Atlas Search Index by Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteIndexByNameApiParams - Parameters for the request
		@return DeleteIndexByNameApiRequest
	*/
	DeleteIndexByNameWithParams(ctx context.Context, args *DeleteIndexByNameApiParams) DeleteIndexByNameApiRequest

	// Method available only for mocking purposes
	DeleteIndexByNameExecute(r DeleteIndexByNameApiRequest) (*http.Response, error)

	/*
		GetClusterFtsIndex Return One Atlas Search Index

		Returns one Atlas Search index in the specified project. You identify this index using its unique ID. Atlas Search index contains the indexed fields and the analyzers used to create the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@param indexId Unique 24-hexadecimal digit string that identifies the Application Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Application Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Application Search indexes.
		@return GetClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	GetClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string) GetClusterFtsIndexApiRequest
	/*
		GetClusterFtsIndex Return One Atlas Search Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterFtsIndexApiParams - Parameters for the request
		@return GetClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	GetClusterFtsIndexWithParams(ctx context.Context, args *GetClusterFtsIndexApiParams) GetClusterFtsIndexApiRequest

	// Method available only for mocking purposes
	GetClusterFtsIndexExecute(r GetClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error)

	/*
		GetClusterSearchDeployment Return All Search Nodes

		Returns the Search Nodes for the specified cluster. Deprecated versions: v2-{2024-05-30}, v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Label that identifies the cluster to return the Search Nodes for.
		@return GetClusterSearchDeploymentApiRequest
	*/
	GetClusterSearchDeployment(ctx context.Context, groupId string, clusterName string) GetClusterSearchDeploymentApiRequest
	/*
		GetClusterSearchDeployment Return All Search Nodes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterSearchDeploymentApiParams - Parameters for the request
		@return GetClusterSearchDeploymentApiRequest
	*/
	GetClusterSearchDeploymentWithParams(ctx context.Context, args *GetClusterSearchDeploymentApiParams) GetClusterSearchDeploymentApiRequest

	// Method available only for mocking purposes
	GetClusterSearchDeploymentExecute(r GetClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error)

	/*
		GetClusterSearchIndex Return One Atlas Search Index by ID

		Returns one Atlas Search index in the specified project. You identify this index using its unique ID. Atlas Search index contains the indexed fields and the analyzers used to create the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@param indexId Unique 24-hexadecimal digit string that identifies the Application Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Application Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Application Search indexes.
		@return GetClusterSearchIndexApiRequest
	*/
	GetClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string) GetClusterSearchIndexApiRequest
	/*
		GetClusterSearchIndex Return One Atlas Search Index by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetClusterSearchIndexApiParams - Parameters for the request
		@return GetClusterSearchIndexApiRequest
	*/
	GetClusterSearchIndexWithParams(ctx context.Context, args *GetClusterSearchIndexApiParams) GetClusterSearchIndexApiRequest

	// Method available only for mocking purposes
	GetClusterSearchIndexExecute(r GetClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error)

	/*
		GetIndexByName Return One Atlas Search Index by Name

		Returns one Atlas Search index in the specified project. You identify this index using its database, collection and name. Atlas Search index contains the indexed fields and the analyzers used to create the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@param collectionName Name of the collection that contains one or more Atlas Search indexes.
		@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
		@param indexName Name of the Atlas Search index to return.
		@return GetIndexByNameApiRequest
	*/
	GetIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string) GetIndexByNameApiRequest
	/*
		GetIndexByName Return One Atlas Search Index by Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetIndexByNameApiParams - Parameters for the request
		@return GetIndexByNameApiRequest
	*/
	GetIndexByNameWithParams(ctx context.Context, args *GetIndexByNameApiParams) GetIndexByNameApiRequest

	// Method available only for mocking purposes
	GetIndexByNameExecute(r GetIndexByNameApiRequest) (*SearchIndexResponse, *http.Response, error)

	/*
		ListClusterFtsIndex Return All Atlas Search Indexes for One Collection

		Returns all Atlas Search indexes on the specified collection. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@param collectionName Name of the collection that contains one or more Atlas Search indexes.
		@param databaseName Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.
		@return ListClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	ListClusterFtsIndex(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string) ListClusterFtsIndexApiRequest
	/*
		ListClusterFtsIndex Return All Atlas Search Indexes for One Collection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterFtsIndexApiParams - Parameters for the request
		@return ListClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	ListClusterFtsIndexWithParams(ctx context.Context, args *ListClusterFtsIndexApiParams) ListClusterFtsIndexApiRequest

	// Method available only for mocking purposes
	ListClusterFtsIndexExecute(r ListClusterFtsIndexApiRequest) ([]ClusterSearchIndex, *http.Response, error)

	/*
		ListClusterSearchIndexes Return All Atlas Search Indexes for One Cluster

		Returns all Atlas Search indexes on the specified cluster. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@return ListClusterSearchIndexesApiRequest
	*/
	ListClusterSearchIndexes(ctx context.Context, groupId string, clusterName string) ListClusterSearchIndexesApiRequest
	/*
		ListClusterSearchIndexes Return All Atlas Search Indexes for One Cluster


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListClusterSearchIndexesApiParams - Parameters for the request
		@return ListClusterSearchIndexesApiRequest
	*/
	ListClusterSearchIndexesWithParams(ctx context.Context, args *ListClusterSearchIndexesApiParams) ListClusterSearchIndexesApiRequest

	// Method available only for mocking purposes
	ListClusterSearchIndexesExecute(r ListClusterSearchIndexesApiRequest) ([]SearchIndexResponse, *http.Response, error)

	/*
		ListSearchIndex Return All Atlas Search Indexes for One Collection

		Returns all Atlas Search indexes on the specified collection. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
		@param collectionName Name of the collection that contains one or more Atlas Search indexes.
		@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
		@return ListSearchIndexApiRequest
	*/
	ListSearchIndex(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string) ListSearchIndexApiRequest
	/*
		ListSearchIndex Return All Atlas Search Indexes for One Collection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListSearchIndexApiParams - Parameters for the request
		@return ListSearchIndexApiRequest
	*/
	ListSearchIndexWithParams(ctx context.Context, args *ListSearchIndexApiParams) ListSearchIndexApiRequest

	// Method available only for mocking purposes
	ListSearchIndexExecute(r ListSearchIndexApiRequest) ([]SearchIndexResponse, *http.Response, error)

	/*
		UpdateClusterFtsIndex Update One Atlas Search Index

		Updates one Atlas Search index that you identified with its unique ID. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection whose Atlas Search index to update.
		@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
		@param clusterSearchIndex Details to update on the Atlas Search index.
		@return UpdateClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	UpdateClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string, clusterSearchIndex *ClusterSearchIndex) UpdateClusterFtsIndexApiRequest
	/*
		UpdateClusterFtsIndex Update One Atlas Search Index


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterFtsIndexApiParams - Parameters for the request
		@return UpdateClusterFtsIndexApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for AtlasSearchApi
	*/
	UpdateClusterFtsIndexWithParams(ctx context.Context, args *UpdateClusterFtsIndexApiParams) UpdateClusterFtsIndexApiRequest

	// Method available only for mocking purposes
	UpdateClusterFtsIndexExecute(r UpdateClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error)

	/*
		UpdateClusterSearchDeployment Update Search Nodes

		Updates the Search Nodes for the specified cluster. Deprecated versions: v2-{2023-01-01}

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Label that identifies the cluster to update the Search Nodes for.
		@param apiSearchDeploymentRequest Updates the Search Nodes for the specified cluster.
		@return UpdateClusterSearchDeploymentApiRequest
	*/
	UpdateClusterSearchDeployment(ctx context.Context, groupId string, clusterName string, apiSearchDeploymentRequest *ApiSearchDeploymentRequest) UpdateClusterSearchDeploymentApiRequest
	/*
		UpdateClusterSearchDeployment Update Search Nodes


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterSearchDeploymentApiParams - Parameters for the request
		@return UpdateClusterSearchDeploymentApiRequest
	*/
	UpdateClusterSearchDeploymentWithParams(ctx context.Context, args *UpdateClusterSearchDeploymentApiParams) UpdateClusterSearchDeploymentApiRequest

	// Method available only for mocking purposes
	UpdateClusterSearchDeploymentExecute(r UpdateClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error)

	/*
		UpdateClusterSearchIndex Update One Atlas Search Index by ID

		Updates one Atlas Search index that you identified with its unique ID. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection whose Atlas Search index you want to update.
		@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
		@param searchIndexUpdateRequest Details to update on the Atlas Search index.
		@return UpdateClusterSearchIndexApiRequest
	*/
	UpdateClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string, searchIndexUpdateRequest *SearchIndexUpdateRequest) UpdateClusterSearchIndexApiRequest
	/*
		UpdateClusterSearchIndex Update One Atlas Search Index by ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateClusterSearchIndexApiParams - Parameters for the request
		@return UpdateClusterSearchIndexApiRequest
	*/
	UpdateClusterSearchIndexWithParams(ctx context.Context, args *UpdateClusterSearchIndexApiParams) UpdateClusterSearchIndexApiRequest

	// Method available only for mocking purposes
	UpdateClusterSearchIndexExecute(r UpdateClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error)

	/*
		UpdateIndexByName Update One Atlas Search Index by Name

		Updates one Atlas Search index that you identified with its database, collection name, and index name. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param clusterName Name of the cluster that contains the collection whose Atlas Search index you want to update.
		@param collectionName Name of the collection that contains one or more Atlas Search indexes.
		@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
		@param indexName Name of the Atlas Search index to update.
		@param searchIndexUpdateRequest Details to update the Atlas Search index with.
		@return UpdateIndexByNameApiRequest
	*/
	UpdateIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string, searchIndexUpdateRequest *SearchIndexUpdateRequest) UpdateIndexByNameApiRequest
	/*
		UpdateIndexByName Update One Atlas Search Index by Name


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateIndexByNameApiParams - Parameters for the request
		@return UpdateIndexByNameApiRequest
	*/
	UpdateIndexByNameWithParams(ctx context.Context, args *UpdateIndexByNameApiParams) UpdateIndexByNameApiRequest

	// Method available only for mocking purposes
	UpdateIndexByNameExecute(r UpdateIndexByNameApiRequest) (*SearchIndexResponse, *http.Response, error)
}

// AtlasSearchApiService AtlasSearchApi service
type AtlasSearchApiService service

type CreateClusterFtsIndexApiRequest struct {
	ctx                context.Context
	ApiService         AtlasSearchApi
	groupId            string
	clusterName        string
	clusterSearchIndex *ClusterSearchIndex
}

type CreateClusterFtsIndexApiParams struct {
	GroupId            string
	ClusterName        string
	ClusterSearchIndex *ClusterSearchIndex
}

func (a *AtlasSearchApiService) CreateClusterFtsIndexWithParams(ctx context.Context, args *CreateClusterFtsIndexApiParams) CreateClusterFtsIndexApiRequest {
	return CreateClusterFtsIndexApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		clusterName:        args.ClusterName,
		clusterSearchIndex: args.ClusterSearchIndex,
	}
}

func (r CreateClusterFtsIndexApiRequest) Execute() (*ClusterSearchIndex, *http.Response, error) {
	return r.ApiService.CreateClusterFtsIndexExecute(r)
}

/*
CreateClusterFtsIndex Create One Atlas Search Index

Creates one Atlas Search index on the specified collection. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index. Only clusters running MongoDB v4.2 or later can use Atlas Search.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection on which to create an Atlas Search index.
	@return CreateClusterFtsIndexApiRequest

Deprecated
*/
func (a *AtlasSearchApiService) CreateClusterFtsIndex(ctx context.Context, groupId string, clusterName string, clusterSearchIndex *ClusterSearchIndex) CreateClusterFtsIndexApiRequest {
	return CreateClusterFtsIndexApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            groupId,
		clusterName:        clusterName,
		clusterSearchIndex: clusterSearchIndex,
	}
}

// CreateClusterFtsIndexExecute executes the request
//
//	@return ClusterSearchIndex
//
// Deprecated
func (a *AtlasSearchApiService) CreateClusterFtsIndexExecute(r CreateClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterSearchIndex
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.CreateClusterFtsIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes"
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
	if r.clusterSearchIndex == nil {
		return localVarReturnValue, nil, reportError("clusterSearchIndex is required and must be specified")
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
	localVarPostBody = r.clusterSearchIndex
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

type CreateClusterSearchDeploymentApiRequest struct {
	ctx                        context.Context
	ApiService                 AtlasSearchApi
	groupId                    string
	clusterName                string
	apiSearchDeploymentRequest *ApiSearchDeploymentRequest
}

type CreateClusterSearchDeploymentApiParams struct {
	GroupId                    string
	ClusterName                string
	ApiSearchDeploymentRequest *ApiSearchDeploymentRequest
}

func (a *AtlasSearchApiService) CreateClusterSearchDeploymentWithParams(ctx context.Context, args *CreateClusterSearchDeploymentApiParams) CreateClusterSearchDeploymentApiRequest {
	return CreateClusterSearchDeploymentApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		clusterName:                args.ClusterName,
		apiSearchDeploymentRequest: args.ApiSearchDeploymentRequest,
	}
}

func (r CreateClusterSearchDeploymentApiRequest) Execute() (*ApiSearchDeploymentResponse, *http.Response, error) {
	return r.ApiService.CreateClusterSearchDeploymentExecute(r)
}

/*
CreateClusterSearchDeployment Create Search Nodes

Creates Search Nodes for the specified cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Label that identifies the cluster to create Search Nodes for.
	@return CreateClusterSearchDeploymentApiRequest
*/
func (a *AtlasSearchApiService) CreateClusterSearchDeployment(ctx context.Context, groupId string, clusterName string, apiSearchDeploymentRequest *ApiSearchDeploymentRequest) CreateClusterSearchDeploymentApiRequest {
	return CreateClusterSearchDeploymentApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterName:                clusterName,
		apiSearchDeploymentRequest: apiSearchDeploymentRequest,
	}
}

// CreateClusterSearchDeploymentExecute executes the request
//
//	@return ApiSearchDeploymentResponse
func (a *AtlasSearchApiService) CreateClusterSearchDeploymentExecute(r CreateClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiSearchDeploymentResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.CreateClusterSearchDeployment")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/deployment"
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
	if r.apiSearchDeploymentRequest == nil {
		return localVarReturnValue, nil, reportError("apiSearchDeploymentRequest is required and must be specified")
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
	localVarPostBody = r.apiSearchDeploymentRequest
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

type CreateClusterSearchIndexApiRequest struct {
	ctx                      context.Context
	ApiService               AtlasSearchApi
	groupId                  string
	clusterName              string
	searchIndexCreateRequest *SearchIndexCreateRequest
}

type CreateClusterSearchIndexApiParams struct {
	GroupId                  string
	ClusterName              string
	SearchIndexCreateRequest *SearchIndexCreateRequest
}

func (a *AtlasSearchApiService) CreateClusterSearchIndexWithParams(ctx context.Context, args *CreateClusterSearchIndexApiParams) CreateClusterSearchIndexApiRequest {
	return CreateClusterSearchIndexApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		clusterName:              args.ClusterName,
		searchIndexCreateRequest: args.SearchIndexCreateRequest,
	}
}

func (r CreateClusterSearchIndexApiRequest) Execute() (*SearchIndexResponse, *http.Response, error) {
	return r.ApiService.CreateClusterSearchIndexExecute(r)
}

/*
CreateClusterSearchIndex Create One Atlas Search Index

Creates one Atlas Search index on the specified collection. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index. Only clusters running MongoDB v4.2 or later can use Atlas Search.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection on which to create an Atlas Search index.
	@return CreateClusterSearchIndexApiRequest
*/
func (a *AtlasSearchApiService) CreateClusterSearchIndex(ctx context.Context, groupId string, clusterName string, searchIndexCreateRequest *SearchIndexCreateRequest) CreateClusterSearchIndexApiRequest {
	return CreateClusterSearchIndexApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		clusterName:              clusterName,
		searchIndexCreateRequest: searchIndexCreateRequest,
	}
}

// CreateClusterSearchIndexExecute executes the request
//
//	@return SearchIndexResponse
func (a *AtlasSearchApiService) CreateClusterSearchIndexExecute(r CreateClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.CreateClusterSearchIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes"
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
	if r.searchIndexCreateRequest == nil {
		return localVarReturnValue, nil, reportError("searchIndexCreateRequest is required and must be specified")
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
	localVarPostBody = r.searchIndexCreateRequest
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

type DeleteClusterFtsIndexApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
	indexId     string
}

type DeleteClusterFtsIndexApiParams struct {
	GroupId     string
	ClusterName string
	IndexId     string
}

func (a *AtlasSearchApiService) DeleteClusterFtsIndexWithParams(ctx context.Context, args *DeleteClusterFtsIndexApiParams) DeleteClusterFtsIndexApiRequest {
	return DeleteClusterFtsIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		indexId:     args.IndexId,
	}
}

func (r DeleteClusterFtsIndexApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterFtsIndexExecute(r)
}

/*
DeleteClusterFtsIndex Remove One Atlas Search Index

Removes one Atlas Search index that you identified with its unique ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
	@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search index. Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
	@return DeleteClusterFtsIndexApiRequest

Deprecated
*/
func (a *AtlasSearchApiService) DeleteClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string) DeleteClusterFtsIndexApiRequest {
	return DeleteClusterFtsIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		indexId:     indexId,
	}
}

// DeleteClusterFtsIndexExecute executes the request
// Deprecated
func (a *AtlasSearchApiService) DeleteClusterFtsIndexExecute(r DeleteClusterFtsIndexApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.DeleteClusterFtsIndex")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes/{indexId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

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

type DeleteClusterSearchDeploymentApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
}

type DeleteClusterSearchDeploymentApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *AtlasSearchApiService) DeleteClusterSearchDeploymentWithParams(ctx context.Context, args *DeleteClusterSearchDeploymentApiParams) DeleteClusterSearchDeploymentApiRequest {
	return DeleteClusterSearchDeploymentApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r DeleteClusterSearchDeploymentApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterSearchDeploymentExecute(r)
}

/*
DeleteClusterSearchDeployment Delete Search Nodes

Deletes the Search Nodes for the specified cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Label that identifies the cluster to delete.
	@return DeleteClusterSearchDeploymentApiRequest
*/
func (a *AtlasSearchApiService) DeleteClusterSearchDeployment(ctx context.Context, groupId string, clusterName string) DeleteClusterSearchDeploymentApiRequest {
	return DeleteClusterSearchDeploymentApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// DeleteClusterSearchDeploymentExecute executes the request
func (a *AtlasSearchApiService) DeleteClusterSearchDeploymentExecute(r DeleteClusterSearchDeploymentApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.DeleteClusterSearchDeployment")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/deployment"
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

type DeleteClusterSearchIndexApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
	indexId     string
}

type DeleteClusterSearchIndexApiParams struct {
	GroupId     string
	ClusterName string
	IndexId     string
}

func (a *AtlasSearchApiService) DeleteClusterSearchIndexWithParams(ctx context.Context, args *DeleteClusterSearchIndexApiParams) DeleteClusterSearchIndexApiRequest {
	return DeleteClusterSearchIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		indexId:     args.IndexId,
	}
}

func (r DeleteClusterSearchIndexApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClusterSearchIndexExecute(r)
}

/*
DeleteClusterSearchIndex Remove One Atlas Search Index by ID

Removes one Atlas Search index that you identified with its unique ID. This deletion is eventually consistent.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
	@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search index. Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
	@return DeleteClusterSearchIndexApiRequest
*/
func (a *AtlasSearchApiService) DeleteClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string) DeleteClusterSearchIndexApiRequest {
	return DeleteClusterSearchIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		indexId:     indexId,
	}
}

// DeleteClusterSearchIndexExecute executes the request
func (a *AtlasSearchApiService) DeleteClusterSearchIndexExecute(r DeleteClusterSearchIndexApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.DeleteClusterSearchIndex")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{indexId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

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

type DeleteIndexByNameApiRequest struct {
	ctx            context.Context
	ApiService     AtlasSearchApi
	groupId        string
	clusterName    string
	collectionName string
	databaseName   string
	indexName      string
}

type DeleteIndexByNameApiParams struct {
	GroupId        string
	ClusterName    string
	CollectionName string
	DatabaseName   string
	IndexName      string
}

func (a *AtlasSearchApiService) DeleteIndexByNameWithParams(ctx context.Context, args *DeleteIndexByNameApiParams) DeleteIndexByNameApiRequest {
	return DeleteIndexByNameApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		collectionName: args.CollectionName,
		databaseName:   args.DatabaseName,
		indexName:      args.IndexName,
	}
}

func (r DeleteIndexByNameApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteIndexByNameExecute(r)
}

/*
DeleteIndexByName Remove One Atlas Search Index by Name

Removes one Atlas Search index that you identified with its database, collection, and name. This deletion is eventually consistent.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the database and collection with one or more Application Search indexes.
	@param collectionName Name of the collection that contains one or more Atlas Search indexes.
	@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
	@param indexName Name of the Atlas Search index to delete.
	@return DeleteIndexByNameApiRequest
*/
func (a *AtlasSearchApiService) DeleteIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string) DeleteIndexByNameApiRequest {
	return DeleteIndexByNameApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		collectionName: collectionName,
		databaseName:   databaseName,
		indexName:      indexName,
	}
}

// DeleteIndexByNameExecute executes the request
func (a *AtlasSearchApiService) DeleteIndexByNameExecute(r DeleteIndexByNameApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.DeleteIndexByName")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{databaseName}/{collectionName}/{indexName}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.collectionName == "" {
		return nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.databaseName == "" {
		return nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.indexName == "" {
		return nil, reportError("indexName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexName"+"}", url.PathEscape(r.indexName), -1)

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

type GetClusterFtsIndexApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
	indexId     string
}

type GetClusterFtsIndexApiParams struct {
	GroupId     string
	ClusterName string
	IndexId     string
}

func (a *AtlasSearchApiService) GetClusterFtsIndexWithParams(ctx context.Context, args *GetClusterFtsIndexApiParams) GetClusterFtsIndexApiRequest {
	return GetClusterFtsIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		indexId:     args.IndexId,
	}
}

func (r GetClusterFtsIndexApiRequest) Execute() (*ClusterSearchIndex, *http.Response, error) {
	return r.ApiService.GetClusterFtsIndexExecute(r)
}

/*
GetClusterFtsIndex Return One Atlas Search Index

Returns one Atlas Search index in the specified project. You identify this index using its unique ID. Atlas Search index contains the indexed fields and the analyzers used to create the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@param indexId Unique 24-hexadecimal digit string that identifies the Application Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Application Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Application Search indexes.
	@return GetClusterFtsIndexApiRequest

Deprecated
*/
func (a *AtlasSearchApiService) GetClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string) GetClusterFtsIndexApiRequest {
	return GetClusterFtsIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		indexId:     indexId,
	}
}

// GetClusterFtsIndexExecute executes the request
//
//	@return ClusterSearchIndex
//
// Deprecated
func (a *AtlasSearchApiService) GetClusterFtsIndexExecute(r GetClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterSearchIndex
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.GetClusterFtsIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes/{indexId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return localVarReturnValue, nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

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

type GetClusterSearchDeploymentApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
}

type GetClusterSearchDeploymentApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *AtlasSearchApiService) GetClusterSearchDeploymentWithParams(ctx context.Context, args *GetClusterSearchDeploymentApiParams) GetClusterSearchDeploymentApiRequest {
	return GetClusterSearchDeploymentApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r GetClusterSearchDeploymentApiRequest) Execute() (*ApiSearchDeploymentResponse, *http.Response, error) {
	return r.ApiService.GetClusterSearchDeploymentExecute(r)
}

/*
GetClusterSearchDeployment Return All Search Nodes

Returns the Search Nodes for the specified cluster. Deprecated versions: v2-{2024-05-30}, v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Label that identifies the cluster to return the Search Nodes for.
	@return GetClusterSearchDeploymentApiRequest
*/
func (a *AtlasSearchApiService) GetClusterSearchDeployment(ctx context.Context, groupId string, clusterName string) GetClusterSearchDeploymentApiRequest {
	return GetClusterSearchDeploymentApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// GetClusterSearchDeploymentExecute executes the request
//
//	@return ApiSearchDeploymentResponse
func (a *AtlasSearchApiService) GetClusterSearchDeploymentExecute(r GetClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiSearchDeploymentResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.GetClusterSearchDeployment")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/deployment"
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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

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

type GetClusterSearchIndexApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
	indexId     string
}

type GetClusterSearchIndexApiParams struct {
	GroupId     string
	ClusterName string
	IndexId     string
}

func (a *AtlasSearchApiService) GetClusterSearchIndexWithParams(ctx context.Context, args *GetClusterSearchIndexApiParams) GetClusterSearchIndexApiRequest {
	return GetClusterSearchIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
		indexId:     args.IndexId,
	}
}

func (r GetClusterSearchIndexApiRequest) Execute() (*SearchIndexResponse, *http.Response, error) {
	return r.ApiService.GetClusterSearchIndexExecute(r)
}

/*
GetClusterSearchIndex Return One Atlas Search Index by ID

Returns one Atlas Search index in the specified project. You identify this index using its unique ID. Atlas Search index contains the indexed fields and the analyzers used to create the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@param indexId Unique 24-hexadecimal digit string that identifies the Application Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Application Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Application Search indexes.
	@return GetClusterSearchIndexApiRequest
*/
func (a *AtlasSearchApiService) GetClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string) GetClusterSearchIndexApiRequest {
	return GetClusterSearchIndexApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
		indexId:     indexId,
	}
}

// GetClusterSearchIndexExecute executes the request
//
//	@return SearchIndexResponse
func (a *AtlasSearchApiService) GetClusterSearchIndexExecute(r GetClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.GetClusterSearchIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{indexId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return localVarReturnValue, nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

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

type GetIndexByNameApiRequest struct {
	ctx            context.Context
	ApiService     AtlasSearchApi
	groupId        string
	clusterName    string
	collectionName string
	databaseName   string
	indexName      string
}

type GetIndexByNameApiParams struct {
	GroupId        string
	ClusterName    string
	CollectionName string
	DatabaseName   string
	IndexName      string
}

func (a *AtlasSearchApiService) GetIndexByNameWithParams(ctx context.Context, args *GetIndexByNameApiParams) GetIndexByNameApiRequest {
	return GetIndexByNameApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		collectionName: args.CollectionName,
		databaseName:   args.DatabaseName,
		indexName:      args.IndexName,
	}
}

func (r GetIndexByNameApiRequest) Execute() (*SearchIndexResponse, *http.Response, error) {
	return r.ApiService.GetIndexByNameExecute(r)
}

/*
GetIndexByName Return One Atlas Search Index by Name

Returns one Atlas Search index in the specified project. You identify this index using its database, collection and name. Atlas Search index contains the indexed fields and the analyzers used to create the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@param collectionName Name of the collection that contains one or more Atlas Search indexes.
	@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
	@param indexName Name of the Atlas Search index to return.
	@return GetIndexByNameApiRequest
*/
func (a *AtlasSearchApiService) GetIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string) GetIndexByNameApiRequest {
	return GetIndexByNameApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		collectionName: collectionName,
		databaseName:   databaseName,
		indexName:      indexName,
	}
}

// GetIndexByNameExecute executes the request
//
//	@return SearchIndexResponse
func (a *AtlasSearchApiService) GetIndexByNameExecute(r GetIndexByNameApiRequest) (*SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.GetIndexByName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{databaseName}/{collectionName}/{indexName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.indexName == "" {
		return localVarReturnValue, nil, reportError("indexName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexName"+"}", url.PathEscape(r.indexName), -1)

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

type ListClusterFtsIndexApiRequest struct {
	ctx            context.Context
	ApiService     AtlasSearchApi
	groupId        string
	clusterName    string
	collectionName string
	databaseName   string
}

type ListClusterFtsIndexApiParams struct {
	GroupId        string
	ClusterName    string
	CollectionName string
	DatabaseName   string
}

func (a *AtlasSearchApiService) ListClusterFtsIndexWithParams(ctx context.Context, args *ListClusterFtsIndexApiParams) ListClusterFtsIndexApiRequest {
	return ListClusterFtsIndexApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		collectionName: args.CollectionName,
		databaseName:   args.DatabaseName,
	}
}

func (r ListClusterFtsIndexApiRequest) Execute() ([]ClusterSearchIndex, *http.Response, error) {
	return r.ApiService.ListClusterFtsIndexExecute(r)
}

/*
ListClusterFtsIndex Return All Atlas Search Indexes for One Collection

Returns all Atlas Search indexes on the specified collection. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@param collectionName Name of the collection that contains one or more Atlas Search indexes.
	@param databaseName Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.
	@return ListClusterFtsIndexApiRequest

Deprecated
*/
func (a *AtlasSearchApiService) ListClusterFtsIndex(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string) ListClusterFtsIndexApiRequest {
	return ListClusterFtsIndexApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		collectionName: collectionName,
		databaseName:   databaseName,
	}
}

// ListClusterFtsIndexExecute executes the request
//
//	@return []ClusterSearchIndex
//
// Deprecated
func (a *AtlasSearchApiService) ListClusterFtsIndexExecute(r ListClusterFtsIndexApiRequest) ([]ClusterSearchIndex, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []ClusterSearchIndex
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.ListClusterFtsIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes/{databaseName}/{collectionName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)

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

type ListClusterSearchIndexesApiRequest struct {
	ctx         context.Context
	ApiService  AtlasSearchApi
	groupId     string
	clusterName string
}

type ListClusterSearchIndexesApiParams struct {
	GroupId     string
	ClusterName string
}

func (a *AtlasSearchApiService) ListClusterSearchIndexesWithParams(ctx context.Context, args *ListClusterSearchIndexesApiParams) ListClusterSearchIndexesApiRequest {
	return ListClusterSearchIndexesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		clusterName: args.ClusterName,
	}
}

func (r ListClusterSearchIndexesApiRequest) Execute() ([]SearchIndexResponse, *http.Response, error) {
	return r.ApiService.ListClusterSearchIndexesExecute(r)
}

/*
ListClusterSearchIndexes Return All Atlas Search Indexes for One Cluster

Returns all Atlas Search indexes on the specified cluster. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@return ListClusterSearchIndexesApiRequest
*/
func (a *AtlasSearchApiService) ListClusterSearchIndexes(ctx context.Context, groupId string, clusterName string) ListClusterSearchIndexesApiRequest {
	return ListClusterSearchIndexesApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		clusterName: clusterName,
	}
}

// ListClusterSearchIndexesExecute executes the request
//
//	@return []SearchIndexResponse
func (a *AtlasSearchApiService) ListClusterSearchIndexesExecute(r ListClusterSearchIndexesApiRequest) ([]SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.ListClusterSearchIndexes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes"
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

type ListSearchIndexApiRequest struct {
	ctx            context.Context
	ApiService     AtlasSearchApi
	groupId        string
	clusterName    string
	collectionName string
	databaseName   string
}

type ListSearchIndexApiParams struct {
	GroupId        string
	ClusterName    string
	CollectionName string
	DatabaseName   string
}

func (a *AtlasSearchApiService) ListSearchIndexWithParams(ctx context.Context, args *ListSearchIndexApiParams) ListSearchIndexApiRequest {
	return ListSearchIndexApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        args.GroupId,
		clusterName:    args.ClusterName,
		collectionName: args.CollectionName,
		databaseName:   args.DatabaseName,
	}
}

func (r ListSearchIndexApiRequest) Execute() ([]SearchIndexResponse, *http.Response, error) {
	return r.ApiService.ListSearchIndexExecute(r)
}

/*
ListSearchIndex Return All Atlas Search Indexes for One Collection

Returns all Atlas Search indexes on the specified collection. Atlas Search indexes contain the indexed fields and the analyzers used to create the indexes.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection with one or more Atlas Search indexes.
	@param collectionName Name of the collection that contains one or more Atlas Search indexes.
	@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
	@return ListSearchIndexApiRequest
*/
func (a *AtlasSearchApiService) ListSearchIndex(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string) ListSearchIndexApiRequest {
	return ListSearchIndexApiRequest{
		ApiService:     a,
		ctx:            ctx,
		groupId:        groupId,
		clusterName:    clusterName,
		collectionName: collectionName,
		databaseName:   databaseName,
	}
}

// ListSearchIndexExecute executes the request
//
//	@return []SearchIndexResponse
func (a *AtlasSearchApiService) ListSearchIndexExecute(r ListSearchIndexApiRequest) ([]SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.ListSearchIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{databaseName}/{collectionName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)

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

type UpdateClusterFtsIndexApiRequest struct {
	ctx                context.Context
	ApiService         AtlasSearchApi
	groupId            string
	clusterName        string
	indexId            string
	clusterSearchIndex *ClusterSearchIndex
}

type UpdateClusterFtsIndexApiParams struct {
	GroupId            string
	ClusterName        string
	IndexId            string
	ClusterSearchIndex *ClusterSearchIndex
}

func (a *AtlasSearchApiService) UpdateClusterFtsIndexWithParams(ctx context.Context, args *UpdateClusterFtsIndexApiParams) UpdateClusterFtsIndexApiRequest {
	return UpdateClusterFtsIndexApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            args.GroupId,
		clusterName:        args.ClusterName,
		indexId:            args.IndexId,
		clusterSearchIndex: args.ClusterSearchIndex,
	}
}

func (r UpdateClusterFtsIndexApiRequest) Execute() (*ClusterSearchIndex, *http.Response, error) {
	return r.ApiService.UpdateClusterFtsIndexExecute(r)
}

/*
UpdateClusterFtsIndex Update One Atlas Search Index

Updates one Atlas Search index that you identified with its unique ID. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection whose Atlas Search index to update.
	@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
	@return UpdateClusterFtsIndexApiRequest

Deprecated
*/
func (a *AtlasSearchApiService) UpdateClusterFtsIndex(ctx context.Context, groupId string, clusterName string, indexId string, clusterSearchIndex *ClusterSearchIndex) UpdateClusterFtsIndexApiRequest {
	return UpdateClusterFtsIndexApiRequest{
		ApiService:         a,
		ctx:                ctx,
		groupId:            groupId,
		clusterName:        clusterName,
		indexId:            indexId,
		clusterSearchIndex: clusterSearchIndex,
	}
}

// UpdateClusterFtsIndexExecute executes the request
//
//	@return ClusterSearchIndex
//
// Deprecated
func (a *AtlasSearchApiService) UpdateClusterFtsIndexExecute(r UpdateClusterFtsIndexApiRequest) (*ClusterSearchIndex, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ClusterSearchIndex
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.UpdateClusterFtsIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes/{indexId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return localVarReturnValue, nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.clusterSearchIndex == nil {
		return localVarReturnValue, nil, reportError("clusterSearchIndex is required and must be specified")
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
	localVarPostBody = r.clusterSearchIndex
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

type UpdateClusterSearchDeploymentApiRequest struct {
	ctx                        context.Context
	ApiService                 AtlasSearchApi
	groupId                    string
	clusterName                string
	apiSearchDeploymentRequest *ApiSearchDeploymentRequest
}

type UpdateClusterSearchDeploymentApiParams struct {
	GroupId                    string
	ClusterName                string
	ApiSearchDeploymentRequest *ApiSearchDeploymentRequest
}

func (a *AtlasSearchApiService) UpdateClusterSearchDeploymentWithParams(ctx context.Context, args *UpdateClusterSearchDeploymentApiParams) UpdateClusterSearchDeploymentApiRequest {
	return UpdateClusterSearchDeploymentApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    args.GroupId,
		clusterName:                args.ClusterName,
		apiSearchDeploymentRequest: args.ApiSearchDeploymentRequest,
	}
}

func (r UpdateClusterSearchDeploymentApiRequest) Execute() (*ApiSearchDeploymentResponse, *http.Response, error) {
	return r.ApiService.UpdateClusterSearchDeploymentExecute(r)
}

/*
UpdateClusterSearchDeployment Update Search Nodes

Updates the Search Nodes for the specified cluster. Deprecated versions: v2-{2023-01-01}

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Label that identifies the cluster to update the Search Nodes for.
	@return UpdateClusterSearchDeploymentApiRequest
*/
func (a *AtlasSearchApiService) UpdateClusterSearchDeployment(ctx context.Context, groupId string, clusterName string, apiSearchDeploymentRequest *ApiSearchDeploymentRequest) UpdateClusterSearchDeploymentApiRequest {
	return UpdateClusterSearchDeploymentApiRequest{
		ApiService:                 a,
		ctx:                        ctx,
		groupId:                    groupId,
		clusterName:                clusterName,
		apiSearchDeploymentRequest: apiSearchDeploymentRequest,
	}
}

// UpdateClusterSearchDeploymentExecute executes the request
//
//	@return ApiSearchDeploymentResponse
func (a *AtlasSearchApiService) UpdateClusterSearchDeploymentExecute(r UpdateClusterSearchDeploymentApiRequest) (*ApiSearchDeploymentResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *ApiSearchDeploymentResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.UpdateClusterSearchDeployment")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/deployment"
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
	if r.apiSearchDeploymentRequest == nil {
		return localVarReturnValue, nil, reportError("apiSearchDeploymentRequest is required and must be specified")
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
	localVarPostBody = r.apiSearchDeploymentRequest
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

type UpdateClusterSearchIndexApiRequest struct {
	ctx                      context.Context
	ApiService               AtlasSearchApi
	groupId                  string
	clusterName              string
	indexId                  string
	searchIndexUpdateRequest *SearchIndexUpdateRequest
}

type UpdateClusterSearchIndexApiParams struct {
	GroupId                  string
	ClusterName              string
	IndexId                  string
	SearchIndexUpdateRequest *SearchIndexUpdateRequest
}

func (a *AtlasSearchApiService) UpdateClusterSearchIndexWithParams(ctx context.Context, args *UpdateClusterSearchIndexApiParams) UpdateClusterSearchIndexApiRequest {
	return UpdateClusterSearchIndexApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		clusterName:              args.ClusterName,
		indexId:                  args.IndexId,
		searchIndexUpdateRequest: args.SearchIndexUpdateRequest,
	}
}

func (r UpdateClusterSearchIndexApiRequest) Execute() (*SearchIndexResponse, *http.Response, error) {
	return r.ApiService.UpdateClusterSearchIndexExecute(r)
}

/*
UpdateClusterSearchIndex Update One Atlas Search Index by ID

Updates one Atlas Search index that you identified with its unique ID. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection whose Atlas Search index you want to update.
	@param indexId Unique 24-hexadecimal digit string that identifies the Atlas Search [index](https://dochub.mongodb.org/core/index-definitions-fts). Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.
	@return UpdateClusterSearchIndexApiRequest
*/
func (a *AtlasSearchApiService) UpdateClusterSearchIndex(ctx context.Context, groupId string, clusterName string, indexId string, searchIndexUpdateRequest *SearchIndexUpdateRequest) UpdateClusterSearchIndexApiRequest {
	return UpdateClusterSearchIndexApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		clusterName:              clusterName,
		indexId:                  indexId,
		searchIndexUpdateRequest: searchIndexUpdateRequest,
	}
}

// UpdateClusterSearchIndexExecute executes the request
//
//	@return SearchIndexResponse
func (a *AtlasSearchApiService) UpdateClusterSearchIndexExecute(r UpdateClusterSearchIndexApiRequest) (*SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.UpdateClusterSearchIndex")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{indexId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.indexId == "" {
		return localVarReturnValue, nil, reportError("indexId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexId"+"}", url.PathEscape(r.indexId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.searchIndexUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("searchIndexUpdateRequest is required and must be specified")
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
	localVarPostBody = r.searchIndexUpdateRequest
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

type UpdateIndexByNameApiRequest struct {
	ctx                      context.Context
	ApiService               AtlasSearchApi
	groupId                  string
	clusterName              string
	collectionName           string
	databaseName             string
	indexName                string
	searchIndexUpdateRequest *SearchIndexUpdateRequest
}

type UpdateIndexByNameApiParams struct {
	GroupId                  string
	ClusterName              string
	CollectionName           string
	DatabaseName             string
	IndexName                string
	SearchIndexUpdateRequest *SearchIndexUpdateRequest
}

func (a *AtlasSearchApiService) UpdateIndexByNameWithParams(ctx context.Context, args *UpdateIndexByNameApiParams) UpdateIndexByNameApiRequest {
	return UpdateIndexByNameApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  args.GroupId,
		clusterName:              args.ClusterName,
		collectionName:           args.CollectionName,
		databaseName:             args.DatabaseName,
		indexName:                args.IndexName,
		searchIndexUpdateRequest: args.SearchIndexUpdateRequest,
	}
}

func (r UpdateIndexByNameApiRequest) Execute() (*SearchIndexResponse, *http.Response, error) {
	return r.ApiService.UpdateIndexByNameExecute(r)
}

/*
UpdateIndexByName Update One Atlas Search Index by Name

Updates one Atlas Search index that you identified with its database, collection name, and index name. Atlas Search indexes define the fields on which to create the index and the analyzers to use when creating the index.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param clusterName Name of the cluster that contains the collection whose Atlas Search index you want to update.
	@param collectionName Name of the collection that contains one or more Atlas Search indexes.
	@param databaseName Label that identifies the database that contains the collection with one or more Atlas Search indexes.
	@param indexName Name of the Atlas Search index to update.
	@return UpdateIndexByNameApiRequest
*/
func (a *AtlasSearchApiService) UpdateIndexByName(ctx context.Context, groupId string, clusterName string, collectionName string, databaseName string, indexName string, searchIndexUpdateRequest *SearchIndexUpdateRequest) UpdateIndexByNameApiRequest {
	return UpdateIndexByNameApiRequest{
		ApiService:               a,
		ctx:                      ctx,
		groupId:                  groupId,
		clusterName:              clusterName,
		collectionName:           collectionName,
		databaseName:             databaseName,
		indexName:                indexName,
		searchIndexUpdateRequest: searchIndexUpdateRequest,
	}
}

// UpdateIndexByNameExecute executes the request
//
//	@return SearchIndexResponse
func (a *AtlasSearchApiService) UpdateIndexByNameExecute(r UpdateIndexByNameApiRequest) (*SearchIndexResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SearchIndexResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AtlasSearchApiService.UpdateIndexByName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{databaseName}/{collectionName}/{indexName}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.clusterName == "" {
		return localVarReturnValue, nil, reportError("clusterName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(r.clusterName), -1)
	if r.collectionName == "" {
		return localVarReturnValue, nil, reportError("collectionName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"collectionName"+"}", url.PathEscape(r.collectionName), -1)
	if r.databaseName == "" {
		return localVarReturnValue, nil, reportError("databaseName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"databaseName"+"}", url.PathEscape(r.databaseName), -1)
	if r.indexName == "" {
		return localVarReturnValue, nil, reportError("indexName is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"indexName"+"}", url.PathEscape(r.indexName), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.searchIndexUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("searchIndexUpdateRequest is required and must be specified")
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
	localVarPostBody = r.searchIndexUpdateRequest
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
