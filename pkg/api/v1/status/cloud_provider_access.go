package status

type CloudProviderAccessRole struct {
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
	CloudProviderAccessStatusNew                 = "NEW"
	CloudProviderAccessStatusCreated             = "CREATED"
	CloudProviderAccessStatusAuthorized          = "AUTHORIZED"
	CloudProviderAccessStatusDeAuthorize         = "DEAUTHORIZE"
	CloudProviderAccessStatusFailedToCreate      = "FAILED_TO_CREATE"
	CloudProviderAccessStatusFailedToAuthorize   = "FAILED_TO_AUTHORIZE"
	CloudProviderAccessStatusFailedToDeAuthorize = "FAILED_TO_DEAUTHORIZE"

	StatusFailed   = "FAILED"
	StatusCreated  = "CREATED"
	StatusReady    = "READY"
	StatusEmptyARN = "EMPTY_ARN"
)

func NewCloudProviderAccessRole(providerName, assumedRoleArn string) CloudProviderAccessRole {
	return CloudProviderAccessRole{
		ProviderName:      providerName,
		IamAssumedRoleArn: assumedRoleArn,
		Status:            CloudProviderAccessStatusNew,
	}
}
