package workflow

type ConditionReason string

// TODO move 'ConditionReason' to 'api' package?

// General reasons
const (
	AtlasCredentialsNotProvided   ConditionReason = "AtlasCredentialsNotProvided"
	Internal                      ConditionReason = "InternalError"
	AtlasResourceVersionMismatch  ConditionReason = "AtlasResourceVersionMismatch"
	AtlasResourceVersionIsInvalid ConditionReason = "AtlasResourceVersionIsInvalid"
)

// Atlas Project reasons
const (
	ProjectNotCreatedInAtlas                   ConditionReason = "ProjectNotCreatedInAtlas"
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
	ProjectCloudAccessRolesIsNotReadyInAtlas   ConditionReason = "ProjectCloudAccessRolesIsNotReadyInAtlas"
	ProjectAuditingReady                       ConditionReason = "ProjectAuditingReady"
	ProjectSettingsReady                       ConditionReason = "ProjectSettingsReady"
	ProjectAlertConfigurationIsNotReadyInAtlas ConditionReason = "ProjectAlertConfigurationIsNotReadyInAtlas"
	ProjectCustomRolesReady                    ConditionReason = "ProjectCustomRolesReady"
	ProjectTeamUnavailable                     ConditionReason = "ProjectTeamUnavailable"
)

// Atlas Cluster reasons
const (
	DeploymentNotCreatedInAtlas           ConditionReason = "DeploymentNotCreatedInAtlas"
	DeploymentNotUpdatedInAtlas           ConditionReason = "DeploymentNotUpdatedInAtlas"
	DeploymentCreating                    ConditionReason = "DeploymentCreating"
	DeploymentUpdating                    ConditionReason = "DeploymentUpdating"
	DeploymentDeleting                    ConditionReason = "DeploymentDeleting"
	DeploymentDeleted                     ConditionReason = "DeploymentDeleted"
	DeploymentConnectionSecretsNotCreated ConditionReason = "DeploymentConnectionSecretsNotCreated"
	DeploymentAdvancedOptionsReady        ConditionReason = "DeploymentAdvancedOptionsReady"
	DeploymentAdvancedOptionsAreNotReady  ConditionReason = "DeploymentAdvancedOptionsAreNotReady"
	ServerlessPrivateEndpointReady        ConditionReason = "ServerlessPrivateEndpointReady"
	ManagedNamespacesReady                ConditionReason = "ManagedNamespacesReady"
	CustomZoneMappingReady                ConditionReason = "CustomZoneMappingReady"
	ClusterNotCreatedInAtlas              ConditionReason = "ClusterNotCreatedInAtlas"
	ClusterNotUpdatedInAtlas              ConditionReason = "ClusterNotUpdatedInAtlas"
	ClusterCreating                       ConditionReason = "ClusterCreating"
	ClusterUpdating                       ConditionReason = "ClusterUpdating"
	ClusterDeleting                       ConditionReason = "ClusterDeleting"
	ClusterDeleted                        ConditionReason = "ClusterDeleted"
	ClusterConnectionSecretsNotCreated    ConditionReason = "ClusterConnectionSecretsNotCreated"
	ClusterAdvancedOptionsAreNotReady     ConditionReason = "ClusterAdvancedOptionsAreNotReady"
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

const (
	TeamNotCreatedInAtlas ConditionReason = "TeamNotCreatedInAtlas"
	TeamNotUpdatedInAtlas ConditionReason = "TeamNotUpdatedInAtlas"
	TeamInvalidSpec       ConditionReason = "TeamInvalidSpec"
	TeamUsersNotReady     ConditionReason = "TeamUsersNotReady"
	TeamDoesNotExist      ConditionReason = "TeamDoesNotExist"
)

// MongoDBAtlasInventory reasons
const (
	MongoDBAtlasInventorySyncOK              ConditionReason = "SyncOK"
	MongoDBAtlasInventoryInputError          ConditionReason = "InputError"
	MongoDBAtlasInventoryBackendError        ConditionReason = "BackendError"
	MongoDBAtlasInventoryEndpointUnreachable ConditionReason = "EndpointUnreachable"
	MongoDBAtlasInventoryAuthenticationError ConditionReason = "AuthenticationError"
)

// MongoDBAtlasConnection reasons
const (
	MongoDBAtlasConnectionReady               ConditionReason = "Ready"
	MongoDBAtlasConnectionAtlasUnreachable    ConditionReason = "Unreachable"
	MongoDBAtlasConnectionInventoryNotReady   ConditionReason = "InventoryNotReady"
	MongoDBAtlasConnectionInventoryNotFound   ConditionReason = "InventoryNotFound"
	MongoDBAtlasConnectionInstanceIDNotFound  ConditionReason = "InstanceIDNotFound"
	MongoDBAtlasConnectionBackendError        ConditionReason = "BackendError"
	MongoDBAtlasConnectionAuthenticationError ConditionReason = "AuthenticationError"
	MongoDBAtlasConnectionInprogress          ConditionReason = "Inprogress"
)

// MongoDBAtlasInstance reasons
const (
	MongoDBAtlasInstanceReady               ConditionReason = "Ready"
	MongoDBAtlasInstanceAtlasUnreachable    ConditionReason = "Unreachable"
	MongoDBAtlasInstanceInventoryNotFound   ConditionReason = "InventoryNotFound"
	MongoDBAtlasInstanceNotReady            ConditionReason = "InstanceNotReady"
	MongoDBAtlasInstanceClusterNotFound     ConditionReason = "AtlasClusterNotFound"
	MongoDBAtlasInstanceBackendError        ConditionReason = "BackendError"
	MongoDBAtlasInstanceAuthenticationError ConditionReason = "AuthenticationError"
	MongoDBAtlasInstanceInprogress          ConditionReason = "Inprogress"
)
