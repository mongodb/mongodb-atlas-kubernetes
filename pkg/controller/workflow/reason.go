package workflow

const (
	// General reasons
	AtlasCredentialsNotProvided ConditionReason = "AtlasCredentialsNotProvided"
	Internal                    ConditionReason = "InternalError"

	// Atlas Project
	ProjectNotCreatedInAtlas ConditionReason = "ProjectNotCreatedInAtlas"
)
