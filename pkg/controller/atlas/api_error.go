package atlas

const (
	// Error codes that Atlas may return that we are concerned about
	GroupExistsAPIErrorCode = "GROUP_ALREADY_EXISTS"

	// The error that Atlas API returns if the GET request is sent to read the project that either doesn't exist
	// or the user doesn't have permissions for
	NotInGroup = "NOT_IN_GROUP"

	// Error indicates that the project is being removed while it still has clusters
	CannotCloseGroupActiveAtlasDeployment = "CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS"

	// Error indicates that the database user doesn't exist
	UsernameNotFound = "USERNAME_NOT_FOUND"

	// Error indicates that the cluster doesn't exist
	ClusterNotFound = "CLUSTER_NOT_FOUND"

	// ServerlessClusterNotFound indicates that the serverless cluster doesn't exist
	ServerlessInstanceNotFound = "SERVERLESS_INSTANCE_NOT_FOUND"

	// ServerlessClusterFromClusterAPI indicates that we are trying to access
	// a serverless instance from the cluster API, which is not allowed
	ServerlessInstanceFromClusterAPI = "CANNOT_USE_SERVERLESS_INSTANCE_IN_CLUSTER_API"

	// Resource not found
	ResourceNotFound = "RESOURCE_NOT_FOUND"
)
