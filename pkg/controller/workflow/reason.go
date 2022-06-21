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
	ProjectNotCreatedInAtlas                ConditionReason = "ProjectNotCreatedInAtlas"
	ProjectIPAccessInvalid                  ConditionReason = "ProjectIPAccessListInvalid"
	ProjectIPNotCreatedInAtlas              ConditionReason = "ProjectIPAccessListNotCreatedInAtlas"
	ProjectWindowInvalid                    ConditionReason = "ProjectWindowInvalid"
	ProjectWindowNotCreatedInAtlas          ConditionReason = "ProjectWindowNotCreatedInAtlas"
	ProjectWindowNotDeletedInAtlas          ConditionReason = "projectWindowNotDeletedInAtlas"
	ProjectPEServiceIsNotReadyInAtlas       ConditionReason = "ProjectPrivateEndpointServiceIsNotReadyInAtlas"
	ProjectPrivateEndpointIsNotReadyInAtlas ConditionReason = "ProjectPrivateEndpointIsNotReadyInAtlas"
	ProjectIPAccessListNotActive            ConditionReason = "ProjectIPAccessListNotActive"
	ProjectIntegrationInternal              ConditionReason = "ProjectIntegrationInternalError"
	ProjectIntegrationRequest               ConditionReason = "ProjectIntegrationRequestError"
	ProjectIntegrationReady                 ConditionReason = "ProjectIntegrationReady"
)

// Atlas Cluster reasons
const (
	DeploymentNotCreatedInAtlas           ConditionReason = "DeploymentNotCreatedInAtlas"
	DeploymentNotUpdatedInAtlas           ConditionReason = "DeploymentNotUpdatedInAtlas"
	DeploymentCreating                    ConditionReason = "DeploymentCreating"
	DeploymentUpdating                    ConditionReason = "DeploymentUpdating"
	DeploymentConnectionSecretsNotCreated ConditionReason = "DeploymentConnectionSecretsNotCreated"
	DeploymentAdvancedOptionsAreNotReady  ConditionReason = "DeploymentAdvancedOptionsAreNotReady"
)

// Atlas Database User reasons
const (
	DatabaseUserNotCreatedInAtlas           ConditionReason = "DatabaseUserNotCreatedInAtlas"
	DatabaseUserNotUpdatedInAtlas           ConditionReason = "DatabaseUserNotUpdatedInAtlas"
	DatabaseUserConnectionSecretsNotCreated ConditionReason = "DatabaseUserConnectionSecretsNotCreated"
	DatabaseUserStaleConnectionSecrets      ConditionReason = "DatabaseUserStaleConnectionSecrets"
	DatabaseUserDeploymentAppliedChanges    ConditionReason = "DeploymentAppliedDatabaseUsersChanges"
	DatabaseUserInvalidSpec                 ConditionReason = "DatabaseUserInvalidSpec"
	DatabaseUserExpired                     ConditionReason = "DatabaseUserExpired"
)
