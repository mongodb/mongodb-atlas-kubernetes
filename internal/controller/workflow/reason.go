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
	AtlasUnsupportedFeature       ConditionReason = "AtlasUnsupportedFeature"
)

// Atlas Project reasons
const (
	ProjectNotCreatedInAtlas                   ConditionReason = "ProjectNotCreatedInAtlas"
	ProjectBeingConfiguredInAtlas              ConditionReason = "ProjectBeingConfiguredInAtlas"
	ProjectIPAccessInvalid                     ConditionReason = "ProjectIPAccessListInvalid"
	ProjectIPNotCreatedInAtlas                 ConditionReason = "ProjectIPAccessListNotCreatedInAtlas"
	ProjectWindowInvalid                       ConditionReason = "ProjectWindowInvalid"
	ProjectWindowNotObtainedFromAtlas          ConditionReason = "ProjectWindowNotObtainedFromAtlas"
	ProjectWindowNotCreatedInAtlas             ConditionReason = "ProjectWindowNotCreatedInAtlas"
	ProjectWindowNotDeletedInAtlas             ConditionReason = "projectWindowNotDeletedInAtlas"
	ProjectWindowNotDeferredInAtlas            ConditionReason = "ProjectWindowNotDeferredInAtlas"
	ProjectWindowNotAutoDeferredInAtlas        ConditionReason = "ProjectWindowNotAutoDeferredInAtlas"
	ProjectPEServiceIsNotReadyInAtlas          ConditionReason = "ProjectPrivateEndpointServiceIsNotReadyInAtlas"
	ProjectPEInterfaceIsNotReadyInAtlas        ConditionReason = "ProjectPrivateEndpointIsNotReadyInAtlas"
	ProjectIPAccessListNotActive               ConditionReason = "ProjectIPAccessListNotActive"
	ProjectIntegrationInternal                 ConditionReason = "ProjectIntegrationInternalError"
	ProjectIntegrationRequest                  ConditionReason = "ProjectIntegrationRequestError"
	ProjectIntegrationReady                    ConditionReason = "ProjectIntegrationReady"
	ProjectPrivateEndpointIsNotReadyInAtlas    ConditionReason = "ProjectPrivateEndpointIsNotReadyInAtlas"
	ProjectNetworkPeerIsNotReadyInAtlas        ConditionReason = "ProjectNetworkPeerIsNotReadyInAtlas"
	ProjectEncryptionAtRestReady               ConditionReason = "ProjectEncryptionAtRestReady"
	ProjectCloudIntegrationsIsNotReadyInAtlas  ConditionReason = "ProjectCloudIntegrationsIsNotReadyInAtlas"
	ProjectAuditingReady                       ConditionReason = "ProjectAuditingReady"
	ProjectSettingsReady                       ConditionReason = "ProjectSettingsReady"
	ProjectAlertConfigurationIsNotReadyInAtlas ConditionReason = "ProjectAlertConfigurationIsNotReadyInAtlas"
	ProjectCustomRolesReady                    ConditionReason = "ProjectCustomRolesReady"
	ProjectTeamUnavailable                     ConditionReason = "ProjectTeamUnavailable"
	ProjectX509NotConfigured                   ConditionReason = "ProjectX509NotConfigured"
)

// Atlas Backup Compliance Policy reasons
const (
	ProjectBackupCompliancePolicyNotMet            ConditionReason = "ProjectBackupCompliancePolicyNotMet"
	ProjectBackupCompliancePolicyNotCreatedInAtlas ConditionReason = "ProjectBackupCompliancePolicyNotCreatedInAtlas"
	ProjectBackupCompliancePolicyCannotDelete      ConditionReason = "ProjectBackupCompliancePolicyCannotDelete"
	ProjectBackupCompliancePolicyUpdating          ConditionReason = "ProjectBackupCompliancePolicyUpdating"
	ProjectBackupCompliancePolicyOperationAborted  ConditionReason = "ProjectBackupCompliancePolicyOperationAborted"
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
	ServerlessPrivateEndpointFailed       ConditionReason = "ServerlessPrivateEndpointFailed"
	ServerlessPrivateEndpointInProgress   ConditionReason = "ServerlessPrivateEndpointInProgress"
	ManagedNamespacesReady                ConditionReason = "ManagedNamespacesReady"
	CustomZoneMappingReady                ConditionReason = "CustomZoneMappingReady"
)

// Atlas SearchNodes reasons
const (
	SearchNodesUpdating ConditionReason = "SearchNodesUpdating"
	SearchNodesCreating ConditionReason = "SearchNodesCreating"
	SearchNodesDeleting ConditionReason = "SearchNodesDeleting"

	ErrorSearchNodesNotUpsertedInAtlas ConditionReason = "SearchNodesNotUpsertedInAtlas"
	ErrorSearchNodesNotDeletedInAtlas  ConditionReason = "SearchNodesNotDeletedInAtlas"
	ErrorSearchNodesOperationAborted   ConditionReason = "SearchNodesOperationAborted"
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
	TeamNotCleaned        ConditionReason = "TeamCleanupFailed"
)

// Atlas Federated Auth reasons
const (
	FederatedAuthNotAvailable     ConditionReason = "FederatedAuthNotAvailable"
	FederatedAuthIsNotEnabledInCR ConditionReason = "FederatedAuthNotEnabledInCR"
	FederatedAuthOrgNotConnected  ConditionReason = "FederatedAuthOrgIsNotConnected"
	FederatedAuthUsersConflict    ConditionReason = "FederatedAuthUsersConflict"
)

// Atlas Streams reasons
const (
	StreamInstanceSetupInProgress ConditionReason = "StreamInstanceSetupInProgress"
	StreamInstanceNotCreated      ConditionReason = "StreamInstanceNotCreated"
	StreamInstanceNotRemoved      ConditionReason = "StreamInstanceNotRemoved"
	StreamInstanceNotUpdated      ConditionReason = "StreamInstanceNotUpdated"
	StreamConnectionNotConfigured ConditionReason = "StreamConnectionNotConfigured"
	StreamConnectionNotCreated    ConditionReason = "StreamConnectionNotCreated"
	StreamConnectionNotRemoved    ConditionReason = "StreamConnectionNotRemoved"
	StreamConnectionNotUpdated    ConditionReason = "StreamConnectionNotUpdated"
)

const (
	AtlasCustomRoleNotCreated ConditionReason = "CustomRoleNotCreated"
	AtlasCustomRoleNotUpdated ConditionReason = "CustomRoleNotUpdated"
	AtlasCustomRoleNotDeleted ConditionReason = "CustomRoleNotDeleted"
)

// Atlas Private Endpoint reasons
const (
	PrivateEndpointServiceCreated           ConditionReason = "PrivateEndpointServiceCreated"
	PrivateEndpointServiceFailedToCreate    ConditionReason = "PrivateEndpointServiceFailedToCreate"
	PrivateEndpointServiceFailedToConfigure ConditionReason = "PrivateEndpointServiceFailedToConfigure"
	PrivateEndpointServiceInitializing      ConditionReason = "PrivateEndpointServiceInitializing"
	PrivateEndpointServiceDeleting          ConditionReason = "PrivateEndpointServiceDeleting"
	PrivateEndpointFailedToCreate           ConditionReason = "PrivateEndpointFailedToCreate"
	PrivateEndpointUpdating                 ConditionReason = "PrivateEndpointUpdating"
	PrivateEndpointConfigurationPending     ConditionReason = "PrivateEndpointConfigurationPending"
	PrivateEndpointFailedToConfigure        ConditionReason = "PrivateEndpointFailedToConfigure"
	PrivateEndpointFailedToDelete           ConditionReason = "PrivateEndpointFailedToDelete"
)

// Atlas IP Access List reasons
const (
	IPAccessListFailedToCreate   ConditionReason = "IPAccessListFailedToCreate"
	IPAccessListFailedToDelete   ConditionReason = "IPAccessListFailedToDelete"
	IPAccessListFailedToGetState ConditionReason = "IPAccessListFailedToGetState"
	IPAccessListPending          ConditionReason = "IPAccessListPending"
)

// Atlas Network Peering reasons
const (
	NetworkPeeringConnectionCreating ConditionReason = "NetworkPeeringConnectionCreating"
	NetworkPeeringConnectionPending  ConditionReason = "NetworkPeeringConnectionPending"
	NetworkPeeringRemovingContainer  ConditionReason = "NetworkPeeringRemovingContainer"
)

// Atlas Network Container reasons
const (
	NetworkContainerNotConfigured ConditionReason = "NetworkContainerNotConfigured"
)
