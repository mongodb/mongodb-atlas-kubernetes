package workflow

type ConditionReason string

// General reasons
const (
	AtlasCredentialsNotProvided ConditionReason = "AtlasCredentialsNotProvided"
	Internal                    ConditionReason = "InternalError"
)

// Atlas Project reasons
const (
	ProjectNotCreatedInAtlas ConditionReason = "ProjectNotCreatedInAtlas"
)

// Atlas Cluster reasons
const (
	ClusterNotCreatedInAtlas ConditionReason = "ClusterNotCreatedInAtlas"
	ClusterCreating          ConditionReason = "ClusterCreating"
	ClusterUpdating          ConditionReason = "ClusterUpdating"
	ClusterNotUpToDate       ConditionReason = "ClusterNotUpToDate"
)
