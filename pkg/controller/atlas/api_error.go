package atlas

const (
	// Error codes that Atlas may return that we are concerned about
	GroupExistsAPIErrorCode = "GROUP_ALREADY_EXISTS"

	// The error that Atlas API returns if the GET request is sent to read the project that either doesn't exist
	// or the user doesn't have permissions for
	NotInGroup = "NOT_IN_GROUP"

	// Error indicates that the project is being removed while it still has clusters
	CannotCloseGroupActiveAtlasCluster = "CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS"
)
