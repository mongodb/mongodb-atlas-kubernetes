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

	// Resource not found
	ResourceNotFound = "RESOURCE_NOT_FOUND"

	// Instance for the passed {groupId, tenantName} pair does not exist
	DataFederationTenantNotFound = "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_NAME"
)
