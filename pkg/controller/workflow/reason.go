package workflow

type ConditionReason string

// TODO move 'ConditionReason' to 'api' package?

// General reasons
const (
	Internal                      ConditionReason = "InternalError"
	AtlasResourceVersionMismatch  ConditionReason = "AtlasResourceVersionMismatch"
	AtlasResourceVersionIsInvalid ConditionReason = "AtlasResourceVersionIsInvalid"
	AtlasFinalizerNotSet          ConditionReason = "AtlasFinalizerNotSet"
	AtlasFinalizerNotRemoved      ConditionReason = "AtlasFinalizerNotRemoved"
	AtlasDeletionProtection       ConditionReason = "AtlasDeletionProtection"
	AtlasGovUnsupported           ConditionReason = "AtlasGovUnsupported"
	AtlasAPIAccessNotConfigured   ConditionReason = "AtlasAPIAccessNotConfigured"
)

// Atlas Project reasons
const (
	ProjectNotCreatedInAtlas                       ConditionReason = "ProjectNotCreatedInAtlas"
	ProjectIPAccessInvalid                         ConditionReason = "ProjectIPAccessListInvalid"
	ProjectIPNotCreatedInAtlas                     ConditionReason = "ProjectIPAccessListNotCreatedInAtlas"
	ProjectWindowInvalid                           ConditionReason = "ProjectWindowInvalid"
	ProjectWindowNotObtainedFromAtlas              ConditionReason = "ProjectWindowNotObtainedFromAtlas"
	ProjectWindowNotCreatedInAtlas                 ConditionReason = "ProjectWindowNotCreatedInAtlas"
	ProjectWindowNotDeletedInAtlas                 ConditionReason = "projectWindowNotDeletedInAtlas"
	ProjectWindowNotDeferredInAtlas                ConditionReason = "ProjectWindowNotDeferredInAtlas"
	ProjectWindowNotAutoDeferredInAtlas            ConditionReason = "ProjectWindowNotAutoDeferredInAtlas"
	ProjectPEServiceIsNotReadyInAtlas              ConditionReason = "ProjectPrivateEndpointServiceIsNotReadyInAtlas"
	ProjectPEInterfaceIsNotReadyInAtlas            ConditionReason = "ProjectPrivateEndpointIsNotReadyInAtlas"
	ProjectIPAccessListNotActive                   ConditionReason = "ProjectIPAccessListNotActive"
	ProjectIntegrationInternal                     ConditionReason = "ProjectIntegrationInternalError"
	ProjectIntegrationRequest                      ConditionReason = "ProjectIntegrationRequestError"
	ProjectIntegrationReady                        ConditionReason = "ProjectIntegrationReady"
	ProjectPrivateEndpointIsNotReadyInAtlas        ConditionReason = "ProjectPrivateEndpointIsNotReadyInAtlas"
	ProjectNetworkPeerIsNotReadyInAtlas            ConditionReason = "ProjectNetworkPeerIsNotReadyInAtlas"
	ProjectEncryptionAtRestReady                   ConditionReason = "ProjectEncryptionAtRestReady"
	ProjectCloudIntegrationsIsNotReadyInAtlas      ConditionReason = "ProjectCloudIntegrationsIsNotReadyInAtlas"
	ProjectAuditingReady                           ConditionReason = "ProjectAuditingReady"
	ProjectSettingsReady                           ConditionReason = "ProjectSettingsReady"
	ProjectAlertConfigurationIsNotReadyInAtlas     ConditionReason = "ProjectAlertConfigurationIsNotReadyInAtlas"
	ProjectCustomRolesReady                        ConditionReason = "ProjectCustomRolesReady"
	ProjectTeamUnavailable                         ConditionReason = "ProjectTeamUnavailable"
	ProjectBackupCompliancePolicyUnavailable       ConditionReason = "ProjectBackupCompliancePolicyUnavailable"
	ProjectBackupCompliancePolicyNotMet            ConditionReason = "ProjectBackupCompliancePolicyNotMet"
	ProjectBackupCompliancePolicyNotCreatedInAtlas ConditionReason = "ProjectBackupCompliancePolicyNotCreatedInAtlas"
)

// Atlas Deployment reasons
const (
	DeploymentNotCreatedInAtlas           ConditionReason = "DeploymentNotCreatedInAtlas"
	DeploymentNotUpdatedInAtlas           ConditionReason = "DeploymentNotUpdatedInAtlas"
	DeploymentCreating                    ConditionReason = "DeploymentCreating"
	DeploymentUpdating                    ConditionReason = "DeploymentUpdating"
	DeploymentConnectionSecretsNotCreated ConditionReason = "DeploymentConnectionSecretsNotCreated"
	DeploymentAdvancedOptionsReady        ConditionReason = "DeploymentAdvancedOptionsReady"
	ServerlessPrivateEndpointReady        ConditionReason = "ServerlessPrivateEndpointReady"
	ManagedNamespacesReady                ConditionReason = "ManagedNamespacesReady"
	CustomZoneMappingReady                ConditionReason = "CustomZoneMappingReady"
)

// Atlas Database User reasons
const (
	DatabaseUserNotCreatedInAtlas           ConditionReason = "DatabaseUserNotCreatedInAtlas"
	DatabaseUserNotUpdatedInAtlas           ConditionReason = "DatabaseUserNotUpdatedInAtlas"
	DatabaseUserNotDeletedInAtlas           ConditionReason = "DatabaseUserNotDeletedInAtlas"
	DatabaseUserConnectionSecretsNotCreated ConditionReason = "DatabaseUserConnectionSecretsNotCreated"
	DatabaseUserConnectionSecretsNotDeleted ConditionReason = "DatabaseUserConnectionSecretsNotDeleted"
	DatabaseUserStaleConnectionSecrets      ConditionReason = "DatabaseUserStaleConnectionSecrets"
	DatabaseUserDeploymentAppliedChanges    ConditionReason = "DeploymentAppliedDatabaseUsersChanges"
	DatabaseUserInvalidSpec                 ConditionReason = "DatabaseUserInvalidSpec"
	DatabaseUserExpired                     ConditionReason = "DatabaseUserExpired"
)

// Atlas Data Federation reasons
const (
	DataFederationNotCreatedInAtlas ConditionReason = "DataFederationNotCreatedInAtlas"
	DataFederationNotUpdatedInAtlas ConditionReason = "DataFederationNotUpdatedInAtlas"
	DataFederationCreating          ConditionReason = "DataFederationCreating"
	DataFederationUpdating          ConditionReason = "DataFederationUpdating"
)

// Atlas Teams reasons
const (
	TeamNotCreatedInAtlas ConditionReason = "TeamNotCreatedInAtlas"
	TeamNotUpdatedInAtlas ConditionReason = "TeamNotUpdatedInAtlas"
	TeamInvalidSpec       ConditionReason = "TeamInvalidSpec"
	TeamUsersNotReady     ConditionReason = "TeamUsersNotReady"
	TeamDoesNotExist      ConditionReason = "TeamDoesNotExist"
)

// Atlas Federated Auth reasons
const (
	FederatedAuthNotAvailable     ConditionReason = "FederatedAuthNotAvailable"
	FederatedAuthIsNotEnabledInCR ConditionReason = "FederatedAuthNotEnabledInCR"
	FederatedAuthOrgNotConnected  ConditionReason = "FederatedAuthOrgIsNotConnected"
	FederatedAuthUsersConflict    ConditionReason = "FederatedAuthUsersConflict"
)
