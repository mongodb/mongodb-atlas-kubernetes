package status

type CloudProviderIntegration struct {
	AtlasAWSAccountArn         string         `json:"atlasAWSAccountArn,omitempty"`
	AtlasAssumedRoleExternalID string         `json:"atlasAssumedRoleExternalId"`
	AuthorizedDate             string         `json:"authorizedDate,omitempty"`
	CreatedDate                string         `json:"createdDate,omitempty"`
	FeatureUsages              []FeatureUsage `json:"featureUsages,omitempty"`
	IamAssumedRoleArn          string         `json:"iamAssumedRoleArn,omitempty"`
	ProviderName               string         `json:"providerName"`
	RoleID                     string         `json:"roleId,omitempty"`
	Status                     string         `json:"status,omitempty"`
	ErrorMessage               string         `json:"errorMessage,omitempty"`
}

type FeatureUsage struct {
	FeatureType string `json:"featureType,omitempty"`
	FeatureID   string `json:"featureId,omitempty"`
}

const (
	CloudProviderIntegrationStatusNew                 = "NEW"
	CloudProviderIntegrationStatusCreated             = "CREATED"
	CloudProviderIntegrationStatusAuthorized          = "AUTHORIZED"
	CloudProviderIntegrationStatusDeAuthorize         = "DEAUTHORIZE"
	CloudProviderIntegrationStatusFailedToCreate      = "FAILED_TO_CREATE"
	CloudProviderIntegrationStatusFailedToAuthorize   = "FAILED_TO_AUTHORIZE"
	CloudProviderIntegrationStatusFailedToDeAuthorize = "FAILED_TO_DEAUTHORIZE"

	StatusFailed = "FAILED"
	StatusReady  = "READY"
)

func NewCloudProviderIntegration(providerName, assumedRoleArn string) CloudProviderIntegration {
	return CloudProviderIntegration{
		ProviderName:      providerName,
		IamAssumedRoleArn: assumedRoleArn,
		Status:            CloudProviderIntegrationStatusNew,
	}
}
