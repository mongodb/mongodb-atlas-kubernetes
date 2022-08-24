package v1

type CloudProviderAccessRole struct {
	// ProviderName is the name of the cloud provider. Currently only AWS is supported.
	ProviderName string `json:"providerName"`
	// IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.
	IamAssumedRoleArn string `json:"iamAssumedRoleArn"`
}
