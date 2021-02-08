package workflow

type ConditionReason string

// TODO move 'ConditionReason' to 'api' package?

// General reasons
const (
	AtlasCredentialsNotProvided ConditionReason = "AtlasCredentialsNotProvided"
	Internal                    ConditionReason = "InternalError"
)

// Atlas Project reasons
const (
	ProjectNotCreatedInAtlas   ConditionReason = "ProjectNotCreatedInAtlas"
	ProjectIPAccessInvalid     ConditionReason = "ProjectIPAccessListInvalid"
	ProjectIPNotCreatedInAtlas ConditionReason = "ProjectIPAccessListNotCreatedInAtlas"
)

// Atlas Cluster reasons
const (
	ClusterNotCreatedInAtlas ConditionReason = "ClusterNotCreatedInAtlas"
	ClusterCreating          ConditionReason = "ClusterCreating"
	ClusterUpdating          ConditionReason = "ClusterUpdating"
)
