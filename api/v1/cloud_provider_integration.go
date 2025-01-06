package v1

// CloudProviderIntegration define an integration to a cloud provider
type CloudProviderIntegration struct {
	// ProviderName is the name of the cloud provider. Currently only AWS is supported.
	ProviderName string `json:"providerName"`
	// IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.
	// +optional
	IamAssumedRoleArn string `json:"iamAssumedRoleArn"`
}

// CloudProviderAccessRole define an integration to a cloud provider
// Deprecated: This type is deprecated in favor of CloudProviderIntegration
type CloudProviderAccessRole struct {
	// ProviderName is the name of the cloud provider. Currently only AWS is supported.
	ProviderName string `json:"providerName"`
	// IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.
	// +optional
	IamAssumedRoleArn string `json:"iamAssumedRoleArn"`
}
